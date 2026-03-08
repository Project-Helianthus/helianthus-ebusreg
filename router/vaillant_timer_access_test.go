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

func TestVaillantSystem_ReadTimer_RequestEncodesPayload(t *testing.T) {
	t.Parallel()

	planes := system.NewProvider().CreatePlanes(registry.DeviceInfo{Manufacturer: "Vaillant", Address: 0x15})
	plane := planes[0].(router.Plane)

	bus := &vaillantMockBus{
		response: &protocol.Frame{
			Source:    0x15,
			Target:    0x71,
			Primary:   0xB5,
			Secondary: 0x24,
			Data:      []byte{0x01, 0x02, 0x03},
		},
	}
	eventRouter := router.NewBusEventRouter(bus)

	_, err := eventRouter.Invoke(context.Background(), plane, "read_timer", map[string]any{
		"source":  byte(0x71),
		"sel1":    byte(0x00),
		"sel2":    byte(0x01),
		"sel3":    byte(0x00),
		"weekday": byte(0x02),
	})
	if err != nil {
		t.Fatalf("Invoke error = %v", err)
	}

	want := []byte{0x03, 0x00, 0x01, 0x00, 0x02}
	if !bytes.Equal(bus.lastRequest.Data, want) {
		t.Fatalf("request data = %v; want %v", bus.lastRequest.Data, want)
	}
	if bus.lastRequest.Source != 0x71 {
		t.Fatalf("source = 0x%02X; want 0x71", bus.lastRequest.Source)
	}
	if bus.lastRequest.Target != 0x15 {
		t.Fatalf("target = 0x%02X; want 0x15", bus.lastRequest.Target)
	}
	if bus.lastRequest.Primary != 0xB5 || bus.lastRequest.Secondary != 0x24 {
		t.Fatalf("primary/secondary = 0x%02X/0x%02X; want 0xB5/0x24",
			bus.lastRequest.Primary, bus.lastRequest.Secondary)
	}
}

func TestVaillantSystem_ReadTimer_ResponseDecode(t *testing.T) {
	t.Parallel()

	planes := system.NewProvider().CreatePlanes(registry.DeviceInfo{Manufacturer: "Vaillant", Address: 0x15})
	plane := planes[0].(router.Plane)

	t.Run("normal schedule payload", func(t *testing.T) {
		t.Parallel()

		timerData := []byte{0x19, 0x1E, 0x23, 0x28, 0x2D, 0x2D, 0x2D}
		bus := &vaillantMockBus{
			response: &protocol.Frame{
				Source:    0x15,
				Target:    0x71,
				Primary:   0xB5,
				Secondary: 0x24,
				Data:      timerData,
			},
		}
		eventRouter := router.NewBusEventRouter(bus)

		result, err := eventRouter.Invoke(context.Background(), plane, "read_timer", map[string]any{
			"source":  byte(0x71),
			"sel1":    byte(0x00),
			"sel2":    byte(0x01),
			"sel3":    byte(0x00),
			"weekday": byte(0x00),
		})
		if err != nil {
			t.Fatalf("Invoke error = %v", err)
		}

		values := result.(map[string]types.Value)
		if got := values["opcode"]; !got.Valid || got.Value != byte(0x03) {
			t.Fatalf("opcode = %+v; want 0x03 valid", got)
		}
		if got := values["sel1"]; !got.Valid || got.Value != byte(0x00) {
			t.Fatalf("sel1 = %+v; want 0x00 valid", got)
		}
		if got := values["sel2"]; !got.Valid || got.Value != byte(0x01) {
			t.Fatalf("sel2 = %+v; want 0x01 valid", got)
		}
		if got := values["sel3"]; !got.Valid || got.Value != byte(0x00) {
			t.Fatalf("sel3 = %+v; want 0x00 valid", got)
		}
		if got := values["weekday"]; !got.Valid || got.Value != byte(0x00) {
			t.Fatalf("weekday = %+v; want 0x00 valid", got)
		}
		if got := values["value"]; !got.Valid || !bytes.Equal(got.Value.([]byte), timerData) {
			t.Fatalf("value = %+v; want %v valid", got, timerData)
		}
		if got := values["slot_count"]; !got.Valid || got.Value != 7 {
			t.Fatalf("slot_count = %+v; want 7 valid", got)
		}
	})

	t.Run("empty response", func(t *testing.T) {
		t.Parallel()

		bus := &vaillantMockBus{
			response: &protocol.Frame{
				Source:    0x15,
				Target:    0x71,
				Primary:   0xB5,
				Secondary: 0x24,
				Data:      nil,
			},
		}
		eventRouter := router.NewBusEventRouter(bus)

		result, err := eventRouter.Invoke(context.Background(), plane, "read_timer", map[string]any{
			"source":  byte(0x71),
			"sel1":    byte(0x00),
			"sel2":    byte(0x00),
			"sel3":    byte(0x00),
			"weekday": byte(0x03),
		})
		if err != nil {
			t.Fatalf("Invoke error = %v", err)
		}

		values := result.(map[string]types.Value)
		if got := values["value"]; got.Valid {
			t.Fatalf("value = %+v; want invalid for empty response", got)
		}
	})
}

