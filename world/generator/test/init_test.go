package test

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestInit_EmptyOperation(t *testing.T) {
	model := defaultTestModel()

	expected := defaultTestModel()
	operation := testModelInitOperation{}
	model.InitOperation(&operation)

	require.Equal(t, expected, model)
}

func TestInit_Int(t *testing.T) {
	model := defaultTestModel()

	newValue := 100
	operation := testModelInitOperation{NewInt: &newValue}
	model.InitOperation(&operation)

	require.Equal(t, newValue, model.int)
}

func TestInit_IntPointer(t *testing.T) {
	model := defaultTestModel()

	newValue := 100
	operation := testModelInitOperation{NewIntPointer: &newValue}
	model.InitOperation(&operation)

	require.Equal(t, &newValue, model.intPointer)
}

func TestInit_IntArray_Set(t *testing.T) {
	model := defaultTestModel()

	newValue := []int{100, 200, 300}
	operation := testModelInitOperation{NewIntArray: newValue}
	model.InitOperation(&operation)

	require.Equal(t, newValue, model.intArray)
}

func TestInit_IntArray_Add(t *testing.T) {
	model := defaultTestModel()

	toAdd := []int{100, 200, 300}
	operation := testModelInitOperation{AdditionalIntArray: toAdd}
	model.InitOperation(&operation)

	expected := []int{2, 3, 4, 100, 200, 300}
	require.ElementsMatch(t, expected, model.intArray)
}

func TestInit_IntArray_Remove(t *testing.T) {
	model := defaultTestModel()

	toRemove := []int{2, 3}
	operation := testModelInitOperation{ToRemoveIntArray: toRemove}
	model.InitOperation(&operation)

	require.Equal(t, []int{4}, model.intArray)
}
