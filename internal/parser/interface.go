package parser

import (
	"reflect"

	"github.com/vpakhuchyi/censor/internal/models"
)

// Interface parses an interface and returns an Interface.
//
//nolint:exhaustive
func (p *Parser) Interface(rv reflect.Value) models.Interface {
	var v models.Value

	switch rv.Elem().Kind() {
	case reflect.Struct:
		v = models.Value{Value: p.Struct(rv.Elem()), Kind: reflect.Struct}
	case reflect.Pointer:
		v = models.Value{Value: p.Ptr(rv.Elem()), Kind: reflect.Pointer}
	case reflect.Slice, reflect.Array:
		v = models.Value{Value: p.Slice(rv.Elem()), Kind: rv.Elem().Kind()}
	case reflect.Map:
		v = models.Value{Value: p.Map(rv.Elem()), Kind: rv.Elem().Kind()}
	case reflect.String:
		v = p.String(rv.Elem())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		v = p.Integer(rv.Elem())
	case reflect.Float32, reflect.Float64:
		v = p.Float(rv.Elem())
	case reflect.Bool:
		v = p.Bool(rv.Elem())
	case reflect.Complex64, reflect.Complex128:
		v = p.Complex(rv.Elem())
	}

	return models.Interface{
		Name:  rv.Type().Name(),
		Value: v,
	}
}
