package heat

import (
	"strings"
	"testing"

	semreg "github.com/Project-Helianthus/helianthus-semreg/semreg/v1"
	"github.com/Project-Helianthus/helianthus-semreg/semreg/v1/packs/thermal"
)

// Identity fields, and the production mutation each assertion catches:
//   - pack id and pack version: the fact key. A key that names another pack still carries the quantity.
//   - fact id: the pack element. A row that swaps two elements keeps the old quantity on the wrong fact.
//   - dimension id: the pack subject dimension. A key that uses the subject word keeps the old resource class.
//   - dimension text: the resource. A key that drops the resource text collides two resources.
//   - group: the native group. A match that ignores the group binds a hot-water row to a zone.
//   - register: the native register. A match that ignores the register binds the neighbouring row.
//   - instance: the native instance bounds. A match that drops the upper bound binds instance 11.
//   - family: the controller family table. A match that ignores the family uses the other family's codes.
//   - selector: the heat-pump locator. A match that also requires a group misses a selector reading.
//   - native unit: the evidence unit. An evidence string that omits the unit publishes a bare number.
//   - celsius offset: kelvin is the native number plus 273.15. Dropping the offset leaves 0 °C as 0 K.
//   - percent scale: a ratio is the percent divided by 100. Leaving the percent unchanged publishes 50 as 50.
//   - valve canonical form: 1.0 is ratio 1. A valve check that rejects the canonical coefficient publishes nothing.
//   - bar scale: one bar is 100000 Pa. The held pressure rows must not emit that quantity.
//   - grade: Hypothesis and Proven both feed. A grade gate drops one of them.
//   - source rank: a stale, failed, absent or unmapped higher source stays. A lower source that replaces it hides the gap.
//   - absence: a missing zone withdraws. An absent reply stays unknown. A feed that still publishes keeps a departed zone.
//   - failure: a failed reply has no quantity. A feed that keeps the text publishes a value that was not read.
//   - deactivation: a deactivated circuit withdraws. A zone reading with the same flag still feeds. A missing flag on another subject still feeds.

