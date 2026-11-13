package vaillant

import (
	"encoding/hex"
	"encoding/json"
	"os"
	"reflect"
	"testing"
)

type b525SpectrumBlockFixtureBlock struct {
	B525SpectrumBlockID       string `json:"id"`
	B525SpectrumBlockSelector int    `json:"selector"`
	B525SpectrumBlockDir      string `json:"direction"`
	B525SpectrumBlockValues   string `json:"values"`
	B525SpectrumBlockExpected struct {
		B525SpectrumBlockState     string  `json:"state"`
		B525SpectrumBlockFirst     int     `json:"first_index"`
		B525SpectrumBlockLast      int     `json:"last_index"`
		B525SpectrumBlockCount     int     `json:"count"`
		B525SpectrumBlockHead      [][]int `json:"head"`
		B525SpectrumBlockTail      []int   `json:"tail"`
		B525SpectrumBlockRemainder string  `json:"remainder"`
	} `json:"expected"`
}

type b525SpectrumBlockFixtureJoinBlock struct {
	B525SpectrumBlockSelector int    `json:"selector"`
	B525SpectrumBlockFamily   string `json:"family"`
	B525SpectrumBlockVariant  string `json:"variant"`
	B525SpectrumBlockOp       int    `json:"operation"`
	B525SpectrumBlockValues   string `json:"values"`
}

type b525SpectrumBlockFixtureJoin struct {
	B525SpectrumBlockID       string                              `json:"id"`
	B525SpectrumBlocks        []b525SpectrumBlockFixtureJoinBlock `json:"blocks"`
	B525SpectrumBlockExpected struct {
		B525SpectrumBlockState     string `json:"state"`
		B525SpectrumBlockComplete  bool   `json:"complete"`
		B525SpectrumBlockCount     int    `json:"count"`
		B525SpectrumBlockMissing   []int  `json:"missing"`
		B525SpectrumBlockConflicts []int  `json:"conflicts"`
	} `json:"expected"`
}

type b525SpectrumBlockFixture struct {
	B525SpectrumBlocks     []b525SpectrumBlockFixtureBlock `json:"blocks"`
	B525SpectrumBlockJoins []b525SpectrumBlockFixtureJoin  `json:"joins"`
}

func b525SpectrumBlockLoad(b525SpectrumBlockT *testing.T) b525SpectrumBlockFixture {
	b525SpectrumBlockT.Helper()
	b525SpectrumBlockData, b525SpectrumBlockErr := os.ReadFile("testdata/ab82_b525_spectrum_blocks.json")
	if b525SpectrumBlockErr != nil {
		b525SpectrumBlockT.Fatal(b525SpectrumBlockErr)
	}
	var b525SpectrumBlockValue b525SpectrumBlockFixture
	if b525SpectrumBlockErr := json.Unmarshal(b525SpectrumBlockData, &b525SpectrumBlockValue); b525SpectrumBlockErr != nil {
		b525SpectrumBlockT.Fatal(b525SpectrumBlockErr)
	}
	return b525SpectrumBlockValue
}

func b525SpectrumBlockHex(b525SpectrumBlockT *testing.T, b525SpectrumBlockText string) []byte {
	b525SpectrumBlockT.Helper()
	b525SpectrumBlockBytes, b525SpectrumBlockErr := hex.DecodeString(b525SpectrumBlockText)
	if b525SpectrumBlockErr != nil {
		b525SpectrumBlockT.Fatal(b525SpectrumBlockErr)
	}
	return b525SpectrumBlockBytes
}

func b525SpectrumBlockCtx(b525SpectrumBlockSelector int, b525SpectrumBlockDirection string) B525SpectrumBlockContext {
	return B525SpectrumBlockContext{B525SpectrumBlockPrimary: 0xB5, B525SpectrumBlockSecondary: 0x25, B525SpectrumBlockFamily: "fam-a", B525SpectrumBlockVariant: "var-a", B525SpectrumBlockOperation: 1, B525SpectrumBlockDirection: b525SpectrumBlockDirection, B525SpectrumBlockSelector: b525SpectrumBlockSelector}
}

func b525SpectrumBlockRepeat(b525SpectrumBlockByte byte, b525SpectrumBlockCount int) []byte {
	b525SpectrumBlockOut := make([]byte, b525SpectrumBlockCount)
	for b525SpectrumBlockIndex := range b525SpectrumBlockOut {
		b525SpectrumBlockOut[b525SpectrumBlockIndex] = b525SpectrumBlockByte
	}
	return b525SpectrumBlockOut
}

