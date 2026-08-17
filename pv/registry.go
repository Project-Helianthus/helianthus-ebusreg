package pv

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
	"sync"
)

type projectionKey struct {
	sourceRef Digest
	outputRef Digest
}

type assetState struct {
	generation  uint64
	source      Provenance
	facts       map[FactKey]Fact
	origins     map[Digest]Provenance
	requested   map[projectionKey]RequestedOutput
	projections map[projectionKey]Projection
	sourceTime  SourceTimeState
}

type Registry struct {
	mu       sync.RWMutex
	resolver SourceResolver
	catalog  Catalog
	assets   map[string]*assetState
}

func NewRegistry(resolver SourceResolver) *Registry {
	return &Registry{
		resolver: resolver,
		catalog:  CatalogV1(),
		assets:   make(map[string]*assetState),
	}
}

func (registry *Registry) Apply(update Update) (Snapshot, error) {
	if registry == nil || !validAssetRef(update.AssetRef) {
		return Snapshot{}, ErrInvalidFact
	}
	if err := registry.validateProvenance(update.Source); err != nil {
		return Snapshot{}, err
	}
	if duplicateFactInputs(update.Facts) {
		return Snapshot{}, ErrInvalidFact
	}
	if err := validateUpdateAccounting(update); err != nil {
		return Snapshot{}, err
	}
	if update.Evaluated < 0 {
		return Snapshot{}, ErrInvalidMonotonicTime
	}

	registry.mu.Lock()
	defer registry.mu.Unlock()

	state := cloneAssetState(registry.assets[update.AssetRef])
	if state == nil {
		state = &assetState{facts: make(map[FactKey]Fact), origins: make(map[Digest]Provenance), requested: make(map[projectionKey]RequestedOutput), projections: make(map[projectionKey]Projection)}
	}
	if err := evaluateRetainedFacts(state.facts, update.Evaluated); err != nil {
		return Snapshot{}, err
	}
	state.source = update.Source
	state.sourceTime = update.SourceTimeState
	if existing, ok := state.origins[update.Source.SourceObservationRef]; ok && existing != update.Source {
		return Snapshot{}, ErrSourceNotAdmitted
	}
	state.origins[update.Source.SourceObservationRef] = update.Source

	seen := make(map[FactKey]struct{}, len(update.Facts))
	for _, input := range update.Facts {
		definition, ok := registry.catalog.Facts[input.Candidate.ID]
		if !ok || registry.catalog.ValidateCandidate(input.Candidate) != nil || definition.validatePolicy(input.Policy) != nil {
			return Snapshot{}, ErrInvalidFact
		}
		if input.Quality != QualityGood && input.Quality != QualitySuspect && input.Quality != QualityBad {
			return Snapshot{}, ErrInvalidFact
		}
		policy, ok := PolicyByID(input.Policy)
		if !ok {
			return Snapshot{}, ErrInvalidFact
		}
		temporal, err := EvaluateTemporal(policy, input.Receipt, update.Evaluated)
		if err != nil {
			return Snapshot{}, err
		}

		key := NewFactKey(input.Candidate.ID, input.Candidate.Dimensions)
		if _, duplicate := seen[key]; duplicate {
			return Snapshot{}, ErrInvalidFact
		}
		seen[key] = struct{}{}
		var continuity *Continuity
		if definition.Accumulator {
			var previous *Decimal
			prior, hasPrior := state.facts[key]
			if hasPrior && prior.Value.Decimal != nil {
				copy := *prior.Value.Decimal
				previous = &copy
			}
			result := Continuity{State: ContinuityDiscontinuity}
			if !hasPrior || prior.Availability == AvailabilityAvailable && sameSourceIdentity(state.origins[prior.OriginRef], update.Source) {
				var err error
				result, err = EvaluateCounter(previous, *input.Candidate.Value.Decimal, input.Counter)
				if err != nil {
					return Snapshot{}, err
				}
			}
			continuity = &result
		} else if !input.Counter.empty() {
			return Snapshot{}, ErrInvalidCounter
		}

		state.facts[key] = Fact{
			ID:           input.Candidate.ID,
			Dimensions:   input.Candidate.Dimensions,
			Value:        cloneFactValue(input.Candidate.Value),
			Unit:         input.Candidate.Unit,
			Quality:      input.Quality,
			Availability: temporal.Availability,
			Freshness:    temporal.Freshness,
			Temporal:     temporal.Temporal,
			OriginRef:    update.Source.SourceObservationRef,
			Continuity:   continuity,
		}
	}

	if err := evaluateRetainedFacts(state.facts, update.Evaluated); err != nil {
		return Snapshot{}, err
	}
	capabilitySatisfied := registry.catalog.ThreePhaseTelemetrySatisfied(state.facts)
	wantOutcome := CapabilityNotSatisfied
	if capabilitySatisfied {
		wantOutcome = CapabilitySatisfied
	}
	if update.Capability.ID != CapabilityThreePhaseTelemetryV1 || update.Capability.Outcome != wantOutcome {
		return Snapshot{}, ErrInvalidCapability
	}
	mergeAccounting(state, update)
	state.generation++
	pruneOrigins(state)
	registry.assets[update.AssetRef] = state
	return registry.snapshotLocked(update.AssetRef, state, update.Evaluated), nil
}

