package formatter

import (
	"fmt"

	"github.com/vpakhuchyi/sanitiser/internal/models"
)

// Bool formats a value as a boolean.
func (f *Formatter) Bool(v models.Value) string {
	return fmt.Sprintf(`%v`, v.Value)
}
