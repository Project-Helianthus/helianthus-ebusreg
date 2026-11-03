package ebus_standard_catalog

import (
	ab76json "encoding/json"
	ab76fmt "fmt"
	ab76os "os"
	ab76reflect "reflect"
	ab76testing "testing"
)

func ab76Byte(ab76Value uint8) *uint8 { return &ab76Value }

func ab76Selection(ab76Operation uint8, ab76Selector, ab76Direction string) Ab76Selection {
	ab76Encoding := Ab76UnsignedU8
	if ab76Operation == 4 {
		ab76Encoding = Ab76DecimalChunks
	}
	return Ab76Selection{Ab76Family: "standard-counter", Ab76Variant: ab76fmt.Sprintf("03/%02d", ab76Operation), Ab76Function: ab76Byte(3), Ab76Operation: ab76Byte(ab76Operation), Ab76Selector: ab76Selector, Ab76Direction: ab76Direction, Ab76Encoding: ab76Encoding}
}

func ab76NoProjection(ab76T *ab76testing.T, ab76View Ab76View) {
	ab76T.Helper()
	if ab76View.Ab76LocalCounter != nil || ab76View.Ab76LocalHours != nil || ab76View.Ab76LocalMinute != nil {
		ab76T.Fatalf("unexpected projection: %+v", ab76View)
	}
}

func TestAb76Fixtures(ab76T *ab76testing.T) {
	ab76Data, ab76Err := ab76os.ReadFile("testdata/ab76_standard03_candidate_views.json")
	if ab76Err != nil {
		ab76T.Fatal(ab76Err)
	}
	var ab76Cases []struct {
		Ab76Label     string
		Ab76Operation uint8
		Ab76Body      []uint8
		Ab76Names     []string
		Ab76Values    []uint8
		Ab76Counter   *int
		Ab76Hours     *int
		Ab76Minute    *int
	}
	if ab76Err = ab76json.Unmarshal(ab76Data, &ab76Cases); ab76Err != nil {
		ab76T.Fatal(ab76Err)
	}
	if len(ab76Cases) != 3 {
		ab76T.Fatal("ab76 fixture inventory")
	}
	for _, ab76Case := range ab76Cases {
		ab76T.Run(ab76Case.Ab76Label, func(ab76T *ab76testing.T) {
			for _, ab76Mode := range []struct{ Ab76Selector, Ab76Direction string }{{"candidate", "request"}, {"standard", "positive"}} {
				if ab76Case.Ab76Operation == 5 && ab76Mode.Ab76Selector == "standard" {
					continue
				}
				ab76View := Ab76Interpret(ab76Selection(ab76Case.Ab76Operation, ab76Mode.Ab76Selector, ab76Mode.Ab76Direction), ab76Case.Ab76Body)
				if ab76View.Ab76State != Ab76Complete || !ab76reflect.DeepEqual(ab76View.Ab76Raw, ab76Case.Ab76Body) || len(ab76View.Ab76Components) != len(ab76Case.Ab76Names) || len(ab76View.Ab76Remainder) != 0 {
					ab76T.Fatalf("ab76 shape: %+v", ab76View)
				}
				for ab76Index, ab76Component := range ab76View.Ab76Components {
					if ab76Component.Ab76Name != ab76Case.Ab76Names[ab76Index] || ab76Component.Ab76State != Ab76Valid || ab76Component.Ab76Value == nil || *ab76Component.Ab76Value != ab76Case.Ab76Values[ab76Index] || ab76Component.Ab76Raw == nil || *ab76Component.Ab76Raw != ab76Case.Ab76Body[ab76Index] || ab76Component.Ab76Width != 1 || ab76Component.Ab76Encoding != ab76View.Ab76Selection.Ab76Encoding {
						ab76T.Fatalf("ab76 component: %+v", ab76Component)
					}
					if ab76Case.Ab76Operation == 4 {
						if ab76Component.Ab76Default == nil || *ab76Component.Ab76Default != 0 {
							ab76T.Fatal("ab76 decimal default")
						}
					} else if ab76Component.Ab76Default != nil {
						ab76T.Fatal("ab76 invented default")
					}
				}
				if !ab76reflect.DeepEqual(ab76View.Ab76LocalCounter, ab76Case.Ab76Counter) || !ab76reflect.DeepEqual(ab76View.Ab76LocalHours, ab76Case.Ab76Hours) || !ab76reflect.DeepEqual(ab76View.Ab76LocalMinute, ab76Case.Ab76Minute) {
					ab76T.Fatalf("ab76 independent projections: %+v", ab76View)
				}
				if ab76View.Ab76ImplementationGrade != "Proven" || ab76View.Ab76MeaningGrade != "Hypothesis" || len(ab76View.Ab76NativeProperties) != 10 {
					ab76T.Fatal("ab76 grades")
				}
				for _, ab76Grade := range ab76View.Ab76NativeProperties {
					if ab76Grade != "Unknown" {
						ab76T.Fatal("ab76 native qualification")
					}
				}
			}
		})
	}
}

