package ebusstandard

import (
	standard09BodyViewEncodingJSON "encoding/json"
	standard09BodyViewOs "os"
	standard09BodyViewReflect "reflect"
	standard09BodyViewTesting "testing"
)

type standard09BodyViewVector struct {
	Standard09BodyViewLabel    string  `json:"ab25an25Label"`
	Standard09BodyViewService  uint16  `json:"ab25an25Service"`
	Standard09BodyViewRaw      []byte  `json:"ab25an25Raw"`
	Standard09BodyViewState    string  `json:"ab25an25State"`
	Standard09BodyViewData     []byte  `json:"ab25an25Data"`
	Standard09BodyViewTail     []byte  `json:"ab25an25Tail"`
	Standard09BodyViewLE       *uint16 `json:"ab25an25LE"`
	Standard09BodyViewLow      *uint8  `json:"ab25an25Low"`
	Standard09BodyViewHigh     *uint8  `json:"ab25an25High"`
	Standard09BodyViewValue    uint8   `json:"ab25an25Value"`
	Standard09BodyViewMaximum  *uint8  `json:"ab25an25Maximum"`
	Standard09BodyViewConflict bool    `json:"ab25an25Conflict"`
}

func standard09BodyViewKey(standard09BodyViewService uint16) Standard09BodyViewKey {
	return Standard09BodyViewKey{Standard09BodyViewFamily: "synthetic-family", Standard09BodyViewProfile: "synthetic-profile", Standard09BodyViewService: standard09BodyViewService, Standard09BodyViewDirection: "positive"}
}

func standard09BodyViewRequestKey(standard09BodyViewService uint16) Standard09BodyViewKey {
	standard09BodyViewResult := standard09BodyViewKey(standard09BodyViewService)
	standard09BodyViewResult.Standard09BodyViewDirection = "abstract-request"
	return standard09BodyViewResult
}

func standard09BodyViewCheckState(standard09BodyViewT *standard09BodyViewTesting.T, standard09BodyViewGot Standard09BodyViewState, standard09BodyViewKey Standard09BodyViewKey, standard09BodyViewRaw []byte, standard09BodyViewState string) {
	standard09BodyViewT.Helper()
	if standard09BodyViewGot.Standard09BodyViewKey != standard09BodyViewKey || !standard09BodyViewReflect.DeepEqual(standard09BodyViewGot.Standard09BodyViewRaw, append([]byte(nil), standard09BodyViewRaw...)) || standard09BodyViewGot.Standard09BodyViewState != standard09BodyViewState || standard09BodyViewGot.Standard09BodyViewContractGrade != "Proven" || standard09BodyViewGot.Standard09BodyViewMeaningGrade != "Hypothesis" || standard09BodyViewGot.Standard09BodyViewExecutionGrade != "Unknown" || standard09BodyViewGot.Standard09BodyViewPhysicalSupportGrade != "Unknown" || standard09BodyViewGot.Standard09BodyViewPhysicalOffsetsGrade != "Unknown" || standard09BodyViewGot.Standard09BodyViewPositivePrefixRelevanceGrade != "Unknown" {
		standard09BodyViewT.Fatalf("ab25an25 state: %+v; want key %+v raw %v state %s and independent grades", standard09BodyViewGot, standard09BodyViewKey, standard09BodyViewRaw, standard09BodyViewState)
	}
}

