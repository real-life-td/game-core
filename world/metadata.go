package world

type Metadata struct {
	width, height          int
	lat1, lon1, lat2, lon2 float64 // latitudes and longitudes bounding the game world from least to greatest
}

func NewMetadata(width, height int, lat1, lon1, lat2, lon2 float64) *Metadata {
	m := new(Metadata)

	m.width = width
	m.height = height

	if lat1 > lat2 {
		m.lat1 = lat2
		m.lat2 = lat1
	} else {
		m.lat1 = lat1
		m.lat2 = lat2
	}

	if lon1 > lon2 {
		m.lon1 = lon2
		m.lon2 = lon1
	} else {
		m.lon1 = lon1
		m.lon2 = lon2
	}

	return m
}

func (m *Metadata) Width() int {
	return m.width
}

func (m *Metadata) Height() int {
	return m.height
}

// The smallest latitude of the game world's bounds
func (m *Metadata) Lat1() float64 {
	return m.lat1
}

// The smallest longitude of the game world's bounds
func (m *Metadata) Lon1() float64 {
	return m.lon1
}

// The greatest latitude of the game world's bounds
func (m *Metadata) Lat2() float64 {
	return m.lat2
}

// The greatest longitude of the game world's bounds
func (m *Metadata) Lon2() float64 {
	return m.lon2
}
