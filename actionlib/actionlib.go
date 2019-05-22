package actionlib

import (
	"github.com/fetchrobotics/rosgo/ros"
)

func NewActionClient(node ros.Node, action string, actionType ActionType) ActionClient {
	return newDefaultActionClient(node, action, actionType)
}

func NewActionServer(
	node ros.Node,
	action string,
	actionType ActionType,
	goalCb interface{},
	cancelCb interface{},
	autoStart bool,
) ActionServer {
	return newDefaultActionServer(node, action, actionType, goalCb, cancelCb, autoStart)
}

func NewSimpleActionServer(
	node ros.Node,
	action string,
	actionType ActionType,
	executeCb interface{},
	autoStart bool,
) SimpleActionServer {
	return newSimpleActionServer(node, action, actionType, executeCb, autoStart)
}

type ActionServer interface {
	Start()
	Shutdown()
	RegisterGoalCallback(interface{})
	RegisterCancelCallback(interface{})
}

type ActionClient interface {
	WaitForServer()
	SendGoal(ros.Message)
	WaitForResult()
	GetResult() ros.Message
	Shutdown()
}

type SimpleActionServer interface {
	Start()

	IsNewGoalAvailable() bool
	IsPreemptRequested() bool
	IsActive() bool

	SetSucceeded(ActionResult, string) error
	SetAborted(ActionResult, string) error
	SetPreempted(ActionResult, string) error

	AcceptNewGoal() (ActionGoal, error)
	PublishFeedback(ActionFeedback)
	GetDefaultResult() ActionResult

	RegisterGoalCallback(interface{}) error
	RegisterPreemptCallback(interface{})
}
