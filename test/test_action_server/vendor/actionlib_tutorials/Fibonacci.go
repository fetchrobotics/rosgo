
// Automatically generated from the message definition "actionlib_tutorials/Fibonacci.action"
package actionlib_tutorials
import (
    "github.com/fetchrobotics/rosgo/ros"
)

// Service type metadata
type _ActionFibonacci struct {
    name string
    md5sum string
    text string
    goalType ros.MessageType
    feedbackType ros.MessageType
    resultType ros.MessageType
}

func (t *_ActionFibonacci) Name() string { return t.name }
func (t *_ActionFibonacci) MD5Sum() string { return t.md5sum }
func (t *_ActionFibonacci) Text() string { return t.text }
func (t *_ActionFibonacci) GoalType() ros.MessageType { return t.goalType }
func (t *_ActionFibonacci) FeedbackType() ros.MessageType { return t.feedbackType }
func (t *_ActionFibonacci) ResultType() ros.MessageType { return t.resultType }
func (t *_ActionFibonacci) NewAction() ros.Action {
    return new(Fibonacci)
}

var (
    ActionFibonacci = &_ActionFibonacci {
        "actionlib_tutorials/Fibonacci",
        "00a5fc530b1d04d07f7b99ac88531c80",
        ``,
        MsgFibonacciActionGoal,
        MsgFibonacciActionFeedback,
        MsgFibonacciActionResult,
    }
)


type Fibonacci struct {
    Goal FibonacciActionGoal
    Feedback FibonacciActionFeedback
    Result FibonacciActionResult
}

func (s *Fibonacci) GoalMessage() ros.Message { return &s.Goal }
func (s *Fibonacci) FeedbackMessage() ros.Message { return &s.Feedback }
func (s *Fibonacci) ResultMessage() ros.Message { return &s.Result }
