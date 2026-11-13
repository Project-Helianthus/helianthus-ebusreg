package vaillant

import "sort"

type B525SpectrumBlockState string

const (
	B525SpectrumBlockStateValid      B525SpectrumBlockState = "valid"
	B525SpectrumBlockStateShort      B525SpectrumBlockState = "short"
	B525SpectrumBlockStateExtra      B525SpectrumBlockState = "extra"
	B525SpectrumBlockStatePrefixOnly B525SpectrumBlockState = "prefix-only"
	B525SpectrumBlockStateRawOnly    B525SpectrumBlockState = "raw-only"
)

type B525SpectrumBlockJoinState string

const (
	B525SpectrumBlockJoinEmpty           B525SpectrumBlockJoinState = "empty"
	B525SpectrumBlockJoinPartial         B525SpectrumBlockJoinState = "partial"
	B525SpectrumBlockJoinConflict        B525SpectrumBlockJoinState = "conflict"
	B525SpectrumBlockJoinContextMismatch B525SpectrumBlockJoinState = "context-mismatch"
	B525SpectrumBlockJoinComplete        B525SpectrumBlockJoinState = "complete"
)

const (
	B525SpectrumBlockGradeHypothesis   = "Hypothesis"
	B525SpectrumBlockDirectionRequest  = "request"
	B525SpectrumBlockDirectionPositive = "positive"
	B525SpectrumBlockSelectorFirst     = 10
	B525SpectrumBlockSelectorLast      = 17
	B525SpectrumBlockLength            = 32
	B525SpectrumBlockIndexCount        = 256
)

type B525SpectrumBlockContext struct {
	B525SpectrumBlockPrimary   byte   `json:"Ab82Primary"`
	B525SpectrumBlockSecondary byte   `json:"Ab82Secondary"`
	B525SpectrumBlockFamily    string `json:"Ab82Family"`
	B525SpectrumBlockVariant   string `json:"Ab82Variant"`
	B525SpectrumBlockOperation int    `json:"Ab82Operation"`
	B525SpectrumBlockDirection string `json:"Ab82Direction"`
	B525SpectrumBlockSelector  int    `json:"Ab82Selector"`
}

type B525SpectrumBlockInput struct {
	B525SpectrumBlockContext B525SpectrumBlockContext `json:"Ab82Context"`
	B525SpectrumBlockValues  []byte                   `json:"Ab82Values"`
}

type B525SpectrumBlockEntry struct {
	B525SpectrumBlockIndex   int  `json:"Ab82Index"`
	B525SpectrumBlockRawByte byte `json:"Ab82RawByte"`
	B525SpectrumBlockSigned  int8 `json:"Ab82Signed"`
}

type B525SpectrumBlock struct {
	B525SpectrumBlockContext    B525SpectrumBlockContext `json:"Ab82Context"`
	B525SpectrumBlockState      B525SpectrumBlockState   `json:"Ab82State"`
	B525SpectrumBlockGrade      string                   `json:"Ab82Grade"`
	B525SpectrumBlockStart      int                      `json:"Ab82Start"`
	B525SpectrumBlockRaw        []byte                   `json:"Ab82Raw"`
	B525SpectrumBlockEntries    []B525SpectrumBlockEntry `json:"Ab82Entries"`
	B525SpectrumBlockMask       [4]uint64                `json:"Ab82Mask"`
	B525SpectrumBlockIncomplete bool                     `json:"Ab82Incomplete"`
	B525SpectrumBlockExtra      bool                     `json:"Ab82Extra"`
	B525SpectrumBlockRemainder  []byte                   `json:"Ab82Remainder"`
}

type B525SpectrumJoinBlock struct {
	B525SpectrumBlockSelector    int                    `json:"Ab82Selector"`
	B525SpectrumBlockState       B525SpectrumBlockState `json:"Ab82State"`
	B525SpectrumBlockContributed bool                   `json:"Ab82Contributed"`
	B525SpectrumBlockConflict    bool                   `json:"Ab82Conflict"`
}

