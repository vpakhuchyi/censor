package encoder

import (
	"bytes"
	"math"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/vpakhuchyi/censor/internal/cache"
)

func TestJSONEncoder_NewJSONEncoder(t *testing.T) {
	got := NewJSONEncoder(Config{
		MaskValue: "[CENSORED]",
	})
	exp := &JSONEncoder{
		baseEncoder: baseEncoder{
			CensorFieldTag:      defaultCensorFieldTag,
			MaskValue:           "[CENSORED]",
			structFieldsCache:   cache.NewTypeCache[[]Field](cache.DefaultMaxCacheSize),
			escapedStringsCache: cache.New[string](cache.DefaultMaxCacheSize),
			regexpCache:         cache.New[string](cache.DefaultMaxCacheSize),
		},
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
		String:           `st"ring`,
		StringMasked:     "string",
		StringWithRegexp: "%bla-bla-example@example.com#",
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

	t.Run("escaping", func(t *testing.T) {
		e := NewJSONEncoder(Config{
			ExcludePatterns: []string{`\b[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Z|a-z]{2,}\b`},
			MaskValue:       "[CENSORED]",
		})
		var b bytes.Buffer
		defer b.Reset()

		// WHEN.
		e.Encode(&b, reflect.ValueOf(p))

		// THEN.
		exp := `{"String": "st\"ring","StringMasked": "[CENSORED]","StringWithRegexp": "%[CENSORED]#",` +
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
			var b bytes.Buffer
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
			var b bytes.Buffer
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
			// Unexported fields are now skipped entirely, so we get an empty struct
			exp := "{}"
			require.Equal(t, exp, b.String())
		})
	})
}

func TestJSONEncoder_Struct_UnexportedFieldsBehavior(t *testing.T) {
	t.Run("unexported fields are skipped like stdlib", func(t *testing.T) {
		// GIVEN: A struct with both exported and unexported fields
		type MixedStruct struct {
			PublicField    string `censor:"display"`
			privateField   string
			AnotherPublic  int `censor:"display"`
			anotherPrivate bool
		}

		e := NewJSONEncoder(Config{})
		var b bytes.Buffer
		defer b.Reset()

		s := MixedStruct{
			PublicField:    "visible",
			privateField:   "hidden",
			AnotherPublic:  42,
			anotherPrivate: true,
		}

		// WHEN
		e.Struct(&b, reflect.ValueOf(s))

		// THEN: Only exported fields appear in output
		exp := `{"PublicField": "visible","AnotherPublic": 42}`
		require.Equal(t, exp, b.String())
	})

	t.Run("all unexported fields results in empty struct", func(t *testing.T) {
		// GIVEN: A struct with only unexported fields
		type AllPrivate struct {
			field1 string
			field2 int
		}

		e := NewJSONEncoder(Config{})
		var b bytes.Buffer
		defer b.Reset()

		s := AllPrivate{field1: "hidden", field2: 123}

		// WHEN
		e.Struct(&b, reflect.ValueOf(s))

		// THEN: Empty struct
		exp := "{}"
		require.Equal(t, exp, b.String())
	})

	t.Run("masked exported fields still appear", func(t *testing.T) {
		// GIVEN: Struct with exported masked field and unexported field
		type SecureConfig struct {
			APIKey     string // No display tag = masked by default
			privateKey string // Unexported = skipped
		}

		e := NewJSONEncoder(Config{MaskValue: "[CENSORED]"})
		var b bytes.Buffer
		defer b.Reset()

		s := SecureConfig{APIKey: "secret123", privateKey: "hidden"}

		// WHEN
		e.Struct(&b, reflect.ValueOf(s))

		// THEN: APIKey is masked but present, privateKey is completely skipped
		exp := `{"APIKey": "[CENSORED]"}`
		require.Equal(t, exp, b.String())
	})
}

func TestJSONEncoder_Map(t *testing.T) {
	t.Run("invalid value kind", func(t *testing.T) {
		require.Panics(t, func() {
			// GIVEN.
			e := NewJSONEncoder(Config{})
			var b bytes.Buffer
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
			var b bytes.Buffer
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
			var b bytes.Buffer
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
			var b bytes.Buffer
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
			var b bytes.Buffer
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
			var b bytes.Buffer
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
			var b bytes.Buffer
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
			var b bytes.Buffer
			defer b.Reset()
			var v *string

			// WHEN.
			e.Ptr(&b, reflect.ValueOf(v))

			// THEN.
			// Panic.
		})
	})
}

