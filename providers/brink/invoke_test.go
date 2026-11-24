package brink

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"math"
	"testing"

	ebuserrors "github.com/Project-Helianthus/helianthus-ebusgo/errors"
	"github.com/Project-Helianthus/helianthus-ebusgo/protocol"
	brinkcat "github.com/Project-Helianthus/helianthus-ebusreg/catalog/brink"
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

func TestBrinkInvokeDocumentedReadAndWrite(t *testing.T) {
	planes := NewProvider().CreatePlanes(registry.DeviceInfo{Address: 0x3C})
	plane, ok := planes[0].(router.Plane)
	if !ok {
		t.Fatal("brink plane does not implement router.Plane")
	}

	readBus := &recordingBus{reply: &protocol.Frame{Primary: 0x40, Secondary: 0x22, Data: []byte{0xC8, 0x00}}}
	read, err := router.NewBusEventRouter(readBus).Invoke(context.Background(), plane, methodSensorRead, map[string]any{
		"source": byte(0x10),
		"id":     byte(0x07),
	})
	if err != nil {
		t.Fatalf("sensor read: %v", err)
	}
	assertRequest(t, readBus, 0x3C, 0x40, 0x22, []byte{0x07})
	reading, ok := read.(brinkcat.Reading)
	if !ok || !reading.Known || reading.Name != "inside_temperature" || reading.Value != 20 {
		t.Fatalf("sensor read = %#v", read)
	}

	writeBus := &recordingBus{reply: &protocol.Frame{Primary: 0x40, Secondary: 0x80}}
	if _, err = router.NewBusEventRouter(writeBus).Invoke(context.Background(), plane, methodParameterWrite, map[string]any{
		"source": byte(0x10), "id": byte(0x0A), "value": 4320, "payload": []byte{0x99},
	}); err != nil {
		t.Fatalf("parameter write: %v", err)
	}
	assertRequest(t, writeBus, 0x3C, 0x40, 0x80, []byte{0x0A, 0xE0, 0x10})
}

func TestBrinkInvokeEveryDocumentedRequest(t *testing.T) {
	parameterReply := []byte{0x64, 0x00, 0x32, 0x00, 0x2C, 0x01, 0x01, 0x00, 0x64, 0x00}
	result, bus := invokeOK(t, 0x3C, methodParameterRead, map[string]any{
		"source": byte(0x10), "id": byte(0x01),
	}, 0x40, 0x50, parameterReply)
	assertRequest(t, bus, 0x3C, 0x40, 0x50, []byte{0x01})
	parameter, ok := result.(brinkcat.ParameterReading)
	if !ok || parameter.Name != "flow_mode_1" || parameter.Current != 100 || parameter.Min != 50 || parameter.Default != 100 {
		t.Fatalf("parameter read = %#v", result)
	}

	result, bus = invokeOK(t, 0x3C, methodErrorHistoryRead, map[string]any{"source": byte(0x10)}, 0x40, 0x90, []byte{0xAB, 0xCD})
	assertRequest(t, bus, 0x3C, 0x40, 0x90, []byte{0x00})
	raw, ok := result.([]byte)
	if !ok || !bytes.Equal(raw, []byte{0xAB, 0xCD}) {
		t.Fatalf("error history = %X", result)
	}

	_, bus = invokeOK(t, 0x3C, methodFactoryReset, map[string]any{
		"source": byte(0x10), "payload": []byte{0x99, 0x88},
	}, 0x40, 0xFF, nil)
	assertRequest(t, bus, 0x3C, 0x40, 0xFF, []byte{})

	_, bus = invokeOK(t, 0x3C, methodResetNotifications, map[string]any{"source": byte(0x10)}, 0x40, 0x91, nil)
	assertRequest(t, bus, 0x3C, 0x40, 0x91, []byte{0x00})

	_, bus = invokeOK(t, 0x3C, methodFanModeTemporary, map[string]any{"source": byte(0x10), "mode": byte(0x03)}, 0x40, 0xA1, nil)
	assertRequest(t, bus, 0x3C, 0x40, 0xA1, []byte{0x00, 0x03})

	_, bus = invokeOK(t, 0x3C, methodFanModePersistent, map[string]any{"source": byte(0x10), "mode": byte(0x01)}, 0x40, 0xCB, nil)
	assertRequest(t, bus, 0x3C, 0x40, 0xCB, []byte{0x00, 0x01})
}

