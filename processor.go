package censor

import (
	"fmt"
	"reflect"

	"github.com/vpakhuchyi/censor/config"
	"github.com/vpakhuchyi/censor/internal/formatter"
	"github.com/vpakhuchyi/censor/internal/models"
	"github.com/vpakhuchyi/censor/internal/parser"
)

// Processor is used to censor any value and format it into a string representation.
type Processor struct {
	formatter *formatter.Formatter
	parser    *parser.Parser
}

// Censor pkg contains a global instance of Processor.
// This globalInstance is used by the package-level functions.
var globalInstance = New()

// New returns a new instance of Processor with default configuration.
func New() *Processor {
	return &Processor{
		formatter: formatter.New(),
		parser:    parser.New(),
	}
}

// NewWithConfig returns a new instance of Processor with given configuration.
func NewWithConfig(c config.Config) *Processor {
	return &Processor{
		formatter: formatter.NewWithConfig(c.Formatter),
		parser:    parser.NewWithConfig(c.Parser),
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
func (p *Processor) Format(val any) string {
	if reflect.TypeOf(val) == nil {
		return "nil"
	}

	v := reflect.ValueOf(val)

	return p.format(v.Kind(), p.parse(v))
}

// SetMaskValue sets a value that will be used to mask struct fields.
// The default value is stored in the config.DefaultMaskValue constant.
func (p *Processor) SetMaskValue(maskValue string) {
	p.formatter.SetMaskValue(maskValue)
}

// UseJSONTagName sets whether to use the `json` tag to get the name of the struct field.
// If no `json` tag is present, the name of struct will be an empty string.
// By default, this option is disabled.
func (p *Processor) UseJSONTagName(v bool) {
	p.parser.UseJSONTagName(v)
}

// DisplayPointerSymbol sets whether to display the '&' (pointer symbol) before the pointed value.
// By default, this option is disabled.
func (p *Processor) DisplayPointerSymbol(v bool) {
	p.formatter.DisplayPointerSymbol(v)
}

// DisplayStructName sets whether to display the name of the struct.
// By default, this option is disabled.
func (p *Processor) DisplayStructName(v bool) {
	p.formatter.DisplayStructName(v)
}

// DisplayMapType sets whether to display map type in the output.
// By default, this option is disabled.
func (p *Processor) DisplayMapType(v bool) {
	p.formatter.DisplayMapType(v)
}

// AddExcludePatterns adds regexp patterns that are used for the selection of strings that must be masked.
// Regexp patterns compilation will be triggered automatically after adding new patterns.
// Note: this method may panic if regexp pattern is invalid.
func (p *Processor) AddExcludePatterns(patterns ...string) {
	p.formatter.AddExcludePatterns(patterns...)
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
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
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
		return models.Value{Value: "", Kind: k}
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
