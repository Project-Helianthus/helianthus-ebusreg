package owner

import (
	"context"
	"errors"

	"github.com/Project-Helianthus/helianthus-ebusreg/driver/vaillant/owner/b503session"
)

var errNotSupported = errors.New("not supported")

// RPCDispatcher is the rpcdispatcher type of the Vaillant owner.
type RPCDispatcher interface {
	Invoke(ctx context.Context, target byte, payload []byte) ([]byte, error)
}

// B503DispatchOutcome is the b503 dispatch outcome type of the Vaillant owner.
type B503DispatchOutcome struct {
	Response  []byte
	Emitted   bool
	Native    b503session.NativeResult
	Err       error
	Transport b503session.TransportKey
}

// B503OutcomeDispatcher is the b503 outcome dispatcher type of the Vaillant owner.
type B503OutcomeDispatcher interface {
	RPCDispatcher
	InvokeB503Outcome(ctx context.Context, target byte, payload []byte) B503DispatchOutcome
}

// InvokeB503Operation implements the invoke b503 operation step of the Vaillant owner.
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
