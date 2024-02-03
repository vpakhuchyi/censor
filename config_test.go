package censor

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func TestDefault(t *testing.T) {
	got := Default()
	exp := Config{
		General: General{
			PrintConfigOnInit: true,
		},
		Parser: ParserConfig{
			UseJSONTagName: false,
		},
		Formatter: FormatterConfig{
			MaskValue:            DefaultMaskValue,
			DisplayPointerSymbol: false,
			DisplayStructName:    false,
			DisplayMapType:       false,
			ExcludePatterns:      nil,
		},
	}

	require.EqualValues(t, exp, got)
}

func TestFromFile(t *testing.T) {
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
				Parser: ParserConfig{
					UseJSONTagName: true,
				},
				Formatter: FormatterConfig{
					MaskValue:            "[CENSORED]",
					DisplayPointerSymbol: true,
					DisplayStructName:    true,
					DisplayMapType:       true,
					ExcludePatterns:      []string{`\d`, `^\w$`},
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

func TestYAMLMarshal(t *testing.T) {
	t.Run("default", func(t *testing.T) {
		want, err := os.ReadFile("./testdata/default.yml")
		require.NoError(t, err)

		got, err := yaml.Marshal(Default())
		require.NoError(t, err)
		require.EqualValues(t, want, got)
		require.YAMLEq(t, string(want), string(got))
	})

	t.Run("custom", func(t *testing.T) {
		want, err := os.ReadFile("./testdata/cfg.yml")
		require.NoError(t, err)

		cfg := Config{
			General: General{
				PrintConfigOnInit: true,
			},
			Parser: ParserConfig{
				UseJSONTagName: true,
			},
			Formatter: FormatterConfig{
				MaskValue:            "[CENSORED]",
				DisplayPointerSymbol: true,
				DisplayStructName:    true,
				DisplayMapType:       true,
				ExcludePatterns:      []string{`\d`, `^\w$`},
			},
		}

		got, err := yaml.Marshal(cfg)
		require.NoError(t, err)
		require.EqualValues(t, want, got)
		require.YAMLEq(t, string(want), string(got))
	})
}
