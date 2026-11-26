package input

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"testing"
	"time"

	ebuserrors "github.com/Project-Helianthus/helianthus-ebusgo/errors"
)

type historyRowsFrame struct {
	Source, Target, Primary, Secondary byte
	RequestHex                         string `json:"request_hex"`
}

type historyRowsCase struct {
	Case        string             `json:"case"`
	Class       string             `json:"class"`
	Row         string             `json:"row"`
	Qualifier   string             `json:"owner_qualifier"`
	Target      byte               `json:"target"`
	PayloadHex  string             `json:"payload_hex"`
	ReplyHex    *string            `json:"reply_hex"`
	SendError   string             `json:"send_error"`
	TypedNAK    bool               `json:"typed_nak"`
	Frames      []historyRowsFrame `json:"frames"`
	ResponseHex *string            `json:"response_hex"`
	Err         string             `json:"err"`
	ErrKind     string             `json:"err_kind"`
	Emitted     bool               `json:"emitted"`
	Native      int                `json:"native"`
	Transport   struct {
		AdapterInstanceID string `json:"adapter_instance_id"`
		TransportEpoch    uint64 `json:"transport_epoch"`
	} `json:"transport"`
	Decoded     any    `json:"decoded"`
	DecodeError string `json:"decode_error"`
}

type historyRowsFile struct {
	Format       string            `json:"format"`
	DonorSources map[string]string `json:"donor_sources"`
	Cases        []historyRowsCase
}

func historyRowsSelector(t *testing.T, c historyRowsCase, body []byte) Selector {
	t.Helper()
	if c.Qualifier != "" {
		return LiveMonitorSelector{Action: LiveMonitorAction(c.Qualifier)}
	}
	s := DiagnosticSelector{Family: body[0], Selector: body[1]}
	if len(body) == 3 {
		i := body[2]
		s.Index = &i
	}
	return s
}

func historyRowsSendError(t *testing.T, c historyRowsCase) error {
	t.Helper()
	switch {
	case c.SendError == "":
		return nil
	case c.TypedNAK:
		err := fmt.Errorf("target rejected request: %w", ebuserrors.ErrNACK)
		if err.Error() != c.SendError {
			t.Fatalf("typed NACK text %q, recorded %q", err.Error(), c.SendError)
		}
		return err
	}
	return errors.New(c.SendError)
}

func historyRowsDispatch(t *testing.T, c historyRowsCase) DiagnosticDispatch {
	t.Helper()
	issued := TransportStamp{AdapterInstanceID: c.Transport.AdapterInstanceID, TransportEpoch: c.Transport.TransportEpoch}
	d := DiagnosticDispatch{Issued: issued, Current: issued, Emitted: c.Emitted, SendError: historyRowsSendError(t, c), Classified: c.Err}
	if c.ErrKind == "stale_epoch" {
		d.Current.TransportEpoch++
		d.Classified = ""
		if d.SendError != nil {
			d.Classified = "b503dispatcher: upstream rpc failed: " + d.SendError.Error()
		}
	}
	return d
}

