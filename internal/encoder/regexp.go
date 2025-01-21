package encoder

import (
	"regexp"
	"strings"
)

// compileExcludePatterns returns compiled regexp patterns joined by "|".
// Note: this method may panic if any regexp pattern is invalid.
func compileRegexpPatterns(patterns []string) *regexp.Regexp {
	if len(patterns) == 0 {
		return nil
	}

	return regexp.MustCompile(strings.Join(patterns, "|"))
}
