package formatter

import (
	"fmt"
	"reflect"

	"github.com/vpakhuchyi/censor/internal/models"
)

// Float formats a value as a float.
// The value is formatted with up to 7 g places for float32 and up to 15 decimal places for float64.
func (f *Formatter) Float(v models.Value) string {
	if v.Kind != reflect.Float32 && v.Kind != reflect.Float64 {
		panic("provided value is not a float")
	}

	if v.Kind == reflect.Float32 {
		return fmt.Sprintf(`%.7g`, v.Value)
	}

	return fmt.Sprintf(`%.15g`, v.Value)
}
