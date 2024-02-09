<p align="center"><img src="https://github.com/vpakhuchyi/censor/blob/main/static/logo.png?raw=true" width="260"></p>

<p align="center">
  <a href="https://goreportcard.com/report/github.com/vpakhuchyi/censor"><img src="https://goreportcard.com/badge/github.com/vpakhuchyi/censor" alt="PkgGoDev"></a>
  <img src="https://raw.githubusercontent.com/vpakhuchyi/censor/badges/.badges/main/coverage.svg">
  <a href="https://godoc.org/github.com/vpakhuchyi/censor"><img src="https://godoc.org/github.com/vpakhuchyi/censor?status.svg" alt="Go Report Card" /></a>
</p>

**Censor** is a Go library focused on formatting values into strings, emphasizing the protection
of sensitive information. Through advanced reflection and specialized formatters, it provides precise,
easily readable output. Ideal for enhancing data presentation in Go projects.

## Installation

```bash
go get -u github.com/vpakhuchyi/censor
```

## Features

- [x] Struct formatting with a default values masking of all the fields (recursively).
- [x] Strings values masking based on provided regexp patterns.
- [x] Wide range of supported types:
    - `struct`, `map`, `slice`, `array`, `pointer`, `string`
    - `float64/float32`, `int/int8/int16/int32/int64/rune`
    - `uint/uint8/uint16/uint32/uint64/uintptr/byte`, `bool`, `interface`
- [x] Support encoding.TextMarshaler interface for custom types. 
- [x] Customizable configuration:
    - Using `.ymal` file
    - Passing `censor.Config` struct

## Usage

*Censor* can be seamlessly integrated into your code to enhance security, particularly in scenarios like logging
where inadvertent exposure of sensitive data is a concern.

### Structs formatting

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
  // Here is what we'll see in the log:
  // Output: `2038/10/25 12:00:01 INFO Request payload={UserID: 123, Address: {City: Kharkiv, Country: UA, Street: [CENSORED], Zip: [CENSORED]}, Email: [CENSORED], FullName: [CENSORED]}`
}

```

### Strings formatting based on provided regexp patterns

By default, *censor* operates without any exclude patterns. However, you have the flexibility to enhance its
functionality by incorporating exclude patterns through the `ExcludePatterns` configuration option.

When exclude patterns are introduced, encapsulated formatters compares all string values against the specified patterns,
masking matched parts of such strings. This capability proves invaluable when, for instance, you want to conceal all email addresses
within a set of strings. Simply add an email address pattern, and `censor.Format()` will automatically mask all values
conforming to that pattern. This ensures that even as new values are added in the future, they will be seamlessly
handled.

**Note**: All added patterns are conveniently stored within the corresponding `censor.Processor` instance, whether it's
global or local, ensuring seamless pattern management.

```go
package main

import (
  "log/slog"

  "github.com/vpakhuchyi/censor"
)

type email struct {
  Text string `censor:"display"`
}

func main() {
  // This is a regular expression that matches email addresses.
  const emailPattern = `(?P<email>[\w.%+\-]+@[A-Za-z0-9.\-]+\.[A-Za-z]{2,4})`

  // This regular expression matches IBANs.
  const ibanPattern = `([A-Z]{2}[0-9]{27})`

  cfg := censor.Config{
    Formatter: censor.FormatterConfig{
      MaskValue: censor.DefaultMaskValue,
      // We can add exclude patterns to censor to make sure that the values that match the pattern will be masked.
      ExcludePatterns: []string{emailPattern, ibanPattern},
    },
  }

  // Create a new instance of censor.Processor with the specified configuration and set it as a global processor.
  censor.NewWithOpts(censor.WithConfig(&cfg))

  const msg = `Here are the details of your account: email=user.example.email@ggmail.com, IBAN=UA123456789123456789123456789`

  v := email{Text: msg}

  slog.Info("Request", "payload", censor.Format(v))
  // Here is what we'll see in the log:
  // Output: `2038/10/25 12:00:01 INFO Request payload={Text: Here are the details of your account: email=[CENSORED], IBAN=[CENSORED]}`
}

