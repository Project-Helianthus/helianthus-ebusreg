package vaillant

import (
	b512DHWViewJSON "encoding/json"
	b512DHWViewOS "os"
	b512DHWViewReflect "reflect"
	b512DHWViewStrconv "strconv"
	b512DHWViewTesting "testing"
)

type b512DHWViewScalars struct {
	Raw      int
	Value    float64
	Conflict bool
	Named    bool
}

type b512DHWViewExpect struct {
	Admitted   bool
	Complete   bool
	Bits       map[string]bool
	Scalars    map[string]b512DHWViewScalars
	Absent     []string
	Unmodelled map[string]int
}

type b512DHWViewVector struct {
	B512DHWViewLabel     string            `json:"Ab13An13Label"`
	B512DHWViewFamily    string            `json:"Ab13An13Family"`
	B512DHWViewVariant   string            `json:"Ab13An13Variant"`
	B512DHWViewDirection string            `json:"Ab13An13Direction"`
	B512DHWViewImage     []int             `json:"Ab13An13Image"`
	B512DHWViewExpected  b512DHWViewExpect `json:"Ab13An13Expected"`
}

func b512DHWViewBits(b512DHWViewResult B512DHWViewResult) map[string]*bool {
	b512DHWViewMap := map[string]*bool{}
	if b512DHWViewRequest := b512DHWViewResult.B512DHWViewRequest; b512DHWViewRequest != nil {
		b512DHWViewMap["task"] = b512DHWViewRequest.B512DHWViewTaskActive
		b512DHWViewMap["heat"] = b512DHWViewRequest.B512DHWViewHeatDemand
		b512DHWViewMap["loading"] = b512DHWViewRequest.B512DHWViewLoadingActive
		b512DHWViewMap["pumpfault"] = b512DHWViewRequest.B512DHWViewLoadingPumpFault
		for b512DHWViewIndex, b512DHWViewBit := range b512DHWViewRequest.B512DHWViewSensorConditionBits {
			b512DHWViewMap["sens"+b512DHWViewStrconv.Itoa(b512DHWViewIndex)] = b512DHWViewBit
		}
		for b512DHWViewIndex, b512DHWViewBit := range b512DHWViewRequest.B512DHWViewInletConditionBits {
			b512DHWViewMap["inlet"+b512DHWViewStrconv.Itoa(b512DHWViewIndex)] = b512DHWViewBit
		}
	}
	if b512DHWViewPositive := b512DHWViewResult.B512DHWViewPositive; b512DHWViewPositive != nil {
		b512DHWViewMap["disabled"] = b512DHWViewPositive.B512DHWViewLoadingDisabled
		b512DHWViewMap["test"] = b512DHWViewPositive.B512DHWViewInternalTestActive
		b512DHWViewMap["ready"] = b512DHWViewPositive.B512DHWViewReady
		b512DHWViewMap["frost"] = b512DHWViewPositive.B512DHWViewFrostProtectionActive
		b512DHWViewMap["legionella"] = b512DHWViewPositive.B512DHWViewLegionellaDemand
		b512DHWViewMap["reset"] = b512DHWViewPositive.B512DHWViewReset
	}
	return b512DHWViewMap
}

func b512DHWViewScalar13s(b512DHWViewResult B512DHWViewResult) map[string]*B512DHWViewScalar {
	b512DHWViewMap := map[string]*B512DHWViewScalar{}
	if b512DHWViewRequest := b512DHWViewResult.B512DHWViewRequest; b512DHWViewRequest != nil {
		b512DHWViewMap["modulation"] = b512DHWViewRequest.B512DHWViewModulation
		b512DHWViewMap["flow"] = b512DHWViewRequest.B512DHWViewFlow
		for b512DHWViewIndex, b512DHWViewValue := range b512DHWViewRequest.B512DHWViewTemperatures {
			b512DHWViewMap["temp"+b512DHWViewStrconv.Itoa(b512DHWViewIndex)] = b512DHWViewValue
		}
	}
	if b512DHWViewPositive := b512DHWViewResult.B512DHWViewPositive; b512DHWViewPositive != nil {
		b512DHWViewMap["setpoint"] = b512DHWViewPositive.B512DHWViewSetpoint
		b512DHWViewMap["return"] = b512DHWViewPositive.B512DHWViewReturnTemperature
		b512DHWViewMap["adjustment"] = b512DHWViewPositive.B512DHWViewAdjustment
	}
	return b512DHWViewMap
}

