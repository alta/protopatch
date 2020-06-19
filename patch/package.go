package patch

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"log"
	"strings"

	"google.golang.org/protobuf/compiler/protogen"
)

// Package represents a Go package for patching.
type Package struct {
	pkg         *types.Package
	files       []*ast.File
	filesByName map[string]*ast.File
}

// NewPackage returns an initialized Package.
func NewPackage(path, name string) *Package {
	log.Printf("Go package:\t%s %q", name, path)
	return &Package{
		pkg:         types.NewPackage(path, name),
		filesByName: make(map[string]*ast.File),
	}
}

// File returns the ast.File for filename, if it exists.
func (pkg *Package) File(filename string) *ast.File {
	return pkg.filesByName[filename]
}

// AddFile adds a parsed file to pkg.
func (pkg *Package) AddFile(filename string, f *ast.File) error {
	if _, ok := pkg.filesByName[filename]; ok {
		return fmt.Errorf("package %s: file already added: %s", pkg.pkg.Name(), filename)
	}
	pkg.files = append(pkg.files, f)
	pkg.filesByName[filename] = f
	return nil
}

// Reset resets pkg type-checks.
func (pkg *Package) Reset() {
	pkg.pkg = types.NewPackage(pkg.pkg.Path(), pkg.pkg.Name())
}

// Check type-checks pkg.
func (pkg *Package) Check(importer types.Importer, fset *token.FileSet, info *types.Info) error {
	log.Printf("Type-check:\t%s", pkg.pkg.Path())

	cfg := &types.Config{
		Error: func(err error) {
			// log.Printf("Warning: %v", err)
		},
		Importer: importer,
	}

	checker := types.NewChecker(cfg, fset, pkg.pkg, info)
	_ = checker.Files(pkg.files)
	return nil // TODO: return an actual error?
}

// Find finds id in Package pkg, and any ancestor(s), or nil if the id is not found in pkg.
func (pkg *Package) Find(id protogen.GoIdent) (obj types.Object, ancestors []types.Object) {
	for _, name := range strings.Split(id.GoName, ".") {
		if obj == nil {
			obj = pkg.pkg.Scope().Lookup(name)
		} else {
			ancestors = append(ancestors, obj)
			obj, _, _ = types.LookupFieldOrMethod(obj.Type(), true, obj.Pkg(), name)
		}
		if obj == nil {
			break
		}
	}
	if obj == nil {
		log.Printf("Warning: unable to find declaration %s.%s", id.GoImportPath, id.GoName)
	}
	return
}
