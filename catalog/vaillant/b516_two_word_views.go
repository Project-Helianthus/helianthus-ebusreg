package vaillant

import (
	b516TwoWordViewBinary "encoding/binary"
	b516TwoWordViewMath "math"
)

const (
	B516TwoWordViewFamily  = "vaillant-b516-family"
	B516TwoWordViewVariant = "b516-op14-candidate"
)

type B516TwoWordViewResult struct {
	B516TwoWordViewDeclaredCount       int      `json:"Ab17An17DeclaredCount"`
	B516TwoWordViewAvailableExtent     int      `json:"Ab17An17AvailableExtent"`
	B516TwoWordViewRaw                 []byte   `json:"Ab17An17Raw"`
	B516TwoWordViewTrailing            []byte   `json:"Ab17An17Trailing"`
	B516TwoWordViewAdmitted            bool     `json:"Ab17An17Admitted"`
	B516TwoWordViewStatusAvailable     bool     `json:"Ab17An17StatusAvailable"`
	B516TwoWordViewStatus              byte     `json:"Ab17An17Status"`
	B516TwoWordViewWord0Available      bool     `json:"Ab17An17Word0Available"`
	B516TwoWordViewWord0Raw            [4]byte  `json:"Ab17An17Word0Raw"`
	B516TwoWordViewWord0Bits           uint32   `json:"Ab17An17Word0Bits"`
	B516TwoWordViewWord0Finite         bool     `json:"Ab17An17Word0Finite"`
	B516TwoWordViewWord0Value          *float32 `json:"Ab17An17Word0Value"`
	B516TwoWordViewWord1Available      bool     `json:"Ab17An17Word1Available"`
	B516TwoWordViewWord1Raw            [4]byte  `json:"Ab17An17Word1Raw"`
	B516TwoWordViewCompleteTwoWord     bool     `json:"Ab17An17CompleteTwoWord"`
	B516TwoWordViewContractGrade       string   `json:"Ab17An17ContractGrade"`
	B516TwoWordViewInterpretationGrade string   `json:"Ab17An17InterpretationGrade"`
	B516TwoWordViewStatusMeaning       string   `json:"Ab17An17StatusMeaning"`
	B516TwoWordViewWord0NativeUnit     string   `json:"Ab17An17Word0NativeUnit"`
	B516TwoWordViewWord0NativeSign     string   `json:"Ab17An17Word0NativeSign"`
	B516TwoWordViewWord1Type           string   `json:"Ab17An17Word1Type"`
	B516TwoWordViewWord1Scale          string   `json:"Ab17An17Word1Scale"`
	B516TwoWordViewWord1Sign           string   `json:"Ab17An17Word1Sign"`
	B516TwoWordViewQualification       string   `json:"Ab17An17Qualification"`
}

func B516TwoWordViewInterpret(b516TwoWordViewFamily, b516TwoWordViewVariant string, b516TwoWordViewOperation byte, b516TwoWordViewDirection string, b516TwoWordViewDeclaredCount int, b516TwoWordViewBody []byte) B516TwoWordViewResult {
	b516TwoWordViewResult := B516TwoWordViewResult{
		B516TwoWordViewDeclaredCount:       b516TwoWordViewDeclaredCount,
		B516TwoWordViewAvailableExtent:     len(b516TwoWordViewBody),
		B516TwoWordViewRaw:                 append([]byte(nil), b516TwoWordViewBody...),
		B516TwoWordViewContractGrade:       "Proven",
		B516TwoWordViewInterpretationGrade: "Hypothesis",
		B516TwoWordViewStatusMeaning:       "Unknown",
		B516TwoWordViewWord0NativeUnit:     "Unknown",
		B516TwoWordViewWord0NativeSign:     "Unknown",
		B516TwoWordViewWord1Type:           "Unknown",
		B516TwoWordViewWord1Scale:          "Unknown",
		B516TwoWordViewWord1Sign:           "Unknown",
		B516TwoWordViewQualification:       "Unknown",
	}
	if len(b516TwoWordViewBody) > 9 {
		b516TwoWordViewResult.B516TwoWordViewTrailing = append([]byte(nil), b516TwoWordViewBody[9:]...)
	}
	if b516TwoWordViewFamily != B516TwoWordViewFamily || b516TwoWordViewVariant != B516TwoWordViewVariant || b516TwoWordViewOperation != 0x14 || b516TwoWordViewDirection != "positive" || b516TwoWordViewDeclaredCount < 0 {
		return b516TwoWordViewResult
	}
	b516TwoWordViewResult.B516TwoWordViewAdmitted = true
	if b516TwoWordViewDeclaredCount >= 1 && len(b516TwoWordViewBody) >= 1 {
		b516TwoWordViewResult.B516TwoWordViewStatusAvailable = true
		b516TwoWordViewResult.B516TwoWordViewStatus = b516TwoWordViewBody[0]
	}
	if b516TwoWordViewDeclaredCount >= 5 && len(b516TwoWordViewBody) >= 5 {
		b516TwoWordViewResult.B516TwoWordViewWord0Available = true
		copy(b516TwoWordViewResult.B516TwoWordViewWord0Raw[:], b516TwoWordViewBody[1:5])
		b516TwoWordViewResult.B516TwoWordViewWord0Bits = b516TwoWordViewBinary.LittleEndian.Uint32(b516TwoWordViewBody[1:5])
		b516TwoWordViewValue := b516TwoWordViewMath.Float32frombits(b516TwoWordViewResult.B516TwoWordViewWord0Bits)
		if !b516TwoWordViewMath.IsNaN(float64(b516TwoWordViewValue)) && !b516TwoWordViewMath.IsInf(float64(b516TwoWordViewValue), 0) {
			b516TwoWordViewResult.B516TwoWordViewWord0Finite = true
			b516TwoWordViewResult.B516TwoWordViewWord0Value = &b516TwoWordViewValue
		}
	}
	if b516TwoWordViewDeclaredCount >= 9 && len(b516TwoWordViewBody) >= 9 {
		b516TwoWordViewResult.B516TwoWordViewWord1Available = true
		copy(b516TwoWordViewResult.B516TwoWordViewWord1Raw[:], b516TwoWordViewBody[5:9])
		b516TwoWordViewResult.B516TwoWordViewCompleteTwoWord = true
	}
	return b516TwoWordViewResult
}
