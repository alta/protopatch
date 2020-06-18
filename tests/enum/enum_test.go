package enum

import (
	"testing"

	"github.com/alta/protopatch/tests"
)

func TestBasic(t *testing.T) {
	tests.ValidateEnum(t, Basic(0), Basic_name, Basic_value)
}

func TestNested(t *testing.T) {
	tests.ValidateEnum(t, Outer_Inner(0), Outer_Inner_name, Outer_Inner_value)
}

func TestRenamed(t *testing.T) {
	tests.ValidateEnum(t, Renamed(0), Renamed_name, Renamed_value)
}
