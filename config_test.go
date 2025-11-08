package censor

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func TestConfig_Default(t *testing.T) {
	// WHEN.
	got := DefaultConfig()
	exp := Config{
		General: General{
			OutputFormat:      OutputFormatJSON,
			PrintConfigOnInit: false,
		},
		Encoder: EncoderConfig{
			DisplayMapType:       false,
			DisplayPointerSymbol: false,
			DisplayStructName:    false,
			ExcludePatterns:      nil,
			MaskValue:            DefaultMaskValue,
			UseJSONTagName:       false,
		},
	}

	// THEN.
	require.EqualValues(t, exp, got)
}

func TestConfig_FromFile(t *testing.T) {
	type args struct {
		path string
	}
	tests := map[string]struct {
		args    args
		want    Config
		wantErr bool
	}{
		"successful": {
			args: args{
				path: "./testdata/cfg.yml",
			},
			want: Config{
				General: General{
					OutputFormat:      OutputFormatText,
					PrintConfigOnInit: true,
				},
				Encoder: EncoderConfig{
					DisplayMapType:       true,
					DisplayPointerSymbol: true,
					DisplayStructName:    true,
					ExcludePatterns:      []string{`\d`, `^\w$`},
					MaskValue:            "[CENSORED]",
					UseJSONTagName:       true,
				},
			},
			wantErr: false,
		},
		"empty_cfg_file_path": {
			args: args{
				path: "",
			},
			want:    Config{},
			wantErr: true,
		},
		"empty_cfg_file_content": {
			args: args{
				path: "./testdata/empty.yml",
			},
			want:    Config{},
			wantErr: false,
		},
		"invalid_file_path": {
			args: args{
				path: "./testdata/non-existing.yml",
			},
			want:    Config{},
			wantErr: true,
		},
		"invalid_file_content": {
			args: args{
				path: "./testdata/invalid_cfg.yml",
			},
			want:    Config{},
			wantErr: true,
		},
		"directory_traversal_attack": {
			args: args{
				path: "../../../etc/passwd",
			},
			want:    Config{},
			wantErr: true,
		},
		"directory_traversal_with_dots": {
			args: args{
				path: "./testdata/../../../etc/passwd",
			},
			want:    Config{},
			wantErr: true,
		},
		"invalid_extension_txt": {
			args: args{
				path: "./testdata/config.txt",
			},
			want:    Config{},
			wantErr: true,
		},
		"invalid_extension_json": {
			args: args{
				path: "./testdata/config.json",
			},
			want:    Config{},
			wantErr: true,
		},
		"valid_yaml_extension": {
			args: args{
				path: "./testdata/cfg.yaml",
			},
			want: Config{
				General: General{
					OutputFormat:      OutputFormatText,
					PrintConfigOnInit: true,
				},
				Encoder: EncoderConfig{
					DisplayMapType:       true,
					DisplayPointerSymbol: true,
					DisplayStructName:    true,
					ExcludePatterns:      []string{`\d`, `^\w$`},
					MaskValue:            "[CENSORED]",
					UseJSONTagName:       true,
				},
			},
			wantErr: false,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := ConfigFromFile(tt.args.path)
			require.Equal(t, tt.want, got)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestConfig_Validate(t *testing.T) {
	tests := map[string]struct {
		cfg             Config
		wantErr         string
		wantErrContains bool
	}{
		"default": {
			cfg: DefaultConfig(),
		},
		"text_format": {
			cfg: func() Config {
				cfg := DefaultConfig()
				cfg.General.OutputFormat = OutputFormatText
				return cfg
			}(),
		},
		"empty_output_format": {
			cfg: func() Config {
				cfg := DefaultConfig()
				cfg.General.OutputFormat = ""
				return cfg
			}(),
			wantErr: "invalid output format: \"\", must be \"text\" or \"json\"",
		},
		"mask_value_empty": {
			cfg: func() Config {
				cfg := DefaultConfig()
				cfg.Encoder.MaskValue = ""
				return cfg
			}(),
			wantErr: "mask value cannot be empty",
		},
		"invalid_output_format": {
			cfg: func() Config {
				cfg := DefaultConfig()
				cfg.General.OutputFormat = "xml"
				return cfg
			}(),
			wantErr: "invalid output format: \"xml\", must be \"text\" or \"json\"",
		},
		"too_many_patterns": {
			cfg: func() Config {
				cfg := DefaultConfig()
				patterns := make([]string, 51)
				for i := range patterns {
					patterns[i] = "pattern"
				}
				cfg.Encoder.ExcludePatterns = patterns
				return cfg
			}(),
			wantErr: "too many exclude patterns (max 50): 51",
		},
		"invalid_pattern": {
			cfg: func() Config {
				cfg := DefaultConfig()
				cfg.Encoder.ExcludePatterns = []string{"("}
				return cfg
			}(),
			wantErr:         "invalid exclude pattern \"(\"",
			wantErrContains: true,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			err := tt.cfg.Validate()
			if tt.wantErr == "" {
				require.NoError(t, err)
				return
			}

			require.Error(t, err)
			if tt.wantErrContains {
				require.ErrorContains(t, err, tt.wantErr)
			} else {
				require.EqualError(t, err, tt.wantErr)
			}
		})
	}
}

func TestConfig_YAMLMarshal(t *testing.T) {
	t.Run("default", func(t *testing.T) {
		// GIVEN.
		want, err := os.ReadFile("./testdata/default.yml")
		require.NoError(t, err)

		// WHEN.
		got, err := yaml.Marshal(DefaultConfig())

		// THEN.
		require.NoError(t, err)
		require.YAMLEq(t, string(want), string(got))
	})

	t.Run("custom", func(t *testing.T) {
		// GIVEN.
		cfg := Config{
			General: General{
				OutputFormat:      OutputFormatText,
				PrintConfigOnInit: true,
			},
			Encoder: EncoderConfig{
				DisplayMapType:       true,
				DisplayPointerSymbol: true,
				DisplayStructName:    true,
				ExcludePatterns:      []string{`\d`, `^\w$`},
				MaskValue:            "[CENSORED]",
				UseJSONTagName:       true,
			},
		}

		want, err := os.ReadFile("./testdata/cfg.yml")
		require.NoError(t, err)

		// WHEN.
		got, err := yaml.Marshal(cfg)

		// THEN.
		require.NoError(t, err)
		require.EqualValues(t, want, got)
		require.YAMLEq(t, string(want), string(got))
	})
}

func TestConfig_ToString(t *testing.T) {
	// GIVEN.
	cfg := Config{
		General: General{
			OutputFormat:      OutputFormatText,
			PrintConfigOnInit: true,
		},
		Encoder: EncoderConfig{
			DisplayMapType:       true,
			DisplayPointerSymbol: true,
			DisplayStructName:    true,
			ExcludePatterns:      []string{`\d`, `^\w$`},
			MaskValue:            "[CENSORED]",
			UseJSONTagName:       true,
		},
	}

	// WHEN.
	got := cfg.ToString()

	// THEN.
	exp := "---------------------------------------------------------------------\n" +
		"          Censor is configured with the following settings:\n" +
		"---------------------------------------------------------------------\n" +
		"general:\n" +
		"    output-format: text\n" +
		"    print-config-on-init: true\n" +
		"encoder:\n" +
		"    display-map-type: true\n" +
		"    display-pointer-symbol: true\n" +
		"    display-struct-name: true\n" +
		"    exclude-patterns:\n" +
		"        - \\d\n" +
		"        - ^\\w$\n" +
		"    mask-value: '[CENSORED]'\n" +
		"    use-json-tag-name: true\n" +
		"---------------------------------------------------------------------\n"
	require.EqualValues(t, exp, got)
}
