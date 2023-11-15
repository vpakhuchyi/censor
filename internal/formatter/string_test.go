package formatter

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/vpakhuchyi/censor/internal/models"
)

func TestFormatter_String(t *testing.T) {
	f := Formatter{
		maskValue:         DefaultMaskValue,
		displayStructName: false,
		displayMapType:    false,
	}

	t.Run("string", func(t *testing.T) {
		require.NotPanics(t, func() {
			v := models.Value{Value: "hello", Kind: reflect.String}
			got := f.String(v)
			exp := "hello"
			require.Equal(t, exp, got)
		})
	})

	t.Run("non_string_value", func(t *testing.T) {
		require.PanicsWithValue(t, "provided value is not a string", func() {
			f.String(models.Value{Value: 44, Kind: reflect.Int})
		})
	})
}
