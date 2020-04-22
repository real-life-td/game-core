package world

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"math"
	"testing"
)

const delta = 0.0000000001

func TestCreateConverters(t *testing.T) {
	_, _, err := CreateConverters(nil)
	require.Error(t, err, "Nil metadata should create error")

	metadata := NewMetadata(100, 100, 0.0, 0.0, 1.0, 1.0)

	to, from, err := CreateConverters(metadata)
	require.NoError(t, err)

	// Test that converting out of and back into game coordinates results in the same number
	testOutIn := func(x, y int, msg string) {
		lat, lon := from(x, y)
		newX, newY := to(lat, lon)
		require.Equal(t, x, newX, msg)
		require.Equal(t, y, newY, msg)
	}

	for _, x := range [3]int{0, 50, 100} {
		for _, y := range [3]int{0, 50, 100} {
			testOutIn(x, y, fmt.Sprintf("out of and back in game coords x: %d and y: %d", x, y))
		}
	}
}

func TestLatLonToGame(t *testing.T) {
	_, _, err := LatLonToGame(nil, 0.0, 0.0)
	require.Error(t, err, "Nil metadata should create error")

	test := func(metadata *Metadata, lat, lon float64, expectedX, expectedY int) {
		msg := fmt.Sprintf("test with lat: %f lon: %f expecting x: %d y: %d", lat, lon, expectedX, expectedY)

		x, y, err := LatLonToGame(metadata, lat, lon)
		require.NoError(t, err, msg)
		require.Equal(t, expectedX, x, msg)
		require.Equal(t, expectedY, y, msg)
	}

	metadata := NewMetadata(100, 100, 0.0, 0.0, 1.0, 1.0)
	for _, lat := range []float64{-0.25, 0.0, 0.25, 0.50, 0.75, 1.0, 1.25} {
		for _, lon := range []float64{-0.25, 0.0, 0.25, 0.50, 0.75, 1.0, 1.25} {
			// We are close enough to the equator and the game is small enough that this simple conversion works
			x, y := int(math.Round(lon * 100)), int(math.Round(lat * 100))
			test(metadata, lat, lon, x, y)
		}
	}

	metadata = NewMetadata(50, 200, 0.0, 45.0, 1.0, 46.0)
	for _, lat := range []float64{-0.25, 0.0, 0.25, 0.50, 0.75, 1.0, 1.25} {
		for _, lon := range []float64{-0.25, 0.0, 0.25, 0.50, 0.75, 1.0, 1.25} {
			// We are close enough to the equator and the game is small enough that this simple conversion works
			x, y := int(math.Round(lon * 50)), int(math.Round(lat * 200))
			test(metadata, lat, lon + 45.0, x, y)
		}
	}
}

func TestGameToLatLon(t *testing.T) {
	_, _, err := GameToLatLon(nil, 0, 0)
	require.Error(t, err)

	test := func(metadata *Metadata, x, y int, expectedLat, expectedLon float64) {
		msg := fmt.Sprintf("test with x: %d y: %d expecting lat: %f lon: %f", x, y, expectedLat, expectedLon)

		lat, lon, err := GameToLatLon(metadata, x, y)
		require.NoError(t, err)
		require.InDelta(t, expectedLat, lat, delta, msg)
		require.InDelta(t, expectedLon, lon, delta, msg)
	}

	metadata := NewMetadata(100000, 100000, 0.0, 0.0, 0.01, 0.01)
	for _, x := range []int{-25000, 0, 25000, 50000, 75000, 100000} {
		for _, y := range []int{-25000, 0, 25000, 50000, 75000, 100000} {
			// We are close enough to the equator and the game is big enough that this simple conversion works
			lat, lon := float64(y) / 10000000.0, float64(x) / 10000000.0
			test(metadata, x, y, lat, lon)
		}
	}

	metadata = NewMetadata(200000, 400000, 0.0, 45.0, 0.01, 45.01)
	for _, x := range []int{-25000, 0, 25000, 50000, 75000, 100000} {
		for _, y := range []int{-25000, 0, 25000, 50000, 75000, 100000} {
			// We are close enough to the equator and the game is big enough that this simple conversion works
			lat, lon := float64(y) / 10000000.0, float64(x) / 10000000.0
			test(metadata, x * 2, y * 4, lat, lon + 45.0)
		}
	}
}

func TestMercator(t *testing.T) {
	// Number determined using https://www.desmos.com/calculator and formulas from Wikipedia article
	// https://en.wikipedia.org/wiki/Mercator_projection#Derivation_of_the_Mercator_projection
	x, y := mercator(10.0, 10.0, 1.0)
	require.Equal(t, 0.17453292519943295, x)
	require.Equal(t, 0.17542582965181813, y)

	x, y = mercator(60.0, 50.0, 1.0)
	require.Equal(t, 0.8726646259971648, x)
	require.Equal(t, 1.316957896924816, y)
}

func TestInvMercator(t *testing.T) {
	// Number determined using https://www.desmos.com/calculator and formulas from Wikipedia article
	// https://en.wikipedia.org/wiki/Mercator_projection#Inverse_transformations
	lat, lon := invMercator(1.0, 1.0, 1.0)
	require.Equal(t, lat, 49.60493742085468)
	require.Equal(t, lon, 57.29577951308232)

	lat, lon = invMercator(-0.2, -0.5, 1.0)
	require.Equal(t, lat, -27.523808392302705)
	require.Equal(t, lon, -11.459155902616466)
}