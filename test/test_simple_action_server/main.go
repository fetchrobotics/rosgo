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
	as     actionlib.SimpleActionServer
	logger ros.Logger
}

func (c *callbacks) createServer(node ros.Node) {
	c.as = actionlib.NewSimpleActionServer(
		node, "fibonacci",
		actionlib_tutorials.ActionFibonacci,
		c.executeCallback,
		false)
	go c.as.Start()
}

func (c *callbacks) executeCallback(msg *actionlib_tutorials.FibonacciActionGoal) {
	feed := actionlib_tutorials.FibonacciFeedback{}
	feed.Sequence = append(feed.Sequence, 0)
	feed.Sequence = append(feed.Sequence, 1)
	success := true

	for i := 1; i < int(msg.Goal.Order); i++ {
		if c.as.IsPreemptRequested() {
			success = false
			if err := c.as.SetPreempted(nil, ""); err != nil {
				c.logger.Fatal(err)
			}
			break
		}

		val := feed.Sequence[i] + feed.Sequence[i-1]
		feed.Sequence = append(feed.Sequence, val)

		c.as.PublishFeedback(&actionlib_tutorials.FibonacciActionFeedback{})
		time.Sleep(10 * time.Millisecond)
	}

	if success {
		ar := &actionlib_tutorials.FibonacciActionResult{}
		ar.Result.Sequence = feed.Sequence
		if err := c.as.SetSucceeded(ar, "goal"); err != nil {
			c.logger.Fatal(err)
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
	node.Spin()
}
