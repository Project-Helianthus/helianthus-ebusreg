package system

import b524UnscopedHintSort "sort"

type B524UnscopedHintB524HintKey struct {
	B524UnscopedHintOpcode   byte   `json:"opcode"`
	B524UnscopedHintGroup    byte   `json:"group"`
	B524UnscopedHintRegister uint16 `json:"register"`
}

type B524UnscopedHintB524HintStatus string

const (
	B524UnscopedHintAbsent B524UnscopedHintB524HintStatus = "absent"
	B524HintUnscoped       B524UnscopedHintB524HintStatus = "unscoped_catalog"
)

type B524UnscopedHintGrade string

const B524UnscopedHintApplicabilityUnknown B524UnscopedHintGrade = "Unknown"

type B524UnscopedHintProperties struct {
	B524UnscopedHintType string  `json:"type"`
	B524UnscopedHintMin  float64 `json:"min"`
	B524UnscopedHintMax  float64 `json:"max"`
	B524UnscopedHintStep float64 `json:"step"`
}

type B524UnscopedHintResult struct {
	B524UnscopedHintKey           B524UnscopedHintB524HintKey    `json:"key"`
	B524UnscopedHintStatus        B524UnscopedHintB524HintStatus `json:"status"`
	B524UnscopedHintQualified     bool                           `json:"qualified"`
	B524UnscopedHintApplicability B524UnscopedHintGrade          `json:"applicability"`
	B524UnscopedHint              B524UnscopedHintProperties     `json:"hint"`
}

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
