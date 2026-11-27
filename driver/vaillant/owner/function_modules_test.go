package owner

import (
	"context"
	"encoding/json"
	"errors"
	"reflect"
	"strings"
	"testing"

	"github.com/Project-Helianthus/helianthus-ebusgo/protocol"

	catalogvaillant "github.com/Project-Helianthus/helianthus-ebusreg/catalog/vaillant"
	"github.com/Project-Helianthus/helianthus-ebusreg/driver/vaillant/output"
)

func TestVaillantFunctionModuleIdentification(t *testing.T) {
	payload := moduleIdentification(0xB5, "VR_71", [2]byte{1, 2}, [2]byte{3, 4})
	var frames []protocol.Frame
	poller := modulePoller(&frames, []byte{0x11})
	record, ok := moduleFromIdentification(payload, 0x11)
	if !ok || record.Family != moduleFamilyTwelve || record.Address != 0x11 || len(frames) != 0 {
		t.Fatalf("record=%+v ok=%v frames=%d", record, ok, len(frames))
	}
	if _, ok := moduleFromIdentification(moduleIdentification(0xB5, "ZZZZZ", [2]byte{1, 2}, [2]byte{3, 4}), 0x26); ok {
		t.Fatal("address alone identified a module")
	}
	if _, ok := moduleFromIdentification(payload[:9], 0x26); ok {
		t.Fatal("short identification became a module")
	}
	otherAddress, ok := moduleFromIdentification(payload, 0x03)
	if !ok || otherAddress.Family != record.Family {
		t.Fatal("address mutation changed the family")
	}
	manufacturer := moduleIdentification(0x00, "VR_71", [2]byte{1, 2}, [2]byte{3, 4})
	changed, ok := moduleFromIdentification(manufacturer, 0x11)
	if !ok || changed.Family != moduleFamilyTwelve || changed.Manufacturer != 0 {
		t.Fatal("manufacturer mutation changed the family")
	}
	for index := 0; index < 5; index++ {
		mutated := moduleIdentification(0xB5, "VR_71", [2]byte{1, 2}, [2]byte{3, 4})
		mutated[1+index] = 'Z'
		if _, ok := moduleFromIdentification(mutated, 0x11); ok {
			t.Fatalf("device id byte %d mutation still matched", index)
		}
	}
	for _, index := range []int{6, 7} {
		mutated := append([]byte(nil), payload...)
		mutated[index] ^= 0x7F
		got, ok := moduleFromIdentification(mutated, 0x11)
		if !ok || got.Family != moduleFamilyTwelve {
			t.Fatalf("software byte %d changed the family", index)
		}
	}
	for _, index := range []int{8, 9} {
		mutated := append([]byte(nil), payload...)
		mutated[index] ^= 0x7F
		got, ok := moduleFromIdentification(mutated, 0x11)
		if !ok || got.Family != moduleFamilyTwelve {
			t.Fatalf("hardware byte %d changed the family", index)
		}
	}
	bare, ok := moduleFromIdentification(moduleIdentification(0xB5, "VR66", [2]byte{9, 8}, [2]byte{7, 6}), 0x52)
	if !ok || bare.Family != moduleFamilyIdentityOnly || len(bare.Operations) != 0 || len(bare.Sensors) != 1 {
		t.Fatalf("bare module = %+v", bare)
	}
	padded := moduleIdentification(0xB5, "VR66", [2]byte{9, 8}, [2]byte{7, 6})
	padded[5] = 0
	nulled, ok := moduleFromIdentification(padded, 0x52)
	if !ok || nulled.DeviceID != "VR66" {
		t.Fatalf("nul padding = %+v", nulled)
	}
	if len(frames) != 0 {
		t.Fatal("identification sent a frame")
	}
	_ = poller
}

