package vaillant

import (
	"encoding/hex"
	"encoding/json"
	"os"
	"reflect"
	"testing"
)

type ab82FixtureBlock struct {
	Ab82ID       string `json:"id"`
	Ab82Selector int    `json:"selector"`
	Ab82Dir      string `json:"direction"`
	Ab82Values   string `json:"values"`
	Ab82Expected struct {
		Ab82State     string  `json:"state"`
		Ab82First     int     `json:"first_index"`
		Ab82Last      int     `json:"last_index"`
		Ab82Count     int     `json:"count"`
		Ab82Head      [][]int `json:"head"`
		Ab82Tail      []int   `json:"tail"`
		Ab82Remainder string  `json:"remainder"`
	} `json:"expected"`
}

type ab82FixtureJoinBlock struct {
	Ab82Selector int    `json:"selector"`
	Ab82Family   string `json:"family"`
	Ab82Variant  string `json:"variant"`
	Ab82Op       int    `json:"operation"`
	Ab82Values   string `json:"values"`
}

type ab82FixtureJoin struct {
	Ab82ID       string                 `json:"id"`
	Ab82Blocks   []ab82FixtureJoinBlock `json:"blocks"`
	Ab82Expected struct {
		Ab82State     string `json:"state"`
		Ab82Complete  bool   `json:"complete"`
		Ab82Count     int    `json:"count"`
		Ab82Missing   []int  `json:"missing"`
		Ab82Conflicts []int  `json:"conflicts"`
	} `json:"expected"`
}

type ab82Fixture struct {
	Ab82Blocks []ab82FixtureBlock `json:"blocks"`
	Ab82Joins  []ab82FixtureJoin  `json:"joins"`
}

func ab82Load(ab82T *testing.T) ab82Fixture {
	ab82T.Helper()
	ab82Data, ab82Err := os.ReadFile("testdata/ab82_b525_spectrum_blocks.json")
	if ab82Err != nil {
		ab82T.Fatal(ab82Err)
	}
	var ab82Value ab82Fixture
	if ab82Err := json.Unmarshal(ab82Data, &ab82Value); ab82Err != nil {
		ab82T.Fatal(ab82Err)
	}
	return ab82Value
}

func ab82Hex(ab82T *testing.T, ab82Text string) []byte {
	ab82T.Helper()
	ab82Bytes, ab82Err := hex.DecodeString(ab82Text)
	if ab82Err != nil {
		ab82T.Fatal(ab82Err)
	}
	return ab82Bytes
}

func ab82Ctx(ab82Selector int, ab82Direction string) Ab82Context {
	return Ab82Context{Ab82Primary: 0xB5, Ab82Secondary: 0x25, Ab82Family: "fam-a", Ab82Variant: "var-a", Ab82Operation: 1, Ab82Direction: ab82Direction, Ab82Selector: ab82Selector}
}

func ab82Repeat(ab82Byte byte, ab82Count int) []byte {
	ab82Out := make([]byte, ab82Count)
	for ab82Index := range ab82Out {
		ab82Out[ab82Index] = ab82Byte
	}
	return ab82Out
}

func ab82Full(ab82Selector int) Ab82Block {
	return Ab82View(Ab82Input{Ab82Context: ab82Ctx(ab82Selector, Ab82DirectionRequest), Ab82Values: ab82Repeat(byte(ab82Selector), 32)})
}

