package formatter

import (
	"reflect"

	"github.com/shopspring/decimal"

	"github.com/vpakhuchyi/censor/internal/models"
)

// Float formats a value as a float.
// To keep a consistent formatting output values are formatted according to
// the https://github.com/shopspring/decimal package rules.
// Note: this method panics if the provided value is not a float.
func (f *Formatter) Float(v models.Value) string {
	if v.Kind != reflect.Float32 && v.Kind != reflect.Float64 {
		panic("provided value is not a float")
	}

	if v.Kind == reflect.Float32 {
		return decimal.NewFromFloat32(v.Value.(float32)).String()
	}

	return decimal.NewFromFloat(v.Value.(float64)).String()
}