```

### Instance-Level Usage

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
  p := censor.New()

  v := address{City: "Kharkiv", Country: "UA", Street: "Nauky Avenue"}

  slog.Info("Request", "payload", p.Format(v))
  // Here is what we'll see in the log:
  // Output: `2038/10/25 12:00:01 INFO Request payload={City: Kharkiv, Country: UA, Street: [CENSORED]}`
}

```

Both approaches, global and instance usage offer the same powerful functionality, allowing you to choose the level of
integration that best suits your application's requirements.

## Configuration

Censor supports two ways of configuration: using the `censor.Config` struct and providing a `.yml` configuration file.
All the configuration options are available in both ways.

Table below shows the names of the configuration options:

| Go name              | YML name               | Default value   | Description                                                                                                                                                  |
|----------------------|------------------------|-----------------|--------------------------------------------------------------------------------------------------------------------------------------------------------------|
| PrintConfigOnInit    | print-config-on-init   | true            | If true, the configuration will be printed when any of available constructors is used.                                                                       |
| UseJSONTagName       | use-json-tag-name      | false           | If true, the JSON tag name will be used instead of the Go struct field name.                                                                                 |
| MaskValue            | mask-value             | [CENSORED]      | The value that will be used to mask the sensitive information.                                                                                               |
| DisplayStructName    | display-struct-name    | false           | If true, the struct name will be displayed in the output.                                                                                                    |
| DisplayMapType       | display-map-type       | false           | If true, the map type will be displayed in the output.                                                                                                       |
| DisplayPointerSymbol | display-pointer-symbol | false           | If true, '&' (the pointer symbol) will be displayed in the output.                                                                                           |
| ExcludePatterns      | exclude-patterns       | []              | A list of regular expressions that will be compared against all the string values. <br/>If a value matches any of the patterns, that section will be masked. |

### Using the `censor.Config` struct

It's possible to define a configuration using `censor.Config` struct:

```go
package main

import (
  "github.com/vpakhuchyi/censor"
)

func main() {
  cfg := censor.Config{
    Parser: censor.ParserConfig{
      UseJSONTagName: false,
    },
    Formatter: censor.FormatterConfig{
      MaskValue:         "[####]",
      DisplayStructName: false,
      DisplayMapType:    false,
      ExcludePatterns:   []string{`\b[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Z|a-z]{2,}\b`},
    },
  }

  p, err  := censor.NewWithOpts(censor.WithConfig(&cfg))
}

```

### Providing `.yml` configuration file

It's also possible to provide a configuration file in `.yml` format. In this case, you can use the `./cfg_example.yml`
file as an example:

```go
package main

import (
  "github.com/vpakhuchyi/censor"
)

func main() {
  pathToConfigFile := "./cfg_example.yml"

  // Create a new instance of censor.Processor with the configuration file usage.
  p, err := censor.NewWithOpts(censor.WithConfigPath(pathToConfigFile))
  if err != nil {
    // Handle error.
  }
}

```

## Censor handlers for loggers

### Handler for `go.uber.org/zap`

Package `github.com/vpakhuchyi/censor/handlers/zap` provides a configurable logging handler for `go.uber.org/zap` library, 
allowing users to apply censoring to log entries, overriding the original values before passing them to the core.

Example of Censor handler initialization:
```go
package main

import (
  censorlog "github.com/vpakhuchyi/censor/handlers/zap"
  "go.uber.org/zap"
  "go.uber.org/zap/zapcore"
)

type User struct {
  Name  string `censor:"display"`
  Email string
}

func main() {
  // Create a new Censor handler with the default configuration.
  o := zap.WrapCore(func(core zapcore.Core) zapcore.Core {
    return censorlog.NewHandler(core)
  })

  // Created zap core wraps the original core that allows to keep the original logger configuration.
  l, err := zap.NewProduction(o)
  if err != nil {
    // handle error
  }

  u := User{Name: "John Doe", Email: "example@gmail.com"}

  // Once the handler is initialized and the logger is created with the wrapped core,
  // you can use the logger as usual and the censor handler will process the log entries.
  l.Info("user", zap.Any("payload", u))
  
  // Output: {"level":"info",...,"msg":"user","payload":"{Name: John Doe, Email: [CENSORED]}"}
}

```
Due to the diversity of the `go.uber.org/zap` library usage, please pay attention to the way you use it with the censor handler.

