package formatter

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/vpakhuchyi/censor/internal/models"
)

func TestFormatter_Complex(t *testing.T) {
	f := Formatter{
		maskValue:         DefaultMaskValue,
		displayStructName: false,
		displayMapType:    false,
	}

	t.Run("complex64", func(t *testing.T) {
		require.NotPanics(t, func() {
			v := models.Value{Value: complex(float32(3.11111111111111), float32(-2345.2352325353)), Kind: reflect.Complex64}
			got := f.Complex(v)
			exp := "(3.111111-2345.235i)"
			require.Equal(t, exp, got)
		})
	})

	t.Run("complex128", func(t *testing.T) {
		require.NotPanics(t, func() {
			v := models.Value{Value: complex(float64(343.11111111111111), float64(-2345.9992352325353)), Kind: reflect.Complex128}
			got := f.Complex(v)
			exp := "(343.111111111111-2345.99923523254i)"
			require.Equal(t, exp, got)
		})
	})

	t.Run("non_complex_value", func(t *testing.T) {
		require.PanicsWithValue(t, "provided value is not a complex", func() {
			f.Complex(models.Value{Value: 44, Kind: reflect.Int})
		})
	})
}
