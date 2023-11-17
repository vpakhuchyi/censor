package parser

import (
	"os"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/vpakhuchyi/censor/internal/models"
	"github.com/vpakhuchyi/censor/internal/options"
)

var mainPkgStruct any

func TestParser_Struct(t *testing.T) {
	type address struct {
		City   string `json:"city" censor:"display"`
		State  string `json:"state" censor:"display"`
		Street string `json:"street"`
		Zip    string `json:"zip"`
	}

	p := Parser{
		useJSONTagName: false,
		censorFieldTag: DefaultCensorFieldTag,
	}

	t.Run("struct_with_strings", func(t *testing.T) {
		require.NotPanics(t, func() {
			v := address{Street: "451 Main St", City: "San Francisco", State: "CA", Zip: "55501"}
			got := p.Struct(reflect.ValueOf(v))
			exp := models.Struct{
				Name: "parser.address",
				Fields: []models.Field{
					{Name: "City", Value: models.Value{Value: "San Francisco", Kind: reflect.String}, Opts: options.FieldOptions{Display: true}, Kind: reflect.String},
					{Name: "State", Value: models.Value{Value: "CA", Kind: reflect.String}, Opts: options.FieldOptions{Display: true}, Kind: reflect.String},
					{Name: "Street", Value: models.Value{Value: "451 Main St", Kind: reflect.String}, Opts: options.FieldOptions{Display: false}, Kind: reflect.String},
					{Name: "Zip", Value: models.Value{Value: "55501", Kind: reflect.String}, Opts: options.FieldOptions{Display: false}, Kind: reflect.String},
				}}

			require.Equal(t, exp, got)
		})
	})

	t.Run("struct_with_integers", func(t *testing.T) {
		type integers struct {
			Int    int    `json:"int" censor:"display"`
			Int8   int8   `json:"int8" censor:"display"`
			Int16  int16  `json:"int16" censor:"display"`
			Int32  int32  `json:"int32" censor:"display"`
			Int64  int64  `json:"int64" censor:"display"`
			Uint   uint   `json:"uint" censor:"display"`
			Uint8  uint8  `json:"uint8" censor:"display"`
			Uint16 uint16 `json:"uint16" censor:"display"`
			Uint32 uint32 `json:"uint32" censor:"display"`
			Uint64 uint64 `json:"uint64" censor:"display"`
			Byte   byte   `json:"byte" censor:"display"`
			Rune   rune   `json:"rune" censor:"display"`
		}

		require.NotPanics(t, func() {
			v := integers{Int: 1, Int8: 2, Int16: 3, Int32: 4, Int64: 5, Uint: 6, Uint8: 7, Uint16: 8, Uint32: 9, Uint64: 10, Byte: 11, Rune: 'y'}
			got := p.Struct(reflect.ValueOf(v))
			exp := models.Struct{
				Name: "parser.integers",
				Fields: []models.Field{
					{Name: "Int", Value: models.Value{Value: 1, Kind: reflect.Int}, Opts: options.FieldOptions{Display: true}, Kind: reflect.Int},
					{Name: "Int8", Value: models.Value{Value: int8(2), Kind: reflect.Int8}, Opts: options.FieldOptions{Display: true}, Kind: reflect.Int8},
					{Name: "Int16", Value: models.Value{Value: int16(3), Kind: reflect.Int16}, Opts: options.FieldOptions{Display: true}, Kind: reflect.Int16},
					{Name: "Int32", Value: models.Value{Value: int32(4), Kind: reflect.Int32}, Opts: options.FieldOptions{Display: true}, Kind: reflect.Int32},
					{Name: "Int64", Value: models.Value{Value: int64(5), Kind: reflect.Int64}, Opts: options.FieldOptions{Display: true}, Kind: reflect.Int64},
					{Name: "Uint", Value: models.Value{Value: uint(6), Kind: reflect.Uint}, Opts: options.FieldOptions{Display: true}, Kind: reflect.Uint},
					{Name: "Uint8", Value: models.Value{Value: uint8(7), Kind: reflect.Uint8}, Opts: options.FieldOptions{Display: true}, Kind: reflect.Uint8},
					{Name: "Uint16", Value: models.Value{Value: uint16(8), Kind: reflect.Uint16}, Opts: options.FieldOptions{Display: true}, Kind: reflect.Uint16},
					{Name: "Uint32", Value: models.Value{Value: uint32(9), Kind: reflect.Uint32}, Opts: options.FieldOptions{Display: true}, Kind: reflect.Uint32},
					{Name: "Uint64", Value: models.Value{Value: uint64(10), Kind: reflect.Uint64}, Opts: options.FieldOptions{Display: true}, Kind: reflect.Uint64},
					{Name: "Byte", Value: models.Value{Value: byte(11), Kind: reflect.Uint8}, Opts: options.FieldOptions{Display: true}, Kind: reflect.Uint8},
					{Name: "Rune", Value: models.Value{Value: rune(121), Kind: reflect.Int32}, Opts: options.FieldOptions{Display: true}, Kind: reflect.Int32},
				},
			}
			require.Equal(t, exp, got)
		})
	})

	t.Run("struct_with_floats", func(t *testing.T) {
		type person struct {
			Name   string  `json:"name" censor:"display"`
			Height float32 `json:"height"`
			Weight float64 `json:"weight"`
		}

		require.NotPanics(t, func() {
			v := person{Name: "John", Height: 1.82, Weight: 82.5}
			got := p.Struct(reflect.ValueOf(v))
			exp := models.Struct{
				Name: "parser.person",
				Fields: []models.Field{
					{Name: "Name", Value: models.Value{Value: "John", Kind: reflect.String}, Opts: options.FieldOptions{Display: true}, Kind: reflect.String},
					{Name: "Height", Value: models.Value{Value: float32(1.82), Kind: reflect.Float32}, Opts: options.FieldOptions{Display: false}, Kind: reflect.Float32},
					{Name: "Weight", Value: models.Value{Value: 82.5, Kind: reflect.Float64}, Opts: options.FieldOptions{Display: false}, Kind: reflect.Float64},
				}}

			require.Equal(t, exp, got)
		})
	})

	t.Run("struct_with_bools", func(t *testing.T) {
		type person struct {
			Active bool `json:"active"`
		}

		require.NotPanics(t, func() {
			v := person{Active: true}
			got := p.Struct(reflect.ValueOf(v))
			exp := models.Struct{
				Name: "parser.person",
				Fields: []models.Field{
					{Name: "Active", Value: models.Value{Value: true, Kind: reflect.Bool}, Opts: options.FieldOptions{Display: false}, Kind: reflect.Bool},
				}}

			require.Equal(t, exp, got)
		})
	})

	t.Run("struct_with_complexes", func(t *testing.T) {
		type person struct {
			Height complex64 `json:"height"`
			Weight complex128
		}

		require.NotPanics(t, func() {
			v := person{Height: complex(1.82, 0), Weight: complex(82.5, 0)}
			got := p.Struct(reflect.ValueOf(v))
			exp := models.Struct{
				Name: "parser.person",
				Fields: []models.Field{
					{Name: "Height", Value: models.Value{Value: complex64(complex(1.82, 0)), Kind: reflect.Complex64}, Opts: options.FieldOptions{Display: false}, Kind: reflect.Complex64},
					{Name: "Weight", Value: models.Value{Value: complex128(complex(82.5, 0)), Kind: reflect.Complex128}, Opts: options.FieldOptions{Display: false}, Kind: reflect.Complex128},
				}}

			require.Equal(t, exp, got)
		})
	})

	t.Run("struct_with_pointer", func(t *testing.T) {
		type person struct {
			Weight *float64 `json:"weight"`
		}

		require.NotPanics(t, func() {
			f := 43.4
			v := person{Weight: &f}
			got := p.Struct(reflect.ValueOf(v))
			exp := models.Struct{
				Name: "parser.person",
				Fields: []models.Field{
					{
						Name:  "Weight",
						Value: models.Value{Value: models.Ptr{Value: 43.4, Kind: reflect.Float64}, Kind: reflect.Pointer},
						Opts:  options.FieldOptions{Display: false},
						Kind:  reflect.Pointer,
					},
				},
			}

			require.Equal(t, exp, got)
		})
	})

	t.Run("struct_with_struct", func(t *testing.T) {
		type contact struct {
			Email string `json:"email" censor:"display"`
			Phone string `json:"phone" censor:"display"`
		}

		type employee struct {
			Contact contact `json:"contact"`
		}

		require.NotPanics(t, func() {
			v := employee{Contact: contact{Email: "example", Phone: "555-555-5555"}}
			got := p.Struct(reflect.ValueOf(v))
			exp := models.Struct{
				Name: "parser.employee",
				Fields: []models.Field{
					{
						Name: "Contact",
						Value: models.Value{
							Value: models.Struct{
								Name: "parser.contact",
								Fields: []models.Field{
									{Name: "Email", Value: models.Value{Value: "example", Kind: reflect.String}, Opts: options.FieldOptions{Display: true}, Kind: reflect.String},
									{Name: "Phone", Value: models.Value{Value: "555-555-5555", Kind: reflect.String}, Opts: options.FieldOptions{Display: true}, Kind: reflect.String},
								},
							},
							Kind: reflect.Struct,
						},
						Opts: options.FieldOptions{Display: false},
						Kind: reflect.Struct,
					},
				},
			}

			require.Equal(t, exp, got)
		})
	})

	t.Run("struct_with_slice", func(t *testing.T) {
		type person struct {
			Names []string `json:"names" censor:"display"`
		}

		require.NotPanics(t, func() {
			v := person{Names: []string{"John", "Doe"}}
			got := p.Struct(reflect.ValueOf(v))
			exp := models.Struct{
				Name: "parser.person",
				Fields: []models.Field{
					{
						Name: "Names",
						Value: models.Value{
							Value: models.Slice{
								Values: []models.Value{
									{Value: "John", Kind: reflect.String},
									{Value: "Doe", Kind: reflect.String},
								},
							},
							Kind: reflect.Slice,
						},
						Opts: options.FieldOptions{Display: true},
						Kind: reflect.Slice,
					},
				},
			}

			require.Equal(t, exp, got)
		})
	})

	t.Run("struct_with_array", func(t *testing.T) {
		type person struct {
			Names [2]string `json:"names" censor:"display"`
		}

		require.NotPanics(t, func() {
			v := person{Names: [2]string{"John", "Doe"}}
			got := p.Struct(reflect.ValueOf(v))
			exp := models.Struct{
				Name: "parser.person",
				Fields: []models.Field{
					{
						Name: "Names",
						Value: models.Value{
							Value: models.Slice{
								Values: []models.Value{
									{Value: "John", Kind: reflect.String},
									{Value: "Doe", Kind: reflect.String},
								},
							},
							Kind: reflect.Array,
						},
						Opts: options.FieldOptions{Display: true},
						Kind: reflect.Array,
					},
				},
			}

			require.Equal(t, exp, got)
		})
	})

	t.Run("struct_with_map", func(t *testing.T) {
		type person struct {
			Names map[string]string `json:"names" censor:"display"`
		}

		require.NotPanics(t, func() {
			v := person{Names: map[string]string{"first": "John", "last": "Doe"}}
			got := p.Struct(reflect.ValueOf(v))
			exp := models.Struct{
				Name: "parser.person",
				Fields: []models.Field{
					{
						Name: "Names",
						Value: models.Value{
							Value: models.Map{
								Values: []models.KV{
									{Key: models.Value{Value: "first", Kind: reflect.String}, Value: models.Value{Value: "John", Kind: reflect.String}, SortValue: "first"},
									{Key: models.Value{Value: "last", Kind: reflect.String}, Value: models.Value{Value: "Doe", Kind: reflect.String}, SortValue: "last"},
								},
								Type: "map[string]string",
							},
							Kind: reflect.Map,
						},
						Opts: options.FieldOptions{Display: true},
						Kind: reflect.Map,
					},
				},
			}

			require.Equal(t, exp, got)
		})
	})

	t.Run("struct_with_interface", func(t *testing.T) {
		type person struct {
			Names interface{} `json:"names" censor:"display"`
		}

		require.NotPanics(t, func() {
			v := person{Names: []string{"John", "Doe"}}
			got := p.Struct(reflect.ValueOf(v))
			exp := models.Struct{
				Name: "parser.person",
				Fields: []models.Field{
					{
						Name: "Names",
						Value: models.Value{
							Value: models.Value{
								Value: models.Slice{
									Values: []models.Value{
										{Value: "John", Kind: reflect.String},
										{Value: "Doe", Kind: reflect.String},
									},
								},
								Kind: reflect.Slice,
							},
							Kind: reflect.Interface,
						},
						Opts: options.FieldOptions{Display: true}, Kind: reflect.Interface,
					},
				},
			}

			require.Equal(t, exp, got)
		})
	})

	t.Run("struct_with_interface_with_struct_value", func(t *testing.T) {
		type contact struct {
			Email string `json:"email" censor:"display"`
		}

		type person struct {
			Contact interface{} `json:"contact"`
		}

		require.NotPanics(t, func() {
			v := person{Contact: contact{Email: "example"}}
			got := p.Struct(reflect.ValueOf(v))
			exp := models.Struct{
				Name: "parser.person",
				Fields: []models.Field{
					{
						Name: "Contact",
						Value: models.Value{
							Value: models.Value{
								Value: models.Struct{
									Name: "parser.contact",
									Fields: []models.Field{
										{
											Name:  "Email",
											Value: models.Value{Value: "example", Kind: reflect.String},
											Opts:  options.FieldOptions{Display: true},
											Kind:  reflect.String,
										},
									},
								},
								Kind: reflect.Struct,
							},
							Kind: reflect.Interface,
						},
						Opts: options.FieldOptions{Display: false},
						Kind: reflect.Interface,
					},
				},
			}

			require.Equal(t, exp, got)
		})
	})

	t.Run("struct_with_interface_with_map_value", func(t *testing.T) {
		type person struct {
			Names interface{} `json:"names" censor:"display"`
		}

		require.NotPanics(t, func() {
			v := person{Names: map[string]string{"first": "John", "last": "Doe"}}
			got := p.Struct(reflect.ValueOf(v))
			exp := models.Struct{
				Name: "parser.person",
				Fields: []models.Field{
					{
						Name: "Names",
						Value: models.Value{
							Value: models.Value{
								Value: models.Map{
									Type: "map[string]string",
									Values: []models.KV{
										{
											SortValue: "first",
											Key:       models.Value{Value: "first", Kind: reflect.String},
											Value:     models.Value{Value: "John", Kind: reflect.String},
										},
										{
											SortValue: "last",
											Key:       models.Value{Value: "last", Kind: reflect.String},
											Value:     models.Value{Value: "Doe", Kind: reflect.String},
										},
									},
								},
								Kind: reflect.Map,
							},
							Kind: reflect.Interface,
						},
						Opts: options.FieldOptions{Display: true},
						Kind: reflect.Interface,
					},
				},
			}

			require.Equal(t, exp, got)
		})
	})

	t.Run("struct_with_interface_with_pointer_value", func(t *testing.T) {
		type person struct {
			Names interface{} `json:"names" censor:"display"`
		}

		require.NotPanics(t, func() {
			f := 43.4
			v := person{Names: &f}
			got := p.Struct(reflect.ValueOf(v))
			exp := models.Struct{
				Name: "parser.person",
				Fields: []models.Field{
					{
						Name: "Names",
						Value: models.Value{
							Value: models.Value{
								Value: models.Ptr{Value: 43.4, Kind: reflect.Float64},
								Kind:  reflect.Pointer,
							},
							Kind: reflect.Interface,
						},
						Opts: options.FieldOptions{Display: true},
						Kind: reflect.Interface,
					},
				},
			}

			require.Equal(t, exp, got)
		})
	})

	t.Run("struct_with_interface_with_complex_value", func(t *testing.T) {
		type person struct {
			Names interface{} `json:"names" censor:"display"`
		}

		require.NotPanics(t, func() {
			v := person{Names: complex(1.82, 0)}
			got := p.Struct(reflect.ValueOf(v))
			exp := models.Struct{
				Name: "parser.person",
				Fields: []models.Field{
					{
						Name: "Names",
						Value: models.Value{
							Value: models.Value{Value: (1.82 + 0i), Kind: reflect.Complex128},
							Kind:  reflect.Interface,
						},
						Opts: options.FieldOptions{Display: true},
						Kind: reflect.Interface,
					},
				},
			}

			require.Equal(t, exp, got)
		})
	})

	t.Run("struct_with_unexported_field", func(t *testing.T) {
		type person struct {
			Name   string
			height float32
			weight float64
		}

		require.NotPanics(t, func() {
			v := person{Name: "John", height: 1.82, weight: 82.5}
			got := p.Struct(reflect.ValueOf(v))
			exp := models.Struct{
				Name: "parser.person",
				Fields: []models.Field{
					{Name: "Name", Value: models.Value{Value: "John", Kind: reflect.String}, Opts: options.FieldOptions{Display: false}, Kind: reflect.String},
				}}

			require.Equal(t, exp, got)
		})
	})

	t.Run("struct_with_all_unexported_fields", func(t *testing.T) {
		type person struct {
			name   string
			height float32
			weight float64
		}

		require.NotPanics(t, func() {
			v := person{name: "John", height: 1.82, weight: 82.5}
			got := p.Struct(reflect.ValueOf(v))
			exp := models.Struct{
				Name:   "parser.person",
				Fields: []models.Field{},
			}

			require.Equal(t, exp, got)
		})
	})

	t.Run("non_struct_value", func(t *testing.T) {
		require.PanicsWithValue(t, "provided value is not a struct", func() { p.Struct(reflect.ValueOf(5.234)) })
	})
}

