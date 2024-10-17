package builderpool

import (
	"sync"

	"github.com/vpakhuchyi/censor/internal/builder"
)

// Define a pool of strings.Builder instances.
var builderPool = sync.Pool{
	New: func() interface{} {
		return builder.New()
	},
}

// Get fetches a builder from the pool.
func Get() *builder.Builder {
	return builderPool.Get().(*builder.Builder)
}

// Put resets the builder and puts it back into the pool.
func Put(builder *builder.Builder) {
	builder.Reset()
	builderPool.Put(builder)
}
