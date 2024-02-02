package formatter

import (
	"reflect"
	"regexp"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/vpakhuchyi/censor/internal/models"
)

func TestFormatter_Map(t *testing.T) {
	f := Formatter{
		maskValue:         "[CENSORED]",
		displayStructName: false,
		displayMapType:    false,
	}

	t.Run("successful", func(t *testing.T) {
		require.NotPanics(t, func() {
			v := models.Map{
				Type: "map[string]int",
				Values: []models.KV{
					{Key: models.Value{Value: "foo", Kind: reflect.String}, Value: models.Value{Value: 1, Kind: reflect.Int}},
					{Key: models.Value{Value: "bar", Kind: reflect.String}, Value: models.Value{Value: 2, Kind: reflect.Int}},
				},
			}
			got := f.Map(v)
			exp := `map[foo: 1, bar: 2]`
			require.Equal(t, exp, got)
		})
	})

	t.Run("with_exclude_patterns", func(t *testing.T) {
		f := Formatter{
			maskValue:               "[CENSORED]",
			displayStructName:       false,
			displayMapType:          false,
			excludePatterns:         []string{`\d`, `\b[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Z|a-z]{2,}\b`},
			excludePatternsCompiled: []*regexp.Regexp{compiledRegExpDigit, compiledRegExpEmail},
		}

		require.NotPanics(t, func() {
			v := models.Map{
				Type: "map[string]int",
				Values: []models.KV{
					{Key: models.Value{Value: "foo", Kind: reflect.String}, Value: models.Value{Value: "hell0", Kind: reflect.String}},
					{Key: models.Value{Value: "test@exaample.com", Kind: reflect.String}, Value: models.Value{Value: "hello", Kind: reflect.String}},
				},
			}
			got := f.Map(v)
			exp := `map[foo: hell[CENSORED], [CENSORED]: hello]`
			require.Equal(t, exp, got)
		})
	})

	t.Run("successful_with_display_map_type", func(t *testing.T) {
		f := Formatter{
			maskValue:         "[CENSORED]",
			displayStructName: false,
			displayMapType:    true,
		}

		require.NotPanics(t, func() {
			v := models.Map{
				Type: "map[string]int",
				Values: []models.KV{
					{Key: models.Value{Value: "foo", Kind: reflect.String}, Value: models.Value{Value: 1, Kind: reflect.Int}},
					{Key: models.Value{Value: "bar", Kind: reflect.String}, Value: models.Value{Value: 2, Kind: reflect.Int}},
				},
			}
			got := f.Map(v)
			exp := `map[string]int[foo: 1, bar: 2]`
			require.Equal(t, exp, got)
		})
	})

	t.Run("empty", func(t *testing.T) {
		require.NotPanics(t, func() {
			v := models.Map{
				Type:   "map[string]int",
				Values: []models.KV{},
			}
			got := f.Map(v)
			exp := `map[]`
			require.Equal(t, exp, got)
		})
	})
}
