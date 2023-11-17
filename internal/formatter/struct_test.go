package formatter

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/vpakhuchyi/censor/config"
	"github.com/vpakhuchyi/censor/internal/models"
	"github.com/vpakhuchyi/censor/internal/options"
)

func TestFormatter_Struct(t *testing.T) {
	f := Formatter{
		maskValue:         config.DefaultMaskValue,
		displayStructName: false,
		displayMapType:    false,
	}

	t.Run("successful", func(t *testing.T) {
		require.NotPanics(t, func() {
			v := models.Struct{
				Name: "Foo",
				Fields: []models.Field{
					{Name: "Foo", Value: models.Value{Value: "foo", Kind: reflect.String}, Opts: options.FieldOptions{Display: true}, Kind: reflect.String},
					{Name: "Bar", Value: models.Value{Value: 1, Kind: reflect.Int}, Opts: options.FieldOptions{Display: false}, Kind: reflect.Int},
				},
			}
			got := f.Struct(v)
			exp := `{Foo: foo, Bar: [******]}`
			require.Equal(t, exp, got)
		})
	})

	t.Run("with_display_struct_name", func(t *testing.T) {
		f := Formatter{
			maskValue:         config.DefaultMaskValue,
			displayStructName: true,
			displayMapType:    false,
		}

		require.NotPanics(t, func() {
			v := models.Struct{
				Name: "Foo",
				Fields: []models.Field{
					{Name: "Foo", Value: models.Value{Value: "foo", Kind: reflect.String}, Opts: options.FieldOptions{Display: true}, Kind: reflect.String},
					{Name: "Bar", Value: models.Value{Value: 1, Kind: reflect.Int}, Opts: options.FieldOptions{Display: true}, Kind: reflect.Int},
				},
			}
			got := f.Struct(v)
			exp := `Foo{Foo: foo, Bar: 1}`
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
			v := models.Struct{
				Name: "Foo",
				Fields: []models.Field{
					{Name: "Foo", Value: models.Value{Value: "testuser@exxxample.com", Kind: reflect.String}, Opts: options.FieldOptions{Display: true}, Kind: reflect.String},
					{Name: "Bar", Value: models.Value{Value: 1, Kind: reflect.Int}, Opts: options.FieldOptions{Display: false}, Kind: reflect.Int},
				},
			}
			got := f.Struct(v)
			exp := `{Foo: [******], Bar: [******]}`
			require.Equal(t, exp, got)
		})
	})
}
