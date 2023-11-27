package formatter

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/vpakhuchyi/censor/config"
	"github.com/vpakhuchyi/censor/internal/models"
)

func TestFormatter_Ptr(t *testing.T) {
	f := Formatter{
		maskValue:         config.DefaultMaskValue,
		displayStructName: false,
		displayMapType:    false,
	}

	t.Run("successful", func(t *testing.T) {
		require.NotPanics(t, func() {
			v := models.Ptr{Value: 1, Kind: reflect.Int}
			got := f.Ptr(v)
			exp := "&1"
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
			v := models.Ptr{Value: "hell0", Kind: reflect.String}
			got := f.Ptr(v)
			exp := "&[CENSORED]"
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
