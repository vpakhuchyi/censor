package censor

import (
	"reflect"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/require"

	"github.com/vpakhuchyi/censor/config"
	"github.com/vpakhuchyi/censor/internal/formatter"
	"github.com/vpakhuchyi/censor/internal/models"
	"github.com/vpakhuchyi/censor/internal/options"
	"github.com/vpakhuchyi/censor/internal/parser"
)

func TestProcessor_parse(t *testing.T) {
	p := New()

	t.Run("struct", func(t *testing.T) {
		require.NotPanics(t, func() {
			type s struct {
				Int64  int64  `censor:"display"`
				String string `censor:"display"`
			}

			val := s{Int64: 123456789, String: "string"}
			exp := models.Struct{
				Name: "censor.s",
				Fields: []models.Field{
					{
						Name:  "Int64",
						Value: models.Value{Value: int64(123456789), Kind: reflect.Int64},
						Opts:  options.FieldOptions{Display: true},
						Kind:  reflect.Int64,
					},
					{
						Name:  "String",
						Value: models.Value{Value: "string", Kind: reflect.String},
						Opts:  options.FieldOptions{Display: true},
						Kind:  reflect.String,
					},
				},
			}

			got := p.parse(reflect.ValueOf(val))
			require.Equal(t, exp, got)
		})
	})

	t.Run("pointer_to_struct", func(t *testing.T) {
		require.NotPanics(t, func() {
			type s struct {
				Int64  int64  `censor:"display"`
				String string `censor:"display"`
			}

			val := &s{Int64: 123456789, String: "string"}
			exp := models.Ptr{
				Value: models.Struct{
					Name: "censor.s",
					Fields: []models.Field{
						{
							Name:  "Int64",
							Value: models.Value{Value: int64(123456789), Kind: reflect.Int64},
							Opts:  options.FieldOptions{Display: true},
							Kind:  reflect.Int64,
						},
						{
							Name:  "String",
							Value: models.Value{Value: "string", Kind: reflect.String},
							Opts:  options.FieldOptions{Display: true},
							Kind:  reflect.String,
						},
					},
				},
				Kind: reflect.Struct,
			}

			got := p.parse(reflect.ValueOf(val))
			require.Equal(t, exp, got)
		})
	})

	t.Run("slice_of_strings", func(t *testing.T) {
		require.NotPanics(t, func() {
			val := []string{"salo", "hlib"}
			exp := models.Slice{
				Values: []models.Value{
					{Value: "salo", Kind: reflect.String},
					{Value: "hlib", Kind: reflect.String},
				},
			}

			got := p.parse(reflect.ValueOf(val))
			require.Equal(t, exp, got)
		})
	})

	t.Run("array_of_strings", func(t *testing.T) {
		require.NotPanics(t, func() {
			val := [2]string{"salo", "hlib"}
			exp := models.Slice{
				Values: []models.Value{
					{Value: "salo", Kind: reflect.String},
					{Value: "hlib", Kind: reflect.String},
				},
			}

			got := p.parse(reflect.ValueOf(val))
			require.Equal(t, exp, got)
		})
	})

	t.Run("map_of_strings", func(t *testing.T) {
		require.NotPanics(t, func() {
			val := map[string]int{"odyn": 1, "p'yatʹ": 5}
			exp := models.Map{
				Values: []models.KV{
					{
						SortValue: "odyn",
						Key:       models.Value{Value: "odyn", Kind: reflect.String},
						Value:     models.Value{Value: 1, Kind: reflect.Int},
					},
					{
						SortValue: "p'yatʹ",
						Key:       models.Value{Value: "p'yatʹ", Kind: reflect.String},
						Value:     models.Value{Value: 5, Kind: reflect.Int},
					},
				},
				Type: "map[string]int",
			}

			got := p.parse(reflect.ValueOf(val))
			require.Equal(t, exp, got)
		})
	})

	t.Run("complex64", func(t *testing.T) {
		require.NotPanics(t, func() {
			val := complex(float32(-45.234), float32(11.933))
			exp := models.Value{Value: complex(float32(-45.234), float32(11.933)), Kind: reflect.Complex64}

			got := p.parse(reflect.ValueOf(val))
			require.Equal(t, exp, got)
		})
	})

	t.Run("complex128", func(t *testing.T) {
		require.NotPanics(t, func() {
			val := complex(float64(-445.2366664), float64(121.93767763))
			exp := models.Value{Value: complex(float64(-445.2366664), float64(121.93767763)), Kind: reflect.Complex128}

			got := p.parse(reflect.ValueOf(val))
			require.Equal(t, exp, got)
		})
	})

	t.Run("float32", func(t *testing.T) {
		require.NotPanics(t, func() {
			val := float32(-45.234444444)
			exp := models.Value{Value: float32(-45.234444), Kind: reflect.Float32}

			got := p.parse(reflect.ValueOf(val))
			require.Equal(t, exp, got)
		})
	})

	t.Run("float64", func(t *testing.T) {
		require.NotPanics(t, func() {
			val := float64(-445.236666433333333333)
			exp := models.Value{Value: float64(-445.2366664333333), Kind: reflect.Float64}

			got := p.parse(reflect.ValueOf(val))
			require.Equal(t, exp, got)
		})
	})

	t.Run("int", func(t *testing.T) {
		require.NotPanics(t, func() {
			val := int(-445)
			exp := models.Value{Value: int(-445), Kind: reflect.Int}

			got := p.parse(reflect.ValueOf(val))
			require.Equal(t, exp, got)
		})
	})

	t.Run("int8", func(t *testing.T) {
		require.NotPanics(t, func() {
			val := int8(-45)
			exp := models.Value{Value: int8(-45), Kind: reflect.Int8}

			got := p.parse(reflect.ValueOf(val))
			require.Equal(t, exp, got)
		})
	})

	t.Run("int16", func(t *testing.T) {
		require.NotPanics(t, func() {
			val := int16(-445)
			exp := models.Value{Value: int16(-445), Kind: reflect.Int16}

			got := p.parse(reflect.ValueOf(val))
			require.Equal(t, exp, got)
		})
	})

	t.Run("int32", func(t *testing.T) {
		require.NotPanics(t, func() {
			val := int32(-445)
			exp := models.Value{Value: int32(-445), Kind: reflect.Int32}

			got := p.parse(reflect.ValueOf(val))
			require.Equal(t, exp, got)
		})
	})

	t.Run("int64", func(t *testing.T) {
		require.NotPanics(t, func() {
			val := int64(-445)
			exp := models.Value{Value: int64(-445), Kind: reflect.Int64}

			got := p.parse(reflect.ValueOf(val))
			require.Equal(t, exp, got)
		})
	})

	t.Run("uint", func(t *testing.T) {
		require.NotPanics(t, func() {
			val := uint(445)
			exp := models.Value{Value: uint(445), Kind: reflect.Uint}

			got := p.parse(reflect.ValueOf(val))
			require.Equal(t, exp, got)
		})
	})

	t.Run("uint8", func(t *testing.T) {
		require.NotPanics(t, func() {
			val := uint8(45)
			exp := models.Value{Value: uint8(45), Kind: reflect.Uint8}

			got := p.parse(reflect.ValueOf(val))
			require.Equal(t, exp, got)
		})
	})

	t.Run("uint16", func(t *testing.T) {
		require.NotPanics(t, func() {
			val := uint16(445)
			exp := models.Value{Value: uint16(445), Kind: reflect.Uint16}

			got := p.parse(reflect.ValueOf(val))
			require.Equal(t, exp, got)
		})
	})

	t.Run("uint32", func(t *testing.T) {
		require.NotPanics(t, func() {
			val := uint32(445)
			exp := models.Value{Value: uint32(445), Kind: reflect.Uint32}

			got := p.parse(reflect.ValueOf(val))
			require.Equal(t, exp, got)
		})
	})

	t.Run("uint64", func(t *testing.T) {
		require.NotPanics(t, func() {
			val := uint64(445)
			exp := models.Value{Value: uint64(445), Kind: reflect.Uint64}

			got := p.parse(reflect.ValueOf(val))
			require.Equal(t, exp, got)
		})
	})

	t.Run("rune", func(t *testing.T) {
		require.NotPanics(t, func() {
			val := 'U'
			exp := models.Value{Value: rune(85), Kind: reflect.Int32}

			got := p.parse(reflect.ValueOf(val))
			require.Equal(t, exp, got)
		})
	})

	t.Run("byte", func(t *testing.T) {
		require.NotPanics(t, func() {
			val := byte(45)
			exp := models.Value{Value: byte(45), Kind: reflect.Uint8}

			got := p.parse(reflect.ValueOf(val))
			require.Equal(t, exp, got)
		})
	})

	t.Run("bool", func(t *testing.T) {
		require.NotPanics(t, func() {
			val := true
			exp := models.Value{Value: true, Kind: reflect.Bool}

			got := p.parse(reflect.ValueOf(val))
			require.Equal(t, exp, got)
		})
	})

	t.Run("string", func(t *testing.T) {
		require.NotPanics(t, func() {
			val := "string"
			exp := models.Value{Value: "string", Kind: reflect.String}

			got := p.parse(reflect.ValueOf(val))
			require.Equal(t, exp, got)
		})
	})

	t.Run("unsupported_type_chan", func(t *testing.T) {
		require.NotPanics(t, func() {
			val := make(chan int)
			exp := models.Value{Value: "", Kind: reflect.Chan}

			got := p.parse(reflect.ValueOf(val))
			require.Equal(t, exp, got)
		})
	})

	t.Run("unsupported_type_func", func(t *testing.T) {
		require.NotPanics(t, func() {
			val := func() {}
			exp := models.Value{Value: "", Kind: reflect.Func}

			got := p.parse(reflect.ValueOf(val))
			require.Equal(t, exp, got)
		})
	})

	t.Run("unsupported_type_unsafe_pointer", func(t *testing.T) {
		require.NotPanics(t, func() {
			val := unsafe.Pointer(nil)
			exp := models.Value{Value: "", Kind: reflect.UnsafePointer}

			got := p.parse(reflect.ValueOf(val))
			require.Equal(t, exp, got)
		})
	})

	t.Run("unsupported_type_uintptr", func(t *testing.T) {
		require.NotPanics(t, func() {
			val := uintptr(0)
			exp := models.Value{Value: "", Kind: reflect.Uintptr}

			got := p.parse(reflect.ValueOf(val))
			require.Equal(t, exp, got)
		})
	})
}

