package operations

import (
	"errors"
	. "github.com/dave/jennifer/jen"
)

func writeGameFunc(file *File, structName string, gameOperations []*operation) {
	receiverId := receiverId(structName)

	operationCode := make([]Code, 0)

	deltaStructName := deltaStructName(structName)
	operationCode = append(operationCode, Id("delta").Op(":=").New(Id(deltaStructName)))

	for _, operation := range gameOperations {
		var actionCode []Code

		switch operation.action {
		case addAction:
			actionCode = gameAddAction(operation, receiverId)
		case removeAction:
			actionCode = gameRemoveAction(operation, receiverId)
		case setAction:
			actionCode = gameSetConnections(operation, receiverId)
		}

		fieldName := operationFieldName(operation)
		operationCode = append(operationCode,
			If(Id("o").Dot(fieldName).Op("!=").Nil()).Block(actionCode...))
	}

	operationCode = append(operationCode, Return(Id("delta")))

	operationStruct := Id("o").Op("*").Id(operationStructName(initStage, structName))
	file.Func().Params(Id(receiverId).Op("*").Id(structName)).Id("Operation").Params(operationStruct).Op("*").Id(deltaStructName).Block(operationCode...)
}

func gameAddAction(operation *operation, receiverName string) []Code {
	if !operation.fieldType.IsArray {
		panic(errors.New("game-func: Cannot create add function for non-array type"))
	}

	fieldName := operationFieldName(operation)
	structField := Id(receiverName).Dot(operation.field)

	ifNotNilOrEmpty := If(Len(Id("o").Dot(fieldName)).Op("!=").Lit(0)).Block(
		If(structField.Clone().Op("==").Nil()).Block(
			structField.Clone().Op("=").Make(Id(operation.fieldType.Value), Lit(0))),
		structField.Clone().Op("=").Append(structField.Clone(), Id("o").Dot(fieldName).Op("...")),
		Id("delta").Dot(operation.field).Op("=").Add(structField.Clone()))

	return []Code{ifNotNilOrEmpty}
}

func gameRemoveAction(operation *operation, receiverName string) []Code {
	if !operation.fieldType.IsArray {
		panic(errors.New("game-func: Cannot create remove function for non-array type"))
	}

	fieldName := operationFieldName(operation)
	structField := Id(receiverName).Dot(operation.field)

	ifNotNilOrEmpty := If(structField.Clone().Op("!=").Nil()).Block(
		Id("removedAny").Op(":=").Lit(false),
		For(List(Id("_"), Id("toRemove")).Op(":=").Range().Id("o").Dot(fieldName)).Block(
			Id("indexOf").Op(":=").Lit(-1),
			For(List(Id("i"), Id("elm")).Op(":=").Range().Add(structField.Clone())).Block(
				If(Id("elm").Op("==").Id("toRemove")).Block(
					Id("indexOf").Op("=").Id("i"),
					Break())),
			If(Id("indexOf").Op("!=").Lit(-1)).Block(
				Id("lastIndex").Op(":=").Len(structField.Clone()).Op("-").Lit(1),
				structField.Clone().Index(Id("indexOf")).Op("=").Add(structField.Clone()).Index(Id("lastIndex")),
				structField.Clone().Op("=").Add(structField.Clone()).Index(Op(":").Id("lastIndex")),
				Id("removedAny").Op("=").Lit(true))),
		If(Id("removedAny")).Block(
			Id("delta").Dot(operation.field).Op("=").Add(structField.Clone())))

	return []Code{ifNotNilOrEmpty}
}

func gameSetConnections(operation *operation, receiverName string) []Code {
	fieldName := operationFieldName(operation)
	structField := Id(receiverName).Dot(operation.field)
	deltaField := Id("delta").Dot(operation.field)

	if operation.fieldType.Nillable {
		return []Code{
			structField.Clone().Op("=").Id("o").Dot(fieldName),
			deltaField.Op("=").Add(structField.Clone()),
		}
	} else {
		return []Code{
			structField.Clone().Op("=").Op("*").Id("o").Dot(fieldName),
			deltaField.Op("=").Op("&").Add(structField.Clone()),
		}
	}
}
