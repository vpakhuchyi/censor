package zerologhandler

import (
	"os"

	"github.com/rs/zerolog"

	"github.com/vpakhuchyi/censor"
)

// New creates a new Zerolog with Censor processor.
// It accepts a list of options that can be used to configure the Zerolog Handler.
// If no options are provided, the default Censor processor and Zerolog logger will be used.
//
// Note: Currently, the Zerolog handler supports sanitizing data only for types `any`/`interface{}`.
// All other types and primitives will be logged as-is. This limitation exists because Zerolog
// provides an interface for modifying input data only for the `any`/`interface{}` type.
func New(opts ...Option) *zerolog.Logger {
	h := new(handler)

	for _, o := range opts {
		o(h)
	}

	if h.censor == nil {
		h.censor = censor.New()
	}

	if h.log == nil {
		log := zerolog.New(os.Stdout)
		h.log = &log
	}

	zerolog.InterfaceMarshalFunc = h.anyMarshal

	return h.log
}

type handler struct {
	censor *censor.Processor
	log    *zerolog.Logger
}

func (h *handler) anyMarshal(v interface{}) ([]byte, error) {
	return []byte(h.censor.Any(v)), nil
}
