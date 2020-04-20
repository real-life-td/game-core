package world

type Container struct {
	roads []*Road
}

func NewContainer(r []*Road) *Container {
	c := new(Container)
	c.roads = r
	return c
}

func (c *Container) Roads() []*Road {
	return c.roads
}
