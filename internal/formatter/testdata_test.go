package formatter

import "regexp"

var (
	compiledRegExpDigit = regexp.MustCompile(`\d`)
	compiledRegExpEmail = regexp.MustCompile(`\b[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Z|a-z]{2,}\b`)
)
