package patch

import (
	"github.com/alta/protopatch/patch/gopb"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func enumOptions(e *protogen.Enum) *gopb.Options {
	return proto.GetExtension(e.Desc.Options(), gopb.E_Enum).(*gopb.Options)
}

func valueOptions(v *protogen.EnumValue) *gopb.Options {
	return proto.GetExtension(v.Desc.Options(), gopb.E_Value).(*gopb.Options)
}

func messageOptions(m *protogen.Message) *gopb.Options {
	return proto.GetExtension(m.Desc.Options(), gopb.E_Message).(*gopb.Options)
}

func fieldOptions(f *protogen.Field) *gopb.Options {
	return proto.GetExtension(f.Desc.Options(), gopb.E_Field).(*gopb.Options)
}

func oneofOptions(o *protogen.Oneof) *gopb.Options {
	return proto.GetExtension(o.Desc.Options(), gopb.E_Oneof).(*gopb.Options)
}

func fileValuesOptions(fileDescriptor protoreflect.FileDescriptor) *gopb.Options {
	return proto.GetExtension(fileDescriptor.Options(), gopb.E_Values).(*gopb.Options)
}

func fileInitialisms(fileDescriptor protoreflect.FileDescriptor) []string {
	return proto.GetExtension(fileDescriptor.Options(), gopb.E_Initialisms).([]string)
}

func fileInitialismsAsMap(fileDescriptor protoreflect.FileDescriptor) map[string]bool {
	result := make(map[string]bool)
	for _, s := range fileInitialisms(fileDescriptor) {
		result[s] = true
	}

	return result
}
