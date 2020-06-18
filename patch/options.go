package patch

import (
	"github.com/alta/protopatch/patch/go/enum"
	"github.com/alta/protopatch/patch/go/message"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
)

func enumOptions(e *protogen.Enum) *enum.Options {
	return proto.GetExtension(e.Desc.Options(), enum.E_Options).(*enum.Options)
}

func messageOptions(m *protogen.Message) *message.Options {
	return proto.GetExtension(m.Desc.Options(), message.E_Options).(*message.Options)
}
