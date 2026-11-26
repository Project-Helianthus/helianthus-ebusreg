package output

import (
	"encoding/json"
	"reflect"
	"testing"
)

// goFieldValues renders a record the way the recorded old outputs do: keyed by
// the Go field name, pointers dereferenced, nil pointers as null.
func goFieldValues(value reflect.Value) map[string]any {
	out := make(map[string]any, value.NumField())
	for i := 0; i < value.NumField(); i++ {
		field := value.Field(i)
		if field.Kind() == reflect.Pointer {
			if field.IsNil() {
				out[value.Type().Field(i).Name] = nil
				continue
			}
			field = field.Elem()
		}
		switch field.Kind() {
		case reflect.String:
			out[value.Type().Field(i).Name] = field.String()
		case reflect.Bool:
			out[value.Type().Field(i).Name] = field.Bool()
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			out[value.Type().Field(i).Name] = float64(field.Int())
		case reflect.Float32, reflect.Float64:
			out[value.Type().Field(i).Name] = field.Float()
		}
	}
	return out
}

// The recorded old radio-device output (testdata/radio-devices-recorded-output.json,
// copied unchanged from the recorded full-site public output) is compared exactly
// with what this package builds. The recording holds the published records, not the
// native snapshots, so the native snapshots below carry the recorded values in
// their native widths.
func TestVaillantRadioDeviceOutputReplaysRecordedOutput(t *testing.T) {
	recorded := readFixture[[]map[string]any](t, "radio-devices-recorded-output.json")
	boolPtr := func(v bool) *bool { return &v }
	byteRef := func(v uint8) *uint8 { return &v }
	strPtr := func(v string) *string { return &v }
	floatPtr := func(v float64) *float64 { return &v }
	hardware := func(v uint16) *uint16 { return &v }

	native := []*NativeRadioDeviceSnapshot{
		{
			Group: 10, Instance: 1, SlotMode: "active", DeviceConnected: boolPtr(true), DeviceClassAddress: byteRef(58), DeviceModel: "Unknown (0x3A)",
			FirmwareVersion: strPtr("01.02.03"), HardwareIdentifier: hardware(258), RemoteControlAddress: byteRef(1), DevicePaired: boolPtr(true),
			ReceptionStrength: byteRef(7), ZoneAssignment: byteRef(0), RoomTemperatureC: floatPtr(20.75), RoomHumidityPct: floatPtr(48),
		},
		{
			Group: 12, Instance: 0, SlotMode: "active", DeviceConnected: boolPtr(true), DeviceClassAddress: byteRef(38), DeviceModel: "VR71/FM5",
			FirmwareVersion: strPtr("02.01.00"), HardwareIdentifier: hardware(513),
		},
	}
	graphql := make([]RadioDevice, 0, len(native))
	for _, device := range native {
		graphql = append(graphql, RadioDeviceFromSnapshot(device))
	}
	mcp := MCPRadioDevicesFromGraphQL(graphql)
	if len(mcp) != len(recorded) || len(graphql) != len(recorded) {
		t.Fatalf("record counts: got %d GraphQL and %d MCP; recorded %d", len(graphql), len(mcp), len(recorded))
	}
	for i, want := range recorded {
		if got := goFieldValues(reflect.ValueOf(graphql[i])); !reflect.DeepEqual(got, want) {
			t.Errorf("GraphQL radio device %d:\ngot  %v\nwant %v", i, got, want)
		}
		if got := goFieldValues(reflect.ValueOf(mcp[i])); !reflect.DeepEqual(got, want) {
			t.Errorf("MCP radio device %d:\ngot  %v\nwant %v", i, got, want)
		}
	}
}

// The recorded old MCP responses of the radio-devices tool
// (testdata/mcp-radio-devices-recorded.json, the two radio records of the pre-port
// typed MCP recording, extracted unchanged) are compared with the serialized MCP
// records this package builds, by response key. The inputs are the provider values
// of the two recorded fixtures: the all-fields fixture fills every field with a
// synthetic value, the second fixture is the pin's semantic-snapshots provider.
func TestVaillantRadioDeviceMCPRecordsReplayRecordedResponses(t *testing.T) {
	type recordedResponse struct {
		Name     string
		Fixture  string
		Envelope struct {
			Data json.RawMessage
		}
	}
	recorded := readFixture[[]recordedResponse](t, "mcp-radio-devices-recorded.json")
	boolPtr := func(v bool) *bool { return &v }
	intPtr := func(v int) *int { return &v }
	strPtr := func(v string) *string { return &v }
	floatPtr := func(v float64) *float64 { return &v }
	inputs := map[string][]RadioDevice{
		"complete-synthetic": {{
			Group: 129, Instance: 130, SlotMode: "synthetic-SlotMode", DeviceConnected: boolPtr(true), DeviceClassAddress: intPtr(131),
			DeviceModel: "synthetic-DeviceModel", FirmwareVersion: strPtr("synthetic-FirmwareVersion"), HardwareIdentifier: intPtr(132),
			RemoteControlAddress: intPtr(133), DevicePaired: boolPtr(true), ReceptionStrength: intPtr(134), ZoneAssignment: intPtr(135),
			RoomTemperatureC: floatPtr(46.5), RoomHumidityPct: floatPtr(47.5),
		}},
		"the provider of the pin's TestServer_ToolsCallSemanticSnapshots": {
			{Group: 9, Instance: 1, DeviceModel: "VRC720"},
			{Group: 10, Instance: 1, DeviceModel: "VR92"},
		},
	}
	if len(recorded) != len(inputs) {
		t.Fatalf("recorded responses: got %d, want %d", len(recorded), len(inputs))
	}
	for _, record := range recorded {
		input, ok := inputs[record.Fixture]
		if !ok {
			t.Fatalf("no input for recorded fixture %q", record.Fixture)
		}
		got, err := json.Marshal(MCPRadioDevicesFromGraphQL(input))
		if err != nil {
			t.Fatal(err)
		}
		var gotValue, wantValue any
		if err := json.Unmarshal(got, &gotValue); err != nil {
			t.Fatal(err)
		}
		if err := json.Unmarshal(record.Envelope.Data, &wantValue); err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(gotValue, wantValue) {
			t.Errorf("%s (%s):\ngot  %s\nwant %s", record.Name, record.Fixture, got, record.Envelope.Data)
		}
	}
}