func TestProcessor_format(t *testing.T) {
	p := New()

	t.Run("struct", func(t *testing.T) {
		require.NotPanics(t, func() {
			val := models.Struct{
				Name: "censor.s",
				Fields: []models.Field{
					{
						Name:  "Int64",
						Value: models.Value{Value: int64(123456789), Kind: reflect.Int64},
						Opts:  options.FieldOptions{Display: true},
						Kind:  reflect.Int64,
					},
					{
						Name:  "String",
						Value: models.Value{Value: "string", Kind: reflect.String},
						Opts:  options.FieldOptions{Display: true},
						Kind:  reflect.String,
					},
				},
			}
			exp := `{Int64: 123456789, String: string}`

			got := p.format(reflect.Struct, val)
			require.Equal(t, exp, got)
		})
	})

	t.Run("pointer_to_struct", func(t *testing.T) {
		require.NotPanics(t, func() {
			val := models.Ptr{
				Value: models.Struct{
					Name: "censor.s",
					Fields: []models.Field{
						{
							Name:  "Int64",
							Value: models.Value{Value: int64(123456789), Kind: reflect.Int64},
							Opts:  options.FieldOptions{Display: true},
							Kind:  reflect.Int64,
						},
						{
							Name:  "String",
							Value: models.Value{Value: "string", Kind: reflect.String},
							Opts:  options.FieldOptions{Display: true},
							Kind:  reflect.String,
						},
					},
				},
				Kind: reflect.Struct,
			}
			exp := `{Int64: 123456789, String: string}`

			got := p.format(reflect.Pointer, val)
			require.Equal(t, exp, got)
		})
	})

	t.Run("slice_of_strings", func(t *testing.T) {
		require.NotPanics(t, func() {
			val := models.Slice{
				Values: []models.Value{
					{Value: "salo", Kind: reflect.String},
					{Value: "hlib", Kind: reflect.String},
				},
			}
			exp := `[salo, hlib]`

			got := p.format(reflect.Slice, val)
			require.Equal(t, exp, got)
		})
	})

	t.Run("array_of_strings", func(t *testing.T) {
		require.NotPanics(t, func() {
			val := models.Slice{
				Values: []models.Value{
					{Value: "salo", Kind: reflect.String},
					{Value: "hlib", Kind: reflect.String},
				},
			}
			exp := `[salo, hlib]`

			got := p.format(reflect.Array, val)
			require.Equal(t, exp, got)
		})
	})

	t.Run("map_of_strings", func(t *testing.T) {
		require.NotPanics(t, func() {
			val := models.Map{
				Values: []models.KV{
					{
						SortValue: "odyn",
						Key:       models.Value{Value: "odyn", Kind: reflect.String},
						Value:     models.Value{Value: 1, Kind: reflect.Int},
					},
					{
						SortValue: "p'yatʹ",
						Key:       models.Value{Value: "p'yatʹ", Kind: reflect.String},
						Value:     models.Value{Value: 5, Kind: reflect.Int},
					},
				},
				Type: "map[string]int",
			}
			exp := `map[odyn: 1, p'yatʹ: 5]`

			got := p.format(reflect.Map, val)
			require.Equal(t, exp, got)
		})
	})

	t.Run("complex64", func(t *testing.T) {
		require.NotPanics(t, func() {
			val := models.Value{Value: complex(float32(-45.234), float32(11.933)), Kind: reflect.Complex64}
			exp := `(-45.234+11.933i)`

			got := p.format(reflect.Complex64, val)
			require.Equal(t, exp, got)
		})
	})

	t.Run("complex128", func(t *testing.T) {
		require.NotPanics(t, func() {
			val := models.Value{Value: complex(float64(-445.2366664), float64(121.93767763)), Kind: reflect.Complex128}
			exp := `(-445.2366664+121.93767763i)`

			got := p.format(reflect.Complex128, val)
			require.Equal(t, exp, got)
		})
	})

	t.Run("float32", func(t *testing.T) {
		require.NotPanics(t, func() {
			val := models.Value{Value: float32(-45.23444444444444), Kind: reflect.Float32}
			exp := `-45.23444`

			got := p.format(reflect.Float32, val)
			require.Equal(t, exp, got)
		})
	})

	t.Run("float64", func(t *testing.T) {
		require.NotPanics(t, func() {
			val := models.Value{Value: float64(-445.236666433333333334), Kind: reflect.Float64}
			exp := `-445.236666433333`

			got := p.format(reflect.Float64, val)
			require.Equal(t, exp, got)
		})
	})

	t.Run("int", func(t *testing.T) {
		require.NotPanics(t, func() {
			val := models.Value{Value: int(-445), Kind: reflect.Int}
			exp := `-445`

			got := p.format(reflect.Int, val)
			require.Equal(t, exp, got)
		})
	})

	t.Run("int8", func(t *testing.T) {
		require.NotPanics(t, func() {
			val := models.Value{Value: int8(-45), Kind: reflect.Int8}
			exp := `-45`

			got := p.format(reflect.Int8, val)
			require.Equal(t, exp, got)
		})
	})

	t.Run("int16", func(t *testing.T) {
		require.NotPanics(t, func() {
			val := models.Value{Value: int16(-445), Kind: reflect.Int16}
			exp := `-445`

			got := p.format(reflect.Int16, val)
			require.Equal(t, exp, got)
		})
	})

	t.Run("int32", func(t *testing.T) {
		require.NotPanics(t, func() {
			val := models.Value{Value: int32(-445), Kind: reflect.Int32}
			exp := `-445`

			got := p.format(reflect.Int32, val)
			require.Equal(t, exp, got)
		})
	})

	t.Run("int64", func(t *testing.T) {
		require.NotPanics(t, func() {
			val := models.Value{Value: int64(-445), Kind: reflect.Int64}
			exp := `-445`

			got := p.format(reflect.Int64, val)
			require.Equal(t, exp, got)
		})
	})

	t.Run("uint", func(t *testing.T) {
		require.NotPanics(t, func() {
			val := models.Value{Value: uint(445), Kind: reflect.Uint}
			exp := `445`

			got := p.format(reflect.Uint, val)
			require.Equal(t, exp, got)
		})
	})

	t.Run("uint8", func(t *testing.T) {
		require.NotPanics(t, func() {
			val := models.Value{Value: uint8(45), Kind: reflect.Uint8}
			exp := `45`

			got := p.format(reflect.Uint8, val)
			require.Equal(t, exp, got)
		})
	})

	t.Run("uint16", func(t *testing.T) {
		require.NotPanics(t, func() {
			val := models.Value{Value: uint16(445), Kind: reflect.Uint16}
			exp := `445`

			got := p.format(reflect.Uint16, val)
			require.Equal(t, exp, got)
		})
	})

	t.Run("uint32", func(t *testing.T) {
		require.NotPanics(t, func() {
			val := models.Value{Value: uint32(445), Kind: reflect.Uint32}
			exp := `445`

			got := p.format(reflect.Uint32, val)
			require.Equal(t, exp, got)
		})
	})

	t.Run("uint64", func(t *testing.T) {
		require.NotPanics(t, func() {
			val := models.Value{Value: uint64(445), Kind: reflect.Uint64}
			exp := `445`

			got := p.format(reflect.Uint64, val)
			require.Equal(t, exp, got)
		})
	})

	t.Run("rune", func(t *testing.T) {
		require.NotPanics(t, func() {
			val := models.Value{Value: rune(85), Kind: reflect.Int32}
			exp := `85`

			got := p.format(reflect.Int32, val)
			require.Equal(t, exp, got)
		})
	})

	t.Run("byte", func(t *testing.T) {
		require.NotPanics(t, func() {
			val := models.Value{Value: byte(45), Kind: reflect.Uint8}
			exp := `45`

			got := p.format(reflect.Uint8, val)
			require.Equal(t, exp, got)
		})
	})

	t.Run("bool", func(t *testing.T) {
		require.NotPanics(t, func() {
			val := models.Value{Value: true, Kind: reflect.Bool}
			exp := `true`

			got := p.format(reflect.Bool, val)
			require.Equal(t, exp, got)
		})
	})

	t.Run("string", func(t *testing.T) {
		require.NotPanics(t, func() {
			val := models.Value{Value: "kotleta", Kind: reflect.String}
			exp := `kotleta`

			got := p.format(reflect.String, val)
			require.Equal(t, exp, got)
		})
	})

	t.Run("unsupported_type_chan", func(t *testing.T) {
		require.NotPanics(t, func() {
			val := models.Value{Value: "", Kind: reflect.Chan}
			exp := ``

			got := p.format(reflect.Chan, val)
			require.Equal(t, exp, got)
		})
	})

	t.Run("unsupported_type_func", func(t *testing.T) {
		require.NotPanics(t, func() {
			val := models.Value{Value: "", Kind: reflect.Func}
			exp := ``

			got := p.format(reflect.Func, val)
			require.Equal(t, exp, got)
		})
	})

	t.Run("unsupported_type_unsafe_pointer", func(t *testing.T) {
		require.NotPanics(t, func() {
			val := models.Value{Value: "", Kind: reflect.UnsafePointer}
			exp := ``

			got := p.format(reflect.UnsafePointer, val)
			require.Equal(t, exp, got)
		})
	})

	t.Run("unsupported_type_uintptr", func(t *testing.T) {
		require.NotPanics(t, func() {
			val := models.Value{Value: "", Kind: reflect.Uintptr}
			exp := ``

			got := p.format(reflect.Uintptr, val)
			require.Equal(t, exp, got)
		})
	})
}

