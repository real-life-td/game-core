//go:generate go run generator/generator.go -- building.go

package world

import (
	"github.com/real-life-td/math/constants"
	"github.com/real-life-td/math/primitives"
)

type Connection struct {
	road            *Road
	distance        float64
	pointOnBuilding *primitives.Point
}

func NewConnection(road *Road, distance float64, pointOnBuilding *primitives.Point) *Connection {
	c := new(Connection)
	c.road = road
	c.distance = distance
	c.pointOnBuilding = pointOnBuilding
	return c
}

func (c *Connection) Road() *Road {
	return c.road
}

func (c *Connection) Distance() float64 {
	return c.distance
}

func (c *Connection) PointOnBuilding() *primitives.Point {
	return c.pointOnBuilding
}

type Building struct {
	id          Id // GEN: INIT_STAGE(SET);GAME_STAGE(SET)
	points      []*Node
	bounds      *primitives.Rectangle
	connections []*Connection // GEN: INIT_STAGE(SET, ADD, REMOVE);GAME_STAGE(ADD,REMOVE,SET)
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

func (b *Building) Connections() []*Connection {
	return b.connections
}
