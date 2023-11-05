package parser

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/vpakhuchyi/censor/internal/models"
)

func TestParser_Bool(t *testing.T) {
	p := Parser{
		UseJSONTagName: false,
		CensorFieldTag: DefaultCensorFieldTag,
	}

	t.Run("successful", func(t *testing.T) {
		require.NotPanics(t, func() {
			got := p.Bool(reflect.ValueOf(true))
			exp := models.Value{Value: true, Kind: reflect.Bool}
			require.Equal(t, exp, got)
		})
	})

	t.Run("non_bool_value", func(t *testing.T) {
		require.PanicsWithValue(t, "provided value is not a boolean", func() { p.Bool(reflect.ValueOf("true")) })
	})
}
