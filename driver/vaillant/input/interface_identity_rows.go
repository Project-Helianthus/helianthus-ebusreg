package input

import (
	"encoding/binary"
	"errors"
	"math"
	"strconv"
	"time"
)

// ErrInterfaceIdentitySelector is returned when a record was acquired for another identity request than the row's.
var ErrInterfaceIdentitySelector = errors.New("record selector is not the identity request of the row")

// InterfaceIdentityReason is why an identity read produced no qualifier.
type InterfaceIdentityReason string

const (
	// InterfaceIdentityReasonNone is the reason of a read that produced a qualifier.
	InterfaceIdentityReasonNone InterfaceIdentityReason = ""
	// InterfaceIdentityReasonSendError is a send that returned an error.
	InterfaceIdentityReasonSendError InterfaceIdentityReason = "send_error"
	// InterfaceIdentityReasonNoReply is a send that returned neither a frame nor an error.
	InterfaceIdentityReasonNoReply InterfaceIdentityReason = "no_reply"
	// InterfaceIdentityReasonNotQualified is a reply that is not a four-byte whole number in the admitted range.
	InterfaceIdentityReasonNotQualified InterfaceIdentityReason = "not_qualified"
)

// InterfaceIdentityRow is one three-byte interface identity request: the identity id and the field of the stored record it fills.
type InterfaceIdentityRow struct {
	ID          string
	Identity    byte
	NativeField string
}

// InterfaceIdentityRows lists the two identity requests the poller sends, in the order it sends them.
var InterfaceIdentityRows = []InterfaceIdentityRow{
	{ID: "b524.interface_identity.api_version", Identity: InterfaceIdentityAPI, NativeField: "b524InterfaceIdentityRecord.identity.APIVersion"},
	{ID: "b524.interface_identity.revision", Identity: InterfaceIdentityRevision, NativeField: "b524InterfaceIdentityRecord.identity.Revision"},
}

// Selector returns the identity request of the row.
func (r InterfaceIdentityRow) Selector() InterfaceIdentitySelector {
	return InterfaceIdentitySelector{ID: r.Identity}
}

// InterfaceIdentityRead is the native result of one acquired identity reply: the raw reply and its decimal qualifier.
type InterfaceIdentityRead struct {
	Row        string
	Request    []byte
	Raw        []byte
	Outcome    ReadOutcome
	Reason     InterfaceIdentityReason
	Qualifier  string
	AcquiredAt time.Time
	Source     byte
	Target     byte
}

// InterfaceIdentityProvenance is what the stored record keeps besides the qualifiers. The instants are inputs: they
// are taken as given and never derived from the clock of the read.
type InterfaceIdentityProvenance struct {
	Address, Source      byte
	Session              uint64
	ControllerObservedAt time.Time
	AcquiredAt           time.Time
}

// InterfaceIdentityRecord is the stored identity of one controller: both qualifiers, the raw replies they came from and the provenance.
type InterfaceIdentityRecord struct {
	APIVersion, Revision string
	Raw                  [2][4]byte
	Provenance           InterfaceIdentityProvenance
}

// InterfaceIdentityQualifier returns the decimal text of a four-byte little-endian float32 that holds a whole number
// from zero to 2^24, and false for any other reply.
func InterfaceIdentityQualifier(reply []byte) (string, bool) {
	if len(reply) != 4 {
		return "", false
	}
	value := float64(math.Float32frombits(binary.LittleEndian.Uint32(reply)))
	if math.IsNaN(value) || math.IsInf(value, 0) || value < 0 || value > 1<<24 || math.Trunc(value) != value {
		return "", false
	}
	return strconv.FormatUint(uint64(value), 10), true
}

// Read evaluates one acquired identity reply against the row.
func (r InterfaceIdentityRow) Read(rec Record) (InterfaceIdentityRead, error) {
	if s, ok := rec.Selector.(InterfaceIdentitySelector); !ok || s != r.Selector() {
		return InterfaceIdentityRead{}, ErrInterfaceIdentitySelector
	}
	out := InterfaceIdentityRead{Row: r.ID, Request: rec.Request, Raw: rec.Reply, Outcome: OutcomeNoValue, AcquiredAt: rec.Clock, Source: rec.Source, Target: rec.Target}
	switch {
	case rec.Diagnostics.SendError != "":
		out.Reason = InterfaceIdentityReasonSendError
	case !rec.Replied:
		out.Reason = InterfaceIdentityReasonNoReply
	default:
		qualifier, ok := InterfaceIdentityQualifier(rec.Reply)
		if !ok {
			out.Reason = InterfaceIdentityReasonNotQualified
			break
		}
		out.Qualifier, out.Outcome = qualifier, OutcomeValue
	}
	return out, nil
}

// AssembleInterfaceIdentity builds the stored record from the two reads of one acquisition. It returns false unless
// both reads are the API and revision rows' and both produced a qualifier.
func AssembleInterfaceIdentity(api, revision InterfaceIdentityRead, provenance InterfaceIdentityProvenance) (InterfaceIdentityRecord, bool) {
	if api.Row != InterfaceIdentityRows[0].ID || revision.Row != InterfaceIdentityRows[1].ID || api.Outcome != OutcomeValue || revision.Outcome != OutcomeValue {
		return InterfaceIdentityRecord{}, false
	}
	out := InterfaceIdentityRecord{APIVersion: api.Qualifier, Revision: revision.Qualifier, Provenance: provenance}
	copy(out.Raw[0][:], api.Raw)
	copy(out.Raw[1][:], revision.Raw)
	return out, true
}
