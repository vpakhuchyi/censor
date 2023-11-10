package formatter

import (
	"fmt"
	"reflect"

	"github.com/vpakhuchyi/censor/internal/models"
)

// String formats a value as a string.
func (f *Formatter) String(v models.Value) string {
	if v.Kind != reflect.String {
		panic("provided value is not a string")
	}

	return fmt.Sprintf(`%s`, v.Value)
}
