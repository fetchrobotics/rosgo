
// Automatically generated from the message definition "actionlib_tutorials/FibonacciActionResult.msg"
package actionlib_tutorials
import (
    "bytes"
    "github.com/fetchrobotics/rosgo/ros"
	"std_msgs"
	"actionlib_msgs"
)


type _MsgFibonacciActionResult struct {
    text string
    name string
    md5sum string
}

func (t *_MsgFibonacciActionResult) Text() string {
    return t.text
}

func (t *_MsgFibonacciActionResult) Name() string {
    return t.name
}

func (t *_MsgFibonacciActionResult) MD5Sum() string {
    return t.md5sum
}

func (t *_MsgFibonacciActionResult) NewMessage() ros.Message {
    m := new(FibonacciActionResult)
	m.Header = std_msgs.Header{}
	m.Status = actionlib_msgs.GoalStatus{}
	m.Result = FibonacciResult{}
    return m
}

var (
    MsgFibonacciActionResult = &_MsgFibonacciActionResult {
        `Header header
actionlib_msgs/GoalStatus status
actionlib_tutorials/FibonacciResult result`,
        "actionlib_tutorials/FibonacciActionResult",
        "bee73a9fe29ae25e966e105f5553dd03",
    }
)

type FibonacciActionResult struct {
	Header std_msgs.Header `rosmsg:"header:Header"`
	Status actionlib_msgs.GoalStatus `rosmsg:"status:GoalStatus"`
	Result FibonacciResult `rosmsg:"result:FibonacciResult"`
}

func (m *FibonacciActionResult) Type() ros.MessageType {
	return MsgFibonacciActionResult
}

func (m *FibonacciActionResult) Serialize(buf *bytes.Buffer) error {
    var err error = nil
    if err = m.Header.Serialize(buf); err != nil {
        return err
    }
    if err = m.Status.Serialize(buf); err != nil {
        return err
    }
    if err = m.Result.Serialize(buf); err != nil {
        return err
    }
    return err
}


func (m *FibonacciActionResult) Deserialize(buf *bytes.Reader) error {
    var err error = nil
    if err = m.Header.Deserialize(buf); err != nil {
        return err
    }
    if err = m.Status.Deserialize(buf); err != nil {
        return err
    }
    if err = m.Result.Deserialize(buf); err != nil {
        return err
    }
    return err
}
