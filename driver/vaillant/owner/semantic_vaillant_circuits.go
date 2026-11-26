package owner

import (
	"context"
	"slices"
	"time"
)

func (p *vaillantSemanticPoller) refreshCircuits(ctx context.Context) {
	if p == nil {
		return
	}

	now := p.now()
	p.mu.Lock()
	controller := p.controller
	instances, fullScan := p.circuitRefreshInstancesLocked(now)
	p.mu.Unlock()
	if controller == 0 {
		return
	}

	discovered, inactive, probeSuccess := p.probeCircuitInstances(ctx, controller, instances)
	if !probeSuccess {
		if !fullScan {
			p.mu.Lock()
			p.lastCircuitFullScanAt = time.Time{}
			p.lastCircuitFullScanComplete = false
			p.mu.Unlock()
		}
		return
	}

	updates := make(map[byte]*vaillantCircuitSnapshot, len(discovered))
	anyRead := false
	for instance, discoveredSnapshot := range discovered {
		snapshot := cloneCircuitSnapshot(discoveredSnapshot)
		if snapshot == nil {
			continue
		}
		anyRead = true
		p.readCircuitDetails(ctx, instance, snapshot, now)
		updates[instance] = snapshot
	}

	p.mu.Lock()
	if p.circuits == nil {
		p.circuits = make(map[byte]*vaillantCircuitSnapshot)
	}
	for instance := range inactive {
		delete(p.circuits, instance)
	}
	for instance, incoming := range updates {
		p.circuits[instance] = mergeCircuitSnapshotNonDestructive(p.circuits[instance], incoming)
	}
	if fullScan {
		p.lastCircuitFullScanAt = now
		p.lastCircuitFullScanComplete = len(discovered)+len(inactive) == len(instances)
	}
	p.mu.Unlock()

	source := semanticSnapshotSourceCache
	if anyRead || len(inactive) > 0 {
		source = semanticSnapshotSourceLive
	}
	p.publishCircuits(source)
}

// probeCircuitInstances reads the circuit type of each instance in order. It
// returns the active circuits, the instances that reported no circuit, and
// whether any type read succeeded.
func (p *vaillantSemanticPoller) probeCircuitInstances(ctx context.Context, controller byte, instances []byte) (map[byte]*vaillantCircuitSnapshot, map[byte]struct{}, bool) {
	discovered := make(map[byte]*vaillantCircuitSnapshot, 4)
	inactive := make(map[byte]struct{})
	probeSuccess := false

	for _, instance := range instances {
		circuitTypeRaw, ok := p.readB524Uint16(ctx, localCircuits.opcode, localCircuits.group, instance, circuitTypeRegister)
		if !ok || circuitTypeRaw == nil {
			continue
		}
		probeSuccess = true
		switch *circuitTypeRaw {
		case 0x0000, 0x00FF, 0xFFFF:
			if snapshot, known := p.readDhwPseudoCircuitEvidence(ctx, instance); snapshot != nil {
				snapshot.Controller = controller
				discovered[instance] = snapshot
				continue
			} else if !known && instance == dhwPseudoCircuitInstance {
				continue
			}
			inactive[instance] = struct{}{}
		default:
			discovered[instance] = &vaillantCircuitSnapshot{
				Instance:       instance,
				Active:         true,
				Controller:     controller,
				CircuitTypeRaw: cloneUint16Ptr(circuitTypeRaw),
			}
		}
	}
	return discovered, inactive, probeSuccess
}

