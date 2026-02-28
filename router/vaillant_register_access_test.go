package router_test

import (
	"bytes"
	"context"
	"testing"

	"github.com/Project-Helianthus/helianthus-ebusgo/protocol"
	"github.com/Project-Helianthus/helianthus-ebusgo/types"
	"github.com/Project-Helianthus/helianthus-ebusreg/registry"
	"github.com/Project-Helianthus/helianthus-ebusreg/router"
	"github.com/Project-Helianthus/helianthus-ebusreg/vaillant/system"
)

func TestVaillantSystem_GetRegister_RequestEncodesAddr(t *testing.T) {
	t.Parallel()

	planes := system.NewProvider().CreatePlanes(registry.DeviceInfo{Manufacturer: "Vaillant", Address: 0x08})
	plane := planes[0].(router.Plane)

	cases := []struct {
		name string
		addr any
	}{
		{name: "uint16", addr: uint16(0xF600)},
		{name: "hex string", addr: "F600"},
		{name: "0x-prefixed hex string", addr: "0xF600"},
		{name: "bytes", addr: []byte{0xF6, 0x00}},
		{name: "ints", addr: []int{0xF6, 0x00}},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			bus := &vaillantMockBus{
				response: &protocol.Frame{
					Source:    0x08,
					Target:    0x10,
					Primary:   0xB5,
					Secondary: 0x09,
					Data:      []byte{0x12, 0x34},
				},
			}
			eventRouter := router.NewBusEventRouter(bus)

			_, err := eventRouter.Invoke(context.Background(), plane, "get_register", map[string]any{
				"source": byte(0x10),
				"addr":   tc.addr,
			})
			if err != nil {
				t.Fatalf("Invoke error = %v", err)
			}

			if bus.lastRequest.Primary != 0xB5 || bus.lastRequest.Secondary != 0x09 {
				t.Fatalf("unexpected request type: %+v", bus.lastRequest)
			}
			if !bytes.Equal(bus.lastRequest.Data, []byte{0x0D, 0xF6, 0x00}) {
				t.Fatalf("unexpected request data %v; want [0d f6 00]", bus.lastRequest.Data)
			}
		})
	}
}

func TestVaillantSystem_GetRegister_ResponseDecode(t *testing.T) {
	t.Parallel()

	planes := system.NewProvider().CreatePlanes(registry.DeviceInfo{Manufacturer: "Vaillant", Address: 0x08})
	plane := planes[0].(router.Plane)

	t.Run("normal payload", func(t *testing.T) {
		t.Parallel()

		bus := &vaillantMockBus{
			response: &protocol.Frame{
				Source:    0x08,
				Target:    0x10,
				Primary:   0xB5,
				Secondary: 0x09,
				Data:      []byte{0x12, 0x34},
			},
		}
		eventRouter := router.NewBusEventRouter(bus)

		result, err := eventRouter.Invoke(context.Background(), plane, "get_register", map[string]any{
			"source": byte(0x10),
			"addr":   uint16(0xF600),
		})
		if err != nil {
			t.Fatalf("Invoke error = %v", err)
		}

		values := result.(map[string]types.Value)
		if got := values["cmd"]; !got.Valid || got.Value != byte(0x0D) {
			t.Fatalf("cmd = %+v; want 0x0D valid", got)
		}
		if got := values["addr"]; !got.Valid || got.Value != uint16(0xF600) {
			t.Fatalf("addr = %+v; want 0xF600 valid", got)
		}
		if got := values["addr_hex"]; !got.Valid || got.Value != "F600" {
			t.Fatalf("addr_hex = %+v; want F600 valid", got)
		}
		if got := values["payload"]; !got.Valid || !bytes.Equal(got.Value.([]byte), []byte{0x12, 0x34}) {
			t.Fatalf("payload = %+v; want [12 34] valid", got)
		}
		if got := values["value"]; !got.Valid || !bytes.Equal(got.Value.([]byte), []byte{0x12, 0x34}) {
			t.Fatalf("value = %+v; want [12 34] valid", got)
		}
	})

	t.Run("short zero payload", func(t *testing.T) {
		t.Parallel()

		bus := &vaillantMockBus{
			response: &protocol.Frame{
				Source:    0x08,
				Target:    0x10,
				Primary:   0xB5,
				Secondary: 0x09,
				Data:      []byte{0x00},
			},
		}
		eventRouter := router.NewBusEventRouter(bus)

		result, err := eventRouter.Invoke(context.Background(), plane, "get_register", map[string]any{
			"source": byte(0x10),
			"addr":   uint16(0xF600),
		})
		if err != nil {
			t.Fatalf("Invoke error = %v", err)
		}

		values := result.(map[string]types.Value)
		if got := values["payload"]; !got.Valid || !bytes.Equal(got.Value.([]byte), []byte{0x00}) {
			t.Fatalf("payload = %+v; want [00] valid", got)
		}
		if got := values["value"]; got.Valid {
			t.Fatalf("value = %+v; want invalid", got)
		}
	})
}

func TestVaillantSystem_SetRegister_SendsPayload(t *testing.T) {
	t.Parallel()

	planes := system.NewProvider().CreatePlanes(registry.DeviceInfo{Manufacturer: "Vaillant", Address: 0x08})
	plane := planes[0].(router.Plane)

	bus := &vaillantMockBus{
		response: &protocol.Frame{
			Source:    0x08,
			Target:    0x10,
			Primary:   0xB5,
			Secondary: 0x09,
			Data:      nil,
		},
	}
	eventRouter := router.NewBusEventRouter(bus)

	_, err := eventRouter.Invoke(context.Background(), plane, "set_register", map[string]any{
		"source": byte(0x10),
		"addr":   "F600",
		"data":   []byte{0xAA, 0xBB},
	})
	if err != nil {
		t.Fatalf("Invoke error = %v", err)
	}

	if bus.lastRequest.Primary != 0xB5 || bus.lastRequest.Secondary != 0x09 {
		t.Fatalf("unexpected request type: %+v", bus.lastRequest)
	}
	if !bytes.Equal(bus.lastRequest.Data, []byte{0x0E, 0xF6, 0x00, 0xAA, 0xBB}) {
		t.Fatalf("unexpected request data %v; want [0e f6 00 aa bb]", bus.lastRequest.Data)
	}
}
