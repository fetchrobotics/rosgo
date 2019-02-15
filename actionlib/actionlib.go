package actionlib

//go:generate gengo msg actionlib_msgs/GoalStatusArray
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
	autoStart bool) ActionServer {
	return newDefaultActionServer(node, action, actionType, goalCb, cancelCb, autoStart)
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
