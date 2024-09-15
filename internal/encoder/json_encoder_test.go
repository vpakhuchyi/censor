package encoder

import (
	"math"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestJSONEncoder_NewJSONEncoder(t *testing.T) {
	got := NewJSONEncoder(Config{
		MaskValue: "[CENSORED]",
	})
	exp := &JSONEncoder{
		baseEncoder: baseEncoder{
			CensorFieldTag: DefaultCensorFieldTag,
			MaskValue:      "\"[CENSORED]\"",
		},
		structFieldsCache: newFieldsCache(defaultMaxCacheSize),
	}
	require.EqualValues(t, exp, got)
}

func TestJSONEncoder_Encode(t *testing.T) {
	type generic[T string | int] struct {
		GenericField T `censor:"display"`
	}

	type nested struct {
		String    string      `censor:"display"`
		Interface interface{} `censor:"display"`
	}

	type payload struct {
		String           string `censor:"display"`
		StringMasked     string
		StringWithRegexp string            `censor:"display"`
		Int              int               `json:"IntTag" censor:"display"`
		Byte             byte              `censor:"display"`
		Int8             int8              `censor:"display"`
		Int16            int16             `censor:"display"`
		Int32            int32             `censor:"display"`
		Int64            int64             `censor:"display"`
		Uint             uint              `censor:"display"`
		Uint8            uint8             `censor:"display"`
		Uint16           uint16            `censor:"display"`
		Uint32           uint32            `censor:"display"`
		Uint64           uint64            `censor:"display"`
		Rune             rune              `censor:"display"`
		Float32          float32           `censor:"display"`
		Float64          float64           `censor:"display"`
		Bool             bool              `censor:"display"`
		Interface        interface{}       `censor:"display"`
		Struct           nested            `censor:"display"`
		AnonymousStruct  interface{}       `censor:"display"`
		GenericString    generic[string]   `censor:"display"`
		GenericInt       generic[int]      `censor:"display"`
		Slice            []nested          `censor:"display"`
		Array            [2]nested         `censor:"display"`
		Map              map[string]nested `censor:"display"`
		Pointer          *nested           `censor:"display"`
		Time             time.Time         `censor:"display"`
		Func             func()            `censor:"display"`
	}

	// GIVEN.
	p := payload{
		String:           "string\\",
		StringMasked:     "string",
		StringWithRegexp: "bla-bla-example@example.com",
		Int:              1,
		Byte:             'a',
		Int8:             2,
		Int16:            3,
		Int32:            4,
		Int64:            5,
		Uint:             6,
		Uint8:            7,
		Uint16:           8,
		Uint32:           9,
		Uint64:           10,
		Rune:             'y',
		Float32:          1.1,
		Float64:          2.2,
		Bool:             true,
		Interface:        nested{String: "string", Interface: "interface"},
		Struct:           nested{String: "string", Interface: "interface"},
		AnonymousStruct: struct {
			String string `censor:"display"`
		}{String: "string"},
		GenericString: generic[string]{
			GenericField: "string",
		},
		GenericInt: generic[int]{
			GenericField: 123,
		},
		Slice: []nested{
			{String: "string", Interface: "interface1"},
			{String: "string", Interface: "interface2"},
		},
		Array: [2]nested{
			{String: "string", Interface: "interface1"},
			{String: "string", Interface: "interface2"},
		},
		Map: map[string]nested{
			"1": {String: "string", Interface: "interface1"},
		},
		Pointer: &nested{String: "string", Interface: "interface"},
		Time:    time.Date(1861, 2, 19, 0, 0, 0, 0, time.UTC),
		Func:    func() {},
	}

	t.Run("escaping is disabled", func(t *testing.T) {
		e := NewJSONEncoder(Config{
			ExcludePatterns: []string{
				`\b[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Z|a-z]{2,}\b`,
			},
			MaskValue: "[CENSORED]",
		})
		var b strings.Builder
		defer b.Reset()

		// WHEN.
		e.Encode(&b, reflect.ValueOf(p))

		// THEN.
		exp := `{"String": "string\","StringMasked": "[CENSORED]","StringWithRegexp": "[CENSORED]",` +
			`"Int": 1,"Byte": 97,"Int8": 2,"Int16": 3,"Int32": 4,"Int64": 5,"Uint": 6,"Uint8": 7,` +
			`"Uint16": 8,"Uint32": 9,"Uint64": 10,"Rune": 121,"Float32": 1.1,"Float64": 2.2,"Bool": true,` +
			`"Interface": {"String": "string","Interface": "interface"},` +
			`"Struct": {"String": "string","Interface": "interface"},"AnonymousStruct": {"String": "string"},` +
			`"GenericString": {"GenericField": "string"},"GenericInt": {"GenericField": 123},` +
			`"Slice": [{"String": "string","Interface": "interface1"}, {"String": "string","Interface": "interface2"}],` +
			`"Array": [{"String": "string","Interface": "interface1"}, {"String": "string","Interface": "interface2"}],` +
			`"Map": {"1":{"String": "string","Interface": "interface1"}},` +
			`"Pointer": {"String": "string","Interface": "interface"},"Time": "1861-02-19T00:00:00Z","Func": "unsupported type=func"}`
		require.Equal(t, exp, b.String())
	})

	t.Run("escaping is enabled", func(t *testing.T) {
		e := NewJSONEncoder(Config{
			EnableJSONEscaping: true,
			ExcludePatterns: []string{
				`\b[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Z|a-z]{2,}\b`,
			},
			MaskValue: "[CENSORED]",
		})
		var b strings.Builder
		defer b.Reset()

		// WHEN.
		e.Encode(&b, reflect.ValueOf(p))

		// THEN.
		exp := `{"String": "string\\","StringMasked": "[CENSORED]","StringWithRegexp": "[CENSORED]",` +
			`"Int": 1,"Byte": 97,"Int8": 2,"Int16": 3,"Int32": 4,"Int64": 5,"Uint": 6,"Uint8": 7,` +
			`"Uint16": 8,"Uint32": 9,"Uint64": 10,"Rune": 121,"Float32": 1.1,"Float64": 2.2,"Bool": true,` +
			`"Interface": {"String": "string","Interface": "interface"},` +
			`"Struct": {"String": "string","Interface": "interface"},"AnonymousStruct": {"String": "string"},` +
			`"GenericString": {"GenericField": "string"},"GenericInt": {"GenericField": 123},` +
			`"Slice": [{"String": "string","Interface": "interface1"}, {"String": "string","Interface": "interface2"}],` +
			`"Array": [{"String": "string","Interface": "interface1"}, {"String": "string","Interface": "interface2"}],` +
			`"Map": {"1":{"String": "string","Interface": "interface1"}},` +
			`"Pointer": {"String": "string","Interface": "interface"},"Time": "1861-02-19T00:00:00Z","Func": "unsupported type=func"}`
		require.Equal(t, exp, b.String())
	})
}

func TestJSONEncoder_Struct(t *testing.T) {
	t.Run("invalid value kind", func(t *testing.T) {
		require.Panics(t, func() {
			// GIVEN.
			e := NewJSONEncoder(Config{})
			var b strings.Builder
			defer b.Reset()
			v := 26

			// WHEN.
			e.Struct(&b, reflect.ValueOf(v))

			// THEN.
			// Panic.
		})
	})

	t.Run("struct field with CanInterface != true", func(t *testing.T) {
		require.NotPanics(t, func() {
			// GIVEN.
			e := NewJSONEncoder(Config{})
			var b strings.Builder
			defer b.Reset()

			s := struct {
				nested struct {
					String    string      `censor:"display"`
					Interface interface{} `censor:"display"`
				} `censor:"display"`
			}{}

			// WHEN.
			e.Struct(&b, reflect.ValueOf(s))

			// THEN.
			exp := "{}"
			require.Equal(t, exp, b.String())
		})
	})
}

func TestJSONEncoder_Map(t *testing.T) {
	t.Run("invalid value kind", func(t *testing.T) {
		require.Panics(t, func() {
			// GIVEN.
			e := NewJSONEncoder(Config{})
			var b strings.Builder
			defer b.Reset()
			v := 26

			// WHEN.
			e.Map(&b, reflect.ValueOf(v))

			// THEN.
			// Panic.
		})
	})

	t.Run("math.NaN as a key", func(t *testing.T) {
		require.Panics(t, func() {
			// GIVEN.
			e := NewJSONEncoder(Config{})
			var b strings.Builder
			defer b.Reset()

			// There is no way to create a decimal value from NaN.
			v := map[float64]string{
				math.NaN(): "",
			}

			// WHEN.
			e.Map(&b, reflect.ValueOf(v))

			// THEN.
			// Panic.
		})
	})

	t.Run("multiple k-v pairs", func(t *testing.T) {
		require.NotPanics(t, func() {
			// GIVEN.
			e := NewJSONEncoder(Config{})
			var b strings.Builder
			defer b.Reset()

			v := map[string]string{
				"key1": "value1",
				"key2": "value2",
			}

			// WHEN.
			e.Map(&b, reflect.ValueOf(v))

			// THEN.
			got := b.String()
			require.True(t, `{"key1":"value1","key2":"value2"}` == got ||
				`{"key2":"value2","key1":"value1"}` == got)
		})
	})

	t.Run("nil map", func(t *testing.T) {
		require.NotPanics(t, func() {
			// GIVEN.
			e := NewJSONEncoder(Config{})
			var b strings.Builder
			defer b.Reset()

			var v map[string]string

			// WHEN.
			e.Map(&b, reflect.ValueOf(v))

			// THEN.
			exp := "null"
			require.Equal(t, exp, b.String())
		})
	})
}

func TestJSONEncoder_Slice(t *testing.T) {
	t.Run("invalid value kind", func(t *testing.T) {
		require.Panics(t, func() {
			// GIVEN.
			e := NewJSONEncoder(Config{})
			var b strings.Builder
			defer b.Reset()
			v := 26

			// WHEN.
			e.Slice(&b, reflect.ValueOf(v))

			// THEN.
			// Panic.
		})
	})
}

func TestJSONEncoder_Interface(t *testing.T) {
	t.Run("invalid value kind", func(t *testing.T) {
		require.Panics(t, func() {
			// GIVEN.
			e := NewJSONEncoder(Config{})
			var b strings.Builder
			defer b.Reset()
			v := 26

			// WHEN.
			e.Interface(&b, reflect.ValueOf(v))

			// THEN.
			// Panic.
		})
	})

	// TODO: Implement the test.
	t.Run("nil interface", func(t *testing.T) {})
}

func TestJSONEncoder_Ptr(t *testing.T) {
	t.Run("invalid value kind", func(t *testing.T) {
		require.Panics(t, func() {
			// GIVEN.
			e := NewJSONEncoder(Config{})
			var b strings.Builder
			defer b.Reset()
			v := 26

			// WHEN.
			e.Ptr(&b, reflect.ValueOf(v))

			// THEN.
			// Panic.
		})
	})

	t.Run("nil value", func(t *testing.T) {
		require.NotPanics(t, func() {
			// GIVEN.
			e := NewJSONEncoder(Config{})
			var b strings.Builder
			defer b.Reset()
			var v *string

			// WHEN.
			e.Ptr(&b, reflect.ValueOf(v))

			// THEN.
			// Panic.
		})
	})
}

func Test_escapeString(t *testing.T) {
	tests := map[string]struct {
		input string
		exp   string
	}{
		"1": {input: `\`, exp: `\\`},
		"2": {input: `"`, exp: `\"`},
		"3": {input: `\n`, exp: `\\n`},
		"4": {input: `\r`, exp: `\\r`},
		"5": {input: `\t`, exp: `\\t`},
		"6": {input: `\f`, exp: `\\f`},
		"7": {input: `\b`, exp: `\\b`},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got := escapeString(tt.input)
			require.Equal(t, tt.exp, got)
		})
	}
}
