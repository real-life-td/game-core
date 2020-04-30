package world

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewRoad(t *testing.T) {
	n := NewNode(1, 2, 3)
	r := NewRoad(4, n)

	require.Equal(t, Id(4), r.Id())
	require.Equal(t, n, r.Node)
}

func TestRoad_Connections(t *testing.T) {
	r := NewRoad(1, nil)
	r.connections = []*Road{r}

	require.Equal(t, []*Road{r}, r.Connections())
}