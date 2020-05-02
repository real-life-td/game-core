package world

import (
	"github.com/real-life-td/math/primitives"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewBuilding(t *testing.T) {
	expectedPoints := []*Node{NewNode(1, 2, 3), NewNode(4, 5, 6)}
	b := NewBuilding(0, expectedPoints)

	require.Equal(t, b.Id(), Id(0))
	require.Equal(t, expectedPoints, b.points)
	require.Equal(t, primitives.NewRectangle(2, 3, 5, 6), b.Bounds())
}

func TestBuilding_Connections(t *testing.T) {
	b := NewBuilding(0, nil)
	b.connections = []*Road{NewRoad(1, nil)}

	require.Equal(t, b.Connections(), []*Road{NewRoad(1, nil)})
}
