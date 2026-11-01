package system

import (
	ab09an09hex "encoding/hex"
	ab09an09json "encoding/json"
	ab09an09os "os"
	ab09an09reflect "reflect"
	ab09an09testing "testing"
)

type ab09an09Vector struct {
	Ab09An09Name      string
	Ab09An09Identity  Ab09An09Selection
	Ab09An09Selection Ab09An09Selection
	Ab09An09Selected  bool
	Ab09An09Raw       []*string
	Ab09An09States    []string
	Ab09An09Values    []*float64
}

func ab09an09Bytes(ab09an09T *ab09an09testing.T, ab09an09Hex *string) []byte {
	ab09an09T.Helper()
	if ab09an09Hex == nil {
		return nil
	}
	ab09an09Bytes, ab09an09Err := ab09an09hex.DecodeString(*ab09an09Hex)
	if ab09an09Err != nil {
		ab09an09T.Fatal(ab09an09Err)
	}
	return ab09an09Bytes
}

func ab09an09Candidates(ab09an09View Ab09An09QuantityView) []Ab09An09Candidate {
	return []Ab09An09Candidate{ab09an09View.Ab09An09Temperature, ab09an09View.Ab09An09Flow, ab09an09View.Ab09An09Volume, ab09an09View.Ab09An09MinimumVolume, ab09an09View.Ab09An09MaximumVolume}
}

func TestAb09An09IndependentVectors(ab09an09T *ab09an09testing.T) {
	ab09an09Data, ab09an09Err := ab09an09os.ReadFile("testdata/ab09an09_b509_quantity_views.json")
	if ab09an09Err != nil {
		ab09an09T.Fatal(ab09an09Err)
	}
	var ab09an09Vectors []ab09an09Vector
	if ab09an09Err = ab09an09json.Unmarshal(ab09an09Data, &ab09an09Vectors); ab09an09Err != nil {
		ab09an09T.Fatal(ab09an09Err)
	}
	if len(ab09an09Vectors) != 41 {
		ab09an09T.Fatalf("ab09an09 vector count = %d", len(ab09an09Vectors))
	}
	for _, ab09an09Case := range ab09an09Vectors {
		ab09an09T.Run(ab09an09Case.Ab09An09Name, func(ab09an09T *ab09an09testing.T) {
			if len(ab09an09Case.Ab09An09Raw) != 5 || len(ab09an09Case.Ab09An09States) != 5 || len(ab09an09Case.Ab09An09Values) != 5 {
				ab09an09T.Fatal("ab09an09 incomplete oracle")
			}
			ab09an09Input := Ab09An09Fields{
				ab09an09Bytes(ab09an09T, ab09an09Case.Ab09An09Raw[0]),
				ab09an09Bytes(ab09an09T, ab09an09Case.Ab09An09Raw[1]),
				ab09an09Bytes(ab09an09T, ab09an09Case.Ab09An09Raw[2]),
				ab09an09Bytes(ab09an09T, ab09an09Case.Ab09An09Raw[3]),
				ab09an09Bytes(ab09an09T, ab09an09Case.Ab09An09Raw[4]),
			}
			ab09an09Before := Ab09An09Fields{ab09an09Bytes(ab09an09T, ab09an09Case.Ab09An09Raw[0]), ab09an09Bytes(ab09an09T, ab09an09Case.Ab09An09Raw[1]), ab09an09Bytes(ab09an09T, ab09an09Case.Ab09An09Raw[2]), ab09an09Bytes(ab09an09T, ab09an09Case.Ab09An09Raw[3]), ab09an09Bytes(ab09an09T, ab09an09Case.Ab09An09Raw[4])}
			ab09an09View := Ab09An09B509QuantityViews(ab09an09Case.Ab09An09Identity, ab09an09Case.Ab09An09Selection, ab09an09Input)
			if !ab09an09reflect.DeepEqual(ab09an09Input, ab09an09Before) {
				ab09an09T.Fatal("ab09an09 input mutated")
			}
			if ab09an09View.Ab09An09Selected != ab09an09Case.Ab09An09Selected || ab09an09View.Ab09An09Identity != ab09an09Case.Ab09An09Identity || ab09an09View.Ab09An09Selection != ab09an09Case.Ab09An09Selection || ab09an09View.Ab09An09ContractGrade != "Proven" {
				ab09an09T.Fatalf("ab09an09 selection/contract = %+v", ab09an09View)
			}
			for ab09an09Index, ab09an09Got := range ab09an09Candidates(ab09an09View) {
				ab09an09Grade := "Hypothesis"
				if ab09an09Index >= 3 {
					ab09an09Grade = "Unknown"
				}
				if !ab09an09reflect.DeepEqual(ab09an09Got.Ab09An09Raw, ab09an09Bytes(ab09an09T, ab09an09Case.Ab09An09Raw[ab09an09Index])) || ab09an09Got.Ab09An09State != ab09an09Case.Ab09An09States[ab09an09Index] || !ab09an09reflect.DeepEqual(ab09an09Got.Ab09An09Value, ab09an09Case.Ab09An09Values[ab09an09Index]) || ab09an09Got.Ab09An09SemanticGrade != ab09an09Grade || ab09an09Got.Ab09An09NativeQualification != "Unknown" || ab09an09Got.Ab09An09Unit != []string{"degrees Celsius", "litres/minute", "litres", "", ""}[ab09an09Index] {
					ab09an09T.Errorf("ab09an09 field %d = %+v; expected state %s value %v", ab09an09Index, ab09an09Got, ab09an09Case.Ab09An09States[ab09an09Index], ab09an09Case.Ab09An09Values[ab09an09Index])
				}
			}
		})
	}
}

