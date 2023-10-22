package sanitiser

import (
	"fmt"
	"reflect"

	"github.com/vpakhuchyi/sanitiser/internal/formatters"
	"github.com/vpakhuchyi/sanitiser/internal/models"
	"github.com/vpakhuchyi/sanitiser/internal/parsers"
)

// unsupportedKind is returned when the value is of unsupported kind.
const unsupportedKind = `[unsupported kind of value: %v]`

// Format takes any value and returns a string representation of it.
// It uses reflection to parse the value and then uses formatters to format it.
// Examples can be found here https://github.com/vpakhuchyi/sanitiser#readme
//
// Supported types:
//
// |----------------------------------------------------------------------------------------|
// |         		   	| By default, all fields values will be masked. 				 	|
// | Struct		   		| To override this behaviour, use the `log:"display"` tag. 	 		|
// |				   	| All nested fields must be tagged as well. 					 	|
// |				  	| Struct/Slice/Array/Pointer/Map values will be parsed recursively	|
// |----------------------------------------------------------------------------------------|
// | Slice/Array    	| Struct/Slice/Array/Pointer/Map values will be parsed recursively 	|
// |----------------------------------------------------------------------------------------|
// | Pointer          	| Struct/Slice/Array/Pointer/Map values will be parsed recursively 	|
// |----------------------------------------------------------------------------------------|
// | Map              	| Struct/Slice/Array/Pointer/Map values will be parsed recursively 	|
// |----------------------------------------------------------------------------------------|
// | String           	| Formatted value will be wrapped in double quotes. 			  	|
// |----------------------------------------------------------------------------------------|
// | Float64          	| Formatted value will have up to 15 precision digits. 		  		|
// |----------------------------------------------------------------------------------------|
// | Float32          	| Formatted value will have up to 7 precision digits. 		  		|
// |----------------------------------------------------------------------------------------|
// | Int/Int8/Int16/  	| 															  		|
// | Int32/Int64/Rune 	| 														 	   		|
// | Uint/Uint8/Uint16	| 	Default fmt package formatting is used. 				   		|
// | Uint32/Uint64/   	| 																	|
// | Byte             	| 																	|
// |----------------------------------------------------------------------------------------|
// | Bool             	| Formatted value will be either "true" or "false". 				|
// |----------------------------------------------------------------------------------------|
//
// Unsupported types:
// |------------------------------------------------|
// | Chan         	| Complex64   	| Complex128	|
// | Interface    	| Func     		|       	 	|
// |------------------------------------------------|
// Note: unsupported types will be replaced with "[unsupported type]" string.
func Format(val any) string {
	return sanitise(val)
}

func sanitise(val any) string {
	v := reflect.ValueOf(val)

	var parsed any
	switch v.Kind() {
	case reflect.Struct:
		parsed = parsers.ParseStruct(v)
	case reflect.Slice, reflect.Array:
		parsed = parsers.ParseSlice(v)
	case reflect.Ptr:
		parsed = parsers.ParsePtr(v)
	case reflect.Map:
		parsed = parsers.ParseMap(v)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64, reflect.String, reflect.Bool:
		parsed = models.Value{Value: v.Interface(), Kind: v.Kind()}
	case reflect.Chan, reflect.Func, reflect.UnsafePointer,
		reflect.Complex64, reflect.Complex128:
		/*
			Note: this case covers all unsupported types.
			In such a case, we return a string with a message about the unsupported type.
		*/
		return fmt.Sprintf(unsupportedKind, v.Kind())
	}

	switch v.Kind() {
	case reflect.Struct:
		return formatters.FormatStruct(parsed.(models.Struct))
	case reflect.Slice, reflect.Array:
		return formatters.FormatSlice(parsed.(models.Slice))
	case reflect.Pointer:
		return formatters.FormatPtr(parsed.(models.Ptr))
	case reflect.String:
		return formatters.FormatString(parsed.(models.Value))
	case reflect.Float32, reflect.Float64:
		return formatters.FormatFloat(parsed.(models.Value))
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return formatters.FormatInteger(parsed.(models.Value))
	case reflect.Bool:
		return formatters.FormatBool(parsed.(models.Value))
	case reflect.Map:
		return formatters.FormatMap(parsed.(models.Map))
	default:
		return formatters.FormatDefault(parsed.(models.Value))
	}
}
