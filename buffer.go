package go_socketio_parser

import (
	"reflect"
)

// Buffer is a binary buffer handler used in emit args. All buffers will be
// sent as binary in the transport layer.
type Buffer struct {
	IsBinary bool   `json:"_placeholder"`
	Num      uint64 `json:"num"`

	Data []byte `json:"-"`
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
