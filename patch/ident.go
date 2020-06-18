package patch

import (
	"strings"

	"google.golang.org/protobuf/compiler/protogen"
)

// WithChild returns ident with a dotted child on its name part.
// Example: Message -> Message.Field
func WithChild(ident protogen.GoIdent, child string) protogen.GoIdent {
	return WithSuffix(ident, "."+child)
}

// WithPrefix returns ident with a prefix on its name part.
// Example: Message -> PrefixedMessage
func WithPrefix(ident protogen.GoIdent, prefix string) protogen.GoIdent {
	ident.GoName = prefix + ident.GoName
	return ident
}

// WithSuffix returns ident with a suffix on its name part.
// Example: Message becomes MessageSuffix
func WithSuffix(ident protogen.GoIdent, suffix string) protogen.GoIdent {
	ident.GoName = ident.GoName + suffix
	return ident
}

// LeafName returns the leaf name (after the last dot, if any) in ident.
func LeafName(ident protogen.GoIdent) string {
	names := strings.Split(ident.GoName, ".")
	return names[len(names)-1]
}
