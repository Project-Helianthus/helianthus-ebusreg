package vaillant

import (
	ab19an19Hex "encoding/hex"
	ab19an19JSON "encoding/json"
	ab19an19OS "os"
	ab19an19Reflect "reflect"
	ab19an19Testing "testing"
)

func ab19an19Raw(ab19an19T *ab19an19Testing.T, ab19an19Text string) []byte {
	ab19an19T.Helper()
	ab19an19Bytes, ab19an19Error := ab19an19Hex.DecodeString(ab19an19Text)
	if ab19an19Error != nil {
		ab19an19T.Fatal(ab19an19Error)
	}
	return ab19an19Bytes
}

func TestAb19An19Vectors(ab19an19T *ab19an19Testing.T) {
	var ab19an19Fixture struct {
		Ab19An19Limits []struct {
			Ab19An19Label                              string
			Ab19An19MinimumRaw, Ab19An19MaximumRaw     string
			Ab19An19MinimumValue, Ab19An19MaximumValue uint32
		}
		Ab19An19Getters []struct {
			Ab19An19Label    string
			Ab19An19Raw      uint32
			Ab19An19Expected uint16
		}
	}
	ab19an19Bytes, ab19an19Error := ab19an19OS.ReadFile("testdata/ab19an19_b531_limit_numeric_views.json")
	if ab19an19Error != nil {
		ab19an19T.Fatal(ab19an19Error)
	}
	if ab19an19Error = ab19an19JSON.Unmarshal(ab19an19Bytes, &ab19an19Fixture); ab19an19Error != nil {
		ab19an19T.Fatal(ab19an19Error)
	}
	if len(ab19an19Fixture.Ab19An19Limits) != 6 || len(ab19an19Fixture.Ab19An19Getters) != 5 {
		ab19an19T.Fatal("ab19an19 incomplete vectors")
	}
	for _, ab19an19Vector := range ab19an19Fixture.Ab19An19Limits {
		ab19an19T.Run(ab19an19Vector.Ab19An19Label, func(ab19an19T *ab19an19Testing.T) {
			ab19an19Minimum := ab19an19Raw(ab19an19T, ab19an19Vector.Ab19An19MinimumRaw)
			ab19an19Maximum := ab19an19Raw(ab19an19T, ab19an19Vector.Ab19An19MaximumRaw)
			ab19an19Result := Ab19An19Interpret(Ab19An19Family, Ab19An19Variant, 2, Ab19An19Direction, ab19an19Minimum, ab19an19Maximum, true, true, Ab19An19Pair{})
			ab19an19Expected := Ab19An19Pair{
				Ab19An19Minimum: Ab19An19Field{Ab19An19Raw: ab19an19Minimum, Ab19An19State: "success", Ab19An19Available: true, Ab19An19Value: ab19an19Vector.Ab19An19MinimumValue, Ab19An19Fresh: true},
				Ab19An19Maximum: Ab19An19Field{Ab19An19Raw: ab19an19Maximum, Ab19An19State: "success", Ab19An19Available: true, Ab19An19Value: ab19an19Vector.Ab19An19MaximumValue, Ab19An19Fresh: true}, Ab19An19Fresh: true,
			}
			if !ab19an19Result.Ab19An19Admitted || !ab19an19Reflect.DeepEqual(ab19an19Result.Ab19An19Attempt, ab19an19Expected) || !ab19an19Reflect.DeepEqual(ab19an19Result.Ab19An19LastKnown, ab19an19Expected) {
				ab19an19T.Fatalf("ab19an19 result: %+v", ab19an19Result)
			}
			if ab19an19Result.Ab19An19ContractGrade != "Proven" || ab19an19Result.Ab19An19InterpretationGrade != "Hypothesis" || ab19an19Result.Ab19An19NativeUnit != "Unknown" || ab19an19Result.Ab19An19Qualification != "Unknown" || ab19an19Result.Ab19An19Presence != "Unknown" || ab19an19Result.Ab19An19Execution != "Unknown" {
				ab19an19T.Fatal("ab19an19 grade conflation")
			}
		})
	}
	for _, ab19an19Vector := range ab19an19Fixture.Ab19An19Getters {
		ab19an19T.Run(ab19an19Vector.Ab19An19Label, func(ab19an19T *ab19an19Testing.T) {
			if ab19an19Actual := Ab19An19Getter(ab19an19Vector.Ab19An19Raw); ab19an19Actual != ab19an19Vector.Ab19An19Expected {
				ab19an19T.Fatalf("ab19an19 got %d want %d", ab19an19Actual, ab19an19Vector.Ab19An19Expected)
			}
		})
	}
}

