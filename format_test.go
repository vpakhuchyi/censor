package censor

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFormat(t *testing.T) {
	type S struct {
		String1 string `censor:"display"`
		String2 string
		Array   [2]string `censor:"display"`
	}

	// GIVEN a default processor instance and a payload.
	p := S{
		String1: "value",
		String2: "value",
		Array:   [2]string{"value", "value"},
	}

	// WHEN the Format func is called.
	got := Format(p)

	// THEN the returned string contains the masked values.
	want := `{"String1":"value","String2":"[CENSORED]","Array":["value","value"]}`
	require.JSONEq(t, want, got)
}
