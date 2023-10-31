package formatter

import (
	"strings"

	"github.com/vpakhuchyi/censor/internal/models"
)

// Ptr formats a pointer into a string.
// It adds the `&` prefix to the formatted value to indicate that it is a pointer.
// Pointer value is formatted according to the rules of the underlying type.
func (f *Formatter) Ptr(ptr models.Ptr) string {
	if ptr.Value.Value == nil {
		return "nil"
	}

	buf := strings.Builder{}
	buf.WriteString("&")
	f.writeValue(&buf, ptr.Value)

	return buf.String()
}
