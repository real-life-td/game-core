package parsing

import (
	"errors"
	"go/ast"
)

type GoStruct struct {
	Ast *ast.StructType
	Name string
}

type GoType struct {
	Value string
	IsArray bool
	Nillable bool
}

func FindStructures(fileAST *ast.File) []*GoStruct {
	structures := make([]*GoStruct, 0)

	for _, d := range fileAST.Decls {
		genDecl, ok := d.(*ast.GenDecl)
		if !ok {
			continue
		}

		for _, spec := range genDecl.Specs {
			typeSpec, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}

			structType, ok := typeSpec.Type.(*ast.StructType)
			if !ok {
				continue
			}

			structures = append(structures, &GoStruct{structType, typeSpec.Name.String()})
		}
	}

	return structures
}

// Convert a expression from the AST back into its goType representation
func GoTypeFromExpr(e ast.Expr) GoType {
	switch v := e.(type) {
	case *ast.Ident:
		return GoType{v.Name, false, false}
	case *ast.ArrayType:
		return GoType{"[]" + GoTypeFromExpr(v.Elt).Value, true, true}
	case *ast.StarExpr:
		return GoType{"*" + GoTypeFromExpr(v.X).Value, false, true}
	case *ast.SelectorExpr:
		return GoType{GoTypeFromExpr(v.X).Value + "." + v.Sel.Name,false, false}
	default:
		panic(errors.New("unknown type"))
	}
}