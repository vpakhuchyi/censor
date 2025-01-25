package cache

import "container/list"

// DefaultMaxCacheSize is the default maximum cache size.
const DefaultMaxCacheSize = 500

// New creates a new cache for type T.
func New[T comparable](size int) *Cache[T] {
	return &Cache[T]{
		size:  size,
		keys:  list.New(),
		cache: make(map[string]T),
	}
}

// Cache is a cache for type T.
type Cache[T comparable] struct {
	size  int
	keys  *list.List
	cache map[string]T
}

// Set adds a new key-value pair to the cache.
// If the cache size exceeds the defaultMaxCacheSize, the oldest key-value pair is removed.
// If the key already exists in the cache, the function will return immediately.
// If the key does not exist in the cache, the key-value pair is added.
func (c Cache[T]) Set(key string, value T) {
	if _, found := c.cache[key]; found {
		return
	}

	if c.keys.Len() >= c.size {
		oldestKey := c.keys.Front().Value.(string) //nolint
		delete(c.cache, oldestKey)
		c.keys.Remove(c.keys.Front())
	}

	c.keys.PushBack(key)
	c.cache[key] = value
}

// Get returns the value for the given key.
// If the key does not exist in the cache, the second return value is false.
func (c Cache[T]) Get(key string) (T, bool) {
	value, found := c.cache[key]

	return value, found
}
