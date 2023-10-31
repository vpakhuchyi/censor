package censor

/*
	This package provides a function to censor any value into a string representation.
	By default, all struct fields are masked (including any kind of nested structs),
	unless the field has the `display` tag.

	Examples can be found here: https://github.com/vpakhuchyi/censor#readme.

	Supported types:

	|----------------------|-----------------------------------------------------------------------|
	|                      | By default, all fields values will be masked.	                       |
	| Struct               | To override this behaviour, use the `censor:"display"` tag.        |
	|                      | All nested fields must be tagged as well.                             |
	|                      | Struct/Slice/Array/Pointer/Map values will be parsed recursively.     |
	|----------------------|-----------------------------------------------------------------------|
	| Slice/Array          | Struct/Slice/Array/Pointer/Map values will be parsed recursively.     |
	|----------------------|-----------------------------------------------------------------------|
	| Pointer              | Struct/Slice/Array/Pointer/Map values will be parsed recursively.     |
	|----------------------|-----------------------------------------------------------------------|
	| Map                  | Struct/Slice/Array/Pointer/Map values will be parsed recursively.     |
	|----------------------|-----------------------------------------------------------------------|
	| Interface            | Will be formatted using the same rules as for the value it contains.  |
	|----------------------|-----------------------------------------------------------------------|
	| String               | Default fmt package formatting is used.                               |
	|----------------------|-----------------------------------------------------------------------|
	| Float64              | Formatted value will have up to 15 precision digits.                  |
	|----------------------|-----------------------------------------------------------------------|
	| Float32              | Formatted value will have up to 7 precision digits.                   |
	|----------------------|-----------------------------------------------------------------------|
	| Int/Int8/Int16/      |                                                                       |
	| Int32/Int64/Rune     |                                                                       |
	| Uint/Uint8/Uint16    | Default fmt package formatting is used.                               |
	| Uint32/Uint64/Byte   |                                                                       |
	|----------------------|-----------------------------------------------------------------------|
	| Bool                 | Default fmt package formatting is used.                               |
	|----------------------|-----------------------------------------------------------------------|

	Unsupported types:

	|------------|------------|------------|
	| Chan       | Complex64  | Complex128 |
	|------------|------------|------------|
	| Uintptr    | Func       |            |
	|------------|------------|------------|

	Note: if a value of unsupported type is provided, an empty string will be returned.
*/

import (
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
}

// Censor pkg contains a global instance of Processor.
// This globalInstance is used by the package-level functions.
var globalInstance = &Processor{
	formatter: formatter.New(),
	parser:    parser.New(),
}

// New returns a new instance of Processor with default configuration.
func New() *Processor {
	return &Processor{
		formatter: formatter.New(),
		parser:    parser.New(),
	}
}

// Format takes any value and returns a string representation of it masking struct fields by default.
// To override this behaviour, use the `censor:"display"` tag.
// Formatting is done recursively for all nested structs/slices/arrays/pointers/maps/interfaces.
func (p *Processor) Format(val any) string {
	return p.sanitise(val)
}

// SetMaskValue sets a value that will be used to mask struct fields.
// The default value is stored in the formatter.DefaultMaskValue constant.
func (p *Processor) SetMaskValue(maskValue string) {
	p.formatter.MaskValue = maskValue
}

// UseJSONTagName sets whether to use the `json` tag to get the name of the struct field.
// If no `json` tag is present, the name of struct will be an empty string.
// By default, this option is disabled.
func (p *Processor) UseJSONTagName(v bool) {
	p.parser.UseJSONTagName = v
}

// DisplayStructName sets whether to display the name of the struct.
// By default, this option is disabled.
func (p *Processor) DisplayStructName(v bool) {
	p.formatter.DisplayStructName = v
}

// SetFieldTag sets a tag name for censor fields.
// The default value is stored in the parser.DefaultCensorFieldTag constant.
func (p *Processor) SetFieldTag(tag string) {
	p.parser.CensorFieldTag = tag
}

// DisplayMapType sets whether to display map type in the output.
// By default, this option is disabled.
func (p *Processor) DisplayMapType(v bool) {
	p.formatter.DisplayMapType = v
}

/*
	Pkg-level functions that work with the global instance of Processor.
*/

// Format takes any value and returns a string representation of it masking struct fields by default.
// It uses the global instance of Processor.
// To override this behaviour, use the `censor:"display"` tag.
// Formatting is done recursively for all nested structs/slices/arrays/pointers/maps/interfaces.
func Format(val any) string {
	return globalInstance.sanitise(val)
}

// SetGlobalInstance sets a given Processor as a global instance.
func SetGlobalInstance(p *Processor) {
	globalInstance = p
}

// GetGlobalInstance returns a global instance of Processor.
func GetGlobalInstance() *Processor {
	return globalInstance
}

// SetMaskValue sets a value that will be used to mask struct fields.
// It applies this change to the global instance of Processor.
// The default value is stored in the formatter.DefaultMaskValue constant.
func SetMaskValue(maskValue string) {
	globalInstance.formatter.MaskValue = maskValue
}

// UseJSONTagName sets whether to use the `json` tag to get the name of the struct field.
// It applies this change to the global instance of Processor.
// If no `json` tag is present, the name of struct will be an empty string.
// By default, this option is disabled.
func UseJSONTagName(v bool) {
	globalInstance.parser.UseJSONTagName = v
}

// DisplayStructName sets whether to display the name of the struct.
// It applies this change to the global instance of Processor.
// By default, this option is disabled.
func DisplayStructName(v bool) {
	globalInstance.formatter.DisplayStructName = v
}

// SetFieldTag sets a tag name for censor fields.
// It applies this change to the global instance of Processor.
// The default value is stored in the parser.DefaultCensorFieldTag constant.
func SetFieldTag(tag string) {
	globalInstance.parser.CensorFieldTag = tag
}

// DisplayMapType sets whether to display map type in the output.
// It applies this change to the global instance of Processor.
// By default, this option is disabled.
func DisplayMapType(v bool) {
	globalInstance.formatter.DisplayMapType = v
}

func (p *Processor) sanitise(val any) string {
	if reflect.TypeOf(val) == nil {
		return "nil"
	}

	v := reflect.ValueOf(val)

	// Handle a case when provided value is of one of non-supported types.
	// In such a case, we return an empty string.
	parsed := p.parse(v)
	if v, ok := parsed.(string); ok && v == "" {
		return v
	}

	return p.format(v.Kind(), parsed)
}

//nolint:exhaustive
func (p *Processor) parse(v reflect.Value) any {
	var parsed any
	switch v.Kind() {
	case reflect.Struct:
		parsed = p.parser.Struct(v)
	case reflect.Slice, reflect.Array:
		parsed = p.parser.Slice(v)
	case reflect.Ptr:
		parsed = p.parser.Ptr(v)
	case reflect.Map:
		parsed = p.parser.Map(v)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64, reflect.String, reflect.Bool:
		parsed = models.Value{Value: v.Interface(), Kind: v.Kind()}
	case reflect.Chan, reflect.Func, reflect.UnsafePointer,
		reflect.Complex64, reflect.Complex128, reflect.Uintptr:
		/*
			Note: this case covers all unsupported types.
			In such a case, we return an empty string.
		*/
		return ""
	}

	return parsed
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
