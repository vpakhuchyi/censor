package censor

/*
Package censor provides a function for formatting supported values into string representations.

The primary purpose of this package is to accept a supported value (as listed in the provided table) and
format it as a string. The main feature is that when a struct is passed (potentially nested within pointers,
interfaces, slices, or arrays), all its field values are masked by default, except for those fields explicitly
tagged with `censor:"display"`. On top of that, it's possible to specify regexp patterns that will be used to
identify strings that must be masked (including nested string).

This functionality is particularly useful for scenarios such as logging, where the result of `censor.Format()`
can be employed without concerns about exposing sensitive data. By default, this package ensures that sensitive
information within struct fields remains masked, enhancing security in scenarios where data privacy is crucial.

Note: this package includes several functions that may potentially result in panics.
While extensive testing covers most common use cases, if you encounter any bugs or issues,
please feel free to report them in https://github.com/vpakhuchyi/censor as new issues.

Examples: https://github.com/vpakhuchyi/censor#readme

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

Note: If a value of an unsupported type is provided, an empty string is returned.
*/

// Format takes a supported value and format it into a string, with default struct field masking.
// Structs, slices, arrays, pointers, maps, and interfaces are parsed recursively.
// The `censor:"display"` tag can be used to override default struct field masking for specific fields.
//
// Note: The function may utilize reflection, introducing potential performance overhead.
// Avoid in performance-critical scenarios.
func Format(val any) string {
	return globalInstance.Format(val)
}
