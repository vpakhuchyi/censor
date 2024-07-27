package builderpool

import (
	"strings"
	"sync"
)

// Define a pool of strings.Builder instances.
var builderPool = sync.Pool{
	New: func() interface{} {
		return &strings.Builder{}
	},
}

// Get fetches a builder from the pool.
func Get() *strings.Builder {
	return builderPool.Get().(*strings.Builder)
}

// Put resets the builder and puts it back into the pool.
func Put(builder *strings.Builder) {
	builder.Reset()
	builderPool.Put(builder)
}
