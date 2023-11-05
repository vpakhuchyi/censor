package parser

import (
	"reflect"

	"github.com/vpakhuchyi/censor/internal/models"
)

// Float parses a float and returns a models.Value.
// Note: this method panics if the provided value is not a float.
func (p *Parser) Float(floatValue reflect.Value) models.Value {
	if floatValue.Kind() != reflect.Float32 && floatValue.Kind() != reflect.Float64 {
		panic("provided value is not a float")
	}

	return models.Value{
		Value: floatValue.Interface(),
		Kind:  floatValue.Kind(),
	}
}
