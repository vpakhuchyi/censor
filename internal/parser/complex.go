package parser

import (
	"reflect"

	"github.com/vpakhuchyi/censor/internal/models"
)

// Complex parses a complex and returns a models.Value.
// Note: this method panics if the provided value is not a complex.
func (p *Parser) Complex(rv reflect.Value) models.Value {
	if rv.Kind() != reflect.Complex64 && rv.Kind() != reflect.Complex128 {
		panic("provided value is not a complex")
	}

	return models.Value{
		Value: rv.Interface(),
		Kind:  rv.Kind(),
	}
}
