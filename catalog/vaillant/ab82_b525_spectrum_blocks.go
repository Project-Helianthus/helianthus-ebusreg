package vaillant

import "sort"

type Ab82State string

const (
	Ab82StateValid      Ab82State = "valid"
	Ab82StateShort      Ab82State = "short"
	Ab82StateExtra      Ab82State = "extra"
	Ab82StatePrefixOnly Ab82State = "prefix-only"
	Ab82StateRawOnly    Ab82State = "raw-only"
)

type Ab82JoinState string

const (
	Ab82JoinEmpty           Ab82JoinState = "empty"
	Ab82JoinPartial         Ab82JoinState = "partial"
	Ab82JoinConflict        Ab82JoinState = "conflict"
	Ab82JoinContextMismatch Ab82JoinState = "context-mismatch"
	Ab82JoinComplete        Ab82JoinState = "complete"
)

const (
	Ab82GradeHypothesis   = "Hypothesis"
	Ab82DirectionRequest  = "request"
	Ab82DirectionPositive = "positive"
	Ab82SelectorFirst     = 10
	Ab82SelectorLast      = 17
	Ab82BlockLength       = 32
	Ab82IndexCount        = 256
)

type Ab82Context struct {
	Ab82Primary   byte
	Ab82Secondary byte
	Ab82Family    string
	Ab82Variant   string
	Ab82Operation int
	Ab82Direction string
	Ab82Selector  int
}

type Ab82Input struct {
	Ab82Context Ab82Context
	Ab82Values  []byte
}

type Ab82Entry struct {
	Ab82Index   int
	Ab82RawByte byte
	Ab82Signed  int8
}

type Ab82Block struct {
	Ab82Context    Ab82Context
	Ab82State      Ab82State
	Ab82Grade      string
	Ab82Start      int
	Ab82Raw        []byte
	Ab82Entries    []Ab82Entry
	Ab82Mask       [4]uint64
	Ab82Incomplete bool
	Ab82Extra      bool
	Ab82Remainder  []byte
}

type Ab82JoinBlock struct {
	Ab82Selector    int
	Ab82State       Ab82State
	Ab82Contributed bool
	Ab82Conflict    bool
}

type Ab82Joined struct {
	Ab82State     Ab82JoinState
	Ab82Grade     string
	Ab82Complete  bool
	Ab82Blocks    []Ab82JoinBlock
	Ab82Entries   []Ab82Entry
	Ab82Mask      [4]uint64
	Ab82Missing   []int
	Ab82Conflicts []int
}

func ab82Bytes(ab82Source []byte) []byte {
	if ab82Source == nil {
		return nil
	}
	ab82Copy := make([]byte, len(ab82Source))
	copy(ab82Copy, ab82Source)
	return ab82Copy
}

func ab82KnownIdentity(ab82Value Ab82Context) bool {
	return ab82Value.Ab82Primary == 0xB5 && ab82Value.Ab82Secondary == 0x25 &&
		ab82Value.Ab82Family != "" && ab82Value.Ab82Variant != "" &&
		ab82Value.Ab82Operation >= 0 && ab82Value.Ab82Operation <= 255
}

func ab82KnownSelector(ab82Selector int) bool {
	return ab82Selector >= Ab82SelectorFirst && ab82Selector <= Ab82SelectorLast
}

func ab82SetMask(ab82Mask *[4]uint64, ab82Index int) {
	ab82Mask[ab82Index/64] |= 1 << uint(ab82Index%64)
}

func (ab82Value Ab82Block) Ab82Present(ab82Index int) bool {
	if ab82Index < 0 || ab82Index >= Ab82IndexCount {
		return false
	}
	return ab82Value.Ab82Mask[ab82Index/64]&(1<<uint(ab82Index%64)) != 0
}

func (ab82Value Ab82Joined) Ab82Present(ab82Index int) bool {
	if ab82Index < 0 || ab82Index >= Ab82IndexCount {
		return false
	}
	return ab82Value.Ab82Mask[ab82Index/64]&(1<<uint(ab82Index%64)) != 0
}

