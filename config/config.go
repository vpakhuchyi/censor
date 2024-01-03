package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// DefaultMaskValue is used to mask struct fields by default.
const DefaultMaskValue = "[CENSORED]"

// Config describes the available parser.Parser and formatter.Formatter configuration.
type Config struct {
	Parser    Parser    `yaml:"parser"`
	Formatter Formatter `yaml:"formatter"`
}

// Parser describes parser.Parser configuration.
type Parser struct {
	// UseJSONTagName sets whether to use the `json` tag to get the name of the struct field.
	// If no `json` tag is present, the name of the struct field is used.
	UseJSONTagName bool `yaml:"use-json-tag-name"`
}

// Formatter describes formatter.Formatter configuration.
type Formatter struct {
	// MaskValue is used to mask struct fields with sensitive data.
	// The default value is stored in DefaultMaskValue constant.
	MaskValue string `yaml:"mask-value"`
	// DisplayPointerSymbol is used to display '&' (pointer symbol) in the output.
	// The default value is false.
	DisplayPointerSymbol bool `yaml:"display-pointer-symbol"`
	// DisplayStructName is used to display struct name in the output.
	// A struct name includes the last part of the package path.
	// The default value is false.
	DisplayStructName bool `yaml:"display-struct-name"`
	// DisplayMapType is used to display map type in the output.
	// The default value is false.
	DisplayMapType bool `yaml:"display-map-type"`
	// ExcludePatterns contains regexp patterns that are used for the selection
	// of strings that must be masked.
	ExcludePatterns []string `yaml:"exclude-patterns"`
}

// Default returns a default configuration.
func Default() Config {
	return Config{
		Parser: Parser{
			UseJSONTagName: false,
		},
		Formatter: Formatter{
			MaskValue:            DefaultMaskValue,
			DisplayPointerSymbol: false,
			DisplayStructName:    false,
			DisplayMapType:       false,
			ExcludePatterns:      nil,
		},
	}
}

// GetParserConfig returns Parser configuration from Config.
func (c Config) GetParserConfig() Parser {
	return c.Parser
}

// GetFormatterConfig returns Formatter configuration from Config.
func (c Config) GetFormatterConfig() Formatter {
	return c.Formatter
}

// FromFile reads a configuration from the given .yml file.
// It returns an error if the file cannot be read or unmarshalled.
func FromFile(path string) (Config, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return Config{}, fmt.Errorf("failed to read file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(file, &cfg); err != nil {
		return Config{}, fmt.Errorf("failed to unmarshal yaml: %w", err)
	}

	return cfg, nil
}
