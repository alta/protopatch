package lint

import (
	"testing"

	"github.com/alta/protopatch/tests"
	"google.golang.org/protobuf/proto"
)

func TestProtocol(t *testing.T) {
	tests.ValidateEnum(t, Protocol(0), Protocol_name, Protocol_value)
	tests.ValidateEnum(t, ProtocolInvalid, Protocol_name, Protocol_value)
	tests.ValidateEnum(t, ProtocolIP, Protocol_name, Protocol_value)
	tests.ValidateEnum(t, ProtocolUDP, Protocol_name, Protocol_value)
	tests.ValidateEnum(t, ProtocolTCP, Protocol_name, Protocol_value)
}

func TestURL(t *testing.T) {
	tests.ValidateMessage(t, &URL{})
}

func TestID(t *testing.T) {
	tests.ValidateMessage(t, &ID{})
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
	var _ *OuterMessageInnerID = m.GetID()
	var _ *OuterMessageInnerURL = m.GetURL()
}

func TestExtendedMessage(t *testing.T) {
	m := &ExtendedMessage{}
	tests.ValidateMessage(t, m)
	_ = proto.GetExtension(m, ExtAlpha).(string)
	_ = proto.GetExtension(m, ExtBeta).(string)
	_ = proto.GetExtension(m, ExtGamma).(string)
	_ = proto.GetExtension(m, ExtDelta).(string)
}
