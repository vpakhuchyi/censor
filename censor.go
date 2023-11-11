package censor

/*
	This package provides a function to censor any value into a string representation.
	By default, all struct fields are masked (including any kind of nested structs),
	unless the field has the `display` tag.

	Examples can be found here: https://github.com/vpakhuchyi/censor#readme.

	Supported types:
	|----------------------|-----------------------------------------------------------------------|
	|                      | By default, all fields values will be masked.	                       |
	| Struct               | To override this behaviour, use the `censor:"display"` tag.        |
	|                      | All nested fields must be tagged as well.                             |
	|                      | Struct/Slice/Array/Pointer/Map values will be parsed recursively.     |
	|----------------------|-----------------------------------------------------------------------|
	| Slice/Array          | Struct/Slice/Array/Pointer/Map values will be parsed recursively.     |
	|----------------------|-----------------------------------------------------------------------|
	| Pointer              | Struct/Slice/Array/Pointer/Map values will be parsed recursively.     |
	|----------------------|-----------------------------------------------------------------------|
	| Map                  | Struct/Slice/Array/Pointer/Map values will be parsed recursively.     |
	|----------------------|-----------------------------------------------------------------------|
	| Interface            | Will be formatted using the same rules as for the value it contains.  |
	|----------------------|-----------------------------------------------------------------------|
	| String               | Default fmt package formatting is used.                               |
	|----------------------|-----------------------------------------------------------------------|
	| Float64              | Formatted value will have up to 15 precision digits.                  |
	|----------------------|-----------------------------------------------------------------------|
	| Float32              | Formatted value will have up to 7 precision digits.                   |
	|----------------------|-----------------------------------------------------------------------|
	| Int/Int8/Int16/      |                                                                       |
	| Int32/Int64/Rune     |                                                                       |
	| Uint/Uint8/Uint16    | Default fmt package formatting is used.                               |
	| Uint32/Uint64/Byte   |                                                                       |
	|----------------------|-----------------------------------------------------------------------|
	| Bool                 | Default fmt package formatting is used.                               |
	|----------------------|-----------------------------------------------------------------------|
	| Complex64/Complex128 | Default fmt package formatting is used.                               |
	|----------------------|-----------------------------------------------------------------------|

	Unsupported types:
	|------------|------------|------------|
	| Chan       | Uintptr    | Func       |
	|------------|------------|------------|

	Note: if a value of unsupported type is provided, an empty string will be returned.
*/

// Format takes any value and returns a string representation of it masking struct fields by default.
// It uses the global instance of Processor.
// To override this behaviour, use the `censor:"display"` tag.
// Formatting is done recursively for all nested structs/slices/arrays/pointers/maps/interfaces.
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
