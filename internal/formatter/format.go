package formatter

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/vpakhuchyi/censor/internal/models"
)

// DefaultMaskValue is used to mask struct fields by default.
const DefaultMaskValue = "[******]"

// Formatter is used to format values.
type Formatter struct {
	// MaskValue is used to mask struct fields with sensitive data.
	// The default value is stored in DefaultMaskValue constant.
	MaskValue string
	// DisplayStructName is used to hide struct name in the output.
	// The default value is false.
	DisplayStructName bool
	// DisplayMapType is used to display map type in the output.
	// The default value is false.
	DisplayMapType bool
}

// New returns a new instance of Formatter with default configuration.
func New() *Formatter {
	return &Formatter{
		MaskValue:         DefaultMaskValue,
		DisplayStructName: false,
		DisplayMapType:    false,
	}
}

//nolint:exhaustive
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
		buf.WriteString(f.Bool(models.Bool(v)))
	default:
		buf.WriteString(fmt.Sprintf(`%v`, v.Value))
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
		buf.WriteString(formatField(field.Name, f.Bool(field.Value.Value.(models.Bool))))
	case reflect.Map:
		buf.WriteString(formatField(field.Name, f.Map(field.Value.Value.(models.Map))))
	case reflect.Interface:
		buf.WriteString(formatField(field.Name, f.Interface(field.Value.Value.(models.Interface))))
	}
}
