package router_test

import (
	"bytes"
	"context"
	"errors"
	"testing"

	ebuserrors "github.com/d3vi1/helianthus-ebusgo/errors"
	"github.com/d3vi1/helianthus-ebusgo/protocol"
	"github.com/d3vi1/helianthus-ebusgo/types"
	"github.com/d3vi1/helianthus-ebusreg/registry"
	"github.com/d3vi1/helianthus-ebusreg/router"
	"github.com/d3vi1/helianthus-ebusreg/vaillant/system"
)

type vaillantMockBus struct {
	lastRequest protocol.Frame
	response    *protocol.Frame
	err         error
}

func (bus *vaillantMockBus) Send(ctx context.Context, frame protocol.Frame) (*protocol.Frame, error) {
	bus.lastRequest = frame
	return bus.response, bus.err
}

func TestVaillantSystem_GetOperationalData_DateTime(t *testing.T) {
	t.Parallel()

	planes := system.NewProvider().CreatePlanes(registry.DeviceInfo{Manufacturer: "Vaillant", Address: 0x08})
	if len(planes) != 1 {
		t.Fatalf("expected 1 plane, got %d", len(planes))
	}

	plane, ok := planes[0].(router.Plane)
	if !ok {
		t.Fatalf("expected system plane to implement router.Plane")
	}

	payload := []byte{
		0x01,       // dcfstate
		0x12,       // hour (BCD 12)
		0x34,       // minute (BCD 34)
		0x03,       // day (BCD 03)
		0x02,       // month (BCD 02)
		0x26,       // year (BCD 26)
		0x80, 0x14, // temp (DATA2b 20.5)
	}

	bus := &vaillantMockBus{
		response: &protocol.Frame{
			Source:    0x08,
			Target:    0x10,
			Primary:   0xB5,
			Secondary: 0x04,
			Data:      payload,
		},
	}
	eventRouter := router.NewBusEventRouter(bus)

	result, err := eventRouter.Invoke(context.Background(), plane, "get_operational_data", map[string]any{
		"source": byte(0x10),
		"op":     byte(0x00),
	})
	if err != nil {
		t.Fatalf("Invoke error = %v", err)
	}

	if bus.lastRequest.Source != 0x10 || bus.lastRequest.Target != 0x08 {
		t.Fatalf("unexpected request addresses: %+v", bus.lastRequest)
	}
	if bus.lastRequest.Primary != 0xB5 || bus.lastRequest.Secondary != 0x04 {
		t.Fatalf("unexpected request type: %+v", bus.lastRequest)
	}
	if !bytes.Equal(bus.lastRequest.Data, []byte{0x00}) {
		t.Fatalf("unexpected request data %v", bus.lastRequest.Data)
	}

	values, ok := result.(map[string]types.Value)
	if !ok {
		t.Fatalf("expected map[string]types.Value, got %T", result)
	}

	if op := values["op"]; !op.Valid || op.Value != byte(0x00) {
		t.Fatalf("op = %+v; want 0x00 valid", op)
	}
	if got := values["payload"]; !got.Valid || !bytes.Equal(got.Value.([]byte), payload) {
		t.Fatalf("payload = %+v; want %v", got, payload)
	}
	if got := values["dcfstate"]; !got.Valid || got.Value != uint8(0x01) {
		t.Fatalf("dcfstate = %+v; want 1 valid", got)
	}
	if got := values["time_hour"]; !got.Valid || got.Value != uint8(12) {
		t.Fatalf("time_hour = %+v; want 12 valid", got)
	}
	if got := values["time_minute"]; !got.Valid || got.Value != uint8(34) {
		t.Fatalf("time_minute = %+v; want 34 valid", got)
	}
	if got := values["date_day"]; !got.Valid || got.Value != uint8(3) {
		t.Fatalf("date_day = %+v; want 3 valid", got)
	}
	if got := values["date_month"]; !got.Valid || got.Value != uint8(2) {
		t.Fatalf("date_month = %+v; want 2 valid", got)
	}
	if got := values["date_year"]; !got.Valid || got.Value != uint8(26) {
		t.Fatalf("date_year = %+v; want 26 valid", got)
	}
	if got := values["temp"]; !got.Valid || got.Value != float64(20.5) {
		t.Fatalf("temp = %+v; want 20.5 valid", got)
	}
}

