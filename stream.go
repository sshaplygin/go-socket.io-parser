package go_socketio_parser

import "io"

type Encoder struct {
	fw FrameWriter
}

func NewEncoder(fw FrameWriter) *Encoder {
	return &Encoder{
		fw: fw,
	}
}

func (e *Encoder) Encode(h Header, attach []interface{}) error {
	//w, err := e.fw.NextWriter(Text)
	//if err != nil {
	//	return err
	//}
	//
	//buffers, err := writePacket(w, h, attach)
	//if err != nil {
	//	return err
	//}
	//
	//for _, b := range buffers {
	//	w, err = e.fw.NextWriter(Binary)
	//	if err != nil {
	//		return err
	//	}
	//	_, err = w.Write(b)
	//	if err != nil {
	//		return err
	//	}
	//}

	return nil
}

type Decoder struct {
	fr           FrameReader
	lastFrame    io.ReadCloser
	packetReader byteReader

	bufferCount uint64
	isEvent     bool
}

func NewDecoder(fr FrameReader) *Decoder {
	return &Decoder{
		fr: fr,
	}
}

func (d *Decoder) Close() error {
	var err error
	if d.lastFrame != nil {
		err = d.lastFrame.Close()
		d.lastFrame = nil
	}
	return err
}

type byteReader interface {
	io.Reader
	ReadByte() (byte, error)
	UnreadByte() error
}

func (d *Decoder) DiscardLast() error {
	if d.lastFrame != nil {
		err := d.lastFrame.Close()
		d.lastFrame = nil
		return err
	}
	return nil
}
