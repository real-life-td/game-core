package test

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGame_EmptyOperation(t *testing.T) {
	model := defaultTestModel()

	expected := defaultTestModel()
	operation := testModelOperationStruct{}
	delta := model.Operation(&operation)

	require.Equal(t, expected, model)
	require.Equal(t, &testModelDelta{}, delta)
}

func TestGame_Int(t *testing.T) {
	model := defaultTestModel()

	newValue := 100
	operation := testModelOperationStruct{intSet: &newValue}
	delta := model.Operation(&operation)

	require.Equal(t, newValue, model.int)
	require.Equal(t, &testModelDelta{IntNew: &newValue}, delta)
}

// Structures with int with set operations defined for them have to return an int pointer in the delta
// in order to make it nillable. This test makes sure that this pointer is unique so that further changes
// to the structure doesn't change the meaning of the delta's output.
func TestGame_IntDeltaPointerUnique(t *testing.T) {
	model := defaultTestModel()

	operation := testModelOperationStruct{intSet: intPointer(100)}
	delta := model.Operation(&operation)
	*operation.intSet = 200

	require.Equal(t, 100, *delta.IntNew)
}

func TestGame_IntPointer(t *testing.T) {
	model := defaultTestModel()

	newValue := 100
	operation := testModelOperationStruct{intPointerSet: &newValue}
	delta := model.Operation(&operation)

	require.Equal(t, &newValue, model.intPointer)
	require.Equal(t, &testModelDelta{IntPointerNew: &newValue}, delta)
}

func TestGame_IntArray_Set(t *testing.T) {
	model := defaultTestModel()

	newValue := []int{100, 200, 300}
	operation := testModelOperationStruct{intArraySet: newValue}
	delta := model.Operation(&operation)

	require.Equal(t, newValue, model.intArray)
	require.Equal(t, &testModelDelta{IntArrayNew: newValue}, delta)
}

func TestGame_IntArray_Add(t *testing.T) {
	model := defaultTestModel()

	toAdd := []int{100, 200, 300}
	operation := testModelOperationStruct{intArrayAdd: toAdd}
	delta := model.Operation(&operation)

	expected := []int{2, 3, 4, 100, 200, 300}
	require.ElementsMatch(t, expected, model.intArray)
	require.Equal(t, toAdd, delta.IntArrayAdded)
}

func TestGame_IntArray_Remove(t *testing.T) {
	model := defaultTestModel()

	toRemove := []int{2, 3}
	operation := testModelOperationStruct{intArrayRemove: toRemove}
	delta := model.Operation(&operation)

	require.Equal(t, []int{4}, model.intArray)
	require.Equal(t, []int{0, 0}, delta.IntArrayRemovedIndices)
	require.Equal(t, 1, *delta._IntArrayArrayLength)
}

func TestGame_IntMap_Put(t *testing.T) {
	model := defaultTestModel()

	operation := testModelOperationStruct{intMapIntPutKey: intPointer(100), intMapIntPutValue: intPointer(200)}
	delta := model.Operation(&operation)

	require.Equal(t, map[int]int{5: 6, 7: 8, 100: 200}, model.intMapInt)
	require.Equal(t, &testModelDelta{IntMapIntAdded: map[int]int{100: 200}}, delta)
}

func TestGame_IntMap_PutMultiple(t *testing.T) {
	model := defaultTestModel()

	newValues := map[int]int{100: 200, 300: 400}
	operation := testModelOperationStruct{intMapIntPutMultiple: newValues}
	delta := model.Operation(&operation)

	require.Equal(t, map[int]int{5: 6, 7: 8, 100: 200, 300: 400}, model.intMapInt)
	require.Equal(t, &testModelDelta{IntMapIntAdded: newValues}, delta)
}

func TestGame_IntMap_Delete(t *testing.T) {
	model := defaultTestModel()

	operation := testModelOperationStruct{intMapIntDelete: []int{5}}
	delta := model.Operation(&operation)

	require.Equal(t, map[int]int{7: 8}, model.intMapInt)
	require.Equal(t, &testModelDelta{IntMapIntDeleted: []int{5}}, delta)
}