type B525SpectrumBlockJoined struct {
	B525SpectrumBlockState     B525SpectrumBlockJoinState `json:"Ab82State"`
	B525SpectrumBlockGrade     string                     `json:"Ab82Grade"`
	B525SpectrumBlockComplete  bool                       `json:"Ab82Complete"`
	B525SpectrumBlocks         []B525SpectrumJoinBlock    `json:"Ab82Blocks"`
	B525SpectrumBlockEntries   []B525SpectrumBlockEntry   `json:"Ab82Entries"`
	B525SpectrumBlockMask      [4]uint64                  `json:"Ab82Mask"`
	B525SpectrumBlockMissing   []int                      `json:"Ab82Missing"`
	B525SpectrumBlockConflicts []int                      `json:"Ab82Conflicts"`
}

func b525SpectrumBlockBytes(b525SpectrumBlockSource []byte) []byte {
	if b525SpectrumBlockSource == nil {
		return nil
	}
	b525SpectrumBlockCopy := make([]byte, len(b525SpectrumBlockSource))
	copy(b525SpectrumBlockCopy, b525SpectrumBlockSource)
	return b525SpectrumBlockCopy
}

func b525SpectrumBlockKnownIdentity(b525SpectrumBlockValue B525SpectrumBlockContext) bool {
	return b525SpectrumBlockValue.B525SpectrumBlockPrimary == 0xB5 && b525SpectrumBlockValue.B525SpectrumBlockSecondary == 0x25 &&
		b525SpectrumBlockValue.B525SpectrumBlockFamily != "" && b525SpectrumBlockValue.B525SpectrumBlockVariant != "" &&
		b525SpectrumBlockValue.B525SpectrumBlockOperation >= 0 && b525SpectrumBlockValue.B525SpectrumBlockOperation <= 255
}

func b525SpectrumBlockKnownSelector(b525SpectrumBlockSelector int) bool {
	return b525SpectrumBlockSelector >= B525SpectrumBlockSelectorFirst && b525SpectrumBlockSelector <= B525SpectrumBlockSelectorLast
}

func b525SpectrumBlockSetMask(b525SpectrumBlockMask *[4]uint64, b525SpectrumBlockIndex int) {
	b525SpectrumBlockMask[b525SpectrumBlockIndex/64] |= 1 << uint(b525SpectrumBlockIndex%64)
}

func (b525SpectrumBlockValue B525SpectrumBlock) B525SpectrumBlockPresent(b525SpectrumBlockIndex int) bool {
	if b525SpectrumBlockIndex < 0 || b525SpectrumBlockIndex >= B525SpectrumBlockIndexCount {
		return false
	}
	return b525SpectrumBlockValue.B525SpectrumBlockMask[b525SpectrumBlockIndex/64]&(1<<uint(b525SpectrumBlockIndex%64)) != 0
}

func (b525SpectrumBlockValue B525SpectrumBlockJoined) B525SpectrumBlockPresent(b525SpectrumBlockIndex int) bool {
	if b525SpectrumBlockIndex < 0 || b525SpectrumBlockIndex >= B525SpectrumBlockIndexCount {
		return false
	}
	return b525SpectrumBlockValue.B525SpectrumBlockMask[b525SpectrumBlockIndex/64]&(1<<uint(b525SpectrumBlockIndex%64)) != 0
}

func (b525SpectrumBlockValue B525SpectrumBlock) B525SpectrumBlockClone() B525SpectrumBlock {
	b525SpectrumBlockCopy := b525SpectrumBlockValue
	b525SpectrumBlockCopy.B525SpectrumBlockRaw = b525SpectrumBlockBytes(b525SpectrumBlockValue.B525SpectrumBlockRaw)
	b525SpectrumBlockCopy.B525SpectrumBlockRemainder = b525SpectrumBlockBytes(b525SpectrumBlockValue.B525SpectrumBlockRemainder)
	if b525SpectrumBlockValue.B525SpectrumBlockEntries != nil {
		b525SpectrumBlockCopy.B525SpectrumBlockEntries = append([]B525SpectrumBlockEntry{}, b525SpectrumBlockValue.B525SpectrumBlockEntries...)
	}
	return b525SpectrumBlockCopy
}

