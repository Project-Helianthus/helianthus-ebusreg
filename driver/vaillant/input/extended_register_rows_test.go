package input

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/Project-Helianthus/helianthus-ebusreg/vaillant/nativepoints"
)

type vectorValue struct{ Kind, Text string }

type vectorExpect struct {
	Outcome    string       `json:"outcome"`
	Stage      string       `json:"stage"`
	Value      *vectorValue `json:"value"`
	Raw        string       `json:"raw"`
	Request    string       `json:"request"`
	Mismatch   string       `json:"mismatch"`
	MismatchN  int          `json:"mismatch_lines"`
	WidthLog   string       `json:"width_log"`
	Sends      int          `json:"sends"`
	DonorError string       `json:"donor_error"`
}

type vector struct {
	Row       string       `json:"row"`
	Purpose   string       `json:"purpose"`
	Codec     string       `json:"codec"`
	Case      string       `json:"case"`
	Instance  int          `json:"instance"`
	Reply     *string      `json:"reply"`
	SendError string       `json:"send_error"`
	WriteData string       `json:"write_value"`
	Expect    vectorExpect `json:"expect"`
}

type vectorFile struct {
	Clock       string `json:"clock"`
	FirstSendAt string `json:"first_send_at"`
	SourceFiles []struct {
		Path   string `json:"path"`
		Digest string `json:"sha256"`
	} `json:"source_files"`
	Vectors []vector `json:"vectors"`
}

func readJSON(t *testing.T, path string, into any) {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	if err := json.Unmarshal(data, into); err != nil {
		t.Fatalf("decode %s: %v", path, err)
	}
}

func loadVectors(t *testing.T) vectorFile {
	t.Helper()
	var file vectorFile
	readJSON(t, "testdata/vectors/extended_register.json", &file)
	if len(file.Vectors) == 0 {
		t.Fatal("the vector file holds no vectors")
	}
	return file
}

func rowsByID(t *testing.T) map[string]Row {
	t.Helper()
	rows := map[string]Row{}
	for _, row := range ExtendedRows() {
		if _, dup := rows[row.ID]; dup {
			t.Fatalf("duplicate row %s", row.ID)
		}
		rows[row.ID] = row
	}
	return rows
}

func mustHex(t *testing.T, s string) []byte {
	t.Helper()
	b, err := hex.DecodeString(s)
	if err != nil {
		t.Fatal(err)
	}
	return b
}

func stageOf(res Result) string {
	switch res.Reason {
	case "":
		return "value"
	case ReasonNoReply:
		return "no_reply"
	case ReasonRejected:
		return "reply_rejected"
	case ReasonWidth:
		return "width_mismatch"
	case ReasonNotSent:
		return "not_sent"
	}
	return "codec_no_value"
}

func selectorOf(row Row, instance int) B524Selector {
	return B524Selector{Opcode: row.Opcode, Operation: row.Operation, Group: row.Group, Instance: byte(instance), Register: row.Register}
}

func decodeOf(t *testing.T, row Row, purpose, codec string) Decode {
	t.Helper()
	for _, dec := range row.Decodes {
		if string(dec.Purpose) == purpose && string(dec.Codec) == codec {
			return dec
		}
	}
	t.Fatalf("row %s has no decode %s/%s", row.ID, purpose, codec)
	return Decode{}
}

func recordOf(t *testing.T, row Row, v vector, clock time.Time) Record {
	t.Helper()
	sel := selectorOf(row, v.Instance)
	rec := Record{Selector: sel, Request: sel.Payload(), Source: 0x7f, Target: 0x15, Clock: clock,
		Diagnostics: Diagnostics{SendError: v.SendError}}
	if v.WriteData != "" {
		rec.Request = sel.WithValue(mustHex(t, v.WriteData))
	}
	if v.Reply != nil {
		rec.Reply, rec.Replied = mustHex(t, *v.Reply), true
	}
	return rec
}

func widthLine(res Result) string {
	if res.Reason != ReasonWidth {
		return ""
	}
	return fmt.Sprintf("semantic: b524 reply width mismatch opcode=0x%02X group=0x%02X register=0x%04X want=%d got=%d",
		res.Selector.Opcode, res.Selector.Group, res.Selector.Register, res.WidthWant, len(res.Raw))
}