func TestVaillantFunctionModuleReads(t *testing.T) {
	var frames []protocol.Frame
	poller := modulePoller(&frames, []byte{0x06, 0x07, 0x08})
	readings, err := readModuleSensors(context.Background(), poller, 0x11, moduleIDTwelve)
	if err != nil || len(readings) != 2 || len(frames) != 2 {
		t.Fatalf("readings=%d frames=%d err=%v", len(readings), len(frames), err)
	}
	if frames[0].Primary != modulePrimary || frames[0].Secondary != moduleSecondary || !reflect.DeepEqual(frames[0].Data, []byte{0x06}) {
		t.Fatalf("first sensor frame = %+v", frames[0])
	}
	if !reflect.DeepEqual(frames[1].Data, []byte{0x07}) || readings[0].Level != moduleLevelUnknown || readings[0].Placement != modulePlacementUnknown || readings[0].Unresolved != "unresolved" {
		t.Fatalf("sensor reading = %+v", readings[0])
	}
	before := len(frames)
	if err := callModuleOperation(context.Background(), poller, 0x11, moduleIDSix, "SetActorState", nil, make([]byte, 6)); !errors.Is(err, errModuleLayout) || len(frames) != before {
		t.Fatal("short actor body was sent")
	}
	raw := []byte{1, 2, 3, 4, 5, 6, 7}
	if err := callModuleOperation(context.Background(), poller, 0x11, moduleIDSix, "SetActorState", nil, raw); err != nil {
		t.Fatal(err)
	}
	sent := frames[len(frames)-1]
	if sent.Data[0] != 0x01 || !reflect.DeepEqual(sent.Data[1:], raw) {
		t.Fatalf("actor payload = %x", sent.Data)
	}
	before = len(frames)
	circuit := byte(0)
	if err := callModuleOperation(context.Background(), poller, 0x11, moduleIDSix, "Mc1Operation", &circuit, []byte{1, 2}); !errors.Is(err, errModuleLayout) || len(frames) != before {
		t.Fatal("open mixer body was sent")
	}
	if err := callModuleOperation(context.Background(), poller, 0x11, moduleIDSix, "Mc1Operation", &circuit, nil); err != nil {
		t.Fatal(err)
	}
	mixer := frames[len(frames)-1].Data
	if !reflect.DeepEqual(mixer, []byte{0x02, 0x00}) {
		t.Fatalf("mixer selector = %x", mixer)
	}
	missing := byte(2)
	before = len(frames)
	if err := callModuleOperation(context.Background(), poller, 0x11, moduleIDSix, "Mc3Operation", &missing, nil); !errors.Is(err, errModuleLayout) || len(frames) != before {
		t.Fatal("unknown circuit was sent")
	}
	before = len(frames)
	if err := callModuleOperation(context.Background(), poller, 0x11, moduleIDBare, "SetActorState", nil, raw); !errors.Is(err, errModuleLayout) || len(frames) != before {
		t.Fatal("identity-only module sent a write")
	}
	record, _ := moduleFromIdentification(moduleIdentification(0xB5, moduleIDTwelve, [2]byte{1, 2}, [2]byte{3, 4}), 0x11)
	var actor output.ModuleOperation
	for _, operation := range record.Operations {
		if operation.Name == "SetActorState" {
			actor = operation
		}
	}
	if actor.Namespace != moduleNamespace || !reflect.DeepEqual(actor.Alternatives, []string{"numeric_codec", "behavior"}) {
		t.Fatalf("actor operation = %+v", actor)
	}
	var filledFrames []protocol.Frame
	filled, err := fillFunctionModule(context.Background(), modulePoller(&filledFrames, []byte{0x10}), moduleIdentification(0xB5, moduleIDTwelve, [2]byte{1, 2}, [2]byte{3, 4}), 0x11)
	if err != nil || len(filled.Sensors) != 2 || len(filledFrames) != 2 || filledFrames[0].Data[0] != 0x06 || filledFrames[1].Data[0] != 0x07 {
		t.Fatalf("filled module sensors=%d frames=%d err=%v", len(filled.Sensors), len(filledFrames), err)
	}
	if retainedFunctionModule(modulePoller(&filledFrames, nil)) != nil {
		t.Fatal("a second poller inherited the filled module")
	}
}

