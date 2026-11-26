package owner

import (
	"fmt"
	"slices"
)

// circuitsViewLocked builds the native view of the active circuits in
// instance order with the managing-device decision taken from the retained
// system snapshot and FM5 mode. The caller holds p.mu.
func (p *vaillantSemanticPoller) circuitsViewLocked() vaillantCircuitsView {
	instances := make([]byte, 0, len(p.circuits))
	for instance := range p.circuits {
		instances = append(instances, instance)
	}
	slices.Sort(instances)
	circuits := make([]*vaillantCircuitSnapshot, 0, len(instances))
	for _, instance := range instances {
		entry := p.circuits[instance]
		if entry == nil || !entry.Active {
			continue
		}
		circuits = append(circuits, cloneCircuitSnapshot(entry))
	}
	return vaillantCircuitsView{
		circuits:                circuits,
		managedByFunctionModule: circuitsManagedByFunctionModule(p.system, p.functionModuleMode),
	}
}

func cloneCircuitSnapshot(snapshot *vaillantCircuitSnapshot) *vaillantCircuitSnapshot {
	if snapshot == nil {
		return nil
	}
	return &vaillantCircuitSnapshot{
		Instance:            snapshot.Instance,
		Active:              snapshot.Active,
		Controller:          snapshot.Controller,
		CircuitTypeRaw:      cloneUint16Ptr(snapshot.CircuitTypeRaw),
		CoolingEnabledRaw:   cloneUint16Ptr(snapshot.CoolingEnabledRaw),
		FlowSetpointC:       cloneFloat64Ptr(snapshot.FlowSetpointC),
		FlowTemperatureC:    cloneFloat64Ptr(snapshot.FlowTemperatureC),
		HeatingCurve:        cloneFloat64Ptr(snapshot.HeatingCurve),
		FlowTempMaxC:        cloneFloat64Ptr(snapshot.FlowTempMaxC),
		FlowTempMinC:        cloneFloat64Ptr(snapshot.FlowTempMinC),
		SummerLimitC:        cloneFloat64Ptr(snapshot.SummerLimitC),
		RoomTempControlRaw:  cloneUint16Ptr(snapshot.RoomTempControlRaw),
		CircuitStateRaw:     cloneUint16Ptr(snapshot.CircuitStateRaw),
		CircuitStateLiveAt:  snapshot.CircuitStateLiveAt,
		SpecialFunctionMode: cloneUint8Ptr(snapshot.SpecialFunctionMode),
		PumpStatusRaw:       cloneUint16Ptr(snapshot.PumpStatusRaw),
		PumpStatusLiveAt:    snapshot.PumpStatusLiveAt,
		// CalcFlowTempC stays unavailable, including snapshots retained before the RR0020 correction.
		MixerPositionPct: cloneFloat64Ptr(snapshot.MixerPositionPct),
		HumidityPct:      cloneFloat64Ptr(snapshot.HumidityPct),
		DewPointC:        cloneFloat64Ptr(snapshot.DewPointC),
		PumpHoursRaw:     cloneUint32Ptr(snapshot.PumpHoursRaw),
		PumpStartsRaw:    cloneUint32Ptr(snapshot.PumpStartsRaw),
	}
}

// DecodeHeatingCircuitTypeToken implements the decode heating circuit type token step of the Vaillant owner.
func DecodeHeatingCircuitTypeToken(raw *uint16) string {
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

// DecodeRoomTempControlToken implements the decode room temp control token step of the Vaillant owner.
func DecodeRoomTempControlToken(raw *uint16) string {
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

// CircuitStateNames is the circuit state names value of the Vaillant owner.
var CircuitStateNames = map[uint16]string{
	0: CircuitStateStandby,
	1: CircuitStateHeating,
	2: CircuitStateCooling,
	3: CircuitStateDHW,
}

// DecodeCircuitStateToken implements the decode circuit state token step of the Vaillant owner.
func DecodeCircuitStateToken(raw *uint16) string {
	if raw == nil {
		return ""
	}
	if name, ok := CircuitStateNames[*raw]; ok {
		return name
	}
	return fmt.Sprintf("unknown_%d", *raw)
}

const (
	circuitManagingDeviceFunctionModuleControllerID      = "VR_71"
	circuitManagingDeviceFunctionModuleControllerAddress = 0x26
)

// circuitsManagedByFunctionModule reports whether the system scheme, the
// module configuration and the interpreted FM5 mode place the heating
// circuits under the function module.
func circuitsManagedByFunctionModule(system *vaillantSystemSnapshot, functionModuleMode vaillantFunctionModuleMode) bool {
	return system != nil &&
		system.SystemScheme != nil && *system.SystemScheme == 1 &&
		system.ModuleConfigurationFunctionModuleController != nil && *system.ModuleConfigurationFunctionModuleController == 2 &&
		functionModuleMode == vaillantFunctionModuleModeInterpreted
}

// DecodeValvePositionPct implements the decode valve position pct step of the Vaillant owner.
func DecodeValvePositionPct(raw *uint16) *float64 {
	if raw == nil {
		return nil
	}
	pct := float64(*raw) / 655.35
	return &pct
}

// DecodeCircuitType implements the decode circuit type step of the Vaillant owner.
func DecodeCircuitType(raw *uint16) string {
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

// DecodeAssociatedCircuit implements the decode associated circuit step of the Vaillant owner.
func DecodeAssociatedCircuit(raw *uint16) *int {
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
