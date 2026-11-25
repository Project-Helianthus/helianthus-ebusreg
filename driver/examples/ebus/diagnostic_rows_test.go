package ebus

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"sync"
	"testing"

	"github.com/Project-Helianthus/helianthus-ebusgo/protocol"
	"github.com/Project-Helianthus/helianthus-ebusgo/protocol/vaillant/b503"
)

func checkDiagnosticResponses(t *testing.T, vectors map[string]nativeResponseVector) {
	t.Helper()
	for _, c := range []struct {
		name, request string
		target        byte
	}{
		{"b503_current_error", "0001", 8}, {"b503_current_service", "0002", 8}, {"b503_error_history", "0101", 8},
		{"b503_service_history", "0102", 8}, {"b503_live_monitor", "0003", 8}, {"b503_ack_decode_failed", "0001", 8},
		{"b503_duplicate_envelope", "b5030001", 8}, {"b503_pre_cancelled", "0001", 8},
		{"b503_nak", "0003", 0x15}, {"b503_ambiguous", "0003", 0x15}, {"b503_stale_epoch", "0001", 8},
	} {
		nativeCase(t, vectors, c.name, func(t *testing.T, v nativeResponseVector) {
			checkDiagnosticCase(t, c.name, c.request, c.target, v)
		})
	}
}

func checkDiagnosticCase(t *testing.T, name, request string, target byte, v nativeResponseVector) {
	t.Helper()
	if v.Request != request || v.Target != target || v.Transport != (diagnosticStamp{AdapterInstanceID: "test", TransportEpoch: 1}) {
		t.Fatalf("request scope or issuing transport changed: %+v", v)
	}
	s, shapeErr := diagnosticBody(octets(t, request))
	if (shapeErr != nil) != (name == "b503_duplicate_envelope") {
		t.Fatalf("body shape: %v", shapeErr)
	}
	c := diagnosticCompletion{Issued: v.Transport, Current: v.CompletionTransport}
	if v.Error != "" {
		c.ClassifiedError = errors.New(v.Error)
	}
	if shapeErr == nil && name != "b503_pre_cancelled" {
		c.Emitted = true
		if v.Frame == nil {
			t.Fatal("record missing emitted frame")
		}
		f, err := s.request(0x7f, target)
		if err != nil {
			t.Fatal(err)
		}
		c.Reply, c.SendError = replayDiagnosticExchange(t, name, v, f)
	} else if v.Frame != nil || v.WireReply != "" || v.Response != nil {
		t.Fatal("no-send record has wire fields")
	}
	x := c.observation()
	if x.Emitted != v.Emitted || x.Transport != v.Transport || x.Native != v.Native || errText(x.ClassifiedError) != v.Error {
		t.Fatalf("native dispatch result %+v, recorded %+v", x, v)
	}
	if (x.Response != nil) != (v.Response != nil) || (v.Response != nil && !bytes.Equal(x.Response, octets(t, *v.Response))) {
		t.Fatalf("unexpected admitted response %x", x.Response)
	}
	if x.Native == diagnosticACK {
		d, err := s.decode(x)
		if errText(err) != v.DecodeError {
			t.Fatalf("decoder error %v, recorded %q", err, v.DecodeError)
		}
		if err == nil {
			compareDecoded(t, d, v.Decoded)
		}
	}
	if name == "b503_stale_epoch" {
		checkDiagnosticFence(t, name, v)
	}
}

func replayDiagnosticExchange(t *testing.T, name string, v nativeResponseVector, f protocol.Frame) (*protocol.Frame, error) {
	t.Helper()
	rec := *v.Frame
	rec.Reply = &v.WireReply
	if name == "b503_nak" || name == "b503_ambiguous" {
		if v.WireReply != "" {
			t.Fatal("failed Send unexpectedly returned bytes")
		}
		failure := "nack"
		if name == "b503_ambiguous" {
			failure = "timeout"
		}
		r, err := simulateExchange(t, f, rec.Source, rec.Target, nativeBody(t, rec), nil, failure)
		if name == "b503_ambiguous" {
			err = errors.New(v.SendError)
		}
		return r, err
	}
	return onSimulatedBus(t, f, rec, true)
}

// liveStamp is the transport stamp the caller reads once an exchange has completed.
type liveStamp struct {
	mu sync.Mutex
	s  diagnosticStamp
}

func (l *liveStamp) set(s diagnosticStamp) { l.mu.Lock(); l.s = s; l.mu.Unlock() }
func (l *liveStamp) get() diagnosticStamp  { l.mu.Lock(); defer l.mu.Unlock(); return l.s }

