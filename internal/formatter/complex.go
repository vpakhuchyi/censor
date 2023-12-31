package formatter

import (
	"fmt"
	"reflect"

	"github.com/vpakhuchyi/censor/internal/models"
)

// Complex formats a value as a complex64/128.
// The value is formatted with up to 7 decimal places for both parts of complex64 and
// up to 15 decimal places for both parts of complex128.
// Note: this method panics if the provided value is not a complex.
func (f *Formatter) Complex(v models.Value) string {
	if v.Kind != reflect.Complex64 && v.Kind != reflect.Complex128 {
		panic("provided value is not a complex")
	}

	if v.Kind == reflect.Complex64 {
		return fmt.Sprintf(`%.7g`, v.Value)
	}

	return fmt.Sprintf(`%.15g`, v.Value)
}
