package wolf

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"testing"

	ebuserrors "github.com/Project-Helianthus/helianthus-ebusgo/errors"
	"github.com/Project-Helianthus/helianthus-ebusgo/protocol"
	wolfcat "github.com/Project-Helianthus/helianthus-ebusreg/catalog/wolf"
	"github.com/Project-Helianthus/helianthus-ebusreg/registry"
	"github.com/Project-Helianthus/helianthus-ebusreg/router"
)

type recordingBus struct {
	sends   int
	request protocol.Frame
	reply   *protocol.Frame
}

func (bus *recordingBus) Send(ctx context.Context, frame protocol.Frame) (*protocol.Frame, error) {
	bus.sends++
	bus.request = frame
	bus.request.Data = append([]byte(nil), frame.Data...)
	return bus.reply, nil
}

type invokeCase struct {
	name      string
	method    string
	params    map[string]any
	primary   byte
	secondary byte
	payload   []byte
}

func TestWolfInvokeDocumentedReadAndWrite(t *testing.T) {
	reg := registry.NewDeviceRegistry([]registry.PlaneProvider{NewProvider()})
	reg.Register(registry.DeviceInfo{Address: 0x08, Manufacturer: "0x19", DeviceID: "WOLF0"})
	entry, ok := reg.Lookup(0x08)
	if !ok || len(entry.Planes()) != 1 {
		t.Fatalf("registry planes missing: %+v ok %v", entry, ok)
	}
	plane, ok := entry.Planes()[0].(router.Plane)
	if !ok {
		t.Fatal("registered wolf plane does not implement router.Plane")
	}

	readBus := &recordingBus{reply: &protocol.Frame{Primary: 0x50, Secondary: 0x22, Data: []byte{0x64, 0x00}}}
	read, err := router.NewBusEventRouter(readBus).Invoke(context.Background(), plane, methodServiceRead, map[string]any{
		"source": byte(0x10),
		"id":     []byte{0x84, 0x22, 0x00},
	})
	if err != nil {
		t.Fatalf("service read: %v", err)
	}
	assertRequest(t, readBus, 0x08, 0x50, 0x22, []byte{0x84, 0x22, 0x00})
	checkServiceRead(t, read)

	writeBus := &recordingBus{reply: &protocol.Frame{Primary: 0x50, Secondary: 0x23}}
	if _, err = router.NewBusEventRouter(writeBus).Invoke(context.Background(), plane, methodServiceWrite, map[string]any{
		"source":  byte(0x10),
		"id":      []byte{0x84, 0x22, 0x00},
		"value":   uint16(0x0064),
		"trailer": []byte{0x01, 0x02, 0x03, 0x04},
	}); err != nil {
		t.Fatalf("service write: %v", err)
	}
	assertRequest(t, writeBus, 0x08, 0x50, 0x23, []byte{0x84, 0x22, 0x64, 0x00, 0x01, 0x02, 0x03, 0x04})
}

func TestWolfInvokeEveryDocumentedRequest(t *testing.T) {
	cases := []invokeCase{
		{
			name: "burner status", method: methodBurnerStatusSend, primary: 0x05, secondary: 0x03,
			params: map[string]any{
				"source": byte(0x10), "state": byte(0x05), "modulation": byte(0xFF),
				"kessel_temperature": []byte{0x3C}, "payload": []byte{0x00},
			},
			payload: []byte{0x01, 0x05, 0xFF, 0x3C},
		},
		{
			name: "modulation", method: methodModulationSend, primary: 0x05, secondary: 0x04,
			params:  map[string]any{"source": byte(0x10), "actual": byte(40), "maximum": byte(100)},
			payload: []byte{40, 100},
		},
		{
			name: "heat request", method: methodHeatRequestSend, primary: 0x05, secondary: 0x07,
			params:  map[string]any{"source": byte(0x10), "state": byte(0xAA), "setpoint": byte(60), "pressure": []byte{0x10}},
			payload: []byte{0xAA, 60, 0x10},
		},
		{
			name: "boiler setpoint", method: methodBoilerSetpointSend, primary: 0x08, secondary: 0x00,
			params:  map[string]any{"source": byte(0x10), "celsius": 40.0, "forced_output": byte(0x02), "status": byte(0x03)},
			payload: []byte{0x80, 0x02, 0x02, 0x03},
		},
		{
			name: "solar status", method: methodSolarStatusSend, primary: 0x50, secondary: 0x14,
			params:  map[string]any{"source": byte(0x10), "status": byte(0x05), "ignored": byte(0x00), "rest": []byte{0x01}},
			payload: []byte{0x05, 0x00, 0x01},
		},
		{
			name: "solar pump", method: methodSolarPumpSend, primary: 0x50, secondary: 0x17,
			params:  map[string]any{"source": byte(0x10), "on": true, "ignored": byte(0x00)},
			payload: []byte{0xBD, 0x00},
		},
		{
			name: "solar power", method: methodSolarPowerSend, primary: 0x50, secondary: 0x18,
			params:  map[string]any{"source": byte(0x10), "watts": 1.5, "yield_high": byte(0x07), "yield_sum_low": byte(0x09)},
			payload: []byte{0x80, 0x01, 0x07, 0x09},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, bus := invokeOK(t, 0x08, tc.method, tc.params, tc.primary, tc.secondary, nil)
			assertRequest(t, bus, 0x08, tc.primary, tc.secondary, tc.payload)
		})
	}
}