func TestAb82FixtureBlocks(ab82T *testing.T) {
	ab82Cases := ab82Load(ab82T).Ab82Blocks
	if len(ab82Cases) != 17 {
		ab82T.Fatalf("fixture count %d", len(ab82Cases))
	}
	for _, ab82Case := range ab82Cases {
		ab82T.Run(ab82Case.Ab82ID, func(ab82T *testing.T) {
			ab82Raw := ab82Hex(ab82T, ab82Case.Ab82Values)
			ab82Got := Ab82View(Ab82Input{Ab82Context: ab82Ctx(ab82Case.Ab82Selector, ab82Case.Ab82Dir), Ab82Values: ab82Raw})
			if string(ab82Got.Ab82State) != ab82Case.Ab82Expected.Ab82State || len(ab82Got.Ab82Entries) != ab82Case.Ab82Expected.Ab82Count {
				ab82T.Fatalf("state %s entries %d", ab82Got.Ab82State, len(ab82Got.Ab82Entries))
			}
			if ab82Got.Ab82Grade != "Hypothesis" || !reflect.DeepEqual(ab82Got.Ab82Raw, ab82Raw) {
				ab82T.Fatal("grade or raw lost")
			}
			if hex.EncodeToString(ab82Got.Ab82Remainder) != ab82Case.Ab82Expected.Ab82Remainder {
				ab82T.Fatalf("remainder %x", ab82Got.Ab82Remainder)
			}
			if ab82Case.Ab82Expected.Ab82State == "short" != ab82Got.Ab82Incomplete || ab82Case.Ab82Expected.Ab82State == "extra" != ab82Got.Ab82Extra {
				ab82T.Fatal("incomplete/extra flags")
			}
			if ab82Case.Ab82Expected.Ab82Count == 0 {
				if ab82Got.Ab82Mask != [4]uint64{} || len(ab82Got.Ab82Entries) != 0 {
					ab82T.Fatal("values from non-value state")
				}
				return
			}
			ab82Exp := ab82Case.Ab82Expected
			ab82Start := (ab82Case.Ab82Selector - 10) * 32
			for ab82Position, ab82Entry := range ab82Got.Ab82Entries {
				if ab82Entry.Ab82Index != ab82Start+ab82Position {
					ab82T.Fatalf("entry %d index %d", ab82Position, ab82Entry.Ab82Index)
				}
			}
			if ab82Exp.Ab82Count == 32 {
				if ab82Got.Ab82Entries[0].Ab82Index != ab82Exp.Ab82First || ab82Got.Ab82Entries[31].Ab82Index != ab82Exp.Ab82Last {
					ab82T.Fatal("index range")
				}
			}
			for _, ab82Pair := range ab82Exp.Ab82Head {
				ab82Entry := ab82Got.Ab82Entries[ab82Pair[0]-ab82Exp.Ab82First]
				if ab82Entry.Ab82Index != ab82Pair[0] || int(ab82Entry.Ab82Signed) != ab82Pair[1] {
					ab82T.Fatalf("entry %+v want %v", ab82Entry, ab82Pair)
				}
			}
			if len(ab82Exp.Ab82Tail) == 2 && int(ab82Got.Ab82Entries[31].Ab82Signed) != ab82Exp.Ab82Tail[1] {
				ab82T.Fatal("tail value")
			}
		})
	}
}

func TestAb82FixtureJoins(ab82T *testing.T) {
	ab82Cases := ab82Load(ab82T).Ab82Joins
	if len(ab82Cases) != 11 {
		ab82T.Fatalf("fixture count %d", len(ab82Cases))
	}
	for _, ab82Case := range ab82Cases {
		ab82T.Run(ab82Case.Ab82ID, func(ab82T *testing.T) {
			var ab82Supplied []Ab82Block
			for _, ab82Item := range ab82Case.Ab82Blocks {
				ab82Context := ab82Ctx(ab82Item.Ab82Selector, Ab82DirectionRequest)
				ab82Context.Ab82Family, ab82Context.Ab82Variant, ab82Context.Ab82Operation = ab82Item.Ab82Family, ab82Item.Ab82Variant, ab82Item.Ab82Op
				ab82Supplied = append(ab82Supplied, Ab82View(Ab82Input{Ab82Context: ab82Context, Ab82Values: ab82Hex(ab82T, ab82Item.Ab82Values)}))
			}
			ab82Got := Ab82Join(ab82Supplied)
			ab82Exp := ab82Case.Ab82Expected
			if string(ab82Got.Ab82State) != ab82Exp.Ab82State || ab82Got.Ab82Complete != ab82Exp.Ab82Complete || len(ab82Got.Ab82Entries) != ab82Exp.Ab82Count {
				ab82T.Fatalf("state %s complete %v entries %d", ab82Got.Ab82State, ab82Got.Ab82Complete, len(ab82Got.Ab82Entries))
			}
			if !reflect.DeepEqual(append([]int{}, ab82Got.Ab82Missing...), append([]int{}, ab82Exp.Ab82Missing...)) || !reflect.DeepEqual(append([]int{}, ab82Got.Ab82Conflicts...), append([]int{}, ab82Exp.Ab82Conflicts...)) {
				ab82T.Fatalf("missing %v conflicts %v", ab82Got.Ab82Missing, ab82Got.Ab82Conflicts)
			}
			for ab82Position, ab82Entry := range ab82Got.Ab82Entries {
				if ab82Position > 0 && ab82Entry.Ab82Index <= ab82Got.Ab82Entries[ab82Position-1].Ab82Index {
					ab82T.Fatal("entries not in index order")
				}
				if !ab82Got.Ab82Present(ab82Entry.Ab82Index) {
					ab82T.Fatal("mask lacks entry")
				}
			}
		})
	}
}