func TestAb13An13FiniteVectors(b512DHWViewT *b512DHWViewTesting.T) {
	b512DHWViewData, b512DHWViewError := b512DHWViewOS.ReadFile("testdata/ab13an13_b512_dhw_views.json")
	if b512DHWViewError != nil {
		b512DHWViewT.Fatal(b512DHWViewError)
	}
	var b512DHWViewVectors []b512DHWViewVector
	if b512DHWViewError = b512DHWViewJSON.Unmarshal(b512DHWViewData, &b512DHWViewVectors); b512DHWViewError != nil {
		b512DHWViewT.Fatal(b512DHWViewError)
	}
	if len(b512DHWViewVectors) < 100 {
		b512DHWViewT.Fatalf("ab13an13 fixture count: %d", len(b512DHWViewVectors))
	}
	for _, b512DHWViewVector := range b512DHWViewVectors {
		b512DHWViewT.Run(b512DHWViewVector.B512DHWViewLabel, func(b512DHWViewT *b512DHWViewTesting.T) {
			b512DHWViewImage := make([]byte, len(b512DHWViewVector.B512DHWViewImage))
			for b512DHWViewIndex, b512DHWViewValue := range b512DHWViewVector.B512DHWViewImage {
				b512DHWViewImage[b512DHWViewIndex] = byte(b512DHWViewValue)
			}
			b512DHWViewOriginal := append([]byte(nil), b512DHWViewImage...)
			b512DHWViewResult := B512DHWViewInterpret(b512DHWViewVector.B512DHWViewFamily, b512DHWViewVector.B512DHWViewVariant, b512DHWViewVector.B512DHWViewDirection, b512DHWViewImage)
			b512DHWViewWant := b512DHWViewVector.B512DHWViewExpected
			if !b512DHWViewReflect.DeepEqual(b512DHWViewResult.B512DHWViewRaw, b512DHWViewOriginal) && len(b512DHWViewOriginal) > 0 {
				b512DHWViewT.Fatalf("ab13an13 raw %v want %v", b512DHWViewResult.B512DHWViewRaw, b512DHWViewOriginal)
			}
			if b512DHWViewResult.B512DHWViewAdmitted != b512DHWViewWant.Admitted || b512DHWViewResult.B512DHWViewComplete != b512DHWViewWant.Complete {
				b512DHWViewT.Fatalf("ab13an13 admitted/complete %v/%v want %v/%v", b512DHWViewResult.B512DHWViewAdmitted, b512DHWViewResult.B512DHWViewComplete, b512DHWViewWant.Admitted, b512DHWViewWant.Complete)
			}
			b512DHWViewUnmodelled := map[string]int{}
			for b512DHWViewIndex, b512DHWViewValue := range b512DHWViewResult.B512DHWViewUnmodelled {
				b512DHWViewUnmodelled[b512DHWViewStrconv.Itoa(b512DHWViewIndex)] = int(b512DHWViewValue)
			}
			if !b512DHWViewReflect.DeepEqual(b512DHWViewUnmodelled, b512DHWViewNonnil(b512DHWViewWant.Unmodelled)) {
				b512DHWViewT.Fatalf("ab13an13 unmodelled %v want %v", b512DHWViewUnmodelled, b512DHWViewWant.Unmodelled)
			}
			b512DHWViewBits13 := b512DHWViewBits(b512DHWViewResult)
			b512DHWViewScalars13 := b512DHWViewScalar13s(b512DHWViewResult)
			if !b512DHWViewWant.Admitted && (b512DHWViewResult.B512DHWViewRequest != nil || b512DHWViewResult.B512DHWViewPositive != nil) {
				b512DHWViewT.Fatal("ab13an13 typed values attached without admission")
			}
			for b512DHWViewName, b512DHWViewExpected := range b512DHWViewWant.Bits {
				b512DHWViewActual := b512DHWViewBits13[b512DHWViewName]
				if b512DHWViewActual == nil || *b512DHWViewActual != b512DHWViewExpected {
					b512DHWViewT.Fatalf("ab13an13 bit %s actual %v want %v", b512DHWViewName, b512DHWViewActual, b512DHWViewExpected)
				}
			}
			for b512DHWViewName, b512DHWViewExpected := range b512DHWViewWant.Scalars {
				b512DHWViewActual := b512DHWViewScalars13[b512DHWViewName]
				if b512DHWViewActual == nil || int(b512DHWViewActual.B512DHWViewRaw) != b512DHWViewExpected.Raw || b512DHWViewActual.B512DHWViewValue != b512DHWViewExpected.Value ||
					b512DHWViewActual.B512DHWViewUnitConflict != b512DHWViewExpected.Conflict || b512DHWViewActual.B512DHWViewNamedDomain != b512DHWViewExpected.Named {
					b512DHWViewT.Fatalf("ab13an13 scalar %s actual %+v want %+v", b512DHWViewName, b512DHWViewActual, b512DHWViewExpected)
				}
			}
			for _, b512DHWViewName := range b512DHWViewWant.Absent {
				if b512DHWViewBit, b512DHWViewOk := b512DHWViewBits13[b512DHWViewName]; b512DHWViewOk && b512DHWViewBit != nil {
					b512DHWViewT.Fatalf("ab13an13 bit %s present", b512DHWViewName)
				}
				if b512DHWViewScalar, b512DHWViewOk := b512DHWViewScalars13[b512DHWViewName]; b512DHWViewOk && b512DHWViewScalar != nil {
					b512DHWViewT.Fatalf("ab13an13 scalar %s present", b512DHWViewName)
				}
			}
			if len(b512DHWViewImage) > 0 {
				b512DHWViewImage[0] ^= 255
				if b512DHWViewResult.B512DHWViewRaw[0] != b512DHWViewOriginal[0] {
					b512DHWViewT.Fatal("ab13an13 raw aliases caller")
				}
			}
		})
	}
}

