# Censor

<p align="center"><img src="https://github.com/vpakhuchyi/censor/blob/main/static/logo.png?raw=true" width="260"></p>

<p align="center">
  <a href="https://goreportcard.com/report/github.com/vpakhuchyi/censor"><img src="https://goreportcard.com/badge/github.com/vpakhuchyi/censor" alt="PkgGoDev"></a>
  <img src="https://raw.githubusercontent.com/vpakhuchyi/censor/badges/.badges/main/coverage.svg">
  <a href="https://godoc.org/github.com/vpakhuchyi/censor"><img src="https://godoc.org/github.com/vpakhuchyi/censor?status.svg" alt="Go Report Card" /></a>
</p>

**Censor** is a Go library focused on formatting values into strings, emphasizing the protection
of sensitive information. Through advanced reflection and specialized formatters, it provides precise,
easily readable output. Ideal for safeguarding confidential data or enhancing data presentation in Go projects.

<!-- TOC -->

* [Censor](#censor)
    * [Installation](#installation)
    * [Usage](#usage)
    * [Global Package-Level Usage](#global-package-level-usage)
    * [Instance-Level Usage](#instance-level-usage)
    * [Configuration](#configuration)
    * [Supported Types](#supported-types)

<!-- TOC -->

### Installation

```bash
go get -u github.com/vpakhuchyi/censor
```

### Usage

**Censor** is a versatile tool designed to mask sensitive information in your Go applications, ensuring that
only specified data is displayed. It can be seamlessly integrated into your code to enhance security,
particularly in scenarios like logging where inadvertent exposure of sensitive data is a concern.

**Note**: this package uses reflection, which can be slow. It is not recommended to use this package
in performance-critical scenarios.

#### Global Package-Level Usage

You can use **Censor** as a global package-level variable, allowing for easy and widespread integration across your
codebase. This approach ensures consistent field masking throughout your application.

```go
package main

import (
	"log/slog"

	"github.com/vpakhuchyi/censor"
)

type request struct {
	UserID   string  `censor:"display"` // Display value.
	Address  address `censor:"display"`
	Email    string  // Mask value.
	FullName string
}

type address struct {
	City    string `censor:"display"`
	Country string `censor:"display"`
	Street  string
	Zip     int
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

#### Instance-Level Usage

Alternatively, you have the flexibility to create a new instance of `censor.Processor` for specific use cases
or to customize behavior.

```go
package main

import "log/slog"

type address struct {
	City    string `censor:"display"`
	Country string `censor:"display"`
	Street  string
}

func main() {
	// Create a new instance of censor.Processor.
	p := censor.NewProcessor()

	v := address{City: "Kharkiv", Country: "UA", Street: "Nauky Avenue"}

	slog.Info("Request", "payload", p.Format(v))
}

// Here is what we'll see in the log:
Output: `2038/10/25 12:00:01 INFO Request payload="{City: Kharkiv, Country: UA, Street: [******]}`

```

Both approaches offer the same powerful functionality, allowing you to choose the level of integration that best suits
your application's requirements. Whether you opt for the global package-level variable or create a custom instance,
censor empowers you to confidently manage and safeguard sensitive information.

### Configuration

All configuration options can be set using the package-level functions as shown below.
At the same time you can create a new instance of `censor.Processor` and use its methods to configure it.

| Global option                                 | Description                                             |
|-----------------------------------------------|---------------------------------------------------------|
| censor.SetMaskValue(s string)                 | Set custom mask value instead of default `[******]`.    |
| censor.UseJSONTagName(b bool)                 | Use JSON tag name instead of struct field name.         |
| censor.DisplayStructName(b bool)              | Display struct name in the output.                      |
| censor.DisplayMapType(b bool)                 | Display map type in the output.                         |
| censor.AddExcludePatterns(patterns ...string) | Add regexp patterns for matched strings values masking. |
| censor.SetExcludePatterns(patterns ...string) | Set regexp patterns for matched strings values masking. |
|                                               | Note: it overrides the old patterns.                    |

Apart from this, it's possible to define a configuration using `config.Config` struct.

```go
package main

import (
	"log/slog"

	"github.com/vpakhuchyi/censor"
	"github.com/vpakhuchyi/censor/config"
)

type user struct {
	ID    string `censor:"display"`
	Email string `censor:"display"`
}

func main() {
	// Describe the configuration.
	cfg := config.Config{
		Parser: config.Parser{
			UseJSONTagName: false,
		},
		Formatter: config.Formatter{
			MaskValue:         "[####]",
			DisplayStructName: false,
			DisplayMapType:    false,
			// ExcludePatterns is a list of regular expressions that will be used to exclude fields from masking.
			// In this example we're masking all the email addresses.
			ExcludePatterns: []string{`\b[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Z|a-z]{2,}\b`},
		},
	}

	// Create a new instance of censor.Processor with the specified configuration.
	p := censor.NewWithConfig(cfg)

	v := []user{
		{ID: "123", Email: "user1@exxample.com"},
		{ID: "456", Email: "user2@exxample.com"},
	}

	slog.Info("Request", "payload", p.Format(v))
}

// Here is what we'll see in the log:
Output: `2038/10/25 12:00:01 INFO Request payload="[{ID: 123, Email: [####]}, {ID: 456, Email: [####]}]"

```

### Supported Types

| Type                                                                     | Description                                                                                                                                                                                                                                                                           |
|--------------------------------------------------------------------------|---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| struct                                                                   | By default, all fields within a struct will be masked. If you need to override this behavior for specific fields, you can use the `censor:"display"` tag. It's important to note that all nested fields must also be tagged for proper displaying. All unexported fields are ignored. |
| map                                                                      | Map values are recursively parsed, ensuring the output is properly formatted.                                                                                                                                                                                                         |
| slice/array                                                              | These types are recursively parsed, ensuring the output is properly formatted.                                                                                                                                                                                                        |
| pointer                                                                  | Pointer values, just like slices and arrays, are recursively parsed.                                                                                                                                                                                                                  |
| string                                                                   | String values are handled with no additional formatting.                                                                                                                                                                                                                              |
| float64/float32                                                          | Floating-point types are formatted to include up to 15 (float64) and 7 (float32) significant figures respectively.                                                                                                                                                                    |
| int/int8/int16/int32/int64/rune<br/>uint/uint8/uint16/uint32/uint64/byte | All integer types are supported, offering a wide range of options for your data.                                                                                                                                                                                                      |
| bool                                                                     | Boolean values are handled with no additional formatting.                                                                                                                                                                                                                             |
| interface                                                                | Will be formatted using the same rules as its underlying type.                                                                                                                                                                                                                        |
| complex64/complex128                                                     | Both parts are formatted to include up to 15 (complex128) and 7 (complex64) precision digits respectively.                                                                                                                                                                            |
