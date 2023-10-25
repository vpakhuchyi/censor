package sanitiser

/*
	This package provides a function to sanitise any value into a string representation.
	By default, all struct fields are masked (including any kind of nested structs),
	unless the field has the `display` tag.

	Examples can be found here: https://github.com/vpakhuchyi/sanitiser#readme.

	Supported types:

	|-------------------------------------------------------------------------------------------|
	|         		| By default, all fields values will be masked. 		    |
	| Struct		| To override this behaviour, use the `sanitiser:"display"` tag. 	    |
	|			| All nested fields must be tagged as well. 			    |
	|			| Struct/Slice/Array/Pointer/Map values will be parsed recursively  |
	|-------------------------------------------------------------------------------------------|
	| Slice/Array    	| Struct/Slice/Array/Pointer/Map values will be parsed recursively  |
	|-------------------------------------------------------------------------------------------|
	| Pointer          	| Struct/Slice/Array/Pointer/Map values will be parsed recursively  |
	|-------------------------------------------------------------------------------------------|
	| Map              	| Struct/Slice/Array/Pointer/Map values will be parsed recursively  |
	|-------------------------------------------------------------------------------------------|
	| String           	| Default fmt package formatting is used. 		   				|
	|-------------------------------------------------------------------------------------------|
	| Float64          	| Formatted value will have up to 15 precision digits. 		    |
	|-------------------------------------------------------------------------------------------|
	| Float32          	| Formatted value will have up to 7 precision digits. 		    |
	|-------------------------------------------------------------------------------------------|
	| Int/Int8/Int16/  	| 								    |
	| Int32/Int64/Rune 	| 								    |
	| Uint/Uint8/Uint16	| Default fmt package formatting is used. 		    	    |
	| Uint32/Uint64/   	| 								    |
	| Byte             	| 								    |
	|-------------------------------------------------------------------------------------------|
	| Bool             	| Default fmt package formatting is used. 		    			|
	|-------------------------------------------------------------------------------------------|

	Unsupported types:

	|-----------------------------------------------|
	| Chan         	| Complex64   	| Complex128	|
	| Interface    	| Func     	|       	|
	|-----------------------------------------------|

	Note: unsupported types will be replaced with "[unsupported type]" string.
*/

import (
	"fmt"
	"reflect"

	"github.com/vpakhuchyi/sanitiser/internal/formatter"
	"github.com/vpakhuchyi/sanitiser/internal/models"
	"github.com/vpakhuchyi/sanitiser/internal/parser"
)

// unsupportedKind is returned when the value is of unsupported kind.
const unsupportedKind = `[unsupported kind of value: %v]`

// Processor is used to sanitise any value into a string representation.
// Any configuration changes could be applied to the global globalInstance or to the globalInstance of this struct
// using the corresponding methods or package-level functions.
type Processor struct {
	formatter *formatter.Formatter
	parser    *parser.Parser
}

// Sanitiser pkg contains a global globalInstance of Processor.
// This globalInstance is used by the package-level functions.
var globalInstance = &Processor{
	formatter: formatter.New(),
	parser:    parser.New(),
}

// New returns a new globalInstance of Processor with default configuration.
func New() *Processor {
	return &Processor{
		formatter: formatter.New(),
		parser:    parser.New(),
	}
}

// Format takes any value and returns a string representation of it masking struct fields by default.
// To override this behaviour, use the `sanitiser:"display"` tag.
// Formatting is done recursively for all nested structs/slices/arrays/pointers/maps.
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

// HideStructName sets whether to hide the name of the struct.
// By default, this option is disabled.
func (p *Processor) HideStructName(v bool) {
	p.formatter.HideStructName = v
}

// SetFieldTag sets a tag name for sanitiser fields.
// The default value is stored in the parser.DefaultSanitiserFieldTag constant.
func (p *Processor) SetFieldTag(tag string) {
	p.parser.SanitiserFieldTag = tag
}

/*
	Pkg-level functions that work with the global globalInstance of Processor.
*/

// Format takes any value and returns a string representation of it masking struct fields by default.
// It uses the global globalInstance of Processor.
// To override this behaviour, use the `sanitiser:"display"` tag.
// Formatting is done recursively for all nested structs/slices/arrays/pointers/maps.
func Format(val any) string {
	return globalInstance.sanitise(val)
}

// SetGlobalInstance sets a given Processor as a global globalInstance.
func SetGlobalInstance(p *Processor) {
	globalInstance = p
}

// GetGlobalInstance returns a global globalInstance of Processor.
func GetGlobalInstance() *Processor {
	return globalInstance
}

// SetMaskValue sets a value that will be used to mask struct fields.
// It applies this change to the global globalInstance of Processor.
// The default value is stored in the formatter.DefaultMaskValue constant.
func SetMaskValue(maskValue string) {
	globalInstance.formatter.MaskValue = maskValue
}

// UseJSONTagName sets whether to use the `json` tag to get the name of the struct field.
// It applies this change to the global globalInstance of Processor.
// If no `json` tag is present, the name of struct will be an empty string.
// By default, this option is disabled.
func UseJSONTagName(v bool) {
	globalInstance.parser.UseJSONTagName = v
}

// HideStructName sets whether to hide the name of the struct.
// It applies this change to the global globalInstance of Processor.
// By default, this option is disabled.
func HideStructName(v bool) {
	globalInstance.formatter.HideStructName = v
}

// SetFieldTag sets a tag name for sanitiser fields.
// It applies this change to the global globalInstance of Processor.
// The default value is stored in the parser.DefaultSanitiserFieldTag constant.
func SetFieldTag(tag string) {
	globalInstance.parser.SanitiserFieldTag = tag
}

func (p *Processor) sanitise(val any) string {
	v := reflect.ValueOf(val)

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
		reflect.Complex64, reflect.Complex128:
		/*
			Note: this case covers all unsupported types.
			In such a case, we return a string with a message about the unsupported type.
		*/
		return fmt.Sprintf(unsupportedKind, v.Kind())
	}

	switch v.Kind() {
	case reflect.Struct:
		return p.formatter.Struct(parsed.(models.Struct))
	case reflect.Slice, reflect.Array:
		return p.formatter.Slice(parsed.(models.Slice))
	case reflect.Pointer:
		return p.formatter.Ptr(parsed.(models.Ptr))
	case reflect.String:
		return p.formatter.String(parsed.(models.Value))
	case reflect.Float32, reflect.Float64:
		return p.formatter.Float(parsed.(models.Value))
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return p.formatter.Integer(parsed.(models.Value))
	case reflect.Bool:
		return p.formatter.Bool(parsed.(models.Value))
	case reflect.Map:
		return p.formatter.Map(parsed.(models.Map))
	default:
		return fmt.Sprintf(`%v`, parsed.(models.Value).Value)
	}
}
