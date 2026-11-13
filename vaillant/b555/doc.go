// Package b555 provides the timer plane for the Vaillant B555 weekly
// timer and schedule service: reading a timer's configuration, its
// per-weekday slot counts and single slots, and writing a single slot.
//
// Matching. The provider matches on the manufacturer name only. It does
// not decide whether a device implements B555: the service is
// limited to one controller family, and another family keeps its timers
// in B524, so callers check the
// device family before choosing the timer transport.
//
// Requests and replies. The plane implements the router's plane contract
// and has no broadcast subscriptions. A request goes to the device's own
// address and needs the initiator address as the source parameter; the
// target parameter is consulted only when the device has no address. A
// reply decodes to a value map that always carries the raw reply payload
// and the request coordinates next to the decoded fields. A reply shorter
// than its layout is not an error: the fields it lacks are marked invalid,
// and an empty write reply reports the write as not accepted. Missing
// parameters, an unknown method or a reply with other command bytes fail
// with the transport library's invalid-payload error.
//
// The package also exports helpers that bound reply extents and narrow
// numeric values into graded views. They only classify the values they are
// given and send nothing; each result keeps the raw input and carries an
// explicit grade (Proven, Hypothesis or Unknown) instead of a guessed
// value.
package b555
