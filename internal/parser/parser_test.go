package parser

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/vpakhuchyi/censor/config"
)

func Test(t *testing.T) {
	got := New(config.Parser{UseJSONTagName: true})
	exp := &Parser{useJSONTagName: true, censorFieldTag: DefaultCensorFieldTag}
	require.EqualValues(t, exp, got)
}
