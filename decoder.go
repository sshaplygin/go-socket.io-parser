package go_socketio_parser

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"strconv"
)

const binarySep = byte('-')
const nsSep = byte('/')
const nsEndSep = byte(',')
const dataOpenSep = byte('[')
const dataCloseSep = byte(']')
const payloadSep = byte(',')
const attachBinarySep = byte('\n')
const bufferOpenDataSep = byte('{')
const bufferCloseDataSep = byte('}')

func Unmarshal(data []byte, message *Packet) error {
	if len(data) == 0 {
		return errors.New("empty input data")
	}
	if message == nil {
		return errors.New("empty output header destination")
	}

	r := bytes.NewReader(data)

	// read <packet type>
	nextByte, err := r.ReadByte()
	if err != nil {
		return err
	}

	ht := Type(nextByte - zeroNumberByte)
	if !ht.IsValid() {
		return ErrInvalidPackageType
	}
	message.Header.Type = ht

	num, err := readUint64(r)
	if err != nil && err != io.EOF {
		return err
	}

	if ht.IsBinary() {
		// count binary attachments: <count of binary attachments>
		if num != 0 {
			nextByte, err = r.ReadByte()
			if err != nil && err != io.EOF {
				return err
			}

			// "-"
			if nextByte == binarySep {
				message.Header.Type -= binaryTypeShift
			}
		}
	} else if num != 0 {
		message.Header.ID = num
	}

	nextByte, err = r.ReadByte()
	if err != nil && err != io.EOF {
		return err
	}

	if err == io.EOF {
		return nil
	}

	// read namespace
	if nextByte == nsSep {
		_ = r.UnreadByte()

		ns, err := readString(r)
		if err != nil && err != io.EOF {
			return err
		}

		message.Header.Namespace = ns

		if err == io.EOF {
			return nil
		}
	}

	if message.Header.ID == 0 {
		//notice: this case handle example: binary packet type, buffer and without nsp.
		if isNumberByte(nextByte) {
			_ = r.UnreadByte()
		}

		// read acknowledgment id
		id, err := readUint64(r)
		if err != nil && err != io.EOF {
			return err
		}

		if id != 0 {
			message.Header.ID = id
		}
	}

	if err == io.EOF {
		return nil
	}

	// hotfix
	if nextByte == '[' {
		_ = r.UnreadByte()
	}

	decodeData, err := decodeData(r)
	if err != nil && err != io.EOF {
		return err
	}

	message.Data = decodeData

	return nil
}

const zeroNumberByte = byte('0')
const nineNumberByte = byte('9')

func isNumberByte(b byte) bool {
	return zeroNumberByte <= b && b <= nineNumberByte
}

func readUint64(r *bytes.Reader) (uint64, error) {
	var res uint64

	for {
		b, err := r.ReadByte()
		if err != nil && err != io.EOF {
			return 0, err
		}
		if err == io.EOF {
			return res, nil
		}

		if !(zeroNumberByte <= b && b <= nineNumberByte) {
			_ = r.UnreadByte()
			return res, nil
		}

		res = res*10 + uint64(b-zeroNumberByte)
	}
}

func readString(r *bytes.Reader) (string, error) {
	var ret bytes.Buffer

	for {
		b, err := r.ReadByte()
		if err != nil && err != io.EOF {
			return "", err
		}
		if err == io.EOF || b == nsEndSep {
			return ret.String(), nil
		}

		if err = ret.WriteByte(b); err != nil {
			return "", err
		}
	}
}

func readJSONPayload(r *bytes.Reader) (interface{}, error) {
	var buf bytes.Buffer
	var attachJSON Buffer

	for {
		b, err := r.ReadByte()
		if err != nil {
			return nil, err
		}

		buf.WriteByte(b)

		if b == bufferCloseDataSep {
			break
		}
	}

	if err := json.Unmarshal(buf.Bytes(), &attachJSON); err != nil {
		return nil, err
	}

	return &attachJSON, nil
}

func decodeData(r *bytes.Reader) ([]interface{}, error) {
	b, err := r.ReadByte()
	if err != nil && err != io.EOF {
		return nil, err
	}

	if err == io.EOF {
		return nil, err
	}

	if b != dataOpenSep {
		return nil, errors.New("invalid data segment")
	}

	bufferIdx := map[int]struct{}{}
	var data []interface{}

	//parse JSON-stringified payload.
	var countData int
	var buf bytes.Buffer

	for {
		b, err := r.ReadByte()
		if err != nil && err != io.EOF {
			return nil, err
		}
		if b == '"' {
			continue
		}
		if b == bufferOpenDataSep {
			_ = r.UnreadByte()
			jsonPayload, err := readJSONPayload(r)
			if err != nil {
				return nil, err
			}
			bufferIdx[countData] = struct{}{}

			data = append(data, jsonPayload)
			countData++

			continue
		}
		if b == dataCloseSep {
			// todo action with buf
			break
		}

		if b == payloadSep && buf.Len() == 0 {
			continue
		}

		if b == payloadSep {
			data = append(data, buf.String())
			countData++

			buf.Reset()

			continue
		}

		buf.WriteByte(b)
	}

	nextByte, err := r.ReadByte()
	if err != nil && err != io.EOF {
		return nil, err
	}
	if len(bufferIdx) > 0 && (err == io.EOF || nextByte != attachBinarySep) {
		return nil, errors.New("not found binary attachments")
	}

	if buf.Len() > 0 {
		str := buf.String()
		val, err := strconv.Atoi(str)
		if err != nil {
			data = append(data, str)
		} else {
			data = append(data, val)
		}
	}

	for idx := range bufferIdx {
		var buf bytes.Buffer
		for {
			b, err := r.ReadByte()
			if err != nil && err != io.EOF {
				return nil, err
			}

			if buf.Len() == 0 && err == io.EOF {
				break
			}

			if err == io.EOF || b == attachBinarySep {
				bufferData, ok := data[idx].(*Buffer)
				// notice: so strange behaviour
				if !ok {
					return nil, errors.New("")
				}

				bufferData.Data = buf.Bytes()
				data[idx] = bufferData

				buf.Reset()

				break
			}

			buf.WriteByte(b)
		}
	}

	return data, nil
}
