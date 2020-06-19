package message

import (
	"testing"

	"github.com/alta/protopatch/tests"
)

func TestBasicMessage(t *testing.T) {
	tests.ValidateMessage(t, &BasicMessage{})
}

func TestNestedMessage(t *testing.T) {
	tests.ValidateMessage(t, &OuterMessage{})
	tests.ValidateMessage(t, &OuterMessage_InnerMessage{})
}

func TestRenamedMessage(t *testing.T) {
	tests.ValidateMessage(t, &RenamedMessage{})
}

func TestRenamedNestedMessage(t *testing.T) {
	tests.ValidateMessage(t, &RenamedOuterMessage{})
	tests.ValidateMessage(t, &RenamedOuterMessage_InnerMessage{})
}

func TestRenamedField(t *testing.T) {
	m := &MessageWithRenamedField{}
	tests.ValidateMessage(t, m)
	var _ int32 = m.ID
	var _ int32 = m.GetID()
}
