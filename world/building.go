//go:generate go run tools/operation-generator.go -- building.go

package world

import (
	"github.com/real-life-td/math/constants"
	"github.com/real-life-td/math/primitives"
)

type Building struct {
	id          Id
	points      []*Node
	bounds      *primitives.Rectangle
	connections []*Road // GEN: INIT_STAGE(SET, ADD)
}

func NewBuilding(id Id, points []*Node) *Building {
	b := new(Building)
	b.id = id
	b.points = points

	// compute bounds
	minX, minY := constants.MaxInt, constants.MaxInt
	maxX, maxY := constants.MinInt, constants.MinInt
	for _, p := range points {
		if p.X() < minX {
			minX = p.X()
		} else if p.X() > maxX {
			maxX = p.X()
		}

		if p.Y() < minY {
			minY = p.Y()
		} else if p.Y() > maxY {
			maxY = p.Y()
		}
	}

	b.bounds = primitives.NewRectangle(minX, minY, maxX, maxY)

	return b
}

func (b *Building) Id() Id {
	return b.id
}

func (b *Building) Points() []*Node {
	return b.points
}

func (b *Building) Bounds() *primitives.Rectangle {
	return b.bounds
}

func (b *Building) Connections() []*Road {
	return b.connections
}
