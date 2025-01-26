package encoder

import (
	"encoding"
	"reflect"
	"strconv"
	"strings"

	"github.com/shopspring/decimal"

	"github.com/vpakhuchyi/censor/internal/cache"
)

// TextEncoder is a struct that contains options for parsing.
type TextEncoder struct {
	baseEncoder
	// DisplayMapType is used to display map type in the output.
	// The default value is false.
	DisplayMapType bool
	// DisplayPointerSymbol is used to display '&' (pointer symbol) in the output.
	// The default value is false.
	DisplayPointerSymbol bool
	// DisplayStructName is used to display struct name in the output.
	// A struct name includes the last part of the package path.
	// The default value is false.
	DisplayStructName bool
}

// NewTextEncoder returns a new instance of TextEncoder with given configuration.
func NewTextEncoder(c Config) *TextEncoder {
	p := TextEncoder{
		baseEncoder: baseEncoder{
			CensorFieldTag:    DefaultCensorFieldTag,
			ExcludePatterns:   c.ExcludePatterns,
			MaskValue:         c.MaskValue,
			UseJSONTagName:    c.UseJSONTagName,
			structFieldsCache: cache.NewSlice[Field](cache.DefaultMaxCacheSize),
			regexpCache:       cache.New[string](cache.DefaultMaxCacheSize),
		},
		DisplayMapType:       c.DisplayMapType,
		DisplayPointerSymbol: c.DisplayPointerSymbol,
		DisplayStructName:    c.DisplayStructName,
	}
	if len(p.ExcludePatterns) != 0 {
		p.ExcludePatternsCompiled = compileRegexpPatterns(p.ExcludePatterns)
	}

	return &p
}

//nolint:exhaustive,gocyclo
func (e *TextEncoder) Encode(b *strings.Builder, f reflect.Value) {
	switch k := f.Kind(); k {
	case reflect.Struct:
		// If a field implements encoding.TextMarshaler interface, then it should be marshaled to string.
		if v, ok := f.Interface().(encoding.TextMarshaler); ok {
			b.WriteString(PrepareTextMarshalerValue(v))
		} else {
			e.Struct(b, f)
		}
	case reflect.Slice, reflect.Array:
		e.Slice(b, f)
	case reflect.Ptr:
		e.Ptr(b, f)
	case reflect.Map:
		e.Map(b, f)
	case reflect.Interface:
		e.Interface(b, f)
	case reflect.Bool:
		b.WriteString(strconv.FormatBool(f.Bool()))
	case reflect.String:
		e.String(b, f.String())
	case reflect.Float32:
		b.WriteString(decimal.NewFromFloat32(float32(f.Float())).String())
	case reflect.Float64:
		b.WriteString(decimal.NewFromFloat(f.Float()).String())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		b.WriteString(strconv.FormatInt(f.Int(), 10))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		b.WriteString(strconv.FormatUint(f.Uint(), 10))
	default:
		b.WriteString(UnsupportedTypeTmpl + k.String())
	}
}

// Struct encodes a struct value to TEXT format.
//
//nolint:gocyclo,gocognit
func (e *TextEncoder) Struct(b *strings.Builder, v reflect.Value) {
	if v.Kind() != reflect.Struct {
		panic("provided value is not a struct")
	}

	t := v.Type()

	structPath := v.Type().PkgPath()
	structName := t.Name()

	var fields []Field
	key := structPath + structName
	if key == "" {
		fields = e.getStructFields(v, t)
	} else {
		var found bool
		fields, found = e.structFieldsCache.Get(key)
		if !found {
			fields = e.getStructFields(v, t)
			e.structFieldsCache.Set(key, fields)
		}
	}

	if e.DisplayStructName {
		var pkg string
		// This custom logic is used instead of strings.Split to avoid unnecessary allocations.
		for i := len(structPath) - 1; i >= 0; i-- {
			// We iterate over the package path in reverse order until we find the last slash,
			// which separates the package name from the package path.
			if structPath[i] == '/' {
				pkg = structPath[i+1:]

				break
			}
			// If there is no slash in the package path, then the package name is equal to the package path.
			// Example: "main" package.
			if i == 0 {
				pkg = structPath[i:]

				break
			}
		}

		b.WriteString(pkg + "." + structName)
	}

	b.WriteString("{")

	for i := 0; i < len(fields); i++ {
		b.WriteString(fields[i].Name)

		if fields[i].IsMasked {
			b.WriteString(e.MaskValue)
		} else {
			e.Encode(b, v.Field(i))
		}

		if i < len(fields)-1 {
			b.WriteString(`, `)
		}
	}

	b.WriteString("}")
}