func TestAb25An25Vectors(standard09BodyViewT *standard09BodyViewTesting.T) {
	standard09BodyViewJSON, standard09BodyViewErr := standard09BodyViewOs.ReadFile("testdata/ab25an25_standard09_body_views.json")
	if standard09BodyViewErr != nil {
		standard09BodyViewT.Fatal(standard09BodyViewErr)
	}
	var standard09BodyViewFixture struct {
		Standard09BodyViewBodies    []standard09BodyViewVector `json:"ab25an25Bodies"`
		Standard09BodyViewAddresses []standard09BodyViewVector `json:"ab25an25Addresses"`
		Standard09BodyViewCounts    []standard09BodyViewVector `json:"ab25an25Counts"`
	}
	if standard09BodyViewErr = standard09BodyViewEncodingJSON.Unmarshal(standard09BodyViewJSON, &standard09BodyViewFixture); standard09BodyViewErr != nil {
		standard09BodyViewT.Fatal(standard09BodyViewErr)
	}
	if len(standard09BodyViewFixture.Standard09BodyViewBodies) != 21 || len(standard09BodyViewFixture.Standard09BodyViewAddresses) != 24 || len(standard09BodyViewFixture.Standard09BodyViewCounts) != 16 {
		standard09BodyViewT.Fatal("ab25an25 missing finite vectors")
	}
	for _, standard09BodyViewCase := range standard09BodyViewFixture.Standard09BodyViewBodies {
		standard09BodyViewT.Run(standard09BodyViewCase.Standard09BodyViewLabel, func(standard09BodyViewT *standard09BodyViewTesting.T) {
			standard09BodyViewKey := standard09BodyViewKey(standard09BodyViewCase.Standard09BodyViewService)
			standard09BodyViewGot := Standard09ViewBody(standard09BodyViewKey, standard09BodyViewKey, standard09BodyViewCase.Standard09BodyViewRaw)
			standard09BodyViewCheckState(standard09BodyViewT, standard09BodyViewGot.Standard09BodyViewState, standard09BodyViewKey, standard09BodyViewCase.Standard09BodyViewRaw, standard09BodyViewCase.Standard09BodyViewState)
			if !standard09BodyViewReflect.DeepEqual(standard09BodyViewGot.Standard09BodyViewData, standard09BodyViewCase.Standard09BodyViewData) || !standard09BodyViewReflect.DeepEqual(standard09BodyViewGot.Standard09BodyViewUnresolvedTail, standard09BodyViewCase.Standard09BodyViewTail) {
				standard09BodyViewT.Fatalf("ab25an25 body: %+v", standard09BodyViewGot)
			}
		})
	}
	for _, standard09BodyViewCase := range standard09BodyViewFixture.Standard09BodyViewAddresses {
		standard09BodyViewT.Run(standard09BodyViewCase.Standard09BodyViewLabel, func(standard09BodyViewT *standard09BodyViewTesting.T) {
			standard09BodyViewKey := standard09BodyViewRequestKey(standard09BodyViewCase.Standard09BodyViewService)
			standard09BodyViewGot := Standard09ViewAddress(standard09BodyViewKey, standard09BodyViewKey, standard09BodyViewCase.Standard09BodyViewRaw)
			standard09BodyViewCheckState(standard09BodyViewT, standard09BodyViewGot.Standard09BodyViewState, standard09BodyViewKey, standard09BodyViewCase.Standard09BodyViewRaw, standard09BodyViewCase.Standard09BodyViewState)
			if !standard09BodyViewReflect.DeepEqual(standard09BodyViewGot.Standard09BodyViewLE, standard09BodyViewCase.Standard09BodyViewLE) || !standard09BodyViewReflect.DeepEqual(standard09BodyViewGot.Standard09BodyViewLow, standard09BodyViewCase.Standard09BodyViewLow) || !standard09BodyViewReflect.DeepEqual(standard09BodyViewGot.Standard09BodyViewHigh, standard09BodyViewCase.Standard09BodyViewHigh) {
				standard09BodyViewT.Fatalf("ab25an25 address: %+v", standard09BodyViewGot)
			}
		})
	}
	for _, standard09BodyViewCase := range standard09BodyViewFixture.Standard09BodyViewCounts {
		standard09BodyViewT.Run(standard09BodyViewCase.Standard09BodyViewLabel, func(standard09BodyViewT *standard09BodyViewTesting.T) {
			standard09BodyViewKey := standard09BodyViewRequestKey(standard09BodyViewCase.Standard09BodyViewService)
			standard09BodyViewGot := Standard09ViewCount(standard09BodyViewKey, standard09BodyViewKey, standard09BodyViewCase.Standard09BodyViewValue)
			standard09BodyViewCheckState(standard09BodyViewT, standard09BodyViewGot.Standard09BodyViewState, standard09BodyViewKey, []byte{standard09BodyViewCase.Standard09BodyViewValue}, standard09BodyViewCase.Standard09BodyViewState)
			if standard09BodyViewGot.Standard09BodyViewValue != standard09BodyViewCase.Standard09BodyViewValue || !standard09BodyViewReflect.DeepEqual(standard09BodyViewGot.Standard09BodyViewCandidateMaximum, standard09BodyViewCase.Standard09BodyViewMaximum) || standard09BodyViewGot.Standard09BodyViewConflict != standard09BodyViewCase.Standard09BodyViewConflict {
				standard09BodyViewT.Fatalf("ab25an25 count: %+v", standard09BodyViewGot)
			}
		})
	}
}