func (ab82Value Ab82Block) Ab82Clone() Ab82Block {
	ab82Copy := ab82Value
	ab82Copy.Ab82Raw = ab82Bytes(ab82Value.Ab82Raw)
	ab82Copy.Ab82Remainder = ab82Bytes(ab82Value.Ab82Remainder)
	if ab82Value.Ab82Entries != nil {
		ab82Copy.Ab82Entries = append([]Ab82Entry{}, ab82Value.Ab82Entries...)
	}
	return ab82Copy
}

func Ab82View(ab82Input Ab82Input) Ab82Block {
	ab82Block := Ab82Block{
		Ab82Context: ab82Input.Ab82Context,
		Ab82Grade:   Ab82GradeHypothesis,
		Ab82Start:   -1,
		Ab82Raw:     ab82Bytes(ab82Input.Ab82Values),
	}
	ab82Context := ab82Input.Ab82Context
	if !ab82KnownIdentity(ab82Context) || !ab82KnownSelector(ab82Context.Ab82Selector) ||
		(ab82Context.Ab82Direction != Ab82DirectionRequest && ab82Context.Ab82Direction != Ab82DirectionPositive) {
		ab82Block.Ab82State = Ab82StateRawOnly
		return ab82Block
	}
	if ab82Context.Ab82Direction == Ab82DirectionPositive {
		ab82Block.Ab82State = Ab82StatePrefixOnly
		return ab82Block
	}
	ab82Block.Ab82Start = (ab82Context.Ab82Selector - Ab82SelectorFirst) * Ab82BlockLength
	ab82Count := len(ab82Block.Ab82Raw)
	if ab82Count > Ab82BlockLength {
		ab82Count = Ab82BlockLength
	}
	ab82Block.Ab82Entries = make([]Ab82Entry, 0, ab82Count)
	for ab82Offset := 0; ab82Offset < ab82Count; ab82Offset++ {
		ab82Raw := ab82Block.Ab82Raw[ab82Offset]
		ab82Index := ab82Block.Ab82Start + ab82Offset
		ab82Block.Ab82Entries = append(ab82Block.Ab82Entries, Ab82Entry{Ab82Index: ab82Index, Ab82RawByte: ab82Raw, Ab82Signed: int8(ab82Raw)})
		ab82SetMask(&ab82Block.Ab82Mask, ab82Index)
	}
	switch {
	case len(ab82Block.Ab82Raw) < Ab82BlockLength:
		ab82Block.Ab82State = Ab82StateShort
		ab82Block.Ab82Incomplete = true
	case len(ab82Block.Ab82Raw) > Ab82BlockLength:
		ab82Block.Ab82State = Ab82StateExtra
		ab82Block.Ab82Extra = true
		ab82Block.Ab82Remainder = ab82Bytes(ab82Block.Ab82Raw[Ab82BlockLength:])
	default:
		ab82Block.Ab82State = Ab82StateValid
	}
	return ab82Block
}

type ab82JoinKey struct {
	ab82Primary, ab82Secondary byte
	ab82Family, ab82Variant    string
	ab82Operation              int
	ab82Direction              string
}

func ab82KeyOf(ab82Value Ab82Context) ab82JoinKey {
	return ab82JoinKey{ab82Value.Ab82Primary, ab82Value.Ab82Secondary, ab82Value.Ab82Family, ab82Value.Ab82Variant, ab82Value.Ab82Operation, ab82Value.Ab82Direction}
}

func ab82SameBytes(ab82Left, ab82Right []byte) bool {
	if len(ab82Left) != len(ab82Right) {
		return false
	}
	for ab82Index := range ab82Left {
		if ab82Left[ab82Index] != ab82Right[ab82Index] {
			return false
		}
	}
	return true
}

