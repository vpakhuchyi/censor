package parser

import (
	"reflect"

	"github.com/vpakhuchyi/sanitiser/internal/models"
)

// Ptr parses a given value and returns a Ptr.
// If the value is nil, it returns a Ptr with a nil Value.
//
//nolint:exhaustive
func (p *Parser) Ptr(ptrValue reflect.Value) models.Ptr {
	if ptrValue.IsNil() {
		return models.Ptr{Value: models.Value{Value: nil, Kind: reflect.Ptr}}
	}

	switch ptrValue.Elem().Kind() {
	case reflect.Struct:
		return models.Ptr{Value: models.Value{Value: p.Struct(ptrValue.Elem()), Kind: ptrValue.Elem().Kind()}}
	case reflect.Slice, reflect.Array:
		return models.Ptr{Value: models.Value{Value: p.Slice(ptrValue.Elem()), Kind: ptrValue.Elem().Kind()}}
	case reflect.Ptr:
		return models.Ptr{Value: models.Value{Value: p.Ptr(ptrValue.Elem()), Kind: ptrValue.Elem().Kind()}}
	default:
		return models.Ptr{Value: models.Value{Value: ptrValue.Elem().Interface(), Kind: ptrValue.Elem().Kind()}}
	}
}
