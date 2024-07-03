package censor

import (
	"regexp"
	"unsafe"

	jsoniter "github.com/json-iterator/go"
)

func newStringEncoder(maskValue string, excludePatterns []string) *stringEncoder {
	se := stringEncoder{
		maskValue: maskValue,
	}

	se.compileExcludePatterns(excludePatterns)

	return &se
}

// stringEncoder is a struct that holds the configuration for encoding strings,
// a list of compiled regular expressions to exclude from encoding,
// and a jsoniter.ValEncoder for encoding values.
type stringEncoder struct {
	maskValue               string
	excludePatternsCompiled []*regexp.Regexp
	jsoniter.ValEncoder
}

// IsEmpty checks if the string pointed to by ptr is empty.
// It returns true if the string is empty, and false otherwise.
func (s stringEncoder) IsEmpty(ptr unsafe.Pointer) bool {
	return *((*string)(ptr)) == ""
}

// Encode encodes the string pointed to by v into the provided jsoniter.Stream.
// If the string matches any of the compiled exclude patterns, it is replaced with the mask value from the configuration.
// If the string does not match any exclude patterns, it is encoded using the default jsoniter.ValEncoder.
func (s stringEncoder) Encode(v unsafe.Pointer, stream *jsoniter.Stream) {
	if len(s.excludePatternsCompiled) != 0 {
		str := *(*string)(v)

		for _, pattern := range s.excludePatternsCompiled {
			if pattern.MatchString(str) {
				stream.WriteString(pattern.ReplaceAllString(str, s.maskValue))

				return
			}
		}
	}

	s.ValEncoder.Encode(v, stream)
}

// compileExcludePatterns compiles regexp patterns from excludePatterns.
// Note: this method may panic if regexp pattern is invalid.
func (s *stringEncoder) compileExcludePatterns(excludePatterns []string) {
	if len(excludePatterns) > 0 {
		s.excludePatternsCompiled = make([]*regexp.Regexp, len(excludePatterns))
		for i, pattern := range excludePatterns {
			s.excludePatternsCompiled[i] = regexp.MustCompile(pattern)
		}
	}
}
