package parser

import (
	"reflect"

	"github.com/vpakhuchyi/sanitiser/internal/models"
)

// Slice parses a given value and returns a Slice.
// If value is a struct/pointer/slice/array, it will be parsed recursively.
//
//nolint:exhaustive
func (p *Parser) Slice(sliceValue reflect.Value) models.Slice {
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
		default:
			slice.Values = append(slice.Values, models.Value{Value: elem.Interface(), Kind: elem.Kind()})
		}
	}

	return slice
}
