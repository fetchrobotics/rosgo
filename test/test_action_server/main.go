package main

//go:generate gengo action actionlib_tutorials/Fibonacci Fibonacci.action
import (
	"actionlib_tutorials"
	"fmt"
	"os"

	"github.com/fetchrobotics/rosgo/actionlib"
	"github.com/fetchrobotics/rosgo/ros"
)

func goalCallback(msg *actionlib_tutorials.FibonacciActionGoal) {
	fmt.Printf("Goal Recieved: %s\n", msg.Goal)
}

func cancelCallback(msg *actionlib_msgs.GoalID) {
	fmt.Printf("Cancel Recieved: %s\n", msg.Id)
}

func main() {
	node, err := ros.NewNode("talker", os.Args)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	defer node.Shutdown()
	node.Logger().SetSeverity(ros.LogLevelDebug)

	actionlib.NewActionServer(
		node, "fibonacci",
		actionlib_tutorials.ActionFibonacci,
		goalCallback, cancelCallback, false)
	node.Spin()
}
