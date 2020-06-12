package operations

import (
	"fmt"
	"github.com/dave/jennifer/jen"
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
}

func writeOperationStructs(file *jen.File, structName string, structOperations StageOperations) {
	for stage, operations := range structOperations.iterable() {
		fields := make([]jen.Code, 0)
		for _, operation := range operations {
			fields = append(fields, jen.Id(operationFieldName(operation)).Id(operation.fieldType.Value))
		}

		file.Type().Id(operationStructName(stage, structName)).Struct(fields...)
	}
}

func operationStructName(stage stage, structName string) string {
	return fmt.Sprintf("%s%sOperation", structName, stagePrefix[stage])
}

func operationFieldName(op *operation) string {
	return fmt.Sprintf("%s%s", actionFieldPrefix[op.action], strings.Title(op.field))
}

func receiverId(structName string) string {
	return strings.ToLower(structName[:1])
}