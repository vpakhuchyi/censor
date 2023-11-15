package parser

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	require.EqualValues(t, &Parser{useJSONTagName: false, CensorFieldTag: DefaultCensorFieldTag}, New())
}

func TestParser_UseJSONTagName(t *testing.T) {
	f := &Parser{useJSONTagName: false, CensorFieldTag: gDefaultCensorFieldTag}
	f.UseJSONTagName(true)
	require.EqualValues(t, f, &Parser{useJSONTagName: true, CensorFieldTag: DefaultCensorFieldTag})
}
