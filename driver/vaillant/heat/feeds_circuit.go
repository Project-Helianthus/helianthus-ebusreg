package heat

// CircuitFeeds are the circuit table rows. A deactivated circuit withdraws
// the facts of every circuit row.
func CircuitFeeds() []Feed {
	min, max := indexedInstance()
	families := families700720()
	return []Feed{
		{
			ID: "circuit.flow_temperature", Subject: "circuit", Group: 0x02, Register: 0x0008,
			InstanceMin: min, InstanceMax: max, NativeName: "circuit_flow_temperature", NativeUnit: "°C",
			Element: "thermal.measurement.circuit_flow_temperature", Conversion: "celsius", Grade: "Hypothesis",
			Families: families, SourceRank: sourceRankPrimary,
		},
		{
			ID: "circuit.flow_rate", Subject: "circuit", NativeName: "circuit_flow_rate",
			Conversion: "namespace", Grade: "Unknown", NoMatch: true, SourceRank: sourceRankPrimary,
			Gap: "thermal.measurement.flow_rate has no documented circuit register",
		},
		{
			ID: "circuit.cooling_enabled", Subject: "circuit", Group: 0x02, Register: 0x0006,
			InstanceMin: min, InstanceMax: max, NativeName: "circuit_cooling_enabled",
			Conversion: "namespace", Grade: "Hypothesis", Families: families, SourceRank: sourceRankPrimary,
			Gap: "circuit cooling enabled has no SemReg element",
		},
		{
			ID: "circuit.read_capability", Subject: "circuit", Group: 0x02, Register: 0x0008,
			InstanceMin: min, InstanceMax: max, NativeName: "circuit_flow_temperature",
			Element: "thermal.capability.read.circuit", Conversion: "capability", Grade: "Hypothesis",
			Families: families, SourceRank: sourceRankPrimary,
		},
		{
			ID: "circuit.return_setpoint", Subject: "circuit", Group: 0x02, Register: 0x0004,
			InstanceMin: min, InstanceMax: max, NativeName: "circuit_return_temperature_setpoint", NativeUnit: "°C",
			Conversion: "namespace", Grade: "Hypothesis", Families: families, SourceRank: sourceRankPrimary,
			Gap: "return setpoint has no pack element on the circuit subject",
		},
		{
			ID: "circuit.external_demand", Subject: "circuit", Group: 0x02, Register: 0x0019,
			InstanceMin: min, InstanceMax: max, NativeName: "circuit_external_demand_input",
			Conversion: "namespace", Grade: "Hypothesis", Families: families, SourceRank: sourceRankPrimary,
			Gap: "external demand stays native and is not a ventilation ratio",
		},
		{
			ID: "circuit.outdoor_limit", Subject: "circuit", Group: 0x02, Register: 0x001F,
			InstanceMin: min, InstanceMax: max, NativeName: "cooling_circuit_outdoor_temperature_limit_minimum", NativeUnit: "°C",
			Conversion: "namespace", Grade: "Hypothesis", Families: families, SourceRank: sourceRankPrimary,
			Gap: "register is an outdoor limit, not a room setpoint",
		},
		{
			ID: "circuit.flow_setpoint", Subject: "circuit", Group: 0x02, Register: 0x0007,
			InstanceMin: min, InstanceMax: max, NativeName: "circuit_flow_temperature_setpoint", NativeUnit: "°C",
			Conversion: "namespace", Grade: "Hypothesis", Families: families, SourceRank: sourceRankPrimary,
			Gap: "circuit flow setpoint has no pack element",
		},
		{
			ID: "circuit.state", Subject: "circuit", Group: 0x02, Register: 0x001E,
			InstanceMin: min, InstanceMax: max, NativeName: "circuit_state",
			Conversion: "namespace", Grade: "Hypothesis", Families: families, SourceRank: sourceRankPrimary,
			Gap: "circuit state has no pack element",
		},
		{
			ID: "circuit.pump_status", Subject: "circuit", Group: 0x02, Register: 0x001B,
			InstanceMin: min, InstanceMax: max, NativeName: "circuit_pump_unknown_8",
			Conversion: "namespace", Grade: "Unknown", Families: families, SourceRank: sourceRankPrimary,
			Gap: "circuit pump status has no pack element",
		},
		{
			ID: "circuit.mixer", Subject: "circuit", Group: 0x02, Register: 0x0021,
			InstanceMin: min, InstanceMax: max, NativeName: "circuit_mixer_position", NativeUnit: "%",
			Conversion: "namespace", Grade: "Hypothesis", Families: families, SourceRank: sourceRankPrimary,
			Gap: "circuit mixer position has no pack element",
		},
		{
			ID: "circuit.humidity", Subject: "circuit", Group: 0x02, Register: 0x0022,
			InstanceMin: min, InstanceMax: max, NativeName: "circuit_room_humidity", NativeUnit: "%",
			Conversion: "namespace", Grade: "Hypothesis", Families: families, SourceRank: sourceRankPrimary,
			Gap: "circuit humidity has no pack element on the circuit subject",
		},
		{
			ID: "circuit.dew_point", Subject: "circuit", Group: 0x02, Register: 0x0023,
			InstanceMin: min, InstanceMax: max, NativeName: "circuit_dew_point_temperature", NativeUnit: "°C",
			Conversion: "namespace", Grade: "Hypothesis", Families: families, SourceRank: sourceRankPrimary,
			Gap: "circuit dew point has no pack element",
		},
		{
			ID: "circuit.pump_hours", Subject: "circuit", Group: 0x02, Register: 0x0024,
			InstanceMin: min, InstanceMax: max, NativeName: "circuit_pump_operating_time_total", NativeUnit: "h",
			Conversion: "namespace", Grade: "Unknown", Families: families, SourceRank: sourceRankPrimary,
			Gap: "circuit pump hours have no pack element",
		},
		{
			ID: "circuit.pump_starts", Subject: "circuit", Group: 0x02, Register: 0x0025,
			InstanceMin: min, InstanceMax: max, NativeName: "circuit_pump_start_count_total", NativeUnit: "count",
			Conversion: "namespace", Grade: "Unknown", Families: families, SourceRank: sourceRankPrimary,
			Gap: "circuit pump starts have no pack element",
		},
	}
}

// DeriveCircuitAction maps one documented circuit-state label to the system
// action token. An undocumented label returns an empty token.
func DeriveCircuitAction(label string) (string, bool) {
	switch label {
	case "StandBy":
		return "idle", true
	case "Heating":
		return "heating", true
	case "Cooling":
		return "cooling", true
	case "DHW":
		return "dhw", true
	default:
		return "", false
	}
}
