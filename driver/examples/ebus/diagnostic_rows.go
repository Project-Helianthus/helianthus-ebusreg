package ebus

import (
	"bytes"
	"errors"

	ebuserrors "github.com/Project-Helianthus/helianthus-ebusgo/errors"
	"github.com/Project-Helianthus/helianthus-ebusgo/protocol"
	"github.com/Project-Helianthus/helianthus-ebusgo/protocol/vaillant/b503"
)

type diagnosticSelector struct {
	Family, Selector byte
	Index            *byte
}

type liveMonitorAction string

const (
	liveMonitorEnable  liveMonitorAction = "enable"
	liveMonitorRead    liveMonitorAction = "read"
	liveMonitorDisable liveMonitorAction = "disable"
)

type liveMonitorRequest struct {
	Target byte
	Action liveMonitorAction
}

func (r liveMonitorRequest) request(source byte) (protocol.Frame, error) {
	return diagnosticSelector{Family: 0, Selector: 3}.request(source, r.Target)
}

var errDiagnosticShape = errors.New("not a B5/03 exemplar body")

func diagnosticBody(body []byte) (diagnosticSelector, error) {
	if len(body) < 2 || len(body) > 3 || bytes.HasPrefix(body, []byte{0xb5, 0x03}) {
		return diagnosticSelector{}, errDiagnosticShape
	}
	s := diagnosticSelector{Family: body[0], Selector: body[1]}
	if len(body) == 3 {
		i := body[2]
		s.Index = &i
	}
	if _, err := s.request(0, 0); err != nil {
		return diagnosticSelector{}, err
	}
	return s, nil
}

func (s diagnosticSelector) request(source, target byte) (protocol.Frame, error) {
	b, err := b503.EncodeRequest(s.Family, s.Selector)
	if err != nil {
		return protocol.Frame{}, err
	}
	if s.Index != nil {
		if s.Family != 1 || (s.Selector != 1 && s.Selector != 2) {
			return protocol.Frame{}, errDiagnosticShape
		}
		b = append(b, *s.Index)
	}
	return protocol.Frame{Source: source, Target: target, Primary: 0xb5, Secondary: 0x03, Data: b}, nil
}

type diagnosticStamp struct {
	AdapterInstanceID string
	TransportEpoch    uint64
}

type diagnosticNative byte

const (
	diagnosticAmbiguous diagnosticNative = iota
	diagnosticACK
	diagnosticNAK
)

type diagnosticCompletion struct {
	Issued, Current diagnosticStamp
	Emitted         bool
	Reply           *protocol.Frame
	SendError       error
	ClassifiedError error
}

type diagnosticValue struct {
	Transport       diagnosticStamp
	Emitted         bool
	Native          diagnosticNative
	Response        []byte
	ClassifiedError error
}

func (c diagnosticCompletion) observation() diagnosticValue {
	v := diagnosticValue{Transport: c.Issued, Emitted: c.Emitted, ClassifiedError: c.ClassifiedError}
	if !c.Emitted || c.Issued != c.Current {
		return v
	}
	if errors.Is(c.SendError, ebuserrors.ErrNACK) {
		v.Native = diagnosticNAK
	} else if c.SendError == nil && c.ClassifiedError == nil && c.Reply != nil {
		v.Native, v.Response = diagnosticACK, c.Reply.Data
	}
	return v
}

func (s diagnosticSelector) decode(v diagnosticValue) (any, error) {
	if v.Native != diagnosticACK {
		return nil, v.ClassifiedError
	}
	switch [2]byte{s.Family, s.Selector} {
	case [2]byte{0, 1}:
		return b503.DecodeCurrentError(v.Response)
	case [2]byte{0, 2}:
		return b503.DecodeCurrentService(v.Response)
	case [2]byte{1, 1}:
		return b503.DecodeErrorHistory(v.Response)
	case [2]byte{1, 2}:
		return b503.DecodeServiceHistory(v.Response)
	case [2]byte{0, 3}:
		return b503.DecodeLiveMonitorMain(v.Response)
	default:
		return nil, errDiagnosticShape
	}
}
