package formatter

import (
	"strings"

	"github.com/vpakhuchyi/censor/internal/models"
)

// Map formats a map into a string.
// If the map contains nested complex types, they are formatted recursively.
// Keys are sorted to make the output deterministic.
// The map type can be added to the output to indicate the map type using the DisplayMapType option.
//
// Supported types:
//
// [basic types]
// - string
// - int, int8, int16, int32, int64
// - uint, uint8, uint16, uint32, uint64
// - float32, float64
// - bool
// - byte - represented as uint8
// - rune - represented as int32
//
// [complex types]
// - struct - formatted recursively
// - slice - struct values are formatted recursively
// - array - struct values are formatted recursively
// - pointer - pointed structure/arrays/slices are formatted recursively.
// - map - struct/slice/array/pointer values are formatted recursively.
// - interface - formatted recursively.
func (f *Formatter) Map(m models.Map) string {
	var buf strings.Builder

	if f.DisplayMapType {
		buf.WriteString(m.Type + "[")
	} else {
		buf.WriteString("map[")
	}

	for i := 0; i < len(m.Values); i++ {
		f.writeValue(&buf, m.Values[i].Key)
		buf.WriteString(": ")
		f.writeValue(&buf, m.Values[i].Value)

		if i < len(m.Values)-1 {
			buf.WriteString(", ")
		}
	}

	buf.WriteString("]")

	return buf.String()
}
