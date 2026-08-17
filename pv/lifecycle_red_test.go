package pv

import (
	"errors"
	"math"
	"testing"
)

func TestFreshnessPoliciesAndBoundaryTransitions(t *testing.T) {
	cases := []struct {
		policy PolicyID
		fresh  MonotonicNanos
		retain MonotonicNanos
	}{
		{PolicyTelemetryFastV1, 30_000_000_000, 300_000_000_000},
		{PolicyStatusV1, 60_000_000_000, 600_000_000_000},
		{PolicyAccumulatorV1, 900_000_000_000, 86_400_000_000_000},
		{PolicyRatingV1, 86_400_000_000_000, 2_592_000_000_000_000},
	}
	for _, tc := range cases {
		policy, ok := PolicyByID(tc.policy)
		if !ok {
			t.Fatalf("missing policy %q", tc.policy)
		}
		if policy.FreshFor != tc.fresh || policy.RetainFor != tc.retain {
			t.Errorf("policy %q = %#v", tc.policy, policy)
		}
	}

	receipt := MonotonicNanos(1_000)
	policy, _ := PolicyByID(PolicyTelemetryFastV1)
	checks := []struct {
		evaluated    MonotonicNanos
		availability Availability
		freshness    Freshness
	}{
		{receipt, AvailabilityAvailable, FreshnessFresh},
		{receipt + policy.FreshFor - 1, AvailabilityAvailable, FreshnessFresh},
		{receipt + policy.FreshFor, AvailabilityAvailable, FreshnessStale},
		{receipt + policy.RetainFor - 1, AvailabilityAvailable, FreshnessStale},
		{receipt + policy.RetainFor, AvailabilityUnavailable, FreshnessExpired},
	}
	for _, check := range checks {
		state, err := EvaluateTemporal(policy, receipt, check.evaluated)
		if err != nil {
			t.Fatalf("EvaluateTemporal: %v", err)
		}
		if state.Availability != check.availability || state.Freshness != check.freshness {
			t.Errorf("at %d got %s/%s, want %s/%s", check.evaluated, state.Availability, state.Freshness, check.availability, check.freshness)
		}
	}
}

func TestFreshnessRejectsDifferentClockDomain(t *testing.T) {
	policy, _ := PolicyByID(PolicyTelemetryFastV1)
	if _, err := EvaluateTemporal(policy, 100, 99); !errors.Is(err, ErrInvalidMonotonicTime) {
		t.Fatalf("error = %v, want ErrInvalidMonotonicTime", err)
	}
}

func TestFreshnessRejectsDeadlineOverflow(t *testing.T) {
	policy, _ := PolicyByID(PolicyTelemetryFastV1)
	receipt := MonotonicNanos(math.MaxInt64) - policy.RetainFor + 1
	if _, err := EvaluateTemporal(policy, receipt, receipt); !errors.Is(err, ErrInvalidMonotonicTime) {
		t.Fatalf("overflow error = %v, want ErrInvalidMonotonicTime", err)
	}
}
