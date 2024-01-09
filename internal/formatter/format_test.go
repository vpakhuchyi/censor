package formatter

import (
	"reflect"
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/vpakhuchyi/censor/config"
	"github.com/vpakhuchyi/censor/internal/models"
	"github.com/vpakhuchyi/censor/internal/options"
)

func TestFormatter_writeValue(t *testing.T) {
	f := Formatter{
		maskValue:               config.DefaultMaskValue,
		displayPointerSymbol:    false,
		displayStructName:       false,
		displayMapType:          false,
		excludePatterns:         nil,
		excludePatternsCompiled: nil,
	}
	var buf strings.Builder

	t.Run("string", func(t *testing.T) {
		t.Cleanup(func() { buf.Reset() })
		require.NotPanics(t, func() {
			v := models.Value{Value: "Kholodetsʹ", Kind: reflect.String}
			f.writeValue(&buf, v)
			exp := "Kholodetsʹ"
			require.Equal(t, exp, buf.String())
		})
	})

	t.Run("int", func(t *testing.T) {
		t.Cleanup(func() { buf.Reset() })
		require.NotPanics(t, func() {
			v := models.Value{Value: 44, Kind: reflect.Int}
			f.writeValue(&buf, v)
			exp := "44"
			require.Equal(t, exp, buf.String())
		})
	})

	t.Run("int8", func(t *testing.T) {
		t.Cleanup(func() { buf.Reset() })
		require.NotPanics(t, func() {
			v := models.Value{Value: int8(44), Kind: reflect.Int8}
			f.writeValue(&buf, v)
			exp := "44"
			require.Equal(t, exp, buf.String())
		})
	})

	t.Run("int16", func(t *testing.T) {
		t.Cleanup(func() { buf.Reset() })
		require.NotPanics(t, func() {
			v := models.Value{Value: int16(44), Kind: reflect.Int16}
			f.writeValue(&buf, v)
			exp := "44"
			require.Equal(t, exp, buf.String())
		})
	})

	t.Run("int32", func(t *testing.T) {
		t.Cleanup(func() { buf.Reset() })
		require.NotPanics(t, func() {
			v := models.Value{Value: int32(44), Kind: reflect.Int32}
			f.writeValue(&buf, v)
			exp := "44"
			require.Equal(t, exp, buf.String())
		})
	})

	t.Run("int64", func(t *testing.T) {
		t.Cleanup(func() { buf.Reset() })
		require.NotPanics(t, func() {
			v := models.Value{Value: int64(44), Kind: reflect.Int64}
			f.writeValue(&buf, v)
			exp := "44"
			require.Equal(t, exp, buf.String())
		})
	})

	t.Run("uint", func(t *testing.T) {
		t.Cleanup(func() { buf.Reset() })
		require.NotPanics(t, func() {
			v := models.Value{Value: uint(44), Kind: reflect.Uint}
			f.writeValue(&buf, v)
			exp := "44"
			require.Equal(t, exp, buf.String())
		})
	})

	t.Run("uint8", func(t *testing.T) {
		t.Cleanup(func() { buf.Reset() })
		require.NotPanics(t, func() {
			v := models.Value{Value: uint8(44), Kind: reflect.Uint8}
			f.writeValue(&buf, v)
			exp := "44"
			require.Equal(t, exp, buf.String())
		})
	})

	t.Run("uint16", func(t *testing.T) {
		t.Cleanup(func() { buf.Reset() })
		require.NotPanics(t, func() {
			v := models.Value{Value: uint16(44), Kind: reflect.Uint16}
			f.writeValue(&buf, v)
			exp := "44"
			require.Equal(t, exp, buf.String())
		})
	})

	t.Run("uint32", func(t *testing.T) {
		t.Cleanup(func() { buf.Reset() })
		require.NotPanics(t, func() {
			v := models.Value{Value: uint32(44), Kind: reflect.Uint32}
			f.writeValue(&buf, v)
			exp := "44"
			require.Equal(t, exp, buf.String())
		})
	})

	t.Run("uint64", func(t *testing.T) {
		t.Cleanup(func() { buf.Reset() })
		require.NotPanics(t, func() {
			v := models.Value{Value: uint64(44), Kind: reflect.Uint64}
			f.writeValue(&buf, v)
			exp := "44"
			require.Equal(t, exp, buf.String())
		})
	})

	t.Run("byte", func(t *testing.T) {
		t.Cleanup(func() { buf.Reset() })
		require.NotPanics(t, func() {
			var b byte = 44
			v := models.Value{Value: b, Kind: reflect.Uint8}
			f.writeValue(&buf, v)
			exp := "44"
			require.Equal(t, exp, buf.String())
		})
	})

	t.Run("rune", func(t *testing.T) {
		t.Cleanup(func() { buf.Reset() })
		require.NotPanics(t, func() {
			var r rune = 'A'
			v := models.Value{Value: r, Kind: reflect.Int32}
			f.writeValue(&buf, v)
			exp := "65"
			require.Equal(t, exp, buf.String())
		})
	})

	t.Run("float32", func(t *testing.T) {
		t.Cleanup(func() { buf.Reset() })
		require.NotPanics(t, func() {
			var fl float32 = 3.11111111111111
			v := models.Value{Value: fl, Kind: reflect.Float32}
			f.writeValue(&buf, v)
			exp := "3.1111112"
			require.Equal(t, exp, buf.String())
		})
	})

	t.Run("float64", func(t *testing.T) {
		t.Cleanup(func() { buf.Reset() })
		require.NotPanics(t, func() {
			var fl float64 = 3.11111111111111
			v := models.Value{Value: fl, Kind: reflect.Float64}
			f.writeValue(&buf, v)
			exp := "3.11111111111111"
			require.Equal(t, exp, buf.String())
		})
	})

	t.Run("complex64", func(t *testing.T) {
		t.Cleanup(func() { buf.Reset() })
		require.NotPanics(t, func() {
			var c complex64 = 3.11111111111111 + 3.11111111111111i
			v := models.Value{Value: c, Kind: reflect.Complex64}
			f.writeValue(&buf, v)
			exp := "(3.111111+3.111111i)"
			require.Equal(t, exp, buf.String())
		})
	})

	t.Run("complex128", func(t *testing.T) {
		t.Cleanup(func() { buf.Reset() })
		require.NotPanics(t, func() {
			var c complex128 = 3.11111111111111 + 3.11111111111111i
			v := models.Value{Value: c, Kind: reflect.Complex128}
			f.writeValue(&buf, v)
			exp := "(3.11111111111111+3.11111111111111i)"
			require.Equal(t, exp, buf.String())
		})
	})

	t.Run("bool", func(t *testing.T) {
		t.Cleanup(func() { buf.Reset() })
		require.NotPanics(t, func() {
			v := models.Value{Value: true, Kind: reflect.Bool}
			f.writeValue(&buf, v)
			exp := "true"
			require.Equal(t, exp, buf.String())
		})
	})

	t.Run("interface", func(t *testing.T) {
		t.Cleanup(func() { buf.Reset() })
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
			f.writeValue(&buf, v)
			exp := `[Kholodetsʹ, Halushky]`
			require.Equal(t, exp, buf.String())
		})
	})

	t.Run("slice", func(t *testing.T) {
		t.Cleanup(func() { buf.Reset() })
		require.NotPanics(t, func() {
			v := models.Value{
				Value: models.Slice{
					Values: []models.Value{
						{Value: "Kholodetsʹ", Kind: reflect.String},
						{Value: "Halushky", Kind: reflect.String},
					},
				},
				Kind: reflect.Slice,
			}
			f.writeValue(&buf, v)
			exp := `[Kholodetsʹ, Halushky]`
			require.Equal(t, exp, buf.String())
		})
	})

	t.Run("array", func(t *testing.T) {
		t.Cleanup(func() { buf.Reset() })
		require.NotPanics(t, func() {
			v := models.Value{
				Value: models.Slice{
					Values: []models.Value{
						{Value: "Kholodetsʹ", Kind: reflect.String},
						{Value: "Halushky", Kind: reflect.String},
					},
				},
				Kind: reflect.Array,
			}
			f.writeValue(&buf, v)
			exp := `[Kholodetsʹ, Halushky]`
			require.Equal(t, exp, buf.String())
		})
	})

	t.Run("struct", func(t *testing.T) {
		t.Cleanup(func() { buf.Reset() })
		require.NotPanics(t, func() {
			v := models.Value{
				Value: models.Struct{
					Fields: []models.Field{
						{
							Name: "Name",
							Value: models.Value{
								Value: "Kholodetsʹ",
								Kind:  reflect.String,
							},
							Kind: reflect.String,
						},
						{
							Name: "Ingredients",
							Value: models.Value{
								Value: models.Slice{
									Values: []models.Value{
										{Value: "Pork", Kind: reflect.String},
										{Value: "Garlic", Kind: reflect.String},
										{Value: "Black pepper", Kind: reflect.String},
										{Value: "Bay leaf", Kind: reflect.String},
										{Value: "Salt", Kind: reflect.String},
									},
								},
								Kind: reflect.Slice,
							},
							Kind: reflect.Slice,
							Opts: options.FieldOptions{Display: true},
						},
					},
				},
				Kind: reflect.Struct,
			}
			f.writeValue(&buf, v)
			exp := `{Name: [CENSORED], Ingredients: [Pork, Garlic, Black pepper, Bay leaf, Salt]}`
			require.Equal(t, exp, buf.String())
		})
	})

	t.Run("map", func(t *testing.T) {
		t.Cleanup(func() { buf.Reset() })
		require.NotPanics(t, func() {
			v := models.Value{
				Value: models.Map{
					Values: []models.KV{
						{
							Key: models.Value{
								Value: "Best dish ever",
								Kind:  reflect.String,
							},
							Value: models.Value{
								Value: "Kholodetsʹ",
								Kind:  reflect.String,
							},
						},
					},
					Type: "map[string]string",
				},
				Kind: reflect.Map,
			}
			f.writeValue(&buf, v)
			exp := `map[Best dish ever: Kholodetsʹ]`
			require.Equal(t, exp, buf.String())
		})
	})

	t.Run("pointer", func(t *testing.T) {
		t.Cleanup(func() { buf.Reset() })
		require.NotPanics(t, func() {
			var s = "Kholodetsʹ"
			v := models.Value{
				Value: models.Ptr{
					Value: s,
					Kind:  reflect.String,
				},
				Kind: reflect.Ptr,
			}
			f.writeValue(&buf, v)
			exp := `Kholodetsʹ`
			require.Equal(t, exp, buf.String())
		})
	})
}

