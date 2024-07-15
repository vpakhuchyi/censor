package censor

import (
	"fmt"
	"reflect"

	jsoniter "github.com/json-iterator/go"
)

const jsonFormatErrorPattern = `{"error":"%s"}`

// Processor is used to censor any value and format it into a string representation.
type Processor struct {
	cfg       Config
	formatter jsoniter.API
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
	} else if cfg.PrintConfigOnInit {
		// We want to print the configuration only if it differs from the default one
		// and the corresponding flag is set to true.
		fmt.Print(cfg.ToString())
	}

	return newProcessor(*cfg), nil
}

func newProcessor(cfg Config) *Processor {
	formatter := jsoniter.Config{EscapeHTML: true}.Froze()
	formatter.RegisterExtension(newExtension(&cfg))

	p := Processor{
		cfg:       cfg,
		formatter: formatter,
	}

	return &p
}

// PrintConfig prints the configuration of the censor Processor.
func (p *Processor) PrintConfig() {
	fmt.Print(p.cfg.ToString())
}

// Format takes any value and returns a string representation of it masking struct fields by default.
// To override this behaviour, use the `censor:"display"` tag.
// Formatting is done recursively for all nested structs/slices/arrays/pointers/maps/interfaces.
// If a value implements [encoding.TextMarshaler], the result of MarshalText is written.
func (p *Processor) Format(val any) string {
	if val == nil || reflect.TypeOf(val) == nil {
		return "null"
	}

	r, err := p.formatter.Marshal(val)
	if err != nil {
		return fmt.Sprintf(jsonFormatErrorPattern, err.Error())
	}

	return string(r)
}

// SetGlobalInstance sets a given Processor as a global instance.
func SetGlobalInstance(p *Processor) {
	globalInstance = p
}

// GetGlobalInstance returns a global instance of Processor.
func GetGlobalInstance() *Processor {
	return globalInstance
}
