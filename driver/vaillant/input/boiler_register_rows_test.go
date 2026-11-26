package input

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/Project-Helianthus/helianthus-ebusreg/catalog/vaillant"
)

var boilerRowsBase = time.Date(2001, 2, 3, 4, 5, 6, 0, time.UTC)

type boilerRowsAttempt struct {
	Source, Target, Primary, Secondary byte
	RequestHex                         string  `json:"request_hex"`
	AtMs                               int64   `json:"at_ms"`
	ReplyHex                           *string `json:"reply_hex"`
	SendError                          string  `json:"send_error"`
	ParseOK                            *bool   `json:"donor_parse_ok"`
	TailHex                            string  `json:"donor_tail_hex"`
	AckOK                              *bool   `json:"donor_ack_ok"`
}

type boilerRowsReadVector struct {
	Case      string              `json:"case"`
	Class     string              `json:"class"`
	Row       string              `json:"row"`
	Attempts  []boilerRowsAttempt `json:"attempts"`
	FinalOK   bool                `json:"final_ok"`
	FinalKind string              `json:"final_kind"`
	Final     any                 `json:"final_value"`
}

type boilerRowsWriteVector struct {
	Case       string              `json:"case"`
	Class      string              `json:"class"`
	Row        string              `json:"row"`
	ValueHex   string              `json:"value_hex"`
	Attempts   []boilerRowsAttempt `json:"attempts"`
	FinalError string              `json:"final_error"`
}

type boilerRowsFile struct {
	Format       string            `json:"format"`
	DonorSources map[string]string `json:"donor_sources"`
	Reads        []boilerRowsReadVector
	Writes       []boilerRowsWriteVector
}

type boilerRowsInventory struct {
	SourceFiles []struct {
		Path   string
		Digest string `json:"sha256"`
	} `json:"source_files"`
	Selectors []struct {
		ID, Service, Operation string
		GatewayConstant        string `json:"gateway_constant"`
		Register               string
		RequestHex             string `json:"request_hex"`
		NativeViewField        string `json:"native_view_field"`
	}
}

func boilerRowsLoad[T any](t *testing.T, path string) T {
	t.Helper()
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var v T
	if err := json.Unmarshal(raw, &v); err != nil {
		t.Fatal(err)
	}
	return v
}

func boilerRowsBytes(t *testing.T, h string) []byte {
	t.Helper()
	b, err := hex.DecodeString(h)
	if err != nil {
		t.Fatal(err)
	}
	return b
}

func boilerRowsRecord(t *testing.T, sel B509Selector, a boilerRowsAttempt) Record {
	t.Helper()
	rec := Record{
		Selector: sel, Request: boilerRowsBytes(t, a.RequestHex), Source: a.Source, Target: a.Target,
		Clock: boilerRowsBase.Add(time.Duration(a.AtMs) * time.Millisecond),
	}
	if a.SendError != "" {
		rec.Diagnostics = Diagnostics{Outcome: OutcomeNoValue, SendError: a.SendError}
	}
	if a.ReplyHex != nil {
		rec.Reply, rec.Replied = boilerRowsBytes(t, *a.ReplyHex), true
	}
	return rec
}

