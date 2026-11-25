package ebus

import (
	"bytes"
	"testing"

	"github.com/Project-Helianthus/helianthus-ebusgo/protocol"
	"github.com/Project-Helianthus/helianthus-ebusreg/vaillant/nativepoints"
)

func checkLocalRegisterResponses(t *testing.T, vectors map[string]nativeResponseVector) {
	t.Helper()
	for _, name := range []string{"identity_matched", "identity_invalid_reply", "identity_missing_companion"} {
		nativeCase(t, vectors, name, checkInterfaceIdentity)
	}
	for _, c := range []struct {
		name     string
		purpose  registerReadPurpose
		register uint16
		echoed   *byte
	}{
		{"b524_dhw_discovery_header", discoveryProbe, 3, testPointer(byte(0))},
		{"b524_dhw_scalar_header", scalarRead, 3, testPointer(byte(0))},
		{"b524_dhw_discovery_wrong_group", discoveryProbe, 3, nil},
		{"b524_dhw_discovery_wrong_register", discoveryProbe, 3, nil},
		{"b524_explicit_instance_scalar", scalarRead, 4, testPointer(byte(0))},
		{"b524_wrong_instance_scalar", scalarRead, 4, testPointer(byte(2))},
	} {
		nativeCase(t, vectors, c.name, func(t *testing.T, v nativeResponseVector) {
			s := B524Selector{Primary: 0xb5, Secondary: 0x24, Opcode: 2, Group: 1, Register: c.register}
			f := protocol.Frame{Source: 0xf0, Target: 0x15, Primary: s.Primary, Secondary: s.Secondary, Data: s.payload()}
			reply, err := onSimulatedBus(t, f, v.nativeFrame, true)
			x := correlateRegisterRead(s, c.purpose, reply.Data, err)
			ok := len(x.Value) > 0
			if c.purpose == discoveryProbe {
				ok = x.Coherent
			}
			if x.Purpose != v.Purpose || x.Selector != s || !bytes.Equal(x.Raw, reply.Data) || ok != v.OK || !bytes.Equal(x.Value, octets(t, v.Tail)) {
				t.Fatalf("B5/24 typed reply %+v, recorded %+v", x, v)
			}
			if (x.InstanceEcho == nil) != (c.echoed == nil) || (c.echoed != nil && *x.InstanceEcho != *c.echoed) {
				t.Fatalf("unexpected instance echo %v", x.InstanceEcho)
			}
			if v.Celsius != nil {
				p := point("dhw_temperature_setpoint")
				n, err := nativepoints.DecodeNative(p, x.Value)
				if err != nil || n.Native != v.Celsius.String() || n.NativeUnit != "°C" {
					t.Fatalf("native codec %+v, %v vs %s °C", n, err, v.Celsius)
				}
			}
		})
	}
	for _, c := range []struct {
		name     string
		attempts int
	}{{"b524_write_short_ack", 1}, {"b524_write_absent_frame", 3}} {
		nativeCase(t, vectors, c.name, func(t *testing.T, v nativeResponseVector) {
			w := localRegisterWrite{Opcode: 2, Register: 1, Encoded: []byte{0, 0, 0xa0, 0x41}}
			replayWriteAttempts(t, w.request(0xf0, 0x15), v, c.attempts, localRegisterAcknowledgement)
		})
	}
}

func testPointer[T any](v T) *T { return &v }

func checkInterfaceIdentity(t *testing.T, v nativeResponseVector) {
	t.Helper()
	if len(v.Steps) != 2 || v.Session != 1 {
		t.Fatalf("interface identity session %+v", v)
	}
	scope := interfaceIdentityScope{Source: 0x71, Target: 0x15, Observation: 1, Session: 1}
	var records [2]interfaceIdentityValue
	for i, field := range []interfaceIdentityField{interfaceAPI, interfaceRevision} {
		s := interfaceIdentityRead{Scope: scope, Field: field}
		step := v.Steps[i]
		if !step.Present && step.Reply != nil && *step.Reply == "" {
			step.Reply = nil
		}
		reply, _ := onSimulatedBus(t, s.request(), step, step.Present)
		records[i] = s.qualify(reply)
		var raw []byte
		if reply != nil {
			raw = reply.Data
		}
		if records[i].Read != s || records[i].Present != step.Present || records[i].Qualifier != step.Qualified || !bytes.Equal(records[i].Raw, raw) {
			t.Fatalf("interface record %+v, recorded %+v", records[i], step)
		}
	}
	if matchingInterfaceIdentity(records[0], records[1]) != v.Complete || v.ObservationMatched != v.Complete {
		t.Fatal("interface companion coherence differs")
	}
	if !v.Complete {
		return
	}
	api, rev := records[0], records[1]
	for _, c := range []struct {
		name   string
		change func(*interfaceIdentityRead)
	}{
		{"source", func(s *interfaceIdentityRead) { s.Scope.Source++ }},
		{"target", func(s *interfaceIdentityRead) { s.Scope.Target++ }},
		{"observation", func(s *interfaceIdentityRead) { s.Scope.Observation++ }},
		{"session", func(s *interfaceIdentityRead) { s.Scope.Session++ }},
		{"companion_field", func(s *interfaceIdentityRead) { s.Field = interfaceAPI }},
	} {
		t.Run(c.name, func(t *testing.T) {
			bad := rev
			c.change(&bad.Read)
			if matchingInterfaceIdentity(api, bad) {
				t.Fatalf("mismatched %s qualified a pair", c.name)
			}
		})
	}
	api.Read.Field = interfaceRevision
	if matchingInterfaceIdentity(api, rev) {
		t.Fatal("API field was not checked")
	}
}
