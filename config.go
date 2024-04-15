package censor

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// DefaultMaskValue is used to mask struct fields by default.
const DefaultMaskValue = "[CENSORED]"

// Config describes the available parser.Parser and formatter.Formatter configuration.
type Config struct {
	// PrintConfigOnInit sets whether to print the configuration on initialization stage.
	// If true, on Processor initialization, the configuration will be printed to stdout.
	// The default value is true.
	PrintConfigOnInit bool `yaml:"print-config-on-init"`

	// MaskValue is used to mask struct fields with sensitive data.
	// The default value is stored in DefaultMaskValue constant.
	MaskValue string `yaml:"mask-value"`

	// ExcludePatterns contains regexp patterns that are used for the selection
	// of strings that must be masked.
	ExcludePatterns []string `yaml:"exclude-patterns"`
}

// Default returns a default configuration.
func Default() Config {
	return Config{
		PrintConfigOnInit: true,
		MaskValue:         DefaultMaskValue,
		ExcludePatterns:   nil,
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
