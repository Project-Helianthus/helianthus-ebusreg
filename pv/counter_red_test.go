package pv

import (
	"testing"
)

func TestCounterContinuityLifecycle(t *testing.T) {
	current := MustDecimal("10000", -2)
	baseline, err := EvaluateCounter(nil, current, CounterEvidence{})
	if err != nil || baseline.State != ContinuityBaseline || baseline.Delta != nil {
		t.Fatalf("baseline = %#v, err=%v", baseline, err)
	}

	previous := MustDecimal("9999", -2)
	contiguous, err := EvaluateCounter(&previous, current, CounterEvidence{})
	if err != nil || contiguous.State != ContinuityContiguous {
		t.Fatalf("contiguous = %#v, err=%v", contiguous, err)
	}
	if contiguous.Delta == nil || contiguous.Delta.Coefficient != "1" || contiguous.Delta.Scale != -2 {
		t.Fatalf("delta = %#v, want 1e-2", contiguous.Delta)
	}

	decreased := MustDecimal("10", 0)
	prior := MustDecimal("100", 0)
	discontinuous, err := EvaluateCounter(&prior, decreased, CounterEvidence{})
	if err != nil || discontinuous.State != ContinuityDiscontinuity || discontinuous.Delta != nil {
		t.Fatalf("discontinuity = %#v, err=%v", discontinuous, err)
	}
}

func TestCounterResetAndRolloverRequireExplicitEvidence(t *testing.T) {
	prior := MustDecimal("999", 0)
	current := MustDecimal("2", 0)

	reset, err := EvaluateCounter(&prior, current, CounterEvidence{
		Kind:        CounterEventReset,
		EvidenceRef: Digest("sha256:aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"),
	})
	if err != nil || reset.State != ContinuityReset || reset.Delta != nil {
		t.Fatalf("reset = %#v, err=%v", reset, err)
	}

	rollover, err := EvaluateCounter(&prior, current, CounterEvidence{
		Kind:             CounterEventRollover,
		Modulus:          ptrDecimal(MustDecimal("1000", 0)),
		BoundaryVerified: true,
		EvidenceRef:      Digest("sha256:bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"),
	})
	if err != nil || rollover.State != ContinuityRollover {
		t.Fatalf("rollover = %#v, err=%v", rollover, err)
	}
	if rollover.Delta == nil || rollover.Delta.Coefficient != "3" || rollover.Delta.Scale != 0 {
		t.Fatalf("rollover delta = %#v, want 3", rollover.Delta)
	}

	for _, evidence := range []CounterEvidence{
		{Kind: CounterEventReset},
		{Kind: CounterEventRollover, Modulus: ptrDecimal(MustDecimal("1000", 0)), EvidenceRef: Digest("sha256:bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb")},
		{Kind: CounterEventRollover, Modulus: ptrDecimal(MustDecimal("0", 0)), BoundaryVerified: true, EvidenceRef: Digest("sha256:bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb")},
	} {
		result, err := EvaluateCounter(&prior, current, evidence)
		if err == nil || result.State != ContinuityDiscontinuity || result.Delta != nil {
			t.Errorf("evidence %#v result=%#v err=%v, want discontinuity error", evidence, result, err)
		}
	}
}

func ptrDecimal(value Decimal) *Decimal { return &value }