func TestGame_IntPointerMap_Put(t *testing.T) {
	model := defaultTestModel()

	newKey, newValue := 100, 200
	operation := testModelOperationStruct{intPointerMapPutKey: &newKey, intPointerMapPutValue: &newValue}
	delta := model.Operation(&operation)

	require.Equal(t, map[*int]*int{&newKey: &newValue}, model.intPointerMap)
	require.Equal(t, &testModelDelta{IntPointerMapAdded: map[*int]*int{&newKey: &newValue}}, delta)
}

func TestGame_IntPointerMap_PutMultiple(t *testing.T) {
	model := defaultTestModel()

	newValues := map[*int]*int{intPointer(100): intPointer(200), intPointer(300): intPointer(400)}
	operation := testModelOperationStruct{intPointerMapPutMultiple: newValues}
	delta := model.Operation(&operation)

	require.Equal(t, newValues, model.intPointerMap)
	require.Equal(t, &testModelDelta{IntPointerMapAdded: newValues}, delta)
}

func TestGame_IntPointerMap_Delete(t *testing.T) {
	model := defaultTestModel()

	key1, value1 := 100, 200
	key2, value2 := 300, 400
	model.intPointerMap = map[*int]*int{&key1: &value1, &key2: &value2}

	operation := testModelOperationStruct{intPointerMapDelete: []*int{&key1}}
	delta := model.Operation(&operation)

	require.Equal(t, map[*int]*int{&key2: &value2}, model.intPointerMap)
	require.Equal(t, &testModelDelta{IntPointerMapDeleted: []*int{&key1}}, delta)
}

// Tests:
// delta1: remove = nil 	add = nil
// delta2: remove = nil		add = nil
func TestGame_Merge_Int_Test0(t *testing.T) {
	model := defaultTestModel()
	model.intArray = []int{1, 2, 3}

	operation1 := testModelOperationStruct{}
	operation2 := testModelOperationStruct{}
	delta1 := model.Operation(&operation1)
	delta2 := model.Operation(&operation2)
	delta1.Merge(delta2)

	require.Nil(t, delta1.IntArrayAdded)
	require.Nil(t, delta1.IntArrayRemovedIndices)
}


// Tests:
// delta1: remove = nil 	add = nil
// delta2: remove = nil		add = not nil
func TestGame_Merge_Int_Test1(t *testing.T) {
	model := defaultTestModel()
	model.intArray = []int{1, 2, 3}

	operation1 := testModelOperationStruct{}
	operation2 := testModelOperationStruct{intArrayRemove: []int{2, 3}}
	delta1 := model.Operation(&operation1)
	delta2 := model.Operation(&operation2)
	delta1.Merge(delta2)

	require.Nil(t, delta1.IntArrayAdded)
	require.Equal(t, []int{1, 1}, delta1.IntArrayRemovedIndices)
}

// Tests:
// delta1: remove = nil 	add = nil
// delta2: remove = not nil	add = nil
func TestGame_Merge_Int_Test2(t *testing.T) {
	model := defaultTestModel()
	model.intArray = []int{1, 2, 3}

	operation1 := testModelOperationStruct{}
	operation2 := testModelOperationStruct{intArrayAdd: []int{4}}
	delta1 := model.Operation(&operation1)
	delta2 := model.Operation(&operation2)
	delta1.Merge(delta2)

	require.Equal(t, []int{4}, delta1.IntArrayAdded)
	require.Nil(t, delta1.IntArrayRemovedIndices)
}

// Tests:
// delta1: remove = nil 	add = nil
// delta2: remove = not nil	add = not nil
func TestGame_Merge_Int_Test3(t *testing.T) {
	model := defaultTestModel()
	model.intArray = []int{1, 2, 3}

	operation1 := testModelOperationStruct{}
	operation2 := testModelOperationStruct{
		intArrayAdd: []int{4, 5},
		intArrayRemove: []int{2, 3},
	}
	delta1 := model.Operation(&operation1)
	delta2 := model.Operation(&operation2)
	delta1.Merge(delta2)

	require.Equal(t, []int{4, 5}, delta1.IntArrayAdded)
	require.Equal(t, []int{1, 1}, delta1.IntArrayRemovedIndices)
}

