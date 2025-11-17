package zerologhandler

import (
	"os"
	"sync"

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
var marshalMu sync.Mutex

// MarshalFunc represents a function capable of marshaling values before zerolog encodes them.
type MarshalFunc func(any) ([]byte, error)

// New returns a zerolog.Logger configured with the supplied options.
// Note: New does not install the marshal func; call InstallMarshalFunc to wire it globally.
func New(opts ...Option) *zerolog.Logger {
	return resolveOptions(opts...).logger
}

// GetMarshalFunc returns a MarshalFunc that applies Censor before delegating to zerolog.
func GetMarshalFunc(opts ...Option) MarshalFunc {
	cfg := resolveOptions(opts...)

	return func(v any) ([]byte, error) {
		return cfg.censor.Any(v), nil
	}
}

// InstallMarshalFunc sets fn as zerolog.InterfaceMarshalFunc and returns a restore function.
func InstallMarshalFunc(fn MarshalFunc) func() {
	if fn == nil {
		return func() {}
	}

	marshalMu.Lock()
	previous := zerolog.InterfaceMarshalFunc
	zerolog.InterfaceMarshalFunc = func(v any) ([]byte, error) {
		return fn(v)
	}
	marshalMu.Unlock()

	return func() {
		marshalMu.Lock()
		zerolog.InterfaceMarshalFunc = previous
		marshalMu.Unlock()
	}
}

func resolveOptions(opts ...Option) options {
	cfg := options{}
	for _, opt := range opts {
		opt(&cfg)
	}

	if cfg.censor == nil {
		cfg.censor = censor.New()
	}

	if cfg.censor.OutputFormat() != censor.OutputFormatJSON {
		panic("zerologhandler: censor processor must use json output format")
	}

	if cfg.logger == nil {
		l := zerolog.New(os.Stdout)
		cfg.logger = &l
	}

	return cfg
}
