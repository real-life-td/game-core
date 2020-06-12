package operations

import (
	"github.com/dave/jennifer/jen"
	"github.com/real-life-td/game-core/world/generator/parsing"
)

func GenerateOperations(file *jen.File, structs []*parsing.GoStruct) {
	for _, s := range structs {
		stageOperations := findStructureOperations(s.Ast)
		if len(stageOperations) != 0 {
			writeOperationStructs(file, s.Name, stageOperations)

			writeInitFunc(file, s.Name, stageOperations.iterable()[initStage])
		}
	}
}