func TestThermalZoneFeed(t *testing.T) {
	feeds := ZoneFeeds()
	air := mustFeed(t, feeds, "zone.air_temperature")
	manual := mustFeed(t, feeds, "zone.manual_setpoint")
	effective := mustFeed(t, feeds, "zone.effective_setpoint")
	humidity := mustFeed(t, feeds, "zone.humidity")
	heatingMode := mustFeed(t, feeds, "zone.heating_mode")
	coolingMode := mustFeed(t, feeds, "zone.cooling_mode")
	comfort := mustFeed(t, feeds, "zone.comfort_setpoint")
	comfortOther := mustFeed(t, feeds, "zone.comfort_setpoint_other")
	setback := mustFeed(t, feeds, "zone.unoccupied_heating")
	coolingManual := mustFeed(t, feeds, "zone.occupied_cooling_manual")
	coolingEffective := mustFeed(t, feeds, "zone.occupied_cooling_effective")
	coolingNative := mustFeed(t, feeds, "zone.cooling_setpoint_native")
	special := mustFeed(t, feeds, "zone.special_function_native")
	valve := mustFeed(t, feeds, "zone.valve_position")

	t.Run("identity and scale", func(t *testing.T) {
		east := textAt(air, "0", "zone-east", "controller-700", 0)
		east.Poll = "slow"
		got := Apply(air, east)
		mustQuantity(t, got, "thermal.measurement.air_temperature", "thermal.dimension.zone", "zone-east", "27315", -2, "unit.kelvin")
		if got.Grade != "Hypothesis" {
			t.Fatalf("grade %s", got.Grade)
		}
		requireEvidence(t, got, "native=0", "unit=°C", "grade=Hypothesis", "rank=3", "poll=slow", "freshness=fresh")
		west := Apply(air, textAt(air, "0", "zone-west", "controller-700", 10))
		mustQuantity(t, west, "thermal.measurement.air_temperature", "thermal.dimension.zone", "zone-west", "27315", -2, "unit.kelvin")
		if dimensionText(got.Key) == dimensionText(west.Key) {
			t.Fatal("two resources share one dimension text")
		}
		below := Apply(air, textAt(air, "-10", "zone-east", "controller-720", 1))
		mustQuantity(t, below, "thermal.measurement.air_temperature", "thermal.dimension.zone", "zone-east", "26315", -2, "unit.kelvin")
		same := Apply(manual, textAt(manual, "21.0", "zone-east", "controller-700", 2))
		plain := Apply(manual, textAt(manual, "21", "zone-east", "controller-700", 2))
		mustQuantity(t, same, "thermal.setpoint.temperature_setpoint", "thermal.dimension.zone", "zone-east", "29415", -2, "unit.kelvin")
		mustQuantity(t, plain, "thermal.setpoint.temperature_setpoint", "thermal.dimension.zone", "zone-east", "29415", -2, "unit.kelvin")
		ratio := Apply(humidity, textAt(humidity, "50", "zone-east", "controller-700", 0))
		mustQuantity(t, ratio, "thermal.measurement.humidity", "thermal.dimension.zone", "zone-east", "5", -1, "unit.ratio")
		full := Apply(humidity, textAt(humidity, "100", "zone-east", "controller-700", 0))
		mustQuantity(t, full, "thermal.measurement.humidity", "thermal.dimension.zone", "zone-east", "1", 0, "unit.ratio")
		over := Apply(humidity, textAt(humidity, "150", "zone-east", "controller-700", 0))
		mustNamespace(t, over, "150")
		if over.Gap != "value outside the pack range" || over.Key.FactID == "" {
			t.Fatalf("range status key %s gap %s", over.Key.FactID, over.Gap)
		}
	})

	t.Run("locator", func(t *testing.T) {
		base := textAt(air, "0", "zone-east", "controller-700", 0)
		if got := Apply(air, base); got.Status != StatusFed {
			t.Fatalf("base status %s", got.Status)
		}
		wrongGroup := base
		wrongGroup.Group = 0x01
		if got := Apply(air, wrongGroup); got.Status != StatusUnmatched {
			t.Fatalf("group status %s", got.Status)
		}
		wrongRegister := base
		wrongRegister.Register = base.Register + 1
		if got := Apply(air, wrongRegister); got.Status != StatusUnmatched {
			t.Fatalf("register status %s", got.Status)
		}
		high := base
		high.Instance = 11
		if got := Apply(air, high); got.Status != StatusUnmatched {
			t.Fatalf("instance status %s", got.Status)
		}
		otherFamily := base
		otherFamily.Family = "controller-other"
		if got := Apply(air, otherFamily); got.Status != StatusUnmatched {
			t.Fatalf("family status %s", got.Status)
		}
		marked := base
		marked.Deactivated = true
		if got := Apply(air, marked); got.Status != StatusFed {
			t.Fatalf("zone deactivation status %s", got.Status)
		}
	})

	t.Run("modes", func(t *testing.T) {
		off := Apply(heatingMode, codeAt(heatingMode, 0, "zone-east", "controller-700", 1))
		mustSymbol(t, off, "thermal.mode.zone", "off", "thermal.dimension.zone", "zone-east")
		auto := Apply(heatingMode, codeAt(heatingMode, 1, "zone-east", "controller-700", 1))
		mustSymbol(t, auto, "thermal.mode.zone", "auto", "thermal.dimension.zone", "zone-east")
		day := Apply(heatingMode, codeAt(heatingMode, 2, "zone-east", "controller-700", 1))
		mustProjection(t, day, "Day")
		if day.Key.FactID != "thermal.mode.zone" {
			t.Fatal("day projection lost its fact")
		}
		setbackCode := Apply(heatingMode, codeAt(heatingMode, 3, "zone-east", "controller-700", 1))
		mustProjection(t, setbackCode, "Setback")
		timed := Apply(heatingMode, codeAt(heatingMode, 1, "zone-east", "controller-720", 1))
		mustProjection(t, timed, "TimeControlled")
		manualCode := Apply(heatingMode, codeAt(heatingMode, 2, "zone-east", "controller-720", 1))
		mustProjection(t, manualCode, "Manual")
		if manualCode.Value != nil {
			t.Fatal("zone manual became a pack token")
		}
		coolingOff := Apply(coolingMode, codeAt(coolingMode, 0, "zone-east", "controller-720", 1))
		mustSymbol(t, coolingOff, "thermal.mode.zone", "off", "thermal.dimension.zone", "zone-east")
		swapped := Apply(heatingMode, codeAt(heatingMode, 1, "zone-east", "controller-720", 1))
		if swapped.Status == StatusFed {
			t.Fatal("family 720 code 1 used the other family table")
		}
	})

	t.Run("setpoints that stay native", func(t *testing.T) {
		got := Apply(effective, textAt(effective, "21", "zone-east", "controller-700", 3))
		mustNamespace(t, got, "21")
		if got.Key.FactID != "" || strings.Contains(got.Evidence, "unit.kelvin") {
			t.Fatalf("effective key %s evidence %s", got.Key.FactID, got.Evidence)
		}
		requireEvidence(t, got, "unit=°C", "native=21")
		comfortGot := Apply(comfort, textAt(comfort, "20", "zone-east", "controller-700", 0))
		mustNamespace(t, comfortGot, "20")
		if comfortGot.Value != nil {
			t.Fatal("comfort setpoint published occupied heating")
		}
		other := Apply(comfortOther, textAt(comfortOther, "20", "zone-east", "controller-720", 0))
		mustNamespace(t, other, "20")
		if Apply(comfort, textAt(comfort, "20", "zone-east", "controller-720", 0)).Status != StatusUnmatched {
			t.Fatal("family 700 comfort row matched the other family")
		}
		nativeCooling := Apply(coolingNative, textAt(coolingNative, "19", "zone-east", "controller-700", 0))
		mustNamespace(t, nativeCooling, "19")
		if nativeCooling.Key.FactID == "thermal.setpoint.unoccupied_cooling" || nativeCooling.Key.FactID == "thermal.setpoint.occupied_cooling" {
			t.Fatal("cooling setpoint gained an unoccupied or occupied role")
		}
		specialGot := Apply(special, textAt(special, "3", "zone-east", "controller-700", 0))
		mustNamespace(t, specialGot, "3")
	})

	t.Run("cooling rank", func(t *testing.T) {
		manualGot := Apply(coolingManual, textAt(coolingManual, "21", "zone-east", "controller-700", 4))
		effectiveGot := Apply(coolingEffective, textAt(coolingEffective, "18", "zone-east", "controller-700", 4))
		mustQuantity(t, manualGot, "thermal.setpoint.occupied_cooling", "thermal.dimension.zone", "zone-east", "29415", -2, "unit.kelvin")
		mustQuantity(t, effectiveGot, "thermal.setpoint.occupied_cooling", "thermal.dimension.zone", "zone-east", "29115", -2, "unit.kelvin")
		selected := Select([]Outcome{manualGot, effectiveGot})
		kept := fedFor(selected, "thermal.setpoint.occupied_cooling")
		if len(kept) != 1 || kept[0].Rank != sourceRankPrimary {
			t.Fatalf("selected %#v", kept)
		}
		stale := manualGot
		stale.Freshness = "stale"
		staleReading := textAt(coolingManual, "21", "zone-east", "controller-700", 4)
		staleReading.Stale = true
		stale = Apply(coolingManual, staleReading)
		if stale.Freshness != "stale" || stale.Status != StatusFed {
			t.Fatalf("stale %+v", stale)
		}
		blocked := Select([]Outcome{stale, effectiveGot})
		foundStale := false
		for _, outcome := range blocked {
			if outcome.Key.FactID != "thermal.setpoint.occupied_cooling" {
				continue
			}
			if outcome.Freshness != "stale" || outcome.Value == nil || outcome.Value.Quantity.Number.Coefficient != "29415" {
				t.Fatal("fresh lower source replaced a stale higher source")
			}
			foundStale = true
		}
		if !foundStale {
			t.Fatal("stale higher source was dropped")
		}
		failedReading := textAt(coolingManual, "21", "zone-east", "controller-700", 4)
		failedReading.Failed = true
		failed := Apply(coolingManual, failedReading)
		if failed.Status != StatusUnknown || failed.Value != nil || failed.Key.FactID == "" {
			t.Fatalf("failed %+v", failed)
		}
		if len(fedFor(Select([]Outcome{failed, effectiveGot}), "thermal.setpoint.occupied_cooling")) != 0 {
			t.Fatal("lower source replaced an unknown higher source")
		}
		absent := textAt(coolingManual, "21", "zone-east", "controller-700", 4)
		absent.Present = false
		unknownAbsent := Apply(coolingManual, absent)
		if unknownAbsent.Status != StatusUnknown || unknownAbsent.Value != nil || unknownAbsent.Key.FactID == "" {
			t.Fatalf("absent %+v", unknownAbsent)
		}
		keptAbsent := Select([]Outcome{unknownAbsent, effectiveGot})
		foundAbsent := false
		for _, outcome := range keptAbsent {
			if outcome.Key.FactID != "thermal.setpoint.occupied_cooling" {
				continue
			}
			if outcome.Status != StatusUnknown || outcome.Value != nil || outcome.Rank != sourceRankPrimary {
				t.Fatal("absent higher source did not keep the fact unknown")
			}
			foundAbsent = true
		}
		if !foundAbsent {
			t.Fatal("absent higher source was dropped")
		}
		if len(fedFor(keptAbsent, "thermal.setpoint.occupied_cooling")) != 0 {
			t.Fatal("lower source replaced an absent higher source")
		}
	})

	t.Run("missing zone", func(t *testing.T) {
		gone := textAt(air, "0", "zone-east", "controller-700", 0)
		gone.Missing = true
		withdrawn := Apply(air, gone)
		if withdrawn.Status != StatusWithdrawn || withdrawn.Value != nil {
			t.Fatalf("missing zone %+v", withdrawn)
		}
		absentReply := textAt(air, "0", "zone-east", "controller-700", 0)
		absentReply.Present = false
		if got := Apply(air, absentReply); got.Status != StatusUnknown || got.Value != nil || got.Key.FactID == "" {
			t.Fatalf("absent zone reply %+v", got)
		}
		valveReading := textAt(valve, "1.0", "terminal-east", "controller-700", 5)
		valveReading.Missing = true
		if got := Apply(valve, valveReading); got.Status != StatusFed {
			t.Fatalf("missing flag withdrew a terminal %+v", got)
		}
	})

	t.Run("mode conflict and projection", func(t *testing.T) {
		heatOff := Apply(heatingMode, codeAt(heatingMode, 0, "zone-east", "controller-700", 1))
		coolAuto := Apply(coolingMode, codeAt(coolingMode, 1, "zone-east", "controller-700", 1))
		mixed := Select([]Outcome{heatOff, coolAuto})
		conflict := false
		for _, outcome := range mixed {
			if outcome.Status == StatusConflict && outcome.Value == nil && outcome.Key.FactID == "thermal.mode.zone" {
				conflict = true
			}
		}
		if !conflict {
			t.Fatal("equal ranks with different tokens published one mode")
		}
		coolOff := Apply(coolingMode, codeAt(coolingMode, 0, "zone-east", "controller-700", 1))
		same := Select([]Outcome{heatOff, coolOff})
		kept := fedFor(same, "thermal.mode.zone")
		if len(kept) != 1 || kept[0].Value.Symbol.Token != "off" {
			t.Fatal("equal tokens were not kept as one fact")
		}
		day := Apply(heatingMode, codeAt(heatingMode, 2, "zone-east", "controller-700", 1))
		lower := day
		lower.Rank = sourceRankSecondary
		lower.Status = StatusFed
		lower.Freshness = "fresh"
		lower.Value = symbolValue("thermal.mode.zone", "auto")
		lower.Gap = ""
		if len(fedFor(Select([]Outcome{day, lower}), "thermal.mode.zone")) != 0 {
			t.Fatal("lower source replaced a projection")
		}
	})

	t.Run("valve and setback", func(t *testing.T) {
		open := Apply(valve, textAt(valve, "1.0", "terminal-east", "controller-700", 5))
		mustQuantity(t, open, "thermal.measurement.valve_position", "thermal.dimension.terminal_unit", "terminal-east", "1", 0, "unit.ratio")
		closed := Apply(valve, textAt(valve, "0", "terminal-east", "controller-700", 5))
		mustQuantity(t, closed, "thermal.measurement.valve_position", "thermal.dimension.terminal_unit", "terminal-east", "0", 0, "unit.ratio")
		if got := Apply(valve, textAt(valve, "2", "terminal-east", "controller-700", 5)); got.Status != StatusUnchanged || got.Value != nil {
			t.Fatalf("valve %+v", got)
		}
		got := Apply(setback, textAt(setback, "17", "zone-east", "controller-720", 6))
		mustQuantity(t, got, "thermal.setpoint.unoccupied_heating", "thermal.dimension.zone", "zone-east", "29015", -2, "unit.kelvin")
		broken := textAt(setback, "17 C", "zone-east", "controller-700", 6)
		if got := Apply(setback, broken); got.Status != StatusUnchanged || got.Value != nil {
			t.Fatalf("malformed %+v", got)
		}
	})

	t.Run("manual failure", func(t *testing.T) {
		failed := textAt(manual, "21", "zone-east", "controller-700", 0)
		failed.Failed = true
		got := Apply(manual, failed)
		if got.Status != StatusUnknown || got.Value != nil {
			t.Fatalf("failed manual %+v", got)
		}
		beside := Apply(effective, textAt(effective, "21", "zone-east", "controller-700", 0))
		for _, outcome := range Select([]Outcome{got, beside}) {
			if outcome.Status == StatusFed && outcome.Key.FactID == "thermal.setpoint.temperature_setpoint" {
				t.Fatal("effective setpoint replaced an unknown manual setpoint")
			}
		}
	})
}

