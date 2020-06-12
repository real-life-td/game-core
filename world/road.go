//go:generate go run generator/generator.go -- road.go

package world

type Road struct {
	id Id
	*Node
	connections []*Road // GEN: INIT_STAGE(SET, ADD)
}

func NewRoad(id Id, pos *Node) *Road {
	road := new(Road)
	road.id = id
	road.Node = pos
	return road
}

func (r *Road) Id() Id {
	return r.id
}

func (r *Road) Connections() []*Road {
	return r.connections
}
