package encoder

import (
	"bytes"
	"encoding"
	"encoding/json"
	"reflect"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/shopspring/decimal"

	"github.com/vpakhuchyi/censor/internal/builderpool"
)

// NewJSONEncoder returns a new instance of JSONEncoder with given configuration.
func NewJSONEncoder(c Config) *JSONEncoder {
	e := &JSONEncoder{
		baseEncoder: baseEncoder{
			CensorFieldTag:    DefaultCensorFieldTag,
			ExcludePatterns:   c.ExcludePatterns,
			MaskValue:         strconv.Quote(c.MaskValue),
			structFieldsCache: newFieldsCache(defaultMaxCacheSize),
		},
		enableEscaping: c.EnableJSONEscaping,
	}

	if len(e.ExcludePatterns) != 0 {
		e.ExcludePatternsCompiled = compileRegexpPatterns(e.ExcludePatterns)
	}

	return e
}

// JSONEncoder is used to encode data to JSON format.
type JSONEncoder struct {
	baseEncoder

	enableEscaping bool
}

// Field is a struct that contains information about a struct field.
type Field struct {
	Name     string
	IsMasked bool
}

//nolint:exhaustive,gocyclo
func (e *JSONEncoder) Encode(b *strings.Builder, f reflect.Value) {
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
		b.WriteString(`"` + UnsupportedTypeTmpl + k.String() + `"`)
	}
}

// Struct encodes a struct value to JSON format.
// Note: this method panics if the provided value is not a map.
func (e *JSONEncoder) Struct(b *strings.Builder, v reflect.Value) {
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
			b.WriteString(e.MaskValue)
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
		b.WriteString(`"`)

		if e.enableEscaping {
			b.WriteString(escapeString(field.Name))
		} else {
			b.WriteString(field.Name)
		}
		b.WriteString(`": `)

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
func (e *JSONEncoder) Map(b *strings.Builder, v reflect.Value) {
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
func (e *JSONEncoder) Slice(b *strings.Builder, v reflect.Value) {
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
// In case of a pointer to unsupported type of value, a string built from UnsupportedTypeTmpl
// is used instead of the real value. That string contains a type of the value.
// Note: this method panics if the provided value is not an interface.
func (e *JSONEncoder) Interface(b *strings.Builder, v reflect.Value) {
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
// In case of a pointer to unsupported type of value, a string built from UnsupportedTypeTmpl
// is used instead of the real value. That string contains a type of the value.
// Note: this method panics if the provided value is not a pointer.
func (e *JSONEncoder) Ptr(b *strings.Builder, v reflect.Value) {
	if v.Kind() != reflect.Ptr {
		panic("provided value is not a pointer")
	}

	if v.IsNil() {
		b.WriteString("null")

		return
	}

	e.Encode(b, v.Elem())
}

// String encodes a string value to JSON format.
// If the string matches one of the ExcludePatterns, it will be masked with the MaskValue.
func (e *JSONEncoder) String(b *strings.Builder, s string) {
	if len(e.ExcludePatterns) != 0 && e.ExcludePatternsCompiled != nil {
		matches := e.ExcludePatternsCompiled.FindAllStringIndex(s, -1)
		if len(matches) > 0 {
			lastIndex := 0
			for _, match := range matches {
				start, end := match[0], match[1]
				b.WriteString(s[lastIndex:start] + e.MaskValue)
				lastIndex = end
			}

			b.WriteString(s[lastIndex:])

			return
		}
	}

	if e.enableEscaping {
		b.WriteString(`"`)
		b.WriteString(escapeString(s))
		b.WriteString(`"`)
	} else {
		b.WriteString(`"`)
		b.WriteString(s)
		b.WriteString(`"`)
	}
}

//nolint:exhaustive
func (e *JSONEncoder) encodeMapKey(b *strings.Builder, f reflect.Value) {
	switch k := f.Kind(); k {
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
		if v, ok := f.Interface().(encoding.TextMarshaler); ok {
			b.WriteString(PrepareTextMarshalerValue(v))
		} else {
			b.WriteString(UnsupportedTypeTmpl + k.String())
		}
	}
}

//nolint:gocyclo,mnd,gocognit
func escapeString(s string) string {
	var buf bytes.Buffer
	for _, r := range s {
		switch r {
		case '"':
			buf.WriteString(`\"`)
		case '\\':
			buf.WriteString(`\\`)
		case '\b':
			buf.WriteString(`\b`)
		case '\f':
			buf.WriteString(`\f`)
		case '\n':
			buf.WriteString(`\n`)
		case '\r':
			buf.WriteString(`\r`)
		case '\t':
			buf.WriteString(`\t`)
		default:
			if (r >= 0x00 && r <= 0x1F) || r == 0x7F {
				buf.WriteString(escapeControlChar(r))

				continue
			}

			if r > 0x7F {
				if utf8.ValidRune(r) && r != utf8.RuneError {
					buf.WriteString(escapeControlChar(r))
				} else {
					buf.WriteString(`\uFFFD`)
				}

				continue
			}

			buf.WriteRune(r)
		}
	}

	return buf.String()
}

func escapeControlChar(r rune) string {
	hexStr := strconv.FormatInt(int64(r), 16)
	for len(hexStr) < 4 {
		hexStr = "0" + hexStr
	}

	return `\u` + hexStr
}
