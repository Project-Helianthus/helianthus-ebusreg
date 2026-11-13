// Package system provides the system plane for Vaillant-group devices and
// the controller-side knowledge it relies on: operational data (B504 reads
// and B505 writes), register access (B509), extended register and timer
// access (B524), and passive decoding of energy statistics broadcasts
// (B516).
//
// Matching. The provider matches on the manufacturer name only and adds
// the same plane to every matching device; whether a register or opcode
// exists on a given device is not decided at match time.
//
// Requests and replies. The plane implements the router's plane contract.
// A request goes to the device's own address and needs the initiator
// address as the source parameter; the target parameter is consulted only
// when the device has no address. Decoded replies keep the raw payload
// next to the decoded values, and a reply with other command bytes fails
// with the transport library's invalid-payload error. The raw B524 method
// sends caller-supplied bytes verbatim (at most 16), so it is declared
// mutating: the router never caches it and method metadata treats it as
// dangerous.
//
// Broadcasts. Of the broadcasts the plane subscribes to, it decodes only
// B516 frames whose selector it can parse; other frames are left unhandled
// rather than partially published. Selector codes are kept raw next to
// their names, and a code without a known name yields an invalid name, not
// a guess. The packed selector date is absolute and independent of the
// wall clock; an invalid date leaves the date fields invalid without
// discarding the other decoded fields. Energy values are in watt-hours.
//
// Bounds. B524 register bounds are scoped to an exact
// controller identity (device, software and hardware revisions, API
// version and revision) and to the exact opcode, group, instance and
// register. A lookup outside a seeded profile or selector reports the
// bound as unknown and never substitutes a default.
package system
