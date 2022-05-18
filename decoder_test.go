package go_socketio_parser

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDecoder(t *testing.T) {
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			var message Packet
			err := Unmarshal([]byte(test.Tmpl), &message)
			require.NoError(t, err)
			require.Equal(t, len(test.Data), len(message.Data))

			assert.Equal(t, test.Header, message.Header)
			for idx, data := range test.Data {
				assert.Equal(t, data, message.Data[idx])
			}
		})
	}
}

func Test_readJSONPayload(t *testing.T) {
	t.Run("empty json", func(t *testing.T) {
		data := []byte(`{}`)
		r := bytes.NewReader(data)

		buf, err := readJSONPayload(r)
		require.NoError(t, err)

		assert.Empty(t, buf)
	})
}

func Test_decodeData(t *testing.T) {

	t.Run("invalid ", func(t *testing.T) {
		data := []byte(`["hello",{"_placeholder":true,"num":0},"word",{"_placeholder":true,"num":1}]`)

		r := bytes.NewReader(data)

		decodedData, err := decodeData(r)
		require.Error(t, err, "not found binary attachments")

		assert.Empty(t, decodedData)
	})

	t.Run("data with double binary payloads", func(t *testing.T) {
		data := []byte(`["hello",{"_placeholder":true,"num":0},"word",{"_placeholder":true,"num":1}]`)

		data = append(data, '\n')
		firstData := []byte{1, 2, 3}
		data = append(data, firstData...)

		data = append(data, '\n')
		secondData := []byte{4, 5, 6}
		data = append(data, secondData...)

		r := bytes.NewReader(data)

		decodedData, err := decodeData(r)
		require.NoError(t, err)
		require.Len(t, decodedData, 4)

		assert.Equal(t, "hello", decodedData[0])
		firstBuffer, ok := decodedData[1].(*Buffer)
		require.True(t, ok)
		assert.Equal(t, firstData, firstBuffer.Data)
		assert.Equal(t, "word", decodedData[2])
		secondBuffer, ok := decodedData[3].(*Buffer)
		require.True(t, ok)
		assert.Equal(t, secondData, secondBuffer.Data)
	})
}

func BenchmarkUnmarshal(b *testing.B) {
	data := []byte(`51-/woot,1["msg",{"_placeholder":true,"num":0}]` + string('\n') + string([]byte{2, 3, 4}))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var message Packet
		_ = Unmarshal(data, &message)
	}
}
