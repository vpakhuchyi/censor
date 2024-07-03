package censor

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestProcessor_Format(t *testing.T) {
	type Nested struct {
		String  string `censor:"display"`
		Float64 float64
		Array   [2]string `censor:"display"`
	}

	type S struct {
		String1 string `censor:"display"`
		String2 string
		Map     map[string]string `censor:"display"`
		Slice   []string
		Nested
	}

	t.Run("with exclude pattern", func(t *testing.T) {
		// GIVEN a config instance.
		cfg := Config{
			PrintConfigOnInit: true,
			MaskValue:         "[CENSORED]",
			ExcludePatterns:   []string{`\b[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Z|a-z]{2,}\b`},
		}
		p, err := NewWithOpts(WithConfig(&cfg))
		if err != nil {
			t.Fatal(err)
		}

		// WHEN the Format method is called on the processor instance.
		payload := S{
			String1: "Hello",
			String2: "exampe@gmail.com",
			Map:     map[string]string{"key": "value"},
			Slice:   []string{"one", "two"},
			Nested: Nested{
				String:  "Nothing",
				Float64: 3.14,
				Array:   [2]string{"one", "exampe@gmail.com"},
			},
		}
		got := p.Format(payload)

		// THEN the returned string contains the configuration details.
		want := `{"String1":"Hello","String2":"[CENSORED]","Map":{"key":"value"},"Slice":"[CENSORED]","String":"Nothing","Float64":"[CENSORED]","Array":["one","[CENSORED]"]}`
		require.JSONEq(t, want, got)
	})

	t.Run("with config file option", func(t *testing.T) {
		// GIVEN a config instance.
		p, err := NewWithOpts(WithConfigPath("testdata/cfg.yml"))
		if err != nil {
			t.Fatal(err)
		}

		// WHEN the Format method is called on the processor instance.
		payload := S{
			String1: "Hello",
			String2: "exampe@gmail.com",
			Map:     map[string]string{"key": "value"},
			Slice:   []string{"one", "two"},
		}
		got := p.Format(payload)

		// THEN the returned string contains the configuration details.
		want := `{"String1":"Hello", "String2":"[CENSORED]", "Map": {"key":"value"}, "Slice":"[CENSORED]", "Array": ["", ""], "Float64":"[CENSORED]",  "String":"" }`
		require.JSONEq(t, want, got)
	})
}