func b512DHWViewNonnil(b512DHWViewMap map[string]int) map[string]int {
	if b512DHWViewMap == nil {
		return map[string]int{}
	}
	return b512DHWViewMap
}

func TestAb13An13NamedFieldsAreExactlyThoseClaimed(b512DHWViewT *b512DHWViewTesting.T) {
	b512DHWViewRequest := B512DHWViewInterpret(B512DHWViewFamily, B512DHWViewVariant, "request", []byte{0xb5, 0x12, 0x06, 0xff, 0, 0, 0, 0, 0, 0xff})
	b512DHWViewNames := b512DHWViewBits(b512DHWViewRequest)
	if len(b512DHWViewNames) != 10 {
		b512DHWViewT.Fatalf("ab13an13 request bit count %d", len(b512DHWViewNames))
	}
	for b512DHWViewName, b512DHWViewBit := range b512DHWViewNames {
		if b512DHWViewBit == nil || !*b512DHWViewBit {
			b512DHWViewT.Fatalf("ab13an13 %s not set", b512DHWViewName)
		}
	}
	if b512DHWViewRequest.B512DHWViewUnmodelled[9] != 0xfc {
		b512DHWViewT.Fatalf("ab13an13 reserved inlet bits lost: %v", b512DHWViewRequest.B512DHWViewUnmodelled)
	}
	if !b512DHWViewRequest.B512DHWViewRequest.B512DHWViewModulation.B512DHWViewUnitConflict {
		b512DHWViewT.Fatal("ab13an13 modulation unit conflict lost")
	}
}
