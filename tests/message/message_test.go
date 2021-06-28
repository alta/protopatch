package message

import (
	"testing"

	"google.golang.org/protobuf/proto"

	"github.com/alta/protopatch/tests"
)

func TestBasicMessage(t *testing.T) {
	tests.ValidateMessage(t, &Basic{})
}

func TestOneofMessage(t *testing.T) {
	m := &Union{}
	tests.ValidateMessage(t, m)
	var _ isUnion_Contents = &Union_Id{}
	var _ isUnion_Contents = &Union_Name{}
	var _ int32 = m.GetId()
	var _ string = m.GetName()
}

func TestNestedMessage(t *testing.T) {
	m := &Outer{}
	tests.ValidateMessage(t, m)
	tests.ValidateMessage(t, &Outer_Middle{})
	tests.ValidateMessage(t, &Outer_Middle_Inner{})
	var _ *Outer_Middle = m.GetMiddle()
	var _ *Outer_Middle_Inner = m.GetInner()
	var _ *Outer_Middle_Inner = m.GetMiddle().GetInner()
}

func TestRenamedMessage(t *testing.T) {
	tests.ValidateMessage(t, &Frank{})
}

func TestRenamedOneofMessage(t *testing.T) {
	m := &RenamedOneofMessage{}
	tests.ValidateMessage(t, m)
	var _ isRenamedOneofMessage_Contents = &RenamedOneofMessage_Id{}
	var _ isRenamedOneofMessage_Contents = &RenamedOneofMessage_Name{}
	var _ int32 = m.GetId()
	var _ string = m.GetName()
}

func TestRenamedOuterMessage(t *testing.T) {
	tests.ValidateMessage(t, &RenamedOuterMessage{})
	tests.ValidateMessage(t, &RenamedOuterMessage_InnerMessage{})
}

func TestRenamedInnerMessage(t *testing.T) {
	tests.ValidateMessage(t, &OuterMessageWithRenamedInnerMessage{})
	tests.ValidateMessage(t, &RenamedInnerMessage{})
}

func TestMessageWithRenamedField(t *testing.T) {
	m := &MessageWithRenamedField{}
	tests.ValidateMessage(t, m)
	var _ int32 = m.ID
	var _ int32 = m.GetID()
}

func TestMessageWithEmbeddedField(t *testing.T) {
	m := &MessageWithEmbeddedField{
		RenamedOuterMessage: &RenamedOuterMessage{
			Inner: &RenamedOuterMessage_InnerMessage{

			},
		},
	}
	tests.ValidateMessage(t, m)
	if &m.Inner != &m.RenamedOuterMessage.Inner {
		t.Error("RenamedOuterMessage message is not embedded")
	}
}

func TestMessageWithStructTags(t *testing.T) {
	m := &MessageWithTags{}
	tests.ValidateTag(t, m, "Value", "test", "value")
}

func TestNestedMessageWithStructTags(t *testing.T) {
	m := &OuterMessageWithTags_InnerMessage{}
	tests.ValidateTag(t, m, "Value", "test", "value")
	tests.ValidateTag(t, m, "Value", "json", "value,omitempty")
}

func TestMessageWithJSONTags(t *testing.T) {
	m := &MessageWithJSONTags{}
	tests.ValidateTag(t, m, "Value", "json", "custom_value")
	tests.ValidateTag(t, m, "Empty", "json", "custom_empty,omitempty")
}

func TestMessageWithRedundantTags(t *testing.T) {
	m := &MessageWithRedundantTags{}
	tests.ValidateTag(t, m, "Value", "test", "3")
	tests.ValidateTag(t, m, "Value", "json", "value,omitempty")
}

func TestExtendedMessage(t *testing.T) {
	m := &ExtendedMessage{}
	tests.ValidateMessage(t, m)
	_ = proto.GetExtension(m, E_Alpha).(string)
	_ = proto.GetExtension(m, E_Beta).(string)
	_ = proto.GetExtension(m, ExtGamma).(string)
	_ = proto.GetExtension(m, ExtDelta).(string)
}
