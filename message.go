package socket

import (
	"io"
)

type Message interface {
	MessageNumber() int32
	Serialize() ([]byte, error)
}

// Codec is the interface for message coder and decoder.
// Application programmer can define a custom codec themselves.
type Codec interface {
	Decode(io.Reader) (Message, error)
	Encode(Message) ([]byte, error)
}
