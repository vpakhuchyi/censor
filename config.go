package censor

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/vpakhuchyi/censor/internal/encoder"
)

const (
	// DefaultMaskValue is used to mask struct fields by default.
	DefaultMaskValue = "[CENSORED]"

	// OutputFormatJSON is used to set the output format to JSON.
	OutputFormatJSON = "json"
	// OutputFormatText is used to set the output format to text.
	OutputFormatText = "text"
)

// Config describes the available encoder.Encoder and formatter.Formatter configuration.
type Config struct {
	General General       `yaml:"general"`
	Encoder EncoderConfig `yaml:"encoder"`
}

// General describes general configuration settings.
type General struct {
	// OutputFormat sets the output format: "text" or "json".
	// The default value is "text".
	OutputFormat string `yaml:"output-format"`
	// PrintConfigOnInit sets whether to print the configuration on initialization stage.
	// If true, on Processor initialization, the configuration will be printed to stdout.
	// The default value is true.
	PrintConfigOnInit bool `yaml:"print-config-on-init"`
}

// EncoderConfig describes censor Encoder configuration.
type EncoderConfig = encoder.Config

// DefaultConfig returns a default configuration.
func DefaultConfig() Config {
	return Config{
		General: General{
			OutputFormat:      OutputFormatText,
			PrintConfigOnInit: true,
		},
		Encoder: EncoderConfig{
			DisplayMapType:       false,
			DisplayPointerSymbol: false,
			DisplayStructName:    false,
			EnableJSONEscaping:   true,
			ExcludePatterns:      nil,
			MaskValue:            DefaultMaskValue,
			UseJSONTagName:       false,
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
//nolint:mnd
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