func TestAb76SelectionAxes(ab76T *ab76testing.T) {
	for _, ab76Operation := range []uint8{4, 5, 6} {
		for _, ab76Mode := range []struct{ Ab76Selector, Ab76Direction string }{{"candidate", "request"}, {"standard", "positive"}} {
			if ab76Operation == 5 && ab76Mode.Ab76Selector == "standard" {
				continue
			}
			ab76Control := ab76Selection(ab76Operation, ab76Mode.Ab76Selector, ab76Mode.Ab76Direction)
			ab76Body := []byte{1, 2, 3}
			if ab76Operation == 6 {
				ab76Body = []byte{59, 1, 2, 3}
			}
			if Ab76Interpret(ab76Control, ab76Body).Ab76State != Ab76Complete {
				ab76T.Fatal("ab76 control")
			}
			for _, ab76Case := range []struct {
				Ab76Label  string
				Ab76Mutate func(*Ab76Selection)
			}{
				{"ab76-family", func(ab76S *Ab76Selection) { ab76S.Ab76Family = "ab76-unknown" }},
				{"ab76-family-absent", func(ab76S *Ab76Selection) { ab76S.Ab76Family = "" }},
				{"ab76-variant", func(ab76S *Ab76Selection) { ab76S.Ab76Variant = "03/07" }},
				{"ab76-variant08", func(ab76S *Ab76Selection) { ab76S.Ab76Variant = "03/08" }},
				{"ab76-variant10", func(ab76S *Ab76Selection) { ab76S.Ab76Variant = "03/10" }},
				{"ab76-variant-absent", func(ab76S *Ab76Selection) { ab76S.Ab76Variant = "" }},
				{"ab76-variant05-06", func(ab76S *Ab76Selection) {
					ab76S.Ab76Variant = "03/05"
					if ab76Operation == 5 {
						ab76S.Ab76Variant = "03/06"
					}
				}},
				{"ab76-function", func(ab76S *Ab76Selection) { ab76S.Ab76Function = ab76Byte(9) }},
				{"ab76-function-absent", func(ab76S *Ab76Selection) { ab76S.Ab76Function = nil }},
				{"ab76-function-zero", func(ab76S *Ab76Selection) { ab76S.Ab76Function = ab76Byte(0) }},
				{"ab76-operation", func(ab76S *Ab76Selection) { ab76S.Ab76Operation = ab76Byte(7) }},
				{"ab76-operation-absent", func(ab76S *Ab76Selection) { ab76S.Ab76Operation = nil }},
				{"ab76-operation-zero", func(ab76S *Ab76Selection) { ab76S.Ab76Operation = ab76Byte(0) }},
				{"ab76-direction", func(ab76S *Ab76Selection) { ab76S.Ab76Direction = "ab76-unknown" }},
				{"ab76-direction-absent", func(ab76S *Ab76Selection) { ab76S.Ab76Direction = "" }},
				{"ab76-direction-opposite", func(ab76S *Ab76Selection) {
					ab76S.Ab76Direction = "positive"
					if ab76Mode.Ab76Direction == "positive" {
						ab76S.Ab76Direction = "request"
					}
				}},
				{"ab76-selector", func(ab76S *Ab76Selection) { ab76S.Ab76Selector = "ab76-unknown" }},
				{"ab76-selector-absent", func(ab76S *Ab76Selection) { ab76S.Ab76Selector = "" }},
				{"ab76-selector-opposite", func(ab76S *Ab76Selection) {
					ab76S.Ab76Selector = "standard"
					if ab76Mode.Ab76Selector == "standard" {
						ab76S.Ab76Selector = "candidate"
					}
				}},
				{"ab76-encoding-unspecified", func(ab76S *Ab76Selection) { ab76S.Ab76Encoding = Ab76EncodingUnspecified }},
				{"ab76-encoding-unknown", func(ab76S *Ab76Selection) { ab76S.Ab76Encoding = "ab76-unknown" }},
				{"ab76-encoding-opposite", func(ab76S *Ab76Selection) {
					ab76S.Ab76Encoding = Ab76DecimalChunks
					if ab76Operation == 4 {
						ab76S.Ab76Encoding = Ab76UnsignedU8
					}
				}},
			} {
				ab76T.Run(ab76fmt.Sprintf("ab76-%d-%s-%s", ab76Operation, ab76Mode.Ab76Selector, ab76Case.Ab76Label), func(ab76T *ab76testing.T) {
					ab76Changed := ab76Control
					ab76Case.Ab76Mutate(&ab76Changed)
					ab76View := Ab76Interpret(ab76Changed, ab76Body)
					if ab76View.Ab76State != Ab76RawOnly || len(ab76View.Ab76Components) != 0 || !ab76reflect.DeepEqual(ab76View.Ab76Raw, ab76Body) || !ab76reflect.DeepEqual(ab76View.Ab76Selection, ab76Changed) {
						ab76T.Fatalf("ab76 axis admitted: %+v", ab76View)
					}
					ab76NoProjection(ab76T, ab76View)
				})
			}
		}
	}
	ab76Unsupported := ab76Selection(5, "standard", "positive")
	if Ab76Interpret(ab76Unsupported, []byte{1, 2, 3}).Ab76State != Ab76RawOnly {
		ab76T.Fatal("ab76 05 standard admitted")
	}
}

