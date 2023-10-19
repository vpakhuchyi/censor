package parsers

import (
	"reflect"

	"sanitiser/internal/models"
)

func ParseSlice(sliceValue reflect.Value) models.Slice {
	var slice models.Slice
	for i := 0; i < sliceValue.Len(); i++ {
		elem := sliceValue.Index(i)
		switch elem.Kind() {
		case reflect.Struct:
			slice.Values = append(slice.Values, models.Value{Value: ParseStruct(elem), Kind: reflect.Struct})
		case reflect.Pointer:
			slice.Values = append(slice.Values, models.Value{Value: ParsePtr(elem), Kind: reflect.Pointer})
		case reflect.Slice, reflect.Array:
			slice.Values = append(slice.Values, models.Value{Value: ParseSlice(elem), Kind: elem.Kind()})
		default:
			slice.Values = append(slice.Values, models.Value{Value: elem.Interface(), Kind: elem.Kind()})
		}
	}

	return slice
}
