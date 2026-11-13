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
// The functions fail closed. Catalogue loaders reject unknown fields,
// trailing data, duplicate keys and incomplete finite catalogues instead
// of returning a partial list, and lookups report "not found" separately
// from a catalogue error. Interpreters read JSON evidence strictly,
// rejecting duplicate keys, trailing data and unexpected shapes; where
// they take raw bytes they keep them next to any interpretation, and an
// unsupported context is reported as such rather than interpreted.
package vaillant
