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
	require.Equal(t, &testModelDelta{int: &newValue}, delta)
}

func TestGame_IntPointer(t *testing.T) {
	model := defaultTestModel()

	newValue := 100
	operation := testModelOperation{NewIntPointer: &newValue}
	delta := model.Operation(&operation)

	require.Equal(t, &newValue, model.intPointer)
	require.Equal(t, &testModelDelta{intPointer: &newValue}, delta)
}

func TestGame_IntArray_Set(t *testing.T) {
	model := defaultTestModel()

	newValue := []int{100, 200, 300}
	operation := testModelOperation{NewIntArray: newValue}
	delta := model.Operation(&operation)

	require.Equal(t, newValue, model.intArray)
	require.Equal(t, &testModelDelta{intArray: newValue}, delta)
}

func TestGame_IntArray_Add(t *testing.T) {
	model := defaultTestModel()

	toAdd := []int{100, 200, 300}
	operation := testModelOperation{AdditionalIntArray: toAdd}
	delta := model.Operation(&operation)

	expected := []int{2, 3, 4, 100, 200, 300}
	require.ElementsMatch(t, expected, model.intArray)
	require.Equal(t, &testModelDelta{intArray: expected}, delta)
}

func TestGame_IntArray_Remove(t *testing.T) {
	model := defaultTestModel()

	toRemove := []int{2, 3}
	operation := testModelOperation{ToRemoveIntArray: toRemove}
	delta := model.Operation(&operation)

	require.Equal(t, []int{4}, model.intArray)
	require.Equal(t, &testModelDelta{intArray: []int{4}}, delta)
}