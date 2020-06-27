package patch

import (
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
)

func enumOptions(e *protogen.Enum) *Options {
	return proto.GetExtension(e.Desc.Options(), E_Enum).(*Options)
}

func valueOptions(v *protogen.EnumValue) *Options {
	return proto.GetExtension(v.Desc.Options(), E_Value).(*Options)
}

func messageOptions(m *protogen.Message) *Options {
	return proto.GetExtension(m.Desc.Options(), E_Message).(*Options)
}

func fieldOptions(f *protogen.Field) *Options {
	return proto.GetExtension(f.Desc.Options(), E_Field).(*Options)
}

func oneofOptions(o *protogen.Oneof) *Options {
	return proto.GetExtension(o.Desc.Options(), E_Oneof).(*Options)
}
