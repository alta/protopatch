package patch

import (
	"go/ast"
	"go/types"
	"log"
	"strings"
)

func (p *Patcher) patchTypeDef(id *ast.Ident, obj types.Object) {
	fieldType, ok := p.fieldTypes[obj]
	if !ok {
		return
	}

	castDecl := func(v *ast.Field) bool {
		switch t := v.Type.(type) {
		case *ast.Ident:
			t.Name = fieldType
			return true
		case *ast.ArrayType:
			v.Type = &ast.Ident{
				Name: fieldType,
			}
			return true
		case *ast.StarExpr:
			t.X = &ast.Ident{
				Name: fieldType,
			}
			return true
		default:
			return false
		}
	}

	// Cast Field definition
	if id.Obj != nil && id.Obj.Decl != nil {
		v, ok := id.Obj.Decl.(*ast.Field)
		if !ok {
			log.Printf("Warning: fieldType declared for non-field object: %v `%s`", obj, fieldType)
			return
		}
		if !castDecl(v) {
			log.Printf("Warning: unsupported fieldType type: %T `%s`", v.Type, fieldType)
		}
		return
	}
	switch obj.Type().(type) {
	// Cast Getter signature
	case *types.Signature:
		parent := p.findParentNode(id)
		n, ok := parent.(*ast.FuncDecl)
		if !ok {
			log.Printf("Warning: unexpected type for getter: %v `%T`", obj, parent)
			break
		}
		if l := len(n.Type.Results.List); l != 1 {
			log.Printf("Warning: unexpected return count for getter: %v `%d`", obj, l)
			return
		}
		if !castDecl(n.Type.Results.List[0]) {
			log.Printf("Warning: unsupported fieldType type: %T `%s`", n.Type.Results.List[0].Type, fieldType)
		}
		return
	}
}

func (p *Patcher) patchTypeUsage(id *ast.Ident, obj types.Object) {
	desiredType, ok := p.fieldTypes[obj]
	if !ok {
		return
	}
	var originalType string
	switch t := obj.Type().(type) {
	case *types.Basic:
		originalType = t.Name()
	case *types.Signature:
		if t.Results().Len() != 1 {
			return
		}
		originalType = t.Results().At(0).Type().String()
	}
	cast := func(as string, expr ast.Expr) ast.Expr {
		return &ast.CallExpr{
			Fun: &ast.Ident{
				Name: as,
			},
			Args: []ast.Expr{expr},
		}
	}

	usageNode := p.findParentNode(id)
	parentNode := p.findParentNode(usageNode)

	switch usage := usageNode.(type) {
	case *ast.SelectorExpr:
		switch parentExpr := parentNode.(type) {
		case *ast.AssignStmt:
			if len(parentExpr.Lhs) != len(parentExpr.Rhs) {
				return
			}
			for i := range parentExpr.Lhs {
				if parentExpr.Lhs[i] == usage {
					parentExpr.Rhs[i] = cast(desiredType, parentExpr.Rhs[i])
					return
				}
			}
			for i := range parentExpr.Rhs {
				if parentExpr.Rhs[i] == usage {
					parentExpr.Rhs[i] = cast(originalType, parentExpr.Rhs[i])
					return
				}
			}
		case *ast.CallExpr:
			for i := range parentExpr.Args {
				if parentExpr.Args[i] == usage {
					parentExpr.Args[i] = cast(originalType, parentExpr.Args[i])
					return
				}
			}
			parent := p.findParentNode(parentExpr)
			assign, isAssign := parent.(*ast.AssignStmt)
			if parentExpr.Fun == usage && isAssign {
				for i := range assign.Rhs {
					if assign.Rhs[i] == parentExpr {
						assign.Rhs[i] = cast(originalType, assign.Rhs[i])
						return
					}
				}
			}
			call, isCall := parent.(*ast.CallExpr)
			if isCall {
				for i := range call.Args {
					if call.Args[i] == parentExpr {
						call.Args[i] = cast(originalType, call.Args[i])
						return
					}
				}
			}
			for i, v := range parentExpr.Args {
				if v == usage {
					parentExpr.Args[i] = cast(originalType, usage)
					return
				}
			}
		case *ast.BinaryExpr:
			if parentExpr.X == usage {
				parentExpr.X = cast(originalType, parentExpr.X)
			}
			if parentExpr.Y == usage {
				parentExpr.Y = cast(originalType, parentExpr.Y)
			}
		}
	case *ast.KeyValueExpr:
		if usage.Key == id {
			usage.Value = cast(desiredType, usage.Value)
			return
		}
		if usage.Value == id {
			usage.Value = cast(originalType, usage.Value)
			return
		}
	}
}

func packageAndName(fqn string) (pkg string, name string, isSlice bool) {
	isSlice = isSliceType(fqn)
	fqn = strings.TrimPrefix(fqn, "[]")
	parts := strings.Split(fqn, ".")
	if len(parts) == 1 {
		return "", fqn, isSlice
	}
	return strings.Join(parts[:len(parts)-1], "."), parts[len(parts)-1], isSlice
}

func isSliceType(typeName string) bool {
	return strings.HasPrefix(strings.TrimSpace(typeName), "[]")
}

