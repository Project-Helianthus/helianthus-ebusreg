package input

import (
	"encoding/binary"
	"errors"
	"fmt"
	"time"
)

const (
	scheduleConfigReplyLength = 9
	scheduleCountsReplyLength = 9
	scheduleTimerReplyLength  = 7
	scheduleSentinelByte      = 0xff
	scheduleSentinelWord      = 0xffff
	scheduleNoResponseText    = "b555 timer write: no response"
	scheduleEmptyResponseText = "b555 timer write: empty response"
	scheduleUnknownStatusText = "unknown error 0x%02X"
	scheduleTemperatureUnit   = "degC"
	scheduleWeekdays          = 7
)

var scheduleStatusText = map[byte]string{
	0x00: "accepted",
	0x01: "parameter out of range",
	0x03: "timer unavailable",
	0x06: "validation failure",
}

// ErrScheduleSelector is returned when a record was acquired for another B555 request than the row's.
var ErrScheduleSelector = errors.New("record selector is not the B555 request of the row")

// ScheduleOperation is the kind of B555 exchange a row stands for.
type ScheduleOperation string

const (
	// ScheduleConfigRead reads the program configuration of one zone and circuit.
	ScheduleConfigRead ScheduleOperation = "config_read"
	// ScheduleCountsRead reads the slots used per weekday of one zone and circuit.
	ScheduleCountsRead ScheduleOperation = "counts_read"
	// ScheduleTimerRead reads one timer slot.
	ScheduleTimerRead ScheduleOperation = "timer_read"
	// ScheduleTimerWrite writes one timer slot and reads back a status byte.
	ScheduleTimerWrite ScheduleOperation = "timer_write"
)

// ScheduleReason is why a B555 exchange produced no value or no accepted status.
type ScheduleReason string

const (
	// ScheduleReasonNone is the reason of an exchange that produced a value or a status.
	ScheduleReasonNone ScheduleReason = ""
	// ScheduleReasonNoReply is an exchange that returned no frame.
	ScheduleReasonNoReply ScheduleReason = "no_reply"
	// ScheduleReasonSendError is an exchange whose send returned an error.
	ScheduleReasonSendError ScheduleReason = "send_error"
	// ScheduleReasonShort is a reply shorter than the layout of the row.
	ScheduleReasonShort ScheduleReason = "too_short"
	// ScheduleReasonStatus is a reply whose status byte is not zero.
	ScheduleReasonStatus ScheduleReason = "status_not_zero"
	// ScheduleReasonEmpty is a write reply without any byte.
	ScheduleReasonEmpty ScheduleReason = "empty_response"
)

// ScheduleRow is one B555 request of the poller: its opcode, the shortest reply its layout needs and the native field it fills.
type ScheduleRow struct {
	ID          string
	Operation   ScheduleOperation
	Opcode      byte
	MinReply    int
	NativeField string
}

// ScheduleRows lists the four B555 requests the poller sends: the three reads and the timer write.
var ScheduleRows = []ScheduleRow{
	{ID: "b555.schedule_config_read", Operation: ScheduleConfigRead, Opcode: scheduleConfigOpcode, MinReply: scheduleConfigReplyLength, NativeField: "vaillantScheduleConfig"},
	{ID: "b555.schedule_counts_read", Operation: ScheduleCountsRead, Opcode: scheduleCountsOpcode, MinReply: scheduleCountsReplyLength, NativeField: "vaillantScheduleProgram.SlotsUsed"},
	{ID: "b555.schedule_timer_read", Operation: ScheduleTimerRead, Opcode: scheduleReadOpcode, MinReply: scheduleTimerReplyLength, NativeField: "vaillantScheduleDay.Slots (vaillantScheduleTimerSlot)"},
	{ID: "b555.schedule_timer_write", Operation: ScheduleTimerWrite, Opcode: scheduleWriteOpcode, NativeField: "vaillantTimeProgramSlotResult (status byte) and the program re-read"},
}

// ScheduleConfig is the native result of a configuration read. The limits are whole degrees Celsius and nil when the reply carries the unavailable byte.
type ScheduleConfig struct {
	MaxSlots, TimeResolution, MinDuration int
	HasTemperature                        bool
	TempSlots                             int
	MinTempC, MaxTempC                    *float64
}

// ScheduleTimerSlot is the native result of a timer read. The temperature is nil when the reply carries the unavailable word.
type ScheduleTimerSlot struct {
	StartHour, StartMinute, EndHour, EndMinute int
	TemperatureC                               *float64
	TemperatureRaw                             *int
}

// ScheduleReadResult is the native result of one acquired B555 read reply: the evidence as received and either a value or the reason there is none.
type ScheduleReadResult struct {
	Row         string
	Request     []byte
	Raw         []byte
	Outcome     ReadOutcome
	Reason      ScheduleReason
	Config      *ScheduleConfig
	SlotsPerDay []int
	Timer       *ScheduleTimerSlot
	Unit        string
	AcquiredAt  time.Time
	Source      byte
	Target      byte
}

