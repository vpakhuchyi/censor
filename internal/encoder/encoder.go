package encoder

import (
	"reflect"
	"regexp"
	"strings"

	"github.com/vpakhuchyi/censor/internal/builderpool"
	"github.com/vpakhuchyi/censor/internal/cache"
)

// Encoder is an interface that describes the behavior of the encoder.
type Encoder interface {
	Struct(b *strings.Builder, rv reflect.Value)
	Ptr(b *strings.Builder, rv reflect.Value)
	Slice(b *strings.Builder, rv reflect.Value)
	Map(b *strings.Builder, rv reflect.Value)
	Interface(b *strings.Builder, rv reflect.Value)
	String(b *strings.Builder, s string)
	Encode(b *strings.Builder, f reflect.Value)
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
	// The default value is stored in the DefaultCensorFieldTag constant.
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
	structFieldsCache *cache.SliceCache[Field]
	// escapedStringsCache is used to cache escaped strings, to improve performance.
	escapedStringsCache *cache.Cache[string]
	// regexpCache is used to cache compiled regexp patterns, to improve performance.
	regexpCache *cache.Cache[string]
}

// String processes the input string by masking any substrings that match the configured exclusion patterns.
// It replaces matched segments with a predefined mask value to censor sensitive information.
func (e *baseEncoder) String(s string) string {
	res := s
	if len(e.ExcludePatterns) != 0 && e.ExcludePatternsCompiled != nil {
		cached, ok := e.regexpCache.Get(s)
		if ok {
			return cached
		}

		matches := e.ExcludePatternsCompiled.FindAllStringIndex(s, -1)
		if len(matches) > 0 {
			bb := builderpool.Get()
			lastIndex := 0
			for _, m := range matches {
				start, end := m[0], m[1]
				bb.WriteString(s[lastIndex:start] + e.MaskValue)
				lastIndex = end
			}

			bb.WriteString(s[lastIndex:])
			res = bb.String()
		}
	}

	e.regexpCache.Set(s, res)

	return res
}
