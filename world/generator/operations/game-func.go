package operations

import (
	"errors"
	. "github.com/dave/jennifer/jen"
	"strings"
)

func writeGameFunc(file *File, structName string, gameOperations []*operation) {
	if len(gameOperations) == 0 {
		return
	}

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
			actionCode = gameSetAction(operation, receiverId)
		}

		operationCode = append(operationCode, actionCode...)
	}

	operationCode = append(operationCode, Return(Id("delta")))

	operationStruct := Id("o").Op("*").Id(operationStructName(gameStage, structName))
	file.Func().Params(Id(receiverId).Op("*").Id(structName)).Id("Operation").Params(operationStruct).Op("*").Id(deltaStructName).Block(operationCode...)
}

func gameAddAction(operation *operation, receiverName string) []Code {
	if !operation.fieldType.IsArray {
		panic(errors.New("game-func: Cannot create add function for non-array type"))
	}

	fieldName := operationFieldName(operation)
	structField := Id(receiverName).Dot(operation.field)
	deltaField := Id("delta").Dot(deltaArrayAddFieldName(operation.field))

	ifNotNilOrEmpty := If(Len(Id("o").Dot(fieldName)).Op("!=").Lit(0)).Block(
		If(structField.Clone().Op("==").Nil()).Block(
			structField.Clone().Op("=").Make(Id(operation.fieldType.Value), Lit(0))),
		structField.Clone().Op("=").Append(structField.Clone(), Id("o").Dot(fieldName).Op("...")),
		If(deltaField.Clone().Op("==").Nil()).Block(
			deltaField.Clone().Op("=").Make(Id(operation.fieldType.Value), Lit(0))),
		deltaField.Clone().Op("=").Append(deltaField.Clone(), Id("o").Dot(fieldName).Op("...")))

	return []Code{If(Id("o").Dot(fieldName).Op("!=").Nil()).Block(ifNotNilOrEmpty)}
}

func gameRemoveAction(operation *operation, receiverName string) []Code {
	if !operation.fieldType.IsArray {
		panic(errors.New("game-func: Cannot create remove function for non-array type"))
	}

	fieldName := operationFieldName(operation)
	structField := Id(receiverName).Dot(operation.field)
	deltaField := Id("delta").Dot(deltaArrayRemoveFieldName(operation.field))

	ifNotNilOrEmpty := If(structField.Clone().Op("!=").Nil()).Block(
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
				If(deltaField.Clone().Op("==").Nil()).Block(
					deltaField.Clone().Op("=").Make(Index().Int(), Lit(0))),
				deltaField.Clone().Op("=").Append(deltaField.Clone(), Id("indexOf")))))

	return []Code{If(Id("o").Dot(fieldName).Op("!=").Nil()).Block(ifNotNilOrEmpty)}
}

func gameSetAction(operation *operation, receiverName string) []Code {
	fieldName := operationFieldName(operation)
	structField := Id(receiverName).Dot(operation.field)
	operationField := valueReference(fieldName, operation.fieldType)
	deltaField := Id("delta").Dot(strings.Title(operation.field))

	return []Code{
		If(Id("o").Dot(fieldName).Op("!=").Nil()).Block(
			structField.Clone().Op("=").Add(operationField),
			deltaField.Op("=").Add(Id("o").Dot(fieldName)))}
}
	}
}
