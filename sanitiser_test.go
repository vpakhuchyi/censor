package sanitiser

import (
	"testing"
	"unsafe"

	"github.com/stretchr/testify/require"
)

func Test_FormatStruct(t *testing.T) {
	tests := map[string]struct {
		val any
		exp string
	}{
		"struct_with_primitive_fields": {
			val: structWithPrimitives{
				Int64:   -53645354,
				Int32:   -346456,
				Int16:   -23452,
				Int8:    -101,
				Int:     -456345655,
				Uint64:  53645354,
				Uint32:  346456,
				Uint16:  23452,
				Uint8:   101,
				Uint:    456345655,
				Bool:    true,
				Rune:    'a',
				Byte:    1,
				Float64: 1.12341,
				Float32: 1.389,
				String:  "string",
			},
			exp: `sanitiser.structWithPrimitives{"Int64": -53645354, "Int32": -346456, "Int16": -23452, "Int8": -101, "Int": -456345655, "Uint64": 53645354, "Uint32": 346456, "Uint16": 23452, "Uint8": 101, "Uint": 456345655, "Bool": true, "Rune": 97, "Byte": 1, "Float64": 1.12341, "Float32": 1.389, "String": "string"}`,
		},
		"struct_with_complex_fields": {
			val: structWithComplexFields{
				Slice:       []address{{City: "San Francisco", State: "CA", Street: "451 Main St", Zip: "55501"}, {City: "Denver", State: "DN", Street: "65 Best St", Zip: "55502"}},
				MaskedSlice: []address{{City: "San Francisco", State: "CA", Street: "451 Main St", Zip: "55501"}, {City: "Denver", State: "DN", Street: "65 Best St", Zip: "55502"}},
				Map: map[string]address{
					"address1": {City: "San Francisco", State: "CA", Street: "451 Main St", Zip: "55501"},
					"address2": {City: "Denver", State: "DN", Street: "65 Best St", Zip: "55502"},
				},
				Array:  [2]address{{City: "San Francisco", State: "CA", Street: "451 Main St", Zip: "55501"}, {City: "Denver", State: "DN", Street: "65 Best St", Zip: "55502"}},
				Ptr:    &address{City: "San Francisco", State: "CA", Street: "451 Main St", Zip: "55501"},
				Struct: address{City: "San Francisco", State: "CA", Street: "451 Main St", Zip: "55501"},
			},
			exp: `sanitiser.structWithComplexFields{"Slice": [sanitiser.address{"City": "San Francisco", "State": "CA", "Street": "[******]", "Zip": "[******]"}, sanitiser.address{"City": "Denver", "State": "DN", "Street": "[******]", "Zip": "[******]"}], "MaskedSlice": "[******]", "Map": map[string]sanitiser.address["address1": sanitiser.address{"City": "San Francisco", "State": "CA", "Street": "[******]", "Zip": "[******]"}, "address2": sanitiser.address{"City": "Denver", "State": "DN", "Street": "[******]", "Zip": "[******]"}], "Array": [sanitiser.address{"City": "San Francisco", "State": "CA", "Street": "[******]", "Zip": "[******]"}, sanitiser.address{"City": "Denver", "State": "DN", "Street": "[******]", "Zip": "[******]"}], "Ptr": &sanitiser.address{"City": "San Francisco", "State": "CA", "Street": "[******]", "Zip": "[******]"}, "Struct": sanitiser.address{"City": "San Francisco", "State": "CA", "Street": "[******]", "Zip": "[******]"}}`,
		},
		"struct_with_containers_fields": {
			val: structWithContainersFields{
				StringSlice:  []string{"tag1", "tag2"},
				IntSlice:     []int{1, 2, 3},
				FloatSlice:   []float64{1.1, 2.2, 3.30022},
				BoolSlice:    []bool{true, false},
				StructSlice:  []address{{City: "San Francisco", State: "CA", Street: "451 Main St", Zip: "55501"}, {City: "Denver", State: "DN", Street: "65 Best St", Zip: "55502"}},
				PointerSlice: []*int{new(int), new(int)},
				ArraySlice:   [2]string{"tag1", "tag2"}},
			exp: `sanitiser.structWithContainersFields{"StringSlice": ["tag1", "tag2"], "IntSlice": [1, 2, 3], "FloatSlice": [1.1, 2.2, 3.30022], "BoolSlice": [true, false], "StructSlice": [sanitiser.address{"City": "San Francisco", "State": "CA", "Street": "[******]", "Zip": "[******]"}, sanitiser.address{"City": "Denver", "State": "DN", "Street": "[******]", "Zip": "[******]"}], "PointerSlice": [&0, &0], "ArraySlice": ["tag1", "tag2"]}`,
		},
		"empty_struct": {
			val: struct{}{},
			exp: `{}`,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got := Format(tt.val)
			require.Equal(t, tt.exp, got)
		})
	}
}