func TestParser_StructWithJSONTags(t *testing.T) {
	type address struct {
		City   string `json:"city" censor:"display"`
		State  string `json:"state" censor:"display"`
		Street string `json:"street"`
		Zip    string
	}

	p := Parser{
		useJSONTagName: true,
		censorFieldTag: DefaultCensorFieldTag,
	}

	t.Run("struct_with_strings", func(t *testing.T) {
		require.NotPanics(t, func() {
			v := address{Street: "451 Main St", City: "San Francisco", State: "CA", Zip: "55501"}
			got := p.Struct(reflect.ValueOf(v))
			exp := models.Struct{
				Name: "parser.address",
				Fields: []models.Field{
					{Name: "city", Value: models.Value{Value: "San Francisco", Kind: reflect.String}, Opts: options.FieldOptions{Display: true}, Kind: reflect.String},
					{Name: "state", Value: models.Value{Value: "CA", Kind: reflect.String}, Opts: options.FieldOptions{Display: true}, Kind: reflect.String},
					{Name: "street", Value: models.Value{Value: "451 Main St", Kind: reflect.String}, Opts: options.FieldOptions{Display: false}, Kind: reflect.String},
				}}

			require.Equal(t, exp, got)
		})
	})
}

func TestMain(m *testing.M) {
	type address struct {
		City   string `censor:"display"`
		State  string `censor:"display"`
		Street string
	}

	mainPkgStruct = address{Street: "451 Main St", City: "San Francisco", State: "CA"}

	exitVal := m.Run()

	os.Exit(exitVal)
}