func Ab82Join(ab82Supplied []Ab82Block) Ab82Joined {
	ab82Result := Ab82Joined{Ab82Grade: Ab82GradeHypothesis, Ab82State: Ab82JoinEmpty}
	ab82Views := make([]Ab82Block, len(ab82Supplied))
	for ab82Index, ab82Item := range ab82Supplied {
		ab82Views[ab82Index] = Ab82View(Ab82Input{Ab82Context: ab82Item.Ab82Context, Ab82Values: ab82Item.Ab82Raw})
		ab82Result.Ab82Blocks = append(ab82Result.Ab82Blocks, Ab82JoinBlock{Ab82Selector: ab82Item.Ab82Context.Ab82Selector, Ab82State: ab82Views[ab82Index].Ab82State})
	}
	var ab82Key ab82JoinKey
	ab82HaveKey := false
	ab82Mismatch := false
	for _, ab82Item := range ab82Views {
		if ab82Item.Ab82Start < 0 {
			continue
		}
		if !ab82HaveKey {
			ab82Key, ab82HaveKey = ab82KeyOf(ab82Item.Ab82Context), true
		} else if ab82KeyOf(ab82Item.Ab82Context) != ab82Key {
			ab82Mismatch = true
		}
	}
	if ab82Mismatch {
		ab82Result.Ab82State = Ab82JoinContextMismatch
		for ab82Selector := Ab82SelectorFirst; ab82Selector <= Ab82SelectorLast; ab82Selector++ {
			ab82Result.Ab82Missing = append(ab82Result.Ab82Missing, ab82Selector)
		}
		return ab82Result
	}
	ab82Chosen := map[int]int{}
	ab82Conflicted := map[int]bool{}
	for ab82Index, ab82Item := range ab82Views {
		if ab82Item.Ab82Start < 0 {
			continue
		}
		ab82Selector := ab82Item.Ab82Context.Ab82Selector
		ab82Prior, ab82Seen := ab82Chosen[ab82Selector]
		if !ab82Seen {
			ab82Chosen[ab82Selector] = ab82Index
			continue
		}
		if !ab82SameBytes(ab82Views[ab82Prior].Ab82Raw, ab82Item.Ab82Raw) {
			ab82Conflicted[ab82Selector] = true
		}
	}
	ab82Selectors := make([]int, 0, len(ab82Chosen))
	for ab82Selector := range ab82Chosen {
		ab82Selectors = append(ab82Selectors, ab82Selector)
	}
	sort.Ints(ab82Selectors)
	ab82ValidCount := 0
	for _, ab82Selector := range ab82Selectors {
		if ab82Conflicted[ab82Selector] {
			ab82Result.Ab82Conflicts = append(ab82Result.Ab82Conflicts, ab82Selector)
			continue
		}
		ab82Item := ab82Views[ab82Chosen[ab82Selector]]
		if ab82Item.Ab82State == Ab82StateValid {
			ab82ValidCount++
		}
		for _, ab82Entry := range ab82Item.Ab82Entries {
			ab82Result.Ab82Entries = append(ab82Result.Ab82Entries, ab82Entry)
			ab82SetMask(&ab82Result.Ab82Mask, ab82Entry.Ab82Index)
		}
	}
	for ab82Index := range ab82Result.Ab82Blocks {
		ab82Selector := ab82Result.Ab82Blocks[ab82Index].Ab82Selector
		if ab82Views[ab82Index].Ab82Start < 0 {
			continue
		}
		ab82Result.Ab82Blocks[ab82Index].Ab82Conflict = ab82Conflicted[ab82Selector]
		ab82Result.Ab82Blocks[ab82Index].Ab82Contributed = !ab82Conflicted[ab82Selector] && ab82Chosen[ab82Selector] == ab82Index
	}
	for ab82Selector := Ab82SelectorFirst; ab82Selector <= Ab82SelectorLast; ab82Selector++ {
		if _, ab82Has := ab82Chosen[ab82Selector]; !ab82Has || ab82Conflicted[ab82Selector] {
			ab82Result.Ab82Missing = append(ab82Result.Ab82Missing, ab82Selector)
		}
	}
	switch {
	case len(ab82Result.Ab82Conflicts) > 0:
		ab82Result.Ab82State = Ab82JoinConflict
	case ab82ValidCount == Ab82SelectorLast-Ab82SelectorFirst+1:
		ab82Result.Ab82State = Ab82JoinComplete
		ab82Result.Ab82Complete = true
	case len(ab82Result.Ab82Entries) > 0:
		ab82Result.Ab82State = Ab82JoinPartial
	}
	return ab82Result
}
