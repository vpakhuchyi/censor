package parser

import (
	"reflect"

	"github.com/vpakhuchyi/censor/internal/models"
)

// Float parses a float and returns a models.Value.
// Note: this method panics if the provided value is not a float.
func (p *Parser) Float(rv reflect.Value) models.Value {
	if rv.Kind() != reflect.Float32 && rv.Kind() != reflect.Float64 {
		panic("provided value is not a float")
	}

	return models.Value{
		Value: rv.Interface(),
		Kind:  rv.Kind(),
	}
}
