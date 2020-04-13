package world

type Road struct {
	id           uint64
	node1, node2 *Node
}

func NewRoad(id uint64, node1, node2 *Node) *Road {
	road := new(Road)
	road.id = id
	road.node1 = node1
	road.node2 = node2
	return road
}

func (r *Road) Id() uint64 {
	return r.id
}

func (r *Road) Node1() *Node {
	return r.node1
}

func (r *Road) Node2() *Node {
	return r.node2
}
