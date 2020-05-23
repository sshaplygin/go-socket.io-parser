package parser

import (
	"bytes"
	"io"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type fakeReader struct {
	index int
	datas [][]byte
	buf   *bytes.Buffer
}

func (r *fakeReader) NextReader() (FrameType, io.ReadCloser, error) {
	if r.index >= len(r.datas) {
		return 0, nil, io.EOF
	}
	r.buf = bytes.NewBuffer(r.datas[r.index])
	ft := BINARY
	if r.index == 0 {
		ft = TEXT
	}
	return ft, r, nil
}

func (r *fakeReader) Read(p []byte) (int, error) {
	return r.buf.Read(p)
}

func (r *fakeReader) Close() error {
	r.index++
	return nil
}

type noBufferStruct struct {
	Str   string            `json:"str"`
	I     int               `json:"i"`
	Array []int             `json:"array"`
	Map   map[string]string `json:"map"`
}

type bufferStruct struct {
	I      int     `json:"i"`
	Buffer *Buffer `json:"buf"`
}

type bufferInnerStruct struct {
	I      int                `json:"i"`
	Buffer *Buffer            `json:"buf"`
	Inner  *bufferInnerStruct `json:"inner,omitempty"`
}

type testCase struct {
	Name   string
	Header Header
	Event  string
	Var    []interface{}
	Data   [][]byte
}

var tests = []testCase{
	{
		"Empty", Header{Connect, "", 0, false}, "", nil, [][]byte{
			[]byte("0"),
		}},
	{"Data", Header{Error, "", 0, false}, "", []interface{}{"error"}, [][]byte{
		[]byte("4[\"error\"]\n"),
		}
	},
	{"BData", Header{Event, "", 0, false}, "msg", []interface{}{
		&Buffer{Data: []byte{1, 2, 3}},
	}, [][]byte{
		[]byte("51-[\"msg\",{\"_placeholder\":true,\"num\":0}]\n"),
		[]byte{1, 2, 3},
	}},
	{"ID", Header{Connect, "", 0, true}, "", nil, [][]byte{
		[]byte("00"),
	}},
	{"IDData", Header{Ack, "", 13, true}, "", []interface{}{"error"}, [][]byte{
		[]byte("313[\"error\"]\n"),
	}},
	{"IDBData", Header{Ack, "", 13, true}, "", []interface{}{
		&Buffer{
			Data: []byte{1, 2, 3},
		}}, [][]byte{
		[]byte("61-13[{\"_placeholder\":true,\"num\":0}]\n"),
		[]byte{1, 2, 3},
	}},
	{"Namespace", Header{Disconnect, "/woot", 0, false}, "", nil, [][]byte{
		[]byte("1/woot"),
	}},
	{"NamespaceData", Header{Event, "/woot", 0, false}, "msg", []interface{}{
		1,
	}, [][]byte{
		[]byte("2/woot,[\"msg\",1]\n"),
	}},
	{"NamespaceBData", Header{Event, "/woot", 0, false}, "msg", []interface{}{
		&Buffer{Data: []byte{2, 3, 4}},
	}, [][]byte{
		[]byte("51-/woot,[\"msg\",{\"_placeholder\":true,\"num\":0}]\n"),
		[]byte{2, 3, 4},
	}},
	{"NamespaceID", Header{Disconnect, "/woot", 1, true}, "", nil, [][]byte{
		[]byte("1/woot,1"),
	}},
	{"NamespaceIDData", Header{Event, "/woot", 1, true}, "msg", []interface{}{
		1,
	}, [][]byte{
		[]byte("2/woot,1[\"msg\",1]\n"),
	}},
	{"NamespaceIDBData", Header{Event, "/woot", 1, true}, "msg", []interface{}{
		&Buffer{Data: []byte{2, 3, 4}},
	}, [][]byte{
		[]byte("51-/woot,1[\"msg\",{\"_placeholder\":true,\"num\":0}]\n"),
		[]byte{2, 3, 4},
	}},
}

func TestDecoder(t *testing.T) {
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			should := assert.New(t)
			must := require.New(t)

			r := fakeReader{datas: test.Datas}
			decoder := NewDecoder(&r)

			defer func() {
				decoder.DiscardLast()
				decoder.Close()
			}()

			var header Header
			var event string
			err := decoder.DecodeHeader(&header, &event)
			must.Nil(err, "decode header error: %s", err)
			should.Equal(test.Header, header)
			should.Equal(test.Event, event)
			types := make([]reflect.Type, len(test.Var))
			for i := range types {
				types[i] = reflect.TypeOf(test.Var[i])
			}
			ret, err := decoder.DecodeArgs(types)
			must.Nil(err, "decode args error: %s", err)
			vars := make([]interface{}, len(ret))
			for i := range vars {
				vars[i] = ret[i].Interface()
			}
			if len(vars) == 0 {
				vars = nil
			}
			should.Equal(test.Var, vars)
		})
	}
}
