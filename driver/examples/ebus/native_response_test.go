package ebus

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"
	"sort"
	"testing"

	"github.com/Project-Helianthus/helianthus-ebusgo/protocol"
)

type nativeFrame struct {
	Source, Target, Primary, Secondary byte
	Request                            string
	Present                            bool
	Reply                              *string
	Qualified                          string
}

type nativeResponseVector struct {
	nativeFrame
	Purpose             registerReadPurpose
	OK                  bool
	Tail                string
	Celsius             *json.Number
	Requests, Steps     []nativeFrame
	Attempts            int
	Complete            bool
	Session             uint64
	ObservationMatched  bool `json:"observation_matched"`
	Decoded             json.RawMessage
	Code                *byte
	SendOutcome         protocol.BusOutcomeClass `json:"send_outcome"`
	Error               string
	Frame               *nativeFrame
	Transport           diagnosticStamp
	CompletionTransport diagnosticStamp `json:"completion_transport"`
	Emitted             bool
	Native              diagnosticNative
	SendError           string `json:"send_error"`

	WireReply   string `json:"wire_reply"`
	Response    *string
	DecodeError string `json:"decode_error"`
}

type nativeFixture struct {
	Gateway, MCP json.RawMessage
	Donor        struct {
		Repository          string
		MergeCommitOf       string `json:"merge_commit_of"`
		SourceListingDigest string `json:"source_listing_sha256"`
	}
	SourceDigest map[string]string `json:"source_sha256"`
	Receipts     map[string]struct {
		ExitStatus *int   `json:"exit_status"`
		Digest     string `json:"sha256"`
	}
}

func TestVaillantNativeResponseExemplars(t *testing.T) {
	f := fixtureJSON[nativeFixture](t, "testdata/vaillant_native_response_vectors.json")
	verifyNativeDonor(t, f)
	verifyNativeReceipt(t, f, "gateway", f.Gateway)
	verifyNativeReceipt(t, f, "mcp", f.MCP)
	var vectors map[string]nativeResponseVector
	if err := json.Unmarshal(f.Gateway, &vectors); err != nil {
		t.Fatal(err)
	}
	checkLocalRegisterResponses(t, vectors)
	checkBoilerRegisterResponses(t, vectors)
	checkScheduleResponses(t, vectors)
	checkDiagnosticResponses(t, vectors)
	if len(vectors) != 0 {
		t.Fatalf("native comparator cases not exercised: %v", reflect.ValueOf(vectors).MapKeys())
	}
	checkDiagnosticCallers(t, f.MCP)
}

func verifyNativeReceipt(t *testing.T, f nativeFixture, label string, raw json.RawMessage) {
	t.Helper()
	r := f.Receipts[label]
	var compact bytes.Buffer
	if err := json.Compact(&compact, raw); err != nil {
		t.Fatal(err)
	}
	sum := fmt.Sprintf("%x", sha256.Sum256(compact.Bytes()))
	if len(raw) == 0 || r.ExitStatus == nil || *r.ExitStatus != 0 || r.Digest != sum {
		t.Fatalf("%s: no successful comparator receipt for sha256 %s", label, sum)
	}
	t.Logf("%s comparison record sha256 %s", label, sum)
}

func nativeCase(t *testing.T, vectors map[string]nativeResponseVector, name string, check func(*testing.T, nativeResponseVector)) {
	t.Helper()
	t.Run(name, func(t *testing.T) {
		v, ok := vectors[name]
		if !ok {
			t.Fatal("missing comparator case")
		}
		delete(vectors, name)
		check(t, v)
	})
}

func onSimulatedBus(t *testing.T, frame protocol.Frame, v nativeFrame, present bool) (*protocol.Frame, error) {
	t.Helper()
	if (v.Reply != nil) != present || frame.Source != v.Source || frame.Target != v.Target || frame.Primary != v.Primary || frame.Secondary != v.Secondary || hex.EncodeToString(frame.Data) != v.Request {
		t.Fatalf("request or reply presence differs: constructed %+v, recorded %+v", frame, v)
	}
	var raw []byte
	failure := "timeout"
	if present {
		raw = octets(t, *v.Reply)
		failure = ""
	}
	return simulateExchange(t, frame, v.Source, v.Target, nativeBody(t, v), raw, failure)
}

func nativeBody(t *testing.T, v nativeFrame) []byte {
	t.Helper()
	return append([]byte{v.Primary, v.Secondary}, octets(t, v.Request)...)
}

func errText(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}

func compareDecoded(t *testing.T, got any, expected json.RawMessage) {
	t.Helper()
	if len(expected) == 0 {
		expected = json.RawMessage("null")
	}
	actual, err := json.Marshal(got)
	if err != nil {
		t.Fatal(err)
	}
	var a, w any
	if err := json.Unmarshal(actual, &a); err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(expected, &w); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(a, w) {
		t.Fatalf("decoded %s, recorded %s", actual, expected)
	}
}

func replayWriteAttempts(t *testing.T, frame protocol.Frame, v nativeResponseVector, count int, ack func(*protocol.Frame) error) {
	t.Helper()
	if len(v.Requests) != count {
		t.Fatalf("native writer issued %d frames; want %d", len(v.Requests), count)
	}
	for _, req := range v.Requests {
		if !bytes.Equal(nativeBody(t, req), nativeBody(t, v.nativeFrame)) || req.Source != v.Source || req.Target != v.Target {
			t.Fatalf("native writer changed request while retrying: %+v", req)
		}
		req.Reply = v.Reply
		reply, _ := onSimulatedBus(t, frame, req, v.Present)
		if err := errText(ack(reply)); err != v.Error {
			t.Fatalf("acknowledgement %q, recorded %q", err, v.Error)
		}
	}
}

func verifyNativeDonor(t *testing.T, f nativeFixture) {
	t.Helper()
	const donorMerge = "https://github.com/Project-Helianthus/helianthus-ebusgateway/pull/1738"
	if f.Donor.Repository != "helianthus-ebusgateway" || f.Donor.MergeCommitOf != donorMerge {
		t.Fatalf("donor record does not name the frozen merge: %+v", f.Donor)
	}
	// Local consistency only: the listing digest is recomputed from the
	// fixture's own per-file digests. It does not verify them against the donor.
	if len(f.SourceDigest) == 0 {
		t.Fatal("donor record lists no source files")
	}
	hexDigest := regexp.MustCompile(`^[0-9a-f]{64}$`)
	for p, d := range f.SourceDigest {
		if p == "" || !hexDigest.MatchString(d) {
			t.Fatalf("source digest of %q is not a SHA-256 hex digest: %q", p, d)
		}
	}
	if !hexDigest.MatchString(f.Donor.SourceListingDigest) {
		t.Fatalf("listing digest is not a SHA-256 hex digest: %q", f.Donor.SourceListingDigest)
	}
	paths := make([]string, 0, len(f.SourceDigest))
	for p := range f.SourceDigest {
		paths = append(paths, p)
	}
	sort.Strings(paths)
	h := sha256.New()
	for _, p := range paths {
		h.Write([]byte(f.SourceDigest[p] + "  " + p + "\n"))
	}
	if fmt.Sprintf("%x", h.Sum(nil)) != f.Donor.SourceListingDigest {
		t.Fatalf("donor record does not identify the frozen source set: %+v", f.Donor)
	}
	t.Logf("donor at merge commit of %s; source listing sha256 %s", f.Donor.MergeCommitOf, f.Donor.SourceListingDigest)
}
