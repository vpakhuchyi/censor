package parser

import (
	"reflect"

	"github.com/vpakhuchyi/censor/internal/models"
)

// Slice parses a given value and returns a Slice.
// This function is also can be used to parse an array.
// All supported complex types will be parsed recursively.
// Note: this method panics if the provided value is not a complex.
//
//nolint:exhaustive,gocyclo
func (p *Parser) Slice(sliceValue reflect.Value) models.Slice {
	if sliceValue.Kind() != reflect.Slice && sliceValue.Kind() != reflect.Array {
		panic("provided value is not a slice/array")
	}

	var slice models.Slice
	for i := 0; i < sliceValue.Len(); i++ {
		elem := sliceValue.Index(i)
		switch elem.Kind() {
		case reflect.Struct:
			slice.Values = append(slice.Values, models.Value{Value: p.Struct(elem), Kind: reflect.Struct})
		case reflect.Pointer:
			slice.Values = append(slice.Values, models.Value{Value: p.Ptr(elem), Kind: reflect.Pointer})
		case reflect.Slice, reflect.Array:
			slice.Values = append(slice.Values, models.Value{Value: p.Slice(elem), Kind: elem.Kind()})
		case reflect.Map:
			slice.Values = append(slice.Values, models.Value{Value: p.Map(elem), Kind: elem.Kind()})
		case reflect.Interface:
			slice.Values = append(slice.Values, models.Value{Value: p.Interface(elem), Kind: elem.Kind()})
		case reflect.String:
			slice.Values = append(slice.Values, p.String(elem))
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			slice.Values = append(slice.Values, p.Integer(elem))
		case reflect.Float32, reflect.Float64:
			slice.Values = append(slice.Values, p.Float(elem))
		case reflect.Bool:
			slice.Values = append(slice.Values, p.Bool(elem))
		case reflect.Complex64, reflect.Complex128:
			slice.Values = append(slice.Values, p.Complex(elem))
		}
	}

	return slice
}
