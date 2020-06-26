package options

import (
	"sync"
)

// RegisterProvider registers a named Provider.
// Registering an already-registered name will overwrite the previously registered Provider.
func RegisterProvider(name string, constructor func() Provider) {
	mu.Lock()
	defer mu.Unlock()
	providers[name] = constructor
}

// NewProvider returns an initialized Provider for name.
// Returns nil if the named Provider is not registered.
func NewProvider(name string) Provider {
	mu.Lock()
	defer mu.Unlock()
	if constructor, ok := providers[name]; ok {
		return constructor()
	}
	return nil
}

var (
	mu        sync.Mutex
	providers = make(map[string]func() Provider)
)
