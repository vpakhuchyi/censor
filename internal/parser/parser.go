package parser

// DefaultCensorFieldTag is a default tag name for censor fields.
const DefaultCensorFieldTag = "censor"

// UnsupportedTypeTmpl is a template for a value that is returned when a given type is not supported.
const UnsupportedTypeTmpl = "[Unsupported type: %s]"

// Parser is a struct that contains options for parsing.
type Parser struct {
	// useJSONTagName sets whether to use the `json` tag to get the name of the struct field.
	// If no `json` tag is present, the name of the struct field is used.
	useJSONTagName bool
	// censorFieldTag is a tag name for censor fields.
	// The default value is stored in the DefaultCensorFieldTag constant.
	censorFieldTag string
}

// New returns a new instance of Parser with given configuration.
func New(c Config) *Parser {
	return &Parser{
		useJSONTagName: c.UseJSONTagName,
		censorFieldTag: DefaultCensorFieldTag,
	}
}
