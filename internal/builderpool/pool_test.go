package builderpool

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGet(t *testing.T) {
	// GIVEN.
	builder := Get()
	builder.WriteString("test")
	Put(builder)

	// WHEN.
	builder2 := Get()

	// THEN.
	require.Equal(t, []byte{}, builder2.Bytes())
	require.Equal(t, 0, builder2.Len())
	require.Equal(t, 32, builder2.Cap())
}
