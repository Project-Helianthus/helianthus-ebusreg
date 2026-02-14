package router_test

import (
	"bytes"
	"context"
	"testing"

	"github.com/d3vi1/helianthus-ebusgo/protocol"
	"github.com/d3vi1/helianthus-ebusgo/types"
	"github.com/d3vi1/helianthus-ebusreg/registry"
	"github.com/d3vi1/helianthus-ebusreg/router"
	"github.com/d3vi1/helianthus-ebusreg/vaillant/system"
)

func TestVaillantSystem_GetExtRegister_RequestEncodesPayload(t *testing.T) {
	t.Parallel()

	planes := system.NewProvider().CreatePlanes(registry.DeviceInfo{Manufacturer: "Vaillant", Address: 0x08})
	plane := planes[0].(router.Plane)

	bus := &vaillantMockBus{
		response: &protocol.Frame{
			Source:    0x08,
			Target:    0x10,
			Primary:   0xB5,
			Secondary: 0x24,
			Data:      []byte{0x01, 0x00, 0x00, 0x5C, 0xAA},
		},
	}
	eventRouter := router.NewBusEventRouter(bus)

	_, err := eventRouter.Invoke(context.Background(), plane, "get_ext_register", map[string]any{
		"source":   byte(0x10),
		"group":    byte(0x00),
		"instance": byte(0x00),
		"addr":     uint16(0x5C00),
	})
	if err != nil {
		t.Fatalf("Invoke error = %v", err)
	}

	if !bytes.Equal(bus.lastRequest.Data, []byte{0x02, 0x00, 0x00, 0x00, 0x00, 0x5C}) {
		t.Fatalf("unexpected request data %v; want [02 00 00 00 00 5c]", bus.lastRequest.Data)
	}
}

