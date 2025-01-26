package cache

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCache_Get(t *testing.T) {
	// GIVEN
	key, value := "key", "value"
	c := New[string](10)
	c.Set(key, value)

	// WHEN
	got, found := c.Get(key)

	// THEN
	require.True(t, found)
	require.Equal(t, value, got)
}

func TestCache_Set(t *testing.T) {
	t.Run("cache already exists", func(t *testing.T) {
		// GIVEN
		key, value := "key", "value"
		c := New[string](10)
		c.Set(key, value)

		// WHEN
		c.Set(key, "value2")

		// THEN
		got, found := c.Get(key)
		require.True(t, found)
		require.Equal(t, value, got)
		require.Equal(t, 1, c.keys.Len())
		require.Len(t, c.cache, 1)
	})

	t.Run("ok", func(t *testing.T) {
		// GIVEN
		key, value := "key", "value"
		c := New[string](10)

		// WHEN
		c.Set(key, value)

		// THEN
		got, found := c.Get(key)
		require.True(t, found)
		require.Equal(t, value, got)
		require.Equal(t, 1, c.keys.Len())
		require.Len(t, c.cache, 1)
	})

	t.Run("cache size exceeds the defaultMaxCacheSize", func(t *testing.T) {
		// GIVEN
		key1, value1 := "key1", "value1"
		key2, value2 := "key2", "value2"
		c := New[string](1)

		// WHEN
		c.Set(key1, value1)
		c.Set(key2, value2)

		// THEN
		got1, found1 := c.Get(key1)
		got2, found2 := c.Get(key2)
		require.Empty(t, got1)
		require.Equal(t, value2, got2)
		require.False(t, found1)
		require.True(t, found2)
		require.Equal(t, 1, c.keys.Len())
		require.Len(t, c.cache, 1)
	})
}
