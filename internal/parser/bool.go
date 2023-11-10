package parser

import (
	"reflect"

	"github.com/vpakhuchyi/censor/internal/models"
)

// Bool parses a boolean and returns a models.Value.
// Note: this method panics if the provided value is not a boolean.
func (p *Parser) Bool(rv reflect.Value) models.Value {
	if rv.Kind() != reflect.Bool {
		panic("provided value is not a boolean")
	}

	return models.Value{
		Value: rv.Interface(),
		Kind:  reflect.Bool,
	}
}
