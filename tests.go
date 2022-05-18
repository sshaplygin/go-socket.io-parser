package go_socketio_parser

type testCase struct {
	Name       string
	Header     Header
	Data       []interface{} // body
	Tmpl       string        // encoded template
	AttachData []byte        // binary attachments extracted.
}

var tests = []testCase{
	//{
	//	Name: "connect",
	//	Header: Header{
	//		Type: Connect,
	//	},
	//	Tmpl: "0",
	//},
	//{
	//	Name: "disconnect",
	//	Header: Header{
	//		Type: Disconnect,
	//	},
	//	Tmpl: "1",
	//},
	//{
	//	Name: "event",
	//	Header: Header{
	//		Type: Event,
	//	},
	//	Tmpl: "2",
	//},
	//{
	//	Name: "ack",
	//	Header: Header{
	//		Type: Ack,
	//	},
	//	Tmpl: "3",
	//},
	//{
	//	Name: "error",
	//	Header: Header{
	//		Type: Error,
	//	},
	//	Tmpl: "4",
	//},
	//{
	//	Name: "binaryEvent",
	//	Header: Header{
	//		Type: BinaryEvent,
	//	},
	//	Tmpl: "50-",
	//},
	//{
	//	Name: "binaryAck",
	//	Header: Header{
	//		Type: BinaryAck,
	//	},
	//	Tmpl: "60-",
	//},
	{
		Name: "event nsp id data",
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
		Name: "error data",
		Header: Header{
			Type: Error,
		},
		Data: []interface{}{
			"error",
		},
		Tmpl: `4["error"]`,
	},
	{
		Name: "event data buffer",
		Header: Header{
			Type: Event,
		},
		Data: []interface{}{
			"msg", &Buffer{
				IsBinary: true,
				Data:     []byte{1, 2, 3},
			},
		},
		Tmpl:       `51-["msg",{"_placeholder":true,"num":0}]` + string('\n') + string([]byte{1, 2, 3}),
		AttachData: []byte{1, 2, 3},
	},
	{
		Name: "connect id",
		Header: Header{
			Type: Connect,
			ID:   145,
		},
		Tmpl: "0145",
	},
	{
		Name: "ack id data",
		Header: Header{
			Type: Ack,

			ID: 13,
		},
		Data: []interface{}{
			"error",
		},
		Tmpl: `313["error"]`,
	},
	{
		Name: "ack id buffer",
		Header: Header{
			Type: Ack,
			ID:   13,
		},
		Data: []interface{}{
			&Buffer{
				IsBinary: true,
				Data:     []byte{1, 2, 3},
			},
		},
		Tmpl: `61-13[{"_placeholder":true,"num":0}]` + string('\n') + string([]byte{1, 2, 3}),
	},
	{
		Name: "disconnect nsp",
		Header: Header{
			Type:      Disconnect,
			Namespace: "/woot",
		},
		Tmpl: "1/woot",
	},
	{
		Name: "event nsp data",
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
		Name: "event namespace data buffer",
		Header: Header{
			Type:      Event,
			Namespace: "/woot",
		},
		Data: []interface{}{
			"msg",
			&Buffer{
				IsBinary: true,
				Data:     []byte{2, 3, 4},
			},
		},
		Tmpl:       `51-/woot,["msg",{"_placeholder":true,"num":0}]` + string('\n') + string([]byte{2, 3, 4}),
		AttachData: []byte{2, 3, 4},
	},
	{
		Name: "disconnect nsp id",
		Header: Header{
			Type:      Disconnect,
			ID:        1,
			Namespace: "/woot",
		},
		Tmpl: "1/woot,1",
	},
	{
		Name: "event nsp id data",
		Header: Header{
			Type:      Event,
			ID:        1,
			Namespace: "/woot",
		},
		Data: []interface{}{
			"msg", 1,
		},
		Tmpl: `2/woot,1["msg",1]`,
	},
	{
		Name: "event namespace id data buffer",
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
		Tmpl:       `51-/woot,1["msg",{"_placeholder":true,"num":0}]` + string('\n') + string([]byte{2, 3, 4}),
		AttachData: []byte{2, 3, 4},
	},
}
