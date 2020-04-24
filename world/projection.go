package world

import (
	"errors"
	"math"
)

const degreeToRad = math.Pi / 180
const radToDegree = 180 / math.Pi

type LatLonToGameFunc func(lat, lon float64) (x, y int)
type GameToLatLonFunc func(x, y int) (lat, lon float64)

// Creates functions for converting into and out of game pixels
func CreateConverters(metadata *Metadata) (to LatLonToGameFunc, from GameToLatLonFunc, err error) {
	if metadata == nil {
		return nil, nil, errors.New("metadata cannot be nil")
	}

	// Define R such that the width of the map is correct
	R := float64(metadata.width) / ((metadata.lon2 * degreeToRad) - (metadata.lon1 * degreeToRad))

	originX, originY := mercator(metadata.Lat1(), metadata.Lon1(), R)
	_, endY := mercator(metadata.Lat2(), metadata.Lon2(), R)

	// Determine the scale factor for height so that it comes out correctly
	scaleFactor := float64(metadata.height) / (endY - originY)

	to = func(lat, lon float64) (x, y int) {
		fx, fy := mercator(lat, lon, R)
		return int(math.Round(fx - originX)), int(math.Round((fy - originY) * scaleFactor))
	}

	from = func(x, y int) (lat, lon float64) {
		lat, lon = invMercator(float64(x)+originX, float64(y)/scaleFactor+originY, R)
		return
	}

	return to, from, nil
}

// Converts from a latitude and longitude into game pixels. For batch conversions the CreateConverters function should
// be used instead
func LatLonToGame(metadata *Metadata, lat, lon float64) (x, y int, err error) {
	to, _, err := CreateConverters(metadata)
	if err != nil {
		return -1, -1, err
	}

	x, y = to(lat, lon)
	return
}

// Converts from game pixels into latitude and longitude. For batch conversions the CreateConverters function should
// be used instead
func GameToLatLon(metadata *Metadata, x, y int) (lat, lon float64, err error) {
	_, from, err := CreateConverters(metadata)
	if err != nil {
		return -1, -1, err
	}

	lat, lon = from(x, y)
	return
}

// https://en.wikipedia.org/wiki/Mercator_projection#Derivation_of_the_Mercator_projection
func mercator(lat, lon, R float64) (x, y float64) {
	return R * lon * degreeToRad, R * math.Log(math.Tan((math.Pi/4)+(lat*degreeToRad/2.0)))
}

// https://en.wikipedia.org/wiki/Mercator_projection#Inverse_transformations
func invMercator(x, y, R float64) (lat, lon float64) {
	return (2*math.Atan(math.Exp(y/R)) - (math.Pi / 2)) * radToDegree, (x / R) * radToDegree
}
