package gogoproto

import (
	patch_go "github.com/alta/protopatch/patch/go"
	"github.com/alta/protopatch/patch/options"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
)

func init() {
	options.RegisterProvider("gogoproto", New)
	options.RegisterProvider("gogo", New)
}

type provider struct {
	options.UnimplementedProvider
}

// New returns a GoGo shim as an options.Provider.
func New() options.Provider {
	return &provider{}
}

func (p *provider) EnumOptions(e *protogen.Enum) *patch_go.Options { return nil }

func (p *provider) ValueOptions(v *protogen.EnumValue) *patch_go.Options { return nil }

func (p *provider) MessageOptions(m *protogen.Message) *patch_go.Options { return nil }

func (p *provider) FieldOptions(f *protogen.Field) *patch_go.Options {
	name := proto.GetExtension(f.Desc.Options(), E_Customname).(string)
	if name != "" {
		return &patch_go.Options{
			Name: &name,
		}
	}
	return nil
}

func (p *provider) OneofOptions(o *protogen.Oneof) *patch_go.Options { return nil }
