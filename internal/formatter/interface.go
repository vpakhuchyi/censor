package formatter

import (
	"strings"

	"github.com/vpakhuchyi/sanitiser/internal/models"
)

// Interface formats a dynamic value of the provided interface as a string.
// Formatting rules depend on the underlying type of the value.
func (f *Formatter) Interface(i models.Interface) string {
	if i.Value.Value == nil {
		return "nil"
	}

	buf := strings.Builder{}
	f.writeValue(&buf, i.Value)

	return buf.String()
}