// Tests:
// delta1: remove = not nil add = nil
// delta2: remove = nil		add = nil
func TestGame_Merge_Int_Test4(t *testing.T) {
	model := defaultTestModel()
	model.intArray = []int{1, 2, 3}

	operation1 := testModelOperationStruct{intArrayRemove: []int{2, 3}}
	operation2 := testModelOperationStruct{}
	delta1 := model.Operation(&operation1)
	delta2 := model.Operation(&operation2)
	delta1.Merge(delta2)

	require.Nil(t, delta1.IntArrayAdded)
	require.Equal(t, []int{1, 1}, delta1.IntArrayRemovedIndices)
}

// Tests:
// delta1: remove = not nil add = nil
// delta2: remove = not nil	add = nil
func TestGame_Merge_Int_Test5(t *testing.T) {
	model := defaultTestModel()
	model.intArray = []int{1, 2, 3}

	operation1 := testModelOperationStruct{intArrayRemove: []int{1, 3}}
	operation2 := testModelOperationStruct{intArrayRemove: []int{2}}
	delta1 := model.Operation(&operation1)
	delta2 := model.Operation(&operation2)
	delta1.Merge(delta2)

	require.Nil(t, delta1.IntArrayAdded)
	require.Equal(t, []int{0, 1, 0}, delta1.IntArrayRemovedIndices)
}

// Tests:
// delta1: remove = not nil add = nil
// delta2: remove = nil		add = not nil
func TestGame_Merge_Int_Test6(t *testing.T) {
	model := defaultTestModel()
	model.intArray = []int{1, 2, 3}

	operation1 := testModelOperationStruct{intArrayRemove: []int{3, 2}}
	operation2 := testModelOperationStruct{intArrayAdd: []int{4, 5}}
	delta1 := model.Operation(&operation1)
	delta2 := model.Operation(&operation2)
	delta1.Merge(delta2)

	require.Equal(t, []int{4, 5}, delta1.IntArrayAdded)
	require.Equal(t, []int{2, 1}, delta1.IntArrayRemovedIndices)
}

// Tests:
// delta1: remove = not nil add = nil
// delta2: remove = not nil	add = not nil
func TestGame_Merge_Int_Test7(t *testing.T) {
	model := defaultTestModel()
	model.intArray = []int{1, 2, 3}

	operation1 := testModelOperationStruct{intArrayRemove: []int{3, 2}}
	operation2 := testModelOperationStruct{
		intArrayRemove: []int{1},
		intArrayAdd: []int{4, 5},
	}
	delta1 := model.Operation(&operation1)
	delta2 := model.Operation(&operation2)
	delta1.Merge(delta2)

	require.Equal(t, []int{4, 5}, delta1.IntArrayAdded)
	require.Equal(t, []int{2, 1, 0}, delta1.IntArrayRemovedIndices)
}

// Tests:
// delta1: remove = nil 	add = not nil
// delta2: remove = nil		add = nil
func TestGame_Merge_Int_Test8(t *testing.T) {
	model := defaultTestModel()
	model.intArray = []int{1, 2, 3}

	operation1 := testModelOperationStruct{intArrayAdd: []int{4, 5}}
	operation2 := testModelOperationStruct{}
	delta1 := model.Operation(&operation1)
	delta2 := model.Operation(&operation2)
	delta1.Merge(delta2)

	require.Equal(t, []int{4, 5}, delta1.IntArrayAdded)
	require.Nil(t, delta1.IntArrayRemovedIndices)
}

