package input

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"reflect"
	"strings"
	"testing"
	"time"
)

type scheduleRowsSeen struct {
	AtMs    int64   `json:"at_ms"`
	Source  byte    `json:"source"`
	Target  byte    `json:"target"`
	Primary byte    `json:"primary"`
	Second  byte    `json:"secondary"`
	Request string  `json:"request_hex"`
	Reply   *string `json:"reply_hex"`
}

type scheduleRowsDonor struct {
	Outcome     string             `json:"outcome"`
	Config      *map[string]any    `json:"config"`
	Slots       []int              `json:"slots"`
	Slot        *map[string]any    `json:"slot"`
	Code        *int               `json:"code"`
	Error       string             `json:"error"`
	Description string             `json:"description"`
	Seen        []scheduleRowsSeen `json:"seen"`
}

type scheduleRowsCase struct {
	Case    string            `json:"case"`
	Class   string            `json:"class"`
	Request string            `json:"request_hex"`
	Reply   *string           `json:"reply_hex"`
	Donor   scheduleRowsDonor `json:"donor"`
}

type scheduleRowsFile struct {
	Format       string            `json:"format"`
	DonorSources map[string]string `json:"donor_sources"`
	Cases        []scheduleRowsCase
}

func scheduleRowsLoad(t *testing.T) scheduleRowsFile {
	t.Helper()
	f := boilerRowsLoad[scheduleRowsFile](t, "testdata/vectors/schedule.json")
	inv := boilerRowsLoad[boilerRowsInventory](t, "testdata/gateway_selectors.json")
	for _, s := range inv.SourceFiles {
		if want, ok := f.DonorSources[s.Path]; ok && want != s.Digest {
			t.Fatalf("vector donor source %s differs from the inventoried digest", s.Path)
		}
	}
	return f
}

func scheduleRowsSelector(t *testing.T, req []byte) Selector {
	t.Helper()
	switch req[0] {
	case scheduleConfigOpcode:
		return ScheduleConfigSelector{Zone: req[1], Circuit: req[2]}
	case scheduleCountsOpcode:
		return ScheduleCountsSelector{Zone: req[1], Circuit: req[2]}
	case scheduleReadOpcode:
		return ScheduleTimerSelector{Zone: req[1], Circuit: req[2], Weekday: req[3], Index: req[4]}
	case scheduleWriteOpcode:
		return ScheduleTimerWriteSelector{Zone: req[1], Circuit: req[2], Weekday: req[3], Index: req[4], Count: req[5],
			StartHour: req[6], StartMinute: req[7], EndHour: req[8], EndMinute: req[9], TemperatureRaw: uint16(req[10]) | uint16(req[11])<<8}
	}
	t.Fatalf("request %x is not a schedule request", req)
	return nil
}

func scheduleRowsRow(t *testing.T, opcode byte) ScheduleRow {
	t.Helper()
	for _, r := range ScheduleRows {
		if r.Opcode == opcode {
			return r
		}
	}
	t.Fatalf("no row for opcode %02x", opcode)
	return ScheduleRow{}
}

func scheduleRowsRecord(t *testing.T, c scheduleRowsCase, base time.Time) Record {
	t.Helper()
	req := boilerRowsBytes(t, c.Request)
	if len(c.Donor.Seen) == 0 {
		t.Fatalf("%s: the donor sent no request", c.Case)
	}
	seen := c.Donor.Seen[0]
	rec := Record{Selector: scheduleRowsSelector(t, req), Request: boilerRowsBytes(t, seen.Request), Source: seen.Source, Target: seen.Target,
		Clock: base.Add(time.Duration(seen.AtMs) * time.Millisecond)}
	if c.Reply != nil {
		rec.Reply, rec.Replied = boilerRowsBytes(t, *c.Reply), true
	}
	return rec
}

func scheduleRowsNormal(t *testing.T, v any) any {
	t.Helper()
	raw, err := json.Marshal(v)
	if err != nil {
		t.Fatal(err)
	}
	var out any
	if err := json.Unmarshal(raw, &out); err != nil {
		t.Fatal(err)
	}
	return out
}

