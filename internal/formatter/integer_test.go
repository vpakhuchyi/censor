package formatter

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/vpakhuchyi/censor/internal/models"
)

func TestFormatter_Integer(t *testing.T) {
	f := Formatter{
		maskValue:         DefaultMaskValue,
		displayStructName: false,
		displayMapType:    false,
	}

	t.Run("int", func(t *testing.T) {
		require.NotPanics(t, func() {
			v := models.Value{Value: int(44), Kind: reflect.Int}
			got := f.Integer(v)
			exp := "44"
			require.Equal(t, exp, got)
		})
	})

	t.Run("int8", func(t *testing.T) {
		require.NotPanics(t, func() {
			v := models.Value{Value: int8(44), Kind: reflect.Int8}
			got := f.Integer(v)
			exp := "44"
			require.Equal(t, exp, got)
		})
	})

	t.Run("int16", func(t *testing.T) {
		require.NotPanics(t, func() {
			v := models.Value{Value: int16(44), Kind: reflect.Int16}
			got := f.Integer(v)
			exp := "44"
			require.Equal(t, exp, got)
		})
	})

	t.Run("int32", func(t *testing.T) {
		require.NotPanics(t, func() {
			v := models.Value{Value: int32(44), Kind: reflect.Int32}
			got := f.Integer(v)
			exp := "44"
			require.Equal(t, exp, got)
		})
	})

	t.Run("int64", func(t *testing.T) {
		require.NotPanics(t, func() {
			v := models.Value{Value: int64(44), Kind: reflect.Int64}
			got := f.Integer(v)
			exp := "44"
			require.Equal(t, exp, got)
		})
	})

	t.Run("uint", func(t *testing.T) {
		require.NotPanics(t, func() {
			v := models.Value{Value: uint(44), Kind: reflect.Uint}
			got := f.Integer(v)
			exp := "44"
			require.Equal(t, exp, got)
		})
	})

	t.Run("uint8", func(t *testing.T) {
		require.NotPanics(t, func() {
			v := models.Value{Value: uint8(44), Kind: reflect.Uint8}
			got := f.Integer(v)
			exp := "44"
			require.Equal(t, exp, got)
		})
	})

	t.Run("uint16", func(t *testing.T) {
		require.NotPanics(t, func() {
			v := models.Value{Value: uint16(44), Kind: reflect.Uint16}
			got := f.Integer(v)
			exp := "44"
			require.Equal(t, exp, got)
		})
	})

	t.Run("uint32", func(t *testing.T) {
		require.NotPanics(t, func() {
			v := models.Value{Value: uint32(44), Kind: reflect.Uint32}
			got := f.Integer(v)
			exp := "44"
			require.Equal(t, exp, got)
		})
	})

	t.Run("uint64", func(t *testing.T) {
		require.NotPanics(t, func() {
			v := models.Value{Value: uint64(44), Kind: reflect.Uint64}
			got := f.Integer(v)
			exp := "44"
			require.Equal(t, exp, got)
		})
	})

	t.Run("non_integer_value", func(t *testing.T) {
		require.PanicsWithValue(t, "provided value is not an integer", func() {
			f.Integer(models.Value{Value: 44.34, Kind: reflect.Float32})
		})
	})
}
