package formatter

import (
	"strings"

	"github.com/vpakhuchyi/sanitiser/internal/models"
)

// Ptr formats a pointer into a string.
// It adds the `&` prefix to the formatted value to indicate that it is a pointer.
// If the pointer points to a struct, it is formatted recursively using FormatStruct function rules.
// If the pointer points to a slice or an array, it is formatted recursively using FormatSlice function rules.
func (f *Formatter) Ptr(ptr models.Ptr) string {
	if ptr.Value.Value == nil {
		return "nil"
	}

	buf := strings.Builder{}
	buf.WriteString("&")
	f.writeValue(&buf, ptr.Value)

	return buf.String()
}