func TestThermalDhwFeed(t *testing.T) {
	feeds := DhwFeeds()
	temperature := mustFeed(t, feeds, "dhw.temperature")
	setpoint := mustFeed(t, feeds, "dhw.setpoint")
	pump := mustFeed(t, feeds, "dhw.pump_active")
	flow := mustFeed(t, feeds, "dhw.flow_setpoint")
	special := mustFeed(t, feeds, "dhw.special_function")
	mode := mustFeed(t, feeds, "dhw.mode")
	top := mustFeed(t, feeds, "store.top_native")
	bottom := mustFeed(t, feeds, "store.bottom_native")
	first := mustFeed(t, feeds, "store.differential_1")
	second := mustFeed(t, feeds, "store.differential_2")

	got := Apply(temperature, textAt(temperature, "48.5", "hot-water", "controller-700", 0))
	mustQuantity(t, got, "thermal.measurement.dhw_temperature", "thermal.dimension.dhw", "hot-water", "32165", -2, "unit.kelvin")
	if got.Grade != "Hypothesis" {
		t.Fatalf("grade %s", got.Grade)
	}
	proven := Apply(setpoint, textAt(setpoint, "50", "hot-water", "controller-720", 0))
	mustQuantity(t, proven, "thermal.setpoint.dhw_temperature_setpoint", "thermal.dimension.dhw", "hot-water", "32315", -2, "unit.kelvin")
	if proven.Grade != "Proven" {
		t.Fatalf("grade %s", proven.Grade)
	}
	otherResource := Apply(setpoint, textAt(setpoint, "50", "hot-water-west", "controller-720", 0))
	if dimensionText(proven.Key) == dimensionText(otherResource.Key) {
		t.Fatal("hot-water resources collided")
	}
	if Apply(setpoint, textAt(setpoint, "50", "hot-water", "controller-700", 1)).Status != StatusUnmatched {
		t.Fatal("hot-water instance 1 matched")
	}
	mustNamespace(t, Apply(pump, textAt(pump, "on", "hot-water", "controller-700", 0)), "on")
	flowGot := Apply(flow, textAt(flow, "55", "hot-water", "controller-700", 0))
	mustNamespace(t, flowGot, "55")
	if flowGot.Key.FactID == "thermal.setpoint.temperature_setpoint" || flowGot.Value != nil {
		t.Fatal("hot-water flow setpoint used the zone setpoint")
	}
	mustNamespace(t, Apply(special, textAt(special, "9", "hot-water", "controller-700", 0)), "9")
	mustSymbol(t, Apply(mode, codeAt(mode, 0, "hot-water", "controller-700", 0)), "thermal.mode.dhw", "off", "thermal.dimension.dhw", "hot-water")
	mustSymbol(t, Apply(mode, codeAt(mode, 1, "hot-water", "controller-700", 0)), "thermal.mode.dhw", "auto", "thermal.dimension.dhw", "hot-water")
	mustProjection(t, Apply(mode, codeAt(mode, 2, "hot-water", "controller-700", 0)), "Day")
	mustSymbol(t, Apply(mode, codeAt(mode, 2, "hot-water", "controller-720", 0)), "thermal.mode.dhw", "manual", "thermal.dimension.dhw", "hot-water")
	timed := Apply(mode, codeAt(mode, 1, "hot-water", "controller-720", 0))
	mustProjection(t, timed, "TimeControlled")
	if timed.Status == StatusFed {
		t.Fatal("time controlled became timed")
	}
	for _, row := range []Feed{top, bottom, first, second} {
		got := Apply(row, textAt(row, "40", "store-candidate", row.Families[0], 0))
		mustNamespace(t, got, "40")
		if got.Key.FactID == "thermal.measurement.store_temperature" || got.Value != nil {
			t.Fatalf("%s published a store temperature", row.ID)
		}
		requireEvidence(t, got, "unit=°C", "name="+row.NativeName)
	}
	if Apply(top, textAt(top, "40", "store-candidate", "controller-700", 0)).NativeText == Apply(bottom, textAt(bottom, "41", "store-candidate", "controller-700", 0)).NativeText {
		t.Fatal("top and bottom readings collapsed")
	}
}

