package heat

import semreg "github.com/Project-Helianthus/helianthus-semreg/semreg/v1"

// SystemFeeds are the system table rows, including the namespace rows for
// values the pack does not define.
func SystemFeeds() []Feed {
	min, max := singleInstance()
	families := families700720()
	return []Feed{
		{
			ID: "system.outdoor_air", Subject: "system", Group: 0x00, Register: 0x0073,
			InstanceMin: min, InstanceMax: max, NativeName: "outdoor_temperature", NativeUnit: "°C",
			Element: "thermal.measurement.outdoor_air_temperature", Conversion: "celsius", Grade: "Hypothesis",
			Families: families, SourceRank: sourceRankPrimary,
		},
		{
			ID: "system.flow_temperature", Subject: "system", Group: 0x00, Register: 0x004B,
			InstanceMin: min, InstanceMax: max, NativeName: "flow_temperature", NativeUnit: "°C",
			Element: "thermal.measurement.system_flow_temperature", Conversion: "celsius", Grade: "Hypothesis",
			Families: families, SourceRank: sourceRankPrimary,
		},
		{
			ID: "system.outdoor_average", Subject: "system", Group: 0x00, Register: 0x0095,
			InstanceMin: min, InstanceMax: max, NativeName: "outdoor_temperature_daily_average", NativeUnit: "°C",
			Conversion: "namespace", Grade: "Hypothesis", Families: families, SourceRank: sourceRankPrimary,
			Gap: "rolling average has no pack element and is not the instantaneous outdoor temperature",
		},
		{
			ID: "system.pressure", Subject: "system", Group: 0x00, Register: 0x0039,
			InstanceMin: min, InstanceMax: max, NativeName: "water_pressure", NativeUnit: "bar",
			Conversion: "namespace", Grade: "Hypothesis", Families: families, SourceRank: sourceRankPrimary,
			Gap: "system pressure has no admitted pack element on the system subject",
		},
		{
			ID: "system.switch", Subject: "system", Group: 0x00, Register: 0x0007,
			InstanceMin: min, InstanceMax: max, NativeName: "unknown_11",
			Conversion: "namespace", Grade: "Hypothesis", Families: families, SourceRank: sourceRankPrimary,
			Gap: "system switch is not a mode fact",
		},
		{
			ID: "system.maintenance", Subject: "system", Group: 0x00, Register: 0x0096,
			InstanceMin: min, InstanceMax: max, NativeName: "maintenance_active",
			Element: "thermal.status.maintenance", Conversion: "maintenance", Grade: "Hypothesis",
			Families: families, SourceRank: sourceRankPrimary,
		},
		{
			ID: "system.emergency", Subject: "system", Group: 0x00, Register: 0x0047,
			InstanceMin: min, InstanceMax: max, NativeName: "emergency_operation_active",
			Conversion: "namespace", Grade: "Unknown", Families: families, SourceRank: sourceRankPrimary,
			Gap: "emergency operation is not a documented fault symbol",
		},
		{
			ID: "system.environmental_yield", Subject: "system", Group: 0x00, Register: 0x003E,
			InstanceMin: min, InstanceMax: max, NativeName: "environmental_yield_total", NativeUnit: "kWh",
			Conversion: "namespace", Grade: "Unknown", Families: families, SourceRank: sourceRankPrimary,
			Gap: "energy yield counter has no pack element",
		},
		{
			ID: "system.solar_yield", Subject: "system", Group: 0x00, Register: 0x003D,
			InstanceMin: min, InstanceMax: max, NativeName: "solar_yield_total", NativeUnit: "kWh",
			Conversion: "namespace", Grade: "Unknown", Families: families, SourceRank: sourceRankPrimary,
			Gap: "energy yield counter has no pack element",
		},
		{
			ID: "system.heat_recovery_yield", Subject: "system", Group: 0x00, Register: 0x005E,
			InstanceMin: min, InstanceMax: max, NativeName: "heat_recovery_yield_total",
			Conversion: "namespace", Grade: "Unknown", Families: families, SourceRank: sourceRankPrimary,
			Gap: "heat recovery yield has no stated unit and no pack element",
		},
		{
			ID: "system.heating_curve", Subject: "circuit", Group: 0x02, Register: 0x000F,
			InstanceMin: 0, InstanceMax: 10, NativeName: "circuit_heating_curve_factor",
			Conversion: "namespace", Grade: "Hypothesis", Families: families, SourceRank: sourceRankPrimary,
			Gap: "heating curve has no pack element",
		},
		{
			ID: "system.hybrid_limit", Subject: "system", Group: 0x00, Register: 0x000F,
			InstanceMin: min, InstanceMax: max, NativeName: "hybrid_control_mode",
			Conversion: "namespace", Grade: "Unknown", Families: families, SourceRank: sourceRankPrimary,
			Gap: "hybrid limit has no pack element",
		},
		{
			ID: "system.holiday_start", Subject: "system", Group: 0x00, Register: 0x00F2,
			InstanceMin: min, InstanceMax: max, NativeName: "holiday_start_date",
			Conversion: "namespace", Grade: "Unknown", Families: families, SourceRank: sourceRankPrimary,
			Gap: "holiday date has no pack element",
		},
		{
			ID: "system.holiday_start_other", Subject: "system", Group: 0x00, Register: 0x0071,
			InstanceMin: min, InstanceMax: max, NativeName: "holiday_start_date",
			Conversion: "namespace", Grade: "Unknown", Families: []string{"controller-700"}, SourceRank: sourceRankPrimary,
			Gap: "holiday date has no pack element",
		},
		{
			ID: "system.ventilation_mode", Subject: "system", Group: 0x09, Register: 0x0002,
			InstanceMin: 0, InstanceMax: 10, NativeName: "ventilation_operating_mode",
			Conversion: "namespace", Grade: "Hypothesis", Families: []string{"controller-720"}, SourceRank: sourceRankPrimary,
			Gap: "ventilation mode has no pack element",
		},
		{
			ID: "system.ventilation_mode_other", Subject: "system", Group: 0x00, Register: 0x0016,
			InstanceMin: min, InstanceMax: max, NativeName: "ventilation_operating_mode",
			Conversion: "namespace", Grade: "Hypothesis", Families: []string{"controller-700"}, SourceRank: sourceRankPrimary,
			Gap: "ventilation mode has no pack element",
		},
	}
}

