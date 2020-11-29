package patch

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"go/types"
	"log"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/alta/protopatch/patch/ident"
	"golang.org/x/tools/go/ast/astutil"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/types/pluginpb"
)

// Patcher patches a set of generated Go Protobuf files with additional features:
// - go_name (Field option) overrides the name of a synthesized struct field and getters.
// - go_tags (Field option) lets you add additional struct tags to a field.
// - go_oneof_name (Oneof option) overrides the name of a oneof field, including wrapper types and getters.
// - go_oneof_tags (Oneof option) lets you specify additional struct tags on a oneof field.
// - go_message_name (Message option) overrides the name of the synthesized struct.
// - go_enum_name (Enum option) overrides the name of an enum type.
// - go_value_name (EnumValue option) overrides the name of an enum const.
type Patcher struct {
	gen            *protogen.Plugin
	fset           *token.FileSet
	files          []*ast.File
	filesByName    map[string]*ast.File
	info           *types.Info
	packages       []*Package
	packagesByPath map[string]*Package
	packagesByName map[string]*Package
	renames        map[protogen.GoIdent]string
	typeRenames    map[protogen.GoIdent]string
	valueRenames   map[protogen.GoIdent]string
	fieldRenames   map[protogen.GoIdent]string
	methodRenames  map[protogen.GoIdent]string
	objectRenames  map[types.Object]string
	tags           map[protogen.GoIdent]string
	fieldTags      map[types.Object]string
}

// NewPatcher returns an initialized Patcher for gen.
func NewPatcher(gen *protogen.Plugin) (*Patcher, error) {
	p := &Patcher{
		gen:            gen,
		packagesByPath: make(map[string]*Package),
		packagesByName: make(map[string]*Package),
		renames:        make(map[protogen.GoIdent]string),
		typeRenames:    make(map[protogen.GoIdent]string),
		valueRenames:   make(map[protogen.GoIdent]string),
		fieldRenames:   make(map[protogen.GoIdent]string),
		methodRenames:  make(map[protogen.GoIdent]string),
		objectRenames:  make(map[types.Object]string),
		tags:           make(map[protogen.GoIdent]string),
		fieldTags:      make(map[types.Object]string),
	}
	return p, p.scan()
}

func (p *Patcher) scan() error {
	for _, f := range p.gen.Files {
		p.scanFile(f)
	}
	return nil
}

func (p *Patcher) scanFile(f *protogen.File) {
	log.Printf("\nScan proto:\t%s", f.Desc.Path())

	_ = p.getPackage(string(f.GoImportPath), string(f.GoPackageName), true)

	for _, e := range f.Enums {
		p.scanEnum(e)
	}

	for _, m := range f.Messages {
		p.scanMessage(m, nil)
	}

	for _, e := range f.Extensions {
		p.scanExtension(e)
	}

	// TODO: scan gRPC services
}

func (p *Patcher) scanEnum(e *protogen.Enum) {
	opts := enumOptions(e)
	newName := opts.GetName()
	if newName != "" {
		p.RenameType(e.GoIdent, newName)                                       // Enum type
		p.RenameValue(ident.WithSuffix(e.GoIdent, "_name"), newName+"_name")   // Enum name map
		p.RenameValue(ident.WithSuffix(e.GoIdent, "_value"), newName+"_value") // Enum value map
	}
	stringerName := opts.GetStringerName()
	if stringerName != "" {
		p.RenameMethod(ident.WithChild(e.GoIdent, "String"), stringerName)
	}
	for _, v := range e.Values {
		p.scanEnumValue(v)
	}
}

func (p *Patcher) scanEnumValue(v *protogen.EnumValue) {
	e := v.Parent
	opts := valueOptions(v)
	newName := opts.GetName()
	if newName == "" && p.isRenamed(e.GoIdent) {
		newName = replacePrefix(v.GoIdent.GoName, e.GoIdent.GoName, p.nameFor(e.GoIdent))
	}
	if newName != "" {
		p.RenameValue(v.GoIdent, newName) // Value const
	}
}

