package enum

import (
	"fmt"
	"testing"

	"github.com/alta/protopatch/tests"
)

func TestKind(t *testing.T) {
	tests.ValidateEnum(t, Kind(0), Kind_name, Kind_value)
	enums := []Kind{
		Kind_INVALID,
		Kind_CHEAP,
		Kind_FAST,
		Kind_GOOD,
	}
	for index, enum := range enums {
		if got, want := enum, Kind(index); got != want {
			t.Errorf("%T(%d) != %v", got, got, want)
		}
	}
}

func TestNestedEnums(t *testing.T) {
	tests.ValidateEnum(t, Outer_Route(0), Outer_Route_name, Outer_Route_value)
	tests.ValidateEnum(t, Outer_Middle_Flavor(0), Outer_Middle_Flavor_name, Outer_Middle_Flavor_value)
	tests.ValidateEnum(t, Outer_Middle_Inner_Arch(0), Outer_Middle_Inner_Arch_name, Outer_Middle_Inner_Arch_value)
}

func TestRenamedEnum(t *testing.T) {
	tests.ValidateEnum(t, Flavour(0), Flavour_name, Flavour_value)
	enums := []Flavour{
		Flavour_INVALID,
		Flavour_SWEET,
		Flavour_SALTY,
		Flavour_SOUR,
		Flavour_BITTER,
	}
	for index, enum := range enums {
		if got, want := enum, Flavour(index); got != want {
			t.Errorf("%T(%d) != %v", got, got, want)
		}
	}
}

func TestRenamedEnumValue(t *testing.T) {
	tests.ValidateEnum(t, Level(0), Level_name, Level_value)
	enums := []Level{
		LevelSimple,
		Level_COMPLEX,
	}
	for index, enum := range enums {
		if got, want := enum, Level(index); got != want {
			t.Errorf("%T(%d) != %v", got, got, want)
		}
	}
}

func TestRenamedNestedEnumValue(t *testing.T) {
	tests.ValidateEnum(t, RenamedNested(0), RenamedNested_name, RenamedNested_value)
	enums := []RenamedNested{
		RenamedValueInvalid,
		RenamedValueA,
		RenamedValueB,
		RenamedValueC,
	}
	for index, enum := range enums {
		if got, want := enum, RenamedNested(index); got != want {
			t.Errorf("%T(%d) != %v", got, got, want)
		}
	}
}

func TestRenamedOuterMessage(t *testing.T) {
	m := &Wrapper{}
	tests.ValidateMessage(t, m)
	tests.ValidateEnum(t, Holiday_Route(0), Holiday_Route_name, Holiday_Route_value)
	enums := []Holiday_Route{
		Holiday_INVALID,
		Holiday_FAST,
		Holiday_SLOW,
		Holiday_SCENIC,
	}
	for index, enum := range enums {
		if got, want := enum, Holiday_Route(index); got != want {
			t.Errorf("%T(%d) != %v", got, got, want)
		}
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

func TestDeprecatedStringerEnum(t *testing.T) {
	tests := []struct {
		enum     DeprecatedStringerEnum
		original string
		patched  string
	}{
		{0, "DEPRECATED_STRINGER_INVALID", "deprecated_stringer_invalid"},
		{1, "DEPRECATED_STRINGER_A", "deprecated_stringer_a"},
		{2, "DEPRECATED_STRINGER_B", "deprecated_stringer_b"},
		{3, "DEPRECATED_STRINGER_C", "deprecated_stringer_c"},
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
