package parser

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/vpakhuchyi/censor/internal/models"
	"github.com/vpakhuchyi/censor/internal/options"
)

func TestParser_Slice(t *testing.T) {
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

	t.Run("slice_of_pointers_to_structs", func(t *testing.T) {
		require.NotPanics(t, func() {
			v := []*address{
				{Street: "451 Main St", City: "San Francisco", State: "CA", Zip: "55501"},
				{Street: "65 Best St", City: "Denver", State: "DN", Zip: "55502"},
			}
			got := p.Slice(reflect.ValueOf(v))
			exp := models.Slice{
				Values: []models.Value{
					{
						Value: models.Ptr{
							Value: models.Struct{
								Name: "parser.address",
								Fields: []models.Field{
									{Name: "City", Value: models.Value{Value: "San Francisco", Kind: reflect.String}, Opts: options.FieldOptions{Display: true}, Kind: reflect.String},
									{Name: "State", Value: models.Value{Value: "CA", Kind: reflect.String}, Opts: options.FieldOptions{Display: true}, Kind: reflect.String},
									{Name: "Street", Value: models.Value{Value: "451 Main St", Kind: reflect.String}, Opts: options.FieldOptions{Display: false}, Kind: reflect.String},
									{Name: "Zip", Value: models.Value{Value: "55501", Kind: reflect.String}, Opts: options.FieldOptions{Display: false}, Kind: reflect.String},
								}},
							Kind: reflect.Struct,
						},
						Kind: reflect.Pointer,
					},
					{
						Value: models.Ptr{
							Value: models.Struct{
								Name: "parser.address",
								Fields: []models.Field{
									{Name: "City", Value: models.Value{Value: "Denver", Kind: reflect.String}, Opts: options.FieldOptions{Display: true}, Kind: reflect.String},
									{Name: "State", Value: models.Value{Value: "DN", Kind: reflect.String}, Opts: options.FieldOptions{Display: true}, Kind: reflect.String},
									{Name: "Street", Value: models.Value{Value: "65 Best St", Kind: reflect.String}, Opts: options.FieldOptions{Display: false}, Kind: reflect.String},
									{Name: "Zip", Value: models.Value{Value: "55502", Kind: reflect.String}, Opts: options.FieldOptions{Display: false}, Kind: reflect.String},
								},
							}, Kind: reflect.Struct,
						}, Kind: reflect.Pointer,
					},
				},
			}

			require.Equal(t, exp, got)
		})
	})

	t.Run("slice_of_structs", func(t *testing.T) {
		require.NotPanics(t, func() {
			v := []address{
				{Street: "451 Main St", City: "San Francisco", State: "CA", Zip: "55501"},
				{Street: "65 Best St", City: "Denver", State: "DN", Zip: "55502"},
			}
			got := p.Slice(reflect.ValueOf(v))
			exp := models.Slice{
				Values: []models.Value{
					{
						Value: models.Struct{
							Name: "parser.address",
							Fields: []models.Field{
								{Name: "City", Value: models.Value{Value: "San Francisco", Kind: reflect.String}, Opts: options.FieldOptions{Display: true}, Kind: reflect.String},
								{Name: "State", Value: models.Value{Value: "CA", Kind: reflect.String}, Opts: options.FieldOptions{Display: true}, Kind: reflect.String},
								{Name: "Street", Value: models.Value{Value: "451 Main St", Kind: reflect.String}, Opts: options.FieldOptions{Display: false}, Kind: reflect.String},
								{Name: "Zip", Value: models.Value{Value: "55501", Kind: reflect.String}, Opts: options.FieldOptions{Display: false}, Kind: reflect.String},
							},
						},
						Kind: reflect.Struct,
					},
					{
						Value: models.Struct{
							Name: "parser.address",
							Fields: []models.Field{
								{Name: "City", Value: models.Value{Value: "Denver", Kind: reflect.String}, Opts: options.FieldOptions{Display: true}, Kind: reflect.String},
								{Name: "State", Value: models.Value{Value: "DN", Kind: reflect.String}, Opts: options.FieldOptions{Display: true}, Kind: reflect.String},
								{Name: "Street", Value: models.Value{Value: "65 Best St", Kind: reflect.String}, Opts: options.FieldOptions{Display: false}, Kind: reflect.String},
								{Name: "Zip", Value: models.Value{Value: "55502", Kind: reflect.String}, Opts: options.FieldOptions{Display: false}, Kind: reflect.String},
							},
						},
						Kind: reflect.Struct,
					},
				},
			}

			require.Equal(t, exp, got)
		})
	})

	t.Run("slices_of_strings", func(t *testing.T) {
		require.NotPanics(t, func() {
			v := []string{"tag1", "tag2"}
			got := p.Slice(reflect.ValueOf(v))
			exp := models.Slice{
				Values: []models.Value{
					{Value: "tag1", Kind: reflect.String},
					{Value: "tag2", Kind: reflect.String},
				},
			}

			require.Equal(t, exp, got)
		})
	})

	t.Run("slice_of_integers", func(t *testing.T) {
		require.NotPanics(t, func() {
			v := []int{1, 2, 3}
			got := p.Slice(reflect.ValueOf(v))
			exp := models.Slice{
				Values: []models.Value{
					{Value: 1, Kind: reflect.Int},
					{Value: 2, Kind: reflect.Int},
					{Value: 3, Kind: reflect.Int},
				},
			}

			require.Equal(t, exp, got)
		})
	})

	t.Run("slice_of_floats", func(t *testing.T) {
		require.NotPanics(t, func() {
			v := []float64{1.1235235245353, -2.2, 3325354.30022}
			got := p.Slice(reflect.ValueOf(v))
			exp := models.Slice{
				Values: []models.Value{
					{Value: 1.1235235245353, Kind: reflect.Float64},
					{Value: -2.2, Kind: reflect.Float64},
					{Value: 3325354.30022, Kind: reflect.Float64},
				},
			}

			require.Equal(t, exp, got)
		})
	})

	t.Run("slice_of_bools", func(t *testing.T) {
		require.NotPanics(t, func() {
			v := []bool{true, false}
			got := p.Slice(reflect.ValueOf(v))
			exp := models.Slice{
				Values: []models.Value{
					{Value: true, Kind: reflect.Bool},
					{Value: false, Kind: reflect.Bool},
				},
			}

			require.Equal(t, exp, got)
		})
	})

	t.Run("empty_slice", func(t *testing.T) {
		require.NotPanics(t, func() {
			v := []string{}
			got := p.Slice(reflect.ValueOf(v))
			exp := models.Slice{Values: nil}

			require.Equal(t, exp, got)
		})
	})

	t.Run("slice_of_slices", func(t *testing.T) {
		require.NotPanics(t, func() {
			v := [][]string{{"tag1", "tag2"}, {"tag3", "tag4"}}
			got := p.Slice(reflect.ValueOf(v))
			exp := models.Slice{
				Values: []models.Value{
					{
						Value: models.Slice{Values: []models.Value{{Value: "tag1", Kind: reflect.String}, {Value: "tag2", Kind: reflect.String}}},
						Kind:  reflect.Slice,
					},
					{
						Value: models.Slice{Values: []models.Value{{Value: "tag3", Kind: reflect.String}, {Value: "tag4", Kind: reflect.String}}},
						Kind:  reflect.Slice,
					},
				},
			}

			require.Equal(t, exp, got)
		})
	})

	t.Run("slice_of_interfaces", func(t *testing.T) {
		require.NotPanics(t, func() {
			v := []interface{}{1, "string", true, 1.1, []string{"tag1", "tag2"}}
			got := p.Slice(reflect.ValueOf(v))
			exp := models.Slice{
				Values: []models.Value{
					{Value: models.Value{Value: 1, Kind: reflect.Int}, Kind: reflect.Interface},
					{Value: models.Value{Value: "string", Kind: reflect.String}, Kind: reflect.Interface},
					{Value: models.Value{Value: true, Kind: reflect.Bool}, Kind: reflect.Interface},
					{Value: models.Value{Value: 1.1, Kind: reflect.Float64}, Kind: reflect.Interface},
					{Value: models.Value{
						Value: models.Slice{
							Values: []models.Value{
								{Value: "tag1", Kind: reflect.String},
								{Value: "tag2", Kind: reflect.String},
							},
						},
						Kind: reflect.Slice,
					},
						Kind: reflect.Interface,
					},
				},
			}

			require.Equal(t, exp, got)
		})
	})

	t.Run("slice_of_complex64", func(t *testing.T) {
		require.NotPanics(t, func() {
			v := []complex64{complex(float32(1), float32(2.2436)), complex(float32(-33241), float32(322.4265436))}
			got := p.Slice(reflect.ValueOf(v))
			exp := models.Slice{
				Values: []models.Value{
					{Value: complex(float32(1), float32(2.2436)), Kind: reflect.Complex64},
					{Value: complex(float32(-33241), float32(322.4265436)), Kind: reflect.Complex64},
				},
			}

			require.Equal(t, exp, got)
		})
	})

	t.Run("slice_of_complex128", func(t *testing.T) {
		require.NotPanics(t, func() {
			v := []complex128{complex(float64(1), float64(2.2436)), complex(float64(-33241), float64(322.4265436))}
			got := p.Slice(reflect.ValueOf(v))
			exp := models.Slice{
				Values: []models.Value{
					{Value: complex(float64(1), float64(2.2436)), Kind: reflect.Complex128},
					{Value: complex(float64(-33241), float64(322.4265436)), Kind: reflect.Complex128},
				},
			}

			require.Equal(t, exp, got)
		})
	})

	t.Run("slice_of_map_string_int", func(t *testing.T) {
		require.NotPanics(t, func() {
			v := []map[string]int{{"tag1": 1, "tag2": 2}, {"tag3": 3, "tag4": 4}}
			got := p.Slice(reflect.ValueOf(v))
			exp := models.Slice{
				Values: []models.Value{
					{
						Value: models.Map{
							Type: "map[string]int",
							Values: []models.KV{
								{SortValue: "tag1", Key: models.Value{Value: "tag1", Kind: reflect.String}, Value: models.Value{Value: 1, Kind: reflect.Int}},
								{SortValue: "tag2", Key: models.Value{Value: "tag2", Kind: reflect.String}, Value: models.Value{Value: 2, Kind: reflect.Int}}}},
						Kind: reflect.Map,
					},
					{
						Value: models.Map{
							Type: "map[string]int",
							Values: []models.KV{
								{SortValue: "tag3", Key: models.Value{Value: "tag3", Kind: reflect.String}, Value: models.Value{Value: 3, Kind: reflect.Int}},
								{SortValue: "tag4", Key: models.Value{Value: "tag4", Kind: reflect.String}, Value: models.Value{Value: 4, Kind: reflect.Int}}}},
						Kind: reflect.Map,
					},
				},
			}

			require.Equal(t, exp, got)
		})
	})

	t.Run("non_slice_value", func(t *testing.T) {
		require.PanicsWithValue(t, "provided value is not a slice/array", func() { p.Slice(reflect.ValueOf(5.234)) })
	})
}
