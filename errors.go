package parser

import "errors"

var (
	// ErrInvalidPackageType type
	ErrInvalidPackageType = errors.New("invalid packet type")
	// ErrShouldBinaryPackageType
	ErrShouldBinaryPackageType = errors.New("packet should be BINARY")

	// ErrBufferAddress
	ErrBufferAddress = errors.New("can't get buffer address")
)
