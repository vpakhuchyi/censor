package encoder

import "regexp"

// compileExcludePatterns returns compiled regexp patterns.
// Note: this method may panic if any regexp pattern is invalid.
func compileRegexpPatterns(patterns []string) []*regexp.Regexp {
	var compiled []*regexp.Regexp
	if patterns != nil {
		compiled = make([]*regexp.Regexp, len(patterns))
		for i, pattern := range patterns {
			compiled[i] = regexp.MustCompile(pattern)
		}
	}

	return compiled
}
