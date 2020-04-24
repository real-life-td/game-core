package world

type Building struct {
	id     Id
	points []*Node
}

func NewBuilding(id Id, points []*Node) *Building {
	b := new(Building)
	b.id = id
	b.points = points
	return b
}

func (b *Building) Id() Id {
	return b.id
}

func (b *Building) Points() []*Node {
	return b.points
}
