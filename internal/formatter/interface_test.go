package formatter

import (
	"reflect"
	"regexp"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/vpakhuchyi/censor/config"
	"github.com/vpakhuchyi/censor/internal/models"
)

func TestFormatter_Interface(t *testing.T) {
	f := Formatter{
		maskValue:         config.DefaultMaskValue,
		displayStructName: false,
		displayMapType:    false,
	}

	t.Run("successful", func(t *testing.T) {
		require.NotPanics(t, func() {
			v := models.Value{
				Value: models.Value{
					Value: models.Slice{
						Values: []models.Value{
							{Value: "Kholodetsʹ", Kind: reflect.String},
							{Value: "Halushky", Kind: reflect.String},
						},
					},
					Kind: reflect.Slice,
				},
				Kind: reflect.Interface,
			}

			got := f.Interface(v)
			exp := `[Kholodetsʹ, Halushky]`
			require.Equal(t, exp, got)
		})
	})

	t.Run("with_exclude_patterns", func(t *testing.T) {
		f := Formatter{
			maskValue:               config.DefaultMaskValue,
			displayStructName:       false,
			displayMapType:          false,
			excludePatterns:         []string{`\d`, `\b[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Z|a-z]{2,}\b`},
			excludePatternsCompiled: []*regexp.Regexp{compiledRegExpDigit, compiledRegExpEmail},
		}

		require.NotPanics(t, func() {
			v := models.Value{
				Value: models.Value{
					Value: models.Slice{
						Values: []models.Value{
							{Value: "Kholodetsʹ", Kind: reflect.String},
							{Value: "Halushky1", Kind: reflect.String},
						},
					},
					Kind: reflect.Slice,
				},
				Kind: reflect.Interface,
			}

			got := f.Interface(v)
			exp := `[Kholodetsʹ, Halushky[CENSORED]]`
			require.Equal(t, exp, got)
		})
	})

	t.Run("nil_value", func(t *testing.T) {
		require.NotPanics(t, func() {
			v := models.Value{Value: nil, Kind: reflect.Interface}
			got := f.Interface(v)
			exp := `nil`
			require.Equal(t, exp, got)
		})
	})

	t.Run("non_interface_value", func(t *testing.T) {
		require.PanicsWithValue(t, "provided value is not an interface", func() {
			f.Interface(models.Value{Value: 44, Kind: reflect.Int})
		})
	})
}
