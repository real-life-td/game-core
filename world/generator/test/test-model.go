//go:generate go run ../generator.go -- $GOFILE

package test

type derived float64

type testModel struct {
	int        int   // GEN: INIT_STAGE(SET);GAME_STAGE(SET)
	intPointer *int  // GEN: INIT_STAGE(SET);GAME_STAGE(SET)
	intArray   []int // GEN: INIT_STAGE(SET,ADD,REMOVE);GAME_STAGE(SET,ADD,REMOVE)
	intMapInt map[int]int // GEN: INIT_STAGE(PUT,PUT_MULTIPLE,DELETE);GAME_STAGE(PUT,PUT_MULTIPLE,DELETE)
	intPointerMap map[*int]*int // GEN: INIT_STAGE(PUT,PUT_MULTIPLE,DELETE);GAME_STAGE(PUT,PUT_MULTIPLE,DELETE)
}

func intPointer(value int) *int {
	return &value
}

func defaultTestModel() *testModel {
	return &testModel{
		int:        0,
		intPointer: intPointer(1),
		intArray:   []int{2, 3, 4},
		intMapInt: 	map[int]int {5: 6, 7: 8},
		intPointerMap: make(map[*int]*int),
	}
}
