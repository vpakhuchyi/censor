package formatter

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/vpakhuchyi/censor/internal/models"
)

func TestFormatter_Bool(t *testing.T) {
	f := Formatter{
		MaskValue:         DefaultMaskValue,
		DisplayStructName: false,
		DisplayMapType:    false,
	}

	tests := map[string]struct {
		value models.Value
		exp   string
	}{
		"bool": {
			value: models.Value{
				Value: true,
				Kind:  reflect.Bool,
			},
			exp: "true",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got := f.Bool(models.Bool(tt.value))
			require.Equal(t, tt.exp, got)
		})
	}
}