Describing the logger usage we can operate with a few keywords: `msg`, `key` and `value`.
For example, in a call to `l.Info("payload", zap.Any("addresses", []string{"address1", "address2"}))`:
- "payload" is a `msg`
- "addresses" is a `key`
- []string{"address1", "address2"} is a `value`

By default, Censor processes only `value` to minimize the overhead. The `msg` and `key` rarely contain sensitive data.
However, there are predefined configuration options that can be passed to `NewHandler()` function to customize the
censor handler behavior.

The following options are available:
1. `WithCensor(censor *censor.Processor)` - sets the censor processor instance for the logger handler.
   If not provided, a default censor processor is used.
2. `WithMessagesFormat()` - enables the censoring of a log "msg" values if such are present.
3. `WithKeysFormat()` - enables the censoring of log "key" values.

After the handler is initialized, it can be used as a regular logger. Because the censor handler is a wrapper around
`go.uber.org/zap` library logic, it may not be compatible with all the possible ways of the logger usage.

That's why it's recommended to use it with the following constructions:

- `l.Info(msg string, fields ...zap.Field)`

- `l.With(fields ...zap.Field).Info(msg string, fields ...zap.Field)`

  In both cases, Info could be replaced with Debug, Warn, Error, Panic, Fatal.

With sugared logger, the following constructions are supported:

- `l.Info(args ...interface{})`

- `l.Infof(template string, args ...interface{})`

- `l.Infow(msg string, keysAndValues ...interface{})`

- `l.Infoln(args ...interface{})`

- `l.With(args ...interface{}).Info(args ...interface{})`

  In all cases, Info could be replaced with Debug, Warn, Error, Panic, Fatal.

Methods ending in `f`, `ln` and `log.Print-style (l.Info)` in a sugared logger can't use all the
censor handler features. Due to the nature of the zap sugared logger, Censor accepts formatted strings and has no
possibility to use its parsing. However, a features like regexp matching are still available.


## Supported Types

Here are examples of the types that are supported by *censor*. All of these types are handled recursively whenever that
is possible.

### Struct

By default, all fields within a struct will be masked. If you need to override this behavior for specific fields, you
can use the `censor:"display"` tag. It's important to note that all nested fields must also be tagged for a proper
displaying.

**Note**: All unexported fields are ignored.

```go
package main

import (
  "log/slog"

  "github.com/vpakhuchyi/censor"
)

type address struct {
  City    string `censor:"display"`
  Country string `censor:"display"`
  Street  string
}

func main() {
  v := address{City: "Kharkiv", Country: "UA", Street: "Nauky Avenue"}

  slog.Info("Request", "payload", censor.Format(v))
  // Here is what we'll see in the log:
  // Output: `2038/10/25 12:00:01 INFO Request payload={City: Kharkiv, Country: UA, Street: [CENSORED]}`
}

```

### Map

Both keys and values are recursively parsed, ensuring the output is properly formatted. Rules are the same as for
key/value types.

In the example below, we're using a map with a struct as a key and a slice of strings as a value. All the fields of the
struct are masked by default, except those that have the `censor:"display"` tag. So, in this example, only the `ID` and
`Balance` fields will be displayed.

Regarding the slice of strings, all the values will be masked because of the exclude pattern that we've added.

```go
package main

import (
  "log/slog"

  "github.com/vpakhuchyi/censor"
)

type user struct {
  ID       string `censor:"display"`
  Balance  int    `censor:"display"`
  FullName string
}

func main() {
  cfg := censor.Config{
    Formatter: censor.FormatterConfig{
      MaskValue: censor.DefaultMaskValue,
      // This is a regular expression that matches email addresses.
      ExcludePatterns: []string{`(?P<email>[\w.%+\-]+@[A-Za-z0-9.\-]+\.[A-Za-z]{2,4})`},
    },
  }

  p := censor.NewWithConfig(cfg)

  v := map[user][]string{
    user{ID: "aaf42135", Balance: 1101, FullName: "Viktor Pakhuchyi"}: []string{"example1.email@ggmail.com", "example2.email@ggmail.com"},
    user{ID: "456mlkn", Balance: 999999, FullName: "Bruce Wayne"}:     []string{"wayne.day.email@ggmail.com", "wayne.night.email@ggmail.com"},
  }

  slog.Info("Request", "payload", p.Format(v))
  // Here is what we'll see in the log:
  //Output: `2038/10/25 12:00:01 INFO Request payload=map[{ID: 456mlkn, Balance: 999999, FullName: [CENSORED]}: [[CENSORED], [CENSORED]], {ID: aaf42135, Balance: 1101, FullName: [CENSORED]}: [[CENSORED], [CENSORED]]]`
}

```

