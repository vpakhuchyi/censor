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
			PrintConfigOnInit: true,
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

func TestConfig_YAMLMarshal(t *testing.T) {
	t.Run("default", func(t *testing.T) {
		// GIVEN.
		want, err := os.ReadFile("./testdata/default.yml")
		require.NoError(t, err)

		// WHEN.
		got, err := yaml.Marshal(DefaultConfig())

		// THEN.
		require.NoError(t, err)
		require.EqualValues(t, want, got)
		require.YAMLEq(t, string(want), string(got))
	})

	t.Run("custom", func(t *testing.T) {
		// GIVEN.
		cfg := Config{
			General: General{
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
