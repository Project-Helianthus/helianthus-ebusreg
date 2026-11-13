package vaillant

import (
	b503HistoryBranchJSON "encoding/json"
	b503HistoryBranchOs "os"
	b503HistoryBranchReflect "reflect"
	b503HistoryBranchTesting "testing"
)

type b503HistoryBranchVector struct {
	B503HistoryBranchLabel    string                  `json:"Ab05An05Label"`
	B503HistoryBranchInput    B503HistoryBranchFields `json:"Ab05An05Input"`
	B503HistoryBranchExpected B503HistoryBranchView   `json:"Ab05An05Expected"`
}

func b503HistoryBranchVectors(b503HistoryBranchT *b503HistoryBranchTesting.T) []b503HistoryBranchVector {
	b503HistoryBranchT.Helper()
	b503HistoryBranchBytes, b503HistoryBranchError := b503HistoryBranchOs.ReadFile("testdata/ab05an05_b503_history_branches.json")
	if b503HistoryBranchError != nil {
		b503HistoryBranchT.Fatal(b503HistoryBranchError)
	}
	var b503HistoryBranchCases []b503HistoryBranchVector
	if b503HistoryBranchError = b503HistoryBranchJSON.Unmarshal(b503HistoryBranchBytes, &b503HistoryBranchCases); b503HistoryBranchError != nil {
		b503HistoryBranchT.Fatal(b503HistoryBranchError)
	}
	if len(b503HistoryBranchCases) != 7 {
		b503HistoryBranchT.Fatal("ab05an05 expected seven finite vectors")
	}
	return b503HistoryBranchCases
}

func b503HistoryBranchEqual(b503HistoryBranchT *b503HistoryBranchTesting.T, b503HistoryBranchGot, b503HistoryBranchWant B503HistoryBranchView) {
	b503HistoryBranchT.Helper()
	if !b503HistoryBranchReflect.DeepEqual(b503HistoryBranchGot, b503HistoryBranchWant) {
		b503HistoryBranchT.Fatalf("ab05an05 got %#v want %#v", b503HistoryBranchGot, b503HistoryBranchWant)
	}
}

func TestAb05An05IndependentVectors(b503HistoryBranchT *b503HistoryBranchTesting.T) {
	for _, b503HistoryBranchCase := range b503HistoryBranchVectors(b503HistoryBranchT) {
		b503HistoryBranchT.Run(b503HistoryBranchCase.B503HistoryBranchLabel, func(b503HistoryBranchT *b503HistoryBranchTesting.T) {
			b503HistoryBranchEqual(b503HistoryBranchT, B503HistoryBranchHistory(b503HistoryBranchCase.B503HistoryBranchInput), b503HistoryBranchCase.B503HistoryBranchExpected)
		})
	}
}

func TestAb05An05AllRemainderBytes(b503HistoryBranchT *b503HistoryBranchTesting.T) {
	b503HistoryBranchCase := b503HistoryBranchVectors(b503HistoryBranchT)[0]
	for b503HistoryBranchIndex := 0; b503HistoryBranchIndex < 7; b503HistoryBranchIndex++ {
		for b503HistoryBranchValue := 0; b503HistoryBranchValue < 255; b503HistoryBranchValue++ {
			b503HistoryBranchInput := b503HistoryBranchCase.B503HistoryBranchInput
			b503HistoryBranchInput.B503HistoryBranchRemainder[b503HistoryBranchIndex] = uint8(b503HistoryBranchValue)
			b503HistoryBranchWant := b503HistoryBranchCase.B503HistoryBranchExpected
			b503HistoryBranchWant.B503HistoryBranchFields = b503HistoryBranchInput
			b503HistoryBranchWant.B503HistoryBranchState = "unknown-branch"
			b503HistoryBranchWant.B503HistoryBranchMeaning = "Unknown"
			b503HistoryBranchEqual(b503HistoryBranchT, B503HistoryBranchHistory(b503HistoryBranchInput), b503HistoryBranchWant)
		}
	}
	b503HistoryBranchInput := b503HistoryBranchCase.B503HistoryBranchInput
	b503HistoryBranchInput.B503HistoryBranchRemainderPresent = false
	if B503HistoryBranchHistory(b503HistoryBranchInput).B503HistoryBranchState != "unknown-branch" {
		b503HistoryBranchT.Fatal("ab05an05 absent remainder admitted")
	}
}

