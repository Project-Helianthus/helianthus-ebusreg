package pv

import (
	"errors"
	"fmt"
	"testing"
)

type staticResolver map[SourceIdentity]Digest

func (resolver staticResolver) ResolveSource(identity SourceIdentity) (Digest, bool) {
	digest, ok := resolver[identity]
	return digest, ok
}

func TestRegistryPartialMergePreservesFactOrigins(t *testing.T) {
	identity := SourceIdentity{
		Protocol:       "sunspec_modbus",
		ProfileID:      "sunspec.inverter.three_phase.monitoring@1.0.0",
		ProfileVersion: "1.0.0",
		Validity:       SourceTerminalVerified,
	}
	registryRef := Digest("sha256:1111111111111111111111111111111111111111111111111111111111111111")
	registry := NewRegistry(staticResolver{identity: registryRef})

	originOne := provenance(identity, registryRef, '1')
	first, err := registry.Apply(accountedUpdate(Update{
		AssetRef:  "pv-asset-01",
		Evaluated: 10,
		Source:    originOne,
		Facts: []FactInput{
			decimalInput(FactACActivePower, Dimensions{Scope: ScopeTotal}, UnitWatt, "7000", 0, PolicyTelemetryFastV1, 10),
			decimalInput(FactEnergyActiveExportTotal, Dimensions{Scope: ScopeTotal}, UnitWattHour, "10000", 0, PolicyAccumulatorV1, 10),
		},
	}))
	if err != nil {
		t.Fatalf("first Apply: %v", err)
	}
	if first.Generation != 1 {
		t.Fatalf("first generation = %d", first.Generation)
	}

	originTwo := provenance(identity, registryRef, '2')
	second, err := registry.Apply(accountedUpdate(Update{
		AssetRef:  "pv-asset-01",
		Evaluated: 40_000_000_010,
		Source:    originTwo,
		Facts: []FactInput{
			decimalInput(FactACActivePower, Dimensions{Scope: ScopeTotal}, UnitWatt, "7100", 0, PolicyTelemetryFastV1, 40_000_000_010),
		},
	}))
	if err != nil {
		t.Fatalf("second Apply: %v", err)
	}
	if second.Generation != 2 {
		t.Fatalf("second generation = %d", second.Generation)
	}
	power := lookupFact(t, second, FactACActivePower, Dimensions{Scope: ScopeTotal})
	energy := lookupFact(t, second, FactEnergyActiveExportTotal, Dimensions{Scope: ScopeTotal})
	if power.OriginRef != originTwo.SourceObservationRef {
		t.Fatalf("power origin = %s, want O2", power.OriginRef)
	}
	if energy.OriginRef != originOne.SourceObservationRef {
		t.Fatalf("retained energy origin = %s, want O1", energy.OriginRef)
	}
	if len(second.Origins) != 2 {
		t.Fatalf("origins = %d, want 2", len(second.Origins))
	}
	if len(second.RequestedOutputs) != 2 || len(second.ProjectionReport) != 2 {
		t.Fatalf("accounting requested=%d projections=%d, want 2/2", len(second.RequestedOutputs), len(second.ProjectionReport))
	}
	mappedOrigins := make(map[Digest]FactID)
	for _, projection := range second.ProjectionReport {
		if projection.Outcome != ProjectionMapped {
			t.Fatalf("projection outcome = %s", projection.Outcome)
		}
		mappedOrigins[projection.SourceRef] = projection.FactID
	}
	if mappedOrigins[originOne.SourceObservationRef] != FactEnergyActiveExportTotal || mappedOrigins[originTwo.SourceObservationRef] != FactACActivePower {
		t.Fatalf("mixed-origin accounting = %#v", mappedOrigins)
	}
	if energy.Freshness != FreshnessFresh {
		t.Fatalf("energy freshness = %s", energy.Freshness)
	}
}

