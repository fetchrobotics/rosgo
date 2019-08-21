package actionlib

import (
	"actionlib_msgs"
	"fmt"
	"reflect"
	"time"

	"github.com/fetchrobotics/rosgo/ros"
)

const (
	SimpleStatePending uint8 = 0
	SimpleStateActive  uint8 = 1
	SimpleStateDone    uint8 = 2
)

type simpleActionClient struct {
	ac          *defaultActionClient
	simpleState uint8
	gh          ClientGoalHandler
	doneCb      interface{}
	activeCb    interface{}
	feedbackCb  interface{}
	doneChan    chan struct{}
	logger      ros.Logger
}

func newSimpleActionClient(node ros.Node, action string, actionType ActionType) *simpleActionClient {
	return &simpleActionClient{
		ac:          newDefaultActionClient(node, action, actionType),
		simpleState: SimpleStateDone,
		doneChan:    make(chan struct{}, 10),
		logger:      node.Logger(),
	}
}

func (sc *simpleActionClient) SendGoal(goal ros.Message, doneCb, activeCb, feedbackCb interface{}) {
	sc.StopTrackingGoal()
	sc.doneCb = doneCb
	sc.activeCb = activeCb
	sc.feedbackCb = feedbackCb

	sc.simpleState = SimpleStatePending
	sc.gh = sc.ac.SendGoal(goal, sc.transitionHandler, sc.feedbackHandler)
}

func (sc *simpleActionClient) SendGoalAndWait(goal ros.Message, executeTimeout, preeptTimeout ros.Duration) (uint8, error) {
	sc.SendGoal(goal, nil, nil, nil)
	if !sc.WaitForResult(executeTimeout) {
		sc.logger.Debug("Cancelling goal")
		sc.CancelGoal()
		if sc.WaitForResult(preeptTimeout) {
			sc.logger.Debug("Preempt finished within specified timeout")
		} else {
			sc.logger.Debug("Preempt did not finish within specified timeout")
		}
	}

	return sc.GetState()
}

func (sc *simpleActionClient) WaitForServer(timeout ros.Duration) bool {
	return sc.ac.WaitForServer(timeout)
}

func (sc *simpleActionClient) WaitForResult(timeout ros.Duration) bool {
	if sc.gh == nil {
		sc.logger.Errorf("Called WaitForResult when no goal exists")
		return false
	}

	waitStart := ros.Now()
	waitStart = waitStart.Add(timeout)

LOOP:
	for {
		select {
		case <-sc.doneChan:
			break LOOP
		case <-time.After(100 * time.Millisecond):
		}

		if waitStart.Cmp(ros.Now()) <= 0 {
			break LOOP
		}
	}

	return sc.simpleState == SimpleStateDone
}

func (sc *simpleActionClient) GetResult() (ros.Message, error) {
	if sc.gh == nil {
		return nil, fmt.Errorf("called get result when no goal running")
	}

	return sc.gh.GetResult()
}

func (sc *simpleActionClient) GetState() (uint8, error) {
	if sc.gh == nil {
		return actionlib_msgs.LOST, fmt.Errorf("called get state when no goal running")
	}

	status, err := sc.gh.GetGoalStatus()
	if err != nil {
		return actionlib_msgs.LOST, err
	}

	if status == actionlib_msgs.RECALLING {
		status = actionlib_msgs.PENDING
	} else if status == actionlib_msgs.PREEMPTING {
		status = actionlib_msgs.ACTIVE
	}

	return status, nil
}

func (sc *simpleActionClient) GetGoalStatusText() (string, error) {
	if sc.gh == nil {
		return "", fmt.Errorf("called GetGoalStatusText when no goal is running")
	}

	return sc.gh.GetGoalStatusText()
}

func (sc *simpleActionClient) CancelAllGoals() {
	sc.ac.CancelAllGoals()
}

func (sc *simpleActionClient) CancelAllGoalsBeforeTime(stamp ros.Time) {
	sc.ac.CancelAllGoalsBeforeTime(stamp)
}

func (sc *simpleActionClient) CancelGoal() error {
	if sc.gh == nil {
		return nil
	}

	return sc.gh.Cancel()
}

func (sc *simpleActionClient) StopTrackingGoal() {
	sc.gh = nil
}

func (sc *simpleActionClient) transitionHandler(gh ClientGoalHandler) {
	commState, err := gh.GetCommState()
	if err != nil {
		sc.logger.Errorf("Error getting CommState: %v", err)
		return
	}

	switch commState {
	case Active:
		switch sc.simpleState {
		case SimpleStatePending:
			sc.setSimpleState(SimpleStateActive)
			sc.runCallback("active", []reflect.Value{})

		case SimpleStateDone:
			sc.logger.Error("")
		}
	case Recalling:
		switch sc.simpleState {
		case SimpleStateActive, SimpleStateDone:
			sc.logger.Error("")
		}
	case Preempting:
		switch sc.simpleState {
		case SimpleStatePending:
			sc.setSimpleState(SimpleStateActive)
			sc.runCallback("active", []reflect.Value{})

		case SimpleStateDone:
			sc.logger.Error("")
		}
	case Done:
		switch sc.simpleState {
		case SimpleStatePending, SimpleStateActive:
			sc.setSimpleState(SimpleStateDone)

			status, err := sc.gh.GetGoalStatus()
			if err != nil {
				sc.logger.Errorf("Error getting status: %v", err)
			}

			result, err := sc.gh.GetResult()
			if err != nil {
				sc.logger.Errorf("Error getting result: %v", err)
			}

			sc.runCallback("done", []reflect.Value{reflect.ValueOf(status), reflect.ValueOf(result)})
			select {
			case sc.doneChan <- struct{}{}:
			default:
				sc.logger.Error("Error sending done notification. Channel full.")
			}

		case SimpleStateDone:
			sc.logger.Error("")
		}
	}
}

func (sc *simpleActionClient) feedbackHandler(gh ClientGoalHandler, msg ros.Message) {
	if sc.gh == nil || sc.gh != gh {
		return
	}

	sc.runCallback("feedback", []reflect.Value{reflect.ValueOf(msg)})
}

func (sc *simpleActionClient) setSimpleState(state uint8) {
	sc.simpleState = state
}

func (sc *simpleActionClient) runCallback(cbType string, args []reflect.Value) error {
	var callback interface{}
	switch cbType {
	case "active":
		callback = sc.activeCb
	case "feedback":
		callback = sc.feedbackCb
	case "done":
		callback = sc.doneCb
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
		return fmt.Errorf("unexepcted number of arguments for callback")
	}

	return nil
}
