package ebus

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"testing"

	"github.com/Project-Helianthus/helianthus-ebusgo/protocol"
	"github.com/Project-Helianthus/helianthus-ebusreg/vaillant/nativepoints"
	semreg "github.com/Project-Helianthus/helianthus-semreg/semreg/v1"
	"github.com/Project-Helianthus/helianthus-semreg/semreg/v1/roweval"
)

type dhwCaller struct {
	Template    semreg.FactCandidate
	Roles       roweval.Roles
	SubjectKind string `json:"subject_kind"`
}

type comparisonRef struct {
	Case, Path, Origin string
	Digest             string `json:"sha256"`
}

type specRef struct {
	comparisonRef
	Index struct {
		Repository, Path string
		Digest           string `json:"sha256"`
		MergedBy         string `json:"merged_by"`
		Entry            struct {
			Case, Kind, Path, Origin, Recorder string
			Digest                             string `json:"sha256"`
		}
	}
}

type dhwVector struct {
	Name, Selector, Response, Direction, Presentation, Caller string
	GateResponse                                              *string `json:"gate_response"`
	GateError                                                 string  `json:"gate_error"`
	ResponseError                                             string  `json:"response_error"`
	Domain                                                    *struct{ Min, Max, Source, Scope string }
	Layout                                                    struct {
		Point, Unit  string
		Codec        nativepoints.Codec
		ScaleDivisor int `json:"scale_divisor"`
		Width        uint16
	}
	Expected struct {
		ValueTail     *string `json:"value_tail"`
		Gate          string
		ReadReason    string `json:"read_reason"`
		GateReason    string `json:"gate_reason"`
		Reason        roweval.Reason
		DecodeFailure nativepoints.DecodeFailure `json:"decode_failure"`
		InstanceEcho  *byte                      `json:"instance_echo"`
		NativeRow     *string                    `json:"native_row"`
		Native        *string
		NativeUnit    string `json:"native_unit"`
		Binding       *struct {
			NativeID                string `json:"native_id"`
			Direction, Presentation string
		}
		Key         *semreg.FactKey
		Subject     *roweval.Subject
		SubjectKind string `json:"subject_kind"`
		SI          *string
		SIUnit      semreg.DefinitionID `json:"si_unit"`
	}
	OldPath  *comparisonRef `json:"old_path"`
	Spec     *specRef       `json:"specification"`
	Verbatim *comparisonRef `json:"verbatim_reply_comparison"`
}

func TestDHWExemplar(t *testing.T) {
	f := fixtureJSON[struct {
		Callers map[string]dhwCaller
		Vectors []dhwVector
	}](t, "testdata/dhw_vectors.json")
	for _, v := range f.Vectors {
		t.Run(v.Name, func(t *testing.T) { checkDHWVector(t, v, f.Callers[v.Caller]) })
	}
	for _, pair := range []struct {
		code  uint16
		label string
	}{{0, "off"}, {1, "on"}, {2, ""}} {
		if got := circuitTypeLabel(pair.code); got != pair.label {
			t.Errorf("circuit type code %d: got %q, want %q", pair.code, got, pair.label)
		}
	}
}

func checkDHWVector(t *testing.T, v dhwVector, caller dhwCaller) {
	t.Helper()
	gate := acquireGate(t, v)
	read := acquireRegister(t, v.Selector, selector(t, v.Selector), v.Response, v.ResponseError)
	w := v.Expected
	var tail []byte
	if w.ValueTail != nil {
		tail = octets(t, *w.ValueTail)
	}
	if read.Reason != w.ReadReason || !reflect.DeepEqual(read.InstanceEcho, w.InstanceEcho) || !bytes.Equal(read.Value, tail) {
		t.Fatalf("correlation: got %+v (tail %x), want %s, instance %v, tail %x", read, read.Value, w.ReadReason, w.InstanceEcho, tail)
	}
	p, ok := nativepoints.Lookup(nativepoints.DeviceClassController, v.Layout.Point)
	if !ok || p.Codec != v.Layout.Codec || int(p.ScaleDivisor) != v.Layout.ScaleDivisor || p.NativeUnit != v.Layout.Unit {
		t.Fatalf("native catalog layout: %+v; want %+v", p, v.Layout)
	}
	a := acquiredDHW(t, v, caller, p)
	table := dhwTable(a.Template.FreshnessPolicy.PolicyID)
	for _, row := range table.Native {
		if row.ID == p.Name && (row.Wire.Width != v.Layout.Width || row.Wire.Codec != string(v.Layout.Codec)) {
			t.Fatalf("row layout %+v; want %+v", row.Wire, v.Layout)
		}
	}
	got, err := evaluateDHW(table, a, read, gate)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got.Read, read) || !reflect.DeepEqual(got.Gate, gate) || !bytes.Equal(got.Read.Raw, octets(t, v.Response)) {
		t.Fatal("raw evidence or request scope was lost")
	}
	if v.GateResponse != nil && !bytes.Equal(got.Gate.Raw, octets(t, *v.GateResponse)) {
		t.Fatal("raw gate reply was lost")
	}
	if got.Evaluation.Reason != w.Reason {
		t.Fatalf("reason %q (%s), want %q", got.Evaluation.Reason, got.Evaluation.Detail, w.Reason)
	}
	assertDHWNative(t, v, p, got, a.Template.Evidence)
	assertDHWBinding(t, v, caller, table, a, got)
	if v.OldPath != nil {
		compareDHWNative(t, *v.OldPath, got, a)
	}
	if v.Verbatim != nil {
		d := readComparison(t, *v.Verbatim)
		if !bytes.Equal(got.Read.Raw, d.Input.Measurement.Attempt.RawReply) || !bytes.Equal(got.Gate.Raw, d.Input.Presence.Attempt.RawReply) {
			t.Fatal("negative case did not retain both recorded replies verbatim")
		}
	}
}

