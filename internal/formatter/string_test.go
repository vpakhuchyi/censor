package formatter

import (
	"reflect"
	"regexp"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/vpakhuchyi/censor/internal/models"
)

func TestFormatter_String(t *testing.T) {
	f := Formatter{
		maskValue:         "[CENSORED]",
		displayStructName: false,
		displayMapType:    false,
		excludePatterns:   nil,
	}

	t.Run("string", func(t *testing.T) {
		require.NotPanics(t, func() {
			v := models.Value{Value: "hello", Kind: reflect.String}
			got := f.String(v)
			exp := "hello"
			require.Equal(t, exp, got)
		})
	})

	t.Run("non_string_value", func(t *testing.T) {
		require.PanicsWithValue(t, "provided value is not a string", func() {
			f.String(models.Value{Value: 44, Kind: reflect.Int})
		})
	})

	t.Run("string_with_exclude_patterns", func(t *testing.T) {
		f := Formatter{
			maskValue:               "[CENSORED]",
			displayStructName:       false,
			displayMapType:          false,
			excludePatterns:         []string{`\d`, `\b[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Z|a-z]{2,}\b`},
			excludePatternsCompiled: []*regexp.Regexp{compiledRegExpDigit, compiledRegExpEmail},
		}

		require.NotPanics(t, func() {
			v := models.Value{Value: "some long text with multiple emails: ttttest@exaample.com and ttttest@exaample.com", Kind: reflect.String}
			got := f.String(v)
			exp := "some long text with multiple emails: [CENSORED] and [CENSORED]"
			require.Equal(t, exp, got)
		})
	})
}