func TestThermalCircuitFeed(t *testing.T) {
	feeds := CircuitFeeds()
	flow := mustFeed(t, feeds, "circuit.flow_temperature")
	rate := mustFeed(t, feeds, "circuit.flow_rate")
	capability := mustFeed(t, feeds, "circuit.read_capability")
	coolingEnabled := mustFeed(t, feeds, "circuit.cooling_enabled")
	returned := mustFeed(t, feeds, "circuit.return_setpoint")
	demand := mustFeed(t, feeds, "circuit.external_demand")
	limit := mustFeed(t, feeds, "circuit.outdoor_limit")
	flowSet := mustFeed(t, feeds, "circuit.flow_setpoint")
	state := mustFeed(t, feeds, "circuit.state")
	mixer := mustFeed(t, feeds, "circuit.mixer")

	got := Apply(flow, textAt(flow, "35", "circuit-east", "controller-700", 0))
	mustQuantity(t, got, "thermal.measurement.circuit_flow_temperature", "thermal.dimension.circuit", "circuit-east", "30815", -2, "unit.kelvin")
	other := Apply(flow, textAt(flow, "35", "circuit-west", "controller-720", 10))
	mustQuantity(t, other, "thermal.measurement.circuit_flow_temperature", "thermal.dimension.circuit", "circuit-west", "30815", -2, "unit.kelvin")
	if dimensionText(got.Key) == dimensionText(other.Key) {
		t.Fatal("circuits collided")
	}
	high := textAt(flow, "35", "circuit-east", "controller-700", 11)
	if Apply(flow, high).Status != StatusUnmatched {
		t.Fatal("instance 11 matched")
	}
	if Apply(rate, textAt(flow, "1", "circuit-east", "controller-700", 0)).Status != StatusUnmatched {
		t.Fatal("missing flow-rate register matched a reading")
	}
	if !gapNamed(feeds, "thermal.measurement.flow_rate has no documented circuit register") {
		t.Fatal("flow rate gap missing")
	}
	active := Apply(capability, textAt(capability, "35", "circuit-east", "controller-700", 2))
	if active.Status != StatusFed || active.Value != nil || active.Key.FactID != "" {
		t.Fatalf("capability %+v", active)
	}
	if active.NativeText != "thermal.capability.read.circuit" {
		t.Fatalf("capability text %s", active.NativeText)
	}
	off := textAt(flow, "35", "circuit-east", "controller-700", 2)
	off.Deactivated = true
	if got := Apply(flow, off); got.Status != StatusWithdrawn || got.Value != nil {
		t.Fatalf("deactivated flow %+v", got)
	}
	missingCircuit := textAt(flow, "35", "circuit-east", "controller-700", 2)
	missingCircuit.Missing = true
	if got := Apply(flow, missingCircuit); got.Status != StatusFed {
		t.Fatalf("missing flag withdrew a circuit %+v", got)
	}
	absentCircuit := textAt(flow, "35", "circuit-east", "controller-700", 2)
	absentCircuit.Present = false
	if got := Apply(flow, absentCircuit); got.Status != StatusUnknown || got.Value != nil || got.Key.FactID == "" {
		t.Fatalf("absent circuit reply %+v", got)
	}
	if got := Apply(capability, off); got.Status != StatusWithdrawn {
		t.Fatalf("deactivated capability %+v", got)
	}
	offReturn := textAt(returned, "30", "circuit-east", "controller-700", 2)
	offReturn.Deactivated = true
	if got := Apply(returned, offReturn); got.Status != StatusWithdrawn {
		t.Fatalf("deactivated return %+v", got)
	}
	returnGot := Apply(returned, textAt(returned, "30", "circuit-east", "controller-700", 2))
	mustNamespace(t, returnGot, "30")
	if returnGot.Key.FactID == "thermal.setpoint.temperature_setpoint" || returnGot.Value != nil {
		t.Fatal("return setpoint became a zone setpoint")
	}
	demandGot := Apply(demand, textAt(demand, "4", "circuit-east", "controller-700", 2))
	mustNamespace(t, demandGot, "4")
	if demandGot.Key.FactID == "thermal.demand.level" || (demandGot.Value != nil && demandGot.Value.Quantity != nil) {
		t.Fatal("external demand became a ventilation ratio")
	}
	limitGot := Apply(limit, textAt(limit, "15", "circuit-east", "controller-720", 2))
	mustNamespace(t, limitGot, "15")
	if limitGot.Key.FactID != "" || limitGot.Value != nil {
		t.Fatal("outdoor limit became a room setpoint")
	}
	mustNamespace(t, Apply(flowSet, textAt(flowSet, "40", "circuit-east", "controller-700", 2)), "40")
	mustNamespace(t, Apply(state, textAt(state, "Heating", "circuit-east", "controller-700", 2)), "Heating")
	mixerGot := Apply(mixer, textAt(mixer, "50", "circuit-east", "controller-700", 2))
	mustNamespace(t, mixerGot, "50")
	if mixerGot.Value != nil {
		t.Fatal("mixer percent became a ratio")
	}
	requireEvidence(t, mixerGot, "unit=%", "native=50")
	if !gapNamed([]Feed{coolingEnabled}, "circuit cooling enabled has no SemReg element") {
		t.Fatal("cooling enabled gap missing")
	}
	mustNamespace(t, Apply(coolingEnabled, textAt(coolingEnabled, "on", "circuit-east", "controller-700", 2)), "on")
	if token, ok := DeriveCircuitAction("Heating"); !ok || token != "heating" {
		t.Fatalf("heating action %s %v", token, ok)
	}
	if _, ok := DeriveCircuitAction("Standby"); ok {
		t.Fatal("undocumented standby label mapped")
	}
}