func (p *Patcher) scanMessage(m *protogen.Message, parent *protogen.Message) {
	opts := messageOptions(m)
	newName := opts.GetName()
	if newName == "" && parent != nil && p.isRenamed(parent.GoIdent) {
		newName = replacePrefix(m.GoIdent.GoName, parent.GoIdent.GoName, p.nameFor(parent.GoIdent))
	}
	if newName != "" {
		p.RenameType(m.GoIdent, newName) // Message struct
	}
	for _, o := range m.Oneofs {
		p.scanOneof(o)
	}
	for _, f := range m.Fields {
		p.scanField(f)
	}
	for _, e := range m.Enums {
		p.scanEnum(e)
	}
	for _, mm := range m.Messages {
		p.scanMessage(mm, m)
	}
}

func replacePrefix(s, prefix, with string) string {
	return with + strings.TrimPrefix(s, prefix)
}

func (p *Patcher) scanOneof(o *protogen.Oneof) {
	m := o.Parent
	opts := oneofOptions(o)
	newName := opts.GetName()
	if newName == "" && p.isRenamed(m.GoIdent) {
		// Implicitly rename this oneof field because its parent message was renamed.
		newName = o.GoName
	}
	if newName != "" {
		p.RenameField(ident.WithChild(m.GoIdent, o.GoName), newName)              // Oneof
		p.RenameMethod(ident.WithChild(m.GoIdent, "Get"+o.GoName), "Get"+newName) // Getter
		ifName := ident.WithPrefix(o.GoIdent, "is")
		newIfName := "is" + p.nameFor(m.GoIdent) + "_" + newName
		p.RenameType(ifName, newIfName)                                   // Interface type (e.g. isExample_Person)
		p.RenameMethod(ident.WithChild(ifName, ifName.GoName), newIfName) // Interface method
	}
	tags := opts.GetTags()
	if tags != "" {
		p.Tag(ident.WithChild(m.GoIdent, o.GoName), tags)
	}
}

func (p *Patcher) scanField(f *protogen.Field) {
	m := f.Parent
	o := f.Oneof
	opts := fieldOptions(f)
	newName := opts.GetName()
	if newName == "" && o != nil && (p.isRenamed(m.GoIdent) || p.isRenamed(o.GoIdent)) {
		// Implicitly rename this oneof field because its parent(s) were renamed.
		newName = f.GoName
	}
	if newName != "" {
		if o != nil {
			p.RenameType(f.GoIdent, p.nameFor(m.GoIdent)+"_"+newName)    // Oneof wrapper struct
			p.RenameField(ident.WithChild(f.GoIdent, f.GoName), newName) // Oneof wrapper field
			ifName := ident.WithPrefix(o.GoIdent, "is")
			p.RenameMethod(ident.WithChild(f.GoIdent, ifName.GoName), p.nameFor(ifName)) // Oneof interface method
		} else {
			p.RenameField(ident.WithChild(m.GoIdent, f.GoName), newName) // Field
		}
		p.RenameMethod(ident.WithChild(m.GoIdent, "Get"+f.GoName), "Get"+newName) // Getter
	}
	tags := opts.GetTags()
	if tags != "" {
		if o != nil {
			p.Tag(ident.WithChild(f.GoIdent, f.GoName), tags) // Oneof wrapper field tags
		} else {
			p.Tag(ident.WithChild(m.GoIdent, f.GoName), tags) // Field tags
		}
	}
}

func (p *Patcher) scanExtension(f *protogen.Field) {
	opts := fieldOptions(f)
	newName := opts.GetName()
	if newName != "" {
		id := f.GoIdent
		id.GoName = f.GoName
		p.RenameValue(ident.WithPrefix(id, "E_"), newName)
	}
}

// RenameType renames the Go type specified by id to newName.
// The id argument specifies a GoName from GoImportPath, e.g.: "github.com/org/repo/example".FooMessage
// To rename a package-level identifier such as a type, var, or const, specify just the name, e.g. "Message" or "Enum_VALUE".
// newName should be the unqualified name.
// The value of id.GoName should be the original generated type name, not a renamed type.
func (p *Patcher) RenameType(id protogen.GoIdent, newName string) {
	p.renames[id] = newName
	p.typeRenames[id] = newName
	log.Printf("Rename type:\t%s.%s → %s", id.GoImportPath, id.GoName, newName)
}

