package operations

import (
	. "github.com/dave/jennifer/jen"
	"github.com/real-life-td/game-core/world/generator/parsing"
	"strings"
)

// Fields in delta structures can be shared by different actions. The delta type helps identify when
// that should happen and when new fields need to be generated.
type deltaType int

const (
	normalDeltaType deltaType = iota
	arrayDeltaType
	mapDeltaType
)

var actionDeltaType = map[action]deltaType{
	setAction:         normalDeltaType,
	addAction:         arrayDeltaType,
	removeAction:      arrayDeltaType,
	putAction:         mapDeltaType,
	putMultipleAction: mapDeltaType,
	deleteAction:      mapDeltaType,
}

type deltaField struct {
	deltaType deltaType
	fieldName string
}

type typedDeltaField struct {
	deltaField
	fieldType parsing.GoType
}

func writeDeltaStructs(file *File, structName string, fields []typedDeltaField) {
	deltaFields := make([]Code, 0)

	for _, deltaField := range fields {
		switch deltaField.deltaType {
		case normalDeltaType:
			newField := Id(deltaNormalFieldName(deltaField.fieldName)).Id(nillableType(deltaField.fieldType))
			deltaFields = append(deltaFields, newField)
		case arrayDeltaType:
			removedField := Id(deltaArrayRemoveFieldName(deltaField.fieldName)).Index().Int()
			addedField := Id(deltaArrayAddFieldName(deltaField.fieldName)).Id(deltaField.fieldType.Value)
			arrayLengthField := Id(deltaArrayLengthFieldName(deltaField.fieldName)).Op("*").Int()
			deltaFields = append(deltaFields, removedField, addedField, arrayLengthField)
		case mapDeltaType:
			deletedField := Id(deltaMapDeleteFieldName(deltaField.fieldName)).Index().Id(deltaField.fieldType.MapKey.Value)
			newField := Id(deltaMapNewFieldName(deltaField.fieldName)).Id(deltaField.fieldType.Value)
			deltaFields = append(deltaFields, deletedField, newField)
		}
	}

	file.Type().Id(deltaStructName(structName)).Struct(deltaFields...)
	writeDeltaMerge(file, structName, fields)
}

func writeDeltaMerge(file *File, structName string, fields []typedDeltaField) {
	mergeCode := make([]Code, 0)

	for _, deltaField := range fields {
		switch deltaField.deltaType {
		case normalDeltaType:
			mergeCode = append(mergeCode, deltaMergeNormal(deltaField)...)
		case arrayDeltaType:
			mergeCode = append(mergeCode, deltaMergeArray(deltaField)...)
		case mapDeltaType:
			mergeCode = append(mergeCode, deltaMergeMap(deltaField)...)
		}
	}

	file.Func().Params(Id("c").Op("*").Id(deltaStructName(structName))).Id("Merge").Params(Id("n").Op("*").Id(deltaStructName(structName))).Block(
		mergeCode...)
}

func deltaMergeNormal(field typedDeltaField) []Code {
	curField := Id("c").Dot(deltaNormalFieldName(field.fieldName))
	newField := Id("n").Dot(deltaNormalFieldName(field.fieldName))

	return []Code{If(newField.Clone().Op("!=").Nil()).Block(
		curField.Op("=").Add(newField.Clone()))}
}

