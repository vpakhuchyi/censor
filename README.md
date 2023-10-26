# Sanitiser

[![GoReportCard example](https://goreportcard.com/badge/github.com/vpakhuchyi/sanitiser)](https://goreportcard.com/report/github.com/vpakhuchyi/sanitiser)
![coverage](https://raw.githubusercontent.com/vpakhuchyi/sanitiser/badges/.badges/main/coverage.svg)
[![GoDoc](https://godoc.org/github.com/vpakhuchyi/sanitiser?status.svg)](https://godoc.org/github.com/vpakhuchyi/sanitiser)

**Sanitiser** is a Go library with the primary objective of formatting any given value into a string while
effectively masking sensitive information. Leveraging reflection for in-depth analysis and employing formatters, it
ensures accurate and readable output.

### Installation

```bash
go get -u github.com/vpakhuchyi/sanitiser
```

### Usage

The `Format` function is at the heart of this library, providing a versatile method to convert various types into a
formatted string.

Most popular use case is to use `Format` function with logging tools:

```go
package main

import (
	"log/slog"

	"github.com/vpakhuchyi/sanitiser"
)

// Let's imagine that we have a request that looks like this:
type request struct {
	UserID   string `sanitiser:"display"` // We want to display this field in the log.
	Email    string // We don't want to display this field in the log.
	FullName string // We don't want to display this field in the log.
	Password string // We don't want to display this field in the log.
}

func main() {
	// Our request contain personal information that we don't want to log.
	// So we can use sanitiser to hide sensitive information but still be able to log the request.
	r := request{
		UserID:   "123",
		Email:    "example@ggmail.com",
		FullName: "Frodo Smith",
		Password: "encoded_password",
	}

	// In case we use slog.Logger, we can use sanitiser to format the request before logging it.
	slog.Info("Request", "payload", sanitiser.Format(r))
}

Output: `2038/10/25 12:00:01 INFO Request payload="main.request{UserID: 123, Email: [******], FullName: [******], Password: [******]}"`

```


### Examples

Here are some examples of how to use the `Format` function.

#### 1. Simple Struct

By default, all fields within a struct will be masked.
In case you use Sanitiser for logging, such an approach will help you to hide sensitive information.
Even if you add new fields to the struct, they will be masked automatically

```go
package main

import (
	"fmt"

	"github.com/vpakhuchyi/sanitiser"
)

type address struct {
	City  string
	State string
	Zip   int
}

func main() {
	a := address{City: "Kharkiv", State: "KH", Zip: 55501}
	fmt.Println(sanitiser.Format(a))
}

Output: `{City: [******], State: [******], Zip: [******]}`
```

If you want to display some fields, you have to use `sanitiser:"display"` tag:

```go
package main

import (
	"fmt"

	"github.com/vpakhuchyi/sanitiser"
)

type address struct {
	City  string `sanitiser:"display"`
	State string `sanitiser:"display"`
	Zip   int
}

func main() {
	a := address{City: "Kharkiv", State: "KH", Zip: 55501}
	fmt.Println(sanitiser.Format(a))
}

Output: `{City: Kharkiv, State: KH, Zip: [******]}`
```

It's possible to override the default tag to use your own:

```go
package main

import (
	"fmt"

	"github.com/vpakhuchyi/sanitiser"
)

type address struct {
	City  string `custom:"display"`
	State string `custom:"display"`
	Zip   int
}

func main() {
	sanitiser.SetFieldTag("custom")
	a := address{City: "Kharkiv", State: "KH", Zip: 55501}
	fmt.Println(sanitiser.Format(a))
}

Output: `{City: Kharkiv, State: KH, Zip: [******]}`
```

#### 2. Struct with Complex Types

```go
package main

import (
	"fmt"

	"github.com/vpakhuchyi/sanitiser"
)

type address struct {
	City  string `sanitiser:"display"`
	State string `sanitiser:"display"`
	Zip   int
}

type structWithComplexFields struct {
	Slice       []address `sanitiser:"display"`
	MaskedSlice []address
	Ptr         *address `sanitiser:"display"`
	Struct      address  `sanitiser:"display"`
}

func main() {
	v := structWithComplexFields{
		Slice: []address{
			{City: "Kharkiv", State: "KH", Zip: 55501},
			{City: "Dnipro", State: "DN", Zip: 55502},
		},
		MaskedSlice: []address{
			{City: "Lviv", State: "LV", Zip: 10001},
			{City: "Kyiv", State: "KY", Zip: 60601},
		},
		Ptr:    &address{City: "Kharkiv", State: "KH", Zip: 55501},
		Struct: address{City: "Kharkiv", State: "KH", Zip: 55501},
	}

	fmt.Println(sanitiser.Format(v))
}

Output: `{Slice: [{City: Kharkiv, State: KH, Zip: [******]}, {City: Dnipro, State: DN, Zip: [******]}], MaskedSlice: [******], Ptr: &{City: Kharkiv, State: KH, Zip: [******]}, Struct: {City: Kharkiv, State: KH, Zip: [******]}}`

```

#### 3. Local instance usage

In previous examples, we used package-level functions. But it's possible to create local instance of sanitiser and use
it:

```go
package main

import (
	"fmt"

	"github.com/vpakhuchyi/sanitiser"
)

type address struct {
	City  string `sanitiser:"display"`
	State string `sanitiser:"display"`
	Zip   int
}

func main() {
	s := sanitiser.New()
	a := address{City: "Kharkiv", State: "KH", Zip: 55501}
	fmt.Println(s.Format(a))
}

Output: `{City: Kharkiv, State: KH, Zip: [******]}`
```

### Configuration

```go
package main

import (
	"fmt"

	"github.com/vpakhuchyi/sanitiser"
)

type address struct {
	City  string `sanitiser:"display"`
	State string `sanitiser:"display"`
	Zip   int
}

type structWithComplexFields struct {
	Slice       []address `json:"slice" custom:"display"`
	MaskedSlice []address
	Ptr         *address `json:"ptr" sanitiser:"display"`
	Struct      address  `sanitiser:"display"`
}

func main() {
	// Create local instance of sanitiser.
	s := sanitiser.New()

	// Set custom mask value.
	s.SetMaskValue("[REDACTED]") // pkg-level function: sanitiser.SetMaskValue("[REDACTED]")

	// Use JSON tag name instead of struct field name.
	// Note: fields with absent JSON tag will be ignored.
	s.UseJSONTagName(true) // pkg-level function: sanitiser.UseJSONTagName(true)

	// Display struct name.
	// It displays the struct name (including a package name) in the output.
	s.DisplayStructName(true) // pkg-level function: sanitiser.DisplayStructName(false)

	// Set custom field tag.
	// It may be useful if a default tag (`sanitiser`) is already used in your project.
	s.SetFieldTag("custom") // pkg-level function: sanitiser.SetFieldTag("custom")

	v := structWithComplexFields{
		Slice: []address{
			{City: "Kharkiv", State: "KH", Zip: 55501},
			{City: "Dnipro", State: "DN", Zip: 55502},
		},
		MaskedSlice: []address{
			{City: "Lviv", State: "LV", Zip: 10001},
			{City: "Kyiv", State: "KY", Zip: 60601},
		},
		Ptr:    &address{City: "Kharkiv", State: "KH", Zip: 55501},
		Struct: address{City: "Kharkiv", State: "KH", Zip: 55501},
	}

	fmt.Println(s.Format(v))
}

Output: `main.structWithComplexFields{slice: [main.address{city: Kharkiv, state: KH, zip: [REDACTED]}, main.address{city: Dnipro, state: DN, zip: [REDACTED]}], ptr: [REDACTED]}`

```

### Supported Types

| Type                                                                     | Description                                                                                                                                                                                                                                           |
|--------------------------------------------------------------------------|-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| struct                                                                   | By default, all fields within a struct will be masked. If you need to override this behavior for specific fields, you can use the `sanitiser:"display"` tag. It's important to note that all nested fields must also be tagged for proper displaying. |
| map                                                                      | Map values are recursively parsed, ensuring the output is properly formatted.                                                                                                                                                                         |
| slice/array                                                              | These types are recursively parsed, ensuring the output is properly formatted.                                                                                                                                                                        |
| pointer                                                                  | Pointer values, just like slices and arrays, are recursively parsed.                                                                                                                                                                                  |
| string                                                                   | String values are handled with no additional formatting.                                                                                                                                                                                              |
| float64/float32                                                          | Floating-point types are formatted to include up to 15 (float64) and 7 (float32) precision digits respectively.                                                                                                                                       |
| int/int8/int16/int32/int64/rune<br/>uint/uint8/uint16/uint32/uint64/byte | All integer types are supported, offering a wide range of options for your data.                                                                                                                                                                      |
| bool                                                                     | Boolean values are handled with no additional formatting.                                                                                                                                                                                             |
