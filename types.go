package go_socketio_parser

import "io"

type FrameReader interface {
	NextReader() (FrameType, io.ReadCloser, error)
}

type FrameWriter interface {
	NextWriter(ft FrameType) (io.Writer, error)
}
