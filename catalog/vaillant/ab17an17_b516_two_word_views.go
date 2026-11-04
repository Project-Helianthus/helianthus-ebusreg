package vaillant

import (
	ab17an17Binary "encoding/binary"
	ab17an17Math "math"
)

const (
	Ab17An17Family  = "vaillant-b516-family"
	Ab17An17Variant = "b516-op14-candidate"
)

type Ab17An17Result struct {
	Ab17An17DeclaredCount       int
	Ab17An17AvailableExtent     int
	Ab17An17Raw                 []byte
	Ab17An17Trailing            []byte
	Ab17An17Admitted            bool
	Ab17An17StatusAvailable     bool
	Ab17An17Status              byte
	Ab17An17Word0Available      bool
	Ab17An17Word0Raw            [4]byte
	Ab17An17Word0Bits           uint32
	Ab17An17Word0Finite         bool
	Ab17An17Word0Value          *float32
	Ab17An17Word1Available      bool
	Ab17An17Word1Raw            [4]byte
	Ab17An17CompleteTwoWord     bool
	Ab17An17ContractGrade       string
	Ab17An17InterpretationGrade string
	Ab17An17StatusMeaning       string
	Ab17An17Word0NativeUnit     string
	Ab17An17Word0NativeSign     string
	Ab17An17Word1Type           string
	Ab17An17Word1Scale          string
	Ab17An17Word1Sign           string
	Ab17An17Qualification       string
}

func Ab17An17Interpret(ab17an17Family, ab17an17Variant string, ab17an17Operation byte, ab17an17Direction string, ab17an17DeclaredCount int, ab17an17Body []byte) Ab17An17Result {
	ab17an17Result := Ab17An17Result{
		Ab17An17DeclaredCount:       ab17an17DeclaredCount,
		Ab17An17AvailableExtent:     len(ab17an17Body),
		Ab17An17Raw:                 append([]byte(nil), ab17an17Body...),
		Ab17An17ContractGrade:       "Proven",
		Ab17An17InterpretationGrade: "Hypothesis",
		Ab17An17StatusMeaning:       "Unknown",
		Ab17An17Word0NativeUnit:     "Unknown",
		Ab17An17Word0NativeSign:     "Unknown",
		Ab17An17Word1Type:           "Unknown",
		Ab17An17Word1Scale:          "Unknown",
		Ab17An17Word1Sign:           "Unknown",
		Ab17An17Qualification:       "Unknown",
	}
	if len(ab17an17Body) > 9 {
		ab17an17Result.Ab17An17Trailing = append([]byte(nil), ab17an17Body[9:]...)
	}
	if ab17an17Family != Ab17An17Family || ab17an17Variant != Ab17An17Variant || ab17an17Operation != 0x14 || ab17an17Direction != "positive" || ab17an17DeclaredCount < 0 {
		return ab17an17Result
	}
	ab17an17Result.Ab17An17Admitted = true
	if ab17an17DeclaredCount >= 1 && len(ab17an17Body) >= 1 {
		ab17an17Result.Ab17An17StatusAvailable = true
		ab17an17Result.Ab17An17Status = ab17an17Body[0]
	}
	if ab17an17DeclaredCount >= 5 && len(ab17an17Body) >= 5 {
		ab17an17Result.Ab17An17Word0Available = true
		copy(ab17an17Result.Ab17An17Word0Raw[:], ab17an17Body[1:5])
		ab17an17Result.Ab17An17Word0Bits = ab17an17Binary.LittleEndian.Uint32(ab17an17Body[1:5])
		ab17an17Value := ab17an17Math.Float32frombits(ab17an17Result.Ab17An17Word0Bits)
		if !ab17an17Math.IsNaN(float64(ab17an17Value)) && !ab17an17Math.IsInf(float64(ab17an17Value), 0) {
			ab17an17Result.Ab17An17Word0Finite = true
			ab17an17Result.Ab17An17Word0Value = &ab17an17Value
		}
	}
	if ab17an17DeclaredCount >= 9 && len(ab17an17Body) >= 9 {
		ab17an17Result.Ab17An17Word1Available = true
		copy(ab17an17Result.Ab17An17Word1Raw[:], ab17an17Body[5:9])
		ab17an17Result.Ab17An17CompleteTwoWord = true
	}
	return ab17an17Result
}
