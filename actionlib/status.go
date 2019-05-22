package actionlib

import (
	"actionlib_msgs"
	"fmt"
	"hash/fnv"
	"strings"

	"github.com/fetchrobotics/rosgo/ros"
)

type status struct {
	actionServer *defaultActionServer
	goal         ActionGoal
	goalStatus   *actionlib_msgs.GoalStatus
	destroyTime  ros.Time
}

func newStatusWithActionGoal(as *defaultActionServer, goal ActionGoal) *status {
	st := new(status)
	st.actionServer = as
	st.goal = goal

	timeNow := ros.Now()
	st.destroyTime = timeNow.Add(as.statusListTimeout)

	st.goalStatus = &actionlib_msgs.GoalStatus{}
	st.goalStatus.GoalId = goal.GetGoalId()
	st.goalStatus.Status = actionlib_msgs.PENDING

	if st.goalStatus.GoalId.Id == "" {
		var strs []string
		strs = append(strs, "nodeName"+"-")
		strs = append(strs, string(as.statusList.Len())+"-")
		strs = append(strs, string(timeNow.Sec)+"."+string(timeNow.NSec))

		st.goalStatus.GoalId.Id = strings.Join(strs, "")
		st.goalStatus.GoalId.Stamp = timeNow
	}

	if st.goalStatus.GoalId.Stamp.IsZero() {
		st.goalStatus.GoalId.Stamp = timeNow
	}

	return st
}

func newStatusWithGoalStatus(as *defaultActionServer, goalStatus actionlib_msgs.GoalStatus) *status {
	st := new(status)
	st.actionServer = as
	st.goalStatus = &goalStatus
	timeNow := ros.Now()
	st.destroyTime = timeNow.Add(as.statusListTimeout)
	return st
}

func (st *status) setAccepted(text string) error {
	if st.goal == nil {
		return fmt.Errorf("attempt to set status on an uninitialized status handler")
	}

	status := st.goalStatus.Status

	if status == actionlib_msgs.PENDING {
		st.goalStatus.Status = actionlib_msgs.ACTIVE
		st.goalStatus.Text = text
		st.actionServer.publishStatus()
	} else if status == actionlib_msgs.RECALLING {
		st.goalStatus.Status = actionlib_msgs.PREEMPTING
		st.goalStatus.Text = text
		st.actionServer.publishStatus()
	} else {
		return fmt.Errorf("to transition to an active state, the goal must be in a pending"+
			"or recalling state, it is currently in state: %d", st.goalStatus.Status)
	}

	return nil
}

func (st *status) setCancelled(result ActionResult, text string) error {
	if st.goal == nil {
		return fmt.Errorf("attempt to set status on an uninitialized status handler")
	}

	status := st.goalStatus.Status

	if status == actionlib_msgs.PENDING || status == actionlib_msgs.RECALLING {
		st.goalStatus.Status = actionlib_msgs.RECALLED
		st.goalStatus.Text = text
		st.destroyTime = ros.Now()
		st.actionServer.publishStatus()
	} else if status == actionlib_msgs.ACTIVE || status == actionlib_msgs.PREEMPTING {
		st.goalStatus.Status = actionlib_msgs.PREEMPTED
		st.goalStatus.Text = text
		st.destroyTime = ros.Now()
		st.actionServer.publishStatus()
	} else {
		return fmt.Errorf("to transition to an active state, the goal must be in a pending"+
			"recalling, active or prempting state, it is currently in state: %d", st.goalStatus.Status)
	}

	return nil
}

func (st *status) setRejected(result ActionResult, text string) error {
	if st.goal == nil {
		return fmt.Errorf("attempt to set status on an uninitialized status handler")
	}

	status := st.goalStatus.Status

	if status == actionlib_msgs.PENDING || status == actionlib_msgs.RECALLING {
		st.goalStatus.Status = actionlib_msgs.REJECTED
		st.goalStatus.Text = text
		st.destroyTime = ros.Now()
		st.actionServer.publishStatus()
	} else {
		return fmt.Errorf("to transition to an active state, the goal must be in a pending"+
			"or recalling state, it is currently in state: %d", st.goalStatus.Status)
	}

	return nil
}

func (st *status) setAborted(result ActionResult, text string) error {
	if st.goal == nil {
		return fmt.Errorf("attempt to set status on an uninitialized status handler")
	}

	status := st.goalStatus.Status

	if status == actionlib_msgs.PREEMPTING || status == actionlib_msgs.ACTIVE {
		st.goalStatus.Status = actionlib_msgs.ABORTED
		st.goalStatus.Text = text
		st.destroyTime = ros.Now()
		st.actionServer.publishStatus()
	} else {
		return fmt.Errorf("to transition to an active state, the goal must be in a prempting"+
			"or active state, it is currently in state: %d", st.goalStatus.Status)
	}

	return nil
}

func (st *status) setSucceeded(result ActionResult, text string) error {
	if st.goal == nil {
		return fmt.Errorf("attempt to set status on an uninitialized status handler")
	}

	status := st.goalStatus.Status

	if status == actionlib_msgs.PREEMPTING || status == actionlib_msgs.ACTIVE {
		st.goalStatus.Status = actionlib_msgs.SUCCEEDED
		st.goalStatus.Text = text
		st.destroyTime = ros.Now()
		st.actionServer.publishResult(st.getGoalStatus(), result.GetResult())
	} else {
		return fmt.Errorf("to transition to an active state, the goal must be in a prempting"+
			"or active state, it is currently in state: %d", st.goalStatus.Status)
	}

	return nil
}

func (st *status) setCancelRequested() bool {
	if st.goal == nil {
		fmt.Errorf("attempt to set status on an uninitialized status handler")
		return false
	}

	status := st.goalStatus.Status

	if status == actionlib_msgs.PENDING {
		st.goalStatus.Status = actionlib_msgs.RECALLING
		st.destroyTime = ros.Now()
		st.actionServer.publishStatus()
	}

	if status == actionlib_msgs.ACTIVE {
		st.goalStatus.Status = actionlib_msgs.PREEMPTING
		st.destroyTime = ros.Now()
		st.actionServer.publishStatus()
	}

	return true
}

func (st *status) publishFeedback(feedback ActionFeedback) {
	st.actionServer.publishFeedback(st.getGoalStatus(), feedback.GetFeedback())
}

func (st *status) getGoal() ActionGoal {
	if st.goal != nil {
		return st.goal
	}

	return nil
}

func (st *status) getGoalID() actionlib_msgs.GoalID {
	if st.goal != nil {
		return st.goalStatus.GoalId
	}

	return actionlib_msgs.GoalID{}
}

func (st *status) getGoalStatus() actionlib_msgs.GoalStatus {
	if st.goal != nil || st.goalStatus != nil {
		return *st.goalStatus
	}

	return actionlib_msgs.GoalStatus{}
}

func (st *status) equal(other *status) bool {
	if st.goal == nil || other.goal == nil {
		return false
	}

	return st.getGoalID().Id == other.getGoalID().Id
}

func (st *status) notEqual(other *status) bool {
	return !st.equal(other)
}

func (st *status) hash() uint32 {
	id := st.getGoalID().Id
	h := fnv.New32a()
	h.Write([]byte(id))
	return h.Sum32()
}
