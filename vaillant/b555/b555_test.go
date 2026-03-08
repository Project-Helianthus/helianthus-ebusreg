package b555

import (
	"testing"

	"github.com/Project-Helianthus/helianthus-ebusgo/protocol"
	"github.com/Project-Helianthus/helianthus-ebusgo/types"
	"github.com/Project-Helianthus/helianthus-ebusreg/registry"
)

func TestProvider_MatchesVaillant(t *testing.T) {
	p := NewProvider()
	tests := []struct {
		manufacturer string
		match        bool
	}{
		{"Vaillant", true},
		{"vaillant", true},
		{"VAILLANT", true},
		{"Saunier Duval", true},
		{"saunier", true},
		{"AWB", true},
		{"awb", true},
		{"Bosch", false},
		{"", false},
		{"  ", false},
	}
	for _, tt := range tests {
		got := p.Match(registry.DeviceInfo{Manufacturer: tt.manufacturer})
		if got != tt.match {
			t.Errorf("Match(%q) = %v, want %v", tt.manufacturer, got, tt.match)
		}
	}
}

func TestProvider_CreatesPlanesWithTimerMethods(t *testing.T) {
	p := NewProvider()
	planes := p.CreatePlanes(registry.DeviceInfo{Address: 0x15, Manufacturer: "Vaillant"})
	if len(planes) != 1 {
		t.Fatalf("expected 1 plane, got %d", len(planes))
	}
	if planes[0].Name() != "timer" {
		t.Errorf("plane name = %q, want %q", planes[0].Name(), "timer")
	}
	methods := planes[0].Methods()
	if len(methods) != 4 {
		t.Fatalf("expected 4 methods, got %d", len(methods))
	}

	wantMethods := []struct {
		name     string
		readOnly bool
	}{
		{"read_timer_config", true},
		{"read_timer_slots", true},
		{"read_timer", true},
		{"write_timer", false},
	}
	for i, want := range wantMethods {
		if methods[i].Name() != want.name {
			t.Errorf("method[%d].Name() = %q, want %q", i, methods[i].Name(), want.name)
		}
		if methods[i].ReadOnly() != want.readOnly {
			t.Errorf("method[%d].ReadOnly() = %v, want %v", i, methods[i].ReadOnly(), want.readOnly)
		}
	}
}

// --- Template Build tests ---

func TestConfigReadTemplate_Build(t *testing.T) {
	tmpl := configReadTemplate{primary: 0xB5, secondary: 0x55}
	payload, err := tmpl.Build(map[string]any{"zone": 0, "hc": 0})
	if err != nil {
		t.Fatal(err)
	}
	if len(payload) != 3 {
		t.Fatalf("payload len = %d, want 3", len(payload))
	}
	if payload[0] != 0xA3 || payload[1] != 0x00 || payload[2] != 0x00 {
		t.Errorf("payload = %X, want A30000", payload)
	}
}

func TestConfigReadTemplate_Build_Zone2DHW(t *testing.T) {
	tmpl := configReadTemplate{primary: 0xB5, secondary: 0x55}
	payload, err := tmpl.Build(map[string]any{"zone": 0xFF, "hc": 2})
	if err != nil {
		t.Fatal(err)
	}
	if payload[0] != 0xA3 || payload[1] != 0xFF || payload[2] != 0x02 {
		t.Errorf("payload = %X, want A3FF02", payload)
	}
}

func TestConfigReadTemplate_RejectsHCOver4(t *testing.T) {
	tmpl := configReadTemplate{primary: 0xB5, secondary: 0x55}
	_, err := tmpl.Build(map[string]any{"zone": 0, "hc": 5})
	if err == nil {
		t.Fatal("expected error for hc=5")
	}
}

func TestSlotsReadTemplate_Build(t *testing.T) {
	tmpl := slotsReadTemplate{primary: 0xB5, secondary: 0x55}
	payload, err := tmpl.Build(map[string]any{"zone": 1, "hc": 0})
	if err != nil {
		t.Fatal(err)
	}
	if payload[0] != 0xA4 || payload[1] != 0x01 || payload[2] != 0x00 {
		t.Errorf("payload = %X, want A40100", payload)
	}
}

