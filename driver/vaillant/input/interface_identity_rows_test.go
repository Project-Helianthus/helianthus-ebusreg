package input

import (
	"encoding/hex"
	"strings"
	"testing"
	"time"
)

type identityRowsRequest struct {
	Source, Target, Primary, Secondary byte
	RequestHex                         string  `json:"request_hex"`
	ReplyHex                           *string `json:"reply_hex"`
	SendError                          string  `json:"send_error"`
	Qualifier                          *string `json:"donor_qualifier"`
	QualifierOK                        *bool   `json:"donor_qualifier_ok"`
}

type identityRowsCase struct {
	Case      string                `json:"case"`
	Class     string                `json:"class"`
	Requests  []identityRowsRequest `json:"requests"`
	Stored    bool                  `json:"record_stored"`
	API       string                `json:"api_version"`
	Revision  string                `json:"revision"`
	RawAPI    string                `json:"raw_api_hex"`
	RawRev    string                `json:"raw_revision_hex"`
	Source    byte                  `json:"source"`
	Address   byte                  `json:"address"`
	Session   uint64                `json:"session"`
	NotBefore bool                  `json:"acquired_not_before_controller_observation"`
}

type identityRowsFile struct {
	Format       string            `json:"format"`
	DonorSources map[string]string `json:"donor_sources"`
	Cases        []identityRowsCase
}

// The instants are synthetic inputs for both paths: the donor stamps its record from the wall clock, so the
// vector records only that the stamp was not before the controller observation.
var (
	identityRowsObservedAt = time.Date(2002, 3, 4, 5, 6, 7, 0, time.UTC)
	identityRowsAcquiredAt = identityRowsObservedAt.Add(1500 * time.Millisecond)
)

func identityRowsRecord(t *testing.T, row InterfaceIdentityRow, r identityRowsRequest) Record {
	t.Helper()
	rec := Record{Selector: row.Selector(), Request: boilerRowsBytes(t, r.RequestHex), Source: r.Source, Target: r.Target, Clock: identityRowsAcquiredAt}
	if r.SendError != "" {
		rec.Diagnostics = Diagnostics{Outcome: OutcomeNoValue, SendError: r.SendError}
	}
	if r.ReplyHex != nil {
		rec.Reply, rec.Replied = boilerRowsBytes(t, *r.ReplyHex), true
	}
	return rec
}

func TestVaillantInputRowsInterfaceIdentityInventory(t *testing.T) {
	inv := boilerRowsLoad[boilerRowsInventory](t, "testdata/gateway_selectors.json")
	byID := map[string]InterfaceIdentityRow{}
	for _, r := range InterfaceIdentityRows {
		byID[r.ID] = r
	}
	seen := 0
	for _, s := range inv.Selectors {
		if !strings.HasPrefix(s.ID, "b524.interface_identity.") {
			continue
		}
		r, ok := byID[s.ID]
		if !ok {
			t.Fatalf("inventoried selector %s has no row", s.ID)
		}
		seen++
		if hex.EncodeToString(r.Selector().Payload()) != s.RequestHex || !strings.HasPrefix(s.NativeViewField, r.NativeField) {
			t.Fatalf("row %s differs from its inventory entry %+v", r.ID, s)
		}
		if svc := r.Selector().Service(); svc != (Service{Primary: 0xb5, Secondary: 0x24}) {
			t.Fatalf("row %s service %+v", r.ID, svc)
		}
	}
	if seen != len(InterfaceIdentityRows) || seen != 2 {
		t.Fatalf("%d rows, %d inventoried identity selectors", len(InterfaceIdentityRows), seen)
	}
}

func TestVaillantInputRowsInterfaceIdentityVectors(t *testing.T) {
	f := boilerRowsLoad[identityRowsFile](t, "testdata/vectors/interface_identity.json")
	inv := boilerRowsLoad[boilerRowsInventory](t, "testdata/gateway_selectors.json")
	for _, s := range inv.SourceFiles {
		if want, ok := f.DonorSources[s.Path]; ok && want != s.Digest {
			t.Fatalf("vector donor source %s differs from the inventoried digest", s.Path)
		}
	}
	if len(f.Cases) != 14 {
		t.Fatalf("%d cases", len(f.Cases))
	}
	stored := 0
	for _, c := range f.Cases {
		t.Run(c.Case, func(t *testing.T) {
			if identityRowsCheck(t, c) {
				stored++
			}
		})
	}
	if stored != 3 {
		t.Fatalf("%d cases stored a record", stored)
	}
}

