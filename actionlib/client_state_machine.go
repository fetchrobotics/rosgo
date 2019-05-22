package actionlib

import (
	"actionlib_msgs"
	"container/list"
	"fmt"
	"reflect"

	"github.com/fetchrobotics/rosgo/ros"
)

const (
	WaitingForGoalAck   uint8 = 0
	Pending             uint8 = 1
	Active              uint8 = 2
	WaitingForResult    uint8 = 3
	WaitingForCancelAck uint8 = 4
	Recalling           uint8 = 5
	Preempting          uint8 = 6
	Done                uint8 = 7
	Lost                uint8 = 8
)

type clientStateMachine struct {
	actionClient     *defaultActionClient
	actionGoal       ActionGoal
	actionType       ActionType
	actionGoalID     string
	state            uint8
	transitionCb     interface{}
	feedbackCb       interface{}
	latestGoalStatus *actionlib_msgs.GoalStatus
	latestResult     ActionResult
	logger           ros.Logger
}

func newClientStateMachine(ac *defaultActionClient, actionGoal ActionGoal, actionType ActionType, transitionCb, feedbackCb interface{}) *clientStateMachine {
	return &clientStateMachine{
		actionClient: ac,
		actionGoal:   actionGoal,
		actionType:   actionType,
		state:        WaitingForGoalAck,
		logger:       ac.logger,
		transitionCb: transitionCb,
		feedbackCb:   feedbackCb,
		latestGoalStatus: &actionlib_msgs.GoalStatus{
			Status: actionlib_msgs.PENDING,
		},
	}
}

func (sm *clientStateMachine) getActionGoal() ActionGoal {
	return sm.actionGoal
}

func (sm *clientStateMachine) getState() uint8 {
	return sm.state
}

func (sm *clientStateMachine) getGoalStatus() actionlib_msgs.GoalStatus {
	return *sm.latestGoalStatus
}

func (sm *clientStateMachine) getResult() ros.Message {
	if sm.latestResult != nil {
		return sm.latestResult.GetResult()
	}

	return nil
}

func (sm *clientStateMachine) findGoalStatus(statusArr *actionlib_msgs.GoalStatusArray) *actionlib_msgs.GoalStatus {
	var status actionlib_msgs.GoalStatus
	for _, st := range statusArr.StatusList {
		if st.GoalId.Id == sm.actionGoalID {
			status = st
			break
		}
	}

	return &status
}

func (sm *clientStateMachine) updateFeedback(af ActionFeedback, gh *clientGoalHandler) {
	if sm.actionGoal.GetGoalId().Id != af.GetStatus().GoalId.Id {
		return
	}

	if sm.feedbackCb != nil && sm.state != Done {
		fun := reflect.ValueOf(sm.feedbackCb)
		args := []reflect.Value{reflect.ValueOf(gh), reflect.ValueOf(af.GetFeedback())}
		numArgsNeeded := fun.Type().NumIn()

		if numArgsNeeded == 2 {
			fun.Call(args)
		}
	}
}

func (sm *clientStateMachine) updateResult(result ActionResult, gh *clientGoalHandler) error {
	if sm.actionGoal.GetGoalId().Id != result.GetStatus().GoalId.Id {
		return nil
	}

	status := result.GetStatus()
	sm.latestGoalStatus = &status
	sm.latestResult = result

	if sm.state == WaitingForGoalAck ||
		sm.state == WaitingForCancelAck ||
		sm.state == Pending ||
		sm.state == Active ||
		sm.state == WaitingForResult ||
		sm.state == Recalling ||
		sm.state == Preempting {

		statusArr := new(actionlib_msgs.GoalStatusArray)
		statusArr.StatusList = append(statusArr.StatusList, result.GetStatus())
		if err := sm.updateStatus(statusArr, gh); err != nil {
			return err
		}
		if err := sm.transition(Done, gh); err != nil {
			fmt.Print("got resut", err)
			return err
		}
		return nil
	} else if sm.state == Done {
		return fmt.Errorf("Got a result when we are in the `DONE` state")
	} else {
		return fmt.Errorf("Unknown state %v", sm.state)
	}
}

func (sm *clientStateMachine) updateStatus(statusArr *actionlib_msgs.GoalStatusArray, gh *clientGoalHandler) error {
	if sm.state == Done {
		return nil
	}

	status := sm.findGoalStatus(statusArr)
	if status == nil {
		if sm.state != WaitingForGoalAck ||
			sm.state != WaitingForResult ||
			sm.state != Done {

			sm.logger.Warn("Transitioning goal to `Lost`")
			sm.latestGoalStatus.Status = Lost
			return sm.transition(Done, gh)
		}

		return nil
	}

	sm.latestGoalStatus = status
	return sm.transition(status.Status, gh)
}