func Test_FormatSlice(t *testing.T) {
	tests := map[string]struct {
		val any
		exp string
	}{
		"slice_of_pointers_to_structs": {
			val: []*address{
				{Street: "451 Main St", City: "San Francisco", State: "CA", Zip: "55501"},
				{Street: "65 Best St", City: "Denver", State: "DN", Zip: "55502"},
			},
			exp: `[&sanitiser.address{"City": "San Francisco", "State": "CA", "Street": "[******]", "Zip": "[******]"}, &sanitiser.address{"City": "Denver", "State": "DN", "Street": "[******]", "Zip": "[******]"}]`,
		},
		"slice_of_structs": {
			val: []address{
				{Street: "451 Main St", City: "San Francisco", State: "CA", Zip: "55501"},
				{Street: "65 Best St", City: "Denver", State: "DN", Zip: "55502"},
			},
			exp: `[sanitiser.address{"City": "San Francisco", "State": "CA", "Street": "[******]", "Zip": "[******]"}, sanitiser.address{"City": "Denver", "State": "DN", "Street": "[******]", "Zip": "[******]"}]`,
		},
		"slices_of_strings": {
			val: []string{"tag1", "tag2"},
			exp: `["tag1", "tag2"]`,
		},
		"slice_of_integers": {
			val: []int{1, 2, 3},
			exp: `[1, 2, 3]`,
		},
		"slice_of_floats": {
			val: []float64{1.1, 2.2, 3.30022},
			exp: `[1.1, 2.2, 3.30022]`,
		},
		"slice_of_bools": {
			val: []bool{true, false},
			exp: `[true, false]`,
		},
		"slice_of_pointers_to_structs_with_nil": {
			val: []*address{
				{Street: "451 Main St", City: "San Francisco", State: "CA", Zip: "55501"},
				nil,
			},
			exp: `[&sanitiser.address{"City": "San Francisco", "State": "CA", "Street": "[******]", "Zip": "[******]"}, nil]`,
		},
		"empty_slice": {
			val: []string{},
			exp: `[]`,
		},
		"slice_of_slices": {
			val: [][]string{{"tag1", "tag2"}, {"tag3", "tag4"}},
			exp: `[["tag1", "tag2"], ["tag3", "tag4"]]`,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got := Format(tt.val)
			require.Equal(t, tt.exp, got)
		})
	}
}