func TestVaillantInputRowsVectors(t *testing.T) {
	file := loadVectors(t)
	rows := rowsByID(t)
	firstSend, err := time.Parse(time.RFC3339Nano, file.FirstSendAt)
	if err != nil {
		t.Fatal(err)
	}
	clock, err := time.Parse(time.RFC3339Nano, file.Clock)
	if err != nil {
		t.Fatal(err)
	}
	for _, v := range file.Vectors {
		row, found := rows[v.Row]
		if !found {
			t.Fatalf("vector %s/%s names no row", v.Row, v.Case)
		}
		name := v.Row + "/" + v.Purpose + "/" + v.Codec + "/" + v.Case
		rec := recordOf(t, row, v, clock)
		res := row.Decode(rec, decodeOf(t, row, v.Purpose, v.Codec))
		compareVector(t, name, v, rec, res, firstSend)
	}
}

func compareVector(t *testing.T, name string, v vector, rec Record, res Result, firstSend time.Time) {
	t.Helper()
	ex := v.Expect
	if ex.Sends > 0 && !res.AcquiredAt.Equal(firstSend) {
		t.Errorf("%s: acquired at %s; donor first send %s", name, res.AcquiredAt, firstSend)
	}
	if string(rec.Request) != string(mustHex(t, ex.Request)) && ex.Sends > 0 {
		t.Errorf("%s: row request %x; donor sent %s", name, rec.Request, ex.Request)
	}
	if res.Diagnostics.SendError != v.SendError {
		t.Errorf("%s: send error %q; want %q", name, res.Diagnostics.SendError, v.SendError)
	}
	if string(res.Reply) != string(rec.Reply) {
		t.Errorf("%s: raw reply evidence %x; want %x", name, res.Reply, rec.Reply)
	}
	if v.Purpose == string(PurposeDiscoveryProbe) {
		coherent := res.Outcome == OutcomeValue && res.Value.Text == "true"
		if want := ex.Value.Text == "true"; coherent != want {
			t.Errorf("%s: coherent %v; donor %v", name, coherent, want)
		}
		return
	}
	if got := stageOf(res); got != ex.Stage {
		t.Errorf("%s: stage %s (%s); donor %s", name, got, res.Reason, ex.Stage)
	}
	if want, named := reasonOfCase(v.Case); named && ex.Stage == "codec_no_value" && res.Reason != want {
		t.Errorf("%s: reason %s; the case names %s", name, res.Reason, want)
	}
	if string(res.Outcome) != ex.Outcome || string(res.Diagnostics.Outcome) != ex.Outcome {
		t.Errorf("%s: outcome %s/%s; donor %s", name, res.Outcome, res.Diagnostics.Outcome, ex.Outcome)
	}
	if ex.Value != nil && (res.Value.Kind != ex.Value.Kind || res.Value.Text != ex.Value.Text) {
		t.Errorf("%s: value %+v; donor %+v", name, res.Value, *ex.Value)
	}
	if ex.Value == nil && res.Value != (Value{}) {
		t.Errorf("%s: value %+v where the donor produced none", name, res.Value)
	}
	if hex.EncodeToString(res.Raw) != ex.Raw {
		t.Errorf("%s: payload %x; donor %s", name, res.Raw, ex.Raw)
	}
	if res.Diagnostics.MismatchText != ex.Mismatch {
		t.Errorf("%s: mismatch line %q; donor %q", name, res.Diagnostics.MismatchText, ex.Mismatch)
	}
	if ex.Mismatch != "" && ex.MismatchN != ex.Sends {
		t.Errorf("%s: donor logged %d mismatch lines for %d sends", name, ex.MismatchN, ex.Sends)
	}
	if widthLine(res) != ex.WidthLog {
		t.Errorf("%s: width line %q; donor %q", name, widthLine(res), ex.WidthLog)
	}
	if res.Unit != "" && res.Unit != unitOfLocation(res.Selector) {
		t.Errorf("%s: unit %q is not the catalog unit", name, res.Unit)
	}
}

