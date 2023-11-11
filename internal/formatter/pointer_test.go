package formatter

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/vpakhuchyi/censor/internal/models"
)

func TestFormatter_Ptr(t *testing.T) {
	f := Formatter{
		MaskValue:         DefaultMaskValue,
		DisplayStructName: false,
		DisplayMapType:    false,
	}

	t.Run("successful", func(t *testing.T) {
		require.NotPanics(t, func() {
			v := models.Ptr{Value: 1, Kind: reflect.Int}
			got := f.Ptr(v)
			exp := "&1"
			require.Equal(t, exp, got)
		})
	})

	t.Run("nil", func(t *testing.T) {
		require.NotPanics(t, func() {
			v := models.Ptr{Value: nil, Kind: reflect.Ptr}
			got := f.Ptr(v)
			exp := "nil"
			require.Equal(t, exp, got)
		})
	})
}
