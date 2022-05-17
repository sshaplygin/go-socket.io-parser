package go_socketio_parser

// Type of packet.
type Type byte

// Protocol packet types.
const (
	// Connect this event is sent:
	// * by the client when requesting access to a namespace
	// * by the server when accepting the connection to a namespace
	Connect Type = iota
	// Disconnect this event is used when one side wants to disconnect from a namespace.
	// It does not contain any payload nor acknowledgement id.
	Disconnect
	// Event this event is used when one side wants to transmit some data (without binary) to the other side.
	// It does contain a payload, and an optional acknowledgement id.
	Event
	// Ack this event is used when one side has received an EVENT or a BINARY_EVENT with an acknowledgement id.
	// It contains the acknowledgement id received in the previous packet, and may contain a payload (without binary).
	Ack
	// Error this event is sent by the server when the connection to a namespace is refused.
	// It contains a payload with a "message" and an optional "data" fields.
	Error
	// BinaryEvent Note: Both BINARY_EVENT and BINARY_ACK are used by the built-in parser,
	// in order to make a distinction between packets that contain binary content and those which don't.
	// They may not be used by other custom parsers.
	// This event is used when one side wants to transmit some data (including binary) to the other side.
	// It does contain a payload, and an optional acknowledgement id.
	BinaryEvent
	// BinaryAck This event is used when one side has received an EVENT or a BINARY_EVENT with an acknowledgement id.
	// It contains the acknowledgement id received in the previous packet, and contain a payload including binary.
	BinaryAck
)

func (t Type) IsValid() bool {
	return t > BinaryAck
}

// Header of packet.
type Header struct {
	Type      Type   `json:"type"`
	ID        uint64 `json:"id,omitempty"`
	Namespace string `json:"nsp,omitempty"`
	NeedAck   bool   `json:"-"`
}

// FrameType is the type of frames.
type FrameType byte

// FrameType aliases.
const (
	// Text is text type message.
	Text FrameType = iota
	// Binary is binary type message.
	Binary
)
