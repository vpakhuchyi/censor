package zerologhandler

import (
	"strings"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vpakhuchyi/censor"
	"github.com/vpakhuchyi/censor/internal/encoder"
)

func TestNewHandler(t *testing.T) {
	c, err := censor.NewWithOpts(censor.WithConfig(&censor.Config{
		Encoder: encoder.Config{
			MaskValue:            censor.DefaultMaskValue,
			DisplayPointerSymbol: false,
			DisplayStructName:    false,
			DisplayMapType:       false,
			ExcludePatterns:      []string{`#sensitive#`},
		},
	}))
	require.NoError(t, err)

	msg := "#sensitive# msg"
	key := "#sensitive# key"

	value := struct {
		Name  string `censor:"display"`
		Text  string `censor:"display"`
		Email string
	}{
		Name:  "Petro Perekotypole",
		Text:  "some text with #sensitive# data",
		Email: "example@example.com",
	}

	var b strings.Builder

	handler := New(WithCensor(c))

	oldInterfaceMarshalFunc := zerolog.InterfaceMarshalFunc
	defer func() {
		zerolog.InterfaceMarshalFunc = oldInterfaceMarshalFunc
	}()

	zerolog.InterfaceMarshalFunc = handler.NewInterfaceMarshal
	l := zerolog.New(handler.Writer(&b))

	tests := map[string]struct {
		fn  func(log zerolog.Logger) *zerolog.Event
		exp string
	}{
		"sensitive key": {
			fn:  func(log zerolog.Logger) *zerolog.Event { return l.Warn().Str(key, "some msg") },
			exp: `{"level":"warn","[CENSORED] key":"some msg"}` + "\n",
		},
		"sensitive msg": {
			fn:  func(log zerolog.Logger) *zerolog.Event { return l.Warn().Str("some key", msg) },
			exp: `{"level":"warn","some key":"[CENSORED] msg"}` + "\n",
		},
		"object": {
			fn:  func(log zerolog.Logger) *zerolog.Event { return l.Warn().Any("o", value) },
			exp: `{"level":"warn","o":{Name: Petro Perekotypole, Text: some text with [CENSORED] data, Email: [CENSORED]}}` + "\n",
		},
	}

	for testCase, tt := range tests {
		t.Run(testCase, func(t *testing.T) {
			b.Reset()

			tt.fn(l).Send()

			assert.Equal(t, tt.exp, b.String())
		})
	}
}
