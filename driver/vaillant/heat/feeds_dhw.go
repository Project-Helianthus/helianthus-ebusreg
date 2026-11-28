package heat

// DhwFeeds are the hot-water table rows. Store temperatures are not created
// from the system or differential sensors: those rows have no cylinder identity.
func DhwFeeds() []Feed {
	min, max := singleInstance()
	families := families700720()
	return []Feed{
		{
			ID: "dhw.temperature", Subject: "dhw", Group: 0x01, Register: 0x0005,
			InstanceMin: min, InstanceMax: max, NativeName: "dhw_temperature", NativeUnit: "°C",
			Element: "thermal.measurement.dhw_temperature", Conversion: "celsius", Grade: "Hypothesis",
			Families: families, SourceRank: sourceRankPrimary,
		},
		{
			ID: "dhw.setpoint", Subject: "dhw", Group: 0x01, Register: 0x0004,
			InstanceMin: min, InstanceMax: max, NativeName: "dhw_temperature_setpoint", NativeUnit: "°C",
			Element: "thermal.setpoint.dhw_temperature_setpoint", Conversion: "celsius", Grade: "Proven",
			Families: families, SourceRank: sourceRankPrimary,
		},
		{
			ID: "dhw.pump_active", Subject: "dhw", Group: 0x01, Register: 0x0006,
			InstanceMin: min, InstanceMax: max, NativeName: "dhw_pump_active",
			Conversion: "namespace", Grade: "Hypothesis", Families: families, SourceRank: sourceRankPrimary,
			Gap: "hot-water pump active stays native",
		},
		{
			ID: "dhw.flow_setpoint", Subject: "dhw", Group: 0x01, Register: 0x0008,
			InstanceMin: min, InstanceMax: max, NativeName: "dhw_flow_temperature_setpoint", NativeUnit: "°C",
			Conversion: "namespace", Grade: "Hypothesis", Families: families, SourceRank: sourceRankPrimary,
			Gap: "hot-water flow setpoint has no pack element on the hot-water subject",
		},
		{
			ID: "dhw.special_function", Subject: "dhw", Group: 0x01, Register: 0x000F,
			InstanceMin: min, InstanceMax: max, NativeName: "dhw_special_function_state",
			Conversion: "namespace", Grade: "Hypothesis", Families: families, SourceRank: sourceRankPrimary,
			Gap: "hot-water special function state stays native",
		},
		{
			ID: "dhw.mode", Subject: "dhw", Group: 0x01, Register: 0x0003,
			InstanceMin: min, InstanceMax: max, NativeName: "dhw_operating_mode",
			Element: "thermal.mode.dhw", Conversion: "symbol", Grade: "Hypothesis",
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
					{Value: 2, Label: "Manual", Token: "manual"},
				},
			},
		},
		{
			ID: "store.top_native", Subject: "store", Group: 0x00, Register: 0x009F,
			InstanceMin: min, InstanceMax: max, NativeName: "heating_cylinder_top_temperature", NativeUnit: "°C",
			Conversion: "namespace", Grade: "Hypothesis", Families: families, SourceRank: sourceRankPrimary,
			Gap: "top sensor has no cylinder identity",
		},
		{
			ID: "store.bottom_native", Subject: "store", Group: 0x00, Register: 0x00A0,
			InstanceMin: min, InstanceMax: max, NativeName: "heating_cylinder_bottom_temperature", NativeUnit: "°C",
			Conversion: "namespace", Grade: "Hypothesis", Families: families, SourceRank: sourceRankPrimary,
			Gap: "bottom sensor has no cylinder identity",
		},
		{
			ID: "store.differential_1", Subject: "store", Group: 0x08, Register: 0x0005,
			InstanceMin: min, InstanceMax: max, NativeName: "differential_control_sensor_1_temperature", NativeUnit: "°C",
			Conversion: "namespace", Grade: "Hypothesis", Families: families, SourceRank: sourceRankPrimary,
			Gap: "differential sensor has no cylinder identity",
		},
		{
			ID: "store.differential_2", Subject: "store", Group: 0x08, Register: 0x0006,
			InstanceMin: min, InstanceMax: max, NativeName: "differential_control_sensor_2_temperature", NativeUnit: "°C",
			Conversion: "namespace", Grade: "Hypothesis", Families: families, SourceRank: sourceRankPrimary,
			Gap: "differential sensor has no cylinder identity",
		},
	}
}
