package main

import (
	"errors"
	"github.com/real-life-td/game-core/world/generator/operations"
	"github.com/real-life-td/game-core/world/generator/parsing"
	"go/parser"
	"go/token"
	"os"
	"path"
	"strings"

	"github.com/dave/jennifer/jen"
)

func main() {
	// Remove the argument separator character so that this program will work with `go run -- <path>`
	if os.Args[1] == "--" {
		os.Args = append(os.Args[:1], os.Args[2:]...)
	}

	if len(os.Args) != 2 {
		panic(errors.New("invalid number of args: expecting path of file to generate operations for"))
	}

	fSet := token.NewFileSet()
	fAST, err := parser.ParseFile(fSet, os.Args[1], nil, parser.ParseComments)
	if err != nil {
		panic(err)
	}

	file := jen.NewFile("world")

	file.Comment("!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!")
	file.Comment("Code generated by tools/operation-generator.go DO NOT EDIT")
	file.Comment("!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!")
	file.Line()

	structs := parsing.FindStructures(fAST)
	operations.GenerateOperations(file, structs)

	fileName := strings.TrimSuffix(os.Args[1], path.Ext(os.Args[1])) + "-operations.go"
	err = file.Save(fileName)
	if err != nil {
		panic(err)
	}
}
