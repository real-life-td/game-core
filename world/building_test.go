package world

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewBuilding(t *testing.T) {
	expectedPoints := []*Node{NewNode(1, 2, 3), NewNode(4, 5, 6)}
	b := NewBuilding(0, expectedPoints)
	require.Equal(t, b.Id(), Id(0))
	require.Equal(t, expectedPoints, b.points)
}