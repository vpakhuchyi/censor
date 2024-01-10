package formatter

import (
	"fmt"
	"reflect"

	"github.com/vpakhuchyi/censor/internal/models"
)

// String formats a value as a string.
func (f *Formatter) String(v models.Value) string {
	if v.Kind != reflect.String {
		panic("provided value is not a string")
	}

	if len(f.excludePatterns) != 0 {
		for _, pattern := range f.excludePatternsCompiled {
			if s, ok := v.Value.(string); ok && pattern.MatchString(s) {
				return pattern.ReplaceAllString(s, f.maskValue)
			}
		}
	}

	return fmt.Sprintf(`%s`, v.Value)
}
