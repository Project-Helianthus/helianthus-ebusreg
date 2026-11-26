package input

import (
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"time"
)

const (
	replayIndexedSession = "current/sessions/vaillant/poller-full-site/session.json"
	replayDirectory      = "testdata/session"
)

type replaySessionEntry struct {
	Case            string `json:"case"`
	Scenario        string `json:"scenario"`
	IndexPath       string `json:"index_path"`
	InputFile       string `json:"input_file"`
	InputDigest     string `json:"input_sha256"`
	InputRole       string `json:"input_role"`
	OutputIndexPath string `json:"output_index_path"`
	OutputDigest    string `json:"output_sha256"`
	NativeFile      string `json:"native_file"`
	Acquisitions    int    `json:"acquisitions"`
	Sends           int    `json:"sends"`
}

type replaySessionIndex struct {
	Format   string               `json:"format"`
	Index    string               `json:"index"`
	Sessions []replaySessionEntry `json:"sessions"`
}

type replayDecode struct {
	Purpose string `json:"purpose"`
	Codec   string `json:"codec"`
	Outcome string `json:"outcome"`
	Stage   string `json:"stage"`
	Kind    string `json:"kind"`
	Text    string `json:"text"`
}

type replayAcquisition struct {
	Seq       int             `json:"seq"`
	Class     string          `json:"class"`
	Step      string          `json:"step"`
	AtMs      int64           `json:"at_ms"`
	Via       string          `json:"via"`
	Source    byte            `json:"source"`
	Target    byte            `json:"target"`
	Primary   byte            `json:"primary"`
	Secondary byte            `json:"secondary"`
	Request   string          `json:"request_hex"`
	Reply     *string         `json:"reply_hex"`
	SendError string          `json:"send_error"`
	Seam      json.RawMessage `json:"seam"`
	Decodes   []replayDecode  `json:"decodes"`
	Mismatch  string          `json:"mismatch"`
}

type replayNative struct {
	Session      string              `json:"session"`
	Scenario     string              `json:"scenario"`
	Base         string              `json:"base"`
	Acquisitions []replayAcquisition `json:"acquisitions"`
}

type replaySeam struct {
	Accepted   bool   `json:"accepted"`
	PayloadHex string `json:"payload_hex"`
	Error      string `json:"error"`
}

type replayStats struct {
	sessions, acquisitions, decodes, differences, noRow int
	b555                                                int
	byFamily                                            map[string]int
}

func (s *replayStats) diff(t *testing.T, format string, args ...any) {
	t.Helper()
	s.differences++
	if s.differences <= 20 {
		t.Errorf(format, args...)
	}
}

func replayStage(r NoValueReason) string {
	switch r {
	case ReasonNoReply:
		return "no_reply"
	case ReasonRejected:
		return "reply_rejected"
	}
	return "decoded_none"
}

func replayBoilerStage(r BoilerReadReason) string {
	switch r {
	case BoilerReasonNoReply, BoilerReasonSendError:
		return "no_reply"
	case BoilerReasonPayloadRejected:
		return "reply_rejected"
	}
	return "decoded_none"
}

func replayBoilerText(v any) string {
	switch x := v.(type) {
	case float64:
		return strconv.FormatFloat(x, 'g', -1, 64)
	case int:
		return strconv.Itoa(x)
	case bool:
		return strconv.FormatBool(x)
	case string:
		return x
	}
	return fmt.Sprintf("%v", v)
}

func replayRecord(t *testing.T, a replayAcquisition, selector Selector, base time.Time) Record {
	t.Helper()
	rec := Record{Selector: selector, Request: boilerRowsBytes(t, a.Request), Source: a.Source, Target: a.Target,
		Clock: base.Add(time.Duration(a.AtMs) * time.Millisecond), Diagnostics: Diagnostics{SendError: a.SendError}}
	if a.Reply != nil {
		rec.Reply, rec.Replied = boilerRowsBytes(t, *a.Reply), true
	}
	return rec
}

func replaySeamOf(t *testing.T, a replayAcquisition) replaySeam {
	t.Helper()
	var seam replaySeam
	if err := json.Unmarshal(a.Seam, &seam); err != nil {
		t.Fatalf("acquisition %d seam: %v", a.Seq, err)
	}
	return seam
}

