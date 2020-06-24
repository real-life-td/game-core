//go:generate go run ../generator.go -- $GOFILE

package test

type testModel struct {
	int           int           // GEN: init(SET);delta(SET)
	intPointer    *int          // GEN: init(SET);delta(SET)
	intArray      []int         // GEN: init(SET,ADD,REMOVE);delta(SET,ADD,REMOVE)
	intMapInt     map[int]int   // GEN: init(PUT,PUT_MULTIPLE,DELETE);delta(PUT,PUT_MULTIPLE,DELETE)
	intPointerMap map[*int]*int // GEN: init(PUT,PUT_MULTIPLE,DELETE);delta(PUT,PUT_MULTIPLE,DELETE)
}

func intPointer(value int) *int {
	return &value
}

func defaultTestModel() *testModel {
	return &testModel{
		int:           0,
		intPointer:    intPointer(1),
		intArray:      []int{2, 3, 4},
		intMapInt:     map[int]int{5: 6, 7: 8},
		intPointerMap: make(map[*int]*int),
	}
}
