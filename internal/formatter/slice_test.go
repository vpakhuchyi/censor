package formatter

import (
	"reflect"
	"regexp"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/vpakhuchyi/censor/internal/models"
)

func TestFormatter_Slice(t *testing.T) {
	f := Formatter{
		maskValue:         "[CENSORED]",
		displayStructName: false,
		displayMapType:    false,
	}

	t.Run("successful", func(t *testing.T) {
		require.NotPanics(t, func() {
			v := models.Slice{
				Values: []models.Value{
					{Value: 1, Kind: reflect.Int},
					{Value: 2, Kind: reflect.Int},
				},
			}
			got := f.Slice(v)
			exp := "[1, 2]"
			require.Equal(t, exp, got)
		})
	})

	t.Run("with_exclude_patterns", func(t *testing.T) {
		f := Formatter{
			maskValue:               "[CENSORED]",
			displayStructName:       false,
			displayMapType:          false,
			excludePatterns:         []string{`[A-Z|a-z]{2}\d{10}`, `\b[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Z|a-z]{2,}\b`},
			excludePatternsCompiled: []*regexp.Regexp{regexp.MustCompile(`[A-Z|a-z]{2}\d{10}`), compiledRegExpEmail},
		}

		require.NotPanics(t, func() {
			v := models.Slice{
				Values: []models.Value{
					{Value: "text with IBAN: UA1234567890 that must not be shown", Kind: reflect.String},
					{Value: "hello", Kind: reflect.String},
				},
			}
			got := f.Slice(v)
			exp := "[text with IBAN: [CENSORED] that must not be shown, hello]"
			require.Equal(t, exp, got)
		})
	})

	t.Run("empty", func(t *testing.T) {
		require.NotPanics(t, func() {
			v := models.Slice{
				Values: []models.Value{},
			}
			got := f.Slice(v)
			exp := "[]"
			require.Equal(t, exp, got)
		})
	})
}
