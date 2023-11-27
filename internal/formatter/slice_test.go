package formatter

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/vpakhuchyi/censor/config"
	"github.com/vpakhuchyi/censor/internal/models"
)

func TestFormatter_Slice(t *testing.T) {
	f := Formatter{
		maskValue:         config.DefaultMaskValue,
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
			maskValue:               config.DefaultMaskValue,
			displayStructName:       false,
			displayMapType:          false,
			excludePatterns:         []string{`\d`, `\b[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Z|a-z]{2,}\b`},
			excludePatternsCompiled: excludePatternsCompiled,
		}

		require.NotPanics(t, func() {
			v := models.Slice{
				Values: []models.Value{
					{Value: "hell0", Kind: reflect.String},
					{Value: "hello", Kind: reflect.String},
				},
			}
			got := f.Slice(v)
			exp := "[[CENSORED], hello]"
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
