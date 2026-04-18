// Package ebus_standard_catalog implements the catalog schema and loader for
// the cross-vendor eBUS Standard L7 services namespace (ebus_standard).
//
// The catalog is the single source of truth for methods in this namespace and
// is expressed as a SHA-pinned, version-tagged YAML data file embedded at
// compile time. Each method carries a full 14-tuple identity key used to
// detect collisions and ambiguous length-selector branches at generation
// time.
//
// Plan reference: ebus-standard-l7-services-w16-26.locked/00-canonical.md
// Canonical SHA-256: 9e0a29bb76d99f551904b05749e322aafd3972621858aa6d1acbe49b9ef37305
package ebus_standard_catalog
