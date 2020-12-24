package validate

import (
	"testing"

	"github.com/alta/protopatch/tests"
)

var (
	loopback4 = &IPAddress{Address: &IPAddress_IPV4{"127.0.0.1"}}
	loopback6 = &IPAddress{Address: &IPAddress_IPV6{"0:0:0:0:0:0:0:1"}}
	bogus4    = &IPAddress{Address: &IPAddress_IPV4{"999.999.999.999"}}
	bogus6    = &IPAddress{Address: &IPAddress_IPV6{"not.an.ip.address"}}
)

func TestInterfaceStatus(t *testing.T) {
	tests.ValidateEnum(t, InterfaceStatus(0), InterfaceStatus_name, InterfaceStatus_value)
	if got, want := StatusUnknown, InterfaceStatus(0); got != want {
		t.Errorf("%T(%d) != %v", got, got, want)
	}
	if got, want := StatusUp, InterfaceStatus(1); got != want {
		t.Errorf("%T(%d) != %v", got, got, want)
	}
	if got, want := StatusDown, InterfaceStatus(2); got != want {
		t.Errorf("%T(%d) != %v", got, got, want)
	}
}

func TestInterfaceValidate(t *testing.T) {
	tests := []struct {
		name    string
		i       *Interface
		wantErr bool
	}{
		{"nil", nil, false}, // Weird, but OK
		{"unknown", &Interface{Status: StatusUnknown}, true},
		{"up", &Interface{Status: StatusUp, Addresses: nil}, false},
		{"down", &Interface{Status: StatusDown, Addresses: nil}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.i.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestAddressValidate(t *testing.T) {
	tests := []struct {
		name    string
		ip      *IPAddress
		wantErr bool
	}{
		{"nil", nil, false}, // Weird, but OK
		{"loopback IPv4", loopback4, false},
		{"loopback IPv6", loopback4, false},
		{"bogus IPv4", bogus4, true},
		{"bogus IPv6", bogus6, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.ip.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