func deltaMergeArray(field typedDeltaField) []Code {
	curAddField := Id("c").Dot(deltaArrayAddFieldName(field.fieldName))
	curRemoveField := Id("c").Dot(deltaArrayRemoveFieldName(field.fieldName))
	curArrayLengthFieldValue := Op("*").Id("c").Dot(deltaArrayLengthFieldName(field.fieldName))
	newAddField := Id("n").Dot(deltaArrayAddFieldName(field.fieldName))
	newRemoveField := Id("n").Dot(deltaArrayRemoveFieldName(field.fieldName))

	curAddFieldIsNil := curAddField.Clone().Op("==").Nil()
	curRemoveFieldIsNil := curRemoveField.Clone().Op("==").Nil()
	newAddFieldIsNil := newAddField.Clone().Op("==").Nil()
	newRemoveFieldIsNil := newRemoveField.Clone().Op("==").Nil()

	copyOverAdd := curAddField.Clone().Op("=").Add(newAddField.Clone())
	appendAdd := curAddField.Clone().Op("=").Append(curAddField.Clone(), newAddField.Clone().Op("..."))
	copyOverRemove := curRemoveField.Clone().Op("=").Add(newRemoveField.Clone())
	appendRemove := curRemoveField.Clone().Op("=").Append(curRemoveField.Clone(), newRemoveField.Clone().Op("..."))
	copyOverBoth := copyOverAdd.Clone().Line().Add(copyOverRemove.Clone())

	initCurAdd := curAddField.Clone().Op("=").Make(Id(field.fieldType.Value), Lit(0))
	initCurRemove := curRemoveField.Clone().Op("=").Make(Index().Int(), Lit(0))

	skip := Id("skips").Op(":=").Make(Index().Int(), Lit(0))
	remove := For(List(Id("_"), Id("toRemove")).Op(":=").Range().Add(newRemoveField.Clone())).Block(
		If(Id("toRemove").Op(">=").Add(curArrayLengthFieldValue)).Block(
			Id("toSkip").Op(":=").Id("toRemove").Op("-").Add(curArrayLengthFieldValue.Clone()),
			Id("skips").Op("=").Append(Id("skips"), Id("toSkip"))).Else().Block(
			// ELSE
			curArrayLengthFieldValue.Clone().Op("--"),
			curRemoveField.Clone().Op("=").Append(curRemoveField.Clone(), Id("toRemove"))))

	// Special case when it's known that the skipMap won't be needed (nothing from curAdd will need to be removed)
	removeNoSkipMap := For(List(Id("_"), Id("toRemove")).Op(":=").Range().Add(newRemoveField.Clone())).Block(
		curArrayLengthFieldValue.Clone().Op("--"),
		curRemoveField.Clone().Op("=").Append(curRemoveField.Clone(), Id("toRemove")))

	removeFromCur := For(List(Id("_"), Id("toSkip")).Op(":=").Range().Id("skips")).Block(
		curAddField.Clone().Op("=").Append(curAddField.Clone().Index(Op(":").Id("toSkip")), curAddField.Clone().Index(Id("toSkip").Op("+").Lit(1).Op(":")).Op("...")))

	add := curAddField.Clone().Op("=").Append(curAddField.Clone(), newAddField.Clone().Op("..."))

	normal := lines(
		skip.Clone(),
		remove.Clone(),
		removeFromCur.Clone(),
		add.Clone())

	// Occurs when nothing is nil but curAdd. Initialize that field as an empty slice and bypass
	// the remove current function
	specialCase1 := lines(
		initCurAdd.Clone(),
		removeNoSkipMap.Clone(),
		add.Clone())

	// Occurs when curAdd and newRemove are not nil. We only need to worry about removing
	specialCase2 := lines(
		skip.Clone(),
		remove.Clone(),
		removeFromCur.Clone())

	// The same as normal except the curRemove slice needs to be initialized
	specialCase3 := lines(
		initCurRemove,
		normal)

	// When there are no new elements to add
	specialCase4 := lines(
		skip.Clone(),
		remove.Clone(),
		removeFromCur.Clone())

	return []Code{branch(curAddFieldIsNil,
		branch(curRemoveFieldIsNil,
			branch(newAddFieldIsNil,
				branch(newRemoveFieldIsNil,
					nil,
					copyOverRemove),
				branch(newRemoveFieldIsNil,
					copyOverAdd,
					copyOverBoth)),
			branch(newAddFieldIsNil,
				branch(newRemoveFieldIsNil,
					nil,
					appendRemove),
				branch(newRemoveFieldIsNil,
					copyOverAdd,
					specialCase1))), // SPECIAL CASE #1
		branch(curRemoveFieldIsNil,
			branch(newAddFieldIsNil,
				branch(newRemoveFieldIsNil,
					nil,
					specialCase2), // SPECIAL CASE #2
				branch(newRemoveFieldIsNil,
					appendAdd,
					specialCase3)), // SPECIAL CASE #3
			branch(newAddFieldIsNil,
				branch(newRemoveFieldIsNil,
					nil,
					specialCase4), // SPECIAL CASE #5
				branch(newRemoveFieldIsNil,
					appendAdd,
					normal))))} // Normal
}

