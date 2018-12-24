
// Automatically generated from the message definition "actionlib_tutorials/FibonacciResult.msg"
package actionlib_tutorials
import (
    "bytes"
    "encoding/binary"
    "github.com/fetchrobotics/rosgo/ros"
)


type _MsgFibonacciResult struct {
    text string
    name string
    md5sum string
}

func (t *_MsgFibonacciResult) Text() string {
    return t.text
}

func (t *_MsgFibonacciResult) Name() string {
    return t.name
}

func (t *_MsgFibonacciResult) MD5Sum() string {
    return t.md5sum
}

func (t *_MsgFibonacciResult) NewMessage() ros.Message {
    m := new(FibonacciResult)
	m.Sequence = []int32{}
    return m
}

var (
    MsgFibonacciResult = &_MsgFibonacciResult {
        `
#feedback
int32[] sequence
`,
        "actionlib_tutorials/FibonacciResult",
        "b81e37d2a31925a0e8ae261a8699cb79",
    }
)

type FibonacciResult struct {
	Sequence []int32 `rosmsg:"sequence:int32[]"`
}

func (m *FibonacciResult) Type() ros.MessageType {
	return MsgFibonacciResult
}

func (m *FibonacciResult) Serialize(buf *bytes.Buffer) error {
    var err error = nil
    binary.Write(buf, binary.LittleEndian, uint32(len(m.Sequence)))
    for _, e := range m.Sequence {
        binary.Write(buf, binary.LittleEndian, e)
    }
    return err
}


func (m *FibonacciResult) Deserialize(buf *bytes.Reader) error {
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