func TestAb19An19Admission(ab19an19T *ab19an19Testing.T) {
	for _, ab19an19Case := range []struct {
		Ab19An19Label, Ab19An19Family, Ab19An19Variant, Ab19An19Direction string
		Ab19An19Operation                                                 byte
	}{
		{"ab19an19-family", "other-family", Ab19An19Variant, Ab19An19Direction, 2},
		{"ab19an19-variant", Ab19An19Family, "other-variant", Ab19An19Direction, 2},
		{"ab19an19-direction", Ab19An19Family, Ab19An19Variant, "request", 2},
		{"ab19an19-operation", Ab19An19Family, Ab19An19Variant, Ab19An19Direction, 3},
	} {
		ab19an19T.Run(ab19an19Case.Ab19An19Label, func(ab19an19T *ab19an19Testing.T) {
			ab19an19Result := Ab19An19Interpret(ab19an19Case.Ab19An19Family, ab19an19Case.Ab19An19Variant, ab19an19Case.Ab19An19Operation, ab19an19Case.Ab19An19Direction, []byte{1, 0}, []byte{2, 0}, true, true, Ab19An19Pair{})
			if ab19an19Result.Ab19An19Admitted || ab19an19Result.Ab19An19Attempt.Ab19An19Fresh || ab19an19Result.Ab19An19LastKnown.Ab19An19Fresh || ab19an19Result.Ab19An19Attempt.Ab19An19Minimum.Ab19An19Available || ab19an19Result.Ab19An19Attempt.Ab19An19Maximum.Ab19An19Available || ab19an19Result.Ab19An19Attempt.Ab19An19Minimum.Ab19An19State != "not-admitted" || ab19an19Result.Ab19An19Attempt.Ab19An19Maximum.Ab19An19State != "not-admitted" {
				ab19an19T.Fatalf("ab19an19 admission: %+v", ab19an19Result)
			}
		})
	}
}

func TestAb19An19Retention(ab19an19T *ab19an19Testing.T) {
	for _, ab19an19Side := range []string{"ab19an19-minimum", "ab19an19-maximum"} {
		for _, ab19an19Case := range []struct {
			Ab19An19Label, Ab19An19Raw, Ab19An19State string
			Ab19An19ReadOK                            bool
		}{
			{"ab19an19-failed", "0100", "read-failed", false},
			{"ab19an19-missing", "", "missing", true},
			{"ab19an19-short", "01", "malformed", true},
			{"ab19an19-long", "010000", "malformed", true},
		} {
			for _, ab19an19Prior := range []Ab19An19Field{
				{},
				{Ab19An19Raw: []byte{3, 0}, Ab19An19State: "success", Ab19An19Available: true, Ab19An19Value: 30, Ab19An19Fresh: true},
				{Ab19An19Raw: []byte{4, 0}, Ab19An19State: "success", Ab19An19Available: true, Ab19An19Value: 40},
				{Ab19An19Raw: []byte{9}, Ab19An19State: "malformed", Ab19An19Value: 77},
			} {
				ab19an19T.Run(ab19an19Side+"/"+ab19an19Case.Ab19An19Label, func(ab19an19T *ab19an19Testing.T) {
					ab19an19Pair := Ab19An19Pair{Ab19An19Minimum: ab19an19Prior, Ab19An19Maximum: ab19an19Prior, Ab19An19Fresh: true}
					ab19an19BadRaw := ab19an19Raw(ab19an19T, ab19an19Case.Ab19An19Raw)
					ab19an19Min, ab19an19Max := []byte{5, 0}, []byte{5, 0}
					ab19an19MinOK, ab19an19MaxOK := true, true
					if ab19an19Side == "ab19an19-minimum" {
						ab19an19Min, ab19an19MinOK = ab19an19BadRaw, ab19an19Case.Ab19An19ReadOK
					} else {
						ab19an19Max, ab19an19MaxOK = ab19an19BadRaw, ab19an19Case.Ab19An19ReadOK
					}
					ab19an19Result := Ab19An19Interpret(Ab19An19Family, Ab19An19Variant, 2, Ab19An19Direction, ab19an19Min, ab19an19Max, ab19an19MinOK, ab19an19MaxOK, ab19an19Pair)
					ab19an19Failed, ab19an19Retained, ab19an19Good := ab19an19Result.Ab19An19Attempt.Ab19An19Minimum, ab19an19Result.Ab19An19LastKnown.Ab19An19Minimum, ab19an19Result.Ab19An19LastKnown.Ab19An19Maximum
					if ab19an19Side == "ab19an19-maximum" {
						ab19an19Failed, ab19an19Retained, ab19an19Good = ab19an19Result.Ab19An19Attempt.Ab19An19Maximum, ab19an19Result.Ab19An19LastKnown.Ab19An19Maximum, ab19an19Result.Ab19An19LastKnown.Ab19An19Minimum
					}
					ab19an19ExpectedFailed := Ab19An19Field{Ab19An19Raw: append([]byte(nil), ab19an19BadRaw...), Ab19An19State: ab19an19Case.Ab19An19State}
					ab19an19ExpectedRetained := ab19an19Prior
					ab19an19ExpectedRetained.Ab19An19Fresh = false
					ab19an19ExpectedGood := Ab19An19Field{Ab19An19Raw: []byte{5, 0}, Ab19An19State: "success", Ab19An19Available: true, Ab19An19Value: 50, Ab19An19Fresh: true}
					if !ab19an19Reflect.DeepEqual(ab19an19Failed, ab19an19ExpectedFailed) || !ab19an19Reflect.DeepEqual(ab19an19Retained, ab19an19ExpectedRetained) || !ab19an19Reflect.DeepEqual(ab19an19Good, ab19an19ExpectedGood) || !ab19an19Reflect.DeepEqual(ab19an19Result.Ab19An19Prior, ab19an19Pair) || ab19an19Result.Ab19An19Attempt.Ab19An19Fresh || ab19an19Result.Ab19An19LastKnown.Ab19An19Fresh {
						ab19an19T.Fatalf("ab19an19 partial: %+v", ab19an19Result)
					}
				})
			}
		}
	}
}

