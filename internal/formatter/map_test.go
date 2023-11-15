package formatter

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/vpakhuchyi/censor/config"
	"github.com/vpakhuchyi/censor/internal/models"
)

func TestFormatter_Map(t *testing.T) {
	f := Formatter{
		maskValue:         config.DefaultMaskValue,
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

	t.Run("successful_with_display_map_type", func(t *testing.T) {
		f := Formatter{
			maskValue:         config.DefaultMaskValue,
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