func Test_FormatArray(t *testing.T) {
	tests := map[string]struct {
		val any
		exp string
	}{
		"array_of_pointers_to_structs": {
			val: [2]*address{
				{Street: "451 Main St", City: "San Francisco", State: "CA", Zip: "55501"},
				{Street: "65 Best St", City: "Denver", State: "DN", Zip: "55502"},
			},
			exp: `[&sanitiser.address{"City": "San Francisco", "State": "CA", "Street": "[******]", "Zip": "[******]"}, &sanitiser.address{"City": "Denver", "State": "DN", "Street": "[******]", "Zip": "[******]"}]`,
		},
		"array_of_structs": {
			val: [2]address{
				{Street: "451 Main St", City: "San Francisco", State: "CA", Zip: "55501"},
				{Street: "65 Best St", City: "Denver", State: "DN", Zip: "55502"},
			},
			exp: `[sanitiser.address{"City": "San Francisco", "State": "CA", "Street": "[******]", "Zip": "[******]"}, sanitiser.address{"City": "Denver", "State": "DN", "Street": "[******]", "Zip": "[******]"}]`,
		},
		"array_of_strings": {
			val: [2]string{"tag1", "tag2"},
			exp: `["tag1", "tag2"]`,
		},
		"array_of_integers": {
			val: [3]int{1, 2, 3},
			exp: `[1, 2, 3]`,
		},
		"array_of_integers_with_bigger_length": {
			val: [4]int{1, 2, 3},
			exp: `[1, 2, 3, 0]`,
		},
		"array_of_floats": {
			val: [3]float64{1.1, 2.2, 3.30022},
			exp: `[1.1, 2.2, 3.30022]`,
		},
		"array_of_bools": {
			val: [2]bool{true, false},
			exp: `[true, false]`,
		},
		"array_of_pointers_to_structs_with_nil": {
			val: [2]*address{
				{Street: "451 Main St", City: "San Francisco", State: "CA", Zip: "55501"},
				nil,
			},
			exp: `[&sanitiser.address{"City": "San Francisco", "State": "CA", "Street": "[******]", "Zip": "[******]"}, nil]`,
		},
		"empty_slice": {
			val: [0]string{},
			exp: `[]`,
		},
		"array_of_arrays": {
			val: [2][2]string{{"tag1", "tag2"}, {"tag3", "tag4"}},
			exp: `[["tag1", "tag2"], ["tag3", "tag4"]]`,
		},
		"array_of_slices": {
			val: [2][]float32{{1.1, 2.2, 3.30022}, {1.1, 2.2, 3.30022}},
			exp: `[[1.1, 2.2, 3.30022], [1.1, 2.2, 3.30022]]`,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got := Format(tt.val)
			require.Equal(t, tt.exp, got)
		})
	}
}

func Test_FormatPointer(t *testing.T) {
	tests := map[string]struct {
		val any
		exp string
	}{
		"pointer_to_struct": {
			val: &address{Street: "451 Main St", City: "San Francisco", State: "CA", Zip: "55501"},
			exp: `&sanitiser.address{"City": "San Francisco", "State": "CA", "Street": "[******]", "Zip": "[******]"}`,
		},
		"nil_pointer": {
			val: (*address)(nil),
			exp: `nil`,
		},
		"pointer_to_int": {
			val: new(int),
			exp: `&0`,
		},
		"pointer_to_string": {
			val: new(string),
			exp: `&""`,
		},
		"pointer_to_bool": {
			val: new(bool),
			exp: `&false`,
		},
		"pointer_to_float64": {
			val: new(float64),
			exp: `&0`,
		},
		"pointer_to_float32": {
			val: func() *float32 {
				var v float32 = 1.2459
				return &v
			}(),
			exp: `&1.2459`,
		},
		"pointer_to_slice": {
			val: &[]string{"tag1", "tag2"},
			exp: `&["tag1", "tag2"]`,
		},
		"pointer_to_array": {
			val: &[2]bool{true, false},
			exp: `&[true, false]`,
		},
		"pointer_to_pointer_to_struct": {
			val: func() **address {
				a := &address{Street: "451 Main St", City: "San Francisco", State: "CA", Zip: "55501"}
				return &a
			}(),
			exp: `&&sanitiser.address{"City": "San Francisco", "State": "CA", "Street": "[******]", "Zip": "[******]"}`,
		},
		"pointer_to_pointer_to_slice": {
			val: func() **[]string {
				a := &[]string{"tag1", "tag2"}
				return &a
			}(),
			exp: `&&["tag1", "tag2"]`,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got := Format(tt.val)
			require.Equal(t, tt.exp, got)
		})
	}
}