func selector(t *testing.T, wire string) B524Selector {
	t.Helper()
	b := octets(t, wire)
	if len(b) != 8 {
		t.Fatalf("not a register request: %x", b)
	}
	return B524Selector{Primary: b[0], Secondary: b[1], Opcode: b[2], Operation: b[3], Group: b[4], Instance: b[5], Register: binary.LittleEndian.Uint16(b[6:])}
}

func acquireRegister(t *testing.T, request string, s B524Selector, raw, failure string) registerRead {
	t.Helper()
	frame := protocol.Frame{Source: fixtureSource, Target: fixtureTarget, Primary: s.Primary, Secondary: s.Secondary, Data: s.payload()}
	b, err := fixtureReply(t, frame, octets(t, request), octets(t, raw), failure)
	return correlateRegisterRead(s, scalarRead, b, err)
}

func acquireGate(t *testing.T, v dhwVector) registerRead {
	t.Helper()
	var gate registerRead
	if v.GateResponse != nil || v.GateError != "" {
		var raw string
		if v.GateResponse != nil {
			raw = *v.GateResponse
		}
		gate = acquireRegister(t, "b524020001000100", dhwGateSelector(), raw, v.GateError)
	}
	present, known := dhwFeatures(gate)[dhwFeature]
	state := "unknown"
	if known {
		state = "not_present"
	}
	if known && present {
		state = "present"
	}
	if state != v.Expected.Gate || gate.Reason != v.Expected.GateReason {
		t.Fatalf("gate %+v: got %s, want %s (%s)", gate, state, v.Expected.Gate, v.Expected.GateReason)
	}
	return gate
}

func acquiredDHW(t *testing.T, v dhwVector, c dhwCaller, p nativepoints.Point) roweval.Acquired[B524Selector] {
	t.Helper()
	q := map[string]string{}
	if v.Direction != "" {
		q["direction"] = v.Direction
	}
	if v.Presentation != "" {
		q["presentation"] = v.Presentation
	}
	if d := v.Domain; d != nil {
		if d.Source == "" || d.Scope == "" {
			t.Fatal("domain is not qualified by source and scope")
		}
		q["native_min"], q["native_max"] = d.Min, d.Max
	}
	return roweval.Acquired[B524Selector]{Qualifiers: roweval.Qualifiers{Values: q}, Template: c.Template, Roles: c.Roles, Descriptors: dhwDescriptors(p),
		Receipt: roweval.Receipt{ReceivedAt: c.Template.Times.ReceivedAt, Monotonic: c.Template.Times.ReceiptMonotonic}}
}

