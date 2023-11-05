package parser

import (
	"reflect"

	"github.com/vpakhuchyi/censor/internal/models"
)

// String parses a string and returns a models.Value.
// Note: this method panics if the provided value is not a string.
func (p *Parser) String(stringValue reflect.Value) models.Value {
	if stringValue.Kind() != reflect.String {
		panic("provided value is not a string")
	}

	return models.Value{
		Value: stringValue.Interface(),
		Kind:  reflect.String,
	}
}
