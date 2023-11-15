package parser

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/vpakhuchyi/censor/internal/models"
)

func TestParser_String(t *testing.T) {
	p := Parser{
		useJSONTagName: false,
		CensorFieldTag: DefaultCensorFieldTag,
	}

	t.Run("successful", func(t *testing.T) {
		require.NotPanics(t, func() {
			got := p.String(reflect.ValueOf("I've tried so hard"))
			exp := models.Value{Value: "I've tried so hard", Kind: reflect.String}
			require.Equal(t, exp, got)
		})

	})

	t.Run("non_string_value", func(t *testing.T) {
		require.PanicsWithValue(t, "provided value is not a string", func() { p.String(reflect.ValueOf(5.234)) })
	})
}