func replayExtended(t *testing.T, st *replayStats, a replayAcquisition, base time.Time) {
	t.Helper()
	data := boilerRowsBytes(t, a.Request)
	sel := B524Selector{Opcode: data[0], Operation: data[1], Group: data[2], Instance: data[3], Register: binary.LittleEndian.Uint16(data[4:6])}
	rec := replayRecord(t, a, sel, base)
	if sel.Operation == OperationWrite {
		replayExtendedWrite(t, st, a, rec, sel)
		return
	}
	seam := replaySeamOf(t, a)
	var payload []byte
	if rec.Replied {
		var line string
		var accepted bool
		payload, line, accepted = ParseReply(rec.Reply, sel)
		if accepted != seam.Accepted || (accepted && hex.EncodeToString(payload) != seam.PayloadHex) || line != a.Mismatch {
			st.diff(t, "acquisition %d %s: header rules accepted=%v payload=%x line=%q, donor accepted=%v payload=%s line=%q",
				a.Seq, a.Request, accepted, payload, line, seam.Accepted, seam.PayloadHex, a.Mismatch)
		}
	}
	matched := 0
	for _, row := range ExtendedRows() {
		if !row.Matches(sel) {
			continue
		}
		matched++
		want := 0
		for _, dec := range row.Decodes {
			if dec.Codec == CodecWriteAck || dec.Codec == CodecNotSent {
				continue
			}
			want++
			donor, found := replayDonorDecode(a, dec)
			if !found {
				st.diff(t, "acquisition %d %s row %s: no donor outcome for %s %s", a.Seq, a.Request, row.ID, dec.Purpose, dec.Codec)
				continue
			}
			st.decodes++
			replayCompareExtended(t, st, a, rec, row, dec, donor, seam)
		}
		if got := len(replayDonorDecodes(a, row)); got != want {
			st.diff(t, "acquisition %d %s row %s: %d donor outcomes for %d row decodes", a.Seq, a.Request, row.ID, got, want)
		}
	}
	if matched == 0 {
		st.diff(t, "acquisition %d: no row for %s", a.Seq, a.Request)
	}
}

func replayDonorDecode(a replayAcquisition, dec Decode) (replayDecode, bool) {
	for _, d := range a.Decodes {
		if d.Purpose == string(dec.Purpose) && d.Codec == string(dec.Codec) {
			return d, true
		}
	}
	return replayDecode{}, false
}

func replayDonorDecodes(a replayAcquisition, row Row) []replayDecode {
	var out []replayDecode
	for _, dec := range row.Decodes {
		if d, ok := replayDonorDecode(a, dec); ok {
			out = append(out, d)
		}
	}
	return out
}

func replayCompareExtended(t *testing.T, st *replayStats, a replayAcquisition, rec Record, row Row, dec Decode, donor replayDecode, seam replaySeam) {
	t.Helper()
	res := row.Decode(rec, dec)
	label := fmt.Sprintf("acquisition %d %s row %s %s/%s", a.Seq, a.Request, row.ID, dec.Purpose, dec.Codec)
	if !res.AcquiredAt.Equal(rec.Clock) || res.Selector != rec.Selector.(B524Selector) || string(res.Reply) != string(rec.Reply) || res.Unit != row.Catalog().Unit {
		st.diff(t, "%s: evidence, selector, time or unit not kept: %+v", label, res)
	}
	if dec.Codec == CodecProbe {
		rowCoherent := res.Outcome == OutcomeValue && res.Value.Text == "true"
		donorCoherent := donor.Outcome == "value" && donor.Text == "true"
		if rowCoherent != donorCoherent || (rec.Replied && (donor.Outcome != string(res.Outcome) || donor.Kind != res.Value.Kind || donor.Text != res.Value.Text)) {
			st.diff(t, "%s: coherent %v, donor %v", label, rowCoherent, donorCoherent)
		}
		return
	}
	if string(res.Outcome) != donor.Outcome {
		st.diff(t, "%s: outcome %s, donor %s", label, res.Outcome, donor.Outcome)
		return
	}
	if res.Diagnostics.MismatchText != a.Mismatch {
		st.diff(t, "%s: mismatch line %q, donor %q", label, res.Diagnostics.MismatchText, a.Mismatch)
	}
	if res.Outcome == OutcomeValue {
		if res.Value.Kind != donor.Kind || res.Value.Text != donor.Text || (seam.Accepted && hex.EncodeToString(res.Raw) != seam.PayloadHex) {
			st.diff(t, "%s: value %+v raw %x, donor %s %q raw %s", label, res.Value, res.Raw, donor.Kind, donor.Text, seam.PayloadHex)
		}
		return
	}
	if stage := replayStage(res.Reason); stage != donor.Stage {
		st.diff(t, "%s: no-value stage %s (%s), donor %s", label, stage, res.Reason, donor.Stage)
	}
}

