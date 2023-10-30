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

We can use `sanitiser` to mask all the fields values by default and only display those fields that has
specified `sanitiser:"display"` tag.
This approach will help us to be sure that we don't log sensitive information by mistake.

Sanitiser can be used as a global package-level variable or as a new instance of `sanitiser.Sanitiser`.
Both approaches offer the same functionality. In this example we're using sanitiser as a global package-level variable.

```go
package main

import (
  "log/slog"

  "github.com/vpakhuchyi/sanitiser"
)

type request struct {
  UserID   string  `sanitiser:"display"` // Display value.
  Address  address `sanitiser:"display"` // Display value.
  Email    string  // Mask value.
  FullName string  // Mask value.
}

type address struct {
  City    string `json:"city" sanitiser:"display"`    // Display value.
  Country string `json:"country" sanitiser:"display"` // Display value.
  Street  string `json:"street"`                      // Mask value.
  Zip     int    `json:"zip"`                         // Mask value.
}

// Here is a request struct that contains sensitive information: Email, FullName and Password.
// We could log only UserID, but it's much easier to control what we're logging by using sanitiser 
// instead of checking each log line and making sure that we're not logging sensitive information.
func main() {
  r := request{
    UserID:   "123",
    Address:  address{City: "Kharkiv", Country: "UA", Street: "Nauky Avenue", Zip: 23335},
    Email:    "viktor.example.email@ggmail.com",
    FullName: "Viktor Pakhuchyi",
  }

  // In this example we're using sanitiser as a global package-level variable with default configuration.
  slog.Info("Request", "payload", sanitiser.Format(r))
}

// Here is what we'll see in the log:
Output: `2038/10/25 12:00:01 INFO Request payload="{UserID: 123, Address: {City: Kharkiv, Country: UA, Street: [******], Zip: [******]}, Email: [******], FullName: [******]}`

// All the fields values are sanitised by default (recursively) except 
// those fields that has specified `sanitiser:"display"` tag.

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
| interface                                                                | Will be formatted using the same rules as its underlying type.                                                                                                                                                                                        |
