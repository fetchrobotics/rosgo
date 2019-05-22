package main

//go:generate gengo action actionlib_tutorials/Fibonacci Fibonacci.action
import (
	"actionlib_tutorials"
	"fmt"
	"os"
	"time"

	"github.com/fetchrobotics/rosgo/actionlib"
	"github.com/fetchrobotics/rosgo/ros"
)

type callbacks struct {
	as actionlib.SimpleActionServer
}

func (c *callbacks) createServer(node ros.Node) {
	c.as = actionlib.NewSimpleActionServer(
		node, "fibonacci",
		actionlib_tutorials.ActionFibonacci,
		c.executeCallback,
		false)
}

func (c *callbacks) executeCallback(msg *actionlib_tutorials.FibonacciActionGoal) {
	feed := actionlib_tutorials.FibonacciFeedback{}
	feed.Sequence = append(feed.Sequence, 0)
	feed.Sequence = append(feed.Sequence, 1)
	success := true

	for i := 1; i < int(msg.Goal.Order); i++ {
		if c.as.IsPreemptRequested() {
			success = false
			break
		}

		val := feed.Sequence[i] + feed.Sequence[i-1]
		feed.Sequence = append(feed.Sequence, val)

		c.as.PublishFeedback(&actionlib_tutorials.FibonacciActionFeedback{
			Feedback: feed,
		})

		time.Sleep(100 * time.Millisecond)
		fmt.Println("working")
	}

	if success {
		result := &actionlib_tutorials.FibonacciActionResult{}
		result.Result.Sequence = feed.Sequence
		if err := c.as.SetSucceeded(result, "goal"); err != nil {
			fmt.Error(err)
		}
	}
}

func main() {
	node, err := ros.NewNode("fibonacci", os.Args)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	defer node.Shutdown()
	node.Logger().SetSeverity(ros.LogLevelInfo)
	c := callbacks{}
	c.createServer(node)

	go c.as.Start()
	node.Spin()
}