There is also an option to display the map type in the output. To do this, you need to enable the `DisplayMapType`
option in the configuration struct or file. In such a case, the output will look like this:

```go
package main

import (
  "log/slog"

  "github.com/vpakhuchyi/censor"
)

func main() {
  cfg := censor.Config{
    Formatter: censor.FormatterConfig{
      MaskValue:      censor.DefaultMaskValue,
      DisplayMapType: true,
    },
  }

  p := censor.NewWithConfig(cfg)

  v := map[string]int{
    "one": 1,
    "two": 2,
  }

  slog.Info("Request", "payload", p.Format(v))
  // Here is what we'll see in the log:
  //Output: `2038/10/25 12:00:01 INFO Request payload=map[string]int[one: 1, two: 2]`
}

```

### Slice/Array

For formatting of `slices` and `arrays` the same rules are applied.
Each element is formatted using its underlying type rules.

```go
package main

import (
  "log/slog"

  "github.com/vpakhuchyi/censor"
)

func main() {
  v := [2][]int{[]int{1, 2, 3}, []int{4, 5, 6}}

  slog.Info("Request", "payload", censor.Format(v))
  // Here is what we'll see in the log:
  //Output: `2038/10/25 12:00:01 INFO Request payload=[[1, 2, 3], [4, 5, 6]]`
}

```

### Pointer

Then the underlying value is formatted using its type rules.
In case of nil pointer, the output will be just `nil`.

```go
package main

import (
  "log/slog"

  "github.com/vpakhuchyi/censor"
)

func main() {
  s := []int{1, 2, 3}
  v := &s

  slog.Info("Request", "payload", censor.Format(v))
  // Here is what we'll see in the log:
  //Output: `2038/10/25 12:00:01 INFO Request payload=[1, 2, 3]`

  cfg := censor.Config{
    Formatter: censor.FormatterConfig{
      MaskValue: censor.DefaultMaskValue,
      // If you want to display the pointer symbol before the pointed value in the output,
      // you can use the `DisplayPointerSymbol` configuration option. 
      // In this case, the output will look like this:
      DisplayPointerSymbol: true,
    },
  }

  p := censor.NewWithConfig(cfg)

  slog.Info("Request", "payload", p.Format(v))
  // Here is what we'll see in the log:
  //Output: `2038/10/25 12:00:01 INFO Request payload=&[1, 2, 3]`
}

```

### String

By default, string values are handled with no additional formatting.

```go
package main

import (
  "log/slog"

  "github.com/vpakhuchyi/censor"
)

// This is a regular expression that matches email addresses.
const emailPattern = `(?P<email>[\w.%+\-]+@[A-Za-z0-9.\-]+\.[A-Za-z]{2, 4})`

func main() {
  v := "example1.email@ggmail.com"

  slog.Info("Request", "payload", censor.Format(v))
  // Here is what we'll see in the log:
  //Output: `2038/10/25 12:00:01 INFO Request payload = example1.email@ggmail.com`

  // If you want to mask specific values, you can use the `ExcludePatterns` configuration option
  // to add exclude patterns. All string values sections that match the specified patterns will be masked.
  cfg := censor.Config{
    Formatter: censor.FormatterConfig{
      MaskValue:       censor.DefaultMaskValue,
      ExcludePatterns: []string{emailPattern},
    },
  }

  p, err := censor.NewWithOpts(censor.WithConfig(&cfg))

  slog.Info("Request", "payload", p.Format(v))
  // Here is what we'll see in the log:
  //Output: `2038/10/25 12:00:01 INFO Request payload = "[CENSORED]"`
}

```

### Float64/Float32

Due to the way Go runtime works, `float64` and `float32` types values are not always kept as the original values.
For example, the value `99.123456789123456789` will be stored as `99.12345678912345` for `float64` type and as
`99.12346` for `float32` type. That's before any formatting is applied.

