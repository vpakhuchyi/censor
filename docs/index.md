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
    - `struct`, `map`, `slice`, `array`, `pointer`, `string`,
    - `float64/float32`, `int/int8/int16/int32/int64/rune`,
    - `uint/uint8/uint16/uint32/uint64/byte`, `bool`,
    - `interface`, `complex64/complex128`.
- [x] Customizable configuration.

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
functionality by incorporating exclude patterns through the `censor.AddExcludePatterns` function.

When exclude patterns are introduced, encapsulated formatters compares all string values against the specified patterns,
masking those that match. This capability proves invaluable when, for instance, you want to conceal all email addresses
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

type request struct {
  UserID   string `censor:"display"`
  Email    string
  FullName string
}

func main() {
  // This is a regular expression that matches email addresses.
  const emailPattern = `(?P<email>[\w.%+\-]+@[A-Za-z0-9.\-]+\.[A-Za-z]{2,4})`

  // We can add exclude patterns to censor to make sure that the values that match the pattern will be masked.
  censor.AddExcludePatterns(emailPattern)

  v := []interface{}{
    "Buzyna",
    request{
      UserID:   "001",
      Email:    "viktor.example.email@ggmail.com", // Masked by default because of the lack of "censor:"display"" tag.
      FullName: "Viktor Pakhuchyi",
    },
    []string{"one", "two", "ivan.example.email@ggmail.com"}, // Masked because of the exclude pattern.
    "mykola.more.example.email@ggmail.com",                  // Masked because of the exclude pattern.
  }

  slog.Info("Request", "payload", censor.Format(v))
  // Here is what we'll see in the log:
  // Output: `2038/10/25 12:00:01 INFO Request payload=[Buzyna, {UserID: 001, Email: [CENSORED], FullName: [CENSORED]}, [one, two, [CENSORED]], [CENSORED]]`
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

All configuration options can be set using the package-level functions as shown below. At the same time, you can create
a new instance of `censor.Processor` and use its methods to configure it. All of these options are available with both
local and global instances.

| Global option                                 | Description                                                          |
|-----------------------------------------------|----------------------------------------------------------------------|
| censor.SetMaskValue(s string)                 | Set custom mask value instead of default `[CENSORED]`.               |
| censor.UseJSONTagName(b bool)                 | Use JSON tag name instead of struct field name.                      |
| censor.DisplayPointerSymbol(b bool)           | Display '&' (pointer symbol) before the pointed value in the output. |
| censor.DisplayStructName(b bool)              | Display struct name in the output.                                   |
| censor.DisplayMapType(b bool)                 | Display map type in the output.                                      |
| censor.AddExcludePatterns(patterns ...string) | Add regexp patterns for matched strings values masking.              |

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
  // Here is what we'll see in the log:
  // Output: `2038/10/25 12:00:01 INFO Request payload=[{ID: 123, Email: [####]}, {ID: 456, Email: [####]}]`
}

```

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

// This is a regular expression that matches email addresses.
const emailPattern = `(?P<email>[\w.%+\-]+@[A-Za-z0-9.\-]+\.[A-Za-z]{2,4})`

func main() {
  censor.AddExcludePatterns(emailPattern)

  v := map[user][]string{
    user{ID: "aaf42135", Balance: 1101, FullName: "Viktor Pakhuchyi"}: []string{"example1.email@ggmail.com", "example2.email@ggmail.com"},
    user{ID: "456mlkn", Balance: 999999, FullName: "Bruce Wayne"}:     []string{"wayne.day.email@ggmail.com", "wayne.night.email@ggmail.com"},
  }

  slog.Info("Request", "payload", censor.Format(v))
  // Here is what we'll see in the log:
  //Output: `2038/10/25 12:00:01 INFO Request payload=map[{ID: 456mlkn, Balance: 999999, FullName: [CENSORED]}: [[CENSORED], [CENSORED]], {ID: aaf42135, Balance: 1101, FullName: [CENSORED]}: [[CENSORED], [CENSORED]]]`
}

```

There is also an option to display the map type in the output. To do this, you need to set
the `censor.DisplayMapType(b bool)` option. In this case, the output will look like this:

```go
package main

import (
  "log/slog"

  "github.com/vpakhuchyi/censor"
)

func main() {
  censor.DisplayMapType(true)

  v := map[string]int{
    "one": 1,
    "two": 2,
  }

  slog.Info("Request", "payload", censor.Format(v))
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

  // If you want to display the pointer symbol before the pointed value in the output,
  // you can use the `censor.DisplayPointerSymbol(b bool)` option. In this case, the output will look like this:

  censor.DisplayPointerSymbol(true)
  slog.Info("Request", "payload", censor.Format(v))
  // Here is what we'll see in the log:
  //Output: `2038/10/25 12:00:01 INFO Request payload=&[1, 2, 3]`
}

```

### String

By default, string values are handled with no additional formatting.

However, if you want to mask specific values, you can use the `censor.AddExcludePatterns` or `censor.SetExcludePatterns`
functions to add exclude patterns. All string values that match the specified patterns will be masked.

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

  // Add the pattern to the list of patterns to be excluded from the log.
  censor.AddExcludePatterns(emailPattern)

  slog.Info("Request", "payload", censor.Format(v))
  // Here is what we'll see in the log:
  //Output: `2038/10/25 12:00:01 INFO Request payload = "[CENSORED]"`
}

```

### Float64/Float32

Floating-point types are formatted to include up to 15 (float64) and 7 (float32) significant figures respectively.

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
  //Output: `2038/10/25 12:00:01 INFO Request payload=9.12345678912346`
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

### Complex64/Complex128

Both parts are formatted to include up to 15 (complex128) and 7 (complex64) precision digits respectively that is
similar to float types.

```go
package main

import (
  "log/slog"

  "github.com/vpakhuchyi/censor"
)

func main() {
  var v complex64 = complex(99.123456789, 22.2222)

  slog.Info("Request", "payload", censor.Format(v))
  // Here is what we'll see in the log:
  //Output: `2038/10/25 12:00:01 INFO Request payload=(99.12346+22.2222i)`

  var v2 complex128 = complex(9.123456789123456, 1.050595950)

  slog.Info("Request", "payload", censor.Format(v2))
  // Here is what we'll see in the log:
  //Output: `2038/10/25 12:00:01 INFO Request payload=(9.12345678912346+1.05059595i)`
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

### Other types

There is no specific formatting for the following types:

- Int/Int8/Int16/Int32/Int64
- Uint/Uint8/Byte/Uint16/Uint32/Uint64
- Bool

Rules of fmt package are applied to them.