func replayExtendedWrite(t *testing.T, st *replayStats, a replayAcquisition, rec Record, sel B524Selector) {
	t.Helper()
	seam := replaySeamOf(t, a)
	for _, row := range ExtendedRows() {
		if !row.Matches(sel) {
			continue
		}
		res := row.Decode(rec, Decode{PurposeWrite, CodecWriteAck})
		st.decodes++
		if res.Reason == ReasonOutsideRow || (res.Outcome == OutcomeValue) != (seam.Error == "") || !res.AcquiredAt.Equal(rec.Clock) {
			st.diff(t, "acquisition %d %s row %s: write outcome %s, donor error %q", a.Seq, a.Request, row.ID, res.Outcome, seam.Error)
		}
		return
	}
	st.noRow++
}

func replayBoiler(t *testing.T, st *replayStats, a replayAcquisition, base time.Time) {
	t.Helper()
	data := boilerRowsBytes(t, a.Request)
	register := binary.BigEndian.Uint16(data[1:3])
	sel := B509Selector{Opcode: data[0], Register: register}
	rec := replayRecord(t, a, sel, base)
	seam := replaySeamOf(t, a)
	if data[0] == BoilerWrite {
		for _, row := range BoilerWriteRows {
			if row.Register != register {
				continue
			}
			res, err := row.Acknowledge(rec)
			if err != nil || res.Accepted != (seam.Error == "") || (!res.Accepted && res.Error != seam.Error) || !res.AcquiredAt.Equal(rec.Clock) || string(res.Raw) != string(rec.Reply) {
				st.diff(t, "acquisition %d %s row %s: %+v %v, donor error %q", a.Seq, a.Request, row.ID, res, err, seam.Error)
			}
			st.decodes++
			return
		}
		st.diff(t, "acquisition %d: no write row for %s", a.Seq, a.Request)
		return
	}
	matched := 0
	for _, row := range BoilerRegisterRows {
		if row.Register != register {
			continue
		}
		matched++
		res, err := row.Read(rec)
		st.decodes++
		label := fmt.Sprintf("acquisition %d %s row %s", a.Seq, a.Request, row.ID)
		donor, found := replayDonorDecode(a, Decode{PurposeValueRead, Codec(row.Codec)})
		if err != nil || !found {
			st.diff(t, "%s: %v, donor outcome found %v", label, err, found)
			continue
		}
		if !res.AcquiredAt.Equal(rec.Clock) || string(res.Raw) != string(rec.Reply) || res.Unit != row.Unit || res.Source != a.Source || res.Target != a.Target {
			st.diff(t, "%s: evidence, time, unit or addresses not kept: %+v", label, res)
		}
		if rec.Replied {
			if tail, ok := ParseBoilerReadPayload(rec.Reply, register); ok != seam.Accepted || (ok && hex.EncodeToString(tail) != seam.PayloadHex) {
				st.diff(t, "%s: payload %x accepted %v, donor %s accepted %v", label, tail, ok, seam.PayloadHex, seam.Accepted)
			}
		}
		if string(res.Outcome) != donor.Outcome {
			st.diff(t, "%s: outcome %s, donor %s", label, res.Outcome, donor.Outcome)
			continue
		}
		if res.Outcome == OutcomeValue {
			if string(res.Kind) != donor.Kind || replayBoilerText(res.Value) != donor.Text || hex.EncodeToString(res.Tail) != seam.PayloadHex {
				st.diff(t, "%s: value %v/%v tail %x, donor %s %q tail %s", label, res.Kind, res.Value, res.Tail, donor.Kind, donor.Text, seam.PayloadHex)
			}
		} else if stage := replayBoilerStage(res.Reason); stage != donor.Stage {
			st.diff(t, "%s: no-value stage %s (%s), donor %s", label, stage, res.Reason, donor.Stage)
		}
	}
	if matched == 0 {
		st.diff(t, "acquisition %d: no row for %s", a.Seq, a.Request)
	}
}

func replayScheduleClass(row ScheduleRow, rec Record) string {
	switch {
	case !rec.Replied:
		return "no_reply"
	case row.Operation == ScheduleTimerWrite && len(rec.Reply) == 0:
		return "empty"
	case row.Operation == ScheduleTimerWrite:
		return "status"
	case len(rec.Reply) < row.MinReply:
		return "short"
	case rec.Reply[0] != 0:
		return "status"
	}
	return "value"
}

