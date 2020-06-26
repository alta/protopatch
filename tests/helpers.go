package tests

import (
	"reflect"
	"testing"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// EnumNames is a type alias for map[int32]string.
type EnumNames = map[int32]string

// EnumValues is a type alias for map[string]int32.
type EnumValues = map[string]int32

// ValidateEnum performs basic validation of an enum.
func ValidateEnum(t *testing.T, e protoreflect.Enum, names EnumNames, values EnumValues) {
	for name, value := range values {
		if value != values[name] {
			t.Errorf("enum %s: %v != names[%s] (%v)", reflect.TypeOf(e), value, name, values[name])
		}
	}
}

// ValidateMessage performs basic validation of a message.
func ValidateMessage(t *testing.T, m proto.Message) {
	// TODO: add some validation
}

// ValidateTag performs basic validation of a struct tag.
func ValidateTag(t *testing.T, m proto.Message, field, tag, value string) {
	f, ok := reflect.TypeOf(m).Elem().FieldByName(field)
	if !ok {
		t.Errorf("%T: expected struct tag, but none found", m)
		return
	}
	if got := f.Tag.Get(tag); got != value {
		t.Errorf("%T.%s tag `%s` = %q, expected %q", m, field, tag, got, value)
	}
}
