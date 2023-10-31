package formatter

import (
	"fmt"

	"github.com/vpakhuchyi/censor/internal/models"
)

// String formats a value as a string.
func (f *Formatter) String(v models.Value) string {
	return fmt.Sprintf(`%s`, v.Value)
}
