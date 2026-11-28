package heat

// HeatPumpFeeds are the heat-pump table rows. A register whose reference,
// codec or unit is not stated does not publish a pack fact.
func HeatPumpFeeds() []Feed {
	min, max := singleInstance()
	return []Feed{
		{
			ID: "heatpump.high_pressure", Subject: "refrigerant", Selector: "14 00",
			InstanceMin: min, InstanceMax: max, NativeName: "refrigerant_high_pressure", NativeUnit: "bar",
			Element: "thermal.measurement.refrigerant_high_side_pressure", Conversion: "bar", Grade: "Hypothesis",
			SourceRank: sourceRankPrimary, Hold: "pressure reference is not stated",
		},
		{
			ID: "heatpump.low_pressure", Subject: "refrigerant", Selector: "15 00",
			InstanceMin: min, InstanceMax: max, NativeName: "refrigerant_low_pressure", NativeUnit: "bar",
			Element: "thermal.measurement.refrigerant_low_side_pressure", Conversion: "bar", Grade: "Hypothesis",
			SourceRank: sourceRankPrimary, Hold: "pressure reference is not stated",
		},
		{
			ID: "heatpump.high_pressure_interface", Subject: "refrigerant", Selector: "54 02 00 6A 09",
			InstanceMin: min, InstanceMax: max, NativeName: "refrigerant_high_pressure", NativeUnit: "bar",
			Element: "thermal.measurement.refrigerant_high_side_pressure", Conversion: "bar", Grade: "Hypothesis",
			SourceRank: sourceRankSecondary, Hold: "pressure reference is not stated",
		},
		{
			ID: "heatpump.condensing", Subject: "refrigerant", Selector: "54 02 00 A6 0E",
			InstanceMin: min, InstanceMax: max, NativeName: "refrigerant_condensation_temperature", NativeUnit: "°C",
			Element: "thermal.measurement.refrigerant_condensing_temperature", Conversion: "celsius", Grade: "Unknown",
			SourceRank: sourceRankPrimary, Hold: "condensation name is not final",
		},
		{
			ID: "heatpump.evaporating", Subject: "refrigerant", Selector: "54 02 00 A7 0E",
			InstanceMin: min, InstanceMax: max, NativeName: "refrigerant_evaporation_temperature", NativeUnit: "°C",
			Element: "thermal.measurement.refrigerant_evaporating_temperature", Conversion: "celsius", Grade: "Unknown",
			SourceRank: sourceRankPrimary, Hold: "evaporation codec is not stated",
		},
		{
			ID: "heatpump.compressor_speed", Subject: "compressor", Selector: "54 02 00 0D 0A",
			InstanceMin: min, InstanceMax: max, NativeName: "compressor_speed",
			Conversion: "namespace", Grade: "Unknown", SourceRank: sourceRankPrimary,
			Gap: "compressor speed name is not final and no unit is stated",
		},
		{
			ID: "heatpump.compressor_state", Subject: "generator", Selector: "E4 00",
			InstanceMin: min, InstanceMax: max, NativeName: "compressor_state",
			Conversion: "namespace", Grade: "Unknown", SourceRank: sourceRankPrimary,
			Gap: "compressor state table is not stated, so heat-generator operation stays unfed",
		},
		{
			ID: "heatpump.heat_output", Subject: "generator", NativeName: "heat_output",
			NoMatch: true, Grade: "Unknown", SourceRank: sourceRankPrimary,
			Gap: "thermal.energy.heat_output has no documented register",
		},
		{
			ID: "heatpump.thermal_power", Subject: "generator", NativeName: "thermal_power",
			NoMatch: true, Grade: "Unknown", SourceRank: sourceRankPrimary,
			Gap: "thermal.measurement.thermal_power has no documented register",
		},
	}
}
