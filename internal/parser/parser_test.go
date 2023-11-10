package parser

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	t.Run("successful", func(t *testing.T) {
		require.EqualValues(t, &Parser{UseJSONTagName: false, CensorFieldTag: DefaultCensorFieldTag}, New())
	})
}
