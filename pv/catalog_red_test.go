package pv

import (
	"errors"
	"testing"
)

func TestCatalogV1IsClosedAndSourceNeutral(t *testing.T) {
	catalog := CatalogV1()
	if catalog.ContractID != ContractV1 {
		t.Fatalf("ContractID = %q, want %q", catalog.ContractID, ContractV1)
	}
	want := []FactID{
		FactACActivePower,
		FactACApparentPower,
		FactACReactivePower,
		FactACPowerFactor,
		FactACFrequency,
		FactACCurrent,
		FactACVoltageLineToNeutral,
		FactACVoltageLineToLine,
		FactEnergyActiveExportTotal,
		FactDCCurrent,
		FactDCVoltage,
		FactDCActivePower,
		FactDCEnergyActiveTotal,
		FactTemperature,
		FactOperatingState,
		FactEventFlags,
		FactRatingACActivePower,
	}
	if len(catalog.Facts) != len(want) {
		t.Fatalf("Facts count = %d, want %d", len(catalog.Facts), len(want))
	}
	for _, id := range want {
		if _, ok := catalog.Facts[id]; !ok {
			t.Errorf("missing fact %q", id)
		}
	}
	if catalog.AdditiveFactsAllowed {
		t.Fatal("V1 catalog unexpectedly allows additive facts")
	}
	if catalog.WriteAuthority {
		t.Fatal("canonical PV catalog grants write authority")
	}
	for _, sourceID := range []string{
		"sunspec.phase1@1.0.0",
		"sunspec.inverter.three_phase.monitoring@1.0.0",
	} {
		if alias := catalog.SourceCompatibility[sourceID]; alias != "" {
			t.Errorf("source ID %q aliases canonical ID %q", sourceID, alias)
		}
	}
}

func TestDecimalPreservesValuesBeyondJSONSafeInteger(t *testing.T) {
	decimal, err := NewDecimal("9007199254740993", 0)
	if err != nil {
		t.Fatalf("NewDecimal: %v", err)
	}
	if decimal.Coefficient != "9007199254740993" || decimal.Scale != 0 {
		t.Fatalf("decimal = %#v", decimal)
	}
	for _, input := range []struct {
		coefficient string
		scale       int
	}{
		{"01", 0},
		{"+1", 0},
		{"1.0", 0},
		{"1", -19},
		{"1", 19},
	} {
		if _, err := NewDecimal(input.coefficient, input.scale); !errors.Is(err, ErrInvalidDecimal) {
			t.Errorf("NewDecimal(%q, %d) error = %v, want ErrInvalidDecimal", input.coefficient, input.scale, err)
		}
	}
}

func TestCatalogValidatesDimensionsUnitsAndValueDomains(t *testing.T) {
	catalog := CatalogV1()
	good := FactCandidate{
		ID:         FactACCurrent,
		Dimensions: Dimensions{Phase: PhaseL1},
		Value:      DecimalFactValue(MustDecimal("1041", -2)),
		Unit:       UnitAmpere,
	}
	if err := catalog.ValidateCandidate(good); err != nil {
		t.Fatalf("valid candidate: %v", err)
	}

	cases := []FactCandidate{
		{ID: FactID("pv.unlisted.metric"), Dimensions: Dimensions{Scope: ScopeTotal}, Value: DecimalFactValue(MustDecimal("1", 0)), Unit: UnitWatt},
		{ID: FactACFrequency, Dimensions: Dimensions{Scope: ScopeTotal}, Value: DecimalFactValue(MustDecimal("50", 0)), Unit: UnitWatt},
		{ID: FactACCurrent, Dimensions: Dimensions{Phase: Phase("L4")}, Value: DecimalFactValue(MustDecimal("1", 0)), Unit: UnitAmpere},
		{ID: FactDCCurrent, Dimensions: Dimensions{InputID: "192.168.100.21"}, Value: DecimalFactValue(MustDecimal("1", 0)), Unit: UnitAmpere},
		{ID: FactOperatingState, Dimensions: Dimensions{Scope: ScopeTotal}, Value: EnumFactValue("VENDOR_MAGIC"), Unit: UnitOne},
	}
	for index, candidate := range cases {
		if err := catalog.ValidateCandidate(candidate); !errors.Is(err, ErrInvalidFact) {
			t.Errorf("case %d error = %v, want ErrInvalidFact", index, err)
		}
	}
}

func TestThreePhaseCapabilityRequiresExactIdentities(t *testing.T) {
	catalog := CatalogV1()
	facts := requiredThreePhaseFacts(t)
	if !catalog.ThreePhaseTelemetrySatisfied(facts) {
		t.Fatal("complete three-phase fact set was not satisfied")
	}
	delete(facts, NewFactKey(FactACCurrent, Dimensions{Phase: PhaseL3}))
	if catalog.ThreePhaseTelemetrySatisfied(facts) {
		t.Fatal("capability satisfied without L3 current")
	}

	facts = requiredThreePhaseFacts(t)
	key := NewFactKey(FactACCurrent, Dimensions{Phase: PhaseL3})
	fact := facts[key]
	fact.Availability = ""
	facts[key] = fact
	if catalog.ThreePhaseTelemetrySatisfied(facts) {
		t.Fatal("capability satisfied with an invalid availability state")
	}
}
