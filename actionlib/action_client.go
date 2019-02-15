package actionlib

import "github.com/fetchrobotics/rosgo/ros"

type defaultActionClient struct {
}

func newDefaultActionClient(node ros.Node, action string, actType ActionType) *defaultActionClient {
	return nil
}

func (ac *defaultActionClient) WaitForServer() {

}

func (ac *defaultActionClient) SendGoal(msg ros.Message) {

}

func (ac *defaultActionClient) WaitForResult() {

}

func (ac *defaultActionClient) GetResult() ros.Message {
	return nil
}

func (ac *defaultActionClient) Shutdown() {

}