func TestAb25An25Selection(standard09BodyViewT *standard09BodyViewTesting.T) {
	for _, standard09BodyViewService := range []uint16{0x0900, 0x0901, 0x0902, 0x0903} {
		standard09BodyViewSelected := standard09BodyViewKey(standard09BodyViewService)
		for standard09BodyViewIndex := 0; standard09BodyViewIndex < 4; standard09BodyViewIndex++ {
			standard09BodyViewOther := standard09BodyViewSelected
			switch standard09BodyViewIndex {
			case 0:
				standard09BodyViewOther.Standard09BodyViewFamily = "other-family"
			case 1:
				standard09BodyViewOther.Standard09BodyViewProfile = "other-profile"
			case 2:
				standard09BodyViewOther.Standard09BodyViewService ^= 1
			case 3:
				standard09BodyViewOther.Standard09BodyViewDirection = "request"
			}
			for _, standard09BodyViewPair := range [][2]Standard09BodyViewKey{{standard09BodyViewSelected, standard09BodyViewOther}, {standard09BodyViewOther, standard09BodyViewSelected}} {
				standard09BodyViewRaw := []byte{9, byte(standard09BodyViewService), 170}
				standard09BodyViewBody := Standard09ViewBody(standard09BodyViewPair[0], standard09BodyViewPair[1], standard09BodyViewRaw)
				standard09BodyViewCheckState(standard09BodyViewT, standard09BodyViewBody.Standard09BodyViewState, standard09BodyViewPair[1], standard09BodyViewRaw, "unknown")
				if standard09BodyViewBody.Standard09BodyViewData != nil || standard09BodyViewBody.Standard09BodyViewUnresolvedTail != nil {
					standard09BodyViewT.Fatal("ab25an25 mismatched body admitted")
				}
				standard09BodyViewRPair := [2]Standard09BodyViewKey{standard09BodyViewPair[0], standard09BodyViewPair[1]}
				standard09BodyViewRPair[0].Standard09BodyViewDirection, standard09BodyViewRPair[1].Standard09BodyViewDirection = "abstract-request", "abstract-request"
				if standard09BodyViewIndex == 3 {
					standard09BodyViewRPair[0].Standard09BodyViewDirection, standard09BodyViewRPair[1].Standard09BodyViewDirection = "positive", "positive"
				}
				standard09BodyViewAddress := Standard09ViewAddress(standard09BodyViewRPair[0], standard09BodyViewRPair[1], []byte{52, 18})
				standard09BodyViewCheckState(standard09BodyViewT, standard09BodyViewAddress.Standard09BodyViewState, standard09BodyViewRPair[1], []byte{52, 18}, "unknown")
				if standard09BodyViewAddress.Standard09BodyViewLE != nil || standard09BodyViewAddress.Standard09BodyViewLow != nil || standard09BodyViewAddress.Standard09BodyViewHigh != nil {
					standard09BodyViewT.Fatal("ab25an25 mismatched address admitted")
				}
				standard09BodyViewCount := Standard09ViewCount(standard09BodyViewRPair[0], standard09BodyViewRPair[1], 11)
				standard09BodyViewCheckState(standard09BodyViewT, standard09BodyViewCount.Standard09BodyViewState, standard09BodyViewRPair[1], []byte{11}, "unknown")
				if standard09BodyViewCount.Standard09BodyViewValue != 11 || standard09BodyViewCount.Standard09BodyViewCandidateMaximum != nil || standard09BodyViewCount.Standard09BodyViewConflict {
					standard09BodyViewT.Fatal("ab25an25 mismatched count admitted")
				}
			}
		}
	}
	for standard09BodyViewIndex := 0; standard09BodyViewIndex < 4; standard09BodyViewIndex++ {
		standard09BodyViewInvalid := standard09BodyViewKey(0x0900)
		standard09BodyViewRInvalid := standard09BodyViewRequestKey(0x0900)
		switch standard09BodyViewIndex {
		case 0:
			standard09BodyViewInvalid.Standard09BodyViewFamily, standard09BodyViewRInvalid.Standard09BodyViewFamily = "", ""
		case 1:
			standard09BodyViewInvalid.Standard09BodyViewProfile, standard09BodyViewRInvalid.Standard09BodyViewProfile = "", ""
		case 2:
			standard09BodyViewInvalid.Standard09BodyViewService, standard09BodyViewRInvalid.Standard09BodyViewService = 0x0904, 0x0904
		case 3:
			standard09BodyViewInvalid.Standard09BodyViewDirection, standard09BodyViewRInvalid.Standard09BodyViewDirection = "request", "request"
		}
		standard09BodyViewCheckState(standard09BodyViewT, Standard09ViewBody(standard09BodyViewInvalid, standard09BodyViewInvalid, []byte{9, 0}).Standard09BodyViewState, standard09BodyViewInvalid, []byte{9, 0}, "unknown")
		standard09BodyViewCheckState(standard09BodyViewT, Standard09ViewAddress(standard09BodyViewRInvalid, standard09BodyViewRInvalid, []byte{52, 18}).Standard09BodyViewState, standard09BodyViewRInvalid, []byte{52, 18}, "unknown")
		standard09BodyViewCheckState(standard09BodyViewT, Standard09ViewCount(standard09BodyViewRInvalid, standard09BodyViewRInvalid, 11).Standard09BodyViewState, standard09BodyViewRInvalid, []byte{11}, "unknown")
	}
}

