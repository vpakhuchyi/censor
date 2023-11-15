package censor

import (
	"testing"

	"github.com/stretchr/testify/require"
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
	t.Run("Display struct name", func(t *testing.T) {
		p := New()
		p.DisplayStructName(true)

		type testStruct struct {
			Name string `censor:"display"`
			Age  int
		}

		exp := `censor.testStruct{Name: John, Age: [******]}`
		got := p.Format(testStruct{Name: "John", Age: 30})
		require.Equal(t, exp, got)
	})

	t.Run("Hide struct name", func(t *testing.T) {
		type testStruct struct {
			Name string `censor:"display"`
			Age  int
		}

		exp := `{Name: John, Age: [******]}`
		got := New().Format(testStruct{Name: "John", Age: 30})
		require.Equal(t, exp, got)
	})

	t.Run("Use JSON tag name", func(t *testing.T) {
		p := New()
		p.UseJSONTagName(true)

		type testStruct struct {
			Name  string `json:"name" censor:"display"`
			Age   int    `json:"age"`
			Email string
		}

		exp := `{name: John, age: [******]}`
		got := p.Format(testStruct{Name: "John", Age: 30})
		require.Equal(t, exp, got)
	})

	t.Run("Custom mask value", func(t *testing.T) {
		p := New()
		p.SetMaskValue(`[REDACTED]`)

		type testStruct struct {
			Name string `censor:"display"`
			Age  int
		}

		exp := `{Name: John, Age: [REDACTED]}`
		got := p.Format(testStruct{Name: "John", Age: 30})
		require.Equal(t, exp, got)
	})

	t.Run("Display map type", func(t *testing.T) {
		p := New()
		p.DisplayMapType(true)

		type testStruct struct {
			M map[string]map[string]int `censor:"display"`
		}

		exp := `{M: map[string]map[string]int[key1: map[string]int[key1: 1, key2: 2]]}`
		got := p.Format(testStruct{M: map[string]map[string]int{"key1": {"key1": 1, "key2": 2}}})
		require.Equal(t, exp, got)
	})
}

func Test_GlobalInstanceConfiguration(t *testing.T) {
	t.Run("Display struct name", func(t *testing.T) {
		t.Cleanup(func() { SetGlobalInstance(New()) })

		DisplayStructName(true)

		type testStruct struct {
			Name string `censor:"display"`
			Age  int
		}

		exp := `censor.testStruct{Name: John, Age: [******]}`
		got := Format(testStruct{Name: "John", Age: 30})
		require.Equal(t, exp, got)
	})

	t.Run("Hide struct name", func(t *testing.T) {
		t.Cleanup(func() { SetGlobalInstance(New()) })

		type testStruct struct {
			Name string `censor:"display"`
			Age  int
		}

		exp := `{Name: John, Age: [******]}`
		got := Format(testStruct{Name: "John", Age: 30})
		require.Equal(t, exp, got)
	})

	t.Run("Use JSON tag name", func(t *testing.T) {
		t.Cleanup(func() { SetGlobalInstance(New()) })

		UseJSONTagName(true)

		type testStruct struct {
			Name string `json:"name" censor:"display"`
			Age  int    `json:"age"`
		}

		exp := `{name: John, age: [******]}`
		got := Format(testStruct{Name: "John", Age: 30})
		require.Equal(t, exp, got)
	})

	t.Run("Custom mask value", func(t *testing.T) {
		t.Cleanup(func() { SetGlobalInstance(New()) })

		SetMaskValue(`[REDACTED]`)

		type testStruct struct {
			Name string `censor:"display"`
			Age  int
		}

		exp := `{Name: John, Age: [REDACTED]}`
		got := Format(testStruct{Name: "John", Age: 30})
		require.Equal(t, exp, got)
	})

	t.Run("Display map type", func(t *testing.T) {
		t.Cleanup(func() { SetGlobalInstance(New()) })

		DisplayMapType(true)

		type testStruct struct {
			M map[string]map[string]int `censor:"display"`
		}

		exp := `{M: map[string]map[string]int[key1: map[string]int[key1: 1, key2: 2]]}`
		got := Format(testStruct{M: map[string]map[string]int{"key1": {"key1": 1, "key2": 2}}})
		require.Equal(t, exp, got)
	})
}

func Test_GetGlobalInstance(t *testing.T) {
	require.EqualValues(t, globalInstance, GetGlobalInstance())
}

func Test_SetGlobalInstance(t *testing.T) {
	t.Cleanup(func() { SetGlobalInstance(New()) })

	p := New()
	p.SetMaskValue("[REDACTED]")

	SetGlobalInstance(p)

	require.EqualValues(t, globalInstance, p)
}
