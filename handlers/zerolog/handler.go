package zerologhandler

import (
	"io"

	"github.com/vpakhuchyi/censor"
)

// Handler contain all configuration
type Handler struct {
	censor *censor.Processor
}

// New returns a new configured Handler.
// Options can be provided to configure the Censor processor. If no options are provided,
// a default configuration is used. See the Option documentation for more details.
// By default, the censoring of log fields only is enabled.
func New(opts ...Option) Handler {
	var cfg Handler
	for _, o := range opts {
		o(&cfg)
	}

	if cfg.censor == nil {
		cfg.censor = censor.New()
	}

	return cfg
}

// NewInterfaceMarshal is alternative implementation of zerolog.InterfaceMarshalFunc.
// Should be used for correct work with objects.
// example: `zerolog.InterfaceMarshalFunc = handler.NewInterfaceMarshal`.
func (h Handler) NewInterfaceMarshal(v interface{}) ([]byte, error) {
	return []byte(h.censor.Format(v)), nil
}

// Writer is alternative decorator for zerolog output.
// This implementation replaces all sensitive information in already generated msg.
// example: `l := zerolog.New(handler.Writer(&b))`.
func (h Handler) Writer(out io.Writer) io.Writer {
	return writer{out: out, handler: h}
}