func TestAb25An25Isolation(standard09BodyViewT *standard09BodyViewTesting.T) {
	for _, standard09BodyViewService := range []uint16{0x0900, 0x0901, 0x0902, 0x0903} {
		standard09BodyViewKey := standard09BodyViewKey(standard09BodyViewService)
		standard09BodyViewRaw := []byte{9, byte(standard09BodyViewService), 170, 85}
		standard09BodyViewBody := Standard09ViewBody(standard09BodyViewKey, standard09BodyViewKey, standard09BodyViewRaw)
		standard09BodyViewRaw[2] = 1
		if standard09BodyViewBody.Standard09BodyViewState.Standard09BodyViewRaw[2] != 170 {
			standard09BodyViewT.Fatal("ab25an25 input aliases raw")
		}
		standard09BodyViewBody.Standard09BodyViewState.Standard09BodyViewRaw[2] = 2
		if standard09BodyViewService == 0x0901 {
			if standard09BodyViewBody.Standard09BodyViewUnresolvedTail[0] != 170 {
				standard09BodyViewT.Fatal("ab25an25 tail aliases raw")
			}
			standard09BodyViewBody.Standard09BodyViewUnresolvedTail[0] = 3
		} else {
			standard09BodyViewIndex := 0
			if standard09BodyViewService == 0x0902 {
				standard09BodyViewIndex = 2
			}
			if standard09BodyViewBody.Standard09BodyViewData[standard09BodyViewIndex] != 170 {
				standard09BodyViewT.Fatal("ab25an25 data aliases raw")
			}
			standard09BodyViewBody.Standard09BodyViewData[standard09BodyViewIndex] = 3
		}
		if standard09BodyViewRaw[2] != 1 || standard09BodyViewBody.Standard09BodyViewState.Standard09BodyViewRaw[2] != 2 {
			standard09BodyViewT.Fatal("ab25an25 body aliases input or raw")
		}
		standard09BodyViewField := []byte{52, 18}
		standard09BodyViewRKey := standard09BodyViewRequestKey(standard09BodyViewService)
		standard09BodyViewAddress := Standard09ViewAddress(standard09BodyViewRKey, standard09BodyViewRKey, standard09BodyViewField)
		standard09BodyViewSecond := Standard09ViewAddress(standard09BodyViewRKey, standard09BodyViewRKey, standard09BodyViewField)
		standard09BodyViewField[0] = 0
		if standard09BodyViewAddress.Standard09BodyViewState.Standard09BodyViewRaw[0] != 52 {
			standard09BodyViewT.Fatal("ab25an25 address raw aliases input")
		}
		standard09BodyViewAddress.Standard09BodyViewState.Standard09BodyViewRaw[0] = 1
		if standard09BodyViewService < 0x0902 {
			if *standard09BodyViewAddress.Standard09BodyViewLE != 4660 {
				standard09BodyViewT.Fatal("ab25an25 LE aliases raw")
			}
			*standard09BodyViewAddress.Standard09BodyViewLE = 0
			if *standard09BodyViewSecond.Standard09BodyViewLE != 4660 {
				standard09BodyViewT.Fatal("ab25an25 LE shared between calls")
			}
		} else {
			if *standard09BodyViewAddress.Standard09BodyViewLow != 52 || *standard09BodyViewAddress.Standard09BodyViewHigh != 18 {
				standard09BodyViewT.Fatal("ab25an25 split aliases raw")
			}
			*standard09BodyViewAddress.Standard09BodyViewLow = 0
			*standard09BodyViewAddress.Standard09BodyViewHigh = 0
			if *standard09BodyViewSecond.Standard09BodyViewLow != 52 || *standard09BodyViewSecond.Standard09BodyViewHigh != 18 {
				standard09BodyViewT.Fatal("ab25an25 split shared between calls")
			}
		}
	}
	standard09BodyViewKey := standard09BodyViewRequestKey(0x0902)
	standard09BodyViewCount := Standard09ViewCount(standard09BodyViewKey, standard09BodyViewKey, 11)
	*standard09BodyViewCount.Standard09BodyViewCandidateMaximum = 0
	standard09BodyViewCount.Standard09BodyViewState.Standard09BodyViewRaw[0] = 0
	standard09BodyViewSecond := Standard09ViewCount(standard09BodyViewKey, standard09BodyViewKey, 11)
	if standard09BodyViewCount.Standard09BodyViewValue != 11 || *standard09BodyViewSecond.Standard09BodyViewCandidateMaximum != 10 || standard09BodyViewSecond.Standard09BodyViewState.Standard09BodyViewRaw[0] != 11 {
		standard09BodyViewT.Fatal("ab25an25 count storage shared")
	}
	for _, standard09BodyViewRaw := range [][]byte{{9}, {8, 0}} {
		standard09BodyViewKey.Standard09BodyViewService, standard09BodyViewKey.Standard09BodyViewDirection = 0x0900, "positive"
		standard09BodyViewBody := Standard09ViewBody(standard09BodyViewKey, standard09BodyViewKey, standard09BodyViewRaw)
		standard09BodyViewBody.Standard09BodyViewState.Standard09BodyViewRaw[0] = 255
		if standard09BodyViewRaw[0] == 255 {
			standard09BodyViewT.Fatal("ab25an25 incomplete raw aliases input")
		}
	}
	standard09BodyViewKey.Standard09BodyViewService, standard09BodyViewKey.Standard09BodyViewDirection = 0x0904, "positive"
	standard09BodyViewRaw := []byte{9, 4}
	standard09BodyViewBody := Standard09ViewBody(standard09BodyViewKey, standard09BodyViewKey, standard09BodyViewRaw)
	standard09BodyViewBody.Standard09BodyViewState.Standard09BodyViewRaw[0] = 255
	if standard09BodyViewRaw[0] != 9 {
		standard09BodyViewT.Fatal("ab25an25 unknown raw aliases input")
	}
}

