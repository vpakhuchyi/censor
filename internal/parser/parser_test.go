package parser

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/vpakhuchyi/censor/config"
)

func TestNew(t *testing.T) {
	require.EqualValues(t, &Parser{useJSONTagName: false, censorFieldTag: DefaultCensorFieldTag}, New())
}

func TestNewWithConfig(t *testing.T) {
	got := NewWithConfig(config.Parser{UseJSONTagName: true})
	exp := &Parser{useJSONTagName: true, censorFieldTag: DefaultCensorFieldTag}
	require.EqualValues(t, exp, got)
}

func TestParser_UseJSONTagName(t *testing.T) {
	f := &Parser{useJSONTagName: false, censorFieldTag: DefaultCensorFieldTag}
	f.UseJSONTagName(true)
	require.EqualValues(t, f, &Parser{useJSONTagName: true, censorFieldTag: DefaultCensorFieldTag})
}