func TestVaillantInputRowsScheduleInventory(t *testing.T) {
	inv := boilerRowsLoad[boilerRowsInventory](t, "testdata/gateway_selectors.json")
	byID := map[string]ScheduleRow{}
	for _, r := range ScheduleRows {
		byID[r.ID] = r
	}
	seen := 0
	for _, s := range inv.Selectors {
		if s.Service != "B555" {
			continue
		}
		r, ok := byID[s.ID]
		if !ok {
			t.Fatalf("inventoried selector %s has no row", s.ID)
		}
		seen++
		if r.NativeField != s.NativeViewField || !strings.HasSuffix(s.ID, string(r.Operation)) {
			t.Fatalf("row %s: native field %q or operation %q differ from the inventory %q", r.ID, r.NativeField, r.Operation, s.NativeViewField)
		}
	}
	if seen != len(ScheduleRows) {
		t.Fatalf("%d inventoried B555 selectors, %d rows", seen, len(ScheduleRows))
	}
}

func TestVaillantInputRowsScheduleOpcodes(t *testing.T) {
	raw := boilerRowsLoad[struct {
		Selectors []struct {
			ID, Service, Opcode, Primary, Secondary string
		}
	}](t, "testdata/gateway_selectors.json")
	byID := map[string]ScheduleRow{}
	for _, r := range ScheduleRows {
		byID[r.ID] = r
	}
	n := 0
	for _, s := range raw.Selectors {
		if s.Service != "B555" {
			continue
		}
		r := byID[s.ID]
		n++
		if !strings.EqualFold(s.Opcode, "0x"+hex.EncodeToString([]byte{r.Opcode})) || s.Primary != "0xB5" || s.Secondary != "0x55" {
			t.Fatalf("%s: inventory opcode %s service %s%s, row opcode %02x", s.ID, s.Opcode, s.Primary, s.Secondary, r.Opcode)
		}
	}
	if n != 4 {
		t.Fatalf("%d B555 selectors", n)
	}
}

func TestVaillantInputRowsScheduleVectors(t *testing.T) {
	f := scheduleRowsLoad(t)
	if f.Format != "ef04-schedule-vectors/v1" || len(f.Cases) != 33 {
		t.Fatalf("format %q with %d cases", f.Format, len(f.Cases))
	}
	base := time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)
	classes := map[string]bool{}
	for _, c := range f.Cases {
		classes[c.Class] = true
		req := boilerRowsBytes(t, c.Request)
		row := scheduleRowsRow(t, req[0])
		rec := scheduleRowsRecord(t, c, base)
		if hex.EncodeToString(rec.Selector.Payload()) != c.Request || hex.EncodeToString(rec.Request) != c.Request {
			t.Fatalf("%s: selector payload %x differs from the request the donor sent %s", c.Case, rec.Selector.Payload(), c.Request)
		}
		if row.Operation == ScheduleTimerWrite {
			scheduleRowsCheckWrite(t, c, row, rec)
		} else {
			scheduleRowsCheckRead(t, c, row, rec)
		}
	}
	for _, want := range []string{"value", "sentinel", "boundary", "status", "short", "no_reply", "ack", "empty"} {
		if !classes[want] {
			t.Fatalf("no vector of class %s", want)
		}
	}
}