func duplicateFactInputs(facts []FactInput) bool {
	seen := make(map[FactKey]struct{}, len(facts))
	for _, fact := range facts {
		key := NewFactKey(fact.Candidate.ID, fact.Candidate.Dimensions)
		if _, duplicate := seen[key]; duplicate {
			return true
		}
		seen[key] = struct{}{}
	}
	return false
}

func (registry *Registry) Snapshot(assetRef string, evaluated MonotonicNanos) (Snapshot, error) {
	if registry == nil || evaluated < 0 {
		return Snapshot{}, ErrInvalidMonotonicTime
	}
	registry.mu.RLock()
	state := cloneAssetState(registry.assets[assetRef])
	registry.mu.RUnlock()
	if state == nil {
		return Snapshot{}, ErrAssetNotFound
	}
	if err := evaluateRetainedFacts(state.facts, evaluated); err != nil {
		return Snapshot{}, err
	}
	pruneOrigins(state)
	return registry.snapshotLocked(assetRef, state, evaluated), nil
}

func (registry *Registry) snapshotLocked(assetRef string, state *assetState, evaluated MonotonicNanos) Snapshot {
	facts := cloneFacts(state.facts)
	origins := make(map[Digest]Provenance, len(state.origins))
	for reference, origin := range state.origins {
		origins[reference] = origin
	}
	return Snapshot{
		ContractID:       ContractV1,
		AssetRef:         assetRef,
		Generation:       state.generation,
		Evaluated:        evaluated,
		Source:           state.source,
		Facts:            facts,
		Origins:          origins,
		SourceTimeState:  state.sourceTime,
		Capability:       Capability{ID: CapabilityThreePhaseTelemetryV1, Outcome: capabilityOutcome(registry.catalog.ThreePhaseTelemetrySatisfied(facts))},
		RequestedOutputs: sortedRequested(state.requested),
		ProjectionReport: sortedProjections(state.projections),
	}
}

func (registry *Registry) validateProvenance(provenance Provenance) error {
	if provenance.Validity != SourceTerminalVerified || provenance.Protocol == "" || provenance.ProfileID == "" || provenance.ProfileVersion == "" {
		return ErrSourceNotAdmitted
	}
	if len(provenance.Protocol) > 128 || len(provenance.ProfileID) > 128 || len(provenance.ProfileVersion) > 64 ||
		!sourceTokenPattern.MatchString(provenance.Protocol) || !sourceTokenPattern.MatchString(provenance.ProfileID) || !versionTokenPattern.MatchString(provenance.ProfileVersion) {
		return ErrSourceNotAdmitted
	}
	parts := strings.Split(provenance.ProfileID, "@")
	if len(parts) != 2 || parts[0] == "" || parts[1] != provenance.ProfileVersion {
		return ErrSourceNotAdmitted
	}
	for _, digest := range []Digest{
		provenance.SourceRegistryRef,
		provenance.SourceObservationRef,
		provenance.SourceShadowRef,
		provenance.EvidenceRef,
	} {
		if digest.Validate() != nil {
			return ErrSourceNotAdmitted
		}
	}
	if registry.resolver == nil {
		return ErrSourceNotAdmitted
	}
	expected, ok := registry.resolver.ResolveSource(provenance.SourceIdentity)
	if !ok || expected != provenance.SourceRegistryRef {
		return ErrSourceNotAdmitted
	}
	return nil
}

