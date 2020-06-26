package options

import (
	patch_go "github.com/alta/protopatch/patch/go"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// A Provider provides patch_go.Options for various types.
// A patch.Patcher is an Provider, as well as any compatibility shims.
type Provider interface {
	EnumOptions(*protogen.Enum) *patch_go.Options
	ValueOptions(*protogen.EnumValue) *patch_go.Options
	MessageOptions(*protogen.Message) *patch_go.Options
	FieldOptions(*protogen.Field) *patch_go.Options
	OneofOptions(*protogen.Oneof) *patch_go.Options
}

// UnimplementedProvider can be embedded in a struct to ensure it implements the Provider interface.
type UnimplementedProvider struct{}

// EnumOptions returns nil.
func (UnimplementedProvider) EnumOptions(*protogen.Enum) *patch_go.Options { return nil }

// ValueOptions returns nil.
func (UnimplementedProvider) ValueOptions(*protogen.EnumValue) *patch_go.Options { return nil }

// MessageOptions returns nil.
func (UnimplementedProvider) MessageOptions(*protogen.Message) *patch_go.Options { return nil }

// FieldOptions returns nil.
func (UnimplementedProvider) FieldOptions(*protogen.Field) *patch_go.Options { return nil }

// OneofOptions returns nil.
func (UnimplementedProvider) OneofOptions(*protogen.Oneof) *patch_go.Options { return nil }

// Ensure UnimplementedProvider implements Provider.
var _ Provider = UnimplementedProvider{}

// Get returns patch_go.Options if present on desc, otherwise nil.
func Get(desc protoreflect.Descriptor, xt protoreflect.ExtensionType) *patch_go.Options {
	o := desc.Options()
	if proto.HasExtension(o, xt) {
		return proto.GetExtension(o, xt).(*patch_go.Options)
	}
	return nil
}
