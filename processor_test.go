package censor

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/vpakhuchyi/censor/internal/encoder"
)

func TestNewWithConfig(t *testing.T) {
	cfg := Config{
		General: General{
			OutputFormat:      OutputFormatText,
			PrintConfigOnInit: true,
		},
		Encoder: EncoderConfig{
			DisplayMapType:       false,
			DisplayPointerSymbol: false,
			DisplayStructName:    false,
			ExcludePatterns:      nil,
			MaskValue:            "####",
			UseJSONTagName:       false,
		},
	}

	got, err := NewWithOpts(WithConfig(&cfg))
	require.NoError(t, err)

	encConfig := EncoderConfig{
		DisplayMapType:       cfg.Encoder.DisplayMapType,
		DisplayPointerSymbol: cfg.Encoder.DisplayPointerSymbol,
		DisplayStructName:    cfg.Encoder.DisplayStructName,
		ExcludePatterns:      cfg.Encoder.ExcludePatterns,
		MaskValue:            cfg.Encoder.MaskValue,
	}

	exp := &Processor{
		encoder: encoder.NewTextEncoder(encConfig.toEncoderConfig()),
		cfg:     cfg,
	}

	require.Equal(t, exp, got)

	t.Run("invalid_config", func(t *testing.T) {
		cfg := Config{
			General: General{
				OutputFormat: OutputFormatText,
			},
			Encoder: EncoderConfig{
				MaskValue: "",
			},
		}

		p, err := NewWithOpts(WithConfig(&cfg))
		require.Nil(t, p)
		require.EqualError(t, err, "invalid configuration: mask value cannot be empty")
	})
}

func TestNewWithFileConfig(t *testing.T) {
	t.Run("successful", func(t *testing.T) {
		t.Cleanup(func() { SetGlobalInstance(New()) })

		// GIVEN.
		cfg, err := ConfigFromFile("./testdata/cfg.yml")
		require.NoError(t, err)

		encConfig := EncoderConfig{
			MaskValue:            cfg.Encoder.MaskValue,
			DisplayPointerSymbol: cfg.Encoder.DisplayPointerSymbol,
			DisplayStructName:    cfg.Encoder.DisplayStructName,
			DisplayMapType:       cfg.Encoder.DisplayMapType,
			ExcludePatterns:      cfg.Encoder.ExcludePatterns,
			UseJSONTagName:       cfg.Encoder.UseJSONTagName,
		}

		// WHEN.
		p, err := NewWithOpts(WithConfigPath("./testdata/cfg.yml"))

		// THEN.
		want := Processor{
			encoder: encoder.NewTextEncoder(encConfig.toEncoderConfig()),
			cfg:     cfg,
		}
		require.NoError(t, err)
		require.EqualValues(t, want.encoder, p.encoder)
		require.EqualValues(t, want.cfg, p.cfg)
	})

	t.Run("empty_file_path", func(t *testing.T) {
		t.Cleanup(func() { SetGlobalInstance(New()) })

		// WHEN.
		p, err := NewWithOpts(WithConfigPath(""))

		// THEN.
		want := New()
		require.NoError(t, err)
		require.Equal(t, want, p)
	})

	t.Run("invalid_file_content", func(t *testing.T) {
		t.Cleanup(func() { SetGlobalInstance(New()) })

		// WHEN.
		p, err := NewWithOpts(WithConfigPath("./config/testdata/invalid_cfg.yml"))

		// THEN.
		var want *Processor
		require.Error(t, err)
		require.Equal(t, want, p)
	})

	t.Run("empty_file_content", func(t *testing.T) {
		t.Cleanup(func() { SetGlobalInstance(New()) })

		// WHEN.
		p, err := NewWithOpts(WithConfigPath("./testdata/empty.yml"))

		// THEN.
		require.Error(t, err)
		require.Nil(t, p)
	})
}