func TestVaillantFunctionModuleFields(t *testing.T) {
	logical := make([]int, 18)
	encoded, err := json.Marshal(map[string]any{
		"context": map[string]any{
			"primary": 181, "secondary": 35, "operation": 3, "direction": "positive", "view": "status-zero-valid",
			"family": map[string]any{"opaque": []any{"synthetic-family", 0}}, "variant": []any{"synthetic-variant", nil},
		},
		"fields":    map[string]any{"logical_object": logical, "unmodelled": []int{0, 255, 7}},
		"arguments": map[string]any{},
	})
	if err != nil {
		t.Fatal(err)
	}
	interpreted, err := interpretModuleSensors(encoded)
	if err != nil {
		t.Fatal(err)
	}
	var sensor map[string]any
	if err := json.Unmarshal(interpreted, &sensor); err != nil || sensor["physical"] != modulePlacementUnknown {
		t.Fatalf("sensor physical = %v err=%v", sensor["physical"], err)
	}
	actuatorInput := map[string]any{
		"context":   map[string]any{"primary": 181, "secondary": 35, "operation": 1, "direction": "request", "view": "relay-six", "family": "synthetic-family", "variant": "opaque"},
		"fields":    map[string]any{"relays": []int{0, 0, 0, 0, 0, 0}, "pwm1": []int{201}},
		"arguments": map[string]any{"variant": "two-circuit", "opaque": map[string]any{"false": false, "zero": 0, "null": nil, "array": []any{0, false, nil}}},
	}
	actuatorJSON, err := json.Marshal(actuatorInput)
	if err != nil {
		t.Fatal(err)
	}
	actuator, err := interpretModuleActuators(actuatorJSON)
	if err != nil {
		t.Fatal(err)
	}
	text := string(actuator)
	if !strings.Contains(text, "unchanged") || !strings.Contains(text, "100.5") || !strings.Contains(text, `"physical":"Unknown"`) {
		t.Fatalf("pwm alternatives dropped: %s", text)
	}
	entries, err := catalogvaillant.B523ConfigMixerEntries()
	if err != nil {
		t.Fatal(err)
	}
	var setpoint catalogvaillant.B523ConfigMixerKey
	for _, entry := range entries {
		if entry.B523ConfigMixerKey.B523ConfigMixerVariant == "mixer-smaller" && entry.B523ConfigMixerKey.B523ConfigMixerField == "setpoint" {
			setpoint = entry.B523ConfigMixerKey
		}
	}
	mixer, found, err := moduleMixerField(setpoint)
	if err != nil || !found || len(mixer.B523ConfigMixerCandidates) != 2 || mixer.B523ConfigMixerSelectedInterpretation != modulePlacementUnknown {
		t.Fatalf("setpoint candidates=%d interpretation=%s found=%v err=%v", len(mixer.B523ConfigMixerCandidates), mixer.B523ConfigMixerSelectedInterpretation, found, err)
	}
}

