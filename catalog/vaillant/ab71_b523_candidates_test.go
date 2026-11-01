package vaillant

import (
	ab71Fmt "fmt"
	ab71Reflect "reflect"
	ab71Testing "testing"
)

func ab71Ptr[ab71T any](ab71Value ab71T) *ab71T { return &ab71Value }

func ab71Key(ab71Operation byte, ab71Field string) Ab71Key {
	ab71Direction := "request"
	if ab71Field == "feedback" {
		ab71Direction = "response"
	}
	return Ab71Key{Ab71Family: "vaillant", Ab71Primary: 0xb5, Ab71Secondary: ab71Ptr(byte(0x23)), Ab71Operation: ab71Ptr(ab71Operation), Ab71Direction: ab71Direction, Ab71Field: ab71Field}
}

func ab71Lookup(ab71T *ab71Testing.T, ab71KeyValue Ab71Key) Ab71Entry {
	ab71T.Helper()
	ab71Entry, ab71Found, ab71Error := Ab71Lookup(ab71KeyValue)
	if ab71Error != nil || !ab71Found {
		ab71T.Fatalf("lookup: found=%v error=%v key=%+v", ab71Found, ab71Error, ab71KeyValue)
	}
	return ab71Entry
}

func ab71Rule(ab71Min, ab71Max byte, ab71Result string, ab71Divisor *float64, ab71Duration *string) Ab71Rule {
	return Ab71Rule{Ab71RawMin: ab71Min, Ab71RawMax: ab71Max, Ab71Result: ab71Result, Ab71DivideBy: ab71Divisor, Ab71DurationUnit: ab71Duration, Ab71Grade: "Hypothesis"}
}

func ab71Expected(ab71Operation byte, ab71Field string, ab71Position int, ab71Request *int, ab71Alternatives []Ab71Alternative, ab71Default *string) Ab71Entry {
	return Ab71Entry{Ab71Key: ab71Key(ab71Operation, ab71Field), Ab71LogicalPosition: ab71Ptr(ab71Position), Ab71WidthBits: ab71Ptr(8), Ab71RequestExtent: ab71Request, Ab71PositiveExtent: ab71Ptr(4), Ab71DefaultRaw: ab71Ptr(byte(0)), Ab71DefaultResult: ab71Default, Ab71Grade: "Hypothesis", Ab71RawDomain: []byte{0, 255}, Ab71Alternatives: ab71Alternatives}
}

func ab71Evaluate(ab71Alternative Ab71Alternative, ab71Raw byte) (string, *float64, *int) {
	for _, ab71RuleValue := range ab71Alternative.Ab71Rules {
		if ab71Raw < ab71RuleValue.Ab71RawMin || ab71Raw > ab71RuleValue.Ab71RawMax {
			continue
		}
		var ab71Numeric *float64
		var ab71Duration *int
		if ab71RuleValue.Ab71DivideBy != nil {
			ab71Numeric = ab71Ptr(float64(ab71Raw) / *ab71RuleValue.Ab71DivideBy)
		}
		if ab71RuleValue.Ab71DurationUnit != nil {
			ab71Duration = ab71Ptr(int(ab71Raw))
		}
		return ab71RuleValue.Ab71Result, ab71Numeric, ab71Duration
	}
	return "outside_domain", nil, nil
}

func TestAb71FiniteCatalog(ab71T *ab71Testing.T) {
	ab71Catalog, ab71Error := Ab71B523Candidates()
	if ab71Error != nil {
		ab71T.Fatal(ab71Error)
	}
	if ab71Catalog.Ab71RepresentationGrade != "Proven" || ab71Catalog.Ab71NativeSupport != "Unknown" || ab71Catalog.Ab71RuntimeEffects != "Unknown" || ab71Catalog.Ab71OtherModuleLayouts != "Unknown" {
		ab71T.Fatal("claim boundaries")
	}
	ab71Limitations := []string{
		"Finite candidate metadata; broader family catalog remains incomplete.",
		"Logical positions include the three-position prefix; abstract extents are not physical reply lengths.",
		"Acceptance is not physical ACK, persistence or readback.",
		"No routing, timed relay controller, physical qualification or runtime admission.",
		"Native qualification, freshness and execution authority are not established by candidate metadata.",
	}
	if !ab71Reflect.DeepEqual(ab71Catalog.Ab71Limitations, ab71Limitations) {
		ab71T.Fatal("limitations")
	}
	if len(ab71Catalog.Ab71Entries) != 26 {
		ab71T.Fatalf("entries: %d", len(ab71Catalog.Ab71Entries))
	}
	ab71Seen := map[string]bool{}
	for _, ab71Entry := range ab71Catalog.Ab71Entries {
		ab71ID := ab71Fmt.Sprintf("%d/%s", *ab71Entry.Ab71Key.Ab71Operation, ab71Entry.Ab71Key.Ab71Field)
		if ab71Seen[ab71ID] {
			ab71T.Fatal("duplicate", ab71ID)
		}
		ab71Seen[ab71ID] = true
		if !ab71Reflect.DeepEqual(ab71Entry, ab71Lookup(ab71T, ab71Entry.Ab71Key)) {
			ab71T.Fatal("accessor/lookup disagreement")
		}
	}
}

