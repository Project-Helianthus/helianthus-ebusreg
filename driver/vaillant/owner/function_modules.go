package owner

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"sync"

	"github.com/Project-Helianthus/helianthus-ebusgo/protocol"

	catalogvaillant "github.com/Project-Helianthus/helianthus-ebusreg/catalog/vaillant"
	"github.com/Project-Helianthus/helianthus-ebusreg/driver/vaillant/output"
)

const (
	moduleFamilyTwelve       = "twelve-relay"
	moduleFamilySix          = "six-relay"
	moduleFamilyIdentityOnly = "identification-only"
	moduleNamespace          = "ebusreg"
	moduleAmbiguityUnknown   = "unknown"
	modulePlacementUnknown   = "Unknown"
	moduleLevelUnknown       = "Unknown"
	modulePrimary            = byte(0xB5)
	moduleSecondary          = byte(0x23)
	moduleIDTwelve           = "VR_71"
	moduleIDSix              = "VR_70"
	moduleIDBare             = "VR66"
)

var errModuleLayout = errors.New("function module operation layout is unresolved")

func moduleFromIdentification(payload []byte, address byte) (output.FunctionModuleRecord, bool) {
	deviceID, ok := moduleDeviceID(payload)
	if !ok {
		return output.FunctionModuleRecord{}, false
	}
	family := moduleFamily(deviceID)
	if family == "" {
		return output.FunctionModuleRecord{}, false
	}
	record := output.FunctionModuleRecord{
		DeviceID:          deviceID,
		Address:           address,
		Manufacturer:      payload[0],
		Software:          append([]byte(nil), payload[6:8]...),
		Hardware:          append([]byte(nil), payload[8:10]...),
		RawIdentification: append([]byte(nil), payload...),
		Family:            family,
		Operations:        moduleOperations(deviceID),
		Ambiguity:         moduleAmbiguityUnknown,
	}
	if family == moduleFamilyIdentityOnly {
		record.Sensors = []output.ModuleReading{{
			Name:       "identification",
			Raw:        append([]byte(nil), payload...),
			Level:      moduleLevelUnknown,
			Placement:  modulePlacementUnknown,
			Unresolved: "unresolved",
		}}
	}
	return record, true
}

func moduleDeviceID(payload []byte) (string, bool) {
	if len(payload) < 10 {
		return "", false
	}
	deviceID := strings.TrimRight(string(payload[1:6]), " \x00")
	if deviceID == "" {
		return "", false
	}
	return deviceID, true
}

func moduleFamily(deviceID string) string {
	switch deviceID {
	case moduleIDTwelve:
		return moduleFamilyTwelve
	case moduleIDSix:
		return moduleFamilySix
	case moduleIDBare:
		return moduleFamilyIdentityOnly
	default:
		return ""
	}
}

func moduleOperations(deviceID string) []output.ModuleOperation {
	switch deviceID {
	case moduleIDTwelve:
		return append([]output.ModuleOperation{moduleActorOperation(0x05, modulePWMAlternatives(0x05))}, moduleMixerOperations([]byte{0, 1, 2})...)
	case moduleIDSix:
		return append([]output.ModuleOperation{moduleActorOperation(0x01, modulePWMAlternatives(0x01))}, moduleMixerOperations([]byte{0, 1})...)
	default:
		return nil
	}
}

func moduleActorOperation(operation byte, alternatives []string) output.ModuleOperation {
	return output.ModuleOperation{
		Namespace:    moduleNamespace,
		Name:         "SetActorState",
		Operation:    operation,
		Unresolved:   []string{"pwm numeric and behavior alternatives"},
		Alternatives: alternatives,
	}
}

func moduleMixerOperations(circuits []byte) []output.ModuleOperation {
	operations := make([]output.ModuleOperation, 0, len(circuits))
	for _, circuit := range circuits {
		value := circuit
		operations = append(operations, output.ModuleOperation{
			Namespace:  moduleNamespace,
			Name:       moduleMixerName(value),
			Operation:  0x02,
			Circuit:    &value,
			Unresolved: []string{"mode and setpoint layout"},
		})
	}
	return operations
}

func moduleMixerName(circuit byte) string {
	switch circuit {
	case 0:
		return "Mc1Operation"
	case 1:
		return "Mc2Operation"
	case 2:
		return "Mc3Operation"
	default:
		return ""
	}
}

func modulePWMAlternatives(operation byte) []string {
	secondary := byte(0x23)
	op := operation
	entry, found, err := catalogvaillant.B523CandidateLookup(catalogvaillant.B523CandidateKey{
		B523CandidateFamily:    "vaillant",
		B523CandidatePrimary:   181,
		B523CandidateSecondary: &secondary,
		B523CandidateOperation: &op,
		B523CandidateDirection: "request",
		B523CandidateField:     "pwm_1",
	})
	if err != nil || !found {
		return nil
	}
	names := make([]string, 0, len(entry.B523CandidateAlternatives))
	for _, alternative := range entry.B523CandidateAlternatives {
		names = append(names, alternative.B523CandidateName)
	}
	return names
}

func moduleReadOperations(deviceID string) []byte {
	switch deviceID {
	case moduleIDTwelve:
		return []byte{0x06, 0x07}
	case moduleIDSix:
		return []byte{0x03}
	default:
		return nil
	}
}

func moduleSensorName(operation byte) string {
	switch operation {
	case 0x06:
		return "SensorData1"
	case 0x07:
		return "SensorData2"
	case 0x03:
		return "SensorData"
	default:
		return ""
	}
}

