package sloghandler

import (
	"bufio"
	"bytes"
	"fmt"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tidwall/sjson"
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

	t.Run("without output option", func(t *testing.T) {
		// GIVEN
		handler := NewJSONHandler()
		log := slog.New(handler)
		// WHEN
		log.Info("test", slog.Any("payload", payload))
	})

	t.Run("with default handler options", func(t *testing.T) {
		// GIVEN
		var buf bytes.Buffer
		out := bufio.NewWriter(&buf)
		log := slog.New(NewJSONHandler([]Option{WithOut(out)}...))
		want := `{
					"level":"INFO",
    				"msg":"test",
    				"payload":{
        				"City":"Kyiv",
        				"Country":"Ukraine",
        				"Street":"[CENSORED]",
						"Zip":"[CENSORED]"
    				}
				}`

		// WHEN
		log.Info("test", slog.Any("payload", payload))

		// THEN
		require.NoError(t, out.Flush())
		got, err := sjson.Delete(buf.String(), "time")
		require.NoError(t, err)

		fmt.Println(got)
		require.JSONEq(t, want, got)
	})

	t.Run("with add source option", func(t *testing.T) {
		// GIVEN
		var buf bytes.Buffer
		out := bufio.NewWriter(&buf)
		log := slog.New(NewJSONHandler([]Option{WithOut(out), WithAddSource()}...))
		want := `{
					"level": "INFO",
    				"source": {
        				"function": "github.com/vpakhuchyi/censor/handlers/slog.TestNewHandler.func3"
    				},
    				"msg": "test",
    				"payload": {
        				"City": "Kyiv",
        				"Country": "Ukraine",
        				"Street": "[CENSORED]",
        				"Zip": "[CENSORED]"
    				}
				}`

		// WHEN
		log.Info("test", slog.Any("payload", payload))

		// THEN
		require.NoError(t, out.Flush())
		got, err := sjson.Delete(buf.String(), "time")
		require.NoError(t, err)
		got, err = sjson.Delete(got, "source.file")
		require.NoError(t, err)
		got, err = sjson.Delete(got, "source.line")
		require.NoError(t, err)

		require.JSONEq(t, want, got)
	})

	t.Run("with level error option", func(t *testing.T) {
		// GIVEN
		var buf bytes.Buffer
		out := bufio.NewWriter(&buf)
		log := slog.New(NewJSONHandler([]Option{WithOut(out), WithLevel(slog.LevelError)}...))

		// WHEN
		log.Info("test", slog.Any("payload", payload))

		// THEN
		require.NoError(t, out.Flush())
		got := buf.String()
		require.Equal(t, "", got)
	})

	t.Run("with replace attr option", func(t *testing.T) {
		// GIVEN
		var buf bytes.Buffer
		out := bufio.NewWriter(&buf)
		log := slog.New(NewJSONHandler([]Option{
			WithOut(out),
			WithReplaceAttr(func(groups []string, a slog.Attr) slog.Attr {
				if a.Key == "level" {
					return slog.Any("level", "TEST")
				}

				return a
			})}...),
		)
		want := `{
    				"level": "TEST",
    				"msg": "test",
    				"payload": {
        				"City": "Kyiv",
        				"Country": "Ukraine",
        				"Street": "[CENSORED]",
        				"Zip": "[CENSORED]"
    				}
				}`

		// WHEN
		log.Info("test", slog.Any("payload", payload))

		// THEN
		require.NoError(t, out.Flush())
		got, err := sjson.Delete(buf.String(), "time")
		require.NoError(t, err)

		require.JSONEq(t, want, got)
	})
}
