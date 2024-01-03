package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDefault(t *testing.T) {
	got := Default()
	exp := Config{
		Parser: Parser{
			UseJSONTagName: false,
		},
		Formatter: Formatter{
			MaskValue:            DefaultMaskValue,
			DisplayPointerSymbol: false,
			DisplayStructName:    false,
			DisplayMapType:       false,
			ExcludePatterns:      nil,
		},
	}

	require.EqualValues(t, exp, got)
}

func TestConfig_GetParserConfig(t *testing.T) {
	got := Default().GetParserConfig()
	exp := Parser{
		UseJSONTagName: false,
	}

	require.EqualValues(t, exp, got)
}

func TestConfig_GetFormatterConfig(t *testing.T) {
	got := Default().GetFormatterConfig()
	exp := Formatter{
		MaskValue:            DefaultMaskValue,
		DisplayPointerSymbol: false,
		DisplayStructName:    false,
		DisplayMapType:       false,
		ExcludePatterns:      nil,
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
				Parser: Parser{
					UseJSONTagName: true,
				},
				Formatter: Formatter{
					MaskValue:            "[CENSORED]",
					DisplayPointerSymbol: true,
					DisplayStructName:    true,
					DisplayMapType:       true,
					ExcludePatterns:      []string{`\d`},
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
				path: "./testdata/invalid-cfg.yml",
			},
			want:    Config{},
			wantErr: true,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := FromFile(tt.args.path)
			require.Equal(t, tt.want, got)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
