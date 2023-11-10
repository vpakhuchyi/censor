package formatter

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/vpakhuchyi/censor/internal/models"
)

func TestFormatter_Interface(t *testing.T) {
	f := Formatter{
		MaskValue:         DefaultMaskValue,
		DisplayStructName: false,
		DisplayMapType:    false,
	}

	t.Run("non_interface_value", func(t *testing.T) {
		require.PanicsWithValue(t, "provided value is not a string", func() {
			f.String(models.Value{Value: 44, Kind: reflect.Int})
		})
	})
}
