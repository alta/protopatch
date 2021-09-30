package patch

import (
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/compiler/protogen"
)

const (
	srcDef = `package foo

import (
	"fmt"
)

type String string

type Message struct {
	Content string
}

func (m *Message) GetContent() string {
	if m != nil {
		return m.Content
	}
	return ""
}

func print(s string) {
	fmt.Println(s)
}
`
	wantDef = `package foo

import (
	"fmt"
)

type String string

type Message struct {
	Content String
}

func (m *Message) GetContent() String {
	if m != nil {
		return m.Content
	}
	return ""
}

func print(s string) {
	fmt.Println(s)
}
`
)

const (
	fileName    = "foo.go"
	packageName = "foo"
	fieldName   = "Content"
	msgName   = "Message"
	fieldType = "String"
)

func prepareCastType(src string) (*Patcher, *ast.File, error) {
	p, err := NewPatcher(&protogen.Plugin{})
	if err != nil {
	    return nil, nil, err
	}
	p.filesByName = make(map[string]*ast.File)
	p.info = &types.Info{
		Types: make(map[ast.Expr]types.TypeAndValue),
		Defs:  make(map[*ast.Ident]types.Object),
		Uses:  make(map[*ast.Ident]types.Object),
	}
	p.fset = token.NewFileSet()
	file, err := parser.ParseFile(p.fset, fileName, src, parser.ParseComments)
	if err != nil {
		return nil, nil, err
	}
	pkg := NewPackage(packageName, packageName)
	if err := pkg.AddFile(fileName, file); err != nil {
		return nil, nil, err
	}
	if err := pkg.Check(basicImporter{p}, p.fset, p.info); err != nil {
		return nil, nil, err
	}
	p.filesByName[fileName] = file
	p.packagesByPath[packageName] = pkg
	p.packagesByName[packageName] = pkg
	p.packages = append(p.packages, pkg)
	p.Type(protogen.GoIdent{GoName: msgName + "." + fieldName, GoImportPath: packageName}, fieldType)
	p.Type(protogen.GoIdent{GoName: msgName + "." + "Get" + fieldName, GoImportPath: packageName}, fieldType)
	// Map cast types
	for id, typ := range p.types {
		obj, _ := p.find(id)
		if obj == nil {
			continue
		}
		p.fieldTypes[obj] = typ
	}
	return p, file, nil
}

func TestScalarCastType(t *testing.T) {
	tests := []struct{
		name string
		src string
		want string
	}{
		{
			name: "cast definition",
			src: srcDef,
			want: wantDef,
		},
		{
			name: "cast field initialization",
			src: srcDef+`
func useContent() {
	s := "ok"
	msg := &Message{Content: s}
}
`,
			want: wantDef+`
func useContent() {
	s := "ok"
	msg := &Message{Content: String(s)}
}
`,
		},
		{
			name: "cast field assignation",
			src: srcDef+`
func useContent() {
	s := "ok"
	msg := &Message{}
	msg.Content = s
}
`,
			want: wantDef+`
func useContent() {
	s := "ok"
	msg := &Message{}
	msg.Content = String(s)
}
`,
		},
		{
			name: "cast field assignation usage",
			src: srcDef+`
func useContent() {
	var s string
	msg := &Message{Content: "ok"}
	s = msg.Content
}
`,
			want: wantDef+`
func useContent() {
	var s string
	msg := &Message{Content: String("ok")}
	s = string(msg.Content)
}
`,
		},
		{
			name: "cast field used as argument",
			src: srcDef+`
func useContent() {
	msg := &Message{Content: "ok"}
	var s string
	s = msg.Content
}
`,
			want: wantDef+`
func useContent() {
	msg := &Message{Content: String("ok")}
	var s string
	s = string(msg.Content)
}
`,
		},
		{
			name: "cast field already casted",
			src: srcDef+`
type OtherString string

func multi(s1, s2, s3 OtherString) {
}

func useContent() {
	msg := &Message{}
	multi("", "", OtherString(msg.Content))
}
`,
			want: wantDef+`
type OtherString string

func multi(s1, s2, s3 OtherString) {
}

func useContent() {
	msg := &Message{}
	multi("", "", OtherString(string(msg.Content)))
}
`,
		},
		{
			name: "cast full code",
			src: srcDef+`
func useContent() {
	s := "ok"
	msg := &Message{Content: s}
	print(msg.Content)
	print(msg.GetContent())
	msg.Content = s
	useContentGetter(msg)
}

func useContentGetter(msg *Message) {
	var s string
	s = msg.GetContent()
	s = msg.Content
	s = msg.Content + "..."
	print(s)
}
`,
			want: wantDef+`
func useContent() {
	s := "ok"
	msg := &Message{Content: String(s)}
	print(string(msg.Content))
	print(string(msg.GetContent()))
	msg.Content = String(s)
	useContentGetter(msg)
}

func useContentGetter(msg *Message) {
	var s string
	s = string(msg.GetContent())
	s = string(msg.Content)
	s = string(msg.Content) + "..."
	print(s)
}
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p, file, err := prepareCastType(tt.src)
			if err != nil {
				t.Errorf("failed to initilialize test")
				return
			}

			if err := p.patchGoFiles(); err != nil {
				t.Fatal(err)
			}
			got := p.nodeToString(file)
			assert.Equal(t, tt.want, got)
		})
	}

}
