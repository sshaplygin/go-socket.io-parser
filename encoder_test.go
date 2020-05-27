package parser

import (
	"bytes"
	"io"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type fakeWriter struct {
	current *bytes.Buffer

	ftype   FrameType
	types   []FrameType
	buffers []*bytes.Buffer
}

func (w *fakeWriter) NextWriter(ft FrameType) (io.WriteCloser, error) {
	w.current = bytes.NewBuffer(nil)
	w.ftype = ft
	return w, nil
}

func (w *fakeWriter) Write(p []byte) (int, error) {
	return w.current.Write(p)
}

func (w *fakeWriter) Close() error {
	w.types = append(w.types, w.ftype)
	w.buffers = append(w.buffers, w.current)
	return nil
}

func TestEncoder(t *testing.T) {
	t.Parallel()

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			should := assert.New(t)
			must := require.New(t)

			w := fakeWriter{}
			encoder := NewEncoder(&w)
			v := test.Var
			if test.Header.Type == Event {
				v = append([]interface{}{test.Event}, test.Var...)
			}
			err := encoder.Encode(test.Header, v)
			must.Nil(err)
			must.Equal(len(test.Data), len(w.types))
			must.Equal(len(test.Data), len(w.buffers))
			for i := range w.types {
				if i == 0 {
					should.Equal(TEXT, w.types[i])
					should.Equal(string(test.Data[i]), string(w.buffers[i].Bytes()))
					continue
				}
				should.Equal(BINARY, w.types[i])
				should.Equal(test.Data[i], w.buffers[i].Bytes())
			}
		})
	}
}

type attachTestCase struct {
	name     string
	max      uint64
	binaries [][]byte
	data     interface{}
}

type bufferStruct struct {
	I      int     `json:"i"`
	Buffer *Buffer `json:"buf"`
}

type noBufferStruct struct {
	Str   string            `json:"str"`
	I     int               `json:"i"`
	Array []int             `json:"array"`
	Map   map[string]string `json:"map"`
}

type bufferInnerStruct struct {
	I      int                `json:"i"`
	Buffer *Buffer            `json:"buf"`
	Inner  *bufferInnerStruct `json:"inner,omitempty"`
}

var attachTests = []attachTestCase{
	{
		"&Buffer",
		1,
		[][]byte{[]byte{1, 2}},
		&Buffer{Data: []byte{1, 2}},
	},
	{"[]interface{}{Buffer}",
		1,
		[][]byte{[]byte{1, 2}},
		[]interface{}{&Buffer{Data: []byte{1, 2}}},
	},
	{"[]interface{}{Buffer,Buffer}",
		2,
		[][]byte{[]byte{1, 2}, []byte{3, 4}},
		[]interface{}{
			&Buffer{Data: []byte{1, 2}},
			&Buffer{Data: []byte{3, 4}},
		}},
	{
		"[1]interface{}{Buffer}",
		1,
		[][]byte{[]byte{1, 2}},
		[...]interface{}{&Buffer{Data: []byte{1, 2}}},
	},
	{
		"[2]interface{}{Buffer,Buffer}",
		2,
		[][]byte{[]byte{1, 2}, []byte{3, 4}},
		[...]interface{}{
			&Buffer{Data: []byte{1, 2}},
			&Buffer{Data: []byte{3, 4}},
		},
	},
	{
		"Struct{Buffer}",
		1,
		[][]byte{[]byte{1, 2}},
		bufferStruct{
			3,
			&Buffer{Data: []byte{1, 2}},
		},
	},
	{
		"map{Buffer}",
		1,
		[][]byte{[]byte{1, 2}},
		map[string]interface{}{
			"data": &Buffer{Data: []byte{1, 2}},
			"i":    3,
		},
	},
}

func TestAttachBuffer(t *testing.T) {
	t.Parallel()

	e := Encoder{}
	for _, test := range attachTests {
		t.Run(test.name, func(t *testing.T) {
			should := assert.New(t)
			must := require.New(t)
			index := uint64(0)
			b, err := e.attachBuffer(reflect.ValueOf(test.data), &index)
			must.Nil(err)
			should.Equal(test.max, index)
			should.Equal(test.binaries, b)
		})
	}
}