// ScheduleWriteResult is the native result of one acquired timer write reply. Status is 0xFF when the exchange produced no status byte.
type ScheduleWriteResult struct {
	Row         string
	Request     []byte
	Raw         []byte
	Status      byte
	Accepted    bool
	Description string
	Reason      ScheduleReason
	Error       string
	AcquiredAt  time.Time
	Source      byte
	Target      byte
}

func (r ScheduleRow) accepts(selector Selector) bool {
	switch s := selector.(type) {
	case ScheduleConfigSelector:
		return r.Opcode == s.Payload()[0]
	case ScheduleCountsSelector:
		return r.Opcode == s.Payload()[0]
	case ScheduleTimerSelector:
		return r.Opcode == s.Payload()[0]
	case ScheduleTimerWriteSelector:
		return r.Opcode == s.Payload()[0]
	}
	return false
}

// ScheduleStatusText returns the text the poller gives a B555 status byte.
func ScheduleStatusText(status byte) string {
	if text, ok := scheduleStatusText[status]; ok {
		return text
	}
	return fmt.Sprintf(scheduleUnknownStatusText, status)
}

// Read evaluates one acquired read reply against the row. The reply is the frame data of the exchange, without any echo of the request.
func (r ScheduleRow) Read(rec Record) (ScheduleReadResult, error) {
	if r.Operation == ScheduleTimerWrite || !r.accepts(rec.Selector) {
		return ScheduleReadResult{}, ErrScheduleSelector
	}
	out := ScheduleReadResult{Row: r.ID, Request: rec.Request, Raw: rec.Reply, Outcome: OutcomeNoValue, AcquiredAt: rec.Clock, Source: rec.Source, Target: rec.Target}
	switch {
	case rec.Diagnostics.SendError != "":
		out.Reason = ScheduleReasonSendError
		return out, nil
	case !rec.Replied:
		out.Reason = ScheduleReasonNoReply
		return out, nil
	case len(rec.Reply) < r.MinReply:
		out.Reason = ScheduleReasonShort
		return out, nil
	case rec.Reply[0] != 0x00:
		out.Reason = ScheduleReasonStatus
		return out, nil
	}
	out.Outcome = OutcomeValue
	switch r.Operation {
	case ScheduleConfigRead:
		out.Config, out.Unit = scheduleConfig(rec.Reply), scheduleTemperatureUnit
	case ScheduleCountsRead:
		out.SlotsPerDay = scheduleCounts(rec.Reply)
	default:
		out.Timer, out.Unit = scheduleTimer(rec.Reply), scheduleTemperatureUnit
	}
	return out, nil
}

func scheduleLimitOf(b byte) *float64 {
	if b == scheduleSentinelByte {
		return nil
	}
	v := float64(b)
	return &v
}

func scheduleConfig(data []byte) *ScheduleConfig {
	return &ScheduleConfig{
		MaxSlots: int(data[1]), TimeResolution: int(data[2]), MinDuration: int(data[3]),
		HasTemperature: data[4] != 0, TempSlots: int(data[5]),
		MinTempC: scheduleLimitOf(data[6]), MaxTempC: scheduleLimitOf(data[7]),
	}
}

func scheduleCounts(data []byte) []int {
	slots := make([]int, scheduleWeekdays)
	for i := range slots {
		slots[i] = int(data[i+1])
	}
	return slots
}

func scheduleTimer(data []byte) *ScheduleTimerSlot {
	slot := &ScheduleTimerSlot{StartHour: int(data[1]), StartMinute: int(data[2]), EndHour: int(data[3]), EndMinute: int(data[4])}
	if raw := binary.LittleEndian.Uint16(data[5:7]); raw != scheduleSentinelWord {
		degrees, tenths := float64(raw)/10.0, int(raw)
		slot.TemperatureC, slot.TemperatureRaw = &degrees, &tenths
	}
	return slot
}

// Acknowledge evaluates one acquired timer write exchange against the write row. A write whose status byte is not zero is a result, not an error.
func (r ScheduleRow) Acknowledge(rec Record) (ScheduleWriteResult, error) {
	if r.Operation != ScheduleTimerWrite || !r.accepts(rec.Selector) {
		return ScheduleWriteResult{}, ErrScheduleSelector
	}
	out := ScheduleWriteResult{Row: r.ID, Request: rec.Request, Raw: rec.Reply, Status: scheduleSentinelByte, AcquiredAt: rec.Clock, Source: rec.Source, Target: rec.Target}
	switch {
	case rec.Diagnostics.SendError != "":
		out.Reason, out.Error = ScheduleReasonSendError, scheduleNoResponseText
	case !rec.Replied:
		out.Reason, out.Error = ScheduleReasonNoReply, scheduleNoResponseText
	case len(rec.Reply) == 0:
		out.Reason, out.Error = ScheduleReasonEmpty, scheduleEmptyResponseText
	default:
		out.Status = rec.Reply[0]
		out.Accepted = out.Status == 0x00
		out.Description = ScheduleStatusText(out.Status)
	}
	return out, nil
}
