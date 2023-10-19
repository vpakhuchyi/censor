package parsers

import (
	"reflect"

	"sanitiser/internal/models"
)

func ParsePtr(ptrValue reflect.Value) models.Ptr {
	if ptrValue.IsNil() {
		return models.Ptr{Value: models.Value{Value: nil, Kind: reflect.Ptr}}
	}

	switch ptrValue.Elem().Kind() {
	case reflect.Struct:
		return models.Ptr{Value: models.Value{Value: ParseStruct(ptrValue.Elem()), Kind: ptrValue.Elem().Kind()}}
	case reflect.Slice, reflect.Array:
		return models.Ptr{Value: models.Value{Value: ParseSlice(ptrValue.Elem()), Kind: ptrValue.Elem().Kind()}}
	case reflect.Ptr:
		return models.Ptr{Value: models.Value{Value: ParsePtr(ptrValue.Elem()), Kind: ptrValue.Elem().Kind()}}
	default:
		return models.Ptr{Value: models.Value{Value: ptrValue.Elem().Interface(), Kind: ptrValue.Elem().Kind()}}
	}
}
