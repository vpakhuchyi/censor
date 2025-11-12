package encoder

import (
	"bytes"
	"encoding"
	"encoding/json"
	"reflect"
	"strconv"
	"unicode/utf8"

	"github.com/shopspring/decimal"

	"github.com/vpakhuchyi/censor/internal/cache"
)

// NewJSONEncoder returns a new instance of JSONEncoder with given configuration.
func NewJSONEncoder(c Config) *JSONEncoder {
	e := &JSONEncoder{
		baseEncoder: baseEncoder{
			CensorFieldTag:      defaultCensorFieldTag,
			ExcludePatterns:     c.ExcludePatterns,
			MaskValue:           c.MaskValue,
			structFieldsCache:   cache.NewTypeCache[[]Field](cache.DefaultMaxCacheSize),
			escapedStringsCache: cache.New[string](cache.DefaultMaxCacheSize),
			regexpCache:         cache.New[string](cache.DefaultMaxCacheSize),
		},
	}

	if len(e.ExcludePatterns) != 0 {
		e.ExcludePatternsCompiled = compileRegexpPatterns(e.ExcludePatterns)
	}

	return e
}

// JSONEncoder is used to encode data to JSON format.
type JSONEncoder struct {
	baseEncoder
}

//nolint:exhaustive,gocyclo
func (e *JSONEncoder) Encode(b *bytes.Buffer, f reflect.Value) {
	if !f.IsValid() {
		b.WriteString("null")

		return
	}

	switch k := f.Kind(); k {
	case reflect.Struct:
		if f.CanInterface() {
			// If a field implements json.Marshaler interface, then it should be marshaled to string.
			v, ok := f.Interface().(json.Marshaler)
			if ok {
				b.WriteString(PrepareJSONMarshalerValue(v))

				return
			}
		}

		e.Struct(b, f)
	case reflect.Slice, reflect.Array:
		e.Slice(b, f)
	case reflect.Pointer:
		e.Ptr(b, f)
	case reflect.Map:
		e.Map(b, f)
	case reflect.Interface:
		e.Interface(b, f)
	case reflect.Bool:
		b.WriteString(strconv.FormatBool(f.Bool()))
	case reflect.String:
		e.StringEscaped(b, f.String())
	case reflect.Float32:
		b.WriteString(decimal.NewFromFloat32(float32(f.Float())).String())
	case reflect.Float64:
		b.WriteString(decimal.NewFromFloat(f.Float()).String())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		b.WriteString(strconv.FormatInt(f.Int(), 10))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		b.WriteString(strconv.FormatUint(f.Uint(), 10))
	default:
		b.WriteString(`"` + unsupportedTypeTmpl + k.String() + `"`)
	}
}

// Struct encodes a struct value to JSON format.
// Note: this method panics if the provided value is not a map.
func (e *JSONEncoder) Struct(b *bytes.Buffer, v reflect.Value) {
	if v.Kind() != reflect.Struct {
		panic("provided value is not a struct")
	}

	t := v.Type()
	var fields []Field

	if t.PkgPath() == "" {
		fields = e.getStructFields(v)
	} else {
		var found bool
		fields, found = e.structFieldsCache.Get(t)
		if !found {
			fields = e.getStructFields(v)
			e.structFieldsCache.Set(t, fields)
		}
	}

	b.WriteByte('{')

	firstField := true
	for i, field := range fields {
		if field.Name == "" {
			continue
		}

		if !firstField {
			b.WriteByte(',')
		}
		firstField = false

		b.WriteByte('"')
		b.WriteString(field.Name)
		b.WriteString(`": `)

		if field.IsMasked {
			b.WriteByte('"')
			b.WriteString(e.MaskValue)
			b.WriteByte('"')
		} else {
			e.Encode(b, v.Field(i))
		}
	}
	b.WriteByte('}')
}

func (e *JSONEncoder) getStructFields(v reflect.Value) []Field {
	numFields := v.NumField()
	fields := make([]Field, numFields)

	for i := 0; i < numFields; i++ {
		field := v.Type().Field(i)
		if !v.Field(i).CanInterface() {
			continue
		}

		fields[i] = Field{
			Name:     field.Name,
			IsMasked: field.Tag.Get(e.CensorFieldTag) != display,
		}
	}

	return fields
}

// Map encodes a map value to JSON format.
// Note: this method panics if the provided value is not a map.
func (e *JSONEncoder) Map(b *bytes.Buffer, v reflect.Value) {
	if v.Kind() != reflect.Map {
		panic("provided value is not a map")
	}

	if v.IsNil() {
		b.WriteString("null")

		return
	}

	b.WriteByte('{')

	first := true
	for iter := v.MapRange(); iter.Next(); {
		if !first {
			b.WriteByte(',')
		}
		first = false

		key, value := iter.Key(), iter.Value()

		e.encodeMapKey(b, key)
		b.WriteByte(':')
		e.Encode(b, value)
	}
	b.WriteByte('}')
}

// Slice encodes a slice value to JSON format.
// This function is also can be used to parse an array.
// Note: this method panics if the provided value is not a slice or array.
func (e *JSONEncoder) Slice(b *bytes.Buffer, v reflect.Value) {
	if v.Kind() != reflect.Slice && v.Kind() != reflect.Array {
		panic("provided value is not a slice/array")
	}

	b.WriteByte('[')
	length := v.Len()
	for i := 0; i < length; i++ {
		e.Encode(b, v.Index(i))

		if i < length-1 {
			b.WriteByte(',')
			b.WriteByte(' ')
		}
	}
	b.WriteByte(']')
}

