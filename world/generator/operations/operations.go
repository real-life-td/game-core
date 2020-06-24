package operations

import (
	"github.com/dave/jennifer/jen"
	"github.com/real-life-td/game-core/world/generator/parsing"
	"sort"
)

func GenerateOperations(file *jen.File, structs []*parsing.GoStruct) {
	for _, s := range structs {
		stageOperations := findStructureOperations(s.Ast)
		if len(stageOperations) != 0 {
			writeOperationStructs(file, s.Name, stageOperations)

			writeInitFunc(file, s.Name, sortByPrecedence(stageOperations.get(initStage)))
			writeEditors(file, s.Name, initStage, stageOperations.get(initStage))
			writeGameFunc(file, s.Name, sortByPrecedence(stageOperations.get(gameStage)))
			writeEditors(file, s.Name, gameStage, stageOperations.get(gameStage))
		}
	}
}

func sortByPrecedence(operations []*operation) []*operation {
	sort.Slice(operations, func(i, j int) bool { return operations[i].action < operations[j].action })
	return operations
}
