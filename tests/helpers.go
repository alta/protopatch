package tests

import (
	"reflect"
	"testing"

	"google.golang.org/protobuf/reflect/protoreflect"
)

// EnumNames is a type alias for map[int32]string.
type EnumNames = map[int32]string

// EnumValues is a type alias for map[string]int32.
type EnumValues = map[string]int32

// ValidateEnum validates an enum type.
func ValidateEnum(t *testing.T, e protoreflect.Enum, names EnumNames, values EnumValues) {
	for name, value := range values {
		if value != values[name] {
			t.Errorf("enum %s: %v != names[%s] (%v)", reflect.TypeOf(e), value, name, values[name])
		}
	}
}
