package owner

import (
	"time"

	"github.com/Project-Helianthus/helianthus-ebusreg/driver/vaillant/output"
)

// boilerFlowMark is the boiler flow temperature kept beside the native snapshot.
// Cached is set when the value was restored. observedAt is the cache write time
// or the live reply time, and it is zero when the age is unknown.
type boilerFlowMark struct {
	celsius    *float64
	cached     bool
	observedAt time.Time
}

func applyOwnBoilerFlow(selector uint16, mark boilerFlowMark, payload []byte, now time.Time) (boilerFlowMark, bool) {
	previous := output.BoilerFlowReading{Celsius: mark.celsius, Cached: mark.cached, ObservedAt: mark.observedAt}
	next, applied := output.ApplyBoilerOwnFlow(selector, previous, payload, now)
	if !applied {
		return mark, false
	}
	return boilerFlowMark{celsius: next.Celsius, observedAt: next.ObservedAt}, true
}

func restoreBoilerFlowMark(celsius *float64, persistedAt time.Time) boilerFlowMark {
	reading := output.RestoreBoilerFlow(celsius, persistedAt)
	return boilerFlowMark{celsius: reading.Celsius, cached: reading.Cached, observedAt: reading.ObservedAt}
}

func boilerFlowWithoutSystem(systemFlow, boilerFlow *float64) *float64 {
	return output.BoilerFlowWithoutSystem(systemFlow, boilerFlow)
}
