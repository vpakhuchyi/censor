package formatter

import (
	"strings"

	"github.com/vpakhuchyi/censor/internal/models"
)

// PointerSymbol is used to display the pointer to the value in the output.
const PointerSymbol = "&"

// Ptr formats a pointer into a string.
// Pointer value is formatted according to the rules of the underlying type.
// If the pointer is nil, the string "nil" is returned.
// It's possible to display the pointer symbol in the output by setting the DisplayPointerSymbol field to true
// in the Formatter configuration.
func (f *Formatter) Ptr(ptr models.Ptr) string {
	if ptr.Value == nil {
		return "nil"
	}

	buf := strings.Builder{}

	if f.displayPointerSymbol {
		buf.WriteString(PointerSymbol)
	}

	f.writeValue(&buf, models.Value(ptr))

	return buf.String()
}
