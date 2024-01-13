package formatter

import (
	"reflect"
	"regexp"
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
					{Name: "Foo", Value: models.Value{Value: "foo", Kind: reflect.String}, Opts: options.FieldOptions{Display: true}},
					{Name: "Bar", Value: models.Value{Value: 1, Kind: reflect.Int}, Opts: options.FieldOptions{Display: false}},
				},
			}
			got := f.Struct(v)
			exp := `{Foo: foo, Bar: [CENSORED]}`
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
					{Name: "Foo", Value: models.Value{Value: "foo", Kind: reflect.String}, Opts: options.FieldOptions{Display: true}},
					{Name: "Bar", Value: models.Value{Value: 1, Kind: reflect.Int}, Opts: options.FieldOptions{Display: true}},
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
			excludePatternsCompiled: []*regexp.Regexp{compiledRegExpDigit, compiledRegExpEmail},
		}

		require.NotPanics(t, func() {
			v := models.Struct{
				Name: "Foo",
				Fields: []models.Field{
					{Name: "Foo", Value: models.Value{Value: "testuser@exxxample.com", Kind: reflect.String}, Opts: options.FieldOptions{Display: true}},
					{Name: "Bar", Value: models.Value{Value: 1, Kind: reflect.Int}, Opts: options.FieldOptions{Display: false}},
				},
			}
			got := f.Struct(v)
			exp := `{Foo: [CENSORED], Bar: [CENSORED]}`
			require.Equal(t, exp, got)
		})
	})

	t.Run("with_unsupported_types", func(t *testing.T) {
		f := Formatter{
			maskValue:               config.DefaultMaskValue,
			displayStructName:       false,
			displayMapType:          false,
			excludePatterns:         []string{},
			excludePatternsCompiled: nil,
		}

		require.NotPanics(t, func() {
			v := models.Struct{
				Name: "parser.structWithUnsupportedTypes",
				Fields: []models.Field{
					{Name: "ChanWithCensorTag", Value: models.Value{Value: "[Unsupported type: chan]", Kind: reflect.Chan}, Opts: options.FieldOptions{Display: true}},
					{Name: "Chan", Value: models.Value{Value: "[Unsupported type: chan]", Kind: reflect.Chan}, Opts: options.FieldOptions{Display: false}},
					{Name: "FuncWithCensorTag", Value: models.Value{Value: "[Unsupported type: func]", Kind: reflect.Func}, Opts: options.FieldOptions{Display: true}},
					{Name: "Func", Value: models.Value{Value: "[Unsupported type: func]", Kind: reflect.Func}, Opts: options.FieldOptions{Display: false}},
					{Name: "UnsafeWithCensorTag", Value: models.Value{Value: "[Unsupported type: unsafe.Pointer]", Kind: reflect.UnsafePointer}, Opts: options.FieldOptions{Display: true}},
					{Name: "Unsafe", Value: models.Value{Value: "[Unsupported type: unsafe.Pointer]", Kind: reflect.UnsafePointer}, Opts: options.FieldOptions{Display: false}},
				},
			}
			got := f.Struct(v)
			exp := `{ChanWithCensorTag: [Unsupported type: chan], Chan: [CENSORED], FuncWithCensorTag: [Unsupported type: func], Func: [CENSORED], UnsafeWithCensorTag: [Unsupported type: unsafe.Pointer], Unsafe: [CENSORED]}`
			require.Equal(t, exp, got)
		})
	})

}
