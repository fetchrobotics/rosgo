
// Automatically generated from the message definition "actionlib_msgs/GoalStatus.msg"
// This message has been customized to run within the standard library
package ros
import (
    "bytes"
    "encoding/binary"
)
const (
	PENDING uint8 = 0
	ACTIVE uint8 = 1
	PREEMPTED uint8 = 2
	SUCCEEDED uint8 = 3
	ABORTED uint8 = 4
	REJECTED uint8 = 5
	PREEMPTING uint8 = 6
	RECALLING uint8 = 7
	RECALLED uint8 = 8
	LOST uint8 = 9
)


type _MsgGoalStatus struct {
    text string
    name string
    md5sum string
}

func (t *_MsgGoalStatus) Text() string {
    return t.text
}

func (t *_MsgGoalStatus) Name() string {
    return t.name
}

func (t *_MsgGoalStatus) MD5Sum() string {
    return t.md5sum
}

func (t *_MsgGoalStatus) NewMessage() Message {
    m := new(GoalStatus)
	m.GoalId = GoalID{}
	m.Status = 0
	m.Text = ""
    return m
}

var (
    MsgGoalStatus = &_MsgGoalStatus {
        `GoalID goal_id
uint8 status
uint8 PENDING         = 0   # The goal has yet to be processed by the action server
uint8 ACTIVE          = 1   # The goal is currently being processed by the action server
uint8 PREEMPTED       = 2   # The goal received a cancel request after it started executing
                            #   and has since completed its execution (Terminal State)
uint8 SUCCEEDED       = 3   # The goal was achieved successfully by the action server (Terminal State)
uint8 ABORTED         = 4   # The goal was aborted during execution by the action server due
                            #    to some failure (Terminal State)
uint8 REJECTED        = 5   # The goal was rejected by the action server without being processed,
                            #    because the goal was unattainable or invalid (Terminal State)
uint8 PREEMPTING      = 6   # The goal received a cancel request after it started executing
                            #    and has not yet completed execution
uint8 RECALLING       = 7   # The goal received a cancel request before it started executing,
                            #    but the action server has not yet confirmed that the goal is canceled
uint8 RECALLED        = 8   # The goal received a cancel request before it started executing
                            #    and was successfully cancelled (Terminal State)
uint8 LOST            = 9   # An action client can determine that a goal is LOST. This should not be
                            #    sent over the wire by an action server

#Allow for the user to associate a string with GoalStatus for debugging
string text

`,
        "actionlib_msgs/GoalStatus",
        "d388f9b87b3c471f784434d671988d4a",
    }
)

type GoalStatus struct {
	GoalId GoalID `rosmsg:"goal_id:GoalID"`
	Status uint8 `rosmsg:"status:uint8"`
	Text string `rosmsg:"text:string"`
}

func (m *GoalStatus) Type() MessageType {
	return MsgGoalStatus
}

func (m *GoalStatus) Serialize(buf *bytes.Buffer) error {
    var err error = nil
    if err = m.GoalId.Serialize(buf); err != nil {
        return err
    }
    binary.Write(buf, binary.LittleEndian, m.Status)
    binary.Write(buf, binary.LittleEndian, uint32(len([]byte(m.Text))))
    buf.Write([]byte(m.Text))
    return err
}


func (m *GoalStatus) Deserialize(buf *bytes.Reader) error {
    var err error = nil
    if err = m.GoalId.Deserialize(buf); err != nil {
        return err
    }
    if err = binary.Read(buf, binary.LittleEndian, &m.Status); err != nil {
        return err
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
        m.Text = string(data)
    }
    return err
}
