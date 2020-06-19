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

func TestRenamedField(t *testing.T) {
	tests.ValidateMessage(t, &MessageWithRenamedField{})
	// TODO: validate renamed field
}
