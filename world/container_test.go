package world

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewContainer(t *testing.T) {
	metadata, roads, buildings := &Metadata{}, make([]*Road, 0), make([]*Building, 0)
	c := NewContainer(metadata, roads, buildings)
	require.Same(t, metadata, c.Meta())
	require.Equal(t, roads, c.Roads())
	require.Equal(t, buildings, c.Buildings())
}
