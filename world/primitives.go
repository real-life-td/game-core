package world

type Node struct {
	id   Id
	x, y int
}

func NewNode(id Id, x, y int) *Node {
	node := new(Node)
	node.id = id
	node.x = x
	node.y = y
	return node
}

func (n *Node) Id() Id {
	return n.id
}

func (n *Node) X() int {
	return n.x
}

func (n *Node) Y() int {
	return n.y
}
