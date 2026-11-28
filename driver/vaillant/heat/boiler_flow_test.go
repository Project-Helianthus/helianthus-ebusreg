package heat

import (
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"time"

	"github.com/Project-Helianthus/helianthus-ebusreg/driver/vaillant/output"
)

// Identity fields, and the production mutation each assertion catches:
//   - selector 0x1800: ApplyBoilerOwnFlow stores the DATA2c reply. A check for another selector keeps the previous temperature.
//   - selector 0x1801: the same DATA2c bytes are not stored. A check that ignores the selector stores 40.
//   - payload 0x80 0x02: the stored temperature is 40. A scale other than 1/16 stores another number.
//   - payload 0x00 0x00: the stored temperature is 0 and the cached mark clears. A rejection of valid zero keeps the previous temperature.
//   - payload 0x00 0x80: the invalid marker leaves the previous temperature and its observation time. A decode of that marker stores -2048.
//   - short payload: the previous temperature stays. A one-byte payload stored as a temperature replaces it.
//   - cached: a restored reading is cached and a live reply clears it. Leaving the mark set after a live reply keeps the cached flag.
//   - observed time: the age is now minus that time. Using the clock of the invalid reply, or a fixed age, changes the duration.
//   - persisted time: a second cache write time changes the age. Ignoring persisted_at makes both ages equal.
//   - zero persisted time: the restored reading is cached and its age is unknown. Inventing an observation time makes the age known.
//   - system flow: BoilerFlowWithoutSystem returns the boiler value. Copying the system value returns 42.5 when the boiler value is absent or 40.
//   - flow_temperature_c: a cache load restores 50.5 and 54.5. Dropping the field on load leaves the temperature nil.
//   - heating_status_raw and the pump flag: a cache load leaves both unset. Copying either field publishes the cached raw status or the pump.

func TestBoilerFlowTemperatureRestored(t *testing.T) {
	if output.BoilerFlowSelector != 0x1800 {
		t.Fatalf("boiler flow selector = %#x, want 0x1800", output.BoilerFlowSelector)
	}
	persisted := time.Date(2030, 4, 10, 8, 0, 0, 0, time.UTC)
	laterPersist := time.Date(2030, 4, 10, 8, 30, 0, 0, time.UTC)
	now := time.Date(2030, 4, 10, 9, 30, 0, 0, time.UTC)
	liveAt := time.Date(2030, 4, 10, 10, 0, 0, 0, time.UTC)

	restored := loadBoilerFlow(t, boilerFlowCache(50.5, 4, true, persisted), persisted)
	if restored.Celsius == nil || *restored.Celsius != 50.5 || !restored.Cached {
		t.Fatalf("restored = %+v, want cached 50.5", restored)
	}
	age, known := restored.Age(now)
	if !known || age != 90*time.Minute || !restored.ObservedAt.Equal(persisted) {
		t.Fatalf("age = %v known=%v observed=%s", age, known, restored.ObservedAt)
	}
	other := loadBoilerFlow(t, boilerFlowCache(54.5, 9, false, laterPersist), laterPersist)
	otherAge, otherKnown := other.Age(now)
	if other.Celsius == nil || *other.Celsius != 54.5 || !otherKnown || otherAge != 60*time.Minute {
		t.Fatalf("second restore = %+v age=%v", other, otherAge)
	}
	if restored.StateHeating != nil || restored.Pump != nil || other.StateHeating != nil || other.Pump != nil {
		t.Fatal("cache load copied heating status or the pump into the boiler flow")
	}

	missing := loadBoilerFlow(t, []byte(`{"schema_version":3,"metadata":{"persisted_at":"2030-04-10T08:00:00Z"},"boiler":{"return_temperature_c":31.5,"heating_status_raw":4}}`), persisted)
	if missing.Celsius != nil || missing.Cached || missing.Return == nil || *missing.Return != 31.5 || missing.StateHeating != nil {
		t.Fatalf("cache without flow = %+v, want return 31.5 and no flow", missing)
	}
	noTime := output.RestoreBoilerFlow(floatPtr(36), time.Time{})
	if noAge, ok := noTime.Age(now); noTime.Celsius == nil || *noTime.Celsius != 36 || !noTime.Cached || ok || !noTime.ObservedAt.IsZero() {
		t.Fatalf("restore without a write time = %+v ageKnown=%v age=%v", noTime, ok, noAge)
	}
	if got := output.RestoreBoilerFlow(nil, persisted); got.Celsius != nil || got.Cached {
		t.Fatalf("restore of no temperature = %+v", got)
	}

	system, boiler := 42.5, 40.0
	if got := output.BoilerFlowWithoutSystem(&system, nil); got != nil {
		t.Fatalf("system flow filled an empty boiler flow with %v", *got)
	}
	if got := output.BoilerFlowWithoutSystem(&system, &boiler); got == nil || *got != 40 {
		t.Fatalf("boiler flow with a system flow = %v, want 40", got)
	}

	live, applied := output.ApplyBoilerOwnFlow(0x1800, restored.BoilerFlowReading, []byte{0x80, 0x02}, liveAt)
	if !applied || live.Cached || live.Celsius == nil || *live.Celsius != 40 || !live.ObservedAt.Equal(liveAt) {
		t.Fatalf("live reply = %+v applied=%v", live, applied)
	}
	if age, known = live.Age(liveAt.Add(15 * time.Minute)); !known || age != 15*time.Minute {
		t.Fatalf("live age = %v known=%v", age, known)
	}
	kept, applied := output.ApplyBoilerOwnFlow(0x1800, live, []byte{0x00, 0x80}, liveAt.Add(time.Hour))
	if applied || kept.Celsius == nil || *kept.Celsius != 40 || !kept.ObservedAt.Equal(liveAt) || kept.Cached {
		t.Fatalf("invalid marker = %+v applied=%v", kept, applied)
	}
	kept, applied = output.ApplyBoilerOwnFlow(0x1801, live, []byte{0x80, 0x02}, liveAt.Add(time.Hour))
	if applied || kept.Celsius == nil || *kept.Celsius != 40 || !kept.ObservedAt.Equal(liveAt) {
		t.Fatalf("other selector = %+v applied=%v", kept, applied)
	}
	kept, applied = output.ApplyBoilerOwnFlow(0x1800, live, []byte{0x80}, liveAt.Add(time.Hour))
	if applied || kept.Celsius == nil || *kept.Celsius != 40 {
		t.Fatalf("short payload = %+v applied=%v", kept, applied)
	}
	zero, applied := output.ApplyBoilerOwnFlow(0x1800, live, []byte{0x00, 0x00}, liveAt.Add(2*time.Hour))
	if !applied || zero.Cached || zero.Celsius == nil || *zero.Celsius != 0 || !zero.ObservedAt.Equal(liveAt.Add(2*time.Hour)) {
		t.Fatalf("valid zero = %+v applied=%v", zero, applied)
	}
}