func TestWolfInvokeJSONDecodedServiceReadAndWrite(t *testing.T) {
	planes := NewProvider().CreatePlanes(registry.DeviceInfo{Address: 0x08})
	plane, ok := planes[0].(router.Plane)
	if !ok {
		t.Fatal("wolf plane does not implement router.Plane")
	}

	readBus := &recordingBus{reply: &protocol.Frame{Primary: 0x50, Secondary: 0x22, Data: []byte{0x64, 0x00}}}
	read, err := router.NewBusEventRouter(readBus).Invoke(context.Background(), plane, methodServiceRead, jsonParams(t, `{"source":16,"id":[132,34,0]}`))
	if err != nil {
		t.Fatalf("service read: %v", err)
	}
	assertRequest(t, readBus, 0x08, 0x50, 0x22, []byte{0x84, 0x22, 0x00})
	checkServiceRead(t, read)

	writeBus := &recordingBus{reply: &protocol.Frame{Primary: 0x50, Secondary: 0x23}}
	if _, err = router.NewBusEventRouter(writeBus).Invoke(context.Background(), plane, methodServiceWrite, jsonParams(t, `{"source":16,"id":[132,34,0],"value":100,"trailer":[1,2,3,4]}`)); err != nil {
		t.Fatalf("service write: %v", err)
	}
	assertRequest(t, writeBus, 0x08, 0x50, 0x23, []byte{0x84, 0x22, 0x64, 0x00, 0x01, 0x02, 0x03, 0x04})

	assertNoSend(t, 0x08, methodServiceRead, jsonParams(t, `{"source":16,"id":[132,34,0.5]}`), ebuserrors.ErrInvalidPayload)
	assertNoSend(t, 0x08, methodServiceRead, jsonParams(t, `{"source":16,"id":[132,34,256]}`), ebuserrors.ErrInvalidPayload)
	assertNoSend(t, 0x08, methodServiceRead, jsonParams(t, `{"source":16,"id":[132,34,-1]}`), ebuserrors.ErrInvalidPayload)
	assertNoSend(t, 0x08, methodServiceWrite, jsonParams(t, `{"source":16,"id":[132,34,0],"value":100,"trailer":[1,2,3,4.5]}`), ebuserrors.ErrInvalidPayload)
	assertNoSend(t, 0x08, methodServiceWrite, jsonParams(t, `{"source":16,"id":[132,34,0],"value":100,"trailer":[1,2,3,256]}`), ebuserrors.ErrInvalidPayload)
}