func TestFormatter_writeField(t *testing.T) {
	f := Formatter{
		maskValue:         config.DefaultMaskValue,
		displayStructName: false,
		displayMapType:    false,
	}
	var buf strings.Builder

	t.Run("string", func(t *testing.T) {
		t.Cleanup(func() { buf.Reset() })
		require.NotPanics(t, func() {
			field := models.Field{
				Name: "Name",
				Value: models.Value{
					Value: "Kholodetsʹ",
					Kind:  reflect.String,
				},
				Kind: reflect.String,
			}
			f.writeField(field, &buf)
			exp := `Name: Kholodetsʹ`
			require.Equal(t, exp, buf.String())
		})
	})

	t.Run("int", func(t *testing.T) {
		t.Cleanup(func() { buf.Reset() })
		require.NotPanics(t, func() {
			field := models.Field{
				Name: "Age",
				Value: models.Value{
					Value: 44,
					Kind:  reflect.Int,
				},
				Kind: reflect.Int,
			}
			f.writeField(field, &buf)
			exp := `Age: 44`
			require.Equal(t, exp, buf.String())
		})
	})

	t.Run("int8", func(t *testing.T) {
		t.Cleanup(func() { buf.Reset() })
		require.NotPanics(t, func() {
			field := models.Field{
				Name: "Age",
				Value: models.Value{
					Value: int8(44),
					Kind:  reflect.Int8,
				},
				Kind: reflect.Int8,
			}
			f.writeField(field, &buf)
			exp := `Age: 44`
			require.Equal(t, exp, buf.String())
		})
	})

	t.Run("int16", func(t *testing.T) {
		t.Cleanup(func() { buf.Reset() })
		require.NotPanics(t, func() {
			field := models.Field{
				Name: "Age",
				Value: models.Value{
					Value: int16(44),
					Kind:  reflect.Int16,
				},
				Kind: reflect.Int16,
			}
			f.writeField(field, &buf)
			exp := `Age: 44`
			require.Equal(t, exp, buf.String())
		})
	})

	t.Run("int32", func(t *testing.T) {
		t.Cleanup(func() { buf.Reset() })
		require.NotPanics(t, func() {
			field := models.Field{
				Name: "Age",
				Value: models.Value{
					Value: int32(44),
					Kind:  reflect.Int32,
				},
				Kind: reflect.Int32,
			}
			f.writeField(field, &buf)
			exp := `Age: 44`
			require.Equal(t, exp, buf.String())
		})
	})

	t.Run("int64", func(t *testing.T) {
		t.Cleanup(func() { buf.Reset() })
		require.NotPanics(t, func() {
			field := models.Field{
				Name: "Age",
				Value: models.Value{
					Value: int64(44),
					Kind:  reflect.Int64,
				},
				Kind: reflect.Int64,
			}
			f.writeField(field, &buf)
			exp := `Age: 44`
			require.Equal(t, exp, buf.String())
		})
	})

	t.Run("uint", func(t *testing.T) {
		t.Cleanup(func() { buf.Reset() })
		require.NotPanics(t, func() {
			field := models.Field{
				Name: "Age",
				Value: models.Value{
					Value: uint(44),
					Kind:  reflect.Uint,
				},
				Kind: reflect.Uint,
			}
			f.writeField(field, &buf)
			exp := `Age: 44`
			require.Equal(t, exp, buf.String())
		})
	})

	t.Run("uint8", func(t *testing.T) {
		t.Cleanup(func() { buf.Reset() })
		require.NotPanics(t, func() {
			var b byte = 44
			field := models.Field{
				Name: "Age",
				Value: models.Value{
					Value: b,
					Kind:  reflect.Uint8,
				},
				Kind: reflect.Uint8,
			}
			f.writeField(field, &buf)
			exp := `Age: 44`
			require.Equal(t, exp, buf.String())
		})
	})

	t.Run("uint16", func(t *testing.T) {
		t.Cleanup(func() { buf.Reset() })
		require.NotPanics(t, func() {
			var b uint16 = 44
			field := models.Field{
				Name: "Age",
				Value: models.Value{
					Value: b,
					Kind:  reflect.Uint16,
				},
				Kind: reflect.Uint16,
			}
			f.writeField(field, &buf)
			exp := `Age: 44`
			require.Equal(t, exp, buf.String())
		})
	})

	t.Run("uint32", func(t *testing.T) {
		t.Cleanup(func() { buf.Reset() })
		require.NotPanics(t, func() {
			var b uint32 = 44
			field := models.Field{
				Name: "Age",
				Value: models.Value{
					Value: b,
					Kind:  reflect.Uint32,
				},
				Kind: reflect.Uint32,
			}
			f.writeField(field, &buf)
			exp := `Age: 44`
			require.Equal(t, exp, buf.String())
		})
	})

	t.Run("uint64", func(t *testing.T) {
		t.Cleanup(func() { buf.Reset() })
		require.NotPanics(t, func() {
			var b uint64 = 44
			field := models.Field{
				Name: "Age",
				Value: models.Value{
					Value: b,
					Kind:  reflect.Uint64,
				},
				Kind: reflect.Uint64,
			}
			f.writeField(field, &buf)
			exp := `Age: 44`
			require.Equal(t, exp, buf.String())
		})
	})

	t.Run("rune", func(t *testing.T) {
		t.Cleanup(func() { buf.Reset() })
		require.NotPanics(t, func() {
			var r rune = 'A'
			field := models.Field{
				Name: "Letter",
				Value: models.Value{
					Value: r,
					Kind:  reflect.Int32,
				},
				Kind: reflect.Int32,
			}
			f.writeField(field, &buf)
			exp := `Letter: 65`
			require.Equal(t, exp, buf.String())
		})
	})

	t.Run("byte", func(t *testing.T) {
		t.Cleanup(func() { buf.Reset() })
		require.NotPanics(t, func() {
			var b byte = 44
			field := models.Field{
				Name: "Age",
				Value: models.Value{
					Value: b,
					Kind:  reflect.Uint8,
				},
				Kind: reflect.Uint8,
			}
			f.writeField(field, &buf)
			exp := `Age: 44`
			require.Equal(t, exp, buf.String())
		})
	})

	t.Run("float32", func(t *testing.T) {
		t.Cleanup(func() { buf.Reset() })
		require.NotPanics(t, func() {
			var fl float32 = 3.11111111111111
			field := models.Field{
				Name: "Float",
				Value: models.Value{
					Value: fl,
					Kind:  reflect.Float32,
				},
				Kind: reflect.Float32,
			}
			f.writeField(field, &buf)
			exp := `Float: 3.1111112`
			require.Equal(t, exp, buf.String())
		})
	})

	t.Run("float64", func(t *testing.T) {
		t.Cleanup(func() { buf.Reset() })
		require.NotPanics(t, func() {
			var fl float64 = 3.11111111111111
			field := models.Field{
				Name: "Float",
				Value: models.Value{
					Value: fl,
					Kind:  reflect.Float64,
				},
				Kind: reflect.Float64,
			}
			f.writeField(field, &buf)
			exp := `Float: 3.11111111111111`
			require.Equal(t, exp, buf.String())
		})
	})

	t.Run("complex64", func(t *testing.T) {
		t.Cleanup(func() { buf.Reset() })
		require.NotPanics(t, func() {
			var c complex64 = 3.11111111111111 + 3.11111111111111i
			field := models.Field{
				Name: "Test",
				Value: models.Value{
					Value: c,
					Kind:  reflect.Complex64,
				},
				Kind: reflect.Complex64,
			}
			f.writeField(field, &buf)
			exp := `Test: (3.111111+3.111111i)`
			require.Equal(t, exp, buf.String())
		})
	})

	t.Run("complex128", func(t *testing.T) {
		t.Cleanup(func() { buf.Reset() })
		require.NotPanics(t, func() {
			var c complex128 = 3.11111111111111 + 3.11111111111111i
			field := models.Field{
				Name: "Test",
				Value: models.Value{
					Value: c,
					Kind:  reflect.Complex128,
				},
				Kind: reflect.Complex128,
			}
			f.writeField(field, &buf)
			exp := `Test: (3.11111111111111+3.11111111111111i)`
			require.Equal(t, exp, buf.String())
		})
	})

	t.Run("bool", func(t *testing.T) {
		t.Cleanup(func() { buf.Reset() })
		require.NotPanics(t, func() {
			field := models.Field{
				Name: "IsAwesome",
				Value: models.Value{
					Value: true,
					Kind:  reflect.Bool,
				},
				Kind: reflect.Bool,
			}
			f.writeField(field, &buf)
			exp := `IsAwesome: true`
			require.Equal(t, exp, buf.String())
		})
	})

	t.Run("interface", func(t *testing.T) {
		t.Cleanup(func() { buf.Reset() })
		require.NotPanics(t, func() {
			field := models.Field{
				Name: "Dishes",
				Value: models.Value{
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
				},
				Kind: reflect.Interface,
			}
			f.writeField(field, &buf)
			exp := `Dishes: [Kholodetsʹ, Halushky]`
			require.Equal(t, exp, buf.String())
		})
	})

	t.Run("slice", func(t *testing.T) {
		t.Cleanup(func() { buf.Reset() })
		require.NotPanics(t, func() {
			field := models.Field{
				Name: "Dishes",
				Value: models.Value{
					Value: models.Slice{
						Values: []models.Value{
							{Value: "Kholodetsʹ", Kind: reflect.String},
							{Value: "Halushky", Kind: reflect.String},
						},
					},
					Kind: reflect.Slice,
				},
				Kind: reflect.Slice,
			}
			f.writeField(field, &buf)
			exp := `Dishes: [Kholodetsʹ, Halushky]`
			require.Equal(t, exp, buf.String())
		})
	})

	t.Run("array", func(t *testing.T) {
		t.Cleanup(func() { buf.Reset() })
		require.NotPanics(t, func() {
			field := models.Field{
				Name: "Dishes",
				Value: models.Value{
					Value: models.Slice{
						Values: []models.Value{
							{Value: "Kholodetsʹ", Kind: reflect.String},
							{Value: "Halushky", Kind: reflect.String},
						},
					},
					Kind: reflect.Slice,
				},
				Kind: reflect.Array,
			}
			f.writeField(field, &buf)
			exp := `Dishes: [Kholodetsʹ, Halushky]`
			require.Equal(t, exp, buf.String())
		})
	})

	t.Run("struct", func(t *testing.T) {
		t.Cleanup(func() { buf.Reset() })
		require.NotPanics(t, func() {
			field := models.Field{
				Name: "Dish",
				Value: models.Value{
					Value: models.Struct{
						Fields: []models.Field{
							{
								Name: "Name",
								Value: models.Value{
									Value: "Kholodetsʹ",
									Kind:  reflect.String,
								},
								Kind: reflect.String,
							},
							{
								Name: "Ingredients",
								Value: models.Value{
									Value: models.Slice{
										Values: []models.Value{
											{Value: "Pork", Kind: reflect.String},
											{Value: "Garlic", Kind: reflect.String},
											{Value: "Black pepper", Kind: reflect.String},
											{Value: "Bay leaf", Kind: reflect.String},
											{Value: "Salt", Kind: reflect.String},
										},
									},
									Kind: reflect.Slice,
								},
								Kind: reflect.Slice,
								Opts: options.FieldOptions{Display: true},
							},
						},
					},
					Kind: reflect.Struct,
				},
				Kind: reflect.Struct,
			}
			f.writeField(field, &buf)
			exp := `Dish: {Name: [CENSORED], Ingredients: [Pork, Garlic, Black pepper, Bay leaf, Salt]}`
			require.Equal(t, exp, buf.String())
		})
	})

	t.Run("map", func(t *testing.T) {
		t.Cleanup(func() { buf.Reset() })
		require.NotPanics(t, func() {
			field := models.Field{
				Name: "Dish",
				Value: models.Value{
					Value: models.Map{
						Values: []models.KV{
							{
								Key: models.Value{
									Value: "Best dish ever",
									Kind:  reflect.String,
								},
								Value: models.Value{
									Value: "Kholodetsʹ",
									Kind:  reflect.String,
								},
							},
						},
						Type: "map[string]string",
					},
					Kind: reflect.Map,
				},
				Kind: reflect.Map,
			}
			f.writeField(field, &buf)
			exp := `Dish: map[Best dish ever: Kholodetsʹ]`
			require.Equal(t, exp, buf.String())
		})
	})

	t.Run("pointer", func(t *testing.T) {
		t.Cleanup(func() { buf.Reset() })
		require.NotPanics(t, func() {
			var s = "Kholodetsʹ"
			field := models.Field{
				Name: "Dish",
				Value: models.Value{
					Value: models.Ptr{
						Value: s,
						Kind:  reflect.String,
					},
					Kind: reflect.Ptr,
				},
				Kind: reflect.Ptr,
			}
			f.writeField(field, &buf)
			exp := `Dish: Kholodetsʹ`
			require.Equal(t, exp, buf.String())
		})
	})
}

