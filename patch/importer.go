package patch

import "go/types"

type importer struct {
	p *Patcher
}

// Import implements the types.Importer interface.
// TODO: move this to a non-exported type.
func (i importer) Import(path string) (*types.Package, error) {
	return i.ImportFrom(path, "", 0)
}

// ImportFrom implements the types.ImporterFrom interface.
// TODO: move this to a non-exported type.
func (i importer) ImportFrom(path, dir string, mode types.ImportMode) (*types.Package, error) {
	pkg := i.p.getPackage(path, "", true)
	return pkg.pkg, nil
}
