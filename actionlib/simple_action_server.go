package actionlib

import (
	"actionlib_msgs"
	"fmt"
	"reflect"
	"sync"
	"time"

	"github.com/fetchrobotics/rosgo/ros"
)

type simpleActionServer struct {
	actionServer          *defaultActionServer
	currentGoal           *status
	nextGoal              *status
	newGoal               bool
	preemptRequest        bool
	newGoalPreemptRequest bool
	goalMutex             sync.RWMutex
	logger                ros.Logger
	goalCallback          interface{}
	preemptCallback       interface{}
	executeCb             interface{}
	executorCh            chan struct{}
}

func newSimpleActionServer(node ros.Node, action string, actType ActionType, executeCb interface{}, start bool) *simpleActionServer {
	s := new(simpleActionServer)
	s.actionServer = newDefaultActionServer(node, action, actType, s.internalGoalCallback, s.internalPreemptCallback, start)
	s.newGoal = false
	s.preemptRequest = false
	s.newGoalPreemptRequest = false
	s.executeCb = executeCb
	s.logger = node.Logger()
	s.executorCh = make(chan struct{}, 10)
	return s
}

func (s *simpleActionServer) Start() {
	if s.executeCb != nil {
		go s.goalExecutor()
	}
	s.actionServer.Start()
}

func (s *simpleActionServer) internalGoalCallback(goal ActionGoal) {
	s.logger.Infof("Simple action server received new goal with id %s", goal.GetGoalId().Id)

	goalStamp := goal.GetGoalId().Stamp
	nextGoalStamp := goal.GetGoalId().Stamp
	newGoal := newStatusWithActionGoal(s.actionServer, goal)

	s.goalMutex.Lock()
	defer s.goalMutex.Unlock()

	if (s.currentGoal == nil || goalStamp.Cmp(s.currentGoal.getGoalID().Stamp) >= 0) &&
		(s.currentGoal == nil || nextGoalStamp.Cmp(s.currentGoal.getGoalID().Stamp) >= 0) {

		if s.nextGoal != nil && (s.currentGoal == nil || s.nextGoal != s.currentGoal) {
			s.nextGoal.setCancelled(nil,
				"This goal was canceled because another goal was received by the simple action server")
		}

		s.nextGoal = newGoal
		s.newGoal = true
		s.newGoalPreemptRequest = false
		args := []reflect.Value{reflect.ValueOf(goal)}

		if s.IsActive() {
			s.preemptRequest = true
			if err := s.runCallback("preempt", args); err != nil {
				s.logger.Error(err)
			}
		}

		if err := s.runCallback("goal", args); err != nil {
			s.logger.Error(err)
		}

		// notify executor that a new goal is available
		s.executorCh <- struct{}{}
	} else {
		newGoal.setCancelled(nil,
			"This goal was canceled because another goal was received by the simple action server")
	}
}

func (s *simpleActionServer) internalPreemptCallback(preempt ActionGoal) {
	s.goalMutex.Lock()
	defer s.goalMutex.Unlock()

	s.logger.Infof("Simple action server received preempt call for goal with id %s",
		preempt.GetGoalId().Id)

	if preempt == s.currentGoal.getGoal() {
		s.preemptRequest = true
		args := []reflect.Value{reflect.ValueOf(preempt)}
		if err := s.runCallback("preempt", args); err != nil {
			s.logger.Error(err)
		}
	} else {
		s.newGoalPreemptRequest = true
	}
}

func (s *simpleActionServer) IsNewGoalAvailable() bool {
	s.goalMutex.Lock()
	defer s.goalMutex.Unlock()

	return s.newGoal
}

func (s *simpleActionServer) IsPreemptRequested() bool {
	s.goalMutex.Lock()
	defer s.goalMutex.Unlock()

	return s.preemptRequest
}

func (s *simpleActionServer) AcceptNewGoal() (ActionGoal, error) {
	s.goalMutex.Lock()
	defer s.goalMutex.Unlock()

	if !s.newGoal || s.nextGoal.getGoal() == nil {
		return nil, fmt.Errorf("Attempting to accept the next goal when a new goal is not available")
	}

	if s.IsActive() && s.currentGoal.getGoal() != nil && s.currentGoal != s.nextGoal {
		s.currentGoal.setCancelled(nil,
			"This goal was canceled because another goal was received by the simple action server")
	}

	fmt.Println("Accepting new goal")
	s.currentGoal = s.nextGoal
	s.newGoal = false
	s.preemptRequest = false
	s.newGoalPreemptRequest = false

	s.currentGoal.setAccepted("This goal has been accepted by the simple action server")

	return s.currentGoal.getGoal(), nil
}

