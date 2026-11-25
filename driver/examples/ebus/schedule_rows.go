package ebus

import (
	"encoding/binary"
	"errors"

	"github.com/Project-Helianthus/helianthus-ebusgo/protocol"
)

type scheduleProgram struct{ Zone, Circuit byte }

type scheduleSlot struct {
	Program        scheduleProgram
	Weekday, Index byte
}

type scheduleWrite struct {
	Slot                                       scheduleSlot
	Count                                      byte
	StartHour, StartMinute, EndHour, EndMinute byte
	TemperatureRaw                             uint16
}

func (p scheduleProgram) configRequest() []byte { return []byte{0xa3, p.Zone, p.Circuit} }
func (p scheduleProgram) countsRequest() []byte { return []byte{0xa4, p.Zone, p.Circuit} }
func (s scheduleSlot) readRequest() []byte {
	return []byte{0xa5, s.Program.Zone, s.Program.Circuit, s.Weekday, s.Index}
}
func (w scheduleWrite) writeRequest() []byte {
	s := w.Slot
	return binary.LittleEndian.AppendUint16([]byte{0xa6, s.Program.Zone, s.Program.Circuit, s.Weekday, s.Index, w.Count,
		w.StartHour, w.StartMinute, w.EndHour, w.EndMinute}, w.TemperatureRaw)
}

type scheduleConfig struct {
	MaxSlots, TimeResolution, MinDuration int
	HasTemperature                        bool
	TempSlots                             int
	MinTempC, MaxTempC                    *int
}

type scheduleTimer struct {
	StartHour, StartMinute, EndHour, EndMinute int
	TemperatureRaw                             *uint16
	TemperatureC                               *float64
}

func scheduleConfigReply(raw []byte) *scheduleConfig {
	if len(raw) < 9 || raw[0] != 0 {
		return nil
	}
	return &scheduleConfig{MaxSlots: int(raw[1]), TimeResolution: int(raw[2]), MinDuration: int(raw[3]), HasTemperature: raw[4] != 0,
		TempSlots: int(raw[5]), MinTempC: scheduleLimit(raw[6]), MaxTempC: scheduleLimit(raw[7])}
}

func scheduleLimit(b byte) *int {
	if b == 0xff {
		return nil
	}
	n := int(b)
	return &n
}

func scheduleCountsReply(raw []byte) *[7]int {
	if len(raw) < 9 || raw[0] != 0 {
		return nil
	}
	var counts [7]int
	for i, b := range raw[1:8] {
		counts[i] = int(b)
	}
	return &counts
}

func scheduleTimerReply(raw []byte) *scheduleTimer {
	if len(raw) < 7 || raw[0] != 0 {
		return nil
	}
	s := &scheduleTimer{StartHour: int(raw[1]), StartMinute: int(raw[2]), EndHour: int(raw[3]), EndMinute: int(raw[4])}
	if x := binary.LittleEndian.Uint16(raw[5:7]); x != 0xffff {
		s.TemperatureRaw = &x
		degrees := float64(x) / 10
		s.TemperatureC = &degrees
	}
	return s
}

func scheduleWriteAcknowledgement(reply *protocol.Frame) (byte, error) {
	switch {
	case reply == nil:
		return 0xff, errors.New("b555 timer write: no response")
	case len(reply.Data) == 0:
		return 0xff, errors.New("b555 timer write: empty response")
	default:
		return reply.Data[0], nil
	}
}
