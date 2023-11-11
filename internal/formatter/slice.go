package formatter

import (
	"strings"

	"github.com/vpakhuchyi/censor/internal/models"
)

// Slice formats a slice or an array into a string.
// The formatting rules depend on the underlying type of the slice/array elements.
func (f *Formatter) Slice(slice models.Slice) string {
	var buf strings.Builder
	buf.WriteString("[")

	for i := 0; i < len(slice.Values); i++ {
		f.writeValue(&buf, slice.Values[i])
		if i < len(slice.Values)-1 {
			buf.WriteString(", ")
		}
	}

	buf.WriteString("]")

	return buf.String()
}
