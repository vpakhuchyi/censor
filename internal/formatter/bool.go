package formatter

import (
	"fmt"

	"github.com/vpakhuchyi/censor/internal/models"
)

// Bool formats a value as a boolean.
func (f *Formatter) Bool(v models.Bool) string {
	return fmt.Sprintf(`%v`, v.Value)
}
