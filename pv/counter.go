package pv

import (
	"fmt"
	"math/big"
)

func EvaluateCounter(previous *Decimal, current Decimal, evidence CounterEvidence) (Continuity, error) {
	if current.Validate() != nil {
		return Continuity{State: ContinuityDiscontinuity}, ErrInvalidCounter
	}
	if previous == nil {
		if !evidence.empty() {
			return Continuity{State: ContinuityDiscontinuity}, ErrInvalidCounter
		}
		return Continuity{State: ContinuityBaseline}, nil
	}
	if previous.Validate() != nil {
		return Continuity{State: ContinuityDiscontinuity}, ErrInvalidCounter
	}
	comparison := compareDecimal(current, *previous)
	if comparison >= 0 && evidence.empty() {
		delta := subtractDecimal(current, *previous)
		return Continuity{State: ContinuityContiguous, Delta: &delta}, nil
	}

	switch evidence.Kind {
	case CounterEventNone:
		if !evidence.empty() {
			return Continuity{State: ContinuityDiscontinuity}, ErrInvalidCounter
		}
		return Continuity{State: ContinuityDiscontinuity}, nil
	case CounterEventReset:
		if evidence.Modulus != nil || evidence.BoundaryVerified || evidence.EvidenceRef.Validate() != nil {
			return Continuity{State: ContinuityDiscontinuity}, ErrInvalidCounter
		}
		return Continuity{State: ContinuityReset, EvidenceRef: evidence.EvidenceRef}, nil
	case CounterEventRollover:
		if comparison >= 0 || evidence.Modulus == nil || !evidence.BoundaryVerified || evidence.EvidenceRef.Validate() != nil {
			return Continuity{State: ContinuityDiscontinuity}, ErrInvalidCounter
		}
		if evidence.Modulus.Validate() != nil || decimalSign(*evidence.Modulus) <= 0 || compareDecimal(*previous, *evidence.Modulus) >= 0 {
			return Continuity{State: ContinuityDiscontinuity}, ErrInvalidCounter
		}
		remaining := subtractDecimal(*evidence.Modulus, *previous)
		delta := addDecimal(remaining, current)
		modulus := *evidence.Modulus
		return Continuity{State: ContinuityRollover, Delta: &delta, Modulus: &modulus, EvidenceRef: evidence.EvidenceRef}, nil
	default:
		return Continuity{State: ContinuityDiscontinuity}, ErrInvalidCounter
	}
}

func compareDecimal(left, right Decimal) int {
	scale := left.Scale
	if right.Scale < scale {
		scale = right.Scale
	}
	return scaledCoefficient(left, scale).Cmp(scaledCoefficient(right, scale))
}

func subtractDecimal(left, right Decimal) Decimal {
	scale := left.Scale
	if right.Scale < scale {
		scale = right.Scale
	}
	coefficient := new(big.Int).Sub(scaledCoefficient(left, scale), scaledCoefficient(right, scale))
	return Decimal{Coefficient: coefficient.String(), Scale: scale}
}

func addDecimal(left, right Decimal) Decimal {
	scale := left.Scale
	if right.Scale < scale {
		scale = right.Scale
	}
	coefficient := new(big.Int).Add(scaledCoefficient(left, scale), scaledCoefficient(right, scale))
	return Decimal{Coefficient: coefficient.String(), Scale: scale}
}

func scaledCoefficient(decimal Decimal, targetScale int) *big.Int {
	coefficient, ok := new(big.Int).SetString(decimal.Coefficient, 10)
	if !ok {
		panic(fmt.Sprintf("validated decimal %q could not be parsed", decimal.Coefficient))
	}
	exponent := decimal.Scale - targetScale
	if exponent == 0 {
		return coefficient
	}
	factor := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(exponent)), nil)
	return coefficient.Mul(coefficient, factor)
}

func decimalSign(decimal Decimal) int {
	coefficient, ok := new(big.Int).SetString(decimal.Coefficient, 10)
	if !ok {
		return 0
	}
	return coefficient.Sign()
}