func TestVaillantSystem_ReadTimer_RejectsInvalidWeekday(t *testing.T) {
	t.Parallel()

	planes := system.NewProvider().CreatePlanes(registry.DeviceInfo{Manufacturer: "Vaillant", Address: 0x15})
	plane := planes[0].(router.Plane)

	bus := &vaillantMockBus{
		response: &protocol.Frame{
			Source:    0x15,
			Target:    0x71,
			Primary:   0xB5,
			Secondary: 0x24,
			Data:      nil,
		},
	}
	eventRouter := router.NewBusEventRouter(bus)

	_, err := eventRouter.Invoke(context.Background(), plane, "read_timer", map[string]any{
		"source":  byte(0x71),
		"sel1":    byte(0x00),
		"sel2":    byte(0x01),
		"sel3":    byte(0x00),
		"weekday": byte(0x07),
	})
	if err == nil {
		t.Fatal("expected error for weekday=7; got nil")
	}
}

func TestVaillantSystem_ReadTimer_RejectsMissingParams(t *testing.T) {
	t.Parallel()

	planes := system.NewProvider().CreatePlanes(registry.DeviceInfo{Manufacturer: "Vaillant", Address: 0x15})
	plane := planes[0].(router.Plane)

	bus := &vaillantMockBus{
		response: &protocol.Frame{
			Source:    0x15,
			Target:    0x71,
			Primary:   0xB5,
			Secondary: 0x24,
			Data:      nil,
		},
	}
	eventRouter := router.NewBusEventRouter(bus)

	for _, missing := range []string{"sel1", "sel2", "sel3", "weekday"} {
		t.Run("missing_"+missing, func(t *testing.T) {
			params := map[string]any{
				"source":  byte(0x71),
				"sel1":    byte(0x00),
				"sel2":    byte(0x01),
				"sel3":    byte(0x00),
				"weekday": byte(0x00),
			}
			delete(params, missing)

			_, err := eventRouter.Invoke(context.Background(), plane, "read_timer", params)
			if err == nil {
				t.Fatalf("expected error for missing %s; got nil", missing)
			}
		})
	}
}

