package world

import "errors"

type Id uint64
type Type uint8

const (
	// IMPORTANT: Do not change the order of elements
	// Add line to test at ids_test.go when adding a new line
	RoadType Type = iota
	NodeType
	typeEnd
)

func typeValid(t Type) bool {
	return t < typeEnd
}

// Adds type information. This requires the 8 most significant bits of the basedId to be empty or error will be thrown
func BaseIdToId(baseId uint64, t Type) (id Id, err error) {
	if baseId >> 56 != 0 {
		return 0, errors.New("8 most significant bits of base id must be 0")
	}

	if !typeValid(t) {
		return 0, errors.New("invalid type")
	}

	return Id((uint64(t) << 56) | baseId), nil
}

func IdType(id Id) Type {
	return Type(id >> 56)
}

func BaseId(id Id) uint64 {
	return uint64(id) << 8 >> 8
}