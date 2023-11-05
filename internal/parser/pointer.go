package parser

import (
	"reflect"

	"github.com/vpakhuchyi/censor/internal/models"
)

// Ptr parses a given value and returns a Ptr.
// If the value is nil, it returns a Ptr with a nil Value.
//
//nolint:exhaustive,gocyclo
func (p *Parser) Ptr(ptrValue reflect.Value) models.Ptr {
	if ptrValue.Kind() != reflect.Pointer {
		panic("provided value is not a pointer")
	}

	if ptrValue.IsNil() {
		return models.Ptr{Value: nil, Kind: reflect.Ptr}
	}

	switch ptrValue.Elem().Kind() {
	case reflect.Struct:
		return models.Ptr{Value: p.Struct(ptrValue.Elem()), Kind: ptrValue.Elem().Kind()}
	case reflect.Slice, reflect.Array:
		return models.Ptr{Value: p.Slice(ptrValue.Elem()), Kind: ptrValue.Elem().Kind()}
	case reflect.Ptr:
		return models.Ptr{Value: p.Ptr(ptrValue.Elem()), Kind: ptrValue.Elem().Kind()}
	case reflect.Map:
		return models.Ptr{Value: p.Map(ptrValue.Elem()), Kind: ptrValue.Elem().Kind()}
	case reflect.Interface:
		return models.Ptr{Value: p.Interface(ptrValue.Elem()), Kind: ptrValue.Elem().Kind()}
	case reflect.String:
		return models.Ptr(p.String(ptrValue.Elem()))
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return models.Ptr(p.Integer(ptrValue.Elem()))
	case reflect.Float32, reflect.Float64:
		return models.Ptr(p.Float(ptrValue.Elem()))
	case reflect.Bool:
		return models.Ptr(p.Bool(ptrValue.Elem()))
	case reflect.Complex64, reflect.Complex128:
		return models.Ptr(p.Complex(ptrValue.Elem()))
	default:
		return models.Ptr{Value: ptrValue.Elem().Interface(), Kind: ptrValue.Elem().Kind()}
	}
}
