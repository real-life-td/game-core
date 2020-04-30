package main

import (
	"bufio"
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"os/exec"
	"strings"
)

var requiredImports = []string{
	"github.com/real-life-td/math/primitives",
}

type stage uint8
const (
	initStage stage = iota
	gameStage
)

var commandToStage = map[string]stage{
	"INIT_STAGE": initStage,
	"GAME_STAGE": gameStage,
}

var stagePrefix = map[stage]string {
	initStage: "Init",
	gameStage: "",
}

type action uint8
const (
	setAction action = iota
)

var commandToAction = map[string]action {
	"SET": setAction,
}

var actionFieldPrefix = map[action]string {
	setAction: "New",
}

type operation struct {
	field, goType string
	stage stage
	action action
}

func main() {
	if len(os.Args) != 2 {
		panic(errors.New("invalid number of args: expecting path of file to generate operations for"))
	}

	fset := token.NewFileSet()
	fast, err := parser.ParseFile(fset, os.Args[1], nil, parser.ParseComments)
	if err != nil {
		panic(err)
	}

	f, err := os.OpenFile("world/operations.go", os.O_WRONLY | os.O_CREATE | os.O_TRUNC, 0777)
	if err != nil {
		panic(err)
	}

	b := bufio.NewWriter(f)
	write := func(s string) {
		_, err := b.WriteString(s)
		if err != nil {
			panic(err)
		}
	}

	write("package world\n")
	write("\n")

	write("import (\n")
	for _, requiredImport := range requiredImports {
		write("\t\"" + requiredImport + "\"\n")
	}
	write(")\n")
	write("\n")

	structures, names := findStructures(fast)
	for i, s := range structures {
		operations := findStructureOperations(s)
		write(operationFunctions(names[i], operations))
	}

	err = b.Flush()
	if err != nil {
		panic(err)
	}

	f.Close()

	// Rather than try to make the output of this generator follow the correct formatting. Just run the go formatter
	// on the generated file.
	err = exec.Command("gofmt","-w", "world/operations.go").Run()
	if err != nil {
		panic(err)
	}
}

func findStructureOperations(structType *ast.StructType) []*operation {
	operations := make([]*operation, 0)

	for _, field := range structType.Fields.List {
		if field.Comment != nil && strings.HasPrefix(field.Comment.Text(), "GEN:"){
			// Remove all spaces (including newlines) and the "GEN:" prefix from the string
			trimmedComment := strings.TrimSpace(field.Comment.Text())
			trimmedComment = strings.ReplaceAll(trimmedComment, " ", "")
			trimmedComment = strings.Replace(trimmedComment, "GEN:", "", 1)

			// Different operations stages will be separated by a semi-colon
			stageStrings := strings.Split(trimmedComment, ";")
			for _, stageString := range stageStrings {
				var stage stage
				var actionsString string

				for potentialCommand, potentialStage := range commandToStage {
					if strings.HasPrefix(stageString, potentialCommand) {
						stage = potentialStage

						// Check that the stage string is followed by an open parentheses
						if stageString[len(potentialCommand) : len(potentialCommand) + 1] != "(" {
							panic(errors.New("stage string: '" + stageString + "' missing open parentheses"))
						}

						// Check that the stage string ends in a closing parentheses
						if stageString[len(stageString) - 1:] != ")" {
							panic(errors.New("stage string: '" + stageString + "' missing closing parentheses"))
						}

						// Get text between the parentheses
						actionsString = stageString[len(potentialCommand) + 1 : len(stageString) - 1]
						break
					}
				}

				for _, actionString := range strings.Split(actionsString, ",") {
					for potentialActionString, potentialAction := range commandToAction {
						if actionString == potentialActionString {
							// Multiple names can be comma-separated in one go field
							for _, fieldName := range field.Names {
								o := new(operation)
								o.field = fieldName.String()
								o.stage = stage
								o.action = potentialAction
								o.goType = goTypeFromExpr(field.Type)
								operations = append(operations, o)
							}
						}
					}
				}
			}
		}
	}

	return operations
}

// Convert a expression from the AST back into the string representation of the type
func goTypeFromExpr(e ast.Expr) string {
	switch v := e.(type) {
	case *ast.Ident:
		return v.Name
	case *ast.ArrayType:
		return "[]" + goTypeFromExpr(v.Elt)
	case *ast.StarExpr:
		return "*" + goTypeFromExpr(v.X)
	case *ast.SelectorExpr:
		return goTypeFromExpr(v.X) + "." + v.Sel.Name
	default:
		panic(errors.New("unknown type"))
	}
}

func findStructures(fileAST *ast.File) (structures []*ast.StructType, names []string) {
	structures = make([]*ast.StructType, 0)
	names = make([]string, 0)

	for _, d := range fileAST.Decls {
		genDecl, ok := d.(*ast.GenDecl)
		if ok {
			for _, spec := range genDecl.Specs {
				typeSpec, ok := spec.(*ast.TypeSpec)
				if ok {
					structType, ok := typeSpec.Type.(*ast.StructType)
					if ok {
						structures = append(structures, structType)
						names = append(names, typeSpec.Name.String())
					}
				}
			}
		}
	}
	return
}

func separateStages(operations []*operation) map[stage][]*operation {
	separated := make(map[stage][]*operation)
	for _, op := range operations {
		s, ok := separated[op.stage]
		if ok {
			separated[op.stage] = append(s, op)
		} else {
			separated[op.stage] = []*operation{op}
		}
	}

	return separated
}

// Generates valid go code that implements the operations on a structure. Ends with a newline.
func operationFunctions(structName string, operations []*operation) string {
	var funcString strings.Builder

	recieverName := strings.ToLower(structName[0: 1])

	separated := separateStages(operations)
	for stage, stageOperations := range separated {
		operationStructName := fmt.Sprintf("%s%sOperation", structName, stagePrefix[stage])

		// Generate the structure for the operation
		funcString.WriteString(fmt.Sprintf("type %s struct {\n", operationStructName))

		for _, op := range stageOperations {
			funcString.WriteString(fmt.Sprintf("\t%s%s %s\n", actionFieldPrefix[op.action], strings.Title(op.field), op.goType))
		}

		funcString.WriteString("}\n")
		funcString.WriteString("\n")

		// Generate the function for the operation
		funcString.WriteString(fmt.Sprintf("func (%s *%s) %sOperation(o *%s) {\n", recieverName, structName, stagePrefix[stage], operationStructName))
		for i, op := range stageOperations {
			operationFieldName := fmt.Sprintf("%s%s", actionFieldPrefix[op.action], strings.Title(op.field))

			funcString.WriteString(fmt.Sprintf("\tif o.%s != nil {\n", operationFieldName))

			switch op.action {
			case setAction:
				funcString.WriteString(fmt.Sprintf("\t\t%s.%s = o.%s\n", recieverName, op.field, operationFieldName))
			}

			funcString.WriteString("\t}\n")
			if i != len(stageOperations) - 1 {
				funcString.WriteString("\n")
			}
		}

		funcString.WriteString("}\n")
		funcString.WriteString("\n")
	}

	return funcString.String()
}