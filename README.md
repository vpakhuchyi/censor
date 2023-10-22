# Sanitiser

[![GoReportCard example](https://goreportcard.com/badge/github.com/vpakhuchyi/sanitiser)](https://goreportcard.com/report/github.com/vpakhuchyi/sanitiser)
![coverage](https://raw.githubusercontent.com/vpakhuchyi/sanitiser/badges/.badges/main/coverage.svg)
[![GoDoc](https://godoc.org/github.com/vpakhuchyi/sanitiser?status.svg)](https://godoc.org/github.com/vpakhuchyi/sanitiser)

**Sanitiser** is a powerful Go library with the primary objective of formatting any given value into a string while
effectively masking sensitive information. Leveraging reflection for in-depth analysis and employing formatters, it
ensures accurate and readable output.

### Installation

```bash
go get -u github.com/vpakhuchyi/sanitiser
```

### Supported Types

| Type                                                                     | Description                                                                                                                                                                                                                                     |
|--------------------------------------------------------------------------|-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| struct                                                                   | By default, all fields within a struct will be masked. If you need to override this behavior for specific fields, you can use the `log:"display"` tag. It's important to note that all nested fields must also be tagged for proper displaying. |
| map                                                                      | Map values are recursively parsed, ensuring the output is properly formatted.                                                                                                                                                                   |
| slice/array                                                              | These types are recursively parsed, ensuring the output is properly formatted.                                                                                                                                                                  |
| pointer                                                                  | Pointer values, just like slices and arrays, are recursively parsed.                                                                                                                                                                            |
| string                                                                   | String values are wrapped in double quotes.                                                                                                                                                                                                     |
| float64/float32                                                          | Floating-point types are formatted to include up to 15 (float64) and 7 (float32) precision digits respectively.                                                                                                                                 |
| int/int8/int16/int32/int64/rune<br/>uint/uint8/uint16/uint32/uint64/byte | All integer types are supported, offering a wide range of options for your data.                                                                                                                                                                |
| bool                                                                     | Boolean values are handled with no additional formatting.                                                                                                                                                                                       |

### Usage

The `Format` function is at the heart of this library, providing a versatile method to convert various types into a
formatted string.

### Examples

Here are some examples of how to use the `Format` function.

#### 1. Simple Struct

By default, all fields within a struct will be masked.
This allows to avoid accidental logging of newly added fields that might contain sensitive information.
That's why `Street` field is masked - it doesn't have `log:"display"` tag.

```go
package main

import (
	"fmt"

	"github.com/vpakhuchyi/sanitiser"
)

type address struct {
	City   string `log:"display"`
	State  string `log:"display"`
	Street string
	Zip    int `log:"display"`
}

func main() {
	v := address{
		City:   "San Francisco",
		State:  "CA",
		Street: "451 Main St",
		Zip:    "55501",
	}

	fmt.Println(sanitiser.Format(v))
}

Output: `main.address{"City": "San Francisco", "State": "CA", "Street": "[******]", "Zip": 55501}`

```

#### 2. Struct with Complex Types

```go
package main

import (
	"fmt"

	"github.com/vpakhuchyi/sanitiser"
)

type address struct {
	City   string `log:"display"`
	State  string `log:"display"`
	Street string
	Zip    string `log:"display"`
}

type structWithComplexFields struct {
	Slice       []address `log:"display"`
	MaskedSlice []address
	Map         map[string]address `log:"display"`
	Array       [2]address         `log:"display"`
	Ptr         *address           `log:"display"`
	Struct      address            `log:"display"`
}

func main() {
	v := structWithComplexFields{
		Slice: []address{
			{City: "San Francisco", State: "CA", Street: "451 Main St", Zip: "55501"},
			{City: "Denver", State: "DN", Street: "65 Best St", Zip: "55502"},
		},
		MaskedSlice: []address{
			{City: "New York", State: "NY", Street: "123 Park Ave", Zip: "10001"},
			{City: "Chicago", State: "IL", Street: "789 Lake St", Zip: "60601"},
		},
		Map: map[string]address{
			"home":   {City: "San Francisco", State: "CA", Street: "451 Main St", Zip: "55501"},
			"office": {City: "Denver", State: "DN", Street: "65 Best St", Zip: "55502"},
		},
		Array: [2]address{
			{City: "San Francisco", State: "CA", Street: "451 Main St", Zip: "55501"},
			{City: "Denver", State: "DN", Street: "65 Best St", Zip: "55502"},
		},
		Ptr:    &address{City: "San Francisco", State: "CA", Street: "451 Main St", Zip: "55501"},
		Struct: address{City: "San Francisco", State: "CA", Street: "451 Main St", Zip: "55501"},
	}

	fmt.Println(sanitiser.Format(v))
}

Output: `main.structWithComplexFields{"Slice": [main.address{"City": "San Francisco", "State": "CA", "Street": "[******]", "Zip": "[******]"}, main.address{"City": "Denver", "State": "DN", "Street": "[******]", "Zip": "[******]"}], "MaskedSlice": "[******]", "Map": map[string]main.address["home": main.address{"City": "San Francisco", "State": "CA", "Street": "[******]", "Zip": "[******]"}, "office": main.address{"City": "Denver", "State": "DN", "Street": "[******]", "Zip": "[******]"}], "Array": [main.address{"City": "San Francisco", "State": "CA", "Street": "[******]", "Zip": "[******]"}, main.address{"City": "Denver", "State": "DN", "Street": "[******]", "Zip": "[******]"}], "Ptr": &main.address{"City": "San Francisco", "State": "CA", "Street": "[******]", "Zip": "[******]"}, "Struct": main.address{"City": "San Francisco", "State": "CA", "Street": "[******]", "Zip": "[******]"}}`

```

#### 3. Struct with Basic Types

```go
package main

import (
	"fmt"

	"github.com/vpakhuchyi/sanitiser"
)

type structWithPrimitives struct {
	Int64   int64   `log:"display"`
	Int32   int32   `log:"display"`
	Int16   int16   `log:"display"`
	Int8    int8    `log:"display"`
	Int     int     `log:"display"`
	Uint64  uint64  `log:"display"`
	Uint32  uint32  `log:"display"`
	Uint16  uint16  `log:"display"`
	Uint8   uint8   `log:"display"`
	Uint    uint    `log:"display"`
	Bool    bool    `log:"display"`
	Rune    rune    `log:"display"`
	Byte    byte    `log:"display"`
	Float64 float64 `log:"display"`
	Float32 float32 `log:"display"`
	String  string  `log:"display"`
}

func main() {
	v := structWithPrimitives{
		Int64:   -53645354,
		Int32:   -346456,
		Int16:   -23452,
		Int8:    -101,
		Int:     -456345655,
		Uint64:  53645354,
		Uint32:  346456,
		Uint16:  23452,
		Uint8:   101,
		Uint:    456345655,
		Bool:    true,
		Rune:    'a',
		Byte:    1,
		Float64: 1.12341,
		Float32: 1.389,
		String:  "string",
	}

	fmt.Println(sanitiser.Format(v))
}

Output: `main.structWithPrimitives{"Int64": -53645354, "Int32": -346456, "Int16": -23452, "Int8": -101, "Int": -456345655, "Uint64": 53645354, "Uint32": 346456, "Uint16": 23452, "Uint8": 101, "Uint": 456345655, "Bool": true, "Rune": 97, "Byte": 1, "Float64": 1.12341, "Float32": 1.389, "String": "string"}`
```