// RenameValue renames the Go value (const or var) specified by id to newName.
// The id argument specifies a GoName from GoImportPath, e.g.: "github.com/org/repo/example".FooValue
// newName should be the unqualified name.
// The value of id.GoName should be the original generated type name, not a renamed type.
func (p *Patcher) RenameValue(id protogen.GoIdent, newName string) {
	p.renames[id] = newName
	p.valueRenames[id] = newName
	log.Printf("Rename value:\t%s.%s → %s", id.GoImportPath, id.GoName, newName)
}

// RenameField renames the Go struct field specified by id to newName.
// The id argument specifies a GoName from GoImportPath, e.g.: "github.com/org/repo/example".FooMessage.BarField
// newName should be the unqualified name (after the dot).
// The value of id.GoName should be the original generated identifier name, not a renamed identifier.
func (p *Patcher) RenameField(id protogen.GoIdent, newName string) {
	p.renames[id] = newName
	p.fieldRenames[id] = newName
	log.Printf("Rename field:\t%s.%s → %s", id.GoImportPath, id.GoName, newName)
}

// RenameMethod renames the Go struct or interface method specified by id to newName.
// The id argument specifies a GoName from GoImportPath, e.g.: "github.com/org/repo/example".FooMessage.GetBarField
// The new name should be the unqualified name (after the dot).
// The value of id.GoName should be the original generated identifier name, not a renamed identifier.
func (p *Patcher) RenameMethod(id protogen.GoIdent, newName string) {
	p.renames[id] = newName
	p.methodRenames[id] = newName
	log.Printf("Rename method:\t%s.%s → %s", id.GoImportPath, id.GoName, newName)
}

func (p *Patcher) isRenamed(id protogen.GoIdent) bool {
	_, ok := p.renames[id]
	return ok
}

func (p *Patcher) nameFor(id protogen.GoIdent) string {
	if name, ok := p.renames[id]; ok {
		return name
	}
	return ident.LeafName(id)
}

// Tag adds the specified struct tags to the field specified by selector,
// in the form of "Message.Field". The tags argument should omit outer backticks (`).
// The value of id.GoName should be the original generated identifier name, not a renamed identifier.
// The struct tags will be applied when Patch is called.
func (p *Patcher) Tag(id protogen.GoIdent, tags string) {
	p.tags[id] = tags
	log.Printf("Tags:\t%s.%s `%s`", id.GoImportPath, id.GoName, tags)
}

// Patch applies the patch(es) in p the Go files in res.
// Clone res before calling Patch if you want to retain an unmodified copy.
// The behavior of calling Patch multiple times is currently undefined.
func (p *Patcher) Patch(res *pluginpb.CodeGeneratorResponse) error {
	if err := p.parseGoFiles(res); err != nil {
		return err
	}

	if err := p.checkGoFiles(); err != nil {
		return err
	}

	if err := p.patchGoFiles(); err != nil {
		return err
	}

	return p.serializeGoFiles(res)
}

func (p *Patcher) parseGoFiles(res *pluginpb.CodeGeneratorResponse) error {
	p.fset = token.NewFileSet()
	p.files = nil
	p.filesByName = make(map[string]*ast.File)

	for _, rf := range res.File {
		if rf.Name == nil || !strings.HasSuffix(*rf.Name, ".go") || rf.Content == nil {
			continue
		}

		f, err := p.parseGoFile(*rf.Name, *rf.Content)
		if err != nil {
			return err
		}

		// TODO: should we cache these?
		p.files = append(p.files, f)
		p.filesByName[*rf.Name] = f

		// FIXME: this will break if the package’s implicit name differs from the types.Package name.
		if pkg, ok := p.packagesByName[f.Name.Name]; ok {
			pkg.AddFile(*rf.Name, f)
		} else {
			return fmt.Errorf("unknown package: %s", f.Name.Name)
		}
	}

	return nil
}

