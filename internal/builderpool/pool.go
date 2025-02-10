package builderpool

import (
	"bytes"
	"sync"
)

// Define a pool of strings.Builder instances.
var builderPool = sync.Pool{
	New: func() interface{} {
		return bytes.NewBuffer(make([]byte, 0, 32))
	},
}

// Get fetches a builder from the pool.
func Get() *bytes.Buffer {
	v, ok := builderPool.Get().(*bytes.Buffer)
	if !ok {
		v = bytes.NewBuffer(make([]byte, 0, 32))
	}

	return v
}

// Put resets the builder and puts it back into the pool.
func Put(builder *bytes.Buffer) {
	builder.Reset()
	builderPool.Put(builder)
}