func replaySchedule(t *testing.T, st *replayStats, a replayAcquisition, base time.Time) {
	t.Helper()
	req := boilerRowsBytes(t, a.Request)
	rec := replayRecord(t, a, scheduleRowsSelector(t, req), base)
	var donor scheduleRowsDonor
	if err := json.Unmarshal(a.Seam, &donor); err != nil {
		t.Fatalf("acquisition %d seam: %v", a.Seq, err)
	}
	row := scheduleRowsRow(t, req[0])
	c := scheduleRowsCase{Case: fmt.Sprintf("acquisition %d %s", a.Seq, a.Request), Class: replayScheduleClass(row, rec), Request: a.Request, Reply: a.Reply, Donor: donor}
	if row.Operation == ScheduleTimerWrite {
		scheduleRowsCheckWrite(t, c, row, rec)
	} else {
		scheduleRowsCheckRead(t, c, row, rec)
	}
	st.b555++
	st.decodes++
}

func replayOpenSessions(t *testing.T) (replaySessionIndex, map[string]replayNative) {
	t.Helper()
	index := boilerRowsLoad[replaySessionIndex](t, filepath.Join(replayDirectory, "sessions.json"))
	natives := map[string]replayNative{}
	for _, e := range index.Sessions {
		raw, err := os.ReadFile(filepath.Join(replayDirectory, e.InputFile))
		if err != nil {
			t.Fatalf("session %s: %v", e.Case, err)
		}
		if sum := sha256.Sum256(raw); hex.EncodeToString(sum[:]) != e.InputDigest || e.InputRole != "input" {
			t.Fatalf("session %s: committed input has sha256 %x, the index entry records %s with role %q", e.Case, sum, e.InputDigest, e.InputRole)
		}
		var named struct{ Session, Scenario string }
		if err := json.Unmarshal(raw, &named); err != nil {
			t.Fatal(err)
		}
		native := boilerRowsLoad[replayNative](t, filepath.Join(replayDirectory, e.NativeFile))
		if native.Session != named.Session || native.Scenario != named.Scenario || e.Scenario != named.Scenario {
			t.Fatalf("session %s: record is for %s/%s, input is %s/%s", e.Case, native.Session, native.Scenario, named.Session, named.Scenario)
		}
		classes := map[string]int{}
		for _, a := range native.Acquisitions {
			classes[a.Class]++
		}
		if classes["acquisition"] != e.Acquisitions || classes["send"] != e.Sends {
			t.Fatalf("session %s: record holds %d acquisitions and %d sends, the index lists %d and %d", e.Case, classes["acquisition"], classes["send"], e.Acquisitions, e.Sends)
		}
		natives[e.Case] = native
	}
	return index, natives
}

func TestVaillantInputSessionReplay(t *testing.T) {
	index, natives := replayOpenSessions(t)
	if len(index.Sessions) == 0 || len(natives) == 0 {
		t.Fatal("no session is replayed")
	}
	indexed := false
	for _, e := range index.Sessions {
		indexed = indexed || e.IndexPath == replayIndexedSession
	}
	if !indexed {
		t.Fatalf("the indexed session %s is not among the replayed sessions", replayIndexedSession)
	}
	st := &replayStats{byFamily: map[string]int{}}
	for _, e := range index.Sessions {
		native := natives[e.Case]
		t.Run(e.Case, func(t *testing.T) {
			base, err := time.Parse(time.RFC3339, native.Base)
			if err != nil {
				t.Fatal(err)
			}
			schedule := 0
			for _, a := range native.Acquisitions {
				st.acquisitions++
				data := boilerRowsBytes(t, a.Request)
				switch {
				case a.Primary == primaryExtendedRegister && a.Secondary == secondaryExtendedRegister:
					st.byFamily["B524"]++
					replayExtended(t, st, a, base)
				case a.Primary == primaryBoilerRegister && a.Secondary == secondaryBoilerRegister:
					st.byFamily["B509"]++
					replayBoiler(t, st, a, base)
				case a.Primary == primarySchedule && a.Secondary == secondarySchedule && len(data) > 0:
					st.byFamily["B555"]++
					schedule++
					replaySchedule(t, st, a, base)
				default:
					st.diff(t, "acquisition %d has service %02x%02x", a.Seq, a.Primary, a.Secondary)
				}
			}
			if e.Case == "poller-full-site" {
				if schedule == 0 {
					t.Fatal("the schedule step recorded no B555 request")
				}
				replayRequireScheduleOperations(t, native)
			}
			st.sessions++
		})
	}
	t.Run("supplemental/schedule_pass", func(t *testing.T) { replaySchedulePass(t, st) })
	t.Run("supplemental/pump_attempts", func(t *testing.T) { replayPumpAttempts(t, st) })
	t.Run("supplemental/mode_discovery_and_value_read", func(t *testing.T) { replayModeDiscovery(t, st) })
	t.Run("supplemental/flow_retained", func(t *testing.T) { replayFlowRetained(t, st) })
	t.Run("supplemental/b503", func(t *testing.T) { replayB503(t, st) })
	t.Run("supplemental/interface_identity", func(t *testing.T) { replayIdentity(t, st) })
	if st.sessions != len(index.Sessions) || st.sessions == 0 {
		t.Fatalf("%d of %d sessions replayed", st.sessions, len(index.Sessions))
	}
	t.Logf("replayed %d sessions, %d sends (B524 %d, B509 %d, B555 %d) and %d row decodes; %d B524 writes outside the inventory; %d differences",
		st.sessions, st.acquisitions, st.byFamily["B524"], st.byFamily["B509"], st.byFamily["B555"], st.decodes, st.noRow, st.differences)
	if st.differences != 0 {
		t.Fatalf("%d differences between the rows and the recorded poller outcomes", st.differences)
	}
}