// Interface encodes an interface value to JSON format.
// In case of a pointer to unsupported type of value, a string built from unsupportedTypeTmpl
// is used instead of the real value. That string contains a type of the value.
// Note: this method panics if the provided value is not an interface.
func (e *JSONEncoder) Interface(b *bytes.Buffer, v reflect.Value) {
	if v.Kind() != reflect.Interface {
		panic("provided value is not an interface")
	}

	if v.IsNil() {
		b.WriteString("null")

		return
	}
	e.Encode(b, v.Elem())
}

// Ptr encodes a pointer value to JSON format.
// In case of a pointer to unsupported type of value, a string built from unsupportedTypeTmpl
// is used instead of the real value. That string contains a type of the value.
// Note: this method panics if the provided value is not a pointer.
func (e *JSONEncoder) Ptr(b *bytes.Buffer, v reflect.Value) {
	if v.Kind() != reflect.Pointer {
		panic("provided value is not a pointer")
	}

	if v.IsNil() {
		b.WriteString("null")

		return
	}

	e.Encode(b, v.Elem())
}

// String encodes the input string by masking any substrings that match the configured exclusion patterns.
// It replaces matched segments with a predefined mask value to censor sensitive information.
func (e *JSONEncoder) String(b *bytes.Buffer, s string) {
	e.WriteString(b, s)
}

// StringEscaped encodes and escapes the input string by masking any substrings that match the configured exclusion patterns.
// It replaces matched segments with a predefined mask value to censor sensitive information.
func (e *JSONEncoder) StringEscaped(b *bytes.Buffer, s string) {
	b.WriteByte('"')
	e.writeEscapedCensoredString(b, s)
	b.WriteByte('"')
}

// writeEscapedCensoredString applies censoring and escaping in one pass.
func (e *JSONEncoder) writeEscapedCensoredString(b *bytes.Buffer, s string) {
	if len(e.ExcludePatterns) == 0 || e.ExcludePatternsCompiled == nil {
		e.escapeString(b, s)

		return
	}

	cached, ok := e.regexpCache.Get(s)
	if ok {
		e.escapeStringWithCache(b, cached)

		return
	}

	matches := e.ExcludePatternsCompiled.FindAllStringIndex(s, -1)
	if len(matches) == 0 {
		e.escapeStringWithCache(b, s)
		e.regexpCache.Set(s, s)

		return
	}

	lastIndex := 0
	for _, m := range matches {
		start, end := m[0], m[1]
		e.escapeString(b, s[lastIndex:start])
		e.escapeString(b, e.MaskValue)
		lastIndex = end
	}
	e.escapeString(b, s[lastIndex:])
}

// escapeString processes the input string by escaping special and control characters to ensure it is safe
// for JSON encoding. It replaces characters like backslashes, quotes, and control characters with their corresponding
// escape sequences. If the string contains non-ASCII characters, it ensures they are properly escaped or replaced
// with the Unicode replacement character if invalid.
//
//nolint:gocyclo,mnd
func (e *JSONEncoder) escapeString(b *bytes.Buffer, s string) {
	for _, r := range s {
		switch r {
		case '\\', '"':
			b.WriteByte('\\')
			b.WriteRune(r)
		case '\b':
			b.WriteByte('\\')
			b.WriteByte('b')
		case '\f':
			b.WriteByte('\\')
			b.WriteByte('f')
		case '\n':
			b.WriteByte('\\')
			b.WriteByte('n')
		case '\r':
			b.WriteByte('\\')
			b.WriteByte('r')
		case '\t':
			b.WriteByte('\\')
			b.WriteByte('t')
		default:
			switch {
			case (r >= 0x00 && r <= 0x1F) || r == 0x7F:
				b.WriteString(escapeControlChar(r))
			case r > 0x7F:
				if utf8.ValidRune(r) && r != utf8.RuneError {
					b.WriteString(escapeControlChar(r))
				} else {
					b.WriteRune('\uFFFD')
				}
			default:
				b.WriteRune(r)
			}
		}
	}
}

// escapeStringWithCache escapes string content with caching.
func (e *JSONEncoder) escapeStringWithCache(b *bytes.Buffer, s string) {
	if cached, ok := e.escapedStringsCache.Get(s); ok {
		b.WriteString(cached)

		return
	}

	startLen := b.Len()
	e.escapeString(b, s)

	escaped := b.String()[startLen:]
	e.escapedStringsCache.Set(s, escaped)
}

//nolint:exhaustive
func (e *JSONEncoder) encodeMapKey(b *bytes.Buffer, f reflect.Value) {
	switch k := f.Kind(); k {
	case reflect.String:
		e.StringEscaped(b, f.String())
	case reflect.Float32:
		b.WriteString(decimal.NewFromFloat32(float32(f.Float())).String())
	case reflect.Float64:
		b.WriteString(decimal.NewFromFloat(f.Float()).String())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		b.WriteString(strconv.FormatInt(f.Int(), 10))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		b.WriteString(strconv.FormatUint(f.Uint(), 10))
	default:
		if v, ok := f.Interface().(encoding.TextMarshaler); ok {
			b.WriteString(PrepareTextMarshalerValue(v))
		} else {
			b.WriteString(unsupportedTypeTmpl + k.String())
		}
	}
}

// escapeControlChar converts a control character rune into its corresponding Unicode escape sequence.
// It formats the rune as a four-digit hexadecimal number prefixed with '\u' to comply with JSON string encoding.
func escapeControlChar(r rune) string {
	hexStr := strconv.FormatInt(int64(r), 16)
	for len(hexStr) < 4 {
		hexStr = "0" + hexStr
	}

	return `\u` + hexStr
}
