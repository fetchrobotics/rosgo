
// Automatically generated from the message definition "actionlib_tutorials/FibonacciActionFeedback.msg"
package actionlib_tutorials
import (
    "bytes"
    "github.com/fetchrobotics/rosgo/ros"
	"std_msgs"
	"actionlib_msgs"
)


type _MsgFibonacciActionFeedback struct {
    text string
    name string
    md5sum string
}

func (t *_MsgFibonacciActionFeedback) Text() string {
    return t.text
}

func (t *_MsgFibonacciActionFeedback) Name() string {
    return t.name
}

func (t *_MsgFibonacciActionFeedback) MD5Sum() string {
    return t.md5sum
}

func (t *_MsgFibonacciActionFeedback) NewMessage() ros.Message {
    m := new(FibonacciActionFeedback)
	m.Header = std_msgs.Header{}
	m.Status = actionlib_msgs.GoalStatus{}
	m.Feedback = FibonacciFeedback{}
    return m
}

var (
    MsgFibonacciActionFeedback = &_MsgFibonacciActionFeedback {
        `Header header
actionlib_msgs/GoalStatus status
actionlib_tutorials/FibonacciFeedback feedback`,
        "actionlib_tutorials/FibonacciActionFeedback",
        "73b8497a9f629a31c0020900e4148f07",
    }
)

type FibonacciActionFeedback struct {
	Header std_msgs.Header `rosmsg:"header:Header"`
	Status actionlib_msgs.GoalStatus `rosmsg:"status:GoalStatus"`
	Feedback FibonacciFeedback `rosmsg:"feedback:FibonacciFeedback"`
}

func (m *FibonacciActionFeedback) Type() ros.MessageType {
	return MsgFibonacciActionFeedback
}

func (m *FibonacciActionFeedback) Serialize(buf *bytes.Buffer) error {
    var err error = nil
    if err = m.Header.Serialize(buf); err != nil {
        return err
    }
    if err = m.Status.Serialize(buf); err != nil {
        return err
    }
    if err = m.Feedback.Serialize(buf); err != nil {
        return err
    }
    return err
}


func (m *FibonacciActionFeedback) Deserialize(buf *bytes.Reader) error {
    var err error = nil
    if err = m.Header.Deserialize(buf); err != nil {
        return err
    }
    if err = m.Status.Deserialize(buf); err != nil {
        return err
    }
    if err = m.Feedback.Deserialize(buf); err != nil {
        return err
    }
    return err
}
