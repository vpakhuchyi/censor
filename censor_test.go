package censor

import (
	"math"
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
	t.Run("display_struct_name", func(t *testing.T) {
		p := New()
		p.DisplayStructName(true)

		type testStruct struct {
			Name string `censor:"display"`
			Age  int
		}

		exp := `censor.testStruct{Name: John, Age: [CENSORED]}`
		got := p.Format(testStruct{Name: "John", Age: 30})
		require.Equal(t, exp, got)
	})

	t.Run("hide_struct_name", func(t *testing.T) {
		type testStruct struct {
			Name string `censor:"display"`
			Age  int
		}

		exp := `{Name: John, Age: [CENSORED]}`
		got := New().Format(testStruct{Name: "John", Age: 30})
		require.Equal(t, exp, got)
	})

	t.Run("use_json_tag_name", func(t *testing.T) {
		p := New()
		p.UseJSONTagName(true)

		type testStruct struct {
			Name  string `json:"name" censor:"display"`
			Age   int    `json:"age"`
			Email string
		}

		exp := `{name: John, age: [CENSORED], Email: [CENSORED]}`
		got := p.Format(testStruct{Name: "John", Age: 30})
		require.Equal(t, exp, got)
	})

	t.Run("custom_mask_value", func(t *testing.T) {
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

	t.Run("display_map_type", func(t *testing.T) {
		p := New()
		p.DisplayMapType(true)

		type testStruct struct {
			M map[string]map[string]int `censor:"display"`
		}

		exp := `{M: map[string]map[string]int[key1: map[string]int[key1: 1, key2: 2]]}`
		got := p.Format(testStruct{M: map[string]map[string]int{"key1": {"key1": 1, "key2": 2}}})
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

	t.Run("display_pointer_symbol", func(t *testing.T) {
		p := New()
		p.DisplayPointerSymbol(true)

		type testStruct struct {
			Names *[]string `censor:"display"`
		}

		exp := `{Names: &[John, Nazar]}`
		got := p.Format(testStruct{Names: &[]string{"John", "Nazar"}})
		require.Equal(t, exp, got)
	})

	t.Run("custom_float32_max_sig_figs", func(t *testing.T) {
		p := New()
		p.SetFloat32MaxSignificantFigures(3)

		type testStruct struct {
			Weight float32 `censor:"display"`
		}

		exp := `{Weight: 58.9}`
		got := p.Format(testStruct{Weight: 58.9350})
		require.Equal(t, exp, got)
	})

	t.Run("custom_float64_max_sig_figs", func(t *testing.T) {
		p := New()
		p.SetFloat64MaxSignificantFigures(10)

		type testStruct struct {
			Pi float64 `censor:"display"`
		}

		exp := `{Pi: 3.141592654}`
		got := p.Format(testStruct{Pi: math.Pi})
		require.Equal(t, exp, got)
	})
}

func Test_GlobalInstanceConfiguration(t *testing.T) {
	t.Run("display_struct_name", func(t *testing.T) {
		t.Cleanup(func() { SetGlobalInstance(New()) })

		DisplayStructName(true)

		type testStruct struct {
			Name string `censor:"display"`
			Age  int
		}

		exp := `censor.testStruct{Name: John, Age: [CENSORED]}`
		got := Format(testStruct{Name: "John", Age: 30})
		require.Equal(t, exp, got)
	})

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

	t.Run("use_json_tag_name", func(t *testing.T) {
		t.Cleanup(func() { SetGlobalInstance(New()) })

		UseJSONTagName(true)

		type testStruct struct {
			Name string `json:"name" censor:"display"`
			Age  int    `json:"age"`
		}

		exp := `{name: John, age: [CENSORED]}`
		got := Format(testStruct{Name: "John", Age: 30})
		require.Equal(t, exp, got)
	})

	t.Run("custom_mask_value", func(t *testing.T) {
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

	t.Run("display_map_type", func(t *testing.T) {
		t.Cleanup(func() { SetGlobalInstance(New()) })

		DisplayMapType(true)

		type testStruct struct {
			M map[string]map[string]int `censor:"display"`
		}

		exp := `{M: map[string]map[string]int[key1: map[string]int[key1: 1, key2: 2]]}`
		got := Format(testStruct{M: map[string]map[string]int{"key1": {"key1": 1, "key2": 2}}})
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

	t.Run("display_pointer_symbol", func(t *testing.T) {
		t.Cleanup(func() { SetGlobalInstance(New()) })

		DisplayPointerSymbol(true)

		type testStruct struct {
			Names *[]string `censor:"display"`
		}

		exp := `{Names: &[John, Nazar]}`
		got := Format(testStruct{Names: &[]string{"John", "Nazar"}})
		require.Equal(t, exp, got)
	})

	t.Run("custom_float32_max_sig_figs", func(t *testing.T) {
		t.Cleanup(func() { SetGlobalInstance(New()) })

		SetFloat32MaxSignificantFigures(3)

		type testStruct struct {
			Weight float32 `censor:"display"`
		}

		exp := `{Weight: 58.9}`
		got := Format(testStruct{Weight: 58.9350})
		require.Equal(t, exp, got)
	})

	t.Run("custom_float64_max_sig_figs", func(t *testing.T) {
		t.Cleanup(func() { SetGlobalInstance(New()) })

		SetFloat64MaxSignificantFigures(10)

		type testStruct struct {
			Pi float64 `censor:"display"`
		}

		exp := `{Pi: 3.141592654}`
		got := Format(testStruct{Pi: math.Pi})
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

func TestGlobalExcludePatterns(t *testing.T) {
	t.Cleanup(func() { SetGlobalInstance(New()) })

	type testStruct struct {
		Name  string `censor:"display"`
		Email string `censor:"display"`
	}

	v := []testStruct{
		{Name: "John", Email: "test@exxample.com"},
		{Name: "John2", Email: "secondtest@exxample.com"},
		{Name: "John Password", Email: "thirdtest@exxample.com"},
	}
	exp := `[{Name: John, Email: test@exxample.com}, {Name: John2, Email: secondtest@exxample.com}, {Name: John Password, Email: thirdtest@exxample.com}]`
	require.Equal(t, exp, Format(v))

	AddExcludePatterns(`\b[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Z|a-z]{2,}\b`)
	exp = `[{Name: John, Email: [CENSORED]}, {Name: John2, Email: [CENSORED]}, {Name: John Password, Email: [CENSORED]}]`
	require.Equal(t, exp, Format(v))

	AddExcludePatterns(`(?i)password`)
	exp = `[{Name: John, Email: [CENSORED]}, {Name: John2, Email: [CENSORED]}, {Name: [CENSORED], Email: [CENSORED]}]`
	require.Equal(t, exp, Format(v))
}

func TestInstanceExcludePatterns(t *testing.T) {
	type testStruct struct {
		Name  string `censor:"display"`
		Email string `censor:"display"`
	}
	p := New()
	v := []testStruct{
		{Name: "John", Email: "test@exxample.com"},
		{Name: "John2", Email: "secondtest@exxample.com"},
		{Name: "John Password", Email: "thirdtest@exxample.com"},
	}
	exp := `[{Name: John, Email: test@exxample.com}, {Name: John2, Email: secondtest@exxample.com}, {Name: John Password, Email: thirdtest@exxample.com}]`
	require.Equal(t, exp, p.Format(v))

	p.AddExcludePatterns(`\b[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Z|a-z]{2,}\b`)
	exp = `[{Name: John, Email: [CENSORED]}, {Name: John2, Email: [CENSORED]}, {Name: John Password, Email: [CENSORED]}]`
	require.Equal(t, exp, p.Format(v))

	p.AddExcludePatterns(`(?i)password`)
	exp = `[{Name: John, Email: [CENSORED]}, {Name: John2, Email: [CENSORED]}, {Name: [CENSORED], Email: [CENSORED]}]`
	require.Equal(t, exp, p.Format(v))
}