func replayRequireScheduleOperations(t *testing.T, native replayNative) {
	t.Helper()
	have := map[byte]bool{}
	for _, a := range native.Acquisitions {
		if a.Step == "refresh_schedules" && a.Primary == primarySchedule && a.Secondary == secondarySchedule {
			have[boilerRowsBytes(t, a.Request)[0]] = true
		}
	}
	for _, op := range []byte{scheduleConfigOpcode, scheduleCountsOpcode, scheduleReadOpcode} {
		if !have[op] {
			t.Fatalf("the schedule step sent no B555 %02x", op)
		}
	}
}

type replaySupplemental struct {
	DonorSources map[string]string `json:"donor_sources"`
	SchedulePass struct {
		Requests []struct {
			AtMs    int64           `json:"at_ms"`
			Source  byte            `json:"source"`
			Target  byte            `json:"target"`
			Request string          `json:"request_hex"`
			Reply   *string         `json:"reply_hex"`
			Donor   json.RawMessage `json:"donor"`
		} `json:"requests"`
		Published int `json:"published_programs"`
	} `json:"schedule_pass"`
	Pump []struct {
		Case     string `json:"case"`
		Target   byte   `json:"target"`
		Source   byte   `json:"source"`
		Attempts []struct {
			AtMs    int64   `json:"at_ms"`
			Request string  `json:"request_hex"`
			Reply   *string `json:"reply_hex"`
			Seam    replaySeam
		} `json:"attempts"`
		Result struct{ Outcome, Kind, Text string } `json:"result"`
	} `json:"pump_attempts"`
	Mode struct {
		Request string `json:"request_hex"`
		Reply   string `json:"reply_hex"`
		Source  byte   `json:"source"`
		Target  byte   `json:"target"`
		Probe   struct {
			Coherent bool `json:"coherent"`
			Sends    int  `json:"sends"`
		} `json:"probe"`
		Value struct {
			Accepted      bool   `json:"accepted"`
			Payload       string `json:"payload_hex"`
			MismatchLines int    `json:"mismatch_lines"`
			Sends         int    `json:"sends"`
		} `json:"value_read"`
		Sends []struct {
			AtMs int64 `json:"at_ms"`
		} `json:"sends"`
	} `json:"mode_discovery"`
	Flow struct {
		Request       string                       `json:"request_hex"`
		Source        byte                         `json:"source"`
		Target        byte                         `json:"target"`
		FlowFirst     float64                      `json:"flow_after_first"`
		FlowAfter     float64                      `json:"flow_after_failure"`
		FailedReadOK  bool                         `json:"failed_read_ok"`
		ReceiptSame   bool                         `json:"receipt_unchanged"`
		ReceiptBefore struct{ Nanoseconds string } `json:"receipt_before"`
		ReceiptAfter  struct{ Nanoseconds string } `json:"receipt_after"`
		Attempts      []struct {
			AtMs  int64   `json:"at_ms"`
			Reply *string `json:"reply_hex"`
			Error string  `json:"send_error"`
			Seam  replaySeam
		} `json:"attempts"`
	} `json:"flow_retained"`
}

