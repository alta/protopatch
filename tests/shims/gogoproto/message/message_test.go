package message

import (
	"testing"

	"github.com/alta/protopatch/tests"
)

func TestMessageWithRenamedField(t *testing.T) {
	m := &MessageWithRenamedField{}
	tests.ValidateMessage(t, m)
	var _ int32 = m.ID
	var _ int32 = m.GetID()
}
