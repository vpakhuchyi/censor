package parser

import (
	"reflect"

	"github.com/vpakhuchyi/censor/internal/models"
)

// Complex parses a complex and returns a models.Value.
// Note: this method panics if the provided value is not a complex.
func (p *Parser) Complex(v reflect.Value) models.Value {
	if v.Kind() != reflect.Complex64 && v.Kind() != reflect.Complex128 {
		panic("provided value is not a complex")
	}

	return models.Value{
		Value: v.Interface(),
		Kind:  v.Kind(),
	}
}
