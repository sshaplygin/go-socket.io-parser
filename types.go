package go_socketio_parser

// Type of packet
type Type byte

// Protocol packet type
const (
	// Connect type
	Connect Type = iota
	// Disconnect type
	Disconnect
	// Event type
	Event
	// Ack type
	Ack
	// Error type
	Error
	// BinaryEvent type
	BinaryEvent
	// BinaryAck type
	BinaryAck
)

// Header of packet
type Header struct {
	Type      Type
	NeedAck   bool
	ID        uint64
	Namespace string
}

// FrameType is the type of frames
type FrameType byte

// FrameType aliases
const (
	// TEXT is text type message.
	TEXT FrameType = iota
	// BINARY is binary type message.
	BINARY
)