func TestAb05An05MarkerAndContext(b503HistoryBranchT *b503HistoryBranchTesting.T) {
	b503HistoryBranchCase := b503HistoryBranchVectors(b503HistoryBranchT)[3]
	for b503HistoryBranchMarker := 3; b503HistoryBranchMarker <= 255; b503HistoryBranchMarker++ {
		b503HistoryBranchInput := b503HistoryBranchCase.B503HistoryBranchInput
		b503HistoryBranchInput.B503HistoryBranchMarker = uint8(b503HistoryBranchMarker)
		b503HistoryBranchWant := b503HistoryBranchCase.B503HistoryBranchExpected
		b503HistoryBranchWant.B503HistoryBranchFields = b503HistoryBranchInput
		b503HistoryBranchWant.B503HistoryBranchState, b503HistoryBranchWant.B503HistoryBranchMeaning, b503HistoryBranchWant.B503HistoryBranchComplete = "unknown-branch", "Unknown", false
		for b503HistoryBranchIndex := range b503HistoryBranchWant.B503HistoryBranchComponents {
			b503HistoryBranchWant.B503HistoryBranchComponents[b503HistoryBranchIndex].B503HistoryBranchState = "unknown"
		}
		b503HistoryBranchEqual(b503HistoryBranchT, B503HistoryBranchHistory(b503HistoryBranchInput), b503HistoryBranchWant)
	}
	for _, b503HistoryBranchField := range []string{"B503HistoryBranchFamily", "B503HistoryBranchVariant", "B503HistoryBranchOperation", "B503HistoryBranchDirection"} {
		b503HistoryBranchT.Run("ab05an05_"+b503HistoryBranchField, func(b503HistoryBranchT *b503HistoryBranchTesting.T) {
			for b503HistoryBranchMode := 0; b503HistoryBranchMode < 2; b503HistoryBranchMode++ {
				b503HistoryBranchInput := b503HistoryBranchCase.B503HistoryBranchInput
				b503HistoryBranchAxis := b503HistoryBranchReflect.ValueOf(&b503HistoryBranchInput.B503HistoryBranchContext).Elem().FieldByName(b503HistoryBranchField)
				if b503HistoryBranchAxis.Kind() == b503HistoryBranchReflect.String {
					if b503HistoryBranchMode == 0 {
						b503HistoryBranchAxis.SetString("")
					} else {
						b503HistoryBranchAxis.SetString("ab05an05_other")
					}
				} else {
					b503HistoryBranchAxis.SetUint(uint64(b503HistoryBranchMode * 7))
				}
				b503HistoryBranchWant := b503HistoryBranchCase.B503HistoryBranchExpected
				b503HistoryBranchWant.B503HistoryBranchFields = b503HistoryBranchInput
				if b503HistoryBranchMode == 0 || b503HistoryBranchField == "B503HistoryBranchOperation" {
					b503HistoryBranchWant.B503HistoryBranchState, b503HistoryBranchWant.B503HistoryBranchMeaning, b503HistoryBranchWant.B503HistoryBranchComplete = "unknown-branch", "Unknown", false
					for b503HistoryBranchIndex := range b503HistoryBranchWant.B503HistoryBranchComponents {
						b503HistoryBranchWant.B503HistoryBranchComponents[b503HistoryBranchIndex].B503HistoryBranchState = "unknown"
					}
				}
				b503HistoryBranchEqual(b503HistoryBranchT, B503HistoryBranchHistory(b503HistoryBranchInput), b503HistoryBranchWant)
			}
		})
	}
	b503HistoryBranchInput := b503HistoryBranchCase.B503HistoryBranchInput
	b503HistoryBranchInput.B503HistoryBranchMarkerPresent = false
	if B503HistoryBranchHistory(b503HistoryBranchInput).B503HistoryBranchState != "unknown-branch" {
		b503HistoryBranchT.Fatal("ab05an05 absent marker admitted")
	}
}

