package pv

import "math"

const (
	PolicyTelemetryFastV1 PolicyID = "pv.telemetry.fast.v1"
	PolicyStatusV1        PolicyID = "pv.status.v1"
	PolicyAccumulatorV1   PolicyID = "pv.accumulator.v1"
	PolicyRatingV1        PolicyID = "pv.rating.v1"
)

type FreshnessPolicy struct {
	ID        PolicyID
	FreshFor  MonotonicNanos
	RetainFor MonotonicNanos
}

var policiesV1 = map[PolicyID]FreshnessPolicy{
	PolicyTelemetryFastV1: {ID: PolicyTelemetryFastV1, FreshFor: 30_000_000_000, RetainFor: 300_000_000_000},
	PolicyStatusV1:        {ID: PolicyStatusV1, FreshFor: 60_000_000_000, RetainFor: 600_000_000_000},
	PolicyAccumulatorV1:   {ID: PolicyAccumulatorV1, FreshFor: 900_000_000_000, RetainFor: 86_400_000_000_000},
	PolicyRatingV1:        {ID: PolicyRatingV1, FreshFor: 86_400_000_000_000, RetainFor: 2_592_000_000_000_000},
}

func PolicyByID(id PolicyID) (FreshnessPolicy, bool) {
	policy, ok := policiesV1[id]
	return policy, ok
}

func EvaluateTemporal(policy FreshnessPolicy, receipt, evaluated MonotonicNanos) (TemporalState, error) {
	if receipt < 0 || evaluated < receipt || policy.FreshFor <= 0 || policy.RetainFor <= policy.FreshFor {
		return TemporalState{}, ErrInvalidMonotonicTime
	}
	if receipt > MonotonicNanos(math.MaxInt64)-policy.RetainFor {
		return TemporalState{}, ErrInvalidMonotonicTime
	}
	temporal := Temporal{
		Receipt:     receipt,
		FreshUntil:  receipt + policy.FreshFor,
		RetainUntil: receipt + policy.RetainFor,
		Policy:      policy.ID,
	}
	state := TemporalState{Availability: AvailabilityAvailable, Freshness: FreshnessFresh, Temporal: temporal}
	if evaluated >= temporal.RetainUntil {
		state.Availability = AvailabilityUnavailable
		state.Freshness = FreshnessExpired
	} else if evaluated >= temporal.FreshUntil {
		state.Freshness = FreshnessStale
	}
	return state, nil
}