func b525SpectrumBlockFull(b525SpectrumBlockSelector int) B525SpectrumBlock {
	return B525SpectrumBlockView(B525SpectrumBlockInput{B525SpectrumBlockContext: b525SpectrumBlockCtx(b525SpectrumBlockSelector, B525SpectrumBlockDirectionRequest), B525SpectrumBlockValues: b525SpectrumBlockRepeat(byte(b525SpectrumBlockSelector), 32)})
}

func TestAb82FixtureBlocks(b525SpectrumBlockT *testing.T) {
	b525SpectrumBlockCases := b525SpectrumBlockLoad(b525SpectrumBlockT).B525SpectrumBlocks
	if len(b525SpectrumBlockCases) != 17 {
		b525SpectrumBlockT.Fatalf("fixture count %d", len(b525SpectrumBlockCases))
	}
	for _, b525SpectrumBlockCase := range b525SpectrumBlockCases {
		b525SpectrumBlockT.Run(b525SpectrumBlockCase.B525SpectrumBlockID, func(b525SpectrumBlockT *testing.T) {
			b525SpectrumBlockRaw := b525SpectrumBlockHex(b525SpectrumBlockT, b525SpectrumBlockCase.B525SpectrumBlockValues)
			b525SpectrumBlockGot := B525SpectrumBlockView(B525SpectrumBlockInput{B525SpectrumBlockContext: b525SpectrumBlockCtx(b525SpectrumBlockCase.B525SpectrumBlockSelector, b525SpectrumBlockCase.B525SpectrumBlockDir), B525SpectrumBlockValues: b525SpectrumBlockRaw})
			if string(b525SpectrumBlockGot.B525SpectrumBlockState) != b525SpectrumBlockCase.B525SpectrumBlockExpected.B525SpectrumBlockState || len(b525SpectrumBlockGot.B525SpectrumBlockEntries) != b525SpectrumBlockCase.B525SpectrumBlockExpected.B525SpectrumBlockCount {
				b525SpectrumBlockT.Fatalf("state %s entries %d", b525SpectrumBlockGot.B525SpectrumBlockState, len(b525SpectrumBlockGot.B525SpectrumBlockEntries))
			}
			if b525SpectrumBlockGot.B525SpectrumBlockGrade != "Hypothesis" || !reflect.DeepEqual(b525SpectrumBlockGot.B525SpectrumBlockRaw, b525SpectrumBlockRaw) {
				b525SpectrumBlockT.Fatal("grade or raw lost")
			}
			if hex.EncodeToString(b525SpectrumBlockGot.B525SpectrumBlockRemainder) != b525SpectrumBlockCase.B525SpectrumBlockExpected.B525SpectrumBlockRemainder {
				b525SpectrumBlockT.Fatalf("remainder %x", b525SpectrumBlockGot.B525SpectrumBlockRemainder)
			}
			if b525SpectrumBlockCase.B525SpectrumBlockExpected.B525SpectrumBlockState == "short" != b525SpectrumBlockGot.B525SpectrumBlockIncomplete || b525SpectrumBlockCase.B525SpectrumBlockExpected.B525SpectrumBlockState == "extra" != b525SpectrumBlockGot.B525SpectrumBlockExtra {
				b525SpectrumBlockT.Fatal("incomplete/extra flags")
			}
			if b525SpectrumBlockCase.B525SpectrumBlockExpected.B525SpectrumBlockCount == 0 {
				if b525SpectrumBlockGot.B525SpectrumBlockMask != [4]uint64{} || len(b525SpectrumBlockGot.B525SpectrumBlockEntries) != 0 {
					b525SpectrumBlockT.Fatal("values from non-value state")
				}
				return
			}
			b525SpectrumBlockExp := b525SpectrumBlockCase.B525SpectrumBlockExpected
			b525SpectrumBlockStart := (b525SpectrumBlockCase.B525SpectrumBlockSelector - 10) * 32
			for b525SpectrumBlockPosition, b525SpectrumBlockEntry := range b525SpectrumBlockGot.B525SpectrumBlockEntries {
				if b525SpectrumBlockEntry.B525SpectrumBlockIndex != b525SpectrumBlockStart+b525SpectrumBlockPosition {
					b525SpectrumBlockT.Fatalf("entry %d index %d", b525SpectrumBlockPosition, b525SpectrumBlockEntry.B525SpectrumBlockIndex)
				}
			}
			if b525SpectrumBlockExp.B525SpectrumBlockCount == 32 {
				if b525SpectrumBlockGot.B525SpectrumBlockEntries[0].B525SpectrumBlockIndex != b525SpectrumBlockExp.B525SpectrumBlockFirst || b525SpectrumBlockGot.B525SpectrumBlockEntries[31].B525SpectrumBlockIndex != b525SpectrumBlockExp.B525SpectrumBlockLast {
					b525SpectrumBlockT.Fatal("index range")
				}
			}
			for _, b525SpectrumBlockPair := range b525SpectrumBlockExp.B525SpectrumBlockHead {
				b525SpectrumBlockEntry := b525SpectrumBlockGot.B525SpectrumBlockEntries[b525SpectrumBlockPair[0]-b525SpectrumBlockExp.B525SpectrumBlockFirst]
				if b525SpectrumBlockEntry.B525SpectrumBlockIndex != b525SpectrumBlockPair[0] || int(b525SpectrumBlockEntry.B525SpectrumBlockSigned) != b525SpectrumBlockPair[1] {
					b525SpectrumBlockT.Fatalf("entry %+v want %v", b525SpectrumBlockEntry, b525SpectrumBlockPair)
				}
			}
			if len(b525SpectrumBlockExp.B525SpectrumBlockTail) == 2 && int(b525SpectrumBlockGot.B525SpectrumBlockEntries[31].B525SpectrumBlockSigned) != b525SpectrumBlockExp.B525SpectrumBlockTail[1] {
				b525SpectrumBlockT.Fatal("tail value")
			}
		})
	}
}

