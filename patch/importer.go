package patch

import "go/types"

type basicImporter struct {
	p *Patcher
}

// Import implements the types.Importer interface.
func (i basicImporter) Import(path string) (*types.Package, error) {
	return i.ImportFrom(path, "", 0)
}

// ImportFrom implements the types.ImporterFrom interface.
func (i basicImporter) ImportFrom(path, dir string, mode types.ImportMode) (*types.Package, error) {
	pkg := i.p.getPackage(path, "", true)
	return pkg.pkg, nil
}