func TestAb71RelayPositions(ab71T *ab71Testing.T) {
	for _, ab71Operation := range []byte{1, 5} {
		ab71Count, ab71Extent := 6, 10
		if ab71Operation == 5 {
			ab71Count, ab71Extent = 12, 17
		}
		for ab71Index := 0; ab71Index < ab71Count; ab71Index++ {
			ab71T.Run(ab71Fmt.Sprintf("ab71_op%d_relay%d", ab71Operation, ab71Index+1), func(ab71T *ab71Testing.T) {
				ab71Field := ab71Fmt.Sprintf("relay_%d", ab71Index+1)
				ab71Entry := ab71Lookup(ab71T, ab71Key(ab71Operation, ab71Field))
				ab71Rules := []Ab71Rule{ab71Rule(0, 0, "off", nil, nil), ab71Rule(1, 253, "on_then_off", nil, ab71Ptr("seconds")), ab71Rule(254, 254, "on", nil, nil), ab71Rule(255, 255, "unchanged", nil, nil)}
				ab71Want := ab71Expected(ab71Operation, ab71Field, 3+ab71Index, ab71Ptr(ab71Extent), []Ab71Alternative{{Ab71Name: "relay_behavior", Ab71Grade: "Hypothesis", Ab71Rules: ab71Rules}}, nil)
				if !ab71Reflect.DeepEqual(ab71Entry, ab71Want) {
					ab71T.Fatalf("relay metadata: %+v", ab71Entry)
				}
				for ab71IndexValue, ab71Raw := range []byte{0, 1, 253, 254, 255} {
					ab71Result, ab71Numeric, ab71Duration := ab71Evaluate(ab71Entry.Ab71Alternatives[0], ab71Raw)
					ab71Results := []string{"off", "on_then_off", "on_then_off", "on", "unchanged"}
					ab71Durations := []*int{nil, ab71Ptr(1), ab71Ptr(253), nil, nil}
					if ab71Result != ab71Results[ab71IndexValue] || ab71Numeric != nil || !ab71Reflect.DeepEqual(ab71Duration, ab71Durations[ab71IndexValue]) {
						ab71T.Fatalf("raw %d: %s %v", ab71Raw, ab71Result, ab71Duration)
					}
				}
			})
		}
	}
}

func TestAb71PWMConflict(ab71T *ab71Testing.T) {
	for _, ab71Operation := range []byte{1, 5} {
		ab71Positions, ab71Extent := []int{9}, 10
		if ab71Operation == 5 {
			ab71Positions, ab71Extent = []int{15, 16}, 17
		}
		for ab71Index, ab71Position := range ab71Positions {
			ab71T.Run(ab71Fmt.Sprintf("ab71_op%d_pwm%d", ab71Operation, ab71Index+1), func(ab71T *ab71Testing.T) {
				ab71Field := ab71Fmt.Sprintf("pwm_%d", ab71Index+1)
				ab71Entry := ab71Lookup(ab71T, ab71Key(ab71Operation, ab71Field))
				ab71Want := ab71Expected(ab71Operation, ab71Field, ab71Position, ab71Ptr(ab71Extent), []Ab71Alternative{
					{Ab71Name: "numeric_codec", Ab71Grade: "Hypothesis", Ab71Rules: []Ab71Rule{ab71Rule(0, 250, "percent", ab71Ptr(2.0), nil)}},
					{Ab71Name: "behavior", Ab71Grade: "Hypothesis", Ab71Rules: []Ab71Rule{ab71Rule(0, 200, "percent", ab71Ptr(2.0), nil), ab71Rule(201, 255, "unchanged", nil, nil)}},
				}, nil)
				if !ab71Reflect.DeepEqual(ab71Entry, ab71Want) {
					ab71T.Fatal("PWM metadata")
				}
				for ab71IndexValue, ab71Raw := range []byte{200, 201, 250, 255} {
					ab71Result, ab71Numeric, ab71Duration := ab71Evaluate(ab71Entry.Ab71Alternatives[0], ab71Raw)
					ab71Numbers := []*float64{ab71Ptr(100.0), ab71Ptr(100.5), ab71Ptr(125.0), nil}
					ab71ExpectedResult := "percent"
					if ab71Raw == 255 {
						ab71ExpectedResult = "outside_domain"
					}
					if ab71Result != ab71ExpectedResult || !ab71Reflect.DeepEqual(ab71Numeric, ab71Numbers[ab71IndexValue]) || ab71Duration != nil {
						ab71T.Fatalf("numeric raw %d", ab71Raw)
					}
					ab71Result, ab71Numeric, ab71Duration = ab71Evaluate(ab71Entry.Ab71Alternatives[1], ab71Raw)
					ab71ExpectedResult = "unchanged"
					var ab71Number *float64
					if ab71Raw == 200 {
						ab71ExpectedResult = "percent"
						ab71Number = ab71Ptr(100.0)
					}
					if ab71Result != ab71ExpectedResult || !ab71Reflect.DeepEqual(ab71Numeric, ab71Number) || ab71Duration != nil {
						ab71T.Fatalf("behavior raw %d", ab71Raw)
					}
				}
			})
		}
	}
}

