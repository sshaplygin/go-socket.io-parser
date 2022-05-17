package go_socketio_parser

import (
	"github.com/stretchr/testify/require"
	"testing"
)

//type fakeReader struct {
//	index int
//	data  [][]byte
//	buf   *bytes.Buffer
//}
//
//func (r *fakeReader) NextReader() (FrameType, io.ReadCloser, error) {
//	if r.index >= len(r.data) {
//		return 0, nil, io.EOF
//	}
//	r.buf = bytes.NewBuffer(r.data[r.index])
//	ft := Binary
//	if r.index == 0 {
//		ft = Text
//	}
//	return ft, r, nil
//}
//
//func (r *fakeReader) Read(p []byte) (int, error) {
//	return r.buf.Read(p)
//}
//
//func (r *fakeReader) Close() error {
//	r.index++
//	return nil
//}

type testCase struct {
	Name       string
	Header     Header
	Data       []interface{} // body
	Tmpl       string        // encoded template
	AttachData []byte        // binary attachments extracted.
}

var tests = []testCase{
	{
		Name: "empty Connect type",
		Header: Header{
			Type: Connect,
		},
		Tmpl: "0",
	},
	{
		Name: "empty Disconnect type",
		Header: Header{
			Type: Disconnect,
		},
		Tmpl: "1",
	},
	{
		Name: "empty Event type",
		Header: Header{
			Type: Event,
		},
		Tmpl: "2",
	},
	{
		Name: "empty Ack type",
		Header: Header{
			Type: Ack,
		},
		Tmpl: "3",
	},
	{
		Name: "empty Error type",
		Header: Header{
			Type: Error,
		},
		Tmpl: "4",
	},
	{
		Name: "empty BinaryEvent type",
		Header: Header{
			Type: BinaryEvent,
		},
		Tmpl: "50-",
	},
	{
		Name: "empty BinaryAck type",
		Header: Header{
			Type: BinaryAck,
		},
		Tmpl: "60-",
	},
	{
		Name: "Event type with nsp, id, data",
		Header: Header{
			Type:      Event,
			Namespace: "/admin",
			ID:        456,
		},
		Data: []interface{}{
			"project:delete", 123,
		},
		Tmpl: `2/admin,456["project:delete",123]`,
	},
	{
		Name: "Data",
		Header: Header{
			Type: Error,
		},
		Data: []interface{}{
			"error",
		},
		Tmpl: `4["error"]`,
	},
	{
		Name: "BData",
		Header: Header{
			Type: Event,
		},
		Data: []interface{}{
			"msg", &Buffer{Data: []byte{1, 2, 3}},
		},
		Tmpl:       `51-["msg",{"_placeholder":true,"num":0}]` + string('\n') + string([]byte{1, 2, 3}),
		AttachData: []byte{1, 2, 3},
	},
	{
		Name: "ID",
		Header: Header{
			Type:    Connect,
			NeedAck: true,
		},
		Tmpl: "00",
	},
	{
		Name: "IDData",
		Header: Header{
			Type:    Ack,
			NeedAck: true,
			ID:      13,
		},
		Data: []interface{}{
			"error",
		},
		Tmpl: `313["error"]`,
	},
	{
		Name: "IDBData",
		Header: Header{
			Type:    Ack,
			NeedAck: true,
			ID:      13,
		},
		Data: []interface{}{
			&Buffer{
				Data: []byte{1, 2, 3},
			},
		},
		Tmpl: `61-13[{"_placeholder":true,"num":0}]` + string('\n') + string([]byte{1, 2, 3}),
	},
	{
		Name: "Namespace",
		Header: Header{
			Type:      Disconnect,
			Namespace: "/woot",
		},
		Tmpl: "1/woot",
	},
	{
		Name: "NamespaceData",
		Header: Header{
			Type:      Event,
			Namespace: "/woot",
		},
		Data: []interface{}{
			"msg", 1,
		},
		Tmpl: `2/woot,["msg",1]`,
	},
	{
		Name: "NamespaceBData",
		Header: Header{
			Type:      Event,
			Namespace: "/woot",
		},
		Data: []interface{}{
			"msg", &Buffer{Data: []byte{2, 3, 4}},
		},
		Tmpl:       `51-/woot,["msg",{"_placeholder":true,"num":0}]` + string('\n') + string([]byte{2, 3, 4}),
		AttachData: []byte{2, 3, 4},
	},
	{
		Name: "NamespaceID",
		Header: Header{
			Type:      Disconnect,
			NeedAck:   true,
			ID:        1,
			Namespace: "/woot",
		},
		Tmpl: "1/woot,1",
	},
	{
		Name: "NamespaceIDData",
		Header: Header{
			Type:      Event,
			NeedAck:   true,
			ID:        1,
			Namespace: "/woot",
		},
		Data: []interface{}{
			"msg", 1,
		},
		Tmpl: `2/woot,1["msg",1]`,
	},
	{
		Name: "NamespaceIDBData",
		Header: Header{
			Type:      Event,
			NeedAck:   true,
			ID:        1,
			Namespace: "/woot",
		},
		Data: []interface{}{
			"msg", &Buffer{Data: []byte{2, 3, 4}},
		},
		Tmpl:       `51-/woot,1["msg",{"_placeholder":true,"num":0}]` + string('\n') + string([]byte{2, 3, 4}),
		AttachData: []byte{2, 3, 4},
	},
}

func TestDecoder(t *testing.T) {
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			var h Header
			var attach []interface{}
			err := Unmarshal([]byte(test.Tmpl), &h, attach)
			t.Log(test.Tmpl, h, attach)
			require.NoError(t, err)
			//should := assert.New(t)
			//must := require.New(t)
			//
			////r := fakeReader{data: test.ExpData}
			//decoder := NewDecoder(&r)
			//
			//defer func() {
			//	_ = decoder.DiscardLast()
			//	_ = decoder.Close()
			//}()
			//
			//var header Header
			//var event string
			//err := decoder.DecodeHeader(&header, &event)
			//must.Nil(err, "decode header error: %s", err)
			//
			//should.Equal(test.Header, header)
			////should.Equal(test.Event, event)
			//
			//types := make([]reflect.Type, len(test.Data))
			//for i := range types {
			//	types[i] = reflect.TypeOf(test.Data[i])
			//}
			//ret, err := decoder.DecodeArgs(types)
			//must.Nil(err, "decode args error: %s", err)
			//vars := make([]interface{}, len(ret))
			//for i := range vars {
			//	vars[i] = ret[i].Interface()
			//}
			//if len(vars) == 0 {
			//	vars = nil
			//}
			//should.Equal(test.Data, vars)
		})
	}
}
