package operations

import (
	"errors"
	. "github.com/dave/jennifer/jen"
)

func writeInitFunc(file *File, structName string, initOperations []*operation) {
	if len(initOperations) == 0 {
		return
	}

	receiverId := receiverId(structName)

	operationCode := make([]Code, 0)
	for _, operation := range initOperations {
		var actionCode []Code

		switch operation.action {
		case addAction:
			actionCode = initAddAction(operation, receiverId)
		case removeAction:
			actionCode = initRemoveAction(operation, receiverId)
		case setAction:
			actionCode = initSetAction(operation, receiverId)
		case putAction:
			actionCode = initPutAction(operation, receiverId)
		case putMultipleAction:
			actionCode = initPutMultipleAction(operation, receiverId)
		case deleteAction:
			actionCode = initDeleteAction(operation, receiverId)
		}

		operationCode = append(operationCode, actionCode...)
	}

	operationStruct := Id("o").Op("*").Id(operationStructName(initStage, structName))
	file.Func().Params(Id(receiverId).Op("*").Id(structName)).Id("InitOperation").Params(operationStruct).Block(operationCode...)
}

func initAddAction(operation *operation, receiverName string) []Code {
	if !operation.fieldType.IsArray {
		panic(errors.New("init-func: Cannot create add function for non-array type"))
	}

	fieldName := operationFieldName(operation)
	structField := Id(receiverName).Dot(operation.field)

	ifNil := If(structField.Clone().Op("==").Nil()).Block(
		structField.Clone().Op("=").Make(Id(operation.fieldType.Value), Lit(0)))

	add := structField.Clone().Op("=").Append(structField.Clone(), Id("o").Dot(fieldName).Op("..."))

	return []Code{If(Id("o").Dot(fieldName).Op("!=").Nil()).Block(ifNil, add)}
}

func initRemoveAction(operation *operation, receiverName string) []Code {
	if !operation.fieldType.IsArray {
		panic(errors.New("init-func: Cannot create remove function for non-array type"))
	}

	fieldName := operationFieldName(operation)
	structField := Id(receiverName).Dot(operation.field)

	ifNotNil := If(structField.Clone().Op("!=").Nil()).Block(
		For(List(Id("_"), Id("toRemove")).Op(":=").Range().Id("o").Dot(fieldName)).Block(
			Id("indexOf").Op(":=").Lit(-1),
			For(List(Id("i"), Id("elm")).Op(":=").Range().Add(structField.Clone())).Block(
				If(Id("elm").Op("==").Id("toRemove")).Block(
					Id("indexOf").Op("=").Id("i"),
					Break())),
			If(Id("indexOf").Op("!=").Lit(-1)).Block(
				Id("lastIndex").Op(":=").Len(structField.Clone()).Op("-").Lit(1),
				structField.Clone().Index(Id("indexOf")).Op("=").Add(structField.Clone()).Index(Id("lastIndex")),
				structField.Clone().Op("=").Add(structField.Clone()).Index(Op(":").Id("lastIndex")))))

	return []Code{If(Id("o").Dot(fieldName).Op("!=").Nil()).Block(ifNotNil)}
}

func initSetAction(operation *operation, receiverName string) []Code {
	fieldName := operationFieldName(operation)
	structField := Id(receiverName).Dot(operation.field)
	operationField := valueReference(fieldName, operation.fieldType)

	return []Code{If(Id("o").Dot(fieldName).Op("!=").Nil()).Block(
		structField.Op("=").Add(operationField))}
}

func initPutAction(operation *operation, receiverName string) []Code {
	if !operation.fieldType.IsMap {
		panic(errors.New("init-func: Cannot create put function for non-map type"))
	}

	keyFieldName, valueFieldName := operationMapFieldNames(operation)
	structField := Id(receiverName).Dot(operation.field)
	keyFieldValue := valueReference(keyFieldName, *operation.fieldType.MapKey)
	valueFieldValue := valueReference(valueFieldName, *operation.fieldType.MapValue)

	return []Code{If(Id("o").Dot(keyFieldName).Op("!=").Nil()).Block(
		structField.Index(keyFieldValue).Op("=").Add(valueFieldValue))}
}

func initPutMultipleAction(operation *operation, receiverName string) []Code {
	if !operation.fieldType.IsMap {
		panic(errors.New("init-func: Cannot create put_multiple function for non-map type"))
	}

	fieldName := operationFieldName(operation)
	structField := Id(receiverName).Dot(operation.field)

	return []Code{If(Id("o").Dot(fieldName).Op("!=").Nil()).Block(
		For(List(Id("key"), Id("value")).Op(":=").Range().Id("o").Dot(fieldName)).Block(
			structField.Index(Id("key")).Op("=").Add(Id("value"))))}
}

func initDeleteAction(operation *operation, receiverName string) []Code {
	if !operation.fieldType.IsMap {
		panic(errors.New("init-func: Cannot create delete function for non-map type"))
	}

	fieldName := operationFieldName(operation)
	structField := Id(receiverName).Dot(operation.field)

	return []Code{If(Id("o").Dot(fieldName).Op("!=").Nil()).Block(
		For(List(Id("_"), Id("toDelete")).Op(":=").Range().Id("o").Dot(fieldName)).Block(
			Id("delete").Call(structField, Id("toDelete"))))}
}
