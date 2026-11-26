package owner

import (
	"context"
	"fmt"
	"time"

	"github.com/Project-Helianthus/helianthus-ebusgo/engine"
	"github.com/Project-Helianthus/helianthus-ebusgo/protocol"
)

func semanticReadB524WatchDescriptorPolicy(key engine.B524WatchKey) (engine.WatchSemanticClass, engine.WatchFreshnessProfile, engine.WatchDirectApplyPolicy) {
	switch key.Opcode {
	case vaillantB524OpcodeLocal:
		return semanticReadB524LocalDescriptorPolicy(key.Group, key.RegisterAddress)
	case vaillantB524OpcodeRead:
		return semanticReadB524RemoteDescriptorPolicy(key.Group, key.RegisterAddress)
	default:
		return semanticReadStateFastDescriptorPolicy()
	}
}

func semanticReadB524LocalDescriptorPolicy(group byte, addr uint16) (engine.WatchSemanticClass, engine.WatchFreshnessProfile, engine.WatchDirectApplyPolicy) {
	switch group {
	case localZones.group:
		return semanticReadB524ZoneDescriptorPolicy(addr)
	case localDHW.group:
		return semanticReadB524HotWaterDescriptorPolicy(addr)
	case localCircuits.group:
		return semanticReadB524CircuitDescriptorPolicy(addr)
	case localRegulator.group:
		return semanticReadB524RegulatorDescriptorPolicy(addr)
	case localSolar.group:
		return semanticReadB524SolarDescriptorPolicy(addr)
	case localCylinders.group:
		return semanticReadB524CylinderDescriptorPolicy(addr)
	default:
		return semanticReadStateFastDescriptorPolicy()
	}
}

func semanticReadB524RegulatorDescriptorPolicy(addr uint16) (engine.WatchSemanticClass, engine.WatchFreshnessProfile, engine.WatchDirectApplyPolicy) {
	switch addr {
	case systemOff,
		systemWaterPressure,
		systemFlowTemperature,
		systemOutdoorTemperature,
		systemMaintenanceDue,
		systemHwcCylinderTemperatureTop,
		systemHwcCylinderTemperatureBottom:
		return semanticReadStateFastDescriptorPolicy()
	case systemOutdoorTemperatureDailyAverage,
		energyFuelSumHc,
		energyElectricitySumHc,
		energyElectricitySumHwc,
		energyFuelSumHwc,
		energyFuelSumHcThisMonth,
		energyElectricitySumHcThisMonth,
		energyElectricitySumHwcThisMonth,
		energyFuelSumHwcThisMonth,
		energyFuelSumHcLastMonth,
		energyElectricitySumHcLastMonth,
		energyElectricitySumHwcLastMonth,
		energyFuelSumHwcLastMonth:
		return semanticReadStateSlowDescriptorPolicy()
	case systemAdaptiveHeatingCurve,
		systemAlternativePoint,
		systemHeatingCircuitBivalencePoint,
		systemDhwBivalencePoint,
		systemHcEmergencyTemperature,
		systemHwcMaxFlowTempDesired,
		systemMaxRoomHumidity,
		systemMaintenanceDate,
		systemInstallerName1,
		systemInstallerName2,
		systemInstallerPhone1,
		systemInstallerPhone2,
		systemInstallerMenuCode:
		return semanticReadConfigDescriptorPolicy()
	case systemScheme, systemModuleConfigurationFunctionModule:
		return semanticReadDiscoveryDescriptorPolicy()
	default:
		return semanticReadStateFastDescriptorPolicy()
	}
}

