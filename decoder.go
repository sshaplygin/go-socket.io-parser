package parser

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"reflect"
	"strings"
)

// FrameReader
type FrameReader interface {
	NextReader() (FrameType, io.ReadCloser, error)
}

// Decoder
type Decoder struct {
	fr           FrameReader
	lastFrame    io.ReadCloser
	packetReader byteReader

	bufferCount uint64
	isEvent     bool
}

// NewDecoder
func NewDecoder(fr FrameReader) *Decoder {
	return &Decoder{
		fr: fr,
	}
}

// Close
func (d *Decoder) Close() error {
	if d.lastFrame != nil {
		d.lastFrame.Close()
		d.lastFrame = nil
	}
	return nil
}

type byteReader interface {
	io.Reader
	ReadByte() (byte, error)
	UnreadByte() error
}

// DiscardLast
func (d *Decoder) DiscardLast() error {
	if d.lastFrame != nil {
		err := d.lastFrame.Close()
		d.lastFrame = nil
		return err
	}
	return nil
}

// DecodeHeader
func (d *Decoder) DecodeHeader(header *Header, event *string) error {
	ft, r, err := d.fr.NextReader()
	if err != nil {
		return err
	}
	if ft != TEXT {
		return errors.New("first packet should be TEXT frame")
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
		if err := d.readEvent(event); err != nil {
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
		if err := d.detachBuffer(ret[i], buffers); err != nil {
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
		d.packetReader.UnreadByte()
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
		d.packetReader.UnreadByte()
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
		err = d.packetReader.UnreadByte()
		return err
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
			d.packetReader.UnreadByte()
			break
		}
		buf.WriteByte(b)
	}
	return json.Unmarshal(buf.Bytes(), event)
}

func (d *Decoder) readBuffer(ft FrameType, r io.ReadCloser) ([]byte, error) {
	defer r.Close()
	if ft != BINARY {
		return nil, ErrShouldBinaryPackageType
	}
	return ioutil.ReadAll(r)
}

func (d *Decoder) detachBuffer(v reflect.Value, buffers []Buffer) error {
	for v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface {
		v = v.Elem()
	}
	switch v.Kind() {
	case reflect.Struct:
		if v.Type().Name() == "Buffer" {
			if !v.CanAddr() {
				return ErrBufferAddress
			}
			buffer := v.Addr().Interface().(*Buffer)
			if buffer.IsBinary {
				*buffer = buffers[buffer.Num]
			}
			return nil
		}
		for i := 0; i < v.NumField(); i++ {
			if err := d.detachBuffer(v.Field(i), buffers); err != nil {
				return err
			}
		}
	case reflect.Map:
		for _, key := range v.MapKeys() {
			if err := d.detachBuffer(v.MapIndex(key), buffers); err != nil {
				return err
			}
		}
	case reflect.Array:
		fallthrough
	case reflect.Slice:
		for i := 0; i < v.Len(); i++ {
			if err := d.detachBuffer(v.Index(i), buffers); err != nil {
				return err
			}
		}
	}
	return nil
}