func TestProcessor_PrintConfig(t *testing.T) {
	t.Run("successful", func(t *testing.T) {
		// GIVEN.
		r, w, err := os.Pipe()
		require.NoError(t, err)

		// Store previous stdout and replace it with our pipe.
		stdout := os.Stdout
		os.Stdout = w

		cfg := Config{
			General: General{
				OutputFormat:      OutputFormatText,
				PrintConfigOnInit: true,
			},
			Encoder: EncoderConfig{
				DisplayMapType:       true,
				DisplayPointerSymbol: true,
				DisplayStructName:    true,
				ExcludePatterns:      []string{`\d`, `.+@.+`},
				MaskValue:            DefaultMaskValue,
				UseJSONTagName:       true,
			},
		}

		p, err := NewWithOpts(WithConfig(&cfg))
		require.NoError(t, err)

		// WHEN.
		p.PrintConfig()

		// THEN.
		// Restore stdout.
		os.Stdout = stdout
		// Open file with valid output.
		f, err := os.Open("./testdata/valid_config_console_output.txt")

		// Read the expected output.
		want := make([]byte, 515)
		_, err = f.Read(want)
		require.NoError(t, err)

		// Read from the pipe.
		got := make([]byte, 515)
		_, err = r.Read(got)
		require.NoError(t, err)
		require.Equal(t, want, got)
	})

	t.Run("not initialized instance", func(t *testing.T) {
		// GIVEN.
		r, w, err := os.Pipe()
		require.NoError(t, err)

		// Store previous stdout and replace it with our pipe.
		stdout := os.Stdout
		os.Stdout = w

		var p *Processor

		// WHEN.
		p.PrintConfig()

		// THEN.
		// Restore stdout.
		os.Stdout = stdout
		want := []byte(censorIsNotInitializedMsg)

		// Read from the pipe.
		got := make([]byte, 25)
		_, err = r.Read(got)
		require.NoError(t, err)
		require.Equal(t, want, got)
	})
}

func TestTestProcessor_Clone(t *testing.T) {
	// GIVEN.
	cfg := Config{
		General: General{
			OutputFormat:      OutputFormatText,
			PrintConfigOnInit: false,
		},
		Encoder: EncoderConfig{
			DisplayMapType:       true,
			DisplayPointerSymbol: true,
			DisplayStructName:    true,
			ExcludePatterns:      nil,
			MaskValue:            "####",
			UseJSONTagName:       true,
		},
	}

	type s struct {
		Int64  int64
		String string `censor:"display"`
		Ptr    *int   `censor:"display"`
	}

	original, err := NewWithOpts(WithConfig(&cfg))
	require.NoError(t, err)

	// WHEN.
	clone, err := original.Clone()

	// THEN.
	require.NoError(t, err)
	// Check if the original and clone have the same configuration.
	require.Equal(t, original.cfg, clone.cfg)
	v := []s{{Int64: 123456789, String: "string", Ptr: new(int)}}
	// Ensure that the original and clone have the same behavior.
	require.Equal(t, original.Any(v), clone.Any(v))
}

func TestProcessor_Format(t *testing.T) {
	t.Run("test", func(t *testing.T) {
		// GIVEN.

		type address struct {
			City    string `censor:"display"`
			Country string `censor:"display"`
			Street  string
			Zip     int
		}
		type request struct {
			UserID   string  `censor:"display"` // Display value.
			Address  address `censor:"display"`
			Email    string  // Mask value.
			FullName string
		}
		r := request{
			UserID:   "123",
			Address:  address{City: "Kharkiv", Country: "UA", Street: "Nauky Avenue", Zip: 23335},
			Email:    "viktor.example.email@ggmail.com",
			FullName: "Viktor Pakhuchyi",
		}

		fmt.Println(string(Any(r)))
		//// THEN.
		//require.Equal(t, "null", string(got))
	})

	t.Run("nil value", func(t *testing.T) {
		// GIVEN.
		var v *string = nil

		// WHEN.
		got := Any(v)

		// THEN.
		require.Equal(t, "null", string(got))
	})

	t.Run("nil type", func(t *testing.T) {
		// GIVEN.
		var v interface{} = nil

		// WHEN.
		got := Any(v)

		// THEN.
		require.Equal(t, "null", string(got))
	})

	t.Run("processor method nil value", func(t *testing.T) {
		// GIVEN.
		processor := New()

		// WHEN.
		got := processor.Any(nil)

		// THEN.
		require.Equal(t, "null", string(got))
	})

	t.Run("nil type with text encoder", func(t *testing.T) {
		// GIVEN.
		processor, err := NewWithOpts(WithConfig(&Config{
			General: General{
				OutputFormat: OutputFormatText,
			},
			Encoder: EncoderConfig{
				MaskValue: DefaultMaskValue,
			},
		}))
		require.NoError(t, err)

		// WHEN.
		got := processor.Any(nil)

		// THEN.
		require.Equal(t, "nil", string(got))
	})
}

func TestProcessor_GetGlobalInstance(t *testing.T) {
	t.Run("successful", func(t *testing.T) {
		// GIVEN.
		want, err := NewWithOpts(
			WithConfig(
				&Config{
					General: General{
						OutputFormat:      OutputFormatText,
						PrintConfigOnInit: true,
					},
					Encoder: EncoderConfig{
						MaskValue: DefaultMaskValue,
					},
				},
			),
		)
		require.NoError(t, err)

		SetGlobalInstance(want)

		// WHEN.
		got := GetGlobalInstance()

		// THEN.
		require.Equal(t, want, got)
	})
}