func (sm *clientStateMachine) transition(goalStatus uint8, gh *clientGoalHandler) error {
	nextStates, err := sm.getTransition(goalStatus)
	if err != nil {
		return err
	}

	for e := nextStates.Front(); e != nil; e = e.Next() {
		state := e.Value.(uint8)
		sm.state = state

		if sm.transitionCb != nil {
			fun := reflect.ValueOf(sm.transitionCb)
			args := []reflect.Value{reflect.ValueOf(gh)}
			numArgsNeeded := fun.Type().NumIn()

			if numArgsNeeded <= 1 {
				fun.Call(args[:numArgsNeeded])
			} else {
				return fmt.Errorf("Unexepcted number of arguments for transition callback")
			}
		}
	}

	return nil
}

func (sm *clientStateMachine) getTransition(goalStatus uint8) (stateList list.List, err error) {
	switch sm.state {
	case WaitingForGoalAck:
		switch goalStatus {
		case actionlib_msgs.PENDING:
			stateList.PushBack(Pending)
			break
		case actionlib_msgs.ACTIVE:
			stateList.PushBack(Active)
			break
		case actionlib_msgs.REJECTED:
			stateList.PushBack(Pending)
			stateList.PushBack(WaitingForCancelAck)
			break
		case actionlib_msgs.RECALLING:
			stateList.PushBack(Pending)
			stateList.PushBack(Recalling)
			break
		case actionlib_msgs.RECALLED:
			stateList.PushBack(Pending)
			stateList.PushBack(WaitingForResult)
			break
		case actionlib_msgs.PREEMPTED:
			stateList.PushBack(Active)
			stateList.PushBack(Preempting)
			stateList.PushBack(WaitingForResult)
			break
		case actionlib_msgs.SUCCEEDED:
			stateList.PushBack(Active)
			stateList.PushBack(WaitingForResult)
			break
		case actionlib_msgs.ABORTED:
			stateList.PushBack(Active)
			stateList.PushBack(WaitingForResult)
			break
		case actionlib_msgs.PREEMPTING:
			stateList.PushBack(Active)
			stateList.PushBack(Preempting)
			break
		}
		break

	case Pending:
		switch goalStatus {
		case actionlib_msgs.PENDING:
			break
		case actionlib_msgs.ACTIVE:
			stateList.PushBack(Active)
			break
		case actionlib_msgs.REJECTED:
			stateList.PushBack(WaitingForResult)
			break
		case actionlib_msgs.RECALLING:
			stateList.PushBack(Recalling)
			break
		case actionlib_msgs.RECALLED:
			stateList.PushBack(Recalling)
			stateList.PushBack(WaitingForResult)
			break
		case actionlib_msgs.PREEMPTED:
			stateList.PushBack(Active)
			stateList.PushBack(Preempting)
			stateList.PushBack(WaitingForResult)
			break
		case actionlib_msgs.SUCCEEDED:
			stateList.PushBack(Active)
			stateList.PushBack(WaitingForResult)
			break
		case actionlib_msgs.ABORTED:
			stateList.PushBack(Active)
			stateList.PushBack(WaitingForResult)
			break
		case actionlib_msgs.PREEMPTING:
			stateList.PushBack(Active)
			stateList.PushBack(Preempting)
			break
		}
		break
	case Active:
		switch goalStatus {
		case actionlib_msgs.PENDING:
			err = fmt.Errorf("[actionlib] Invalid transition from Active to Pending")
			break
		case actionlib_msgs.ACTIVE:
			break
		case actionlib_msgs.REJECTED:
			err = fmt.Errorf("[actionlib] Invalid transition from Active to Rejected")
			break
		case actionlib_msgs.RECALLING:
			err = fmt.Errorf("[actionlib] Invalid transition from Active to Recalling")
			break
		case actionlib_msgs.RECALLED:
			err = fmt.Errorf("[actionlib] Invalid transition from Active to Recalled")
			break
		case actionlib_msgs.PREEMPTED:
			stateList.PushBack(Preempting)
			stateList.PushBack(WaitingForResult)
			break
		case actionlib_msgs.SUCCEEDED:
			stateList.PushBack(WaitingForResult)
			break
		case actionlib_msgs.ABORTED:
			stateList.PushBack(WaitingForResult)
			break
		case actionlib_msgs.PREEMPTING:
			stateList.PushBack(Preempting)
			break
		}
		break
	case WaitingForResult:
		switch goalStatus {
		case actionlib_msgs.PENDING:
			err = fmt.Errorf("[actionlib] Invalid transition from WaitingForResult to Pending")
			break
		case actionlib_msgs.ACTIVE:
			break
		case actionlib_msgs.REJECTED:
			break
		case actionlib_msgs.RECALLING:
			err = fmt.Errorf("[actionlib] Invalid transition from WaitingForResult to Recalling")
			break
		case actionlib_msgs.RECALLED:
			break
		case actionlib_msgs.PREEMPTED:
			break
		case actionlib_msgs.SUCCEEDED:
			break
		case actionlib_msgs.ABORTED:
			break
		case actionlib_msgs.PREEMPTING:
			err = fmt.Errorf("[actionlib] Invalid transition from WaitingForResult to Preempting")
			break
		}
		break
	case WaitingForCancelAck:
		switch goalStatus {
		case actionlib_msgs.PENDING:
			break
		case actionlib_msgs.ACTIVE:
			break
		case actionlib_msgs.REJECTED:
			stateList.PushBack(WaitingForResult)
			break
		case actionlib_msgs.RECALLING:
			stateList.PushBack(Recalling)
			break
		case actionlib_msgs.RECALLED:
			stateList.PushBack(Recalling)
			stateList.PushBack(WaitingForResult)
			break
		case actionlib_msgs.PREEMPTED:
			stateList.PushBack(Preempting)
			stateList.PushBack(WaitingForResult)
			break
		case actionlib_msgs.SUCCEEDED:
			stateList.PushBack(Recalling)
			stateList.PushBack(WaitingForResult)
			break
		case actionlib_msgs.ABORTED:
			stateList.PushBack(Recalling)
			stateList.PushBack(WaitingForResult)
			break
		case actionlib_msgs.PREEMPTING:
			stateList.PushBack(Preempting)
			break
		}
		break
	case Recalling:
		switch goalStatus {
		case actionlib_msgs.PENDING:
			err = fmt.Errorf("[actionlib] Invalid transition from Recalling to Pending")
			break
		case actionlib_msgs.ACTIVE:
			err = fmt.Errorf("[actionlib] Invalid transition from Recalling to Active")
			break
		case actionlib_msgs.REJECTED:
			stateList.PushBack(WaitingForResult)
			break
		case actionlib_msgs.RECALLING:
			break
		case actionlib_msgs.RECALLED:
			stateList.PushBack(WaitingForResult)
			break
		case actionlib_msgs.PREEMPTED:
			stateList.PushBack(Preempting)
			stateList.PushBack(WaitingForResult)
			break
		case actionlib_msgs.SUCCEEDED:
			stateList.PushBack(Preempting)
			stateList.PushBack(WaitingForResult)
			break
		case actionlib_msgs.ABORTED:
			stateList.PushBack(Preempting)
			stateList.PushBack(WaitingForResult)
			break
		case actionlib_msgs.PREEMPTING:
			stateList.PushBack(Preempting)
			break
		}
		break
	case Preempting:
		switch goalStatus {
		case actionlib_msgs.PENDING:
			err = fmt.Errorf("[actionlib] Invalid transition from Preempting to Pending")
			break
		case actionlib_msgs.ACTIVE:
			err = fmt.Errorf("[actionlib] Invalid transition from Preempting to Active")
			break
		case actionlib_msgs.REJECTED:
			err = fmt.Errorf("[actionlib] Invalid transition from Preempting to Rejected")
			break
		case actionlib_msgs.RECALLING:
			err = fmt.Errorf("[actionlib] Invalid transition from Preempting to Recalling")
			break
		case actionlib_msgs.RECALLED:
			err = fmt.Errorf("[actionlib] Invalid transition from Preempting to Recalled")
			break
		case actionlib_msgs.PREEMPTED:
			stateList.PushBack(WaitingForResult)
			break
		case actionlib_msgs.SUCCEEDED:
			stateList.PushBack(WaitingForResult)
			break
		case actionlib_msgs.ABORTED:
			stateList.PushBack(WaitingForResult)
			break
		case actionlib_msgs.PREEMPTING:
			break
		}
		break
	case Done:
		switch goalStatus {
		case actionlib_msgs.PENDING:
			err = fmt.Errorf("[actionlib] Invalid transition from Done to Pending")
			break
		case actionlib_msgs.ACTIVE:
			err = fmt.Errorf("[actionlib] Invalid transition from Done to Active")
			break
		case actionlib_msgs.REJECTED:
			break
		case actionlib_msgs.RECALLING:
			err = fmt.Errorf("[actionlib] Invalid transition from Done to Recalling")
			break
		case actionlib_msgs.RECALLED:
			break
		case actionlib_msgs.PREEMPTED:
			break
		case actionlib_msgs.SUCCEEDED:
			break
		case actionlib_msgs.ABORTED:
			break
		case actionlib_msgs.PREEMPTING:
			err = fmt.Errorf("[actionlib] Invalid transition from Done to Preempting")
			break
		}
		break
	}

	return
}

func stateToString(state uint8) string {
	switch state {
	case WaitingForGoalAck:
		return "WAITING_FOR_GOAL_ACK"
	case Pending:
		return "PENDING"
	case Active:
		return "ACTIVE"
	case WaitingForResult:
		return "WAITING_FOR_RESULT"
	case WaitingForCancelAck:
		return "WAITING_FOR_CANCEL_ACK"
	case Recalling:
		return "RECALLING"
	case Preempting:
		return "PREEMPTING"
	case Done:
		return "DONE"
	case Lost:
		return "LOST"
	default:
		return "UNKNOWN_STATE"
	}
}
