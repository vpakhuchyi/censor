package censor

import (
	"encoding"
	"fmt"
	"reflect"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/vpakhuchyi/censor/config"
	"github.com/vpakhuchyi/censor/internal/formatter"
	"github.com/vpakhuchyi/censor/internal/models"
	"github.com/vpakhuchyi/censor/internal/parser"
)

// Processor is used to censor any value and format it into a string representation.
type Processor struct {
	formatter *formatter.Formatter
	parser    *parser.Parser
	cfg       config.Config
}

// Censor pkg contains a global instance of Processor.
// This globalInstance is used by the package-level functions.
var globalInstance = New()

// New returns a new instance of Processor with default configuration.
func New() *Processor {
	c := config.Default()
	p := Processor{
		formatter: formatter.New(c.Formatter),
		parser:    parser.New(c.Parser),
		cfg:       c,
	}

	p.PrintConfig()

	return &p
}

// NewWithConfig returns a new instance of Processor with given configuration.
func NewWithConfig(c config.Config) *Processor {
	p := Processor{
		formatter: formatter.New(c.Formatter),
		parser:    parser.New(c.Parser),
		cfg:       c,
	}

	p.PrintConfig()

	return &p
}

// NewWithFileConfig returns a new instance of Processor with configuration from a given file.
// It returns an error if the file cannot be read or unmarshalled.
func NewWithFileConfig(path string) (*Processor, error) {
	cfg, err := config.FromFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read the configuration: %w", err)
	}

	return NewWithConfig(cfg), nil
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

// PrintConfig prints the configuration of the censor Processor.
func (p *Processor) PrintConfig() {
	const (
		lineLength = 69
		padLength  = 10
	)

	line := strings.Repeat("-", lineLength) + "\n"

	var b strings.Builder
	b.WriteString(line)
	b.WriteString(strings.Repeat(" ", padLength) + "Censor is configured with the following settings:" + "\n")
	b.WriteString(line)

	cfg := config.Config{
		General:   p.cfg.General,
		Parser:    p.cfg.Parser,
		Formatter: p.cfg.Formatter,
	}

	// config.Config and its nested fields contain only supported YAML types, so it must be marshalled successfully.
	// However, if the configuration is changed, it may contain unsupported types. To handle this case,
	// we have tests that check whether the configuration can be marshalled.
	// Because such kind of changes can happen only in the development process and considering that such an issue
	// can't be fixed automatically or by the user, we don't want to fail the application in this case.
	//
	//nolint:errcheck
	d, _ := yaml.Marshal(cfg)

	b.Write(d)

	b.WriteString(line)
	fmt.Print(b.String())
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
	case reflect.Complex64, reflect.Complex128:
		return p.parser.Complex(v)
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
		return models.Value{Value: fmt.Sprintf(config.UnsupportedTypeTmpl, k), Kind: k}
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
	case reflect.Complex64, reflect.Complex128:
		return p.formatter.Complex(v.(models.Value))
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
