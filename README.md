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

### Usage

**Censor** is a versatile tool designed to mask sensitive information in your Go applications, ensuring that
only specified data is displayed. It can be seamlessly integrated into your code to enhance security,
particularly in scenarios like logging where inadvertent exposure of sensitive data is a concern.

```go
package main

import (
  "fmt"
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
  fmt.Println(censor.Format(r))
}

// Here is what we'll see in the log:
Output: `{UserID: 123, Address: {City: Kharkiv, Country: UA, Street: [CENSORED], Zip: [CENSORED]}, Email: [CENSORED], FullName: [CENSORED]}`

// All the fields values are masked by default (recursively) except 
// those fields that has specified `censor:"display"` tag.
```

### Handler for "log/slog" 

Censor instance can be used as a handler for the `log/slog` package.

```go
  // Define the configuration.  
  cfg := censor.Config{
    Encoder: censor.EncoderConfig{
      DisplayMapType:       true,
      MaskValue:      "[CENSORED]", 
      // Other configuration options...
    },
  }
  
  // Initialize a Censor instance with the specified configuration.
  c, err := censor.NewWithOpts(censor.WithConfig(&cfg))
  if err != nil {
    // Handle error.
  }
  
  // Create and register a new slog handler with the initialized instance.
  opts := []sloghandler.Option{sloghandler.WithCensor(c)}
  log := slog.New(sloghandler.NewJSONHandler(opts...))

  // Use logger as usually.
  log.Info("user", slog.Any("payload", payload))
```

### Handler for "go.uber.org/zap" 

Censor also can be used as a handler for the `go.uber.org/zap` package.

```go
  // Define the configuration.  
  cfg := censor.Config{
    Encoder: censor.EncoderConfig{
      DisplayMapType:       true,
      MaskValue:      "[CENSORED]",
      // Other configuration options...
    },
  }

  // Initialize a Censor instance with the specified configuration.
  c, err := censor.NewWithOpts(censor.WithConfig(&cfg))
  if err != nil {
    // Handle error.
  }
  
  // Wrap the zap core with the Censor handler.
  o := zap.WrapCore(func(core zapcore.Core) zapcore.Core {
    return zaphandler.NewHandler(core, zaphandler.WithCensor(c))
  })

  // Initialize a new zap logger instance.
  l, err := zap.NewProduction(o)
  if err != nil {
    // Handle error.
  }
  
  // Use logger as usually.
  l.Info("user", zap.Any("payload", payload))
```