func TestAb76LengthsAndFields(ab76T *ab76testing.T) {
	for _, ab76Operation := range []uint8{4, 5, 6} {
		ab76S := ab76Selection(ab76Operation, "candidate", "request")
		ab76Maxima := []byte{99, 99, 99}
		if ab76Operation == 5 {
			ab76Maxima = []byte{59, 59, 59}
		}
		if ab76Operation == 6 {
			ab76Maxima = []byte{59, 99, 99, 99}
		}
		for ab76Length := 0; ab76Length <= len(ab76Maxima)+1; ab76Length++ {
			ab76T.Run(ab76fmt.Sprintf("ab76-%d-length-%d", ab76Operation, ab76Length), func(ab76T *ab76testing.T) {
				ab76Body := append(append([]byte(nil), ab76Maxima...), 42)[:ab76Length]
				ab76View := Ab76Interpret(ab76S, ab76Body)
				ab76WantState := Ab76Complete
				if ab76Length < len(ab76Maxima) {
					ab76WantState = Ab76Incomplete
				}
				if ab76Length > len(ab76Maxima) {
					ab76WantState = Ab76Extra
				}
				if ab76View.Ab76State != ab76WantState || !ab76reflect.DeepEqual(ab76View.Ab76Raw, ab76Body) {
					ab76T.Fatalf("ab76 length state: %+v", ab76View)
				}
				if ab76WantState != Ab76Complete {
					ab76NoProjection(ab76T, ab76View)
				}
				if ab76WantState == Ab76Extra && !ab76reflect.DeepEqual(ab76View.Ab76Remainder, []byte{42}) {
					ab76T.Fatal("ab76 remainder")
				}
				for ab76Index, ab76Component := range ab76View.Ab76Components {
					if ab76Index < ab76Length {
						if ab76Component.Ab76State != Ab76Valid || ab76Component.Ab76Value == nil || *ab76Component.Ab76Value != ab76Maxima[ab76Index] {
							ab76T.Fatal("ab76 dropped partial component")
						}
					} else if ab76Component.Ab76State != Ab76Absent || ab76Component.Ab76Raw != nil || ab76Component.Ab76Value != nil {
						ab76T.Fatal("ab76 absent became zero")
					}
				}
			})
		}
		for ab76Index, ab76Maximum := range ab76Maxima {
			for _, ab76Value := range []byte{0, ab76Maximum, ab76Maximum + 1, 255} {
				ab76T.Run(ab76fmt.Sprintf("ab76-%d-field-%d-value-%d", ab76Operation, ab76Index, ab76Value), func(ab76T *ab76testing.T) {
					ab76Body := append([]byte(nil), ab76Maxima...)
					ab76Body[ab76Index] = ab76Value
					ab76View := Ab76Interpret(ab76S, ab76Body)
					ab76Valid := ab76Value <= ab76Maximum
					if (ab76View.Ab76State == Ab76Complete) != ab76Valid {
						ab76T.Fatal("ab76 boundary")
					}
					if !ab76Valid {
						if ab76View.Ab76State != Ab76Invalid {
							ab76T.Fatal("ab76 invalid state")
						}
						ab76NoProjection(ab76T, ab76View)
					}
					for ab76Other, ab76Component := range ab76View.Ab76Components {
						if ab76Component.Ab76Maximum != ab76Maxima[ab76Other] || ab76Component.Ab76Raw == nil || *ab76Component.Ab76Raw != ab76Body[ab76Other] {
							ab76T.Fatal("ab76 raw or domain lost")
						}
						if ab76Other == ab76Index && !ab76Valid {
							if ab76Component.Ab76State != Ab76Invalid || ab76Component.Ab76Value != nil {
								ab76T.Fatal("ab76 invalid accepted")
							}
						} else if ab76Component.Ab76State != Ab76Valid || ab76Component.Ab76Value == nil || *ab76Component.Ab76Value != ab76Body[ab76Other] {
							ab76T.Fatal("ab76 unaffected field lost")
						}
					}
				})
			}
		}
	}
}

