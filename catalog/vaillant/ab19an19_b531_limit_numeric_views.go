package vaillant

const (
	Ab19An19Family    = "vaillant-b531-family"
	Ab19An19Variant   = "b531-op02-candidate"
	Ab19An19Direction = "positive"
)

type Ab19An19Field struct {
	Ab19An19Raw       []byte
	Ab19An19State     string
	Ab19An19Available bool
	Ab19An19Value     uint32
	Ab19An19Fresh     bool
}

type Ab19An19Pair struct {
	Ab19An19Minimum Ab19An19Field
	Ab19An19Maximum Ab19An19Field
	Ab19An19Fresh   bool
}

type Ab19An19Result struct {
	Ab19An19Admitted            bool
	Ab19An19Attempt             Ab19An19Pair
	Ab19An19Prior               Ab19An19Pair
	Ab19An19LastKnown           Ab19An19Pair
	Ab19An19ContractGrade       string
	Ab19An19InterpretationGrade string
	Ab19An19NativeUnit          string
	Ab19An19Qualification       string
	Ab19An19Presence            string
	Ab19An19Execution           string
}

func ab19an19Clone(ab19an19Field Ab19An19Field) Ab19An19Field {
	ab19an19Field.Ab19An19Raw = append([]byte(nil), ab19an19Field.Ab19An19Raw...)
	return ab19an19Field
}

func ab19an19Decode(ab19an19Raw []byte, ab19an19ReadOK, ab19an19Admitted bool) Ab19An19Field {
	ab19an19Field := Ab19An19Field{Ab19An19Raw: append([]byte(nil), ab19an19Raw...), Ab19An19State: "not-admitted"}
	if !ab19an19Admitted {
		return ab19an19Field
	}
	switch {
	case !ab19an19ReadOK:
		ab19an19Field.Ab19An19State = "read-failed"
	case len(ab19an19Raw) == 0:
		ab19an19Field.Ab19An19State = "missing"
	case len(ab19an19Raw) != 2:
		ab19an19Field.Ab19An19State = "malformed"
	default:
		ab19an19Field.Ab19An19State = "success"
		ab19an19Field.Ab19An19Available = true
		ab19an19Field.Ab19An19Fresh = true
		ab19an19Field.Ab19An19Value = (uint32(ab19an19Raw[0]) | uint32(ab19an19Raw[1])<<8) * 10
	}
	return ab19an19Field
}

func ab19an19Retain(ab19an19Attempt, ab19an19Prior Ab19An19Field) Ab19An19Field {
	if ab19an19Attempt.Ab19An19Available {
		return ab19an19Clone(ab19an19Attempt)
	}
	ab19an19Prior = ab19an19Clone(ab19an19Prior)
	ab19an19Prior.Ab19An19Fresh = false
	return ab19an19Prior
}

func Ab19An19Interpret(ab19an19Family, ab19an19Variant string, ab19an19Operation byte, ab19an19Direction string, ab19an19MinimumRaw, ab19an19MaximumRaw []byte, ab19an19MinimumReadOK, ab19an19MaximumReadOK bool, ab19an19Prior Ab19An19Pair) Ab19An19Result {
	ab19an19Result := Ab19An19Result{
		Ab19An19Admitted:      ab19an19Family == Ab19An19Family && ab19an19Variant == Ab19An19Variant && ab19an19Operation == 0x02 && ab19an19Direction == Ab19An19Direction,
		Ab19An19Prior:         Ab19An19Pair{Ab19An19Minimum: ab19an19Clone(ab19an19Prior.Ab19An19Minimum), Ab19An19Maximum: ab19an19Clone(ab19an19Prior.Ab19An19Maximum), Ab19An19Fresh: ab19an19Prior.Ab19An19Fresh},
		Ab19An19ContractGrade: "Proven", Ab19An19InterpretationGrade: "Hypothesis", Ab19An19NativeUnit: "Unknown", Ab19An19Qualification: "Unknown", Ab19An19Presence: "Unknown", Ab19An19Execution: "Unknown",
	}
	ab19an19Result.Ab19An19Attempt.Ab19An19Minimum = ab19an19Decode(ab19an19MinimumRaw, ab19an19MinimumReadOK, ab19an19Result.Ab19An19Admitted)
	ab19an19Result.Ab19An19Attempt.Ab19An19Maximum = ab19an19Decode(ab19an19MaximumRaw, ab19an19MaximumReadOK, ab19an19Result.Ab19An19Admitted)
	ab19an19Result.Ab19An19Attempt.Ab19An19Fresh = ab19an19Result.Ab19An19Attempt.Ab19An19Minimum.Ab19An19Fresh && ab19an19Result.Ab19An19Attempt.Ab19An19Maximum.Ab19An19Fresh
	ab19an19Result.Ab19An19LastKnown = Ab19An19Pair{
		Ab19An19Minimum: ab19an19Retain(ab19an19Result.Ab19An19Attempt.Ab19An19Minimum, ab19an19Prior.Ab19An19Minimum),
		Ab19An19Maximum: ab19an19Retain(ab19an19Result.Ab19An19Attempt.Ab19An19Maximum, ab19an19Prior.Ab19An19Maximum),
		Ab19An19Fresh:   ab19an19Result.Ab19An19Attempt.Ab19An19Fresh,
	}
	return ab19an19Result
}

func Ab19An19Getter(ab19an19Raw uint32) uint16 {
	return uint16(ab19an19Raw / 10)
}