func TestAb82FixtureJoins(b525SpectrumBlockT *testing.T) {
	b525SpectrumBlockCases := b525SpectrumBlockLoad(b525SpectrumBlockT).B525SpectrumBlockJoins
	if len(b525SpectrumBlockCases) != 11 {
		b525SpectrumBlockT.Fatalf("fixture count %d", len(b525SpectrumBlockCases))
	}
	for _, b525SpectrumBlockCase := range b525SpectrumBlockCases {
		b525SpectrumBlockT.Run(b525SpectrumBlockCase.B525SpectrumBlockID, func(b525SpectrumBlockT *testing.T) {
			var b525SpectrumBlockSupplied []B525SpectrumBlock
			for _, b525SpectrumBlockItem := range b525SpectrumBlockCase.B525SpectrumBlocks {
				b525SpectrumBlockContext := b525SpectrumBlockCtx(b525SpectrumBlockItem.B525SpectrumBlockSelector, B525SpectrumBlockDirectionRequest)
				b525SpectrumBlockContext.B525SpectrumBlockFamily, b525SpectrumBlockContext.B525SpectrumBlockVariant, b525SpectrumBlockContext.B525SpectrumBlockOperation = b525SpectrumBlockItem.B525SpectrumBlockFamily, b525SpectrumBlockItem.B525SpectrumBlockVariant, b525SpectrumBlockItem.B525SpectrumBlockOp
				b525SpectrumBlockSupplied = append(b525SpectrumBlockSupplied, B525SpectrumBlockView(B525SpectrumBlockInput{B525SpectrumBlockContext: b525SpectrumBlockContext, B525SpectrumBlockValues: b525SpectrumBlockHex(b525SpectrumBlockT, b525SpectrumBlockItem.B525SpectrumBlockValues)}))
			}
			b525SpectrumBlockGot := B525SpectrumBlockJoin(b525SpectrumBlockSupplied)
			b525SpectrumBlockExp := b525SpectrumBlockCase.B525SpectrumBlockExpected
			if string(b525SpectrumBlockGot.B525SpectrumBlockState) != b525SpectrumBlockExp.B525SpectrumBlockState || b525SpectrumBlockGot.B525SpectrumBlockComplete != b525SpectrumBlockExp.B525SpectrumBlockComplete || len(b525SpectrumBlockGot.B525SpectrumBlockEntries) != b525SpectrumBlockExp.B525SpectrumBlockCount {
				b525SpectrumBlockT.Fatalf("state %s complete %v entries %d", b525SpectrumBlockGot.B525SpectrumBlockState, b525SpectrumBlockGot.B525SpectrumBlockComplete, len(b525SpectrumBlockGot.B525SpectrumBlockEntries))
			}
			if !reflect.DeepEqual(append([]int{}, b525SpectrumBlockGot.B525SpectrumBlockMissing...), append([]int{}, b525SpectrumBlockExp.B525SpectrumBlockMissing...)) || !reflect.DeepEqual(append([]int{}, b525SpectrumBlockGot.B525SpectrumBlockConflicts...), append([]int{}, b525SpectrumBlockExp.B525SpectrumBlockConflicts...)) {
				b525SpectrumBlockT.Fatalf("missing %v conflicts %v", b525SpectrumBlockGot.B525SpectrumBlockMissing, b525SpectrumBlockGot.B525SpectrumBlockConflicts)
			}
			for b525SpectrumBlockPosition, b525SpectrumBlockEntry := range b525SpectrumBlockGot.B525SpectrumBlockEntries {
				if b525SpectrumBlockPosition > 0 && b525SpectrumBlockEntry.B525SpectrumBlockIndex <= b525SpectrumBlockGot.B525SpectrumBlockEntries[b525SpectrumBlockPosition-1].B525SpectrumBlockIndex {
					b525SpectrumBlockT.Fatal("entries not in index order")
				}
				if !b525SpectrumBlockGot.B525SpectrumBlockPresent(b525SpectrumBlockEntry.B525SpectrumBlockIndex) {
					b525SpectrumBlockT.Fatal("mask lacks entry")
				}
			}
		})
	}
}

