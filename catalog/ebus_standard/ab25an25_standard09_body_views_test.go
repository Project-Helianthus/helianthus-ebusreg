package ebus_standard_catalog

import (
	ab25an25json "encoding/json"
	ab25an25os "os"
	ab25an25reflect "reflect"
	ab25an25testing "testing"
)

type ab25an25Vector struct {
	Ab25An25Label    string  `json:"ab25an25Label"`
	Ab25An25Service  uint16  `json:"ab25an25Service"`
	Ab25An25Raw      []byte  `json:"ab25an25Raw"`
	Ab25An25State    string  `json:"ab25an25State"`
	Ab25An25Data     []byte  `json:"ab25an25Data"`
	Ab25An25Tail     []byte  `json:"ab25an25Tail"`
	Ab25An25LE       *uint16 `json:"ab25an25LE"`
	Ab25An25Low      *uint8  `json:"ab25an25Low"`
	Ab25An25High     *uint8  `json:"ab25an25High"`
	Ab25An25Value    uint8   `json:"ab25an25Value"`
	Ab25An25Maximum  *uint8  `json:"ab25an25Maximum"`
	Ab25An25Conflict bool    `json:"ab25an25Conflict"`
}

func ab25an25Key(ab25an25Service uint16) Ab25An25Key {
	return Ab25An25Key{Ab25An25Family: "synthetic-family", Ab25An25Profile: "synthetic-profile", Ab25An25Service: ab25an25Service, Ab25An25Direction: "positive"}
}

func ab25an25CheckState(ab25an25T *ab25an25testing.T, ab25an25Got Ab25An25State, ab25an25Key Ab25An25Key, ab25an25Raw []byte, ab25an25State string) {
	ab25an25T.Helper()
	if ab25an25Got.Ab25An25Key != ab25an25Key || !ab25an25reflect.DeepEqual(ab25an25Got.Ab25An25Raw, append([]byte(nil), ab25an25Raw...)) || ab25an25Got.Ab25An25State != ab25an25State || ab25an25Got.Ab25An25ContractGrade != "Proven" || ab25an25Got.Ab25An25MeaningGrade != "Hypothesis" || ab25an25Got.Ab25An25ExecutionGrade != "Unknown" || ab25an25Got.Ab25An25PhysicalSupportGrade != "Unknown" || ab25an25Got.Ab25An25PhysicalOffsetsGrade != "Unknown" || ab25an25Got.Ab25An25PositivePrefixRelevanceGrade != "Unknown" {
		ab25an25T.Fatalf("ab25an25 state: %+v; want key %+v raw %v state %s and independent grades", ab25an25Got, ab25an25Key, ab25an25Raw, ab25an25State)
	}
}

