package world

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewNode(t *testing.T) {
	n := NewNode(1, 2, 3)

	require.Equal(t, uint64(1), n.Id())
	require.Equal(t, 2, n.X())
	require.Equal(t, 3, n.Y())
}
