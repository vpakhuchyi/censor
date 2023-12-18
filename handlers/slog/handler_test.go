package sloghandler

import (
	"bufio"
	"bytes"
	"encoding/json"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/require"
)

type source struct {
	Func string `json:"function"`
	File string `json:"-"`
}

type logEntry struct {
	Level   string  `json:"level"`
	Msg     string  `json:"msg"`
	Payload string  `json:"payload"`
	Source  *source `json:"source,omitempty"`
}

type address struct {
	City    string `censor:"display"`
	Country string `censor:"display"`
	Street  string
	Zip     int
}

func TestNewHandler(t *testing.T) {
	payload := address{
		City:    "Kyiv",
		Country: "Ukraine",
		Street:  "Khreshchatyk",
		Zip:     12345,
	}

	type args struct {
		opts []Option
	}

	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "with_default_handler_options",
			want: `{"level": "INFO",
    				 "msg": "test",
					 "payload": "{City: Kyiv, Country: Ukraine, Street: [CENSORED], Zip: [CENSORED]}"}`,
		},
		{
			name: "with_add_source_option",
			want: `{"level": "INFO",
    				 "msg": "test",
    				 "payload": "{City: Kyiv, Country: Ukraine, Street: [CENSORED], Zip: [CENSORED]}",
    				 "source": {"function": "github.com/vpakhuchyi/censor/handlers/slog.TestNewHandler.func2"}}`,
			args: args{
				opts: []Option{WithAddSource()},
			},
		},
		{
			name: "with_level_error_option",
			want: "",
			args: args{
				opts: []Option{WithLevel(slog.LevelError)},
			},
		},
		{
			name: "with_replace_attr_option",
			want: `{"level": "TEST",
    				 "msg": "test",
					 "payload": "{City: Kyiv, Country: Ukraine, Street: [CENSORED], Zip: [CENSORED]}"}`,
			args: args{
				opts: []Option{WithReplaceAttr(func(groups []string, a slog.Attr) slog.Attr {
					if a.Key == "level" {
						return slog.Any("level", "TEST")
					}

					return a
				})},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			out := bufio.NewWriter(&buf)
			options := []Option{WithOut(out)}

			opts := append(options, tt.args.opts...)
			handler := NewJSONHandler(opts...)
			log := slog.New(handler)
			log.Info("test", slog.Any("payload", payload))

			require.NoError(t, out.Flush())
			got := buf.String()
			if tt.want == "" {
				require.Equal(t, tt.want, got)
				return
			}

			require.JSONEq(t, tt.want, prepareLogEntry(t, got))
		})
	}
}

func prepareLogEntry(t *testing.T, s string) string {
	if s == "" {
		return s
	}

	logE := logEntry{}
	err := json.Unmarshal([]byte(s), &logE)
	require.NoError(t, err)

	log, err := json.Marshal(logE)
	require.NoError(t, err)

	return string(log)
}
