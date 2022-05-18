package go_socketio_parser

import "errors"

var (
	// ErrInvalidPackageType type
	ErrInvalidPackageType = errors.New("invalid packet type")
	// ErrShouldBinaryPackageType
	ErrShouldBinaryPackageType = errors.New("packet should be BINARY")
	//ErrShouldTextPackageType
	ErrShouldTextPackageType = errors.New("first packet should be TEXT frame")

	// ErrBufferAddress
	ErrBufferAddress = errors.New("invalid buffer address")
)
