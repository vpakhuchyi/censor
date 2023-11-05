package parser

import (
	"reflect"

	"github.com/vpakhuchyi/censor/internal/models"
)

// Bool parses a boolean and returns a Bool.
// Note: this method panics if the provided value is not a boolean.
func (p *Parser) Bool(boolValue reflect.Value) models.Bool {
	if boolValue.Kind() != reflect.Bool {
		panic("provided value is not a boolean")
	}

	return models.Bool{
		Value: boolValue.Interface(),
		Kind:  reflect.Bool,
	}
}
