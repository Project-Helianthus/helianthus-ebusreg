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
	case system_off,
		system_water_pressure,
		system_flow_temperature,
		system_outdoor_temperature,
		system_maintenance_due,
		system_hwc_cylinder_temperature_top,
		system_hwc_cylinder_temperature_bottom:
		return semanticReadStateFastDescriptorPolicy()
	case system_outdoor_temperature_avg_24h,
		energy_fuel_sum_hc,
		energy_electricity_sum_hc,
		energy_electricity_sum_hwc,
		energy_fuel_sum_hwc,
		energy_fuel_sum_hc_this_month,
		energy_electricity_sum_hc_this_month,
		energy_electricity_sum_hwc_this_month,
		energy_fuel_sum_hwc_this_month,
		energy_fuel_sum_hc_last_month,
		energy_electricity_sum_hc_last_month,
		energy_electricity_sum_hwc_last_month,
		energy_fuel_sum_hwc_last_month:
		return semanticReadStateSlowDescriptorPolicy()
	case system_adaptive_heating_curve,
		system_alternative_point,
		system_heating_circuit_bivalence_point,
		system_dhw_bivalence_point,
		system_hc_emergency_temperature,
		system_hwc_max_flow_temp_desired,
		system_max_room_humidity,
		system_maintenance_date,
		system_installer_name_1,
		system_installer_name_2,
		system_installer_phone_1,
		system_installer_phone_2,
		system_installer_menu_code:
		return semanticReadConfigDescriptorPolicy()
	case system_scheme, system_module_configuration_function_module:
		return semanticReadDiscoveryDescriptorPolicy()
	default:
		return semanticReadStateFastDescriptorPolicy()
	}
}

func semanticReadB524ZoneDescriptorPolicy(addr uint16) (engine.WatchSemanticClass, engine.WatchFreshnessProfile, engine.WatchDirectApplyPolicy) {
	switch addr {
	case zone_current_temp, zone_special_function, zone_valve_status:
		return semanticReadStateFastDescriptorPolicy()
	case zone_current_humidity, zone_quick_veto_end_time, zone_quick_veto_end_date:
		return semanticReadStateSlowDescriptorPolicy()
	case zone_heating_op_mode,
		zone_target_temp,
		zone_fallback_manual_temp,
		zone_room_temperature_zone_mapping_raw,
		zone_quick_veto_temp,
		zone_quick_veto_duration,
		zone_holiday_start_date,
		zone_holiday_end_date,
		zone_holiday_setpoint,
		zone_holiday_end_time,
		zone_holiday_start_time:
		return semanticReadConfigDescriptorPolicy()
	case zone_name, zone_name_prefix, zone_name_suffix, zone_index:
		return semanticReadDiscoveryDescriptorPolicy()
	default:
		return semanticReadStateFastDescriptorPolicy()
	}
}

func semanticReadB524HotWaterDescriptorPolicy(addr uint16) (engine.WatchSemanticClass, engine.WatchFreshnessProfile, engine.WatchDirectApplyPolicy) {
	switch addr {
	case dhw_current_temp, dhw_special_function:
		return semanticReadStateFastDescriptorPolicy()
	case dhw_operation_mode, dhw_target_temp, dhw_holiday_start_date, dhw_holiday_end_date:
		return semanticReadConfigDescriptorPolicy()
	default:
		return semanticReadStateFastDescriptorPolicy()
	}
}

func semanticReadB524CircuitDescriptorPolicy(addr uint16) (engine.WatchSemanticClass, engine.WatchFreshnessProfile, engine.WatchDirectApplyPolicy) {
	switch addr {
	case circuit_flow_setpoint,
		circuit_flow_temp,
		circuit_circuit_state,
		circuit_pump_status,
		circuit_calc_flow_temp,
		circuit_mixer_position:
		return semanticReadStateFastDescriptorPolicy()
	case circuit_humidity,
		circuit_dew_point,
		circuit_pump_hours,
		circuit_pump_starts:
		return semanticReadStateSlowDescriptorPolicy()
	case circuit_type,
		circuit_cooling_enabled,
		circuit_heating_curve,
		circuit_flow_temp_max,
		circuit_flow_temp_min,
		circuit_summer_limit,
		circuit_room_temp_control,
		circuit_special_function_mode:
		return semanticReadConfigDescriptorPolicy()
	default:
		return semanticReadStateFastDescriptorPolicy()
	}
}

func semanticReadB524SolarDescriptorPolicy(addr uint16) (engine.WatchSemanticClass, engine.WatchFreshnessProfile, engine.WatchDirectApplyPolicy) {
	switch addr {
	case solar_collector_temp, solar_return_temp, solar_pump_active, solar_current_yield:
		return semanticReadStateFastDescriptorPolicy()
	case solar_pump_hours:
		return semanticReadStateSlowDescriptorPolicy()
	case solar_enabled, solar_function_mode:
		return semanticReadConfigDescriptorPolicy()
	default:
		return semanticReadStateFastDescriptorPolicy()
	}
}

func semanticReadB524CylinderDescriptorPolicy(addr uint16) (engine.WatchSemanticClass, engine.WatchFreshnessProfile, engine.WatchDirectApplyPolicy) {
	switch addr {
	case cylinder_temperature:
		return semanticReadStateFastDescriptorPolicy()
	case cylinder_max_setpoint, cylinder_charge_hysteresis, cylinder_charge_offset:
		return semanticReadConfigDescriptorPolicy()
	default:
		return semanticReadStateFastDescriptorPolicy()
	}
}

func semanticReadB524RemoteDescriptorPolicy(group byte, addr uint16) (engine.WatchSemanticClass, engine.WatchFreshnessProfile, engine.WatchDirectApplyPolicy) {
	switch group {
	case remoteRegulators.group, remoteThermostats.group, remoteFunctionalModules.group:
		switch addr {
		case device_slot_connected:
			return semanticReadStateSlowDescriptorPolicy()
		case device_slot_room_humidity,
			device_slot_room_temperature,
			device_slot_paired,
			device_slot_reception_strength:
			return semanticReadStateSlowDescriptorPolicy()
		case device_slot_class_address,
			device_slot_firmware,
			device_slot_remote_control_address,
			device_slot_hardware_identifier,
			device_slot_zone_assignment:
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

func SemanticReadBreakerKey(target, opcode, group, instance byte, addr uint16) string {
	return engine.NewB524WatchKey(target, opcode, group, instance, addr).Canonical()
}
