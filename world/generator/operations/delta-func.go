package operations

import (
	"errors"
	. "github.com/dave/jennifer/jen"
)

func writeDeltaFunc(file *File, structName string, deltaOperations []*operation) {
	if len(deltaOperations) == 0 {
		return
	}

	receiverId := receiverId(structName)

	operationCode := make([]Code, 0)

	deltaStructName := deltaStructName(structName)
	operationCode = append(operationCode, Id("delta").Op(":=").New(Id(deltaStructName)))

	for _, operation := range deltaOperations {
		var actionCode []Code

		switch operation.action {
		case addAction:
			actionCode = deltaAddAction(operation, receiverId)
		case removeAction:
			actionCode = deltaRemoveAction(operation, receiverId)
		case setAction:
			actionCode = deltaSetAction(operation, receiverId)
		case putAction:
			actionCode = deltaPutAction(operation, receiverId)
		case putMultipleAction:
			actionCode = deltaPutMultipleAction(operation, receiverId)
		case deleteAction:
			actionCode = deltaDeleteAction(operation, receiverId)
		}

		operationCode = append(operationCode, actionCode...)
	}

	operationCode = append(operationCode, Return(Id("delta")))

	operationStruct := Id("o").Op("*").Id(operationStructName(deltaStage, structName))
	file.Func().Params(Id(receiverId).Op("*").Id(structName)).Id("Operation").Params(operationStruct).Op("*").Id(deltaStructName).Block(operationCode...)
}

func deltaAddAction(operation *operation, receiverName string) []Code {
	if !operation.fieldType.IsArray {
		panic(errors.New("delta-func: Cannot create add function for non-array type"))
	}

	fieldName := operationFieldName(operation)
	structField := Id(receiverName).Dot(operation.field)
	deltaField := Id("delta").Dot(deltaArrayAddFieldName(operation.field))
	deltaArrayLengthField := Id("delta").Dot(deltaArrayLengthFieldName(operation.field))

	assignInitialArrayLength := If(deltaArrayLengthField.Clone().Op("==").Nil()).Block(
		Id("length").Op(":=").Len(structField.Clone()),
		deltaArrayLengthField.Clone().Op("=").Op("&").Id("length"))

	ifNotNilOrEmpty := If(Len(Id("o").Dot(fieldName)).Op("!=").Lit(0)).Block(
		If(structField.Clone().Op("==").Nil()).Block(
			structField.Clone().Op("=").Make(Id(operation.fieldType.Value), Lit(0))),
		structField.Clone().Op("=").Append(structField.Clone(), Id("o").Dot(fieldName).Op("...")),
		If(deltaField.Clone().Op("==").Nil()).Block(
			deltaField.Clone().Op("=").Make(Id(operation.fieldType.Value), Lit(0))),
		deltaField.Clone().Op("=").Append(deltaField.Clone(), Id("o").Dot(fieldName).Op("...")))

	return []Code{If(Id("o").Dot(fieldName).Op("!=").Nil()).Block(assignInitialArrayLength, ifNotNilOrEmpty)}
}

func deltaRemoveAction(operation *operation, receiverName string) []Code {
	if !operation.fieldType.IsArray {
		panic(errors.New("delta-func: Cannot create remove function for non-array type"))
	}

	fieldName := operationFieldName(operation)
	structField := Id(receiverName).Dot(operation.field)
	deltaField := Id("delta").Dot(deltaArrayRemoveFieldName(operation.field))
	deltaArrayLengthField := Id("delta").Dot(deltaArrayLengthFieldName(operation.field))

	ifNotNilOrEmpty := If(structField.Clone().Op("!=").Nil()).Block(
		For(List(Id("_"), Id("toRemove")).Op(":=").Range().Id("o").Dot(fieldName)).Block(
			Id("indexOf").Op(":=").Lit(-1),
			For(List(Id("i"), Id("elm")).Op(":=").Range().Add(structField.Clone())).Block(
				If(Id("elm").Op("==").Id("toRemove")).Block(
					Id("indexOf").Op("=").Id("i"),
					Break())),
			If(Id("indexOf").Op("!=").Lit(-1)).Block(
				structField.Clone().Op("=").Append(structField.Clone().Index(Op(":").Id("indexOf")), structField.Clone().Index(Id("indexOf").Op("+").Lit(1).Op(":")).Op("...")),
				If(deltaField.Clone().Op("==").Nil()).Block(
					deltaField.Clone().Op("=").Make(Index().Int(), Lit(0))),
				deltaField.Clone().Op("=").Append(deltaField.Clone(), Id("indexOf")))))

	assignInitialArrayLength := If(deltaArrayLengthField.Clone().Op("==").Nil()).Block(
		Id("length").Op(":=").Len(structField.Clone()),
		deltaArrayLengthField.Clone().Op("=").Op("&").Id("length"))

	return []Code{If(Id("o").Dot(fieldName).Op("!=").Nil()).Block(ifNotNilOrEmpty, assignInitialArrayLength)}
}

