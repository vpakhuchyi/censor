package parser

// DefaultSanitiserFieldTag is a default tag name for sanitiser fields.
const DefaultSanitiserFieldTag = "sanitiser"

// Parser is a struct that contains options for parsing.
type Parser struct {
	// UseJSONTagName sets whether to use the `json` tag to get the name of the struct field.
	// If no `json` tag is present, the name of struct will be an empty string.
	UseJSONTagName bool
	// SanitiserFieldTag is a tag name for sanitiser fields.
	// The default value is stored in the DefaultSanitiserFieldTag constant.
	SanitiserFieldTag string
}

// New returns a new instance of Parser with default configuration.
func New() *Parser {
	return &Parser{
		UseJSONTagName:    false,
		SanitiserFieldTag: DefaultSanitiserFieldTag,
	}
}
