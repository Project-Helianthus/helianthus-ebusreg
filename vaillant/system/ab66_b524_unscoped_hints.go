package system

import ab66sort "sort"

type Ab66B524HintKey struct {
	Ab66Opcode   byte
	Ab66Group    byte
	Ab66Register uint16
}

type Ab66B524HintStatus string

const (
	Ab66B524HintAbsent   Ab66B524HintStatus = "absent"
	Ab66B524HintUnscoped Ab66B524HintStatus = "unscoped_catalog"
)

type Ab66B524HintGrade string

const Ab66B524ApplicabilityUnknown Ab66B524HintGrade = "Unknown"

type Ab66B524HintProperties struct {
	Ab66Type string
	Ab66Min  float64
	Ab66Max  float64
	Ab66Step float64
}

type Ab66B524HintResult struct {
	Ab66Key           Ab66B524HintKey
	Ab66Status        Ab66B524HintStatus
	Ab66Qualified     bool
	Ab66Applicability Ab66B524HintGrade
	Ab66Hint          Ab66B524HintProperties
}

func Ab66LookupB524UnscopedHint(ab66Key Ab66B524HintKey) Ab66B524HintResult {
	ab66Result := Ab66B524HintResult{
		Ab66Key:           ab66Key,
		Ab66Status:        Ab66B524HintAbsent,
		Ab66Applicability: Ab66B524ApplicabilityUnknown,
	}
	ab66Constraint, ab66OK := lookupB524Constraint(ab66Key.Ab66Opcode, ab66Key.Ab66Group, ab66Key.Ab66Register)
	if !ab66OK {
		return ab66Result
	}
	ab66Result.Ab66Status = Ab66B524HintUnscoped
	ab66Result.Ab66Hint = Ab66B524HintProperties{
		Ab66Type: ab66Constraint.Type,
		Ab66Min:  ab66Constraint.Min,
		Ab66Max:  ab66Constraint.Max,
		Ab66Step: ab66Constraint.Step,
	}
	return ab66Result
}

func Ab66B524UnscopedHints() []Ab66B524HintResult {
	ab66Results := make([]Ab66B524HintResult, 0, len(b524StaticConstraintsCatalog))
	for ab66Selector := range b524StaticConstraintsCatalog {
		ab66Results = append(ab66Results, Ab66LookupB524UnscopedHint(Ab66B524HintKey{
			Ab66Opcode:   ab66Selector.Opcode,
			Ab66Group:    ab66Selector.Group,
			Ab66Register: ab66Selector.Register,
		}))
	}
	ab66sort.Slice(ab66Results, func(ab66I, ab66J int) bool {
		ab66Left, ab66Right := ab66Results[ab66I].Ab66Key, ab66Results[ab66J].Ab66Key
		if ab66Left.Ab66Opcode != ab66Right.Ab66Opcode {
			return ab66Left.Ab66Opcode < ab66Right.Ab66Opcode
		}
		if ab66Left.Ab66Group != ab66Right.Ab66Group {
			return ab66Left.Ab66Group < ab66Right.Ab66Group
		}
		return ab66Left.Ab66Register < ab66Right.Ab66Register
	})
	return ab66Results
}
