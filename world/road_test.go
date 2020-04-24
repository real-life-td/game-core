package world

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewRoad(t *testing.T) {
	n1, n2 := NewNode(1, 2, 3), NewNode(4, 5, 6)
	r := NewRoad(7, n1, n2, 8)

	require.Equal(t, Id(7), r.Id())
	require.Equal(t, n1, r.Node1())
	require.Equal(t, n2, r.Node2())
	require.Equal(t, 8, r.Cost())
}
