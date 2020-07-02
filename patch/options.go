package patch

import (
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
)

func enumOptions(e *protogen.Enum) *GoOptions {
	return proto.GetExtension(e.Desc.Options(), E_Enum).(*GoOptions)
}

func valueOptions(v *protogen.EnumValue) *GoOptions {
	return proto.GetExtension(v.Desc.Options(), E_Value).(*GoOptions)
}

func messageOptions(m *protogen.Message) *GoOptions {
	return proto.GetExtension(m.Desc.Options(), E_Message).(*GoOptions)
}

func fieldOptions(f *protogen.Field) *GoOptions {
	return proto.GetExtension(f.Desc.Options(), E_Field).(*GoOptions)
}

func oneofOptions(o *protogen.Oneof) *GoOptions {
	return proto.GetExtension(o.Desc.Options(), E_Oneof).(*GoOptions)
}