func scheduleRowsCheckRead(t *testing.T, c scheduleRowsCase, row ScheduleRow, rec Record) {
	t.Helper()
	res, err := row.Read(rec)
	if err != nil {
		t.Fatalf("%s: %v", c.Case, err)
	}
	if !reflect.DeepEqual(res.Raw, rec.Reply) || !reflect.DeepEqual(res.Request, rec.Request) || res.Source != rec.Source || res.Target != rec.Target || !res.AcquiredAt.Equal(rec.Clock) {
		t.Fatalf("%s: evidence, request, addresses or time not kept: %+v", c.Case, res)
	}
	if string(res.Outcome) != c.Donor.Outcome {
		t.Fatalf("%s: outcome %s, donor %s", c.Case, res.Outcome, c.Donor.Outcome)
	}
	if res.Outcome == OutcomeNoValue {
		want := map[string]ScheduleReason{"no_reply": ScheduleReasonNoReply, "short": ScheduleReasonShort, "status": ScheduleReasonStatus}[c.Class]
		if res.Reason != want || res.Config != nil || res.SlotsPerDay != nil || res.Timer != nil || res.Unit != "" {
			t.Fatalf("%s: reason %q (want %q) or a value kept for a read without one: %+v", c.Case, res.Reason, want, res)
		}
		return
	}
	switch row.Operation {
	case ScheduleConfigRead:
		if !reflect.DeepEqual(scheduleRowsNormal(t, res.Config), scheduleRowsNormal(t, c.Donor.Config)) || res.Unit != "degC" {
			t.Fatalf("%s: config %+v, donor %+v", c.Case, res.Config, c.Donor.Config)
		}
	case ScheduleCountsRead:
		if !reflect.DeepEqual(res.SlotsPerDay, c.Donor.Slots) || res.Unit != "" {
			t.Fatalf("%s: slots %v, donor %v", c.Case, res.SlotsPerDay, c.Donor.Slots)
		}
	default:
		if !reflect.DeepEqual(scheduleRowsNormal(t, res.Timer), scheduleRowsNormal(t, c.Donor.Slot)) || res.Unit != "degC" {
			t.Fatalf("%s: timer %+v, donor %+v", c.Case, res.Timer, c.Donor.Slot)
		}
	}
	if res.Reason != ScheduleReasonNone {
		t.Fatalf("%s: reason %q on a value", c.Case, res.Reason)
	}
}

func scheduleRowsCheckWrite(t *testing.T, c scheduleRowsCase, row ScheduleRow, rec Record) {
	t.Helper()
	res, err := row.Acknowledge(rec)
	if err != nil {
		t.Fatalf("%s: %v", c.Case, err)
	}
	if !reflect.DeepEqual(res.Raw, rec.Reply) || !reflect.DeepEqual(res.Request, rec.Request) || res.Source != rec.Source || res.Target != rec.Target || !res.AcquiredAt.Equal(rec.Clock) {
		t.Fatalf("%s: evidence, request, addresses or time not kept: %+v", c.Case, res)
	}
	if c.Donor.Code == nil || int(res.Status) != *c.Donor.Code || res.Error != c.Donor.Error {
		t.Fatalf("%s: status %02x error %q, donor %v error %q", c.Case, res.Status, res.Error, c.Donor.Code, c.Donor.Error)
	}
	switch c.Donor.Outcome {
	case "status":
		if res.Description != c.Donor.Description || res.Accepted != (res.Status == 0) || res.Reason != ScheduleReasonNone {
			t.Fatalf("%s: description %q accepted %v reason %q, donor %q", c.Case, res.Description, res.Accepted, res.Reason, c.Donor.Description)
		}
	default:
		want := map[string]ScheduleReason{"no_reply": ScheduleReasonNoReply, "empty": ScheduleReasonEmpty}[c.Class]
		if res.Accepted || res.Reason != want || res.Description != "" {
			t.Fatalf("%s: accepted %v reason %q (want %q) description %q", c.Case, res.Accepted, res.Reason, want, res.Description)
		}
	}
}

func TestVaillantInputRowsScheduleWriteEmptyNegative(t *testing.T) {
	f := scheduleRowsLoad(t)
	var one, empty *scheduleRowsCase
	for i := range f.Cases {
		switch f.Cases[i].Case {
		case "write_status_accepted":
			one = &f.Cases[i]
		case "write_empty_reply":
			empty = &f.Cases[i]
		}
	}
	if one == nil || empty == nil {
		t.Fatal("write vectors missing")
	}
	row := scheduleRowsRow(t, scheduleWriteOpcode)
	base := time.Unix(0, 0).UTC()
	a, _ := row.Acknowledge(scheduleRowsRecord(t, *one, base))
	b, _ := row.Acknowledge(scheduleRowsRecord(t, *empty, base))
	if a.Status != 0x00 || !a.Accepted || b.Status != 0xff || b.Accepted || b.Error != "b555 timer write: empty response" || a.Error != "" {
		t.Fatalf("one-byte status %+v, empty reply %+v", a, b)
	}
	if len(b.Raw) != 0 || b.Reason != ScheduleReasonEmpty {
		t.Fatalf("empty reply evidence %+v", b)
	}
}