func TestAb25An25Extent(standard09BodyViewT *standard09BodyViewTesting.T) {
	for _, standard09BodyViewService := range []uint16{0x0900, 0x0902, 0x0903} {
		standard09BodyViewKey := standard09BodyViewKey(standard09BodyViewService)
		standard09BodyViewRaw := make([]byte, 4096)
		for standard09BodyViewIndex := range standard09BodyViewRaw {
			standard09BodyViewRaw[standard09BodyViewIndex] = 170
		}
		standard09BodyViewWant := append([]byte{}, standard09BodyViewRaw...)
		if standard09BodyViewService != 0x0902 {
			standard09BodyViewRaw = append([]byte{9, byte(standard09BodyViewService)}, standard09BodyViewRaw...)
		}
		standard09BodyViewGot := Standard09ViewBody(standard09BodyViewKey, standard09BodyViewKey, standard09BodyViewRaw)
		standard09BodyViewCheckState(standard09BodyViewT, standard09BodyViewGot.Standard09BodyViewState, standard09BodyViewKey, standard09BodyViewRaw, "candidate")
		if !standard09BodyViewReflect.DeepEqual(standard09BodyViewGot.Standard09BodyViewData, standard09BodyViewWant) {
			standard09BodyViewT.Fatal("ab25an25 inferred extent limit")
		}
	}
}

func TestAb25An25DirectionByView(standard09BodyViewT *standard09BodyViewTesting.T) {
	for _, standard09BodyViewService := range []uint16{0x0900, 0x0901, 0x0902, 0x0903} {
		standard09BodyViewPositive, standard09BodyViewRequest := standard09BodyViewKey(standard09BodyViewService), standard09BodyViewRequestKey(standard09BodyViewService)
		standard09BodyViewCheckState(standard09BodyViewT, Standard09ViewBody(standard09BodyViewRequest, standard09BodyViewRequest, []byte{9, byte(standard09BodyViewService)}).Standard09BodyViewState, standard09BodyViewRequest, []byte{9, byte(standard09BodyViewService)}, "unknown")
		standard09BodyViewCheckState(standard09BodyViewT, Standard09ViewAddress(standard09BodyViewPositive, standard09BodyViewPositive, []byte{52, 18}).Standard09BodyViewState, standard09BodyViewPositive, []byte{52, 18}, "unknown")
		standard09BodyViewCheckState(standard09BodyViewT, Standard09ViewCount(standard09BodyViewPositive, standard09BodyViewPositive, 10).Standard09BodyViewState, standard09BodyViewPositive, []byte{10}, "unknown")
		standard09BodyViewCheckState(standard09BodyViewT, Standard09ViewAddress(standard09BodyViewRequest, standard09BodyViewRequest, []byte{52, 18}).Standard09BodyViewState, standard09BodyViewRequest, []byte{52, 18}, "candidate")
	}
}
