package formatter

import (
	"fmt"

	"github.com/vpakhuchyi/sanitiser/internal/models"
)

// String formats a value as a string.
// The value is wrapped in double quotes.
func (f *Formatter) String(v models.Value) string {
	return fmt.Sprintf(`"%s"`, v.Value)
}