func TestRegistryEvaluatesStaleExpiredAndRecoveryWithoutDeletion(t *testing.T) {
	identity := SourceIdentity{Protocol: "test_source", ProfileID: "test.profile@1.0.0", ProfileVersion: "1.0.0", Validity: SourceTerminalVerified}
	registryRef := Digest("sha256:2222222222222222222222222222222222222222222222222222222222222222")
	registry := NewRegistry(staticResolver{identity: registryRef})
	origin := provenance(identity, registryRef, '3')
	_, err := registry.Apply(accountedUpdate(Update{
		AssetRef:  "pv-asset-02",
		Evaluated: 0,
		Source:    origin,
		Facts: []FactInput{
			decimalInput(FactACActivePower, Dimensions{Scope: ScopeTotal}, UnitWatt, "1", 0, PolicyTelemetryFastV1, 0),
		},
	}))
	if err != nil {
		t.Fatal(err)
	}

	stale, err := registry.Snapshot("pv-asset-02", 30_000_000_000)
	if err != nil {
		t.Fatal(err)
	}
	fact := lookupFact(t, stale, FactACActivePower, Dimensions{Scope: ScopeTotal})
	if fact.Availability != AvailabilityAvailable || fact.Freshness != FreshnessStale {
		t.Fatalf("stale fact = %s/%s", fact.Availability, fact.Freshness)
	}

	expired, err := registry.Snapshot("pv-asset-02", 300_000_000_000)
	if err != nil {
		t.Fatal(err)
	}
	fact = lookupFact(t, expired, FactACActivePower, Dimensions{Scope: ScopeTotal})
	if fact.Availability != AvailabilityUnavailable || fact.Freshness != FreshnessExpired {
		t.Fatalf("expired fact = %s/%s", fact.Availability, fact.Freshness)
	}

	recoveredOrigin := provenance(identity, registryRef, '4')
	recovered, err := registry.Apply(accountedUpdate(Update{
		AssetRef:  "pv-asset-02",
		Evaluated: 300_000_000_001,
		Source:    recoveredOrigin,
		Facts: []FactInput{
			decimalInput(FactACActivePower, Dimensions{Scope: ScopeTotal}, UnitWatt, "2", 0, PolicyTelemetryFastV1, 300_000_000_001),
		},
	}))
	if err != nil {
		t.Fatal(err)
	}
	fact = lookupFact(t, recovered, FactACActivePower, Dimensions{Scope: ScopeTotal})
	if fact.Availability != AvailabilityAvailable || fact.Freshness != FreshnessFresh {
		t.Fatalf("recovered fact = %s/%s", fact.Availability, fact.Freshness)
	}
}

func TestRegistryFailsClosedOnUnresolvedOrMismatchedSource(t *testing.T) {
	identity := SourceIdentity{Protocol: "test_source", ProfileID: "test.profile@1.0.0", ProfileVersion: "1.0.0", Validity: SourceTerminalVerified}
	registry := NewRegistry(staticResolver{})
	input := Update{AssetRef: "pv-asset-03", Source: provenance(identity, Digest("sha256:aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"), '5')}
	if _, err := registry.Apply(input); !errors.Is(err, ErrSourceNotAdmitted) {
		t.Fatalf("unresolved source error = %v", err)
	}

	resolved := NewRegistry(staticResolver{identity: Digest("sha256:bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb")})
	if _, err := resolved.Apply(input); !errors.Is(err, ErrSourceNotAdmitted) {
		t.Fatalf("mismatched source error = %v", err)
	}
}

