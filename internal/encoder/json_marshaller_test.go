package encoder

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

type fakeJSONMarshaler struct{}

func (f fakeJSONMarshaler) MarshalJSON() ([]byte, error) {
	return nil, errors.New("test error")
}

func TestJSONMarshaler(t *testing.T) {
	t.Run("invalid value kind", func(t *testing.T) {
		require.NotPanics(t, func() {
			// GIVEN.
			v := fakeJSONMarshaler{}

			// WHEN.
			got := PrepareJSONMarshalerValue(v)

			// THEN.
			require.Equal(t, "!ERROR:test error", got)
		})
	})
}
