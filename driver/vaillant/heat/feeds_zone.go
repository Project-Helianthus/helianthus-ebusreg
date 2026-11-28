package heat

// ZoneFeeds are the zone table rows. A row whose pack subject does not match
// stays in the driver namespace and is listed in its Gap.
func ZoneFeeds() []Feed {
	min, max := indexedInstance()
	families := families700720()
	return []Feed{
		{
			ID: "zone.air_temperature", Subject: "zone", Group: 0x03, Register: 0x000F,
			InstanceMin: min, InstanceMax: max, NativeName: "zone_room_temperature", NativeUnit: "°C",
			Element: "thermal.measurement.air_temperature", Conversion: "celsius", Grade: "Hypothesis",
			Families: families, SourceRank: sourceRankPrimary,
		},
		{
			ID: "zone.manual_setpoint", Subject: "zone", Group: 0x03, Register: 0x0022,
			InstanceMin: min, InstanceMax: max, NativeName: "heating_zone_room_manual_temperature_setpoint", NativeUnit: "°C",
			Element: "thermal.setpoint.temperature_setpoint", Conversion: "celsius", Grade: "Hypothesis",
			Families: families, SourceRank: sourceRankPrimary,
		},
		{
			ID: "zone.effective_setpoint", Subject: "zone", Group: 0x03, Register: 0x0014,
			InstanceMin: min, InstanceMax: max, NativeName: "heating_zone_room_effective_temperature_setpoint", NativeUnit: "°C",
			Conversion: "namespace", Grade: "Hypothesis", Families: families, SourceRank: sourceRankSecondary,
			Gap: "heating computed effective is a separate value and the pack has no element for it",
		},
		{
			ID: "zone.humidity", Subject: "zone", Group: 0x03, Register: 0x0028,
			InstanceMin: min, InstanceMax: max, NativeName: "zone_room_humidity", NativeUnit: "%",
			Element: "thermal.measurement.humidity", Conversion: "percent", Grade: "Hypothesis",
			Families: families, SourceRank: sourceRankPrimary,
		},
		{
			ID: "zone.heating_mode", Subject: "zone", Group: 0x03, Register: 0x0006,
			InstanceMin: min, InstanceMax: max, NativeName: "heating_zone_operating_mode",
			Element: "thermal.mode.zone", Conversion: "symbol", Grade: "Hypothesis",
			Families: families, SourceRank: sourceRankPrimary,
			Codes: map[string][]FeedCode{
				"controller-700": {
					{Value: 0, Label: "Off", Token: "off"},
					{Value: 1, Label: "Auto", Token: "auto"},
					{Value: 2, Label: "Day"},
					{Value: 3, Label: "Setback"},
				},
				"controller-720": {
					{Value: 0, Label: "Off", Token: "off"},
					{Value: 1, Label: "TimeControlled"},
					{Value: 2, Label: "Manual"},
				},
			},
		},
		{
			ID: "zone.cooling_mode", Subject: "zone", Group: 0x03, Register: 0x0001,
			InstanceMin: min, InstanceMax: max, NativeName: "cooling_zone_operating_mode",
			Element: "thermal.mode.zone", Conversion: "symbol", Grade: "Hypothesis",
			Families: families, SourceRank: sourceRankPrimary,
			Codes: map[string][]FeedCode{
				"controller-700": {
					{Value: 0, Label: "Off", Token: "off"},
					{Value: 1, Label: "Auto", Token: "auto"},
					{Value: 2, Label: "Day"},
				},
				"controller-720": {
					{Value: 0, Label: "Off", Token: "off"},
					{Value: 1, Label: "TimeControlled"},
					{Value: 2, Label: "Manual"},
				},
			},
		},
		{
			ID: "zone.comfort_setpoint", Subject: "zone", Group: 0x03, Register: 0x0007,
			InstanceMin: min, InstanceMax: max, NativeName: "heating_zone_room_comfort_temperature_setpoint", NativeUnit: "°C",
			Conversion: "namespace", Grade: "Unknown", Families: []string{"controller-700"}, SourceRank: sourceRankPrimary,
			Gap: "comfort setpoint name is not final, so occupied heating stays unfed",
		},
		{
			ID: "zone.comfort_setpoint_other", Subject: "zone", Group: 0x03, Register: 0x0007,
			InstanceMin: min, InstanceMax: max, NativeName: "heating_zone_room_comfort_temperature_setpoint",
			Conversion: "namespace", Grade: "Unknown", Families: []string{"controller-720"}, SourceRank: sourceRankPrimary,
			Gap: "comfort setpoint name is not final and the codec is not stated for this family",
		},
		{
			ID: "zone.unoccupied_heating", Subject: "zone", Group: 0x03, Register: 0x0009,
			InstanceMin: min, InstanceMax: max, NativeName: "heating_zone_room_setback_temperature_setpoint", NativeUnit: "°C",
			Element: "thermal.setpoint.unoccupied_heating", Conversion: "celsius", Grade: "Hypothesis",
			Families: families, SourceRank: sourceRankPrimary,
		},
		{
			ID: "zone.occupied_cooling_manual", Subject: "zone", Group: 0x03, Register: 0x0023,
			InstanceMin: min, InstanceMax: max, NativeName: "cooling_zone_room_manual_temperature_setpoint", NativeUnit: "°C",
			Element: "thermal.setpoint.occupied_cooling", Conversion: "celsius", Grade: "Hypothesis",
			Families: families, SourceRank: sourceRankPrimary,
		},
		{
			ID: "zone.occupied_cooling_effective", Subject: "zone", Group: 0x03, Register: 0x0015,
			InstanceMin: min, InstanceMax: max, NativeName: "cooling_zone_room_effective_temperature_setpoint", NativeUnit: "°C",
			Element: "thermal.setpoint.occupied_cooling", Conversion: "celsius", Grade: "Hypothesis",
			Families: families, SourceRank: sourceRankSecondary,
		},
		{
			ID: "zone.cooling_setpoint_native", Subject: "zone", Group: 0x03, Register: 0x0002,
			InstanceMin: min, InstanceMax: max, NativeName: "cooling_zone_room_temperature_setpoint", NativeUnit: "°C",
			Conversion: "namespace", Grade: "Hypothesis", Families: families, SourceRank: sourceRankPrimary,
			Gap: "cooling setpoint states no unoccupied role",
		},
		{
			ID: "zone.special_function_native", Subject: "zone", Group: 0x03, Register: 0x001B,
			InstanceMin: min, InstanceMax: max, NativeName: "zone_special_function_state",
			Conversion: "namespace", Grade: "Hypothesis", Families: families, SourceRank: sourceRankPrimary, KeepNative: true,
			Gap: "zone special function state stays native",
		},
		{
			ID: "zone.valve_position", Subject: "terminal", Group: 0x03, Register: 0x0012,
			InstanceMin: min, InstanceMax: max, NativeName: "zone_valve_position", NativeUnit: "1",
			Element: "thermal.measurement.valve_position", Conversion: "valve", Grade: "Hypothesis",
			Families: families, SourceRank: sourceRankPrimary,
		},
	}
}