// checkDiagnosticFence issues the recorded positive read three times. Each run holds the
// reply until the send is in flight, then moves the live transport (or leaves it), and
// only afterwards reads the live stamp; no stamp is edited after the exchange completes.
func checkDiagnosticFence(t *testing.T, name string, v nativeResponseVector) {
	t.Helper()
	s, err := diagnosticBody(octets(t, v.Request))
	if err != nil {
		t.Fatal(err)
	}
	f, err := s.request(0x7f, v.Target)
	if err != nil {
		t.Fatal(err)
	}
	rec := *v.Frame
	rec.Reply = &v.WireReply
	if v.CompletionTransport.TransportEpoch == v.Transport.TransportEpoch {
		t.Fatal("record does not advance the epoch")
	}
	for _, run := range []struct {
		label    string
		advance  func(*liveStamp)
		admitted bool
	}{
		{"epoch advanced in flight", func(l *liveStamp) { l.set(v.CompletionTransport) }, false},
		{"adapter instance replaced in flight", func(l *liveStamp) {
			l.set(diagnosticStamp{AdapterInstanceID: "different-adapter", TransportEpoch: v.Transport.TransportEpoch})
		}, false},
		{"transport unchanged", func(*liveStamp) {}, true},
	} {
		live := &liveStamp{s: v.Transport}
		var started bool
		reply, sendErr := simulateExchangeInFlight(t, f, rec.Source, rec.Target, nativeBody(t, rec), octets(t, v.WireReply), "", func() {
			started = true
			run.advance(live)
		})
		if !started || reply == nil || sendErr != nil {
			t.Fatalf("%s: exchange did not complete a positive reply in flight", run.label)
		}
		c := diagnosticCompletion{Issued: v.Transport, Current: live.get(), Emitted: true, Reply: reply}
		r := c.observation()
		if run.admitted {
			if r.Native != diagnosticACK || !bytes.Equal(r.Response, reply.Data) {
				t.Fatalf("%s: same-stamp reply not admitted", run.label)
			}
			continue
		}
		if r.Native != diagnosticAmbiguous || r.Response != nil || r.Transport != v.Transport {
			t.Fatalf("%s: obsolete reply admitted or lost its issuing transport: %+v", run.label, r)
		}
	}
}

type diagnosticCallerRecord struct {
	Target                byte
	ProbeTarget           byte `json:"probe_target"`
	Request, Reply, Probe string
	Action                liveMonitorAction
	Data, Error           json.RawMessage
	Meta                  struct {
		Capabilities map[string]struct {
			Available bool
			Reason    string
		}
		Consistency struct{ Mode string }
	}
}

func checkDiagnosticCallers(t *testing.T, source json.RawMessage) {
	t.Helper()
	var records map[string]diagnosticCallerRecord
	if err := json.Unmarshal(source, &records); err != nil {
		t.Fatal(err)
	}
	if len(records) != 7 {
		t.Fatalf("recorder produced %d caller cases, want 7", len(records))
	}
	for _, c := range []struct {
		name     string
		selector byte
		index    *byte
	}{
		{"error_history_index_absent", 1, nil}, {"error_history_index_zero", 1, testPointer(byte(0))},
		{"service_history_index_absent", 2, nil}, {"service_history_index_zero", 2, testPointer(byte(0))},
	} {
		t.Run(c.name, func(t *testing.T) {
			checkHistoryCaller(t, records[c.name], diagnosticSelector{Family: 1, Selector: c.selector, Index: c.index})
		})
	}
	for _, action := range []liveMonitorAction{liveMonitorEnable, liveMonitorRead, liveMonitorDisable} {
		t.Run("live_monitor_"+string(action), func(t *testing.T) { checkLiveCaller(t, records["live_monitor_"+string(action)], action) })
	}
}

func checkHistoryCaller(t *testing.T, v diagnosticCallerRecord, s diagnosticSelector) {
	t.Helper()
	f, err := s.request(0x7f, 8)
	if err != nil {
		t.Fatal(err)
	}
	reply := callerExchange(t, f, v, true)
	x, err := s.decode(diagnosticValue{Native: diagnosticACK, Response: reply})
	if err != nil {
		t.Fatal(err)
	}
	history := x.(b503.ErrorHistoryRecord)
	code, found := history.FirstActive()
	if !found || v.Probe != "0001" || v.ProbeTarget != 8 {
		t.Fatalf("history or following probe %+v, %+v", history, v)
	}
	compareDecoded(t, map[string]any{"index": history.Index, "first_active_error": code, "slots": []any{code, nil, nil, nil, nil}}, v.Data)
	probe, err := (diagnosticSelector{Family: 0, Selector: 1}).request(0x7f, 8)
	if err != nil {
		t.Fatal(err)
	}
	current := nativeFrame{Source: 0x7f, Target: v.ProbeTarget, Primary: 0xb5, Secondary: 3, Request: v.Probe, Reply: testPointer("1901ffffffffffffffff")}
	_, _ = onSimulatedBus(t, probe, current, true)
}

func checkLiveCaller(t *testing.T, v diagnosticCallerRecord, action liveMonitorAction) {
	t.Helper()
	r := liveMonitorRequest{Target: 8, Action: action}
	f, err := r.request(0x7f)
	if err != nil {
		t.Fatal(err)
	}
	if v.Action != r.Action || v.Probe != "" || v.ProbeTarget != 0 {
		t.Fatalf("managed call did not keep action separate: %+v", v)
	}
	reply := callerExchange(t, f, v, false)
	live, err := b503.DecodeLiveMonitorMain(reply)
	if err != nil || live.Status != 1 || live.Function != 0 {
		t.Fatalf("native live monitor result %+v, %v", live, err)
	}
	expected := map[string]any{}
	if action == liveMonitorRead {
		expected["raw_hex"] = hex.EncodeToString(reply)
	}
	if action == liveMonitorDisable {
		expected["disabled"] = true
	}
	compareDecoded(t, expected, v.Data)
}

func callerExchange(t *testing.T, f protocol.Frame, v diagnosticCallerRecord, available bool) []byte {
	t.Helper()
	capability := v.Meta.Capabilities["vaillant_b503"]
	wantReason := "UNKNOWN"
	if available {
		wantReason = "AVAILABLE"
	}
	if v.Target != 8 || string(v.Error) != "null" || capability.Available != available || capability.Reason != wantReason || v.Meta.Consistency.Mode != "LIVE" {
		t.Fatalf("caller envelope state differs: %+v", v)
	}
	rec := nativeFrame{Source: 0x7f, Target: v.Target, Primary: 0xb5, Secondary: 3, Request: v.Request, Reply: &v.Reply}
	reply, _ := onSimulatedBus(t, f, rec, true)
	return reply.Data
}
