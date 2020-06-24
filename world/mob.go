//go:generate go run generator/generator.go -- $GOFILE

package world

type MobAttribute int

const (
	// IMPORTANT: Do not change the order of elements
	// Add line to test at mob_test.go when adding a new line
	TestMobAttribute MobAttribute = iota
	// IMPORTANT: Leave this at end
	MobAttributeLength
)

func (m MobAttribute) ref() *MobAttribute {
	return &m
}

type Mob struct {
	id Id
	*Node
	attributes map[MobAttribute]interface{} // GEN: INIT_STAGE(PUT,PUT_MULTIPLE,DELETE);GAME_STAGE(PUT, PUT_MULTIPLE,DELETE)
}

func (m *Mob) GetId() Id {
	return m.id
}

func (m *Mob) GetAttribute(attribute MobAttribute) (value interface{}, ok bool) {
	value, ok = m.attributes[attribute]
	return
}
