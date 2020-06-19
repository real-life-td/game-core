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

	return []Code{ifNil, add}
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

	return []Code{ifNotNil}
}

func initSetConnections(operation *operation, receiverName string) []Code {
	fieldName := operationFieldName(operation)
	structField := Id(receiverName).Dot(operation.field)

	if operation.fieldType.Nillable {
		return []Code{structField.Op("=").Id("o").Dot(fieldName)}
	} else {
		return []Code{structField.Op("=").Op("*").Id("o").Dot(fieldName)}
	}

}