func TestAb76OwnedMemory(ab76T *ab76testing.T) {
	ab76S := ab76Selection(4, "candidate", "request")
	ab76Body := []byte{99, 1, 2, 42}
	ab76View := Ab76Interpret(ab76S, ab76Body)
	ab76Baseline := Ab76Interpret(ab76S, ab76Body)
	*ab76S.Ab76Function, *ab76S.Ab76Operation = 0, 0
	ab76Body[0], ab76Body[3] = 0, 0
	if !ab76reflect.DeepEqual(ab76View, ab76Baseline) {
		ab76T.Fatal("ab76 input alias")
	}
	ab76View.Ab76Raw[0], ab76View.Ab76Raw[3] = 8, 8
	ab76View.Ab76Remainder[0] = 9
	*ab76View.Ab76Components[0].Ab76Raw = 7
	*ab76View.Ab76Components[0].Ab76Value = 6
	*ab76View.Ab76Components[0].Ab76Default = 5
	*ab76View.Ab76Selection.Ab76Function = 0
	*ab76View.Ab76Selection.Ab76Operation = 0
	ab76View.Ab76NativeProperties["physical-support"] = "changed"
	if *ab76Baseline.Ab76Components[0].Ab76Value != 99 || ab76View.Ab76Raw[0] != 8 || ab76View.Ab76Raw[3] != 8 || *ab76View.Ab76Components[0].Ab76Raw != 7 || *ab76View.Ab76Components[1].Ab76Default != 0 {
		ab76T.Fatal("ab76 return internal alias")
	}
	if !ab76reflect.DeepEqual(ab76Baseline, Ab76Interpret(ab76Selection(4, "candidate", "request"), []byte{99, 1, 2, 42})) {
		ab76T.Fatal("ab76 shared return storage")
	}
	if ab76Body[0] != 0 || ab76Body[3] != 0 {
		ab76T.Fatal("ab76 output aliases input")
	}
	ab76Complete := Ab76Interpret(ab76Selection(6, "candidate", "request"), []byte{59, 1, 2, 3})
	*ab76Complete.Ab76LocalMinute, *ab76Complete.Ab76LocalHours = 0, 0
	ab76Fresh := Ab76Interpret(ab76Selection(6, "candidate", "request"), []byte{59, 1, 2, 3})
	if *ab76Fresh.Ab76LocalMinute != 59 || *ab76Fresh.Ab76LocalHours != 30201 || *ab76Complete.Ab76Components[0].Ab76Value != 59 {
		ab76T.Fatal("ab76 projection alias")
	}
}

