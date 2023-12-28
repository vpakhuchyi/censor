package config

// DefaultMaskValue is used to mask struct fields by default.
const DefaultMaskValue = "[CENSORED]"

// Config describes the available parser.Parser and formatter.Formatter configuration.
type Config struct {
	Parser    Parser
	Formatter Formatter
}

// Parser describes parser.Parser configuration.
type Parser struct {
	// UseJSONTagName sets whether to use the `json` tag to get the name of the struct field.
	// If no `json` tag is present, the name of the struct field is used.
	UseJSONTagName bool
}

// Formatter describes formatter.Formatter configuration.
type Formatter struct {
	// MaskValue is used to mask struct fields with sensitive data.
	// The default value is stored in DefaultMaskValue constant.
	MaskValue string
	// DisplayPointerSymbol is used to display '&' (pointer symbol) in the output.
	// The default value is false.
	DisplayPointerSymbol bool
	// DisplayStructName is used to display struct name in the output.
	// A struct name includes the last part of the package path.
	// The default value is false.
	DisplayStructName bool
	// DisplayMapType is used to display map type in the output.
	// The default value is false.
	DisplayMapType bool
	// ExcludePatterns contains regexp patterns that are used for the selection
	// of strings that must be masked.
	ExcludePatterns []string
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
