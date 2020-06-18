package message

import (
	"testing"

	"github.com/alta/protopatch/tests"
)

func TestBasic(t *testing.T) {
	tests.ValidateMessage(t, &Basic{})
}

func TestNested(t *testing.T) {
	tests.ValidateMessage(t, &Outer{})
	tests.ValidateMessage(t, &Outer_Inner{})
}

func TestRename(t *testing.T) {
	tests.ValidateMessage(t, &NewName{})
}
