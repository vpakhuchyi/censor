package parser

import (
	"fmt"
	"reflect"
	"sort"

	"github.com/vpakhuchyi/censor/internal/models"
)

// Map parses a given value and returns a Map.
// If value is a struct/pointer/slice/array/map/interface, it will be parsed recursively.
// Note: this method panics if the provided value is not a complex.
//
//nolint:exhaustive,funlen,gocyclo
func (p *Parser) Map(mapValue reflect.Value) models.Map {
	if mapValue.Kind() != reflect.Map {
		panic("provided value is not a map")
	}

	m := models.Map{Type: mapValue.Type().String()}
	iter := mapValue.MapRange()

	for iter.Next() {
		key, value := iter.Key(), iter.Value()

		pair := models.KV{SortValue: fmt.Sprintf("%v", key.Interface())}

		switch k := value.Kind(); k {
		case reflect.Struct:
			pair.Value = models.Value{Value: p.Struct(value), Kind: value.Kind()}
		case reflect.Pointer:
			pair.Value = models.Value{Value: p.Ptr(value), Kind: value.Kind()}
		case reflect.Slice, reflect.Array:
			pair.Value = models.Value{Value: p.Slice(value), Kind: value.Kind()}
		case reflect.Map:
			pair.Value = models.Value{Value: p.Map(value), Kind: value.Kind()}
		case reflect.Interface:
			pair.Value = models.Value{Value: p.Interface(value), Kind: value.Kind()}
		case reflect.String:
			pair.Value = p.String(value)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			pair.Value = p.Integer(value)
		case reflect.Float32, reflect.Float64:
			pair.Value = p.Float(value)
		case reflect.Bool:
			pair.Value = p.Bool(value)
		case reflect.Complex64, reflect.Complex128:
			pair.Value = p.Complex(value)
		default:
			pair.Value = models.Value{Value: fmt.Sprintf("unsupported type: %s", k.String()), Kind: reflect.String}
		}

		switch key.Kind() {
		case reflect.Struct:
			pair.Key = models.Value{Value: p.Struct(key), Kind: key.Kind()}
		case reflect.Pointer:
			pair.Key = models.Value{Value: p.Ptr(key), Kind: key.Kind()}
		case reflect.Array:
			pair.Key = models.Value{Value: p.Slice(key), Kind: key.Kind()}
		case reflect.Interface:
			pair.Key = models.Value{Value: p.Interface(key), Kind: key.Kind()}
		case reflect.String:
			pair.Key = p.String(key)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			pair.Key = p.Integer(key)
		case reflect.Float32, reflect.Float64:
			pair.Key = p.Float(key)
		case reflect.Bool:
			pair.Key = p.Bool(key)
		case reflect.Complex64, reflect.Complex128:
			pair.Key = p.Complex(key)
		}

		m.Values = append(m.Values, pair)
	}

	// Sort map keys to make output deterministic.
	sort.SliceStable(m.Values, func(i, j int) bool {
		return m.Values[i].SortValue < m.Values[j].SortValue
	})

	return m
}