func TestVaillantInputRowsBoilerRegisterInventory(t *testing.T) {
	inv := boilerRowsLoad[boilerRowsInventory](t, "testdata/gateway_selectors.json")
	reads, writes := map[string]int{}, map[string]int{}
	for i, r := range BoilerRegisterRows {
		reads[r.ID] = i
	}
	for i, r := range BoilerWriteRows {
		writes[r.ID] = i
	}
	seenReads, seenWrites := 0, 0
	for _, s := range inv.Selectors {
		if s.Service != "B509" {
			continue
		}
		if s.Operation == "read" {
			i, ok := reads[s.ID]
			if !ok {
				t.Fatalf("inventoried selector %s has no row", s.ID)
			}
			seenReads++
			r := BoilerRegisterRows[i]
			if r.Name != s.GatewayConstant || "0d"+fmt.Sprintf("%04x", r.Register) != strings.ToLower(s.RequestHex) ||
				hex.EncodeToString(r.Selector().Payload()) != strings.ToLower(s.RequestHex) ||
				!strings.HasPrefix(s.NativeViewField, r.NativeField) {
				t.Fatalf("row %s differs from its inventory entry %+v", r.ID, s)
			}
			continue
		}
		i, ok := writes[s.ID]
		if !ok {
			t.Fatalf("inventoried write %s has no row", s.ID)
		}
		seenWrites++
		w := BoilerWriteRows[i]
		if w.Name != s.GatewayConstant || hex.EncodeToString(w.Selector().Payload()) != strings.ToLower(s.RequestHex) {
			t.Fatalf("write row %s differs from its inventory entry %+v", w.ID, s)
		}
	}
	if seenReads != len(BoilerRegisterRows) || seenWrites != len(BoilerWriteRows) || seenReads != 36 || seenWrites != 7 {
		t.Fatalf("rows %d/%d, inventoried %d/%d: a row has no inventoried selector", len(BoilerRegisterRows), len(BoilerWriteRows), seenReads, seenWrites)
	}
}

// boilerRowsCatalogCodec maps the codec declarations of the native catalog to the codec of a row; an empty result
// means the declaration does not state the conversion the poller applies.
func boilerRowsCatalogCodec(declaration string) BoilerCodec {
	switch declaration {
	case "DATA2c":
		return BoilerCodecSignedFixed
	case "DATA2b":
		return BoilerCodecUnsignedFixed
	case "UCH":
		return BoilerCodecUCH
	case "UIN":
		return BoilerCodecUIN
	case "SIN / 10":
		return BoilerCodecSIN
	case "UIN / 100":
		return BoilerCodecHundredths
	case "percent0":
		return BoilerCodecPercent0
	}
	return ""
}

func TestVaillantInputRowsBoilerRegisterCatalogReference(t *testing.T) {
	entries, err := vaillant.B509BoilerCodecDictionaryEntries()
	if err != nil {
		t.Fatal(err)
	}
	byRegister := map[uint16]vaillant.B509BoilerCodecDictionaryEntry{}
	for _, e := range entries {
		sel := e.B509BoilerCodecDictionaryKey.B509BoilerCodecDictionarySelector
		if len(sel) == 2 {
			byRegister[uint16(sel[0])<<8|uint16(sel[1])] = e
		}
	}
	referenced, fromCatalog, fromGateway := 0, 0, 0
	for _, r := range BoilerRegisterRows {
		e, found := byRegister[r.Register]
		if (r.Catalog != "") != found || (found && r.Catalog != e.B509BoilerCodecDictionaryQuantity) {
			t.Fatalf("row %s catalog reference %q, catalog entry present %v", r.ID, r.Catalog, found)
		}
		if found {
			referenced++
		}
		same := found && boilerRowsCatalogCodec(e.B509BoilerCodecDeclaration) == r.Codec
		switch {
		case r.CodecSource == CodecSourceCatalog && same:
			fromCatalog++
		case r.CodecSource == CodecSourceGatewayDecoder && !same:
			fromGateway++
		default:
			t.Fatalf("row %s codec source %s with catalog declaration match %v", r.ID, r.CodecSource, same)
		}
	}
	if referenced != len(byRegister) {
		t.Fatalf("%d catalog entries, %d referenced by a row", len(byRegister), referenced)
	}
	t.Logf("rows %d, referencing the catalog %d, codec from the catalog %d, codec from the gateway decoder %d", len(BoilerRegisterRows), referenced, fromCatalog, fromGateway)
}

func boilerRowsRows(t *testing.T) (map[string]BoilerRegisterRow, map[uint16]BoilerRegisterRow) {
	t.Helper()
	byID, byRegister := map[string]BoilerRegisterRow{}, map[uint16]BoilerRegisterRow{}
	for _, r := range BoilerRegisterRows {
		byID[r.ID], byRegister[r.Register] = r, r
	}
	return byID, byRegister
}

