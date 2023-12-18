package sloghandler

import (
	"bytes"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/vpakhuchyi/censor"
)

func TestWithCensor(t *testing.T) {
	t.Run("apply_option_with_censor_instance", func(t *testing.T) {
		// GIVEN a censor instance and a handler options config instance.
		censorInstance := censor.New()
		cfg := config{}

		// WHEN the WithCensor option is applied to the config instance.
		got := WithCensor(censorInstance)
		got(&cfg)

		// THEN the censor instance is set in the config instance.
		require.Equal(t, config{censor: censorInstance}, cfg)
	})
}

func TestWithOut(t *testing.T) {
	t.Run("apply_option_with_out", func(t *testing.T) {
		// GIVEN an out and a handler options config instance.
		out := bytes.NewBuffer(nil)
		cfg := config{}

		// WHEN the WithOut option is applied to the config instance.
		got := WithOut(out)
		got(&cfg)

		// THEN the out is set in the config instance.
		require.Equal(t, config{out: out}, cfg)
	})
}

func TestWithAddSource(t *testing.T) {
	t.Run("apply_option_with_add_source", func(t *testing.T) {
		// GIVEN a handler options config instance.
		cfg := config{}

		// WHEN the WithAddSource option is applied to the config instance.
		got := WithAddSource()
		got(&cfg)

		// THEN the AddSource is set to true in the config instance.
		require.Equal(t, config{HandlerOptions: slog.HandlerOptions{AddSource: true}}, cfg)
	})
}

func TestWithLevel(t *testing.T) {
	t.Run("apply_option_with_level", func(t *testing.T) {
		// GIVEN a level and a handler options config instance.
		level := slog.LevelDebug
		cfg := config{}

		// WHEN the WithLevel option is applied to the config instance.
		got := WithLevel(level)
		got(&cfg)

		// THEN the level is set in the config instance.
		require.Equal(t, config{HandlerOptions: slog.HandlerOptions{Level: level}}, cfg)
	})
}

func TestWithReplaceAttr(t *testing.T) {
	t.Run("apply_option_with_replace_attr", func(t *testing.T) {
		// GIVEN a replaceAttr function and a handler options config instance.
		replaceAttr := func(groups []string, a slog.Attr) slog.Attr {
			return a
		}
		cfg := config{}

		// WHEN the WithReplaceAttr option is applied to the config instance.
		got := WithReplaceAttr(replaceAttr)
		got(&cfg)

		// THEN the replaceAttr function is set in the config instance and is used to replace attributes.
		atrr := slog.Any("key", "value")
		require.Equal(t, replaceAttr(nil, atrr), cfg.ReplaceAttr(nil, atrr))
	})
}
