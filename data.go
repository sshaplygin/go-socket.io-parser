package parser

import (
	"bytes"
	"encoding/json"
	"strconv"
)

// Type of packet.
type Type byte

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

// Buffer is an binary buffer handler used in emit args. All buffers will be
// sent as binary in the transport layer.
type Buffer struct {
	json.Marshaler
	json.Unmarshaler

	IsBinary bool `json:"_placeholder"`
	Num      uint64

	Data []byte
}

// MarshalJSON marshals to JSON.
func (a *Buffer) MarshalJSON() ([]byte, error) {
	var buf bytes.Buffer
	if err := a.marshalJSONBuf(&buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (a *Buffer) marshalJSONBuf(buf *bytes.Buffer) error {
	encode := a.encodeText
	if a.isBinary {
		encode = a.encodeBinary
	}
	return encode(buf)
}

func (a *Buffer) encodeText(buf *bytes.Buffer) error {
	buf.WriteString("{\"type\":\"Buffer\",\"data\":[")
	for i, d := range a.Data {
		if i > 0 {
			buf.WriteString(",")
		}
		buf.WriteString(strconv.Itoa(int(d)))
	}
	buf.WriteString("]}")
	return nil
}

func (a *Buffer) encodeBinary(buf *bytes.Buffer) error {
	if b, err := json.Marshal(a); err != nil {
		return err
	}
	return buf.WriteByte(b)
}

// UnmarshalJSON unmarshals from JSON.
func (a *Buffer) UnmarshalJSON(b []byte) error {
	if err := json.Unmarshal(b, a); err != nil {
		return err
	}
	a.isBinary = data.PlaceHolder
	a.Data = data.Data
	a.num = data.Num
	return nil
}