func reasonOfCase(name string) (NoValueReason, bool) {
	switch {
	case name == "sentinel_nan" || name == "sentinel_infinity":
		return ReasonNonFinite, true
	case strings.HasPrefix(name, "sentinel"):
		return ReasonSentinel, true
	case strings.HasPrefix(name, "invalid"):
		return ReasonInvalid, true
	case strings.HasPrefix(name, "out_of_range"):
		return ReasonRange, true
	case name == "short":
		return ReasonShort, true
	}
	return "", false
}

func unitOfLocation(s B524Selector) string {
	for _, point := range nativepoints.Points() {
		if point.DeviceClass != nativepoints.DeviceClassController {
			continue
		}
		for _, loc := range point.Locations {
			if loc.Opcode == s.Opcode && loc.Group == s.Group && loc.Register == s.Register {
				return loc.NativeUnit
			}
		}
	}
	return ""
}

func TestVaillantInputRowsInventory(t *testing.T) {
	inventory := loadInventory(t)
	rows := rowsByID(t)
	seen := map[string]bool{}
	for _, entry := range inventory.Selectors {
		if entry.Geometry != "B524Selector" {
			continue
		}
		row, found := rows[entry.ID]
		if !found {
			t.Errorf("inventory selector %s has no row", entry.ID)
			continue
		}
		seen[entry.ID] = true
		if row.Opcode != parseByte(t, entry.Opcode) || row.Operation != parseByte(t, entry.OperationByte) || row.Group != parseByte(t, entry.Group) ||
			row.FirstInstance != parseByte(t, entry.InstanceFirst) || row.LastInstance != parseByte(t, entry.InstanceLast) || row.Register != parseWord(t, entry.Register) {
			t.Errorf("row %s differs from its inventory selector", entry.ID)
		}
		have := map[string]bool{}
		for _, dec := range row.Decodes {
			have[string(dec.Purpose)] = true
		}
		for _, purpose := range entry.Purposes {
			if purpose == string(PurposeDiscoveryProbe) && entry.Operation == "read" {
				continue
			}
			if !have[purpose] {
				t.Errorf("row %s lacks inventory purpose %s", entry.ID, purpose)
			}
		}
	}
	for id := range rows {
		if !seen[id] {
			t.Errorf("row %s has no inventoried selector", id)
		}
	}
	if len(rows) != 127 {
		t.Errorf("rows = %d; want 127 extended register inventory selectors", len(rows))
	}
}

func TestVaillantInputRowsDonorFiles(t *testing.T) {
	inventory := loadInventory(t)
	digests := map[string]string{}
	for _, f := range inventory.SourceFiles {
		digests[f.Path] = f.Digest
	}
	file := loadVectors(t)
	if len(file.SourceFiles) == 0 {
		t.Fatal("the vector file records no donor source files")
	}
	for _, f := range file.SourceFiles {
		if digests[f.Path] != f.Digest {
			t.Errorf("donor file %s digest %s differs from the inventory digest %q", f.Path, f.Digest, digests[f.Path])
		}
	}
}

