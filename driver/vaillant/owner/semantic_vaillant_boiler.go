package owner

import (
	"context"
	"log"
	"strconv"
	"strings"
	"time"
)

// refreshBoilerAddress applies the registry's boiler-class heat source without
// a root discovery pass; it reads the registry only and sends nothing. While
// the fallback discovery is suspended it keeps a heat source identified after
// ABSENT on the state and configuration ticks.
func (p *vaillantSemanticPoller) refreshBoilerAddress(ctx context.Context) {
	boilerAddress := p.findBoilerAddress()
	p.mu.Lock()
	prevBoilerAddress := p.boilerAddress
	p.boilerAddress = boilerAddress
	p.mu.Unlock()
	p.refreshHeatSourcePathReason()
	if boilerAddress != prevBoilerAddress && boilerAddress != 0 && boilerAddress != 0x08 {
		log.Printf("semantic_boiler_discovery address=0x%02x source=nonstandard", boilerAddress)
	}
	if boilerAddress != prevBoilerAddress && boilerAddress != 0 {
		p.enqueueBoilerStatusPriming(ctx)
	}
}

type vaillantBoilerSnapshot struct {
	FlowTemperatureC         *float64
	ReturnTemperatureC       *float64
	CentralHeatingPumpActive *bool
	WaterPressureBar         *float64
	ExternalPumpActive       *bool
	CirculationPumpActive    *bool
	GasValveActive           *bool
	FlameActive              *bool
	DiverterValvePositionPct *float64
	FanSpeedRpm              *int
	TargetFanSpeedRpm        *int
	IonisationVoltageUa      *float64
	DhwWaterFlowLpm          *float64
	DhwDemandActive          *bool
	HeatingSwitchActive      *bool
	StorageLoadPumpPct       *float64
	ModulationPct            *float64
	PrimaryCircuitFlowLpm    *float64
	FlowTempDesiredC         *float64
	DhwTempDesiredC          *float64
	StateNumber              *int
	DhwTemperatureC          *float64
	DhwTargetTemperatureC    *float64
	DhwOperatingMode         *int
	FlowsetHcMaxC            *float64
	FlowsetHwcMaxC           *float64
	PartloadHcKW             *float64
	PartloadHwcKW            *float64
	HeatingStatusRaw         *int
	DhwStatusRaw             *int
	CentralHeatingHours      *float64
	DhwHours                 *float64
	CentralHeatingStarts     *int
	DhwStarts                *int
	PumpHours                *float64
	FanHours                 *float64
	DeactivationsIFC         *int
	DeactivationsTemplimiter *int
	InstallerMenuCode        *int
	PhoneNumber              *string
	HoursTillService         *int
}

type boilerStatusField uint8

// Circuit registers (GG=0x02) are not boiler fields: circuit 0's pump status
// (RR=0x001B) and circuit state (RR=0x001E) are published only on the circuit
// surface, so no boiler status field is defined for them.
const (
	boilerStatusFieldFlowTemperature boilerStatusField = iota
)

type boilerStatusRegisterDecoder uint8

const (
	boilerStatusRegisterDecoderFloat32 boilerStatusRegisterDecoder = iota
)

type boilerStatusRegisterDefinition struct {
	field    boilerStatusField
	decoder  boilerStatusRegisterDecoder
	opcode   byte
	group    byte
	instance byte
	addr     uint16
}

// findBoilerAddress is retained for potential future use (e.g. broadcast sniffing).
// Currently unused — boiler data is read via B524 registers on the controller.
func boilerStatusRegisterDefinitionsForTier(tier boilerStatusTier) []boilerStatusRegisterDefinition {
	switch tier {
	case boilerStatusTierFast:
		return []boilerStatusRegisterDefinition{
			{
				field:    boilerStatusFieldFlowTemperature,
				decoder:  boilerStatusRegisterDecoderFloat32,
				opcode:   localRegulator.opcode,
				group:    localRegulator.group,
				instance: regulatorInstance,
				addr:     systemFlowTemperature,
			},
		}
	case boilerStatusTierMedium, boilerStatusTierSlow:
		// Reduced profile: scaffolding only. Additional mappings are added once
		// authoritative B524 semantics are promoted for boiler_status.
		return nil
	default:
		return nil
	}
}

// refreshBoilerStatus reads boiler operational data via B524 registers on the controller.
// The BAI boiler does not respond to direct B504 reads from third-party sources — it only
// accepts requests from its paired controller. Instead, the controller (VRC/BASV2) mirrors
// boiler data in its own B524 register space, which we can read reliably.
func (p *vaillantSemanticPoller) refreshBoilerStatus(ctx context.Context) {
	p.refreshBoilerStatusTier(ctx, boilerStatusTierFast)
}

