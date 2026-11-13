// Package match holds device-matching predicates shared by the plane
// providers of this module and not exported beyond it.
//
// Hardware versions are compared as the decimal number formed by all
// digits in the version string, ignoring every other character. The
// predicates fail closed: a version without digits, or one whose value
// does not fit in 16 bits, never satisfies a lower bound.
package match
