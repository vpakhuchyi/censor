package zerologhandler

import (
	"bufio"
	"bytes"
	"encoding/json"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
	"github.com/vpakhuchyi/censor"
)

type address struct {
	City    string `censor:"display"`
	Country string `censor:"display"`
	Street  string
	Zip     int
}

type logEntry struct {
	Level   string `json:"level"`
	Msg     string `json:"message"`
	Payload any    `json:"payload"`
	Key     string `json:"key,omitempty"`
}

func TestNewHandler(t *testing.T) {
	payload := address{
		City:    "Kyiv",
		Country: "Ukraine",
		Street:  "Khreshchatyk",
		Zip:     12345,
	}

	t.Run("any: without output option", func(t *testing.T) {
		// GIVEN
		var buf bytes.Buffer
		out := bufio.NewWriter(&buf)
		handler := New().Output(out)

		// WHEN
		handler.Info().Any("payload", payload).Msg("test")

		// THEN
		require.NoError(t, out.Flush())
		want := `{
					"level":"info",
					"payload":{"City": "Kyiv","Country": "Ukraine","Street": "[CENSORED]","Zip": "[CENSORED]"},
					"message":"test"
				}`
		require.JSONEq(t, want, prepareLogEntry(t, buf.String()))
	})

	t.Run("interface: without output option", func(t *testing.T) {
		// GIVEN
		var buf bytes.Buffer
		out := bufio.NewWriter(&buf)
		handler := New().Output(out)

		// WHEN
		handler.Info().Interface("payload", payload).Msg("test")

		// THEN
		require.NoError(t, out.Flush())
		want := `{
						"level":"info",
						"payload":{"City": "Kyiv","Country": "Ukraine","Street": "[CENSORED]","Zip": "[CENSORED]"},
						"message":"test"
					}`
		require.JSONEq(t, want, prepareLogEntry(t, buf.String()))
	})

	t.Run("with zerolog option", func(t *testing.T) {
		// GIVEN
		var buf bytes.Buffer
		out := bufio.NewWriter(&buf)
		zl := zerolog.New(&buf).With().Str("key", "value").Logger()
		log := New(WithZerolog(&zl))
		handler := log.Output(out)

		// WHEN
		handler.Info().Any("payload", payload).Msg("test")

		// THEN
		require.NoError(t, out.Flush())
		want := `{
					"level":"info",
					"payload":{"City": "Kyiv","Country": "Ukraine","Street": "[CENSORED]","Zip": "[CENSORED]"},
					"message":"test",
					"key":"value"
				}`
		require.JSONEq(t, want, prepareLogEntry(t, buf.String()))
	})

	t.Run("with censor option", func(t *testing.T) {
		// GIVEN
		var buf bytes.Buffer
		out := bufio.NewWriter(&buf)
		censor := censor.New()
		log := New(WithCensor(censor))
		handler := log.Output(out)

		// WHEN
		handler.Info().Any("payload", payload).Msg("test")

		// THEN
		require.NoError(t, out.Flush())
		want := `{
					"level":"info",
					"payload":{"City": "Kyiv","Country": "Ukraine","Street": "[CENSORED]","Zip": "[CENSORED]"},
					"message":"test"
				}`
		require.JSONEq(t, want, prepareLogEntry(t, buf.String()))
	})
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