func TestAb71Feedback(ab71T *ab71Testing.T) {
	for _, ab71Operation := range []byte{0, 1, 2, 4, 5} {
		ab71T.Run(ab71Fmt.Sprintf("ab71_op%d", ab71Operation), func(ab71T *ab71Testing.T) {
			ab71Entry := ab71Lookup(ab71T, ab71Key(ab71Operation, "feedback"))
			ab71Rejected := "rejected"
			if ab71Operation == 4 || ab71Operation == 5 {
				ab71Rejected = "rejected-or-unsupported"
			}
			var ab71Extent *int
			if ab71Operation == 1 {
				ab71Extent = ab71Ptr(10)
			}
			if ab71Operation == 5 {
				ab71Extent = ab71Ptr(17)
			}
			ab71Want := ab71Expected(ab71Operation, "feedback", 3, ab71Extent, []Ab71Alternative{{Ab71Name: "feedback", Ab71Grade: "Hypothesis", Ab71Rules: []Ab71Rule{ab71Rule(0, 0, ab71Rejected, nil, nil), ab71Rule(1, 1, "accepted", nil, nil), ab71Rule(2, 255, ab71Rejected, nil, nil)}}}, ab71Ptr(ab71Rejected))
			if ab71Extent == nil {
				ab71Want.Ab71PositiveExtent = nil
			}
			if !ab71Reflect.DeepEqual(ab71Entry, ab71Want) {
				ab71T.Fatal("feedback metadata")
			}
			for _, ab71Raw := range []byte{0, 1, 2, 255} {
				ab71WantResult := ab71Rejected
				if ab71Raw == 1 {
					ab71WantResult = "accepted"
				}
				ab71Result, ab71Numeric, ab71Duration := ab71Evaluate(ab71Entry.Ab71Alternatives[0], ab71Raw)
				if ab71Result != ab71WantResult || ab71Numeric != nil || ab71Duration != nil {
					ab71T.Fatalf("feedback %d", ab71Raw)
				}
			}
		})
	}
}

