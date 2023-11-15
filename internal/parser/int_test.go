package parser

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/vpakhuchyi/censor/internal/models"
)

func TestParser_Integer(t *testing.T) {
	p := Parser{
		useJSONTagName: false,
		CensorFieldTag: DefaultCensorFieldTag,
	}

	t.Run("successful_int", func(t *testing.T) {
		require.NotPanics(t, func() {
			var v int = -2435435
			got := p.Integer(reflect.ValueOf(v))
			exp := models.Value{Value: -2435435, Kind: reflect.Int}
			require.Equal(t, exp, got)
		})
	})

	t.Run("successful_int8", func(t *testing.T) {
		require.NotPanics(t, func() {
			var v int8 = 123
			got := p.Integer(reflect.ValueOf(v))
			exp := models.Value{Value: int8(123), Kind: reflect.Int8}
			require.Equal(t, exp, got)
		})
	})

	t.Run("successful_int16", func(t *testing.T) {
		require.NotPanics(t, func() {
			var v int16 = -12334
			got := p.Integer(reflect.ValueOf(v))
			exp := models.Value{Value: int16(-12334), Kind: reflect.Int16}
			require.Equal(t, exp, got)
		})
	})

	t.Run("successful_int32", func(t *testing.T) {
		require.NotPanics(t, func() {
			var v int32 = -1236934
			got := p.Integer(reflect.ValueOf(v))
			exp := models.Value{Value: int32(-1236934), Kind: reflect.Int32}
			require.Equal(t, exp, got)
		})
	})

	t.Run("successful_int64", func(t *testing.T) {
		require.NotPanics(t, func() {
			var v int64 = -91236934
			got := p.Integer(reflect.ValueOf(v))
			exp := models.Value{Value: int64(-91236934), Kind: reflect.Int64}
			require.Equal(t, exp, got)
		})
	})

	t.Run("successful_uint", func(t *testing.T) {
		require.NotPanics(t, func() {
			var v uint = 2435435
			got := p.Integer(reflect.ValueOf(v))
			exp := models.Value{Value: uint(2435435), Kind: reflect.Uint}
			require.Equal(t, exp, got)
		})
	})

	t.Run("successful_uint8", func(t *testing.T) {
		require.NotPanics(t, func() {
			var v uint8 = 122
			got := p.Integer(reflect.ValueOf(v))
			exp := models.Value{Value: uint8(122), Kind: reflect.Uint8}
			require.Equal(t, exp, got)
		})
	})

	t.Run("successful_uint16", func(t *testing.T) {
		require.NotPanics(t, func() {
			var v uint16 = 1322
			got := p.Integer(reflect.ValueOf(v))
			exp := models.Value{Value: uint16(1322), Kind: reflect.Uint16}
			require.Equal(t, exp, got)
		})
	})

	t.Run("successful_uint32", func(t *testing.T) {
		require.NotPanics(t, func() {
			var v uint32 = 551322
			got := p.Integer(reflect.ValueOf(v))
			exp := models.Value{Value: uint32(551322), Kind: reflect.Uint32}
			require.Equal(t, exp, got)
		})
	})

	t.Run("successful_uint64", func(t *testing.T) {
		require.NotPanics(t, func() {
			var v uint64 = 551322
			got := p.Integer(reflect.ValueOf(v))
			exp := models.Value{Value: uint64(551322), Kind: reflect.Uint64}
			require.Equal(t, exp, got)
		})
	})

	t.Run("successful_byte", func(t *testing.T) {
		require.NotPanics(t, func() {
			var v byte = 122
			got := p.Integer(reflect.ValueOf(v))
			exp := models.Value{Value: uint8(122), Kind: reflect.Uint8}
			require.Equal(t, exp, got)
		})
	})

	t.Run("successful_rune", func(t *testing.T) {
		require.NotPanics(t, func() {
			var v rune = 't'
			got := p.Integer(reflect.ValueOf(v))
			exp := models.Value{Value: int32(116), Kind: reflect.Int32}
			require.Equal(t, exp, got)
		})
	})

	t.Run("non_integer_value", func(t *testing.T) {
		require.PanicsWithValue(t, "provided value is not an integer", func() { p.Integer(reflect.ValueOf("hello")) })
	})
}
