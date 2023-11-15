package parser

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/vpakhuchyi/censor/internal/models"
)

func TestParser_Float(t *testing.T) {
	p := Parser{
		useJSONTagName: false,
		CensorFieldTag: DefaultCensorFieldTag,
	}

	t.Run("successful_float32", func(t *testing.T) {
		require.NotPanics(t, func() {
			got := p.Float(reflect.ValueOf(float32(5.234)))
			exp := models.Value{Value: float32(5.234), Kind: reflect.Float32}
			require.Equal(t, exp, got)
		})
	})

	t.Run("successful_float64", func(t *testing.T) {
		require.NotPanics(t, func() {
			got := p.Float(reflect.ValueOf(float64(46546235.2342342)))
			exp := models.Value{Value: float64(46546235.2342342), Kind: reflect.Float64}
			require.Equal(t, exp, got)
		})
	})

	t.Run("non_float_value", func(t *testing.T) {
		require.PanicsWithValue(t, "provided value is not a float", func() { p.Float(reflect.ValueOf("true")) })
	})
}
