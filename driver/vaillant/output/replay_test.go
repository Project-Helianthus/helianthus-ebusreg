package output

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/Project-Helianthus/helianthus-ebusreg/vaillant/productids"
)

type recordedCase struct {
	Clock struct {
		StoreNow time.Time `json:"store_now"`
	}
	Input  map[string]json.RawMessage
	Output struct {
		GraphQL      map[string]json.RawMessage
		MCP          map[string]json.RawMessage
		Cache        map[string]json.RawMessage
		Decisions    []json.RawMessage
		WriteResults []json.RawMessage `json:"write_results"`
		Observations []string
	}
}

func readFixture[T any](t *testing.T, file string) T {
	t.Helper()
	data, err := os.ReadFile(filepath.Join("testdata", file))
	if err != nil {
		t.Fatal(err)
	}
	return decodeRecorded[T](t, data)
}

func decodeRecorded[T any](t *testing.T, data json.RawMessage) T {
	t.Helper()
	var value T
	if err := json.Unmarshal(data, &value); err != nil {
		t.Fatal(err)
	}
	return value
}

func assertRecorded(t *testing.T, label string, got any, want json.RawMessage) {
	t.Helper()
	encoded, err := json.Marshal(got)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(decodeRecorded[any](t, encoded), decodeRecorded[any](t, want)) {
		t.Errorf("%s:\ngot  %s\nwant %s", label, encoded, want)
	}
}

type replayPlanes struct {
	graphql map[string]any
	mcp     map[string]any
}

func (planes replayPlanes) add(key string, graphql, mcp any) {
	planes.graphql[key] = graphql
	if key == "fm5_semantic_mode" {
		key = "fm5_mode"
	}
	planes.mcp[key] = mcp
}

func TestVaillantTypedOutputReplay(t *testing.T) {
	fields := readFixture[map[string]map[string]string](t, "graphql-go-fields.json")
	vectors := readFixture[map[string]recordedCase](t, "typed-view-vectors.json")
	for name, vector := range vectors {
		t.Run(name, func(t *testing.T) {
			planes := replayPlanes{make(map[string]any), make(map[string]any)}
			switch {
			case strings.HasPrefix(name, "cache-"):
				replayCache(t, vector, planes)
			case name == "schedule-write-mappings":
				replayScheduleWrite(t, vector)
			case name == "circuit-managing-device-derivation":
				replayManagingDevice(t, vector)
			case name == "regulator-capability-mapping":
				replayRegulatorCapability(t, vector, planes)
			default:
				replayView(t, vector.Input, planes)
			}
			if len(planes.mcp) != len(vector.Output.MCP) || len(planes.graphql) != len(vector.Output.GraphQL) {
				t.Fatalf("plane counts: got %d MCP and %d GraphQL; want %d and %d", len(planes.mcp), len(planes.graphql), len(vector.Output.MCP), len(vector.Output.GraphQL))
			}
			for key, want := range vector.Output.MCP {
				assertRecorded(t, "MCP "+key, planes.mcp[key], want)
			}
			for key, want := range vector.Output.GraphQL {
				typeName := recordedGraphQLType(key)
				assertRecorded(t, "GraphQL "+key, graphqlReplayValue(reflect.ValueOf(planes.graphql[key]), typeName, fields), want)
			}
		})
	}
}

func recordedGraphQLType(key string) string {
	switch key {
	case "zones":
		return "zone"
	case "circuits":
		return "circuit_status"
	case "dhw":
		return "dhw_status"
	case "solar":
		return "solar_status"
	case "cylinders":
		return "cylinder_status"
	case "system":
		return "system_status"
	case "schedules":
		return "schedule_status"
	default:
		return key
	}
}

func graphqlReplayValue(value reflect.Value, record string, fields map[string]map[string]string) any {
	if !value.IsValid() {
		return nil
	}
	if value.Kind() == reflect.Pointer {
		if value.IsNil() {
			return nil
		}
		value = value.Elem()
	}
	if value.Kind() == reflect.Slice {
		out := make([]any, 0, value.Len())
		for i := 0; i < value.Len(); i++ {
			out = append(out, graphqlReplayValue(value.Index(i), record, fields))
		}
		return out
	}
	if value.Kind() == reflect.String && value.Len() == 0 {
		return nil
	}
	if value.Kind() != reflect.Struct || fields[record] == nil {
		return value.Interface()
	}
	out := make(map[string]any)
	for field, goField := range fields[record] {
		v := value.FieldByName(goField)
		if record == "zone_config" && field == "allowed_modes" && v.Len() == 0 {
			// This default belongs to the old GraphQL projection, not the record.
			out[field] = []string{"off", "auto", "heat"}
			continue
		}
		out[field] = graphqlReplayValue(v, graphqlChildRecord(record, field), fields)
	}
	return out
}

