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

func TestRenamedEnumValue(t *testing.T) {
	tests.ValidateEnum(t, EnumWithRenamedValue(0), EnumWithRenamedValue_name, EnumWithRenamedValue_value)
	if got, want := RenamedValue, EnumWithRenamedValue(0); got != want {
		t.Errorf("%T(%d) != %v", got, got, want)
	}
}

func TestCustomStringerEnum(t *testing.T) {
	tests := []struct {
		enum     CustomStringerEnum
		original string
		patched  string
	}{
		{0, "CUSTOM_STRINGER_INVALID", "custom_stringer_invalid"},
		{1, "CUSTOM_STRINGER_A", "custom_stringer_a"},
		{2, "CUSTOM_STRINGER_B", "custom_stringer_b"},
		{3, "CUSTOM_STRINGER_C", "custom_stringer_c"},
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

func TestDefaultStringerEnum(t *testing.T) {
	e := DefaultStringerEnum(0)
	if got, want := e.String(), "DEFAULT_STRINGER_UNSET"; got != want {
		t.Errorf("%T(%d) incorrect original string %q != %q", e, e, got, want)
	}
}
