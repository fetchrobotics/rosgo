package ros

//go:generate gengo msg actionlib_msgs/GoalStatusArray
//go:generate gengo msg actionlib_msgs/GoalStatus
//go:generate gengo msg actionlib_msgs/GoalID
//go:generate gengo msg std_msgs/Header

import (
	//	"fmt"
	//"reflect"
	"time"
)

type actionResult struct {
	action Action
	err    error
}

type remoteClientActionError struct {
	session *remoteClientSession
	err     error
}

type defaultGoal struct {
	handler interface{}
	ready   chan bool
}

type defaultActionServer struct {
	node         *defaultNode
	action       string
	actionType   ActionType
	goalHandler  []*defaultGoal
	shutdownChan chan bool
	goalChan     chan *MessageType
	goalSub      Subscriber
	cancelSub    Subscriber
	resultPub    Publisher
	feedbackPub  Publisher
	statusPub    Publisher
}

func newDefaultActionServer(node *defaultNode, action string, actType ActionType, handler interface{}, routines int, start bool) *defaultActionServer {
	//logger := node.logger
	server := new(defaultActionServer)
	server.node = node
	gh := new(defaultGoal)
	gh.handler = handler
	gh.ready = make(chan bool, 1)
	server.goalHandler = append(server.goalHandler, gh)
	server.action = action
	server.actionType = actType
	server.shutdownChan = make(chan bool, 1)
	server.goalSub = node.NewSubscriber(action+"/goal", actType.GoalType(), server.internalGoalCallback)
	server.cancelSub = node.NewSubscriber(action+"/cancel", MsgGoalID, server.internalGoalCallback)
	server.resultPub = node.NewPublisher(action+"/result", actType.ResultType())
	server.statusPub = node.NewPublisher(action+"/status", MsgGoalStatusArray)
	server.feedbackPub = node.NewPublisher(action+"/feedback", actType.FeedbackType())
	if start {
		go server.start()
	}
	return server
}

func (s *defaultActionServer) internalStatusPublisher() {
	for s.node.OK() {
		s.node.SpinOnce()
		var msg GoalStatusArray
		s.statusPub.Publish(&msg)
		time.Sleep(time.Second)
	}
}

func (s *defaultActionServer) internalGoalCallback(m *MessageType) {
	s.goalChan <- m
	// args := []reflect.Value{reflect.ValueOf(g), reflect.ValueOf(s.actionType.GoalType)}
	// fun := reflect.ValueOf(s.handler[0].handler)
	// fun.Call(args)
}

func (s *defaultActionServer) Shutdown() {
	s.shutdownChan <- true
}

// event loop
func (s *defaultActionServer) start() {
	logger := s.node.logger
	logger.Debugf("action server '%s' start listen %s.", s.action)
	s.node.waitGroup.Add(1)
	go s.internalStatusPublisher()
	defer func() {
		logger.Debug("defaultActionServer.start exit")
		s.node.waitGroup.Done()
	}()
}
