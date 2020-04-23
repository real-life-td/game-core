package world

import "errors"

// Id contains 16 bits of type information followed by 48 bits of a base id. This approach has 2 benefits
// 1. Simplifies the process of creating unique ids for objects
//		a. For objects with ids from overpass we can just re-use those
//		b. For other objects we can just use a simple counter
// 2. Compactly adds type information to each object
//
// The reason for the 16/48 split is that if we ever want to consume these components in JavaScript it will allow us
// to safely fit the base id into a 64 bit floating point number (max safe int being 53 bits long)
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

// The 16 most significant bits of baseId must be empty or error will be thrown
func NewId(baseId uint64, t Type) (id Id, err error) {
	if baseId >> 48 != 0 {
		return 0, errors.New("16 most significant bits of base id must be 0")
	}

	if !typeValid(t) {
		return 0, errors.New("invalid type")
	}

	return Id((uint64(t) << 48) | baseId), nil
}

func (id Id) Type() Type {
	return Type(id >> 48)
}

func (id Id) BaseId() uint64 {
	return uint64(id) & 0x0000FFFFFFFFFFFF
}