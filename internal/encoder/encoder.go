package encoder

import (
	"reflect"
	"regexp"
	"strings"
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

type baseEncoder struct {
	// CensorFieldTag is a tag name for censor fields.
	// The default value is stored in the DefaultCensorFieldTag constant.
	CensorFieldTag string
	// ExcludePatterns contains regexp patterns that are used to identify strings that must be masked.
	ExcludePatterns []string
	// ExcludePatternsCompiled contains already compiled regexp patterns from ExcludePatterns.
	ExcludePatternsCompiled []*regexp.Regexp
	// MaskValue is used to mask struct fields with sensitive data.
	// The default value is stored in config.DefaultMaskValue constant.
	MaskValue string
	// UseJSONTagName sets whether to use the `json` tag to get the name of the struct field.
	// If no `json` tag is present, the name of the struct field is used.
	UseJSONTagName bool
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
	// EnableJSONEscaping specifies if strings escaping must be performed
	// before marshalling to JSON.
	// The default value is true.
	EnableJSONEscaping bool `yaml:"enable-json-escaping"`
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
