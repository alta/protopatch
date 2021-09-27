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

func TestMessageWithEmbeddedFields(t *testing.T) {
	message := "noop"
	m := &MessageWithEmbeddedField{
		Embedded: &Embedded{
			Message: message,
		},
	}
	tests.ValidateMessage(t, m)
	var _ *Embedded = m.Embedded
	if m.Message != message {
		t.Errorf("inalid Embedded.Message: expected '%s' got '%s'", message, m.Message)
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

func TestMessageWithOptionals(t *testing.T) {
	m := &MessageWithOptionals{
		OptionalString: proto.String("42"),
		OptionalInt32:  proto.Int32(42),
		OptionalInt64:  proto.Int64(42),
		OptionalFloat:  proto.Float32(42),
		OptionalDouble: proto.Float64(42),
		OptionalUInt32: proto.Uint32(42),
		OptionalUInt64: proto.Uint64(42),
		OptionalBool:   proto.Bool(false),
	}
	tests.ValidateMessage(t, m)
	var _ string = m.GetOptionalString()
	var _ int32 = m.GetOptionalInt32()
	var _ int64 = m.GetOptionalInt64()
	var _ float32 = m.GetOptionalFloat()
	var _ float64 = m.GetOptionalDouble()
	var _ uint32 = m.GetOptionalUInt32()
	var _ uint64 = m.GetOptionalUInt64()
	var _ bool = m.GetOptionalBool()

	if v := m.GetOptionalString(); v != "42" {
		t.Errorf("MessageWithOptions.GetOptionalString(): got %v, expected \"42\"", v)
	}
	if v := m.GetOptionalInt32(); v != 42 {
		t.Errorf("MessageWithOptions.GetOptionalInt32(): got %v, expected 42", v)
	}
	if v := m.GetOptionalInt64(); v != 42 {
		t.Errorf("MessageWithOptions.GetOptionalInt64(): got %v, expected 42", v)
	}
	if v := m.GetOptionalFloat(); v != 42 {
		t.Errorf("MessageWithOptions.GetOptionalFloat(): got %v, expected 42", v)
	}
	if v := m.GetOptionalDouble(); v != 42 {
		t.Errorf("MessageWithOptions.GetOptionalDouble(): got %v, expected 42", v)
	}
	if v := m.GetOptionalUInt32(); v != 42 {
		t.Errorf("MessageWithOptions.GetOptionalUInt32(): got %v, expected 42", v)
	}
	if v := m.GetOptionalUInt64(); v != 42 {
		t.Errorf("MessageWithOptions.GetOptionalUInt64(): got %v, expected 42", v)
	}
	if m.OptionalBool == nil || m.GetOptionalBool() {
		t.Errorf("MessageWithOptions.GetOptionalBool() is true expected false")
	}
}
