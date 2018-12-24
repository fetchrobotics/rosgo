
// Automatically generated from the message definition "actionlib_tutorials/FibonacciGoal.msg"
package actionlib_tutorials
import (
    "bytes"
    "encoding/binary"
    "github.com/fetchrobotics/rosgo/ros"
)


type _MsgFibonacciGoal struct {
    text string
    name string
    md5sum string
}

func (t *_MsgFibonacciGoal) Text() string {
    return t.text
}

func (t *_MsgFibonacciGoal) Name() string {
    return t.name
}

func (t *_MsgFibonacciGoal) MD5Sum() string {
    return t.md5sum
}

func (t *_MsgFibonacciGoal) NewMessage() ros.Message {
    m := new(FibonacciGoal)
	m.Order = 0
    return m
}

var (
    MsgFibonacciGoal = &_MsgFibonacciGoal {
        `#goal definition
int32 order
`,
        "actionlib_tutorials/FibonacciGoal",
        "6889063349a00b249bd1661df429d822",
    }
)

type FibonacciGoal struct {
	Order int32 `rosmsg:"order:int32"`
}

func (m *FibonacciGoal) Type() ros.MessageType {
	return MsgFibonacciGoal
}

func (m *FibonacciGoal) Serialize(buf *bytes.Buffer) error {
    var err error = nil
    binary.Write(buf, binary.LittleEndian, m.Order)
    return err
}


func (m *FibonacciGoal) Deserialize(buf *bytes.Reader) error {
    var err error = nil
    if err = binary.Read(buf, binary.LittleEndian, &m.Order); err != nil {
        return err
    }
    return err
}
