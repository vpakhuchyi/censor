package options

import "strings"

const display = "display"

type FieldOptions struct {
	Display bool
}

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