func deltaMergeMap(field typedDeltaField) []Code {
	curNewField := Id("c").Dot(deltaMapNewFieldName(field.fieldName))
	curDeleteField := Id("c").Dot(deltaMapDeleteFieldName(field.fieldName))
	newNewField := Id("n").Dot(deltaMapNewFieldName(field.fieldName))
	newDeleteField := Id("n").Dot(deltaMapDeleteFieldName(field.fieldName))

	curNewFieldIsNil := curNewField.Clone().Op("==").Nil()
	curDeleteFieldIsNil := curDeleteField.Clone().Op("==").Nil()
	newNewFieldIsNil := newNewField.Clone().Op("==").Nil()
	newDeleteFieldIsNil := newDeleteField.Clone().Op("==").Nil()

	copyOverNew := curNewField.Clone().Op("=").Add(newNewField.Clone())
	copyOverDelete := curDeleteField.Clone().Op("=").Add(newDeleteField.Clone())
	copyOverBoth := copyOverNew.Clone().Line().Add(copyOverDelete.Clone())

	appendDelete := curDeleteField.Clone().Op("=").Append(curDeleteField.Clone(), newDeleteField.Clone().Op("..."))
	appendNew := For(List(Id("key"), Id("value")).Op(":=").Range().Add(newNewField)).Block(
		curNewField.Clone().Index(Id("key")).Op("=").Id("value"))

	remove := For(List(Id("_"), Id("toRemove")).Op(":=").Range().Add(newDeleteField.Clone())).Block(
		Delete(curNewField.Clone(), Id("toRemove")),
		curDeleteField.Clone().Op("=").Append(curDeleteField.Clone(), Id("toRemove")))

	// Only removes elements from the current delta and doesn't add the new delete elements
	removeNoAppend := For(List(Id("_"), Id("toRemove")).Op(":=").Range().Add(newDeleteField.Clone())).Block(
		Delete(curNewField.Clone(), Id("toRemove")))

	normal := lines(
		remove.Clone(),
		appendNew.Clone())

	// Occurs when only the current deltas new map is nil
	specialCase1 := lines(
		appendDelete.Clone(),
		copyOverNew.Clone())

	// Occurs when both the current deltas delete and the new deltas new map are nil
	specialCase2 := lines(
		copyOverDelete.Clone(),
		removeNoAppend.Clone())

	// Occurs when only the current deltas delete slice is nil
	specialCase3 := lines(
		removeNoAppend.Clone(),
		copyOverDelete.Clone(),
		appendNew.Clone())

	// Occurs when only the new deltas new map is nil
	specialCase4 := remove.Clone()

	return []Code{branch(curNewFieldIsNil,
		branch(curDeleteFieldIsNil,
			branch(newNewFieldIsNil,
				branch(newDeleteFieldIsNil,
					nil,
					copyOverDelete),
				branch(newDeleteFieldIsNil,
					copyOverNew,
					copyOverBoth)),
			branch(newNewFieldIsNil,
				branch(newDeleteFieldIsNil,
					nil,
					appendDelete),
				branch(newDeleteFieldIsNil,
					copyOverNew,
					specialCase1))), // SPECIAL CASE #1
		branch(curDeleteFieldIsNil,
			branch(newNewFieldIsNil,
				branch(newDeleteFieldIsNil,
					nil,
					specialCase2), // SPECIAL CASE #2
				branch(newDeleteFieldIsNil,
					appendNew,
					specialCase3)), // SPECIAL CASE #3
			branch(newNewFieldIsNil,
				branch(newDeleteFieldIsNil,
					nil,
					specialCase4), // SPECIAL CASE #4
				branch(newDeleteFieldIsNil,
					appendNew,
					normal))))} // Normal
}

func deltaStructName(structName string) string {
	return structName + "Delta"
}

func deltaNormalFieldName(fieldName string) string {
	return strings.Title(fieldName) + "New"
}

func deltaArrayRemoveFieldName(fieldName string) string {
	return strings.Title(fieldName) + "RemovedIndices"
}

func deltaArrayAddFieldName(fieldName string) string {
	return strings.Title(fieldName) + "Added"
}

func deltaArrayLengthFieldName(fieldName string) string {
	return "_" + strings.Title(fieldName) + "ArrayLength"
}

func deltaMapNewFieldName(fieldName string) string {
	return strings.Title(fieldName) + "Added"
}

func deltaMapDeleteFieldName(fieldName string) string {
	return strings.Title(fieldName) + "Deleted"
}

func branch(condition, ifTrue, ifFalse *Statement) *Statement {
	if ifTrue == nil {
		return If(Op("!").Parens(Add(condition.Clone()))).Block(ifFalse.Clone())
	} else if ifFalse == nil {
		return If(condition.Clone()).Block(ifTrue.Clone())
	} else {
		return If(condition.Clone()).Block(ifTrue.Clone()).Else().Block(ifFalse.Clone())
	}
}

func lines(statements... *Statement) *Statement {
	if len(statements) == 0 {
		return Null()
	}

	lines := statements[0]
	for _, statement := range statements[1:] {
		lines = lines.Line().Add(statement)
	}

	return lines
}