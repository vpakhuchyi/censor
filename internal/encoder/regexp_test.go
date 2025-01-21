package encoder

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_compileRegexpPatterns(t *testing.T) {
	tests := map[string]struct {
		patterns []string
		exp      *regexp.Regexp
	}{
		"successful": {
			patterns: []string{"[0-9]", "[a-z]"},
			exp:      regexp.MustCompile("[0-9]|[a-z]"),
		},
		"empty": {
			patterns: []string{},
			exp:      nil,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			require.NotPanics(t, func() {
				got := compileRegexpPatterns(tt.patterns)
				require.Equal(t, tt.exp, got)
			})
		})
	}

	require.Panics(t, func() {
		compileRegexpPatterns([]string{`^[\w._%+-]+@[\w.-]+\.[a-zA-Z]{2,}$`, `*(invalid)`})
	})
}