func boilerRowsSameValue(t *testing.T, kind string, got any, want any) {
	t.Helper()
	switch kind {
	case "float":
		if g, ok := got.(float64); !ok || g != want.(float64) {
			t.Fatalf("value %#v, donor %#v", got, want)
		}
	case "int":
		if g, ok := got.(int); !ok || float64(g) != want.(float64) {
			t.Fatalf("value %#v, donor %#v", got, want)
		}
	case "bool", "text", "raw":
		if got != want {
			t.Fatalf("value %#v, donor %#v", got, want)
		}
	default:
		t.Fatalf("unknown kind %q", kind)
	}
}

func TestVaillantInputRowsBoilerRegisterReadVectors(t *testing.T) {
	f := boilerRowsLoad[boilerRowsFile](t, "testdata/vectors/boiler_register.json")
	inv := boilerRowsLoad[boilerRowsInventory](t, "testdata/gateway_selectors.json")
	for _, s := range inv.SourceFiles {
		if want, ok := f.DonorSources[s.Path]; ok && want != s.Digest {
			t.Fatalf("vector donor source %s differs from the inventoried digest", s.Path)
		}
	}
	if len(f.DonorSources) != 5 || len(f.Reads) != 36*11+5 {
		t.Fatalf("vector file has %d donor sources and %d read cases", len(f.DonorSources), len(f.Reads))
	}
	byID, byRegister := boilerRowsRows(t)
	classes := map[string]map[string]bool{}
	for _, v := range f.Reads {
		if classes[v.Row] == nil {
			classes[v.Row] = map[string]bool{}
		}
		classes[v.Row][v.Class] = true
		t.Run(v.Case, func(t *testing.T) { boilerRowsCheckRead(t, v, byID, byRegister) })
	}
	for id := range byID {
		for _, c := range []string{"value_opcode_echo", "value_compact_echo", "value_bare", "neighbour_register", "sentinel_all_ff", "boundary_zero", "empty_reply", "nil_response", "send_error", "header_only_compact"} {
			if !classes[id][c] {
				t.Fatalf("row %s has no %s case", id, c)
			}
		}
	}
}

func boilerRowsCheckRead(t *testing.T, v boilerRowsReadVector, byID map[string]BoilerRegisterRow, byRegister map[uint16]BoilerRegisterRow) {
	t.Helper()
	if _, ok := byID[v.Row]; !ok || len(v.Attempts) == 0 {
		t.Fatalf("case row %q, %d attempts", v.Row, len(v.Attempts))
	}
	produced := -1
	var value BoilerReadResult
	for i, a := range v.Attempts {
		request := boilerRowsBytes(t, a.RequestHex)
		row, ok := byRegister[uint16(request[1])<<8|uint16(request[2])]
		if !ok || hex.EncodeToString(row.Selector().Payload()) != a.RequestHex {
			t.Fatalf("attempt %d request %s is not the selector of a row", i, a.RequestHex)
		}
		rec := boilerRowsRecord(t, row.Selector(), a)
		got, err := row.Read(rec)
		if err != nil {
			t.Fatal(err)
		}
		boilerRowsCheckAttempt(t, got, a, rec)
		if got.Outcome == OutcomeValue && produced < 0 {
			produced, value = i, got
		}
	}
	if (produced >= 0) != v.FinalOK {
		t.Fatalf("rows produced a value at attempt %d, donor result ok=%v", produced, v.FinalOK)
	}
	if !v.FinalOK {
		return
	}
	if string(value.Kind) != v.FinalKind {
		t.Fatalf("kind %s, donor %s", value.Kind, v.FinalKind)
	}
	boilerRowsSameValue(t, v.FinalKind, value.Value, v.Final)
	if produced != len(v.Attempts)-1 {
		t.Fatalf("value produced at attempt %d of %d; the donor stops at its first value", produced, len(v.Attempts))
	}
}

