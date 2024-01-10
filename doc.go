/*
 Package Censor provides a versatile and secure solution for formatting sensitive data,
 making it suitable for scenarios such as logging where data privacy is crucial.

 The primary function of this package is to accept a supported value, as listed in the provided table,
 and format it as a string. Notably, when a struct is passed (potentially nested within pointers,
 interfaces, slices, etc.) all its field values are masked by default. Exceptions can be made for specific
 fields by tagging them with `censor:"display"`.

 Note: while extensive testing covers most common use cases, if you encounter any bugs or issues,
 please feel free to report them in https://github.com/vpakhuchyi/censor as new issues.

Documentation: https://vpakhuchyi.github.io/censor

 Main Features:

 - Struct formatting: Supports comprehensive struct formatting with default values masking of all fields,
   even when nested within pointers, interfaces, slices, arrays, etc.

 - String values masking: Allows users to define regexp patterns for masking specific string parts,
   providing fine-grained control over sensitive information.

 - Wide type support: Censor accommodates a variety of types, including:
	- struct, map, slice, array, pointer, interface
	- string, bool
	- float64/float32, complex64/complex128
	- int/int8/int16/int32/int64/rune
	- uint/uint8/uint16/uint32/uint64/byte

 - Customizable configuration: Offers flexibility in configuration through the use of a `.yaml` file or by
   directly passing a `config.Config` struct.

 Promoted use-case:

 Censor is designed to seamlessly integrate with popular loggers (slog is already supported).
 Its primary use case involves enhancing existing loggers with secure data masking functionality.
 The Censor logger handler can be effortlessly combined with other loggers, ensuring that sensitive
 information within struct fields remains masked, contributing to a more secure logging environment.
 Developers can benefit from this functionality without incurring additional development costs.
 Censor acts as a safety net, mitigating the risk of overlooking its usage in specific code paths.

 Example:

	import (
		"log/slog"

		censor "github.com/vpakhuchyi/censor/handlers/slog"
	)

	func main() {
		// Create a Censor logger handler.
		censorHandler := censor.NewJSONHandler()

		// Create a new slog instance with the Censor handler.
		logger := slog.New(censorHandler)

		// Linking the Censor handler with the slog logger allows forgetting about the Censor usage in the code.
		// Just use your logger as usual. Censor will take care of the rest.

		type User struct {
			Name  string `censor:"display"`
			Email string
		}
		u := User{
			Name:  "Ivan Mazepa",
			Email: "example@gmail.com",
		}

		logger.Info("Sensitive data:", "payload", u)

	Output:	{"time":"2024-01-10T17:01:39.809303+01:00","level":"INFO","msg":"Sensitive data:","payload":"{Name: Ivan Mazepa, Email: [CENSORED]}"}
}

 By default, Censor ensures that sensitive information remains masked, providing enhanced security in
 scenarios where data privacy is of utmost importance.

 Supported Types:
 |------------------------|-----------------------------------------------------------------------|
 | Type                   | Description                                                           |
 |------------------------|-----------------------------------------------------------------------|
 | Struct                 | By default, all field values are masked.                              |
 |                        | Use the `censor:"display"` tag to override this behavior.             |
 |                        | Struct/Slice/Array/Pointer/Map/Interface are parsed recursively.      |
 |                        | All unexported fields are ignored.                                    |
 |------------------------|-----------------------------------------------------------------------|
 | Slice/Array/Map        | Formatted using the same rules as the value it contains.              |
 | Pointer/Interface      |                                                                       |
 |------------------------|-----------------------------------------------------------------------|
 | Float32/Float64        | Formatted using https://github.com/shopspring/decimal.                |
 |------------------------|-----------------------------------------------------------------------|
 | Int/Int8/Int16/        | Default fmt package formatting is used.                               |
 | Int32/Int64/Rune       |                                                                       |
 | Uint/Uint8/Uint16      |                                                                       |
 | Uint32/Uint64/Byte     |                                                                       |
 | String/Bool            |                                                                       |
 | Complex64/Complex128   |                                                                       |
 |------------------------|-----------------------------------------------------------------------|

 Unsupported Types:
 |------------|------------|------------|
 | Chan       | Uintptr    | Func       |
 |------------|------------|------------|
 | UnsafePtr  |            |            |
 |------------|------------|------------|

 If a value of an unsupported type is provided, a string value with the following format will be returned:
 "[Unsupported type: <type>]" - where `<type>` is the type of the provided value.
*/

package censor
