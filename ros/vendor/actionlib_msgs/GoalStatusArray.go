
// Automatically generated from the message definition "actionlib_msgs/GoalStatusArray.msg"
package actionlib_msgs
import (
    "bytes"
    "github.com/fetchrobotics/rosgo/ros"
	"std_msgs"
	"actionlib_msgs"
)


type _MsgGoalStatusArray struct {
    text string
    name string
    md5sum string
}

func (t *_MsgGoalStatusArray) Text() string {
    return t.text
}

func (t *_MsgGoalStatusArray) Name() string {
    return t.name
}

func (t *_MsgGoalStatusArray) MD5Sum() string {
    return t.md5sum
}

func (t *_MsgGoalStatusArray) NewMessage() ros.Message {
    m := new(GoalStatusArray)
	m.Header = std_msgs.Header{}
	m.StatusList = []actionlib_msgs.GoalStatus{}
    return m
}

var (
    MsgGoalStatusArray = &_MsgGoalStatusArray {
        `# Stores the statuses for goals that are currently being tracked
# by an action server
Header header
GoalStatus[] status_list

`,
        "actionlib_msgs/GoalStatusArray",
        "8b2b82f13216d0a8ea88bd3af735e619",
    }
)

type GoalStatusArray struct {
	Header std_msgs.Header `rosmsg:"header:Header"`
	StatusList []actionlib_msgs.GoalStatus `rosmsg:"status_list:GoalStatus[]"`
}

func (m *GoalStatusArray) Type() ros.MessageType {
	return MsgGoalStatusArray
}

func (m *GoalStatusArray) Serialize(buf *bytes.Buffer) error {
    var err error = nil
    if err = m.Header.Serialize(buf); err != nil {
        return err
    }
    binary.Write(buf, binary.LittleEndian, uint32(len(m.StatusList)))
    for _, e := range m.StatusList {
        if err = e.Serialize(buf); err != nil {
            return err
        }
    }
    return err
}


func (m *GoalStatusArray) Deserialize(buf *bytes.Reader) error {
    var err error = nil
    if err = m.Header.Deserialize(buf); err != nil {
        return err
    }
    {
        var size uint32
        if err = binary.Read(buf, binary.LittleEndian, &size); err != nil {
            return err
        }
        m.StatusList = make([]actionlib_msgs.GoalStatus, int(size))
        for i := 0; i < int(size); i++ {
            if err = m.StatusList[i].Deserialize(buf); err != nil {
                return err
            }
        }
    }
    return err
}
