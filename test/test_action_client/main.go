package main

//go:generate gengo action actionlib_tutorials/Fibonacci Fibonacci.action
import (
	"actionlib_tutorials"
	"fmt"
	"os"
	"os/signal"

	"github.com/fetchrobotics/rosgo/actionlib"
	"github.com/fetchrobotics/rosgo/ros"
)

var (
	resChan = make(chan *actionlib_tutorials.FibonacciResult, 10)
)

func doneCallback() {
	fmt.Println("Done!")
	resChan <- &actionlib_tutorials.FibonacciResult{}
}

func activeCallback() {
	fmt.Println("Active")
}

func feedbackCallback(msg *actionlib_tutorials.FibonacciFeedback) {
	fmt.Printf("Feeback Recieved: %v\n", msg.Sequence)
}

func main() {
	node, err := ros.NewNode("fibonacci-client", os.Args)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	node.Logger().SetSeverity(ros.LogLevelInfo)
	go node.Spin()
	defer node.Shutdown()

	interruptCh := make(chan os.Signal)
	signal.Notify(interruptCh, os.Interrupt)
	go func() {
		<-interruptCh
	}()

	as := actionlib.NewSimpleActionClient(node, "fibonacci", actionlib_tutorials.ActionFibonacci)
	if started := as.WaitForServer(ros.NewDuration(60, 0)); !started {
		fmt.Println("server not started within timeout")
		os.Exit(-1)
	}

	goal := new(actionlib_tutorials.FibonacciGoal)
	goal.Order = 10
	as.SendGoal(goal, doneCallback, activeCallback, feedbackCallback)
	<-resChan
	res, _ := as.GetResult()
	fmt.Printf("Got result: %v", res)
}
