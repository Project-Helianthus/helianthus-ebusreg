package vaillant

import (
	ab13an13JSON "encoding/json"
	ab13an13OS "os"
	ab13an13Reflect "reflect"
	ab13an13Strconv "strconv"
	ab13an13Testing "testing"
)

type ab13an13Scalars struct {
	Raw      int
	Value    float64
	Conflict bool
	Named    bool
}

type ab13an13Expect struct {
	Admitted   bool
	Complete   bool
	Bits       map[string]bool
	Scalars    map[string]ab13an13Scalars
	Absent     []string
	Unmodelled map[string]int
}

type ab13an13Vector struct {
	Ab13An13Label     string
	Ab13An13Family    string
	Ab13An13Variant   string
	Ab13An13Direction string
	Ab13An13Image     []int
	Ab13An13Expected  ab13an13Expect
}

func ab13an13Bits(ab13an13Result Ab13An13Result) map[string]*bool {
	ab13an13Map := map[string]*bool{}
	if ab13an13Request := ab13an13Result.Ab13An13Request; ab13an13Request != nil {
		ab13an13Map["task"] = ab13an13Request.Ab13An13TaskActive
		ab13an13Map["heat"] = ab13an13Request.Ab13An13HeatDemand
		ab13an13Map["loading"] = ab13an13Request.Ab13An13LoadingActive
		ab13an13Map["pumpfault"] = ab13an13Request.Ab13An13LoadingPumpFault
		for ab13an13Index, ab13an13Bit := range ab13an13Request.Ab13An13SensorConditionBits {
			ab13an13Map["sens"+ab13an13Strconv.Itoa(ab13an13Index)] = ab13an13Bit
		}
		for ab13an13Index, ab13an13Bit := range ab13an13Request.Ab13An13InletConditionBits {
			ab13an13Map["inlet"+ab13an13Strconv.Itoa(ab13an13Index)] = ab13an13Bit
		}
	}
	if ab13an13Positive := ab13an13Result.Ab13An13Positive; ab13an13Positive != nil {
		ab13an13Map["disabled"] = ab13an13Positive.Ab13An13LoadingDisabled
		ab13an13Map["test"] = ab13an13Positive.Ab13An13InternalTestActive
		ab13an13Map["ready"] = ab13an13Positive.Ab13An13Ready
		ab13an13Map["frost"] = ab13an13Positive.Ab13An13FrostProtectionActive
		ab13an13Map["legionella"] = ab13an13Positive.Ab13An13LegionellaDemand
		ab13an13Map["reset"] = ab13an13Positive.Ab13An13Reset
	}
	return ab13an13Map
}

func ab13an13Scalar13s(ab13an13Result Ab13An13Result) map[string]*Ab13An13Scalar {
	ab13an13Map := map[string]*Ab13An13Scalar{}
	if ab13an13Request := ab13an13Result.Ab13An13Request; ab13an13Request != nil {
		ab13an13Map["modulation"] = ab13an13Request.Ab13An13Modulation
		ab13an13Map["flow"] = ab13an13Request.Ab13An13Flow
		for ab13an13Index, ab13an13Value := range ab13an13Request.Ab13An13Temperatures {
			ab13an13Map["temp"+ab13an13Strconv.Itoa(ab13an13Index)] = ab13an13Value
		}
	}
	if ab13an13Positive := ab13an13Result.Ab13An13Positive; ab13an13Positive != nil {
		ab13an13Map["setpoint"] = ab13an13Positive.Ab13An13Setpoint
		ab13an13Map["return"] = ab13an13Positive.Ab13An13ReturnTemperature
		ab13an13Map["adjustment"] = ab13an13Positive.Ab13An13Adjustment
	}
	return ab13an13Map
}

