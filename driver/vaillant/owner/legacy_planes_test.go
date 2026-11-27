package owner

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/Project-Helianthus/helianthus-ebusgo/protocol"

	catalogvaillant "github.com/Project-Helianthus/helianthus-ebusreg/catalog/vaillant"
	"github.com/Project-Helianthus/helianthus-ebusreg/driver/vaillant/output"
)

func TestVaillantLegacyPlanesHeating(t *testing.T) {
	var frames []protocol.Frame
	poller := legacyPoller(&frames)
	marker := map[byte]*vaillantZoneSnapshot{0x08: &vaillantZoneSnapshot{}}
	poller.zones = marker
	record, err := fillLegacyPlanes(context.Background(), poller, legacyBaseAdmission(), nil)
	if err != nil {
		t.Fatal(err)
	}
	if poller.zones == nil || len(poller.zones) != 1 {
		t.Fatal("legacy fill replaced the controller zone state")
	}
	if len(frames) != 4 {
		t.Fatalf("frames = %d", len(frames))
	}
	for index, frame := range frames {
		if frame.Primary != 0xB5 || frame.Secondary != 0x04 || frame.Target != 0x08 || frame.Source != adapterSource {
			t.Fatalf("frame %d = %+v", index, frame)
		}
	}
	if !reflect.DeepEqual(frames[0].Data, []byte{0x0D}) || !reflect.DeepEqual(frames[1].Data, []byte{0x09}) || !reflect.DeepEqual(frames[2].Data, []byte{0x0D}) || !reflect.DeepEqual(frames[3].Data, []byte{0x09}) {
		t.Fatalf("payloads = %x %x %x %x", frames[0].Data, frames[1].Data, frames[2].Data, frames[3].Data)
	}
	if len(record.Planes) != 2 || record.Planes[0].Name != "heating" || record.Planes[1].Name != "dhw" {
		t.Fatalf("planes = %+v", record.Planes)
	}
	assertLegacyNumber(t, record.Planes[0].Reads[0].Fields, "status", 4)
	assertLegacyNumber(t, record.Planes[0].Reads[1].Fields, "flow_temp", 1)
	assertLegacyNumber(t, record.Planes[0].Reads[1].Fields, "return_temp", 2)
	assertLegacyNumber(t, record.Planes[0].Reads[1].Fields, "pump_status", 3)
	assertLegacyNumber(t, record.Planes[1].Reads[0].Fields, "status", 4)
	assertLegacyNumber(t, record.Planes[1].Reads[1].Fields, "dhw_temp", 1)
	assertLegacyNumber(t, record.Planes[1].Reads[1].Fields, "target_temp", 2)
	assertLegacyNumber(t, record.Planes[1].Reads[1].Fields, "mode", 3)
	for _, plane := range record.Planes {
		for _, read := range plane.Reads {
			for _, field := range read.Fields {
				if field.Level != legacyLevelUnknown {
					t.Fatalf("field %s level = %s", field.Name, field.Level)
				}
			}
		}
	}
	if len(record.Operations) != 2 || record.Operations[0].Name != legacyTargetMethod || record.Operations[0].Namespace != legacyNamespace || record.Operations[0].Unresolved != legacyUnresolvedBody || record.Operations[1].Plane != "dhw" {
		t.Fatalf("operations = %+v", record.Operations)
	}
	before := len(frames)
	if err := callLegacyOperation(context.Background(), poller, "heating", legacyTargetMethod, 0x08, "7602", []byte{1, 2, 3}); !errors.Is(err, errLegacyLayout) || len(frames) != before {
		t.Fatal("target temperature body was sent")
	}
	for _, name := range []string{"holiday", "program_blocks", "opaque_programming", "legacy_timer"} {
		if !legacyCataloguePresent(record.Catalogues, name) {
			t.Fatalf("catalogue %s missing from %v", name, record.Catalogues)
		}
	}
}

