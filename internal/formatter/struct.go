package formatter

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/vpakhuchyi/sanitiser/internal/models"
)

// Struct formats a struct into a string with masked sensitive fields.
// All fields are masked by default, unless the field has the `display` tag.
// Supported types:
//
// [basic types]
// - string
// - int, int8, int16, int32, int64
// - uint, uint8, uint16, uint32, uint64
// - float32, float64
// - bool
//
// [complex types]
// - struct - formatted recursively
// - slice - struct values are formatted recursively
// - array - struct values are formatted recursively
// - pointer - pointed structure/arrays/slices are formatted recursively.
func (f *Formatter) Struct(s models.Struct) string {
	var buf strings.Builder

	if f.HideStructName {
		s.Name = ""
	}

	buf.WriteString(fmt.Sprintf("%s{", s.Name))

	fields := s.Fields
	for i := 0; i < len(s.Fields); i++ {
		field := fields[i]

		if field.Opts.Display {
			switch field.Kind {
			case reflect.String:
				buf.WriteString(f.formatStringField(field.Name, field.Value.Value))
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
				reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				buf.WriteString(f.formatIntField(field.Name, field.Value.Value))
			case reflect.Float32, reflect.Float64:
				buf.WriteString(f.formatFloatField(field.Name, field.Value))
			case reflect.Struct:
				buf.WriteString(fmt.Sprintf(`"%s": %s`, field.Name, f.Struct(field.Value.Value.(models.Struct))))
			case reflect.Slice, reflect.Array:
				buf.WriteString(fmt.Sprintf(`"%s": %s`, field.Name, f.Slice(field.Value.Value.(models.Slice))))
			case reflect.Pointer:
				buf.WriteString(fmt.Sprintf(`"%s": %s`, field.Name, f.Ptr(field.Value.Value.(models.Ptr))))
			case reflect.Bool:
				buf.WriteString(fmt.Sprintf(`"%s": %v`, field.Name, field.Value.Value))
			case reflect.Map:
				buf.WriteString(fmt.Sprintf(`"%s": %s`, field.Name, f.Map(field.Value.Value.(models.Map))))
			}
		} else {
			buf.WriteString(fmt.Sprintf(`"%s": "%s"`, field.Name, f.MaskValue))
		}

		if i < len(fields)-1 {
			buf.WriteString(", ")
		}
	}

	buf.WriteString("}")

	return buf.String()
}

func (f *Formatter) formatIntField(name string, value any) string {
	return fmt.Sprintf(`"%s": %d`, name, value)
}

func (f *Formatter) formatFloatField(name string, value models.Value) string {
	return fmt.Sprintf(`"%s": %s`, name, f.Float(value))
}

func (f *Formatter) formatStringField(name string, value any) string {
	return fmt.Sprintf(`"%s": "%s"`, name, value)
}
