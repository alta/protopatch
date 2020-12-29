package lint

import (
	"testing"

	"github.com/alta/protopatch/tests"
	"google.golang.org/protobuf/proto"
)

func TestURL(t *testing.T) {
	tests.ValidateMessage(t, &URL{})
}

func TestID(t *testing.T) {
	tests.ValidateMessage(t, &ID{})
}

func TestRGBColor(t *testing.T) {
	tests.ValidateMessage(t, &RGBColor{})
}

func TestOneofMessage(t *testing.T) {
	m := &OneofMessage{}
	tests.ValidateMessage(t, m)
	var _ isOneofMessage_Contents = &OneofMessage_ID{}
	var _ isOneofMessage_Contents = &OneofMessage_URL{}
	var _ *ID = m.GetID()
	var _ *URL = m.GetURL()
}

func TestOuterMessage(t *testing.T) {
	m := &OuterMessage{}
	tests.ValidateMessage(t, &OuterMessage{})
	tests.ValidateMessage(t, &OuterMessageInnerID{})
	tests.ValidateMessage(t, &OuterMessageInnerURL{})
	tests.ValidateEnum(t, OuterMessageFlavor(0), OuterMessageFlavor_name, OuterMessageFlavor_value)
	var _ *OuterMessageInnerID = m.GetID()
	var _ *OuterMessageInnerURL = m.GetURL()
}

func TestColor(t *testing.T) {
	m := &Color{}
	tests.ValidateMessage(t, m)
	var _ string = m.GetRGB()
	var _ string = m.GetRGBA()
	var _ string = m.GetHSV()
}

func TestExtendedMessage(t *testing.T) {
	m := &ExtendedMessage{}
	tests.ValidateMessage(t, m)
	_ = proto.GetExtension(m, ExtAlpha).(string)
	_ = proto.GetExtension(m, ExtBeta).(string)
	_ = proto.GetExtension(m, ExtGamma).(string)
	_ = proto.GetExtension(m, ExtDelta).(string)
}

func TestBasic(t *testing.T) {
	tests.ValidateEnum(t, Basic(0), Basic_name, Basic_value)
	tests.ValidateEnum(t, BasicInvalid, Basic_name, Basic_value)
	tests.ValidateEnum(t, BasicA, Basic_name, Basic_value)
	tests.ValidateEnum(t, BasicB, Basic_name, Basic_value)
	tests.ValidateEnum(t, BasicC, Basic_name, Basic_value)
}

func TestProtocol(t *testing.T) {
	tests.ValidateEnum(t, Protocol(0), Protocol_name, Protocol_value)
	tests.ValidateEnum(t, ProtocolInvalid, Protocol_name, Protocol_value)
	tests.ValidateEnum(t, ProtocolIP, Protocol_name, Protocol_value)
	tests.ValidateEnum(t, ProtocolUDP, Protocol_name, Protocol_value)
	tests.ValidateEnum(t, ProtocolTCP, Protocol_name, Protocol_value)
}