// Tests:
// delta1: remove = nil 	add = not nil
// delta2: remove = not nil	add = nil
func TestGame_Merge_Int_Test9(t *testing.T) {
	model := defaultTestModel()
	model.intArray = []int{1, 2, 3}

	operation1 := testModelOperationStruct{intArrayAdd: []int{4, 5}}
	operation2 := testModelOperationStruct{intArrayRemove: []int{2, 3, 5}}
	delta1 := model.Operation(&operation1)
	delta2 := model.Operation(&operation2)
	delta1.Merge(delta2)

	require.Equal(t, []int{4}, delta1.IntArrayAdded)
	require.Equal(t, []int{1, 1}, delta1.IntArrayRemovedIndices)
}

// Tests:
// delta1: remove = nil 	add = not nil
// delta2: remove = nil		add = not nil
func TestGame_Merge_Int_Test10(t *testing.T) {
	model := defaultTestModel()
	model.intArray = []int{1, 2, 3}

	operation1 := testModelOperationStruct{intArrayAdd: []int{4, 5}}
	operation2 := testModelOperationStruct{intArrayAdd: []int{6}}
	delta1 := model.Operation(&operation1)
	delta2 := model.Operation(&operation2)
	delta1.Merge(delta2)

	require.Nil(t, delta1.IntArrayRemovedIndices)
	require.Equal(t, []int{4, 5, 6}, delta1.IntArrayAdded)
}

// Tests:
// delta1: remove = nil 	add = not nil
// delta2: remove = not nil	add = not nil
func TestGame_Merge_Int_Test11_1(t *testing.T) {
	model := defaultTestModel()
	model.intArray = []int{1, 2, 3}

	operation1 := testModelOperationStruct{intArrayAdd: []int{4, 5}}
	operation2 := testModelOperationStruct{intArrayRemove: []int{2, 1}, intArrayAdd: []int{6, 7, 8}}
	delta1 := model.Operation(&operation1)
	delta2 := model.Operation(&operation2)
	delta1.Merge(delta2)

	require.Equal(t, []int{1, 0}, delta1.IntArrayRemovedIndices)
	require.Equal(t, []int{4, 5, 6, 7, 8}, delta1.IntArrayAdded)
}

// Tests:
// delta1: remove = nil 	add = not nil
// delta2: remove = not nil	add = not nil
func TestGame_Merge_Int_Test11_2(t *testing.T) {
	model := defaultTestModel()
	model.intArray = []int{1, 2, 3}

	operation1 := testModelOperationStruct{intArrayAdd: []int{4, 5, 6}}
	operation2 := testModelOperationStruct{intArrayRemove: []int{5, 2, 4}, intArrayAdd: []int{7}}
	delta1 := model.Operation(&operation1)
	delta2 := model.Operation(&operation2)
	delta1.Merge(delta2)

	require.Equal(t, []int{1}, delta1.IntArrayRemovedIndices)
	require.Equal(t, []int{6, 7}, delta1.IntArrayAdded)
}

// Tests:
// delta1: remove = not nil add = not nil
// delta2: remove = nil		add = nil
func TestGame_Merge_Int_Test12(t *testing.T) {
	model := defaultTestModel()
	model.intArray = []int{1, 2, 3}

	operation1 := testModelOperationStruct{intArrayRemove: []int{1}, intArrayAdd: []int{4, 5}}
	operation2 := testModelOperationStruct{}
	delta1 := model.Operation(&operation1)
	delta2 := model.Operation(&operation2)
	delta1.Merge(delta2)

	require.Equal(t, []int{0}, delta1.IntArrayRemovedIndices)
	require.Equal(t, []int{4, 5}, delta1.IntArrayAdded)
}

// Tests:
// delta1: remove = not nil add = not nil
// delta2: remove = not nil	add = nil
func TestGame_Merge_Int_Test13_1(t *testing.T) {
	model := defaultTestModel()
	model.intArray = []int{1, 2, 3}

	operation1 := testModelOperationStruct{intArrayRemove: []int{1}, intArrayAdd: []int{4, 5}}
	operation2 := testModelOperationStruct{intArrayRemove: []int{2, 3}}
	delta1 := model.Operation(&operation1)
	delta2 := model.Operation(&operation2)
	delta1.Merge(delta2)

	require.Equal(t, []int{0, 0, 0}, delta1.IntArrayRemovedIndices)
	require.Equal(t, []int{4, 5}, delta1.IntArrayAdded)
}

