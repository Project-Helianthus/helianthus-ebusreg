package owner

import (
	"context"
	"fmt"
	"math"
)

// The B555 reads and the time-program writes of the poller use the native
// schedule types below; the publication owner maps them to the GraphQL
// schedules and the writer adapters map them to the MCP time-program forms.

// vaillantScheduleTimerSlot is one timer slot read from B555 A5.
// TemperatureRaw holds the raw deci-degree value; both temperature fields
// are nil when the slot carries the unavailable sentinel.
type vaillantScheduleTimerSlot struct {
	StartHour      int
	StartMinute    int
	EndHour        int
	EndMinute      int
	TemperatureC   *float64
	TemperatureRaw *int
}

// vaillantScheduleDay is the timer slots of one weekday of a program.
type vaillantScheduleDay struct {
	Weekday string
	Slots   []vaillantScheduleTimerSlot
}

// vaillantScheduleConfig is the B555 A3 configuration of a program.
type vaillantScheduleConfig struct {
	MaxSlots       int
	TimeResolution int
	MinDuration    int
	HasTemperature bool
	TempSlots      int
	MinTempC       *float64
	MaxTempC       *float64
}

// vaillantScheduleProgram is one (zone, heating circuit kind) program with
// its configuration, slots per day and day programs.
type vaillantScheduleProgram struct {
	Zone      int
	HC        string
	Config    *vaillantScheduleConfig
	SlotsUsed []int
	Days      []vaillantScheduleDay
}

// vaillantTimeProgramSlot is one slot of a time-program write request.
type vaillantTimeProgramSlot struct {
	StartHour    int
	StartMinute  int
	EndHour      int
	EndMinute    int
	TemperatureC *float64
}

// vaillantTimeProgramSlotResult is the B555 outcome of one written slot.
type vaillantTimeProgramSlotResult struct {
	SlotIndex int
	Accepted  bool
	ErrorCode int
	ErrorDesc string
}

// vaillantTimeProgramWriteResult is the outcome of a time-program write.
type vaillantTimeProgramWriteResult struct {
	Success     bool
	SlotResults []vaillantTimeProgramSlotResult
	Error       string
}

// VaillantScheduleWriter writes zone and hot water time programs in native
// form; the poller implements it.
type VaillantScheduleWriter interface {
	SetZoneTimeProgram(ctx context.Context, zone int, weekday int, slots []vaillantTimeProgramSlot) (*vaillantTimeProgramWriteResult, error)
	SetDhwTimeProgram(ctx context.Context, weekday int, slots []vaillantTimeProgramSlot) (*vaillantTimeProgramWriteResult, error)
}

var b555ErrorDescriptions = map[byte]string{
	0x00: "accepted",
	0x01: "parameter out of range",
	0x03: "timer unavailable",
	0x06: "validation failure",
}

func b555ErrorDesc(code byte) string {
	if desc, ok := b555ErrorDescriptions[code]; ok {
		return desc
	}
	return fmt.Sprintf("unknown error 0x%02X", code)
}

