package encoder

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_compileExcludePatterns(t *testing.T) {
	require.NotPanics(t, func() {
		e := baseEncoder{ExcludePatterns: []string{`^A$`}}
		compileExcludePatterns(&e)

		require.Equal(t, len(e.ExcludePatterns), len(e.ExcludePatternsCompiled))
	})
}
