package cache

import (
	"container/list"
	"reflect"
	"sync"
)

// TypeCache is a generic cache that uses reflect.Type as the key.
// This avoids string concatenation and allocations when building cache keys.
type TypeCache[T any] struct {
	mu    sync.RWMutex
	size  int
	keys  *list.List
	cache map[reflect.Type]T
}

// NewTypeCache creates a new cache indexed by reflect.Type.
func NewTypeCache[T any](size int) *TypeCache[T] {
	return &TypeCache[T]{
		size:  size,
		keys:  list.New(),
		cache: make(map[reflect.Type]T, size),
	}
}

// Set adds a new type-value pair to the cache.
// If the cache size exceeds the limit, the oldest entry is removed.
// If the type already exists in the cache, the function returns immediately.
func (c *TypeCache[T]) Set(t reflect.Type, value T) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, found := c.cache[t]; found {
		return
	}

	if c.keys.Len() >= c.size {
		oldestKey, ok := c.keys.Front().Value.(reflect.Type)
		if !ok {
			c.keys.Remove(c.keys.Front())

			return
		}

		delete(c.cache, oldestKey)
		c.keys.Remove(c.keys.Front())
	}

	c.keys.PushBack(t)
	c.cache[t] = value
}

// Get returns the value for the given type and a boolean indicating if it was found.
func (c *TypeCache[T]) Get(t reflect.Type) (T, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	value, found := c.cache[t]

	return value, found
}