func TestAb82SelectorStartsAndSignedBytes(ab82T *testing.T) {
	ab82Starts := map[int]int{10: 0, 11: 32, 12: 64, 13: 96, 14: 128, 15: 160, 16: 192, 17: 224}
	for ab82Selector, ab82Start := range ab82Starts {
		ab82Raw := ab82Repeat(0x01, 32)
		ab82Raw[0], ab82Raw[1], ab82Raw[2], ab82Raw[30], ab82Raw[31] = 0x00, 0x7F, 0x80, 0xFF, 0x7F
		ab82Got := Ab82View(Ab82Input{Ab82Context: ab82Ctx(ab82Selector, Ab82DirectionRequest), Ab82Values: ab82Raw})
		ab82Want := []struct{ ab82Index, ab82Value int }{{ab82Start, 0}, {ab82Start + 1, 127}, {ab82Start + 2, -128}, {ab82Start + 30, -1}, {ab82Start + 31, 127}}
		for _, ab82Item := range ab82Want {
			ab82Offset := ab82Item.ab82Index - ab82Start
			if ab82Got.Ab82Entries[ab82Offset].Ab82Index != ab82Item.ab82Index || int(ab82Got.Ab82Entries[ab82Offset].Ab82Signed) != ab82Item.ab82Value {
				ab82T.Fatalf("selector %d entry %+v want %+v", ab82Selector, ab82Got.Ab82Entries[ab82Offset], ab82Item)
			}
		}
		if ab82Got.Ab82Start != ab82Start || ab82Got.Ab82Present(ab82Start-1) || ab82Got.Ab82Present(ab82Start+32) || !ab82Got.Ab82Present(ab82Start) || !ab82Got.Ab82Present(ab82Start+31) {
			ab82T.Fatalf("selector %d mask", ab82Selector)
		}
	}
	for ab82Byte := 0; ab82Byte < 256; ab82Byte++ {
		ab82Want := ab82Byte
		if ab82Byte >= 128 {
			ab82Want = ab82Byte - 256
		}
		ab82Got := Ab82View(Ab82Input{Ab82Context: ab82Ctx(10, Ab82DirectionRequest), Ab82Values: []byte{byte(ab82Byte)}})
		if int(ab82Got.Ab82Entries[0].Ab82Signed) != ab82Want || ab82Got.Ab82Entries[0].Ab82RawByte != byte(ab82Byte) {
			ab82T.Fatalf("byte %d", ab82Byte)
		}
	}
}

func TestAb82ExtraDoesNotAdvance(ab82T *testing.T) {
	ab82Raw := append(ab82Repeat(0x05, 32), 0xFF, 0x00)
	ab82Got := Ab82View(Ab82Input{Ab82Context: ab82Ctx(10, Ab82DirectionRequest), Ab82Values: ab82Raw})
	if ab82Got.Ab82State != Ab82StateExtra || len(ab82Got.Ab82Entries) != 32 || ab82Got.Ab82Present(32) || !reflect.DeepEqual(ab82Got.Ab82Remainder, []byte{0xFF, 0x00}) {
		ab82T.Fatal("extra advanced or remainder lost")
	}
}

func TestAb82IdentityAxes(ab82T *testing.T) {
	ab82Raw := ab82Repeat(0x01, 32)
	if Ab82View(Ab82Input{Ab82Context: ab82Ctx(10, Ab82DirectionRequest), Ab82Values: ab82Raw}).Ab82State != Ab82StateValid {
		ab82T.Fatal("valid control")
	}
	ab82Mutators := map[string]func(*Ab82Context){
		"primary":        func(ab82C *Ab82Context) { ab82C.Ab82Primary = 0xB6 },
		"secondary":      func(ab82C *Ab82Context) { ab82C.Ab82Secondary = 0x26 },
		"family":         func(ab82C *Ab82Context) { ab82C.Ab82Family = "" },
		"variant":        func(ab82C *Ab82Context) { ab82C.Ab82Variant = "" },
		"operation-low":  func(ab82C *Ab82Context) { ab82C.Ab82Operation = -1 },
		"operation-high": func(ab82C *Ab82Context) { ab82C.Ab82Operation = 256 },
		"direction":      func(ab82C *Ab82Context) { ab82C.Ab82Direction = "response" },
		"selector-low":   func(ab82C *Ab82Context) { ab82C.Ab82Selector = 9 },
		"selector-high":  func(ab82C *Ab82Context) { ab82C.Ab82Selector = 18 },
	}
	for ab82Name, ab82Mutate := range ab82Mutators {
		ab82Context := ab82Ctx(10, Ab82DirectionRequest)
		ab82Mutate(&ab82Context)
		ab82Got := Ab82View(Ab82Input{Ab82Context: ab82Context, Ab82Values: ab82Raw})
		if ab82Got.Ab82State != Ab82StateRawOnly || len(ab82Got.Ab82Entries) != 0 || !reflect.DeepEqual(ab82Got.Ab82Raw, ab82Raw) || ab82Got.Ab82Context != ab82Context {
			ab82T.Fatalf("%s not raw-only with preserved raw", ab82Name)
		}
	}
	for _, ab82Operation := range []int{0, 1, 255} {
		ab82Context := ab82Ctx(10, Ab82DirectionRequest)
		ab82Context.Ab82Operation = ab82Operation
		if Ab82View(Ab82Input{Ab82Context: ab82Context, Ab82Values: ab82Raw}).Ab82State != Ab82StateValid {
			ab82T.Fatalf("operation %d", ab82Operation)
		}
	}
}