func TestThermalSystemFeed(t *testing.T) {
	feeds := SystemFeeds()
	outdoor := mustFeed(t, feeds, "system.outdoor_air")
	flow := mustFeed(t, feeds, "system.flow_temperature")
	average := mustFeed(t, feeds, "system.outdoor_average")
	pressure := mustFeed(t, feeds, "system.pressure")
	switched := mustFeed(t, feeds, "system.switch")
	maintenance := mustFeed(t, feeds, "system.maintenance")
	emergency := mustFeed(t, feeds, "system.emergency")
	yield := mustFeed(t, feeds, "system.environmental_yield")
	recovery := mustFeed(t, feeds, "system.heat_recovery_yield")
	curve := mustFeed(t, feeds, "system.heating_curve")
	hybrid := mustFeed(t, feeds, "system.hybrid_limit")
	holiday := mustFeed(t, feeds, "system.holiday_start")
	holidayOther := mustFeed(t, feeds, "system.holiday_start_other")
	ventilation := mustFeed(t, feeds, "system.ventilation_mode")
	ventilationOther := mustFeed(t, feeds, "system.ventilation_mode_other")

	air := Apply(outdoor, textAt(outdoor, "0", "system-one", "controller-700", 0))
	mustQuantity(t, air, "thermal.measurement.outdoor_air_temperature", "thermal.dimension.system", "system-one", "27315", -2, "unit.kelvin")
	mean := Apply(average, textAt(average, "0", "system-one", "controller-700", 0))
	mustNamespace(t, mean, "0")
	if mean.Key.FactID == "thermal.measurement.outdoor_air_temperature" || mean.Value != nil {
		t.Fatal("rolling average filled the outdoor temperature")
	}
	systemFlow := Apply(flow, textAt(flow, "42", "system-one", "controller-720", 0))
	mustQuantity(t, systemFlow, "thermal.measurement.system_flow_temperature", "thermal.dimension.system", "system-one", "31515", -2, "unit.kelvin")
	if systemFlow.Key.FactID == air.Key.FactID {
		t.Fatal("system flow used the outdoor fact")
	}
	pressureGot := Apply(pressure, textAt(pressure, "1.2", "system-one", "controller-700", 0))
	mustNamespace(t, pressureGot, "1.2")
	if pressureGot.Value != nil {
		t.Fatal("system pressure published pascal")
	}
	requireEvidence(t, pressureGot, "unit=bar", "native=1.2")
	off := Apply(switched, textAt(switched, "system off", "system-one", "controller-700", 0))
	mustNamespace(t, off, "system off")
	if off.Key.FactID == "thermal.mode.system" || off.Value != nil {
		t.Fatal("system switch became a mode")
	}
	due := Apply(maintenance, textAt(maintenance, "1", "system-one", "controller-700", 0))
	mustSymbol(t, due, "thermal.status.maintenance", "due", "thermal.dimension.system", "system-one")
	clear := Apply(maintenance, textAt(maintenance, "false", "system-one", "controller-700", 0))
	mustSymbol(t, clear, "thermal.status.maintenance", "clear", "thermal.dimension.system", "system-one")
	codeClear := codeAt(maintenance, 0, "system-one", "controller-700", 0)
	codeClear.Text = ""
	mustSymbol(t, Apply(maintenance, codeClear), "thermal.status.maintenance", "clear", "thermal.dimension.system", "system-one")
	unknownCode := codeAt(maintenance, 2, "system-one", "controller-700", 0)
	unknownCode.Text = ""
	mustProjection(t, Apply(maintenance, unknownCode), "2")
	emergencyGot := Apply(emergency, textAt(emergency, "1", "system-one", "controller-700", 0))
	mustNamespace(t, emergencyGot, "1")
	if emergencyGot.Key.FactID == "thermal.status.fault" {
		t.Fatal("emergency operation became a fault")
	}
	yieldGot := Apply(yield, textAt(yield, "12", "system-one", "controller-700", 0))
	mustNamespace(t, yieldGot, "12")
	if yieldGot.Key.FactID == "thermal.energy.heat_output" || yieldGot.Value != nil {
		t.Fatal("yield became heat output")
	}
	requireEvidence(t, yieldGot, "unit=kWh")
	recoveryGot := Apply(recovery, textAt(recovery, "3", "system-one", "controller-700", 0))
	mustNamespace(t, recoveryGot, "3")
	if strings.Contains(recoveryGot.Evidence, "unit=kWh") || recoveryGot.Value != nil {
		t.Fatal("heat recovery gained a kilowatt-hour unit")
	}
	for _, row := range []Feed{curve, hybrid, holiday, holidayOther, ventilation, ventilationOther} {
		got := Apply(row, textAt(row, "1", "system-one", row.Families[0], 0))
		mustNamespace(t, got, "1")
		if got.Value != nil {
			t.Fatalf("%s published a pack fact", row.ID)
		}
	}
	if Apply(holidayOther, textAt(holidayOther, "1", "system-one", "controller-720", 0)).Status != StatusUnmatched {
		t.Fatal("the other holiday location matched the wrong family")
	}
	if Apply(ventilation, textAt(ventilation, "1", "system-one", "controller-700", 0)).Status != StatusUnmatched {
		t.Fatal("ventilation location matched the wrong family")
	}

	t.Run("derived", func(t *testing.T) {
		operation, action := DeriveSystem(nil, "system-one")
		if operation.Status != StatusUnknown || action.Status != StatusUnknown || operation.Value != nil || action.Value != nil {
			t.Fatal("empty circuit states became idle")
		}
		requireEvidence(t, operation, "origin=derived")
		operation, action = DeriveSystem([]string{"StandBy", "StandBy"}, "system-one")
		mustSymbol(t, operation, "thermal.status.operation", "idle", "thermal.dimension.system", "system-one")
		mustSymbol(t, action, "thermal.action.state", "idle", "thermal.dimension.system", "system-one")
		operation, action = DeriveSystem([]string{"Heating"}, "system-west")
		mustSymbol(t, operation, "thermal.status.operation", "active", "thermal.dimension.system", "system-west")
		mustSymbol(t, action, "thermal.action.state", "heating", "thermal.dimension.system", "system-west")
		if dimensionText(operation.Key) == "system-one" {
			t.Fatal("derived resource text was fixed")
		}
		operation, action = DeriveSystem([]string{"Cooling", "StandBy"}, "system-one")
		mustSymbol(t, action, "thermal.action.state", "cooling", "thermal.dimension.system", "system-one")
		operation, action = DeriveSystem([]string{"DHW"}, "system-one")
		mustSymbol(t, action, "thermal.action.state", "dhw", "thermal.dimension.system", "system-one")
		operation, action = DeriveSystem([]string{"Heating", "Cooling"}, "system-one")
		mustSymbol(t, operation, "thermal.status.operation", "active", "thermal.dimension.system", "system-one")
		if action.Status != StatusUnknown || action.Value != nil || action.Gap == "" {
			t.Fatalf("mixed action %+v", action)
		}
		operation, action = DeriveSystem([]string{"StandBy", "later"}, "system-one")
		if operation.Status != StatusUnknown || operation.Value != nil || action.Status != StatusUnknown {
			t.Fatal("an unknown label became idle")
		}
		operation, _ = DeriveSystem([]string{"Heating", "later"}, "system-one")
		mustSymbol(t, operation, "thermal.status.operation", "active", "thermal.dimension.system", "system-one")
	})

	t.Run("fault", func(t *testing.T) {
		present := DeriveFault(ErrorRead{Known: true, Set: true}, ErrorRead{Known: false}, "system-one")
		mustSymbol(t, present, "thermal.status.fault", "present", "thermal.dimension.system", "system-one")
		present = DeriveFault(ErrorRead{Known: false}, ErrorRead{Known: true, Set: true}, "system-west")
		mustSymbol(t, present, "thermal.status.fault", "present", "thermal.dimension.system", "system-west")
		clear := DeriveFault(ErrorRead{Known: true}, ErrorRead{Known: true}, "system-one")
		mustSymbol(t, clear, "thermal.status.fault", "clear", "thermal.dimension.system", "system-one")
		unknown := DeriveFault(ErrorRead{Known: false}, ErrorRead{Known: true}, "system-one")
		if unknown.Status != StatusUnknown || unknown.Value != nil {
			t.Fatal("one unknown side became clear")
		}
		unknown = DeriveFault(ErrorRead{}, ErrorRead{}, "system-one")
		if unknown.Status != StatusUnknown || unknown.Value != nil {
			t.Fatal("two unknown sides became clear")
		}
	})
}

