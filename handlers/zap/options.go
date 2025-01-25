package zaphandler

import (
	"github.com/vpakhuchyi/censor"
)

// Option type represents a function that can be used to configure a Zap Handler.
// These options can be applied during the initialization of the Zap Handler to modify
// its configuration.
type Option func(h *handler)

// WithCensor sets the Censor processor instance for the Zap Handler. If not provided,
// a default Censor processor is used.
func WithCensor(censor *censor.Processor) Option {
	return func(h *handler) {
		h.censor = censor
	}
}
