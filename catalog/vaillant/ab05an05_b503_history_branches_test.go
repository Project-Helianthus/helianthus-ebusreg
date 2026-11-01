package vaillant

import (
	ab05an05json "encoding/json"
	ab05an05os "os"
	ab05an05reflect "reflect"
	ab05an05testing "testing"
)

type ab05an05Vector struct {
	Ab05An05Label    string
	Ab05An05Input    Ab05An05Fields
	Ab05An05Expected Ab05An05View
}

func ab05an05Vectors(ab05an05T *ab05an05testing.T) []ab05an05Vector {
	ab05an05T.Helper()
	ab05an05Bytes, ab05an05Error := ab05an05os.ReadFile("testdata/ab05an05_b503_history_branches.json")
	if ab05an05Error != nil {
		ab05an05T.Fatal(ab05an05Error)
	}
	var ab05an05Cases []ab05an05Vector
	if ab05an05Error = ab05an05json.Unmarshal(ab05an05Bytes, &ab05an05Cases); ab05an05Error != nil {
		ab05an05T.Fatal(ab05an05Error)
	}
	if len(ab05an05Cases) != 7 {
		ab05an05T.Fatal("ab05an05 expected seven finite vectors")
	}
	return ab05an05Cases
}

func ab05an05Equal(ab05an05T *ab05an05testing.T, ab05an05Got, ab05an05Want Ab05An05View) {
	ab05an05T.Helper()
	if !ab05an05reflect.DeepEqual(ab05an05Got, ab05an05Want) {
		ab05an05T.Fatalf("ab05an05 got %#v want %#v", ab05an05Got, ab05an05Want)
	}
}

func TestAb05An05IndependentVectors(ab05an05T *ab05an05testing.T) {
	for _, ab05an05Case := range ab05an05Vectors(ab05an05T) {
		ab05an05T.Run(ab05an05Case.Ab05An05Label, func(ab05an05T *ab05an05testing.T) {
			ab05an05Equal(ab05an05T, Ab05An05History(ab05an05Case.Ab05An05Input), ab05an05Case.Ab05An05Expected)
		})
	}
}

func TestAb05An05AllRemainderBytes(ab05an05T *ab05an05testing.T) {
	ab05an05Case := ab05an05Vectors(ab05an05T)[0]
	for ab05an05Index := 0; ab05an05Index < 7; ab05an05Index++ {
		for ab05an05Value := 0; ab05an05Value < 255; ab05an05Value++ {
			ab05an05Input := ab05an05Case.Ab05An05Input
			ab05an05Input.Ab05An05Remainder[ab05an05Index] = uint8(ab05an05Value)
			ab05an05Want := ab05an05Case.Ab05An05Expected
			ab05an05Want.Ab05An05Fields = ab05an05Input
			ab05an05Want.Ab05An05State = "unknown-branch"
			ab05an05Want.Ab05An05Meaning = "Unknown"
			ab05an05Equal(ab05an05T, Ab05An05History(ab05an05Input), ab05an05Want)
		}
	}
	ab05an05Input := ab05an05Case.Ab05An05Input
	ab05an05Input.Ab05An05RemainderPresent = false
	if Ab05An05History(ab05an05Input).Ab05An05State != "unknown-branch" {
		ab05an05T.Fatal("ab05an05 absent remainder admitted")
	}
}

