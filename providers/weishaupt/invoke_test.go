package weishaupt

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"testing"

	ebuserrors "github.com/Project-Helianthus/helianthus-ebusgo/errors"
	"github.com/Project-Helianthus/helianthus-ebusgo/protocol"
	weishauptcat "github.com/Project-Helianthus/helianthus-ebusreg/catalog/weishaupt"
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

func TestWeishauptInvokeDocumentedReadAndWrite(t *testing.T) {
	reg := registry.NewDeviceRegistry([]registry.PlaneProvider{NewProvider()})
	reg.Register(registry.DeviceInfo{Address: 0x08, Manufacturer: "0x10", DeviceID: "WRSol"})
	entry, ok := reg.Lookup(0x08)
	if !ok || len(entry.Planes()) != 1 {
		t.Fatalf("registry planes missing: %+v ok %v", entry, ok)
	}
	plane, ok := entry.Planes()[0].(router.Plane)
	if !ok {
		t.Fatal("registered weishaupt plane does not implement router.Plane")
	}

	readBus := &recordingBus{reply: &protocol.Frame{Primary: 0x09, Secondary: 0x00, Data: []byte{0xD2, 0x00}}}
	read, err := router.NewBusEventRouter(readBus).Invoke(context.Background(), plane, methodRAMRead, map[string]any{
		"source":  byte(0x10),
		"address": uint16(0xF484),
	})
	if err != nil {
		t.Fatalf("RAM read: %v", err)
	}
	assertRequest(t, readBus, 0x08, 0x09, 0x00, []byte{0x84, 0xF4, 0x02})
	checkRAM(t, read)
}

func TestWeishauptDateTimeFrameIsBroadcast(t *testing.T) {
	// A sent broadcast has no reply. The synthetic bus returns (nil, nil), and
	// Invoke succeeds with a nil result. ZZ stays FEh for a bound address and
	// when the caller supplies a different target.
	want := []byte{0x00, 0x14, 0x56, 0x34, 0x12, 0x02, 0x02, 0x01, 0x26}
	params := []map[string]any{
		{
			"source": byte(0x10), "outdoor_c": 20.0,
			"seconds": 56, "minutes": 34, "hours": 12, "day": 2, "month": 2, "weekday": 1, "year": 26,
		},
		{
			"source": byte(0x10), "target": byte(0x15), "outdoor_c": 20.0,
			"seconds": 56, "minutes": 34, "hours": 12, "day": 2, "month": 2, "weekday": 1, "year": 26,
		},
		jsonParams(t, `{"source":16,"target":21,"outdoor_c":20,"seconds":56,"minutes":34,"hours":12,"day":2,"month":2,"weekday":1,"year":26}`),
	}
	for _, address := range []byte{0x00, 0x08} {
		planes := NewProvider().CreatePlanes(registry.DeviceInfo{Address: address})
		plane, ok := planes[0].(router.Plane)
		if !ok {
			t.Fatal("weishaupt plane does not implement router.Plane")
		}
		method := methodByName(t, plane, methodDateTimeSet)
		for i, param := range params {
			got, err := plane.BuildRequest(method, param)
			if err != nil {
				t.Fatalf("address %02X case %d BuildRequest: %v", address, i, err)
			}
			if got.Source != 0x10 || got.Target != 0xFE || got.Primary != 0x07 || got.Secondary != 0x00 || !bytes.Equal(got.Data, want) {
				t.Fatalf("address %02X case %d frame source=%02X target=%02X pb=%02X sb=%02X data=% X",
					address, i, got.Source, got.Target, got.Primary, got.Secondary, got.Data)
			}

			bus := &recordingBus{}
			result, err := router.NewBusEventRouter(bus).Invoke(context.Background(), plane, methodDateTimeSet, param)
			if err != nil {
				t.Fatalf("address %02X case %d Invoke: %v", address, i, err)
			}
			if result != nil {
				t.Fatalf("address %02X case %d result = %#v; want nil", address, i, result)
			}
			assertRequest(t, bus, 0xFE, 0x07, 0x00, want)
		}
	}
}

