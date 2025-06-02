package patch

import (
	"go/ast"
	"go/token"
	"go/types"
	"log"
	"strings"
)

func (p *Patcher) patchNoGetters(id *ast.Ident, obj types.Object) {
	var found bool
	for _, v := range p.noGetters {
		if found = v == obj; found {
			break
		}
	}
	if !found {
		return
	}
	named, ok := obj.Type().(*types.Named)
	if !ok {
		return
	}
	pos := make(map[string]token.Pos)
	for i := 0; i < named.NumMethods(); i++ {
		m := named.Method(i)
		if !strings.HasPrefix(m.Name(), "Get") {
			continue
		}
		pos[m.Name()] = m.Pos()
	}
	file := p.fileOf(id)
	for k, v := range pos {
		found := false
		for i, vv := range file.Decls {
			vv, ok := vv.(*ast.FuncDecl)
			if !ok || vv.Recv == nil {
				continue
			}
			if found = v == vv.Name.NamePos; found {
				file.Decls = append(file.Decls[:i], file.Decls[i+1:]...)
				break
			}
		}
		if !found {
			log.Printf("Warning: getter not found: %v `%s`", obj, k)
		}
	}
}