func TestVaillantSystem_GetExtRegister_ResponseDecode(t *testing.T) {
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
				Secondary: 0x24,
				Data:      []byte{0x01, 0x00, 0x00, 0x5C, 0xAA, 0xBB},
			},
		}
		eventRouter := router.NewBusEventRouter(bus)

		result, err := eventRouter.Invoke(context.Background(), plane, "get_ext_register", map[string]any{
			"source":   byte(0x10),
			"group":    byte(0x00),
			"instance": byte(0x00),
			"addr":     uint16(0x5C00),
		})
		if err != nil {
			t.Fatalf("Invoke error = %v", err)
		}

		values := result.(map[string]types.Value)
		if got := values["group"]; !got.Valid || got.Value != byte(0x00) {
			t.Fatalf("group = %+v; want 0x00 valid", got)
		}
		if got := values["instance"]; !got.Valid || got.Value != byte(0x00) {
			t.Fatalf("instance = %+v; want 0x00 valid", got)
		}
		if got := values["addr"]; !got.Valid || got.Value != uint16(0x5C00) {
			t.Fatalf("addr = %+v; want 0x5C00 valid", got)
		}
		if got := values["addr_hex"]; !got.Valid || got.Value != "5C00" {
			t.Fatalf("addr_hex = %+v; want 5C00 valid", got)
		}
		if got := values["prefix"]; !got.Valid || !bytes.Equal(got.Value.([]byte), []byte{0x01, 0x00, 0x00, 0x5C}) {
			t.Fatalf("prefix = %+v; want [01 00 00 5c] valid", got)
		}
		if got := values["reply_kind"]; !got.Valid || got.Value != byte(0x01) {
			t.Fatalf("reply_kind = %+v; want 0x01 valid", got)
		}
		if got := values["reply_group"]; !got.Valid || got.Value != byte(0x00) {
			t.Fatalf("reply_group = %+v; want 0x00 valid", got)
		}
		if got := values["reply_addr"]; !got.Valid || got.Value != uint16(0x5C00) {
			t.Fatalf("reply_addr = %+v; want 0x5C00 valid", got)
		}
		if got := values["reply_addr_hex"]; !got.Valid || got.Value != "5C00" {
			t.Fatalf("reply_addr_hex = %+v; want 5C00 valid", got)
		}
		if got := values["value"]; !got.Valid || !bytes.Equal(got.Value.([]byte), []byte{0xAA, 0xBB}) {
			t.Fatalf("value = %+v; want [aa bb] valid", got)
		}
		if got, ok := values["constraints"]; ok {
			t.Fatalf("constraints = %+v; want absent for unknown mapping", got)
		}
	})

	t.Run("known mapping adds constraints", func(t *testing.T) {
		t.Parallel()

		bus := &vaillantMockBus{
			response: &protocol.Frame{
				Source:    0x08,
				Target:    0x10,
				Primary:   0xB5,
				Secondary: 0x24,
				Data:      []byte{0x01, 0x01, 0x00, 0x04, 0x2A},
			},
		}
		eventRouter := router.NewBusEventRouter(bus)

		result, err := eventRouter.Invoke(context.Background(), plane, "get_ext_register", map[string]any{
			"source":   byte(0x10),
			"group":    byte(0x01),
			"instance": byte(0x01),
			"addr":     uint16(0x0400),
		})
		if err != nil {
			t.Fatalf("Invoke error = %v", err)
		}

		values := result.(map[string]types.Value)
		constraintsValue, ok := values["constraints"]
		if !ok || !constraintsValue.Valid {
			t.Fatalf("constraints = %+v; want valid", constraintsValue)
		}

		constraints, ok := constraintsValue.Value.(map[string]any)
		if !ok {
			t.Fatalf("constraints type = %T; want map[string]any", constraintsValue.Value)
		}
		if constraints["type"] != "f32_range" || constraints["min"] != float64(35) || constraints["max"] != float64(70) || constraints["step"] != float64(1) {
			t.Fatalf("constraints map = %#v; want f32_range/min=35/max=70/step=1", constraints)
		}
		if got := values["constraint_record"]; !got.Valid || got.Value != uint16(0x0400) {
			t.Fatalf("constraint_record = %+v; want 0x0400 valid", got)
		}
		if got := values["constraint_record_hex"]; !got.Valid || got.Value != "0400" {
			t.Fatalf("constraint_record_hex = %+v; want 0400 valid", got)
		}
		if got := values["constraint_type"]; !got.Valid || got.Value != "f32_range" {
			t.Fatalf("constraint_type = %+v; want f32_range valid", got)
		}
		if got := values["constraint_min"]; !got.Valid || got.Value != float64(35) {
			t.Fatalf("constraint_min = %+v; want 35 valid", got)
		}
		if got := values["constraint_max"]; !got.Valid || got.Value != float64(70) {
			t.Fatalf("constraint_max = %+v; want 70 valid", got)
		}
		if got := values["constraint_step"]; !got.Valid || got.Value != float64(1) {
			t.Fatalf("constraint_step = %+v; want 1 valid", got)
		}
	})

	t.Run("echoed reply mapping overrides constraint record", func(t *testing.T) {
		t.Parallel()

		bus := &vaillantMockBus{
			response: &protocol.Frame{
				Source:    0x08,
				Target:    0x10,
				Primary:   0xB5,
				Secondary: 0x24,
				Data:      []byte{0x01, 0x03, 0x00, 0x05, 0x2A},
			},
		}
		eventRouter := router.NewBusEventRouter(bus)

		result, err := eventRouter.Invoke(context.Background(), plane, "get_ext_register", map[string]any{
			"source":   byte(0x10),
			"group":    byte(0x01),
			"instance": byte(0x01),
			"addr":     uint16(0x0400),
		})
		if err != nil {
			t.Fatalf("Invoke error = %v", err)
		}

		values := result.(map[string]types.Value)
		if got := values["constraint_record"]; !got.Valid || got.Value != uint16(0x0500) {
			t.Fatalf("constraint_record = %+v; want 0x0500 valid", got)
		}
		if got := values["constraint_record_hex"]; !got.Valid || got.Value != "0500" {
			t.Fatalf("constraint_record_hex = %+v; want 0500 valid", got)
		}
		if got := values["constraint_min"]; !got.Valid || got.Value != float64(5) {
			t.Fatalf("constraint_min = %+v; want 5 valid", got)
		}
		if got := values["constraint_max"]; !got.Valid || got.Value != float64(30) {
			t.Fatalf("constraint_max = %+v; want 30 valid", got)
		}
	})

	t.Run("short zero payload", func(t *testing.T) {
		t.Parallel()

		bus := &vaillantMockBus{
			response: &protocol.Frame{
				Source:    0x08,
				Target:    0x10,
				Primary:   0xB5,
				Secondary: 0x24,
				Data:      []byte{0x00},
			},
		}
		eventRouter := router.NewBusEventRouter(bus)

		result, err := eventRouter.Invoke(context.Background(), plane, "get_ext_register", map[string]any{
			"source":   byte(0x10),
			"group":    byte(0x00),
			"instance": byte(0x00),
			"addr":     uint16(0x5C00),
		})
		if err != nil {
			t.Fatalf("Invoke error = %v", err)
		}

		values := result.(map[string]types.Value)
		if got := values["value"]; got.Valid {
			t.Fatalf("value = %+v; want invalid", got)
		}
	})
}

func TestVaillantSystem_SetExtRegister_SendsPayload(t *testing.T) {
	t.Parallel()

	planes := system.NewProvider().CreatePlanes(registry.DeviceInfo{Manufacturer: "Vaillant", Address: 0x08})
	plane := planes[0].(router.Plane)

	bus := &vaillantMockBus{
		response: &protocol.Frame{
			Source:    0x08,
			Target:    0x10,
			Primary:   0xB5,
			Secondary: 0x24,
			Data:      nil,
		},
	}
	eventRouter := router.NewBusEventRouter(bus)

	_, err := eventRouter.Invoke(context.Background(), plane, "set_ext_register", map[string]any{
		"source":   byte(0x10),
		"group":    0,
		"instance": 0,
		"addr":     "5C00",
		"data":     []byte{0xAA, 0xBB},
	})
	if err != nil {
		t.Fatalf("Invoke error = %v", err)
	}

	if !bytes.Equal(bus.lastRequest.Data, []byte{0x02, 0x01, 0x00, 0x00, 0x00, 0x5C, 0xAA, 0xBB}) {
		t.Fatalf("unexpected request data %v; want [02 01 00 00 00 5c aa bb]", bus.lastRequest.Data)
	}
}