func TestWolfInvokeSendsNothingWhenUnspecified(t *testing.T) {
	assertNoSend(t, 0x08, "not_a_method", map[string]any{"source": byte(0x10)}, ebuserrors.ErrInvalidPayload)
	assertNoSend(t, 0x08, methodBurnerStatusSend, nil, ebuserrors.ErrInvalidPayload)
	assertNoSend(t, 0x08, methodBurnerStatusSend, map[string]any{
		"source": byte(0x10), "payload": []byte{0x01, 0x02, 0x03},
	}, ebuserrors.ErrInvalidPayload)
	assertNoSend(t, 0x08, methodBurnerStatusSend, map[string]any{
		"source": byte(0x10), "state": byte(0x05), "modulation": byte(0xFF), "kessel_temperature": "3C",
	}, ebuserrors.ErrInvalidPayload)
	assertNoSend(t, 0x08, methodModulationSend, map[string]any{
		"source": byte(0x10), "actual": byte(101), "maximum": byte(0),
	}, wolfcat.ErrRange)
	assertNoSend(t, 0x08, methodServiceRead, map[string]any{
		"source": byte(0x10), "id": []byte{0x00, 0x00, 0x00},
	}, ebuserrors.ErrInvalidPayload)
	assertNoSend(t, 0x08, methodServiceWrite, map[string]any{
		"source": byte(0x10), "id": []byte{0xCD, 0x59, 0x01}, "value": uint16(1), "trailer": []byte{0, 0, 0, 0},
	}, wolfcat.ErrReadOnly)
	assertNoSend(t, 0x08, methodSolarStatusSend, map[string]any{
		"source": byte(0x10), "status": byte(0x05),
	}, ebuserrors.ErrInvalidPayload)
	assertNoSend(t, 0, methodModulationSend, map[string]any{
		"source": byte(0x10), "actual": byte(40), "maximum": byte(100),
	}, ebuserrors.ErrInvalidPayload)
}

func TestWolfInvokeUsesTargetWhenAddressUnset(t *testing.T) {
	_, bus := invokeOK(t, 0, methodModulationSend, map[string]any{
		"source": byte(0x10), "target": byte(0x15), "actual": byte(40), "maximum": byte(100),
	}, 0x05, 0x04, nil)
	assertRequest(t, bus, 0x15, 0x05, 0x04, []byte{40, 100})
}

func invokeOK(t *testing.T, address byte, method string, params map[string]any, primary, secondary byte, reply []byte) (any, *recordingBus) {
	t.Helper()
	planes := NewProvider().CreatePlanes(registry.DeviceInfo{Address: address})
	if len(planes) != 1 {
		t.Fatalf("planes = %d", len(planes))
	}
	plane, ok := planes[0].(router.Plane)
	if !ok {
		t.Fatal("plane does not implement router.Plane")
	}
	bus := &recordingBus{reply: &protocol.Frame{Source: address, Target: 0x10, Primary: primary, Secondary: secondary, Data: reply}}
	result, err := router.NewBusEventRouter(bus).Invoke(context.Background(), plane, method, params)
	if err != nil {
		t.Fatalf("Invoke %s: %v", method, err)
	}
	return result, bus
}

func assertRequest(t *testing.T, bus *recordingBus, target, primary, secondary byte, payload []byte) {
	t.Helper()
	if bus.sends != 1 {
		t.Fatalf("sends = %d; want 1", bus.sends)
	}
	got := bus.request
	if got.Source != 0x10 || got.Target != target || got.Primary != primary || got.Secondary != secondary || !bytes.Equal(got.Data, payload) {
		t.Fatalf("request source=%02X target=%02X pb=%02X sb=%02X data=% X; want target=%02X pb=%02X sb=%02X data=% X",
			got.Source, got.Target, got.Primary, got.Secondary, got.Data, target, primary, secondary, payload)
	}
}

func jsonParams(t *testing.T, raw string) map[string]any {
	t.Helper()
	var params map[string]any
	if err := json.Unmarshal([]byte(raw), &params); err != nil {
		t.Fatalf("json: %v", err)
	}
	return params
}

func assertNoSend(t *testing.T, address byte, method string, params map[string]any, want error) {
	t.Helper()
	planes := NewProvider().CreatePlanes(registry.DeviceInfo{Address: address})
	plane, ok := planes[0].(router.Plane)
	if !ok {
		t.Fatal("plane does not implement router.Plane")
	}
	bus := &recordingBus{reply: &protocol.Frame{}}
	_, err := router.NewBusEventRouter(bus).Invoke(context.Background(), plane, method, params)
	if !errors.Is(err, want) {
		t.Fatalf("Invoke %s err = %v; want %v", method, err, want)
	}
	if bus.sends != 0 {
		t.Fatalf("sends = %d; data % X", bus.sends, bus.request.Data)
	}
}

func checkServiceRead(t *testing.T, result any) {
	t.Helper()
	field, ok := result.(wolfcat.Field)
	value, isFloat := field.Value.(float64)
	if !ok || !isFloat || field.Name != "hg01 hysterese_vorlauf" || value != 10 {
		t.Fatalf("service read = %#v", result)
	}
}