func TestTimerReadTemplate_Build(t *testing.T) {
	tmpl := timerReadTemplate{primary: 0xB5, secondary: 0x55}
	payload, err := tmpl.Build(map[string]any{"zone": 0, "hc": 0, "weekday": 6, "slot": 0})
	if err != nil {
		t.Fatal(err)
	}
	if len(payload) != 5 {
		t.Fatalf("payload len = %d, want 5", len(payload))
	}
	// A5 00 00 06 00
	if payload[0] != 0xA5 || payload[1] != 0x00 || payload[2] != 0x00 || payload[3] != 0x06 || payload[4] != 0x00 {
		t.Errorf("payload = %X, want A500000600", payload)
	}
}

func TestTimerReadTemplate_RejectsWeekdayOver6(t *testing.T) {
	tmpl := timerReadTemplate{primary: 0xB5, secondary: 0x55}
	_, err := tmpl.Build(map[string]any{"zone": 0, "hc": 0, "weekday": 7, "slot": 0})
	if err == nil {
		t.Fatal("expected error for weekday=7")
	}
}

func TestTimerReadTemplate_RejectsMissingParams(t *testing.T) {
	tmpl := timerReadTemplate{primary: 0xB5, secondary: 0x55}
	_, err := tmpl.Build(nil)
	if err == nil {
		t.Fatal("expected error for nil params")
	}
	_, err = tmpl.Build(map[string]any{"zone": 0, "hc": 0, "weekday": 0})
	if err == nil {
		t.Fatal("expected error for missing slot")
	}
}