func TestWeishauptInvokeJSONDecodedMemoryReadAndWrite(t *testing.T) {
	planes := NewProvider().CreatePlanes(registry.DeviceInfo{Address: 0x08})
	plane, ok := planes[0].(router.Plane)
	if !ok {
		t.Fatal("weishaupt plane does not implement router.Plane")
	}

	readBus := &recordingBus{reply: &protocol.Frame{Primary: 0x50, Secondary: 0x00, Data: []byte{0x7F}}}
	read, err := router.NewBusEventRouter(readBus).Invoke(context.Background(), plane, methodMemoryRead, jsonParams(t, `{"source":16,"id":[88,1,4]}`))
	if err != nil {
		t.Fatalf("memory read: %v", err)
	}
	assertRequest(t, readBus, 0x08, 0x50, 0x00, []byte{0x58, 0x01, 0x04})
	field, ok := read.(weishauptcat.Field)
	raw, isRaw := field.Value.([]byte)
	if !ok || !isRaw || field.Name != "value_raw" || !bytes.Equal(raw, []byte{0x7F}) {
		t.Fatalf("memory read = %#v", read)
	}

	writeBus := &recordingBus{reply: &protocol.Frame{Primary: 0x50, Secondary: 0x01}}
	if _, err = router.NewBusEventRouter(writeBus).Invoke(context.Background(), plane, methodMemoryWrite, jsonParams(t, `{"source":16,"id":[88,1,4],"value":[42]}`)); err != nil {
		t.Fatalf("memory write: %v", err)
	}
	assertRequest(t, writeBus, 0x08, 0x50, 0x01, []byte{0x58, 0x01, 0x04, 0x2A})

	assertNoSend(t, 0x08, methodMemoryRead, jsonParams(t, `{"source":16,"id":[88,1,4.5]}`), ebuserrors.ErrInvalidPayload)
	assertNoSend(t, 0x08, methodMemoryRead, jsonParams(t, `{"source":16,"id":[88,1,256]}`), ebuserrors.ErrInvalidPayload)
	assertNoSend(t, 0x08, methodMemoryRead, jsonParams(t, `{"source":16,"id":[88,1,-1]}`), ebuserrors.ErrInvalidPayload)
	assertNoSend(t, 0x08, methodMemoryWrite, jsonParams(t, `{"source":16,"id":[88,1,4],"value":[42.5]}`), ebuserrors.ErrInvalidPayload)
	assertNoSend(t, 0x08, methodMemoryWrite, jsonParams(t, `{"source":16,"id":[88,1,4],"value":[256]}`), ebuserrors.ErrInvalidPayload)
	assertNoSend(t, 0x08, methodMemoryWrite, jsonParams(t, `{"source":16,"id":[88,1,4],"value":[-1]}`), ebuserrors.ErrInvalidPayload)
}

