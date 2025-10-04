package encoder

import (
	"bytes"
	"reflect"
	"regexp"

	"github.com/vpakhuchyi/censor/internal/cache"
)

const (
	defaultCensorFieldTag = "censor"
	unsupportedTypeTmpl   = "unsupported type="
	display               = "display"
)

// Encoder is an interface that describes the behavior of the encoder.
type Encoder interface {
	Struct(b *bytes.Buffer, rv reflect.Value)
	Ptr(b *bytes.Buffer, rv reflect.Value)
	Slice(b *bytes.Buffer, rv reflect.Value)
	Map(b *bytes.Buffer, rv reflect.Value)
	Interface(b *bytes.Buffer, rv reflect.Value)
	String(b *bytes.Buffer, s string)
	Encode(b *bytes.Buffer, f reflect.Value)
}

// Config describes censor Encoder configuration.
type Config struct {
	// DisplayMapType is used to display map type in the output.
	// The default value is false.
	DisplayMapType bool `yaml:"display-map-type"`
	// DisplayPointerSymbol is used to display '&' (pointer symbol) in the output.
	// The default value is false.
	DisplayPointerSymbol bool `yaml:"display-pointer-symbol"`
	// DisplayStructName is used to display struct name in the output.
	// A struct name includes the last part of the package path.
	// The default value is false.
	DisplayStructName bool `yaml:"display-struct-name"`
	// ExcludePatterns contains regexp patterns that are used for the selection
	// of strings that must be masked.
	ExcludePatterns []string `yaml:"exclude-patterns"`
	// MaskValue is used to mask struct fields with sensitive data.
	// The default value is stored in DefaultMaskValue constant.
	MaskValue string `yaml:"mask-value"`
	// UseJSONTagName sets whether to use the `json` tag to get the name of the struct field.
	// If no `json` tag is present, the name of the struct field is used.
	UseJSONTagName bool `yaml:"use-json-tag-name"`
}

type baseEncoder struct {
	// CensorFieldTag is a tag name for censor fields.
	// The default value is stored in the defaultCensorFieldTag constant.
	CensorFieldTag string
	// ExcludePatterns contains regexp patterns that are used to identify strings that must be masked.
	ExcludePatterns []string
	// ExcludePatternsCompiled contains already compiled regexp patterns from ExcludePatterns joined using "|".
	ExcludePatternsCompiled *regexp.Regexp
	// MaskValue is used to mask struct fields with sensitive data.
	// The default value is stored in config.DefaultMaskValue constant.
	MaskValue string
	// UseJSONTagName sets whether to use the `json` tag to get the name of the struct field.
	// If no `json` tag is present, the name of the struct field is used.
	UseJSONTagName bool

	// structFieldsCache is used to cache struct fields, so we don't need to use reflection every time.
	// Note: fields of anonymous structs are not cached due to the absence of a name.
	structFieldsCache *cache.TypeCache[[]Field]
	// escapedStringsCache is used to cache escaped strings, to improve performance.
	escapedStringsCache *cache.Cache[string]
	// regexpCache is used to cache compiled regexp patterns, to improve performance.
	regexpCache *cache.Cache[string]
}

// Field is a struct that contains information about a struct field.
type Field struct {
	Name     string
	IsMasked bool
}

// WriteString processes the input string by masking any substrings that match the configured exclusion patterns.
// It replaces matched segments with a predefined mask value and writes the result directly to the buffer.
func (e *baseEncoder) WriteString(b *bytes.Buffer, s string) {
	if len(e.ExcludePatterns) == 0 || e.ExcludePatternsCompiled == nil {
		b.WriteString(s)

		return
	}

	cached, ok := e.regexpCache.Get(s)
	if ok {
		b.WriteString(cached)

		return
	}

	matches := e.ExcludePatternsCompiled.FindAllStringIndex(s, -1)
	if len(matches) == 0 {
		b.WriteString(s)
		e.regexpCache.Set(s, s)

		return
	}

	startLen := b.Len()
	lastIndex := 0
	for _, m := range matches {
		start, end := m[0], m[1]
		b.WriteString(s[lastIndex:start])
		b.WriteString(e.MaskValue)
		lastIndex = end
	}
	b.WriteString(s[lastIndex:])

	result := b.String()[startLen:]
	e.regexpCache.Set(s, result)
}