func evaluateRetainedFacts(facts map[FactKey]Fact, evaluated MonotonicNanos) error {
	for key, fact := range facts {
		policy, ok := PolicyByID(fact.Temporal.Policy)
		if !ok {
			return ErrInvalidFact
		}
		state, err := EvaluateTemporal(policy, fact.Temporal.Receipt, evaluated)
		if err != nil {
			return err
		}
		fact.Availability = state.Availability
		fact.Freshness = state.Freshness
		facts[key] = fact
	}
	return nil
}

func pruneOrigins(state *assetState) {
	referenced := map[Digest]struct{}{state.source.SourceObservationRef: {}}
	for _, fact := range state.facts {
		referenced[fact.OriginRef] = struct{}{}
	}
	for reference := range state.origins {
		if _, ok := referenced[reference]; !ok {
			delete(state.origins, reference)
		}
	}
}

func cloneAssetState(source *assetState) *assetState {
	if source == nil {
		return nil
	}
	origins := make(map[Digest]Provenance, len(source.origins))
	for reference, origin := range source.origins {
		origins[reference] = origin
	}
	requested := make(map[projectionKey]RequestedOutput, len(source.requested))
	for key, value := range source.requested {
		requested[key] = value
	}
	projections := make(map[projectionKey]Projection, len(source.projections))
	for key, value := range source.projections {
		value.Dimensions = cloneDimensions(value.Dimensions)
		projections[key] = value
	}
	return &assetState{
		generation:  source.generation,
		source:      source.source,
		facts:       cloneFacts(source.facts),
		origins:     origins,
		requested:   requested,
		projections: projections,
		sourceTime:  source.sourceTime,
	}
}

func cloneFacts(source map[FactKey]Fact) map[FactKey]Fact {
	facts := make(map[FactKey]Fact, len(source))
	for key, fact := range source {
		fact.Value = cloneFactValue(fact.Value)
		if fact.Continuity != nil {
			copy := *fact.Continuity
			if copy.Delta != nil {
				delta := *copy.Delta
				copy.Delta = &delta
			}
			if copy.Modulus != nil {
				modulus := *copy.Modulus
				copy.Modulus = &modulus
			}
			fact.Continuity = &copy
		}
		facts[key] = fact
	}
	return facts
}

func cloneFactValue(value FactValue) FactValue {
	if value.Decimal != nil {
		decimal := *value.Decimal
		value.Decimal = &decimal
	}
	value.Symbols = append([]string(nil), value.Symbols...)
	return value
}

func cloneDimensions(value *Dimensions) *Dimensions {
	if value == nil {
		return nil
	}
	copy := *value
	return &copy
}

func sameSourceIdentity(left, right Provenance) bool {
	return left.SourceIdentity == right.SourceIdentity && left.SourceRegistryRef == right.SourceRegistryRef
}

