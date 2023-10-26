package formatter

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/vpakhuchyi/sanitiser/internal/models"
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
}

// New returns a new instance of Formatter with default configuration.
func New() *Formatter {
	return &Formatter{
		MaskValue:         DefaultMaskValue,
		DisplayStructName: false,
	}
}

func (f *Formatter) writeValue(buf *strings.Builder, v models.Value) {
	switch v.Kind {
	case reflect.String:
		buf.WriteString(f.String(v))
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		buf.WriteString(f.Integer(v))
	case reflect.Float32, reflect.Float64:
		buf.WriteString(f.Float(v))
	case reflect.Struct:
		buf.WriteString(f.Struct(v.Value.(models.Struct)))
	case reflect.Slice, reflect.Array:
		buf.WriteString(f.Slice(v.Value.(models.Slice)))
	case reflect.Pointer:
		buf.WriteString(f.Ptr(v.Value.(models.Ptr)))
	case reflect.Map:
		buf.WriteString(f.Map(v.Value.(models.Map)))
	default:
		buf.WriteString(fmt.Sprintf(`%v`, v.Value))
	}
}