func graphqlChildRecord(record, field string) string {
	switch record {
	case "zone", "dhw_status", "circuit_status", "boiler_status", "system_status":
		switch field {
		case "state", "config", "properties", "diagnostics", "managing_device":
			return strings.TrimSuffix(record, "_status") + "_" + field
		}
	case "system_config":
		if field == "dhw_outdoor_bivalence_temperature_limit_bound" {
			return "system_config_bound"
		}
	case "schedule_status":
		if field == "programs" {
			return "schedule_program"
		}
	case "schedule_program":
		if field == "config" {
			return "schedule_config"
		}
		if field == "days" {
			return "schedule_day_program"
		}
	case "schedule_day_program":
		if field == "slots" {
			return "schedule_timer_slot"
		}
	}
	return ""
}

func replayManagingDevice(t *testing.T, vector recordedCase) {
	t.Helper()
	inputs := decodeRecorded[[]struct {
		Scheme        *uint16       `json:"system_scheme"`
		Configuration *uint16       `json:"module_configuration_vr71"`
		Mode          NativeFM5Mode `json:"fm5_mode"`
	}](t, vector.Input["pairings"])
	if len(inputs) != len(vector.Output.Decisions) {
		t.Fatal("pairing count mismatch")
	}
	for i, input := range inputs {
		snapshot := &NativeSystemSnapshot{SystemScheme: input.Scheme, ModuleConfigurationVR71: input.Configuration}
		got := map[string]any{
			"managed_by_function_module": CircuitsManagedByFunctionModule(snapshot, input.Mode),
			"managing_device":            DeriveCircuitManagingDevice(snapshot, input.Mode),
		}
		assertRecorded(t, "circuit managing device", got, vector.Output.Decisions[i])
	}
}

func replayRegulatorCapability(t *testing.T, vector recordedCase, planes replayPlanes) {
	t.Helper()
	inputs := decodeRecorded[[]struct {
		Capability productids.ControllerCapability `json:"native_value"`
	}](t, vector.Input["capabilities"])
	if len(inputs) != len(vector.Output.Decisions) {
		t.Fatal("capability count mismatch")
	}
	for i, input := range inputs {
		mapped := RegulatorCapabilityFromNative(input.Capability)
		assertRecorded(t, "regulator capability", map[string]any{"published": mapped}, vector.Output.Decisions[i])
		planes.graphql["vaillant_regulator_capability"] = mapped
	}
}

func replayScheduleWrite(t *testing.T, vector recordedCase) {
	t.Helper()
	requests := decodeRecorded[[][]MCPTimeProgramSlot](t, vector.Input["slot_requests"])
	if len(requests) != len(vector.Output.Decisions) {
		t.Fatal("request count mismatch")
	}
	for i, request := range requests {
		assertRecorded(t, "slot request", map[string]any{"mcp_slots": request, "native_slots": TimeProgramSlotsFromMCP(request)}, vector.Output.Decisions[i])
	}
	results := decodeRecorded[[]struct {
		Name   string
		Error  *string
		Native *NativeTimeProgramWriteResult `json:"native_result"`
	}](t, vector.Input["write_results"])
	if len(results) != len(vector.Output.WriteResults) {
		t.Fatal("result count mismatch")
	}
	for i, result := range results {
		var err error
		if result.Error != nil {
			err = errors.New(*result.Error)
		}
		converted, resultErr := MCPTimeProgramResult(result.Native, err)
		if resultErr != err {
			t.Errorf("%s: error identity changed", result.Name)
		}
		var resultError *string
		if resultErr != nil {
			message := resultErr.Error()
			resultError = &message
		}
		assertRecorded(t, result.Name, map[string]any{"name": result.Name, "error": resultError, "mcp_result": converted}, vector.Output.WriteResults[i])
	}
}
