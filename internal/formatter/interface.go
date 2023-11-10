package formatter

import (
	"reflect"
	"strings"

	"github.com/vpakhuchyi/censor/internal/models"
)

// Interface formats a dynamic value of the provided interface as a string.
// Formatting rules depend on the underlying type of the value.
func (f *Formatter) Interface(v models.Value) string {
	if v.Kind != reflect.Interface {
		panic("provided value is not an interface")
	}

	if v.Value == nil {
		return "nil"
	}

	buf := strings.Builder{}
	f.writeValue(&buf, v.Value.(models.Value))

	return buf.String()
}
