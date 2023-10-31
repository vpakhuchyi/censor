package formatter

import (
	"fmt"

	"github.com/vpakhuchyi/censor/internal/models"
)

// Integer formats a value as an integer.
func (f *Formatter) Integer(v models.Value) string {
	return fmt.Sprintf(`%d`, v.Value)
}