func TestAb05An05MarkerAndContext(ab05an05T *ab05an05testing.T) {
	ab05an05Case := ab05an05Vectors(ab05an05T)[3]
	for ab05an05Marker := 3; ab05an05Marker <= 255; ab05an05Marker++ {
		ab05an05Input := ab05an05Case.Ab05An05Input
		ab05an05Input.Ab05An05Marker = uint8(ab05an05Marker)
		ab05an05Want := ab05an05Case.Ab05An05Expected
		ab05an05Want.Ab05An05Fields = ab05an05Input
		ab05an05Want.Ab05An05State, ab05an05Want.Ab05An05Meaning, ab05an05Want.Ab05An05Complete = "unknown-branch", "Unknown", false
		for ab05an05Index := range ab05an05Want.Ab05An05Components {
			ab05an05Want.Ab05An05Components[ab05an05Index].Ab05An05State = "unknown"
		}
		ab05an05Equal(ab05an05T, Ab05An05History(ab05an05Input), ab05an05Want)
	}
	for _, ab05an05Field := range []string{"Ab05An05Family", "Ab05An05Variant", "Ab05An05Operation", "Ab05An05Direction"} {
		ab05an05T.Run("ab05an05_"+ab05an05Field, func(ab05an05T *ab05an05testing.T) {
			for ab05an05Mode := 0; ab05an05Mode < 2; ab05an05Mode++ {
				ab05an05Input := ab05an05Case.Ab05An05Input
				ab05an05Axis := ab05an05reflect.ValueOf(&ab05an05Input.Ab05An05Context).Elem().FieldByName(ab05an05Field)
				if ab05an05Axis.Kind() == ab05an05reflect.String {
					if ab05an05Mode == 0 {
						ab05an05Axis.SetString("")
					} else {
						ab05an05Axis.SetString("ab05an05_other")
					}
				} else {
					ab05an05Axis.SetUint(uint64(ab05an05Mode * 7))
				}
				ab05an05Want := ab05an05Case.Ab05An05Expected
				ab05an05Want.Ab05An05Fields = ab05an05Input
				if ab05an05Mode == 0 || ab05an05Field == "Ab05An05Operation" {
					ab05an05Want.Ab05An05State, ab05an05Want.Ab05An05Meaning, ab05an05Want.Ab05An05Complete = "unknown-branch", "Unknown", false
					for ab05an05Index := range ab05an05Want.Ab05An05Components {
						ab05an05Want.Ab05An05Components[ab05an05Index].Ab05An05State = "unknown"
					}
				}
				ab05an05Equal(ab05an05T, Ab05An05History(ab05an05Input), ab05an05Want)
			}
		})
	}
	ab05an05Input := ab05an05Case.Ab05An05Input
	ab05an05Input.Ab05An05MarkerPresent = false
	if Ab05An05History(ab05an05Input).Ab05An05State != "unknown-branch" {
		ab05an05T.Fatal("ab05an05 absent marker admitted")
	}
}

func TestAb05An05IndependentCells(ab05an05T *ab05an05testing.T) {
	ab05an05Case := ab05an05Vectors(ab05an05T)[3]
	for ab05an05Index, ab05an05Name := range []string{"ab05an05_minute", "ab05an05_hour", "ab05an05_day", "ab05an05_month", "ab05an05_year"} {
		ab05an05T.Run(ab05an05Name, func(ab05an05T *ab05an05testing.T) {
			for _, ab05an05Value := range []uint64{0, 1, 12, 13, 23, 24, 31, 32, 59, 60, 99, 100, 255, ^uint64(0)} {
				ab05an05Input := ab05an05Case.Ab05An05Input
				ab05an05Input.Ab05An05Calendar[ab05an05Index].Ab05An05Value = ab05an05Value
				ab05an05Want := ab05an05Case.Ab05An05Expected
				ab05an05Want.Ab05An05Fields = ab05an05Input
				ab05an05Want.Ab05An05Components[ab05an05Index].Ab05An05Cell = ab05an05Input.Ab05An05Calendar[ab05an05Index]
				ab05an05Candidate := []bool{ab05an05Value <= 59, ab05an05Value <= 23, ab05an05Value >= 1 && ab05an05Value <= 31, ab05an05Value >= 1 && ab05an05Value <= 12, ab05an05Value <= 99}[ab05an05Index]
				if !ab05an05Candidate {
					ab05an05Want.Ab05An05Complete = false
					ab05an05Want.Ab05An05Components[ab05an05Index].Ab05An05State = "out-of-domain"
				}
				ab05an05Equal(ab05an05T, Ab05An05History(ab05an05Input), ab05an05Want)
			}
			for ab05an05Mode := 0; ab05an05Mode < 4; ab05an05Mode++ {
				ab05an05Input := ab05an05Case.Ab05An05Input
				ab05an05Want := ab05an05Case.Ab05An05Expected
				switch ab05an05Mode {
				case 0:
					ab05an05Input.Ab05An05Calendar[ab05an05Index].Ab05An05RawPresent = false
				case 1:
					ab05an05Input.Ab05An05Calendar[ab05an05Index].Ab05An05NormalizedPresent = false
				case 2:
					ab05an05Input.Ab05An05Calendar[ab05an05Index].Ab05An05Raw = 255
				case 3:
					ab05an05Input.Ab05An05Calendar[ab05an05Index].Ab05An05Raw = 42
				}
				ab05an05Want.Ab05An05Fields = ab05an05Input
				ab05an05Want.Ab05An05Components[ab05an05Index].Ab05An05Cell = ab05an05Input.Ab05An05Calendar[ab05an05Index]
				if ab05an05Mode < 2 {
					ab05an05Want.Ab05An05State, ab05an05Want.Ab05An05Meaning = "unknown-branch", "Unknown"
					ab05an05Want.Ab05An05Components[ab05an05Index].Ab05An05State = "missing"
				} else if ab05an05Mode == 2 {
					ab05an05Want.Ab05An05Components[ab05an05Index].Ab05An05State = "replacement"
				}
				if ab05an05Mode != 3 {
					ab05an05Want.Ab05An05Complete = false
				}
				ab05an05Equal(ab05an05T, Ab05An05History(ab05an05Input), ab05an05Want)
			}
		})
	}
}