func TestAb82SelectorStartsAndSignedBytes(b525SpectrumBlockT *testing.T) {
	b525SpectrumBlockStarts := map[int]int{10: 0, 11: 32, 12: 64, 13: 96, 14: 128, 15: 160, 16: 192, 17: 224}
	for b525SpectrumBlockSelector, b525SpectrumBlockStart := range b525SpectrumBlockStarts {
		b525SpectrumBlockRaw := b525SpectrumBlockRepeat(0x01, 32)
		b525SpectrumBlockRaw[0], b525SpectrumBlockRaw[1], b525SpectrumBlockRaw[2], b525SpectrumBlockRaw[30], b525SpectrumBlockRaw[31] = 0x00, 0x7F, 0x80, 0xFF, 0x7F
		b525SpectrumBlockGot := B525SpectrumBlockView(B525SpectrumBlockInput{B525SpectrumBlockContext: b525SpectrumBlockCtx(b525SpectrumBlockSelector, B525SpectrumBlockDirectionRequest), B525SpectrumBlockValues: b525SpectrumBlockRaw})
		b525SpectrumBlockWant := []struct{ b525SpectrumBlockIndex, b525SpectrumBlockValue int }{{b525SpectrumBlockStart, 0}, {b525SpectrumBlockStart + 1, 127}, {b525SpectrumBlockStart + 2, -128}, {b525SpectrumBlockStart + 30, -1}, {b525SpectrumBlockStart + 31, 127}}
		for _, b525SpectrumBlockItem := range b525SpectrumBlockWant {
			b525SpectrumBlockOffset := b525SpectrumBlockItem.b525SpectrumBlockIndex - b525SpectrumBlockStart
			if b525SpectrumBlockGot.B525SpectrumBlockEntries[b525SpectrumBlockOffset].B525SpectrumBlockIndex != b525SpectrumBlockItem.b525SpectrumBlockIndex || int(b525SpectrumBlockGot.B525SpectrumBlockEntries[b525SpectrumBlockOffset].B525SpectrumBlockSigned) != b525SpectrumBlockItem.b525SpectrumBlockValue {
				b525SpectrumBlockT.Fatalf("selector %d entry %+v want %+v", b525SpectrumBlockSelector, b525SpectrumBlockGot.B525SpectrumBlockEntries[b525SpectrumBlockOffset], b525SpectrumBlockItem)
			}
		}
		if b525SpectrumBlockGot.B525SpectrumBlockStart != b525SpectrumBlockStart || b525SpectrumBlockGot.B525SpectrumBlockPresent(b525SpectrumBlockStart-1) || b525SpectrumBlockGot.B525SpectrumBlockPresent(b525SpectrumBlockStart+32) || !b525SpectrumBlockGot.B525SpectrumBlockPresent(b525SpectrumBlockStart) || !b525SpectrumBlockGot.B525SpectrumBlockPresent(b525SpectrumBlockStart+31) {
			b525SpectrumBlockT.Fatalf("selector %d mask", b525SpectrumBlockSelector)
		}
	}
	for b525SpectrumBlockByte := 0; b525SpectrumBlockByte < 256; b525SpectrumBlockByte++ {
		b525SpectrumBlockWant := b525SpectrumBlockByte
		if b525SpectrumBlockByte >= 128 {
			b525SpectrumBlockWant = b525SpectrumBlockByte - 256
		}
		b525SpectrumBlockGot := B525SpectrumBlockView(B525SpectrumBlockInput{B525SpectrumBlockContext: b525SpectrumBlockCtx(10, B525SpectrumBlockDirectionRequest), B525SpectrumBlockValues: []byte{byte(b525SpectrumBlockByte)}})
		if int(b525SpectrumBlockGot.B525SpectrumBlockEntries[0].B525SpectrumBlockSigned) != b525SpectrumBlockWant || b525SpectrumBlockGot.B525SpectrumBlockEntries[0].B525SpectrumBlockRawByte != byte(b525SpectrumBlockByte) {
			b525SpectrumBlockT.Fatalf("byte %d", b525SpectrumBlockByte)
		}
	}
}

