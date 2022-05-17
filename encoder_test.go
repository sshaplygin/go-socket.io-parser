package go_socketio_parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//type fakeWriter struct {
//	current *bytes.Buffer
//
//	ftype   FrameType
//	types   []FrameType
//	buffers []*bytes.Buffer
//}
//
//func (w *fakeWriter) NextWriter(ft FrameType) (io.WriteCloser, error) {
//	w.current = bytes.NewBuffer(nil)
//	w.ftype = ft
//	return w, nil
//}
//
//func (w *fakeWriter) Write(p []byte) (int, error) {
//	return w.current.Write(p)
//}
//
//func (w *fakeWriter) Close() error {
//	w.types = append(w.types, w.ftype)
//	w.buffers = append(w.buffers, w.current)
//	return nil
//}

func TestEncoder(t *testing.T) {
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			resp, err := Marshal(test.Header, test.Data)
			require.NoError(t, err)

			assert.Equal(t, test.Tmpl, string(resp))
		})
	}
}

//type noBufferStruct struct {
//	Str   string            `json:"str"`
//	I     int               `json:"i"`
//	Array []int             `json:"array"`
//	Map   map[string]string `json:"map"`
//}

type bufferInnerStruct struct {
	I      int                `json:"i"`
	Buffer *Buffer            `json:"buf"`
	Inner  *bufferInnerStruct `json:"inner,omitempty"`
}
