package parser

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test(t *testing.T) {
	got := New(Config{UseJSONTagName: true})
	exp := &Parser{useJSONTagName: true, censorFieldTag: DefaultCensorFieldTag}
	require.EqualValues(t, exp, got)
}
