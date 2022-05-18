package go_socketio_parser

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"reflect"
)

const brByte = byte('\n')

// Marshal packet header with request payload.
func Marshal(packet *Packet) ([]byte, error) {
	if packet == nil {
		return nil, errors.New("empty packet source")
	}

	buf := bytes.NewBuffer(nil)
	bw := bufio.NewWriter(buf)

	buffers, err := writePacket(bw, packet.Header, packet.Data)
	if err != nil {
		return nil, err
	}

	// write binary data.
	if len(buffers) > 0 {
		_ = bw.WriteByte(brByte)

		for _, b := range buffers {
			_, err = bw.Write(b)
			if err != nil {
				return nil, err
			}
		}
	}

	_ = bw.Flush()

	return ioutil.ReadAll(buf)
}

type byteWriter interface {
	io.Writer
	WriteByte(byte) error
}

const binaryTypeShift = 3

func writePacket(bw byteWriter, h Header, data []interface{}) ([][]byte, error) {
	var max uint64
	buffers, err := attachBuffer(reflect.ValueOf(data), &max)
	if err != nil {
		return nil, err
	}

	// if client send data, but use Event or Ack we will upgrade header type to binary.
	if len(buffers) > 0 && (h.Type == Event || h.Type == Ack) {
		h.Type += binaryTypeShift
	}

	// packet type
	if err = bw.WriteByte(byte(h.Type + '0')); err != nil {
		return nil, err
	}

	// type of binary attachments with '-'
	if h.Type == BinaryAck || h.Type == BinaryEvent {
		if err = writeUint64(bw, max); err != nil {
			return nil, err
		}

		if err = bw.WriteByte('-'); err != nil {
			return nil, err
		}
	}

	// namespace
	if h.Namespace != "" {
		if _, err = bw.Write([]byte(h.Namespace)); err != nil {
			return nil, err
		}

		if h.ID != 0 || data != nil {
			if err = bw.WriteByte(','); err != nil {
				return nil, err
			}
		}
	}

	//acknowledgment id
	if h.IsNeedAck() {
		if err = writeUint64(bw, h.ID); err != nil {
			return nil, err
		}
	}

	// JSON-stringified payload without binary
	if data != nil {
		jsonData, err := json.Marshal(data)
		if err != nil {
			return nil, err
		}

		if _, err = bw.Write(jsonData); err != nil {
			return nil, err
		}
	}

	return buffers, nil
}

func writeUint64(w byteWriter, i uint64) error {
	base := uint64(1)
	for i/base >= 10 {
		base *= 10
	}
	for base > 0 {
		p := i / base
		if err := w.WriteByte(byte(p) + '0'); err != nil {
			return err
		}
		i -= p * base
		base /= 10
	}

	return nil
}