func TestVaillantInputRowsScheduleSendErrorIsNoValue(t *testing.T) {
	for _, row := range ScheduleRows {
		var sel Selector
		switch row.Operation {
		case ScheduleConfigRead:
			sel = ScheduleConfigSelector{Zone: 1, Circuit: 0}
		case ScheduleCountsRead:
			sel = ScheduleCountsSelector{Zone: 1, Circuit: 0}
		case ScheduleTimerRead:
			sel = ScheduleTimerSelector{Zone: 1, Circuit: 0, Weekday: 2, Index: 1}
		default:
			sel = ScheduleTimerWriteSelector{Zone: 1}
		}
		rec := Record{Selector: sel, Request: sel.Payload(), Replied: true, Reply: []byte{0, 0, 0, 0, 0, 0, 0, 0, 0}, Diagnostics: Diagnostics{SendError: "synthetic"}}
		if row.Operation == ScheduleTimerWrite {
			res, _ := row.Acknowledge(rec)
			if res.Reason != ScheduleReasonSendError || res.Accepted || res.Status != 0xff || res.Error != "b555 timer write: no response" {
				t.Fatalf("%s: %+v", row.ID, res)
			}
			continue
		}
		res, _ := row.Read(rec)
		if res.Outcome != OutcomeNoValue || res.Reason != ScheduleReasonSendError {
			t.Fatalf("%s: %+v", row.ID, res)
		}
	}
}

func scheduleRowsSelectors() map[ScheduleOperation]Selector {
	return map[ScheduleOperation]Selector{
		ScheduleConfigRead: ScheduleConfigSelector{Zone: 0, Circuit: 0},
		ScheduleCountsRead: ScheduleCountsSelector{Zone: 0, Circuit: 0},
		ScheduleTimerRead:  ScheduleTimerSelector{Zone: 0, Circuit: 0, Weekday: 0, Index: 0},
		ScheduleTimerWrite: ScheduleTimerWriteSelector{},
	}
}

// TestVaillantInputRowsScheduleKey varies every field of the row key and of the selector it claims to cover.
//
//	opcode (row key)             a row accepting another opcode's selector
//	geometry (selector type)     a row accepting a selector of another B5 geometry
//	zone, circuit, weekday,      fields that are not part of the key: a row that compared one of them to a constant
//	index, count, times,         would reject a selector of another program or slot
//	temperature
func TestVaillantInputRowsScheduleKey(t *testing.T) {
	selectors := scheduleRowsSelectors()
	for _, row := range ScheduleRows {
		for op, sel := range selectors {
			rec := Record{Selector: sel, Request: sel.Payload(), Replied: true, Reply: []byte{0, 0, 0, 0, 0, 0, 0, 0, 0}}
			var err error
			if row.Operation == ScheduleTimerWrite {
				_, err = row.Acknowledge(rec)
			} else {
				_, err = row.Read(rec)
			}
			if (op == row.Operation) != (err == nil) {
				t.Fatalf("row %s with the %s selector: error %v", row.ID, op, err)
			}
			if err != nil && !errors.Is(err, ErrScheduleSelector) {
				t.Fatalf("row %s: %v", row.ID, err)
			}
		}
	}
	foreign := []Selector{
		B524Selector{Opcode: OpcodeLocal, Operation: OperationRead, Group: 0, Instance: 0, Register: 0xa3},
		B509Selector{Opcode: BoilerRead, Register: 0xa300},
		InterfaceIdentitySelector{ID: scheduleConfigOpcode},
		DiagnosticSelector{Family: scheduleConfigOpcode, Selector: 0},
		LiveMonitorSelector{Action: LiveMonitorRead},
	}
	for _, row := range ScheduleRows {
		for _, sel := range foreign {
			var err error
			if row.Operation == ScheduleTimerWrite {
				_, err = row.Acknowledge(Record{Selector: sel})
			} else {
				_, err = row.Read(Record{Selector: sel})
			}
			if !errors.Is(err, ErrScheduleSelector) {
				t.Fatalf("row %s accepted a %T selector: %v", row.ID, sel, err)
			}
		}
	}
	reply := []byte{0, 3, 10, 10, 1, 1, 5, 30, 0}
	for _, zone := range []byte{0, 1, 2, 0xff} {
		for _, circuit := range []byte{0, 1, 2, 3, 4} {
			row := scheduleRowsRow(t, scheduleConfigOpcode)
			rec := Record{Selector: ScheduleConfigSelector{Zone: zone, Circuit: circuit}, Replied: true, Reply: reply}
			if res, err := row.Read(rec); err != nil || res.Outcome != OutcomeValue {
				t.Fatalf("zone %02x circuit %02x: %v %+v", zone, circuit, err, res)
			}
		}
	}
	timer := scheduleRowsRow(t, scheduleReadOpcode)
	for _, sel := range []ScheduleTimerSelector{{0, 0, 0, 0}, {2, 1, 6, 2}, {0xff, 4, 3, 9}} {
		rec := Record{Selector: sel, Replied: true, Reply: []byte{0, 6, 30, 21, 0, 205, 0}}
		if res, err := timer.Read(rec); err != nil || res.Outcome != OutcomeValue || res.Timer.StartHour != 6 {
			t.Fatalf("timer %+v: %v %+v", sel, err, res)
		}
	}
	write := scheduleRowsRow(t, scheduleWriteOpcode)
	for _, sel := range []ScheduleTimerWriteSelector{{}, {Zone: 0xff, Circuit: 2, Weekday: 6, Index: 2, Count: 3, StartHour: 24, StartMinute: 59, EndHour: 24, EndMinute: 59, TemperatureRaw: 0xffff}} {
		rec := Record{Selector: sel, Replied: true, Reply: []byte{0}}
		if res, err := write.Acknowledge(rec); err != nil || !res.Accepted {
			t.Fatalf("write %+v: %v %+v", sel, err, res)
		}
	}
}

