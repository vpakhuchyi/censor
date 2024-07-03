package censor

import (
	"reflect"

	jsoniter "github.com/json-iterator/go"
	"github.com/modern-go/reflect2"
)

//nolint:godot
const (
	// defaultCensorFieldTag is a default tag name for censor fields.
	defaultCensorFieldTag = "censor"
	// defaultMaskValue is a default value to mask sensitive data.
	display = "display"
)

func newExtension(c *Config) *extension {
	return &extension{
		structEncoder: newStructEncoder(c.MaskValue),
		StringEncoder: newStringEncoder(c.MaskValue, c.ExcludePatterns),
	}
}

type extension struct {
	structEncoder *structEncoder
	StringEncoder *stringEncoder
	jsoniter.EncoderExtension
}

// UpdateStructDescriptor updates the encoder for each field in the provided structDescriptor.
// If the censor tag for a field is not "display", the field's encoder is replaced with the structEncoder.
func (e extension) UpdateStructDescriptor(structDescriptor *jsoniter.StructDescriptor) {
	for _, f := range structDescriptor.Fields {
		if f.Field.Tag().Get(defaultCensorFieldTag) != display {
			f.Encoder = e.structEncoder
		}
	}
}

// DecorateEncoder is used to decorate the JSON encoder for a specific type.
// In this case, it is used to decorate the string encoder to use the custom one.
func (e extension) DecorateEncoder(typ reflect2.Type, encoder jsoniter.ValEncoder) jsoniter.ValEncoder {
	if typ.Type1().Kind() == reflect.String {
		e.StringEncoder.ValEncoder = encoder

		return e.StringEncoder
	}

	return encoder
}
