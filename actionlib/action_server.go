package actionlib

import (
	"actionlib_msgs"
	"reflect"
	"std_msgs"
	"sync"
	"time"

	"github.com/fetchrobotics/rosgo/ros"
)

type defaultActionServer struct {
	node              ros.Node
	autoStart         bool
	started           bool
	action            string
	actionType        ActionType
	actionResult      ros.MessageType
	actionResultType  ros.MessageType
	actionFeedback    ros.MessageType
	actionGoal        ros.MessageType
	statusList        []*status
	statusMutex       sync.RWMutex
	statusFrequency   ros.Rate
	statusTimer       *time.Ticker
	statusListTimeout ros.Duration
	goalCallback      interface{}
	cancelCallback    interface{}
	lastCancel        ros.Time
	pubQueueSize      int
	subQueueSize      int
	goalSub           ros.Subscriber
	cancelSub         ros.Subscriber
	resultPub         ros.Publisher
	feedbackPub       ros.Publisher
	statusPub         ros.Publisher
	statusPubChan     chan struct{}
	shutdownChan      chan struct{}
}

func newDefaultActionServer(node ros.Node, action string, actType ActionType, goalCb interface{}, cancelCb interface{}, start bool) *defaultActionServer {
	server := new(defaultActionServer)
	server.node = node
	server.autoStart = start
	server.started = false
	server.action = action
	server.actionType = actType
	server.actionResult = actType.ResultType()
	server.actionFeedback = actType.FeedbackType()
	server.actionGoal = actType.GoalType()
	server.statusListTimeout = ros.NewDuration(60, 0)
	server.goalCallback = goalCb
	server.cancelCallback = cancelCb
	server.lastCancel = ros.Now()
	server.goalSub = node.NewSubscriber(action+"/goal", actType.GoalType(), server.internalGoalCallback)
	server.cancelSub = node.NewSubscriber(action+"/cancel", actionlib_msgs.MsgGoalID, server.internalCancelCallback)
	server.resultPub = node.NewPublisher(action+"/result", actType.ResultType())
	server.feedbackPub = node.NewPublisher(action+"/feedback", actType.FeedbackType())
	server.statusPub = node.NewPublisher(action+"/status", actionlib_msgs.MsgGoalStatusArray)
	server.statusPubChan = make(chan struct{}, 100)
	server.shutdownChan = make(chan struct{}, 10)
	// get frequency from ros params
	server.statusFrequency = ros.NewRate(5.0)
	// get queue sizes from ros params
	// queue sizes not implemented by ros.Node yet
	server.pubQueueSize = 50
	server.subQueueSize = 50

	return server
}

func (as *defaultActionServer) Start() {
	logger := as.node.Logger()
	defer func() {
		logger.Debug("defaultActionServer.start exit")
		as.started = false
	}()

	as.statusTimer = time.NewTicker(time.Second / 5.0)
	defer as.statusTimer.Stop()
	as.started = true

	for {
		select {
		case <-as.shutdownChan:
			return

		case <-as.statusTimer.C:
			as.publishStatus()

		case <-as.statusPubChan:
			arr := as.getStatus()
			as.statusPub.Publish(arr)
		}
	}
}

// publishResult publishes action result message
func (as *defaultActionServer) publishResult(status actionlib_msgs.GoalStatus, result ros.Message) {
	msg := as.actionResult.NewMessage().(ActionResult)
	msg.SetHeader(std_msgs.Header{Stamp: ros.Now()})
	msg.SetStatus(status)
	msg.SetResult(result)
	as.resultPub.Publish(msg)
}

// publishFeedback publishes action feedback messages
func (as *defaultActionServer) publishFeedback(status actionlib_msgs.GoalStatus, feedback ros.Message) {
	msg := as.actionFeedback.NewMessage().(ActionFeedback)
	msg.SetHeader(std_msgs.Header{Stamp: ros.Now()})
	msg.SetStatus(status)
	msg.SetFeedback(feedback)
	as.feedbackPub.Publish(msg)
}

