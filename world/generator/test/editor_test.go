package test

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func test(t *testing.T) {

	require.Equal(t, 2, *operation.IntSet)
}