func TestVaillantLegacyPlanesAdmission(t *testing.T) {
	extended := legacyBaseAdmission()
	extended.ExtendedRegister = true
	frames, record, err := legacyFill(extended, nil)
	if err != nil || len(frames) != 0 || len(record.Planes) != 0 || !record.ExtendedRegister {
		t.Fatal("ExtendedRegister mutation still read heating; dropping ExtendedRegister from the guard would keep the original planes")
	}
	heatingOnly := legacyBaseAdmission()
	heatingOnly.AnswersHotWater = false
	frames, record, err = legacyFill(heatingOnly, nil)
	if err != nil || len(record.Planes) != 1 || record.Planes[0].Name != "heating" || len(frames) != 2 {
		t.Fatal("AnswersHotWater mutation still read hot water; dropping AnswersHotWater from the guard would keep the original planes")
	}
	hotWaterOnly := legacyBaseAdmission()
	hotWaterOnly.AnswersHeating = false
	frames, record, err = legacyFill(hotWaterOnly, nil)
	if err != nil || len(record.Planes) != 1 || record.Planes[0].Name != "dhw" || len(frames) != 2 {
		t.Fatal("AnswersHeating mutation still read heating; dropping AnswersHeating from the guard would keep the original planes")
	}
	solarOnly := legacyBaseAdmission()
	solarOnly.ExtendedRegister = true
	solarOnly.AnswersHeating = false
	solarOnly.AnswersHotWater = false
	solarOnly.StandaloneSolar = true
	frames, record, err = legacyFill(solarOnly, map[string]any{"period": 1, "source": 2, "usage": 3})
	if err != nil || len(frames) != 0 || len(record.Planes) != 1 || record.Planes[0].Name != "solar" || len(record.Planes[0].Reads) != 3 {
		t.Fatal("StandaloneSolar mutation published no solar view; dropping StandaloneSolar would leave the station empty")
	}
	for _, read := range record.Planes[0].Reads {
		if read.Sent || read.Unresolved != legacyUnresolvedBody {
			t.Fatalf("solar read was sent: %+v", read)
		}
	}
	version := legacyBaseAdmission()
	version.HardwareVersion = "7603"
	params := map[string]any{"period": 1, "source": 2, "usage": 3}
	frames, record, err = legacyFill(version, params)
	if err != nil || len(frames) != 6 || frames[2].Secondary != 0x16 || !reflect.DeepEqual(frames[2].Data, []byte{1, 2, 3}) || frames[5].Secondary != 0x16 {
		t.Fatal("HardwareVersion mutation omitted energy statistics; dropping HardwareVersion would keep the original read set")
	}
	assertLegacyNumber(t, record.Planes[0].Reads[2].Fields, "energy", 2)
	frames, record, err = legacyFill(version, nil)
	if err != nil || len(frames) != 4 || record.Planes[0].Reads[2].Method != legacyEnergyMethod || record.Planes[0].Reads[2].Sent || record.Planes[0].Reads[2].Unresolved != legacyUnresolvedBody {
		t.Fatal("omitted energy parameters were still sent")
	}
	frames, _, err = legacyFill(version, map[string]any{"period": 1})
	if err == nil || len(frames) != 2 {
		t.Fatalf("partial energy parameters sent %d frames err=%v", len(frames), err)
	}
	above := legacyBaseAdmission()
	above.HardwareVersion = "76030"
	frames, record, err = legacyFill(above, params)
	if err != nil || len(frames) != 4 || legacyRead(record, "heating", legacyEnergyMethod) != nil {
		t.Fatal("HardwareVersion 76030 mutation still selected energy statistics")
	}
	below := legacyBaseAdmission()
	below.HardwareVersion = "7602"
	frames, record, err = legacyFill(below, params)
	if err != nil || len(frames) != 4 || legacyRead(record, "heating", legacyEnergyMethod) != nil {
		t.Fatal("HardwareVersion 7602 mutation still selected energy statistics")
	}
	other := legacyBaseAdmission()
	other.Manufacturer = "Other"
	frames, record, err = legacyFill(other, nil)
	if err != nil || len(frames) != 4 || len(record.Planes) != 2 {
		t.Fatal("Manufacturer mutation suppressed a plane the admission says answers; dropping Manufacturer would keep the original planes")
	}
	controller := legacyBaseAdmission()
	controller.Target = 0x10
	frames, record, err = legacyFill(controller, nil)
	if err != nil || len(frames) != 4 || frames[0].Target != 0x10 {
		t.Fatal("Target mutation left the frame target unchanged; dropping Target would keep the original controller")
	}
	names := legacyFieldNames(record.Planes[0].Reads[1].Fields)
	if !reflect.DeepEqual(names, []string{"room_temp", "target_temp", "mode"}) {
		t.Fatalf("controller parameters = %v", names)
	}
}

