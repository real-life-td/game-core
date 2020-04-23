package world

import (
	"github.com/stretchr/testify/require"
	"math"
	"testing"
)

func TestTypeIota(t *testing.T) {
	// Ensure that the order of the type constant iota hasn't changed
	require.Equal(t, Type(0), RoadType)
	require.Equal(t, Type(1), NodeType)
}

func TestBaseIdToId(t *testing.T) {
	_, err := BaseIdToId(0x01FFFFFFFFFFFFFF, RoadType)
	require.Error(t, err)

	_, err = BaseIdToId(0, Type(math.MaxUint8))
	require.Error(t, err)

	id, err := BaseIdToId(1, NodeType)
	require.NoError(t, err)
	require.Equal(t, Id(0x0100000000000001), id)

	id, err = BaseIdToId(0x00FFFFFFFFFFFFFF, RoadType)
	require.NoError(t, err)
	require.Equal(t, Id(0x00FFFFFFFFFFFFFF), id)
}

func TestBaseId(t *testing.T) {
	require.Equal(t, uint64(1), BaseId(Id(0x0000000000000001)))
	require.Equal(t, uint64(4), BaseId(Id(0x0100000000000004)))
}

func TestIdType(t *testing.T) {
	require.Equal(t, Type(0), IdType(Id(0x0000000000000001)))
	require.Equal(t, Type(1), IdType(Id(0x0100000000000004)))
}