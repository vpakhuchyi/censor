package parser

// DefaultCensorFieldTag is a default tag name for censor fields.
const DefaultCensorFieldTag = "censor"

// Parser is a struct that contains options for parsing.
type Parser struct {
	// useJSONTagName sets whether to use the `json` tag to get the name of the struct field.
	// If no `json` tag is present, the name of struct will be an empty string.
	useJSONTagName bool
	// CensorFieldTag is a tag name for censor fields.
	// The default value is stored in the DefaultCensorFieldTag constant.
	CensorFieldTag string
}

// New returns a new instance of Parser with default configuration.
func New() *Parser {
	return &Parser{
		useJSONTagName: false,
		CensorFieldTag: DefaultCensorFieldTag,
	}
}

// UseJSONTagName sets whether to use the `json` tag to get the name of the struct field.
// If no `json` tag is present, the name of struct will be an empty string.
func (p *Parser) UseJSONTagName(v bool) {
	p.useJSONTagName = v
}
