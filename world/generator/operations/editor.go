package operations

import (
	. "github.com/dave/jennifer/jen"
	"strings"
)

var actionToEditorName = map[action]string {
	setAction: "Set",
	removeAction: "Remove",
	addAction: "Add",
	deleteAction: "Delete",
	putAction: "Put",
	putMultipleAction: "PutMultiple",
}

var actionToEditorParamName = map[action]string {
	setAction: "new",
	removeAction: "toRemove",
	addAction: "toAdd",
	deleteAction: "toDelete",
	putMultipleAction: "toPut",
}

func writeEditors(file *File, structName string, stage stage, operations []*operation) {
	operationStructName := operationStructName(stage, structName)

	file.Func().Id(editorCreateName(structName, stage)).Params().Op("*").Id(operationStructName).Block(
		Return().Op("&").Id(operationStructName).Block())

	fields := groupByField(operations)
	for fieldName, operations := range fields {
		editorName := editorTypeName(operationStructName, fieldName)
		editorMethodName := editorFieldMethodName(fieldName)

		file.Type().Id(editorName).Id(operationStructName)
		file.Func().Params(Id("op").Op("*").Id(operationStructName)).Id(editorMethodName).Params().Params(Op("*").Id(editorName)).Block(
			Return(Parens(Op("*").Id(editorName)).Parens(Id("op"))))

		editorAssignment := func(fieldName, paramName string, dereferencePoint bool) *Statement {
			if dereferencePoint {
				return Id("op").Dot(fieldName).Op("=").Op("&").Id(paramName)
			} else {
				return Id("op").Dot(fieldName).Op("=").Id(paramName)
			}
		}

		for _, operation := range operations {
			var assignmentCode []*Statement
			switch operation.action {
			case putAction:
				keyFieldName, valueFieldName := operationMapFieldNames(operation)
				assignmentCode = append(assignmentCode, editorAssignment(keyFieldName, "key", !operation.fieldType.MapKey.Nillable))
				assignmentCode = append(assignmentCode, editorAssignment(valueFieldName, "value", !operation.fieldType.MapValue.Nillable))
			default:
				fieldName := operationFieldName(operation)
				assignmentCode = append(assignmentCode, editorAssignment(fieldName, actionToEditorParamName[operation.action], !operation.fieldType.Nillable))
			}

			receiver := Id("op").Op("*").Id(editorName)
			file.Func().Params(receiver).Id(actionToEditorName[operation.action]).Params(editorParams(operation)...).Op("*").Id(operationStructName).Block(
				lines(assignmentCode...),
				Return(Parens(Op("*").Id(operationStructName)).Parens(Id("op"))))
		}
	}
}

func editorCreateName(structName string, stage stage) string {
	return structName + stagePrefix[stage] + "Operation"
}

// The name for the type that will contain methods to modify the given field on the given structure
func editorTypeName(operationStructName string, fieldName string) string {
	return operationStructName + strings.Title(fieldName)
}

// The name for the method that will return an editor type for the given field
func editorFieldMethodName(fieldName string) string {
	return strings.Title(fieldName)
}

// Gets an array of Jennifer code representing parameters for a given operations editor method
func editorParams(operation *operation) []Code {
	switch operation.action {
	case putAction:
		keyType, valueType := mapTypes(operation.fieldType)
		return []Code{Id("key").Id(keyType), Id("value").Id(valueType)}
	case deleteAction:
		return []Code{Id(actionToEditorParamName[operation.action]).Index().Id(operation.fieldType.MapKey.Value)}
	default:
		return []Code{Id(actionToEditorParamName[operation.action]).Id(operation.fieldType.Value)}
	}
}

func groupByField(operations []*operation) map[string][]*operation {
	groups := make(map[string][]*operation)
	for _, op := range operations {
		group, ok := groups[op.field]
		if !ok {
			group = make([]*operation, 0)
		}

		groups[op.field] = append(group, op)
	}

	return groups
}