func TestAb71ExactKeyAxes(ab71T *ab71Testing.T) {
	ab71Mutations := map[string]func(*Ab71Key){
		"ab71_family":           func(ab71KeyValue *Ab71Key) { ab71KeyValue.Ab71Family = "other_family" },
		"ab71_variant":          func(ab71KeyValue *Ab71Key) { ab71KeyValue.Ab71Variant = ab71Ptr("other_variant") },
		"ab71_variant_empty":    func(ab71KeyValue *Ab71Key) { ab71KeyValue.Ab71Variant = ab71Ptr("") },
		"ab71_primary":          func(ab71KeyValue *Ab71Key) { ab71KeyValue.Ab71Primary = 0 },
		"ab71_secondary":        func(ab71KeyValue *Ab71Key) { ab71KeyValue.Ab71Secondary = ab71Ptr(byte(0x24)) },
		"ab71_secondary_absent": func(ab71KeyValue *Ab71Key) { ab71KeyValue.Ab71Secondary = nil },
		"ab71_secondary_zero":   func(ab71KeyValue *Ab71Key) { ab71KeyValue.Ab71Secondary = ab71Ptr(byte(0)) },
		"ab71_operation":        func(ab71KeyValue *Ab71Key) { ab71KeyValue.Ab71Operation = ab71Ptr(byte(127)) },
		"ab71_operation_absent": func(ab71KeyValue *Ab71Key) { ab71KeyValue.Ab71Operation = nil },
		"ab71_operation_zero":   func(ab71KeyValue *Ab71Key) { ab71KeyValue.Ab71Operation = ab71Ptr(byte(0)) },
		"ab71_direction":        func(ab71KeyValue *Ab71Key) { ab71KeyValue.Ab71Direction = "response" },
		"ab71_field":            func(ab71KeyValue *Ab71Key) { ab71KeyValue.Ab71Field = "unknown_field" },
	}
	for ab71Name, ab71Mutation := range ab71Mutations {
		ab71T.Run(ab71Name, func(ab71T *ab71Testing.T) {
			ab71Control := ab71Key(1, "relay_1")
			ab71Lookup(ab71T, ab71Control)
			ab71Changed := ab71Key(1, "relay_1")
			ab71Mutation(&ab71Changed)
			if _, ab71Found, ab71Error := Ab71Lookup(ab71Changed); ab71Found || ab71Error != nil {
				ab71T.Fatalf("single-axis mutation found=%v error=%v", ab71Found, ab71Error)
			}
			ab71Lookup(ab71T, ab71Control)
		})
	}
	ab71Zero := ab71Key(0, "feedback")
	ab71Lookup(ab71T, ab71Zero)
	ab71Zero.Ab71Operation = nil
	if _, ab71Found, ab71Error := Ab71Lookup(ab71Zero); ab71Found || ab71Error != nil {
		ab71T.Fatal("absent operation matched zero")
	}
	ab71Allocated := ab71Key(1, "relay_1")
	if !ab71Reflect.DeepEqual(ab71Lookup(ab71T, ab71Allocated), ab71Lookup(ab71T, ab71Key(1, "relay_1"))) {
		ab71T.Fatal("pointer identity used as key")
	}
}

func TestAb71DeepCopy(ab71T *ab71Testing.T) {
	ab71Baseline, ab71Error := Ab71B523Candidates()
	if ab71Error != nil {
		ab71T.Fatal(ab71Error)
	}
	ab71Changed, ab71Error := Ab71B523Candidates()
	if ab71Error != nil {
		ab71T.Fatal(ab71Error)
	}
	ab71Changed.Ab71Limitations[0] = "changed"
	for ab71Index := range ab71Changed.Ab71Entries {
		ab71Entry := &ab71Changed.Ab71Entries[ab71Index]
		*ab71Entry.Ab71Key.Ab71Secondary = 0
		*ab71Entry.Ab71Key.Ab71Operation = 127
		*ab71Entry.Ab71LogicalPosition = 99
		*ab71Entry.Ab71WidthBits = 99
		if ab71Entry.Ab71RequestExtent != nil {
			*ab71Entry.Ab71RequestExtent = 99
		}
		if ab71Entry.Ab71PositiveExtent != nil {
			*ab71Entry.Ab71PositiveExtent = 99
		}
		*ab71Entry.Ab71DefaultRaw = 99
		if ab71Entry.Ab71DefaultResult != nil {
			*ab71Entry.Ab71DefaultResult = "changed"
		}
		ab71Entry.Ab71RawDomain[0] = 99
		for ab71AlternativeIndex := range ab71Entry.Ab71Alternatives {
			ab71Alternative := &ab71Entry.Ab71Alternatives[ab71AlternativeIndex]
			ab71Alternative.Ab71Name = "changed"
			for ab71RuleIndex := range ab71Alternative.Ab71Rules {
				ab71RuleValue := &ab71Alternative.Ab71Rules[ab71RuleIndex]
				ab71RuleValue.Ab71Result = "changed"
				if ab71RuleValue.Ab71DivideBy != nil {
					*ab71RuleValue.Ab71DivideBy = 99
				}
				if ab71RuleValue.Ab71DurationUnit != nil {
					*ab71RuleValue.Ab71DurationUnit = "changed"
				}
			}
		}
	}
	ab71Fresh, ab71Error := Ab71B523Candidates()
	if ab71Error != nil || !ab71Reflect.DeepEqual(ab71Fresh, ab71Baseline) {
		ab71T.Fatal("catalog storage shared")
	}
	for _, ab71ExpectedEntry := range ab71Baseline.Ab71Entries {
		ab71Returned := ab71Lookup(ab71T, ab71ExpectedEntry.Ab71Key)
		*ab71Returned.Ab71LogicalPosition = 99
		ab71Returned.Ab71Alternatives[0].Ab71Rules[0].Ab71Result = "changed"
		if !ab71Reflect.DeepEqual(ab71Lookup(ab71T, ab71ExpectedEntry.Ab71Key), ab71ExpectedEntry) {
			ab71T.Fatal("lookup storage shared")
		}
	}
}