func TestNew(t *testing.T) {
	t.Run("with_exclude_patterns", func(t *testing.T) {
		got := New(config.Formatter{
			MaskValue:         "[censored]",
			DisplayStructName: true,
			DisplayMapType:    true,
			ExcludePatterns:   []string{`\d`},
		})
		exp := &Formatter{
			maskValue:               "[censored]",
			displayStructName:       true,
			displayMapType:          true,
			excludePatterns:         []string{`\d`},
			excludePatternsCompiled: []*regexp.Regexp{regexp.MustCompile(`\d`)},
		}
		require.EqualValues(t, exp, got)
	})

	t.Run("without_exclude_patterns", func(t *testing.T) {
		got := New(config.Formatter{
			MaskValue:         "[censored]",
			DisplayStructName: true,
			DisplayMapType:    true,
			ExcludePatterns:   nil,
		})
		exp := &Formatter{
			maskValue:               "[censored]",
			displayStructName:       true,
			displayMapType:          true,
			excludePatterns:         nil,
			excludePatternsCompiled: nil,
		}
		require.EqualValues(t, exp, got)
	})

	t.Run("with_empty_exclude_patterns", func(t *testing.T) {
		got := New(config.Formatter{
			MaskValue:         "[censored]",
			DisplayStructName: true,
			DisplayMapType:    true,
			ExcludePatterns:   []string{},
		})
		exp := &Formatter{
			maskValue:               "[censored]",
			displayStructName:       true,
			displayMapType:          true,
			excludePatterns:         []string{},
			excludePatternsCompiled: nil,
		}
		require.EqualValues(t, exp, got)
	})
}
