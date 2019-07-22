package main

//go:generate gengo action actionlib_tutorials/Fibonacci Fibonacci.action
import (
	"actionlib_tutorials"
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/fetchrobotics/rosgo/actionlib"
	"github.com/fetchrobotics/rosgo/ros"
)

var (
	resChan = make(chan *actionlib_tutorials.FibonacciActionResult, 10)
)

func transitionCallback(gh actionlib.ClientGoalHandler) {
	state, _ := gh.GetState()
	fmt.Printf("Transition Recieved: %v\n", state)
}

func feedbackCallback(gh actionlib.ClientGoalHandler, msg *actionlib_tutorials.FibonacciFeedback) {
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

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	interruptCh := make(chan os.Signal)
	signal.Notify(interruptCh, os.Interrupt)
	go func() {
		<-interruptCh
		cancel()
	}()

	as := actionlib.NewActionClient(node, "fibonacci", actionlib_tutorials.ActionFibonacci)
	if started := as.WaitForServer(ctx); !started {
		fmt.Println("server not started within timeout")
		os.Exit(-1)
	}

	goal := new(actionlib_tutorials.FibonacciGoal)
	goal.Order = 10
	as.SendGoal(goal, transitionCallback, feedbackCallback)
	result := <-resChan
	fmt.Printf("Got result: %v", result.Result)
}
