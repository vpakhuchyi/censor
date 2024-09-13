package encoder

import (
	"encoding/json"
)

// PrepareJSONMarshalerValue marshals a value that implements [json.Marshaler] interface to string.
func PrepareJSONMarshalerValue(jm json.Marshaler) string {
	data, err := jm.MarshalJSON()
	if err != nil {
		return "!ERROR:" + err.Error()
	}

	return string(data)
}
