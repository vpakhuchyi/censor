package main

import (
	"reflect"

	"github.com/vpakhuchyi/sanitiser/internal/formatters"
	"github.com/vpakhuchyi/sanitiser/internal/models"
	"github.com/vpakhuchyi/sanitiser/internal/parsers"
)

// Format takes any value and returns a string representation of it.
// It uses reflection to parse the value and then uses formatters to format it.
// Examples can be found here https://github.com/vpakhuchyi/sanitiser#readme
//
// Supported types:
//
// - Struct
// By default, all fields values will be masked.
// To override this behaviour, use the `log:"display"` tag.
// All nested fields must be tagged as well.
//
// - Slice/Array
// Struct/Slice/Array/Pointer values will be parsed recursively
//
// - Pointer
// Struct/Slice/Array/Pointer values will be parsed recursively
//
// - String
// - Float64
// Formatted value will have up to 15 precision digits.
//
// - Float32
// Formatted value will have up to 7 precision digits.
//
// - Int/Int8/Int16/Int32/Int64/Rune
// - Uint/Uint8/Uint16/Uint32/Uint64/Byte
// - Bool
func Format(val any) string {
	return sanitise(val)
}

func sanitise(val any) string {
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
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return formatters.FormatInteger(parsed.(models.Value))
	case reflect.Bool:
		return formatters.FormatBool(parsed.(models.Value))
	}

	return formatters.FormatSimple(parsed.(models.Value))
}
