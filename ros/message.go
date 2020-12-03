package ros

import (
	"bytes"
)

type MessageType interface {
	Text() string
	MD5Sum() string
	Name() string
	NewMessage() Message
}

type Message interface {
	GetType() MessageType
	Serialize(buf *bytes.Buffer) error
	Deserialize(buf *Reader) error
}
