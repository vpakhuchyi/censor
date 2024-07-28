package sloghandler

import (
	"io"
	"log/slog"

	"github.com/vpakhuchyi/censor"
)

// Option type represents a function that can be used to configure a Slog Handler.
// Users can create instances of Option to customize the behavior of the Slog Handler.
// These options can be applied during the initialization of the Slog Handler to modify
// its configuration.
type Option func(cfg *config)

// WithCensor sets the Censor processor instance for the Slog Handler. If not provided,
// a default Censor processor is used.
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
