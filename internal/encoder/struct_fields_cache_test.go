package encoder

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFieldsCache_Get(t *testing.T) {
	// GIVEN
	key := "key"
	fields := []Field{{Name: "name"}}
	cache := newFieldsCache(10)
	cache.Set(key, fields)

	// WHEN
	got, found := cache.Get(key)

	// THEN
	require.True(t, found)
	require.Equal(t, fields, got)
}

func TestFieldsCache_Set(t *testing.T) {
	t.Run("cache already exists", func(t *testing.T) {
		// GIVEN
		key := "key"
		fields := []Field{{Name: "name"}}
		cache := newFieldsCache(10)
		cache.Set(key, fields)

		// WHEN
		cache.Set(key, []Field{{Name: "name2"}})

		// THEN
		got, found := cache.Get(key)
		require.True(t, found)
		require.Equal(t, fields, got)
		require.Equal(t, 1, cache.keys.Len())
		require.Len(t, cache.cache, 1)
	})

	t.Run("ok", func(t *testing.T) {
		// GIVEN
		key := "key"
		fields := []Field{{Name: "name"}}
		cache := newFieldsCache(10)

		// WHEN
		cache.Set(key, fields)

		// THEN
		got, found := cache.Get(key)
		require.True(t, found)
		require.Equal(t, fields, got)
		require.Equal(t, 1, cache.keys.Len())
		require.Len(t, cache.cache, 1)
	})

	t.Run("cache size exceeds the defaultMaxCacheSize", func(t *testing.T) {
		// GIVEN
		key1 := "key1"
		key2 := "key2"
		fields1 := []Field{{Name: "name1"}}
		fields2 := []Field{{Name: "name2"}}
		cache := newFieldsCache(1)

		// WHEN
		cache.Set(key1, fields1)
		cache.Set(key2, fields2)

		// THEN
		got1, found1 := cache.Get(key1)
		got2, found2 := cache.Get(key2)
		require.Empty(t, got1)
		require.Equal(t, fields2, got2)
		require.False(t, found1)
		require.True(t, found2)
		require.Equal(t, 1, cache.keys.Len())
		require.Len(t, cache.cache, 1)
	})
}
