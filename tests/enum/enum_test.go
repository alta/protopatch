package enum

import (
	"fmt"
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

func TestStringerName(t *testing.T) {
	tests := []struct {
		enum     CustomStringer
		original string
		patched  string
	}{
		{0, "INVALID", "invalid"},
		{1, "A", "a"},
		{2, "B", "b"},
		{3, "C", "c"},
	}
	for _, tt := range tests {
		name := fmt.Sprintf("enum(%d)/%s/%s", int32(tt.enum), tt.original, tt.patched)
		t.Run(name, func(t *testing.T) {
			if s := tt.enum.OrigString(); s != tt.original {
				t.Errorf("%T(%d) incorrect original string %q != %q", tt.enum, tt.enum, s, tt.original)
			}
			if s := tt.enum.String(); s != tt.patched {
				t.Errorf("%T(%d) incorrect patched string %q != %q", tt.enum, tt.enum, s, tt.patched)
			}
		})
	}
}