func TestAb82ExtraDoesNotAdvance(b525SpectrumBlockT *testing.T) {
	b525SpectrumBlockRaw := append(b525SpectrumBlockRepeat(0x05, 32), 0xFF, 0x00)
	b525SpectrumBlockGot := B525SpectrumBlockView(B525SpectrumBlockInput{B525SpectrumBlockContext: b525SpectrumBlockCtx(10, B525SpectrumBlockDirectionRequest), B525SpectrumBlockValues: b525SpectrumBlockRaw})
	if b525SpectrumBlockGot.B525SpectrumBlockState != B525SpectrumBlockStateExtra || len(b525SpectrumBlockGot.B525SpectrumBlockEntries) != 32 || b525SpectrumBlockGot.B525SpectrumBlockPresent(32) || !reflect.DeepEqual(b525SpectrumBlockGot.B525SpectrumBlockRemainder, []byte{0xFF, 0x00}) {
		b525SpectrumBlockT.Fatal("extra advanced or remainder lost")
	}
}

func TestAb82IdentityAxes(b525SpectrumBlockT *testing.T) {
	b525SpectrumBlockRaw := b525SpectrumBlockRepeat(0x01, 32)
	if B525SpectrumBlockView(B525SpectrumBlockInput{B525SpectrumBlockContext: b525SpectrumBlockCtx(10, B525SpectrumBlockDirectionRequest), B525SpectrumBlockValues: b525SpectrumBlockRaw}).B525SpectrumBlockState != B525SpectrumBlockStateValid {
		b525SpectrumBlockT.Fatal("valid control")
	}
	b525SpectrumBlockMutators := map[string]func(*B525SpectrumBlockContext){
		"primary": func(b525SpectrumBlockC *B525SpectrumBlockContext) { b525SpectrumBlockC.B525SpectrumBlockPrimary = 0xB6 },
		"secondary": func(b525SpectrumBlockC *B525SpectrumBlockContext) {
			b525SpectrumBlockC.B525SpectrumBlockSecondary = 0x26
		},
		"family":        func(b525SpectrumBlockC *B525SpectrumBlockContext) { b525SpectrumBlockC.B525SpectrumBlockFamily = "" },
		"variant":       func(b525SpectrumBlockC *B525SpectrumBlockContext) { b525SpectrumBlockC.B525SpectrumBlockVariant = "" },
		"operation-low": func(b525SpectrumBlockC *B525SpectrumBlockContext) { b525SpectrumBlockC.B525SpectrumBlockOperation = -1 },
		"operation-high": func(b525SpectrumBlockC *B525SpectrumBlockContext) {
			b525SpectrumBlockC.B525SpectrumBlockOperation = 256
		},
		"direction": func(b525SpectrumBlockC *B525SpectrumBlockContext) {
			b525SpectrumBlockC.B525SpectrumBlockDirection = "response"
		},
		"selector-low":  func(b525SpectrumBlockC *B525SpectrumBlockContext) { b525SpectrumBlockC.B525SpectrumBlockSelector = 9 },
		"selector-high": func(b525SpectrumBlockC *B525SpectrumBlockContext) { b525SpectrumBlockC.B525SpectrumBlockSelector = 18 },
	}
	for b525SpectrumBlockName, b525SpectrumBlockMutate := range b525SpectrumBlockMutators {
		b525SpectrumBlockContext := b525SpectrumBlockCtx(10, B525SpectrumBlockDirectionRequest)
		b525SpectrumBlockMutate(&b525SpectrumBlockContext)
		b525SpectrumBlockGot := B525SpectrumBlockView(B525SpectrumBlockInput{B525SpectrumBlockContext: b525SpectrumBlockContext, B525SpectrumBlockValues: b525SpectrumBlockRaw})
		if b525SpectrumBlockGot.B525SpectrumBlockState != B525SpectrumBlockStateRawOnly || len(b525SpectrumBlockGot.B525SpectrumBlockEntries) != 0 || !reflect.DeepEqual(b525SpectrumBlockGot.B525SpectrumBlockRaw, b525SpectrumBlockRaw) || b525SpectrumBlockGot.B525SpectrumBlockContext != b525SpectrumBlockContext {
			b525SpectrumBlockT.Fatalf("%s not raw-only with preserved raw", b525SpectrumBlockName)
		}
	}
	for _, b525SpectrumBlockOperation := range []int{0, 1, 255} {
		b525SpectrumBlockContext := b525SpectrumBlockCtx(10, B525SpectrumBlockDirectionRequest)
		b525SpectrumBlockContext.B525SpectrumBlockOperation = b525SpectrumBlockOperation
		if B525SpectrumBlockView(B525SpectrumBlockInput{B525SpectrumBlockContext: b525SpectrumBlockContext, B525SpectrumBlockValues: b525SpectrumBlockRaw}).B525SpectrumBlockState != B525SpectrumBlockStateValid {
			b525SpectrumBlockT.Fatalf("operation %d", b525SpectrumBlockOperation)
		}
	}
}