func boilerRowsCheckAttempt(t *testing.T, got BoilerReadResult, a boilerRowsAttempt, rec Record) {
	t.Helper()
	wantReason := BoilerReasonNone
	switch {
	case a.SendError != "":
		wantReason = BoilerReasonSendError
	case a.ReplyHex == nil:
		wantReason = BoilerReasonNoReply
	case !*a.ParseOK:
		wantReason = BoilerReasonPayloadRejected
	default:
		if hex.EncodeToString(got.Tail) != a.TailHex {
			t.Fatalf("tail %x, donor parser %s", got.Tail, a.TailHex)
		}
		if got.Outcome == OutcomeNoValue {
			wantReason = BoilerReasonDecodeRejected
		}
	}
	if got.Reason != wantReason || (got.Outcome == OutcomeValue) != (wantReason == BoilerReasonNone) {
		t.Fatalf("outcome %s reason %q, want reason %q", got.Outcome, got.Reason, wantReason)
	}
	if got.Value != nil && wantReason != BoilerReasonNone || (got.Value == nil) == (wantReason == BoilerReasonNone) {
		t.Fatalf("value %#v with reason %q", got.Value, wantReason)
	}
	if !got.AcquiredAt.Equal(rec.Clock) || got.Source != a.Source || got.Target != a.Target || hex.EncodeToString(got.Request) != a.RequestHex {
		t.Fatalf("acquisition fields changed: %+v", got)
	}
	if a.ReplyHex != nil && hex.EncodeToString(got.Raw) != *a.ReplyHex {
		t.Fatalf("raw reply %x, received %s", got.Raw, *a.ReplyHex)
	}
}

func TestVaillantInputRowsBoilerRegisterSpecifiedCalls(t *testing.T) {
	f := boilerRowsLoad[boilerRowsFile](t, "testdata/vectors/boiler_register.json")
	byID, _ := boilerRowsRows(t)
	cases := map[string]boilerRowsReadVector{}
	for _, v := range f.Reads {
		cases[v.Case] = v
	}
	read := func(name string, attempt int) BoilerReadResult {
		t.Helper()
		v, ok := cases[name]
		if !ok {
			t.Fatalf("missing case %s", name)
		}
		row := byID[v.Row]
		got, err := row.Read(boilerRowsRecord(t, row.Selector(), v.Attempts[attempt]))
		if err != nil {
			t.Fatal(err)
		}
		return got
	}
	if g := read("boiler_b509_central_heating_pump.tail_00_true", 0); g.Value != true || hex.EncodeToString(g.Tail) != "00" {
		t.Fatalf("pump tail 00: %+v", g)
	}
	if g := read("boiler_b509_central_heating_pump.tail_f0_false", 0); g.Value != false || g.Outcome != OutcomeValue {
		t.Fatalf("pump tail f0: %+v", g)
	}
	rejected := read("boiler_b509_central_heating_pump.header_only_then_value", 0)
	accepted := read("boiler_b509_central_heating_pump.header_only_then_value", 1)
	if rejected.Outcome != OutcomeNoValue || rejected.Reason != BoilerReasonPayloadRejected || accepted.Value != true ||
		accepted.AcquiredAt.Sub(rejected.AcquiredAt) != 75*time.Millisecond {
		t.Fatalf("pump header-only then value: %+v then %+v", rejected, accepted)
	}
	if g := read("boiler_b509_water_pressure.tail_0080_128_bar", 0); g.Value != 128.0 || g.Unit != "bar" || hex.EncodeToString(g.Tail) != "0080" {
		t.Fatalf("pressure tail 0080: %+v", g)
	}
	for _, c := range []string{"boiler_b509_central_heating_pump.neighbour_register", "boiler_b509_water_pressure.neighbour_register"} {
		if g := read(c, 0); g.Outcome != OutcomeValue || hex.EncodeToString(g.Tail) != hex.EncodeToString(g.Raw) {
			t.Fatalf("neighbour header is not kept as an unchanged bare payload: %+v", g)
		}
	}
	fallback := cases["boiler_b509_flowset_hc_max_c_fallback.value_opcode_echo"]
	if len(fallback.Attempts) != 4 || fallback.Attempts[3].RequestHex != "0da500" || fallback.Final != 55.0 {
		t.Fatalf("fallback case %+v", fallback)
	}
}

