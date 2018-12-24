
// Automatically generated from the message definition "actionlib_tutorials/FibonacciFeedback.msg"
package actionlib_tutorials
import (
    "bytes"
    "encoding/binary"
    "github.com/fetchrobotics/rosgo/ros"
)


type _MsgFibonacciFeedback struct {
    text string
    name string
    md5sum string
}

func (t *_MsgFibonacciFeedback) Text() string {
    return t.text
}

func (t *_MsgFibonacciFeedback) Name() string {
    return t.name
}

func (t *_MsgFibonacciFeedback) MD5Sum() string {
    return t.md5sum
}

func (t *_MsgFibonacciFeedback) NewMessage() ros.Message {
    m := new(FibonacciFeedback)
	m.Sequence = []int32{}
    return m
}

var (
    MsgFibonacciFeedback = &_MsgFibonacciFeedback {
        `
#result definition
int32[] sequence
`,
        "actionlib_tutorials/FibonacciFeedback",
        "b81e37d2a31925a0e8ae261a8699cb79",
    }
)

type FibonacciFeedback struct {
	Sequence []int32 `rosmsg:"sequence:int32[]"`
}

func (m *FibonacciFeedback) Type() ros.MessageType {
	return MsgFibonacciFeedback
}

func (m *FibonacciFeedback) Serialize(buf *bytes.Buffer) error {
    var err error = nil
    binary.Write(buf, binary.LittleEndian, uint32(len(m.Sequence)))
    for _, e := range m.Sequence {
        binary.Write(buf, binary.LittleEndian, e)
    }
    return err
}


func (m *FibonacciFeedback) Deserialize(buf *bytes.Reader) error {
    var err error = nil
    {
        var size uint32
        if err = binary.Read(buf, binary.LittleEndian, &size); err != nil {
            return err
        }
        m.Sequence = make([]int32, int(size))
        for i := 0; i < int(size); i++ {
            if err = binary.Read(buf, binary.LittleEndian, &m.Sequence[i]); err != nil {
                return err
            }
        }
    }
    return err
}
