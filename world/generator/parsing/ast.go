package parsing

import (
	"errors"
	"go/ast"
)

type GoStruct struct {
	Ast  *ast.StructType
	Name string
}

type GoType struct {
	Value    string
	IsArray  bool
	IsMap    bool
	Nillable bool

	MapKey   *GoType
	MapValue *GoType
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
		return GoType{v.Name, false, false, false, nil, nil}
	case *ast.ArrayType:
		return GoType{"[]" + GoTypeFromExpr(v.Elt).Value, true, false, true, nil, nil}
	case *ast.StarExpr:
		return GoType{"*" + GoTypeFromExpr(v.X).Value, false, false, true, nil, nil}
	case *ast.SelectorExpr:
		return GoType{GoTypeFromExpr(v.X).Value + "." + v.Sel.Name, false, false, false, nil, nil}
	case *ast.MapType:
		keyType := GoTypeFromExpr(v.Key)
		valueType := GoTypeFromExpr(v.Value)
		return GoType{"map[" + keyType.Value + "]" + valueType.Value, false, true, true, &keyType, &valueType}
	case *ast.InterfaceType:
		return GoType{"interface{}", false, false, true, nil, nil}
	default:
		panic(errors.New("unknown type"))
	}
}
