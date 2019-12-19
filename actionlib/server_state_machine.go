package actionlib

import (
	"actionlib_msgs"
	"fmt"
	"sync"
)

type Event uint8

const (
	CancelRequest Event = iota + 1
	Cancel
	Reject
	Accept
	Succeed
	Abort
)

func (e Event) String() string {
	switch e {
	case CancelRequest:
		return "CANCEL_REQUEST"
	case Cancel:
		return "CANCEL"
	case Reject:
		return "REJECT"
	case Accept:
		return "ACCEPT"
	case Succeed:
		return "SUCCEED"
	case Abort:
		return "ABORT"
	default:
		return "UNKNOWN"
	}
}

type serverStateMachine struct {
	goalStatus actionlib_msgs.GoalStatus
	mutex      sync.RWMutex
}

func newServerStateMachine(goalID actionlib_msgs.GoalID) *serverStateMachine {
	return &serverStateMachine{
		goalStatus: actionlib_msgs.GoalStatus{
			GoalId: goalID,
			Status: actionlib_msgs.GoalStatus_PENDING,
		},
	}
}

func (sm *serverStateMachine) transition(event Event, text string) (actionlib_msgs.GoalStatus, error) {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	nextState := sm.goalStatus.Status

	switch sm.goalStatus.Status {
	case actionlib_msgs.GoalStatus_PENDING:
		switch event {
		case Reject:
			nextState = actionlib_msgs.GoalStatus_REJECTED
			break
		case CancelRequest:
			nextState = actionlib_msgs.GoalStatus_RECALLING
			break
		case Cancel:
			nextState = actionlib_msgs.GoalStatus_RECALLED
			break
		case Accept:
			nextState = actionlib_msgs.GoalStatus_ACTIVE
			break
		default:
			return sm.goalStatus, fmt.Errorf("invalid transition Event")
		}

	case actionlib_msgs.GoalStatus_RECALLING:
		switch event {
		case Reject:
			nextState = actionlib_msgs.GoalStatus_REJECTED
			break
		case Cancel:
			nextState = actionlib_msgs.GoalStatus_RECALLED
			break
		case Accept:
			nextState = actionlib_msgs.GoalStatus_PREEMPTING
			break
		default:
			return sm.goalStatus, fmt.Errorf("invalid transition Event")
		}

	case actionlib_msgs.GoalStatus_ACTIVE:
		switch event {
		case Succeed:
			nextState = actionlib_msgs.GoalStatus_SUCCEEDED
			break
		case CancelRequest:
			nextState = actionlib_msgs.GoalStatus_PREEMPTING
			break
		case Cancel:
			nextState = actionlib_msgs.GoalStatus_PREEMPTED
			break
		case Abort:
			nextState = actionlib_msgs.GoalStatus_ABORTED
			break
		default:
			return sm.goalStatus, fmt.Errorf("invalid transition Event")
		}

	case actionlib_msgs.GoalStatus_PREEMPTING:
		switch event {
		case Succeed:
			nextState = actionlib_msgs.GoalStatus_SUCCEEDED
			break
		case Cancel:
			nextState = actionlib_msgs.GoalStatus_PREEMPTED
			break
		case Abort:
			nextState = actionlib_msgs.GoalStatus_ABORTED
			break
		default:
			return sm.goalStatus, fmt.Errorf("invalid transition Event")
		}
	case actionlib_msgs.GoalStatus_REJECTED:
		break
	case actionlib_msgs.GoalStatus_RECALLED:
		break
	case actionlib_msgs.GoalStatus_SUCCEEDED:
		break
	case actionlib_msgs.GoalStatus_PREEMPTED:
		break
	case actionlib_msgs.GoalStatus_ABORTED:
		break
	default:
		return sm.goalStatus, fmt.Errorf("invalid state")
	}

	sm.goalStatus.Status = nextState
	sm.goalStatus.Text = text

	return sm.goalStatus, nil
}

func (sm *serverStateMachine) getStatus() actionlib_msgs.GoalStatus {
	sm.mutex.RLock()
	defer sm.mutex.RUnlock()

	return sm.goalStatus
}
