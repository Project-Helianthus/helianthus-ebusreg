package owner

import (
	"context"
	"errors"

	"github.com/Project-Helianthus/helianthus-ebusreg/driver/vaillant/owner/b503session"
)

var errNotSupported = errors.New("not supported")

type RPCDispatcher interface {
	Invoke(ctx context.Context, target byte, payload []byte) ([]byte, error)
}

type B503DispatchOutcome struct {
	Response  []byte
	Emitted   bool
	Native    b503session.NativeResult
	Err       error
	Transport b503session.TransportKey
}

type B503OutcomeDispatcher interface {
	RPCDispatcher
	InvokeB503Outcome(ctx context.Context, target byte, payload []byte) B503DispatchOutcome
}

func InvokeB503Operation(ctx context.Context, dispatcher RPCDispatcher, target byte, payload []byte) b503session.DispatchOutcome {
	if dispatcher == nil {
		return b503session.DispatchOutcome{Err: errNotSupported}
	}
	aware, ok := dispatcher.(B503OutcomeDispatcher)
	if !ok {
		return b503session.DispatchOutcome{Err: errNotSupported}
	}
	outcome := aware.InvokeB503Outcome(ctx, target, payload)
	return b503session.DispatchOutcome{
		Response: outcome.Response, Emitted: outcome.Emitted, Native: outcome.Native,
		Err: outcome.Err, Transport: outcome.Transport,
	}
}
