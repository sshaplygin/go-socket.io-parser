package go_socketio_parser

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"reflect"
	"strings"
)

const binarySep = "-"
const nsSep = "/"
const packetTypeIdx = 0
const sep = ","
const dataSep = "["
const attachSep = byte('\n')

func Unmarshal(data []byte, h *Header, attach []interface{}) error {
	//r := bytes.NewReader(data)

	//for i, b := range data {
	//	// read packet type
	//	if i == 0 {
	//		ht := Type(b - '0')
	//		if !ht.IsValid() {
	//			return ErrInvalidPackageType
	//		}
	//		h.Type = ht
	//
	//		continue
	//	}
	//
	//	//
	//
	//	if b == attachSep && len(data) > i {
	//		attachIdx = i + 1
	//	}
	//}

	return nil
}

func (d *Decoder) DecodeHeader(header *Header, event *string) error {
	ft, r, err := d.fr.NextReader()
	if err != nil {
		return err
	}

	if ft != Text {
		return ErrShouldTextPackageType
	}

	d.lastFrame = r
	br, ok := r.(byteReader)
	if !ok {
		br = bufio.NewReader(r)
	}
	d.packetReader = br

	bufferCount, err := d.readHeader(header)
	if err != nil {
		return err
	}
	d.bufferCount = bufferCount
	if header.Type == BinaryEvent || header.Type == BinaryAck {
		header.Type -= 3
	}
	d.isEvent = header.Type == Event
	if d.isEvent {
		if err = d.readEvent(event); err != nil {
			return err
		}
	}
	return nil
}

func (d *Decoder) DecodeArgs(types []reflect.Type) ([]reflect.Value, error) {
	r := d.packetReader.(io.Reader)
	if d.isEvent {
		r = io.MultiReader(strings.NewReader("["), r)
	}

	ret := make([]reflect.Value, len(types))
	values := make([]interface{}, len(types))
	for i, typ := range types {
		if typ.Kind() == reflect.Ptr {
			typ = typ.Elem()
		}
		ret[i] = reflect.New(typ)
		values[i] = ret[i].Interface()
	}

	if err := json.NewDecoder(r).Decode(&values); err != nil {
		if err == io.EOF {
			err = nil
		}
		_ = d.DiscardLast()
		return nil, err
	}

	//we can't use defer or call DiscardLast before decoding, because
	//there are buffered readers involved and if we invoke .Close() json will encounter unexpected EOF.
	_ = d.DiscardLast()

	for i, typ := range types {
		if typ.Kind() != reflect.Ptr {
			ret[i] = ret[i].Elem()
		}
	}

	buffers := make([]Buffer, d.bufferCount)
	for i := range buffers {
		ft, r, err := d.fr.NextReader()
		if err != nil {
			return nil, err
		}
		buffers[i].Data, err = d.readBuffer(ft, r)
		if err != nil {
			return nil, err
		}
	}
	for i := range ret {
		if err := detachBuffer(ret[i], buffers); err != nil {
			return nil, err
		}
	}
	return ret, nil
}

func (d *Decoder) readUint64FromText(r byteReader) (uint64, bool, error) {
	ret := uint64(0)
	hasRead := false
	for {
		b, err := r.ReadByte()
		if err != nil {
			if hasRead {
				return ret, true, nil
			}
			return 0, false, err
		}
		if !('0' <= b && b <= '9') {
			err = r.UnreadByte()
			return ret, hasRead, err
		}
		hasRead = true
		ret = ret*10 + uint64(b-'0')
	}
}

func (d *Decoder) readString(r byteReader, until byte) (string, error) {
	var ret bytes.Buffer
	hasRead := false
	for {
		b, err := r.ReadByte()
		if err != nil {
			if hasRead {
				return ret.String(), nil
			}
			return "", err
		}
		if b == until {
			return ret.String(), nil
		}
		if err := ret.WriteByte(b); err != nil {
			return "", err
		}
		hasRead = true
	}
}

func (d *Decoder) readHeader(header *Header) (uint64, error) {
	typ, err := d.packetReader.ReadByte()
	if err != nil {
		return 0, err
	}
	header.Type = Type(typ - '0')
	if header.Type > BinaryAck {
		return 0, ErrInvalidPackageType
	}

	num, hasNum, err := d.readUint64FromText(d.packetReader)
	if err != nil {
		if err == io.EOF {
			err = nil
		}
		return 0, err
	}
	nextByte, err := d.packetReader.ReadByte()
	if err != nil {
		header.ID = num
		header.NeedAck = hasNum
		if err == io.EOF {
			err = nil
		}
		return 0, err
	}

	// check if buffer count
	var bufferCount uint64
	if nextByte == '-' {
		bufferCount = num
		hasNum = false
		num = 0
	} else {
		_ = d.packetReader.UnreadByte()
	}

	// check namespace
	nextByte, err = d.packetReader.ReadByte()
	if err != nil {
		if err == io.EOF {
			err = nil
		}
		return bufferCount, err
	}
	if nextByte == '/' {
		_ = d.packetReader.UnreadByte()
		header.Namespace, err = d.readString(d.packetReader, ',')
		if err != nil {
			if err == io.EOF {
				err = nil
			}
			return bufferCount, err
		}
	} else {
		_ = d.packetReader.UnreadByte()
	}

	// read id
	header.ID, header.NeedAck, err = d.readUint64FromText(d.packetReader)
	if err != nil {
		if err == io.EOF {
			err = nil
		}
		return bufferCount, err
	}
	if !header.NeedAck {
		// 313["data"], id has beed read at beginning, need add back.
		header.ID = num
		header.NeedAck = hasNum
	}

	return bufferCount, err
}

func (d *Decoder) readEvent(event *string) error {
	b, err := d.packetReader.ReadByte()
	if err != nil {
		return err
	}
	if b != '[' {
		return d.packetReader.UnreadByte()
	}
	var buf bytes.Buffer
	for {
		b, err := d.packetReader.ReadByte()
		if err != nil {
			return err
		}
		if b == ',' {
			break
		}
		if b == ']' {
			_ = d.packetReader.UnreadByte()
			break
		}
		buf.WriteByte(b)
	}
	return json.Unmarshal(buf.Bytes(), event)
}

func (d *Decoder) readBuffer(ft FrameType, r io.ReadCloser) ([]byte, error) {
	defer r.Close()
	if ft != Binary {
		return nil, ErrShouldBinaryPackageType
	}
	return ioutil.ReadAll(r)
}
