package patch

import (
	patch_go "github.com/alta/protopatch/patch/go"
	"github.com/alta/protopatch/patch/go/enum"
	"github.com/alta/protopatch/patch/go/field"
	"github.com/alta/protopatch/patch/go/message"
	"github.com/alta/protopatch/patch/go/oneof"
	"github.com/alta/protopatch/patch/go/value"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
)

func enumOptions(e *protogen.Enum) *patch_go.Options {
	return proto.GetExtension(e.Desc.Options(), enum.E_Options).(*patch_go.Options)
}

func valueOptions(v *protogen.EnumValue) *patch_go.Options {
	return proto.GetExtension(v.Desc.Options(), value.E_Options).(*patch_go.Options)
}

func messageOptions(m *protogen.Message) *patch_go.Options {
	return proto.GetExtension(m.Desc.Options(), message.E_Options).(*patch_go.Options)
}

func fieldOptions(f *protogen.Field) *patch_go.Options {
	return proto.GetExtension(f.Desc.Options(), field.E_Options).(*patch_go.Options)
}

func oneofOptions(o *protogen.Oneof) *patch_go.Options {
	return proto.GetExtension(o.Desc.Options(), oneof.E_Options).(*patch_go.Options)
}