func (p *vaillantSemanticPoller) setZoneTimeProgram(ctx context.Context, zone int, weekday int, slots []vaillantTimeProgramSlot) (*vaillantTimeProgramWriteResult, error) {
	if p == nil {
		return nil, fmt.Errorf("schedule writer unavailable")
	}
	if zone < 0 || zone > 2 {
		return &vaillantTimeProgramWriteResult{Success: false, Error: "zone must be 0-2"}, nil
	}
	if weekday < 0 || weekday > 6 {
		return &vaillantTimeProgramWriteResult{Success: false, Error: "weekday must be 0-6"}, nil
	}
	if len(slots) == 0 {
		return &vaillantTimeProgramWriteResult{Success: false, Error: "slots array must not be empty"}, nil
	}

	cfg := p.readB555Config(ctx, byte(zone), b555HCHeating)
	if cfg == nil {
		return &vaillantTimeProgramWriteResult{Success: false, Error: "unable to read B555 config for zone heating"}, nil
	}

	maxSlots := cfg.MaxSlots
	if maxSlots <= 0 {
		maxSlots = 3
	}
	if len(slots) > maxSlots {
		return &vaillantTimeProgramWriteResult{
			Success: false,
			Error:   fmt.Sprintf("slot count %d exceeds max_slots %d", len(slots), maxSlots),
		}, nil
	}

	minTemp := 5.0
	maxTemp := 30.0
	if cfg.MinTempC != nil {
		minTemp = *cfg.MinTempC
	}
	if cfg.MaxTempC != nil {
		maxTemp = *cfg.MaxTempC
	}

	for i, slot := range slots {
		if slot.StartHour < 0 || slot.StartHour > 24 {
			return &vaillantTimeProgramWriteResult{Success: false, Error: fmt.Sprintf("slot %d: start_hour %d out of range 0-24", i, slot.StartHour)}, nil
		}
		if slot.StartMinute < 0 || slot.StartMinute > 59 {
			return &vaillantTimeProgramWriteResult{Success: false, Error: fmt.Sprintf("slot %d: start_minute %d out of range 0-59", i, slot.StartMinute)}, nil
		}
		if slot.EndHour < 0 || slot.EndHour > 24 {
			return &vaillantTimeProgramWriteResult{Success: false, Error: fmt.Sprintf("slot %d: end_hour %d out of range 0-24", i, slot.EndHour)}, nil
		}
		if slot.EndMinute < 0 || slot.EndMinute > 59 {
			return &vaillantTimeProgramWriteResult{Success: false, Error: fmt.Sprintf("slot %d: end_minute %d out of range 0-59", i, slot.EndMinute)}, nil
		}
		if slot.TemperatureC == nil {
			return &vaillantTimeProgramWriteResult{Success: false, Error: fmt.Sprintf("slot %d: temperature_c is required for heating", i)}, nil
		}
		if *slot.TemperatureC < minTemp || *slot.TemperatureC > maxTemp {
			return &vaillantTimeProgramWriteResult{Success: false, Error: fmt.Sprintf("slot %d: temperature_c %.1f out of range [%.0f, %.0f]", i, *slot.TemperatureC, minTemp, maxTemp)}, nil
		}
	}

	result := &vaillantTimeProgramWriteResult{
		Success:     true,
		SlotResults: make([]vaillantTimeProgramSlotResult, len(slots)),
	}

	for i, slot := range slots {
		tempRaw := uint16(math.Round(*slot.TemperatureC * 10))
		code, err := p.writeB555TimerSlot(ctx,
			byte(zone), b555HCHeating, byte(weekday),
			byte(i), byte(len(slots)),
			byte(slot.StartHour), byte(slot.StartMinute),
			byte(slot.EndHour), byte(slot.EndMinute),
			tempRaw,
		)
		if err != nil {
			result.SlotResults[i] = vaillantTimeProgramSlotResult{
				SlotIndex: i,
				Accepted:  false,
				ErrorCode: 0xFF,
				ErrorDesc: err.Error(),
			}
			result.Success = false
			continue
		}
		accepted := code == 0x00
		if !accepted {
			result.Success = false
		}
		result.SlotResults[i] = vaillantTimeProgramSlotResult{
			SlotIndex: i,
			Accepted:  accepted,
			ErrorCode: int(code),
			ErrorDesc: b555ErrorDesc(code),
		}
	}

	go p.refreshScheduleForProgram(context.Background(), byte(zone), b555HCHeating, byte(weekday))

	return result, nil
}

