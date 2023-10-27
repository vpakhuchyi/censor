package parser

import (
	"reflect"

	"github.com/vpakhuchyi/sanitiser/internal/models"
)

// Interface parses an interface and returns an Interface.
func (p *Parser) Interface(interfaceValue reflect.Value) models.Interface {
	var v models.Value

	switch interfaceValue.Elem().Kind() {
	case reflect.Struct:
		v = models.Value{Value: p.Struct(interfaceValue.Elem()), Kind: reflect.Struct}
	case reflect.Pointer:
		v = models.Value{Value: p.Ptr(interfaceValue.Elem()), Kind: reflect.Pointer}
	case reflect.Slice, reflect.Array:
		v = models.Value{Value: p.Slice(interfaceValue.Elem()), Kind: interfaceValue.Elem().Kind()}
	case reflect.Map:
		v = models.Value{Value: p.Map(interfaceValue.Elem()), Kind: interfaceValue.Elem().Kind()}
	default:
		v = models.Value{Value: interfaceValue.Elem().Interface(), Kind: interfaceValue.Elem().Kind()}
	}

	return models.Interface{
		Name:  interfaceValue.Type().Name(),
		Value: v,
	}
}
