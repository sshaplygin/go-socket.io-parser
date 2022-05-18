package go_socketio_parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEncoder(t *testing.T) {
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			packet := &Packet{
				Header: test.Header,
				Data:   test.Data,
			}

			resp, err := Marshal(packet)
			require.NoError(t, err)

			assert.Equal(t, test.Tmpl, string(resp))
		})
	}
}

func BenchmarkMarshal(b *testing.B) {
	message := &Packet{
		Header: Header{
			Type:      Event,
			ID:        1,
			Namespace: "/woot",
		},
		Data: []interface{}{
			"msg",
			&Buffer{
				IsBinary: true,
				Data:     []byte{2, 3, 4},
			},
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = Marshal(message)
	}
}