func historyRowsNormal(t *testing.T, v any) any {
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

func historyRowsRow(t *testing.T, id string) DiagnosticRow {
	t.Helper()
	for _, r := range DiagnosticRows {
		if r.ID == id {
			return r
		}
	}
	t.Fatalf("no row %s", id)
	return DiagnosticRow{}
}

func historyRowsLoad(t *testing.T) historyRowsFile {
	t.Helper()
	f := boilerRowsLoad[historyRowsFile](t, "testdata/vectors/history.json")
	inv := boilerRowsLoad[boilerRowsInventory](t, "testdata/gateway_selectors.json")
	for _, s := range inv.SourceFiles {
		if want, ok := f.DonorSources[s.Path]; ok && want != s.Digest {
			t.Fatalf("vector donor source %s differs from the inventoried digest", s.Path)
		}
	}
	return f
}

func TestVaillantInputRowsHistoryInventory(t *testing.T) {
	inv := boilerRowsLoad[boilerRowsInventory](t, "testdata/gateway_selectors.json")
	byID := map[string]DiagnosticRow{}
	for _, r := range DiagnosticRows {
		byID[r.ID] = r
	}
	seen := 0
	for _, s := range inv.Selectors {
		if s.Service != "B503" {
			continue
		}
		r, ok := byID[s.ID]
		if !ok {
			t.Fatalf("inventoried selector %s has no row", s.ID)
		}
		seen++
		if hex.EncodeToString([]byte{r.Family, r.Selector}) != s.RequestHex {
			t.Fatalf("row %s body %02x%02x, inventory %s", r.ID, r.Family, r.Selector, s.RequestHex)
		}
	}
	if seen != len(DiagnosticRows) || seen != 5 {
		t.Fatalf("%d rows, %d inventoried B503 selectors: a row has no inventoried selector", len(DiagnosticRows), seen)
	}
}

func TestVaillantInputRowsHistoryVectors(t *testing.T) {
	f := historyRowsLoad(t)
	if len(f.Cases) != 20 {
		t.Fatalf("%d cases", len(f.Cases))
	}
	for _, c := range f.Cases {
		t.Run(c.Case, func(t *testing.T) { historyRowsCheck(t, c) })
	}
}

func historyRowsCheck(t *testing.T, c historyRowsCase) {
	t.Helper()
	body := boilerRowsBytes(t, c.PayloadHex)
	row := historyRowsRow(t, c.Row)
	rec := Record{Selector: historyRowsSelector(t, c, body), Request: body, Source: 0x7f, Target: c.Target, Clock: boilerRowsBase}
	if c.ReplyHex != nil {
		rec.Reply, rec.Replied = boilerRowsBytes(t, *c.ReplyHex), true
	}
	for _, fr := range c.Frames {
		if fr.Primary != 0xb5 || fr.Secondary != 0x03 || fr.RequestHex != c.PayloadHex || fr.Target != c.Target {
			t.Fatalf("donor frame %+v", fr)
		}
		rec.Source = fr.Source
	}
	got, err := row.Observe(rec, historyRowsDispatch(t, c))
	if c.ErrKind == "malformed_payload" {
		if err != ErrDiagnosticSelectorMismatch || c.Emitted || len(c.Frames) != 0 {
			t.Fatalf("a body that starts with the service bytes reached the row: %v", err)
		}
		return
	}
	if err != nil {
		t.Fatal(err)
	}
	if got.Emitted != c.Emitted || int(got.Native) != c.Native || got.Err != c.Err || got.Transport.AdapterInstanceID != c.Transport.AdapterInstanceID ||
		got.Transport.TransportEpoch != c.Transport.TransportEpoch || !got.AcquiredAt.Equal(boilerRowsBase) || got.Target != c.Target || hex.EncodeToString(got.Request) != c.PayloadHex {
		t.Fatalf("native result %+v, donor %+v", got, c)
	}
	if (c.ResponseHex == nil && len(got.Response) != 0) || (c.ResponseHex != nil && hex.EncodeToString(got.Response) != *c.ResponseHex) || (got.Native != DiagnosticACK && got.Response != nil) {
		t.Fatalf("response %x, donor %v", got.Response, c.ResponseHex)
	}
	if got.DecodeError != c.DecodeError || !reflect.DeepEqual(historyRowsNormal(t, got.Decoded), historyRowsNormal(t, c.Decoded)) {
		t.Fatalf("decoded %+v / %q, donor %+v / %q", got.Decoded, got.DecodeError, c.Decoded, c.DecodeError)
	}
}

func TestVaillantInputRowsHistoryIndexAbsentAndZero(t *testing.T) {
	f := historyRowsLoad(t)
	cases := map[string]historyRowsCase{}
	for _, c := range f.Cases {
		cases[c.Case] = c
	}
	for _, name := range []string{"error_history", "service_history"} {
		absent, zero := cases[name+"_index_absent"], cases[name+"_index_zero"]
		row := historyRowsRow(t, absent.Row)
		var results [2]DiagnosticResult
		for i, c := range []historyRowsCase{absent, zero} {
			body := boilerRowsBytes(t, c.PayloadHex)
			rec := Record{Selector: historyRowsSelector(t, c, body), Request: body, Reply: boilerRowsBytes(t, *c.ReplyHex), Replied: true}
			r, err := row.Observe(rec, historyRowsDispatch(t, c))
			if err != nil {
				t.Fatal(err)
			}
			results[i] = r
		}
		if results[0].Index != nil || results[1].Index == nil || *results[1].Index != 0 ||
			len(results[0].Request) != 2 || len(results[1].Request) != 3 || results[0].Native != DiagnosticACK || results[1].Native != DiagnosticACK {
			t.Fatalf("%s: absent %+v, zero %+v", name, results[0], results[1])
		}
	}
}

func TestVaillantInputRowsHistoryTypedAndTextOnlyNAK(t *testing.T) {
	f := historyRowsLoad(t)
	cases := map[string]historyRowsCase{}
	for _, c := range f.Cases {
		cases[c.Case] = c
	}
	typed, text := cases["typed_peer_nak"], cases["text_only_nak"]
	if !typed.TypedNAK || text.TypedNAK || !strings.Contains(text.SendError, "NAK") || strings.Contains(text.SendError, "NACK") && !strings.Contains(text.SendError, "text only") {
		t.Fatalf("vector pair is not typed against text-only: %+v / %+v", typed, text)
	}
	row := historyRowsRow(t, typed.Row)
	results := map[bool]DiagnosticResult{}
	for _, c := range []historyRowsCase{typed, text} {
		r, err := row.Observe(Record{Selector: LiveMonitorSelector{Action: LiveMonitorAction(c.Qualifier)}, Request: []byte{0, 3}}, historyRowsDispatch(t, c))
		if err != nil {
			t.Fatal(err)
		}
		results[c.TypedNAK] = r
	}
	if results[true].Native != DiagnosticNAK || results[false].Native != DiagnosticAmbiguous || !results[true].Emitted || !results[false].Emitted ||
		results[true].Response != nil || results[false].Response != nil || results[true].Err == results[false].Err {
		t.Fatalf("typed %+v, text-only %+v", results[true], results[false])
	}
}

func TestVaillantInputRowsHistoryStaleEpochKeepsOriginalTransport(t *testing.T) {
	f := historyRowsLoad(t)
	row := historyRowsRow(t, "b503.current_error_slots")
	var positive historyRowsCase
	for _, c := range f.Cases {
		if c.Case == "stale_epoch_after_positive_reply" {
			positive = c
		}
	}
	issued := TransportStamp{AdapterInstanceID: positive.Transport.AdapterInstanceID, TransportEpoch: positive.Transport.TransportEpoch}
	for name, current := range map[string]TransportStamp{
		"epoch advanced":    {AdapterInstanceID: issued.AdapterInstanceID, TransportEpoch: issued.TransportEpoch + 1},
		"epoch rolled back": {AdapterInstanceID: issued.AdapterInstanceID, TransportEpoch: issued.TransportEpoch - 1},
	} {
		body := boilerRowsBytes(t, positive.PayloadHex)
		rec := Record{Selector: DiagnosticSelector{Family: body[0], Selector: body[1]}, Request: body, Reply: boilerRowsBytes(t, *positive.ReplyHex), Replied: true}
		r, err := row.Observe(rec, DiagnosticDispatch{Issued: issued, Current: current, Emitted: true})
		if err != nil || r.Native != DiagnosticAmbiguous || r.Response != nil || !r.Emitted || r.Transport != issued || r.Err != positive.Err || r.Decoded != nil {
			t.Fatalf("%s: %+v, %v", name, r, err)
		}
	}
	// The dispatcher compares the epoch only: an adapter instance replaced at the same epoch keeps the reply.
	for name, current := range map[string]TransportStamp{
		"unchanged transport":       issued,
		"adapter instance replaced": {AdapterInstanceID: issued.AdapterInstanceID + "-replaced", TransportEpoch: issued.TransportEpoch},
	} {
		same, err := row.Observe(Record{Selector: DiagnosticSelector{Family: 0, Selector: 1}, Reply: []byte{0x19, 0x01, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}, Replied: true}, DiagnosticDispatch{Issued: issued, Current: current, Emitted: true})
		if err != nil || same.Native != DiagnosticACK || same.Response == nil || same.Err != "" || same.Transport != issued {
			t.Fatalf("%s: %+v, %v", name, same, err)
		}
	}
}

func TestVaillantInputRowsHistoryLiveMonitorContext(t *testing.T) {
	f := historyRowsLoad(t)
	actions := []LiveMonitorAction{LiveMonitorEnable, LiveMonitorRead, LiveMonitorDisable}
	checked := 0
	for _, c := range f.Cases {
		if c.Qualifier == "" || c.Class != "ack" {
			continue
		}
		body := boilerRowsBytes(t, c.PayloadHex)
		key := "live_monitor_" + c.Qualifier
		if c.Case != key {
			t.Fatalf("case %s has key %s", c.Case, key)
		}
		for _, a := range actions {
			err := LiveMonitorSelector{Action: a}.Resolve(key, LiveMonitorAction(c.Qualifier), body)
			if (err == nil) != (string(a) == c.Qualifier) {
				t.Fatalf("selector %s against context %s: %v", a, key, err)
			}
		}
		own := LiveMonitorSelector{Action: LiveMonitorAction(c.Qualifier)}
		for name, err := range map[string]error{
			"sibling key":       own.Resolve("live_monitor_"+string(actions[(checked+1)%3]), LiveMonitorAction(c.Qualifier), body),
			"sibling qualifier": own.Resolve(key, actions[(checked+1)%3], body),
			"other body":        own.Resolve(key, LiveMonitorAction(c.Qualifier), []byte{0x00, 0x02}),
			"index appended":    own.Resolve(key, LiveMonitorAction(c.Qualifier), append(append([]byte(nil), body...), 0x00)),
		} {
			if err != ErrLiveMonitorContext {
				t.Fatalf("%s resolved against %s: %v", name, key, err)
			}
		}
		checked++
	}
	if checked != 3 {
		t.Fatalf("%d live-monitor cases", checked)
	}
}

// TestVaillantInputRowsHistoryKey varies every field of the request key and of the transport stamp.
//
//	Family:                    a row that ignores the family would accept a history body for a current-slot row.
//	Selector:                  a row that ignores the selector would accept the service body for the error row.
//	Index on a current row:    a row that accepts an index everywhere would accept 00 01 00.
//	Selector type:             a row that ignores the geometry would accept a live-monitor selector for a history row.
//	Epoch:                     a row that ignores the epoch would keep a reply that completed after the epoch advanced.
//	Adapter instance:          a row that compares it would discard a reply the dispatcher keeps at an unchanged epoch.
//	Emitted:                   a row that ignores it would classify a request that was never sent.
func TestVaillantInputRowsHistoryKey(t *testing.T) {
	zero := byte(0)
	for _, r := range DiagnosticRows {
		if r.LiveMonitor {
			continue
		}
		good := DiagnosticSelector{Family: r.Family, Selector: r.Selector}
		if _, err := r.Observe(Record{Selector: good}, DiagnosticDispatch{}); err != nil {
			t.Fatalf("%s: %v", r.ID, err)
		}
		mutations := map[string]Selector{
			"family": DiagnosticSelector{Family: r.Family ^ 1, Selector: r.Selector}, "selector": DiagnosticSelector{Family: r.Family, Selector: r.Selector ^ 3},
			"geometry": LiveMonitorSelector{Action: LiveMonitorRead},
		}
		if !r.Indexed {
			mutations["index on a row without one"] = DiagnosticSelector{Family: r.Family, Selector: r.Selector, Index: &zero}
		}
		for name, sel := range mutations {
			if _, err := r.Observe(Record{Selector: sel}, DiagnosticDispatch{}); err != ErrDiagnosticSelectorMismatch {
				t.Fatalf("%s %s: %v", r.ID, name, err)
			}
		}
	}
	row := historyRowsRow(t, "b503.current_error_slots")
	reply := []byte{0x19, 0x01, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}
	base := TransportStamp{AdapterInstanceID: "a", TransportEpoch: 7}
	rec := Record{Selector: DiagnosticSelector{Family: 0, Selector: 1}, Reply: reply, Replied: true, Clock: time.Unix(5, 0)}
	for name, d := range map[string]DiagnosticDispatch{
		"epoch":   {Issued: base, Current: TransportStamp{AdapterInstanceID: "a", TransportEpoch: 8}, Emitted: true},
		"emitted": {Issued: base, Current: base},
	} {
		if r, _ := row.Observe(rec, d); r.Native == DiagnosticACK || r.Response != nil {
			t.Fatalf("%s: %+v", name, r)
		}
	}
	if r, _ := row.Observe(rec, DiagnosticDispatch{Issued: base, Current: TransportStamp{AdapterInstanceID: "b", TransportEpoch: 7}, Emitted: true}); r.Native != DiagnosticACK || r.Response == nil {
		t.Fatalf("instance replaced at the same epoch: %+v", r)
	}
	if r, _ := row.Observe(rec, DiagnosticDispatch{Issued: base, Current: base, Emitted: true}); r.Native != DiagnosticACK || !r.AcquiredAt.Equal(rec.Clock) {
		t.Fatalf("matching stamp: %+v", r)
	}
}