func TestVaillantInputRowsBoilerRegisterWriteVectors(t *testing.T) {
	f := boilerRowsLoad[boilerRowsFile](t, "testdata/vectors/boiler_register.json")
	if len(f.Writes) != 7*12 {
		t.Fatalf("%d write cases", len(f.Writes))
	}
	rows := map[string]BoilerWriteRow{}
	for _, w := range BoilerWriteRows {
		rows[w.ID] = w
	}
	for _, v := range f.Writes {
		t.Run(v.Case, func(t *testing.T) {
			row, ok := rows[v.Row]
			if !ok {
				t.Fatalf("no row %s", v.Row)
			}
			finalError, accepted := "", false
			for i, a := range v.Attempts {
				if want := hex.EncodeToString(row.Selector().WithValue(boilerRowsBytes(t, v.ValueHex))); a.RequestHex != want {
					t.Fatalf("attempt %d request %s, row layout %s", i, a.RequestHex, want)
				}
				rec := boilerRowsRecord(t, row.Selector(), a)
				got, err := row.Acknowledge(rec)
				if err != nil {
					t.Fatal(err)
				}
				wantReason := BoilerReasonNone
				switch {
				case a.SendError != "":
					wantReason = BoilerReasonSendError
				case a.ReplyHex == nil:
					wantReason = BoilerReasonNoReply
				case !*a.AckOK:
					wantReason = BoilerReasonAckInvalid
				}
				if got.Reason != wantReason || got.Accepted != (wantReason == BoilerReasonNone) || !got.AcquiredAt.Equal(rec.Clock) {
					t.Fatalf("attempt %d: %+v, want reason %q", i, got, wantReason)
				}
				finalError = got.Error
				if got.Accepted {
					accepted = true
					break
				}
			}
			if accepted != (v.FinalError == "") || (!accepted && finalError != v.FinalError) {
				t.Fatalf("rows end with accepted=%v error %q, donor error %q", accepted, finalError, v.FinalError)
			}
		})
	}
}

// TestVaillantInputRowsBoilerRegisterKey varies every field of the row key and of the acquisition identity.
//
//	Opcode:            a row that ignores the opcode would accept the write header for a read.
//	Register high byte: a row keyed on the low byte only would accept the neighbouring register.
//	Register low byte:  a row keyed on the high byte only would accept the neighbouring register.
//	Selector type:      a row that ignores the geometry would accept a B524 selector with the same register.
//	Source, Target, Clock: a row that fixes or drops them would return a constant.
func TestVaillantInputRowsBoilerRegisterKey(t *testing.T) {
	row := BoilerRegisterRows[0]
	base := row.Selector()
	ok := Record{Selector: base, Request: base.Payload(), Replied: true, Reply: []byte{0x0d, 0x02, 0x00, 0x80, 0x01}, Source: 0x71, Target: 0x08, Clock: boilerRowsBase}
	if _, err := row.Read(ok); err != nil {
		t.Fatal(err)
	}
	for name, sel := range map[string]Selector{
		"opcode": B509Selector{Opcode: BoilerWrite, Register: base.Register}, "register high": B509Selector{Opcode: base.Opcode, Register: base.Register ^ 0x0100},
		"register low": B509Selector{Opcode: base.Opcode, Register: base.Register ^ 0x0001}, "geometry": B524Selector{Opcode: OpcodeLocal, Operation: OperationRead, Register: base.Register},
	} {
		bad := ok
		bad.Selector = sel
		if _, err := row.Read(bad); err != ErrBoilerSelectorMismatch {
			t.Fatalf("%s: err %v", name, err)
		}
	}
	w := BoilerWriteRows[0]
	if _, err := w.Acknowledge(Record{Selector: row.Selector()}); err != ErrBoilerSelectorMismatch {
		t.Fatalf("read selector accepted by a write row: %v", err)
	}
	moved := ok
	moved.Source, moved.Target, moved.Clock = 0x72, 0x09, boilerRowsBase.Add(time.Hour)
	a, _ := row.Read(ok)
	b, _ := row.Read(moved)
	if a.Source == b.Source || a.Target == b.Target || a.AcquiredAt.Equal(b.AcquiredAt) || !reflect.DeepEqual(a.Value, b.Value) {
		t.Fatalf("acquisition fields do not follow the record: %+v %+v", a, b)
	}
}
