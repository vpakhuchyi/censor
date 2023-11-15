package formatter

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/vpakhuchyi/censor/internal/models"
)

func TestFormatter_Bool(t *testing.T) {
	f := Formatter{
		maskValue:         DefaultMaskValue,
		displayStructName: false,
		displayMapType:    false,
	}

	t.Run("successful", func(t *testing.T) {
		require.NotPanics(t, func() {
			v := models.Value{Value: true, Kind: reflect.Bool}
			got := f.Bool(v)
			exp := "true"
			require.Equal(t, exp, got)
		})
	})

	t.Run("non_bool_value", func(t *testing.T) {
		require.PanicsWithValue(t, "provided value is not a boolean", func() {
			f.Bool(models.Value{Value: 44, Kind: reflect.Int})
		})
	})
}