func assertDHWNative(t *testing.T, v dhwVector, p nativepoints.Point, got dhwResult, evidence []semreg.EvidenceRef) {
	t.Helper()
	w, n := v.Expected, got.Evaluation.Native
	if w.NativeRow == nil {
		if n != nil || got.DecodeError != nil {
			t.Fatalf("unmatched read acquired native output: %+v", n)
		}
		return
	}
	if n == nil || n.RowID != *w.NativeRow || n.Unit != w.NativeUnit || !reflect.DeepEqual(n.Evidence, evidence) {
		t.Fatalf("native row/unit/evidence = %+v; want %s, %s, %v", n, *w.NativeRow, w.NativeUnit, evidence)
	}
	if w.Native == nil {
		var failure *nativepoints.DecodeError
		if n.Value != nil || n.Reason != w.Reason || !errors.As(got.DecodeError, &failure) || failure.Failure != w.DecodeFailure {
			t.Fatalf("unknown native output %+v, decode error %v; expected %s", n, got.DecodeError, w.DecodeFailure)
		}
		return
	}
	if n.Reason != "" || got.DecodeError != nil {
		t.Fatalf("unexpected native failure: %+v, %v", n, got.DecodeError)
	}
	assertQuantity(t, n.Value, *w.Native, semreg.DefinitionID(w.NativeUnit))
	native, err := nativepoints.DecodeNative(p, got.Read.Value)
	if err != nil || native.Native != *w.Native || native.NativeUnit != w.NativeUnit {
		t.Fatalf("existing catalog decode: %+v, %v", native, err)
	}
}

func assertDHWBinding(t *testing.T, v dhwVector, c dhwCaller, table roweval.Table[B524Selector], a roweval.Acquired[B524Selector], got dhwResult) {
	t.Helper()
	w, e := v.Expected, got.Evaluation
	if w.Key == nil {
		if e.Candidate != nil || e.Subject != nil || w.SI != nil || w.Binding != nil {
			t.Fatalf("unexpected binding or number: %+v", e)
		}
		return
	}
	if e.Candidate == nil || !reflect.DeepEqual(e.Candidate.Key, *w.Key) || !reflect.DeepEqual(e.Subject, w.Subject) || c.SubjectKind != w.SubjectKind {
		t.Fatalf("binding returned candidate %+v and subject %+v; want key %+v and subject %+v", e.Candidate, e.Subject, w.Key, w.Subject)
	}
	assertQuantity(t, e.Candidate.Value, *w.SI, w.SIUnit)
	template := c.Template
	template.Key, template.Value = e.Candidate.Key, e.Candidate.Value
	if !reflect.DeepEqual(*e.Candidate, template) {
		t.Fatalf("metadata changed: got %+v, want %+v", *e.Candidate, template)
	}
	var others []roweval.BindingRow
	matches := 0
	for _, b := range table.Bindings {
		if b.NativeID == w.Binding.NativeID && b.When.Equal["direction"] == w.Binding.Direction && b.When.Equal["presentation"] == w.Binding.Presentation {
			if b.Read == nil || b.Write != nil {
				t.Fatalf("expected read-only binding: %+v", b)
			}
			matches++
		} else {
			others = append(others, b)
		}
	}
	table.Bindings = others
	without, err := evaluateDHW(table, a, got.Read, got.Gate)
	if matches != 1 || err != nil || without.Evaluation.Reason != roweval.ReasonNoBinding || without.Evaluation.Candidate != nil {
		t.Fatalf("expected exactly the fixture binding %+v to produce this candidate", w.Binding)
	}
}

func assertQuantity(t *testing.T, v *semreg.Value, text string, unit semreg.DefinitionID) {
	t.Helper()
	if v == nil || v.Kind != semreg.ValueQuantity || v.Quantity == nil || v.Quantity.Unit != unit {
		t.Fatalf("got %+v, want %s %s", v, text, unit)
	}
	want, ok := new(big.Rat).SetString(text)
	n := v.Quantity.Number
	got, parsed := new(big.Rat).SetString(fmt.Sprintf("%se%d", n.Coefficient, n.Exponent10))
	if !ok || !parsed || got.Cmp(want) != 0 {
		t.Fatalf("got %+v %s, want %s", n, unit, text)
	}
}

type measurementAttempt struct {
	RawReply, Value []byte
	Evidence        semreg.EvidenceRef
	ReceivedAt      semreg.TimePoint
	Receipt         semreg.MonotonicPoint
}

type dhwComparison struct {
	Input struct {
		Measurement, Presence struct{ Attempt measurementAttempt }
	}
	Output struct {
		Celsius     json.Number
		Target      struct{ CanonicalUnit semreg.DefinitionRef }
		Facts       []json.RawMessage
		Withheld    string
		Disposition struct {
			Outcome string
			Loss    []struct{ Description string }
		}
		Sources []struct {
			Observation struct {
				RawReply []byte
				Value    *semreg.Value
				Quality  semreg.Quality
				Times    semreg.Times
				Evidence []semreg.EvidenceRef
			}
		}
	}
}