// Tests:
// delta1: remove = not nil add = not nil
// delta2: remove = not nil	add = nil
func TestGame_Merge_Int_Test13_2(t *testing.T) {
	model := defaultTestModel()
	model.intArray = []int{1, 2, 3}

	operation1 := testModelOperationStruct{intArrayRemove: []int{1}, intArrayAdd: []int{4, 5}}
	operation2 := testModelOperationStruct{intArrayRemove: []int{5, 4}}
	delta1 := model.Operation(&operation1)
	delta2 := model.Operation(&operation2)
	delta1.Merge(delta2)

	require.Equal(t, []int{0}, delta1.IntArrayRemovedIndices)
	require.Equal(t, []int{}, delta1.IntArrayAdded)
}

// Tests:
// delta1: remove = not nil add = not nil
// delta2: remove = nil		add = not nil
func TestGame_Merge_Int_Test14(t *testing.T) {
	model := defaultTestModel()
	model.intArray = []int{1, 2, 3}

	operation1 := testModelOperationStruct{intArrayRemove: []int{2, 1}, intArrayAdd: []int{4, 5}}
	operation2 := testModelOperationStruct{intArrayAdd: []int{6, 7}}
	delta1 := model.Operation(&operation1)
	delta2 := model.Operation(&operation2)
	delta1.Merge(delta2)

	require.Equal(t, []int{1, 0}, delta1.IntArrayRemovedIndices)
	require.Equal(t, []int{4, 5, 6, 7}, delta1.IntArrayAdded)
}

// Tests:
// delta1: remove = not nil add = not nil
// delta2: remove = not nil	add = not nil
func TestGame_Merge_Int_Test15_1(t *testing.T) {
	model := defaultTestModel()
	model.intArray = []int{1, 2, 3}

	operation1 := testModelOperationStruct{intArrayRemove: []int{1, 2}, intArrayAdd: []int{4, 5, 6, 7}}
	operation2 := testModelOperationStruct{intArrayRemove: []int{6, 7}, intArrayAdd: []int{8, 9}}

	delta1 := model.Operation(&operation1)
	delta2 := model.Operation(&operation2)
	delta1.Merge(delta2)

	require.Equal(t, []int{0, 0}, delta1.IntArrayRemovedIndices)
	require.Equal(t, []int{4, 5, 8, 9}, delta1.IntArrayAdded)
}

// Tests:
// delta1: remove = not nil add = not nil
// delta2: remove = not nil	add = not nil
func TestGame_Merge_Int_Test15_2(t *testing.T) {
	model := defaultTestModel()
	model.intArray = []int{1, 2, 3}

	operation1 := testModelOperationStruct{intArrayRemove: []int{1}, intArrayAdd: []int{4}}
	operation2 := testModelOperationStruct{intArrayRemove: []int{3}, intArrayAdd: []int{5}}

	delta1 := model.Operation(&operation1)
	delta2 := model.Operation(&operation2)
	delta1.Merge(delta2)

	require.Equal(t, []int{0, 1}, delta1.IntArrayRemovedIndices)
	require.Equal(t, []int{4, 5}, delta1.IntArrayAdded)
}

// Tests:
// delta1: deleted = nil		new = nil
// delta2: deleted = nil		new = nil
func TestGame_Merge_IntMapInt_Test0(t *testing.T) {
	model := defaultTestModel()
	model.intMapInt = map[int]int{1: 2, 3: 4, 5: 6}

	operation1 := testModelOperationStruct{}
	operation2 := testModelOperationStruct{}

	delta1 := model.Operation(&operation1)
	delta2 := model.Operation(&operation2)
	delta1.Merge(delta2)

	require.Nil(t, delta1.IntMapIntAdded)
	require.Nil(t, delta1.IntMapIntDeleted)
}

