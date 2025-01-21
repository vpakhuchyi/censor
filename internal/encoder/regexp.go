package encoder

import (
	"regexp"
	"strings"
)

// compileExcludePatterns returns compiled regexp patterns.
// Note: this method may panic if any regexp pattern is invalid.
func compileRegexpPatterns(patterns []string) *regexp.Regexp {
	if len(patterns) > 0 {
		return regexp.MustCompile(strings.Join(patterns, "|"))
	}

	return nil
}