func B525SpectrumBlockView(b525SpectrumBlockInput B525SpectrumBlockInput) B525SpectrumBlock {
	b525SpectrumBlockBlock := B525SpectrumBlock{
		B525SpectrumBlockContext: b525SpectrumBlockInput.B525SpectrumBlockContext,
		B525SpectrumBlockGrade:   B525SpectrumBlockGradeHypothesis,
		B525SpectrumBlockStart:   -1,
		B525SpectrumBlockRaw:     b525SpectrumBlockBytes(b525SpectrumBlockInput.B525SpectrumBlockValues),
	}
	b525SpectrumBlockContext := b525SpectrumBlockInput.B525SpectrumBlockContext
	if !b525SpectrumBlockKnownIdentity(b525SpectrumBlockContext) || !b525SpectrumBlockKnownSelector(b525SpectrumBlockContext.B525SpectrumBlockSelector) ||
		(b525SpectrumBlockContext.B525SpectrumBlockDirection != B525SpectrumBlockDirectionRequest && b525SpectrumBlockContext.B525SpectrumBlockDirection != B525SpectrumBlockDirectionPositive) {
		b525SpectrumBlockBlock.B525SpectrumBlockState = B525SpectrumBlockStateRawOnly
		return b525SpectrumBlockBlock
	}
	if b525SpectrumBlockContext.B525SpectrumBlockDirection == B525SpectrumBlockDirectionPositive {
		b525SpectrumBlockBlock.B525SpectrumBlockState = B525SpectrumBlockStatePrefixOnly
		return b525SpectrumBlockBlock
	}
	b525SpectrumBlockBlock.B525SpectrumBlockStart = (b525SpectrumBlockContext.B525SpectrumBlockSelector - B525SpectrumBlockSelectorFirst) * B525SpectrumBlockLength
	b525SpectrumBlockCount := len(b525SpectrumBlockBlock.B525SpectrumBlockRaw)
	if b525SpectrumBlockCount > B525SpectrumBlockLength {
		b525SpectrumBlockCount = B525SpectrumBlockLength
	}
	b525SpectrumBlockBlock.B525SpectrumBlockEntries = make([]B525SpectrumBlockEntry, 0, b525SpectrumBlockCount)
	for b525SpectrumBlockOffset := 0; b525SpectrumBlockOffset < b525SpectrumBlockCount; b525SpectrumBlockOffset++ {
		b525SpectrumBlockRaw := b525SpectrumBlockBlock.B525SpectrumBlockRaw[b525SpectrumBlockOffset]
		b525SpectrumBlockIndex := b525SpectrumBlockBlock.B525SpectrumBlockStart + b525SpectrumBlockOffset
		b525SpectrumBlockBlock.B525SpectrumBlockEntries = append(b525SpectrumBlockBlock.B525SpectrumBlockEntries, B525SpectrumBlockEntry{B525SpectrumBlockIndex: b525SpectrumBlockIndex, B525SpectrumBlockRawByte: b525SpectrumBlockRaw, B525SpectrumBlockSigned: int8(b525SpectrumBlockRaw)})
		b525SpectrumBlockSetMask(&b525SpectrumBlockBlock.B525SpectrumBlockMask, b525SpectrumBlockIndex)
	}
	switch {
	case len(b525SpectrumBlockBlock.B525SpectrumBlockRaw) < B525SpectrumBlockLength:
		b525SpectrumBlockBlock.B525SpectrumBlockState = B525SpectrumBlockStateShort
		b525SpectrumBlockBlock.B525SpectrumBlockIncomplete = true
	case len(b525SpectrumBlockBlock.B525SpectrumBlockRaw) > B525SpectrumBlockLength:
		b525SpectrumBlockBlock.B525SpectrumBlockState = B525SpectrumBlockStateExtra
		b525SpectrumBlockBlock.B525SpectrumBlockExtra = true
		b525SpectrumBlockBlock.B525SpectrumBlockRemainder = b525SpectrumBlockBytes(b525SpectrumBlockBlock.B525SpectrumBlockRaw[B525SpectrumBlockLength:])
	default:
		b525SpectrumBlockBlock.B525SpectrumBlockState = B525SpectrumBlockStateValid
	}
	return b525SpectrumBlockBlock
}

type b525SpectrumBlockJoinKey struct {
	b525SpectrumBlockPrimary, b525SpectrumBlockSecondary byte
	b525SpectrumBlockFamily, b525SpectrumBlockVariant    string
	b525SpectrumBlockOperation                           int
	b525SpectrumBlockDirection                           string
}