// Tests:
// delta1: deleted = nil		new = nil
// delta2: deleted = not nil	new = nil
func TestGame_Merge_IntMapInt_Test1(t *testing.T) {
	model := defaultTestModel()
	model.intMapInt = map[int]int{1: 2, 3: 4, 5: 6}

	operation1 := testModelOperationStruct{}
	operation2 := testModelOperationStruct{intMapIntDelete: []int{5, 3}}

	delta1 := model.Operation(&operation1)
	delta2 := model.Operation(&operation2)
	delta1.Merge(delta2)

	require.Nil(t, delta1.IntMapIntAdded)
	require.Equal(t, []int{5, 3}, delta1.IntMapIntDeleted)
}

// Tests:
// delta1: deleted = nil		new = nil
// delta2: deleted = nil		new = not nil
func TestGame_Merge_IntMapInt_Test2(t *testing.T) {
	model := defaultTestModel()
	model.intMapInt = map[int]int{1: 2, 3: 4, 5: 6}

	operation1 := testModelOperationStruct{}
	operation2 := testModelOperationStruct{intMapIntPutMultiple: map[int]int{7: 8, 9: 10}}

	delta1 := model.Operation(&operation1)
	delta2 := model.Operation(&operation2)
	delta1.Merge(delta2)

	require.Equal(t, map[int]int{7: 8, 9: 10}, delta1.IntMapIntAdded)
	require.Nil(t, delta1.IntMapIntDeleted)
}

// Tests:
// delta1: deleted = nil		new = nil
// delta2: deleted = not nil	new = not nil
func TestGame_Merge_IntMapInt_Test3(t *testing.T) {
	model := defaultTestModel()
	model.intMapInt = map[int]int{1: 2, 3: 4, 5: 6}

	operation1 := testModelOperationStruct{}
	operation2 := testModelOperationStruct{
		intMapIntDelete: []int{5, 3},
		intMapIntPutMultiple: map[int]int{7: 8, 9: 10},
	}

	delta1 := model.Operation(&operation1)
	delta2 := model.Operation(&operation2)
	delta1.Merge(delta2)

	require.Equal(t, map[int]int{7: 8, 9: 10}, delta1.IntMapIntAdded)
	require.Equal(t, []int{5, 3}, delta1.IntMapIntDeleted)
}

// Tests:
// delta1: deleted = not nil	new = nil
// delta2: deleted = nil		new = nil
func TestGame_Merge_IntMapInt_Test4(t *testing.T) {
	model := defaultTestModel()
	model.intMapInt = map[int]int{1: 2, 3: 4, 5: 6}

	operation1 := testModelOperationStruct{intMapIntDelete: []int{1}}
	operation2 := testModelOperationStruct{}

	delta1 := model.Operation(&operation1)
	delta2 := model.Operation(&operation2)
	delta1.Merge(delta2)

	require.Nil(t, delta1.IntMapIntAdded)
	require.Equal(t, []int{1}, delta1.IntMapIntDeleted)
}

// Tests:
// delta1: deleted = not nil	new = nil
// delta2: deleted = not nil	new = nil
func TestGame_Merge_IntMapInt_Test5(t *testing.T) {
	model := defaultTestModel()
	model.intMapInt = map[int]int{1: 2, 3: 4, 5: 6}

	operation1 := testModelOperationStruct{intMapIntDelete: []int{1}}
	operation2 := testModelOperationStruct{intMapIntDelete: []int{3}}

	delta1 := model.Operation(&operation1)
	delta2 := model.Operation(&operation2)
	delta1.Merge(delta2)

	require.Nil(t, delta1.IntMapIntAdded)
	require.Equal(t, []int{1, 3}, delta1.IntMapIntDeleted)
}

// Tests:
// delta1: deleted = not nil	new = nil
// delta2: deleted = nil		new = not nil
func TestGame_Merge_IntMapInt_Test6(t *testing.T) {
	model := defaultTestModel()
	model.intMapInt = map[int]int{1: 2, 3: 4, 5: 6}

	operation1 := testModelOperationStruct{intMapIntDelete: []int{1, 5}}
	operation2 := testModelOperationStruct{intMapIntPutMultiple: map[int]int{7: 8, 9: 10, 11: 12}}

	delta1 := model.Operation(&operation1)
	delta2 := model.Operation(&operation2)
	delta1.Merge(delta2)

	require.Equal(t, map[int]int{7: 8, 9: 10, 11: 12}, delta1.IntMapIntAdded)
	require.Equal(t, []int{1, 5}, delta1.IntMapIntDeleted)
}