func TestAb82AbsentDiffersFromZero(ab82T *testing.T) {
	ab82Zero := Ab82View(Ab82Input{Ab82Context: ab82Ctx(10, Ab82DirectionRequest), Ab82Values: []byte{0x00}})
	ab82Empty := Ab82View(Ab82Input{Ab82Context: ab82Ctx(10, Ab82DirectionRequest), Ab82Values: nil})
	if !ab82Zero.Ab82Present(0) || ab82Zero.Ab82Entries[0].Ab82Signed != 0 || ab82Empty.Ab82Present(0) || len(ab82Empty.Ab82Entries) != 0 {
		ab82T.Fatal("explicit zero and absent conflated")
	}
	ab82Joined := Ab82Join([]Ab82Block{ab82Zero, ab82Full(11)})
	if !ab82Joined.Ab82Present(0) || ab82Joined.Ab82Present(1) || ab82Joined.Ab82Present(31) || !ab82Joined.Ab82Present(32) {
		ab82T.Fatal("short block was zero-filled")
	}
	ab82Missing := Ab82Join([]Ab82Block{ab82Full(10), ab82Full(12)})
	for ab82Index := 32; ab82Index < 64; ab82Index++ {
		if ab82Missing.Ab82Present(ab82Index) {
			ab82T.Fatal("missing block zero-filled")
		}
	}
}

func TestAb82PrefixOnlyNeverValues(ab82T *testing.T) {
	for ab82Selector := 10; ab82Selector <= 17; ab82Selector++ {
		ab82Got := Ab82View(Ab82Input{Ab82Context: ab82Ctx(ab82Selector, Ab82DirectionPositive), Ab82Values: ab82Repeat(0x40, 32)})
		if ab82Got.Ab82State != Ab82StatePrefixOnly || len(ab82Got.Ab82Entries) != 0 || ab82Got.Ab82Mask != [4]uint64{} || len(ab82Got.Ab82Raw) != 32 {
			ab82T.Fatalf("selector %d", ab82Selector)
		}
	}
	ab82Supplied := []Ab82Block{}
	for ab82Selector := 10; ab82Selector <= 17; ab82Selector++ {
		ab82Supplied = append(ab82Supplied, Ab82View(Ab82Input{Ab82Context: ab82Ctx(ab82Selector, Ab82DirectionPositive), Ab82Values: ab82Repeat(0x40, 32)}))
	}
	ab82Joined := Ab82Join(ab82Supplied)
	if ab82Joined.Ab82State != Ab82JoinEmpty || ab82Joined.Ab82Complete || len(ab82Joined.Ab82Entries) != 0 || len(ab82Joined.Ab82Blocks) != 8 {
		ab82T.Fatal("positive blocks joined into values")
	}
}

func TestAb82JoinArrivalOrderAndForgery(ab82T *testing.T) {
	ab82Forward := []Ab82Block{}
	for ab82Selector := 10; ab82Selector <= 17; ab82Selector++ {
		ab82Forward = append(ab82Forward, ab82Full(ab82Selector))
	}
	ab82Reverse := []Ab82Block{}
	for ab82Index := len(ab82Forward) - 1; ab82Index >= 0; ab82Index-- {
		ab82Reverse = append(ab82Reverse, ab82Forward[ab82Index])
	}
	ab82A, ab82B := Ab82Join(ab82Forward), Ab82Join(ab82Reverse)
	if !reflect.DeepEqual(ab82A.Ab82Entries, ab82B.Ab82Entries) || ab82A.Ab82Mask != ab82B.Ab82Mask {
		ab82T.Fatal("arrival order changed entries")
	}
	for ab82Index, ab82Entry := range ab82A.Ab82Entries {
		if ab82Entry.Ab82Index != ab82Index || ab82Entry.Ab82RawByte != byte(10+ab82Index/32) {
			ab82T.Fatalf("entry %d %+v", ab82Index, ab82Entry)
		}
	}
	ab82Forged := ab82Full(17)
	ab82Forged.Ab82Raw = ab82Forged.Ab82Raw[:5]
	ab82Seven := Ab82Join(append(append([]Ab82Block{}, ab82Forward[:7]...), ab82Forged))
	if ab82Seven.Ab82Complete || len(ab82Seven.Ab82Entries) != 229 {
		ab82T.Fatal("forged state trusted")
	}
	ab82Blank := Ab82Block{}
	ab82Mixed := Ab82Join([]Ab82Block{ab82Full(10), ab82Blank})
	if ab82Mixed.Ab82State != Ab82JoinPartial || ab82Mixed.Ab82Blocks[1].Ab82State != Ab82StateRawOnly || ab82Mixed.Ab82Blocks[1].Ab82Contributed {
		ab82T.Fatal("raw-only block disturbed join")
	}
}