func TestThermalOverrunFeed(t *testing.T) {
	feeds := OverrunFeeds()
	zone := mustFeed(t, feeds, "overrun.zone")
	circuit := mustFeed(t, feeds, "overrun.circuit")
	dhw := mustFeed(t, feeds, "overrun.dhw")
	ventilation := mustFeed(t, feeds, "overrun.ventilation")
	ventilationOther := mustFeed(t, feeds, "overrun.ventilation_other")
	for _, id := range []string{"overrun.one_time_ventilation", "overrun.one_day_away", "overrun.one_day_at_home"} {
		row := mustFeed(t, feeds, id)
		if Apply(row, Reading{Present: true, Text: "1", Resource: "system-one"}).Status != StatusUnmatched {
			t.Fatalf("%s matched a reading", id)
		}
		if row.Gap == "" {
			t.Fatalf("%s has no gap", id)
		}
	}
	party := Apply(zone, codeAt(zone, 3, "system-one", "controller-700", 1))
	mustSymbol(t, party, "thermal.overrun.party", "active", "thermal.dimension.system", "system-one")
	if party.Namespace != NamespaceEbusreg {
		t.Fatal("party dropped the native copy")
	}
	requireEvidence(t, party, "token=active", "code=3")
	if Apply(zone, codeAt(zone, 3, "system-one", "controller-720", 1)).Status != StatusProjection {
		t.Fatal("family 720 gained a party code")
	}
	off := Apply(zone, codeAt(zone, 10, "system-one", "controller-720", 2))
	mustSymbol(t, off, "thermal.overrun.hvac_system_off", "active", "thermal.dimension.system", "system-one")
	boost := Apply(zone, codeAt(zone, 1, "system-one", "controller-700", 0))
	mustProjection(t, boost, "VentilationBoost")
	if boost.Key.FactID == "thermal.overrun.one_time_ventilation" || boost.Value != nil {
		t.Fatal("ventilation boost became one-time ventilation")
	}
	none := Apply(zone, codeAt(zone, 0, "system-one", "controller-700", 0))
	if none.Status != StatusNone || none.Value != nil || none.NativeText != "NoQuickMode" {
		t.Fatalf("no quick mode %+v", none)
	}
	raw := Apply(zone, codeAt(zone, 99, "system-one", "controller-700", 0))
	mustProjection(t, raw, "99")
	if raw.Value != nil {
		t.Fatal("undocumented code published an overrun")
	}
	if Apply(zone, codeAt(zone, 3, "system-one", "controller-700", 11)).Status != StatusUnmatched {
		t.Fatal("zone overrun instance 11 matched")
	}
	cylinder := Apply(circuit, codeAt(circuit, 6, "system-one", "controller-700", 0))
	mustProjection(t, cylinder, "cylinder boost")
	if cylinder.Key.FactID == "thermal.overrun.one_time_dhw" || cylinder.Value != nil {
		t.Fatal("family 700 cylinder boost became one-time hot water")
	}
	loading := Apply(circuit, codeAt(circuit, 6, "system-one", "controller-720", 0))
	mustSymbol(t, loading, "thermal.overrun.one_time_dhw", "active", "thermal.dimension.system", "system-one")
	if loading.Namespace != NamespaceEbusreg {
		t.Fatal("circuit overrun dropped the native copy")
	}
	quiet := Apply(circuit, codeAt(circuit, 0, "system-one", "controller-720", 0))
	if quiet.Status != StatusNone || quiet.Value != nil {
		t.Fatalf("no quick mode %+v", quiet)
	}
	away := Apply(circuit, codeAt(circuit, 4, "system-one", "controller-700", 0))
	mustProjection(t, away, "one-day holiday")
	if away.Key.FactID == "thermal.overrun.one_day_away" || away.Key.FactID == "thermal.overrun.one_day_at_home" {
		t.Fatal("one-day holiday chose away or home")
	}
	tank := Apply(dhw, codeAt(dhw, 9, "system-one", "controller-700", 0))
	mustSymbol(t, tank, "thermal.overrun.one_time_dhw", "active", "thermal.dimension.system", "system-one")
	if Apply(dhw, codeAt(dhw, 9, "system-west", "controller-720", 0)).Key.FactID != tank.Key.FactID {
		t.Fatal("families disagree on one-time tank loading")
	}
	if dimensionText(tank.Key) == dimensionText(Apply(dhw, codeAt(dhw, 9, "system-west", "controller-720", 0)).Key) {
		t.Fatal("overrun resources collided")
	}
	vent := Apply(ventilation, codeAt(ventilation, 3, "system-one", "controller-700", 0))
	mustSymbol(t, vent, "thermal.overrun.party", "active", "thermal.dimension.system", "system-one")
	named := Apply(ventilation, codeAt(ventilation, 1, "system-one", "controller-700", 0))
	mustProjection(t, named, "Ventilation boost")
	if Apply(ventilation, codeAt(ventilation, 10, "system-one", "controller-720", 0)).Status != StatusUnmatched {
		t.Fatal("family 700 ventilation row matched the other family")
	}
	otherOff := Apply(ventilationOther, codeAt(ventilationOther, 10, "system-one", "controller-720", 3))
	mustSymbol(t, otherOff, "thermal.overrun.hvac_system_off", "active", "thermal.dimension.system", "system-one")
	if Apply(ventilationOther, codeAt(ventilationOther, 1, "system-one", "controller-720", 3)).Value != nil {
		t.Fatal("the other ventilation boost published an overrun")
	}
}