func (s *simpleActionServer) IsActive() bool {
	if s.currentGoal == nil || s.currentGoal.getGoalID().Id == "" {
		return false
	}

	status := s.currentGoal.getGoalStatus().Status
	if status == actionlib_msgs.ACTIVE || status == actionlib_msgs.PREEMPTING {
		return true
	}

	return false
}

func (s *simpleActionServer) SetSucceeded(result ActionResult, text string) error {
	s.goalMutex.Lock()
	defer s.goalMutex.Unlock()

	if result == nil {
		result = s.GetDefaultResult()
	}

	return s.currentGoal.setSucceeded(result, text)
}

func (s *simpleActionServer) SetAborted(result ActionResult, text string) error {
	s.goalMutex.Lock()
	defer s.goalMutex.Unlock()

	if result == nil {
		result = s.GetDefaultResult()
	}

	return s.currentGoal.setAborted(result, text)
}

func (s *simpleActionServer) SetPreempted(result ActionResult, text string) error {
	s.goalMutex.Lock()
	defer s.goalMutex.Unlock()

	if result == nil {
		result = s.GetDefaultResult()
	}

	return s.currentGoal.setCancelled(result, text)
}

func (s *simpleActionServer) PublishFeedback(feedback ActionFeedback) {
	s.goalMutex.Lock()
	defer s.goalMutex.Unlock()

	s.currentGoal.publishFeedback(feedback)
}

func (s *simpleActionServer) GetDefaultResult() ActionResult {
	return s.actionServer.actionResult.NewMessage().(ActionResult)
}

func (s *simpleActionServer) RegisterGoalCallback(cb interface{}) error {
	if s.executeCb != nil {
		return fmt.Errorf("execute callback if present: not registering goal callback")
	}
	s.goalCallback = cb
	return nil
}

func (s *simpleActionServer) RegisterPreemptCallback(cb interface{}) {
	s.preemptCallback = cb
}

func (s *simpleActionServer) goalExecutor() {
	intervalCh := time.NewTicker(1 * time.Second)
	defer intervalCh.Stop()

	for s.actionServer.node.OK() {
		select {
		case <-s.executorCh:
			if err := s.execute(); err != nil {
				s.logger.Error(err)
				return
			}

		case <-intervalCh.C:
			if err := s.execute(); err != nil {
				s.logger.Error(err)
				return
			}
		}
	}
}

func (s *simpleActionServer) execute() error {
	if s.IsActive() {
		return fmt.Errorf("Should never reach this code with an active goal")
	}

	if s.IsNewGoalAvailable() {
		goal, err := s.AcceptNewGoal()
		if err != nil {
			return err
		}

		if s.executeCb == nil {
			return fmt.Errorf("Execute callback must exist. This is a bug in SimpleActionServer")
		}

		args := []reflect.Value{reflect.ValueOf(goal)}
		if err := s.runCallback("execute", args); err != nil {
			return err
		}

		if s.IsActive() {
			s.logger.Warn("Your executeCallback did not set the goal to a terminal status.  " +
				"This is a bug in your ActionServer implementation. Fix your code!  " +
				"For now, the ActionServer will set this goal to aborted")
			if err := s.SetAborted(nil, ""); err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *simpleActionServer) runCallback(cbType string, args []reflect.Value) error {
	var callback interface{}
	switch cbType {
	case "goal":
		callback = s.goalCallback
	case "preempt":
		callback = s.preemptCallback
	case "execute":
		callback = s.executeCb
	default:
		return fmt.Errorf("unknown callback type called")
	}

	if callback == nil {
		return nil
	}

	fun := reflect.ValueOf(callback)
	numArgsNeeded := fun.Type().NumIn()

	if numArgsNeeded <= 1 {
		fun.Call(args[0:numArgsNeeded])
	} else {
		return fmt.Errorf("Unexepcted number of arguments for callback")
	}

	return nil
}