func readModuleSensors(ctx context.Context, poller *vaillantSemanticPoller, target byte, deviceID string) ([]output.ModuleReading, error) {
	operations := moduleReadOperations(deviceID)
	if len(operations) == 0 {
		return nil, nil
	}
	readings := make([]output.ModuleReading, 0, len(operations))
	for _, operation := range operations {
		response, err := sendModule(ctx, poller, target, []byte{operation})
		reading := output.ModuleReading{
			Name:       moduleSensorName(operation),
			Operation:  operation,
			Level:      moduleLevelUnknown,
			Placement:  modulePlacementUnknown,
			Unresolved: "unresolved",
		}
		if err != nil {
			readings = append(readings, reading)
			return readings, err
		}
		reading.Raw = append([]byte(nil), response...)
		readings = append(readings, reading)
	}
	return readings, nil
}

func callModuleOperation(ctx context.Context, poller *vaillantSemanticPoller, target byte, deviceID, name string, circuit *byte, raw []byte) error {
	operation, ok := moduleOperation(deviceID, name, circuit)
	if !ok {
		return errModuleLayout
	}
	payload, err := modulePayload(deviceID, operation, raw)
	if err != nil {
		return err
	}
	_, err = sendModule(ctx, poller, target, payload)
	return err
}

func moduleOperation(deviceID, name string, circuit *byte) (output.ModuleOperation, bool) {
	for _, operation := range moduleOperations(deviceID) {
		if operation.Name != name {
			continue
		}
		if operation.Circuit == nil && circuit == nil {
			return operation, true
		}
		if operation.Circuit != nil && circuit != nil && *operation.Circuit == *circuit {
			return operation, true
		}
	}
	return output.ModuleOperation{}, false
}

func modulePayload(deviceID string, operation output.ModuleOperation, raw []byte) ([]byte, error) {
	if operation.Name == "SetActorState" {
		width := 7
		if deviceID == moduleIDTwelve {
			width = 14
		}
		if len(raw) != width {
			return nil, errModuleLayout
		}
		payload := make([]byte, 0, 1+len(raw))
		payload = append(payload, operation.Operation)
		payload = append(payload, raw...)
		return payload, nil
	}
	if operation.Circuit == nil || len(raw) != 0 {
		return nil, errModuleLayout
	}
	return []byte{operation.Operation, *operation.Circuit}, nil
}

func sendModule(ctx context.Context, poller *vaillantSemanticPoller, target byte, data []byte) ([]byte, error) {
	if poller == nil {
		return nil, errors.New("function module poller unavailable")
	}
	if ctx == nil {
		ctx = context.Background()
	}
	source, _ := poller.sourceAndCompanion()
	response, err := poller.sendSemanticFrame(ctx, protocol.Frame{
		FrameType: protocol.FrameTypeInitiatorTarget,
		Source:    source,
		Target:    target,
		Primary:   modulePrimary,
		Secondary: moduleSecondary,
		Data:      append([]byte(nil), data...),
	})
	if err != nil {
		return nil, err
	}
	if response == nil {
		return nil, nil
	}
	return append([]byte(nil), response.Data...), nil
}

func interpretModuleSensors(input json.RawMessage) (json.RawMessage, error) {
	return catalogvaillant.B523SensorViewInterpret(input)
}

func interpretModuleActuators(input json.RawMessage) (json.RawMessage, error) {
	return catalogvaillant.B523ActuatorFieldEvaluate(input)
}

func moduleSensorField(key catalogvaillant.B523SensorKey) (catalogvaillant.B523SensorField, bool, error) {
	return catalogvaillant.B523LookupSensorField(key)
}

func moduleCandidateField(key catalogvaillant.B523CandidateKey) (catalogvaillant.B523CandidateEntry, bool, error) {
	return catalogvaillant.B523CandidateLookup(key)
}

func moduleMixerField(key catalogvaillant.B523ConfigMixerKey) (catalogvaillant.B523ConfigMixerEntry, bool, error) {
	return catalogvaillant.B523ConfigMixerLookup(key)
}

func fillFunctionModule(ctx context.Context, poller *vaillantSemanticPoller, payload []byte, address byte) (output.FunctionModuleRecord, error) {
	record, ok := moduleFromIdentification(payload, address)
	if !ok {
		return output.FunctionModuleRecord{}, errModuleLayout
	}
	readings, err := readModuleSensors(ctx, poller, address, record.DeviceID)
	if len(readings) > 0 {
		record.Sensors = append(record.Sensors, readings...)
	}
	retainFunctionModule(poller, record)
	return record, err
}

func retainFunctionModule(p *vaillantSemanticPoller, record output.FunctionModuleRecord) {
	if p == nil {
		return
	}
	moduleMu.Lock()
	defer moduleMu.Unlock()
	cloned := record.GraphQL()
	moduleHeld[p] = &cloned
}

func retainedFunctionModule(p *vaillantSemanticPoller) *output.FunctionModuleRecord {
	if p == nil {
		return nil
	}
	moduleMu.Lock()
	defer moduleMu.Unlock()
	record := moduleHeld[p]
	if record == nil {
		return nil
	}
	cloned := record.GraphQL()
	return &cloned
}

var (
	moduleMu   sync.Mutex
	moduleHeld = map[*vaillantSemanticPoller]*output.FunctionModuleRecord{}
)
