package sloghandler

import (
	"encoding/json"
	"io"
	"log/slog"
	"os"

	"github.com/vpakhuchyi/censor"
)

type config struct {
	out         io.Writer
	censor      *censor.Processor
	ReplaceAttr func(groups []string, a slog.Attr) slog.Attr
	slog.HandlerOptions
}

// NewJSONHandler returns a new log/slog JSONHandler along with a Censor processor. Options can be provided to configure
// the Censor processor and the log/slog Handler. If no options are provided, a default configuration is used.
// See the Option documentation for more details.
func NewJSONHandler(opts ...Option) *slog.JSONHandler {
	var cfg config
	for _, o := range opts {
		o(&cfg)
	}

	if cfg.censor == nil {
		censorCfg := censor.DefaultConfig()
		censorCfg.General.OutputFormat = censor.OutputFormatJSON
		censorCfg.General.PrintConfigOnInit = false

		// Error can be discarded here because the default configuration is always successfully parsed.
		cfg.censor, _ = censor.NewWithOpts(censor.WithConfig(&censorCfg)) //nolint:errcheck
	}

	if cfg.out == nil {
		cfg.out = os.Stdout
	}

	cfg.HandlerOptions.ReplaceAttr = func(groups []string, a slog.Attr) slog.Attr {
		attr := a
		if cfg.ReplaceAttr != nil {
			attr = cfg.ReplaceAttr(groups, a)
		}

		switch attr.Key {
		// These attributes are required by log/slog. We don't want to censor them.
		case slog.TimeKey, slog.LevelKey, slog.SourceKey:
			return attr
		default:
			return slog.Any(attr.Key, json.RawMessage(cfg.censor.Format(attr.Value.Any())))
		}
	}

	return slog.NewJSONHandler(cfg.out, &cfg.HandlerOptions)
}
