package encoder

import (
	"math"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestTextEncoder_NewTextEncoder(t *testing.T) {
	got := NewTextEncoder(Config{UseJSONTagName: true})
	exp := &TextEncoder{
		baseEncoder: baseEncoder{
			CensorFieldTag: DefaultCensorFieldTag,
			UseJSONTagName: true,
		},
	}
	require.EqualValues(t, exp, got)
}

func TestTextEncoder_Encode(t *testing.T) {
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
		Slice            []nested          `censor:"display"`
		Array            [2]nested         `censor:"display"`
		Map              map[string]nested `censor:"display"`
		Pointer          *nested           `censor:"display"`
		Time             time.Time         `censor:"display"`
		Func             func()            `censor:"display"`
	}

	// GIVEN.
	p := payload{
		String:           "string",
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

	e := NewTextEncoder(Config{
		DisplayMapType:       true,
		DisplayPointerSymbol: true,
		DisplayStructName:    true,
		ExcludePatterns: []string{
			`\b[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Z|a-z]{2,}\b`,
		},
		MaskValue:      "[CENSORED]",
		UseJSONTagName: true,
	})
	var b strings.Builder
	defer b.Reset()

	// WHEN.
	e.Encode(&b, reflect.ValueOf(p))

	// THEN.
	exp := `encoder.payload{` +
		`String: string, StringMasked: [CENSORED], StringWithRegexp: [CENSORED], ` +
		`IntTag: 1, Byte: 97, Int8: 2, Int16: 3, Int32: 4, Int64: 5, Uint: 6, Uint8: 7, Uint16: 8, ` +
		`Uint32: 9, Uint64: 10, Rune: 121, Float32: 1.1, Float64: 2.2, Bool: true, ` +
		`Interface: encoder.nested{String: string, Interface: interface}, ` +
		`Struct: encoder.nested{String: string, Interface: interface}, ` +
		`Slice: [encoder.nested{String: string, Interface: interface1}, encoder.nested{String: string, Interface: interface2}], ` +
		`Array: [encoder.nested{String: string, Interface: interface1}, encoder.nested{String: string, Interface: interface2}], ` +
		`Map: map[string]encoder.nested{1: encoder.nested{String: string, Interface: interface1}}, ` +
		`Pointer: &encoder.nested{String: string, Interface: interface}, ` +
		`Time: 1861-02-19T00:00:00Z, ` +
		`Func: unsupported type=func` +
		`}`
	require.Equal(t, exp, b.String())
}

func TestTextEncoder_Struct(t *testing.T) {
	t.Run("invalid value kind", func(t *testing.T) {
		require.Panics(t, func() {
			// GIVEN.
			e := NewTextEncoder(Config{})
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
			e := NewTextEncoder(Config{})
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

func TestTextEncoder_Map(t *testing.T) {
	t.Run("invalid value kind", func(t *testing.T) {
		require.Panics(t, func() {
			// GIVEN.
			e := NewTextEncoder(Config{})
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
			e := NewTextEncoder(Config{})
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
			e := NewTextEncoder(Config{})
			var b strings.Builder
			defer b.Reset()

			// There is no way to create a decimal value from NaN.
			v := map[string]string{
				"key1": "value1",
				"key2": "value2",
			}

			// WHEN.
			e.Map(&b, reflect.ValueOf(v))

			// THEN.
			got := b.String()
			require.True(t, `map[string]string{key1: value1, key2: value2}` == got ||
				`map[string]string{key2: value2, key1: value1}` == got)
		})
	})
}

func TestTextEncoder_Slice(t *testing.T) {
	t.Run("invalid value kind", func(t *testing.T) {
		require.Panics(t, func() {
			// GIVEN.
			e := NewTextEncoder(Config{})
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

func TestTextEncoder_Interface(t *testing.T) {
	t.Run("invalid value kind", func(t *testing.T) {
		require.Panics(t, func() {
			// GIVEN.
			e := NewTextEncoder(Config{})
			var b strings.Builder
			defer b.Reset()
			v := 26

			// WHEN.
			e.Interface(&b, reflect.ValueOf(v))

			// THEN.
			// Panic.
		})
	})
}

func TestTextEncoder_Ptr(t *testing.T) {
	t.Run("invalid value kind", func(t *testing.T) {
		require.Panics(t, func() {
			// GIVEN.
			e := NewTextEncoder(Config{})
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
			e := NewTextEncoder(Config{})
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
