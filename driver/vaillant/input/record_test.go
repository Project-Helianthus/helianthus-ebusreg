package input

import (
	"reflect"
	"strings"
	"testing"
	"time"
)

func fieldNames(v any) []string {
	typ := reflect.TypeOf(v)
	names := make([]string, typ.NumField())
	for i := range names {
		names[i] = typ.Field(i).Name
	}
	return names
}

func TestVaillantInputRecordCarriesTheAcquiredObservationAndNoGeneration(t *testing.T) {
	want := []string{"Selector", "Request", "Reply", "Replied", "Source", "Target", "Qualifiers", "Clock", "Diagnostics"}
	if got := fieldNames(Record{}); !reflect.DeepEqual(got, want) {
		t.Fatalf("record fields %v, want %v", got, want)
	}
}

func TestVaillantInputRecordHoldsEveryFieldItIsGiven(t *testing.T) {
	one, two, three, four := uint16(1), uint16(2), uint16(3), uint16(4)
	clock := time.Date(2030, 1, 2, 3, 4, 5, 6, time.UTC)
	sel := B524Selector{Opcode: 0x02, Group: 0x03, Instance: 0x01, Register: 0x000e}
	record := Record{
		Selector: sel, Request: sel.Payload(), Reply: []byte{0x01, 0x03, 0x0e, 0x00, 0x02}, Replied: true,
		Source: 0x31, Target: 0x15,
		Qualifiers: Qualifiers{
			Zone:       ZoneInputs{OperatingMode: &one, SpecialFunction: &two, CircuitType: &three, HasCircuitType: true},
			CachedZone: ZoneInputs{OperatingMode: &four, SpecialFunction: &one, CircuitType: &two, HasCircuitType: false},
			DHW:        DHWInputs{OperatingMode: &three, SpecialFunction: &four},
			CachedDHW:  DHWInputs{OperatingMode: &two, SpecialFunction: &three},
		},
		Clock:       clock,
		Diagnostics: Diagnostics{Outcome: OutcomeNoValue, SendError: "send failed", MismatchText: "mismatch"},
	}
	if record.Selector.(B524Selector) != sel || string(record.Request) != string(sel.Payload()) {
		t.Errorf("selector or request: %+v %x", record.Selector, record.Request)
	}
	if record.Source != 0x31 || record.Target != 0x15 || !record.Clock.Equal(clock) || !record.Replied {
		t.Errorf("addresses, clock or replied: %+v", record)
	}
	q := record.Qualifiers
	if *q.Zone.OperatingMode != 1 || *q.Zone.SpecialFunction != 2 || *q.Zone.CircuitType != 3 || !q.Zone.HasCircuitType {
		t.Errorf("current zone inputs: %+v", q.Zone)
	}
	if *q.CachedZone.OperatingMode != 4 || *q.CachedZone.SpecialFunction != 1 || *q.CachedZone.CircuitType != 2 || q.CachedZone.HasCircuitType {
		t.Errorf("cached zone inputs: %+v", q.CachedZone)
	}
	if *q.DHW.OperatingMode != 3 || *q.DHW.SpecialFunction != 4 || *q.CachedDHW.OperatingMode != 2 || *q.CachedDHW.SpecialFunction != 3 {
		t.Errorf("hot-water inputs: %+v %+v", q.DHW, q.CachedDHW)
	}
	if record.Diagnostics.Outcome != OutcomeNoValue || record.Diagnostics.SendError != "send failed" || record.Diagnostics.MismatchText != "mismatch" {
		t.Errorf("diagnostics: %+v", record.Diagnostics)
	}
}

func TestVaillantInputRecordQualifiersCarryExactlyWhatTheHooksRead(t *testing.T) {
	if got, want := fieldNames(ZoneInputs{}), []string{"OperatingMode", "SpecialFunction", "CircuitType", "HasCircuitType"}; !reflect.DeepEqual(got, want) {
		t.Errorf("zone inputs %v, want %v", got, want)
	}
	if got, want := fieldNames(DHWInputs{}), []string{"OperatingMode", "SpecialFunction"}; !reflect.DeepEqual(got, want) {
		t.Errorf("hot-water inputs %v, want %v", got, want)
	}
	if got, want := fieldNames(Qualifiers{}), []string{"Zone", "CachedZone", "DHW", "CachedDHW"}; !reflect.DeepEqual(got, want) {
		t.Errorf("qualifiers %v, want %v", got, want)
	}
}

func TestVaillantInputRecordSeparatesNoFrameFromEmptyReply(t *testing.T) {
	noFrame := Record{Selector: B509Selector{Opcode: BoilerRead, Register: 0x4400}}
	emptyReply := noFrame
	emptyReply.Replied = true
	emptyReply.Reply = []byte{}
	if reflect.DeepEqual(noFrame, emptyReply) || noFrame.Replied || len(emptyReply.Reply) != 0 || !emptyReply.Replied {
		t.Fatalf("no frame %+v and empty reply %+v are not distinct", noFrame, emptyReply)
	}
}

