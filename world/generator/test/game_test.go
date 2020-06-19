package test

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGame_EmptyOperation(t *testing.T) {
	model := defaultTestModel()

	expected := defaultTestModel()
	operation := testModelOperation{}
	delta := model.Operation(&operation)

	require.Equal(t, expected, model)
	require.Equal(t, &testModelDelta{}, delta)
}

func TestGame_Int(t *testing.T) {
	model := defaultTestModel()

	newValue := 100
	operation := testModelOperation{NewInt: &newValue}
	delta := model.Operation(&operation)

	require.Equal(t, newValue, model.int)
	require.Equal(t, &testModelDelta{Int: &newValue}, delta)
}

func TestGame_IntPointer(t *testing.T) {
	model := defaultTestModel()

	newValue := 100
	operation := testModelOperation{NewIntPointer: &newValue}
	delta := model.Operation(&operation)

	require.Equal(t, &newValue, model.intPointer)
	require.Equal(t, &testModelDelta{IntPointer: &newValue}, delta)
}

func TestGame_IntArray_Set(t *testing.T) {
	model := defaultTestModel()

	newValue := []int{100, 200, 300}
	operation := testModelOperation{NewIntArray: newValue}
	delta := model.Operation(&operation)

	require.Equal(t, newValue, model.intArray)
	require.Equal(t, &testModelDelta{IntArray: newValue}, delta)
}

func TestGame_IntArray_Add(t *testing.T) {
	model := defaultTestModel()

	toAdd := []int{100, 200, 300}
	operation := testModelOperation{AdditionalIntArray: toAdd}
	delta := model.Operation(&operation)

	expected := []int{2, 3, 4, 100, 200, 300}
	require.ElementsMatch(t, expected, model.intArray)
	require.Equal(t, &testModelDelta{AddedIntArray: toAdd}, delta)
}

func TestGame_IntArray_Remove(t *testing.T) {
	model := defaultTestModel()

	toRemove := []int{2, 3}
	operation := testModelOperation{ToRemoveIntArray: toRemove}
	delta := model.Operation(&operation)

	require.Equal(t, []int{4}, model.intArray)
	require.Equal(t, &testModelDelta{RemovedIntArray: []int{0, 1}}, delta)
}

func TestGame_IntMap_Put(t *testing.T) {
	model := defaultTestModel()

	operation := testModelOperation{PutKeyIntMapInt: intPointer(100), PutValueIntMapInt: intPointer(200)}
	delta := model.Operation(&operation)

	require.Equal(t, map[int]int {5: 6, 7: 8, 100: 200}, model.intMapInt)
	require.Equal(t, &testModelDelta{NewIntMapInt: map[int]int {100: 200}}, delta)
}

func TestGame_IntMap_PutMultiple(t *testing.T) {
	model := defaultTestModel()

	newValues := map[int]int {100: 200, 300: 400}
	operation := testModelOperation{PutMultipleIntMapInt: newValues}
	delta := model.Operation(&operation)

	require.Equal(t, map[int]int {5: 6, 7: 8, 100: 200, 300: 400}, model.intMapInt)
	require.Equal(t, &testModelDelta{NewIntMapInt: newValues}, delta)
}

func TestGame_IntMap_Delete(t *testing.T) {
	model := defaultTestModel()

	operation := testModelOperation{DeleteIntMapInt: []int{5}}
	delta := model.Operation(&operation)

	require.Equal(t, map[int]int {7: 8}, model.intMapInt)
	require.Equal(t, &testModelDelta{DeletedIntMapInt: []int {5}}, delta)
}

func TestGame_IntPointerMap_Put(t *testing.T) {
	model := defaultTestModel()

	newKey, newValue := 100, 200
	operation := testModelOperation{PutKeyIntPointerMap: &newKey, PutValueIntPointerMap: &newValue}
	delta := model.Operation(&operation)

	require.Equal(t, map[*int]*int {&newKey: &newValue}, model.intPointerMap)
	require.Equal(t, &testModelDelta{NewIntPointerMap: map[*int]*int {&newKey: &newValue}}, delta)
}

func TestGame_IntPointerMap_PutMultiple(t *testing.T) {
	model := defaultTestModel()

	newValues := map[*int]*int {intPointer(100): intPointer(200), intPointer(300): intPointer(400)}
	operation := testModelOperation{PutMultipleIntPointerMap: newValues}
	delta := model.Operation(&operation)

	require.Equal(t, newValues, model.intPointerMap)
	require.Equal(t, &testModelDelta{NewIntPointerMap: newValues}, delta)
}

func TestGame_IntPointerMap_Delete(t *testing.T) {
	model := defaultTestModel()

	key1, value1 := 100, 200
	key2, value2 := 300, 400
	model.intPointerMap = map[*int]*int {&key1: &value1, &key2: &value2}

	operation := testModelOperation{DeleteIntPointerMap: []*int {&key1}}
	delta := model.Operation(&operation)

	require.Equal(t, map[*int]*int {&key2: &value2}, model.intPointerMap)
	require.Equal(t, &testModelDelta{DeletedIntPointerMap: []*int {&key1}}, delta)
}