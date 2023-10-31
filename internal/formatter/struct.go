package formatter

import (
	"fmt"
	"strings"

	"github.com/vpakhuchyi/censor/internal/models"
)

// Struct formats a struct into a string with masked sensitive fields.
// All fields are masked by default, unless the field has the `display` tag.
// Supported types could be found in README.md/#supported-types.
func (f *Formatter) Struct(s models.Struct) string {
	var buf strings.Builder

	if f.DisplayStructName {
		buf.WriteString(s.Name)
	}

	buf.WriteString("{")

	fields := s.Fields
	for i := 0; i < len(s.Fields); i++ {
		field := fields[i]

		if field.Opts.Display {
			f.writeField(field, &buf)
		} else {
			buf.WriteString(formatField(field.Name, f.MaskValue))
		}

		if i < len(fields)-1 {
			buf.WriteString(", ")
		}
	}

	buf.WriteString("}")

	return buf.String()
}

func formatField(name, val string) string {
	return fmt.Sprintf(`%s: %s`, name, val)
}
