package zerologhandler

import (
	"github.com/rs/zerolog"
	"github.com/vpakhuchyi/censor"
)

type options struct {
	censor *censor.Processor
	logger *zerolog.Logger
}

// Option configures the zerolog handler via a shared options struct.
type Option func(cfg *options)

// WithCensor stores the provided Censor processor on the options struct.
func WithCensor(processor *censor.Processor) Option {
	return func(cfg *options) {
		cfg.censor = processor
	}
}

// WithZerolog stores the provided zerolog logger on the options struct.
func WithZerolog(logger *zerolog.Logger) Option {
	return func(cfg *options) {
		cfg.logger = logger
	}
}
