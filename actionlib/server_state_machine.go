package actionlib

import (
	"actionlib_msgs"
	"fmt"
	"sync"
)

type serverStateMachine struct {
	goalStatus actionlib_msgs.GoalStatus
	mutex      sync.RWMutex
}

type Event uint8

const (
	CancelRequest Event = 1
	Cancel        Event = 2
	Reject        Event = 3
	Accept        Event = 4
	Succeed       Event = 5
	Abort         Event = 6
)

func newServerStateMachine(goalID actionlib_msgs.GoalID) *serverStateMachine {
	return &serverStateMachine{
		goalStatus: actionlib_msgs.GoalStatus{
			GoalId: goalID,
			Status: actionlib_msgs.PENDING,
		},
	}
}

func (sm *serverStateMachine) transition(event Event, text string) (actionlib_msgs.GoalStatus, error) {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	nextState := sm.goalStatus.Status

	switch sm.goalStatus.Status {
	case actionlib_msgs.PENDING:
		switch event {
		case Reject:
			nextState = actionlib_msgs.REJECTED
			break
		case CancelRequest:
			nextState = actionlib_msgs.RECALLING
			break
		case Cancel:
			nextState = actionlib_msgs.RECALLED
			break
		case Accept:
			nextState = actionlib_msgs.ACTIVE
			break
		default:
			return sm.goalStatus, fmt.Errorf("invalid transition Event")
		}

	case actionlib_msgs.RECALLING:
		switch event {
		case Reject:
			nextState = actionlib_msgs.REJECTED
			break
		case Cancel:
			nextState = actionlib_msgs.RECALLED
			break
		case Accept:
			nextState = actionlib_msgs.PREEMPTING
			break
		default:
			return sm.goalStatus, fmt.Errorf("invalid transition Event")
		}

	case actionlib_msgs.ACTIVE:
		switch event {
		case Succeed:
			nextState = actionlib_msgs.SUCCEEDED
			break
		case CancelRequest:
			nextState = actionlib_msgs.PREEMPTING
			break
		case Cancel:
			nextState = actionlib_msgs.PREEMPTED
			break
		case Abort:
			nextState = actionlib_msgs.ABORTED
			break
		default:
			return sm.goalStatus, fmt.Errorf("invalid transition Event")
		}

	case actionlib_msgs.PREEMPTING:
		switch event {
		case Succeed:
			nextState = actionlib_msgs.SUCCEEDED
			break
		case Cancel:
			nextState = actionlib_msgs.PREEMPTED
			break
		case Abort:
			nextState = actionlib_msgs.ABORTED
			break
		default:
			return sm.goalStatus, fmt.Errorf("invalid transition Event")
		}
	case actionlib_msgs.REJECTED:
		break
	case actionlib_msgs.RECALLED:
		break
	case actionlib_msgs.SUCCEEDED:
		break
	case actionlib_msgs.PREEMPTED:
		break
	case actionlib_msgs.ABORTED:
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
