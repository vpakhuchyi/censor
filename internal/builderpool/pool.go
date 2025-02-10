package builderpool

import (
	"bytes"
	"sync"
)

const defaultBufferCapacity = 32

// Define a pool of bytes.Buffer instances.
var builderPool = sync.Pool{
	New: func() interface{} {
		return bytes.NewBuffer(make([]byte, 0, defaultBufferCapacity))
	},
}

// Get fetches a builder from the pool.
func Get() *bytes.Buffer {
	v, ok := builderPool.Get().(*bytes.Buffer)
	if !ok {
		v = bytes.NewBuffer(make([]byte, 0, defaultBufferCapacity))
	}

	return v
}

// Put resets the builder and puts it back into the pool.
func Put(builder *bytes.Buffer) {
	builder.Reset()
	builderPool.Put(builder)
}
