package system

import (
	ab07an07bytes "bytes"
	ab07an07json "encoding/json"
	ab07an07os "os"
	ab07an07reflect "reflect"
	ab07an07testing "testing"
)

type ab07an07Vector struct {
	Ab07An07Label    string
	Ab07An07Input    Ab07An07Input
	Ab07An07Values   [5]int16
	Ab07An07States   [5]string
	Ab07An07Selected bool
	Ab07An07Equal    bool
}

func TestAb07An07Vectors(ab07an07T *ab07an07testing.T) {
	ab07an07Data, ab07an07Err := ab07an07os.ReadFile("testdata/ab07an07_b509_constraint_word_views.json")
	if ab07an07Err != nil {
		ab07an07T.Fatal(ab07an07Err)
	}
	var ab07an07Vectors []ab07an07Vector
	if ab07an07Err = ab07an07json.Unmarshal(ab07an07Data, &ab07an07Vectors); ab07an07Err != nil {
		ab07an07T.Fatal(ab07an07Err)
	}
	if len(ab07an07Vectors) == 0 {
		ab07an07T.Fatal("ab07an07_no_vectors")
	}
	for _, ab07an07Vector := range ab07an07Vectors {
		ab07an07T.Run(ab07an07Vector.Ab07An07Label, func(ab07an07T *ab07an07testing.T) {
			ab07an07View := Ab07An07Interpret(ab07an07Vector.Ab07An07Input)
			if !ab07an07reflect.DeepEqual(ab07an07View.Ab07An07Context, ab07an07Vector.Ab07An07Input) {
				ab07an07T.Fatal("ab07an07_context_lost")
			}
			if ab07an07View.Ab07An07Selected != ab07an07Vector.Ab07An07Selected || ab07an07View.Ab07An07RawAddressEqual != ab07an07Vector.Ab07An07Equal {
				ab07an07T.Fatal("ab07an07_selection_or_local_equality")
			}
			if ab07an07View.Ab07An07RequestedEditBoundsAdmitted {
				ab07an07T.Fatal("ab07an07_candidate_cannot_admit_bounds")
			}
			ab07an07Words := []Ab07An07Word{ab07an07View.Ab07An07Minimum, ab07an07View.Ab07An07Maximum, ab07an07View.Ab07An07Default, ab07an07View.Ab07An07Step, ab07an07View.Ab07An07Factor}
			ab07an07Raw := [][]byte{ab07an07Vector.Ab07An07Input.Ab07An07Minimum, ab07an07Vector.Ab07An07Input.Ab07An07Maximum, ab07an07Vector.Ab07An07Input.Ab07An07Default, ab07an07Vector.Ab07An07Input.Ab07An07Step, ab07an07Vector.Ab07An07Input.Ab07An07Factor}
			for ab07an07Index, ab07an07Word := range ab07an07Words {
				if ab07an07Word.Ab07An07Value != int(ab07an07Vector.Ab07An07Values[ab07an07Index]) || ab07an07Word.Ab07An07State != ab07an07Vector.Ab07An07States[ab07an07Index] || !ab07an07reflect.DeepEqual(ab07an07Word.Ab07An07Raw, ab07an07Raw[ab07an07Index]) {
					ab07an07T.Fatalf("ab07an07_word_%d: %+v", ab07an07Index, ab07an07Word)
				}
			}
			ab07an07Grades := map[string]string{
				"Ab07An07SoftwareContract": "Proven", "Ab07An07Interpretation": "Hypothesis",
				"Ab07An07NativeEchoRules": "Unknown", "Ab07An07PhysicalQualification": "Unknown",
				"Ab07An07PacketAssociation": "Unknown", "Ab07An07FamilyApplicability": "Unknown",
				"Ab07An07FactorAlgorithm": "Unknown", "Ab07An07Writability": "Unknown",
				"Ab07An07PeerEnforcement": "Unknown", "Ab07An07Freshness": "Unknown", "Ab07An07Execution": "Unknown",
			}
			for ab07an07Name, ab07an07Expected := range ab07an07Grades {
				if ab07an07reflect.ValueOf(ab07an07View).FieldByName(ab07an07Name).String() != ab07an07Expected {
					ab07an07T.Fatal(ab07an07Name)
				}
			}
		})
	}
}

func TestAb07An07DetachedMemory(ab07an07T *ab07an07testing.T) {
	ab07an07Names := []string{"Ab07An07RequestedAddress", "Ab07An07ReturnedAddress", "Ab07An07Minimum", "Ab07An07Maximum", "Ab07An07Default", "Ab07An07Step", "Ab07An07Factor"}
	for _, ab07an07Name := range ab07an07Names {
		ab07an07T.Run(ab07an07Name, func(ab07an07T *ab07an07testing.T) {
			ab07an07Input := Ab07An07Input{Ab07An07Operation: 0x33, Ab07An07Direction: Ab07An07PositiveResponse}
			ab07an07Source := []byte{1, 2}
			ab07an07reflect.ValueOf(&ab07an07Input).Elem().FieldByName(ab07an07Name).SetBytes(ab07an07Source)
			ab07an07View := Ab07An07Interpret(ab07an07Input)
			ab07an07Other := Ab07An07Interpret(ab07an07Input)
			ab07an07ContextRaw := ab07an07reflect.ValueOf(ab07an07View.Ab07An07Context).FieldByName(ab07an07Name).Bytes()
			ab07an07Source[0] = 9
			if !ab07an07bytes.Equal(ab07an07ContextRaw, []byte{1, 2}) {
				ab07an07T.Fatal("ab07an07_input_alias")
			}
			ab07an07ContextRaw[1] = 8
			if ab07an07Source[1] != 2 || ab07an07reflect.ValueOf(ab07an07Other.Ab07An07Context).FieldByName(ab07an07Name).Bytes()[1] != 2 {
				ab07an07T.Fatal("ab07an07_result_alias")
			}
			if ab07an07Name != "Ab07An07RequestedAddress" && ab07an07Name != "Ab07An07ReturnedAddress" {
				ab07an07Word := ab07an07reflect.ValueOf(ab07an07View).FieldByName(ab07an07Name).Interface().(Ab07An07Word)
				if !ab07an07bytes.Equal(ab07an07Word.Ab07An07Raw, []byte{1, 2}) {
					ab07an07T.Fatal("ab07an07_word_alias")
				}
				ab07an07Word.Ab07An07Raw[0] = 7
				if ab07an07Source[0] != 9 || ab07an07ContextRaw[0] != 1 || ab07an07reflect.ValueOf(ab07an07Other).FieldByName(ab07an07Name).Interface().(Ab07An07Word).Ab07An07Raw[0] != 1 {
					ab07an07T.Fatal("ab07an07_word_output_alias")
				}
			}
		})
	}
}