type loadedBoilerFlow struct {
	output.BoilerFlowReading
	Return       *float64
	StateHeating *int
	Pump         *bool
}

func loadBoilerFlow(t *testing.T, body []byte, persisted time.Time) loadedBoilerFlow {
	t.Helper()
	path := filepath.Join(t.TempDir(), "semantic-cache.json")
	if err := os.WriteFile(path, body, 0o600); err != nil {
		t.Fatal(err)
	}
	snapshot, ok := output.NewSemanticCacheStore(path, func(string, ...any) {}).Load()
	if !ok || snapshot.Boiler == nil {
		t.Fatalf("cache load ok=%v boiler=%v", ok, snapshot.Boiler)
	}
	if !snapshot.PersistedAt.Equal(persisted) {
		t.Fatalf("persisted_at = %s, want %s", snapshot.PersistedAt, persisted)
	}
	reading := output.RestoreBoilerFlow(snapshot.Boiler.State.FlowTemperatureC, snapshot.PersistedAt)
	return loadedBoilerFlow{
		BoilerFlowReading: reading,
		Return:            snapshot.Boiler.State.ReturnTemperatureC,
		StateHeating:      snapshot.Boiler.Diagnostics.HeatingStatusRaw,
		Pump:              snapshot.Boiler.State.CentralHeatingPumpActive,
	}
}

func boilerFlowCache(flow float64, heating int, pump bool, persisted time.Time) []byte {
	pumpText := "false"
	if pump {
		pumpText = "true"
	}
	return []byte(`{"schema_version":3,"metadata":{"persisted_at":"` + persisted.Format(time.RFC3339) + `"},"boiler":{"flow_temperature_c":` + strconv.FormatFloat(flow, 'f', -1, 64) + `,"central_heating_pump_active":` + pumpText + `,"heating_status_raw":` + strconv.Itoa(heating) + `}}`)
}

func floatPtr(value float64) *float64 {
	return &value
}