func TestVaillantFunctionModuleIdentity(t *testing.T) {
	sensors, err := catalogvaillant.B523SensorFields()
	if err != nil || len(sensors.B523SensorFieldFields) == 0 {
		t.Fatal(err)
	}
	sensorKey := sensors.B523SensorFieldFields[0].B523SensorFieldKey
	if _, found, err := moduleSensorField(sensorKey); err != nil || !found {
		t.Fatalf("unmutated sensor found=%v err=%v", found, err)
	}
	for _, mutation := range []struct {
		field  string
		mutate func(*catalogvaillant.B523SensorKey)
	}{
		{"family", func(key *catalogvaillant.B523SensorKey) { key.B523SensorFieldFamily = "absent-family" }},
		{"variant", func(key *catalogvaillant.B523SensorKey) { key.B523SensorFieldVariant = "absent-variant" }},
		{"primary", func(key *catalogvaillant.B523SensorKey) {
			key.B523SensorFieldPrimary = mutateOptionalWide(key.B523SensorFieldPrimary)
		}},
		{"secondary", func(key *catalogvaillant.B523SensorKey) {
			key.B523SensorFieldSecondary = mutateOptionalWide(key.B523SensorFieldSecondary)
		}},
		{"operation", func(key *catalogvaillant.B523SensorKey) {
			key.B523SensorFieldOperation = mutateOptionalWide(key.B523SensorFieldOperation)
		}},
		{"direction", func(key *catalogvaillant.B523SensorKey) {
			key.B523SensorFieldDirection = mutateString(key.B523SensorFieldDirection)
		}},
		{"field", func(key *catalogvaillant.B523SensorKey) { key.B523SensorFieldField = "absent-field" }},
		{"opcode", func(key *catalogvaillant.B523SensorKey) {
			key.B523SensorFieldOpcode = mutateOptionalWide(key.B523SensorFieldOpcode)
		}},
		{"group", func(key *catalogvaillant.B523SensorKey) {
			key.B523SensorFieldGroup = mutateOptionalWide(key.B523SensorFieldGroup)
		}},
		{"instance", func(key *catalogvaillant.B523SensorKey) {
			key.B523SensorFieldInstance = mutateOptionalWide(key.B523SensorFieldInstance)
		}},
		{"register", func(key *catalogvaillant.B523SensorKey) {
			key.B523SensorFieldRegister = mutateOptionalWide(key.B523SensorFieldRegister)
		}},
	} {
		mutated := cloneJSON(t, sensorKey)
		mutation.mutate(&mutated)
		if _, found, err := moduleSensorField(mutated); err != nil || found {
			t.Fatalf("%s mutation still matched; dropping %s from equality would keep the original row", mutation.field, mutation.field)
		}
	}
	candidates, err := catalogvaillant.B523Candidates()
	if err != nil || len(candidates.B523CandidateEntries) == 0 {
		t.Fatal(err)
	}
	candidateKey := candidates.B523CandidateEntries[0].B523CandidateKey
	if _, found, err := moduleCandidateField(candidateKey); err != nil || !found {
		t.Fatalf("unmutated candidate found=%v err=%v", found, err)
	}
	for _, mutation := range []struct {
		field  string
		mutate func(*catalogvaillant.B523CandidateKey)
	}{
		{"family", func(key *catalogvaillant.B523CandidateKey) { key.B523CandidateFamily = "absent-family" }},
		{"variant", func(key *catalogvaillant.B523CandidateKey) {
			key.B523CandidateVariant = mutateString(key.B523CandidateVariant)
		}},
		{"primary", func(key *catalogvaillant.B523CandidateKey) { key.B523CandidatePrimary++ }},
		{"secondary", func(key *catalogvaillant.B523CandidateKey) {
			key.B523CandidateSecondary = mutateByte(key.B523CandidateSecondary)
		}},
		{"operation", func(key *catalogvaillant.B523CandidateKey) {
			key.B523CandidateOperation = mutateByte(key.B523CandidateOperation)
		}},
		{"direction", func(key *catalogvaillant.B523CandidateKey) { key.B523CandidateDirection = "sideways" }},
		{"field", func(key *catalogvaillant.B523CandidateKey) { key.B523CandidateField = "absent-field" }},
	} {
		mutated := cloneJSON(t, candidateKey)
		mutation.mutate(&mutated)
		if _, found, err := moduleCandidateField(mutated); err != nil || found {
			t.Fatalf("%s mutation still matched; dropping %s from equality would keep the original row", mutation.field, mutation.field)
		}
	}
	mixers, err := catalogvaillant.B523ConfigMixerEntries()
	if err != nil || len(mixers) == 0 {
		t.Fatal(err)
	}
	mixerKey := mixers[0].B523ConfigMixerKey
	if _, found, err := moduleMixerField(mixerKey); err != nil || !found {
		t.Fatalf("unmutated mixer found=%v err=%v", found, err)
	}
	for _, mutation := range []struct {
		field  string
		mutate func(*catalogvaillant.B523ConfigMixerKey)
	}{
		{"family", func(key *catalogvaillant.B523ConfigMixerKey) { key.B523ConfigMixerFamily = "absent-family" }},
		{"variant", func(key *catalogvaillant.B523ConfigMixerKey) { key.B523ConfigMixerVariant = "absent-variant" }},
		{"primary", func(key *catalogvaillant.B523ConfigMixerKey) { key.B523ConfigMixerPrimary++ }},
		{"secondary", func(key *catalogvaillant.B523ConfigMixerKey) {
			key.B523ConfigMixerSecondary = mutateOptionalOctet(key.B523ConfigMixerSecondary)
		}},
		{"operation", func(key *catalogvaillant.B523ConfigMixerKey) {
			key.B523ConfigMixerOperation = mutateOptionalOctet(key.B523ConfigMixerOperation)
		}},
		{"direction", func(key *catalogvaillant.B523ConfigMixerKey) { key.B523ConfigMixerDirection = "sideways" }},
		{"field", func(key *catalogvaillant.B523ConfigMixerKey) { key.B523ConfigMixerField = "absent-field" }},
		{"group", func(key *catalogvaillant.B523ConfigMixerKey) {
			key.B523ConfigMixerGroup = mutateOptionalOctet(key.B523ConfigMixerGroup)
		}},
		{"requested_instance", func(key *catalogvaillant.B523ConfigMixerKey) {
			key.B523ConfigMixerRequestedInstance = mutateOptionalOctet(key.B523ConfigMixerRequestedInstance)
		}},
		{"register", func(key *catalogvaillant.B523ConfigMixerKey) {
			key.B523ConfigMixerRegister = mutateOptionalWide(key.B523ConfigMixerRegister)
		}},
		{"command_role", func(key *catalogvaillant.B523ConfigMixerKey) {
			key.B523ConfigMixerCommandRole = mutateString(key.B523ConfigMixerCommandRole)
		}},
	} {
		mutated := cloneJSON(t, mixerKey)
		mutation.mutate(&mutated)
		if _, found, err := moduleMixerField(mutated); err != nil || found {
			t.Fatalf("%s mutation still matched; dropping %s from equality would keep the original row", mutation.field, mutation.field)
		}
	}
}