// readCircuitDetails reads the detail registers of one active circuit in a
// fixed order. A failed read keeps the field the snapshot already holds.
func (p *vaillantSemanticPoller) readCircuitDetails(ctx context.Context, instance byte, snapshot *vaillantCircuitSnapshot, now time.Time) {
	if raw := p.readCircuitUint16(ctx, instance, circuitCoolingEnabled); raw != nil {
		snapshot.CoolingEnabledRaw = raw
	}
	if value := p.readCircuitFloat(ctx, instance, circuitFlowSetpoint); value != nil {
		snapshot.FlowSetpointC = value
	}
	if value := p.readCircuitFloat(ctx, instance, circuitFlowTemp); value != nil {
		snapshot.FlowTemperatureC = value
	}
	if value := p.readCircuitFloat(ctx, instance, circuitHeatingCurve); value != nil {
		snapshot.HeatingCurve = value
	}
	if value := p.readCircuitFloat(ctx, instance, circuitFlowTempMax); value != nil {
		snapshot.FlowTempMaxC = value
	}
	if value := p.readCircuitFloat(ctx, instance, circuitFlowTempMin); value != nil {
		snapshot.FlowTempMinC = value
	}
	if value := p.readCircuitFloat(ctx, instance, circuitSummerLimit); value != nil {
		snapshot.SummerLimitC = value
	}
	if raw := p.readCircuitUint16(ctx, instance, circuitRoomTempControl); raw != nil {
		snapshot.RoomTempControlRaw = raw
	}
	if raw := p.readCircuitUint16(ctx, instance, circuitCircuitState); raw != nil {
		snapshot.CircuitStateRaw = raw
		snapshot.CircuitStateLiveAt = now
	}
	if value := p.readB524Uint8(ctx, localCircuits.opcode, localCircuits.group, instance, circuitSpecialFunctionMode); value != nil {
		snapshot.SpecialFunctionMode = cloneUint8Ptr(value)
	}
	if raw := p.readCircuitUint16(ctx, instance, circuitPumpStatus); raw != nil {
		snapshot.PumpStatusRaw = raw
		snapshot.PumpStatusLiveAt = now
	}
	if value := p.readCircuitFloat(ctx, instance, circuitMixerPosition); value != nil {
		snapshot.MixerPositionPct = value
	}
	if value := p.readCircuitFloat(ctx, instance, circuitHumidity); value != nil {
		snapshot.HumidityPct = value
	}
	if value := p.readCircuitFloat(ctx, instance, circuitDewPoint); value != nil {
		snapshot.DewPointC = value
	}
	if value := p.readCircuitUint32(ctx, instance, circuitPumpHours); value != nil {
		snapshot.PumpHoursRaw = value
	}
	if value := p.readCircuitUint32(ctx, instance, circuitPumpStarts); value != nil {
		snapshot.PumpStartsRaw = value
	}
}

// readCircuitFloat reads one float register of a circuit, or nil when the
// read fails.
func (p *vaillantSemanticPoller) readCircuitFloat(ctx context.Context, instance byte, addr uint16) *float64 {
	value, ok := p.readB524Float32LE(ctx, localCircuits.opcode, localCircuits.group, instance, addr)
	if !ok {
		return nil
	}
	return &value
}

// readCircuitUint16 reads one 16-bit register of a circuit as a detached
// value, or nil when the read fails or carries no value.
func (p *vaillantSemanticPoller) readCircuitUint16(ctx context.Context, instance byte, addr uint16) *uint16 {
	raw, ok := p.readB524Uint16(ctx, localCircuits.opcode, localCircuits.group, instance, addr)
	if !ok || raw == nil {
		return nil
	}
	return cloneUint16Ptr(raw)
}

// readCircuitUint32 reads one 32-bit counter register of a circuit, or nil
// when the read fails.
func (p *vaillantSemanticPoller) readCircuitUint32(ctx context.Context, instance byte, addr uint16) *uint32 {
	value, ok := p.readB524Uint32LE(ctx, localCircuits.opcode, localCircuits.group, instance, addr)
	if !ok {
		return nil
	}
	return &value
}

func (p *vaillantSemanticPoller) readDhwPseudoCircuitEvidence(ctx context.Context, instance byte) (*vaillantCircuitSnapshot, bool) {
	if p == nil || instance != dhwPseudoCircuitInstance {
		return nil, false
	}
	snapshot := newDhwPseudoCircuitSnapshot(instance)
	flow, flowKnown := p.readDhwPseudoCircuitTemperatureEvidence(ctx, instance, circuitFlowTemp)
	if isDhwPseudoCircuitTemperatureEvidence(flow) {
		snapshot.FlowTemperatureC = flow
	}
	if snapshot.FlowTemperatureC == nil {
		return nil, flowKnown
	}
	return snapshot, true
}

func (p *vaillantSemanticPoller) readDhwPseudoCircuitTemperatureEvidence(ctx context.Context, instance byte, addr uint16) (*float64, bool) {
	raw, ok := p.readB524Value(ctx, localCircuits.opcode, localCircuits.group, instance, addr)
	if !ok {
		return nil, false
	}
	return decodeB524Float32FromRaw(raw), true
}

func newDhwPseudoCircuitSnapshot(instance byte) *vaillantCircuitSnapshot {
	circuitType := dhwPseudoCircuitTypeRaw
	return &vaillantCircuitSnapshot{
		Instance:       instance,
		Active:         true,
		CircuitTypeRaw: &circuitType,
	}
}

func isDhwPseudoCircuitTemperatureEvidence(value *float64) bool {
	if value == nil {
		return false
	}
	return *value >= 5 && *value <= 95
}

func allCircuitRefreshInstances() []byte {
	instances := make([]byte, 0, 0x0B)
	for instance := byte(0x00); instance <= 0x0A; instance++ {
		instances = append(instances, instance)
	}
	return instances
}

