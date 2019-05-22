package actionlib

import (
	"actionlib_msgs"
	"fmt"

	"github.com/fetchrobotics/rosgo/ros"
)

type clientGoalHandler struct {
	actionClient *defaultActionClient
	stateMachine *clientStateMachine
}

func newClientGoalHandler(ac *defaultActionClient, sm *clientStateMachine) *clientGoalHandler {
	return &clientGoalHandler{
		actionClient: ac,
		stateMachine: sm,
	}
}

func (gh *clientGoalHandler) Shutdown(deleteFromManager bool) {
	gh.stateMachine = nil
	if deleteFromManager {
		gh.actionClient.DeleteGoalHandler(gh)
	}
}

func (gh *clientGoalHandler) IsExpired() bool {
	return gh.stateMachine == nil
}

func (gh *clientGoalHandler) GetStateMachine() *clientStateMachine {
	return gh.stateMachine
}

func (gh *clientGoalHandler) GetState() (uint8, error) {
	if gh.stateMachine == nil {
		return Lost, fmt.Errorf("trying to get state on an inactive ClientGoalHandler")
	}

	return gh.stateMachine.state, nil
}

func (gh *clientGoalHandler) GetTerminalState() (uint8, error) {
	if gh.stateMachine == nil {
		return 0, fmt.Errorf("trying to get goal status on inactive clientGoalHandler")
	}

	if gh.stateMachine.state != Done {
		gh.actionClient.logger.Warnf("Asking for terminal state when we are in state %v", gh.stateMachine.state)
	}

	// implement get status
	goalStatus := gh.stateMachine.latestGoalStatus.Status
	if goalStatus == actionlib_msgs.PREEMPTED ||
		goalStatus == actionlib_msgs.SUCCEEDED ||
		goalStatus == actionlib_msgs.ABORTED ||
		goalStatus == actionlib_msgs.REJECTED ||
		goalStatus == actionlib_msgs.RECALLED ||
		goalStatus == actionlib_msgs.LOST {

		return goalStatus, nil
	}

	gh.actionClient.logger.Warnf("Asking for terminal state when latest goal is in %v", goalStatus)
	return actionlib_msgs.LOST, nil
}

func (gh *clientGoalHandler) GetResult() (ros.Message, error) {
	if gh.stateMachine == nil {
		return nil, fmt.Errorf("trying to get goal status on inactive clientGoalHandler")
	}

	if gh.stateMachine.latestResult == nil {
		return nil, fmt.Errorf("trying to get result when no result has been recieved")
	}

	return gh.stateMachine.getResult(), nil
}

func (gh *clientGoalHandler) Resend() error {
	if gh.stateMachine == nil {
		return fmt.Errorf("trying to call resend on inactive client goal hanlder")
	}

	gh.actionClient.goalPub.Publish(gh.stateMachine.getActionGoal())
	return nil
}

func (gh *clientGoalHandler) Cancel() error {
	if gh.stateMachine == nil {
		return fmt.Errorf("trying to call cancel on inactive client goal hanlder")
	}

	cancelMsg := &actionlib_msgs.GoalID{
		Stamp: ros.Now(),
		Id:    gh.stateMachine.getActionGoal().GetGoalId().Id}

	gh.actionClient.goalPub.Publish(cancelMsg)

	return gh.stateMachine.transition(WaitingForCancelAck, gh)
}