func semanticReadB524ZoneDescriptorPolicy(addr uint16) (engine.WatchSemanticClass, engine.WatchFreshnessProfile, engine.WatchDirectApplyPolicy) {
	switch addr {
	case zoneCurrentTemp, zoneSpecialFunction, zoneValveStatus:
		return semanticReadStateFastDescriptorPolicy()
	case zoneCurrentHumidity, zoneQuickVetoEndTime, zoneQuickVetoEndDate:
		return semanticReadStateSlowDescriptorPolicy()
	case zoneHeatingOpMode,
		zoneTargetTemp,
		zoneFallbackManualTemp,
		zoneRoomTemperatureZoneMappingRawRegister,
		zoneQuickVetoTemp,
		zoneQuickVetoDuration,
		zoneHolidayStartDate,
		zoneHolidayEndDate,
		zoneHolidaySetpoint,
		zoneHolidayEndTime,
		zoneHolidayStartTime:
		return semanticReadConfigDescriptorPolicy()
	case zoneName, zoneNamePrefix, zoneNameSuffix, zoneIndex:
		return semanticReadDiscoveryDescriptorPolicy()
	default:
		return semanticReadStateFastDescriptorPolicy()
	}
}

func semanticReadB524HotWaterDescriptorPolicy(addr uint16) (engine.WatchSemanticClass, engine.WatchFreshnessProfile, engine.WatchDirectApplyPolicy) {
	switch addr {
	case dhwCurrentTemp, dhwSpecialFunction:
		return semanticReadStateFastDescriptorPolicy()
	case dhwOperationMode, dhwTargetTemp, dhwHolidayStartDateRegister, dhwHolidayEndDateRegister:
		return semanticReadConfigDescriptorPolicy()
	default:
		return semanticReadStateFastDescriptorPolicy()
	}
}

func semanticReadB524CircuitDescriptorPolicy(addr uint16) (engine.WatchSemanticClass, engine.WatchFreshnessProfile, engine.WatchDirectApplyPolicy) {
	switch addr {
	case circuitFlowSetpoint,
		circuitFlowTemp,
		circuitCircuitState,
		circuitPumpStatus,
		circuitCalcFlowTemp,
		circuitMixerPosition:
		return semanticReadStateFastDescriptorPolicy()
	case circuitHumidity,
		circuitDewPoint,
		circuitPumpHours,
		circuitPumpStarts:
		return semanticReadStateSlowDescriptorPolicy()
	case circuitTypeRegister,
		circuitCoolingEnabled,
		circuitHeatingCurve,
		circuitFlowTempMax,
		circuitFlowTempMin,
		circuitSummerLimit,
		circuitRoomTempControl,
		circuitSpecialFunctionMode:
		return semanticReadConfigDescriptorPolicy()
	default:
		return semanticReadStateFastDescriptorPolicy()
	}
}

func semanticReadB524SolarDescriptorPolicy(addr uint16) (engine.WatchSemanticClass, engine.WatchFreshnessProfile, engine.WatchDirectApplyPolicy) {
	switch addr {
	case solarCollectorTemp, solarReturnTemp, solarPumpActive, solarCurrentYield:
		return semanticReadStateFastDescriptorPolicy()
	case solarPumpHours:
		return semanticReadStateSlowDescriptorPolicy()
	case solarEnabled, solarFunctionMode:
		return semanticReadConfigDescriptorPolicy()
	default:
		return semanticReadStateFastDescriptorPolicy()
	}
}

func semanticReadB524CylinderDescriptorPolicy(addr uint16) (engine.WatchSemanticClass, engine.WatchFreshnessProfile, engine.WatchDirectApplyPolicy) {
	switch addr {
	case cylinderTemperature:
		return semanticReadStateFastDescriptorPolicy()
	case cylinderMaxSetpoint, cylinderChargeHysteresis, cylinderChargeOffset:
		return semanticReadConfigDescriptorPolicy()
	default:
		return semanticReadStateFastDescriptorPolicy()
	}
}

