package ros

import (
	"bytes"
)

// MessageType defines the interface that a ros message type should implement
type MessageType interface {
	Text() string
	MD5Sum() string
	Name() string
	NewMessage() Message
}

// Message defines the interface that a ros message should implement
type Message interface {
	GetType() MessageType
	Serialize(buf *bytes.Buffer) error
	Deserialize(buf *bytes.Reader) error
}