func TestVaillantInputRowsScheduleBoundaries(t *testing.T) {
	cfg := scheduleRowsRow(t, scheduleConfigOpcode)
	rec := func(reply []byte) Record {
		return Record{Selector: ScheduleConfigSelector{}, Replied: true, Reply: reply}
	}
	if res, _ := cfg.Read(rec(make([]byte, 8))); res.Reason != ScheduleReasonShort {
		t.Fatalf("eight bytes: %+v", res)
	}
	if res, _ := cfg.Read(rec(make([]byte, 9))); res.Outcome != OutcomeValue {
		t.Fatalf("nine bytes: %+v", res)
	}
	if res, _ := cfg.Read(rec(append([]byte{1}, make([]byte, 8)...))); res.Reason != ScheduleReasonStatus {
		t.Fatalf("status: %+v", res)
	}
	if res, _ := cfg.Read(rec(append([]byte{1}, make([]byte, 7)...))); res.Reason != ScheduleReasonShort {
		t.Fatalf("short before status: %+v", res)
	}
	tm := scheduleRowsRow(t, scheduleReadOpcode)
	for raw, wantNil := range map[uint16]bool{0xffff: true, 0xfffe: false, 0: false} {
		res, _ := tm.Read(Record{Selector: ScheduleTimerSelector{}, Replied: true, Reply: []byte{0, 0, 0, 0, 0, byte(raw), byte(raw >> 8)}})
		if (res.Timer.TemperatureC == nil) != wantNil || (res.Timer.TemperatureRaw == nil) != wantNil {
			t.Fatalf("raw %04x: %+v", raw, res.Timer)
		}
	}
}

func TestVaillantInputRowsScheduleStatusText(t *testing.T) {
	f := scheduleRowsLoad(t)
	seen := map[byte]bool{}
	for _, c := range f.Cases {
		if c.Donor.Code != nil && c.Donor.Outcome == "status" {
			seen[byte(*c.Donor.Code)] = true
			if got := ScheduleStatusText(byte(*c.Donor.Code)); got != c.Donor.Description {
				t.Fatalf("%s: %q, donor %q", c.Case, got, c.Donor.Description)
			}
		}
	}
	for _, code := range []byte{0x00, 0x01, 0x03, 0x06, 0x07} {
		if !seen[code] {
			t.Fatalf("no vector for status %02x", code)
		}
	}
}
