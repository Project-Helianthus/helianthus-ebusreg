package system_test

import (
	"bytes"
	"errors"
	"fmt"
	"testing"

	ebuserrors "github.com/Project-Helianthus/helianthus-ebusgo/errors"
	"github.com/Project-Helianthus/helianthus-ebusgo/protocol"
	"github.com/Project-Helianthus/helianthus-ebusgo/types"
	"github.com/Project-Helianthus/helianthus-ebusreg/registry"
	"github.com/Project-Helianthus/helianthus-ebusreg/vaillant/system"
)

type k08RequestPlane interface {
	BuildRequest(registry.Method, map[string]any) (protocol.Frame, error)
	DecodeResponse(registry.Method, protocol.Frame, map[string]any) (any, error)
}

func k08RegisteredRead(t *testing.T) (k08RequestPlane, registry.Method) {
	t.Helper()
	r := registry.NewDeviceRegistry(nil)
	r.RegisterProvider(system.NewProvider())
	r.Register(registry.DeviceInfo{Manufacturer: "Vaillant", Address: 0x15})
	entry, ok := r.Lookup(0x15)
	if !ok {
		t.Fatal("registered synthetic device missing")
	}
	for _, plane := range entry.Planes() {
		if plane.Name() != "system" {
			continue
		}
		for _, method := range plane.Methods() {
			if method.Name() != "get_ext_register" {
				continue
			}
			if !method.ReadOnly() {
				t.Fatal("discovery read method is not read-only")
			}
			p, ok := plane.(k08RequestPlane)
			if !ok {
				t.Fatal("registered system plane lacks request/response API")
			}
			return p, method
		}
	}
	t.Fatal("registered B524 read method missing")
	return nil, nil
}

func TestK08DiscoveryFullKeysThroughRegistration(t *testing.T) {
	t.Parallel()
	plane, method := k08RegisteredRead(t)
	type key struct{ opcode, group byte }
	profiles := make(map[key]system.B524GroupProfile)
	for _, profile := range system.B524DiscoveryProfiles {
		k := key{profile.Opcode, profile.Group}
		if _, exists := profiles[k]; exists {
			t.Fatalf("duplicate full discovery key OP%02X/GG%02X", k.opcode, k.group)
		}
		if k.opcode != system.B524OpcodeLocal && k.opcode != system.B524OpcodeRemote {
			t.Fatalf("unsupported discovery opcode %02X", k.opcode)
		}
		profiles[k] = profile
	}
	// These are exploration endpoints, not qualified bounds or write authority.
	cases := []system.B524GroupProfile{
		{Group: 0x08, Opcode: system.B524OpcodeLocal, InstanceMax: 0x00, RegisterMax: 0x0007},
		{Group: 0x08, Opcode: system.B524OpcodeRemote, InstanceMax: 0x0A, RegisterMax: 0x0004},
		{Group: 0x09, Opcode: system.B524OpcodeLocal, InstanceMax: 0x0A, RegisterMax: 0x000F},
		{Group: 0x09, Opcode: system.B524OpcodeRemote, InstanceMax: 0x0A, RegisterMax: 0x002F},
		{Group: 0x0A, Opcode: system.B524OpcodeLocal, InstanceMax: 0x0A, RegisterMax: 0x004D},
		{Group: 0x0A, Opcode: system.B524OpcodeRemote, InstanceMax: 0x0A, RegisterMax: 0x003F},
	}
	for _, want := range cases {
		t.Run(fmt.Sprintf("OP%02X/GG%02X", want.Opcode, want.Group), func(t *testing.T) {
			got, ok := profiles[key{want.Opcode, want.Group}]
			if !ok || got != want {
				t.Fatalf("discovery profile = %+v, present %v; want %+v", got, ok, want)
			}
			for _, selector := range []string{"opcode", "op"} {
				params := map[string]any{"source": byte(0x31), selector: got.Opcode, "group": got.Group, "instance": got.InstanceMax, "addr": got.RegisterMax}
				request, err := plane.BuildRequest(method, params)
				if err != nil {
					t.Fatal(err)
				}
				expected := []byte{want.Opcode, 0x00, want.Group, want.InstanceMax, byte(want.RegisterMax), byte(want.RegisterMax >> 8)}
				if request.Primary != 0xB5 || request.Secondary != 0x24 || !bytes.Equal(request.Data, expected) {
					t.Fatalf("%s request = %+v; want B524 % X", selector, request, expected)
				}
			}
		})
	}
	// Identical GG/II/RR and reply bytes must still retain distinct opcode identity.
	for _, group := range []byte{0x08, 0x09, 0x0A} {
		for _, opcode := range []byte{system.B524OpcodeLocal, system.B524OpcodeRemote} {
			params := map[string]any{"source": byte(0x31), "opcode": opcode, "group": group, "instance": byte(0), "addr": uint16(1)}
			request, err := plane.BuildRequest(method, params)
			if err != nil {
				t.Fatal(err)
			}
			if request.Data[0] != opcode {
				t.Fatalf("opcode alias in request: % X", request.Data)
			}
			decoded, err := plane.DecodeResponse(method, protocol.Frame{Primary: 0xB5, Secondary: 0x24, Data: []byte{0x01, group, 0x01, 0x00, 0x2A}}, params)
			if err != nil {
				t.Fatal(err)
			}
			values, ok := decoded.(map[string]types.Value)
			if !ok {
				t.Fatalf("decoded type = %T", decoded)
			}
			for name, want := range map[string]any{"opcode": opcode, "group": group, "instance": byte(0), "addr": uint16(1)} {
				if got := values[name]; !got.Valid || got.Value != want {
					t.Fatalf("%s = %+v; want valid %v", name, got, want)
				}
			}
		}
	}
}

func TestK08RegisteredReadRejectsUnsupportedOpcode(t *testing.T) {
	t.Parallel()
	plane, method := k08RegisteredRead(t)
	for _, selector := range []string{"opcode", "op"} {
		for _, opcode := range []byte{0x00, 0x01, 0x03, 0x04, 0xFF} {
			t.Run(fmt.Sprintf("%s/%02X", selector, opcode), func(t *testing.T) {
				params := map[string]any{"source": byte(0x31), selector: opcode, "group": byte(0x08), "instance": byte(0), "addr": uint16(1)}
				if _, err := plane.BuildRequest(method, params); !errors.Is(err, ebuserrors.ErrInvalidPayload) {
					t.Fatalf("BuildRequest error = %v; want invalid payload", err)
				}
				if _, err := plane.DecodeResponse(method, protocol.Frame{Primary: 0xB5, Secondary: 0x24, Data: []byte{0x01, 0x08, 0x01, 0x00, 0x2A}}, params); !errors.Is(err, ebuserrors.ErrInvalidPayload) {
					t.Fatalf("DecodeResponse error = %v; want invalid payload", err)
				}
			})
		}
	}
}