func TestAb82AbsentDiffersFromZero(b525SpectrumBlockT *testing.T) {
	b525SpectrumBlockZero := B525SpectrumBlockView(B525SpectrumBlockInput{B525SpectrumBlockContext: b525SpectrumBlockCtx(10, B525SpectrumBlockDirectionRequest), B525SpectrumBlockValues: []byte{0x00}})
	b525SpectrumBlockEmpty := B525SpectrumBlockView(B525SpectrumBlockInput{B525SpectrumBlockContext: b525SpectrumBlockCtx(10, B525SpectrumBlockDirectionRequest), B525SpectrumBlockValues: nil})
	if !b525SpectrumBlockZero.B525SpectrumBlockPresent(0) || b525SpectrumBlockZero.B525SpectrumBlockEntries[0].B525SpectrumBlockSigned != 0 || b525SpectrumBlockEmpty.B525SpectrumBlockPresent(0) || len(b525SpectrumBlockEmpty.B525SpectrumBlockEntries) != 0 {
		b525SpectrumBlockT.Fatal("explicit zero and absent conflated")
	}
	b525SpectrumBlockJoined := B525SpectrumBlockJoin([]B525SpectrumBlock{b525SpectrumBlockZero, b525SpectrumBlockFull(11)})
	if !b525SpectrumBlockJoined.B525SpectrumBlockPresent(0) || b525SpectrumBlockJoined.B525SpectrumBlockPresent(1) || b525SpectrumBlockJoined.B525SpectrumBlockPresent(31) || !b525SpectrumBlockJoined.B525SpectrumBlockPresent(32) {
		b525SpectrumBlockT.Fatal("short block was zero-filled")
	}
	b525SpectrumBlockMissing := B525SpectrumBlockJoin([]B525SpectrumBlock{b525SpectrumBlockFull(10), b525SpectrumBlockFull(12)})
	for b525SpectrumBlockIndex := 32; b525SpectrumBlockIndex < 64; b525SpectrumBlockIndex++ {
		if b525SpectrumBlockMissing.B525SpectrumBlockPresent(b525SpectrumBlockIndex) {
			b525SpectrumBlockT.Fatal("missing block zero-filled")
		}
	}
}

func TestAb82PrefixOnlyNeverValues(b525SpectrumBlockT *testing.T) {
	for b525SpectrumBlockSelector := 10; b525SpectrumBlockSelector <= 17; b525SpectrumBlockSelector++ {
		b525SpectrumBlockGot := B525SpectrumBlockView(B525SpectrumBlockInput{B525SpectrumBlockContext: b525SpectrumBlockCtx(b525SpectrumBlockSelector, B525SpectrumBlockDirectionPositive), B525SpectrumBlockValues: b525SpectrumBlockRepeat(0x40, 32)})
		if b525SpectrumBlockGot.B525SpectrumBlockState != B525SpectrumBlockStatePrefixOnly || len(b525SpectrumBlockGot.B525SpectrumBlockEntries) != 0 || b525SpectrumBlockGot.B525SpectrumBlockMask != [4]uint64{} || len(b525SpectrumBlockGot.B525SpectrumBlockRaw) != 32 {
			b525SpectrumBlockT.Fatalf("selector %d", b525SpectrumBlockSelector)
		}
	}
	b525SpectrumBlockSupplied := []B525SpectrumBlock{}
	for b525SpectrumBlockSelector := 10; b525SpectrumBlockSelector <= 17; b525SpectrumBlockSelector++ {
		b525SpectrumBlockSupplied = append(b525SpectrumBlockSupplied, B525SpectrumBlockView(B525SpectrumBlockInput{B525SpectrumBlockContext: b525SpectrumBlockCtx(b525SpectrumBlockSelector, B525SpectrumBlockDirectionPositive), B525SpectrumBlockValues: b525SpectrumBlockRepeat(0x40, 32)}))
	}
	b525SpectrumBlockJoined := B525SpectrumBlockJoin(b525SpectrumBlockSupplied)
	if b525SpectrumBlockJoined.B525SpectrumBlockState != B525SpectrumBlockJoinEmpty || b525SpectrumBlockJoined.B525SpectrumBlockComplete || len(b525SpectrumBlockJoined.B525SpectrumBlockEntries) != 0 || len(b525SpectrumBlockJoined.B525SpectrumBlocks) != 8 {
		b525SpectrumBlockT.Fatal("positive blocks joined into values")
	}
}