func b525SpectrumBlockKeyOf(b525SpectrumBlockValue B525SpectrumBlockContext) b525SpectrumBlockJoinKey {
	return b525SpectrumBlockJoinKey{b525SpectrumBlockValue.B525SpectrumBlockPrimary, b525SpectrumBlockValue.B525SpectrumBlockSecondary, b525SpectrumBlockValue.B525SpectrumBlockFamily, b525SpectrumBlockValue.B525SpectrumBlockVariant, b525SpectrumBlockValue.B525SpectrumBlockOperation, b525SpectrumBlockValue.B525SpectrumBlockDirection}
}

func b525SpectrumBlockSameBytes(b525SpectrumBlockLeft, b525SpectrumBlockRight []byte) bool {
	if len(b525SpectrumBlockLeft) != len(b525SpectrumBlockRight) {
		return false
	}
	for b525SpectrumBlockIndex := range b525SpectrumBlockLeft {
		if b525SpectrumBlockLeft[b525SpectrumBlockIndex] != b525SpectrumBlockRight[b525SpectrumBlockIndex] {
			return false
		}
	}
	return true
}

func B525SpectrumBlockJoin(b525SpectrumBlockSupplied []B525SpectrumBlock) B525SpectrumBlockJoined {
	b525SpectrumBlockResult := B525SpectrumBlockJoined{B525SpectrumBlockGrade: B525SpectrumBlockGradeHypothesis, B525SpectrumBlockState: B525SpectrumBlockJoinEmpty}
	b525SpectrumBlockViews := make([]B525SpectrumBlock, len(b525SpectrumBlockSupplied))
	for b525SpectrumBlockIndex, b525SpectrumBlockItem := range b525SpectrumBlockSupplied {
		b525SpectrumBlockViews[b525SpectrumBlockIndex] = B525SpectrumBlockView(B525SpectrumBlockInput{B525SpectrumBlockContext: b525SpectrumBlockItem.B525SpectrumBlockContext, B525SpectrumBlockValues: b525SpectrumBlockItem.B525SpectrumBlockRaw})
		b525SpectrumBlockResult.B525SpectrumBlocks = append(b525SpectrumBlockResult.B525SpectrumBlocks, B525SpectrumJoinBlock{B525SpectrumBlockSelector: b525SpectrumBlockItem.B525SpectrumBlockContext.B525SpectrumBlockSelector, B525SpectrumBlockState: b525SpectrumBlockViews[b525SpectrumBlockIndex].B525SpectrumBlockState})
	}
	var b525SpectrumBlockKey b525SpectrumBlockJoinKey
	b525SpectrumBlockHaveKey := false
	b525SpectrumBlockMismatch := false
	for _, b525SpectrumBlockItem := range b525SpectrumBlockViews {
		if b525SpectrumBlockItem.B525SpectrumBlockStart < 0 {
			continue
		}
		if !b525SpectrumBlockHaveKey {
			b525SpectrumBlockKey, b525SpectrumBlockHaveKey = b525SpectrumBlockKeyOf(b525SpectrumBlockItem.B525SpectrumBlockContext), true
		} else if b525SpectrumBlockKeyOf(b525SpectrumBlockItem.B525SpectrumBlockContext) != b525SpectrumBlockKey {
			b525SpectrumBlockMismatch = true
		}
	}
	if b525SpectrumBlockMismatch {
		b525SpectrumBlockResult.B525SpectrumBlockState = B525SpectrumBlockJoinContextMismatch
		for b525SpectrumBlockSelector := B525SpectrumBlockSelectorFirst; b525SpectrumBlockSelector <= B525SpectrumBlockSelectorLast; b525SpectrumBlockSelector++ {
			b525SpectrumBlockResult.B525SpectrumBlockMissing = append(b525SpectrumBlockResult.B525SpectrumBlockMissing, b525SpectrumBlockSelector)
		}
		return b525SpectrumBlockResult
	}
	b525SpectrumBlockChosen := map[int]int{}
	b525SpectrumBlockConflicted := map[int]bool{}
	for b525SpectrumBlockIndex, b525SpectrumBlockItem := range b525SpectrumBlockViews {
		if b525SpectrumBlockItem.B525SpectrumBlockStart < 0 {
			continue
		}
		b525SpectrumBlockSelector := b525SpectrumBlockItem.B525SpectrumBlockContext.B525SpectrumBlockSelector
		b525SpectrumBlockPrior, b525SpectrumBlockSeen := b525SpectrumBlockChosen[b525SpectrumBlockSelector]
		if !b525SpectrumBlockSeen {
			b525SpectrumBlockChosen[b525SpectrumBlockSelector] = b525SpectrumBlockIndex
			continue
		}
		if !b525SpectrumBlockSameBytes(b525SpectrumBlockViews[b525SpectrumBlockPrior].B525SpectrumBlockRaw, b525SpectrumBlockItem.B525SpectrumBlockRaw) {
			b525SpectrumBlockConflicted[b525SpectrumBlockSelector] = true
		}
	}
	b525SpectrumBlockSelectors := make([]int, 0, len(b525SpectrumBlockChosen))
	for b525SpectrumBlockSelector := range b525SpectrumBlockChosen {
		b525SpectrumBlockSelectors = append(b525SpectrumBlockSelectors, b525SpectrumBlockSelector)
	}
	sort.Ints(b525SpectrumBlockSelectors)
	b525SpectrumBlockValidCount := 0
	for _, b525SpectrumBlockSelector := range b525SpectrumBlockSelectors {
		if b525SpectrumBlockConflicted[b525SpectrumBlockSelector] {
			b525SpectrumBlockResult.B525SpectrumBlockConflicts = append(b525SpectrumBlockResult.B525SpectrumBlockConflicts, b525SpectrumBlockSelector)
			continue
		}
		b525SpectrumBlockItem := b525SpectrumBlockViews[b525SpectrumBlockChosen[b525SpectrumBlockSelector]]
		if b525SpectrumBlockItem.B525SpectrumBlockState == B525SpectrumBlockStateValid {
			b525SpectrumBlockValidCount++
		}
		for _, b525SpectrumBlockEntry := range b525SpectrumBlockItem.B525SpectrumBlockEntries {
			b525SpectrumBlockResult.B525SpectrumBlockEntries = append(b525SpectrumBlockResult.B525SpectrumBlockEntries, b525SpectrumBlockEntry)
			b525SpectrumBlockSetMask(&b525SpectrumBlockResult.B525SpectrumBlockMask, b525SpectrumBlockEntry.B525SpectrumBlockIndex)
		}
	}
	for b525SpectrumBlockIndex := range b525SpectrumBlockResult.B525SpectrumBlocks {
		b525SpectrumBlockSelector := b525SpectrumBlockResult.B525SpectrumBlocks[b525SpectrumBlockIndex].B525SpectrumBlockSelector
		if b525SpectrumBlockViews[b525SpectrumBlockIndex].B525SpectrumBlockStart < 0 {
			continue
		}
		b525SpectrumBlockResult.B525SpectrumBlocks[b525SpectrumBlockIndex].B525SpectrumBlockConflict = b525SpectrumBlockConflicted[b525SpectrumBlockSelector]
		b525SpectrumBlockResult.B525SpectrumBlocks[b525SpectrumBlockIndex].B525SpectrumBlockContributed = !b525SpectrumBlockConflicted[b525SpectrumBlockSelector] && b525SpectrumBlockChosen[b525SpectrumBlockSelector] == b525SpectrumBlockIndex
	}
	for b525SpectrumBlockSelector := B525SpectrumBlockSelectorFirst; b525SpectrumBlockSelector <= B525SpectrumBlockSelectorLast; b525SpectrumBlockSelector++ {
		if _, b525SpectrumBlockHas := b525SpectrumBlockChosen[b525SpectrumBlockSelector]; !b525SpectrumBlockHas || b525SpectrumBlockConflicted[b525SpectrumBlockSelector] {
			b525SpectrumBlockResult.B525SpectrumBlockMissing = append(b525SpectrumBlockResult.B525SpectrumBlockMissing, b525SpectrumBlockSelector)
		}
	}
	switch {
	case len(b525SpectrumBlockResult.B525SpectrumBlockConflicts) > 0:
		b525SpectrumBlockResult.B525SpectrumBlockState = B525SpectrumBlockJoinConflict
	case b525SpectrumBlockValidCount == B525SpectrumBlockSelectorLast-B525SpectrumBlockSelectorFirst+1:
		b525SpectrumBlockResult.B525SpectrumBlockState = B525SpectrumBlockJoinComplete
		b525SpectrumBlockResult.B525SpectrumBlockComplete = true
	case len(b525SpectrumBlockResult.B525SpectrumBlockEntries) > 0:
		b525SpectrumBlockResult.B525SpectrumBlockState = B525SpectrumBlockJoinPartial
	}
	return b525SpectrumBlockResult
}
