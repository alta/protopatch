package enum

import "strings"

// String returns a lower cased representation of the enum value.
func (cs CustomStringerEnum) String() string {
	return strings.ToLower(cs.OrigString())
}
