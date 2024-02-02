package parser

import (
	"encoding"
	"fmt"
	"reflect"
	"strings"

	"github.com/vpakhuchyi/censor/internal/models"
	"github.com/vpakhuchyi/censor/internal/options"
)

// Struct parses a given value and returns a Struct.
// All supported complex types will be parsed recursively.
// Note: all unexported fields will be ignored.
//
//nolint:exhaustive,gocyclo,funlen,gocognit
func (p *Parser) Struct(rv reflect.Value) models.Struct {
	if rv.Kind() != reflect.Struct {
		panic("provided value is not a struct")
	}

	s := models.Struct{
		Name:   getStructName(rv),
		Fields: make([]models.Field, 0, rv.NumField()),
	}

	for i := 0; i < rv.NumField(); i++ {
		f := rv.Field(i)
		if !f.CanInterface() {
			continue
		}

		strField := rv.Type().Field(i)

		field := models.Field{
			Opts: options.Parse(strField.Tag.Get(p.censorFieldTag)),
		}

		if p.useJSONTagName {
			if jsonName, ok := strField.Tag.Lookup("json"); ok {
				field.Name = jsonName
			} else {
				field.Name = strField.Name // If tag is absent, then a struct filed name shall be used.
			}
		} else {
			field.Name = strField.Name
		}

		switch k := f.Kind(); k {
		case reflect.Struct:
			// If a field implements encoding.TextMarshaler interface, then it should be marshaled to string.
			if v, ok := f.Interface().(encoding.TextMarshaler); ok {
				field.Value = models.Value{Value: PrepareTextMarshalerValue(v), Kind: reflect.String}
			} else {
				field.Value = models.Value{Value: p.Struct(f), Kind: reflect.Struct}
			}
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
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
			field.Value = p.Integer(f)
		default:
			field.Value = models.Value{Value: fmt.Sprintf(UnsupportedTypeTmpl, k), Kind: k}
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

// PrepareTextMarshalerValue marshals a value that implements [encoding.TextMarshaler] interface to string.
func PrepareTextMarshalerValue(tm encoding.TextMarshaler) string {
	data, err := tm.MarshalText()
	if err != nil {
		return fmt.Sprintf("!ERROR:%v", err)
	}

	return string(data)
}