func readComparison(t *testing.T, ref comparisonRef) dhwComparison {
	t.Helper()
	path := filepath.Join("testdata/qq01", filepath.Base(ref.Path))
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	digest := fmt.Sprintf("%x", sha256.Sum256(raw))
	if ref.Origin != "recorded-current" || digest != ref.Digest {
		t.Fatalf("record %s: unexpected origin or checksum %s", ref.Case, digest)
	}
	t.Logf("comparison %s (%s), sha256 %s", ref.Case, ref.Path, digest)
	var data dhwComparison
	if err := json.Unmarshal(raw, &data); err != nil {
		t.Fatal(err)
	}
	return data
}

func compareDHWNative(t *testing.T, ref comparisonRef, got dhwResult, a roweval.Acquired[B524Selector]) {
	t.Helper()
	record := readComparison(t, ref)
	d := record.Output
	n := got.Evaluation.Native
	if ref.Case == "vaillant-dhw-0005" {
		m := record.Input.Measurement.Attempt
		if d.Target.CanonicalUnit.ID != "unit.kelvin" || len(d.Facts) != 0 || d.Withheld != "final fact key and decimal policy unresolved" || !bytes.Equal(got.Read.Value, m.Value) {
			t.Fatalf("unexpected measured comparison fields: facts %d, hold %q, tail %x", len(d.Facts), d.Withheld, m.Value)
		}
		assertQuantity(t, n.Value, d.Celsius.String(), "°C")
		if a.Receipt.ReceivedAt != m.ReceivedAt || a.Receipt.Monotonic != m.Receipt || !reflect.DeepEqual(n.Evidence, []semreg.EvidenceRef{m.Evidence}) {
			t.Fatal("measured observation lost shared receipt or provenance")
		}
		return
	}
	if d.Disposition.Outcome != "withheld" || len(d.Disposition.Loss) != 1 || d.Disposition.Loss[0].Description != "hot-water resource identity is not allocated before graph publication" || len(d.Sources) != 1 {
		t.Fatalf("unexpected setpoint disposition: %+v", d.Disposition)
	}
	ob := d.Sources[0].Observation
	if !bytes.Equal(got.Read.Value, ob.RawReply) {
		t.Fatalf("setpoint tails differ: %x versus %x", got.Read.Value, ob.RawReply)
	}
	if n == nil || n.Value == nil || ob.Value == nil {
		return
	}
	if ob.Value.Quantity.Unit != "native.celsius" || n.Unit != "°C" || ob.Value.Quantity.Number != n.Value.Quantity.Number || !reflect.DeepEqual(ob.Evidence, n.Evidence) || !reflect.DeepEqual(ob.Times, a.Template.Times) {
		t.Fatalf("setpoint native value, time or evidence differs: new %+v, recorded %+v", n, ob)
	}
	if got.Evaluation.Candidate != nil {
		q := got.Evaluation.Candidate.Quality
		q.Promotion = ob.Quality.Promotion
		if !reflect.DeepEqual(q, ob.Quality) {
			t.Fatalf("shared quality differs: %+v versus %+v", q, ob.Quality)
		}
	}
}

// Local consistency check of the fixture against its own index entry; the
// index file itself is not staged, so the index digest is checked for form only.
func TestDHWSpecificationOperandIsTraceable(t *testing.T) {
	f := fixtureJSON[struct{ Vectors []dhwVector }](t, "testdata/dhw_vectors.json")
	found := 0
	for _, v := range f.Vectors {
		if v.Spec == nil || v.Spec.Case == "" {
			continue
		}
		found++
		s, e := v.Spec, v.Spec.Index.Entry
		if s.Origin != "specification" || e.Origin != s.Origin || e.Case != s.Case || e.Path != s.Path || e.Digest != s.Digest {
			t.Fatalf("%s: operand differs from its index entry: %+v", v.Name, s)
		}
		if s.Index.Repository == "" || s.Index.Path != "acceptance/vectors/index.json" || !regexp.MustCompile(`^[0-9a-f]{64}$`).MatchString(s.Index.Digest) || !regexp.MustCompile(`^[0-9a-f]{64}$`).MatchString(s.Digest) || s.Index.MergedBy == "" {
			t.Fatalf("%s: index entry is not immutable and traceable: %+v", v.Name, s.Index)
		}
		if v.Name != "setpoint_50" || *v.Expected.SI != "323.15" {
			t.Fatalf("%s: unexpected specification operand", v.Name)
		}
	}
	if found != 1 {
		t.Fatalf("%d specification operands, want the single 50.0 C operand", found)
	}
}