func Test_FormatPrimitives(t *testing.T) {
	tests := map[string]struct {
		val any
		exp string
	}{
		"int":     {val: -33453435, exp: `-33453435`},
		"int64":   {val: -3453435, exp: `-3453435`},
		"int32":   {val: -453435, exp: `-453435`},
		"int16":   {val: -53435, exp: `-53435`},
		"int8":    {val: -101, exp: `-101`},
		"uint":    {val: 33453435, exp: `33453435`},
		"uint64":  {val: 3453435, exp: `3453435`},
		"uint32":  {val: 453435, exp: `453435`},
		"uint16":  {val: 53435, exp: `53435`},
		"uint8":   {val: 101, exp: `101`},
		"rune":    {val: 1234, exp: `1234`},
		"byte":    {val: 89, exp: `89`},
		"bool":    {val: true, exp: `true`},
		"string":  {val: "hello", exp: `"hello"`},
		"float64": {val: 12.235325, exp: `12.235325`},
		"float32": {val: -9.654670, exp: `-9.65467`},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got := globalInstance.sanitise(tt.val)
			require.Equal(t, tt.exp, got)
		})
	}
}

func Test_FormatUnsupported(t *testing.T) {
	tests := map[string]struct {
		val any
		exp string
	}{
		"complex64":  {val: complex64(1.1), exp: `[unsupported kind of value: complex64]`},
		"complex128": {val: complex128(1.1), exp: `[unsupported kind of value: complex128]`},
		"chan":       {val: make(chan int), exp: `[unsupported kind of value: chan]`},
		"func":       {val: func() {}, exp: `[unsupported kind of value: func]`},
		"unsafe.Pointer": {
			val: func() unsafe.Pointer {
				var v int
				return unsafe.Pointer(&v)
			}(),
			exp: `[unsupported kind of value: unsafe.Pointer]`,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got := Format(tt.val)
			require.Equal(t, tt.exp, got)
		})
	}
}

