package encoder

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

type fakeTextMarshaler struct{}

func (f fakeTextMarshaler) MarshalText() ([]byte, error) {
	return nil, errors.New("test error")
}

func TestTextMarshaler(t *testing.T) {
	t.Run("invalid value kind", func(t *testing.T) {
		require.NotPanics(t, func() {
			// GIVEN.
			v := fakeTextMarshaler{}

			// WHEN.
			got := PrepareTextMarshalerValue(v)

			// THEN.
			require.Equal(t, "!ERROR:test error", got)
		})
	})
}
