package formatter

import (
	"reflect"
	"strings"

	"github.com/vpakhuchyi/censor/config"
	"github.com/vpakhuchyi/censor/internal/models"
)

// Formatter is used to format values.
type Formatter struct {
	// maskValue is used to mask struct fields with sensitive data.
	// The default value is stored in config.DefaultMaskValue constant.
	maskValue string
	// displayStructName is used to hide struct name in the output.
	// The default value is false.
	displayStructName bool
	// displayMapType is used to display map type in the output.
	// The default value is false.
	displayMapType bool

	stringsExcludePatterns []string
}

// New returns a new instance of Formatter with default configuration.
func New() *Formatter {
	return &Formatter{
		maskValue:         config.DefaultMaskValue,
		displayStructName: false,
		displayMapType:    false,
	}
}

// NewWithConfig returns a new instance of Formatter with given configuration.
func NewWithConfig(c config.Formatter) *Formatter {
	return &Formatter{
		maskValue:              c.MaskValue,
		displayStructName:      c.DisplayStructName,
		displayMapType:         c.DisplayMapType,
		stringsExcludePatterns: c.StringsExcludePatterns,
	}
}

// SetMaskValue sets a value that will be used to mask struct fields.
func (f *Formatter) SetMaskValue(mask string) {
	f.maskValue = mask
}

// DisplayStructName sets whether to display the name of the struct.
func (f *Formatter) DisplayStructName(v bool) {
	f.displayStructName = v
}

// DisplayMapType sets whether to display map type in the output.
func (f *Formatter) DisplayMapType(v bool) {
	f.displayMapType = v
}

//nolint:exhaustive,gocyclo
func (f *Formatter) writeValue(buf *strings.Builder, v models.Value) {
	switch v.Kind {
	case reflect.String:
		buf.WriteString(f.String(v))
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		buf.WriteString(f.Integer(v))
	case reflect.Float32, reflect.Float64:
		buf.WriteString(f.Float(v))
	case reflect.Complex64, reflect.Complex128:
		buf.WriteString(f.Complex(v))
	case reflect.Struct:
		buf.WriteString(f.Struct(v.Value.(models.Struct)))
	case reflect.Slice, reflect.Array:
		buf.WriteString(f.Slice(v.Value.(models.Slice)))
	case reflect.Pointer:
		buf.WriteString(f.Ptr(v.Value.(models.Ptr)))
	case reflect.Map:
		buf.WriteString(f.Map(v.Value.(models.Map)))
	case reflect.Bool:
		buf.WriteString(f.Bool(v))
	case reflect.Interface:
		buf.WriteString(f.Interface(v))
	}
}

//nolint:exhaustive,gocyclo
func (f *Formatter) writeField(field models.Field, buf *strings.Builder) {
	switch field.Kind {
	case reflect.String:
		buf.WriteString(formatField(field.Name, f.String(field.Value)))
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		buf.WriteString(formatField(field.Name, f.Integer(field.Value)))
	case reflect.Float32, reflect.Float64:
		buf.WriteString(formatField(field.Name, f.Float(field.Value)))
	case reflect.Complex64, reflect.Complex128:
		buf.WriteString(formatField(field.Name, f.Complex(field.Value)))
	case reflect.Struct:
		buf.WriteString(formatField(field.Name, f.Struct(field.Value.Value.(models.Struct))))
	case reflect.Slice, reflect.Array:
		buf.WriteString(formatField(field.Name, f.Slice(field.Value.Value.(models.Slice))))
	case reflect.Pointer:
		buf.WriteString(formatField(field.Name, f.Ptr(field.Value.Value.(models.Ptr))))
	case reflect.Bool:
		buf.WriteString(formatField(field.Name, f.Bool(field.Value)))
	case reflect.Map:
		buf.WriteString(formatField(field.Name, f.Map(field.Value.Value.(models.Map))))
	case reflect.Interface:
		buf.WriteString(formatField(field.Name, f.Interface(field.Value)))
	}
}