func TestVaillantLegacyPlanesReplies(t *testing.T) {
	status, err := decodeLegacyPlane("solar", "get_status", 0x15, "", []byte{0x00, 0x01, 0x00, 0x02})
	if err != nil {
		t.Fatal(err)
	}
	assertLegacyNumber(t, status.Fields, "collector_temp", 1)
	assertLegacyNumber(t, status.Fields, "tank_temp", 2)
	yield, err := decodeLegacyPlane("solar", "get_solar_yield", 0x15, "", []byte{0x02, 0x00})
	if err != nil {
		t.Fatal(err)
	}
	assertLegacyNumber(t, yield.Fields, "yield", 2)
	parameters, err := decodeLegacyPlane("solar", "get_parameters", 0x15, "", []byte{0x04, 0x00, 0x01})
	if err != nil {
		t.Fatal(err)
	}
	assertLegacyNumber(t, parameters.Fields, "pump_speed", 4)
	assertLegacyNumber(t, parameters.Fields, "delta_temp", 1)
	short, err := decodeLegacyPlane("heating", "get_parameters", 0x08, "7602", []byte{0x00})
	if err == nil || short.Unresolved != legacyUnresolvedBody || len(short.Fields) != 0 || !reflect.DeepEqual(short.Raw, []byte{0x00}) {
		t.Fatalf("short parameters repaired the reply: %+v err=%v", short, err)
	}
	if _, err := decodeLegacyPlane("heating", legacyEnergyMethod, 0x08, "76030", []byte{0x02, 0x00}); !errors.Is(err, errLegacyLayout) {
		t.Fatal("an absent energy method still decoded")
	}
}

func TestVaillantLegacyPlanesIdentity(t *testing.T) {
	candidates, err := catalogvaillant.B504Candidates()
	if err != nil || len(candidates.B504CandidateEntries) == 0 {
		t.Fatal(err)
	}
	heatingKey := candidates.B504CandidateEntries[0].B504CandidateKey
	grade, found, err := gradeHeatingCandidate(heatingKey)
	if err != nil || !found || grade == "" {
		t.Fatalf("heating candidate found=%v grade=%s err=%v", found, grade, err)
	}
	fields := []output.LegacyField{{Name: "flow_temp", Level: legacyLevelUnknown}, {Name: heatingKey.B504CandidateField, Level: legacyLevelUnknown}}
	annotated := annotateLegacyFields(fields, heatingKey.B504CandidateField, grade, "degrees Celsius", true)
	if annotated[0].Level != legacyLevelUnknown || annotated[1].Level != grade || annotated[1].Unit != "degrees Celsius" {
		t.Fatalf("annotation crossed field names: %+v", annotated)
	}
	if untouched := annotateLegacyFields(fields, heatingKey.B504CandidateField, grade, "degrees Celsius", false); untouched[1].Level != legacyLevelUnknown {
		t.Fatal("a missing catalogue row annotated a field")
	}
	for _, mutation := range []struct {
		field  string
		mutate func(*catalogvaillant.B504CandidateKey)
	}{
		{"family", func(key *catalogvaillant.B504CandidateKey) { key.B504CandidateFamily = "absent-family" }},
		{"variant", func(key *catalogvaillant.B504CandidateKey) { key.B504CandidateVariant = "absent-variant" }},
		{"primary", func(key *catalogvaillant.B504CandidateKey) { key.B504CandidatePrimary++ }},
		{"secondary", func(key *catalogvaillant.B504CandidateKey) {
			key.B504CandidateSecondary = mutateOptionalOctet(key.B504CandidateSecondary)
		}},
		{"operation", func(key *catalogvaillant.B504CandidateKey) {
			key.B504CandidateOperation = mutateOptionalOctet(key.B504CandidateOperation)
		}},
		{"direction", func(key *catalogvaillant.B504CandidateKey) { key.B504CandidateDirection = "sideways" }},
		{"field", func(key *catalogvaillant.B504CandidateKey) { key.B504CandidateField = "absent-field" }},
	} {
		mutated := cloneJSON(t, heatingKey)
		mutation.mutate(&mutated)
		if _, found, err := gradeHeatingCandidate(mutated); err != nil || found {
			t.Fatalf("%s mutation still matched; dropping %s from equality would keep the original row", mutation.field, mutation.field)
		}
	}
	operatingKey := catalogvaillant.B505Candidates().B505CandidateEntries[0].B505CandidateKey
	if _, found := gradeOperatingCandidate(operatingKey); !found {
		t.Fatal("unmutated operating candidate missed")
	}
	for _, mutation := range []struct {
		field  string
		mutate func(*catalogvaillant.B505CandidateKey)
	}{
		{"family", func(key *catalogvaillant.B505CandidateKey) { key.B505CandidateFamily = "absent-family" }},
		{"variant", func(key *catalogvaillant.B505CandidateKey) { key.B505CandidateVariant = "absent-variant" }},
		{"primary", func(key *catalogvaillant.B505CandidateKey) { key.B505CandidatePrimary++ }},
		{"secondary", func(key *catalogvaillant.B505CandidateKey) {
			key.B505CandidateSecondary = mutateByte(key.B505CandidateSecondary)
		}},
		{"operation", func(key *catalogvaillant.B505CandidateKey) {
			key.B505CandidateOperation = mutateByte(key.B505CandidateOperation)
		}},
		{"direction", func(key *catalogvaillant.B505CandidateKey) {
			key.B505CandidateDirection = mutateString(key.B505CandidateDirection)
		}},
		{"field", func(key *catalogvaillant.B505CandidateKey) { key.B505CandidateField = "absent-field" }},
	} {
		mutated := cloneJSON(t, operatingKey)
		mutation.mutate(&mutated)
		if _, found := gradeOperatingCandidate(mutated); found {
			t.Fatalf("%s mutation still matched; dropping %s from equality would keep the original row", mutation.field, mutation.field)
		}
	}
	assertLegacyRoleMutations(t)
}

