package sloghandler

import (
	"io"
	"log/slog"

	"github.com/vpakhuchyi/censor"
)

// Option represents a set of options for configuring the Slog Handler along with a censor processor.
// It is used as arguments for the NewJSONHandler function. See NewJSONHandler for more details.
type Option func(options *config)

// WithCensor sets the censor processor instance for the Slog Handler. If not provided,
// a default censor processor is used.
func WithCensor(censor *censor.Processor) Option {
	return func(h *config) {
		h.censor = censor
	}
}

// WithOut sets the output destination for the Slog Handler. If not provided, os.Stdout is used.
func WithOut(out io.Writer) Option {
	return func(h *config) {
		h.out = out
	}
}

// WithAddSource enables the addition of source information to log entries.
func WithAddSource() Option {
	return func(h *config) {
		h.AddSource = true
	}
}

// WithLevel sets the log level for the Slog Handler.
func WithLevel(level slog.Leveler) Option {
	return func(h *config) {
		h.Level = level
	}
}

// WithReplaceAttr sets the function for replacing attributes in log entries.
// The provided replaceAttr function should take a slice of attribute groups and an attribute,
// and return a modified attribute.
func WithReplaceAttr(replaceAttr func(groups []string, a slog.Attr) slog.Attr) Option {
	return func(h *config) {
		h.ReplaceAttr = replaceAttr
	}
}
