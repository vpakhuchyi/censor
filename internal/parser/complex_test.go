package parser

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/vpakhuchyi/censor/internal/models"
)

func TestParser_Complex(t *testing.T) {
	p := Parser{
		UseJSONTagName: false,
		CensorFieldTag: DefaultCensorFieldTag,
	}

	t.Run("successful_complex64", func(t *testing.T) {
		require.NotPanics(t, func() {
			got := p.Complex(reflect.ValueOf(complex(float32(-45.234), float32(11.933))))
			exp := models.Value{Value: complex(float32(-45.234), float32(11.933)), Kind: reflect.Complex64}
			require.Equal(t, exp, got)
		})
	})

	t.Run("successful_complex128", func(t *testing.T) {
		require.NotPanics(t, func() {
			got := p.Complex(reflect.ValueOf(complex(float64(-445.2366664), float64(121.93767763))))
			exp := models.Value{Value: complex(float64(-445.2366664), float64(121.93767763)), Kind: reflect.Complex128}
			require.Equal(t, exp, got)
		})
	})

	t.Run("non_complex_value", func(t *testing.T) {
		require.PanicsWithValue(t, "provided value is not a complex", func() { p.Complex(reflect.ValueOf("true")) })
	})
}