func TestAb25An25Vectors(ab25an25T *ab25an25testing.T) {
	ab25an25JSON, ab25an25Err := ab25an25os.ReadFile("testdata/ab25an25_standard09_body_views.json")
	if ab25an25Err != nil {
		ab25an25T.Fatal(ab25an25Err)
	}
	var ab25an25Fixture struct {
		Ab25An25Bodies    []ab25an25Vector `json:"ab25an25Bodies"`
		Ab25An25Addresses []ab25an25Vector `json:"ab25an25Addresses"`
		Ab25An25Counts    []ab25an25Vector `json:"ab25an25Counts"`
	}
	if ab25an25Err = ab25an25json.Unmarshal(ab25an25JSON, &ab25an25Fixture); ab25an25Err != nil {
		ab25an25T.Fatal(ab25an25Err)
	}
	if len(ab25an25Fixture.Ab25An25Bodies) != 21 || len(ab25an25Fixture.Ab25An25Addresses) != 24 || len(ab25an25Fixture.Ab25An25Counts) != 16 {
		ab25an25T.Fatal("ab25an25 missing finite vectors")
	}
	for _, ab25an25Case := range ab25an25Fixture.Ab25An25Bodies {
		ab25an25T.Run(ab25an25Case.Ab25An25Label, func(ab25an25T *ab25an25testing.T) {
			ab25an25Key := ab25an25Key(ab25an25Case.Ab25An25Service)
			ab25an25Got := Ab25An25ViewBody(ab25an25Key, ab25an25Key, ab25an25Case.Ab25An25Raw)
			ab25an25CheckState(ab25an25T, ab25an25Got.Ab25An25State, ab25an25Key, ab25an25Case.Ab25An25Raw, ab25an25Case.Ab25An25State)
			if !ab25an25reflect.DeepEqual(ab25an25Got.Ab25An25Data, ab25an25Case.Ab25An25Data) || !ab25an25reflect.DeepEqual(ab25an25Got.Ab25An25UnresolvedTail, ab25an25Case.Ab25An25Tail) {
				ab25an25T.Fatalf("ab25an25 body: %+v", ab25an25Got)
			}
		})
	}
	for _, ab25an25Case := range ab25an25Fixture.Ab25An25Addresses {
		ab25an25T.Run(ab25an25Case.Ab25An25Label, func(ab25an25T *ab25an25testing.T) {
			ab25an25Key := ab25an25Key(ab25an25Case.Ab25An25Service)
			ab25an25Got := Ab25An25ViewAddress(ab25an25Key, ab25an25Key, ab25an25Case.Ab25An25Raw)
			ab25an25CheckState(ab25an25T, ab25an25Got.Ab25An25State, ab25an25Key, ab25an25Case.Ab25An25Raw, ab25an25Case.Ab25An25State)
			if !ab25an25reflect.DeepEqual(ab25an25Got.Ab25An25LE, ab25an25Case.Ab25An25LE) || !ab25an25reflect.DeepEqual(ab25an25Got.Ab25An25Low, ab25an25Case.Ab25An25Low) || !ab25an25reflect.DeepEqual(ab25an25Got.Ab25An25High, ab25an25Case.Ab25An25High) {
				ab25an25T.Fatalf("ab25an25 address: %+v", ab25an25Got)
			}
		})
	}
	for _, ab25an25Case := range ab25an25Fixture.Ab25An25Counts {
		ab25an25T.Run(ab25an25Case.Ab25An25Label, func(ab25an25T *ab25an25testing.T) {
			ab25an25Key := ab25an25Key(ab25an25Case.Ab25An25Service)
			ab25an25Got := Ab25An25ViewCount(ab25an25Key, ab25an25Key, ab25an25Case.Ab25An25Value)
			ab25an25CheckState(ab25an25T, ab25an25Got.Ab25An25State, ab25an25Key, []byte{ab25an25Case.Ab25An25Value}, ab25an25Case.Ab25An25State)
			if ab25an25Got.Ab25An25Value != ab25an25Case.Ab25An25Value || !ab25an25reflect.DeepEqual(ab25an25Got.Ab25An25CandidateMaximum, ab25an25Case.Ab25An25Maximum) || ab25an25Got.Ab25An25Conflict != ab25an25Case.Ab25An25Conflict {
				ab25an25T.Fatalf("ab25an25 count: %+v", ab25an25Got)
			}
		})
	}
}