// publishStatus publishes action status messages
func (as *defaultActionServer) getStatus() *actionlib_msgs.GoalStatusArray {
	var statArr []actionlib_msgs.GoalStatus
	as.statusMutex.Lock()
	defer as.statusMutex.Unlock()

	if as.node.OK() {
		for i, st := range as.statusList {
			destTime := st.destroyTime.Add(as.statusListTimeout)
			if !st.destroyTime.IsZero() && destTime.Cmp(ros.Now()) <= 0 {
				as.statusList[i] = as.statusList[len(as.statusList)-1]
				as.statusList[len(as.statusList)-1] = nil
				as.statusList = as.statusList[:len(as.statusList)-1]
				continue
			}

			statArr = append(statArr, st.getGoalStatus())
		}
	}

	arr := &actionlib_msgs.GoalStatusArray{}
	arr.Header.Stamp = ros.Now()
	arr.StatusList = statArr
	return arr
}

func (as *defaultActionServer) publishStatus() {
	as.statusPubChan <- struct{}{}
}

// internalCancelCallback recieves cancel message from client
func (as *defaultActionServer) internalCancelCallback(goalID *actionlib_msgs.GoalID) {
	logger := as.node.Logger()
	logger.Debug("action server has received a new cancel request")

	as.statusMutex.Lock()
	defer as.statusMutex.Unlock()

	goalFound := false
	for _, st := range as.statusList {
		cancelAll := (goalID.Id == "" && goalID.Stamp.IsZero())
		cancelCurrent := (goalID.Id == st.goalStatus.GoalId.Id)
		cancelBeforeStamp := (!goalID.Stamp.IsZero() && st.goalStatus.GoalId.Stamp.Cmp(goalID.Stamp) <= 0)

		if cancelAll || cancelCurrent || cancelBeforeStamp {
			if goalID.Id == st.goalStatus.GoalId.Id {
				goalFound = true
			}

			st.destroyTime = ros.Now()
			if st.setCancelRequested() {
				args := []reflect.Value{reflect.ValueOf(goalID)}
				fun := reflect.ValueOf(as.cancelCallback)
				numArgsNeeded := fun.Type().NumIn()

				if numArgsNeeded <= 1 {
					fun.Call(args[0:numArgsNeeded])
				}
			}
		}
	}

	if goalID.Id != "" && !goalFound {
		st := newStatusWithGoalStatus(as, actionlib_msgs.GoalStatus{
			GoalId: *goalID,
			Status: actionlib_msgs.RECALLING,
		})

		as.statusList = append(as.statusList, st)
		st.destroyTime = ros.Now()
	}

	if goalID.Stamp.Cmp(as.lastCancel) > 0 {
		as.lastCancel = goalID.Stamp
	}
}

// internalGoalCallback recieves the goals from client and checks if
// the goalID already exists in the status list. If not, it will call
// server's goalCallback with goal that was recieved from the client.
func (as *defaultActionServer) internalGoalCallback(goal ActionGoal) {
	logger := as.node.Logger()
	goalID := goal.GetGoalId()

	as.statusMutex.Lock()
	defer as.statusMutex.Unlock()

	for _, st := range as.statusList {
		if goalID.Id == st.goalStatus.GoalId.Id {
			logger.Debugf("Goal %s was already in the status list with status %+v", goalID.Id, st.goalStatus)
			if st.goalStatus.Status == actionlib_msgs.RECALLING {
				st.goalStatus.Status = actionlib_msgs.RECALLED
				result := as.actionResultType.NewMessage()
				as.publishResult(*st.goalStatus, result)
			}

			st.destroyTime = ros.Now()
			return
		}
	}

	st := newStatusWithActionGoal(as, goal)
	as.statusList = append(as.statusList, st)

	if !goalID.Stamp.IsZero() && goalID.Stamp.Cmp(as.lastCancel) <= 0 {
		st.setCancelled(nil, "timestamp older than last goal cancel")
		return
	}

	args := []reflect.Value{reflect.ValueOf(goal)}
	fun := reflect.ValueOf(as.goalCallback)
	numArgsNeeded := fun.Type().NumIn()

	if numArgsNeeded <= 1 {
		fun.Call(args[0:numArgsNeeded])
	}
}

// RegisterGoalCallback replaces existing goal callback function with newly
// provided goal callback function.
func (as *defaultActionServer) RegisterGoalCallback(goalCb interface{}) {
	as.goalCallback = goalCb
}

func (as *defaultActionServer) RegisterCancelCallback(cancelCb interface{}) {
	as.cancelCallback = cancelCb
}

func (as *defaultActionServer) Shutdown() {
	as.shutdownChan <- struct{}{}
}
