package go_socketio_parser

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type attachDataTestCase struct {
	name     string
	max      uint64
	binaries [][]byte
	data     interface{}
}

type bufferStruct struct {
	I      int     `json:"i"`
	Buffer *Buffer `json:"buf"`
}

var attachDataTests = []attachDataTestCase{
	{
		name: "&Buffer",
		max:  1,
		binaries: [][]byte{
			[]byte{1, 2},
		},
		data: &Buffer{
			Data: []byte{1, 2},
		},
	},
	{
		name: "[]interface{}{Buffer}",
		max:  1,
		binaries: [][]byte{
			[]byte{1, 2},
		},
		data: []interface{}{
			&Buffer{
				Data: []byte{1, 2},
			},
		},
	},
	{
		name: "[]interface{}{Buffer,Buffer}",
		max:  2,
		binaries: [][]byte{
			[]byte{1, 2},
			[]byte{3, 4},
		},
		data: []interface{}{
			&Buffer{Data: []byte{1, 2}},
			&Buffer{Data: []byte{3, 4}},
		},
	},
	{
		name: "[1]interface{}{Buffer}",
		max:  1,
		binaries: [][]byte{
			[]byte{1, 2},
		},
		data: []interface{}{
			&Buffer{
				Data: []byte{1, 2},
			},
		},
	},
	{
		name: "[2]interface{}{Buffer,Buffer}",
		max:  2,
		binaries: [][]byte{
			[]byte{1, 2},
			[]byte{3, 4},
		},
		data: [...]interface{}{
			&Buffer{
				Data: []byte{1, 2},
			},
			&Buffer{
				Data: []byte{3, 4},
			},
		},
	},
	{
		name: "Struct{Buffer}",
		max:  1,
		binaries: [][]byte{
			[]byte{1, 2},
		},
		data: bufferStruct{
			I: 3,
			Buffer: &Buffer{
				Data: []byte{1, 2},
			},
		},
	},
	{
		name: "map{Buffer}",
		max:  1,
		binaries: [][]byte{
			[]byte{1, 2},
		},
		data: map[string]interface{}{
			"data": &Buffer{
				Data: []byte{1, 2},
			},
			"i": 3,
		},
	},
}

func TestAttachBuffer(t *testing.T) {
	for _, test := range attachDataTests {
		t.Run(test.name, func(t *testing.T) {
			index := uint64(0)

			buf, err := attachBuffer(reflect.ValueOf(test.data), &index)
			require.NoError(t, err)

			assert.Equal(t, test.max, index)
			assert.Equal(t, test.binaries, buf)
		})
	}
}

func BenchmarkAttachBuffer(b *testing.B) {
	for i := 0; i < b.N; i++ {
		index := uint64(0)

		_, _ = attachBuffer(reflect.ValueOf(
			map[string]interface{}{
				"data": &Buffer{
					Data: []byte{1, 2},
				},
				"i": 3,
			},
		), &index)
	}
}
