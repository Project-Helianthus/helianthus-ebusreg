package output

import (
	"fmt"
)

func decodeHeatingCircuitTypeToken(raw *uint16) string {
	if raw == nil {
		return ""
	}
	switch *raw {
	case 1:
		return "heating"
	case 2:
		return "fixed_value"
	case 3:
		return "dhw"
	case 4:
		return "return_increase"
	default:
		return fmt.Sprintf("unknown_%d", *raw)
	}
}

func decodeRoomTempControlToken(raw *uint16) string {
	if raw == nil {
		return ""
	}
	switch *raw {
	case 0:
		return "off"
	case 1:
		return "modulating"
	case 2:
		return "thermostat"
	default:
		return fmt.Sprintf("unknown_%d", *raw)
	}
}

func decodeCircuitStateToken(raw *uint16) string {
	if raw == nil {
		return ""
	}
	if name, ok := circuitStateNames[*raw]; ok {
		return name
	}
	return fmt.Sprintf("unknown_%d", *raw)
}

var circuitStateNames = map[uint16]string{
	0: CircuitStateStandby,
	1: CircuitStateHeating,
	2: CircuitStateCooling,
	3: CircuitStateDHW,
}

const (
	functionModuleManagingDeviceID      = "VR_71"
	functionModuleManagingDeviceAddress = 0x26
)

// CircuitsManagedByFunctionModule reports whether the scheme, module configuration, and mode place circuits under the function module.
func CircuitsManagedByFunctionModule(system *NativeSystemSnapshot, moduleMode NativeFunctionModuleMode) bool {
	return system != nil &&
		system.SystemScheme != nil && *system.SystemScheme == 1 &&
		system.FunctionModuleConfigurationCode != nil && *system.FunctionModuleConfigurationCode == 2 &&
		moduleMode == NativeFunctionModuleModeInterpreted
}

func decodeValvePositionPct(raw *uint16) *float64 {
	if raw == nil {
		return nil
	}
	pct := float64(*raw) / 655.35
	return &pct
}

func decodeCircuitType(raw *uint16) string {
	if raw == nil {
		return ""
	}
	switch *raw {
	case 0:
		return "radiator"
	case 1:
		return "underfloor"
	case 2:
		return "mixed"
	default:
		return fmt.Sprintf("unknown_%d", *raw)
	}
}

func decodeAssociatedCircuit(raw *uint16) *int {
	if raw == nil {
		return nil
	}
	if *raw == 0xFF || *raw == 0xFFFF {
		return nil
	}
	v := int(*raw)
	return &v
}

func decodeRoomTemperatureZoneMapping(raw *uint16) *int {
	if raw == nil {
		return nil
	}
	if *raw == 0xFF || *raw == 0xFFFF {
		return nil
	}
	v := int(*raw)
	return &v
}

// Circuit state labels in the circuit output.
const (
	CircuitStateStandby = "standby"
	CircuitStateHeating = "heating"
	CircuitStateCooling = "cooling"
	CircuitStateDHW     = "dhw"
)

func decodeUint16Bool(raw *uint16) *bool {
	if raw == nil {
		return nil
	}
	parsed := *raw != 0
	return &parsed
}

func decodeUint32Float64(raw *uint32) *float64 {
	if raw == nil {
		return nil
	}
	parsed := float64(*raw)
	return &parsed
}

func decodeUint32Int(raw *uint32) *int {
	if raw == nil {
		return nil
	}
	parsed := int(*raw)
	return &parsed
}

func decodeUint8Int(raw *uint8) *int {
	if raw == nil {
		return nil
	}
	parsed := int(*raw)
	return &parsed
}

func boolPtrEquals(a, b *bool) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return *a == *b
}

func stringPtrEquals(a, b *string) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return *a == *b
}

func intPtrEquals(a, b *int) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return *a == *b
}

func floatPtrEquals(a, b *float64) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return *a == *b
}

// ResolveAssociatedCircuitInstance resolves the associated circuit from the room mapping or zone instance.
func ResolveAssociatedCircuitInstance(roomTemperatureZoneMapping *uint16, zoneInstance byte) byte {
	if roomTemperatureZoneMapping == nil {
		return zoneInstance
	}
	value := *roomTemperatureZoneMapping
	switch value {
	case 0xFF, 0xFFFF, 0:
		return zoneInstance
	default:
		if value >= 1 && value <= 0x20 {
			return byte(value - 1)
		}
		return zoneInstance
	}
}

// HasLiveCylinderEvidence reports whether the native snapshot contains cylinder evidence.
func HasLiveCylinderEvidence(snapshot *NativeCylinderSnapshot) bool {
	return snapshot != nil && snapshot.TemperatureC != nil
}

// DecodeDhwMode maps the native hot water mode to the published label.
func DecodeDhwMode(raw int) string {
	switch raw {
	case 0:
		return "OFF"
	case 1:
		return "ON"
	case 2:
		return "AUTO"
	default:
		return fmt.Sprintf("UNKNOWN(%d)", raw)
	}
}

func uint16ToIntPtr(value *uint16) *int {
	if value == nil {
		return nil
	}
	v := int(*value)
	return &v
}

func uint8ToIntPtr(value *uint8) *int {
	if value == nil {
		return nil
	}
	v := int(*value)
	return &v
}