func TestAb19An19PriorDimensions(ab19an19T *ab19an19Testing.T) {
	for _, ab19an19Case := range []struct {
		Ab19An19Label string
		Ab19An19Field Ab19An19Field
	}{
		{"ab19an19-baseline", Ab19An19Field{[]byte{3, 0}, "success", true, 30, true}},
		{"ab19an19-raw", Ab19An19Field{[]byte{4, 0}, "success", true, 30, true}},
		{"ab19an19-state", Ab19An19Field{[]byte{3, 0}, "malformed", true, 30, true}},
		{"ab19an19-available", Ab19An19Field{[]byte{3, 0}, "success", false, 30, true}},
		{"ab19an19-value", Ab19An19Field{[]byte{3, 0}, "success", true, 40, true}},
		{"ab19an19-fresh", Ab19An19Field{[]byte{3, 0}, "success", true, 30, false}},
	} {
		for _, ab19an19Side := range []string{"ab19an19-minimum", "ab19an19-maximum"} {
			for _, ab19an19Fresh := range []bool{false, true} {
				ab19an19T.Run(ab19an19Case.Ab19An19Label+"/"+ab19an19Side, func(ab19an19T *ab19an19Testing.T) {
					ab19an19Prior := Ab19An19Pair{Ab19An19Minimum: Ab19An19Field{[]byte{3, 0}, "success", true, 30, true}, Ab19An19Maximum: Ab19An19Field{[]byte{3, 0}, "success", true, 30, true}, Ab19An19Fresh: ab19an19Fresh}
					if ab19an19Side == "ab19an19-minimum" {
						ab19an19Prior.Ab19An19Minimum = ab19an19Case.Ab19An19Field
					} else {
						ab19an19Prior.Ab19An19Maximum = ab19an19Case.Ab19An19Field
					}
					ab19an19Expected := ab19an19Prior
					ab19an19Expected.Ab19An19Fresh = false
					ab19an19Expected.Ab19An19Minimum.Ab19An19Fresh = false
					ab19an19Expected.Ab19An19Maximum.Ab19An19Fresh = false
					ab19an19Result := Ab19An19Interpret(Ab19An19Family, Ab19An19Variant, 2, Ab19An19Direction, nil, nil, false, false, ab19an19Prior)
					if !ab19an19Reflect.DeepEqual(ab19an19Result.Ab19An19Prior, ab19an19Prior) || !ab19an19Reflect.DeepEqual(ab19an19Result.Ab19An19LastKnown, ab19an19Expected) {
						ab19an19T.Fatalf("ab19an19 prior dimensions: %+v", ab19an19Result)
					}
				})
			}
		}
	}
}

