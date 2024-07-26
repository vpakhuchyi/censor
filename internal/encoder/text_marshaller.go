package encoder

import "encoding"

// PrepareTextMarshalerValue marshals a value that implements [encoding.TextMarshaler] interface to string.
func PrepareTextMarshalerValue(tm encoding.TextMarshaler) string {
	data, err := tm.MarshalText()
	if err != nil {
		return "!ERROR:" + err.Error()
	}

	return string(data)
}
