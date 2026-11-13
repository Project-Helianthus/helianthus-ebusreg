// Package productids holds embedded product knowledge for Vaillant-group
// devices: a part-number catalogue and an opt-in table of static address
// seeds. It does not talk to the bus.
//
// Catalogue. The catalogue is parsed from embedded CSV on every load.
// Loading fails when a required column is missing. Every row is kept in
// the full list; only rows with complete metadata are indexed by part
// number, and for a repeated part number the first complete row wins.
//
// Controller capability is classified from the catalogue role and fails
// closed: an empty or unlisted part number, a row without a role, or
// duplicate rows that disagree all yield unknown. A regulator or
// thermostat role means a controller is present; any other role means
// none. Incomplete metadata elsewhere in a row does not turn a listed part
// number into unknown.
//
// Static seeds. The seed table is returned only when explicitly enabled.
// Its entries name the addresses a known device family is expected to
// occupy, each with candidate confidence. The records carry only manufacturer
// and device ID, so registering each address does not join faces of one
// physical device; the caller must alias the addresses explicitly. A seed
// never confirms identity on its own.
package productids
