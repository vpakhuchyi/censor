# Censor

[![GoReportCard example](https://goreportcard.com/badge/github.com/vpakhuchyi/censor)](https://goreportcard.com/report/github.com/vpakhuchyi/censor)
![coverage](https://raw.githubusercontent.com/vpakhuchyi/censor/badges/.badges/main/coverage.svg)
[![GoDoc](https://godoc.org/github.com/vpakhuchyi/censor?status.svg)](https://godoc.org/github.com/vpakhuchyi/censor)

**Censor** is a Go library focused on formatting values into strings, emphasizing the protection
of sensitive information. Through advanced reflection and specialized formatters, it provides precise,
easily readable output. Ideal for safeguarding confidential data or enhancing data presentation in Go projects.

<!-- TOC -->

* [Censor](#censor)
    * [Installation](#installation)
    * [Usage](#usage)
    * [Configuration](#configuration)
    * [Supported Types](#supported-types)

<!-- TOC -->

### Installation

```bash
go get -u github.com/vpakhuchyi/censor
```

### Usage

We can use `censor` to mask all the fields values by default and only display those fields that has
specified `censor:"display"` tag.
This approach will help us to be sure that we don't log sensitive information by mistake.

Censor can be used as a global package-level variable or as a new instance of `censor.Processor`.
Both approaches offer the same functionality. In this example we're using censor as a global package-level variable.

```go
package main

import (
	"log/slog"

	"github.com/vpakhuchyi/censor"
)

type request struct {
	UserID   string  `censor:"display"` // Display value.
	Address  address `censor:"display"` // Display value.
	Email    string  // Mask value.
	FullName string  // Mask value.
}

type address struct {
	City    string `json:"city" censor:"display"`    // Display value.
	Country string `json:"country" censor:"display"` // Display value.
	Street  string `json:"street"`                   // Mask value.
	Zip     int    `json:"zip"`                      // Mask value.
}

// Here is a request struct that contains sensitive information: Email, FullName and Password.
// We could log only UserID, but it's much easier to control what we're logging by using censor 
// instead of checking each log line and making sure that we're not logging sensitive information.
func main() {
	r := request{
		UserID:   "123",
		Address:  address{City: "Kharkiv", Country: "UA", Street: "Nauky Avenue", Zip: 23335},
		Email:    "viktor.example.email@ggmail.com",
		FullName: "Viktor Pakhuchyi",
	}

	// In this example we're using censor as a global package-level variable with default configuration.
	slog.Info("Request", "payload", censor.Format(r))
}

// Here is what we'll see in the log:
Output: `2038/10/25 12:00:01 INFO Request payload="{UserID: 123, Address: {City: Kharkiv, Country: UA, Street: [******], Zip: [******]}, Email: [******], FullName: [******]}`

// All the fields values are masked by default (recursively) except 
// those fields that has specified `censor:"display"` tag.

```

### Configuration

All configuration options can be set using the `censor` package-level functions as shown below.
At the same time you can create a new instance of `censor.Processor` and use its methods to configure it.

| Global option                    | Description                                          |
|----------------------------------|------------------------------------------------------|
| censor.SetMaskValue(s string)    | Set custom mask value instead of default `[******]`. |
| censor.SetFieldTag(s string)     | Set custom field tag instead of default `censor`.    |
| censor.UseJSONTagName(b bool)    | Use JSON tag name instead of struct field name.      |
| censor.DisplayStructName(b bool) | Display struct name in the output.                   |
| censor.DisplayMapType(b bool)    | Display map type in the output.                      |

### Supported Types

| Type                                                                     | Description                                                                                                                                                                                                                                        |
|--------------------------------------------------------------------------|----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| struct                                                                   | By default, all fields within a struct will be masked. If you need to override this behavior for specific fields, you can use the `censor:"display"` tag. It's important to note that all nested fields must also be tagged for proper displaying. |
| map                                                                      | Map values are recursively parsed, ensuring the output is properly formatted.                                                                                                                                                                      |
| slice/array                                                              | These types are recursively parsed, ensuring the output is properly formatted.                                                                                                                                                                     |
| pointer                                                                  | Pointer values, just like slices and arrays, are recursively parsed.                                                                                                                                                                               |
| string                                                                   | String values are handled with no additional formatting.                                                                                                                                                                                           |
| float64/float32                                                          | Floating-point types are formatted to include up to 15 (float64) and 7 (float32) precision digits respectively.                                                                                                                                    |
| int/int8/int16/int32/int64/rune<br/>uint/uint8/uint16/uint32/uint64/byte | All integer types are supported, offering a wide range of options for your data.                                                                                                                                                                   |
| bool                                                                     | Boolean values are handled with no additional formatting.                                                                                                                                                                                          |
| interface                                                                | Will be formatted using the same rules as its underlying type.                                                                                                                                                                                     |
| complex64/complex128                                                     | Both parts are formatted to include up to 15 (complex128) and 7 (complex64) precision digits respectively.                                                                                                                                         |