func replaySupplementalLoad(t *testing.T) replaySupplemental {
	t.Helper()
	s := boilerRowsLoad[replaySupplemental](t, filepath.Join(replayDirectory, "supplemental.json"))
	inv := boilerRowsLoad[boilerRowsInventory](t, "testdata/gateway_selectors.json")
	for _, f := range inv.SourceFiles {
		if want, ok := s.DonorSources[f.Path]; ok && want != f.Digest {
			t.Fatalf("supplemental donor source %s differs from the inventoried digest", f.Path)
		}
	}
	return s
}

func replaySchedulePass(t *testing.T, st *replayStats) {
	s := replaySupplementalLoad(t)
	base := time.Date(2026, 1, 5, 6, 0, 0, 0, time.UTC)
	reqs := s.SchedulePass.Requests
	if len(reqs) < 3 {
		t.Fatalf("the schedule pass recorded %d B555 requests", len(reqs))
	}
	for i, op := range []byte{scheduleConfigOpcode, scheduleCountsOpcode, scheduleReadOpcode} {
		if boilerRowsBytes(t, reqs[i].Request)[0] != op {
			t.Fatalf("request %d is %s, want operation %02x in the order A3, A4, A5", i, reqs[i].Request, op)
		}
	}
	var counts struct{ Slots []int }
	if err := json.Unmarshal(reqs[1].Donor, &counts); err != nil || len(counts.Slots) != 7 {
		t.Fatalf("slots per day: %v %v", err, counts.Slots)
	}
	used := -1
	for day, n := range counts.Slots {
		if n > 0 {
			used = day
			break
		}
	}
	timer := boilerRowsBytes(t, reqs[2].Request)
	if used < 0 || int(timer[3]) != used || timer[4] != 0 {
		t.Fatalf("the timer read %x is not the first used slot of weekday %d", timer, used)
	}
	var slot struct {
		Slot *struct{ StartHour, EndHour int }
	}
	if err := json.Unmarshal(reqs[2].Donor, &slot); err != nil || slot.Slot == nil || slot.Slot.StartHour != 6 || slot.Slot.EndHour != 21 {
		t.Fatalf("the used slot has no native value: %v %+v", err, slot)
	}
	if s.SchedulePass.Published < 1 {
		t.Fatalf("the donor published %d programs", s.SchedulePass.Published)
	}
	for i, r := range reqs {
		a := replayAcquisition{Seq: i, Step: "schedule_pass", AtMs: r.AtMs, Source: r.Source, Target: r.Target, Request: r.Request, Reply: r.Reply, Seam: r.Donor}
		replaySchedule(t, st, a, base)
	}
	st.acquisitions += len(reqs)
	st.byFamily["B555"] += len(reqs)
}

func replayPumpAttempts(t *testing.T, st *replayStats) {
	s := replaySupplementalLoad(t)
	base := time.Date(2026, 1, 5, 6, 0, 0, 0, time.UTC)
	var row BoilerRegisterRow
	for _, r := range BoilerRegisterRows {
		if r.Register == 0x4400 {
			row = r
		}
	}
	if len(s.Pump) != 3 {
		t.Fatalf("%d pump cases", len(s.Pump))
	}
	sawPair := false
	for _, c := range s.Pump {
		var accepted *BoilerReadResult
		var last BoilerReadResult
		for i, a := range c.Attempts {
			rec := Record{Selector: row.Selector(), Request: boilerRowsBytes(t, a.Request), Source: c.Source, Target: c.Target, Clock: base.Add(time.Duration(a.AtMs) * time.Millisecond)}
			if a.Reply != nil {
				rec.Reply, rec.Replied = boilerRowsBytes(t, *a.Reply), true
			}
			res, err := row.Read(rec)
			if err != nil || !res.AcquiredAt.Equal(rec.Clock) || string(res.Raw) != string(rec.Reply) {
				st.diff(t, "%s attempt %d: %+v %v", c.Case, i, res, err)
				continue
			}
			wantAccepted := a.Seam.Accepted
			if (res.Outcome == OutcomeValue) != wantAccepted || (!wantAccepted && res.Reason != BoilerReasonPayloadRejected) || (wantAccepted && hex.EncodeToString(res.Tail) != a.Seam.PayloadHex) {
				st.diff(t, "%s attempt %d: outcome %s reason %s, donor accepted %v", c.Case, i, res.Outcome, res.Reason, wantAccepted)
			}
			if i > 0 && !res.AcquiredAt.After(last.AcquiredAt) {
				st.diff(t, "%s: attempt %d is not after attempt %d", c.Case, i, i-1)
			}
			if res.Outcome == OutcomeValue && accepted == nil {
				r := res
				accepted = &r
			}
			last = res
			st.decodes++
		}
		if accepted == nil || c.Result.Outcome != "value" || replayBoilerText(accepted.Value) != c.Result.Text || string(accepted.Kind) != c.Result.Kind {
			st.diff(t, "%s: accepted attempt %+v, donor result %+v", c.Case, accepted, c.Result)
		}
		if len(c.Attempts) == 2 && !c.Attempts[0].Seam.Accepted && c.Attempts[1].Seam.Accepted {
			sawPair = true
			if c.Attempts[0].Request != "0d4400" || c.Attempts[1].AtMs <= c.Attempts[0].AtMs {
				st.diff(t, "%s: attempts %+v", c.Case, c.Attempts)
			}
		}
		st.acquisitions += len(c.Attempts)
		st.byFamily["B509"] += len(c.Attempts)
	}
	if !sawPair {
		t.Fatal("no case with a rejected then an accepted attempt")
	}
}