func (p *Patcher) checkGoFiles() error {
	// Type-check Go packages first to find any missing identifiers.
	if err := p.checkPackages(); err != nil {
		return err
	}

	var recheck bool

	// Find missing type declarations.
	for id := range p.typeRenames {
		if obj, _ := p.find(id); obj != nil {
			continue
		}
		if err := p.synthesize(id); err != nil {
			return err
		}
		recheck = true
	}

	// Find missing value declarations.
	for id := range p.valueRenames {
		if obj, _ := p.find(id); obj != nil {
			continue
		}
		if err := p.synthesize(id); err != nil {
			return err
		}
		recheck = true
	}

	// Re-type-check if necessary.
	if recheck {
		if err := p.checkPackages(); err != nil {
			return err
		}
	}

	// Map renames.
	for id, name := range p.renames {
		obj, _ := p.find(id)
		if obj == nil {
			continue
		}
		p.objectRenames[obj] = name
	}

	// Map struct tags.
	for id, tags := range p.tags {
		obj, _ := p.find(id)
		if obj == nil {
			continue
		}
		p.fieldTags[obj] = tags
	}

	return nil
}

func (p *Patcher) parseGoFile(filename string, src interface{}) (*ast.File, error) {
	f, err := parser.ParseFile(p.fset, filename, src, parser.ParseComments)
	if err != nil {
		return nil, err
	}
	log.Printf("\nParse Go:\t%s\n", filename)
	return f, nil
}