func TestAb05An05IndependentCells(b503HistoryBranchT *b503HistoryBranchTesting.T) {
	b503HistoryBranchCase := b503HistoryBranchVectors(b503HistoryBranchT)[3]
	for b503HistoryBranchIndex, b503HistoryBranchName := range []string{"ab05an05_minute", "ab05an05_hour", "ab05an05_day", "ab05an05_month", "ab05an05_year"} {
		b503HistoryBranchT.Run(b503HistoryBranchName, func(b503HistoryBranchT *b503HistoryBranchTesting.T) {
			for _, b503HistoryBranchValue := range []uint64{0, 1, 12, 13, 23, 24, 31, 32, 59, 60, 99, 100, 255, ^uint64(0)} {
				b503HistoryBranchInput := b503HistoryBranchCase.B503HistoryBranchInput
				b503HistoryBranchInput.B503HistoryBranchCalendar[b503HistoryBranchIndex].B503HistoryBranchValue = b503HistoryBranchValue
				b503HistoryBranchWant := b503HistoryBranchCase.B503HistoryBranchExpected
				b503HistoryBranchWant.B503HistoryBranchFields = b503HistoryBranchInput
				b503HistoryBranchWant.B503HistoryBranchComponents[b503HistoryBranchIndex].B503HistoryBranchCell = b503HistoryBranchInput.B503HistoryBranchCalendar[b503HistoryBranchIndex]
				b503HistoryBranchCandidate := []bool{b503HistoryBranchValue <= 59, b503HistoryBranchValue <= 23, b503HistoryBranchValue >= 1 && b503HistoryBranchValue <= 31, b503HistoryBranchValue >= 1 && b503HistoryBranchValue <= 12, b503HistoryBranchValue <= 99}[b503HistoryBranchIndex]
				if !b503HistoryBranchCandidate {
					b503HistoryBranchWant.B503HistoryBranchComplete = false
					b503HistoryBranchWant.B503HistoryBranchComponents[b503HistoryBranchIndex].B503HistoryBranchState = "out-of-domain"
				}
				b503HistoryBranchEqual(b503HistoryBranchT, B503HistoryBranchHistory(b503HistoryBranchInput), b503HistoryBranchWant)
			}
			for b503HistoryBranchMode := 0; b503HistoryBranchMode < 4; b503HistoryBranchMode++ {
				b503HistoryBranchInput := b503HistoryBranchCase.B503HistoryBranchInput
				b503HistoryBranchWant := b503HistoryBranchCase.B503HistoryBranchExpected
				switch b503HistoryBranchMode {
				case 0:
					b503HistoryBranchInput.B503HistoryBranchCalendar[b503HistoryBranchIndex].B503HistoryBranchRawPresent = false
				case 1:
					b503HistoryBranchInput.B503HistoryBranchCalendar[b503HistoryBranchIndex].B503HistoryBranchNormalizedPresent = false
				case 2:
					b503HistoryBranchInput.B503HistoryBranchCalendar[b503HistoryBranchIndex].B503HistoryBranchRaw = 255
				case 3:
					b503HistoryBranchInput.B503HistoryBranchCalendar[b503HistoryBranchIndex].B503HistoryBranchRaw = 42
				}
				b503HistoryBranchWant.B503HistoryBranchFields = b503HistoryBranchInput
				b503HistoryBranchWant.B503HistoryBranchComponents[b503HistoryBranchIndex].B503HistoryBranchCell = b503HistoryBranchInput.B503HistoryBranchCalendar[b503HistoryBranchIndex]
				if b503HistoryBranchMode < 2 {
					b503HistoryBranchWant.B503HistoryBranchState, b503HistoryBranchWant.B503HistoryBranchMeaning = "unknown-branch", "Unknown"
					b503HistoryBranchWant.B503HistoryBranchComponents[b503HistoryBranchIndex].B503HistoryBranchState = "missing"
				} else if b503HistoryBranchMode == 2 {
					b503HistoryBranchWant.B503HistoryBranchComponents[b503HistoryBranchIndex].B503HistoryBranchState = "replacement"
				}
				if b503HistoryBranchMode != 3 {
					b503HistoryBranchWant.B503HistoryBranchComplete = false
				}
				b503HistoryBranchEqual(b503HistoryBranchT, B503HistoryBranchHistory(b503HistoryBranchInput), b503HistoryBranchWant)
			}
		})
	}
}

func TestAb05An05ErrorPresenceAndValue(b503HistoryBranchT *b503HistoryBranchTesting.T) {
	for _, b503HistoryBranchCaseIndex := range []int{0, 2, 3, 5, 6} {
		b503HistoryBranchCase := b503HistoryBranchVectors(b503HistoryBranchT)[b503HistoryBranchCaseIndex]
		for _, b503HistoryBranchPresent := range []bool{false, true} {
			for _, b503HistoryBranchValue := range []uint64{0, 1, 4660, ^uint64(0)} {
				b503HistoryBranchInput := b503HistoryBranchCase.B503HistoryBranchInput
				b503HistoryBranchInput.B503HistoryBranchErrorPresent, b503HistoryBranchInput.B503HistoryBranchError = b503HistoryBranchPresent, b503HistoryBranchValue
				b503HistoryBranchWant := b503HistoryBranchCase.B503HistoryBranchExpected
				b503HistoryBranchWant.B503HistoryBranchFields = b503HistoryBranchInput
				if b503HistoryBranchCaseIndex == 2 && !b503HistoryBranchPresent {
					b503HistoryBranchWant.B503HistoryBranchState, b503HistoryBranchWant.B503HistoryBranchMeaning = "unknown-branch", "Unknown"
				}
				b503HistoryBranchEqual(b503HistoryBranchT, B503HistoryBranchHistory(b503HistoryBranchInput), b503HistoryBranchWant)
			}
		}
	}
}

