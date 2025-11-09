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
		cfg := options{}

		// WHEN
		got := WithCensor(censorInstance)
		got(&cfg)

		// THEN
		require.Equal(t, options{censor: censorInstance}, cfg)
	})
}

func TestWithZerolog(t *testing.T) {
	t.Run("should apply option with provided zerolog logger instance", func(t *testing.T) {
		// GIVEN
		log := zerolog.New(os.Stdout)
		cfg := options{}

		// WHEN
		got := WithZerolog(&log)
		got(&cfg)

		// THEN
		require.Equal(t, options{logger: &log}, cfg)
	})
}