func (p *Patcher) checkPackages() error {
	p.info = &types.Info{
		Types: make(map[ast.Expr]types.TypeAndValue),
		Defs:  make(map[*ast.Ident]types.Object),
		Uses:  make(map[*ast.Ident]types.Object),
	}

	for _, pkg := range p.packages {
		pkg.Reset()
	}

	for _, pkg := range p.packages {
		err := pkg.Check(importer{p}, p.fset, p.info)
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *Patcher) synthesize(id protogen.GoIdent) error {
	pkg := p.getPackage(string(id.GoImportPath), id.GoName, true)

	// Already synthesized?
	filename := pkg.pkg.Name() + "/" + id.GoName + ".synthetic.go"
	if f := pkg.File(filename); f != nil {
		return nil
	}

	// Synthesize a Go file for this identifier.
	b := &bytes.Buffer{}
	fmt.Fprintf(b, "package %s\n\n", pkg.pkg.Name())
	names := strings.Split(id.GoName, ".")
	if len(names) == 1 {
		// Type or value.
		// Synthesize a Go type as a map so subscript access works, e.g.: foo[key]
		fmt.Fprintf(b, "type %s map[interface{}]interface{}\n", names[0])
	} else {
		// Field or method.
		// Synthesize a Go method so a non-call expr works, e.g.: foo.Method
		fmt.Fprintf(b, "func (%s) %s() {}\n", names[0], names[1])
	}
	log.Printf("\nGenerated Go code: %s\n\n%s\n", filename, b.String())

	// Parse and add it to pkg.
	f, err := p.parseGoFile(filename, b)
	if err != nil {
		return err
	}
	return pkg.AddFile(filename, f)
}

// find finds id in all parsed Go packages, along with any ancestor(s),
// or nil if the id is not found.
func (p *Patcher) find(id protogen.GoIdent) (obj types.Object, ancestors []types.Object) {
	pkg := p.getPackage(string(id.GoImportPath), "", false)
	if pkg == nil {
		return
	}
	return pkg.Find(id)
}

// getPackage finds a getPackage with path, or creates it if it doesn’t exist.
// If name is empty, getPackage will use the last path element as the package name.
func (p *Patcher) getPackage(path, name string, create bool) *Package {
	if pkg, ok := p.packagesByPath[path]; ok {
		return pkg
	}
	if !create {
		return nil
	}
	if name == "" {
		name = filepath.Base(path)
	}
	pkg := NewPackage(path, name)
	name = pkg.pkg.Name() // Get real name
	p.packages = append(p.packages, pkg)
	p.packagesByPath[path] = pkg
	p.packagesByName[name] = pkg
	return pkg
}

func (p *Patcher) serializeGoFiles(res *pluginpb.CodeGeneratorResponse) error {
	for _, rf := range res.File {
		if rf.Name == nil || !strings.HasSuffix(*rf.Name, ".go") || rf.Content == nil {
			continue
		}
		log.Printf("\nSerialize:\t%s\n", *rf.Name)

		f := p.filesByName[*rf.Name]
		if f == nil {
			continue // Should never happen
		}

		var b strings.Builder
		err := format.Node(&b, p.fset, f)
		if err != nil {
			return err
		}

		content := b.String()
		rf.Content = &content
	}
	return nil
}

func (p *Patcher) patchGoFiles() error {
	log.Printf("\nDefs")
	for id, obj := range p.info.Defs {
		p.patchIdent(id, obj)
	}

	log.Printf("\nUses\n")
	for id, obj := range p.info.Uses {
		p.patchIdent(id, obj)
	}

	log.Printf("\nUnresolved\n")
	for _, f := range p.files {
		for _, id := range f.Unresolved {
			p.patchIdent(id, nil)
		}
	}

	return nil
}

func (p *Patcher) patchIdent(id *ast.Ident, obj types.Object) {
	// Renames
	name := p.objectRenames[obj]
	if name != "" {
		log.Printf("Rename %s:\t%q.%s → %s", typeString(obj), obj.Pkg().Path(), id.Name, name)
		p.patchComments(id, name)
		id.Name = name
	}

	// Struct tags
	tags := p.fieldTags[obj]
	if tags != "" && id.Obj != nil {
		v, ok := id.Obj.Decl.(*ast.Field)
		if !ok {
			log.Printf("Warning: struct tags declared for non-field object: %v `%s`", obj, tags)
		} else {
			if v.Tag == nil {
				v.Tag = &ast.BasicLit{}
			}
			v.Tag.Value = "`" + strings.TrimSpace(strings.Trim(v.Tag.Value, "` ")+" "+tags) + "`"
			log.Printf("Add tags:\t%q.%s %s", obj.Pkg().Path(), id.Name, v.Tag.Value)
		}
	}
}

func (p *Patcher) patchComments(id *ast.Ident, repl string) {
	doc, comment := p.findCommentGroups(id)
	if doc == nil && comment == nil {
		return
	}
	x, err := regexp.Compile(`\b` + regexp.QuoteMeta(id.Name) + `\b`)
	if err != nil {
		return
	}
	log.Printf("Comment:\t%v → %s", x, repl)
	patchCommentGroup(doc, x, repl)
	patchCommentGroup(comment, x, repl)
}

// Borrowed from https://github.com/golang/tools/blob/HEAD/refactor/rename/rename.go#L543
func (p *Patcher) findCommentGroups(id *ast.Ident) (doc *ast.CommentGroup, comment *ast.CommentGroup) {
	tf := p.fset.File(id.Pos())
	if tf == nil {
		return
	}
	f := p.filesByName[tf.Name()]
	if f == nil {
		return
	}
	nodes, _ := astutil.PathEnclosingInterval(f, id.Pos(), id.End())
	for _, node := range nodes {
		switch decl := node.(type) {
		case *ast.FuncDecl:
			return decl.Doc, nil
		case *ast.Field:
			return decl.Doc, decl.Comment
		case *ast.GenDecl:
			return decl.Doc, nil
		// For {Type,Value}Spec, if the doc on the spec is absent, search for the enclosing GenDecl.
		case *ast.TypeSpec:
			if decl.Doc != nil {
				return decl.Doc, decl.Comment
			}
		case *ast.ValueSpec:
			if decl.Doc != nil {
				return decl.Doc, decl.Comment
			}
		case *ast.Ident:
		default:
			return
		}
	}
	return
}

func patchCommentGroup(c *ast.CommentGroup, x *regexp.Regexp, repl string) {
	if c == nil {
		return
	}
	for _, c := range c.List {
		c.Text = x.ReplaceAllString(c.Text, repl)
	}
}

func typeString(obj types.Object) string {
	switch obj.(type) {
	case *types.PkgName:
		return "package name"
	case *types.TypeName:
		return "type"
	case *types.Var:
		if obj.Parent() == nil {
			return "field"
		}
		return "var"
	case *types.Const:
		return "const"
	case *types.Func:
		if obj.Parent() == nil {
			return "method"
		}
		return "func"
	case nil:
		return "(nil)"
	}
	return obj.Type().String()
}
