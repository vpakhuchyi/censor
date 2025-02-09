package zerologhandler

import (
	"os"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
	"github.com/vpakhuchyi/censor"
)

func TestWithCensor(t *testing.T) {
	t.Run("should apply option with provided censor instance", func(t *testing.T) {
		// GIVEN
		censorInstance := censor.New()
		h := handler{}

		// WHEN
		got := WithCensor(censorInstance)
		got(&h)

		// THEN
		require.Equal(t, handler{censor: censorInstance}, h)
	})
}

func TestWithZerolog(t *testing.T) {
	t.Run("should apply option with provided zerolog logger instance", func(t *testing.T) {
		// GIVEN
		log := zerolog.New(os.Stdout)
		h := handler{}

		// WHEN
		got := WithZerolog(&log)
		got(&h)

		// THEN
		require.Equal(t, handler{log: &log}, h)
	})
}