func Test_FormatMap(t *testing.T) {
	tests := map[string]struct {
		val any
		exp string
	}{
		"map_string_string": {
			val: map[string]string{"key1": "value1", "key2": "value2"},
			exp: `map[string]string["key1": "value1", "key2": "value2"]`,
		},
		"map_string_int": {
			val: map[string]int{"key1": 1, "key2": 2},
			exp: `map[string]int["key1": 1, "key2": 2]`,
		},
		"map_string_float64": {
			val: map[string]float64{"key1": -1.12342342, "key2": 2.24567},
			exp: `map[string]float64["key1": -1.12342342, "key2": 2.24567]`,
		},
		"map_string_float32": {
			val: map[string]float32{"key1": -6457.2342, "key2": -13.1234},
			exp: `map[string]float32["key1": -6457.234, "key2": -13.1234]`,
		},
		"map_string_bool": {
			val: map[string]bool{"key1": true, "key2": false},
			exp: `map[string]bool["key1": true, "key2": false]`,
		},
		"map_string_struct": {
			val: map[string]address{"key1": {Street: "451 Main St", City: "San Francisco", State: "CA", Zip: "55501"}, "key2": {Street: "65 Best St", City: "Denver", State: "DN", Zip: "55502"}},
			exp: `map[string]sanitiser.address["key1": sanitiser.address{"City": "San Francisco", "State": "CA", "Street": "[******]", "Zip": "[******]"}, "key2": sanitiser.address{"City": "Denver", "State": "DN", "Street": "[******]", "Zip": "[******]"}]`,
		},
		"map_string_slice": {
			val: map[string][]string{"key1": {"tag1", "tag2"}, "key2": {"tag3", "tag4"}},
			exp: `map[string][]string["key1": ["tag1", "tag2"], "key2": ["tag3", "tag4"]]`,
		},
		"map_string_array": {
			val: map[string][2]string{"key1": {"tag1", "tag2"}, "key2": {"tag3", "tag4"}},
			exp: `map[string][2]string["key1": ["tag1", "tag2"], "key2": ["tag3", "tag4"]]`,
		},
		"map_string_pointer_to_struct": {
			val: map[string]*address{"key1": {Street: "451 Main St", City: "San Francisco", State: "CA", Zip: "55501"}, "key2": {Street: "65 Best St", City: "Denver", State: "DN", Zip: "55502"}},
			exp: `map[string]*sanitiser.address["key1": &sanitiser.address{"City": "San Francisco", "State": "CA", "Street": "[******]", "Zip": "[******]"}, "key2": &sanitiser.address{"City": "Denver", "State": "DN", "Street": "[******]", "Zip": "[******]"}]`,
		},
		"map_string_map_string": {
			val: map[string]map[string]string{"key1": {"key1": "value1", "key2": "value2"}, "key2": {"key3": "value3", "key4": "value4"}},
			exp: `map[string]map[string]string["key1": map[string]string["key1": "value1", "key2": "value2"], "key2": map[string]string["key3": "value3", "key4": "value4"]]`,
		},
		"map_string_slice_of_map_string": {
			val: map[string][]map[string]string{"key1": {{"key1": "value1", "key2": "value2"}, {"key3": "value3", "key4": "value4"}}, "key2": {{"key5": "value5", "key6": "value6"}, {"key7": "value7", "key8": "value8"}}},
			exp: `map[string][]map[string]string["key1": [map[string]string["key1": "value1", "key2": "value2"], map[string]string["key3": "value3", "key4": "value4"]], "key2": [map[string]string["key5": "value5", "key6": "value6"], map[string]string["key7": "value7", "key8": "value8"]]]`,
		},
		"map_float64_string": {
			val: map[float64]string{1.1: "value1", 2.2: "value2"},
			exp: `map[float64]string[1.1: "value1", 2.2: "value2"]`,
		},
		"map_float32_string": {
			val: map[float32]string{1.1: "value1", 2.2: "value2"},
			exp: `map[float32]string[1.1: "value1", 2.2: "value2"]`,
		},
		"map_int_string": {
			val: map[int]string{1: "value1", 2: "value2"},
			exp: `map[int]string[1: "value1", 2: "value2"]`,
		},
		"map_rune_string": {
			val: map[rune]string{1: "value1", 2: "value2"},
			exp: `map[int32]string[1: "value1", 2: "value2"]`,
		},
		"map_array_string": {
			val: map[[2]string]string{{"key1", "key2"}: "value1", {"key3", "key4"}: "value2"},
			exp: `map[[2]string]string[["key1", "key2"]: "value1", ["key3", "key4"]: "value2"]`,
		},
		"map_pointer_to_struct_string": {
			val: map[*address]string{{Street: "451 Main St", City: "San Francisco", State: "CA", Zip: "55501"}: "value1", {Street: "65 Best St", City: "Denver", State: "DN", Zip: "55502"}: "value2"},
			exp: `map[*sanitiser.address]string[&sanitiser.address{"City": "Denver", "State": "DN", "Street": "[******]", "Zip": "[******]"}: "value2", &sanitiser.address{"City": "San Francisco", "State": "CA", "Street": "[******]", "Zip": "[******]"}: "value1"]`,
		},
		"map_struct_string": {
			val: map[address]string{{Street: "451 Main St", City: "San Francisco", State: "CA", Zip: "55501"}: "value1", {Street: "65 Best St", City: "Denver", State: "DN", Zip: "55502"}: "value2"},
			exp: `map[sanitiser.address]string[sanitiser.address{"City": "Denver", "State": "DN", "Street": "[******]", "Zip": "[******]"}: "value2", sanitiser.address{"City": "San Francisco", "State": "CA", "Street": "[******]", "Zip": "[******]"}: "value1"]`,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got := Format(tt.val)
			require.Equal(t, tt.exp, got)
		})
	}
}

func Test_InstanceFormatPrimitives(t *testing.T) {
	tests := map[string]struct {
		val any
		exp string
	}{
		"int":     {val: -33453435, exp: `-33453435`},
		"uint":    {val: 33453435, exp: `33453435`},
		"string":  {val: "hello", exp: `"hello"`},
		"float64": {val: 12.235325, exp: `12.235325`},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got := New().Format(tt.val)
			require.Equal(t, tt.exp, got)
		})
	}
}