func TestThermalHeatPumpFeed(t *testing.T) {
	feeds := HeatPumpFeeds()
	high := mustFeed(t, feeds, "heatpump.high_pressure")
	low := mustFeed(t, feeds, "heatpump.low_pressure")
	interfaceHigh := mustFeed(t, feeds, "heatpump.high_pressure_interface")
	condensing := mustFeed(t, feeds, "heatpump.condensing")
	evaporating := mustFeed(t, feeds, "heatpump.evaporating")
	speed := mustFeed(t, feeds, "heatpump.compressor_speed")
	state := mustFeed(t, feeds, "heatpump.compressor_state")
	output := mustFeed(t, feeds, "heatpump.heat_output")
	power := mustFeed(t, feeds, "heatpump.thermal_power")

	number, ok := parseNativeDecimal("1")
	if !ok {
		t.Fatal("decimal 1")
	}
	converted, unit, ok := convertNative("bar", number)
	if !ok || converted.Coefficient != "1" || converted.Exponent10 != 5 || unit != "unit.pascal" {
		t.Fatalf("bar scale %+v %s %v", converted, unit, ok)
	}
	reading := textAt(high, "1.5", "refrigerant-east", "heat-pump", 0)
	reading.Group = 0x09
	reading.Register = 0x0001
	got := Apply(high, reading)
	mustNamespace(t, got, "1.5")
	if got.Value != nil || got.Key.FactID != "thermal.measurement.refrigerant_high_side_pressure" {
		t.Fatalf("high pressure %+v", got)
	}
	requireEvidence(t, got, "unit=bar", "native=1.5", "grade=Hypothesis")
	lowGot := Apply(low, textAt(low, "0.4", "refrigerant-east", "heat-pump", 0))
	mustNamespace(t, lowGot, "0.4")
	if lowGot.Key.FactID != "thermal.measurement.refrigerant_low_side_pressure" || lowGot.Value != nil {
		t.Fatal("low pressure published pascal")
	}
	second := Apply(interfaceHigh, textAt(interfaceHigh, "1.5", "refrigerant-east", "heat-pump", 0))
	mustNamespace(t, second, "1.5")
	if second.Rank != sourceRankSecondary || second.Value != nil {
		t.Fatal("interface pressure published pascal")
	}
	selected := Select([]Outcome{got, second})
	ranked := 0
	for _, outcome := range selected {
		if outcome.Key.FactID == "thermal.measurement.refrigerant_high_side_pressure" {
			ranked++
			if outcome.Rank != sourceRankPrimary || outcome.Value != nil {
				t.Fatalf("selected pressure %+v", outcome)
			}
		}
	}
	if ranked != 1 {
		t.Fatalf("pressure sources %d", ranked)
	}
	west := textAt(high, "1.5", "refrigerant-west", "heat-pump", 0)
	if dimensionText(Apply(high, west).Key) == dimensionText(got.Key) {
		t.Fatal("refrigerant resources collided")
	}
	if Apply(high, textAt(high, "1.5", "refrigerant-east", "heat-pump", 1)).Status != StatusUnmatched {
		t.Fatal("selector instance 1 matched")
	}
	wrong := textAt(high, "1.5", "refrigerant-east", "heat-pump", 0)
	wrong.Selector = "15 00"
	if Apply(high, wrong).Status != StatusUnmatched {
		t.Fatal("high pressure matched the low selector")
	}
	failed := textAt(high, "1.5", "refrigerant-east", "heat-pump", 0)
	failed.Failed = true
	failedGot := Apply(high, failed)
	if failedGot.Status != StatusUnknown || failedGot.Value != nil {
		t.Fatalf("failed pressure %+v", failedGot)
	}
	condensingGot := Apply(condensing, textAt(condensing, "40", "refrigerant-east", "heat-pump", 0))
	mustNamespace(t, condensingGot, "40")
	if condensingGot.Value != nil || condensingGot.Grade != "Unknown" {
		t.Fatal("condensing temperature published kelvin")
	}
	evaporatingGot := Apply(evaporating, textAt(evaporating, "5", "refrigerant-east", "heat-pump", 0))
	mustNamespace(t, evaporatingGot, "5")
	if evaporatingGot.Value != nil || evaporatingGot.Key.FactID != "thermal.measurement.refrigerant_evaporating_temperature" {
		t.Fatal("evaporating temperature published kelvin")
	}
	speedGot := Apply(speed, textAt(speed, "1200", "compressor-east", "heat-pump", 0))
	mustNamespace(t, speedGot, "1200")
	if speedGot.Value != nil || speedGot.Key.FactID == "thermal.measurement.compressor_speed" {
		t.Fatal("compressor speed invented a unit")
	}
	stateGot := Apply(state, codeAt(state, 1, "generator-east", "heat-pump", 0))
	mustNamespace(t, stateGot, "")
	if stateGot.Key.FactID == "thermal.status.heat_generator_operation" || stateGot.Value != nil {
		t.Fatal("compressor state became heat-generator operation")
	}
	if Apply(output, Reading{Present: true, Text: "10", Resource: "generator-east"}).Status != StatusUnmatched {
		t.Fatal("heat output matched a reading")
	}
	if Apply(power, Reading{Present: true, Text: "10", Resource: "generator-east"}).Status != StatusUnmatched {
		t.Fatal("thermal power matched a reading")
	}
	if !gapNamed(feeds, "thermal.energy.heat_output has no documented register") || !gapNamed(feeds, "thermal.measurement.thermal_power has no documented register") {
		t.Fatal("unfed heat-pump elements missing")
	}
}