// Tests:
// delta1: deleted = not nil	new = nil
// delta2: deleted = not nil	new = not nil
func TestGame_Merge_IntMapInt_Test7(t *testing.T) {
	model := defaultTestModel()
	model.intMapInt = map[int]int{1: 2, 3: 4, 5: 6}

	operation1 := testModelOperationStruct{intMapIntDelete: []int{1}}
	operation2 := testModelOperationStruct{
		intMapIntDelete: []int{5, 3},
		intMapIntPutMultiple: map[int]int{7: 8},
	}

	delta1 := model.Operation(&operation1)
	delta2 := model.Operation(&operation2)
	delta1.Merge(delta2)

	require.Equal(t, map[int]int{7: 8}, delta1.IntMapIntAdded)
	require.Equal(t, []int{1, 5, 3}, delta1.IntMapIntDeleted)
}

// Tests:
// delta1: deleted = nil		new = not nil
// delta2: deleted = nil		new = nil
func TestGame_Merge_IntMapInt_Test8(t *testing.T) {
	model := defaultTestModel()
	model.intMapInt = map[int]int{1: 2, 3: 4, 5: 6}

	operation1 := testModelOperationStruct{intMapIntPutMultiple: map[int]int{7: 8, 9: 10}}
	operation2 := testModelOperationStruct{}

	delta1 := model.Operation(&operation1)
	delta2 := model.Operation(&operation2)
	delta1.Merge(delta2)

	require.Equal(t, map[int]int{7: 8, 9: 10}, delta1.IntMapIntAdded)
	require.Nil(t, delta1.IntMapIntDeleted)
}

// Tests:
// delta1: deleted = nil		new = not nil
// delta2: deleted = not nil	new = nil
func TestGame_Merge_IntMapInt_Test9(t *testing.T) {
	model := defaultTestModel()
	model.intMapInt = map[int]int{1: 2, 3: 4, 5: 6}

	operation1 := testModelOperationStruct{intMapIntPutMultiple: map[int]int{7: 8, 9: 10}}
	operation2 := testModelOperationStruct{intMapIntDelete: []int{1, 9}}

	delta1 := model.Operation(&operation1)
	delta2 := model.Operation(&operation2)
	delta1.Merge(delta2)

	require.Equal(t, map[int]int{7: 8}, delta1.IntMapIntAdded)
	require.Equal(t, []int{1, 9}, delta1.IntMapIntDeleted)
}

// Tests:
// delta1: deleted = nil		new = not nil
// delta2: deleted = nil		new = not nil
func TestGame_Merge_IntMapInt_Test10(t *testing.T) {
	model := defaultTestModel()
	model.intMapInt = map[int]int{1: 2, 3: 4, 5: 6}

	operation1 := testModelOperationStruct{intMapIntPutMultiple: map[int]int{7: 8, 9: 10}}
	operation2 := testModelOperationStruct{intMapIntPutMultiple: map[int]int{11: 12}}

	delta1 := model.Operation(&operation1)
	delta2 := model.Operation(&operation2)
	delta1.Merge(delta2)

	require.Equal(t, map[int]int{7: 8, 9: 10, 11: 12}, delta1.IntMapIntAdded)
	require.Nil(t, delta1.IntMapIntDeleted)
}

// Tests:
// delta1: deleted = nil		new = not nil
// delta2: deleted = not nil	new = not nil
func TestGame_Merge_IntMapInt_Test11(t *testing.T) {
	model := defaultTestModel()
	model.intMapInt = map[int]int{1: 2, 3: 4, 5: 6}

	operation1 := testModelOperationStruct{intMapIntPutMultiple: map[int]int{7: 8, 9: 10}}
	operation2 := testModelOperationStruct{
		intMapIntDelete: []int{9, 1, 3},
		intMapIntPutMultiple: map[int]int{11: 12},
	}

	delta1 := model.Operation(&operation1)
	delta2 := model.Operation(&operation2)
	delta1.Merge(delta2)

	require.Equal(t, map[int]int{7: 8, 11: 12}, delta1.IntMapIntAdded)
	require.Equal(t, []int{9, 1, 3}, delta1.IntMapIntDeleted)
}

