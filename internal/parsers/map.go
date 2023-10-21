package parsers

import (
	"fmt"
	"reflect"
	"sort"

	"github.com/vpakhuchyi/sanitiser/internal/models"
)

// ParseMap parses a given value and returns a Map.
// If value is a struct/pointer/slice/array/map, it will be parsed recursively.
func ParseMap(mapValue reflect.Value) models.Map {
	m := models.Map{Type: mapValue.Type().String()}
	iter := mapValue.MapRange()

	for iter.Next() {
		key, value := iter.Key(), iter.Value()

		pair := models.KV{SortValue: fmt.Sprintf("%v", key.Interface())}

		switch value.Kind() {
		case reflect.Struct:
			pair.Value = models.Value{Value: ParseStruct(value), Kind: value.Kind()}
		case reflect.Pointer:
			pair.Value = models.Value{Value: ParsePtr(value), Kind: value.Kind()}
		case reflect.Slice, reflect.Array:
			pair.Value = models.Value{Value: ParseSlice(value), Kind: value.Kind()}
		case reflect.Map:
			pair.Value = models.Value{Value: ParseMap(value), Kind: value.Kind()}
		default:
			pair.Value = models.Value{Value: value.Interface(), Kind: value.Kind()}
		}

		switch key.Kind() {
		case reflect.Struct:
			pair.Key = models.Value{Value: ParseStruct(key), Kind: key.Kind()}
		case reflect.Pointer:
			pair.Key = models.Value{Value: ParsePtr(key), Kind: key.Kind()}
		case reflect.Array:
			pair.Key = models.Value{Value: ParseSlice(key), Kind: key.Kind()}
		default:
			pair.Key = models.Value{Value: key.Interface(), Kind: key.Kind()}
		}

		m.Values = append(m.Values, pair)
	}

	// Sort map keys to make output deterministic.
	sort.SliceStable(m.Values, func(i, j int) bool {
		return m.Values[i].SortValue < m.Values[j].SortValue
	})

	return m
}
