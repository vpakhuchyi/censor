package censor

import (
	"fmt"
	"reflect"

	"github.com/vpakhuchyi/censor/internal/builderpool"
	"github.com/vpakhuchyi/censor/internal/encoder"
)

// Processor is responsible for data encoding according to the specified configuration.
type Processor struct {
	encoder encoder.Encoder
	cfg     Config
}

// Censor pkg contains a global instance of Processor.
// This globalInstance is used by the package-level functions.
var globalInstance = New()

// New returns a new instance of Processor with default configuration.
func New() *Processor {
	return newProcessor(DefaultConfig())
}

// NewWithOpts returns a new instance of Processor, options can be passed to it.
// If no options are passed, the default configuration will be used.
func NewWithOpts(opts ...Option) (*Processor, error) {
	var optCfg OptsConfig
	for _, opt := range opts {
		opt(&optCfg)
	}

	cfg := optCfg.config

	if cfg == nil && optCfg.configPath != "" {
		c, err := ConfigFromFile(optCfg.configPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read the configuration: %w", err)
		}

		cfg = &c
	}

	if cfg == nil {
		c := DefaultConfig()
		cfg = &c
	} else if cfg.General.PrintConfigOnInit {
		// We want to print the configuration only if it differs from the default one
		// and the corresponding flag is set to true.
		fmt.Print(cfg.ToString())
	}

	return newProcessor(*cfg), nil
}

// SetGlobalInstance sets a given Processor as a global instance.
func SetGlobalInstance(p *Processor) {
	globalInstance = p
}

// GetGlobalInstance returns a global instance of Processor.
func GetGlobalInstance() *Processor {
	return globalInstance
}

func newProcessor(cfg Config) *Processor {
	p := &Processor{
		cfg: cfg,
	}

	if cfg.General.OutputFormat == "json" {
		p.encoder = encoder.NewJSONEncoder(cfg.Encoder)
	} else {
		p.encoder = encoder.NewTextEncoder(cfg.Encoder)
	}

	return p
}

// Any accepts an arbitrary value and returns a byte slice representation with sensitive data masked.
//
// The function recursively processes the provided value—handling any type including structs, slices,
// arrays, pointers, maps, interfaces, etc. Every string encountered, regardless of its context, is
// evaluated against the configured regex exclude patterns; if a match is found, the matching segments
// are replaced with a mask.
//
// For struct fields, any field lacking the `censor:"display"` tag is treated as sensitive and will be
// masked. This applies whether the field is a direct member of a struct, an element of a slice, part
// of an interface, or nested within any composite type.
//
// The final output format is determined by the current configuration, yielding either JSON or plain text.
//
// For bug reports or feedback, please contribute to the project at https://github.com/vpakhuchyi/censor.
func Any(val any) []byte {
	return globalInstance.Any(val)
}

// Any returns a byte slice representation of the given value with sensitive data masked.
// It behaves the same as the global Any function — recursively processing and masking values.
func (p *Processor) Any(val any) []byte {
	if val == nil || reflect.TypeOf(val) == nil {
		return []byte("nil")
	}

	b := builderpool.Get()
	defer builderpool.Put(b)

	p.encoder.Encode(b, reflect.ValueOf(val))

	return b.Bytes()
}

// String processes the given string by validating it against the configured regular expressions.
// Any segments matching these patterns are replaced with the mask value, and the resulting string
// is returned as a byte slice.
func String(s string) []byte {
	return globalInstance.String(s)
}

// String returns a byte slice containing the processed version of the input string,
// where segments matching the configured regular expressions are replaced with the mask value.
// It behaves identically to the global String function, using the Processor instance's configuration.
func (p *Processor) String(s string) []byte {
	b := builderpool.Get()
	defer builderpool.Put(b)

	p.encoder.String(b, s)

	return b.Bytes()
}

// Clone returns a new instance of Processor with the same configuration as the original one.
func (p *Processor) Clone() (*Processor, error) {
	return NewWithOpts(WithConfig(&p.cfg))
}

const censorIsNotInitializedMsg = "censor is not initialized"

// PrintConfig prints the configuration of the Processor.
func (p *Processor) PrintConfig() {
	if p == nil {
		fmt.Print(censorIsNotInitializedMsg)
	} else {
		fmt.Print(p.cfg.ToString())
	}
}