func TestVaillantInputRowsCoverage(t *testing.T) {
	file := loadVectors(t)
	rows := rowsByID(t)
	cases := map[string]map[string]bool{}
	for _, v := range file.Vectors {
		key := v.Row + "|" + v.Purpose + "|" + v.Codec
		if cases[key] == nil {
			cases[key] = map[string]bool{}
		}
		if cases[key][v.Case] {
			t.Errorf("duplicate vector %s %s", key, v.Case)
		}
		cases[key][v.Case] = true
	}
	edge := []string{"boundary", "sentinel", "invalid", "out_of_range", "unchecked", "short", "wide", "other", "high_byte", "only_spaces", "leading", "wrong_"}
	for id, row := range rows {
		for _, dec := range row.Decodes {
			have := cases[id+"|"+string(dec.Purpose)+"|"+string(dec.Codec)]
			if len(have) == 0 {
				t.Errorf("row %s decode %s/%s has no vectors", id, dec.Purpose, dec.Codec)
				continue
			}
			if dec.Purpose == PurposeWrite || dec.Purpose == PurposeNotSent || dec.Purpose == PurposeDiscoveryProbe {
				continue
			}
			valid, edged := false, false
			for name := range have {
				valid = valid || strings.HasPrefix(name, "valid")
				for _, prefix := range edge {
					edged = edged || strings.HasPrefix(name, prefix)
				}
			}
			if !valid || !edged {
				t.Errorf("row %s decode %s/%s lacks a valid or an edge case", id, dec.Purpose, dec.Codec)
			}
		}
		if row.Decodes[0].Purpose != PurposeValueRead {
			continue
		}
		primary := cases[id+"|"+string(row.Decodes[0].Purpose)+"|"+string(row.Decodes[0].Codec)]
		for _, name := range []string{"neighbour_register", "neighbour_group", "neighbour_instance", "instance_plus_one_accepted", "header_only", "status_zero", "short_frame", "no_reply"} {
			if !primary[name] {
				t.Errorf("row %s lacks the %s vector", id, name)
			}
		}
	}
	for key := range cases {
		parts := strings.Split(key, "|")
		found := false
		for _, dec := range rows[parts[0]].Decodes {
			found = found || (string(dec.Purpose) == parts[1] && string(dec.Codec) == parts[2])
		}
		if !found {
			t.Errorf("vectors %s match no row decode", key)
		}
	}
}

func vectorFor(t *testing.T, file vectorFile, row, purpose, codec, name string) vector {
	t.Helper()
	for _, v := range file.Vectors {
		if v.Row == row && v.Purpose == purpose && v.Codec == codec && v.Case == name {
			return v
		}
	}
	t.Fatalf("no vector %s/%s/%s/%s", row, purpose, codec, name)
	return vector{}
}

func TestVaillantInputRowsNeighbourAndMiss(t *testing.T) {
	file := loadVectors(t)
	rows := rowsByID(t)
	clock, _ := time.Parse(time.RFC3339Nano, file.Clock)
	mismatchCases := map[string]bool{"neighbour_register": true, "neighbour_group": true, "neighbour_instance": true}
	silentCases := map[string]bool{"header_only": true, "header_only_short_layout": true, "status_zero": true, "short_frame": true, "empty_frame": true}
	checked := map[string]int{}
	for _, v := range file.Vectors {
		isMismatch, isSilent := mismatchCases[v.Case], silentCases[v.Case]
		if (!isMismatch && !isSilent) || v.Purpose == string(PurposeDiscoveryProbe) {
			continue
		}
		row := rows[v.Row]
		res := row.Decode(recordOf(t, row, v, clock), decodeOf(t, row, v.Purpose, v.Codec))
		if res.Outcome != OutcomeNoValue || res.Reason != ReasonRejected || res.Raw != nil {
			t.Errorf("%s/%s: outcome %s reason %s raw %x; want a rejected no-value reply", v.Row, v.Case, res.Outcome, res.Reason, res.Raw)
		}
		if isMismatch && (res.Diagnostics.MismatchText == "" || res.Diagnostics.MismatchText != v.Expect.Mismatch || v.Expect.Mismatch == "") {
			t.Errorf("%s/%s: mismatch line %q; donor %q", v.Row, v.Case, res.Diagnostics.MismatchText, v.Expect.Mismatch)
		}
		if isSilent && (res.Diagnostics.MismatchText != "" || v.Expect.Mismatch != "") {
			t.Errorf("%s/%s: a reply the donor rejects silently produced a line", v.Row, v.Case)
		}
		checked[v.Case]++
	}
	for name := range mismatchCases {
		if checked[name] == 0 {
			t.Errorf("no %s vector was replayed", name)
		}
	}
	accepted := vectorFor(t, file, "b524.zone.zone_target_temp", "value_read", "float32_le_exact_width", "instance_plus_one_accepted")
	row := rows[accepted.Row]
	res := row.Decode(recordOf(t, row, accepted, clock), decodeOf(t, row, accepted.Purpose, accepted.Codec))
	if res.Outcome != OutcomeValue || res.Diagnostics.MismatchText != "" || string(res.Raw) != "\x01\x02\x03\x04" || res.Value.Text != accepted.Expect.Value.Text {
		t.Errorf("a reply that names the next instance is accepted by the donor; got %+v", res)
	}
}