func replayModeDiscovery(t *testing.T, st *replayStats) {
	s := replaySupplementalLoad(t).Mode
	base := time.Date(2026, 1, 5, 6, 0, 0, 0, time.UTC)
	data := boilerRowsBytes(t, s.Request)
	sel := B524Selector{Opcode: data[0], Operation: data[1], Group: data[2], Instance: data[3], Register: binary.LittleEndian.Uint16(data[4:6])}
	if len(s.Sends) == 0 {
		t.Fatal("no send recorded")
	}
	rows := map[string]Row{}
	for _, r := range ExtendedRows() {
		if r.Matches(sel) {
			rows[r.ID] = r
		}
	}
	probe, value := rows["b524.discovery.capability_probe_hot_water_mode"], rows["b524.dhw.dhw_operation_mode"]
	if probe.ID == "" || value.ID == "" {
		t.Fatalf("rows for %x: %v", data, rows)
	}
	probeRec := Record{Selector: sel, Request: data, Reply: boilerRowsBytes(t, s.Reply), Replied: true, Source: s.Source, Target: s.Target, Clock: base.Add(time.Duration(s.Sends[0].AtMs) * time.Millisecond)}
	valueRec := probeRec
	valueRec.Clock = base.Add(time.Duration(s.Sends[s.Probe.Sends].AtMs) * time.Millisecond)
	p := probe.Decode(probeRec, Decode{PurposeDiscoveryProbe, CodecProbe})
	v := value.Decode(valueRec, Decode{PurposeValueRead, CodecUintWidths})
	if !s.Probe.Coherent || p.Outcome != OutcomeValue || p.Value.Text != "true" || !p.AcquiredAt.Equal(probeRec.Clock) {
		t.Fatalf("probe %+v, donor coherent %v", p, s.Probe.Coherent)
	}
	if s.Value.Accepted || s.Value.MismatchLines != 0 || v.Outcome != OutcomeNoValue || v.Reason != ReasonRejected || v.Diagnostics.MismatchText != "" || !v.AcquiredAt.Equal(valueRec.Clock) {
		t.Fatalf("value read %+v, donor %+v", v, s.Value)
	}
	if string(p.Reply) != string(v.Reply) || p.AcquiredAt.Equal(v.AcquiredAt) {
		t.Fatalf("the two outcomes do not keep the same bytes with their own receipts: %+v %+v", p, v)
	}
	st.decodes += 2
}