func deltaSetAction(operation *operation, receiverName string) []Code {
	fieldName := operationFieldName(operation)
	structField := Id(receiverName).Dot(operation.field)
	operationField := valueReference(fieldName, operation.fieldType)
	deltaField := Id("delta").Dot(deltaNormalFieldName(operation.field))

	if operation.fieldType.Nillable {
		return []Code{
			If(Id("o").Dot(fieldName).Op("!=").Nil()).Block(
				structField.Clone().Op("=").Add(operationField),
				deltaField.Clone().Op("=").Add(Id("o").Dot(fieldName)))}
	} else {
		return []Code{
			If(Id("o").Dot(fieldName).Op("!=").Nil()).Block(
				structField.Clone().Op("=").Add(operationField),
				Id("valueCopy").Op(":=").Add(structField.Clone()),
				deltaField.Clone().Op("=").Op("&").Id("valueCopy"))}
	}

}
func deltaPutAction(operation *operation, receiverName string) []Code {
	if !operation.fieldType.IsMap {
		panic(errors.New("delta-func: Cannot create put function for non-map type"))
	}

	keyFieldName, valueFieldName := operationMapFieldNames(operation)
	structField := Id(receiverName).Dot(operation.field)
	keyFieldValue := valueReference(keyFieldName, *operation.fieldType.MapKey)
	valueFieldValue := valueReference(valueFieldName, *operation.fieldType.MapValue)
	deltaField := Id("delta").Dot(deltaMapNewFieldName(operation.field))

	return []Code{If(Id("o").Dot(keyFieldName).Op("!=").Nil()).Block(
		structField.Index(keyFieldValue.Clone()).Op("=").Add(valueFieldValue.Clone()),
		If(deltaField.Clone().Op("==").Nil()).Block(
			deltaField.Clone().Op("=").Make(Id(operation.fieldType.Value))),
		deltaField.Clone().Index(keyFieldValue.Clone()).Op("=").Add(valueFieldValue.Clone()))}
}

func deltaPutMultipleAction(operation *operation, receiverName string) []Code {
	if !operation.fieldType.IsMap {
		panic(errors.New("delta-func: Cannot create put_multiple function for non-map type"))
	}

	fieldName := operationFieldName(operation)
	structField := Id(receiverName).Dot(operation.field)
	deltaField := Id("delta").Dot(deltaMapNewFieldName(operation.field))

	return []Code{If(Id("o").Dot(fieldName).Op("!=").Nil()).Block(
		For(List(Id("key"), Id("value")).Op(":=").Range().Id("o").Dot(fieldName)).Block(
			structField.Index(Id("key")).Op("=").Add(Id("value")),
			If(deltaField.Clone().Op("==").Nil()).Block(
				deltaField.Clone().Op("=").Make(Id(operation.fieldType.Value))),
			deltaField.Clone().Index(Id("key")).Op("=").Add(Id("value").Clone())))}
}

func deltaDeleteAction(operation *operation, receiverName string) []Code {
	if !operation.fieldType.IsMap {
		panic(errors.New("delta-func: Cannot create put function for non-map type"))
	}

	fieldName := operationFieldName(operation)
	structField := Id(receiverName).Dot(operation.field)
	deltaField := Id("delta").Dot(deltaMapDeleteFieldName(operation.field))

	return []Code{If(Id("o").Dot(fieldName).Op("!=").Nil()).Block(
		For(List(Id("_"), Id("toDelete")).Op(":=").Range().Id("o").Dot(fieldName)).Block(
			Id("delete").Call(structField, Id("toDelete")),
			If(deltaField.Clone().Op("==").Nil()).Block(
				deltaField.Clone().Op("=").Make(Index().Id(operation.fieldType.MapKey.Value), Lit(0))),
			deltaField.Clone().Op("=").Append(deltaField.Clone(), Id("toDelete"))))}
}
