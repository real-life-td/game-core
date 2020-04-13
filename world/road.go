package world

type Road struct {
	id uint64
	node1, node2 *Node
}

func NewRoad(id uint64, node1, node2 *Node) *Road {
	road := new(Road)
	road.id = id
	road.node1 = node1
	road.node2 = node2
	return road
}