func semanticReadB524RemoteDescriptorPolicy(group byte, addr uint16) (engine.WatchSemanticClass, engine.WatchFreshnessProfile, engine.WatchDirectApplyPolicy) {
	switch group {
	case remoteRegulators.group, remoteThermostats.group, remoteFunctionalModules.group:
		switch addr {
		case deviceSlotConnected:
			return semanticReadStateSlowDescriptorPolicy()
		case deviceSlotRoomHumidity,
			deviceSlotRoomTemperature,
			deviceSlotPaired,
			deviceSlotReceptionStrength:
			return semanticReadStateSlowDescriptorPolicy()
		case deviceSlotClassAddress,
			deviceSlotFirmware,
			deviceSlotRemoteControlAddress,
			deviceSlotHardwareIdentifier,
			deviceSlotZoneAssignment:
			return semanticReadDiscoveryDescriptorPolicy()
		default:
			return semanticReadStateFastDescriptorPolicy()
		}
	default:
		return semanticReadStateFastDescriptorPolicy()
	}
}

func (p *vaillantSemanticPoller) readB524Value(ctx context.Context, opcode, group, instance byte, addr uint16) ([]byte, bool) {
	return p.readB524ValueWithMaxAge(ctx, opcode, group, instance, addr, -1)
}

// readB524ValueWithMaxAge is readB524Value with a cache-age override: a
// non-negative maxAgeOverride replaces the descriptor TTL (0 never consults the
// scheduler or shadow cache); a negative override keeps the descriptor TTL.
func (p *vaillantSemanticPoller) readB524ValueWithMaxAge(ctx context.Context, opcode, group, instance byte, addr uint16, maxAgeOverride time.Duration) ([]byte, bool) {
	if p == nil || (p.bus == nil && p.sendFrameFn == nil) {
		return nil, false
	}
	if ctx == nil {
		ctx = context.Background()
	}

	p.mu.Lock()
	source := p.source
	target := p.controller
	timeout := p.requestTimeout
	p.mu.Unlock()

	if target == 0 {
		return nil, false
	}
	if timeout <= 0 {
		timeout = 2 * time.Second
	}

	watchKey := engine.NewB524WatchKey(target, opcode, group, instance, addr)
	watchRuntime := p.prepareSemanticReadWatchRuntime(watchKey)
	maxAge := watchRuntime.maxAge
	if maxAgeOverride >= 0 {
		maxAge = maxAgeOverride
	}
	value, stats, err := p.scheduler.GetWatchWithStats(ctx, watchKey, maxAge, func(ctx context.Context) ([]byte, error) {
		var lastErr error
		for attempt := 0; attempt < 3; attempt++ {
			p.readMu.Lock()

			reqCtx, cancel := context.WithTimeout(ctx, timeout)
			request := protocol.Frame{
				FrameType: protocol.FrameTypeInitiatorTarget,
				Source:    source,
				Target:    target,
				Primary:   vaillantExtRegisterPrimary,
				Secondary: vaillantExtRegisterSecondary,
				Data:      buildB524ReadSelector(opcode, group, instance, addr),
			}
			response, err := p.sendSemanticFrame(reqCtx, request)
			cancel()
			p.readMu.Unlock()

			if err != nil {
				lastErr = err
			} else if response == nil {
				lastErr = fmt.Errorf("b524 read returned nil response")
			} else {
				payload, ok := parseB524ReadPayload(response.Data, opcode, group, instance, addr)
				if ok {
					return payload, nil
				}
				lastErr = fmt.Errorf(
					"b524 read failed: opcode=0x%02x group=0x%02x instance=0x%02x addr=0x%04x",
					opcode,
					group,
					instance,
					addr,
				)
			}

			if attempt < 2 && !p.pause(ctx, 75*time.Millisecond) {
				return nil, ctx.Err()
			}
		}
		if lastErr == nil {
			lastErr = fmt.Errorf("b524 read failed")
		}
		return nil, lastErr
	})
	p.emitWatchReadEfficiency(watchRuntime, maxAge, stats)
	if stats.ActiveFetchSucceeded {
		p.pause(ctx, 25*time.Millisecond)
	}
	if err != nil {
		return nil, false
	}
	return value, true
}

// SemanticReadBreakerKey implements the semantic read breaker key step of the Vaillant owner.
func SemanticReadBreakerKey(target, opcode, group, instance byte, addr uint16) string {
	return engine.NewB524WatchKey(target, opcode, group, instance, addr).Canonical()
}
