package options

import (
	patch_go "github.com/alta/protopatch/patch/go"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// Get returns patch_go.Options if present on desc, otherwise nil.
func Get(desc protoreflect.Descriptor, xt protoreflect.ExtensionType) *patch_go.Options {
	o := desc.Options()
	if proto.HasExtension(o, xt) {
		return proto.GetExtension(o, xt).(*patch_go.Options)
	}
	return nil
}
