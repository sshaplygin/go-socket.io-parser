package go_socketio_parser

import (
	"bytes"
	"encoding/json"
	"reflect"
	"strconv"
)

// Buffer is a binary buffer handler used in emit args. All buffers will be
// sent as binary in the transport layer.
type Buffer struct {
	IsBinary bool   `json:"_placeholder"`
	Num      uint64 `json:"num"`

	Data []byte `json:"-"`
}

// Marshal marshals to JSON.
func (b *Buffer) Marshal() ([]byte, error) {
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

func (b *Buffer) Unmarshal(data []byte) error {
	if err := json.Unmarshal(data, b); err != nil {
		return err
	}
	return nil
}

const structBuffer = "Buffer"

func attachBuffer(v reflect.Value, index *uint64) ([][]byte, error) {
	for v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface {
		v = v.Elem()
	}
	var ret [][]byte

	switch v.Kind() {
	case reflect.Struct:
		if v.Type().Name() == structBuffer {
			if !v.CanAddr() {
				return nil, ErrBufferAddress
			}

			buffer, ok := v.Addr().Interface().(*Buffer)
			if !ok {
				return nil, nil
			}

			buffer.Num = *index
			buffer.IsBinary = true
			ret = append(ret, buffer.Data)
			*index++
		} else {
			for i := 0; i < v.NumField(); i++ {
				b, err := attachBuffer(v.Field(i), index)
				if err != nil {
					return nil, err
				}

				ret = append(ret, b...)
			}
		}
	case reflect.Array, reflect.Slice:
		for i := 0; i < v.Len(); i++ {
			b, err := attachBuffer(v.Index(i), index)
			if err != nil {
				return nil, err
			}

			ret = append(ret, b...)
		}
	case reflect.Map:
		for _, key := range v.MapKeys() {
			b, err := attachBuffer(v.MapIndex(key), index)
			if err != nil {
				return nil, err
			}

			ret = append(ret, b...)
		}
	}

	return ret, nil
}

func detachBuffer(v reflect.Value, buffers []Buffer) error {
	for v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface {
		v = v.Elem()
	}
	switch v.Kind() {
	case reflect.Struct:
		if v.Type().Name() == structBuffer {
			if !v.CanAddr() {
				return ErrBufferAddress
			}
			buffer := v.Addr().Interface().(*Buffer)
			if buffer.IsBinary {
				*buffer = buffers[buffer.Num]
			}
			return nil
		}
		for i := 0; i < v.NumField(); i++ {
			if err := detachBuffer(v.Field(i), buffers); err != nil {
				return err
			}
		}
	case reflect.Map:
		for _, key := range v.MapKeys() {
			if err := detachBuffer(v.MapIndex(key), buffers); err != nil {
				return err
			}
		}
	case reflect.Array:
		fallthrough
	case reflect.Slice:
		for i := 0; i < v.Len(); i++ {
			if err := detachBuffer(v.Index(i), buffers); err != nil {
				return err
			}
		}
	}
	return nil
}
