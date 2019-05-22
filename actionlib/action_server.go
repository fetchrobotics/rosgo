package actionlib

import (
	"actionlib_msgs"
	"container/list"
	"reflect"
	"std_msgs"
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
	statusList        *list.List
	statusFrequency   ros.Rate
	statusTimer       *time.Ticker
	statusListTimeout ros.Duration
	goalCallback      interface{}
	cancelCallback    interface{}
	lastCancel        ros.Time
	pubQueueSize      int
	subQueueSize      int
	goalSub           ros.Subscriber
	goalSubChan       chan ActionGoal
	cancelSub         ros.Subscriber
	cancelSubChan     chan *actionlib_msgs.GoalID
	statusPub         ros.Publisher
	statusPubChan     chan *actionlib_msgs.GoalStatusArray
	resultPub         ros.Publisher
	resultPubChan     chan ActionResult
	feedbackPub       ros.Publisher
	feedbackPubChan   chan ActionFeedback
	goalChan          chan *ros.MessageType
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

	server.statusFrequency = ros.NewRate(5.0)

	server.pubQueueSize = 50
	server.subQueueSize = 0

	server.goalSubChan = make(chan ActionGoal, 100)
	server.resultPubChan = make(chan ActionResult, 100)
	server.feedbackPubChan = make(chan ActionFeedback, 100)
	server.cancelSubChan = make(chan *actionlib_msgs.GoalID, 100)
	server.statusPubChan = make(chan *actionlib_msgs.GoalStatusArray, 100)

	server.goalCallback = goalCb
	server.cancelCallback = cancelCb

	server.lastCancel = ros.Now()
	server.statusList = list.New()
	server.statusListTimeout = ros.NewDuration(60, 0)

	server.shutdownChan = make(chan struct{}, 10)

	return server
}

func (as *defaultActionServer) Start() {
	as.init()
	go as.publishStatusRoutine()
	logger := as.node.Logger()
	defer func() { logger.Debug("defaultActionServer.start exit") }()

	for {
		logger.Debug("loop")
		select {
		case <-as.shutdownChan:
			return

		case goal := <-as.goalSubChan:
			as.internalGoalCallback(goal)

		case goalId := <-as.cancelSubChan:
			as.internalCancelCallback(goalId)

		case arr := <-as.statusPubChan:
			as.statusPub.Publish(arr)

		case fb := <-as.feedbackPubChan:
			as.feedbackPub.Publish(fb)

		case res := <-as.resultPubChan:
			as.resultPub.Publish(res)
		}
	}
}

// init intializes action publishers and subscribers
func (as *defaultActionServer) init() {
	node := as.node

	// queue sizes not implemented by ros.Node yet
	// subSize, _ := as.node.GetParam("actionlib_server_sub_queue_size")
	// as.subQueueSize = subSize.(int)

	// pubSize, _ := node.GetParam("actionlib_server_pub_queue_size")
	// as.pubQueueSize = pubSize.(int)

	as.goalSub = node.NewSubscriber(as.action+"/goal", as.actionType.GoalType(),
		func(goal ActionGoal) {
			as.goalSubChan <- goal
		})
	as.cancelSub = node.NewSubscriber(as.action+"/cancel", actionlib_msgs.MsgGoalID,
		func(goalId *actionlib_msgs.GoalID) {
			as.cancelSubChan <- goalId
		})
	as.resultPub = node.NewPublisher(as.action+"/result", as.actionType.ResultType())
	as.statusPub = node.NewPublisher(as.action+"/status", actionlib_msgs.MsgGoalStatusArray)
	as.feedbackPub = node.NewPublisher(as.action+"/feedback", as.actionType.FeedbackType())

	// hp := node.HasParam(as.action + "/status_frequency")
	// if hp {
	// 	logger.Warn("You're using the deprecated status_frequency parameter, please switch to actionlib_status_frequency.")
	// } else {
	// 	statFreq := node.SearchParam("actionlib_status_frequency", "5.0")
	// 	if statFreq != "" {
	// 		newros.Rate := node.GetParam(as.action+"/status_list_timeout", "5.0")
	// 		as.statusFrequency = Newros.Rate(newros.Rate.(float64))
	// 	}
	// }
}

// publishResult publishes action result message
func (as *defaultActionServer) publishResult(status actionlib_msgs.GoalStatus, result ros.Message) {
	msg := as.actionResult.NewMessage().(ActionResult)
	msg.SetHeader(std_msgs.Header{Stamp: ros.Now()})
	msg.SetStatus(status)
	msg.SetResult(result)
	as.resultPubChan <- msg
}

// publishFeedback publishes action feedback messages
func (as *defaultActionServer) publishFeedback(status actionlib_msgs.GoalStatus, feedback ros.Message) {
	msg := as.actionFeedback.NewMessage().(ActionFeedback)
	msg.SetHeader(std_msgs.Header{Stamp: ros.Now()})
	msg.SetStatus(status)
	msg.SetFeedback(feedback)
	as.feedbackPubChan <- msg
}

// publishStatus publishes action status messages
func (as *defaultActionServer) publishStatus() {
	var statArr []actionlib_msgs.GoalStatus
	if as.node.OK() {
		for e := as.statusList.Front(); e != nil; e = e.Next() {
			st := e.Value.(*status)
			destTime := st.destroyTime.Add(as.statusListTimeout)
			if !st.destroyTime.IsZero() && destTime.Cmp(ros.Now()) <= 0 {
				as.statusList.Remove(e)
				continue
			}

			statArr = append(statArr, st.getGoalStatus())
		}
	}

	arr := &actionlib_msgs.GoalStatusArray{}
	arr.Header.Stamp = ros.Now()
	arr.StatusList = statArr
	as.statusPubChan <- arr
}

func (as *defaultActionServer) publishStatusRoutine() {
	as.statusTimer = time.NewTicker(time.Second / 5.0)
	for {
		select {
		case <-as.statusTimer.C:
			as.publishStatus()
		}
	}
}

// internalCancelCallback recieves cancel message from client
func (as *defaultActionServer) internalCancelCallback(goalId *actionlib_msgs.GoalID) {
	logger := as.node.Logger()
	logger.Debug("action server has received a new cancel request")
}

// internalGoalCallback recieves the goals from client and checks if
// the goalID already exists in the status list. If not, it will call
// server's goalCallback with goal that was recieved from the client.
func (as *defaultActionServer) internalGoalCallback(goal ActionGoal) {
	logger := as.node.Logger()
	goalID := goal.GetGoalId()

	for e := as.statusList.Front(); e != nil; e = e.Next() {
		st := e.Value.(*status)
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
	as.statusList.PushBack(st)

	if !goalID.Stamp.IsZero() && goalID.Stamp.Cmp(as.lastCancel) <= 0 {
		// set_cancelled
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
