package main

import (
	"reflect"

	"sanitiser/internal/formatters"
	"sanitiser/internal/models"
	"sanitiser/internal/parsers"
)

func Sanitise(val any) string {
	return sanitized(val)
}

func sanitized(val any) string {
	v := reflect.ValueOf(val)

	var parsed any
	switch v.Kind() {
	case reflect.Struct:
		parsed = parsers.ParseStruct(v)
	case reflect.Slice, reflect.Array:
		parsed = parsers.ParseSlice(v)
	case reflect.Ptr:
		parsed = parsers.ParsePtr(v)
	default:
		parsed = models.Value{Value: v.Interface(), Kind: v.Kind()}
	}

	switch v.Kind() {
	case reflect.Struct:
		return formatters.FormatStruct(parsed.(models.Struct))
	case reflect.Slice, reflect.Array:
		return formatters.FormatSlice(parsed.(models.Slice))
	case reflect.Pointer:
		return formatters.FormatPtr(parsed.(models.Ptr))
	case reflect.String:
		return formatters.FormatString(parsed.(models.Value))
	case reflect.Float32, reflect.Float64:
		return formatters.FormatFloat(parsed.(models.Value))
	default:
		return formatters.FormatSimple(parsed.(models.Value))
	}
}