func TestVaillantInputRecordCurrentAndCachedInputsAreIndependent(t *testing.T) {
	current, cached := uint16(1), uint16(2)
	q := Qualifiers{Zone: ZoneInputs{SpecialFunction: &current}, CachedZone: ZoneInputs{SpecialFunction: &cached}}
	if q.Zone.SpecialFunction == q.CachedZone.SpecialFunction || *q.Zone.SpecialFunction == *q.CachedZone.SpecialFunction {
		t.Fatalf("current and cached special function share a value: %+v", q)
	}
	var none ZoneInputs
	if none.OperatingMode != nil || none.SpecialFunction != nil || none.CircuitType != nil || none.HasCircuitType {
		t.Fatalf("absent zone inputs are not absent: %+v", none)
	}
}

func TestVaillantInputDiagnosticsFormatTheTextsThePollerUses(t *testing.T) {
	cases := []struct{ name, got, want string }{
		{"b509 nil", BoilerReadNilResponse(), "b509 read returned nil response"},
		{"b509 failed", BoilerReadFailed(0x08, 0x4400), "b509 read failed: target=0x08 addr=0x4400"},
		{"b524 nil", ExtendedReadNilResponse(), "b524 read returned nil response"},
		{"b524 failed", ExtendedReadFailed(0x02, 0x03, 0x01, 0x000e), "b524 read failed: opcode=0x02 group=0x03 instance=0x01 addr=0x000e"},
		{
			"b524 instance mismatch", ExtendedInstanceMismatch(0x02, 0x03, 0x01, 0x000e, 0x05, 0x03, 0x000e),
			"b524 read mismatch: want opcode=0x02 group=0x03 instance=0x01 addr=0x000e; got reply-instance=0x05 (group=0x03 addr=0x000e)",
		},
		{
			"b524 header mismatch", ExtendedHeaderMismatch(0x02, 0x03, 0x01, 0x000e, 0x01, 0x04, 0x000f, 6),
			"b524 read mismatch: want opcode=0x02 group=0x03 instance=0x01 addr=0x000e; got kind=0x01 group=0x04 addr=0x000f len=6",
		},
	}
	for _, c := range cases {
		if c.got != c.want {
			t.Errorf("%s: %q, want %q", c.name, c.got, c.want)
		}
	}
}

func TestVaillantInputDiagnosticsAnchorEveryTextWithItsSurface(t *testing.T) {
	wantFile := map[string]string{
		"b509": "cmd/gateway/semantic_vaillant_b509_access.go", "b524": "cmd/gateway/semantic_vaillant_b524_access.go",
	}
	hidden := map[string]bool{
		boilerReadNilResponseText: true, boilerReadFailedText: true, boilerReadFailedFormat: true,
		extendedReadNilResponseText: true, extendedReadFailedText: true, extendedReadFailedFormat: true,
	}
	logged := map[string]bool{extendedInstanceMismatchText: true, extendedHeaderMismatchText: true}
	inv := loadInventory(t)
	files := map[string]bool{}
	for _, f := range inv.SourceFiles {
		files[f.Path] = true
	}
	seen := map[string]bool{}
	for _, a := range FailureAnchors {
		if seen[a.Text] {
			t.Errorf("duplicate anchor %q", a.Text)
		}
		seen[a.Text] = true
		if a.File != wantFile[a.Text[:4]] || !files[a.File] || a.Function == "" || a.Line <= 0 {
			t.Errorf("anchor %+v", a)
		}
		switch {
		case hidden[a.Text] && a.Surface == SurfaceSchedulerCallback:
		case logged[a.Text] && a.Surface == SurfaceLog:
		default:
			t.Errorf("anchor %+v has the wrong surface", a)
		}
	}
	if len(seen) != len(hidden)+len(logged) {
		t.Errorf("anchored texts %d, want %d", len(seen), len(hidden)+len(logged))
	}
}

func TestVaillantInputDiagnosticsAnchorsDifferByMessageNotBySource(t *testing.T) {
	prefixes := map[string]int{}
	for _, a := range FailureAnchors {
		prefixes[strings.SplitN(a.Text, ":", 2)[0]]++
	}
	if prefixes["b509 read failed"] != 2 || prefixes["b524 read failed"] != 2 || prefixes["b524 read mismatch"] != 2 {
		t.Fatalf("anchors per message: %v", prefixes)
	}
}

func TestVaillantInputDiagnosticsZeroValueIsNeitherAValueNorNoValue(t *testing.T) {
	var d Diagnostics
	if d.Outcome == OutcomeValue || d.Outcome == OutcomeNoValue || OutcomeValue == OutcomeNoValue {
		t.Fatalf("zero diagnostics %+v", d)
	}
	if d.SendError != "" || d.MismatchText != "" {
		t.Fatalf("zero diagnostics carry text: %+v", d)
	}
}
