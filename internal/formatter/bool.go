package formatter

import (
	"fmt"
	"reflect"

	"github.com/vpakhuchyi/censor/internal/models"
)

// Bool formats a value as a boolean.
func (f *Formatter) Bool(v models.Value) string {
	if v.Kind != reflect.Bool {
		panic("provided value is not a boolean")
	}

	return fmt.Sprintf(`%v`, v.Value)
}
