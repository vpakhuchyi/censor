package parsers

import (
	"reflect"

	"sanitiser/internal/models"
	"sanitiser/internal/options"
)

// ParseStruct parses a given value and returns a Struct.
// All fields of pointer/slice/array/struct types will be parsed recursively.
func ParseStruct(structValue reflect.Value) models.Struct {
	var s models.Struct
	var v models.Value

	for i := 0; i < structValue.NumField(); i++ {
		field := structValue.Field(i)

		switch field.Kind() {
		case reflect.Struct:
			v = models.Value{Value: ParseStruct(field), Kind: reflect.Struct}
		case reflect.Pointer:
			v = models.Value{Value: ParsePtr(field), Kind: reflect.Pointer}
		case reflect.Slice, reflect.Array:
			v = models.Value{Value: ParseSlice(field), Kind: field.Kind()}
		default:
			v = models.Value{Value: field.Interface(), Kind: field.Kind()}
		}

		tag := structValue.Type().Field(i).Tag.Get(models.FieldTag)

		s.Fields = append(s.Fields, models.Field{
			Name:  structValue.Type().Field(i).Name,
			Tag:   tag,
			Value: v,
			Opts:  options.Parse(tag),
			Kind:  field.Kind(),
		})
	}

	return s
}
