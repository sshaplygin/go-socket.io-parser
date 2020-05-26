package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type attachmentTest struct {
	buffer         Buffer
	textEncoding   string
	binaryEncoding string
}

var attachmentTests = []attachmentTest{
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
			IsBinary: true,
			Num:      0,
			Data:     []byte{1, 255},
		},
		`{"type":"Buffer","data":[1,255]}`,
		`{"_placeholder":true,"num":0}`,
	}, {
		Buffer{
			IsBinary: true,
			Num:      2,
			Data:     []byte{3, 44},
		},
		`{"type":"Buffer","data":[3,44]}`,
		`{"_placeholder":true,"num":2}`,
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
		j, err := b.Marshal()
		must.Nil(err)
		t.Log(test.textEncoding, string(j))
		should.Equal(test.textEncoding, string(j))
	}
}

func TestAttachmentEncodeBinary(t *testing.T) {
	should := assert.New(t)
	must := require.New(t)

	for _, test := range attachmentTests {
		b := test.buffer
		b.IsBinary = false
		j, err := b.Marshal()
		must.Nil(err)
		t.Log(test.textEncoding, string(j))
		should.Equal(test.textEncoding, string(j))
	}
}

func TestAttachmentDecodeText(t *testing.T) {
	should := assert.New(t)
	must := require.New(t)

	for _, test := range attachmentTests {
		var a Buffer
		err := a.Unmarshal([]byte(test.textEncoding))
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
		err := a.Unmarshal([]byte(test.binaryEncoding))
		must.Nil(err)
		should.True(a.IsBinary)
		should.Equal(test.buffer.Num, a.Num)
	}
}
