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