Talking about formatting, there are a few strategies that we could use to display float values.
More details can be found here: https://github.com/golang/go/blob/master/src/fmt/doc.go#L38.

To have a more deterministic output, Censor is using the https://github.com/shopspring/decimal package.
In such a way, we can display float values that are stored in the runtime with no changes in most cases.

```go
package main

import (
  "log/slog"

  "github.com/vpakhuchyi/censor"
)

func main() {
  var v float32 = 99.123456789123456789

  slog.Info("Request", "payload", censor.Format(v))
  // Here is what we'll see in the log:
  //Output: `2038/10/25 12:00:01 INFO Request payload=99.12346`

  var v2 float64 = 9.123456789123456789

  slog.Info("Request", "payload", censor.Format(v2))
  // Here is what we'll see in the log:
  //Output: `2038/10/25 12:00:01 INFO Request payload=99.12345678912345`
}

```

### Interface

Will be formatted using the same rules as its underlying type.

Note: the main goal is to display the value in the output, not the runtime or compile-time type. That's why in the
example below, we will see the runtime values and based on the output it's not possible to determine the specific type
of it.

```go
package main

import (
  "log/slog"

  "github.com/vpakhuchyi/censor"
)

func main() {
  // Variable v is an interface that contains a slice of interfaces with different types values.
  var v interface{} = []interface{}{1, 'v', "kanoe", true}

  slog.Info("Request", "payload", censor.Format(v))
  // Here is what we'll see in the log:
  //Output: `2038/10/25 12:00:01 INFO Request payload=[1, 118, kanoe, true]`
}

```

### Rune

Please pay attention that rune type is formatted as int32 type because, under the hood, it's an alias for int32.
That's why all the runes will be formatted to be displayed as int32 Unicode code points.

```go
package main

import (
  "log/slog"

  "github.com/vpakhuchyi/censor"
)

func main() {
  var v rune = 'A'

  slog.Info("Request", "payload", censor.Format(v))
  // Here is what we'll see in the log:
  //Output: `2038/10/25 12:00:01 INFO Request payload=65`
}

```

### encoding.TextMarshaler

If a type implements the `encoding.TextMarshaler` interface, the `MarshalText()` method will be used to get the string

```go
package main

import (
  "fmt"
  "log/slog"

  "github.com/vpakhuchyi/censor"
)

type address struct {
  City    string `censor:"display"`
  Country string `censor:"display"`
  Street  string
}

func (a address) MarshalText() (text []byte, err error) {
  return []byte(fmt.Sprintf("%s, %s", a.City, a.Country)), nil
}

func main() {
  a := address{
    City:    "London",
    Country: "UK",
  }

  slog.Info("Request", "payload", censor.Format(a))
  // Here is what we'll see in the log:
  //Output: `2024/01/10 21:20:35 INFO Request payload="London, UK"`
}


```

### Other types

There is no specific formatting for the following types:

- Int/Int8/Int16/Int32/Int64
- Uint/Uint8/Byte/Uint16/Uint32/Uint64/Uintptr
- Bool

Rules of fmt package are applied to them.

### Unsupported types

There are a few types that are not supported by Censor. In case of such types, the output will contain the string
of the following format: `[Unsupported type: <type>]` - where `<type>` is the type of the value.

```go
package main

import (
  "log/slog"

  "github.com/vpakhuchyi/censor"
)

type s struct {
  Chan          chan int       `censor:"display"`
  Func          func()         `censor:"display"`
  UnsafePointer unsafe.Pointer `censor:"display"`
  Complex64     complex64      `censor:"display"`
  Complex128    complex128     `censor:"display"`
}

func main() {
  v := s{
    Chan:          make(chan int),
    Func:          func() {},
    UnsafePointer: unsafe.Pointer(uintptr(1)),
	Complex64:     complex(1.11231, 2.034),
	Complex128:    complex(11123.123, 5.5468098889),
  }

  slog.Info("Request", "payload", censor.Format(v))
  // Here is what we'll see in the log:
  //Output: `2038/10/25 12:00:01 INFO Request payload={Chan: [Unsupported type: chan], Func: [Unsupported type: func], UnsafePointer: [Unsupported type: unsafe.Pointer], Complex64: [Unsupported type: complex64], Complex128: [Unsupported type: complex128]}`
}

```