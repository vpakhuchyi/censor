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
//nolint:exhaustive,gocyclo
func (p *Parser) Struct(rv reflect.Value) models.Struct {
	if rv.Kind() != reflect.Struct {
		panic("provided value is not a struct")
	}

	s := models.Struct{
		Name:   getStructName(rv),
		Fields: make([]models.Field, 0, rv.NumField()),
	}

	for i := 0; i < rv.NumField(); i++ {
		field := models.Field{
			Opts: options.Parse(rv.Type().Field(i).Tag.Get(p.censorFieldTag)),
			Kind: rv.Field(i).Kind(),
		}

		if p.useJSONTagName {
			tagValue := rv.Type().Field(i).Tag.Get("json")
			if tagValue == "" {
				// If the tag is not present, then such a field will be ignored.
				continue
			}

			field.Name = tagValue
		} else {
			field.Name = rv.Type().Field(i).Name
		}

		f := rv.Field(i)

		switch k := field.Kind; k {
		case reflect.Struct:
			field.Value = models.Value{Value: p.Struct(f), Kind: reflect.Struct}
		case reflect.Pointer:
			field.Value = models.Value{Value: p.Ptr(f), Kind: reflect.Pointer}
		case reflect.Slice, reflect.Array:
			field.Value = models.Value{Value: p.Slice(f), Kind: k}
		case reflect.Map:
			field.Value = models.Value{Value: p.Map(f), Kind: reflect.Map}
		case reflect.Interface:
			field.Value = models.Value{Value: p.Interface(f), Kind: reflect.Interface}
		case reflect.Bool:
			field.Value = p.Bool(f)
		case reflect.String:
			field.Value = p.String(f)
		case reflect.Float32, reflect.Float64:
			field.Value = p.Float(f)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			field.Value = p.Integer(f)
		case reflect.Complex64, reflect.Complex128:
			field.Value = p.Complex(f)
		}

		s.Fields = append(s.Fields, field)
	}

	return s
}

// getStructName returns a name of the struct.
// It uses the last part of the package path and the struct name.
func getStructName(structValue reflect.Value) string {
	t := structValue.Type()
	pkg := strings.Split(t.PkgPath(), "/")

	return fmt.Sprintf("%s.%s", pkg[len(pkg)-1], t.Name())
}