func identityRowsCheck(t *testing.T, c identityRowsCase) bool {
	t.Helper()
	if len(c.Requests) != 2 {
		t.Fatalf("%d requests", len(c.Requests))
	}
	var reads [2]InterfaceIdentityRead
	for i, r := range c.Requests {
		row := InterfaceIdentityRows[i]
		if r.Primary != 0xb5 || r.Secondary != 0x24 || r.RequestHex != hex.EncodeToString(row.Selector().Payload()) {
			t.Fatalf("request %d %+v is not the identity request of row %s", i, r, row.ID)
		}
		got, err := row.Read(identityRowsRecord(t, row, r))
		if err != nil {
			t.Fatal(err)
		}
		want := InterfaceIdentityReasonNone
		switch {
		case r.SendError != "":
			want = InterfaceIdentityReasonSendError
		case r.ReplyHex == nil:
			want = InterfaceIdentityReasonNoReply
		case !*r.QualifierOK:
			want = InterfaceIdentityReasonNotQualified
		case got.Qualifier != *r.Qualifier:
			t.Fatalf("request %d qualifier %q, donor %q", i, got.Qualifier, *r.Qualifier)
		}
		if got.Reason != want || (got.Outcome == OutcomeValue) != (want == InterfaceIdentityReasonNone) || (want != InterfaceIdentityReasonNone && got.Qualifier != "") ||
			!got.AcquiredAt.Equal(identityRowsAcquiredAt) || got.Source != r.Source || got.Target != r.Target {
			t.Fatalf("request %d: %+v, want reason %q", i, got, want)
		}
		if r.ReplyHex != nil && hex.EncodeToString(got.Raw) != *r.ReplyHex {
			t.Fatalf("request %d raw reply %x", i, got.Raw)
		}
		reads[i] = got
	}
	provenance := InterfaceIdentityProvenance{Address: c.Requests[0].Target, Source: c.Requests[0].Source, Session: c.Session, ControllerObservedAt: identityRowsObservedAt, AcquiredAt: identityRowsAcquiredAt}
	record, ok := AssembleInterfaceIdentity(reads[0], reads[1], provenance)
	if ok != c.Stored {
		t.Fatalf("record assembled %v, donor stored %v", ok, c.Stored)
	}
	if !ok {
		return false
	}
	if record.APIVersion != c.API || record.Revision != c.Revision || hex.EncodeToString(record.Raw[0][:]) != c.RawAPI || hex.EncodeToString(record.Raw[1][:]) != c.RawRev ||
		record.Provenance.Source != c.Source || record.Provenance.Address != c.Address || record.Provenance.Session != c.Session {
		t.Fatalf("record %+v, donor %+v", record, c)
	}
	if !c.NotBefore || record.Provenance.AcquiredAt.Before(record.Provenance.ControllerObservedAt) ||
		!record.Provenance.AcquiredAt.Equal(identityRowsAcquiredAt) || !record.Provenance.ControllerObservedAt.Equal(identityRowsObservedAt) {
		t.Fatalf("provenance instants were not passed through unchanged: %+v", record.Provenance)
	}
	return true
}

// TestVaillantInputRowsInterfaceIdentityKey varies every field of the request key and of the stored provenance.
//
//	Identity id:              a row that ignores the id would qualify the revision reply as the API version.
//	Selector type:            a row that ignores the geometry would accept a six-byte register read with the same bytes.
//	Read order:               an assembly that does not check the rows would store the revision as the API version.
//	Address, Source, Session, ControllerObservedAt, AcquiredAt:
//	                          an assembly that derives or drops one of them would store a zero or a clock value.
func TestVaillantInputRowsInterfaceIdentityKey(t *testing.T) {
	api, revision := InterfaceIdentityRows[0], InterfaceIdentityRows[1]
	reply := []byte{0x00, 0x00, 0x80, 0x3f}
	rec := func(sel Selector) Record {
		return Record{Selector: sel, Request: []byte{0, 0, 0}, Reply: reply, Replied: true}
	}
	for name, sel := range map[string]Selector{
		"other id": InterfaceIdentitySelector{ID: revision.Identity}, "id zero": InterfaceIdentitySelector{ID: 0},
		"geometry": B524Selector{Opcode: 0x00, Operation: api.Identity, Group: 0x00},
	} {
		if _, err := api.Read(rec(sel)); err != ErrInterfaceIdentitySelector {
			t.Fatalf("%s: %v", name, err)
		}
	}
	a, _ := api.Read(rec(api.Selector()))
	r, _ := revision.Read(rec(revision.Selector()))
	base := InterfaceIdentityProvenance{Address: 0x15, Source: 0x71, Session: 4, ControllerObservedAt: identityRowsObservedAt, AcquiredAt: identityRowsAcquiredAt}
	if _, ok := AssembleInterfaceIdentity(r, a, base); ok {
		t.Fatal("reads in the wrong order were assembled")
	}
	if _, ok := AssembleInterfaceIdentity(a, a, base); ok {
		t.Fatal("the API read was assembled as the revision")
	}
	moved := InterfaceIdentityProvenance{Address: 0x26, Source: 0x72, Session: 5, ControllerObservedAt: identityRowsObservedAt.Add(time.Second), AcquiredAt: identityRowsAcquiredAt.Add(time.Minute)}
	x, _ := AssembleInterfaceIdentity(a, r, base)
	y, _ := AssembleInterfaceIdentity(a, r, moved)
	if x.Provenance != base || y.Provenance != moved || x.Provenance.Address == y.Provenance.Address || x.Provenance.Source == y.Provenance.Source ||
		x.Provenance.Session == y.Provenance.Session || x.Provenance.ControllerObservedAt.Equal(y.Provenance.ControllerObservedAt) || x.Provenance.AcquiredAt.Equal(y.Provenance.AcquiredAt) {
		t.Fatalf("provenance does not follow its inputs: %+v %+v", x.Provenance, y.Provenance)
	}
}
