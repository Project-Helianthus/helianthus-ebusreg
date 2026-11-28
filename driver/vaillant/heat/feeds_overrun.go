package heat

// OverrunFeeds are the special-function rows. A documented code feeds an
// overrun only when its label states that overrun. The three overrun elements
// no value table states are listed and are not fed.
func OverrunFeeds() []Feed {
	families := families700720()
	zoneMin, zoneMax := indexedInstance()
	singleMin, singleMax := singleInstance()
	return []Feed{
		{
			ID: "overrun.zone", Subject: "system", Group: 0x03, Register: 0x001B,
			InstanceMin: zoneMin, InstanceMax: zoneMax, NativeName: "zone_special_function_state",
			Conversion: "symbol", Grade: "Hypothesis", Families: families, SourceRank: sourceRankPrimary, KeepNative: true,
			Codes: map[string][]FeedCode{
				"controller-700": zoneOverrun700(),
				"controller-720": zoneOverrun720(),
			},
		},
		{
			ID: "overrun.circuit", Subject: "system", Group: 0x02, Register: 0x001D,
			InstanceMin: zoneMin, InstanceMax: zoneMax, NativeName: "circuit_unknown_9",
			Conversion: "symbol", Grade: "Hypothesis", Families: families, SourceRank: sourceRankPrimary, KeepNative: true,
			Codes: map[string][]FeedCode{
				"controller-700": {
					{Value: 0, Label: "no quick mode", None: true},
					{Value: 2, Label: "party", Token: "active", Element: "thermal.overrun.party"},
					{Value: 4, Label: "one-day holiday"},
					{Value: 5, Label: "one-day bank holiday"},
					{Value: 6, Label: "cylinder boost"},
					{Value: 255, Label: "no action needed", None: true},
				},
				"controller-720": {
					{Value: 0, Label: "no quick mode", None: true},
					{Value: 6, Label: "one-time cylinder loading", Token: "active", Element: "thermal.overrun.one_time_dhw"},
					{Value: 255, Label: "no action needed", None: true},
				},
			},
		},
		{
			ID: "overrun.dhw", Subject: "system", Group: 0x01, Register: 0x000F,
			InstanceMin: singleMin, InstanceMax: singleMax, NativeName: "dhw_special_function_state",
			Conversion: "symbol", Grade: "Hypothesis", Families: families, SourceRank: sourceRankPrimary, KeepNative: true,
			Codes: map[string][]FeedCode{
				"controller-700": dhwOverrunCodes(),
				"controller-720": dhwOverrunCodes(),
			},
		},
		{
			ID: "overrun.ventilation_other", Subject: "system", Group: 0x09, Register: 0x0004,
			InstanceMin: zoneMin, InstanceMax: zoneMax, NativeName: "ventilation_special_function_state",
			Conversion: "symbol", Grade: "Unknown", Families: []string{"controller-720"}, SourceRank: sourceRankPrimary, KeepNative: true,
			Codes: map[string][]FeedCode{
				"controller-720": {
					{Value: 0, Label: "NoSpecialFunctionActive", None: true},
					{Value: 1, Label: "VentilationBoost"},
					{Value: 7, Label: "Holiday"},
					{Value: 10, Label: "SystemOff", Token: "active", Element: "thermal.overrun.hvac_system_off"},
				},
			},
		},
		{
			ID: "overrun.ventilation", Subject: "system", Group: 0x00, Register: 0x00A4,
			InstanceMin: singleMin, InstanceMax: singleMax, NativeName: "ventilation_special_function_state",
			Conversion: "symbol", Grade: "Unknown", Families: []string{"controller-700"}, SourceRank: sourceRankPrimary, KeepNative: true,
			Codes: map[string][]FeedCode{
				"controller-700": {
					{Value: 0, Label: "No special function active", None: true},
					{Value: 1, Label: "Ventilation boost"},
					{Value: 3, Label: "Party", Token: "active", Element: "thermal.overrun.party"},
					{Value: 5, Label: "One day holiday"},
					{Value: 6, Label: "One day bank holiday"},
					{Value: 7, Label: "Holiday"},
					{Value: 8, Label: "BankHoliday"},
					{Value: 10, Label: "SystemOff", Token: "active", Element: "thermal.overrun.hvac_system_off"},
				},
			},
		},
		{
			ID: "overrun.one_time_ventilation", Subject: "system", NativeName: "one_time_ventilation",
			NoMatch: true, Grade: "Unknown", SourceRank: sourceRankPrimary,
			Gap: "thermal.overrun.one_time_ventilation is not stated by a value table",
		},
		{
			ID: "overrun.one_day_away", Subject: "system", NativeName: "one_day_away",
			NoMatch: true, Grade: "Unknown", SourceRank: sourceRankPrimary,
			Gap: "thermal.overrun.one_day_away is not stated by a value table",
		},
		{
			ID: "overrun.one_day_at_home", Subject: "system", NativeName: "one_day_at_home",
			NoMatch: true, Grade: "Unknown", SourceRank: sourceRankPrimary,
			Gap: "thermal.overrun.one_day_at_home is not stated by a value table",
		},
	}
}

func zoneOverrun700() []FeedCode {
	return []FeedCode{
		{Value: 0, Label: "NoQuickMode", None: true},
		{Value: 1, Label: "VentilationBoost"},
		{Value: 3, Label: "Party", Token: "active", Element: "thermal.overrun.party"},
		{Value: 4, Label: "QuickVeto"},
		{Value: 5, Label: "OneDayHoliday"},
		{Value: 6, Label: "OneDayBankHoliday"},
		{Value: 7, Label: "Holiday"},
		{Value: 8, Label: "BankHoliday"},
		{Value: 10, Label: "SystemOff", Token: "active", Element: "thermal.overrun.hvac_system_off"},
	}
}

func zoneOverrun720() []FeedCode {
	return []FeedCode{
		{Value: 0, Label: "NoSpecialFunctionActive", None: true},
		{Value: 1, Label: "VentilationBoost"},
		{Value: 4, Label: "QuickVeto"},
		{Value: 7, Label: "Holiday"},
		{Value: 10, Label: "SystemOff", Token: "active", Element: "thermal.overrun.hvac_system_off"},
	}
}

func dhwOverrunCodes() []FeedCode {
	return []FeedCode{
		{Value: 0, Label: "NoSpecialFunctionActive", None: true},
		{Value: 3, Label: "Party", Token: "active", Element: "thermal.overrun.party"},
		{Value: 5, Label: "One day holiday"},
		{Value: 6, Label: "One day bank holiday"},
		{Value: 7, Label: "Holiday"},
		{Value: 8, Label: "BankHoliday"},
		{Value: 9, Label: "OneTimeTankloading", Token: "active", Element: "thermal.overrun.one_time_dhw"},
		{Value: 10, Label: "SystemOff", Token: "active", Element: "thermal.overrun.hvac_system_off"},
	}
}