func TestAb19An19TotalFailureAndOwnership(ab19an19T *ab19an19Testing.T) {
	ab19an19Min, ab19an19Max := []byte{1, 0}, []byte{2, 0}
	ab19an19Prior := Ab19An19Pair{Ab19An19Minimum: Ab19An19Field{Ab19An19Raw: []byte{3, 0}, Ab19An19State: "success", Ab19An19Available: true, Ab19An19Value: 30, Ab19An19Fresh: true}, Ab19An19Maximum: Ab19An19Field{Ab19An19Raw: []byte{4, 0}, Ab19An19State: "success", Ab19An19Available: true, Ab19An19Value: 40, Ab19An19Fresh: true}, Ab19An19Fresh: true}
	ab19an19First := Ab19An19Interpret(Ab19An19Family, Ab19An19Variant, 2, Ab19An19Direction, ab19an19Min, ab19an19Max, false, false, ab19an19Prior)
	ab19an19Second := Ab19An19Interpret(Ab19An19Family, Ab19An19Variant, 2, Ab19An19Direction, ab19an19Min, ab19an19Max, false, false, ab19an19Prior)
	if ab19an19First.Ab19An19LastKnown.Ab19An19Minimum.Ab19An19Value != 30 || ab19an19First.Ab19An19LastKnown.Ab19An19Maximum.Ab19An19Value != 40 || ab19an19First.Ab19An19LastKnown.Ab19An19Fresh || ab19an19First.Ab19An19LastKnown.Ab19An19Minimum.Ab19An19Fresh || ab19an19First.Ab19An19LastKnown.Ab19An19Maximum.Ab19An19Fresh || ab19an19First.Ab19An19Attempt.Ab19An19Minimum.Ab19An19State != "read-failed" || ab19an19First.Ab19An19Attempt.Ab19An19Maximum.Ab19An19State != "read-failed" {
		ab19an19T.Fatalf("ab19an19 total failure: %+v", ab19an19First)
	}
	ab19an19Min[0], ab19an19Max[0] = 99, 99
	ab19an19Prior.Ab19An19Minimum.Ab19An19Raw[0], ab19an19Prior.Ab19An19Maximum.Ab19An19Raw[0] = 99, 99
	ab19an19First.Ab19An19Attempt.Ab19An19Minimum.Ab19An19Raw[0] = 88
	ab19an19First.Ab19An19Attempt.Ab19An19Maximum.Ab19An19Raw[0] = 88
	ab19an19First.Ab19An19Prior.Ab19An19Minimum.Ab19An19Raw[0] = 88
	ab19an19First.Ab19An19Prior.Ab19An19Maximum.Ab19An19Raw[0] = 88
	if ab19an19Second.Ab19An19Attempt.Ab19An19Minimum.Ab19An19Raw[0] != 1 || ab19an19Second.Ab19An19Attempt.Ab19An19Maximum.Ab19An19Raw[0] != 2 || ab19an19Second.Ab19An19Prior.Ab19An19Minimum.Ab19An19Raw[0] != 3 || ab19an19Second.Ab19An19Prior.Ab19An19Maximum.Ab19An19Raw[0] != 4 || ab19an19First.Ab19An19LastKnown.Ab19An19Minimum.Ab19An19Raw[0] != 3 || ab19an19First.Ab19An19LastKnown.Ab19An19Maximum.Ab19An19Raw[0] != 4 {
		ab19an19T.Fatal("ab19an19 raw buffers alias")
	}
	ab19an19First = Ab19An19Interpret(Ab19An19Family, Ab19An19Variant, 2, Ab19An19Direction, []byte{1, 0}, []byte{2, 0}, true, true, Ab19An19Pair{})
	ab19an19First.Ab19An19Attempt.Ab19An19Minimum.Ab19An19Raw[0], ab19an19First.Ab19An19Attempt.Ab19An19Maximum.Ab19An19Raw[0] = 88, 88
	if ab19an19First.Ab19An19LastKnown.Ab19An19Minimum.Ab19An19Raw[0] != 1 || ab19an19First.Ab19An19LastKnown.Ab19An19Maximum.Ab19An19Raw[0] != 2 {
		ab19an19T.Fatal("ab19an19 successful buffers alias")
	}
}