// ErrorRead is one current-error observation. Known is false when the reply
// failed or was absent.
type ErrorRead struct {
	Known bool
	Set   bool
}

// DeriveSystem publishes the system operation and action from documented
// circuit-state labels. The origin is derived. An empty set stays unknown.
func DeriveSystem(labels []string, resource string) (Outcome, Outcome) {
	operation := derivedOutcome("thermal.status.operation", resource)
	action := derivedOutcome("thermal.action.state", resource)
	if len(labels) == 0 {
		operation.Status = StatusUnknown
		action.Status = StatusUnknown
		return operation, action
	}
	active := false
	unknown := false
	kinds := map[string]struct{}{}
	for _, label := range labels {
		token, ok := DeriveCircuitAction(label)
		if !ok {
			unknown = true
			continue
		}
		if token != "idle" {
			active = true
		}
		kinds[token] = struct{}{}
	}
	switch {
	case active:
		operation = fedSymbol(operation, "thermal.status.operation", "active")
	case unknown:
		operation.Status = StatusUnknown
	default:
		operation = fedSymbol(operation, "thermal.status.operation", "idle")
	}
	delete(kinds, "idle")
	switch len(kinds) {
	case 0:
		if unknown || active {
			action.Status = StatusUnknown
		} else {
			action = fedSymbol(action, "thermal.action.state", "idle")
		}
	case 1:
		for token := range kinds {
			action = fedSymbol(action, "thermal.action.state", token)
		}
	default:
		action.Status = StatusUnknown
		action.Gap = "circuit states name more than one action"
	}
	return operation, action
}

// DeriveFault publishes thermal.status.fault. It is present when either
// current error is set, and clear only when both reads are known and unset.
func DeriveFault(controller, boiler ErrorRead, resource string) Outcome {
	out := derivedOutcome("thermal.status.fault", resource)
	if (controller.Known && controller.Set) || (boiler.Known && boiler.Set) {
		return fedSymbol(out, "thermal.status.fault", "present")
	}
	if controller.Known && boiler.Known {
		return fedSymbol(out, "thermal.status.fault", "clear")
	}
	out.Status = StatusUnknown
	return out
}

func derivedOutcome(element, resource string) Outcome {
	feed := Feed{
		ID: "derived." + element, Subject: "system", Element: element, Grade: "Hypothesis",
		SourceRank: sourceRankPrimary, Derived: true, NativeName: element, NativeUnit: "none",
	}
	return Outcome{
		Feed: feed, Grade: feed.Grade, Rank: feed.SourceRank, Freshness: "fresh",
		Key:      factKey(feed, resource),
		Evidence: "native=derived unit=none name=" + element + " grade=Hypothesis rank=3 poll=unknown freshness=fresh origin=derived",
	}
}

func symbolValue(element, token string) *semreg.Value {
	return &semreg.Value{Kind: semreg.ValueSymbol, Symbol: &semreg.Symbol{Namespace: semreg.DefinitionID(element), Token: token, Known: true}}
}
