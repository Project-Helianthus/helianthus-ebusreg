package output

import (
	"encoding/binary"
	"time"
)

// BoilerFlowSelector is the boiler register whose DATA2c reply is the boiler
// flow temperature. The controller system flow is a different register.
const BoilerFlowSelector uint16 = 0x1800

// BoilerFlowReading is the native boiler flow temperature in degree Celsius,
// marked cached when it was restored from the semantic cache. ObservedAt is
// the cache write time or the live reply time. Age is unknown when it is zero.
type BoilerFlowReading struct {
	Celsius    *float64
	Cached     bool
	ObservedAt time.Time
}

// Age reports the time since ObservedAt. It is unknown when the reading has
// no temperature or no observation time.
func (r BoilerFlowReading) Age(now time.Time) (time.Duration, bool) {
	if r.Celsius == nil || r.ObservedAt.IsZero() {
		return 0, false
	}
	return now.Sub(r.ObservedAt), true
}

// RestoreBoilerFlow marks a cached flow temperature and keeps the cache write
// time as its observation time. A missing temperature restores nothing.
func RestoreBoilerFlow(celsius *float64, persistedAt time.Time) BoilerFlowReading {
	if celsius == nil {
		return BoilerFlowReading{}
	}
	value := *celsius
	reading := BoilerFlowReading{Celsius: &value, Cached: true}
	if !persistedAt.IsZero() {
		reading.ObservedAt = persistedAt.UTC()
	}
	return reading
}

// ApplyBoilerOwnFlow stores a valid DATA2c reply of BoilerFlowSelector.
// The invalid marker and a short payload leave the previous reading. A valid
// zero replaces it and clears the cached mark. Any other selector leaves the
// previous reading.
func ApplyBoilerOwnFlow(selector uint16, previous BoilerFlowReading, payload []byte, now time.Time) (BoilerFlowReading, bool) {
	if selector != BoilerFlowSelector || len(payload) < 2 || (payload[0] == 0x00 && payload[1] == 0x80) {
		return previous, false
	}
	raw := int16(binary.LittleEndian.Uint16(payload[:2]))
	value := float64(raw) / 16.0
	return BoilerFlowReading{Celsius: &value, ObservedAt: now}, true
}

// BoilerFlowWithoutSystem returns the boiler flow. The controller system flow
// is not a source for this field.
func BoilerFlowWithoutSystem(systemFlow, boilerFlow *float64) *float64 {
	if boilerFlow == nil {
		return nil
	}
	value := *boilerFlow
	_ = systemFlow
	return &value
}
