package formatter

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/vpakhuchyi/censor/config"
	"github.com/vpakhuchyi/censor/internal/models"
)

func TestFormatter_Float(t *testing.T) {
	f := Formatter{
		maskValue:                    config.DefaultMaskValue,
		displayStructName:            false,
		displayMapType:               false,
		float32MaxSignificantFigures: config.Float32MaxSignificantFigures,
		float64MaxSignificantFigures: config.Float64MaxSignificantFigures,
	}

	t.Run("float32", func(t *testing.T) {
		require.NotPanics(t, func() {
			v := models.Value{Value: float32(3.11111111111111), Kind: reflect.Float32}
			got := f.Float(v)
			exp := "3.111111"
			require.Equal(t, exp, got)
		})
	})

	t.Run("float64", func(t *testing.T) {
		require.NotPanics(t, func() {
			v := models.Value{Value: float64(3.11111111111111), Kind: reflect.Float64}
			got := f.Float(v)
			exp := "3.11111111111111"
			require.Equal(t, exp, got)
		})
	})

	t.Run("non_float_value", func(t *testing.T) {
		require.PanicsWithValue(t, "provided value is not a float", func() {
			f.Float(models.Value{Value: 44, Kind: reflect.Int})
		})
	})
}