func TestRegistryCounterContinuityBreaksAcrossSourceIdentityOrExpiry(t *testing.T) {
	identityA := SourceIdentity{Protocol: "test_source", ProfileID: "test.profile@1.0.0", ProfileVersion: "1.0.0", Validity: SourceTerminalVerified}
	identityB := SourceIdentity{Protocol: "other_source", ProfileID: "other.profile@1.0.0", ProfileVersion: "1.0.0", Validity: SourceTerminalVerified}
	refA := Digest("sha256:aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
	refB := Digest("sha256:bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb")
	registry := NewRegistry(staticResolver{identityA: refA, identityB: refB})

	first, err := registry.Apply(accountedUpdate(Update{
		AssetRef: "pv-asset-counter-source",
		Source:   provenance(identityA, refA, '1'),
		Facts: []FactInput{
			decimalInput(FactEnergyActiveExportTotal, Dimensions{Scope: ScopeTotal}, UnitWattHour, "100", 0, PolicyAccumulatorV1, 0),
		},
	}))
	if err != nil {
		t.Fatal(err)
	}
	if continuity := lookupFact(t, first, FactEnergyActiveExportTotal, Dimensions{Scope: ScopeTotal}).Continuity; continuity == nil || continuity.State != ContinuityBaseline {
		t.Fatalf("first continuity = %#v", continuity)
	}

	second, err := registry.Apply(accountedUpdate(Update{
		AssetRef:  "pv-asset-counter-source",
		Evaluated: 1,
		Source:    provenance(identityB, refB, '2'),
		Facts: []FactInput{
			decimalInput(FactEnergyActiveExportTotal, Dimensions{Scope: ScopeTotal}, UnitWattHour, "101", 0, PolicyAccumulatorV1, 1),
		},
	}))
	if err != nil {
		t.Fatal(err)
	}
	if continuity := lookupFact(t, second, FactEnergyActiveExportTotal, Dimensions{Scope: ScopeTotal}).Continuity; continuity == nil || continuity.State != ContinuityDiscontinuity || continuity.Delta != nil {
		t.Fatalf("source-change continuity = %#v", continuity)
	}

	policy, _ := PolicyByID(PolicyAccumulatorV1)
	third, err := registry.Apply(accountedUpdate(Update{
		AssetRef:  "pv-asset-counter-source",
		Evaluated: policy.RetainFor + 2,
		Source:    provenance(identityB, refB, '3'),
		Facts: []FactInput{
			decimalInput(FactEnergyActiveExportTotal, Dimensions{Scope: ScopeTotal}, UnitWattHour, "102", 0, PolicyAccumulatorV1, policy.RetainFor+2),
		},
	}))
	if err != nil {
		t.Fatal(err)
	}
	if continuity := lookupFact(t, third, FactEnergyActiveExportTotal, Dimensions{Scope: ScopeTotal}).Continuity; continuity == nil || continuity.State != ContinuityDiscontinuity || continuity.Delta != nil {
		t.Fatalf("post-expiry continuity = %#v", continuity)
	}
}

func TestRegistryRejectsDuplicateFactsAtomicallyAndCopiesCallerValues(t *testing.T) {
	identity := SourceIdentity{Protocol: "test_source", ProfileID: "test.profile@1.0.0", ProfileVersion: "1.0.0", Validity: SourceTerminalVerified}
	registryRef := Digest("sha256:aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
	registry := NewRegistry(staticResolver{identity: registryRef})
	dimensions := Dimensions{Scope: ScopeTotal}
	duplicate := decimalInput(FactACActivePower, dimensions, UnitWatt, "1", 0, PolicyTelemetryFastV1, 0)
	if _, err := registry.Apply(accountedUpdate(Update{
		AssetRef: "pv-asset-duplicate",
		Source:   provenance(identity, registryRef, '4'),
		Facts:    []FactInput{duplicate, duplicate},
	})); !errors.Is(err, ErrInvalidFact) {
		t.Fatalf("duplicate error = %v, want ErrInvalidFact", err)
	}
	if _, err := registry.Snapshot("pv-asset-duplicate", 0); !errors.Is(err, ErrAssetNotFound) {
		t.Fatalf("duplicate update mutated registry: %v", err)
	}

	value := DecimalFactValue(MustDecimal("7", 0))
	snapshot, err := registry.Apply(accountedUpdate(Update{
		AssetRef: "pv-asset-copy",
		Source:   provenance(identity, registryRef, '5'),
		Facts: []FactInput{{
			Candidate: FactCandidate{ID: FactACActivePower, Dimensions: dimensions, Value: value, Unit: UnitWatt},
			Quality:   QualityGood,
			Policy:    PolicyTelemetryFastV1,
		}},
	}))
	if err != nil {
		t.Fatal(err)
	}
	value.Decimal.Coefficient = "999"
	fact := lookupFact(t, snapshot, FactACActivePower, dimensions)
	if fact.Value.Decimal.Coefficient != "7" {
		t.Fatalf("registry retained caller alias: %#v", fact.Value.Decimal)
	}
}

func TestRegistryRequiresCompleteProjectionAccounting(t *testing.T) {
	identity := SourceIdentity{Protocol: "test_source", ProfileID: "test.profile@1.0.0", ProfileVersion: "1.0.0", Validity: SourceTerminalVerified}
	registryRef := Digest("sha256:aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
	registry := NewRegistry(staticResolver{identity: registryRef})
	source := provenance(identity, registryRef, '6')
	requested := RequestedOutput{SourceRef: source.SourceObservationRef, RequestedOutputRef: Digest("sha256:7777777777777777777777777777777777777777777777777777777777777777")}
	dimensions := Dimensions{Scope: ScopeTotal}
	base := Update{
		AssetRef:        "pv-asset-accounting",
		SourceTimeState: SourceTimeUnavailable,
		Source:          source,
		Facts: []FactInput{
			decimalInput(FactACActivePower, dimensions, UnitWatt, "1", 0, PolicyTelemetryFastV1, 0),
		},
		Capability:       Capability{ID: CapabilityThreePhaseTelemetryV1, Outcome: CapabilityNotSatisfied},
		RequestedOutputs: []RequestedOutput{requested},
		ProjectionReport: []Projection{{
			SourceRef:          requested.SourceRef,
			RequestedOutputRef: requested.RequestedOutputRef,
			FactID:             FactACActivePower,
			Dimensions:         &dimensions,
			Outcome:            ProjectionMapped,
		}},
	}
	if _, err := registry.Apply(base); err != nil {
		t.Fatalf("complete accounting: %v", err)
	}

	missing := base
	missing.AssetRef = "pv-asset-accounting-missing"
	missing.ProjectionReport = nil
	if _, err := registry.Apply(missing); !errors.Is(err, ErrInvalidProjection) {
		t.Fatalf("missing projection error = %v", err)
	}

	mismatch := base
	mismatch.AssetRef = "pv-asset-accounting-mismatch"
	mismatch.Capability.Outcome = CapabilitySatisfied
	if _, err := registry.Apply(mismatch); !errors.Is(err, ErrInvalidCapability) {
		t.Fatalf("capability mismatch error = %v", err)
	}
}

func TestRegistryRejectsSchemaInvalidIdentityTokens(t *testing.T) {
	registryRef := Digest("sha256:aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
	for _, identity := range []SourceIdentity{
		{Protocol: "/private/fixture", ProfileID: "test.profile@1.0.0", ProfileVersion: "1.0.0", Validity: SourceTerminalVerified},
		{Protocol: "test source", ProfileID: "test.profile@1.0.0", ProfileVersion: "1.0.0", Validity: SourceTerminalVerified},
		{Protocol: "test_source", ProfileID: "/private/fixture@1.0.0", ProfileVersion: "1.0.0", Validity: SourceTerminalVerified},
		{Protocol: "test_source", ProfileID: "test.profile@1.0.0", ProfileVersion: "v1", Validity: SourceTerminalVerified},
	} {
		registry := NewRegistry(staticResolver{identity: registryRef})
		if _, err := registry.Apply(Update{AssetRef: "pv-asset-token", Source: provenance(identity, registryRef, '8')}); !errors.Is(err, ErrSourceNotAdmitted) {
			t.Errorf("identity %#v error = %v", identity, err)
		}
	}

	identity := SourceIdentity{Protocol: "test_source", ProfileID: "test.profile@1.0.0", ProfileVersion: "1.0.0", Validity: SourceTerminalVerified}
	registry := NewRegistry(staticResolver{identity: registryRef})
	for _, asset := range []string{"pv-asset-../fixture", "pv-asset-space value", "pv-asset-a.b"} {
		if _, err := registry.Apply(Update{AssetRef: asset, Source: provenance(identity, registryRef, '9')}); !errors.Is(err, ErrInvalidFact) {
			t.Errorf("asset %q error = %v", asset, err)
		}
	}
}

func TestRegistryAcceptsFinalProfileVersionDelimiter(t *testing.T) {
	identity := SourceIdentity{Protocol: "test_source", ProfileID: "vendor@family@1.0.0", ProfileVersion: "1.0.0", Validity: SourceTerminalVerified}
	registryRef := Digest("sha256:aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
	registry := NewRegistry(staticResolver{identity: registryRef})
	update := accountedUpdate(Update{
		AssetRef: "pv-asset-final-delimiter",
		Source:   provenance(identity, registryRef, 'a'),
		Facts: []FactInput{
			decimalInput(FactACActivePower, Dimensions{Scope: ScopeTotal}, UnitWatt, "1", 0, PolicyTelemetryFastV1, 0),
		},
	})
	if _, err := registry.Apply(update); err != nil {
		t.Fatalf("schema-valid profile ID rejected: %v", err)
	}
}

func TestRegistryEnforcesEnvelopeCardinalityLimits(t *testing.T) {
	identity := SourceIdentity{Protocol: "test_source", ProfileID: "test.profile@1.0.0", ProfileVersion: "1.0.0", Validity: SourceTerminalVerified}
	registryRef := Digest("sha256:aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
	inputs := make([]FactInput, 257)
	for index := range inputs {
		inputs[index] = decimalInput(FactDCCurrent, Dimensions{InputID: fmt.Sprintf("input-%03d", index)}, UnitAmpere, "1", 0, PolicyTelemetryFastV1, 0)
	}
	registry := NewRegistry(staticResolver{identity: registryRef})
	allowed := accountedUpdate(Update{AssetRef: "pv-asset-limit-256", Source: provenance(identity, registryRef, 'b'), Facts: inputs[:256]})
	if _, err := registry.Apply(allowed); err != nil {
		t.Fatalf("256 facts rejected: %v", err)
	}
	rejected := accountedUpdate(Update{AssetRef: "pv-asset-limit-257", Source: provenance(identity, registryRef, 'c'), Facts: inputs})
	if _, err := registry.Apply(rejected); !errors.Is(err, ErrInvalidProjection) {
		t.Fatalf("257 facts error = %v", err)
	}

	overAccounting := accountedUpdate(Update{
		AssetRef: "pv-asset-limit-513",
		Source:   provenance(identity, registryRef, 'd'),
		Facts: []FactInput{
			decimalInput(FactACActivePower, Dimensions{Scope: ScopeTotal}, UnitWatt, "1", 0, PolicyTelemetryFastV1, 0),
		},
	})
	for index := 1; index < 513; index++ {
		ref := Digest(fmt.Sprintf("sha256:%064x", index+1000))
		output := RequestedOutput{SourceRef: overAccounting.Source.SourceObservationRef, RequestedOutputRef: ref}
		overAccounting.RequestedOutputs = append(overAccounting.RequestedOutputs, output)
		overAccounting.ProjectionReport = append(overAccounting.ProjectionReport, Projection{SourceRef: output.SourceRef, RequestedOutputRef: output.RequestedOutputRef, Outcome: ProjectionWithheld})
	}
	if _, err := registry.Apply(overAccounting); !errors.Is(err, ErrInvalidProjection) {
		t.Fatalf("513 accounting rows error = %v", err)
	}
}

func TestRegistryRejectsOversizedFactsBeforeDuplicateAllocation(t *testing.T) {
	identity := SourceIdentity{Protocol: "test_source", ProfileID: "test.profile@1.0.0", ProfileVersion: "1.0.0", Validity: SourceTerminalVerified}
	registryRef := Digest("sha256:aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
	inputs := make([]FactInput, 257)
	for index := range inputs {
		inputs[index] = decimalInput(FactDCCurrent, Dimensions{InputID: fmt.Sprintf("input-%03d", index)}, UnitAmpere, "1", 0, PolicyTelemetryFastV1, 0)
	}
	update := accountedUpdate(Update{AssetRef: "pv-asset-limit-before-allocation", Source: provenance(identity, registryRef, 'e'), Facts: inputs})
	registry := NewRegistry(staticResolver{identity: registryRef})
	allocations := testing.AllocsPerRun(10, func() {
		if _, err := registry.Apply(update); !errors.Is(err, ErrInvalidProjection) {
			t.Fatalf("oversized facts error = %v", err)
		}
	})
	if allocations != 0 {
		t.Fatalf("oversized facts allocated %.0f times before rejection; want 0", allocations)
	}
}

func TestRegistryCanonicalizesBitfieldSymbolsBeforeSnapshot(t *testing.T) {
	identity := SourceIdentity{Protocol: "test_source", ProfileID: "test.profile@1.0.0", ProfileVersion: "1.0.0", Validity: SourceTerminalVerified}
	registryRef := Digest("sha256:aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
	registry := NewRegistry(staticResolver{identity: registryRef})
	update := accountedUpdate(Update{
		AssetRef: "pv-asset-canonical-symbols",
		Source:   provenance(identity, registryRef, 'f'),
		Facts: []FactInput{{
			Candidate: FactCandidate{ID: FactEventFlags, Dimensions: Dimensions{Scope: ScopeTotal}, Value: BitfieldFactValue("INTERNAL_FAULT", "AC_DISCONNECT"), Unit: UnitOne},
			Quality:   QualityGood,
			Receipt:   0,
			Policy:    PolicyStatusV1,
		}},
	})
	snapshot, err := registry.Apply(update)
	if err != nil {
		t.Fatal(err)
	}
	got := lookupFact(t, snapshot, FactEventFlags, Dimensions{Scope: ScopeTotal}).Value.Symbols
	if len(got) != 2 || got[0] != "AC_DISCONNECT" || got[1] != "INTERNAL_FAULT" {
		t.Fatalf("snapshot symbols = %#v; want canonical order", got)
	}
}

func requiredThreePhaseFacts(t *testing.T) map[FactKey]Fact {
	t.Helper()
	facts := make(map[FactKey]Fact)
	add := func(id FactID, dimensions Dimensions, unit Unit, value FactValue) {
		key := NewFactKey(id, dimensions)
		facts[key] = Fact{ID: id, Dimensions: dimensions, Unit: unit, Value: value, Availability: AvailabilityAvailable}
	}
	add(FactACActivePower, Dimensions{Scope: ScopeTotal}, UnitWatt, DecimalFactValue(MustDecimal("1", 0)))
	add(FactACFrequency, Dimensions{Scope: ScopeTotal}, UnitHertz, DecimalFactValue(MustDecimal("50", 0)))
	for _, phase := range []Phase{PhaseL1, PhaseL2, PhaseL3} {
		add(FactACCurrent, Dimensions{Phase: phase}, UnitAmpere, DecimalFactValue(MustDecimal("1", 0)))
		add(FactACVoltageLineToNeutral, Dimensions{Phase: phase}, UnitVolt, DecimalFactValue(MustDecimal("230", 0)))
	}
	add(FactEnergyActiveExportTotal, Dimensions{Scope: ScopeTotal}, UnitWattHour, DecimalFactValue(MustDecimal("1", 0)))
	add(FactOperatingState, Dimensions{Scope: ScopeTotal}, UnitOne, EnumFactValue(OperatingStateOperating))
	return facts
}

func decimalInput(id FactID, dimensions Dimensions, unit Unit, coefficient string, scale int, policy PolicyID, receipt MonotonicNanos) FactInput {
	return FactInput{Candidate: FactCandidate{ID: id, Dimensions: dimensions, Value: DecimalFactValue(MustDecimal(coefficient, scale)), Unit: unit}, Quality: QualityGood, Receipt: receipt, Policy: policy}
}

func accountedUpdate(update Update) Update {
	update.SourceTimeState = SourceTimeUnavailable
	update.Capability = Capability{ID: CapabilityThreePhaseTelemetryV1, Outcome: CapabilityNotSatisfied}
	for index, fact := range update.Facts {
		requested := RequestedOutput{
			SourceRef:          update.Source.SourceObservationRef,
			RequestedOutputRef: Digest(fmt.Sprintf("sha256:%064x", index+1)),
		}
		dimensions := fact.Candidate.Dimensions
		update.RequestedOutputs = append(update.RequestedOutputs, requested)
		update.ProjectionReport = append(update.ProjectionReport, Projection{
			SourceRef:          requested.SourceRef,
			RequestedOutputRef: requested.RequestedOutputRef,
			FactID:             fact.Candidate.ID,
			Dimensions:         &dimensions,
			Outcome:            ProjectionMapped,
		})
	}
	return update
}

func provenance(identity SourceIdentity, registryRef Digest, marker byte) Provenance {
	hex := make([]byte, 64)
	for index := range hex {
		hex[index] = marker
	}
	return Provenance{
		SourceIdentity:       identity,
		SourceRegistryRef:    registryRef,
		SourceObservationRef: Digest("sha256:" + string(hex)),
		SourceShadowRef:      Digest("sha256:" + string(hex)),
		EvidenceRef:          Digest("sha256:" + string(hex)),
	}
}

func lookupFact(t *testing.T, snapshot Snapshot, id FactID, dimensions Dimensions) Fact {
	t.Helper()
	fact, ok := snapshot.Facts[NewFactKey(id, dimensions)]
	if !ok {
		t.Fatalf("missing fact %s %#v", id, dimensions)
	}
	return fact
}