func mustFeed(t *testing.T, feeds []Feed, id string) Feed {
	t.Helper()
	feed, ok := FindFeed(feeds, id)
	if !ok {
		t.Fatalf("missing %s", id)
	}
	return feed
}

func textAt(feed Feed, text, resource, family string, instance byte) Reading {
	reading := ReadingOf(feed, text, resource)
	reading.Family = family
	reading.Instance = instance
	return reading
}

func codeAt(feed Feed, code int, resource, family string, instance byte) Reading {
	reading := CodeReading(feed, code, resource)
	reading.Family = family
	reading.Instance = instance
	return reading
}

func mustQuantity(t *testing.T, out Outcome, element, dimension, resource, coefficient string, exponent int32, unit string) {
	t.Helper()
	if out.Status != StatusFed || out.Value == nil || out.Value.Quantity == nil {
		t.Fatalf("quantity status %s gap %s", out.Status, out.Gap)
	}
	number := out.Value.Quantity.Number
	if number.Coefficient != coefficient || number.Exponent10 != exponent || string(out.Value.Quantity.Unit) != unit {
		t.Fatalf("number %s e%d unit %s", number.Coefficient, number.Exponent10, out.Value.Quantity.Unit)
	}
	assertKey(t, out.Key, element, dimension, resource)
	if err := thermal.New().ValidateFact(out.Key, out.Value); err != nil {
		t.Fatal(err)
	}
}

func mustSymbol(t *testing.T, out Outcome, element, token, dimension, resource string) {
	t.Helper()
	if out.Status != StatusFed || out.Value == nil || out.Value.Symbol == nil {
		t.Fatalf("symbol status %s gap %s", out.Status, out.Gap)
	}
	if string(out.Value.Symbol.Namespace) != element || out.Value.Symbol.Token != token || !out.Value.Symbol.Known {
		t.Fatalf("symbol %s %s", out.Value.Symbol.Namespace, out.Value.Symbol.Token)
	}
	assertKey(t, out.Key, element, dimension, resource)
	if err := thermal.New().ValidateFact(out.Key, out.Value); err != nil {
		t.Fatal(err)
	}
}

func assertKey(t *testing.T, key semreg.FactKey, element, dimension, resource string) {
	t.Helper()
	if string(key.PackID) != "helianthus.pack.thermal" || string(key.PackVersion) != "1.0.0" {
		t.Fatalf("pack %s %s", key.PackID, key.PackVersion)
	}
	if string(key.FactID) != element {
		t.Fatalf("fact %s", key.FactID)
	}
	if len(key.Dimensions) != 1 || string(key.Dimensions[0].ID) != dimension {
		t.Fatalf("dimension %+v", key.Dimensions)
	}
	if key.Dimensions[0].Value.Text == nil || *key.Dimensions[0].Value.Text != resource {
		t.Fatal("resource text")
	}
}

func mustNamespace(t *testing.T, out Outcome, native string) {
	t.Helper()
	if out.Status != StatusNamespace || out.Namespace != NamespaceEbusreg || out.Value != nil {
		t.Fatalf("namespace status %s value %v gap %s", out.Status, out.Value, out.Gap)
	}
	if native != "" && out.NativeText != native {
		t.Fatalf("native %s", out.NativeText)
	}
}

func mustProjection(t *testing.T, out Outcome, native string) {
	t.Helper()
	if out.Status != StatusProjection || out.Namespace != NamespaceEbusreg || out.Value != nil || out.NativeText != native {
		t.Fatalf("projection status %s native %s gap %s", out.Status, out.NativeText, out.Gap)
	}
}

func requireEvidence(t *testing.T, out Outcome, parts ...string) {
	t.Helper()
	for _, part := range parts {
		if !strings.Contains(out.Evidence, part) {
			t.Fatalf("evidence %s missing %s", out.Evidence, part)
		}
	}
}

func fedFor(outcomes []Outcome, element semreg.DefinitionID) []Outcome {
	var fed []Outcome
	for _, outcome := range outcomes {
		if outcome.Status == StatusFed && outcome.Key.FactID == element {
			fed = append(fed, outcome)
		}
	}
	return fed
}

func gapNamed(feeds []Feed, text string) bool {
	for _, feed := range Gaps(feeds) {
		if strings.Contains(feed.Gap, text) || strings.Contains(feed.Hold, text) {
			return true
		}
	}
	return false
}
