package censor

import (
	"fmt"
	"os"
	"strings"

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

// DefaultConfig returns a default configuration.
func DefaultConfig() Config {
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

// ToString returns a string that contains a description of the Config struct.
// Example:
// ---------------------------------------------------------------------
//
//	Censor is configured with the following settings:
//
// ---------------------------------------------------------------------
// print-config-on-init: true
// mask-value: '[CENSORED]'
// exclude-patterns:
//   - '[0-9]'
//
// ---------------------------------------------------------------------
//
//nolint:gomnd
func (c Config) ToString() string {
	const lineLength = 69

	var b strings.Builder

	writeLine := func() {
		b.WriteString(strings.Repeat("-", lineLength) + "\n")
	}

	const text = "Censor is configured with the following settings:"

	writeLine()
	b.WriteString(strings.Repeat(" ", (lineLength-len(text))/2) + text + "\n")
	writeLine()

	// config.Config and its nested fields contain only supported YAML types, so it must be marshalled successfully.
	// However, if the configuration is changed, it may contain unsupported types. To handle this case,
	// we have tests that check whether the configuration can be marshalled.
	// Because such kind of changes can happen only in the development process and considering that such an issue
	// can't be fixed automatically or by the user, we don't want to fail the application in this case.
	//
	//nolint:errcheck
	d, _ := yaml.Marshal(c)

	b.Write(d)

	writeLine()

	return b.String()
}