func TestVaillantLegacyPlanesSurfaces(t *testing.T) {
	var frames []protocol.Frame
	poller := legacyPoller(&frames)
	record, err := fillLegacyPlanes(context.Background(), poller, legacyBaseAdmission(), nil)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(record.GraphQL(), record.MCP()) || record.Ambiguity != legacyAmbiguity {
		t.Fatal("legacy surfaces diverged")
	}
	status := &output.SystemStatus{}
	status.SetLegacyPlanes(&record)
	mcp := output.MCPSystemFromGraphQL(status)
	if !reflect.DeepEqual(status.LegacyPlanes(), mcp.LegacyPlanes()) {
		t.Fatal("system legacy planes diverged")
	}
	if retainedLegacyPlanes(poller) == nil || retainedLegacyPlanes(poller).Planes[0].Name != "heating" {
		t.Fatal("retained legacy planes lost heating")
	}
}

func legacyBaseAdmission() legacyAdmission {
	return legacyAdmission{
		AnswersHeating:  true,
		AnswersHotWater: true,
		HardwareVersion: "7602",
		Manufacturer:    "Vaillant",
		Target:          0x08,
	}
}

func legacyFill(admission legacyAdmission, energyParams map[string]any) ([]protocol.Frame, output.LegacyPlaneRecord, error) {
	var frames []protocol.Frame
	record, err := fillLegacyPlanes(context.Background(), legacyPoller(&frames), admission, energyParams)
	return frames, record, err
}

func legacyPoller(frames *[]protocol.Frame) *vaillantSemanticPoller {
	return &vaillantSemanticPoller{
		source: adapterSource,
		sendFrameFn: func(_ context.Context, frame protocol.Frame) (*protocol.Frame, error) {
			*frames = append(*frames, frame)
			var data []byte
			switch {
			case frame.Secondary == 0x16:
				data = []byte{0x02, 0x00}
			case len(frame.Data) > 0 && frame.Data[0] == 0x0D:
				data = []byte{0x04}
			case len(frame.Data) > 0 && frame.Data[0] == 0x09:
				data = []byte{0x00, 0x01, 0x00, 0x02, 0x03}
			default:
				data = []byte{0x00}
			}
			return &protocol.Frame{Data: append([]byte(nil), data...)}, nil
		},
	}
}

func assertLegacyNumber(t *testing.T, fields []output.LegacyField, name string, want float64) {
	t.Helper()
	for _, field := range fields {
		if field.Name != name {
			continue
		}
		if field.Value == nil || !field.Valid || *field.Value != want {
			t.Fatalf("%s = %+v", name, field)
		}
		return
	}
	t.Fatalf("missing %s", name)
}

func legacyFieldNames(fields []output.LegacyField) []string {
	names := make([]string, 0, len(fields))
	for _, field := range fields {
		names = append(names, field.Name)
	}
	return names
}

func legacyRead(record output.LegacyPlaneRecord, planeName, method string) *output.LegacyRead {
	for _, plane := range record.Planes {
		if plane.Name != planeName {
			continue
		}
		for index := range plane.Reads {
			if plane.Reads[index].Method == method {
				return &plane.Reads[index]
			}
		}
	}
	return nil
}

func legacyCataloguePresent(names []string, want string) bool {
	for _, name := range names {
		if name == want {
			return true
		}
	}
	return false
}
