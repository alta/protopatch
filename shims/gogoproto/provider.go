package gogoproto

import (
	"github.com/alta/protopatch/patch/options"
)

type provider struct {
	options.UnimplementedProvider
}

// New returns a GoGo shim as an options.Provider.
func New() options.Provider {
	return &provider{}
}

func init() {
	options.RegisterProvider("gogoproto", New)
	options.RegisterProvider("gogo", New)
}
