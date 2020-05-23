package parser

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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
