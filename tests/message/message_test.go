package message

import (
	"testing"

	"github.com/stretchr/testify/assert"
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
	m := &MessageWithRenamedField{
		ID: 66,
	}
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

func TestMessageWithCustomTypes(t *testing.T) {
	m := &MessageWithCustomTypes{
		StringField: "42",
		Int32Field:  42,
		Int64Field:  42,
		FloatField:  42,
		DoubleField: 42,
		Uint32Field: 42,
		Uint64Field: 42,
	}

	tests.ValidateMessage(t, m)
	var _ string = string(m.StringField)
	var _ int32 = int32(m.Int32Field)
	var _ int64 = int64(m.Int64Field)
	var _ float32 = float32(m.FloatField)
	var _ float64 = float64(m.DoubleField)
	var _ uint32 = uint32(m.Uint32Field)
	var _ uint64 = uint64(m.Uint64Field)

	assert.Equal(t, String("42"), m.StringField)
	assert.Equal(t, Int32(42), m.Int32Field)
	assert.Equal(t, Int64(42), m.Int64Field)
	assert.Equal(t, Float(42), m.FloatField)
	assert.Equal(t, Double(42), m.DoubleField)
	assert.Equal(t, Uint32(42), m.Uint32Field)
	assert.Equal(t, Uint64(42), m.Uint64Field)
}

func TestMessageWithOptionalCustomTypes(t *testing.T) {
	var (
		StringValue String = "42"
		Int32Value  Int32  = 42
		Int64Value  Int64  = 42
		FloatValue  Float  = 42
		DoubleValue Double = 42
		Uint32Value Uint32 = 42
		Uint64Value Uint64 = 42
	)
	m := &MessageWithOptionalCustomTypes{
		StringField: &StringValue,
		Int32Field:  &Int32Value,
		Int64Field:  &Int64Value,
		FloatField:  &FloatValue,
		DoubleField: &DoubleValue,
		Uint32Field: &Uint32Value,
		Uint64Field: &Uint64Value,
	}

	tests.ValidateMessage(t, m)
	var _ string = string(m.GetStringField())
	var _ int32 = int32(m.GetInt32Field())
	var _ int64 = int64(m.GetInt64Field())
	var _ float32 = float32(m.GetFloatField())
	var _ float64 = float64(m.GetDoubleField())
	var _ uint32 = uint32(m.GetUint32Field())
	var _ uint64 = uint64(m.GetUint64Field())

	assert.Equal(t, &StringValue, m.StringField)
	assert.Equal(t, &Int32Value, m.Int32Field)
	assert.Equal(t, &Int64Value, m.Int64Field)
	assert.Equal(t, &FloatValue, m.FloatField)
	assert.Equal(t, &DoubleValue, m.DoubleField)
	assert.Equal(t, &Uint32Value, m.Uint32Field)
	assert.Equal(t, &Uint64Value, m.Uint64Field)
}

func TestMessageWithCustomRepeatedType(t *testing.T) {
	slice := Strings{"one", "two"}
	m := &MessageWithCustomRepeatedType{
		RepeatedStringField: slice,
	}
	tests.ValidateMessage(t, m)
	var _ Strings = m.RepeatedStringField
	assert.Equal(t, slice, m.RepeatedStringField)
}