func TestTimerWriteTemplate_Build(t *testing.T) {
	tmpl := timerWriteTemplate{primary: 0xB5, secondary: 0x55}
	payload, err := tmpl.Build(map[string]any{
		"zone":         0,
		"hc":           0,
		"weekday":      0,
		"slot":         0,
		"slot_count":   1,
		"start_hour":   0,
		"start_minute": 0,
		"end_hour":     24,
		"end_minute":   0,
		"temperature":  225, // 22.5°C
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(payload) != 12 {
		t.Fatalf("payload len = %d, want 12", len(payload))
	}
	// A6 00 00 00 00 01 00 00 18 00 E1 00
	expected := []byte{0xA6, 0x00, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x18, 0x00, 0xE1, 0x00}
	for i := range expected {
		if payload[i] != expected[i] {
			t.Errorf("payload[%d] = 0x%02X, want 0x%02X", i, payload[i], expected[i])
		}
	}
}

func TestTimerWriteTemplate_DHWTemp(t *testing.T) {
	tmpl := timerWriteTemplate{primary: 0xB5, secondary: 0x55}
	payload, err := tmpl.Build(map[string]any{
		"zone":         0xFF,
		"hc":           2,
		"weekday":      0,
		"slot":         0,
		"slot_count":   1,
		"start_hour":   0,
		"start_minute": 0,
		"end_hour":     24,
		"end_minute":   0,
		"temperature":  0xFFFF, // no-op sentinel
	})
	if err != nil {
		t.Fatal(err)
	}
	// Temperature bytes: 0xFF, 0xFF (LE)
	if payload[10] != 0xFF || payload[11] != 0xFF {
		t.Errorf("temp bytes = [%02X %02X], want [FF FF]", payload[10], payload[11])
	}
}

func TestTimerWriteTemplate_RejectsHourOver24(t *testing.T) {
	tmpl := timerWriteTemplate{primary: 0xB5, secondary: 0x55}
	_, err := tmpl.Build(map[string]any{
		"zone": 0, "hc": 0, "weekday": 0, "slot": 0, "slot_count": 1,
		"start_hour": 25, "start_minute": 0, "end_hour": 24, "end_minute": 0,
		"temperature": 225,
	})
	if err == nil {
		t.Fatal("expected error for start_hour=25")
	}
}

func TestTimerWriteTemplate_RejectsMinuteOver59(t *testing.T) {
	tmpl := timerWriteTemplate{primary: 0xB5, secondary: 0x55}
	_, err := tmpl.Build(map[string]any{
		"zone": 0, "hc": 0, "weekday": 0, "slot": 0, "slot_count": 1,
		"start_hour": 0, "start_minute": 60, "end_hour": 24, "end_minute": 0,
		"temperature": 225,
	})
	if err == nil {
		t.Fatal("expected error for start_minute=60")
	}
}

func TestTimerWriteTemplate_RejectsZeroSlotCount(t *testing.T) {
	tmpl := timerWriteTemplate{primary: 0xB5, secondary: 0x55}
	_, err := tmpl.Build(map[string]any{
		"zone": 0, "hc": 0, "weekday": 0, "slot": 0, "slot_count": 0,
		"start_hour": 0, "start_minute": 0, "end_hour": 24, "end_minute": 0,
		"temperature": 225,
	})
	if err == nil {
		t.Fatal("expected error for slot_count=0")
	}
}

// --- Decode tests ---

func TestDecodeConfigResponse_Z1Heating(t *testing.T) {
	// Real response from live hardware: status=0, max_slots=12, time_res=10,
	// min_dur=5, has_temp=1, temp_slots=12, min=5, max=30, pad=0
	payload := []byte{0x00, 0x0C, 0x0A, 0x05, 0x01, 0x0C, 0x05, 0x1E, 0x00}
	values := decodeConfigResponse(0x00, 0x00, payload)

	assertValue(t, values, "status", byte(0x00))
	assertValue(t, values, "available", true)
	assertValue(t, values, "max_slots", byte(12))
	assertValue(t, values, "time_resolution", byte(10))
	assertValue(t, values, "min_duration", byte(5))
	assertValue(t, values, "has_temperature", true)
	assertValue(t, values, "temp_slots", byte(12))
	assertValue(t, values, "min_temp_c", float64(5))
	assertValue(t, values, "max_temp_c", float64(30))
}

func TestDecodeConfigResponse_CC(t *testing.T) {
	// CC: status=0, max_slots=3, time_res=10, min_dur=0, has_temp=0, temp_slots=0, min=FF, max=FF
	payload := []byte{0x00, 0x03, 0x0A, 0x00, 0x00, 0x00, 0xFF, 0xFF, 0x00}
	values := decodeConfigResponse(0x00, 0x03, payload)

	assertValue(t, values, "available", true)
	assertValue(t, values, "max_slots", byte(3))
	assertValue(t, values, "has_temperature", false)
	assertValue(t, values, "temp_slots", byte(0))
	assertNotValid(t, values, "min_temp_c")
	assertNotValid(t, values, "max_temp_c")
}

func TestDecodeConfigResponse_Unavailable(t *testing.T) {
	// Cooling: status=0x03 (unavailable)
	payload := []byte{0x03, 0x01, 0x01, 0x00, 0x00, 0x00, 0xFF, 0xFF, 0x00}
	values := decodeConfigResponse(0x00, 0x01, payload)

	assertValue(t, values, "status", byte(0x03))
	assertValue(t, values, "available", false)
}

func TestDecodeConfigResponse_ShortPayload(t *testing.T) {
	values := decodeConfigResponse(0x00, 0x00, []byte{0x00, 0x0C})
	assertNotValid(t, values, "status")
}

func TestDecodeSlotsResponse_Z1Heating(t *testing.T) {
	// status=0, all days have 1 slot
	payload := []byte{0x00, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x00}
	values := decodeSlotsResponse(0x00, 0x00, payload)

	assertValue(t, values, "available", true)
	slotsVal, ok := values["slots_per_day"]
	if !ok || !slotsVal.Valid {
		t.Fatal("slots_per_day missing or invalid")
	}
	slots, ok := slotsVal.Value.([]int)
	if !ok {
		t.Fatalf("slots_per_day type = %T, want []int", slotsVal.Value)
	}
	for i, s := range slots {
		if s != 1 {
			t.Errorf("slots_per_day[%d] = %d, want 1", i, s)
		}
	}
}

func TestDecodeTimerReadResponse_Z1HeatingMonday(t *testing.T) {
	// status=0, 00:00-24:00 @ 22.5°C
	payload := []byte{0x00, 0x00, 0x00, 0x18, 0x00, 0xE1, 0x00}
	values := decodeTimerReadResponse(0x00, 0x00, 0x00, 0x00, payload)

	assertValue(t, values, "available", true)
	assertValue(t, values, "start_hour", byte(0))
	assertValue(t, values, "start_minute", byte(0))
	assertValue(t, values, "end_hour", byte(24))
	assertValue(t, values, "end_minute", byte(0))
	assertValue(t, values, "temperature_raw", uint16(225))
	assertValue(t, values, "temperature_c", float64(22.5))
}

func TestDecodeTimerReadResponse_DHW(t *testing.T) {
	// DHW: 00:00-24:00 @ 61.0°C (raw=610=0x0262)
	payload := []byte{0x00, 0x00, 0x00, 0x18, 0x00, 0x62, 0x02}
	values := decodeTimerReadResponse(0x00, 0x02, 0x00, 0x00, payload)

	assertValue(t, values, "temperature_raw", uint16(610))
	assertValue(t, values, "temperature_c", float64(61.0))
}

func TestDecodeTimerReadResponse_CCNoTemp(t *testing.T) {
	// CC: 00:00-24:00, temp=0xFFFF (no temperature)
	payload := []byte{0x00, 0x00, 0x00, 0x18, 0x00, 0xFF, 0xFF}
	values := decodeTimerReadResponse(0x00, 0x03, 0x00, 0x00, payload)

	assertValue(t, values, "temperature_raw", uint16(0xFFFF))
	assertNotValid(t, values, "temperature_c")
}

func TestDecodeTimerReadResponse_ShortPayload(t *testing.T) {
	values := decodeTimerReadResponse(0x00, 0x00, 0x00, 0x00, []byte{0x00})
	assertNotValid(t, values, "status")
}

func TestDecodeTimerWriteResponse_ACK(t *testing.T) {
	values := decodeTimerWriteResponse(0x00, 0x00, 0x00, 0x00, []byte{0x00})
	assertValue(t, values, "error_code", byte(0x00))
	assertValue(t, values, "accepted", true)
}

func TestDecodeTimerWriteResponse_ParamOutOfRange(t *testing.T) {
	values := decodeTimerWriteResponse(0x00, 0x00, 0x00, 0x00, []byte{0x01})
	assertValue(t, values, "error_code", byte(0x01))
	assertValue(t, values, "accepted", false)
}

func TestDecodeTimerWriteResponse_Unavailable(t *testing.T) {
	values := decodeTimerWriteResponse(0x00, 0x01, 0x00, 0x00, []byte{0x03})
	assertValue(t, values, "error_code", byte(0x03))
	assertValue(t, values, "accepted", false)
}

func TestDecodeTimerWriteResponse_ValidationFailure(t *testing.T) {
	values := decodeTimerWriteResponse(0x00, 0x02, 0x00, 0x00, []byte{0x06})
	assertValue(t, values, "error_code", byte(0x06))
	assertValue(t, values, "accepted", false)
}

func TestDecodeTimerWriteResponse_EmptyPayload(t *testing.T) {
	values := decodeTimerWriteResponse(0x00, 0x00, 0x00, 0x00, nil)
	assertNotValid(t, values, "error_code")
	assertValue(t, values, "accepted", false)
}

// --- Router plane tests ---

func TestBuildRequest_ConfigRead(t *testing.T) {
	pl := newTimerPlane(registry.DeviceInfo{Address: 0x15, Manufacturer: "Vaillant"})
	methods := pl.Methods()
	frame, err := pl.BuildRequest(methods[0], map[string]any{
		"source": 0x71,
		"zone":   0,
		"hc":     0,
	})
	if err != nil {
		t.Fatal(err)
	}
	if frame.Source != 0x71 {
		t.Errorf("Source = 0x%02X, want 0x71", frame.Source)
	}
	if frame.Target != 0x15 {
		t.Errorf("Target = 0x%02X, want 0x15", frame.Target)
	}
	if frame.Primary != 0xB5 || frame.Secondary != 0x55 {
		t.Errorf("PB/SB = %02X/%02X, want B5/55", frame.Primary, frame.Secondary)
	}
	if len(frame.Data) != 3 || frame.Data[0] != 0xA3 {
		t.Errorf("Data = %X, want A3...", frame.Data)
	}
}

func TestBuildRequest_TimerWrite(t *testing.T) {
	pl := newTimerPlane(registry.DeviceInfo{Address: 0x15, Manufacturer: "Vaillant"})
	methods := pl.Methods()
	frame, err := pl.BuildRequest(methods[3], map[string]any{
		"source":       0x71,
		"zone":         0,
		"hc":           0,
		"weekday":      0,
		"slot":         0,
		"slot_count":   1,
		"start_hour":   0,
		"start_minute": 0,
		"end_hour":     24,
		"end_minute":   0,
		"temperature":  225,
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(frame.Data) != 12 {
		t.Fatalf("Data len = %d, want 12", len(frame.Data))
	}
	if frame.Data[0] != 0xA6 {
		t.Errorf("opcode = 0x%02X, want 0xA6", frame.Data[0])
	}
}

func TestDecodeResponse_ConfigRead(t *testing.T) {
	pl := newTimerPlane(registry.DeviceInfo{Address: 0x15, Manufacturer: "Vaillant"})
	methods := pl.Methods()
	result, err := pl.DecodeResponse(methods[0], protocol.Frame{
		Primary:   0xB5,
		Secondary: 0x55,
		Data:      []byte{0x00, 0x0C, 0x0A, 0x05, 0x01, 0x0C, 0x05, 0x1E, 0x00},
	}, map[string]any{"zone": 0, "hc": 0})
	if err != nil {
		t.Fatal(err)
	}
	values, ok := result.(map[string]types.Value)
	if !ok {
		t.Fatalf("result type = %T, want map[string]types.Value", result)
	}
	assertValue(t, values, "max_slots", byte(12))
}

func TestDecodeResponse_TimerRead(t *testing.T) {
	pl := newTimerPlane(registry.DeviceInfo{Address: 0x15, Manufacturer: "Vaillant"})
	methods := pl.Methods()
	result, err := pl.DecodeResponse(methods[2], protocol.Frame{
		Primary:   0xB5,
		Secondary: 0x55,
		Data:      []byte{0x00, 0x00, 0x00, 0x18, 0x00, 0xE1, 0x00},
	}, map[string]any{"zone": 0, "hc": 0, "weekday": 0, "slot": 0})
	if err != nil {
		t.Fatal(err)
	}
	values, ok := result.(map[string]types.Value)
	if !ok {
		t.Fatalf("result type = %T, want map[string]types.Value", result)
	}
	assertValue(t, values, "temperature_c", float64(22.5))
}

func TestDecodeResponse_RejectsWrongSecondary(t *testing.T) {
	pl := newTimerPlane(registry.DeviceInfo{Address: 0x15, Manufacturer: "Vaillant"})
	methods := pl.Methods()
	_, err := pl.DecodeResponse(methods[0], protocol.Frame{
		Primary:   0xB5,
		Secondary: 0x24, // wrong — B524, not B555
		Data:      []byte{0x00},
	}, map[string]any{"zone": 0, "hc": 0})
	if err == nil {
		t.Fatal("expected error for wrong secondary byte")
	}
}

// --- Helpers ---

func assertValue(t *testing.T, values map[string]types.Value, key string, want any) {
	t.Helper()
	v, ok := values[key]
	if !ok {
		t.Errorf("key %q missing", key)
		return
	}
	if !v.Valid {
		t.Errorf("key %q not valid", key)
		return
	}
	if v.Value != want {
		t.Errorf("key %q = %v (%T), want %v (%T)", key, v.Value, v.Value, want, want)
	}
}

func assertNotValid(t *testing.T, values map[string]types.Value, key string) {
	t.Helper()
	v, ok := values[key]
	if !ok {
		return // missing is acceptable for "not valid"
	}
	if v.Valid {
		t.Errorf("key %q should not be valid, got %v", key, v.Value)
	}
}
