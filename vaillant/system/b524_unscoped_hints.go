package system

import b524UnscopedHintSort "sort"

// B524UnscopedHintB524HintKey selects an unscoped hint by opcode, group and
// register; it carries no controller identity and no instance.
type B524UnscopedHintB524HintKey struct {
	B524UnscopedHintOpcode   byte   `json:"opcode"`
	B524UnscopedHintGroup    byte   `json:"group"`
	B524UnscopedHintRegister uint16 `json:"register"`
}

// B524UnscopedHintB524HintStatus says whether the unscoped catalog has a hint
// for a key.
type B524UnscopedHintB524HintStatus string

// Unscoped hint statuses: absent when the catalog has no entry, unscoped
// catalog when it has one that is not bound to a controller identity.
const (
	B524UnscopedHintAbsent B524UnscopedHintB524HintStatus = "absent"
	B524HintUnscoped       B524UnscopedHintB524HintStatus = "unscoped_catalog"
)

// B524UnscopedHintGrade is the applicability grade of an unscoped hint.
type B524UnscopedHintGrade string

// B524UnscopedHintApplicabilityUnknown is the only applicability an unscoped
// hint has: it is not known to apply to any particular controller.
const B524UnscopedHintApplicabilityUnknown B524UnscopedHintGrade = "Unknown"

// B524UnscopedHintProperties is the type name, minimum, maximum and step an
// unscoped catalog entry records.
type B524UnscopedHintProperties struct {
	B524UnscopedHintType string  `json:"type"`
	B524UnscopedHintMin  float64 `json:"min"`
	B524UnscopedHintMax  float64 `json:"max"`
	B524UnscopedHintStep float64 `json:"step"`
}

// B524UnscopedHintResult is the outcome of an unscoped hint lookup.
// Qualified is always false and applicability Unknown: a hint is never a
// qualified bound. Use LookupB524Bound for bounds scoped to a controller
// identity.
type B524UnscopedHintResult struct {
	B524UnscopedHintKey           B524UnscopedHintB524HintKey    `json:"key"`
	B524UnscopedHintStatus        B524UnscopedHintB524HintStatus `json:"status"`
	B524UnscopedHintQualified     bool                           `json:"qualified"`
	B524UnscopedHintApplicability B524UnscopedHintGrade          `json:"applicability"`
	B524UnscopedHint              B524UnscopedHintProperties     `json:"hint"`
}

// B524LookupUnscopedHint returns the unscoped catalog hint for a key, or an
// absent result with zero properties when the catalog has none.
func B524LookupUnscopedHint(b524UnscopedHintKey B524UnscopedHintB524HintKey) B524UnscopedHintResult {
	b524UnscopedHintResult := B524UnscopedHintResult{
		B524UnscopedHintKey:           b524UnscopedHintKey,
		B524UnscopedHintStatus:        B524UnscopedHintAbsent,
		B524UnscopedHintApplicability: B524UnscopedHintApplicabilityUnknown,
	}
	b524UnscopedHintConstraint, b524UnscopedHintOK := lookupB524Constraint(b524UnscopedHintKey.B524UnscopedHintOpcode, b524UnscopedHintKey.B524UnscopedHintGroup, b524UnscopedHintKey.B524UnscopedHintRegister)
	if !b524UnscopedHintOK {
		return b524UnscopedHintResult
	}
	b524UnscopedHintResult.B524UnscopedHintStatus = B524HintUnscoped
	b524UnscopedHintResult.B524UnscopedHint = B524UnscopedHintProperties{
		B524UnscopedHintType: b524UnscopedHintConstraint.Type,
		B524UnscopedHintMin:  b524UnscopedHintConstraint.Min,
		B524UnscopedHintMax:  b524UnscopedHintConstraint.Max,
		B524UnscopedHintStep: b524UnscopedHintConstraint.Step,
	}
	return b524UnscopedHintResult
}

// B524UnscopedHints returns every unscoped catalog hint, sorted by opcode,
// group and register, in a new slice.
func B524UnscopedHints() []B524UnscopedHintResult {
	b524UnscopedHintResults := make([]B524UnscopedHintResult, 0, len(b524StaticConstraintsCatalog))
	for b524UnscopedHintSelector := range b524StaticConstraintsCatalog {
		b524UnscopedHintResults = append(b524UnscopedHintResults, B524LookupUnscopedHint(B524UnscopedHintB524HintKey{
			B524UnscopedHintOpcode:   b524UnscopedHintSelector.Opcode,
			B524UnscopedHintGroup:    b524UnscopedHintSelector.Group,
			B524UnscopedHintRegister: b524UnscopedHintSelector.Register,
		}))
	}
	b524UnscopedHintSort.Slice(b524UnscopedHintResults, func(b524UnscopedHintI, b524UnscopedHintJ int) bool {
		b524UnscopedHintLeft, b524UnscopedHintRight := b524UnscopedHintResults[b524UnscopedHintI].B524UnscopedHintKey, b524UnscopedHintResults[b524UnscopedHintJ].B524UnscopedHintKey
		if b524UnscopedHintLeft.B524UnscopedHintOpcode != b524UnscopedHintRight.B524UnscopedHintOpcode {
			return b524UnscopedHintLeft.B524UnscopedHintOpcode < b524UnscopedHintRight.B524UnscopedHintOpcode
		}
		if b524UnscopedHintLeft.B524UnscopedHintGroup != b524UnscopedHintRight.B524UnscopedHintGroup {
			return b524UnscopedHintLeft.B524UnscopedHintGroup < b524UnscopedHintRight.B524UnscopedHintGroup
		}
		return b524UnscopedHintLeft.B524UnscopedHintRegister < b524UnscopedHintRight.B524UnscopedHintRegister
	})
	return b524UnscopedHintResults
}
