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
	cancelSub         ros.Subscriber
	statusPub         ros.Publisher
	resultPub         ros.Publisher
	feedbackPub       ros.Publisher
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

	server.goalCallback = goalCb
	server.cancelCallback = cancelCb

	server.lastCancel = ros.Now()
	server.statusListTimeout = ros.NewDuration(60, 0)

	server.shutdownChan = make(chan struct{}, 10)

	return server
}

func (as *defaultActionServer) RegisterGoalCallback(goalCb interface{}) {
	as.goalCallback = goalCb
}

func (as *defaultActionServer) RegisterCancelCallback(cancelCb interface{}) {
	as.cancelCallback = cancelCb
}

func (as *defaultActionServer) Start() {
	as.init()
	as.started = true
	as.publishStatus()
	logger := as.node.Logger()

	defer func() {
		logger.Debug("defaultActionServer.start exit")
	}()
}

// init intializes action publishers and subscribers
func (as *defaultActionServer) init() {
	node := as.node

	// queue sizes not implemented by ros.Node yet
	subSize, _ := as.node.GetParam("actionlib_server_sub_queue_size")
	as.subQueueSize = subSize.(int)

	pubSize, _ := node.GetParam("actionlib_server_pub_queue_size")
	as.pubQueueSize = pubSize.(int)

	as.goalSub = node.NewSubscriber(as.action+"/goal", as.actionType.GoalType(), as.internalGoalCallback)
	as.cancelSub = node.NewSubscriber(as.action+"/cancel", actionlib_msgs.MsgGoalID, as.cancelCallback)
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

	as.statusTimer = time.NewTicker(time.Second / 5.0)
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
func (as *defaultActionServer) publishStatus(status actionlib_msgs.GoalStatus) {

}

// internalCancelCallback recieves cancel message from client
func (as *defaultActionServer) internalCancelCallback() {
	logger := as.node.Logger()
	logger.Debug("action server has received a new cancel request")
}

// internalGoalCallback recieves the goals from client and checks if
// the goalID already exists in the status list. If not, it will call
// server's goalCallback with goal that was recieved from the client.
func (as *defaultActionServer) internalGoalCallback(goal ActionGoal) {
	goalID := goal.GetGoalId()
	for e := as.statusList.Front(); e != nil; e = e.Next() {
		st := e.Value.(*Status)
		if goalID.Id == st.goalID.Id {
			logger.Debugf("Goal %s was already in the status list with status %+v", goalID.Id, st.goalStatus)
			if st.goalStatus.Status == actionlib_msgs.RECALLING {
				st.goalStatus.Status = actionlib_msgs.RECALLED
				result := as.actionResultType.NewMessage()
				as.publishResult(st.goalStatus, result)
			}

			st.destroyTime = ros.Now()
			return
		}
	}

	st := NewStatus(goalID, as.statusListTimeout)
	as.statusList.PushBack(st)

	goal := goal.GetGoal()

	args := []reflect.Value{reflect.ValueOf(goal)}
	fun := reflect.ValueOf(as.goalCallback)
	numArgsNeeded := fun.Type().NumIn()

	if numArgsNeeded <= 1 {
		fun.Call(args[0:numArgsNeeded])
	}
}

func (as *defaultActionServer) internalStatusPublisher() {
	for as.node.OK() {
		var msg actionlib_msgs.GoalStatusArray
		as.statusPub.Publish(&msg)
		time.Sleep(time.Second)
	}
}

func (as *defaultActionServer) Shutdown() {
	as.shutdownChan <- struct{}{}
}

type Status struct {
	goalID      actionlib_msgs.GoalID
	goalStatus  actionlib_msgs.GoalStatus
	destroyTime ros.Time
}

func NewStatus(gID actionlib_msgs.GoalID, timeout ros.Duration) *Status {
	dTime := ros.Now()
	gStatus := actionlib_msgs.GoalStatus{
		GoalId: gID,
		Status: actionlib_msgs.PENDING,
	}
	st := &Status{
		goalID:      gID,
		goalStatus:  gStatus,
		destroyTime: dTime.Add(timeout),
	}

	return st
}
