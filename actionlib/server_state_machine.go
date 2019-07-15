package actionlib

import (
	"actionlib_msgs"
	"fmt"
	"sync"
)

type serverStateMachine struct {
	st    actionlib_msgs.GoalStatus
	mutex sync.RWMutex
}

const (
	CancelRequest uint8 = 1
	Cancel        uint8 = 2
	Reject        uint8 = 3
	Accept        uint8 = 4
	Succeed       uint8 = 5
	Abort         uint8 = 6
)

func newServerStateMachine(goalID actionlib_msgs.GoalID) *serverStateMachine {
	sm := new(serverStateMachine)
	sm.st.GoalId = goalID
	sm.st.Status = actionlib_msgs.PENDING

	return sm
}

func (sm *serverStateMachine) transition(event uint8, text string) (actionlib_msgs.GoalStatus, error) {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	nextState := sm.st.Status

	switch sm.st.Status {
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
			return sm.st, fmt.Errorf("invalid transition Event")
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
			return sm.st, fmt.Errorf("invalid transition Event")
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
			return sm.st, fmt.Errorf("invalid transition Event")
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
			return sm.st, fmt.Errorf("invalid transition Event")
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
		return sm.st, fmt.Errorf("invalid state")
	}

	sm.st.Status = nextState
	sm.st.Text = text

	return sm.st, nil
}

func (sm *serverStateMachine) getStatus() actionlib_msgs.GoalStatus {
	sm.mutex.RLock()
	defer sm.mutex.RUnlock()

	return sm.st
}
