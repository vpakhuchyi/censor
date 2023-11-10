package formatter

import (
	"fmt"
	"reflect"

	"github.com/vpakhuchyi/censor/internal/models"
)

// Integer formats a value as an integer.
// Note: this method panics if the provided value is not an integer.
//
//nolint:exhaustive
func (f *Formatter) Integer(v models.Value) string {
	switch v.Kind {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return fmt.Sprintf(`%d`, v.Value)
	default:
		panic("provided value is not an integer")
	}
}
