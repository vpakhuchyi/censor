package formatter

import (
	"reflect"
	"regexp"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/vpakhuchyi/censor/internal/models"
)

func TestFormatter_Ptr(t *testing.T) {
	f := Formatter{
		maskValue:            "[CENSORED]",
		displayPointerSymbol: false,
		displayStructName:    false,
		displayMapType:       false,
	}

	t.Run("successful", func(t *testing.T) {
		require.NotPanics(t, func() {
			v := models.Ptr{Value: 1, Kind: reflect.Int}
			got := f.Ptr(v)
			exp := "1"
			require.Equal(t, exp, got)
		})
	})

	t.Run("successful_with_pointer_symbol", func(t *testing.T) {
		f := Formatter{
			maskValue:            "[CENSORED]",
			displayPointerSymbol: true,
		}

		require.NotPanics(t, func() {
			v := models.Ptr{Value: 1, Kind: reflect.Int}
			got := f.Ptr(v)
			exp := "&1"
			require.Equal(t, exp, got)
		})
	})

	t.Run("with_exclude_patterns", func(t *testing.T) {
		f := Formatter{
			maskValue:               "[CENSORED]",
			displayPointerSymbol:    false,
			displayStructName:       false,
			displayMapType:          false,
			excludePatterns:         []string{`\d`, `\b[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Z|a-z]{2,}\b`, `l{2}`},
			excludePatternsCompiled: []*regexp.Regexp{compiledRegExpDigit, compiledRegExpEmail, regexp.MustCompile(`l{2}`)},
		}

		require.NotPanics(t, func() {
			v := models.Ptr{Value: "hello", Kind: reflect.String}
			got := f.Ptr(v)
			exp := "he[CENSORED]o"
			require.Equal(t, exp, got)
		})
	})

	t.Run("nil", func(t *testing.T) {
		require.NotPanics(t, func() {
			v := models.Ptr{Value: nil, Kind: reflect.Ptr}
			got := f.Ptr(v)
			exp := "nil"
			require.Equal(t, exp, got)
		})
	})
}
