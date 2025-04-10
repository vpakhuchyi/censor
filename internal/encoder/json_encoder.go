package encoder

import (
	"bytes"
	"encoding"
	"encoding/json"
	"reflect"
	"strconv"
	"unicode/utf8"

	"github.com/shopspring/decimal"

	"github.com/vpakhuchyi/censor/internal/builderpool"
	"github.com/vpakhuchyi/censor/internal/cache"
)

// NewJSONEncoder returns a new instance of JSONEncoder with given configuration.
func NewJSONEncoder(c Config) *JSONEncoder {
	e := &JSONEncoder{
		baseEncoder: baseEncoder{
			CensorFieldTag:      defaultCensorFieldTag,
			ExcludePatterns:     c.ExcludePatterns,
			MaskValue:           c.MaskValue,
			structFieldsCache:   cache.NewSlice[Field](cache.DefaultMaxCacheSize),
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

// Field is a struct that contains information about a struct field.
type Field struct {
	Name     string
	IsMasked bool
}

//nolint:exhaustive,gocyclo
func (e *JSONEncoder) Encode(b *bytes.Buffer, f reflect.Value) {
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
	case reflect.Ptr:
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

	var fields []Field
	key := v.Type().PkgPath() + v.Type().Name()
	if key == "" {
		fields = e.getStructFields(v)
	} else {
		var found bool
		fields, found = e.structFieldsCache.Get(key)
		if !found {
			fields = e.getStructFields(v)
			e.structFieldsCache.Set(key, fields)
		}
	}

	b.WriteString(`{`)

	for i := 0; i < len(fields); i++ {
		b.WriteString(fields[i].Name)

		if fields[i].IsMasked {
			b.WriteString(`"` + e.MaskValue + `"`)
		} else {
			e.Encode(b, v.Field(i))
		}

		if i < len(fields)-1 {
			b.WriteString(`,`)
		}
	}

	b.WriteString(`}`)
}

func (e *JSONEncoder) getStructFields(v reflect.Value) []Field {
	fields := make([]Field, v.NumField())

	b := builderpool.Get()
	for i := 0; i < v.NumField(); i++ {
		field := v.Type().Field(i)

		if !v.Field(i).CanInterface() {
			continue
		}
		b.WriteString(`"` + field.Name + `": `)

		fields[i] = Field{
			Name:     b.String(),
			IsMasked: field.Tag.Get(e.CensorFieldTag) != display,
		}
		b.Reset()
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

	b.WriteString("{")

	var addComma bool
	for iter := v.MapRange(); iter.Next(); {
		if addComma {
			b.WriteString(`,`)
		}

		key, value := iter.Key(), iter.Value()

		e.encodeMapKey(b, key)
		b.WriteString(`:`)
		e.Encode(b, value)
		addComma = true
	}

	b.WriteString("}")
}

// Slice encodes a slice value to JSON format.
// This function is also can be used to parse an array.
// Note: this method panics if the provided value is not a slice or array.
func (e *JSONEncoder) Slice(b *bytes.Buffer, v reflect.Value) {
	if v.Kind() != reflect.Slice && v.Kind() != reflect.Array {
		panic("provided value is not a slice/array")
	}

	b.WriteString("[")
	length := v.Len()
	for i := 0; i < length; i++ {
		e.Encode(b, v.Index(i))

		if i < length-1 {
			b.WriteString(", ")
		}
	}
	b.WriteString("]")
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
	if v.Kind() != reflect.Ptr {
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
	b.WriteString(e.baseEncoder.String(s))
}

// StringEscaped encodes and escapes the input string by masking any substrings that match the configured exclusion patterns.
// It replaces matched segments with a predefined mask value to censor sensitive information.
func (e *JSONEncoder) StringEscaped(b *bytes.Buffer, s string) {
	b.WriteString(e.baseEncoder.String(`"`+e.escapeString(s)) + `"`)
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

// escapeString processes the input string by escaping special and control characters to ensure it is safe
// for JSON encoding. It replaces characters like backslashes, quotes, and control characters with their corresponding
// escape sequences. If the string contains non-ASCII characters, it ensures they are properly escaped or replaced
// with the Unicode replacement character if invalid.
//
//nolint:gocyclo,mnd,gocognit
func (e *JSONEncoder) escapeString(s string) string {
	cached, ok := e.escapedStringsCache.Get(s)
	if ok {
		return cached
	}

	buf := builderpool.Get()
	for _, r := range s {
		switch r {
		case '\\', '"':
			buf.WriteByte('\\')
			buf.WriteRune(r)
		case '\b':
			buf.WriteByte('\\')
			buf.WriteByte('b')
		case '\f':
			buf.WriteByte('\\')
			buf.WriteByte('f')
		case '\n':
			buf.WriteByte('\\')
			buf.WriteByte('n')
		case '\r':
			buf.WriteByte('\\')
			buf.WriteByte('r')
		case '\t':
			buf.WriteByte('\\')
			buf.WriteByte('t')
		default:
			if (r >= 0x00 && r <= 0x1F) || r == 0x7F {
				buf.WriteString(escapeControlChar(r))

				continue
			}

			if r > 0x7F {
				if utf8.ValidRune(r) && r != utf8.RuneError {
					buf.WriteString(escapeControlChar(r))
				} else {
					buf.WriteRune('\uFFFD')
				}

				continue
			}

			buf.WriteRune(r)
		}
	}

	res := buf.String()
	builderpool.Put(buf)

	e.escapedStringsCache.Set(s, res)

	return res
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
