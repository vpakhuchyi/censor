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

	values := slice.Values
	for i := 0; i < len(values); i++ {
		f.writeValue(&buf, values[i])
		if i < len(values)-1 {
			buf.WriteString(", ")
		}
	}

	buf.WriteString("]")

	return buf.String()
}
