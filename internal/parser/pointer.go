package parser

import (
	"reflect"

	"github.com/vpakhuchyi/censor/internal/models"
)

// Ptr parses a given value and returns a Ptr.
// If the value is nil, it returns a Ptr with a nil Value.
// Note: this method panics if the provided value is not a complex.
//
//nolint:exhaustive,gocyclo
func (p *Parser) Ptr(rv reflect.Value) models.Ptr {
	if rv.Kind() != reflect.Pointer {
		panic("provided value is not a pointer")
	}

	if rv.IsNil() {
		return models.Ptr{Value: nil, Kind: reflect.Pointer}
	}

	switch rv.Elem().Kind() {
	case reflect.Struct:
		return models.Ptr{Value: p.Struct(rv.Elem()), Kind: rv.Elem().Kind()}
	case reflect.Slice, reflect.Array:
		return models.Ptr{Value: p.Slice(rv.Elem()), Kind: rv.Elem().Kind()}
	case reflect.Pointer:
		return models.Ptr{Value: p.Ptr(rv.Elem()), Kind: rv.Elem().Kind()}
	case reflect.Map:
		return models.Ptr{Value: p.Map(rv.Elem()), Kind: rv.Elem().Kind()}
	case reflect.Interface:
		return models.Ptr{Value: p.Interface(rv.Elem()), Kind: rv.Elem().Kind()}
	case reflect.String:
		return models.Ptr(p.String(rv.Elem()))
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return models.Ptr(p.Integer(rv.Elem()))
	case reflect.Float32, reflect.Float64:
		return models.Ptr(p.Float(rv.Elem()))
	case reflect.Bool:
		return models.Ptr(p.Bool(rv.Elem()))
	case reflect.Complex64, reflect.Complex128:
		return models.Ptr(p.Complex(rv.Elem()))
	default:
		return models.Ptr{Value: rv.Elem().Interface(), Kind: rv.Elem().Kind()}
	}
}
