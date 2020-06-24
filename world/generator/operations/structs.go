package operations

import (
	"fmt"
	. "github.com/dave/jennifer/jen"
	"github.com/real-life-td/game-core/world/generator/parsing"
	"strings"
)

var stagePrefix = map[stage]string{
	initStage: "Init",
	deltaStage: "",
}

var actionFieldSuffix = map[action]string{
	setAction:         "Set",
	addAction:         "Add",
	removeAction:      "Remove",
	putMultipleAction: "PutMultiple",
	deleteAction:      "Delete",
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
			fields = append(fields, operationFields(operation)...)

			if operation.stage == deltaStage {
				typedDeltaField := typedDeltaField{
					deltaField: deltaField{
						deltaType: actionDeltaType[operation.action],
						fieldName: operation.field,
					},
					fieldType: operation.fieldType,
				}

				if !resistedDeltaFields[typedDeltaField.deltaField] {
					resistedDeltaFields[typedDeltaField.deltaField] = true
					neededDeltaFields = append(neededDeltaFields, typedDeltaField)
				}
			}
		}

		file.Type().Id(operationStructName(stage, structName)).Struct(fields...)
	}

	if len(neededDeltaFields) != 0 {
		writeDeltaStructs(file, structName, neededDeltaFields)
	}
}

func operationFields(operation *operation) []Code {
	fields := make([]Code, 0)

	switch operation.action {
	case putAction:
		keyName, valueName := operationMapFieldNames(operation)

		var keyType, valueType string
		keyType, valueType = mapNillableTypes(operation.fieldType)

		fields = append(fields, Id(keyName).Id(keyType))
		fields = append(fields, Id(valueName).Id(valueType))
	case deleteAction:
		fields = append(fields, Id(operationFieldName(operation)).Index().Id(operation.fieldType.MapKey.Value))
	default:
		fields = append(fields, Id(operationFieldName(operation)).Id(nillableType(operation.fieldType)))
	}

	return fields
}

func operationStructName(stage stage, structName string) string {
	return fmt.Sprintf("%s%sOperationStruct", structName, stagePrefix[stage])
}

func operationFieldName(op *operation) string {
	return fmt.Sprintf("%s%s", op.field, actionFieldSuffix[op.action])
}

func operationMapFieldNames(op *operation) (key, value string) {
	key = fmt.Sprintf("%sPutKey", op.field)
	value = fmt.Sprintf("%sPutValue", op.field)
	return
}

func nillableType(fieldType parsing.GoType) string {
	if fieldType.Nillable {
		return fieldType.Value
	} else {
		return "*" + fieldType.Value
	}
}

func mapTypes(fieldType parsing.GoType) (keyType, valueType string) {
	return fieldType.MapKey.Value, fieldType.MapValue.Value
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
