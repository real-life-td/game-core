package operations

import (
	"fmt"
	"github.com/dave/jennifer/jen"
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

func writeOperationStructs(file *jen.File, structName string, structOperations StageOperations) {
	// Map from field name to type for all fields that need to be a part of the delta structure
	neededDeltaFields := make(map[string]parsing.GoType)

	for stage := initStage; stage < numStages; stage++ {
		operations := structOperations.get(stage)
		if len(operations) == 0 {
			continue
		}

		fields := make([]jen.Code, 0)
		for _, operation := range operations {
			fields = append(fields, jen.Id(operationFieldName(operation)).Id(nillableType(operation.fieldType)))

			if operation.stage == gameStage {
				neededDeltaFields[operation.field] = operation.fieldType
			}
		}

		file.Type().Id(operationStructName(stage, structName)).Struct(fields...)
	}

	if len(neededDeltaFields) == 0 {
		return
	}

	deltaFields := make([]jen.Code, 0)
	for fieldName, fieldType := range neededDeltaFields {
		deltaFields = append(deltaFields, jen.Id(fieldName).Id(nillableType(fieldType)))
	}

	file.Type().Id(deltaStructName(structName)).Struct(deltaFields...)
}

func operationStructName(stage stage, structName string) string {
	return fmt.Sprintf("%s%sOperation", structName, stagePrefix[stage])
}

func deltaStructName(structName string) string {
	return fmt.Sprintf("%sDelta", structName)
}

func operationFieldName(op *operation) string {
	return fmt.Sprintf("%s%s", actionFieldPrefix[op.action], strings.Title(op.field))
}

func nillableType(fieldType parsing.GoType) string {
	if fieldType.Nillable {
		return fieldType.Value
	} else {
		return "*" + fieldType.Value
	}
}

func receiverId(structName string) string {
	return strings.ToLower(structName[:1])
}
