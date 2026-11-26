package output

import (
	"encoding/json"
	"reflect"
	"testing"
)

// The recorded old output for the identifier-bearing radio devices is retained by
// digest only, so the expectation below is written from the old mapping
// (radioDeviceFromSnapshot and the MCP adapter) over the same synthetic devices
// as the recording, not produced by this package.
func TestVaillantRadioDeviceOutputMatchesOldMapping(t *testing.T) {
	boolPtr := func(v bool) *bool { return &v }
	byteRef := func(v uint8) *uint8 { return &v }
	intPtr := func(v int) *int { return &v }
	strPtr := func(v string) *string { return &v }
	floatPtr := func(v float64) *float64 { return &v }

	native := []*NativeRadioDeviceSnapshot{
		{
			Group: 0x09, Instance: 0x00, SlotMode: "paired", DeviceConnected: boolPtr(true), DeviceClassAddress: byteRef(0x15), DeviceModel: "SYN-RC-1",
			FirmwareVersion: strPtr("1.2.3"), HardwareIdentifier: func() *uint16 { v := uint16(0x1234); return &v }(), RemoteControlAddress: byteRef(0x10), DevicePaired: boolPtr(true),
			ReceptionStrength: byteRef(4), ZoneAssignment: byteRef(1), RoomTemperatureC: floatPtr(21.25), RoomHumidityPct: floatPtr(47),
		},
		{Group: 0x0A, Instance: 0x01, SlotMode: "free", DeviceConnected: boolPtr(false), DeviceModel: "SYN-RC-2"},
	}
	wantGraphQL := []RadioDevice{
		{
			Group: 9, Instance: 0, SlotMode: "paired", DeviceConnected: boolPtr(true), DeviceClassAddress: intPtr(21), DeviceModel: "SYN-RC-1",
			FirmwareVersion: strPtr("1.2.3"), HardwareIdentifier: intPtr(4660), RemoteControlAddress: intPtr(16), DevicePaired: boolPtr(true),
			ReceptionStrength: intPtr(4), ZoneAssignment: intPtr(1), RoomTemperatureC: floatPtr(21.25), RoomHumidityPct: floatPtr(47),
		},
		{Group: 10, Instance: 1, SlotMode: "free", DeviceConnected: boolPtr(false), DeviceModel: "SYN-RC-2"},
	}
	wantMCP := `[{"group":9,"instance":0,"slot_mode":"paired","device_connected":true,"device_class_address":21,"device_model":"SYN-RC-1","firmware_version":"1.2.3","hardware_identifier":4660,"remote_control_address":16,"device_paired":true,"reception_strength":4,"zone_assignment":1,"room_temperature_c":21.25,"room_humidity_pct":47},{"group":10,"instance":1,"slot_mode":"free","device_connected":false,"device_model":"SYN-RC-2"}]`

	got := make([]RadioDevice, 0, len(native))
	for _, device := range native {
		got = append(got, RadioDeviceFromSnapshot(device))
	}
	if !reflect.DeepEqual(got, wantGraphQL) {
		t.Fatalf("radio device records:\ngot  %#v\nwant %#v", got, wantGraphQL)
	}
	encoded, err := json.Marshal(MCPRadioDevicesFromGraphQL(got))
	if err != nil {
		t.Fatal(err)
	}
	if string(encoded) != wantMCP {
		t.Fatalf("radio device MCP output:\ngot  %s\nwant %s", encoded, wantMCP)
	}
}
