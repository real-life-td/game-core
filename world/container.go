package world

type Container struct {
	meta *Metadata
	roads []*Road
	buildings []*Building
}

func NewContainer(meta *Metadata, r []*Road, b []*Building) *Container {
	c := new(Container)
	c.meta = meta
	c.roads = r
	c.buildings = b
	return c
}

func (c *Container) Meta() *Metadata {
	return c.meta
}

func (c *Container) Roads() []*Road {
	return c.roads
}

func (c *Container) Buildings() []*Building {
	return c.buildings
}