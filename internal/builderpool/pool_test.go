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
	require.Equal(t, "", builder2.String())
	require.Equal(t, 0, builder2.Len())
	require.Equal(t, 0, builder2.Cap())
}
