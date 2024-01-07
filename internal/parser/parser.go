package parser

import (
	"github.com/vpakhuchyi/censor/config"
)

// DefaultCensorFieldTag is a default tag name for censor fields.
const DefaultCensorFieldTag = "censor"

// Parser is a struct that contains options for parsing.
type Parser struct {
	// useJSONTagName sets whether to use the `json` tag to get the name of the struct field.
	// If no `json` tag is present, the name of the struct field is used.
	useJSONTagName bool
	// censorFieldTag is a tag name for censor fields.
	// The default value is stored in the DefaultCensorFieldTag constant.
	censorFieldTag string
}

// New returns a new instance of Parser with default configuration.
func New() *Parser {
	return &Parser{
		useJSONTagName: false,
		censorFieldTag: DefaultCensorFieldTag,
	}
}

// NewWithConfig returns a new instance of Parser with given configuration.
func NewWithConfig(p config.Parser) *Parser {
	return &Parser{
		useJSONTagName: p.UseJSONTagName,
		censorFieldTag: DefaultCensorFieldTag,
	}
}