func TestBrinkInvokeSendsNothingWhenUnspecified(t *testing.T) {
	assertNoSend(t, 0x3C, methodSensorRead, map[string]any{"source": byte(0x10)}, ebuserrors.ErrInvalidPayload)
	assertNoSend(t, 0x3C, methodParameterRead, map[string]any{"source": byte(0x10), "id": byte(0x7F)}, brinkcat.ErrUnknownID)
	assertNoSend(t, 0x3C, methodParameterWrite, map[string]any{
		"source": byte(0x10), "id": byte(0x7F), "value": 1, "payload": []byte{0x7F, 0x01},
	}, brinkcat.ErrUnknownID)
	assertNoSend(t, 0x3C, methodFanModeTemporary, map[string]any{"source": byte(0x10)}, ebuserrors.ErrInvalidPayload)
	assertNoSend(t, 0x3C, methodParameterWrite, nil, ebuserrors.ErrInvalidPayload)
	assertNoSend(t, 0x3C, "not_a_method", map[string]any{"source": byte(0x10)}, ebuserrors.ErrInvalidPayload)
	assertNoSend(t, 0, methodFactoryReset, map[string]any{"source": byte(0x10)}, ebuserrors.ErrInvalidPayload)
}

func TestBrinkInvokeJSONDecodedReadAndWrite(t *testing.T) {
	planes := NewProvider().CreatePlanes(registry.DeviceInfo{Address: 0x3C})
	plane, ok := planes[0].(router.Plane)
	if !ok {
		t.Fatal("brink plane does not implement router.Plane")
	}

	readBus := &recordingBus{reply: &protocol.Frame{Primary: 0x40, Secondary: 0x22, Data: []byte{0xC8, 0x00}}}
	read, err := router.NewBusEventRouter(readBus).Invoke(context.Background(), plane, methodSensorRead, jsonParams(t, `{"source":16,"id":7}`))
	if err != nil {
		t.Fatalf("sensor read: %v", err)
	}
	assertRequest(t, readBus, 0x3C, 0x40, 0x22, []byte{0x07})
	reading, ok := read.(brinkcat.Reading)
	if !ok || !reading.Known || reading.Name != "inside_temperature" || reading.Value != 20 {
		t.Fatalf("sensor read = %#v", read)
	}

	writeBus := &recordingBus{reply: &protocol.Frame{Primary: 0x40, Secondary: 0x80}}
	if _, err = router.NewBusEventRouter(writeBus).Invoke(context.Background(), plane, methodParameterWrite, jsonParams(t, `{"source":16,"id":10,"value":4320}`)); err != nil {
		t.Fatalf("parameter write: %v", err)
	}
	assertRequest(t, writeBus, 0x3C, 0x40, 0x80, []byte{0x0A, 0xE0, 0x10})

	for _, raw := range []string{
		`{"source":16.5,"id":7}`,
		`{"source":256,"id":7}`,
		`{"source":-1,"id":7}`,
		`{"source":16,"id":7.2}`,
		`{"source":16,"id":256}`,
		`{"source":16,"id":-1}`,
	} {
		assertNoSend(t, 0x3C, methodSensorRead, jsonParams(t, raw), ebuserrors.ErrInvalidPayload)
	}
	assertNoSend(t, 0x3C, methodParameterWrite, jsonParams(t, `{"source":16,"id":10,"value":4320.5}`), ebuserrors.ErrInvalidPayload)
}

func TestBrinkUint8ParamRejectsNonIntegral(t *testing.T) {
	for _, value := range []any{16.5, 256.0, -1.0, 255.2, math.NaN(), math.Inf(1), math.Inf(-1), "16"} {
		if _, ok := uint8Param(map[string]any{"source": value}, "source"); ok {
			t.Fatalf("accepted %#v", value)
		}
	}
	if got, ok := uint8Param(map[string]any{"source": 0.0}, "source"); !ok || got != 0 {
		t.Fatalf("0.0 => %d ok=%v", got, ok)
	}
	if got, ok := uint8Param(map[string]any{"source": 255.0}, "source"); !ok || got != 255 {
		t.Fatalf("255.0 => %d ok=%v", got, ok)
	}
	if got, ok := uint8Param(map[string]any{"source": json.Number("16")}, "source"); !ok || got != 16 {
		t.Fatalf("json.Number => %d ok=%v", got, ok)
	}
	if _, ok := uint8Param(map[string]any{"source": json.Number("16.5")}, "source"); ok {
		t.Fatal("accepted fractional json.Number")
	}
}

func TestBrinkInvokeUsesTargetWhenAddressUnset(t *testing.T) {
	_, bus := invokeOK(t, 0, methodFactoryReset, map[string]any{
		"source": byte(0x10), "target": byte(0x3C),
	}, 0x40, 0xFF, nil)
	assertRequest(t, bus, 0x3C, 0x40, 0xFF, []byte{})
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