func TestAb05An05ErrorPresenceAndValue(ab05an05T *ab05an05testing.T) {
	for _, ab05an05CaseIndex := range []int{0, 2, 3, 5, 6} {
		ab05an05Case := ab05an05Vectors(ab05an05T)[ab05an05CaseIndex]
		for _, ab05an05Present := range []bool{false, true} {
			for _, ab05an05Value := range []uint64{0, 1, 4660, ^uint64(0)} {
				ab05an05Input := ab05an05Case.Ab05An05Input
				ab05an05Input.Ab05An05ErrorPresent, ab05an05Input.Ab05An05Error = ab05an05Present, ab05an05Value
				ab05an05Want := ab05an05Case.Ab05An05Expected
				ab05an05Want.Ab05An05Fields = ab05an05Input
				if ab05an05CaseIndex == 2 && !ab05an05Present {
					ab05an05Want.Ab05An05State, ab05an05Want.Ab05An05Meaning = "unknown-branch", "Unknown"
				}
				ab05an05Equal(ab05an05T, Ab05An05History(ab05an05Input), ab05an05Want)
			}
		}
	}
}

func TestAb05An05ValueIsolationAndNoTimestamp(ab05an05T *ab05an05testing.T) {
	ab05an05Case := ab05an05Vectors(ab05an05T)[3]
	ab05an05Input := ab05an05Case.Ab05An05Input
	ab05an05Got := Ab05An05History(ab05an05Input)
	ab05an05Input.Ab05An05Calendar[4].Ab05An05Value = 2000
	ab05an05Input.Ab05An05Remainder[0] = 0
	ab05an05Input.Ab05An05Context.Ab05An05Family = "ab05an05_changed"
	ab05an05Equal(ab05an05T, ab05an05Got, ab05an05Case.Ab05An05Expected)
	ab05an05Got.Ab05An05Fields.Ab05An05Calendar[0].Ab05An05Raw = 255
	ab05an05Got.Ab05An05Components[4].Ab05An05Cell.Ab05An05Value = 2099
	ab05an05Equal(ab05an05T, Ab05An05History(ab05an05Case.Ab05An05Input), ab05an05Case.Ab05An05Expected)
	if ab05an05reflect.TypeOf(ab05an05Got).NumField() != 7 || ab05an05reflect.TypeOf(Ab05An05Cell{}).NumField() != 4 {
		ab05an05T.Fatal("ab05an05 field contract expanded")
	}
	ab05an05Input = ab05an05Case.Ab05An05Input
	ab05an05Input.Ab05An05Calendar[3].Ab05An05Value = 2
	ab05an05Got = Ab05An05History(ab05an05Input)
	if !ab05an05Got.Ab05An05Complete || ab05an05Got.Ab05An05Meaning != "Hypothesis" || ab05an05Got.Ab05An05Qualification != "Unknown" {
		ab05an05T.Fatal("ab05an05 joint calendar or qualification inferred")
	}
}

func TestAb05An05CodeOnlyMissingAndReplacementPrecedence(ab05an05T *ab05an05testing.T) {
	ab05an05Case := ab05an05Vectors(ab05an05T)[2]
	for ab05an05Index := 0; ab05an05Index < 5; ab05an05Index++ {
		for ab05an05Mode := 0; ab05an05Mode < 3; ab05an05Mode++ {
			ab05an05Input := ab05an05Case.Ab05An05Input
			ab05an05Want := ab05an05Case.Ab05An05Expected
			switch ab05an05Mode {
			case 0:
				ab05an05Input.Ab05An05Calendar[ab05an05Index].Ab05An05RawPresent = false
			case 1:
				ab05an05Input.Ab05An05Calendar[ab05an05Index].Ab05An05NormalizedPresent = false
			case 2:
				ab05an05Input.Ab05An05Calendar[ab05an05Index].Ab05An05Raw = 42
				ab05an05Input.Ab05An05Calendar[ab05an05Index].Ab05An05Value = 1
			}
			ab05an05Want.Ab05An05Fields = ab05an05Input
			ab05an05Want.Ab05An05Components[ab05an05Index].Ab05An05Cell = ab05an05Input.Ab05An05Calendar[ab05an05Index]
			if ab05an05Mode == 0 {
				ab05an05Want.Ab05An05Components[ab05an05Index].Ab05An05State = "missing"
			}
			if ab05an05Mode == 2 {
				ab05an05Want.Ab05An05Components[ab05an05Index].Ab05An05State = "unknown"
			}
			ab05an05Equal(ab05an05T, Ab05An05History(ab05an05Input), ab05an05Want)
		}
	}
}