func TestVaillantSystem_GetOperationalData_MissingParams(t *testing.T) {
	t.Parallel()

	planes := system.NewProvider().CreatePlanes(registry.DeviceInfo{Manufacturer: "Vaillant", Address: 0x08})
	plane := planes[0].(router.Plane)

	eventRouter := router.NewBusEventRouter(&vaillantMockBus{})

	_, err := eventRouter.Invoke(context.Background(), plane, "get_operational_data", map[string]any{
		"source": byte(0x10),
	})
	if !errors.Is(err, ebuserrors.ErrInvalidPayload) {
		t.Fatalf("expected ErrInvalidPayload, got %v", err)
	}
}

func TestVaillantSystem_SetOperationalData_SendsPayload(t *testing.T) {
	t.Parallel()

	planes := system.NewProvider().CreatePlanes(registry.DeviceInfo{Manufacturer: "Vaillant", Address: 0x08})
	plane := planes[0].(router.Plane)

	bus := &vaillantMockBus{
		response: &protocol.Frame{
			Source:    0x08,
			Target:    0x10,
			Primary:   0xB5,
			Secondary: 0x05,
		},
	}
	eventRouter := router.NewBusEventRouter(bus)

	result, err := eventRouter.Invoke(context.Background(), plane, "set_operational_data", map[string]any{
		"source": byte(0x10),
		"op":     byte(0x02),
		"data":   []byte{0xAA, 0xBB},
	})
	if err != nil {
		t.Fatalf("Invoke error = %v", err)
	}

	if bus.lastRequest.Source != 0x10 || bus.lastRequest.Target != 0x08 {
		t.Fatalf("unexpected request addresses: %+v", bus.lastRequest)
	}
	if bus.lastRequest.Primary != 0xB5 || bus.lastRequest.Secondary != 0x05 {
		t.Fatalf("unexpected request type: %+v", bus.lastRequest)
	}
	if !bytes.Equal(bus.lastRequest.Data, []byte{0x02, 0xAA, 0xBB}) {
		t.Fatalf("unexpected request data %v", bus.lastRequest.Data)
	}

	values, ok := result.(map[string]types.Value)
	if !ok {
		t.Fatalf("expected map[string]types.Value, got %T", result)
	}
	if op := values["op"]; !op.Valid || op.Value != byte(0x02) {
		t.Fatalf("op = %+v; want 0x02 valid", op)
	}
	if got := values["payload"]; !got.Valid || len(got.Value.([]byte)) != 0 {
		t.Fatalf("payload = %+v; want empty valid", got)
	}
}

func TestVaillantSystem_SetOperationalData_InputValidation(t *testing.T) {
	t.Parallel()

	planes := system.NewProvider().CreatePlanes(registry.DeviceInfo{Manufacturer: "Vaillant", Address: 0x08})
	plane := planes[0].(router.Plane)

	eventRouter := router.NewBusEventRouter(&vaillantMockBus{})

	_, err := eventRouter.Invoke(context.Background(), plane, "set_operational_data", map[string]any{
		"source": byte(0x10),
	})
	if !errors.Is(err, ebuserrors.ErrInvalidPayload) {
		t.Fatalf("expected ErrInvalidPayload, got %v", err)
	}

	_, err = eventRouter.Invoke(context.Background(), plane, "set_operational_data", map[string]any{
		"source": byte(0x10),
		"op":     byte(0x02),
		"data":   "nope",
	})
	if !errors.Is(err, ebuserrors.ErrInvalidPayload) {
		t.Fatalf("expected ErrInvalidPayload, got %v", err)
	}
}

func TestVaillantSystem_SetOperationalData_TypedErrors(t *testing.T) {
	t.Parallel()

	planes := system.NewProvider().CreatePlanes(registry.DeviceInfo{Manufacturer: "Vaillant", Address: 0x08})
	plane := planes[0].(router.Plane)

	tests := []struct {
		name string
		err  error
	}{
		{name: "NACK", err: ebuserrors.ErrNACK},
		{name: "NoDevice", err: ebuserrors.ErrNoSuchDevice},
		{name: "Timeout", err: ebuserrors.ErrTimeout},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			bus := &vaillantMockBus{err: tc.err}
			eventRouter := router.NewBusEventRouter(bus)

			_, err := eventRouter.Invoke(context.Background(), plane, "set_operational_data", map[string]any{
				"source": byte(0x10),
				"op":     byte(0x02),
			})
			if !errors.Is(err, tc.err) {
				t.Fatalf("expected %v, got %v", tc.err, err)
			}
		})
	}
}
