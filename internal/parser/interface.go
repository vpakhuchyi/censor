package parser

import (
	"fmt"
	"reflect"

	"github.com/vpakhuchyi/censor/config"
	"github.com/vpakhuchyi/censor/internal/models"
)

// Interface parses an interface and returns an Interface.
// In case of a pointer to unsupported type of value, a string built from config.UnsupportedTypeTmpl
// is used instead of the real value. That string contains a type of the value.
//
//nolint:exhaustive
func (p *Parser) Interface(rv reflect.Value) models.Value {
	if rv.Kind() != reflect.Interface {
		panic("provided value is not an interface")
	}

	switch k := rv.Elem().Kind(); k {
	case reflect.Struct:
		return models.Value{Value: p.Struct(rv.Elem()), Kind: reflect.Struct}
	case reflect.Pointer:
		return models.Value{Value: p.Ptr(rv.Elem()), Kind: reflect.Pointer}
	case reflect.Slice, reflect.Array:
		return models.Value{Value: p.Slice(rv.Elem()), Kind: k}
	case reflect.Map:
		return models.Value{Value: p.Map(rv.Elem()), Kind: reflect.Map}
	case reflect.String:
		return p.String(rv.Elem())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return p.Integer(rv.Elem())
	case reflect.Float32, reflect.Float64:
		return p.Float(rv.Elem())
	case reflect.Bool:
		return p.Bool(rv.Elem())
	default:
		return models.Value{
			Kind:  k,
			Value: fmt.Sprintf(config.UnsupportedTypeTmpl, k),
		}
	}
}
