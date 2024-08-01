package zerologhandler

import "github.com/vpakhuchyi/censor"

// Option type represents a function that can be used to configure a zerolog Handler.
// These options can be applied during the initialization of the zerolog Handler to modify
// its configuration.
type Option func(h *Handler)

// WithCensor sets the Censor processor instance for the zerolog Handler. If not provided,
// a default Censor processor is used.
func WithCensor(censor *censor.Processor) Option {
	return func(h *Handler) {
		h.censor = censor
	}
}
