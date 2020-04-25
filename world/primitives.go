package world

import "github.com/real-life-td/math/primitives"

type Node struct {
	id Id
	*primitives.Point
}

func NewNode(id Id, x, y int) *Node {
	node := new(Node)
	node.id = id
	node.Point = primitives.NewPoint(x, y)
	return node
}

func (n *Node) Id() Id {
	return n.id
}