func TestAb25An25Selection(ab25an25T *ab25an25testing.T) {
	for _, ab25an25Service := range []uint16{0x0900, 0x0901, 0x0902, 0x0903} {
		ab25an25Selected := ab25an25Key(ab25an25Service)
		for ab25an25Index := 0; ab25an25Index < 4; ab25an25Index++ {
			ab25an25Other := ab25an25Selected
			switch ab25an25Index {
			case 0:
				ab25an25Other.Ab25An25Family = "other-family"
			case 1:
				ab25an25Other.Ab25An25Profile = "other-profile"
			case 2:
				ab25an25Other.Ab25An25Service ^= 1
			case 3:
				ab25an25Other.Ab25An25Direction = "request"
			}
			for _, ab25an25Pair := range [][2]Ab25An25Key{{ab25an25Selected, ab25an25Other}, {ab25an25Other, ab25an25Selected}} {
				ab25an25Raw := []byte{9, byte(ab25an25Service), 170}
				ab25an25Body := Ab25An25ViewBody(ab25an25Pair[0], ab25an25Pair[1], ab25an25Raw)
				ab25an25CheckState(ab25an25T, ab25an25Body.Ab25An25State, ab25an25Pair[1], ab25an25Raw, "unknown")
				if ab25an25Body.Ab25An25Data != nil || ab25an25Body.Ab25An25UnresolvedTail != nil {
					ab25an25T.Fatal("ab25an25 mismatched body admitted")
				}
				ab25an25Address := Ab25An25ViewAddress(ab25an25Pair[0], ab25an25Pair[1], []byte{52, 18})
				ab25an25CheckState(ab25an25T, ab25an25Address.Ab25An25State, ab25an25Pair[1], []byte{52, 18}, "unknown")
				if ab25an25Address.Ab25An25LE != nil || ab25an25Address.Ab25An25Low != nil || ab25an25Address.Ab25An25High != nil {
					ab25an25T.Fatal("ab25an25 mismatched address admitted")
				}
				ab25an25Count := Ab25An25ViewCount(ab25an25Pair[0], ab25an25Pair[1], 11)
				ab25an25CheckState(ab25an25T, ab25an25Count.Ab25An25State, ab25an25Pair[1], []byte{11}, "unknown")
				if ab25an25Count.Ab25An25Value != 11 || ab25an25Count.Ab25An25CandidateMaximum != nil || ab25an25Count.Ab25An25Conflict {
					ab25an25T.Fatal("ab25an25 mismatched count admitted")
				}
			}
		}
	}
	for ab25an25Index := 0; ab25an25Index < 4; ab25an25Index++ {
		ab25an25Invalid := ab25an25Key(0x0900)
		switch ab25an25Index {
		case 0:
			ab25an25Invalid.Ab25An25Family = ""
		case 1:
			ab25an25Invalid.Ab25An25Profile = ""
		case 2:
			ab25an25Invalid.Ab25An25Service = 0x0904
		case 3:
			ab25an25Invalid.Ab25An25Direction = "request"
		}
		ab25an25CheckState(ab25an25T, Ab25An25ViewBody(ab25an25Invalid, ab25an25Invalid, []byte{9, 0}).Ab25An25State, ab25an25Invalid, []byte{9, 0}, "unknown")
		ab25an25CheckState(ab25an25T, Ab25An25ViewAddress(ab25an25Invalid, ab25an25Invalid, []byte{52, 18}).Ab25An25State, ab25an25Invalid, []byte{52, 18}, "unknown")
		ab25an25CheckState(ab25an25T, Ab25An25ViewCount(ab25an25Invalid, ab25an25Invalid, 11).Ab25An25State, ab25an25Invalid, []byte{11}, "unknown")
	}
}

