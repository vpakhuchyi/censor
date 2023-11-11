package formatter

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/vpakhuchyi/censor/internal/models"
)

func TestFormatter_Slice(t *testing.T) {
	f := Formatter{
		MaskValue:         DefaultMaskValue,
		DisplayStructName: false,
		DisplayMapType:    false,
	}

	t.Run("successful", func(t *testing.T) {
		require.NotPanics(t, func() {
			v := models.Slice{
				Values: []models.Value{
					{Value: 1, Kind: reflect.Int},
					{Value: 2, Kind: reflect.Int},
				},
			}
			got := f.Slice(v)
			exp := "[1, 2]"
			require.Equal(t, exp, got)
		})
	})

	t.Run("empty", func(t *testing.T) {
		require.NotPanics(t, func() {
			v := models.Slice{
				Values: []models.Value{},
			}
			got := f.Slice(v)
			exp := "[]"
			require.Equal(t, exp, got)
		})
	})
}
