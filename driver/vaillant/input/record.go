package input

import "time"

// ZoneInputs are the values the zone mode and preset derivation reads.
type ZoneInputs struct {
	OperatingMode, SpecialFunction, CircuitType *uint16
	HasCircuitType                              bool
}

// DHWInputs are the values the hot-water mode and preset derivation reads.
type DHWInputs struct {
	OperatingMode, SpecialFunction *uint16
}

// Qualifiers are the current and cached derivation inputs of one observation.
type Qualifiers struct {
	Zone, CachedZone ZoneInputs
	DHW, CachedDHW   DHWInputs
}

// Record carries one acquired reply: the selector, the request and reply bytes,
// the frame addresses, the qualifiers, the clock of the paired run and the
// diagnostics of the read. Replied is false when the read returned no frame.
type Record struct {
	Selector    Selector
	Request     []byte
	Reply       []byte
	Replied     bool
	Source      byte
	Target      byte
	Qualifiers  Qualifiers
	Clock       time.Time
	Diagnostics Diagnostics
}
