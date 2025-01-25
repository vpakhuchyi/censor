package cache

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type Field struct {
	Name     string
	IsMasked bool
}

func TestFieldsCache_Get(t *testing.T) {
	// GIVEN
	key := "key"
	fields := []Field{{Name: "name"}}
	c := NewSlice[Field](10)
	c.Set(key, fields)

	// WHEN
	got, found := c.Get(key)

	// THEN
	require.True(t, found)
	require.Equal(t, fields, got)
}

func TestFieldsCache_Set(t *testing.T) {
	t.Run("cache already exists", func(t *testing.T) {
		// GIVEN
		key := "key"
		fields := []Field{{Name: "name"}}
		c := NewSlice[Field](10)
		c.Set(key, fields)

		// WHEN
		c.Set(key, []Field{{Name: "name2"}})

		// THEN
		got, found := c.Get(key)
		require.True(t, found)
		require.Equal(t, fields, got)
		require.Equal(t, 1, c.keys.Len())
		require.Len(t, c.cache, 1)
	})

	t.Run("ok", func(t *testing.T) {
		// GIVEN
		key := "key"
		fields := []Field{{Name: "name"}}
		c := NewSlice[Field](10)

		// WHEN
		c.Set(key, fields)

		// THEN
		got, found := c.Get(key)
		require.True(t, found)
		require.Equal(t, fields, got)
		require.Equal(t, 1, c.keys.Len())
		require.Len(t, c.cache, 1)
	})

	t.Run("cache size exceeds the defaultMaxCacheSize", func(t *testing.T) {
		// GIVEN
		key1 := "key1"
		key2 := "key2"
		fields1 := []Field{{Name: "name1"}}
		fields2 := []Field{{Name: "name2"}}
		c := NewSlice[Field](1)

		// WHEN
		c.Set(key1, fields1)
		c.Set(key2, fields2)

		// THEN
		got1, found1 := c.Get(key1)
		got2, found2 := c.Get(key2)
		require.Empty(t, got1)
		require.Equal(t, fields2, got2)
		require.False(t, found1)
		require.True(t, found2)
		require.Equal(t, 1, c.keys.Len())
		require.Len(t, c.cache, 1)
	})
}
