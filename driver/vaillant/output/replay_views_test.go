package output

import (
	"encoding/json"
	"testing"
)

type recordedZoneView struct {
	Instance      byte
	Snapshot      *NativeZoneSnapshot
	RadioTemp     *float64 `json:"radio_room_temp_c"`
	RadioHumidity *float64 `json:"radio_room_humidity_pct"`
	RoomMapping   *uint16  `json:"inferred_room_mapping_raw"`
}

func replayZones(t *testing.T, input json.RawMessage) []Zone {
	t.Helper()
	views := decodeRecorded[[]recordedZoneView](t, input)
	zones := make([]Zone, 0, len(views))
	for _, v := range views {
		if v.Snapshot == nil {
			continue
		}
		zones = append(zones, ZoneFromView(NativeZoneView{
			Instance: v.Instance, Snapshot: v.Snapshot, RadioRoomTempC: v.RadioTemp,
			RadioRoomHumidityPct: v.RadioHumidity, InferredRoomMappingRaw: v.RoomMapping,
		}))
	}
	return zones
}

// The recorded native view uses the original Go field names for these two values.
type recordedSystemSnapshot struct {
	NativeSystemSnapshot
	DailyAverage *float64 `json:"OutdoorTemperatureAvg24h"`
	ModuleCode   *uint16  `json:"ModuleConfigurationVR71"`
}

func (recorded *recordedSystemSnapshot) native() *NativeSystemSnapshot {
	if recorded == nil {
		return nil
	}
	recorded.OutdoorTemperatureDailyAverage = recorded.DailyAverage
	recorded.FunctionModuleConfigurationCode = recorded.ModuleCode
	return &recorded.NativeSystemSnapshot
}

type recordedCircuitsView struct {
	Circuits []*NativeCircuitSnapshot
	Managed  bool `json:"managed_by_function_module"`
}

func replayCircuits(view recordedCircuitsView, planes replayPlanes) {
	circuits := make([]CircuitStatus, 0, len(view.Circuits))
	for _, circuit := range view.Circuits {
		if circuit != nil {
			circuits = append(circuits, CircuitStatusFromSnapshot(circuit, view.Managed))
		}
	}
	planes.add("circuits", circuits, MCPCircuitsFromGraphQL(circuits))
}

func replayView(t *testing.T, input map[string]json.RawMessage, planes replayPlanes) {
	t.Helper()
	for field, payload := range input {
		switch field {
		case "zone_views":
			zones := replayZones(t, payload)
			planes.add("zones", zones, MCPZonesFromGraphQL(zones))
		case "dhw_view":
			view := decodeRecorded[NativeDhwView](t, payload)
			dhw := DhwStatusFromSnapshot(view.Snapshot)
			planes.add("dhw", dhw, MCPDHWFromGraphQL(dhw))
		case "circuits_view":
			replayCircuits(decodeRecorded[recordedCircuitsView](t, payload), planes)
		case "circuits_views":
			// The second input is a cache withdrawal; only the first feeds the mapping.
			views := decodeRecorded[[]recordedCircuitsView](t, payload)
			if len(views) != 2 || len(views[1].Circuits) != 0 {
				t.Fatal("unexpected managed-circuit replay sequence")
			}
			replayCircuits(views[0], planes)
		case "boiler_snapshot":
			boiler := BoilerStatusFromSnapshot(decodeRecorded[*NativeBoilerSnapshot](t, payload))
			planes.add("boiler_status", boiler, MCPBoilerStatusFromGraphQL(boiler))
		case "system_view":
			view := decodeRecorded[struct {
				Snapshot *recordedSystemSnapshot
				Bound    ConfigBound `json:"dhw_bivalence_bound"`
			}](t, payload)
			system := SystemStatusFromView(NativeSystemView{Snapshot: view.Snapshot.native(), DhwBivalenceBound: view.Bound})
			planes.add("system", system, MCPSystemFromGraphQL(system))
		case "fm5_view":
			replayFunctionModule(t, payload, planes)
		case "programs":
			programs := ScheduleProgramsFromNative(decodeRecorded[[]NativeScheduleProgram](t, payload))
			status := &ScheduleStatus{Programs: programs}
			planes.add("schedules", status, MCPSchedulesFromGraphQL(status))
		case "energy_totals":
			// Energy accounting and acquisition are not moved here. Feed the recorded
			// GraphQL view to the pure MCP adapter using the frozen field selectors.
			fields := readFixture[map[string]map[string]string](t, "graphql-go-fields.json")
			view := decodeRecorded[any](t, payload)
			named := energyGoFields(view, "energy_totals", fields)
			encoded, err := json.Marshal(named)
			if err != nil {
				t.Fatal(err)
			}
			totals := decodeRecorded[EnergyTotals](t, encoded)
			planes.mcp[field] = MCPEnergyTotalsFromGraphQL(&totals)
		default:
			t.Fatalf("unhandled native fixture field %q", field)
		}
	}
}

func replayFunctionModule(t *testing.T, payload json.RawMessage, planes replayPlanes) {
	t.Helper()
	view := decodeRecorded[NativeFunctionModuleView](t, payload)
	if !view.Verdict.Valid() {
		t.Fatal("invalid recorded verdict")
	}
	published := PublishedFunctionModuleInterpretation(view.Verdict)
	if err := published.Validate(); err != nil {
		t.Fatal(err)
	}
	mcp := MCPFunctionModuleInterpretationFromGraphQL(published)
	if err := mcp.validate(); err != nil {
		t.Fatal(err)
	}
	planes.add("fm5_interpretation", published, mcp)
	planes.add("fm5_semantic_mode", published.Mode, MCPFunctionModuleSemanticModeFromGraphQL(published.Mode))
	// The recorder primes the provider with its empty structural planes.
	solar, cylinders := &SolarStatus{}, []CylinderStatus{}
	if published.Mode == FunctionModuleSemanticModeInterpreted {
		solar = SolarStatusFromSnapshot(view.Solar)
		cylinders = CylinderStatusesFromSnapshots(view.Cylinders)
	}
	planes.add("solar", solar, MCPSolarFromGraphQL(solar))
	planes.add("cylinders", cylinders, MCPCylindersFromGraphQL(cylinders))
}

func energyGoFields(value any, record string, fields map[string]map[string]string) any {
	switch value := value.(type) {
	case map[string]any:
		out := make(map[string]any, len(value))
		child := map[string]string{"energy_totals": "energy_channel", "energy_channel": "energy_series", "energy_series": "energy_point_meta"}[record]
		for key, field := range value {
			out[fields[record][key]] = energyGoFields(field, child, fields)
		}
		return out
	case []any:
		out := make([]any, len(value))
		for i, field := range value {
			out[i] = energyGoFields(field, record, fields)
		}
		return out
	default:
		return value
	}
}
