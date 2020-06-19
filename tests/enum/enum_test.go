package enum

import (
	"testing"

	"github.com/alta/protopatch/tests"
)

func TestBasicEnum(t *testing.T) {
	tests.ValidateEnum(t, BasicEnum(0), BasicEnum_name, BasicEnum_value)
}

func TestNestedEnum(t *testing.T) {
	tests.ValidateEnum(t, OuterMessage_InnerEnum(0), OuterMessage_InnerEnum_name, OuterMessage_InnerEnum_value)
}

func TestRenamedEnum(t *testing.T) {
	tests.ValidateEnum(t, RenamedEnum(0), RenamedEnum_name, RenamedEnum_value)
}
