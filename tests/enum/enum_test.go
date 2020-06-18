package enum

import (
	"testing"

	"github.com/alta/protopatch/tests"
)

func TestBasic(t *testing.T) {
	tests.ValidateEnum(t, Basic(0), Basic_name, Basic_value)
}

func TestRename(t *testing.T) {
	tests.ValidateEnum(t, NewName(0), NewName_name, NewName_value)
}
