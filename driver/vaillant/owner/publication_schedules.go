package owner

import (
	"github.com/Project-Helianthus/helianthus-ebusreg/driver/vaillant/output"
)

func (o *vaillantSemanticPublication) publishSchedulePrograms(programs []vaillantScheduleProgram) {
	if o == nil || o.state == nil {
		return
	}
	o.state.SetSchedules(&output.ScheduleStatus{Programs: output.ScheduleProgramsFromNative(nativeSchedulePrograms(programs))})
}

func (o *vaillantSemanticPublication) publishScheduleDay(zone int, hc string, weekday int, day vaillantScheduleDay, slotsUsed []int) {
	if o == nil || o.state == nil {
		return
	}
	status := o.state.Schedules()
	if status == nil {
		return
	}
	for i, prog := range status.Programs {
		if prog.Zone != zone || prog.HC != hc {
			continue
		}
		if weekday < len(prog.Days) {
			status.Programs[i].Days[weekday] = output.ScheduleDayFromNative(nativeScheduleDay(day))
		}
		if slotsUsed != nil {
			status.Programs[i].SlotsUsed = slotsUsed
		}
		o.state.SetSchedules(status)
		return
	}
}

func timeProgramSlotsFromMCP(slots []output.MCPTimeProgramSlot) []vaillantTimeProgramSlot {
	if slots == nil {
		return nil
	}
	native := output.TimeProgramSlotsFromMCP(slots)
	out := make([]vaillantTimeProgramSlot, len(native))
	for i, slot := range native {
		out[i] = vaillantTimeProgramSlot(slot)
	}
	return out
}

func mcpTimeProgramWriteResult(result *vaillantTimeProgramWriteResult, err error) (*output.MCPTimeProgramWriteResult, error) {
	if result == nil {
		return nil, err
	}
	native := &output.NativeTimeProgramWriteResult{Success: result.Success, Error: result.Error}
	if result.SlotResults != nil {
		native.SlotResults = make([]output.NativeTimeProgramSlotResult, len(result.SlotResults))
		for i, slot := range result.SlotResults {
			native.SlotResults[i] = output.NativeTimeProgramSlotResult(slot)
		}
	}
	return output.MCPTimeProgramResult(native, err)
}