func Test_InstanceConfiguration(t *testing.T) {
	t.Run("Hide struct name", func(t *testing.T) {
		p := New()
		p.HideStructName(true)

		type testStruct struct {
			Name string `log:"display"`
			Age  int
		}

		exp := `{"Name": "John", "Age": "[******]"}`
		got := p.Format(testStruct{Name: "John", Age: 30})
		require.Equal(t, exp, got)
	})

	t.Run("Show struct name", func(t *testing.T) {
		type testStruct struct {
			Name string `log:"display"`
			Age  int
		}

		exp := `sanitiser.testStruct{"Name": "John", "Age": "[******]"}`
		got := New().Format(testStruct{Name: "John", Age: 30})
		require.Equal(t, exp, got)
	})

	t.Run("Use JSON tag name", func(t *testing.T) {
		p := New()
		p.UseJSONTagName(true)

		type testStruct struct {
			Name string `json:"name" log:"display"`
			Age  int    `json:"age"`
		}

		exp := `sanitiser.testStruct{"name": "John", "age": "[******]"}`
		got := p.Format(testStruct{Name: "John", Age: 30})
		require.Equal(t, exp, got)
	})

	t.Run("Use custom sanitiser field tag", func(t *testing.T) {
		p := New()
		p.SetFieldTag("custom")

		type testStruct struct {
			Name string `json:"name" custom:"display"`
			Age  int    `json:"age"`
		}

		exp := `sanitiser.testStruct{"Name": "John", "Age": "[******]"}`
		got := p.Format(testStruct{Name: "John", Age: 30})
		require.Equal(t, exp, got)
	})

	t.Run("Custom mask value", func(t *testing.T) {
		p := New()
		p.SetMaskValue(`[REDACTED]`)

		type testStruct struct {
			Name string `log:"display"`
			Age  int
		}

		exp := `sanitiser.testStruct{"Name": "John", "Age": "[REDACTED]"}`
		got := p.Format(testStruct{Name: "John", Age: 30})
		require.Equal(t, exp, got)
	})
}

func Test_GlobalInstanceConfiguration(t *testing.T) {
	t.Run("Hide struct name", func(t *testing.T) {
		t.Cleanup(func() { SetGlobalInstance(New()) })

		HideStructName(true)

		type testStruct struct {
			Name string `log:"display"`
			Age  int
		}

		exp := `{"Name": "John", "Age": "[******]"}`
		got := Format(testStruct{Name: "John", Age: 30})
		require.Equal(t, exp, got)
	})

	t.Run("Show struct name", func(t *testing.T) {
		t.Cleanup(func() { SetGlobalInstance(New()) })

		type testStruct struct {
			Name string `log:"display"`
			Age  int
		}

		exp := `sanitiser.testStruct{"Name": "John", "Age": "[******]"}`
		got := Format(testStruct{Name: "John", Age: 30})
		require.Equal(t, exp, got)
	})

	t.Run("Use JSON tag name", func(t *testing.T) {
		t.Cleanup(func() { SetGlobalInstance(New()) })

		UseJSONTagName(true)

		type testStruct struct {
			Name string `json:"name" log:"display"`
			Age  int    `json:"age"`
		}

		exp := `sanitiser.testStruct{"name": "John", "age": "[******]"}`
		got := Format(testStruct{Name: "John", Age: 30})
		require.Equal(t, exp, got)
	})

	t.Run("Use custom sanitiser field tag", func(t *testing.T) {
		t.Cleanup(func() { SetGlobalInstance(New()) })

		SetFieldTag("custom")

		type testStruct struct {
			Name string `json:"name" custom:"display"`
			Age  int    `json:"age"`
		}

		exp := `sanitiser.testStruct{"Name": "John", "Age": "[******]"}`
		got := Format(testStruct{Name: "John", Age: 30})
		require.Equal(t, exp, got)
	})

	t.Run("Custom mask value", func(t *testing.T) {
		t.Cleanup(func() { SetGlobalInstance(New()) })

		SetMaskValue(`[REDACTED]`)

		type testStruct struct {
			Name string `log:"display"`
			Age  int
		}

		exp := `sanitiser.testStruct{"Name": "John", "Age": "[REDACTED]"}`
		got := Format(testStruct{Name: "John", Age: 30})
		require.Equal(t, exp, got)
	})
}

func Test_GetGlobalInstance(t *testing.T) {
	require.EqualValues(t, globalInstance, GetGlobalInstance())
}

func Test_SetGlobalInstance(t *testing.T) {
	p := New()
	p.SetFieldTag("custom")

	SetGlobalInstance(p)

	require.EqualValues(t, globalInstance, p)
}
