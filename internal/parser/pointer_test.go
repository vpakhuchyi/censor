package parser

import (
	"reflect"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/require"

	"github.com/vpakhuchyi/censor/internal/models"
	"github.com/vpakhuchyi/censor/internal/options"
)

func TestParser_Pointer(t *testing.T) {
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

	t.Run("pointer_to_struct", func(t *testing.T) {
		require.NotPanics(t, func() {
			v := &address{Street: "451 Main St", City: "San Francisco", State: "CA", Zip: "55501"}

			got := p.Ptr(reflect.ValueOf(v))
			exp := models.Ptr{
				Value: models.Struct{
					Name: "parser.address",
					Fields: []models.Field{
						{Name: "City", Value: models.Value{Value: "San Francisco", Kind: reflect.String}, Opts: options.FieldOptions{Display: true}},
						{Name: "State", Value: models.Value{Value: "CA", Kind: reflect.String}, Opts: options.FieldOptions{Display: true}},
						{Name: "Street", Value: models.Value{Value: "451 Main St", Kind: reflect.String}, Opts: options.FieldOptions{Display: false}},
						{Name: "Zip", Value: models.Value{Value: "55501", Kind: reflect.String}, Opts: options.FieldOptions{Display: false}},
					}},
				Kind: reflect.Struct,
			}

			require.Equal(t, exp, got)
		})
	})

	t.Run("nil_pointer", func(t *testing.T) {
		require.NotPanics(t, func() {
			var v *address
			got := p.Ptr(reflect.ValueOf(v))
			exp := models.Ptr{Value: nil, Kind: reflect.Pointer}

			require.Equal(t, exp, got)
		})
	})

	t.Run("pointer_to_int", func(t *testing.T) {
		require.NotPanics(t, func() {
			v := 13
			got := p.Ptr(reflect.ValueOf(&v))
			exp := models.Ptr{Value: 13, Kind: reflect.Int}

			require.Equal(t, exp, got)
		})
	})

	t.Run("pointer_to_string", func(t *testing.T) {
		require.NotPanics(t, func() {
			v := "borsch"
			got := p.Ptr(reflect.ValueOf(&v))
			exp := models.Ptr{Value: "borsch", Kind: reflect.String}

			require.Equal(t, exp, got)
		})
	})

	t.Run("pointer_to_bool", func(t *testing.T) {
		require.NotPanics(t, func() {
			v := true
			got := p.Ptr(reflect.ValueOf(&v))
			exp := models.Ptr{Value: true, Kind: reflect.Bool}

			require.Equal(t, exp, got)
		})
	})

	t.Run("pointer_to_float32", func(t *testing.T) {
		require.NotPanics(t, func() {
			var v float32 = 1.2459
			got := p.Ptr(reflect.ValueOf(&v))
			exp := models.Ptr{Value: float32(1.2459), Kind: reflect.Float32}

			require.Equal(t, exp, got)
		})
	})

	t.Run("pointer_to_float64", func(t *testing.T) {
		require.NotPanics(t, func() {
			var v float64 = -231.245359
			got := p.Ptr(reflect.ValueOf(&v))
			exp := models.Ptr{Value: float64(-231.245359), Kind: reflect.Float64}

			require.Equal(t, exp, got)
		})
	})

	t.Run("pointer_to_slice", func(t *testing.T) {
		require.NotPanics(t, func() {
			v := &[]string{"tag1", "tag2"}
			got := p.Ptr(reflect.ValueOf(&v))
			exp := models.Ptr{
				Value: models.Ptr{
					Value: models.Slice{
						Values: []models.Value{{Value: "tag1", Kind: reflect.String}, {Value: "tag2", Kind: reflect.String}},
					},
					Kind: reflect.Slice,
				},
				Kind: reflect.Pointer,
			}

			require.Equal(t, exp, got)
		})
	})

	t.Run("pointer_to_array", func(t *testing.T) {
		require.NotPanics(t, func() {
			v := &[2]string{"tag1", "tag2"}
			got := p.Ptr(reflect.ValueOf(&v))
			exp := models.Ptr{
				Value: models.Ptr{
					Value: models.Slice{
						Values: []models.Value{{Value: "tag1", Kind: reflect.String}, {Value: "tag2", Kind: reflect.String}},
					},
					Kind: reflect.Array,
				},
				Kind: reflect.Pointer,
			}

			require.Equal(t, exp, got)
		})
	})

	t.Run("pointer_to_map", func(t *testing.T) {
		require.NotPanics(t, func() {
			v := &map[string]int{"tag1": 1, "tag2": 2}
			got := p.Ptr(reflect.ValueOf(&v))
			exp := models.Ptr{
				Value: models.Ptr{
					Value: models.Map{
						Type: "map[string]int",
						Values: []models.KV{
							{SortValue: "tag1", Key: models.Value{Value: "tag1", Kind: reflect.String}, Value: models.Value{Value: 1, Kind: reflect.Int}},
							{SortValue: "tag2", Key: models.Value{Value: "tag2", Kind: reflect.String}, Value: models.Value{Value: 2, Kind: reflect.Int}},
						},
					},
					Kind: reflect.Map,
				},
				Kind: reflect.Pointer,
			}

			require.Equal(t, exp, got)
		})
	})

	t.Run("pointer_to_interface", func(t *testing.T) {
		require.NotPanics(t, func() {
			v := interface{}("moloko")
			got := p.Ptr(reflect.ValueOf(&v))
			exp := models.Ptr{Value: models.Value{Value: "moloko", Kind: reflect.String}, Kind: reflect.Interface}

			require.Equal(t, exp, got)
		})
	})

	t.Run("pointer_to_complex64", func(t *testing.T) {
		require.NotPanics(t, func() {
			v := complex(float32(1.2459), float32(-345.345234))
			got := p.Ptr(reflect.ValueOf(&v))
			exp := models.Ptr{Value: complex(float32(1.2459), float32(-345.345234)), Kind: reflect.Complex64}

			require.Equal(t, exp, got)
		})
	})

	t.Run("pointer_to_complex128", func(t *testing.T) {
		require.NotPanics(t, func() {
			v := complex(float64(11.24359), float64(-5345.34523768684))
			got := p.Ptr(reflect.ValueOf(&v))
			exp := models.Ptr{Value: complex(float64(11.24359), float64(-5345.34523768684)), Kind: reflect.Complex128}

			require.Equal(t, exp, got)
		})
	})

	t.Run("pointer_to_func", func(t *testing.T) {
		require.NotPanics(t, func() {
			v := func() {}
			got := p.Ptr(reflect.ValueOf(&v))
			exp := models.Ptr{Value: "[Unsupported type: func]", Kind: reflect.Func}

			require.Equal(t, exp, got)
		})
	})

	t.Run("pointer_to_chan", func(t *testing.T) {
		require.NotPanics(t, func() {
			v := make(chan int)
			got := p.Ptr(reflect.ValueOf(&v))
			exp := models.Ptr{Value: "[Unsupported type: chan]", Kind: reflect.Chan}

			require.Equal(t, exp, got)
		})
	})

	t.Run("pointer_to_unsafe_pointer", func(t *testing.T) {
		require.NotPanics(t, func() {
			v := unsafe.Pointer(nil)
			got := p.Ptr(reflect.ValueOf(&v))
			exp := models.Ptr{Value: "[Unsupported type: unsafe.Pointer]", Kind: reflect.UnsafePointer}

			require.Equal(t, exp, got)
		})
	})

	t.Run("pointer_to_uintptr", func(t *testing.T) {
		require.NotPanics(t, func() {
			v := uintptr(0)
			got := p.Ptr(reflect.ValueOf(&v))
			exp := models.Ptr{Value: uintptr(0), Kind: reflect.Uintptr}

			require.Equal(t, exp, got)
		})
	})

	t.Run("non_pointer_value", func(t *testing.T) {
		require.PanicsWithValue(t, "provided value is not a pointer", func() { p.Ptr(reflect.ValueOf(5.234)) })
	})
}
