package parser

import (
	"reflect"

	"github.com/vpakhuchyi/censor/internal/models"
)

// Integer parses an integer (int/int8/int16/int32/int64/uint/uint8/uint16/uint32/uint64/byte/rune) and returns a models.Value.
// Note: this method panics if the provided value is not an integer.
//
//nolint:exhaustive
func (p *Parser) Integer(rv reflect.Value) models.Value {
	switch k := rv.Kind(); k {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return models.Value{
			Value: rv.Interface(),
			Kind:  k,
		}
	default:
		panic("provided value is not an integer")
	}
}
