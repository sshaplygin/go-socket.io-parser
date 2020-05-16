package parser

import "github.com/googollee/go-engine.io/base"

// FrameType is the type of frames.
type FrameType byte

const (
	// FrameString identifies a string frame.
	FrameString FrameType = iota
	// FrameBinary identifies a binary frame.
	FrameBinary
)

// FrameType is type of a message frame.
type FrameType base.FrameType

const (
	// TEXT is text type message.
	TEXT = FrameType(base.FrameString)
	// BINARY is binary type message.
	BINARY = FrameType(base.FrameBinary)
)
