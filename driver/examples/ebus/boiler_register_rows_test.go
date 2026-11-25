package ebus

import (
	"bytes"
	"testing"
)

func checkBoilerRegisterResponses(t *testing.T, vectors map[string]nativeResponseVector) {
	t.Helper()
	for _, name := range []string{
		"b509_read_opcode_echo", "b509_read_address_echo", "b509_read_bare", "b509_read_other_header_bare",
		"b509_read_header_only", "b509_read_address_only", "b509_read_wrong_width_empty", "b509_read_overwide_native",
	} {
		nativeCase(t, vectors, name, func(t *testing.T, v nativeResponseVector) {
			s := boilerRegisterRead{Target: 0x08, Register: 0x4400}
			reply, _ := onSimulatedBus(t, s.request(0xf0), v.nativeFrame, true)
			x := s.pumpReply(reply.Data)
			if !bytes.Equal(x.Raw, reply.Data) || !bytes.Equal(x.Tail, octets(t, v.Tail)) || (x.PumpOn != nil) != v.OK {
				t.Fatalf("pump result %+v, case %+v", x, v)
			}
			compareDecoded(t, x.PumpOn, v.Decoded)
		})
	}
	for _, c := range []struct {
		name     string
		attempts int
	}{
		{"b509_write_absent_frame", 3}, {"b509_write_present_empty", 1}, {"b509_write_present_zero", 3},
	} {
		nativeCase(t, vectors, c.name, func(t *testing.T, v nativeResponseVector) {
			s := boilerRegisterWrite{Target: 0x08, Register: 0x0200, Encoded: []byte{0x33, 0x44}}
			replayWriteAttempts(t, s.request(0xf0), v, c.attempts, s.acknowledgement)
		})
	}
}
