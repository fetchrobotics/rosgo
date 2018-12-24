
// Automatically generated from the message definition "actionlib_msgs/GoalID.msg"
// This message has been customized to be within the std ros library
package ros
import (
    "bytes"
    "encoding/binary"
)


type _MsgGoalID struct {
    text string
    name string
    md5sum string
}

func (t *_MsgGoalID) Text() string {
    return t.text
}

func (t *_MsgGoalID) Name() string {
    return t.name
}

func (t *_MsgGoalID) MD5Sum() string {
    return t.md5sum
}

func (t *_MsgGoalID) NewMessage() Message {
    m := new(GoalID)
	m.Stamp = Time{}
	m.Id = ""
    return m
}

var (
    MsgGoalID = &_MsgGoalID {
        `# The stamp should store the time at which this goal was requested.
# It is used by an action server when it tries to preempt all
# goals that were requested before a certain time
time stamp

# The id provides a way to associate feedback and
# result message with specific goal requests. The id
# specified must be unique.
string id

`,
        "actionlib_msgs/GoalID",
        "302881f31927c1df708a2dbab0e80ee8",
    }
)

type GoalID struct {
	Stamp Time `rosmsg:"stamp:time"`
	Id string `rosmsg:"id:string"`
}

func (m *GoalID) Type() MessageType {
	return MsgGoalID
}

func (m *GoalID) Serialize(buf *bytes.Buffer) error {
    var err error = nil
    binary.Write(buf, binary.LittleEndian, m.Stamp.Sec)
    binary.Write(buf, binary.LittleEndian, m.Stamp.NSec)
    binary.Write(buf, binary.LittleEndian, uint32(len([]byte(m.Id))))
    buf.Write([]byte(m.Id))
    return err
}


func (m *GoalID) Deserialize(buf *bytes.Reader) error {
    var err error = nil
    {
        if err = binary.Read(buf, binary.LittleEndian, &m.Stamp.Sec); err != nil {
            return err
        }

        if err = binary.Read(buf, binary.LittleEndian, &m.Stamp.NSec); err != nil {
            return err
        }
    }
    {
        var size uint32
        if err = binary.Read(buf, binary.LittleEndian, &size); err != nil {
            return err
        }
        data := make([]byte, int(size))
        if err = binary.Read(buf, binary.LittleEndian, data); err != nil {
            return err
        }
        m.Id = string(data)
    }
    return err
}
