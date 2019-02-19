package main

//go:generate gengo action actionlib_tutorials/Fibonacci Fibonacci.action
import (
	"actionlib_msgs"
	"actionlib_tutorials"
	"fmt"
	"os"

	"github.com/fetchrobotics/rosgo/actionlib"
	"github.com/fetchrobotics/rosgo/ros"
)

func goalCallback(msg *actionlib_tutorials.FibonacciGoal) {
	fmt.Printf("Goal Recieved: %s !!!", msg)
}

func cancelCallback(msg *actionlib_msgs.GoalID) {
	fmt.Printf("Cancel Recieved: %s\n", msg.Id)
}

func main() {
	node, err := ros.NewNode("fibonacci", os.Args)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	defer node.Shutdown()
	node.Logger().SetSeverity(ros.LogLevelInfo)

	as := actionlib.NewActionServer(
		node, "fibonacci",
		actionlib_tutorials.ActionFibonacci,
		goalCallback, cancelCallback, false)

	go as.Start()
	node.Spin()
}
