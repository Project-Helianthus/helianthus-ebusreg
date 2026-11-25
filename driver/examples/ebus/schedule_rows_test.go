package ebus

import (
	"testing"

	"github.com/Project-Helianthus/helianthus-ebusgo/protocol"
)

func checkScheduleResponses(t *testing.T, vectors map[string]nativeResponseVector) {
	t.Helper()
	p := scheduleProgram{Zone: 1, Circuit: 0}
	s := scheduleSlot{Program: p, Weekday: 2, Index: 0}
	w := scheduleWrite{Slot: s, Count: 1, StartHour: 6, EndHour: 7, EndMinute: 30, TemperatureRaw: 200}
	for _, c := range []struct {
		name    string
		request []byte
		decoder func([]byte) any
	}{
		{"a3_config", p.configRequest(), decodeScheduleConfig}, {"a3_wrong_width", p.configRequest(), decodeScheduleConfig},
		{"a4_slots", p.countsRequest(), decodeScheduleCounts}, {"a4_wrong_width", p.countsRequest(), decodeScheduleCounts},
		{"a5_timer", s.readRequest(), decodeScheduleTimer}, {"a5_no_temperature", s.readRequest(), decodeScheduleTimer}, {"a5_wrong_width", s.readRequest(), decodeScheduleTimer},
		{"a6_absent_frame", w.writeRequest(), nil}, {"a6_present_empty", w.writeRequest(), nil},
		{"a6_zero_status", w.writeRequest(), nil}, {"a6_nonzero_status", w.writeRequest(), nil},
	} {
		nativeCase(t, vectors, c.name, func(t *testing.T, v nativeResponseVector) {
			frame := protocol.Frame{Source: 0xf0, Target: 0x15, Primary: 0xb5, Secondary: 0x55, Data: c.request}
			reply, err := onSimulatedBus(t, frame, v.nativeFrame, v.Present)
			outcome := protocol.BusOutcomeSuccess
			if err != nil {
				outcome = protocol.BusOutcomeTimeout
			}
			if outcome != v.SendOutcome || (reply != nil) != v.Present {
				t.Fatalf("Send outcome %v, present %t; recorded %v, %t", outcome, reply != nil, v.SendOutcome, v.Present)
			}
			if c.decoder != nil {
				if reply == nil {
					t.Fatal("no schedule read reply")
				}
				compareDecoded(t, c.decoder(reply.Data), v.Decoded)
				return
			}
			code, ackErr := scheduleWriteAcknowledgement(reply)
			if v.Code == nil || code != *v.Code || errText(ackErr) != v.Error {
				t.Fatalf("schedule acknowledgement code %02x, error %v; recorded %+v", code, ackErr, v)
			}
			compareDecoded(t, nil, v.Decoded)
		})
	}
}

func decodeScheduleConfig(raw []byte) any { return scheduleConfigReply(raw) }
func decodeScheduleCounts(raw []byte) any { return scheduleCountsReply(raw) }
func decodeScheduleTimer(raw []byte) any  { return scheduleTimerReply(raw) }
