package censor

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"

	"github.com/vpakhuchyi/censor/internal/formatter"
	"github.com/vpakhuchyi/censor/internal/parser"
)

// DefaultMaskValue is used to mask struct fields by default.
const DefaultMaskValue = "[CENSORED]"

// Config describes the available parser.Parser and formatter.Formatter configuration.
type Config struct {
	General   General         `yaml:"general"`
	Parser    ParserConfig    `yaml:"parser"`
	Formatter FormatterConfig `yaml:"formatter"`
}

// General describes general configuration settings.
type General struct {
	// PrintConfigOnInit sets whether to print the configuration on initialization stage.
	// If true, on Processor initialization, the configuration will be printed to stdout.
	// The default value is true.
	PrintConfigOnInit bool `yaml:"print-config-on-init"`
}

// ParserConfig describes censor Parser configuration.
type ParserConfig = parser.Config

// FormatterConfig describes censor Formatter configuration.
type FormatterConfig = formatter.Config

// Default returns a default configuration.
func Default() Config {
	return Config{
		Parser: parser.Config{
			UseJSONTagName: false,
		},
		Formatter: formatter.Config{
			MaskValue:            DefaultMaskValue,
			DisplayPointerSymbol: false,
			DisplayStructName:    false,
			DisplayMapType:       false,
			ExcludePatterns:      nil,
		},
		General: General{
			PrintConfigOnInit: true,
		},
	}
}

// ConfigFromFile reads a configuration from the given .yml file.
// It returns an error if the file cannot be read or unmarshalled.
func ConfigFromFile(path string) (Config, error) {
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
