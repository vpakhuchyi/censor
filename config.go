package censor

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
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
	// The default value is false.
	PrintConfigOnInit bool `yaml:"print-config-on-init"`
}

// EncoderConfig describes censor Encoder configuration.
type EncoderConfig struct {
	DisplayMapType       bool     `yaml:"display-map-type"`
	DisplayPointerSymbol bool     `yaml:"display-pointer-symbol"`
	DisplayStructName    bool     `yaml:"display-struct-name"`
	ExcludePatterns      []string `yaml:"exclude-patterns"`
	MaskValue            string   `yaml:"mask-value"`
	UseJSONTagName       bool     `yaml:"use-json-tag-name"`
}

func (c EncoderConfig) toEncoderConfig() encoder.Config {
	return encoder.Config{
		DisplayMapType:       c.DisplayMapType,
		DisplayPointerSymbol: c.DisplayPointerSymbol,
		DisplayStructName:    c.DisplayStructName,
		ExcludePatterns:      c.ExcludePatterns,
		MaskValue:            c.MaskValue,
		UseJSONTagName:       c.UseJSONTagName,
	}
}

// DefaultConfig returns a default configuration.
func DefaultConfig() Config {
	return Config{
		General: General{
			OutputFormat:      OutputFormatJSON,
			PrintConfigOnInit: false,
		},
		Encoder: EncoderConfig{
			DisplayMapType:       false,
			DisplayPointerSymbol: false,
			DisplayStructName:    false,
			ExcludePatterns:      nil,
			MaskValue:            DefaultMaskValue,
			UseJSONTagName:       false,
		},
	}
}

const maxRegExPatterns = 50

// Validate checks whether the configuration is valid.
func (c Config) Validate() error {
	switch c.General.OutputFormat {
	case OutputFormatText, OutputFormatJSON:
	default:
		return fmt.Errorf("invalid output format: %q, must be %q or %q", c.General.OutputFormat, OutputFormatText, OutputFormatJSON)
	}

	if c.Encoder.MaskValue == "" {
		return fmt.Errorf("mask value cannot be empty")
	}

	if len(c.Encoder.ExcludePatterns) > maxRegExPatterns {
		return fmt.Errorf("too many exclude patterns (max %d): %d", maxRegExPatterns, len(c.Encoder.ExcludePatterns))
	}

	for _, pattern := range c.Encoder.ExcludePatterns {
		if _, err := regexp.Compile(pattern); err != nil {
			return fmt.Errorf("invalid exclude pattern %q: %w", pattern, err)
		}
	}

	return nil
}

// ConfigFromFile reads a configuration from the given .yml file.
// It returns an error if the file cannot be read or unmarshalled.
func ConfigFromFile(path string) (Config, error) {
	if path == "" {
		return Config{}, fmt.Errorf("config file path cannot be empty")
	}

	cleanPath := filepath.Clean(path)
	if strings.Contains(cleanPath, "..") {
		return Config{}, fmt.Errorf("invalid config file path: directory traversal detected")
	}

	ext := filepath.Ext(cleanPath)
	if ext != ".yml" && ext != ".yaml" {
		return Config{}, fmt.Errorf("config file must have .yml or .yaml extension")
	}

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