func TestAb13An13FiniteVectors(ab13an13T *ab13an13Testing.T) {
	ab13an13Data, ab13an13Error := ab13an13OS.ReadFile("testdata/ab13an13_b512_dhw_views.json")
	if ab13an13Error != nil {
		ab13an13T.Fatal(ab13an13Error)
	}
	var ab13an13Vectors []ab13an13Vector
	if ab13an13Error = ab13an13JSON.Unmarshal(ab13an13Data, &ab13an13Vectors); ab13an13Error != nil {
		ab13an13T.Fatal(ab13an13Error)
	}
	if len(ab13an13Vectors) < 100 {
		ab13an13T.Fatalf("ab13an13 fixture count: %d", len(ab13an13Vectors))
	}
	for _, ab13an13Vector := range ab13an13Vectors {
		ab13an13T.Run(ab13an13Vector.Ab13An13Label, func(ab13an13T *ab13an13Testing.T) {
			ab13an13Image := make([]byte, len(ab13an13Vector.Ab13An13Image))
			for ab13an13Index, ab13an13Value := range ab13an13Vector.Ab13An13Image {
				ab13an13Image[ab13an13Index] = byte(ab13an13Value)
			}
			ab13an13Original := append([]byte(nil), ab13an13Image...)
			ab13an13Result := Ab13An13Interpret(ab13an13Vector.Ab13An13Family, ab13an13Vector.Ab13An13Variant, ab13an13Vector.Ab13An13Direction, ab13an13Image)
			ab13an13Want := ab13an13Vector.Ab13An13Expected
			if !ab13an13Reflect.DeepEqual(ab13an13Result.Ab13An13Raw, ab13an13Original) && len(ab13an13Original) > 0 {
				ab13an13T.Fatalf("ab13an13 raw %v want %v", ab13an13Result.Ab13An13Raw, ab13an13Original)
			}
			if ab13an13Result.Ab13An13Admitted != ab13an13Want.Admitted || ab13an13Result.Ab13An13Complete != ab13an13Want.Complete {
				ab13an13T.Fatalf("ab13an13 admitted/complete %v/%v want %v/%v", ab13an13Result.Ab13An13Admitted, ab13an13Result.Ab13An13Complete, ab13an13Want.Admitted, ab13an13Want.Complete)
			}
			ab13an13Unmodelled := map[string]int{}
			for ab13an13Index, ab13an13Value := range ab13an13Result.Ab13An13Unmodelled {
				ab13an13Unmodelled[ab13an13Strconv.Itoa(ab13an13Index)] = int(ab13an13Value)
			}
			if !ab13an13Reflect.DeepEqual(ab13an13Unmodelled, ab13an13Nonnil(ab13an13Want.Unmodelled)) {
				ab13an13T.Fatalf("ab13an13 unmodelled %v want %v", ab13an13Unmodelled, ab13an13Want.Unmodelled)
			}
			ab13an13Bits13 := ab13an13Bits(ab13an13Result)
			ab13an13Scalars13 := ab13an13Scalar13s(ab13an13Result)
			if !ab13an13Want.Admitted && (ab13an13Result.Ab13An13Request != nil || ab13an13Result.Ab13An13Positive != nil) {
				ab13an13T.Fatal("ab13an13 typed values attached without admission")
			}
			for ab13an13Name, ab13an13Expected := range ab13an13Want.Bits {
				ab13an13Actual := ab13an13Bits13[ab13an13Name]
				if ab13an13Actual == nil || *ab13an13Actual != ab13an13Expected {
					ab13an13T.Fatalf("ab13an13 bit %s actual %v want %v", ab13an13Name, ab13an13Actual, ab13an13Expected)
				}
			}
			for ab13an13Name, ab13an13Expected := range ab13an13Want.Scalars {
				ab13an13Actual := ab13an13Scalars13[ab13an13Name]
				if ab13an13Actual == nil || int(ab13an13Actual.Ab13An13Raw) != ab13an13Expected.Raw || ab13an13Actual.Ab13An13Value != ab13an13Expected.Value ||
					ab13an13Actual.Ab13An13UnitConflict != ab13an13Expected.Conflict || ab13an13Actual.Ab13An13NamedDomain != ab13an13Expected.Named {
					ab13an13T.Fatalf("ab13an13 scalar %s actual %+v want %+v", ab13an13Name, ab13an13Actual, ab13an13Expected)
				}
			}
			for _, ab13an13Name := range ab13an13Want.Absent {
				if ab13an13Bit, ab13an13Ok := ab13an13Bits13[ab13an13Name]; ab13an13Ok && ab13an13Bit != nil {
					ab13an13T.Fatalf("ab13an13 bit %s present", ab13an13Name)
				}
				if ab13an13Scalar, ab13an13Ok := ab13an13Scalars13[ab13an13Name]; ab13an13Ok && ab13an13Scalar != nil {
					ab13an13T.Fatalf("ab13an13 scalar %s present", ab13an13Name)
				}
			}
			if len(ab13an13Image) > 0 {
				ab13an13Image[0] ^= 255
				if ab13an13Result.Ab13An13Raw[0] != ab13an13Original[0] {
					ab13an13T.Fatal("ab13an13 raw aliases caller")
				}
			}
		})
	}
}

func ab13an13Nonnil(ab13an13Map map[string]int) map[string]int {
	if ab13an13Map == nil {
		return map[string]int{}
	}
	return ab13an13Map
}

func TestAb13An13NamedFieldsAreExactlyThoseClaimed(ab13an13T *ab13an13Testing.T) {
	ab13an13Request := Ab13An13Interpret(Ab13An13Family, Ab13An13Variant, "request", []byte{0xb5, 0x12, 0x06, 0xff, 0, 0, 0, 0, 0, 0xff})
	ab13an13Names := ab13an13Bits(ab13an13Request)
	if len(ab13an13Names) != 10 {
		ab13an13T.Fatalf("ab13an13 request bit count %d", len(ab13an13Names))
	}
	for ab13an13Name, ab13an13Bit := range ab13an13Names {
		if ab13an13Bit == nil || !*ab13an13Bit {
			ab13an13T.Fatalf("ab13an13 %s not set", ab13an13Name)
		}
	}
	if ab13an13Request.Ab13An13Unmodelled[9] != 0xfc {
		ab13an13T.Fatalf("ab13an13 reserved inlet bits lost: %v", ab13an13Request.Ab13An13Unmodelled)
	}
	if !ab13an13Request.Ab13An13Request.Ab13An13Modulation.Ab13An13UnitConflict {
		ab13an13T.Fatal("ab13an13 modulation unit conflict lost")
	}
}