func TestAb82JoinArrivalOrderAndForgery(b525SpectrumBlockT *testing.T) {
	b525SpectrumBlockForward := []B525SpectrumBlock{}
	for b525SpectrumBlockSelector := 10; b525SpectrumBlockSelector <= 17; b525SpectrumBlockSelector++ {
		b525SpectrumBlockForward = append(b525SpectrumBlockForward, b525SpectrumBlockFull(b525SpectrumBlockSelector))
	}
	b525SpectrumBlockReverse := []B525SpectrumBlock{}
	for b525SpectrumBlockIndex := len(b525SpectrumBlockForward) - 1; b525SpectrumBlockIndex >= 0; b525SpectrumBlockIndex-- {
		b525SpectrumBlockReverse = append(b525SpectrumBlockReverse, b525SpectrumBlockForward[b525SpectrumBlockIndex])
	}
	b525SpectrumBlockA, b525SpectrumBlockB := B525SpectrumBlockJoin(b525SpectrumBlockForward), B525SpectrumBlockJoin(b525SpectrumBlockReverse)
	if !reflect.DeepEqual(b525SpectrumBlockA.B525SpectrumBlockEntries, b525SpectrumBlockB.B525SpectrumBlockEntries) || b525SpectrumBlockA.B525SpectrumBlockMask != b525SpectrumBlockB.B525SpectrumBlockMask {
		b525SpectrumBlockT.Fatal("arrival order changed entries")
	}
	for b525SpectrumBlockIndex, b525SpectrumBlockEntry := range b525SpectrumBlockA.B525SpectrumBlockEntries {
		if b525SpectrumBlockEntry.B525SpectrumBlockIndex != b525SpectrumBlockIndex || b525SpectrumBlockEntry.B525SpectrumBlockRawByte != byte(10+b525SpectrumBlockIndex/32) {
			b525SpectrumBlockT.Fatalf("entry %d %+v", b525SpectrumBlockIndex, b525SpectrumBlockEntry)
		}
	}
	b525SpectrumBlockForged := b525SpectrumBlockFull(17)
	b525SpectrumBlockForged.B525SpectrumBlockRaw = b525SpectrumBlockForged.B525SpectrumBlockRaw[:5]
	b525SpectrumBlockSeven := B525SpectrumBlockJoin(append(append([]B525SpectrumBlock{}, b525SpectrumBlockForward[:7]...), b525SpectrumBlockForged))
	if b525SpectrumBlockSeven.B525SpectrumBlockComplete || len(b525SpectrumBlockSeven.B525SpectrumBlockEntries) != 229 {
		b525SpectrumBlockT.Fatal("forged state trusted")
	}
	b525SpectrumBlockBlank := B525SpectrumBlock{}
	b525SpectrumBlockMixed := B525SpectrumBlockJoin([]B525SpectrumBlock{b525SpectrumBlockFull(10), b525SpectrumBlockBlank})
	if b525SpectrumBlockMixed.B525SpectrumBlockState != B525SpectrumBlockJoinPartial || b525SpectrumBlockMixed.B525SpectrumBlocks[1].B525SpectrumBlockState != B525SpectrumBlockStateRawOnly || b525SpectrumBlockMixed.B525SpectrumBlocks[1].B525SpectrumBlockContributed {
		b525SpectrumBlockT.Fatal("raw-only block disturbed join")
	}
}

