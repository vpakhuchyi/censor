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

### Features

- [x] Struct formatting with a default values masking of all the fields (recursively).
- [x] Strings values masking based on provided regexp patterns.
- [x] Wide range of supported types:
    - `struct`, `map`, `slice`, `array`, `pointer`, `string`
    - `float64/float32`, `int/int8/int16/int32/int64/rune`
    - `uint/uint8/uint16/uint32/uint64/byte`, `bool`, `interface`
- [x] Customizable configuration.

### Installation

```bash
go get -u github.com/vpakhuchyi/censor
```

### Documentation

Find detailed documentation and practical examples for **Censor** at https://vpakhuchyi.github.io/censor.   
Explore how to use and configure the library effectively, making the most of its features. 

### Usage

**Censor** is a versatile tool designed to mask sensitive information in your Go applications, ensuring that
only specified data is displayed. It can be seamlessly integrated into your code to enhance security,
particularly in scenarios like logging where inadvertent exposure of sensitive data is a concern.

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
Output: `2038/10/25 12:00:01 INFO Request payload={UserID: 123, Address: {City: Kharkiv, Country: UA, Street: [CENSORED], Zip: [CENSORED]}, Email: [CENSORED], FullName: [CENSORED]}`

// All the fields values are masked by default (recursively) except 
// those fields that has specified `censor:"display"` tag.

```