func (p *vaillantSemanticPoller) refreshBoilerStatusTier(ctx context.Context, tier boilerStatusTier) {
	if p == nil {
		return
	}

	p.mu.Lock()
	controller := p.controller
	boilerAddress := p.boilerAddress
	p.mu.Unlock()
	if controller == 0 && boilerAddress == 0 {
		if p.mergePublishedBoilerFlow() {
			p.publishBoilerStatus(semanticSnapshotSourceLive)
		}
		return
	}

	snapshot := &vaillantBoilerSnapshot{}
	updated := false
	if controller != 0 {
		updated = p.refreshBoilerStatusB524(ctx, tier, snapshot) || updated
	}
	if boilerAddress != 0 {
		updated = p.refreshBoilerStatusB509(ctx, boilerAddress, tier, snapshot) || updated
	}
	dhwUpdated := false
	if tier == boilerStatusTierFast {
		dhwUpdated = p.mergeBoilerDHWFieldsFromSnapshot(snapshot)
	}

	// Preserve last-known values on transient/partial failures.
	if !updated {
		if dhwUpdated {
			p.mu.Lock()
			p.boiler = mergeBoilerSnapshotNonDestructive(p.boiler, snapshot)
			p.mu.Unlock()
		}
		if p.mergePublishedBoilerFlow() {
			p.publishBoilerStatus(semanticSnapshotSourceLive)
		}
		return
	}

	p.mu.Lock()
	p.boiler = mergeBoilerSnapshotNonDestructive(p.boiler, snapshot)
	p.mu.Unlock()

	p.publishBoilerStatus(semanticSnapshotSourceLive)
}

func (p *vaillantSemanticPoller) refreshBoilerStatusB524(ctx context.Context, tier boilerStatusTier, snapshot *vaillantBoilerSnapshot) bool {
	if p == nil || snapshot == nil {
		return false
	}

	if tier == boilerStatusTierFast {
		return p.mergeBoilerMirrorFieldsFromSnapshots(snapshot)
	}

	updated := false
	registers := boilerStatusRegisterDefinitionsForTier(tier)
	for _, register := range registers {
		switch register.decoder {
		case boilerStatusRegisterDecoderFloat32:
			value, ok := p.readB524Float32LE(ctx, register.opcode, register.group, register.instance, register.addr)
			if !ok {
				continue
			}
			updated = true
			switch register.field {
			case boilerStatusFieldFlowTemperature:
				// Controller flow is not qualified boiler flow.
				_ = value
			}
		}
	}

	return updated
}

func (p *vaillantSemanticPoller) mergeBoilerMirrorFieldsFromSnapshots(snapshot *vaillantBoilerSnapshot) bool {
	if p == nil || snapshot == nil {
		return false
	}

	p.mu.Lock()
	controller := p.controller
	system := cloneSystemSnapshot(p.system)
	p.mu.Unlock()
	if controller == 0 {
		return false
	}

	now := p.now()
	maxAge := p.boilerMirrorMaxAge()
	// Fresh controller flow triggers publication without supplying boiler flow.
	updated := system != nil && system.Controller == controller && system.SystemFlowTemperature != nil && semanticSnapshotObservedWithin(system.SystemFlowTemperatureLiveAt, now, maxAge)
	var systemFlow *float64
	if system != nil {
		systemFlow = system.SystemFlowTemperature
	}
	snapshot.FlowTemperatureC = boilerFlowWithoutSystem(systemFlow, snapshot.FlowTemperatureC)
	// Circuit 0 is not a boiler source. Its pump status and circuit state stay
	// on the circuit surface; the boiler pump and heating-status diagnostic
	// remain absent here, so partial merges preserve independent boiler evidence.
	if updated {
		p.mergeBoilerDHWFieldsFromSnapshot(snapshot)
	}
	return updated
}

func (p *vaillantSemanticPoller) boilerMirrorMaxAge() time.Duration {
	interval := p.configInterval
	if interval <= 0 {
		interval = defaultSemanticConfigInterval
	}
	fastInterval := p.boilerFastInterval
	if fastInterval <= 0 {
		fastInterval = 30 * time.Second
	}
	return interval + fastInterval
}

func semanticSnapshotObservedWithin(observedAt, now time.Time, maxAge time.Duration) bool {
	if observedAt.IsZero() {
		return false
	}
	if maxAge <= 0 {
		return true
	}
	if now.Before(observedAt) {
		return true
	}
	return now.Sub(observedAt) <= maxAge
}

func (p *vaillantSemanticPoller) mergeBoilerDHWFieldsFromSnapshot(snapshot *vaillantBoilerSnapshot) bool {
	if p == nil || snapshot == nil {
		return false
	}

	p.mu.Lock()
	dhw := p.dhw
	if dhw == nil {
		p.mu.Unlock()
		return false
	}
	current := cloneFloat64Ptr(dhw.CurrentTempC)
	target := cloneFloat64Ptr(dhw.TargetTempC)
	operationMode := strings.TrimSpace(dhw.ConfigurationDHWOperationMode)
	p.mu.Unlock()

	updated := false
	if current != nil {
		snapshot.DhwTemperatureC = current
		updated = true
	}
	if target != nil {
		snapshot.DhwTargetTemperatureC = target
		updated = true
	}
	if operationMode != "" {
		if parsed, err := strconv.Atoi(operationMode); err == nil {
			snapshot.DhwOperatingMode = &parsed
			updated = true
		}
	}
	return updated
}
