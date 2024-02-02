package parser

// Config describes Parser configuration.
type Config struct {
	// UseJSONTagName sets whether to use the `json` tag to get the name of the struct field.
	// If no `json` tag is present, the name of the struct field is used.
	UseJSONTagName bool `yaml:"use-json-tag-name"`
}
