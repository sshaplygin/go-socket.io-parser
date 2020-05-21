package parser

import (
	"bytes"
	"encoding/json"
	"strconv"
)

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
func (b *Buffer) MarshalJSON() ([]byte, error) {
	var buf bytes.Buffer
	if err := b.marshalJSONBuf(&buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (b *Buffer) marshalJSONBuf(buf *bytes.Buffer) error {
	encode := b.encodeText
	if b.IsBinary {
		encode = b.encodeBinary
	}
	return encode(buf)
}

func (b *Buffer) encodeText(buf *bytes.Buffer) error {
	buf.WriteString("{\"type\":\"Buffer\",\"data\":[")
	for i, d := range b.Data {
		if i > 0 {
			buf.WriteString(",")
		}
		buf.WriteString(strconv.Itoa(int(d)))
	}
	buf.WriteString("]}")
	return nil
}

func (b *Buffer) encodeBinary(buf *bytes.Buffer) error {
	data, err := json.Marshal(b)
	if err != nil {
		return err
	}
	_, err = buf.Write(data)
	return err
}

// UnmarshalJSON unmarshals from JSON.
func (b *Buffer) UnmarshalJSON(data []byte) error {
	if err := json.Unmarshal(data, b); err != nil {
		return err
	}
	return nil
}
