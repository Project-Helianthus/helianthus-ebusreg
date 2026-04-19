// Package safetypolicy holds the invoke-boundary safety-class enforcement
// primitive for the generic ebus_standard provider.
//
// The package is placed under catalog/ebus_standard/internal/ so that Go's
// internal-package visibility rule structurally forbids any importer whose
// path is not rooted at catalog/ebus_standard/. This is the M3 AC6
// BUILD_TIME_NAMESPACE_ISOLATION mechanism (option a): the Vaillant tree
// cannot reach this symbol and therefore cannot reuse the gate by accident.
// A second-line static-import regression test lives at
// catalog/ebus_standard/provider_namespace_isolation_test.go and asserts
// that no file inside catalog/ebus_standard/... imports any vaillant/...
// package.
package safetypolicy

// Class is the subset of catalog safety classes the gate knows about. The
// values MUST equal the string constants declared in identity.go so callers
// can pass catalog SafetyClass values directly (after a type conversion).
type Class string

// Known safety classes (must mirror identity.go — duplicated here because
// this package cannot import the outer package without a cycle).
const (
	ReadOnlySafe    Class = "read_only_safe"
	ReadOnlyBusLoad Class = "read_only_bus_load"
	Mutating        Class = "mutating"
	Destructive     Class = "destructive"
	Broadcast       Class = "broadcast"
	MemoryWrite     Class = "memory_write"
)

// Caller tags the intent of an Invoke call. Future M4+ expansion may
// whitelist mutating for system_nm_runtime; in M3 the gate is identical
// for every caller context.
type Caller int

// Caller contexts.
const (
	CallerUserFacing Caller = iota
	CallerSystemNMRuntime
)

// Allow reports whether a method carrying the given safety class may be
// dispatched by a caller with the given context. Default-deny: every class
// outside the explicit allow-list returns false.
func Allow(class Class, caller Caller) bool {
	_ = caller
	switch class {
	case ReadOnlySafe, ReadOnlyBusLoad:
		return true
	default:
		// mutating | destructive | broadcast | memory_write => deny.
		// Unknown strings also deny, preserving fail-closed semantics.
		return false
	}
}
