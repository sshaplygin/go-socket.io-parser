package parser

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// type noBufferStruct struct {
// 	Str   string            `json:"str"`
// 	I     int               `json:"i"`
// 	Array []int             `json:"array"`
// 	Map   map[string]string `json:"map"`
// }

// type bufferStruct struct {
// 	I      int     `json:"i"`
// 	Buffer *Buffer `json:"buf"`
// }

// type bufferInnerStruct struct {
// 	I      int                `json:"i"`
// 	Buffer *Buffer            `json:"buf"`
// 	Inner  *bufferInnerStruct `json:"inner,omitempty"`
// }

// type bufferTestCase struct {
// 	Header Header
// 	Event  string
// 	Var    []interface{}
// 	Data   [][]byte
// }

// var bufferTests = []bufferTestCase{
// 	{
// 		Header{Connect, false, 0, ""},
// 		"", nil,
// 		[][]byte{
// 			[]byte("0"),
// 		},
// 	},
// 	{
// 		Header{Error, "", 0, false},
// 		"",
// 		[]interface{}{"error"},
// 		[][]byte{
// 			[]byte("4[\"error\"]\n"),
// 		},
// 	},
// 	{
// 		Header{Event, "", 0, false}, "msg", []interface{}{
// 			&Buffer{Data: []byte{1, 2, 3}},
// 		}, [][]byte{
// 			[]byte("51-[\"msg\",{\"_placeholder\":true,\"num\":0}]\n"),
// 			[]byte{1, 2, 3},
// 		}},
// 	{
// 		Header{Connect, "", 0, true}, "", nil, [][]byte{
// 			[]byte("00"),
// 		}},
// 	{
// 		Header{Ack, "", 13, true}, "", []interface{}{"error"}, [][]byte{
// 			[]byte("313[\"error\"]\n"),
// 		}},
// 	{
// 		Header{Ack, "", 13, true}, "", []interface{}{
// 			&Buffer{
// 				Data: []byte{1, 2, 3},
// 			}}, [][]byte{
// 			[]byte("61-13[{\"_placeholder\":true,\"num\":0}]\n"),
// 			[]byte{1, 2, 3},
// 		}},
// 	{
// 		Header{Disconnect, "/woot", 0, false}, "", nil, [][]byte{
// 			[]byte("1/woot"),
// 		}},
// 	{
// 		Header{Event, "/woot", 0, false}, "msg", []interface{}{
// 			1,
// 		}, [][]byte{
// 			[]byte("2/woot,[\"msg\",1]\n"),
// 		}},
// 	{
// 		Header{Event, "/woot", 0, false}, "msg", []interface{}{
// 			&Buffer{Data: []byte{2, 3, 4}},
// 		}, [][]byte{
// 			[]byte("51-/woot,[\"msg\",{\"_placeholder\":true,\"num\":0}]\n"),
// 			[]byte{2, 3, 4},
// 		}},
// 	{
// 		Header{Disconnect, "/woot", 1, true}, "", nil, [][]byte{
// 			[]byte("1/woot,1"),
// 		}},
// 	{
// 		Header{Event, "/woot", 1, true}, "msg", []interface{}{
// 			1,
// 		}, [][]byte{
// 			[]byte("2/woot,1[\"msg\",1]\n"),
// 		}},
// 	{
// 		Header{Event, "/woot", 1, true},
// 		"msg",
// 		[]interface{}{
// 			&Buffer{Data: []byte{2, 3, 4}},
// 		}, [][]byte{
// 			[]byte("51-/woot,1[\"msg\",{\"_placeholder\":true,\"num\":0}]\n"),
// 			[]byte{2, 3, 4},
// 		}},
// }

var attachmentTests = []struct {
	buffer         Buffer
	textEncoding   string
	binaryEncoding string
}{
	{
		Buffer{
			IsBinary: false,
			Num:      0,
			Data:     []byte{1, 255},
		},
		`{"type":"Buffer","data":[1,255]}`,
		`{"_placeholder":true,"num":0}`,
	},
	{
		Buffer{
			IsBinary: false,
			Num:      1,
			Data:     []byte{},
		},
		`{"type":"Buffer","data":[]}`,
		`{"_placeholder":true,"num":1}`,
	},
	{
		Buffer{
			IsBinary: false,
			Num:      2,
			Data:     nil,
		},
		`{"type":"Buffer","data":[]}`,
		`{"_placeholder":true,"num":2}`,
	},
}

func TestAttachmentEncodeText(t *testing.T) {
	should := assert.New(t)
	must := require.New(t)

	for _, test := range attachmentTests {
		b := test.buffer
		b.IsBinary = false
		j, err := json.Marshal(b)
		must.Nil(err)
		should.Equal(test.textEncoding, string(j))
	}
}

func TestAttachmentEncodeBinary(t *testing.T) {
	should := assert.New(t)
	must := require.New(t)

	for _, test := range attachmentTests {
		a := test.buffer
		a.IsBinary = true
		j, err := json.Marshal(a)
		must.Nil(err)
		should.Equal(test.binaryEncoding, string(j))
	}
}

func TestAttachmentDecodeText(t *testing.T) {
	should := assert.New(t)
	must := require.New(t)

	for _, test := range attachmentTests {
		var a Buffer
		err := json.Unmarshal([]byte(test.textEncoding), &a)
		must.Nil(err)
		should.False(a.IsBinary)
		if len(test.buffer.Data) == 0 {
			should.Equal([]byte{}, a.Data)
			continue
		}
		should.Equal(test.buffer.Data, a.Data)
	}
}

func TestAttachmentDecodeBinary(t *testing.T) {
	should := assert.New(t)
	must := require.New(t)

	for _, test := range attachmentTests {
		var a Buffer
		err := json.Unmarshal([]byte(test.binaryEncoding), &a)
		must.Nil(err)
		should.True(a.IsBinary)
		should.Equal(test.buffer.Num, a.Num)
	}
}
