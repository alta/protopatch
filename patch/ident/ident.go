package ident

import (
	"strings"

	"google.golang.org/protobuf/compiler/protogen"
)

// WithChild returns id with a dotted child on its name part.
// Example: Message -> Message.Field
func WithChild(id protogen.GoIdent, child string) protogen.GoIdent {
	return WithSuffix(id, "."+child)
}

// WithPrefix returns id with a prefix on its name part.
// Example: Message -> PrefixedMessage
func WithPrefix(id protogen.GoIdent, prefix string) protogen.GoIdent {
	id.GoName = prefix + id.GoName
	return id
}

// WithSuffix returns id with a suffix on its name part.
// Example: Message becomes MessageSuffix
func WithSuffix(id protogen.GoIdent, suffix string) protogen.GoIdent {
	id.GoName = id.GoName + suffix
	return id
}

// LeafName returns the leaf name (after the last dot, if any) in id.
func LeafName(id protogen.GoIdent) string {
	names := strings.Split(id.GoName, ".")
	return names[len(names)-1]
}
