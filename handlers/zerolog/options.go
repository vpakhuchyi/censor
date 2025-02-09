package zerologhandler

import (
	"github.com/rs/zerolog"
	"github.com/vpakhuchyi/censor"
)

// Option represents a function that can be used to configure a Handler.
// These options can be applied during the initialization of the Handler to modify
// its configuration.
type Option func(h *handler)

// WithCensor sets the Censor processor instance for the Handler. If not provided,
// a default Censor processor will be used.
func WithCensor(censor *censor.Processor) Option {
	return func(h *handler) {
		h.censor = censor
	}
}

// WithZerolog sets the Zerolog logger instance for the Handler. If not provided,
// a default Zerolog logger will be used.
func WithZerolog(log *zerolog.Logger) Option {
	return func(h *handler) {
		h.log = log
	}
}