func TestVaillantInputRowsProbeAndRead(t *testing.T) {
	file := loadVectors(t)
	rows := rowsByID(t)
	clock, _ := time.Parse(time.RFC3339Nano, file.Clock)
	probe := vectorFor(t, file, "b524.discovery.capability_probe_hot_water_mode", "discovery_probe", "reply_header_coherence", "coherent_header_only")
	var read vector
	for _, v := range file.Vectors {
		if v.Row == "b524.dhw.dhw_operation_mode" && v.Case == "header_only" && v.Purpose == "value_read" {
			read = v
		}
	}
	if read.Reply == nil || probe.Reply == nil || *read.Reply != *probe.Reply || read.Instance != probe.Instance {
		t.Fatalf("the probe and the read vectors do not carry the same bytes: %+v %+v", probe.Reply, read.Reply)
	}
	probeRow, readRow := rows[probe.Row], rows[read.Row]
	probeRes := probeRow.Decode(recordOf(t, probeRow, probe, clock), probeRow.Decodes[0])
	readRes := readRow.Decode(recordOf(t, readRow, read, clock), readRow.Decodes[0])
	if probeRes.Outcome != OutcomeValue || probeRes.Value.Text != "true" || probe.Expect.Value.Text != "true" {
		t.Errorf("the probe of a header-only reply is %+v; the donor finds it coherent", probeRes.Value)
	}
	if readRes.Outcome != OutcomeNoValue || read.Expect.Outcome != "no_value" || readRes.Value != (Value{}) {
		t.Errorf("the value read of the same bytes is %+v; the donor reads no measured value", readRes.Value)
	}
	if probeRes.Selector.Payload()[0] != OpcodeLocal || readRes.Selector.Register != probeRes.Selector.Register || readRes.Selector.Group != probeRes.Selector.Group {
		t.Errorf("the probe and the read name different selectors: %+v %+v", probeRes.Selector, readRes.Selector)
	}
}

func TestVaillantInputRowsIdentityFields(t *testing.T) {
	file := loadVectors(t)
	rows := rowsByID(t)
	clock, _ := time.Parse(time.RFC3339Nano, file.Clock)
	for _, id := range []string{"b524.zone.zone_target_temp", "b524.regulator.system_water_pressure", "b524.remote.remoteThermostats.device_slot_room_humidity", "b524.dhw.dhw_target_temp", "b524.cylinder.cylinder_temperature"} {
		row := rows[id]
		valid := vectorFor(t, file, id, "value_read", string(row.Decodes[0].Codec), "valid")
		rec := recordOf(t, row, valid, clock)
		sel := rec.Selector.(B524Selector)
		if res := row.Decode(rec, row.Decodes[0]); res.Reason != "" {
			t.Fatalf("%s: the unvaried selector reads no value: %s", id, res.Reason)
		}
		for _, instance := range []byte{row.FirstInstance, row.LastInstance} {
			s := sel
			s.Instance = instance
			if !row.Matches(s) {
				t.Errorf("%s: instance %d inside the range is refused", id, instance)
			}
		}
		mutations := map[string]B524Selector{}
		add := func(name string, mutate func(s *B524Selector)) {
			s := sel
			mutate(&s)
			mutations[name] = s
		}
		add("opcode", func(s *B524Selector) { s.Opcode ^= 0x04 })
		add("operation", func(s *B524Selector) { s.Operation ^= 0x01 })
		add("group", func(s *B524Selector) { s.Group++ })
		add("register", func(s *B524Selector) { s.Register++ })
		add("instance above range", func(s *B524Selector) { s.Instance = row.LastInstance + 1 })
		for name, s := range mutations {
			rec2 := rec
			rec2.Selector = s
			if row.Matches(s) {
				t.Errorf("%s: the %s mutation still matches the row", id, name)
			}
			if res := row.Decode(rec2, row.Decodes[0]); res.Reason != ReasonOutsideRow || res.Outcome != OutcomeNoValue || res.Value != (Value{}) {
				t.Errorf("%s: the %s mutation decoded as %+v", id, name, res)
			}
		}
		rec3 := rec
		rec3.Selector = B509Selector{Opcode: sel.Opcode, Register: sel.Register}
		if res := row.Decode(rec3, row.Decodes[0]); res.Reason != ReasonOutsideRow {
			t.Errorf("%s: a boiler selector decoded as %+v", id, res)
		}
	}
	for id, row := range rows {
		if row.FirstInstance == 0 && row.LastInstance < 0xff {
			beyond := selectorOf(row, int(row.LastInstance)+1)
			if row.Matches(beyond) {
				t.Errorf("%s: the instance after the range matches", id)
			}
		}
	}
}

