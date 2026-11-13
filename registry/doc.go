// Package registry owns the eBUS device registry: which devices are on the
// bus, under which addresses, with what identity, and which typed planes,
// methods and projections each device exposes.
//
// Address faces and role-based addressing. One physical device can answer
// under several eBUS addresses (an initiator address, its target companion,
// a broadcast face). The registry keeps one slot per address and joins
// faces into one device entry. Each face records its role, how it was first
// discovered and how far it has been verified; the discovery source is
// fixed once known and is not a confidence order, while the verification
// state only advances. Routing must ask an entry for the address of the
// role it needs: the display address is representative only and may name
// the wrong side of an aliased pair, so it is never a wire destination.
//
// Identity admission. Independently observed addresses are merged into one
// entry only on a complete qualified identity (manufacturer, device ID and a
// serial number that is not a known placeholder value) or through an
// explicit alias. Partial or model-only evidence never creates
// cross-address identity, and a contradictory observation does not
// overwrite an identity that is already published. A static seed enters a
// face as a candidate and a passive observation as corroborated; active
// registration confirms a face's identity when the identity is qualified,
// and a valid directed 07 04 response during Scan also confirms the
// responding face even when the optional serial read returns nothing, so a
// scan-created entry can be confirmed without a qualified triple. That
// confirmation spreads to other faces only along explicit source-target or companion
// topology, never through shared identity membership alone. A partial
// observation for a known address keeps the last-known-good values of the
// fields it omits instead of replacing the entry.
//
// Scanning. Scan and ScanDirected probe addresses with the standard 07 04
// identification request through a caller-supplied bus and register what
// answers. ScanDirected requires an explicit, non-empty target list and
// fails without touching the bus otherwise; only Scan falls back to the
// default address range.
//
// Planes and methods. Plane providers are matched against each device's
// identity and create its planes when the device is registered. Method
// safety (mutability, danger, routability) is resolved from optional
// method interfaces with conservative defaults: a method that does not
// declare otherwise is mutating and dangerous unless it reports itself
// read-only.
//
// Projections. Projections describe a device's data as nodes and edges per
// plane; node identifiers are derived from canonical paths, so the same
// node keeps its identifier across planes. The portal and canonical
// indexes give read-only lookups over them. When a device's projections do
// not form one consistent canonical index, the entry carries no
// projections rather than a partial set.
//
// Concurrency. All registry methods are safe for concurrent use. Lookups
// that return pointers expose live state that may change after the call;
// concurrent readers use the snapshot lookups, which return value copies.
// The callbacks that observe the observation generation or a qualified
// identity witness run while the registry lock is held, so they must stay
// bounded and must not call back into the registry.
package registry
