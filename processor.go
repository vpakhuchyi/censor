package censor

import (
	"fmt"
	"reflect"
	"strings"

	jsoniter "github.com/json-iterator/go"
	"gopkg.in/yaml.v3"
)

const jsonFormatErrorPattern = `{"error":"%s"}`

// Processor is used to censor any value and format it into a string representation.
type Processor struct {
	cfg       Config
	formatter jsoniter.API
}

// Censor pkg contains a global instance of Processor.
// This globalInstance is used by the package-level functions.
var globalInstance = New()

// New returns a new instance of Processor with default configuration.
func New() *Processor {
	c := Default()
	formatter := jsoniter.Config{EscapeHTML: true}.Froze()

	formatter.RegisterExtension(newExtension(&c))

	p := Processor{
		cfg:       c,
		formatter: formatter,
	}

	p.PrintConfig()

	return &p
}

// NewWithOpts returns a new instance of Processor, options can be passed to it.
// If no options are passed, the default configuration will be used.
func NewWithOpts(opts ...Option) (*Processor, error) {
	var optCfg OptsConfig
	for _, opt := range opts {
		opt(&optCfg)
	}

	cfg := optCfg.config

	if cfg == nil && optCfg.configPath != "" {
		c, err := ConfigFromFile(optCfg.configPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read the configuration: %w", err)
		}

		cfg = &c
	}

	if cfg == nil {
		c := Default()
		cfg = &c
	}

	formatter := jsoniter.Config{EscapeHTML: true}.Froze()

	formatter.RegisterExtension(newExtension(cfg))

	p := Processor{
		formatter: formatter,
		cfg:       *cfg,
	}

	p.PrintConfig()

	return &p, nil
}

// PrintConfig prints the configuration of the censor Processor.
//
//nolint:gomnd
func (p *Processor) PrintConfig() {
	const lineLength = 69

	var b strings.Builder

	writeLine := func() {
		b.WriteString(strings.Repeat("-", lineLength) + "\n")
	}

	// Handle the case when the censor instance isn't initialized.
	if p == nil {
		const text = "Censor instance isn't initialized"

		writeLine()
		b.WriteString(strings.Repeat(" ", (lineLength-len(text))/2) + text + "\n")
		writeLine()

		fmt.Print(b.String())

		return
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
	d, _ := yaml.Marshal(p.cfg)

	b.Write(d)

	writeLine()

	fmt.Print(b.String())
}

// Format takes any value and returns a string representation of it masking struct fields by default.
// To override this behaviour, use the `censor:"display"` tag.
// Formatting is done recursively for all nested structs/slices/arrays/pointers/maps/interfaces.
// If a value implements [encoding.TextMarshaler], the result of MarshalText is written.
func (p *Processor) Format(val any) string {
	if val == nil || reflect.TypeOf(val) == nil {
		return "nil"
	}

	r, err := p.formatter.Marshal(val)
	if err != nil {
		return fmt.Sprintf(jsonFormatErrorPattern, err.Error())
	}

	return string(r)
}

/*
	Pkg-level functions that work with the global instance of Processor.
*/

// SetGlobalInstance sets a given Processor as a global instance.
func SetGlobalInstance(p *Processor) {
	globalInstance = p
}

// GetGlobalInstance returns a global instance of Processor.
func GetGlobalInstance() *Processor {
	return globalInstance
}
