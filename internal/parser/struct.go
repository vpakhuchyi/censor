package parser

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/vpakhuchyi/sanitiser/internal/models"
	"github.com/vpakhuchyi/sanitiser/internal/options"
)

// Struct parses a given value and returns a Struct.
// All fields of pointer/slice/array/struct/map types will be parsed recursively.
func (p *Parser) Struct(structValue reflect.Value) models.Struct {
	var v models.Value
	s := models.Struct{Name: getStructName(structValue)}
	for i := 0; i < structValue.NumField(); i++ {
		f := structValue.Field(i)

		switch f.Kind() {
		case reflect.Struct:
			v = models.Value{Value: p.Struct(f), Kind: reflect.Struct}
		case reflect.Pointer:
			v = models.Value{Value: p.Ptr(f), Kind: reflect.Pointer}
		case reflect.Slice, reflect.Array:
			v = models.Value{Value: p.Slice(f), Kind: f.Kind()}
		case reflect.Map:
			v = models.Value{Value: p.Map(f), Kind: f.Kind()}
		default:
			v = models.Value{Value: f.Interface(), Kind: f.Kind()}
		}

		tag := structValue.Type().Field(i).Tag.Get(p.SanitiserFieldTag)

		field := models.Field{
			Name:  structValue.Type().Field(i).Name,
			Tag:   tag,
			Value: v,
			Opts:  options.Parse(tag),
			Kind:  f.Kind(),
		}

		if p.UseJSONTagName {
			field.Name = structValue.Type().Field(i).Tag.Get("json")
		}

		s.Fields = append(s.Fields, field)
	}

	return s
}

// getStructName returns a name of the struct.
// It uses the last part of the package path and the struct name.
// If the package is the main package, then the path will be an empty string.
func getStructName(structValue reflect.Value) string {
	t := structValue.Type()
	pkg := strings.Split(t.PkgPath(), "/")

	// If the package is the main package, then the path will be an empty string.
	if len(pkg) == 1 && pkg[0] == "" {
		return t.Name()
	}

	return fmt.Sprintf("%s.%s", pkg[len(pkg)-1], t.Name())
}
