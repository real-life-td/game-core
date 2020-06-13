//go:generate go run ../generator.go -- $GOFILE

package test

type derived float64

type testModel struct {
	int int // GEN: INIT_STAGE(SET);GAME_STAGE(SET)
	intPointer *int // GEN: INIT_STAGE(SET);GAME_STAGE(SET)
	intArray []int // GEN: INIT_STAGE(SET,ADD,REMOVE);GAME_STAGE(SET,ADD,REMOVE)
}

func defaultTestModel() *testModel {
	intPointer := 1

	return &testModel{
		int: 0,
		intPointer: &intPointer,
		intArray: []int{2, 3, 4},
	}
}