func TestAb09An09Ownership(ab09an09T *ab09an09testing.T) {
	ab09an09Identity := Ab09An09Selection{"ab09an09_synthetic_family", "ab09an09_synthetic_variant", 0x59, "ab09an09_response"}
	ab09an09Shared := []byte{0x10, 0}
	ab09an09Input := Ab09An09Fields{ab09an09Shared, ab09an09Shared, ab09an09Shared, ab09an09Shared, ab09an09Shared}
	ab09an09First := Ab09An09B509QuantityViews(ab09an09Identity, ab09an09Identity, ab09an09Input)
	ab09an09Second := Ab09An09B509QuantityViews(ab09an09Identity, ab09an09Identity, ab09an09Input)
	ab09an09Shared[0] = 0xff
	for _, ab09an09Got := range ab09an09Candidates(ab09an09First) {
		if !ab09an09reflect.DeepEqual(ab09an09Got.Ab09An09Raw, []byte{0x10, 0}) {
			ab09an09T.Fatal("ab09an09 input aliases result")
		}
	}
	for ab09an09Index, ab09an09Got := range ab09an09Candidates(ab09an09First) {
		ab09an09Got.Ab09An09Raw[0] = byte(ab09an09Index)
		if ab09an09Got.Ab09An09Value != nil {
			*ab09an09Got.Ab09An09Value = -999
		}
		for ab09an09OtherIndex, ab09an09Other := range ab09an09Candidates(ab09an09First) {
			if ab09an09OtherIndex > ab09an09Index && ab09an09Other.Ab09An09Raw[0] != 0x10 {
				ab09an09T.Fatal("ab09an09 fields alias")
			}
		}
	}
	for ab09an09Index, ab09an09Got := range ab09an09Candidates(ab09an09Second) {
		if !ab09an09reflect.DeepEqual(ab09an09Got.Ab09An09Raw, []byte{0x10, 0}) {
			ab09an09T.Fatal("ab09an09 results alias")
		}
		if ab09an09Index < 3 && (ab09an09Got.Ab09An09Value == nil || *ab09an09Got.Ab09An09Value != []float64{1, 0.16, 1.6}[ab09an09Index]) {
			ab09an09T.Fatal("ab09an09 scalar results alias")
		}
	}
	if ab09an09Shared[0] != 0xff {
		ab09an09T.Fatal("ab09an09 returned raw aliases input")
	}
}
