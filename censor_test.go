package censor

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/vpakhuchyi/censor/config"
)

func Test_InstanceFormatPrimitives(t *testing.T) {
	tests := map[string]struct {
		val any
		exp string
	}{
		"int":     {val: -33453435, exp: `-33453435`},
		"uint":    {val: 33453435, exp: `33453435`},
		"string":  {val: "hello", exp: `hello`},
		"float64": {val: 12.235325, exp: `12.235325`},
		"nil":     {val: nil, exp: `nil`},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got := New().Format(tt.val)
			require.Equal(t, tt.exp, got)
		})
	}
}

func Test_InstanceConfiguration(t *testing.T) {
	t.Run("hide_struct_name", func(t *testing.T) {
		type testStruct struct {
			Name string `censor:"display"`
			Age  int
		}

		exp := `{Name: John, Age: [CENSORED]}`
		got := New().Format(testStruct{Name: "John", Age: 30})
		require.Equal(t, exp, got)
	})

	t.Run("with_provided_configuration", func(t *testing.T) {
		c := config.Config{
			Formatter: config.Formatter{
				MaskValue:         "[redacted]",
				DisplayStructName: true,
				DisplayMapType:    false,
				ExcludePatterns:   nil,
			},
			Parser: config.Parser{
				UseJSONTagName: true,
			},
		}
		p := NewWithConfig(c)

		type testStruct struct {
			Name string `censor:"display"`
			Age  int    `json:"age" censor:"display"`
		}

		exp := `censor.testStruct{Name: John, age: 30}`
		got := p.Format(testStruct{Name: "John", Age: 30})
		require.Equal(t, exp, got)
	})
}

func Test_GlobalInstanceConfiguration(t *testing.T) {
	t.Run("hide_struct_name", func(t *testing.T) {
		t.Cleanup(func() { SetGlobalInstance(New()) })

		type testStruct struct {
			Name string `censor:"display"`
			Age  int
		}

		exp := `{Name: John, Age: [CENSORED]}`
		got := Format(testStruct{Name: "John", Age: 30})
		require.Equal(t, exp, got)
	})

	t.Run("with_provided_configuration", func(t *testing.T) {
		t.Cleanup(func() { SetGlobalInstance(New()) })

		c := config.Config{
			Formatter: config.Formatter{
				MaskValue:         "[redacted]",
				DisplayStructName: true,
				DisplayMapType:    false,
				ExcludePatterns:   nil,
			},
			Parser: config.Parser{
				UseJSONTagName: true,
			},
		}

		type testStruct struct {
			Name string `censor:"display"`
			Age  int    `json:"age" censor:"display"`
		}

		p := NewWithConfig(c)
		SetGlobalInstance(p)

		exp := `censor.testStruct{Name: John, age: 30}`
		got := Format(testStruct{Name: "John", Age: 30})
		require.Equal(t, exp, got)
	})
}

func Test_GetGlobalInstance(t *testing.T) {
	require.EqualValues(t, globalInstance, GetGlobalInstance())
}

func Test_SetGlobalInstance(t *testing.T) {
	t.Cleanup(func() { SetGlobalInstance(New()) })

	p := NewWithConfig(config.Config{
		Formatter: config.Formatter{
			MaskValue: "[censored]",
		},
	})

	SetGlobalInstance(p)

	type testStruct struct {
		Email string
	}

	v := testStruct{Email: "test@exxample.com"}

	require.EqualValues(t, globalInstance.Format(v), p.Format(v))
}

func TestExcludePatterns(t *testing.T) {
	t.Cleanup(func() { SetGlobalInstance(New()) })

	p := NewWithConfig(config.Config{
		Formatter: config.Formatter{
			MaskValue:       "[CENSORED]",
			ExcludePatterns: []string{`\b[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Z|a-z]{2,}\b`},
		},
	})

	SetGlobalInstance(p)

	type testStruct struct {
		Email string `censor:"display"`
	}

	exp := `{Email: [CENSORED]}`
	require.Equal(t, exp, p.Format(testStruct{Email: "test@exxample.com"}))
}