func TestAb76ProjectionAxes(ab76T *ab76testing.T) {
	for _, ab76Case := range []struct {
		Ab76Label     string
		Ab76Operation uint8
		Ab76Body      []byte
		Ab76Counter   int
		Ab76Hours     int
		Ab76Minute    int
	}{
		{"ab76-c0-zero", 4, []byte{0, 1, 2}, 20100, 0, 0},
		{"ab76-c1-zero", 4, []byte{99, 0, 2}, 20099, 0, 0},
		{"ab76-c2-zero", 4, []byte{99, 1, 0}, 199, 0, 0},
		{"ab76-c-all-max", 4, []byte{99, 99, 99}, 999999, 0, 0},
		{"ab76-minute-zero", 6, []byte{0, 1, 2, 3}, 0, 30201, 0},
		{"ab76-h0-zero", 6, []byte{59, 0, 2, 3}, 0, 30200, 59},
		{"ab76-h1-zero", 6, []byte{59, 1, 0, 3}, 0, 30001, 59},
		{"ab76-h2-zero", 6, []byte{59, 1, 2, 0}, 0, 201, 59},
		{"ab76-h-all-max", 6, []byte{59, 99, 99, 99}, 0, 999999, 59},
	} {
		ab76T.Run(ab76Case.Ab76Label, func(ab76T *ab76testing.T) {
			ab76View := Ab76Interpret(ab76Selection(ab76Case.Ab76Operation, "candidate", "request"), ab76Case.Ab76Body)
			if ab76View.Ab76State != Ab76Complete {
				ab76T.Fatal("ab76 zero component is present")
			}
			if ab76Case.Ab76Operation == 4 {
				if ab76View.Ab76LocalCounter == nil || *ab76View.Ab76LocalCounter != ab76Case.Ab76Counter || ab76View.Ab76LocalHours != nil || ab76View.Ab76LocalMinute != nil {
					ab76T.Fatal("ab76 counter weight")
				}
			} else if ab76View.Ab76LocalHours == nil || *ab76View.Ab76LocalHours != ab76Case.Ab76Hours || ab76View.Ab76LocalMinute == nil || *ab76View.Ab76LocalMinute != ab76Case.Ab76Minute || ab76View.Ab76LocalCounter != nil {
				ab76T.Fatal("ab76 hour or minute weight")
			}
		})
	}
}

func TestAb76RawOnlyAndUnselectedFunctions(ab76T *ab76testing.T) {
	for _, ab76Operation := range []uint8{0, 7, 8, 10, 255} {
		ab76S := ab76Selection(ab76Operation, "candidate", "request")
		for ab76Length := 0; ab76Length <= 5; ab76Length++ {
			ab76T.Run(ab76fmt.Sprintf("ab76-raw-%d-length-%d", ab76Operation, ab76Length), func(ab76T *ab76testing.T) {
				ab76Body := make([]byte, ab76Length)
				ab76View := Ab76Interpret(ab76S, ab76Body)
				if ab76View.Ab76State != Ab76RawOnly || len(ab76View.Ab76Raw) != ab76Length || len(ab76View.Ab76Components) != 0 {
					ab76T.Fatal("ab76 unselected body interpreted")
				}
				ab76NoProjection(ab76T, ab76View)
				if ab76Length > 0 {
					ab76Body[0] = 42
					if ab76View.Ab76Raw[0] != 0 {
						ab76T.Fatal("ab76 raw-only input alias")
					}
					ab76View.Ab76Raw[0] = 99
					if ab76Body[0] != 42 {
						ab76T.Fatal("ab76 raw-only output alias")
					}
				}
			})
		}
	}
}
