package parser

import (
	"fmt"
	"reflect"

	"github.com/vpakhuchyi/censor/config"
	"github.com/vpakhuchyi/censor/internal/models"
)

// Slice parses a given value and returns a Slice.
// This function is also can be used to parse an array.
// All supported complex types will be parsed recursively.
// Note: this method panics if the provided value is not a complex.
//
//nolint:exhaustive,gocyclo
func (p *Parser) Slice(rv reflect.Value) models.Slice {
	if rv.Kind() != reflect.Slice && rv.Kind() != reflect.Array {
		panic("provided value is not a slice/array")
	}

	slice := models.Slice{Values: make([]models.Value, 0, rv.Len())}
	for i := 0; i < rv.Len(); i++ {
		elem := rv.Index(i)
		switch k := elem.Kind(); k {
		case reflect.Struct:
			slice.Values = append(slice.Values, models.Value{Value: p.Struct(elem), Kind: reflect.Struct})
		case reflect.Pointer:
			slice.Values = append(slice.Values, models.Value{Value: p.Ptr(elem), Kind: reflect.Pointer})
		case reflect.Slice, reflect.Array:
			slice.Values = append(slice.Values, models.Value{Value: p.Slice(elem), Kind: k})
		case reflect.Map:
			slice.Values = append(slice.Values, models.Value{Value: p.Map(elem), Kind: reflect.Map})
		case reflect.Interface:
			slice.Values = append(slice.Values, models.Value{Value: p.Interface(elem), Kind: reflect.Interface})
		case reflect.String:
			slice.Values = append(slice.Values, p.String(elem))
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
			slice.Values = append(slice.Values, p.Integer(elem))
		case reflect.Float32, reflect.Float64:
			slice.Values = append(slice.Values, p.Float(elem))
		case reflect.Bool:
			slice.Values = append(slice.Values, p.Bool(elem))
		case reflect.Complex64, reflect.Complex128:
			slice.Values = append(slice.Values, p.Complex(elem))
		default:
			slice.Values = append(slice.Values, models.Value{Value: fmt.Sprintf(config.UnsupportedTypeTmpl, k), Kind: k})
		}
	}

	return slice
}
