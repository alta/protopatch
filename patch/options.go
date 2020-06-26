package patch

import (
	patch_go "github.com/alta/protopatch/patch/go"
	"github.com/alta/protopatch/patch/go/enum"
	"github.com/alta/protopatch/patch/go/field"
	"github.com/alta/protopatch/patch/go/message"
	"github.com/alta/protopatch/patch/go/oneof"
	"github.com/alta/protopatch/patch/go/value"
	"github.com/alta/protopatch/patch/options"

	"google.golang.org/protobuf/compiler/protogen"
)

// Ensure Patcher implements options.Provider.
var _ options.Provider = &Patcher{}

// EnumOptions returns patch_go.Options if present, otherwise nil.
func (p *Patcher) EnumOptions(e *protogen.Enum) *patch_go.Options {
	return options.Get(e.Desc, enum.E_Options)
}

// ValueOptions returns patch_go.Options if present, otherwise nil.
func (p *Patcher) ValueOptions(v *protogen.EnumValue) *patch_go.Options {
	return options.Get(v.Desc, value.E_Options)
}

// MessageOptions returns patch_go.Options if present, otherwise nil.
func (p *Patcher) MessageOptions(m *protogen.Message) *patch_go.Options {
	return options.Get(m.Desc, message.E_Options)
}

// FieldOptions returns patch_go.Options if present on e, otherwise nil.
func (p *Patcher) FieldOptions(f *protogen.Field) *patch_go.Options {
	// First try (go.field.options)
	if opts := options.Get(f.Desc, field.E_Options); opts != nil {
		return opts
	}
	// Then try (go.field.options)
	if opts := options.Get(f.Desc, patch_go.E_Options); opts != nil {
		return opts
	}
	return nil
}

// OneofOptions returns patch_go.Options if present on e, otherwise nil.
func (p *Patcher) OneofOptions(o *protogen.Oneof) *patch_go.Options {
	return options.Get(o.Desc, oneof.E_Options)
}
