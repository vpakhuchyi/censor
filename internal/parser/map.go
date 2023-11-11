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
func (p *Parser) Map(rv reflect.Value) models.Map {
	if rv.Kind() != reflect.Map {
		panic("provided value is not a map")
	}

	m := models.Map{
		Type:   rv.Type().String(),
		Values: make([]models.KV, 0, rv.Len()),
	}

	iter := rv.MapRange()
	for iter.Next() {
		key := iter.Key()
		pair := models.KV{SortValue: fmt.Sprintf("%v", key.Interface())}

		switch key.Kind() {
		case reflect.Struct:
			pair.Key = models.Value{Value: p.Struct(key), Kind: reflect.Struct}
		case reflect.Pointer:
			pair.Key = models.Value{Value: p.Ptr(key), Kind: reflect.Pointer}
		case reflect.Array:
			pair.Key = models.Value{Value: p.Slice(key), Kind: reflect.Array}
		case reflect.Interface:
			pair.Key = models.Value{Value: p.Interface(key), Kind: reflect.Interface}
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

		value := iter.Value()
		switch k := value.Kind(); k {
		case reflect.Struct:
			pair.Value = models.Value{Value: p.Struct(value), Kind: reflect.Struct}
		case reflect.Pointer:
			pair.Value = models.Value{Value: p.Ptr(value), Kind: reflect.Pointer}
		case reflect.Slice, reflect.Array:
			pair.Value = models.Value{Value: p.Slice(value), Kind: k}
		case reflect.Map:
			pair.Value = models.Value{Value: p.Map(value), Kind: reflect.Map}
		case reflect.Interface:
			pair.Value = models.Value{Value: p.Interface(value), Kind: reflect.Interface}
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

		m.Values = append(m.Values, pair)
	}

	// Sort map keys to make output deterministic.
	sort.SliceStable(m.Values, func(i, j int) bool {
		return m.Values[i].SortValue < m.Values[j].SortValue
	})

	return m
}
