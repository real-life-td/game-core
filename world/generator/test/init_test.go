package test

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestInit_EmptyOperation(t *testing.T) {
	model := defaultTestModel()

	expected := defaultTestModel()
	operation := testModelInitOperationStruct{}
	model.InitOperation(&operation)

	require.Equal(t, expected, model)
}

func TestInit_Int(t *testing.T) {
	model := defaultTestModel()

	newValue := 100
	operation := testModelInitOperationStruct{intSet: &newValue}
	model.InitOperation(&operation)

	require.Equal(t, newValue, model.int)
}

func TestInit_IntPointer(t *testing.T) {
	model := defaultTestModel()

	newValue := 100
	operation := testModelInitOperationStruct{intPointerSet: &newValue}
	model.InitOperation(&operation)

	require.Equal(t, &newValue, model.intPointer)
}

func TestInit_IntArray_Set(t *testing.T) {
	model := defaultTestModel()

	newValue := []int{100, 200, 300}
	operation := testModelInitOperationStruct{intArraySet: newValue}
	model.InitOperation(&operation)

	require.Equal(t, newValue, model.intArray)
}

func TestInit_IntArray_Add(t *testing.T) {
	model := defaultTestModel()

	toAdd := []int{100, 200, 300}
	operation := testModelInitOperationStruct{intArrayAdd: toAdd}
	model.InitOperation(&operation)

	expected := []int{2, 3, 4, 100, 200, 300}
	require.ElementsMatch(t, expected, model.intArray)
}

func TestInit_IntArray_Remove(t *testing.T) {
	model := defaultTestModel()

	toRemove := []int{2, 3}
	operation := testModelInitOperationStruct{intArrayRemove: toRemove}
	model.InitOperation(&operation)

	require.Equal(t, []int{4}, model.intArray)
}

func TestInit_IntMap_Put(t *testing.T) {
	model := defaultTestModel()

	operation := testModelInitOperationStruct{
		intMapIntPutKey: intPointer(100),
		intMapIntPutValue: intPointer(200),
	}
	model.InitOperation(&operation)

	require.Equal(t, map[int]int{5: 6, 7: 8, 100: 200}, model.intMapInt)
}

func TestInit_IntMap_PutMultiple(t *testing.T) {
	model := defaultTestModel()

	newValues := map[int]int{100: 200, 300: 400}
	operation := testModelInitOperationStruct{intMapIntPutMultiple: newValues}
	model.InitOperation(&operation)

	require.Equal(t, map[int]int{5: 6, 7: 8, 100: 200, 300: 400}, model.intMapInt)
}

func TestInit_IntMap_Delete(t *testing.T) {
	model := defaultTestModel()

	operation := testModelInitOperationStruct{intMapIntDelete: []int{5}}
	model.InitOperation(&operation)

	require.Equal(t, map[int]int{7: 8}, model.intMapInt)
}

func TestInit_IntPointerMap_Put(t *testing.T) {
	model := defaultTestModel()

	newKey, newValue := 100, 200
	operation := testModelInitOperationStruct{intPointerMapPutKey: &newKey, intPointerMapPutValue: &newValue}
	model.InitOperation(&operation)

	require.Equal(t, map[*int]*int{&newKey: &newValue}, model.intPointerMap)
}

func TestInit_IntPointerMap_PutMultiple(t *testing.T) {
	model := defaultTestModel()

	newValues := map[*int]*int{
		intPointer(100): intPointer(200),
		intPointer(300): intPointer(400),
	}
	operation := testModelInitOperationStruct{intPointerMapPutMultiple: newValues}
	model.InitOperation(&operation)

	require.Equal(t, newValues, model.intPointerMap)
}

func TestInit_IntPointerMap_Delete(t *testing.T) {
	model := defaultTestModel()

	key1, value1 := 100, 200
	key2, value2 := 300, 400
	model.intPointerMap = map[*int]*int{&key1: &value1, &key2: &value2}

	operation := testModelInitOperationStruct{intPointerMapDelete: []*int{&key1}}
	model.InitOperation(&operation)

	require.Equal(t, map[*int]*int{&key2: &value2}, model.intPointerMap)
}