func (p *vaillantSemanticPoller) setDhwTimeProgram(ctx context.Context, weekday int, slots []vaillantTimeProgramSlot) (*vaillantTimeProgramWriteResult, error) {
	if p == nil {
		return nil, fmt.Errorf("schedule writer unavailable")
	}
	if weekday < 0 || weekday > 6 {
		return &vaillantTimeProgramWriteResult{Success: false, Error: "weekday must be 0-6"}, nil
	}
	if len(slots) == 0 {
		return &vaillantTimeProgramWriteResult{Success: false, Error: "slots array must not be empty"}, nil
	}

	cfg := p.readB555Config(ctx, b555ZoneDHW, b555HCDHW)
	if cfg == nil {
		return &vaillantTimeProgramWriteResult{Success: false, Error: "unable to read B555 config for DHW"}, nil
	}

	maxSlots := cfg.MaxSlots
	if maxSlots <= 0 {
		maxSlots = 3
	}
	if len(slots) > maxSlots {
		return &vaillantTimeProgramWriteResult{
			Success: false,
			Error:   fmt.Sprintf("slot count %d exceeds max_slots %d", len(slots), maxSlots),
		}, nil
	}

	minTemp := 35.0
	maxTemp := 65.0
	if cfg.MinTempC != nil {
		minTemp = *cfg.MinTempC
	}
	if cfg.MaxTempC != nil {
		maxTemp = *cfg.MaxTempC
	}

	for i, slot := range slots {
		if slot.StartHour < 0 || slot.StartHour > 24 {
			return &vaillantTimeProgramWriteResult{Success: false, Error: fmt.Sprintf("slot %d: start_hour %d out of range 0-24", i, slot.StartHour)}, nil
		}
		if slot.StartMinute < 0 || slot.StartMinute > 59 {
			return &vaillantTimeProgramWriteResult{Success: false, Error: fmt.Sprintf("slot %d: start_minute %d out of range 0-59", i, slot.StartMinute)}, nil
		}
		if slot.EndHour < 0 || slot.EndHour > 24 {
			return &vaillantTimeProgramWriteResult{Success: false, Error: fmt.Sprintf("slot %d: end_hour %d out of range 0-24", i, slot.EndHour)}, nil
		}
		if slot.EndMinute < 0 || slot.EndMinute > 59 {
			return &vaillantTimeProgramWriteResult{Success: false, Error: fmt.Sprintf("slot %d: end_minute %d out of range 0-59", i, slot.EndMinute)}, nil
		}
		if slot.TemperatureC != nil {
			if *slot.TemperatureC < minTemp || *slot.TemperatureC > maxTemp {
				return &vaillantTimeProgramWriteResult{Success: false, Error: fmt.Sprintf("slot %d: temperature_c %.1f out of range [%.0f, %.0f]", i, *slot.TemperatureC, minTemp, maxTemp)}, nil
			}
		}
	}

	result := &vaillantTimeProgramWriteResult{
		Success:     true,
		SlotResults: make([]vaillantTimeProgramSlotResult, len(slots)),
	}

	for i, slot := range slots {
		var tempRaw uint16
		if slot.TemperatureC != nil {
			tempRaw = uint16(math.Round(*slot.TemperatureC * 10))
		} else {
			tempRaw = 0xFFFF
		}
		code, err := p.writeB555TimerSlot(ctx,
			b555ZoneDHW, b555HCDHW, byte(weekday),
			byte(i), byte(len(slots)),
			byte(slot.StartHour), byte(slot.StartMinute),
			byte(slot.EndHour), byte(slot.EndMinute),
			tempRaw,
		)
		if err != nil {
			result.SlotResults[i] = vaillantTimeProgramSlotResult{
				SlotIndex: i,
				Accepted:  false,
				ErrorCode: 0xFF,
				ErrorDesc: err.Error(),
			}
			result.Success = false
			continue
		}
		accepted := code == 0x00
		if !accepted {
			result.Success = false
		}
		result.SlotResults[i] = vaillantTimeProgramSlotResult{
			SlotIndex: i,
			Accepted:  accepted,
			ErrorCode: int(code),
			ErrorDesc: b555ErrorDesc(code),
		}
	}

	go p.refreshScheduleForProgram(context.Background(), b555ZoneDHW, b555HCDHW, byte(weekday))

	return result, nil
}
