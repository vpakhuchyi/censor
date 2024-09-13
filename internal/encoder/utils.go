package encoder

import (
	"regexp"
)

// compileExcludePatterns compiles regexp patterns from ExcludePatterns.
// Note: this method may panic if regexp pattern is invalid.
func compileExcludePatterns(e *baseEncoder) {
	if e.ExcludePatterns != nil {
		e.ExcludePatternsCompiled = make([]*regexp.Regexp, len(e.ExcludePatterns))
		for i, pattern := range e.ExcludePatterns {
			e.ExcludePatternsCompiled[i] = regexp.MustCompile(pattern)
		}
	}
}