// Tests:
// delta1: deleted = not nil	new = not nil
// delta2: deleted = nil		new = nil
func TestGame_Merge_IntMapInt_Test12(t *testing.T) {
	model := defaultTestModel()
	model.intMapInt = map[int]int{1: 2, 3: 4, 5: 6}

	operation1 := testModelOperationStruct{
		intMapIntPutMultiple: map[int]int{7: 8, 9: 10},
		intMapIntDelete: []int{3},
	}
	operation2 := testModelOperationStruct{}

	delta1 := model.Operation(&operation1)
	delta2 := model.Operation(&operation2)
	delta1.Merge(delta2)

	require.Equal(t, map[int]int{7: 8, 9: 10}, delta1.IntMapIntAdded)
	require.Equal(t, []int{3}, delta1.IntMapIntDeleted)
}

// Tests:
// delta1: deleted = not nil	new = not nil
// delta2: deleted = not nil	new = nil
func TestGame_Merge_IntMapInt_Test13(t *testing.T) {
	model := defaultTestModel()
	model.intMapInt = map[int]int{1: 2, 3: 4, 5: 6}

	operation1 := testModelOperationStruct{
		intMapIntPutMultiple: map[int]int{7: 8, 9: 10},
		intMapIntDelete: []int{3},
	}
	operation2 := testModelOperationStruct{intMapIntDelete: []int{1, 9}}

	delta1 := model.Operation(&operation1)
	delta2 := model.Operation(&operation2)
	delta1.Merge(delta2)

	require.Equal(t, map[int]int{7: 8}, delta1.IntMapIntAdded)
	require.Equal(t, []int{3, 1, 9}, delta1.IntMapIntDeleted)
}

// Tests:
// delta1: deleted = not nil	new = not nil
// delta2: deleted = nil		new = not nil
func TestGame_Merge_IntMapInt_Test14(t *testing.T) {
	model := defaultTestModel()
	model.intMapInt = map[int]int{1: 2, 3: 4, 5: 6}

	operation1 := testModelOperationStruct{
		intMapIntPutMultiple: map[int]int{7: 8},
		intMapIntDelete: []int{3, 1},
	}
	operation2 := testModelOperationStruct{intMapIntPutMultiple: map[int]int{9: 10, 11: 12}}

	delta1 := model.Operation(&operation1)
	delta2 := model.Operation(&operation2)
	delta1.Merge(delta2)

	require.Equal(t, map[int]int{7: 8, 9: 10, 11: 12}, delta1.IntMapIntAdded)
	require.Equal(t, []int{3, 1}, delta1.IntMapIntDeleted)
}

// Tests:
// delta1: deleted = not nil	new = not nil
// delta2: deleted = not nil	new = not nil
func TestGame_Merge_IntMapInt_Test15(t *testing.T) {
	model := defaultTestModel()
	model.intMapInt = map[int]int{1: 2, 3: 4, 5: 6}

	operation1 := testModelOperationStruct{
		intMapIntPutMultiple: map[int]int{7: 8, 9: 10},
		intMapIntDelete: []int{3, 1},
	}
	operation2 := testModelOperationStruct{
		intMapIntPutMultiple: map[int]int{11: 12, 13: 14},
		intMapIntDelete: []int{9, 5},
	}

	delta1 := model.Operation(&operation1)
	delta2 := model.Operation(&operation2)
	delta1.Merge(delta2)

	require.Equal(t, map[int]int{7: 8, 11: 12, 13: 14}, delta1.IntMapIntAdded)
	require.Equal(t, []int{3, 1, 9, 5}, delta1.IntMapIntDeleted)
}