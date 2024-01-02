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
			MaskValue:                    DefaultMaskValue,
			DisplayPointerSymbol:         false,
			DisplayStructName:            false,
			DisplayMapType:               false,
			ExcludePatterns:              nil,
			Float32MaxSignificantFigures: Float32MaxSignificantFigures,
			Float64MaxSignificantFigures: Float64MaxSignificantFigures,
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
		MaskValue:                    DefaultMaskValue,
		DisplayPointerSymbol:         false,
		DisplayStructName:            false,
		DisplayMapType:               false,
		ExcludePatterns:              nil,
		Float32MaxSignificantFigures: Float32MaxSignificantFigures,
		Float64MaxSignificantFigures: Float64MaxSignificantFigures,
	}

	require.EqualValues(t, exp, got)
}
