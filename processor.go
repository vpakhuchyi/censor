package censor

import (
	"encoding"
	"fmt"
	"reflect"

	"github.com/vpakhuchyi/censor/internal/formatter"
	"github.com/vpakhuchyi/censor/internal/models"
	"github.com/vpakhuchyi/censor/internal/parser"
)

// Processor is used to censor any value and format it into a string representation.
type Processor struct {
	formatter *formatter.Formatter
	parser    *parser.Parser
	cfg       Config
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
	return &Processor{
		cfg:       cfg,
		formatter: formatter.New(cfg.Formatter),
		parser:    parser.New(cfg.Parser),
	}
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

	// If a value implements [encoding.TextMarshaler] interface, then it should be marshaled to string.
	if tm, ok := val.(encoding.TextMarshaler); ok {
		val = parser.PrepareTextMarshalerValue(tm)
	}

	v := reflect.ValueOf(val)

	return p.format(v.Kind(), p.parse(v))
}

// Clone returns a new instance of Processor with the same configuration as the original one.
func (p *Processor) Clone() (*Processor, error) {
	return NewWithOpts(WithConfig(&p.cfg))
}

const censorIsNotInitializedMsg = "censor is not initialized"

// PrintConfig prints the configuration of the censor Processor.
func (p *Processor) PrintConfig() {
	if p == nil {
		fmt.Print(censorIsNotInitializedMsg)
	} else {
		fmt.Print(p.cfg.ToString())
	}
}

//nolint:exhaustive
func (p *Processor) parse(v reflect.Value) any {
	switch k := v.Kind(); k {
	case reflect.Struct:
		return p.parser.Struct(v)
	case reflect.Slice, reflect.Array:
		return p.parser.Slice(v)
	case reflect.Pointer:
		return p.parser.Ptr(v)
	case reflect.Map:
		return p.parser.Map(v)
	case reflect.Float32, reflect.Float64:
		return p.parser.Float(v)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return p.parser.Integer(v)
	case reflect.Bool:
		return p.parser.Bool(v)
	case reflect.String:
		return p.parser.String(v)
	default:
		/*
			Note: this case covers all unsupported types.
			In such a case, we return an empty string.
		*/
		return models.Value{Value: fmt.Sprintf(parser.UnsupportedTypeTmpl, k), Kind: k}
	}
}

//nolint:exhaustive
func (p *Processor) format(k reflect.Kind, v any) string {
	switch k {
	case reflect.Struct:
		return p.formatter.Struct(v.(models.Struct))
	case reflect.Slice, reflect.Array:
		return p.formatter.Slice(v.(models.Slice))
	case reflect.Pointer:
		return p.formatter.Ptr(v.(models.Ptr))
	case reflect.String:
		return p.formatter.String(v.(models.Value))
	case reflect.Float32, reflect.Float64:
		return p.formatter.Float(v.(models.Value))
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return p.formatter.Integer(v.(models.Value))
	case reflect.Bool:
		return p.formatter.Bool(v.(models.Value))
	case reflect.Map:
		return p.formatter.Map(v.(models.Map))
	default:
		return fmt.Sprintf(`%v`, v.(models.Value).Value)
	}
}