func TestAb05An05ValueIsolationAndNoTimestamp(b503HistoryBranchT *b503HistoryBranchTesting.T) {
	b503HistoryBranchCase := b503HistoryBranchVectors(b503HistoryBranchT)[3]
	b503HistoryBranchInput := b503HistoryBranchCase.B503HistoryBranchInput
	b503HistoryBranchGot := B503HistoryBranchHistory(b503HistoryBranchInput)
	b503HistoryBranchInput.B503HistoryBranchCalendar[4].B503HistoryBranchValue = 2000
	b503HistoryBranchInput.B503HistoryBranchRemainder[0] = 0
	b503HistoryBranchInput.B503HistoryBranchContext.B503HistoryBranchFamily = "ab05an05_changed"
	b503HistoryBranchEqual(b503HistoryBranchT, b503HistoryBranchGot, b503HistoryBranchCase.B503HistoryBranchExpected)
	b503HistoryBranchGot.B503HistoryBranchFields.B503HistoryBranchCalendar[0].B503HistoryBranchRaw = 255
	b503HistoryBranchGot.B503HistoryBranchComponents[4].B503HistoryBranchCell.B503HistoryBranchValue = 2099
	b503HistoryBranchEqual(b503HistoryBranchT, B503HistoryBranchHistory(b503HistoryBranchCase.B503HistoryBranchInput), b503HistoryBranchCase.B503HistoryBranchExpected)
	if b503HistoryBranchReflect.TypeOf(b503HistoryBranchGot).NumField() != 7 || b503HistoryBranchReflect.TypeOf(B503HistoryBranchCell{}).NumField() != 4 {
		b503HistoryBranchT.Fatal("ab05an05 field contract expanded")
	}
	b503HistoryBranchInput = b503HistoryBranchCase.B503HistoryBranchInput
	b503HistoryBranchInput.B503HistoryBranchCalendar[3].B503HistoryBranchValue = 2
	b503HistoryBranchGot = B503HistoryBranchHistory(b503HistoryBranchInput)
	if !b503HistoryBranchGot.B503HistoryBranchComplete || b503HistoryBranchGot.B503HistoryBranchMeaning != "Hypothesis" || b503HistoryBranchGot.B503HistoryBranchQualification != "Unknown" {
		b503HistoryBranchT.Fatal("ab05an05 joint calendar or qualification inferred")
	}
}

func TestAb05An05CodeOnlyMissingAndReplacementPrecedence(b503HistoryBranchT *b503HistoryBranchTesting.T) {
	b503HistoryBranchCase := b503HistoryBranchVectors(b503HistoryBranchT)[2]
	for b503HistoryBranchIndex := 0; b503HistoryBranchIndex < 5; b503HistoryBranchIndex++ {
		for b503HistoryBranchMode := 0; b503HistoryBranchMode < 3; b503HistoryBranchMode++ {
			b503HistoryBranchInput := b503HistoryBranchCase.B503HistoryBranchInput
			b503HistoryBranchWant := b503HistoryBranchCase.B503HistoryBranchExpected
			switch b503HistoryBranchMode {
			case 0:
				b503HistoryBranchInput.B503HistoryBranchCalendar[b503HistoryBranchIndex].B503HistoryBranchRawPresent = false
			case 1:
				b503HistoryBranchInput.B503HistoryBranchCalendar[b503HistoryBranchIndex].B503HistoryBranchNormalizedPresent = false
			case 2:
				b503HistoryBranchInput.B503HistoryBranchCalendar[b503HistoryBranchIndex].B503HistoryBranchRaw = 42
				b503HistoryBranchInput.B503HistoryBranchCalendar[b503HistoryBranchIndex].B503HistoryBranchValue = 1
			}
			b503HistoryBranchWant.B503HistoryBranchFields = b503HistoryBranchInput
			b503HistoryBranchWant.B503HistoryBranchComponents[b503HistoryBranchIndex].B503HistoryBranchCell = b503HistoryBranchInput.B503HistoryBranchCalendar[b503HistoryBranchIndex]
			if b503HistoryBranchMode == 0 {
				b503HistoryBranchWant.B503HistoryBranchComponents[b503HistoryBranchIndex].B503HistoryBranchState = "missing"
			}
			if b503HistoryBranchMode == 2 {
				b503HistoryBranchWant.B503HistoryBranchComponents[b503HistoryBranchIndex].B503HistoryBranchState = "unknown"
			}
			b503HistoryBranchEqual(b503HistoryBranchT, B503HistoryBranchHistory(b503HistoryBranchInput), b503HistoryBranchWant)
		}
	}
}
