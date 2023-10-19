package formatters

import (
	"fmt"
	"reflect"
	"strings"

	"sanitiser/internal/models"
)

const masked = "[******]"

// FormatStruct formats a struct into a string with masked sensitive fields.
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
func FormatStruct(s models.Struct) string {
	var buf strings.Builder
	buf.WriteString("{")

	fields := s.Fields

	for i := 0; i < len(s.Fields); i++ {
		f := fields[i]

		if f.Opts.Display {
			switch f.Kind {
			case reflect.String:
				buf.WriteString(formatStringField(f.Name, f.Value.Value))
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
				reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				buf.WriteString(formatIntField(f.Name, f.Value.Value))
			case reflect.Float32, reflect.Float64:
				buf.WriteString(formatFloatField(f.Name, f.Value))
			case reflect.Struct:
				buf.WriteString(fmt.Sprintf(`"%s": %s`, f.Name, FormatStruct(f.Value.Value.(models.Struct))))
			case reflect.Slice, reflect.Array:
				buf.WriteString(fmt.Sprintf(`"%s": %s`, f.Name, FormatSlice(f.Value.Value.(models.Slice))))
			case reflect.Pointer:
				buf.WriteString(fmt.Sprintf(`"%s": %s`, f.Name, FormatPtr(f.Value.Value.(models.Ptr))))
			case reflect.Bool:
				buf.WriteString(fmt.Sprintf(`"%s": %v`, f.Name, f.Value.Value))
			}
		} else {
			buf.WriteString(fmt.Sprintf(`"%s": "%s"`, f.Name, masked))
		}

		if i < len(fields)-1 {
			buf.WriteString(", ")
		}
	}

	buf.WriteString("}")

	return buf.String()
}

// FormatSlice formats a slice or an array into a string.
// If the slice contains structs, they are formatted recursively using FormatStruct function rules.
func FormatSlice(slice models.Slice) string {
	var buf strings.Builder
	buf.WriteString("[")

	values := slice.Values
	for i := 0; i < len(values); i++ {
		switch values[i].Kind {
		case reflect.String:
			buf.WriteString(FormatString(values[i]))
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			buf.WriteString(FormatInteger(values[i]))
		case reflect.Float32, reflect.Float64:
			buf.WriteString(FormatFloat(values[i]))
		case reflect.Struct:
			buf.WriteString(FormatStruct(values[i].Value.(models.Struct)))
		case reflect.Slice, reflect.Array:
			buf.WriteString(FormatSlice(values[i].Value.(models.Slice)))
		case reflect.Pointer:
			buf.WriteString(FormatPtr(values[i].Value.(models.Ptr)))
		default:
			buf.WriteString(fmt.Sprintf(`%v`, values[i].Value))
		}

		if i < len(values)-1 {
			buf.WriteString(", ")
		}
	}

	buf.WriteString("]")

	return buf.String()
}

// FormatPtr formats a pointer into a string.
// It adds the `&` prefix to the formatted value to indicate that it is a pointer.
// If the pointer points to a struct, it is formatted recursively using FormatStruct function rules.
// If the pointer points to a slice or an array, it is formatted recursively using FormatSlice function rules.
func FormatPtr(ptr models.Ptr) string {
	if ptr.Value.Value == nil {
		return "nil"
	}

	var val string
	switch ptr.Value.Kind {
	case reflect.Struct:
		val = FormatStruct(ptr.Value.Value.(models.Struct))
	case reflect.Slice, reflect.Array:
		val = FormatSlice(ptr.Value.Value.(models.Slice))
	case reflect.String:
		val = fmt.Sprintf(`"%s"`, ptr.Value.Value)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		val = fmt.Sprintf(`%d`, ptr.Value.Value)
	case reflect.Float32, reflect.Float64:
		val = fmt.Sprintf(`%f`, ptr.Value.Value)
	case reflect.Ptr:
		val = FormatPtr(ptr.Value.Value.(models.Ptr))
	case reflect.Bool:
		val = fmt.Sprintf(`%v`, ptr.Value.Value)
	}

	return "&" + val
}

// FormatString formats a value as a string.
// The value is wrapped in double quotes.
func FormatString(v models.Value) string {
	return fmt.Sprintf(`"%s"`, v.Value)
}

// FormatInteger formats a value as an integer.
func FormatInteger(v models.Value) string {
	return fmt.Sprintf(`%d`, v.Value)
}

// FormatFloat formats a value as a float.
// The value is formatted with up to 7 decimal places for float32 and up to 15 decimal places for float64.
func FormatFloat(v models.Value) string {
	if v.Kind == reflect.Float32 {
		return fmt.Sprintf(`%.7g`, v.Value)
	}

	return fmt.Sprintf(`%.15g`, v.Value)
}

// FormatSimple formats a value as a string using the default fmt.Sprintf rules.
func FormatSimple(v models.Value) string {
	return fmt.Sprintf(`%v`, v.Value)
}

// FormatBool formats a value as a boolean.
func FormatBool(v models.Value) string {
	return fmt.Sprintf(`%v`, v.Value)
}

func formatIntField(name string, value any) string {
	return fmt.Sprintf(`"%s": %d`, name, value)
}

func formatFloatField(name string, value models.Value) string {
	return fmt.Sprintf(`"%s": %s`, name, FormatFloat(value))
}

func formatStringField(name string, value any) string {
	return fmt.Sprintf(`"%s": "%s"`, name, value)
}
