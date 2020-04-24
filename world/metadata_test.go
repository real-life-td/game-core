package world

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewMetadata(t *testing.T) {
	m := NewMetadata(1, 2, 3.0, 4.0, 5.0, 6.0)

	require.Equal(t, 1, m.Width())
	require.Equal(t, 2, m.Height())
	require.Equal(t, 3.0, m.Lat1())
	require.Equal(t, 4.0, m.Lon1())
	require.Equal(t, 5.0, m.Lat2())
	require.Equal(t, 6.0, m.Lon2())

	// Should sort latitudes and longitudes
	m = NewMetadata(1, 2, 5.0, 6.0, 3.0, 4.0)
	require.Equal(t, 3.0, m.Lat1())
	require.Equal(t, 4.0, m.Lon1())
	require.Equal(t, 5.0, m.Lat2())
	require.Equal(t, 6.0, m.Lon2())
}
