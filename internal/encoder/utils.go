package encoder

import (
	"reflect"
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

func parseStructName(t reflect.Type) string {
	pkgPath := t.PkgPath()
	// This custom logic is used instead of strings.Split to avoid unnecessary allocations.
	for i := len(pkgPath) - 1; i >= 0; i-- {
		// We iterate over the package path in reverse order until we find the last slash,
		// which separates the package name from the package path.
		if pkgPath[i] == '/' {
			return pkgPath[i+1:]
		}

		// If there is no slash in the package path, then the package name is equal to the package path.
		// Example: "main" package.
		if i == 0 {
			return pkgPath[i:]
		}
	}

	return ""
}