func (p *vaillantSemanticPoller) circuitRefreshInstancesLocked(now time.Time) ([]byte, bool) {
	if p == nil {
		return nil, false
	}
	interval := p.circuitFullScanInterval
	if interval <= 0 {
		interval = semanticCircuitFullScanInterval
	}
	if len(p.circuits) == 0 || p.lastCircuitFullScanAt.IsZero() || now.Sub(p.lastCircuitFullScanAt) >= interval {
		return allCircuitRefreshInstances(), true
	}
	if !p.lastCircuitFullScanComplete && now.Sub(p.lastCircuitFullScanAt) >= semanticCircuitPartialScanInterval {
		return allCircuitRefreshInstances(), true
	}

	instances := make([]byte, 0, len(p.circuits))
	for instance, snapshot := range p.circuits {
		if snapshot == nil || !snapshot.Active {
			continue
		}
		instances = append(instances, instance)
	}
	if len(instances) == 0 {
		return allCircuitRefreshInstances(), true
	}
	slices.Sort(instances)
	return instances, false
}

func mergeCircuitSnapshotNonDestructive(existing, incoming *vaillantCircuitSnapshot) *vaillantCircuitSnapshot {
	merged := cloneCircuitSnapshot(existing)
	if merged == nil {
		merged = &vaillantCircuitSnapshot{}
	}
	if incoming == nil {
		return merged
	}

	merged.Instance = incoming.Instance
	merged.Active = merged.Active || incoming.Active
	controllerChanged := incoming.Controller != 0 && incoming.Controller != merged.Controller
	if controllerChanged && incoming.CircuitStateRaw == nil {
		merged.CircuitStateRaw = nil
		merged.CircuitStateLiveAt = time.Time{}
	}
	if controllerChanged && incoming.PumpStatusRaw == nil {
		merged.PumpStatusRaw = nil
		merged.PumpStatusLiveAt = time.Time{}
	}
	if incoming.Controller != 0 {
		merged.Controller = incoming.Controller
	}
	if incoming.CircuitTypeRaw != nil {
		merged.CircuitTypeRaw = cloneUint16Ptr(incoming.CircuitTypeRaw)
	}
	if incoming.CoolingEnabledRaw != nil {
		merged.CoolingEnabledRaw = cloneUint16Ptr(incoming.CoolingEnabledRaw)
	}
	if incoming.FlowSetpointC != nil {
		merged.FlowSetpointC = cloneFloat64Ptr(incoming.FlowSetpointC)
	}
	if incoming.FlowTemperatureC != nil {
		merged.FlowTemperatureC = cloneFloat64Ptr(incoming.FlowTemperatureC)
	}
	if incoming.HeatingCurve != nil {
		merged.HeatingCurve = cloneFloat64Ptr(incoming.HeatingCurve)
	}
	if incoming.FlowTempMaxC != nil {
		merged.FlowTempMaxC = cloneFloat64Ptr(incoming.FlowTempMaxC)
	}
	if incoming.FlowTempMinC != nil {
		merged.FlowTempMinC = cloneFloat64Ptr(incoming.FlowTempMinC)
	}
	if incoming.SummerLimitC != nil {
		merged.SummerLimitC = cloneFloat64Ptr(incoming.SummerLimitC)
	}
	if incoming.RoomTempControlRaw != nil {
		merged.RoomTempControlRaw = cloneUint16Ptr(incoming.RoomTempControlRaw)
	}
	if incoming.CircuitStateRaw != nil {
		merged.CircuitStateRaw = cloneUint16Ptr(incoming.CircuitStateRaw)
		merged.CircuitStateLiveAt = incoming.CircuitStateLiveAt
	}
	if incoming.SpecialFunctionMode != nil {
		merged.SpecialFunctionMode = cloneUint8Ptr(incoming.SpecialFunctionMode)
	}
	if incoming.PumpStatusRaw != nil {
		merged.PumpStatusRaw = cloneUint16Ptr(incoming.PumpStatusRaw)
		merged.PumpStatusLiveAt = incoming.PumpStatusLiveAt
	}
	if incoming.MixerPositionPct != nil {
		merged.MixerPositionPct = cloneFloat64Ptr(incoming.MixerPositionPct)
	}
	if incoming.HumidityPct != nil {
		merged.HumidityPct = cloneFloat64Ptr(incoming.HumidityPct)
	}
	if incoming.DewPointC != nil {
		merged.DewPointC = cloneFloat64Ptr(incoming.DewPointC)
	}
	if incoming.PumpHoursRaw != nil {
		merged.PumpHoursRaw = cloneUint32Ptr(incoming.PumpHoursRaw)
	}
	if incoming.PumpStartsRaw != nil {
		merged.PumpStartsRaw = cloneUint32Ptr(incoming.PumpStartsRaw)
	}
	return merged
}