func TestNewWithConfig(t *testing.T) {
	cfg := config.Config{
		Parser: config.Parser{
			UseJSONTagName: false,
		},
		Formatter: config.Formatter{
			MaskValue:            "####",
			DisplayPointerSymbol: false,
			DisplayStructName:    false,
			DisplayMapType:       false,
			ExcludePatterns:      nil,
		},
	}
	got := NewWithConfig(cfg)
	exp := &Processor{
		formatter: formatter.NewWithConfig(cfg.Formatter),
		parser:    parser.NewWithConfig(cfg.Parser),
	}

	require.Equal(t, exp, got)
}

func TestNewWithFileConfig(t *testing.T) {
	t.Run("successful", func(t *testing.T) {
		t.Cleanup(func() { SetGlobalInstance(New()) })

		cfg, err := config.FromFile("./config/testdata/cfg.yml")
		require.NoError(t, err)

		want := Processor{
			formatter: formatter.NewWithConfig(cfg.Formatter),
			parser:    parser.NewWithConfig(cfg.Parser),
		}

		p, err := NewWithFileConfig("./config/testdata/cfg.yml")
		require.NoError(t, err)
		require.EqualValues(t, want.formatter, p.formatter)
		require.EqualValues(t, want.parser, p.parser)
	})

	t.Run("empty_file_path", func(t *testing.T) {
		t.Cleanup(func() { SetGlobalInstance(New()) })

		var want *Processor
		p, err := NewWithFileConfig("")
		require.Error(t, err)
		require.Equal(t, want, p)
	})

	t.Run("invalid_file_content", func(t *testing.T) {
		t.Cleanup(func() { SetGlobalInstance(New()) })

		var want *Processor

		p, err := NewWithFileConfig("./config/testdata/invalid-cfg.yml")
		require.Error(t, err)
		require.Equal(t, want, p)
	})

	t.Run("empty_file_content", func(t *testing.T) {
		t.Cleanup(func() { SetGlobalInstance(New()) })

		want := &Processor{
			formatter: &formatter.Formatter{},
			parser:    parser.New(),
		}

		p, err := NewWithFileConfig("./config/testdata/empty.yml")
		require.NoError(t, err)
		require.EqualValues(t, want.formatter, p.formatter)
		require.EqualValues(t, want.parser, p.parser)
	})
}
