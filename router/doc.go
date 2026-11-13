// Package router dispatches method invocations and received broadcasts
// between registry planes and an eBUS transport supplied by the caller. It
// builds request frames through the plane and decodes replies through the
// plane; it owns no framing, arbitration or retry policy of its own.
//
// Invocation. A method is invoked by name on one plane. Mutating methods
// always go to the bus. Read-only invocations are coalesced: concurrent
// calls with the same plane name, method and parameters share one bus
// transaction while it succeeds. When the leading call fails, only that
// caller receives the error; each waiter then starts its own transaction.
// A successful reply frame is reused for later identical
// calls until it is older than the refresh interval of its class (state or
// configuration). The key uses the plane name, not the plane instance, so
// planes of different devices that share a name must differ in a parameter
// to stay apart. The class is inferred from the method name, may be overridden by
// the router's classification option, and a call's own parameters take final
// precedence; the intervals and disabling coalescing are per-router options. Only the raw reply
// frame is cached; every caller decodes it again, and failed transactions
// are never cached. A caller waiting on another caller's transaction
// returns early when its own context ends.
//
// Broadcasts. Planes subscribe to primary and secondary command bytes.
// Every subscribed plane sees each matching broadcast, and planes that
// decode broadcasts publish the decoded values as events on one bounded
// channel. Publishing never blocks the receive path: when the channel is
// full the event is dropped, counted, and reported as an overflow error to
// the caller of the broadcast handler. Errors from one plane do not stop
// delivery to the others.
package router
