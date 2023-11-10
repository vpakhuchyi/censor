package parser

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/vpakhuchyi/censor/internal/models"
	"github.com/vpakhuchyi/censor/internal/options"
)

func TestParser_Interface(t *testing.T) {
	p := Parser{
		UseJSONTagName: false,
		CensorFieldTag: DefaultCensorFieldTag,
	}

	t.Run("struct_with_interface_with_slice_value", func(t *testing.T) {
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
						Tag:  "display",
						Value: models.Value{
							Value: models.Interface{
								Name: "",
								Value: models.Value{
									Value: models.Slice{
										Values: []models.Value{{Value: "John", Kind: 0x18}, {Value: "Doe", Kind: 0x18}},
									},
									Kind: reflect.Slice},
							},
							Kind: reflect.Interface,
						},
						Opts: options.FieldOptions{Display: true}, Kind: 0x14},
				},
			}

			require.Equal(t, exp, got)
		})
	})

	t.Run("struct_with_interface_with_struct_value", func(t *testing.T) {
		type contact struct {
			Email string `json:"email" censor:"display"`
			Phone string `json:"phone" censor:"display"`
		}

		type person struct {
			Contact interface{} `json:"contact"`
		}

		require.NotPanics(t, func() {
			v := person{Contact: contact{Email: "example", Phone: "555-555-5555"}}
			got := p.Struct(reflect.ValueOf(v))
			exp := models.Struct{
				Name: "parser.person",
				Fields: []models.Field{
					{
						Name: "Contact",
						Tag:  "",
						Value: models.Value{
							Value: models.Interface{
								Name: "",
								Value: models.Value{
									Value: models.Struct{
										Name: "parser.contact",
										Fields: []models.Field{
											{Name: "Email", Tag: "display", Value: models.Value{Value: "example", Kind: reflect.String}, Opts: options.FieldOptions{Display: true}, Kind: reflect.String},
											{Name: "Phone", Tag: "display", Value: models.Value{Value: "555-555-5555", Kind: reflect.String}, Opts: options.FieldOptions{Display: true}, Kind: reflect.String},
										},
									},
									Kind: reflect.Struct,
								},
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
						Tag:  "display",
						Value: models.Value{
							Value: models.Interface{
								Name: "",
								Value: models.Value{
									Value: models.Map{
										Values: []models.KV{
											{Key: models.Value{Value: "first", Kind: reflect.String}, Value: models.Value{Value: "John", Kind: reflect.String}, SortValue: "first"},
											{Key: models.Value{Value: "last", Kind: reflect.String}, Value: models.Value{Value: "Doe", Kind: reflect.String}, SortValue: "last"},
										},
										Type: "map[string]string",
									},
									Kind: reflect.Map},
							},
							Kind: reflect.Interface,
						},
						Opts: options.FieldOptions{Display: true}, Kind: reflect.Interface},
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
						Tag:  "display",
						Value: models.Value{
							Value: models.Interface{
								Name: "",
								Value: models.Value{
									Value: models.Ptr{Value: 43.4, Kind: reflect.Float64},
									Kind:  reflect.Ptr,
								},
							},
							Kind: reflect.Interface,
						},
						Opts: options.FieldOptions{Display: true}, Kind: reflect.Interface},
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
						Tag:  "display",
						Value: models.Value{
							Value: models.Interface{
								Name: "",
								Value: models.Value{
									Value: complex(1.82, 0),
									Kind:  reflect.Complex128,
								},
							},
							Kind: reflect.Interface,
						},
						Opts: options.FieldOptions{Display: true}, Kind: reflect.Interface},
				},
			}

			require.Equal(t, exp, got)
		})
	})

	t.Run("struct_with_interface_with_bool_value", func(t *testing.T) {
		type person struct {
			Names interface{} `json:"names" censor:"display"`
		}

		require.NotPanics(t, func() {
			v := person{Names: true}
			got := p.Struct(reflect.ValueOf(v))
			exp := models.Struct{
				Name: "parser.person",
				Fields: []models.Field{
					{
						Name: "Names",
						Tag:  "display",
						Value: models.Value{
							Value: models.Interface{
								Name: "",
								Value: models.Value{
									Value: true,
									Kind:  reflect.Bool,
								},
							},
							Kind: reflect.Interface,
						},
						Opts: options.FieldOptions{Display: true}, Kind: reflect.Interface},
				},
			}

			require.Equal(t, exp, got)
		})
	})

	t.Run("struct_with_interface_with_int_value", func(t *testing.T) {
		type person struct {
			Names interface{} `json:"names" censor:"display"`
		}

		require.NotPanics(t, func() {
			v := person{Names: 13}
			got := p.Struct(reflect.ValueOf(v))
			exp := models.Struct{
				Name: "parser.person",
				Fields: []models.Field{
					{
						Name: "Names",
						Tag:  "display",
						Value: models.Value{
							Value: models.Interface{
								Name: "",
								Value: models.Value{
									Value: 13,
									Kind:  reflect.Int,
								},
							},
							Kind: reflect.Interface,
						},
						Opts: options.FieldOptions{Display: true}, Kind: reflect.Interface},
				},
			}

			require.Equal(t, exp, got)
		})
	})

	t.Run("struct_with_interface_with_float_value", func(t *testing.T) {
		type person struct {
			Names interface{} `json:"names" censor:"display"`
		}

		require.NotPanics(t, func() {
			v := person{Names: 13.5}
			got := p.Struct(reflect.ValueOf(v))
			exp := models.Struct{
				Name: "parser.person",
				Fields: []models.Field{
					{
						Name: "Names",
						Tag:  "display",
						Value: models.Value{
							Value: models.Interface{
								Name: "",
								Value: models.Value{
									Value: 13.5,
									Kind:  reflect.Float64,
								},
							},
							Kind: reflect.Interface,
						},
						Opts: options.FieldOptions{Display: true}, Kind: reflect.Interface},
				},
			}

			require.Equal(t, exp, got)
		})
	})

	t.Run("struct_with_interface_with_string_value", func(t *testing.T) {
		type person struct {
			Names interface{} `json:"names" censor:"display"`
		}

		require.NotPanics(t, func() {
			v := person{Names: "John"}
			got := p.Struct(reflect.ValueOf(v))
			exp := models.Struct{
				Name: "parser.person",
				Fields: []models.Field{
					{
						Name: "Names",
						Tag:  "display",
						Value: models.Value{
							Value: models.Interface{
								Name: "",
								Value: models.Value{
									Value: "John",
									Kind:  reflect.String,
								},
							},
							Kind: reflect.Interface,
						},
						Opts: options.FieldOptions{Display: true}, Kind: reflect.Interface},
				},
			}

			require.Equal(t, exp, got)
		})
	})

	t.Run("non_interface_value", func(t *testing.T) {
		require.PanicsWithValue(t, "provided value is not an interface", func() { p.Interface(reflect.ValueOf(5.234)) })
	})
}
