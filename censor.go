package censor

/*
Package censor provides a function for formatting supported values into string representations.

The primary purpose of this package is to accept a supported value (as listed in the provided table) and
format it as a string. The main feature is that when a struct is passed (potentially nested within pointers,
interfaces, slices, or arrays), all its field values are masked by default, except for those fields explicitly
tagged with `censor:"display"`.

This functionality is particularly useful for scenarios such as logging, where the result of `censor.Format()`
can be employed without concerns about exposing sensitive data. By default, this package ensures that sensitive
information within struct fields remains masked, enhancing security in scenarios where data privacy is crucial.

Note: this package uses reflection, which can be slow. It is not recommended to use this package
in performance-critical scenarios.

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
|------------------------|-----------------------------------------------------------------------|
| Slice/Array/Map        | Formatted using the same rules as the value it contains.              |
| Pointer/Interface      |                                                                       |
|------------------------|-----------------------------------------------------------------------|
| Float64                | Formatted value has up to 15 significant figures.                     |
|------------------------|-----------------------------------------------------------------------|
| Float32                | Formatted value has up to 7 significant figures.                      |
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
	return globalInstance.sanitise(val)
}

// SetMaskValue sets a value that will be used to mask struct fields.
// It applies this change to the global instance of Processor.
// The default value is stored in the formatter.DefaultMaskValue constant.
func SetMaskValue(maskValue string) {
	globalInstance.formatter.MaskValue = maskValue
}

// UseJSONTagName sets whether to use the `json` tag to get the name of the struct field.
// It applies this change to the global instance of Processor.
// If no `json` tag is present, the name of struct will be an empty string.
// By default, this option is disabled.
func UseJSONTagName(v bool) {
	globalInstance.parser.UseJSONTagName = v
}

// DisplayStructName sets whether to display the name of the struct.
// It applies this change to the global instance of Processor.
// By default, this option is disabled.
func DisplayStructName(v bool) {
	globalInstance.formatter.DisplayStructName = v
}

// DisplayMapType sets whether to display map type in the output.
// It applies this change to the global instance of Processor.
// By default, this option is disabled.
func DisplayMapType(v bool) {
	globalInstance.formatter.DisplayMapType = v
}