func replayFlowRetained(t *testing.T, st *replayStats) {
	s := replaySupplementalLoad(t).Flow
	base := time.Date(2026, 1, 5, 6, 0, 0, 0, time.UTC)
	var row BoilerRegisterRow
	for _, r := range BoilerRegisterRows {
		if r.Register == 0x1800 {
			row = r
		}
	}
	flow := row
	flow.Codec, flow.Kind = BoilerCodecSignedFixed, BoilerKindFloat
	var lastGood *BoilerReadResult
	var lastGoodFlow float64
	for i, a := range s.Attempts {
		rec := Record{Selector: row.Selector(), Request: boilerRowsBytes(t, s.Request), Source: s.Source, Target: s.Target, Clock: base.Add(time.Duration(a.AtMs) * time.Millisecond), Diagnostics: Diagnostics{SendError: a.Error}}
		if a.Reply != nil {
			rec.Reply, rec.Replied = boilerRowsBytes(t, *a.Reply), true
		}
		res, err := row.Read(rec)
		if err != nil {
			t.Fatal(err)
		}
		if (res.Outcome == OutcomeValue) != a.Seam.Accepted {
			st.diff(t, "flow attempt %d: outcome %s, donor accepted %v", i, res.Outcome, a.Seam.Accepted)
		}
		if res.Outcome == OutcomeValue {
			decoded, err := flow.Read(rec)
			if err != nil || decoded.Outcome != OutcomeValue {
				t.Fatalf("flow decode of %x: %+v %v", rec.Reply, decoded, err)
			}
			r := res
			lastGood, lastGoodFlow = &r, decoded.Value.(float64)
		} else if res.Value != nil || res.Reason != BoilerReasonSendError {
			st.diff(t, "flow attempt %d: failed attempt %+v carries a value or another reason", i, res)
		}
		st.decodes++
	}
	if lastGood == nil || len(s.Attempts) < 2 || s.FailedReadOK || !s.ReceiptSame || s.ReceiptBefore != s.ReceiptAfter {
		t.Fatalf("retained flow case: last good %+v, donor %+v", lastGood, s)
	}
	if lastGoodFlow != s.FlowFirst || lastGoodFlow != s.FlowAfter || !lastGood.AcquiredAt.Equal(base.Add(time.Duration(s.Attempts[0].AtMs)*time.Millisecond)) ||
		hex.EncodeToString(lastGood.Tail) != s.Attempts[0].Seam.PayloadHex {
		t.Fatalf("retained flow %v at %v, donor flow %v then %v", lastGoodFlow, lastGood.AcquiredAt, s.FlowFirst, s.FlowAfter)
	}
	st.acquisitions += len(s.Attempts)
	st.byFamily["B509"] += len(s.Attempts)
}

func replayB503(t *testing.T, st *replayStats) {
	f := historyRowsLoad(t)
	required := map[string]bool{"current_error_slots": false, "current_service_slots": false, "error_history_index_absent": false, "error_history_index_zero": false,
		"service_history_index_absent": false, "service_history_index_zero": false, "live_monitor_enable": false, "live_monitor_read": false, "live_monitor_disable": false,
		"typed_peer_nak": false, "text_only_nak": false, "pre_cancelled": false, "stale_epoch_after_positive_reply": false}
	for _, c := range f.Cases {
		t.Run(c.Case, func(t *testing.T) { historyRowsCheck(t, c) })
		if _, ok := required[c.Case]; ok {
			required[c.Case] = true
		}
		st.decodes++
	}
	for name, seen := range required {
		if !seen {
			t.Fatalf("the B503 cases lack %s", name)
		}
	}
	actions := []LiveMonitorAction{LiveMonitorEnable, LiveMonitorRead, LiveMonitorDisable}
	for _, c := range f.Cases {
		if c.Qualifier == "" || c.Class != "ack" {
			continue
		}
		body := boilerRowsBytes(t, c.PayloadHex)
		key := "live_monitor_" + c.Qualifier
		if c.Case != key {
			t.Fatalf("case %s has key %s", c.Case, key)
		}
		own := LiveMonitorSelector{Action: LiveMonitorAction(c.Qualifier)}
		if err := own.Resolve(key, LiveMonitorAction(c.Qualifier), body); err != nil {
			t.Fatalf("%s does not resolve its own context: %v", key, err)
		}
		for _, a := range actions {
			if string(a) == c.Qualifier {
				continue
			}
			if err := (LiveMonitorSelector{Action: a}).Resolve(key, LiveMonitorAction(c.Qualifier), body); err != ErrLiveMonitorContext {
				t.Fatalf("sibling action %s with the same bytes resolved against %s: %v", a, key, err)
			}
		}
	}
	st.acquisitions += len(f.Cases)
	st.byFamily["B503"] += len(f.Cases)
}

func replayIdentity(t *testing.T, st *replayStats) {
	f := boilerRowsLoad[identityRowsFile](t, "testdata/vectors/interface_identity.json")
	stored := 0
	for _, c := range f.Cases {
		t.Run(c.Case, func(t *testing.T) {
			if identityRowsCheck(t, c) {
				stored++
			}
		})
		st.decodes += len(c.Requests)
	}
	if stored == 0 {
		t.Fatal("no identity case stored a record")
	}

}