func TestAb25An25Isolation(ab25an25T *ab25an25testing.T) {
	for _, ab25an25Service := range []uint16{0x0900, 0x0901, 0x0902, 0x0903} {
		ab25an25Key := ab25an25Key(ab25an25Service)
		ab25an25Raw := []byte{9, byte(ab25an25Service), 170, 85}
		ab25an25Body := Ab25An25ViewBody(ab25an25Key, ab25an25Key, ab25an25Raw)
		ab25an25Raw[2] = 1
		if ab25an25Body.Ab25An25State.Ab25An25Raw[2] != 170 {
			ab25an25T.Fatal("ab25an25 input aliases raw")
		}
		ab25an25Body.Ab25An25State.Ab25An25Raw[2] = 2
		if ab25an25Service == 0x0901 {
			if ab25an25Body.Ab25An25UnresolvedTail[0] != 170 {
				ab25an25T.Fatal("ab25an25 tail aliases raw")
			}
			ab25an25Body.Ab25An25UnresolvedTail[0] = 3
		} else {
			ab25an25Index := 0
			if ab25an25Service == 0x0902 {
				ab25an25Index = 2
			}
			if ab25an25Body.Ab25An25Data[ab25an25Index] != 170 {
				ab25an25T.Fatal("ab25an25 data aliases raw")
			}
			ab25an25Body.Ab25An25Data[ab25an25Index] = 3
		}
		if ab25an25Raw[2] != 1 || ab25an25Body.Ab25An25State.Ab25An25Raw[2] != 2 {
			ab25an25T.Fatal("ab25an25 body aliases input or raw")
		}
		ab25an25Field := []byte{52, 18}
		ab25an25Address := Ab25An25ViewAddress(ab25an25Key, ab25an25Key, ab25an25Field)
		ab25an25Second := Ab25An25ViewAddress(ab25an25Key, ab25an25Key, ab25an25Field)
		ab25an25Field[0] = 0
		if ab25an25Address.Ab25An25State.Ab25An25Raw[0] != 52 {
			ab25an25T.Fatal("ab25an25 address raw aliases input")
		}
		ab25an25Address.Ab25An25State.Ab25An25Raw[0] = 1
		if ab25an25Service < 0x0902 {
			if *ab25an25Address.Ab25An25LE != 4660 {
				ab25an25T.Fatal("ab25an25 LE aliases raw")
			}
			*ab25an25Address.Ab25An25LE = 0
			if *ab25an25Second.Ab25An25LE != 4660 {
				ab25an25T.Fatal("ab25an25 LE shared between calls")
			}
		} else {
			if *ab25an25Address.Ab25An25Low != 52 || *ab25an25Address.Ab25An25High != 18 {
				ab25an25T.Fatal("ab25an25 split aliases raw")
			}
			*ab25an25Address.Ab25An25Low = 0
			*ab25an25Address.Ab25An25High = 0
			if *ab25an25Second.Ab25An25Low != 52 || *ab25an25Second.Ab25An25High != 18 {
				ab25an25T.Fatal("ab25an25 split shared between calls")
			}
		}
	}
	ab25an25Key := ab25an25Key(0x0902)
	ab25an25Count := Ab25An25ViewCount(ab25an25Key, ab25an25Key, 11)
	*ab25an25Count.Ab25An25CandidateMaximum = 0
	ab25an25Count.Ab25An25State.Ab25An25Raw[0] = 0
	ab25an25Second := Ab25An25ViewCount(ab25an25Key, ab25an25Key, 11)
	if ab25an25Count.Ab25An25Value != 11 || *ab25an25Second.Ab25An25CandidateMaximum != 10 || ab25an25Second.Ab25An25State.Ab25An25Raw[0] != 11 {
		ab25an25T.Fatal("ab25an25 count storage shared")
	}
	for _, ab25an25Raw := range [][]byte{{9}, {8, 0}} {
		ab25an25Key.Ab25An25Service = 0x0900
		ab25an25Body := Ab25An25ViewBody(ab25an25Key, ab25an25Key, ab25an25Raw)
		ab25an25Body.Ab25An25State.Ab25An25Raw[0] = 255
		if ab25an25Raw[0] == 255 {
			ab25an25T.Fatal("ab25an25 incomplete raw aliases input")
		}
	}
	ab25an25Key.Ab25An25Service = 0x0904
	ab25an25Raw := []byte{9, 4}
	ab25an25Body := Ab25An25ViewBody(ab25an25Key, ab25an25Key, ab25an25Raw)
	ab25an25Body.Ab25An25State.Ab25An25Raw[0] = 255
	if ab25an25Raw[0] != 9 {
		ab25an25T.Fatal("ab25an25 unknown raw aliases input")
	}
}

func TestAb25An25Extent(ab25an25T *ab25an25testing.T) {
	for _, ab25an25Service := range []uint16{0x0900, 0x0902, 0x0903} {
		ab25an25Key := ab25an25Key(ab25an25Service)
		ab25an25Raw := make([]byte, 4096)
		for ab25an25Index := range ab25an25Raw {
			ab25an25Raw[ab25an25Index] = 170
		}
		ab25an25Want := append([]byte{}, ab25an25Raw...)
		if ab25an25Service != 0x0902 {
			ab25an25Raw = append([]byte{9, byte(ab25an25Service)}, ab25an25Raw...)
		}
		ab25an25Got := Ab25An25ViewBody(ab25an25Key, ab25an25Key, ab25an25Raw)
		ab25an25CheckState(ab25an25T, ab25an25Got.Ab25An25State, ab25an25Key, ab25an25Raw, "candidate")
		if !ab25an25reflect.DeepEqual(ab25an25Got.Ab25An25Data, ab25an25Want) {
			ab25an25T.Fatal("ab25an25 inferred extent limit")
		}
	}
}
