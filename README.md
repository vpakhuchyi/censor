# Sanitiser

[![GoReportCard example](https://goreportcard.com/badge/github.com/vpakhuchyi/sanitiser)](https://goreportcard.com/report/github.com/vpakhuchyi/sanitiser)
![coverage](https://raw.githubusercontent.com/vpakhuchyi/sanitiser/badges/.badges/main/coverage.svg)
[![GoDoc](https://godoc.org/github.com/vpakhuchyi/sanitiser?status.svg)](https://godoc.org/github.com/vpakhuchyi/sanitiser)

**Sanitiser** is a Go library focused on formatting values into strings, emphasizing the protection
of sensitive information. Through advanced reflection and specialized formatters, it provides precise,
easily readable output. Ideal for safeguarding confidential data or enhancing data presentation in Go projects.

<!-- TOC -->

* [Sanitiser](#sanitiser)
    * [Installation](#installation)
    * [Usage](#usage)
    * [Configuration](#configuration)
    * [Supported Types](#supported-types)

<!-- TOC -->

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

Output: `2038/10/25 12:00:01 INFO Request payload="{UserID: 123, Email: [******], FullName: [******], Password: [******]}"`

```

### Configuration

All configuration options can be set using the `sanitiser` package-level functions as shown below.
At the same time you can create a new instance of `sanitiser.Sanitiser` and use its methods to configure it.

| Global option                       | Description                                          |
|-------------------------------------|------------------------------------------------------|
| sanitiser.SetMaskValue(s string)    | Set custom mask value instead of default `[******]`. |
| sanitiser.SetFieldTag(s string)     | Set custom field tag instead of default `sanitiser`. |
| sanitiser.UseJSONTagName(b bool)    | Use JSON tag name instead of struct field name.      |
| sanitiser.DisplayStructName(b bool) | Display struct name in the output.                   |
| sanitiser.DisplayMapType(b bool)    | Display map type in the output.                      |

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