func TestJSONEncoder_StringEscaped(t *testing.T) {
	t.Run("no escaping needed", func(t *testing.T) {
		// GIVEN: An encoder and a plain string without special characters.
		e := NewJSONEncoder(Config{})
		var b bytes.Buffer
		defer b.Reset()

		// WHEN: escapeString is called.
		e.escapeString(&b, "hello")

		// THEN: The string should be written as-is.
		require.Equal(t, "hello", b.String())
	})

	t.Run("escape double quote", func(t *testing.T) {
		// GIVEN: An encoder and a string with a double quote.
		e := NewJSONEncoder(Config{})
		var b bytes.Buffer
		defer b.Reset()

		// WHEN: escapeString is called.
		e.escapeString(&b, `"`)

		// THEN: The quote should be escaped.
		require.Equal(t, `\"`, b.String())
	})

	t.Run("escape backslash", func(t *testing.T) {
		// GIVEN: An encoder and a string with a backslash.
		e := NewJSONEncoder(Config{})
		var b bytes.Buffer
		defer b.Reset()

		// WHEN: escapeString is called.
		e.escapeString(&b, `\`)

		// THEN: The backslash should be escaped.
		require.Equal(t, `\\`, b.String())
	})

	t.Run("escape backspace", func(t *testing.T) {
		// GIVEN: An encoder and a string with a backspace character.
		e := NewJSONEncoder(Config{})
		var b bytes.Buffer
		defer b.Reset()

		// WHEN: escapeString is called.
		e.escapeString(&b, "foo\bbar")

		// THEN: The backspace should be escaped.
		require.Equal(t, "foo\\bbar", b.String())
	})

	t.Run("escape form feed", func(t *testing.T) {
		// GIVEN: An encoder and a string with a form feed character.
		e := NewJSONEncoder(Config{})
		var b bytes.Buffer
		defer b.Reset()

		// WHEN: escapeString is called.
		e.escapeString(&b, "foo\fbar")

		// THEN: The form feed should be escaped.
		require.Equal(t, "foo\\fbar", b.String())
	})

	t.Run("escape newline", func(t *testing.T) {
		// GIVEN: An encoder and a string with a newline character.
		e := NewJSONEncoder(Config{})
		var b bytes.Buffer
		defer b.Reset()

		// WHEN: escapeString is called.
		e.escapeString(&b, "foo\nbar")

		// THEN: The newline should be escaped.
		require.Equal(t, "foo\\nbar", b.String())
	})

	t.Run("escape carriage return", func(t *testing.T) {
		// GIVEN: An encoder and a string with a carriage return character.
		e := NewJSONEncoder(Config{})
		var b bytes.Buffer
		defer b.Reset()

		// WHEN: escapeString is called.
		e.escapeString(&b, "foo\rbar")

		// THEN: The carriage return should be escaped.
		require.Equal(t, "foo\\rbar", b.String())
	})

	t.Run("escape tab", func(t *testing.T) {
		// GIVEN: An encoder and a string with a tab character.
		e := NewJSONEncoder(Config{})
		var b bytes.Buffer
		defer b.Reset()

		// WHEN: escapeString is called.
		e.escapeString(&b, "foo\tbar")

		// THEN: The tab should be escaped.
		require.Equal(t, "foo\\tbar", b.String())
	})

	t.Run("escape control char 0x01", func(t *testing.T) {
		// GIVEN: An encoder and a string with a low control character.
		e := NewJSONEncoder(Config{})
		var b bytes.Buffer
		defer b.Reset()

		// WHEN: escapeString is called.
		e.escapeString(&b, string([]byte{0x01}))

		// THEN: The control character should be escaped as unicode.
		require.Equal(t, `\u0001`, b.String())
	})

	t.Run("escape control char 0x7F", func(t *testing.T) {
		// GIVEN: An encoder and a string with a DEL control character.
		e := NewJSONEncoder(Config{})
		var b bytes.Buffer
		defer b.Reset()

		// WHEN: escapeString is called.
		e.escapeString(&b, string([]byte{0x7F}))

		// THEN: The DEL character should be escaped as unicode.
		require.Equal(t, `\u007f`, b.String())
	})

	t.Run("escape invalid utf8 char", func(t *testing.T) {
		// GIVEN: An encoder and a string with an invalid UTF-8 byte.
		e := NewJSONEncoder(Config{})
		var b bytes.Buffer
		defer b.Reset()

		// WHEN: escapeString is called.
		e.escapeString(&b, string([]byte{0xC0}))

		// THEN: The invalid byte should be replaced with the replacement character.
		require.Equal(t, string([]byte{0xEF, 0xBF, 0xBD}), b.String())
	})

	t.Run("escape valid non-ascii", func(t *testing.T) {
		// GIVEN: An encoder and a string with a valid non-ASCII character.
		e := NewJSONEncoder(Config{})
		var b bytes.Buffer
		defer b.Reset()

		// WHEN: escapeString is called.
		e.escapeString(&b, "ñ")

		// THEN: The non-ASCII character should be escaped as unicode.
		require.Equal(t, `\u00f1`, b.String())
	})

	t.Run("escape unicode line separator", func(t *testing.T) {
		// GIVEN: An encoder and a string with a unicode line separator.
		e := NewJSONEncoder(Config{})
		var b bytes.Buffer
		defer b.Reset()

		// WHEN: escapeString is called.
		e.escapeString(&b, "\u2028")

		// THEN: The line separator should be escaped as unicode.
		require.Equal(t, `\u2028`, b.String())
	})

	t.Run("cache escaped string", func(t *testing.T) {
		// GIVEN: An encoder with empty cache.
		e := NewJSONEncoder(Config{})
		var b bytes.Buffer
		defer b.Reset()
		value := "foo\bbar"

		// WHEN: escapeStringWithCache is called.
		e.escapeStringWithCache(&b, value)

		// THEN: The escaped string should be cached and written to buffer.
		cacheValue, ok := e.escapedStringsCache.Get(value)
		require.True(t, ok)
		require.Equal(t, "foo\\bbar", cacheValue)
		require.Equal(t, "foo\\bbar", b.String())
	})

	t.Run("cache non-escaped string", func(t *testing.T) {
		// GIVEN: An encoder with empty cache.
		e := NewJSONEncoder(Config{})
		var b bytes.Buffer
		defer b.Reset()
		value := "hello"

		// WHEN: escapeStringWithCache is called with a string that needs no escaping.
		e.escapeStringWithCache(&b, value)

		// THEN: The original string should be cached and written to buffer.
		cacheValue, ok := e.escapedStringsCache.Get(value)
		require.True(t, ok)
		require.Equal(t, "hello", cacheValue)
		require.Equal(t, "hello", b.String())
	})

	t.Run("retrieve from cache", func(t *testing.T) {
		// GIVEN: An encoder with a pre-populated cache.
		e := NewJSONEncoder(Config{})
		var b bytes.Buffer
		defer b.Reset()
		s, cached := "foo\bbar", "cached_value"
		e.escapedStringsCache.Set(s, cached)

		// WHEN: escapeStringWithCache is called.
		e.escapeStringWithCache(&b, s)

		// THEN: The cached value should be used without re-escaping.
		require.Equal(t, cached, b.String())
	})

	t.Run("simple string", func(t *testing.T) {
		// GIVEN: An encoder without censoring patterns.
		e := NewJSONEncoder(Config{})
		var b bytes.Buffer
		defer b.Reset()

		// WHEN: StringEscaped is called with a simple string.
		e.StringEscaped(&b, "hello")

		// THEN: The string should be wrapped in quotes.
		require.Equal(t, `"hello"`, b.String())
	})

	t.Run("empty string", func(t *testing.T) {
		// GIVEN: An encoder without censoring patterns.
		e := NewJSONEncoder(Config{})
		var b bytes.Buffer
		defer b.Reset()

		// WHEN: StringEscaped is called with an empty string.
		e.StringEscaped(&b, "")

		// THEN: The result should be empty quotes.
		require.Equal(t, `""`, b.String())
	})

	t.Run("string with quotes", func(t *testing.T) {
		// GIVEN: An encoder without censoring patterns.
		e := NewJSONEncoder(Config{})
		var b bytes.Buffer
		defer b.Reset()

		// WHEN: StringEscaped is called with a string containing quotes.
		e.StringEscaped(&b, `say "hello"`)

		// THEN: The quotes should be escaped and the string wrapped in quotes.
		require.Equal(t, `"say \"hello\""`, b.String())
	})

	t.Run("string with newline", func(t *testing.T) {
		// GIVEN: An encoder without censoring patterns.
		e := NewJSONEncoder(Config{})
		var b bytes.Buffer
		defer b.Reset()

		// WHEN: StringEscaped is called with a string containing a newline.
		e.StringEscaped(&b, "line1\nline2")

		// THEN: The newline should be escaped and the string wrapped in quotes.
		require.Equal(t, `"line1\nline2"`, b.String())
	})

	t.Run("unicode string", func(t *testing.T) {
		// GIVEN: An encoder without censoring patterns.
		e := NewJSONEncoder(Config{})
		var b bytes.Buffer
		defer b.Reset()

		// WHEN: StringEscaped is called with a unicode string.
		e.StringEscaped(&b, "café")

		// THEN: The unicode characters should be escaped and the string wrapped in quotes.
		require.Equal(t, `"caf\u00e9"`, b.String())
	})

	t.Run("string with email pattern", func(t *testing.T) {
		// GIVEN: An encoder with email censoring pattern.
		e := NewJSONEncoder(Config{
			ExcludePatterns: []string{`\b[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Z|a-z]{2,}\b`},
			MaskValue:       "[CENSORED]",
		})
		var b bytes.Buffer
		defer b.Reset()

		// WHEN: StringEscaped is called with a string containing an email.
		e.StringEscaped(&b, "contact: user@example.com")

		// THEN: The email should be censored and the string wrapped in quotes.
		require.Equal(t, `"contact: [CENSORED]"`, b.String())
	})

	t.Run("string with email and escaping", func(t *testing.T) {
		// GIVEN: An encoder with email censoring pattern.
		e := NewJSONEncoder(Config{
			ExcludePatterns: []string{`\b[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Z|a-z]{2,}\b`},
			MaskValue:       "[CENSORED]",
		})
		var b bytes.Buffer
		defer b.Reset()

		// WHEN: StringEscaped is called with a string containing both email and quotes.
		e.StringEscaped(&b, "say \"email: user@example.com\"")

		// THEN: The email should be censored, quotes escaped, and string wrapped in quotes.
		require.Equal(t, `"say \"email: [CENSORED]\""`, b.String())
	})
}
