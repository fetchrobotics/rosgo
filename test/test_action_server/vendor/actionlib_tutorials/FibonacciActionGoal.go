
// Automatically generated from the message definition "actionlib_tutorials/FibonacciActionGoal.msg"
package actionlib_tutorials
import (
    "bytes"
    "github.com/fetchrobotics/rosgo/ros"
	"std_msgs"
	"actionlib_msgs"
)


type _MsgFibonacciActionGoal struct {
    text string
    name string
    md5sum string
}

func (t *_MsgFibonacciActionGoal) Text() string {
    return t.text
}

func (t *_MsgFibonacciActionGoal) Name() string {
    return t.name
}

func (t *_MsgFibonacciActionGoal) MD5Sum() string {
    return t.md5sum
}

func (t *_MsgFibonacciActionGoal) NewMessage() ros.Message {
    m := new(FibonacciActionGoal)
	m.Header = std_msgs.Header{}
	m.GoalId = actionlib_msgs.GoalID{}
	m.Goal = FibonacciGoal{}
    return m
}

var (
    MsgFibonacciActionGoal = &_MsgFibonacciActionGoal {
        `Header header
actionlib_msgs/GoalID goal_id
actionlib_tutorials/FibonacciGoal goal
`,
        "actionlib_tutorials/FibonacciActionGoal",
        "006871c7fa1d0e3d5fe2226bf17b2a94",
    }
)

type FibonacciActionGoal struct {
	Header std_msgs.Header `rosmsg:"header:Header"`
	GoalId actionlib_msgs.GoalID `rosmsg:"goal_id:GoalID"`
	Goal FibonacciGoal `rosmsg:"goal:FibonacciGoal"`
}

func (m *FibonacciActionGoal) Type() ros.MessageType {
	return MsgFibonacciActionGoal
}

func (m *FibonacciActionGoal) Serialize(buf *bytes.Buffer) error {
    var err error = nil
    if err = m.Header.Serialize(buf); err != nil {
        return err
    }
    if err = m.GoalId.Serialize(buf); err != nil {
        return err
    }
    if err = m.Goal.Serialize(buf); err != nil {
        return err
    }
    return err
}


func (m *FibonacciActionGoal) Deserialize(buf *bytes.Reader) error {
    var err error = nil
    if err = m.Header.Deserialize(buf); err != nil {
        return err
    }
    if err = m.GoalId.Deserialize(buf); err != nil {
        return err
    }
    if err = m.Goal.Deserialize(buf); err != nil {
        return err
    }
    return err
}
