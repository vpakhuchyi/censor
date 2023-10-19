package options

import "strings"

const display = "display"

// FieldOptions is a struct that holds the options for a field.
type FieldOptions struct {
	// Display is a boolean that determines whether the field should be displayed.
	Display bool
}

// Parse parses the tag and returns the options for the field.
func Parse(tag string) FieldOptions {
	tagValues := strings.Split(tag, ",")

	var opts FieldOptions
	for _, v := range tagValues {
		if v == display {
			opts.Display = true
		}
	}

	return opts
}