func TestVaillantFunctionModuleSurfaces(t *testing.T) {
	record, ok := moduleFromIdentification(moduleIdentification(0xB5, moduleIDSix, [2]byte{1, 2}, [2]byte{3, 4}), 0x52)
	if !ok {
		t.Fatal("missing module")
	}
	if !reflect.DeepEqual(record.GraphQL(), record.MCP()) {
		t.Fatal("module surfaces diverged")
	}
	status := &output.SystemStatus{}
	status.SetFunctionModule(&record)
	mcp := output.MCPSystemFromGraphQL(status)
	if !reflect.DeepEqual(status.FunctionModule(), mcp.FunctionModule()) {
		t.Fatal("system module diverged")
	}
	poller := &vaillantSemanticPoller{}
	retainFunctionModule(poller, record)
	if retainedFunctionModule(poller).DeviceID != moduleIDSix {
		t.Fatal("retained module lost its identity")
	}
}

func moduleIdentification(manufacturer byte, deviceID string, software, hardware [2]byte) []byte {
	payload := make([]byte, 10)
	payload[0] = manufacturer
	for index := 0; index < 5; index++ {
		payload[1+index] = ' '
	}
	copy(payload[1:6], deviceID)
	payload[6], payload[7] = software[0], software[1]
	payload[8], payload[9] = hardware[0], hardware[1]
	return payload
}

func modulePoller(frames *[]protocol.Frame, reply []byte) *vaillantSemanticPoller {
	return &vaillantSemanticPoller{
		source: adapterSource,
		sendFrameFn: func(_ context.Context, frame protocol.Frame) (*protocol.Frame, error) {
			*frames = append(*frames, frame)
			return &protocol.Frame{Data: append([]byte(nil), reply...)}, nil
		},
	}
}

func mutateByte(current *byte) *byte {
	if current == nil {
		value := byte(0)
		return &value
	}
	value := *current + 1
	return &value
}
