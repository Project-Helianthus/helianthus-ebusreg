package output

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

// The recorded old MCP response for a populated boiler (testdata/mcp-boiler-populated-recorded.json,
// the pin's MCP boiler case of the pre-port typed MCP recording, extracted unchanged) holds the flow
// temperature, the return temperature, the pump flag and the heating status of one boiler.
// This package builds the MCP record from the GraphQL record with the same four values and the result
// is compared with the recording as a whole. The recording is of the pre-port release, which hands its
// provider record to the server directly; the frozen pin's adapter builds the record itself. Two
// differences follow, both from the frozen pin and both listed here so that nothing else can differ:
//   - the response key of the flow temperature is flow_temperature at the pin (testdata/mcp-record-tags.txt)
//     and flow_temperature_c in the recording;
//   - the pin's adapter always allocates the config record (testdata/mcp-adapter-assignments.txt), so
//     the response holds an empty config object where the recording holds none.
func TestVaillantBoilerStatusPopulatedReplaysRecordedResponse(t *testing.T) {
	type recordedPopulated struct {
		Name     string
		Fixture  string
		Envelope struct {
			Data map[string]any
		}
	}
	recorded := readFixture[[]recordedPopulated](t, "mcp-boiler-populated-recorded.json")
	if len(recorded) != 1 {
		t.Fatalf("recorded cases: got %d, want 1", len(recorded))
	}
	record := recorded[0]

	flow, returned, pump, heating := 52.5, 47.0, true, 2
	input := &BoilerStatus{
		State:       BoilerState{FlowTemperatureC: &flow, ReturnTemperatureC: &returned, CentralHeatingPumpActive: &pump},
		Diagnostics: BoilerDiagnostics{HeatingStatusRaw: &heating},
	}
	encoded, err := json.Marshal(MCPBoilerStatusFromGraphQL(input))
	if err != nil {
		t.Fatal(err)
	}
	var got map[string]any
	if err := json.Unmarshal(encoded, &got); err != nil {
		t.Fatal(err)
	}

	tags, err := os.ReadFile(filepath.Join("testdata", "mcp-record-tags.txt"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(tags), "`json:\"flow_temperature,omitempty\"`") || strings.Contains(string(tags), "`json:\"flow_temperature_c,omitempty\"`") {
		t.Fatal("the pin's MCP boiler state tags do not carry the flow temperature as flow_temperature")
	}
	assignments, err := os.ReadFile(filepath.Join("testdata", "mcp-adapter-assignments.txt"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(assignments), "Config: &mcp.BoilerConfig{") {
		t.Fatal("the pin's adapter does not always allocate the boiler config record")
	}

	want := make(map[string]any, len(record.Envelope.Data)+1)
	for key, value := range record.Envelope.Data {
		want[key] = value
	}
	state, ok := want["state"].(map[string]any)
	if !ok {
		t.Fatalf("recorded state: %v", want["state"])
	}
	renamed := make(map[string]any, len(state))
	for key, value := range state {
		if key == "flow_temperature_c" {
			key = "flow_temperature"
		}
		renamed[key] = value
	}
	want["state"] = renamed
	want["config"] = map[string]any{}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("%s (%s):\ngot  %s\nwant the recording with the two listed pin differences: %v", record.Name, record.Fixture, encoded, want)
	}
}
