package encoder

import (
	"encoding"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/shopspring/decimal"
)

const (
	// DefaultCensorFieldTag is a default tag name for censor fields.
	DefaultCensorFieldTag = "censor"

	// UnsupportedTypeTmpl is a template for a value that is returned when a given type is not supported.
	UnsupportedTypeTmpl = "unsupported type="
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
			CensorFieldTag:  DefaultCensorFieldTag,
			ExcludePatterns: c.ExcludePatterns,
			MaskValue:       c.MaskValue,
			UseJSONTagName:  c.UseJSONTagName,
		},
		DisplayMapType:       c.DisplayMapType,
		DisplayPointerSymbol: c.DisplayPointerSymbol,
		DisplayStructName:    c.DisplayStructName,
	}
	if len(p.ExcludePatterns) != 0 {
		p.compileExcludePatterns()
	}

	return &p
}

// Struct parses a given value and returns a Struct.
// All supported complex types will be parsed recursively.
// Note: all unexported fields will be ignored.
//
//nolint:gocyclo,funlen,gocognit
func (te *TextEncoder) Struct(b *strings.Builder, rv reflect.Value) {
	if rv.Kind() != reflect.Struct {
		panic("provided value is not a struct")
	}

	t := rv.Type()

	if te.DisplayPointerSymbol {
		var pkg string
		pkgPath := t.PkgPath()
		// This custom logic is used instead of strings.Split to avoid unnecessary allocations.
		for i := len(pkgPath) - 1; i >= 0; i-- {
			// We iterate over the package path in reverse order until we find the last slash,
			// which separates the package name from the package path.
			if pkgPath[i] == '/' {
				pkg = pkgPath[i+1:]

				break
			}
			// If there is no slash in the package path, then the package name is equal to the package path.
			// Example: "main" package.
			if i == 0 {
				pkg = pkgPath[i:]

				break
			}
		}

		b.WriteString(pkg + dot + t.Name())
	}

	b.WriteString(openBrace)

	numFields := rv.NumField()
	for i := 0; i < numFields; i++ {
		f := rv.Field(i)
		if !f.CanInterface() {
			continue
		}

		strField := t.Field(i)
		if jsonName, ok := strField.Tag.Lookup("json"); ok && te.UseJSONTagName {
			b.WriteString(jsonName + colon)
		} else {
			b.WriteString(strField.Name + colon) // If tag is absent, then a struct filed name shall be used.
		}

		if strField.Tag.Get(te.CensorFieldTag) != "display" {
			b.WriteString(te.MaskValue)
			if i < numFields-1 {
				b.WriteString(commaWithSpace)
			}

			continue
		}

		te.Encode(b, f)

		if i < numFields-1 {
			b.WriteString(commaWithSpace)
		}
	}
	b.WriteString(closeBrace)
}

// Map parses a given value and returns a Map.
// If value is a struct/pointer/slice/array/map/interface, it will be parsed recursively.
// Note: this method panics if the provided value is not a complex.
func (te *TextEncoder) Map(b *strings.Builder, rv reflect.Value) {
	if rv.Kind() != reflect.Map {
		panic("provided value is not a map")
	}

	b.WriteString(rv.Type().String() + openBrace)

	var addComma bool
	for iter := rv.MapRange(); iter.Next(); {
		if addComma {
			b.WriteString(commaWithSpace)
		}

		key, value := iter.Key(), iter.Value()

		te.Encode(b, key)
		b.WriteString(colon)
		te.Encode(b, value)
		addComma = true
	}

	b.WriteString(closeBrace)
}

// Slice parses a given value and returns a Slice.
// This function is also can be used to parse an array.
// All supported complex types will be parsed recursively.
// Note: this method panics if the provided value is not a complex.
func (te *TextEncoder) Slice(b *strings.Builder, rv reflect.Value) {
	if rv.Kind() != reflect.Slice && rv.Kind() != reflect.Array {
		panic("provided value is not a slice/array")
	}

	b.WriteString(openBracket)

	for i := 0; i < rv.Len(); i++ {
		te.Encode(b, rv.Index(i))

		if i < rv.Len()-1 {
			b.WriteString(commaWithSpace)
		}
	}
	b.WriteString(closeBracket)
}

// Interface parses an interface and returns an Interface.
// In case of a pointer to unsupported type of value, a string built from UnsupportedTypeTmpl
// is used instead of the real value. That string contains a type of the value.
func (te *TextEncoder) Interface(b *strings.Builder, rv reflect.Value) {
	if rv.Kind() != reflect.Interface {
		panic("provided value is not an interface")
	}

	if rv.IsNil() {
		b.WriteString(nilSymbol)

		return
	}
	te.Encode(b, rv.Elem())
}

// Ptr parses a given value and returns a Ptr.
// If the value is nil, it returns a Ptr with a nil Value.
// In case of a pointer to unsupported type of value, a string built from UnsupportedTypeTmpl
// is used instead of the real value. That string contains a type of the value.
func (te *TextEncoder) Ptr(b *strings.Builder, rv reflect.Value) {
	if rv.Kind() != reflect.Ptr {
		panic("provided value is not a pointer")
	}

	if rv.IsNil() {
		b.WriteString(nilSymbol)

		return
	}

	if te.DisplayPointerSymbol {
		b.WriteString(ptrSymbol)
	}
	te.Encode(b, rv.Elem())
}

// String formats a value as a string.
func (te *TextEncoder) String(b *strings.Builder, s string) {
	if len(te.ExcludePatterns) != 0 {
		for _, pattern := range te.ExcludePatternsCompiled {
			if pattern.MatchString(s) {
				b.WriteString(pattern.ReplaceAllString(s, te.MaskValue))

				return
			}
		}
	}

	b.WriteString(s)
}

//nolint:exhaustive,gocyclo
func (te *TextEncoder) Encode(b *strings.Builder, f reflect.Value) {
	switch k := f.Kind(); k {
	case reflect.Struct:
		// If a field implements encoding.TextMarshaler interface, then it should be marshaled to string.
		if v, ok := f.Interface().(encoding.TextMarshaler); ok {
			b.WriteString(PrepareTextMarshalerValue(v))
		} else {
			te.Struct(b, f)
		}
	case reflect.Slice, reflect.Array:
		te.Slice(b, f)
	case reflect.Ptr:
		te.Ptr(b, f)
	case reflect.Map:
		te.Map(b, f)
	case reflect.Interface:
		te.Interface(b, f)
	case reflect.Bool:
		b.WriteString(strconv.FormatBool(f.Bool()))
	case reflect.String:
		te.String(b, f.String())
	case reflect.Float32:
		b.WriteString(decimal.NewFromFloat32(float32(f.Float())).String())
	case reflect.Float64:
		b.WriteString(decimal.NewFromFloat(f.Float()).String())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		b.WriteString(strconv.FormatInt(f.Int(), 10))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		b.WriteString(strconv.FormatInt(int64(f.Uint()), 10))
	default:
		b.WriteString(UnsupportedTypeTmpl + k.String())
	}
}

// compileExcludePatterns compiles regexp patterns from ExcludePatterns.
// Note: this method may panic if regexp pattern is invalid.
func (te *TextEncoder) compileExcludePatterns() {
	if te.ExcludePatterns != nil {
		te.ExcludePatternsCompiled = make([]*regexp.Regexp, len(te.ExcludePatterns))
		for i, pattern := range te.ExcludePatterns {
			te.ExcludePatternsCompiled[i] = regexp.MustCompile(pattern)
		}
	}
}
