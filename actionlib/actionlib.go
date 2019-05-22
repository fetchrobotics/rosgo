package actionlib

import (
	"actionlib_msgs"
	"context"

	"github.com/fetchrobotics/rosgo/ros"
)

func NewActionClient(node ros.Node, action string, actionType ActionType) ActionClient {
	return newDefaultActionClient(node, action, actionType)
}

func NewActionServer(node ros.Node, action string, actionType ActionType, goalCb, cancelCb interface{}, autoStart bool) ActionServer {
	return newDefaultActionServer(node, action, actionType, goalCb, cancelCb, autoStart)
}

func NewSimpleActionServer(node ros.Node, action string, actionType ActionType, executeCb interface{}, autoStart bool) SimpleActionServer {
	return newSimpleActionServer(node, action, actionType, executeCb, autoStart)
}

func NewServerGoalHandlerWithGoal(as ActionServer, goal ActionGoal) ServerGoalHandler {
	return newServerGoalHandlerWithGoal(as, goal)
}

func NewServerGoalHandlerWithGoalId(as ActionServer, goalId *actionlib_msgs.GoalID) ServerGoalHandler {
	return newServerGoalHandlerWithGoalId(as, goalId)
}

type ActionServer interface {
	Start()
	Shutdown()
	PublishResult(actionlib_msgs.GoalStatus, ros.Message)
	PublishFeedback(actionlib_msgs.GoalStatus, ros.Message)
	PublishStatus()
	RegisterGoalCallback(interface{})
	RegisterCancelCallback(interface{})
}

type ActionClient interface {
	WaitForServer(context.Context) bool
	SendGoal(ros.Message, interface{}, interface{}) ClientGoalHandler
	CancelAllGoals()
	CancelAllGoalsBeforeTime(ros.Time)
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

type ClientGoalHandler interface {
	IsExpired() bool
	GetState() (uint8, error)
	GetTerminalState() (uint8, error)
	GetResult() (ros.Message, error)
	Resend() error
	Cancel() error
}

type ServerGoalHandler interface {
	SetAccepted(string) error
	SetCancelled(ActionResult, string) error
	SetRejected(ActionResult, string) error
	SetAborted(ActionResult, string) error
	SetSucceeded(ActionResult, string) error
	SetCancelRequested() bool
	PublishFeedback(ActionFeedback)
	GetGoal() ActionGoal
	GetGoalId() actionlib_msgs.GoalID
	GetGoalStatus() actionlib_msgs.GoalStatus
	Equal(ServerGoalHandler) bool
	NotEqual(ServerGoalHandler) bool
	Hash() uint32
	GetHandlerDestructionTime() ros.Time
	SetHandlerDestructionTime(ros.Time)
}
