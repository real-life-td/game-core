package operations

import (
	"fmt"
	. "github.com/dave/jennifer/jen"
	"github.com/real-life-td/game-core/world/generator/parsing"
	"strings"
)

var stagePrefix = map[stage]string{
	initStage: "Init",
	gameStage: "",
}

var actionFieldPrefix = map[action]string{
	setAction:    "New",
	addAction:    "Additional",
	removeAction: "ToRemove",
	putMultipleAction: "PutMultiple",
	deleteAction: "Delete",
}

// Fields in delta structures can be shared by different actions. The delta type helps identify when
// that should happen and when new fields need to be generated.
type deltaType int
const (
	normal deltaType = iota
	arrayRemoved
	arrayAdded
	mapDelete
	mapNew
)

var actionDeltaType = map[action]deltaType{
	setAction:    normal,
	addAction:    arrayAdded,
	removeAction: arrayRemoved,
	putAction: mapNew,
	putMultipleAction: mapNew,
	deleteAction: mapDelete,
}

type deltaField struct {
	deltaType deltaType
	fieldName string
}

type typedDeltaField struct {
	deltaField
	fieldType parsing.GoType
}

func writeOperationStructs(file *File, structName string, structOperations StageOperations) {
	neededDeltaFields := make([]typedDeltaField, 0)
	// Keeps track of what field name + delta type combos are already in neededDeltaFields
	resistedDeltaFields := make(map[deltaField]bool)

	for stage := initStage; stage < numStages; stage++ {
		operations := structOperations.get(stage)
		if len(operations) == 0 {
			continue
		}

		fields := make([]Code, 0)
		for _, operation := range operations {

			if operation.action == putAction {
				keyName, valueName := operationMapFieldNames(operation)
				keyType, valueType := mapNillableTypes(operation.fieldType)

				fields = append(fields, Id(keyName).Id(keyType))
				fields = append(fields, Id(valueName).Id(valueType))
			} else if operation.action == deleteAction {
				fields = append(fields, Id(operationFieldName(operation)).Index().Id(operation.fieldType.MapKey.Value))
			} else {
				fields = append(fields, Id(operationFieldName(operation)).Id(nillableType(operation.fieldType)))
			}

			if operation.stage == gameStage {
				typedDeltaField := typedDeltaField{
					deltaField: deltaField{
						deltaType: actionDeltaType[operation.action],
						fieldName: operation.field,
					},
					fieldType: operation.fieldType,
				}

				if !resistedDeltaFields[typedDeltaField.deltaField] {
					resistedDeltaFields[typedDeltaField.deltaField]	= true
					neededDeltaFields = append(neededDeltaFields, typedDeltaField)
				}
			}
		}

		file.Type().Id(operationStructName(stage, structName)).Struct(fields...)
	}

	if len(neededDeltaFields) == 0 {
		return
	}

	deltaFields := make([]Code, 0)
	for _, deltaField := range neededDeltaFields {
		switch deltaField.deltaType {
		case normal:
			newField := Id(strings.Title(deltaField.fieldName)).Id(nillableType(deltaField.fieldType))
			deltaFields = append(deltaFields, newField)
		case arrayRemoved:
			removedField := Id(deltaArrayRemoveFieldName(deltaField.fieldName)).Index().Int()
			deltaFields = append(deltaFields, removedField)
		case arrayAdded:
			addedField := Id(deltaArrayAddFieldName(deltaField.fieldName)).Id(deltaField.fieldType.Value)
			deltaFields = append(deltaFields, addedField)
		case mapDelete:
			deletedField := Id(deltaMapDeleteFieldName(deltaField.fieldName)).Index().Id(deltaField.fieldType.MapKey.Value)
			deltaFields = append(deltaFields, deletedField)
		case mapNew:
			newField := Id(deltaMapNewFieldName(deltaField.fieldName)).Id(deltaField.fieldType.Value)
			deltaFields = append(deltaFields, newField)
		}
	}

	file.Type().Id(deltaStructName(structName)).Struct(deltaFields...)
}

func operationStructName(stage stage, structName string) string {
	return fmt.Sprintf("%s%sOperation", structName, stagePrefix[stage])
}

func deltaStructName(structName string) string {
	return structName + "Delta"
}

func operationFieldName(op *operation) string {
	return fmt.Sprintf("%s%s", actionFieldPrefix[op.action], strings.Title(op.field))
}

func operationMapFieldNames(op *operation) (key, value string) {
	key = fmt.Sprintf("PutKey%s", strings.Title(op.field))
	value = fmt.Sprintf("PutValue%s", strings.Title(op.field))
	return
}

func deltaArrayRemoveFieldName(fieldName string) string {
	return "Removed" + strings.Title(fieldName)
}

func deltaArrayAddFieldName(fieldName string) string {
	return "Added" + strings.Title(fieldName)
}

func deltaMapNewFieldName(fieldName string) string {
	return "New" + strings.Title(fieldName)
}

func deltaMapDeleteFieldName(fieldName string) string {
	return "Deleted" + strings.Title(fieldName)
}

func nillableType(fieldType parsing.GoType) string {
	if fieldType.Nillable {
		return fieldType.Value
	} else {
		return "*" + fieldType.Value
	}
}

func mapNillableTypes(fieldType parsing.GoType) (keyType, valueType string) {
	return nillableType(*fieldType.MapKey), nillableType(*fieldType.MapValue)
}

func receiverId(structName string) string {
	return strings.ToLower(structName[:1])
}

// Returns code that references the value of field of a operation structure.
// If that field has been made into a pointer so that it is nillable then
// this reference will include the code to de-reference the pointer
func valueReference(fieldName string, fieldType parsing.GoType) *Statement {
	if fieldType.Nillable {
		return Id("o").Dot(fieldName)
	} else {
		return Op("*").Id("o").Dot(fieldName)
	}
}