func (e *TextEncoder) getStructFields(v reflect.Value, t reflect.Type) []Field {
	var fields []Field
	var name string
	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)

		if !v.Field(i).CanInterface() {
			continue
		}

		if e.UseJSONTagName {
			name = field.Tag.Get("json")
		} else {
			name = field.Name
		}

		fields = append(fields, Field{
			Name:     name + `: `,
			IsMasked: field.Tag.Get(e.CensorFieldTag) != display,
		})
	}

	return fields
}

// Map encodes a map value to TEXT format.
// Note: this method panics if the provided value is not a map.
func (e *TextEncoder) Map(b *strings.Builder, rv reflect.Value) {
	if rv.Kind() != reflect.Map {
		panic("provided value is not a map")
	}

	b.WriteString(rv.Type().String() + "{")

	var addComma bool
	for iter := rv.MapRange(); iter.Next(); {
		if addComma {
			b.WriteString(", ")
		}

		key, value := iter.Key(), iter.Value()

		e.Encode(b, key)
		b.WriteString(": ")
		e.Encode(b, value)
		addComma = true
	}

	b.WriteString("}")
}

// Slice encodes a slice value to TEXT format.
// This function is also can be used to parse an array.
// Note: this method panics if the provided value is not a slice or array.
func (e *TextEncoder) Slice(b *strings.Builder, rv reflect.Value) {
	if rv.Kind() != reflect.Slice && rv.Kind() != reflect.Array {
		panic("provided value is not a slice/array")
	}

	b.WriteString("[")

	for i := 0; i < rv.Len(); i++ {
		e.Encode(b, rv.Index(i))

		if i < rv.Len()-1 {
			b.WriteString(", ")
		}
	}
	b.WriteString("]")
}

// Interface encodes an interface value to TEXT format.
// In case of a pointer to unsupported type of value, a string built from UnsupportedTypeTmpl
// is used instead of the real value. That string contains a type of the value.
// Note: this method panics if the provided value is not an interface.
func (e *TextEncoder) Interface(b *strings.Builder, rv reflect.Value) {
	if rv.Kind() != reflect.Interface {
		panic("provided value is not an interface")
	}

	if rv.IsNil() {
		b.WriteString("nil")

		return
	}
	e.Encode(b, rv.Elem())
}

// Ptr encodes a pointer value to TEXT format.
// In case of a pointer to unsupported type of value, a string built from UnsupportedTypeTmpl
// is used instead of the real value. That string contains a type of the value.
// Note: this method panics if the provided value is not a pointer.
func (e *TextEncoder) Ptr(b *strings.Builder, rv reflect.Value) {
	if rv.Kind() != reflect.Ptr {
		panic("provided value is not a pointer")
	}

	if rv.IsNil() {
		b.WriteString("nil")

		return
	}

	if e.DisplayPointerSymbol {
		b.WriteString("&")
	}
	e.Encode(b, rv.Elem())
}

// String formats a value as a string.
// If the string matches one of the ExcludePatterns, it will be masked with the MaskValue.
func (e *TextEncoder) String(b *strings.Builder, s string) {
	b.WriteString(e.baseEncoder.String(s))
}
