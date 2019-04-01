package actionlib

import (
	"actionlib_msgs"
	"fmt"

	"github.com/fetchrobotics/rosgo/ros"
)

type simpleActionServer struct {
	actionServer          *defaultActionServer
	currentGoal           *status
	nextGoal              *status
	newGoal               bool
	preemptRequest        bool
	newGoalPreemptRequest bool
	goalCallback          interface{}
	preemptCallback       interface{}
}

func newSimpleActionServer(node ros.Node, action string, actType ActionType, executeCb interface{}, start bool) *simpleActionServer {
	s := new(simpleActionServer)
	s.actionServer = newDefaultActionServer(node, action, actType, s.internalGoalCallback, s.internalPreemptCallback, start)
	s.newGoal = false
	s.preemptRequest = false
	s.newGoalPreemptRequest = false

	return s
}

func (s *simpleActionServer) Start() {
	s.actionServer.Start()
}

func (s *simpleActionServer) internalGoalCallback(goal ActionGoal) {
	goalStamp := goal.GetGoalId().Stamp
	nextGoalStamp := goal.GetGoalId().Stamp
	newGoal := newStatusWithActionGoal(s.actionServer, goal)

	if (s.currentGoal.getGoal() == nil || goalStamp.Cmp(s.currentGoal.getGoalId().Stamp) >= 0) &&
		(s.currentGoal.getGoal() == nil || nextGoalStamp.Cmp(s.currentGoal.getGoalId().Stamp) >= 0) {

		if s.nextGoal.getGoal() != nil && (s.currentGoal.getGoal() == nil || s.nextGoal != s.currentGoal) {
			s.nextGoal.setCancelled(nil, "This goal was canceled because another goal was received by the simple action server")
		}

		s.nextGoal = newGoal
		s.newGoal = true
		s.newGoalPreemptRequest = false

		if s.IsActive() {
			s.preemptRequest = true
			if s.preemptCallback != nil {
				// call preempt callback
			}
		}

		if s.goalCallback != nil {
			// call goal callback
		}

		// run execute call back asynchronously

	} else {
		newGoal.setCancelled(nil, "This goal was canceled because another goal was received by the simple action server")
	}

}

func (s *simpleActionServer) internalPreemptCallback() {

}

func (s *simpleActionServer) IsNewGoalAvailable() bool {
	return false
}

func (s *simpleActionServer) IsPreemptRequested() bool {
	return false
}

func (s *simpleActionServer) AcceptNewGoal() (ActionGoal, error) {
	if !s.newGoal || s.nextGoal.getGoal() == nil {
		return nil, fmt.Errorf("Attempting to accept the next goal when a new goal is not available")
	}

	if s.IsActive() && s.currentGoal.getGoal() != nil && s.currentGoal != s.nextGoal {
		s.currentGoal.setCancelled(nil, "This goal was canceled because another goal was received by the simple action server")
	}

	fmt.Println("Accepting new goal")
	s.currentGoal = s.nextGoal
	s.newGoal = false
	s.preemptRequest = false
	s.newGoalPreemptRequest = false

	s.currentGoal.setAccepted("This goal has been accepted by the simple action server")

	return s.currentGoal.getGoal(), nil
}

func (s simpleActionServer) IsActive() bool {
	if s.currentGoal.getGoalId().Id == "" {
		return false
	}

	status := s.currentGoal.getGoalStatus().Status
	if status == actionlib_msgs.ACTIVE || status == actionlib_msgs.PREEMPTING {
		return true
	}

	return false
}

func (s simpleActionServer) SetSucceeded(result ActionResult, text string) error {
	if result == nil {
		result = s.GetDefaultResult()
	}

	return s.currentGoal.setSucceeded(result, text)
}

func (s simpleActionServer) SetAborted(result ActionResult, text string) error {
	if result == nil {
		result = s.GetDefaultResult()
	}

	return s.currentGoal.setAborted(result, text)
}

func (s simpleActionServer) SetPreempted(result ActionResult, text string) error {
	if result == nil {
		result = s.GetDefaultResult()
	}

	return s.currentGoal.setCancelled(result, text)
}

func (s simpleActionServer) PublishFeedback(feedback ActionFeedback) error {
	return fmt.Errorf("Not implemented")
}

func (s simpleActionServer) GetDefaultResult() ActionResult {
	return s.actionServer.actionResult.NewMessage().(ActionResult)
}

func (s *simpleActionServer) RegisterGoalCallback(cb interface{}) {
	s.goalCallback = cb
}

func (s *simpleActionServer) RegisterPreemptCallback(cb interface{}) {
	s.preemptCallback = cb
}