func TestWeishauptInvokeEveryDocumentedRequest(t *testing.T) {
	identification := []byte{0x10, 'W', 'R', 'S', 'o', 'l', 0x24, 0x22, 0x11, 0x12}
	result, bus := invokeOK(t, 0x08, methodIdentificationRead, map[string]any{
		"source": byte(0x10), "payload": []byte{0xAB},
	}, 0x07, 0x04, identification)
	assertRequest(t, bus, 0x08, 0x07, 0x04, []byte{})
	id, ok := result.(weishauptcat.Identification)
	if !ok || id.Manufacturer != 0x10 || id.DeviceID != "WRSol" || id.SoftwareVersion != 24 || id.HardwareRev != 12 {
		t.Fatalf("identification = %#v", result)
	}

	result, bus = invokeOK(t, 0x08, methodMemoryRead, map[string]any{
		"source": byte(0x10), "id": []byte{0x58, 0x01, 0x04},
	}, 0x50, 0x00, []byte{0x7F})
	assertRequest(t, bus, 0x08, 0x50, 0x00, []byte{0x58, 0x01, 0x04})
	field, ok := result.(weishauptcat.Field)
	raw, isRaw := field.Value.([]byte)
	if !ok || !isRaw || field.Name != "value_raw" || !bytes.Equal(raw, []byte{0x7F}) {
		t.Fatalf("memory read = %#v", result)
	}

	_, bus = invokeOK(t, 0x08, methodMemoryWrite, map[string]any{
		"source": byte(0x10), "id": []byte{0x58, 0x01, 0x04}, "value": []byte{0x2A},
	}, 0x50, 0x01, nil)
	assertRequest(t, bus, 0x08, 0x50, 0x01, []byte{0x58, 0x01, 0x04, 0x2A})

	result, bus = invokeOK(t, 0x08, methodRAMRead, map[string]any{
		"source": byte(0x10), "address": uint16(0x0001),
	}, 0x09, 0x00, []byte{0x00, 0x01})
	assertRequest(t, bus, 0x08, 0x09, 0x00, []byte{0x01, 0x00, 0x02})
	raw, ok = result.([]byte)
	if !ok || !bytes.Equal(raw, []byte{0x00, 0x01}) {
		t.Fatalf("unknown RAM reply = %#v", result)
	}
}

func TestWeishauptInvokeSendsNothingWhenUnspecified(t *testing.T) {
	raw := map[string]any{"source": byte(0x10), "payload": []byte{0x01, 0x02, 0x03, 0x04}}
	assertNoSend(t, 0x08, methodSecondaryRead, raw, ebuserrors.ErrInvalidPayload)
	assertNoSend(t, 0x08, methodSecondaryWrite, raw, ebuserrors.ErrInvalidPayload)
	assertNoSend(t, 0x08, methodMemoryWrite, map[string]any{
		"source": byte(0x10), "address": uint16(0x0004), "value": []byte{0x2A},
	}, ebuserrors.ErrInvalidPayload)
	assertNoSend(t, 0x08, methodMemoryWrite, map[string]any{
		"source": byte(0x10), "id": []byte{0x58, 0x01, 0x04}, "value": []byte{1, 2, 3},
	}, weishauptcat.ErrValueSize)
	assertNoSend(t, 0x08, methodDateTimeSet, map[string]any{
		"source": byte(0x10), "outdoor_c": 50.5,
		"seconds": 0, "minutes": 0, "hours": 0, "day": 1, "month": 1, "weekday": 1, "year": 0,
	}, weishauptcat.ErrRange)
	assertNoSend(t, 0x08, methodRAMRead, map[string]any{"source": byte(0x10)}, ebuserrors.ErrInvalidPayload)
	assertNoSend(t, 0x08, "not_a_method", map[string]any{"source": byte(0x10)}, ebuserrors.ErrInvalidPayload)
	assertNoSend(t, 0, methodIdentificationRead, map[string]any{"source": byte(0x10)}, ebuserrors.ErrInvalidPayload)
}

func TestWeishauptInvokeUsesTargetWhenAddressUnset(t *testing.T) {
	_, bus := invokeOK(t, 0, methodIdentificationRead, map[string]any{
		"source": byte(0x10), "target": byte(0x15),
	}, 0x07, 0x04, []byte{0x10, 'W', 'R', 'S', 'o', 'l', 0x24, 0x22, 0x11, 0x12})
	assertRequest(t, bus, 0x15, 0x07, 0x04, []byte{})
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

func methodByName(t *testing.T, plane router.Plane, name string) registry.Method {
	t.Helper()
	for _, method := range plane.Methods() {
		if method.Name() == name {
			return method
		}
	}
	t.Fatalf("method %s missing", name)
	return nil
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

func checkRAM(t *testing.T, result any) {
	t.Helper()
	field, ok := result.(weishauptcat.Field)
	value, isFloat := field.Value.(float64)
	if !ok || !isFloat || field.Name != "TSO" || field.Unit != "C" || value != 21 {
		t.Fatalf("RAM read = %#v", result)
	}
}
