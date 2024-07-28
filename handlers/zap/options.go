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

// WithMessagesFormat enables the censoring of a log message if such is present.
// It's useful when the log message is built using the following patterns:
//
//	log.With(keysAndValues).Info(msg)
//	log.Infow(msg, keysAndValues)
//
// This option allows censoring of the msg value.
func WithMessagesFormat() Option {
	return func(h *handler) {
		h.formatMessages = true
	}
}

// WithKeysFormat enables the censoring of log keys.
// Keys are recognized using the zap logic. Here are some examples:
//
//	log.Infow("msg", "key", "val")
//	log.With("key", "val").Info("msg")
//
// For more details, see the zap documentation, please.
func WithKeysFormat() Option {
	return func(h *handler) {
		h.formatKeys = true
	}
}
