// Package vaillant catalogues field layouts, candidate fields and value
// views of the Vaillant vendor services (primary command 0xB5) as embedded,
// qualification-graded data, together with pure functions that read those
// catalogues and interpret caller-supplied evidence against them.
//
// The package has no bus access and no registry or router coupling; plane
// providers and consumers decide what to send and when. It does not
// establish protocol facts: entries carry their qualification and grade
// (Proven, Hypothesis or Unknown), and nothing here promotes a Hypothesis
// or Unknown entry.
//
// Validation is per catalogue, not package-wide. Most loaders and
// interpreters read their embedded or supplied JSON strictly, rejecting
// unknown fields, trailing data, duplicate keys and incomplete entries
// instead of returning a partial list, and report "not found" separately
// from a catalogue error. A few older candidate loaders decode their
// embedded list without that validation and return it as is, so callers
// must not assume strictness for a loader whose signature has no error
// result. Where raw bytes are taken they are kept next to any
// interpretation, and an unsupported context is reported as such rather
// than interpreted.
package vaillant
