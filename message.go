package socket

type Message interface {
	Length() int
	Body() []byte
}

// Codec is the interface for message coder and decoder.
// Application programmer can define a custom codec themselves.
type Codec interface {
	Decode([]byte) (Message, error)
	Encode(Message) ([]byte, error)
}