func validateUpdateAccounting(update Update) error {
	if update.SourceTimeState != SourceTimeUnavailable && update.SourceTimeState != SourceTimeValid && update.SourceTimeState != SourceTimeInvalid {
		return ErrInvalidProjection
	}
	if len(update.Facts) == 0 || len(update.RequestedOutputs) == 0 || len(update.ProjectionReport) != len(update.RequestedOutputs) {
		return ErrInvalidProjection
	}
	requested := make(map[projectionKey]struct{}, len(update.RequestedOutputs))
	for _, output := range update.RequestedOutputs {
		key := projectionKey{output.SourceRef, output.RequestedOutputRef}
		if output.SourceRef.Validate() != nil || output.RequestedOutputRef.Validate() != nil {
			return ErrInvalidProjection
		}
		if _, duplicate := requested[key]; duplicate {
			return ErrInvalidProjection
		}
		requested[key] = struct{}{}
	}
	mappedFacts := make(map[FactKey]int, len(update.Facts))
	seenProjection := make(map[projectionKey]struct{}, len(update.ProjectionReport))
	for _, projection := range update.ProjectionReport {
		key := projectionKey{projection.SourceRef, projection.RequestedOutputRef}
		if _, ok := requested[key]; !ok || projection.SourceRef != update.Source.SourceObservationRef {
			return ErrInvalidProjection
		}
		if _, duplicate := seenProjection[key]; duplicate {
			return ErrInvalidProjection
		}
		seenProjection[key] = struct{}{}
		switch projection.Outcome {
		case ProjectionMapped:
			if projection.FactID == "" || projection.Dimensions == nil {
				return ErrInvalidProjection
			}
			mappedFacts[NewFactKey(projection.FactID, *projection.Dimensions)]++
		case ProjectionWithheld, ProjectionUnrepresentable:
			if projection.FactID != "" || projection.Dimensions != nil {
				return ErrInvalidProjection
			}
		default:
			return ErrInvalidProjection
		}
	}
	for _, fact := range update.Facts {
		if mappedFacts[NewFactKey(fact.Candidate.ID, fact.Candidate.Dimensions)] != 1 {
			return ErrInvalidProjection
		}
	}
	if len(mappedFacts) != len(update.Facts) {
		return ErrInvalidProjection
	}
	return nil
}

func mergeAccounting(state *assetState, update Update) {
	updatedFacts := make(map[FactKey]struct{}, len(update.Facts))
	for _, fact := range update.Facts {
		updatedFacts[NewFactKey(fact.Candidate.ID, fact.Candidate.Dimensions)] = struct{}{}
	}
	for key, projection := range state.projections {
		remove := projection.Outcome != ProjectionMapped
		if projection.Dimensions != nil {
			_, remove = updatedFacts[NewFactKey(projection.FactID, *projection.Dimensions)]
		}
		if remove {
			delete(state.projections, key)
			delete(state.requested, key)
		}
	}
	for _, output := range update.RequestedOutputs {
		key := projectionKey{output.SourceRef, output.RequestedOutputRef}
		state.requested[key] = output
	}
	for _, projection := range update.ProjectionReport {
		key := projectionKey{projection.SourceRef, projection.RequestedOutputRef}
		projection.Dimensions = cloneDimensions(projection.Dimensions)
		state.projections[key] = projection
	}
}

func capabilityOutcome(satisfied bool) CapabilityOutcome {
	if satisfied {
		return CapabilitySatisfied
	}
	return CapabilityNotSatisfied
}

func sortedRequested(values map[projectionKey]RequestedOutput) []RequestedOutput {
	result := make([]RequestedOutput, 0, len(values))
	for _, value := range values {
		result = append(result, value)
	}
	sort.Slice(result, func(i, j int) bool {
		if result[i].SourceRef == result[j].SourceRef {
			return result[i].RequestedOutputRef < result[j].RequestedOutputRef
		}
		return result[i].SourceRef < result[j].SourceRef
	})
	return result
}

func sortedProjections(values map[projectionKey]Projection) []Projection {
	result := make([]Projection, 0, len(values))
	for _, value := range values {
		value.Dimensions = cloneDimensions(value.Dimensions)
		result = append(result, value)
	}
	sort.Slice(result, func(i, j int) bool {
		if result[i].SourceRef == result[j].SourceRef {
			return result[i].RequestedOutputRef < result[j].RequestedOutputRef
		}
		return result[i].SourceRef < result[j].SourceRef
	})
	return result
}

var assetRefPattern = regexp.MustCompile(`^pv-asset-[A-Za-z0-9_-]{1,96}$`)
var sourceTokenPattern = regexp.MustCompile(`^[A-Za-z][A-Za-z0-9._@-]*$`)
var versionTokenPattern = regexp.MustCompile(`^[0-9][0-9A-Za-z._-]*$`)

func validAssetRef(value string) bool { return assetRefPattern.MatchString(value) }

func (snapshot Snapshot) String() string {
	return fmt.Sprintf("%s asset=%s generation=%d facts=%d", snapshot.ContractID, snapshot.AssetRef, snapshot.Generation, len(snapshot.Facts))
}
