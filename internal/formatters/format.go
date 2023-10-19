package formatters

import (
	"fmt"
	"reflect"
	"strings"

	"sanitiser/internal/models"
)

const masked = "[******]"

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
				buf.WriteString(formatFloatField(f.Name, f.Value.Value))
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

func FormatString(v models.Value) string {
	return fmt.Sprintf(`"%s"`, v.Value)
}

func FormatInteger(v models.Value) string {
	return fmt.Sprintf(`%d`, v.Value)
}

func FormatFloat(v models.Value) string {
	return fmt.Sprintf(`%f`, v.Value)
}

func FormatSimple(v models.Value) string {
	return fmt.Sprintf(`%v`, v.Value)
}

func formatIntField(name string, value any) string {
	return fmt.Sprintf(`"%s": %d`, name, value)
}

func formatFloatField(name string, value any) string {
	return fmt.Sprintf(`"%s": %f`, name, value)
}

func formatStringField(name string, value any) string {
	return fmt.Sprintf(`"%s": "%s"`, name, value)
}