func TestAb82JoinPerBlockValidityAndConflict(ab82T *testing.T) {
	ab82Short := Ab82View(Ab82Input{Ab82Context: ab82Ctx(11, Ab82DirectionRequest), Ab82Values: []byte{1, 2}})
	ab82Other := ab82Full(11)
	ab82Joined := Ab82Join([]Ab82Block{ab82Full(10), ab82Short, ab82Other})
	if ab82Joined.Ab82State != Ab82JoinConflict || !reflect.DeepEqual(ab82Joined.Ab82Conflicts, []int{11}) || len(ab82Joined.Ab82Entries) != 32 {
		ab82T.Fatal("short vs full duplicate not conflicting")
	}
	if !ab82Joined.Ab82Blocks[1].Ab82Conflict || !ab82Joined.Ab82Blocks[2].Ab82Conflict || ab82Joined.Ab82Blocks[1].Ab82State != Ab82StateShort || !ab82Joined.Ab82Blocks[0].Ab82Contributed || ab82Joined.Ab82Blocks[0].Ab82Conflict {
		ab82T.Fatal("per-block validity")
	}
	ab82Cross := ab82Full(12)
	ab82Cross.Ab82Context.Ab82Family = "fam-b"
	if Ab82Join([]Ab82Block{ab82Full(10), ab82Cross}).Ab82State != Ab82JoinContextMismatch {
		ab82T.Fatal("cross family")
	}
	ab82Cross = ab82Full(12)
	ab82Cross.Ab82Context.Ab82Primary = 0xB6
	if Ab82Join([]Ab82Block{ab82Full(10), ab82Cross}).Ab82State != Ab82JoinPartial {
		ab82T.Fatal("unknown identity block must stay raw-only, not mismatch")
	}
}

func TestAb82OwnedMemory(ab82T *testing.T) {
	ab82Raw := append(ab82Repeat(0x01, 32), 0x09)
	ab82Got := Ab82View(Ab82Input{Ab82Context: ab82Ctx(10, Ab82DirectionRequest), Ab82Values: ab82Raw})
	ab82Raw[0], ab82Raw[32] = 0x77, 0x77
	if ab82Got.Ab82Raw[0] != 0x01 || ab82Got.Ab82Entries[0].Ab82RawByte != 0x01 || ab82Got.Ab82Remainder[0] != 0x09 {
		ab82T.Fatal("input aliased")
	}
	ab82Clone := ab82Got.Ab82Clone()
	ab82Clone.Ab82Raw[0], ab82Clone.Ab82Remainder[0], ab82Clone.Ab82Entries[0].Ab82RawByte = 0x55, 0x55, 0x55
	if ab82Got.Ab82Raw[0] != 0x01 || ab82Got.Ab82Remainder[0] != 0x09 || ab82Got.Ab82Entries[0].Ab82RawByte != 0x01 {
		ab82T.Fatal("clone aliased")
	}
	ab82Source := ab82Full(10)
	ab82Joined := Ab82Join([]Ab82Block{ab82Source})
	ab82Source.Ab82Raw[0] = 0x66
	ab82Joined.Ab82Entries[1].Ab82RawByte = 0x66
	if ab82Joined.Ab82Entries[0].Ab82RawByte != 10 || Ab82Join([]Ab82Block{ab82Full(10)}).Ab82Entries[1].Ab82RawByte != 10 {
		ab82T.Fatal("join shares memory")
	}
	if Ab82View(Ab82Input{Ab82Context: ab82Ctx(10, Ab82DirectionRequest)}).Ab82Raw != nil {
		ab82T.Fatal("nil raw invented")
	}
}
