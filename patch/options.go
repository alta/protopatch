package patch

import (
	"github.com/alta/protopatch/patch/go/enum"
	"github.com/alta/protopatch/patch/go/field"
	"github.com/alta/protopatch/patch/go/message"
	"github.com/alta/protopatch/patch/go/oneof"
	"github.com/alta/protopatch/patch/go/value"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
)

func enumOptions(e *protogen.Enum) *enum.Options {
	return proto.GetExtension(e.Desc.Options(), enum.E_Options).(*enum.Options)
}

func valueOptions(v *protogen.EnumValue) *value.Options {
	return proto.GetExtension(v.Desc.Options(), value.E_Options).(*value.Options)
}

func messageOptions(m *protogen.Message) *message.Options {
	return proto.GetExtension(m.Desc.Options(), message.E_Options).(*message.Options)
}

func fieldOptions(f *protogen.Field) *field.Options {
	return proto.GetExtension(f.Desc.Options(), field.E_Options).(*field.Options)
}

func oneofOptions(o *protogen.Oneof) *oneof.Options {
	return proto.GetExtension(o.Desc.Options(), oneof.E_Options).(*oneof.Options)
}
