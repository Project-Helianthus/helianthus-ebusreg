package output

// ScheduleProgramsFromNative maps native programs to the schedule output.
func ScheduleProgramsFromNative(programs []NativeScheduleProgram) []ScheduleProgram {
	if programs == nil {
		return nil
	}
	out := make([]ScheduleProgram, len(programs))
	for i, prog := range programs {
		out[i] = ScheduleProgram{
			Zone:      prog.Zone,
			HC:        prog.HC,
			Config:    scheduleConfigFromNative(prog.Config),
			SlotsUsed: prog.SlotsUsed,
		}
		if prog.Days != nil {
			out[i].Days = make([]ScheduleDayProgram, len(prog.Days))
			for d, day := range prog.Days {
				out[i].Days[d] = ScheduleDayFromNative(day)
			}
		}
	}
	return out
}

func scheduleConfigFromNative(cfg *NativeScheduleConfig) *ScheduleConfig {
	if cfg == nil {
		return nil
	}
	return &ScheduleConfig{
		MaxSlots:       cfg.MaxSlots,
		TimeResolution: cfg.TimeResolution,
		MinDuration:    cfg.MinDuration,
		HasTemperature: cfg.HasTemperature,
		TempSlots:      cfg.TempSlots,
		MinTempC:       CloneFloatPtr(cfg.MinTempC),
		MaxTempC:       CloneFloatPtr(cfg.MaxTempC),
	}
}

// ScheduleDayFromNative maps one native day of intervals.
func ScheduleDayFromNative(day NativeScheduleDay) ScheduleDayProgram {
	out := ScheduleDayProgram{Weekday: day.Weekday}
	if day.Slots != nil {
		out.Slots = make([]ScheduleTimerSlot, len(day.Slots))
		for i, slot := range day.Slots {
			out.Slots[i] = ScheduleTimerSlot{
				StartHour:      slot.StartHour,
				StartMinute:    slot.StartMinute,
				EndHour:        slot.EndHour,
				EndMinute:      slot.EndMinute,
				TemperatureC:   CloneFloatPtr(slot.TemperatureC),
				TemperatureRaw: CloneIntPtr(slot.TemperatureRaw),
			}
		}
	}
	return out
}

// TimeProgramSlotsFromMCP maps MCP write request slots to native slots.
func TimeProgramSlotsFromMCP(slots []MCPTimeProgramSlot) []NativeTimeProgramSlot {
	if slots == nil {
		return nil
	}
	out := make([]NativeTimeProgramSlot, len(slots))
	for i, slot := range slots {
		out[i] = NativeTimeProgramSlot(slot)
	}
	return out
}

// MCPTimeProgramResult maps a native result while preserving the write error.
func MCPTimeProgramResult(result *NativeTimeProgramWriteResult, err error) (*MCPTimeProgramWriteResult, error) {
	if result == nil {
		return nil, err
	}
	out := &MCPTimeProgramWriteResult{Success: result.Success, Error: result.Error}
	if result.SlotResults != nil {
		out.SlotResults = make([]MCPTimeProgramSlotResult, len(result.SlotResults))
		for i, slot := range result.SlotResults {
			out.SlotResults[i] = MCPTimeProgramSlotResult(slot)
		}
	}
	return out, err
}
