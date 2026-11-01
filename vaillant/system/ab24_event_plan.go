package system

import "fmt"

type Ab24EventIdentity struct {
	Ab24Family     string
	Ab24SystemType int
	Ab24Instance   int
	Ab24Address    int
	Ab24Code       int
}

type Ab24EventBaseline struct {
	Ab24Identity Ab24EventIdentity
	Ab24Opcode   int
	Ab24Success  bool
	Ab24Raw      []byte
}

type Ab24EventEdit struct {
	Ab24Identity Ab24EventIdentity
	Ab24Opcode   int
	Ab24Values   []int
	Ab24OP09     *Ab24EventBaseline
	Ab24OP0B     *Ab24EventBaseline
}

type Ab24EventPlan struct {
	Ab24Identity Ab24EventIdentity
	Ab24Payload  [12]byte
	Ab24Values   [7]byte
	Ab24OP09     [8]byte
	Ab24OP0B     [8]byte
}

func Ab24BuildEventPlan(ab24Edit Ab24EventEdit) (*Ab24EventPlan, error) {
	ab24Identity := ab24Edit.Ab24Identity
	if ab24Identity.Ab24Family == "" {
		return nil, fmt.Errorf("ab24: missing family")
	}
	for _, ab24Field := range []int{ab24Identity.Ab24SystemType, ab24Identity.Ab24Instance, ab24Identity.Ab24Address, ab24Identity.Ab24Code, ab24Edit.Ab24Opcode} {
		if ab24Field < 0 || ab24Field > 255 {
			return nil, fmt.Errorf("ab24: field outside raw u8 range")
		}
	}
	ab24MaxAddress := 2
	switch ab24Identity.Ab24SystemType {
	case 0:
		ab24MaxAddress = 3
	case 1, 3:
	default:
		return nil, fmt.Errorf("ab24: unsupported system type")
	}
	if ab24Identity.Ab24Address < 1 || ab24Identity.Ab24Address > ab24MaxAddress {
		return nil, fmt.Errorf("ab24: invalid address")
	}
	if ab24Edit.Ab24Opcode != 0x0a && ab24Edit.Ab24Opcode != 0x0c {
		return nil, fmt.Errorf("ab24: invalid edit opcode")
	}
	if len(ab24Edit.Ab24Values) != 7 {
		return nil, fmt.Errorf("ab24: exactly seven raw values required")
	}
	for _, ab24Value := range ab24Edit.Ab24Values {
		if ab24Value < 0 || ab24Value > 255 {
			return nil, fmt.Errorf("ab24: value outside raw u8 range")
		}
	}
	for ab24Index, ab24Baseline := range []*Ab24EventBaseline{ab24Edit.Ab24OP09, ab24Edit.Ab24OP0B} {
		ab24Opcode := 0x09 + ab24Index*2
		if ab24Baseline == nil || !ab24Baseline.Ab24Success || ab24Baseline.Ab24Opcode != ab24Opcode || ab24Baseline.Ab24Identity != ab24Identity || len(ab24Baseline.Ab24Raw) != 8 {
			return nil, fmt.Errorf("ab24: missing, failed or uncorrelated baseline %02x", ab24Opcode)
		}
	}
	ab24Plan := &Ab24EventPlan{Ab24Identity: ab24Identity}
	ab24Plan.Ab24Payload = [12]byte{byte(ab24Edit.Ab24Opcode), byte(ab24Identity.Ab24SystemType), byte(ab24Identity.Ab24Instance), byte(ab24Identity.Ab24Address), byte(ab24Identity.Ab24Code)}
	for ab24Index, ab24Value := range ab24Edit.Ab24Values {
		ab24Plan.Ab24Values[ab24Index] = byte(ab24Value)
		ab24Plan.Ab24Payload[5+ab24Index] = byte(ab24Value)
	}
	copy(ab24Plan.Ab24OP09[:], ab24Edit.Ab24OP09.Ab24Raw)
	copy(ab24Plan.Ab24OP0B[:], ab24Edit.Ab24OP0B.Ab24Raw)
	return ab24Plan, nil
}