func TestAb82JoinPerBlockValidityAndConflict(b525SpectrumBlockT *testing.T) {
	b525SpectrumBlockShort := B525SpectrumBlockView(B525SpectrumBlockInput{B525SpectrumBlockContext: b525SpectrumBlockCtx(11, B525SpectrumBlockDirectionRequest), B525SpectrumBlockValues: []byte{1, 2}})
	b525SpectrumBlockOther := b525SpectrumBlockFull(11)
	b525SpectrumBlockJoined := B525SpectrumBlockJoin([]B525SpectrumBlock{b525SpectrumBlockFull(10), b525SpectrumBlockShort, b525SpectrumBlockOther})
	if b525SpectrumBlockJoined.B525SpectrumBlockState != B525SpectrumBlockJoinConflict || !reflect.DeepEqual(b525SpectrumBlockJoined.B525SpectrumBlockConflicts, []int{11}) || len(b525SpectrumBlockJoined.B525SpectrumBlockEntries) != 32 {
		b525SpectrumBlockT.Fatal("short vs full duplicate not conflicting")
	}
	if !b525SpectrumBlockJoined.B525SpectrumBlocks[1].B525SpectrumBlockConflict || !b525SpectrumBlockJoined.B525SpectrumBlocks[2].B525SpectrumBlockConflict || b525SpectrumBlockJoined.B525SpectrumBlocks[1].B525SpectrumBlockState != B525SpectrumBlockStateShort || !b525SpectrumBlockJoined.B525SpectrumBlocks[0].B525SpectrumBlockContributed || b525SpectrumBlockJoined.B525SpectrumBlocks[0].B525SpectrumBlockConflict {
		b525SpectrumBlockT.Fatal("per-block validity")
	}
	b525SpectrumBlockCross := b525SpectrumBlockFull(12)
	b525SpectrumBlockCross.B525SpectrumBlockContext.B525SpectrumBlockFamily = "fam-b"
	if B525SpectrumBlockJoin([]B525SpectrumBlock{b525SpectrumBlockFull(10), b525SpectrumBlockCross}).B525SpectrumBlockState != B525SpectrumBlockJoinContextMismatch {
		b525SpectrumBlockT.Fatal("cross family")
	}
	b525SpectrumBlockCross = b525SpectrumBlockFull(12)
	b525SpectrumBlockCross.B525SpectrumBlockContext.B525SpectrumBlockPrimary = 0xB6
	if B525SpectrumBlockJoin([]B525SpectrumBlock{b525SpectrumBlockFull(10), b525SpectrumBlockCross}).B525SpectrumBlockState != B525SpectrumBlockJoinPartial {
		b525SpectrumBlockT.Fatal("unknown identity block must stay raw-only, not mismatch")
	}
}

func TestAb82OwnedMemory(b525SpectrumBlockT *testing.T) {
	b525SpectrumBlockRaw := append(b525SpectrumBlockRepeat(0x01, 32), 0x09)
	b525SpectrumBlockGot := B525SpectrumBlockView(B525SpectrumBlockInput{B525SpectrumBlockContext: b525SpectrumBlockCtx(10, B525SpectrumBlockDirectionRequest), B525SpectrumBlockValues: b525SpectrumBlockRaw})
	b525SpectrumBlockRaw[0], b525SpectrumBlockRaw[32] = 0x77, 0x77
	if b525SpectrumBlockGot.B525SpectrumBlockRaw[0] != 0x01 || b525SpectrumBlockGot.B525SpectrumBlockEntries[0].B525SpectrumBlockRawByte != 0x01 || b525SpectrumBlockGot.B525SpectrumBlockRemainder[0] != 0x09 {
		b525SpectrumBlockT.Fatal("input aliased")
	}
	b525SpectrumBlockClone := b525SpectrumBlockGot.B525SpectrumBlockClone()
	b525SpectrumBlockClone.B525SpectrumBlockRaw[0], b525SpectrumBlockClone.B525SpectrumBlockRemainder[0], b525SpectrumBlockClone.B525SpectrumBlockEntries[0].B525SpectrumBlockRawByte = 0x55, 0x55, 0x55
	if b525SpectrumBlockGot.B525SpectrumBlockRaw[0] != 0x01 || b525SpectrumBlockGot.B525SpectrumBlockRemainder[0] != 0x09 || b525SpectrumBlockGot.B525SpectrumBlockEntries[0].B525SpectrumBlockRawByte != 0x01 {
		b525SpectrumBlockT.Fatal("clone aliased")
	}
	b525SpectrumBlockSource := b525SpectrumBlockFull(10)
	b525SpectrumBlockJoined := B525SpectrumBlockJoin([]B525SpectrumBlock{b525SpectrumBlockSource})
	b525SpectrumBlockSource.B525SpectrumBlockRaw[0] = 0x66
	b525SpectrumBlockJoined.B525SpectrumBlockEntries[1].B525SpectrumBlockRawByte = 0x66
	if b525SpectrumBlockJoined.B525SpectrumBlockEntries[0].B525SpectrumBlockRawByte != 10 || B525SpectrumBlockJoin([]B525SpectrumBlock{b525SpectrumBlockFull(10)}).B525SpectrumBlockEntries[1].B525SpectrumBlockRawByte != 10 {
		b525SpectrumBlockT.Fatal("join shares memory")
	}
	if B525SpectrumBlockView(B525SpectrumBlockInput{B525SpectrumBlockContext: b525SpectrumBlockCtx(10, B525SpectrumBlockDirectionRequest)}).B525SpectrumBlockRaw != nil {
		b525SpectrumBlockT.Fatal("nil raw invented")
	}
}
