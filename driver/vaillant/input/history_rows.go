package input

import (
	"bytes"
	"errors"
	"time"

	ebuserrors "github.com/Project-Helianthus/helianthus-ebusgo/errors"
	"github.com/Project-Helianthus/helianthus-ebusgo/protocol/vaillant/b503"
)

const diagnosticStaleEpochText = "b503dispatcher: stale-epoch reply discarded"

// ErrDiagnosticSelectorMismatch is returned when a record was acquired for another B503 selector than the row's.
var ErrDiagnosticSelectorMismatch = errors.New("record selector is not the selector of the row")

// ErrLiveMonitorContext is returned when the key, the owner qualifier or the body of a live-monitor request does not resolve to the action.
var ErrLiveMonitorContext = errors.New("live-monitor request does not resolve to the managed context")

// TransportStamp is the transport generation a request was issued on.
type TransportStamp struct {
	AdapterInstanceID string
	TransportEpoch    uint64
}

// DiagnosticNative is the wire disposition of a B503 request.
type DiagnosticNative uint8

const (
	// DiagnosticAmbiguous is a request whose wire outcome is not established: not sent, failed without a typed peer
	// rejection, answered with no frame, or answered on a transport that has since advanced.
	DiagnosticAmbiguous DiagnosticNative = iota
	// DiagnosticACK is a request the peer answered on the transport it was issued on.
	DiagnosticACK
	// DiagnosticNAK is a request the peer rejected with a typed NACK.
	DiagnosticNAK
)

// DiagnosticRow is one B503 request body. Indexed rows may carry the optional history index byte.
type DiagnosticRow struct {
	ID               string
	Family, Selector byte
	Indexed          bool
	LiveMonitor      bool
}

// DiagnosticRows lists the B503 request bodies the poller's raw-frame dispatcher carries.
var DiagnosticRows = []DiagnosticRow{
	{ID: "b503.current_error_slots", Family: 0x00, Selector: 0x01},
	{ID: "b503.current_service_slots", Family: 0x00, Selector: 0x02},
	{ID: "b503.error_history", Family: 0x01, Selector: 0x01, Indexed: true},
	{ID: "b503.service_history", Family: 0x01, Selector: 0x02, Indexed: true},
	{ID: "b503.live_monitor_main", Family: 0x00, Selector: 0x03, LiveMonitor: true},
}

// DiagnosticDispatch is what the dispatcher knew when the request finished: the transport it was issued on, the
// transport current at completion, whether the send was entered, the error the send returned and the text the
// dispatcher classified a failure with.
type DiagnosticDispatch struct {
	Issued, Current TransportStamp
	Emitted         bool
	SendError       error
	Classified      string
}

// DiagnosticResult is the native result of one B503 request: the disposition, the response the transport
// delivered and the transport it was issued on. Index is the history index of the request when it carried one.
type DiagnosticResult struct {
	Row         string
	Request     []byte
	Index       *byte
	Emitted     bool
	Native      DiagnosticNative
	Response    []byte
	Err         string
	Transport   TransportStamp
	Decoded     any
	DecodeError string
	AcquiredAt  time.Time
	Source      byte
	Target      byte
}

func (r DiagnosticRow) accepts(selector Selector) (*byte, bool) {
	if r.LiveMonitor {
		_, ok := selector.(LiveMonitorSelector)
		return nil, ok
	}
	s, ok := selector.(DiagnosticSelector)
	if !ok || s.Family != r.Family || s.Selector != r.Selector || (s.Index != nil && !r.Indexed) {
		return nil, false
	}
	return s.Index, true
}

// Observe evaluates one dispatched request against the row. A reply that completes after the transport epoch differs
// from the epoch the request was issued on is discarded, whatever it was; the dispatcher compares the epoch only.
func (r DiagnosticRow) Observe(rec Record, dispatch DiagnosticDispatch) (DiagnosticResult, error) {
	index, ok := r.accepts(rec.Selector)
	if !ok {
		return DiagnosticResult{}, ErrDiagnosticSelectorMismatch
	}
	out := DiagnosticResult{
		Row: r.ID, Request: rec.Request, Index: index, Emitted: dispatch.Emitted, Err: dispatch.Classified,
		Transport: dispatch.Issued, AcquiredAt: rec.Clock, Source: rec.Source, Target: rec.Target,
	}
	switch {
	case !dispatch.Emitted:
	case dispatch.Issued.TransportEpoch != dispatch.Current.TransportEpoch:
		out.Err = diagnosticStaleEpochText
	case errors.Is(dispatch.SendError, ebuserrors.ErrNACK):
		out.Native = DiagnosticNAK
	case dispatch.SendError == nil && rec.Replied:
		out.Native, out.Response = DiagnosticACK, rec.Reply
	}
	if out.Native == DiagnosticACK {
		out.Decoded, out.DecodeError = r.decode(out.Response)
	}
	return out, nil
}

func (r DiagnosticRow) decode(response []byte) (any, string) {
	var decoded any
	var err error
	switch [2]byte{r.Family, r.Selector} {
	case [2]byte{0x00, 0x01}:
		decoded, err = b503.DecodeCurrentError(response)
	case [2]byte{0x00, 0x02}:
		decoded, err = b503.DecodeCurrentService(response)
	case [2]byte{0x01, 0x01}:
		decoded, err = b503.DecodeErrorHistory(response)
	case [2]byte{0x01, 0x02}:
		decoded, err = b503.DecodeServiceHistory(response)
	default:
		decoded, err = b503.DecodeLiveMonitorMain(response)
	}
	if err != nil {
		return nil, err.Error()
	}
	return decoded, ""
}

// Resolve checks that a live-monitor request belongs to the managed context of the selector's action: the replay key,
// the qualifier the owner supplied and the body on the wire must all name it. All three actions send the same body, so
// the key and the qualifier alone tell them apart.
func (s LiveMonitorSelector) Resolve(key string, qualifier LiveMonitorAction, body []byte) error {
	if key != s.Key() || qualifier != s.Action || !bytes.Equal(body, s.Payload()) {
		return ErrLiveMonitorContext
	}
	return nil
}