func TestVaillantSystem_ReadRaw_RequestPassthroughPayload(t *testing.T) {
	t.Parallel()

	planes := system.NewProvider().CreatePlanes(registry.DeviceInfo{Manufacturer: "Vaillant", Address: 0x15})
	plane := planes[0].(router.Plane)

	rawPayload := []byte{0x0B, 0x06, 0x00, 0x01}
	bus := &vaillantMockBus{
		response: &protocol.Frame{
			Source:    0x15,
			Target:    0x71,
			Primary:   0xB5,
			Secondary: 0x24,
			Data:      []byte{0xAA, 0xBB, 0xCC},
		},
	}
	eventRouter := router.NewBusEventRouter(bus)

	_, err := eventRouter.Invoke(context.Background(), plane, "read_raw", map[string]any{
		"source":  byte(0x71),
		"payload": rawPayload,
	})
	if err != nil {
		t.Fatalf("Invoke error = %v", err)
	}

	if !bytes.Equal(bus.lastRequest.Data, rawPayload) {
		t.Fatalf("request data = %v; want %v", bus.lastRequest.Data, rawPayload)
	}
}

func TestVaillantSystem_ReadRaw_ResponseDecode(t *testing.T) {
	t.Parallel()

	planes := system.NewProvider().CreatePlanes(registry.DeviceInfo{Manufacturer: "Vaillant", Address: 0x15})
	plane := planes[0].(router.Plane)

	responseData := []byte{0xAA, 0xBB, 0xCC, 0xDD}
	rawPayload := []byte{0x0B, 0x06, 0x00}
	bus := &vaillantMockBus{
		response: &protocol.Frame{
			Source:    0x15,
			Target:    0x71,
			Primary:   0xB5,
			Secondary: 0x24,
			Data:      responseData,
		},
	}
	eventRouter := router.NewBusEventRouter(bus)

	result, err := eventRouter.Invoke(context.Background(), plane, "read_raw", map[string]any{
		"source":  byte(0x71),
		"payload": rawPayload,
	})
	if err != nil {
		t.Fatalf("Invoke error = %v", err)
	}

	values := result.(map[string]types.Value)
	if got := values["request_payload"]; !got.Valid || !bytes.Equal(got.Value.([]byte), rawPayload) {
		t.Fatalf("request_payload = %+v; want %v valid", got, rawPayload)
	}
	if got := values["response_payload"]; !got.Valid || !bytes.Equal(got.Value.([]byte), responseData) {
		t.Fatalf("response_payload = %+v; want %v valid", got, responseData)
	}
	if got := values["value"]; !got.Valid || !bytes.Equal(got.Value.([]byte), responseData) {
		t.Fatalf("value = %+v; want %v valid", got, responseData)
	}
}

func TestVaillantSystem_ReadRaw_RejectsEmptyPayload(t *testing.T) {
	t.Parallel()

	planes := system.NewProvider().CreatePlanes(registry.DeviceInfo{Manufacturer: "Vaillant", Address: 0x15})
	plane := planes[0].(router.Plane)

	bus := &vaillantMockBus{
		response: &protocol.Frame{
			Source:    0x15,
			Target:    0x71,
			Primary:   0xB5,
			Secondary: 0x24,
			Data:      nil,
		},
	}
	eventRouter := router.NewBusEventRouter(bus)

	_, err := eventRouter.Invoke(context.Background(), plane, "read_raw", map[string]any{
		"source":  byte(0x71),
		"payload": []byte{},
	})
	if err == nil {
		t.Fatal("expected error for empty payload; got nil")
	}
}

func TestVaillantSystem_ReadRaw_RejectsOversizedPayload(t *testing.T) {
	t.Parallel()

	planes := system.NewProvider().CreatePlanes(registry.DeviceInfo{Manufacturer: "Vaillant", Address: 0x15})
	plane := planes[0].(router.Plane)

	bus := &vaillantMockBus{
		response: &protocol.Frame{
			Source:    0x15,
			Target:    0x71,
			Primary:   0xB5,
			Secondary: 0x24,
			Data:      nil,
		},
	}
	eventRouter := router.NewBusEventRouter(bus)

	oversized := make([]byte, 17)
	_, err := eventRouter.Invoke(context.Background(), plane, "read_raw", map[string]any{
		"source":  byte(0x71),
		"payload": oversized,
	})
	if err == nil {
		t.Fatal("expected error for oversized payload; got nil")
	}
}
