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

/*
	Pkg-level functions that work with the global instance of Processor.
*/

// SetGlobalInstance sets a given Processor as a global instance.
func SetGlobalInstance(p *Processor) {
	globalInstance = p
}

// GetGlobalInstance returns a global instance of Processor.
func GetGlobalInstance() *Processor {
	return globalInstance
}

/*
	Scoped methods that work with a specific instance of Processor.
*/

// Format takes any value and returns a string representation of it masking struct fields by default.
// To override this behaviour, use the `censor:"display"` tag.
// Formatting is done recursively for all nested structs/slices/arrays/pointers/maps/interfaces.
// If a value implements [encoding.TextMarshaler], the result of MarshalText is written.
func (p *Processor) Format(val any) string {
	if val == nil || reflect.TypeOf(val) == nil {
		return "nil"
	}

	return p.encode(val)
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

func (p *Processor) encode(v any) string {
	b := builderpool.Get()

	p.encoder.Encode(b, reflect.ValueOf(v))

	res := b.String()
	builderpool.Put(b)

	return res
}
