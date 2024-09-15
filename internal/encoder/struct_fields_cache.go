package encoder

import (
	"container/list"
)

const defaultMaxCacheSize = 500

// newFieldsCache creates a new cache for struct fields.
func newFieldsCache(size int) *fieldsCache {
	return &fieldsCache{
		size:  size,
		keys:  list.New(),
		cache: make(map[string][]Field),
	}
}

// fieldsCache is a cache for struct fields.
type fieldsCache struct {
	size  int
	keys  *list.List
	cache map[string][]Field
}

// Set adds a new key-value pair to the cache.
// If the cache size exceeds the defaultMaxCacheSize, the oldest key-value pair is removed.
// If the key already exists in the cache, the function will return immediately.
// If the key does not exist in the cache, the key-value pair is added.
func (c fieldsCache) Set(key string, value []Field) {
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
func (c fieldsCache) Get(key string) ([]Field, bool) {
	value, found := c.cache[key]

	return value, found
}