var catalogCodecEquivalents = map[Codec]nativepoints.Codec{
	CodecFloatStrict: nativepoints.CodecFloat32LE, CodecFloatRaw: nativepoints.CodecFloat32LE,
	CodecUintWidths: nativepoints.CodecUint8OrUint16LE, CodecUintWidthsSent: nativepoints.CodecUint8OrUint16LE, CodecUintRaw: nativepoints.CodecUint8OrUint16LE,
	CodecByteStrict: nativepoints.CodecUint8, CodecByteRaw: nativepoints.CodecUint8, CodecByteClass: nativepoints.CodecUint8,
	CodecCounterStrict: nativepoints.CodecUint32LE, CodecCounterRaw: nativepoints.CodecUint32LE,
	CodecVersion: nativepoints.CodecMajorMinorPatch, CodecTextSanitized: nativepoints.CodecText, CodecTextTrimmed: nativepoints.CodecText,
}

func TestVaillantInputRowsCatalog(t *testing.T) {
	rows := rowsByID(t)
	referenced, fromCatalog, fromGateway, unreferenced := 0, 0, 0, []string{}
	for id, row := range rows {
		ref := row.Catalog()
		if ref.Name == "" {
			unreferenced = append(unreferenced, id)
			continue
		}
		referenced++
		point, found := nativepoints.Lookup(nativepoints.DeviceClassController, ref.Name)
		if !found {
			t.Errorf("%s: catalog point %s does not exist", id, ref.Name)
			continue
		}
		located := false
		for _, loc := range point.Locations {
			if loc.Opcode == row.Opcode && loc.Group == row.Group && loc.Register == row.Register && loc.NativeUnit == ref.Unit && loc.Codec == ref.Codec {
				located = true
			}
		}
		if !located {
			t.Errorf("%s: point %s has no location with the row's selector, unit and codec", id, ref.Name)
		}
		if want, known := catalogCodecEquivalents[row.Decodes[0].Codec]; known && want == ref.Codec {
			fromCatalog++
		} else {
			fromGateway++
		}
	}
	t.Logf("rows=%d referencing the catalog=%d codec equal to the catalog's=%d codec stated from the gateway decoder=%d without catalog metadata=%v", len(rows), referenced, fromCatalog, fromGateway, unreferenced)
	byID := func(id string) Row { return rows[id] }
	water := byID("b524.regulator.system_water_pressure").Catalog()
	if water.Name != "water_pressure" || water.Unit != "bar" || water.Codec != nativepoints.CodecFloat32LE {
		t.Errorf("system water pressure catalog reference = %+v", water)
	}
	if off := byID("b524.regulator.system_off").Catalog(); off.Codec != nativepoints.CodecUint8 {
		t.Errorf("system off catalog codec = %q; the poller reads it with a different codec", off.Codec)
	}
	if byID("b524.regulator.system_off").Decodes[0].Codec != CodecUintWidths {
		t.Error("system off keeps the codec of the gateway call site")
	}
	for _, id := range unreferenced {
		if !strings.HasPrefix(id, "b524.remote.remoteFunctionalModules.") && !strings.HasPrefix(id, "b524.discovery.") && !strings.Contains(id, ".write.") {
			t.Errorf("%s has no catalog metadata", id)
		}
	}
}
