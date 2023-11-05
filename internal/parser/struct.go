package parser

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/vpakhuchyi/censor/internal/models"
	"github.com/vpakhuchyi/censor/internal/options"
)

// Struct parses a given value and returns a Struct.
// All supported complex types will be parsed recursively.
//
//nolint:exhaustive
func (p *Parser) Struct(structValue reflect.Value) models.Struct {
	var v models.Value
	s := models.Struct{Name: getStructName(structValue)}

	for i := 0; i < structValue.NumField(); i++ {
		var fieldName string
		if p.UseJSONTagName {
			tagValue := structValue.Type().Field(i).Tag.Get("json")

			// If the tag is not present, then such a field will be ignored.
			if tagValue == "" {
				continue
			}

			fieldName = tagValue
		} else {
			fieldName = structValue.Type().Field(i).Name
		}

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
		case reflect.Interface:
			v = models.Value{Value: p.Interface(f), Kind: f.Kind()}
		case reflect.Bool:
			v = p.Bool(f)
		default:
			v = models.Value{Value: f.Interface(), Kind: f.Kind()}
		}

		tag := structValue.Type().Field(i).Tag.Get(p.CensorFieldTag)

		s.Fields = append(s.Fields, models.Field{
			Name:  fieldName,
			Tag:   tag,
			Value: v,
			Opts:  options.Parse(tag),
			Kind:  f.Kind(),
		})
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
