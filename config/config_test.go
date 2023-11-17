package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDefault(t *testing.T) {
	got := Default()
	exp := Config{
		Parser: Parser{
			UseJSONTagName: false,
		},
		Formatter: Formatter{
			MaskValue:         DefaultMaskValue,
			DisplayStructName: false,
			DisplayMapType:    false,
			ExcludePatterns:   nil,
		},
	}

	require.EqualValues(t, exp, got)
}

func TestConfig_GetParserConfig(t *testing.T) {
	got := Default().GetParserConfig()
	exp := Parser{
		UseJSONTagName: false,
	}

	require.EqualValues(t, exp, got)
}

func TestConfig_GetFormatterConfig(t *testing.T) {
	got := Default().GetFormatterConfig()
	exp := Formatter{
		MaskValue:         DefaultMaskValue,
		DisplayStructName: false,
		DisplayMapType:    false,
		ExcludePatterns:   nil,
	}

	require.EqualValues(t, exp, got)
}
