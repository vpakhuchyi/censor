package parser

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/vpakhuchyi/censor/internal/models"
	"github.com/vpakhuchyi/censor/internal/options"
)

func TestParser_Map(t *testing.T) {
	type address struct {
		City   string `json:"city" censor:"display"`
		State  string `json:"state" censor:"display"`
		Street string `json:"street"`
		Zip    string `json:"zip"`
	}

	p := Parser{
		UseJSONTagName: false,
		CensorFieldTag: DefaultCensorFieldTag,
	}

	// Test values.
	t.Run("map_string_string", func(t *testing.T) {
		require.NotPanics(t, func() {
			v := map[string]string{"key1": "value1", "key2": "value2"}
			got := p.Map(reflect.ValueOf(v))
			exp := models.Map{
				Type: "map[string]string",
				Values: []models.KV{
					{Key: models.Value{Value: "key1", Kind: reflect.String}, Value: models.Value{Value: "value1", Kind: reflect.String}, SortValue: "key1"},
					{Key: models.Value{Value: "key2", Kind: reflect.String}, Value: models.Value{Value: "value2", Kind: reflect.String}, SortValue: "key2"},
				},
			}

			require.Equal(t, exp, got)
		})
	})

	t.Run("map_string_int", func(t *testing.T) {
		require.NotPanics(t, func() {
			v := map[string]int{"key1": 1, "key2": 2}
			got := p.Map(reflect.ValueOf(v))
			exp := models.Map{
				Type: "map[string]int",
				Values: []models.KV{
					{Key: models.Value{Value: "key1", Kind: reflect.String}, Value: models.Value{Value: 1, Kind: reflect.Int}, SortValue: "key1"},
					{Key: models.Value{Value: "key2", Kind: reflect.String}, Value: models.Value{Value: 2, Kind: reflect.Int}, SortValue: "key2"},
				},
			}

			require.Equal(t, exp, got)
		})
	})

	t.Run("map_string_float32", func(t *testing.T) {
		require.NotPanics(t, func() {
			v := map[string]float32{"key1": 1.234, "key2": 2.345}
			got := p.Map(reflect.ValueOf(v))
			exp := models.Map{
				Type: "map[string]float32",
				Values: []models.KV{
					{Key: models.Value{Value: "key1", Kind: reflect.String}, Value: models.Value{Value: float32(1.234), Kind: reflect.Float32}, SortValue: "key1"},
					{Key: models.Value{Value: "key2", Kind: reflect.String}, Value: models.Value{Value: float32(2.345), Kind: reflect.Float32}, SortValue: "key2"},
				},
			}

			require.Equal(t, exp, got)
		})
	})

	t.Run("map_string_float64", func(t *testing.T) {
		require.NotPanics(t, func() {
			v := map[string]float64{"key1": 1.234, "key2": 2.345}
			got := p.Map(reflect.ValueOf(v))
			exp := models.Map{
				Type: "map[string]float64",
				Values: []models.KV{
					{Key: models.Value{Value: "key1", Kind: reflect.String}, Value: models.Value{Value: 1.234, Kind: reflect.Float64}, SortValue: "key1"},
					{Key: models.Value{Value: "key2", Kind: reflect.String}, Value: models.Value{Value: 2.345, Kind: reflect.Float64}, SortValue: "key2"},
				},
			}

			require.Equal(t, exp, got)
		})
	})

	t.Run("map_string_bool", func(t *testing.T) {
		require.NotPanics(t, func() {
			v := map[string]bool{"key1": true, "key2": false}
			got := p.Map(reflect.ValueOf(v))
			exp := models.Map{
				Type: "map[string]bool",
				Values: []models.KV{
					{Key: models.Value{Value: "key1", Kind: reflect.String}, Value: models.Value{Value: true, Kind: reflect.Bool}, SortValue: "key1"},
					{Key: models.Value{Value: "key2", Kind: reflect.String}, Value: models.Value{Value: false, Kind: reflect.Bool}, SortValue: "key2"},
				},
			}

			require.Equal(t, exp, got)
		})
	})

	t.Run("map_string_complex64", func(t *testing.T) {
		require.NotPanics(t, func() {
			v := map[string]complex64{"key1": complex(float32(1.234), float32(2.345)), "key2": complex(float32(3.456), float32(4.567))}
			got := p.Map(reflect.ValueOf(v))
			exp := models.Map{
				Type: "map[string]complex64",
				Values: []models.KV{
					{Key: models.Value{Value: "key1", Kind: reflect.String}, Value: models.Value{Value: complex(float32(1.234), float32(2.345)), Kind: reflect.Complex64}, SortValue: "key1"},
					{Key: models.Value{Value: "key2", Kind: reflect.String}, Value: models.Value{Value: complex(float32(3.456), float32(4.567)), Kind: reflect.Complex64}, SortValue: "key2"},
				},
			}

			require.Equal(t, exp, got)
		})
	})

	t.Run("map_string_complex128", func(t *testing.T) {
		require.NotPanics(t, func() {
			v := map[string]complex128{"key1": complex(float64(1.234), float64(2.345)), "key2": complex(float64(3.456), float64(4.567))}
			got := p.Map(reflect.ValueOf(v))
			exp := models.Map{
				Type: "map[string]complex128",
				Values: []models.KV{
					{Key: models.Value{Value: "key1", Kind: reflect.String}, Value: models.Value{Value: complex(float64(1.234), float64(2.345)), Kind: reflect.Complex128}, SortValue: "key1"},
					{Key: models.Value{Value: "key2", Kind: reflect.String}, Value: models.Value{Value: complex(float64(3.456), float64(4.567)), Kind: reflect.Complex128}, SortValue: "key2"},
				},
			}

			require.Equal(t, exp, got)
		})
	})

	t.Run("map_string_struct", func(t *testing.T) {
		require.NotPanics(t, func() {
			v := map[string]address{"key1": {Street: "451 Main St", City: "San Francisco", State: "CA", Zip: "55501"}, "key2": {Street: "65 Best St", City: "Denver", State: "DN", Zip: "55502"}}
			got := p.Map(reflect.ValueOf(v))
			exp := models.Map{
				Type: "map[string]parser.address",
				Values: []models.KV{
					{Key: models.Value{Value: "key1", Kind: reflect.String}, Value: models.Value{Value: models.Struct{
						Name: "parser.address",
						Fields: []models.Field{
							{Name: "City", Tag: "display", Value: models.Value{Value: "San Francisco", Kind: reflect.String}, Kind: reflect.String, Opts: options.FieldOptions{Display: true}},
							{Name: "State", Tag: "display", Value: models.Value{Value: "CA", Kind: reflect.String}, Kind: reflect.String, Opts: options.FieldOptions{Display: true}},
							{Name: "Street", Tag: "", Value: models.Value{Value: "451 Main St", Kind: reflect.String}, Kind: reflect.String, Opts: options.FieldOptions{Display: false}},
							{Name: "Zip", Tag: "", Value: models.Value{Value: "55501", Kind: reflect.String}, Kind: reflect.String, Opts: options.FieldOptions{Display: false}}}}, Kind: reflect.Struct},
						SortValue: "key1"},

					{Key: models.Value{Value: "key2", Kind: reflect.String}, Value: models.Value{Value: models.Struct{
						Name: "parser.address",
						Fields: []models.Field{
							{Name: "City", Tag: "display", Value: models.Value{Value: "Denver", Kind: reflect.String}, Kind: reflect.String, Opts: options.FieldOptions{Display: true}},
							{Name: "State", Tag: "display", Value: models.Value{Value: "DN", Kind: reflect.String}, Kind: reflect.String, Opts: options.FieldOptions{Display: true}},
							{Name: "Street", Tag: "", Value: models.Value{Value: "65 Best St", Kind: reflect.String}, Kind: reflect.String, Opts: options.FieldOptions{Display: false}},
							{Name: "Zip", Tag: "", Value: models.Value{Value: "55502", Kind: reflect.String}, Kind: reflect.String, Opts: options.FieldOptions{Display: false}}}}, Kind: reflect.Struct},
						SortValue: "key2"},
				},
			}

			require.Equal(t, exp, got)
		})
	})

	t.Run("map_string_slice", func(t *testing.T) {
		require.NotPanics(t, func() {
			v := map[string][]string{"key1": {"tag1", "tag2"}, "key2": {"tag3", "tag4"}}
			got := p.Map(reflect.ValueOf(v))
			exp := models.Map{
				Type: "map[string][]string",
				Values: []models.KV{
					{Key: models.Value{Value: "key1", Kind: reflect.String}, Value: models.Value{Value: models.Slice{
						Values: []models.Value{{Value: "tag1", Kind: reflect.String}, {Value: "tag2", Kind: reflect.String}},
					}, Kind: reflect.Slice}, SortValue: "key1"},

					{Key: models.Value{Value: "key2", Kind: reflect.String}, Value: models.Value{Value: models.Slice{
						Values: []models.Value{{Value: "tag3", Kind: reflect.String}, {Value: "tag4", Kind: reflect.String}},
					}, Kind: reflect.Slice}, SortValue: "key2"},
				},
			}

			require.Equal(t, exp, got)
		})
	})

	t.Run("map_string_array", func(t *testing.T) {
		require.NotPanics(t, func() {
			v := map[string][2]string{"key1": {"tag1", "tag2"}, "key2": {"tag3", "tag4"}}
			got := p.Map(reflect.ValueOf(v))
			exp := models.Map{
				Type: "map[string][2]string",
				Values: []models.KV{
					{Key: models.Value{Value: "key1", Kind: reflect.String}, Value: models.Value{Value: models.Slice{
						Values: []models.Value{{Value: "tag1", Kind: reflect.String}, {Value: "tag2", Kind: reflect.String}},
					}, Kind: reflect.Array}, SortValue: "key1"},

					{Key: models.Value{Value: "key2", Kind: reflect.String}, Value: models.Value{Value: models.Slice{
						Values: []models.Value{{Value: "tag3", Kind: reflect.String}, {Value: "tag4", Kind: reflect.String}},
					}, Kind: reflect.Array}, SortValue: "key2"},
				},
			}

			require.Equal(t, exp, got)
		})
	})

	t.Run("map_string_pointer_to_struct", func(t *testing.T) {
		require.NotPanics(t, func() {
			v := map[string]*address{"key1": {Street: "451 Main St", City: "San Francisco", State: "CA", Zip: "55501"}, "key2": {Street: "65 Best St", City: "Denver", State: "DN", Zip: "55502"}}
			got := p.Map(reflect.ValueOf(v))
			exp := models.Map{
				Type: "map[string]*parser.address",
				Values: []models.KV{
					{Key: models.Value{Value: "key1", Kind: reflect.String}, Value: models.Value{
						Value: models.Ptr{
							Kind: reflect.Struct,
							Value: models.Struct{
								Name: "parser.address",
								Fields: []models.Field{
									{Name: "City", Tag: "display", Value: models.Value{Value: "San Francisco", Kind: reflect.String}, Kind: reflect.String, Opts: options.FieldOptions{Display: true}},
									{Name: "State", Tag: "display", Value: models.Value{Value: "CA", Kind: reflect.String}, Kind: reflect.String, Opts: options.FieldOptions{Display: true}},
									{Name: "Street", Tag: "", Value: models.Value{Value: "451 Main St", Kind: reflect.String}, Kind: reflect.String, Opts: options.FieldOptions{Display: false}},
									{Name: "Zip", Tag: "", Value: models.Value{Value: "55501", Kind: reflect.String}, Kind: reflect.String, Opts: options.FieldOptions{Display: false}},
								},
							},
						}, Kind: reflect.Ptr},
						SortValue: "key1"},
					{Key: models.Value{Value: "key2", Kind: reflect.String}, Value: models.Value{
						Value: models.Ptr{
							Kind: reflect.Struct,
							Value: models.Struct{
								Name: "parser.address",
								Fields: []models.Field{
									{Name: "City", Tag: "display", Value: models.Value{Value: "Denver", Kind: reflect.String}, Kind: reflect.String, Opts: options.FieldOptions{Display: true}},
									{Name: "State", Tag: "display", Value: models.Value{Value: "DN", Kind: reflect.String}, Kind: reflect.String, Opts: options.FieldOptions{Display: true}},
									{Name: "Street", Tag: "", Value: models.Value{Value: "65 Best St", Kind: reflect.String}, Kind: reflect.String, Opts: options.FieldOptions{Display: false}},
									{Name: "Zip", Tag: "", Value: models.Value{Value: "55502", Kind: reflect.String}, Kind: reflect.String, Opts: options.FieldOptions{Display: false}},
								},
							},
						},
						Kind: reflect.Ptr},
						SortValue: "key2"},
				},
			}

			require.Equal(t, exp, got)
		})
	})

	t.Run("map_string_map_string", func(t *testing.T) {
		require.NotPanics(t, func() {
			v := map[string]map[string]string{"key1": {"key1": "value1", "key2": "value2"}, "key2": {"key3": "value3", "key4": "value4"}}
			got := p.Map(reflect.ValueOf(v))
			exp := models.Map{
				Type: "map[string]map[string]string",
				Values: []models.KV{
					{Key: models.Value{Value: "key1", Kind: reflect.String}, Value: models.Value{
						Value: models.Map{
							Type: "map[string]string",
							Values: []models.KV{
								{Key: models.Value{Value: "key1", Kind: reflect.String}, Value: models.Value{Value: "value1", Kind: reflect.String}, SortValue: "key1"},
								{Key: models.Value{Value: "key2", Kind: reflect.String}, Value: models.Value{Value: "value2", Kind: reflect.String}, SortValue: "key2"},
							},
						},
						Kind: reflect.Map},
						SortValue: "key1"},
					{Key: models.Value{Value: "key2", Kind: reflect.String}, Value: models.Value{
						Value: models.Map{
							Type: "map[string]string",
							Values: []models.KV{
								{Key: models.Value{Value: "key3", Kind: reflect.String}, Value: models.Value{Value: "value3", Kind: reflect.String}, SortValue: "key3"},
								{Key: models.Value{Value: "key4", Kind: reflect.String}, Value: models.Value{Value: "value4", Kind: reflect.String}, SortValue: "key4"},
							},
						},
						Kind: reflect.Map},
						SortValue: "key2"},
				},
			}

			require.Equal(t, exp, got)
		})
	})

	t.Run("map_string_interface", func(t *testing.T) {
		require.NotPanics(t, func() {
			v := map[string]interface{}{"key1": 1, "key2": "value2"}
			got := p.Map(reflect.ValueOf(v))
			exp := models.Map{
				Type: "map[string]interface {}",
				Values: []models.KV{
					{
						Key:   models.Value{Value: "key1", Kind: reflect.String},
						Value: models.Value{Value: models.Interface{Value: models.Value{Value: 1, Kind: reflect.Int}}, Kind: reflect.Interface}, SortValue: "key1",
					},
					{
						Key:   models.Value{Value: "key2", Kind: reflect.String},
						Value: models.Value{Value: models.Interface{Value: models.Value{Value: "value2", Kind: reflect.String}}, Kind: reflect.Interface}, SortValue: "key2",
					},
				},
			}

			require.Equal(t, exp, got)
		})
	})

	t.Run("map_string_func", func(t *testing.T) {
		require.NotPanics(t, func() {
			v := map[string]func(){"key1": func() {}, "key2": func() {}}
			got := p.Map(reflect.ValueOf(v))
			exp := models.Map{
				Type: "map[string]func()",
				Values: []models.KV{
					{Key: models.Value{Value: "key1", Kind: reflect.String}, Value: models.Value{Value: "unsupported type: func", Kind: reflect.String}, SortValue: "key1"},
					{Key: models.Value{Value: "key2", Kind: reflect.String}, Value: models.Value{Value: "unsupported type: func", Kind: reflect.String}, SortValue: "key2"},
				},
			}

			require.Equal(t, exp, got)
		})
	})

	// Test keys.
	t.Run("map_int_string", func(t *testing.T) {
		require.NotPanics(t, func() {
			v := map[int]string{1: "value1", 2: "value2"}
			got := p.Map(reflect.ValueOf(v))
			exp := models.Map{
				Type: "map[int]string",
				Values: []models.KV{
					{Key: models.Value{Value: 1, Kind: reflect.Int}, Value: models.Value{Value: "value1", Kind: reflect.String}, SortValue: "1"},
					{Key: models.Value{Value: 2, Kind: reflect.Int}, Value: models.Value{Value: "value2", Kind: reflect.String}, SortValue: "2"},
				},
			}

			require.Equal(t, exp, got)
		})
	})

	t.Run("map_float32_string", func(t *testing.T) {
		require.NotPanics(t, func() {
			v := map[float32]string{1.1: "value1", 2.2: "value2"}
			got := p.Map(reflect.ValueOf(v))
			exp := models.Map{
				Type: "map[float32]string",
				Values: []models.KV{
					{Key: models.Value{Value: float32(1.1), Kind: reflect.Float32}, Value: models.Value{Value: "value1", Kind: reflect.String}, SortValue: "1.1"},
					{Key: models.Value{Value: float32(2.2), Kind: reflect.Float32}, Value: models.Value{Value: "value2", Kind: reflect.String}, SortValue: "2.2"},
				},
			}

			require.Equal(t, exp, got)
		})
	})

	t.Run("map_float64_string", func(t *testing.T) {
		require.NotPanics(t, func() {
			v := map[float64]string{1.1: "value1", 2.2: "value2"}
			got := p.Map(reflect.ValueOf(v))
			exp := models.Map{
				Type: "map[float64]string",
				Values: []models.KV{
					{Key: models.Value{Value: 1.1, Kind: reflect.Float64}, Value: models.Value{Value: "value1", Kind: reflect.String}, SortValue: "1.1"},
					{Key: models.Value{Value: 2.2, Kind: reflect.Float64}, Value: models.Value{Value: "value2", Kind: reflect.String}, SortValue: "2.2"},
				},
			}

			require.Equal(t, exp, got)
		})
	})

	t.Run("map_bool_string", func(t *testing.T) {
		require.NotPanics(t, func() {
			v := map[bool]string{true: "value1", false: "value2"}
			got := p.Map(reflect.ValueOf(v))
			exp := models.Map{
				Type: "map[bool]string",
				Values: []models.KV{
					{Key: models.Value{Value: false, Kind: reflect.Bool}, Value: models.Value{Value: "value2", Kind: reflect.String}, SortValue: "false"},
					{Key: models.Value{Value: true, Kind: reflect.Bool}, Value: models.Value{Value: "value1", Kind: reflect.String}, SortValue: "true"},
				},
			}

			require.Equal(t, exp, got)
		})
	})

	t.Run("map_complex64_string", func(t *testing.T) {
		require.NotPanics(t, func() {
			v := map[complex64]string{complex(float32(1.234), float32(2.345)): "value1", complex(float32(3.456), float32(4.567)): "value2"}
			got := p.Map(reflect.ValueOf(v))
			exp := models.Map{
				Type: "map[complex64]string",
				Values: []models.KV{
					{Key: models.Value{Value: complex(float32(1.234), float32(2.345)), Kind: reflect.Complex64}, Value: models.Value{Value: "value1", Kind: reflect.String}, SortValue: "(1.234+2.345i)"},
					{Key: models.Value{Value: complex(float32(3.456), float32(4.567)), Kind: reflect.Complex64}, Value: models.Value{Value: "value2", Kind: reflect.String}, SortValue: "(3.456+4.567i)"},
				},
			}

			require.Equal(t, exp, got)
		})
	})

	t.Run("map_complex128_string", func(t *testing.T) {
		require.NotPanics(t, func() {
			v := map[complex128]string{complex(float64(1.234), float64(2.345)): "value1", complex(float64(3.456), float64(4.567)): "value2"}
			got := p.Map(reflect.ValueOf(v))
			exp := models.Map{
				Type: "map[complex128]string",
				Values: []models.KV{
					{Key: models.Value{Value: complex(float64(1.234), float64(2.345)), Kind: reflect.Complex128}, Value: models.Value{Value: "value1", Kind: reflect.String}, SortValue: "(1.234+2.345i)"},
					{Key: models.Value{Value: complex(float64(3.456), float64(4.567)), Kind: reflect.Complex128}, Value: models.Value{Value: "value2", Kind: reflect.String}, SortValue: "(3.456+4.567i)"},
				},
			}

			require.Equal(t, exp, got)
		})
	})

	t.Run("map_rune_string", func(t *testing.T) {
		require.NotPanics(t, func() {
			v := map[rune]string{1: "value1", 2: "value2"}
			got := p.Map(reflect.ValueOf(v))
			exp := models.Map{
				Type: "map[int32]string",
				Values: []models.KV{
					{Key: models.Value{Value: int32(1), Kind: reflect.Int32}, Value: models.Value{Value: "value1", Kind: reflect.String}, SortValue: "1"},
					{Key: models.Value{Value: int32(2), Kind: reflect.Int32}, Value: models.Value{Value: "value2", Kind: reflect.String}, SortValue: "2"},
				},
			}

			require.Equal(t, exp, got)
		})
	})

	t.Run("map_array_string", func(t *testing.T) {
		require.NotPanics(t, func() {
			v := map[[2]string]string{{"key1", "key2"}: "value1", {"key3", "key4"}: "value2"}
			got := p.Map(reflect.ValueOf(v))
			exp := models.Map{
				Type: "map[[2]string]string",
				Values: []models.KV{
					{
						Key:   models.Value{Value: models.Slice{Values: []models.Value{{Value: "key1", Kind: reflect.String}, {Value: "key2", Kind: reflect.String}}}, Kind: reflect.Array},
						Value: models.Value{Value: "value1", Kind: reflect.String}, SortValue: "[key1 key2]",
					},
					{
						Key:   models.Value{Value: models.Slice{Values: []models.Value{{Value: "key3", Kind: reflect.String}, {Value: "key4", Kind: reflect.String}}}, Kind: reflect.Array},
						Value: models.Value{Value: "value2", Kind: reflect.String}, SortValue: "[key3 key4]",
					},
				},
			}

			require.Equal(t, exp, got)
		})
	})

	t.Run("map_pointer_to_string", func(t *testing.T) {
		require.NotPanics(t, func() {
			key1, key2 := "key1", "key2"
			v := map[*string]string{&key1: "value1", &key2: "value2"}
			got := p.Map(reflect.ValueOf(v))
			exp := models.Map{
				Type: "map[*string]string",
				Values: []models.KV{
					{Key: models.Value{Value: models.Ptr{Kind: reflect.String, Value: "key1"}, Kind: reflect.Ptr}, Value: models.Value{Value: "value1", Kind: reflect.String}, SortValue: fmt.Sprintf("%p", &key1)},
					{Key: models.Value{Value: models.Ptr{Kind: reflect.String, Value: "key2"}, Kind: reflect.Ptr}, Value: models.Value{Value: "value2", Kind: reflect.String}, SortValue: fmt.Sprintf("%p", &key2)},
				},
			}

			require.Equal(t, exp, got)
		})
	})

	t.Run("map_interface_string", func(t *testing.T) {
		require.NotPanics(t, func() {
			v := map[interface{}]string{"key1": "value1", "key2": "value2"}
			got := p.Map(reflect.ValueOf(v))
			exp := models.Map{
				Type: "map[interface {}]string",
				Values: []models.KV{
					{Key: models.Value{Value: models.Interface{Value: models.Value{Value: "key1", Kind: reflect.String}}, Kind: reflect.Interface}, Value: models.Value{Value: "value1", Kind: reflect.String}, SortValue: "key1"},
					{Key: models.Value{Value: models.Interface{Value: models.Value{Value: "key2", Kind: reflect.String}}, Kind: reflect.Interface}, Value: models.Value{Value: "value2", Kind: reflect.String}, SortValue: "key2"},
				},
			}

			require.Equal(t, exp, got)
		})
	})

	t.Run("map_struct_string", func(t *testing.T) {
		require.NotPanics(t, func() {
			v := map[address]string{{Street: "451 Main St", City: "San Francisco", State: "CA", Zip: "55501"}: "value1", {Street: "65 Best St", City: "Denver", State: "DN", Zip: "55502"}: "value2"}
			got := p.Map(reflect.ValueOf(v))
			exp := models.Map{
				Type: "map[parser.address]string",
				Values: []models.KV{
					{
						Key: models.Value{
							Value: models.Struct{
								Name: "parser.address",
								Fields: []models.Field{
									{Name: "City", Tag: "display", Value: models.Value{Value: "Denver", Kind: reflect.String}, Kind: reflect.String, Opts: options.FieldOptions{Display: true}},
									{Name: "State", Tag: "display", Value: models.Value{Value: "DN", Kind: reflect.String}, Kind: reflect.String, Opts: options.FieldOptions{Display: true}},
									{Name: "Street", Tag: "", Value: models.Value{Value: "65 Best St", Kind: reflect.String}, Kind: reflect.String, Opts: options.FieldOptions{Display: false}},
									{Name: "Zip", Tag: "", Value: models.Value{Value: "55502", Kind: reflect.String}, Kind: reflect.String, Opts: options.FieldOptions{Display: false}},
								},
							},
							Kind: reflect.Struct,
						},
						Value:     models.Value{Value: "value2", Kind: reflect.String},
						SortValue: `{Denver DN 65 Best St 55502}`,
					},
					{
						Key: models.Value{
							Value: models.Struct{
								Name: "parser.address",
								Fields: []models.Field{
									{Name: "City", Tag: "display", Value: models.Value{Value: "San Francisco", Kind: reflect.String}, Kind: reflect.String, Opts: options.FieldOptions{Display: true}},
									{Name: "State", Tag: "display", Value: models.Value{Value: "CA", Kind: reflect.String}, Kind: reflect.String, Opts: options.FieldOptions{Display: true}},
									{Name: "Street", Tag: "", Value: models.Value{Value: "451 Main St", Kind: reflect.String}, Kind: reflect.String, Opts: options.FieldOptions{Display: false}},
									{Name: "Zip", Tag: "", Value: models.Value{Value: "55501", Kind: reflect.String}, Kind: reflect.String, Opts: options.FieldOptions{Display: false}},
								},
							},
							Kind: reflect.Struct,
						},
						Value:     models.Value{Value: "value1", Kind: reflect.String},
						SortValue: `{San Francisco CA 451 Main St 55501}`,
					},
				},
			}

			require.Equal(t, exp, got)
		})
	})

	t.Run("non_map_value", func(t *testing.T) {
		require.PanicsWithValue(t, "provided value is not a map", func() { p.Map(reflect.ValueOf(5.234)) })
	})
}
