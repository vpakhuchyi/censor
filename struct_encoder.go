package censor

import (
	"unsafe"

	jsoniter "github.com/json-iterator/go"
)

func newStructEncoder(maskValue string) *structEncoder {
	return &structEncoder{
		maskValue: maskValue,
	}
}

// structEncoder is a struct that holds the configuration for encoding structures.
type structEncoder struct {
	maskValue string
}

// IsEmpty is a method that always returns false. It's used to satisfy the jsoniter.ValEncoder interface.
func (structEncoder) IsEmpty(_ unsafe.Pointer) bool {
	return false
}

// Encode encodes the structure into the provided jsoniter.Stream.
// It replaces the structure with the mask value from the configuration.
func (s structEncoder) Encode(_ unsafe.Pointer, stream *jsoniter.Stream) {
	stream.WriteString(s.maskValue)
}
