package vaillant

import (
	b512HeatpumpIOBytesPackage "bytes"
	b512HeatpumpIOJSON "encoding/json"
	b512HeatpumpIOReflect "reflect"
	b512HeatpumpIOTesting "testing"
)

type b512HeatpumpIORow struct {
	B512HeatpumpIOID          string            `json:"id"`
	B512HeatpumpIOKey         B512HeatpumpIOKey `json:"key"`
	B512HeatpumpIOPosition    int               `json:"position"`
	B512HeatpumpIOWidth       int               `json:"width_bits"`
	B512HeatpumpIOEncoding    string            `json:"encoding"`
	B512HeatpumpIOScale       float64           `json:"scale"`
	B512HeatpumpIOReplacement *struct {
		B512HeatpumpIORaw         uint16 `json:"raw_word"`
		B512HeatpumpIOSelection   string `json:"selection"`
		B512HeatpumpIOMeaning     string `json:"meaning"`
		B512HeatpumpIOUnavailable bool   `json:"unavailable_assertion"`
	} `json:"replacement"`
	B512HeatpumpIOBits []struct {
		B512HeatpumpIOBit    uint8  `json:"bit"`
		B512HeatpumpIOLabel  string `json:"label"`
		B512HeatpumpIOStatus string `json:"status"`
	} `json:"bits"`
	B512HeatpumpIOValues map[string]string `json:"values"`
}

func b512HeatpumpIOPointer[b512HeatpumpIOType any](b512HeatpumpIOValue b512HeatpumpIOType) *b512HeatpumpIOType {
	return &b512HeatpumpIOValue
}
func b512HeatpumpIOEqual(b512HeatpumpIOT *b512HeatpumpIOTesting.T, b512HeatpumpIOGot, b512HeatpumpIOWant any) {
	b512HeatpumpIOT.Helper()
	if !b512HeatpumpIOReflect.DeepEqual(b512HeatpumpIOGot, b512HeatpumpIOWant) {
		b512HeatpumpIOT.Fatalf("got %#v; want %#v", b512HeatpumpIOGot, b512HeatpumpIOWant)
	}
}
func b512HeatpumpIOGet(b512HeatpumpIOT *b512HeatpumpIOTesting.T, b512HeatpumpIOOp uint8, b512HeatpumpIODirection, b512HeatpumpIOField string) (b512HeatpumpIORow, b512HeatpumpIOJSON.RawMessage) {
	b512HeatpumpIOT.Helper()
	b512HeatpumpIORaw, b512HeatpumpIOFound, b512HeatpumpIOErr := B512HeatpumpIOLookup(B512HeatpumpIOKey{B512HeatpumpIOFamily: "vaillant", B512HeatpumpIOVariant: "B512-heatpump-io", B512HeatpumpIOPrimary: 181, B512HeatpumpIOSecondary: b512HeatpumpIOPointer(uint8(18)), B512HeatpumpIOOperation: b512HeatpumpIOPointer(b512HeatpumpIOOp), B512HeatpumpIODirection: b512HeatpumpIODirection, B512HeatpumpIOField: "ab19_" + b512HeatpumpIOField, B512HeatpumpIOAlternative: "none", B512HeatpumpIOInstanceScope: "none"})
	if b512HeatpumpIOErr != nil || !b512HeatpumpIOFound {
		b512HeatpumpIOT.Fatalf("lookup %d/%s/%s: %v %v", b512HeatpumpIOOp, b512HeatpumpIODirection, b512HeatpumpIOField, b512HeatpumpIOFound, b512HeatpumpIOErr)
	}
	var b512HeatpumpIOEntry b512HeatpumpIORow
	if b512HeatpumpIOErr := b512HeatpumpIOJSON.Unmarshal(b512HeatpumpIORaw, &b512HeatpumpIOEntry); b512HeatpumpIOErr != nil {
		b512HeatpumpIOT.Fatal(b512HeatpumpIOErr)
	}
	return b512HeatpumpIOEntry, b512HeatpumpIORaw
}
func b512HeatpumpIOObject(b512HeatpumpIOT *b512HeatpumpIOTesting.T, b512HeatpumpIORaw b512HeatpumpIOJSON.RawMessage) map[string]any {
	b512HeatpumpIOT.Helper()
	var b512HeatpumpIOFields map[string]any
	if b512HeatpumpIOErr := b512HeatpumpIOJSON.Unmarshal(b512HeatpumpIORaw, &b512HeatpumpIOFields); b512HeatpumpIOErr != nil {
		b512HeatpumpIOT.Fatal(b512HeatpumpIOErr)
	}
	return b512HeatpumpIOFields
}
func b512HeatpumpIOMarshal(b512HeatpumpIOT *b512HeatpumpIOTesting.T, b512HeatpumpIOValue any) string {
	b512HeatpumpIOT.Helper()
	b512HeatpumpIORaw, b512HeatpumpIOErr := b512HeatpumpIOJSON.Marshal(b512HeatpumpIOValue)
	if b512HeatpumpIOErr != nil {
		b512HeatpumpIOT.Fatal(b512HeatpumpIOErr)
	}
	return string(b512HeatpumpIORaw)
}

func TestAb19CatalogContract(b512HeatpumpIOT *b512HeatpumpIOTesting.T) {
	b512HeatpumpIORaw, b512HeatpumpIOErr := B512HeatpumpIOEntries()
	if b512HeatpumpIOErr != nil {
		b512HeatpumpIOT.Fatal(b512HeatpumpIOErr)
	}
	var b512HeatpumpIORows []b512HeatpumpIOJSON.RawMessage
	if b512HeatpumpIOErr := b512HeatpumpIOJSON.Unmarshal(b512HeatpumpIORaw, &b512HeatpumpIORows); b512HeatpumpIOErr != nil {
		b512HeatpumpIOT.Fatal(b512HeatpumpIOErr)
	}
	b512HeatpumpIOEqual(b512HeatpumpIOT, len(b512HeatpumpIORows), 34)
	for _, b512HeatpumpIORaw := range b512HeatpumpIORows {
		var b512HeatpumpIOEntry b512HeatpumpIORow
		if b512HeatpumpIOErr := b512HeatpumpIOJSON.Unmarshal(b512HeatpumpIORaw, &b512HeatpumpIOEntry); b512HeatpumpIOErr != nil {
			b512HeatpumpIOT.Fatal(b512HeatpumpIOErr)
		}
		b512HeatpumpIOLookupRaw, b512HeatpumpIOFound, b512HeatpumpIOErr := B512HeatpumpIOLookup(b512HeatpumpIOEntry.B512HeatpumpIOKey)
		if b512HeatpumpIOErr != nil || !b512HeatpumpIOFound {
			b512HeatpumpIOT.Fatalf("whole-row lookup: %v %v", b512HeatpumpIOFound, b512HeatpumpIOErr)
		}
		b512HeatpumpIOEqual(b512HeatpumpIOT, b512HeatpumpIOObject(b512HeatpumpIOT, b512HeatpumpIOLookupRaw), b512HeatpumpIOObject(b512HeatpumpIOT, b512HeatpumpIORaw))
		b512HeatpumpIOFields := b512HeatpumpIOObject(b512HeatpumpIOT, b512HeatpumpIORaw)
		b512HeatpumpIOEqual(b512HeatpumpIOT, b512HeatpumpIOFields["grade"], "Hypothesis")
		b512HeatpumpIOEqual(b512HeatpumpIOT, b512HeatpumpIOFields["metadata_presence"], "Proven")
		b512HeatpumpIOEqual(b512HeatpumpIOT, b512HeatpumpIOFields["native_contract"], "https://github.com/Project-Helianthus/helianthus-docs-ebus/blob/main/protocols/ebus-services/ebus-overview.md")
		b512HeatpumpIOEqual(b512HeatpumpIOT, b512HeatpumpIOFields["unknown"], map[string]any{"native_support": "Unknown", "physical_effects": "Unknown", "persistence": "Unknown", "delivery": "Unknown", "physical_reply_form": "Unknown", "qualification": "Unknown"})
		b512HeatpumpIOEqual(b512HeatpumpIOT, b512HeatpumpIOFields["raw"], map[string]any{"retained": true, "short": "retain_raw", "extra": "retain_raw", "physical_reply_selection": false})
	}
}

func TestAb19IdentityAxes(b512HeatpumpIOT *b512HeatpumpIOTesting.T) {
	b512HeatpumpIOBase, b512HeatpumpIORaw := b512HeatpumpIOGet(b512HeatpumpIOT, 15, "request", "outputs")
	b512HeatpumpIOChanges := []struct {
		b512HeatpumpIOName   string
		b512HeatpumpIOChange func(*B512HeatpumpIOKey)
	}{
		{"ab19Family", func(b512HeatpumpIOKey *B512HeatpumpIOKey) {
			b512HeatpumpIOKey.B512HeatpumpIOFamily = "synthetic-family"
		}},
		{"ab19Variant", func(b512HeatpumpIOKey *B512HeatpumpIOKey) {
			b512HeatpumpIOKey.B512HeatpumpIOVariant = "synthetic-variant"
		}},
		{"ab19Primary", func(b512HeatpumpIOKey *B512HeatpumpIOKey) { b512HeatpumpIOKey.B512HeatpumpIOPrimary = 1 }},
		{"ab19Secondary", func(b512HeatpumpIOKey *B512HeatpumpIOKey) {
			b512HeatpumpIOKey.B512HeatpumpIOSecondary = b512HeatpumpIOPointer(uint8(0))
		}},
		{"ab19Operation", func(b512HeatpumpIOKey *B512HeatpumpIOKey) {
			b512HeatpumpIOKey.B512HeatpumpIOOperation = b512HeatpumpIOPointer(uint8(0))
		}},
		{"ab19Direction", func(b512HeatpumpIOKey *B512HeatpumpIOKey) { b512HeatpumpIOKey.B512HeatpumpIODirection = "positive" }},
		{"ab19Field", func(b512HeatpumpIOKey *B512HeatpumpIOKey) { b512HeatpumpIOKey.B512HeatpumpIOField = "ab19_unlisted" }},
		{"ab19Alternative", func(b512HeatpumpIOKey *B512HeatpumpIOKey) {
			b512HeatpumpIOKey.B512HeatpumpIOAlternative = "ab19_alternative"
		}},
		{"ab19Group", func(b512HeatpumpIOKey *B512HeatpumpIOKey) {
			b512HeatpumpIOKey.B512HeatpumpIOGroup = b512HeatpumpIOPointer(uint8(0))
		}},
		{"ab19InstanceScope", func(b512HeatpumpIOKey *B512HeatpumpIOKey) { b512HeatpumpIOKey.B512HeatpumpIOInstanceScope = "dynamic" }},
		{"ab19Register", func(b512HeatpumpIOKey *B512HeatpumpIOKey) {
			b512HeatpumpIOKey.B512HeatpumpIORegister = b512HeatpumpIOPointer(uint16(0))
		}},
	}
	for _, b512HeatpumpIOCase := range b512HeatpumpIOChanges {
		b512HeatpumpIOT.Run(b512HeatpumpIOCase.b512HeatpumpIOName, func(b512HeatpumpIOT *b512HeatpumpIOTesting.T) {
			b512HeatpumpIOKey := b512HeatpumpIOBase.B512HeatpumpIOKey
			b512HeatpumpIOCase.b512HeatpumpIOChange(&b512HeatpumpIOKey)
			if _, b512HeatpumpIOFound, b512HeatpumpIOErr := B512HeatpumpIOLookup(b512HeatpumpIOKey); b512HeatpumpIOFound || b512HeatpumpIOErr != nil {
				b512HeatpumpIOT.Fatalf("mismatched axis admitted: %v %v", b512HeatpumpIOFound, b512HeatpumpIOErr)
			}
			b512HeatpumpIOFirst := b512HeatpumpIOObject(b512HeatpumpIOT, b512HeatpumpIORaw)
			b512HeatpumpIOSecond := b512HeatpumpIOObject(b512HeatpumpIOT, b512HeatpumpIORaw)
			b512HeatpumpIOSecond["key"] = b512HeatpumpIOKey
			b512HeatpumpIOSecond["id"] = "ab19_control"
			b512HeatpumpIOInput := b512HeatpumpIOMarshal(b512HeatpumpIOT, []any{b512HeatpumpIOFirst, b512HeatpumpIOSecond})
			for _, b512HeatpumpIOControl := range []struct {
				b512HeatpumpIOKey B512HeatpumpIOKey
				b512HeatpumpIOID  string
			}{{b512HeatpumpIOBase.B512HeatpumpIOKey, b512HeatpumpIOBase.B512HeatpumpIOID}, {b512HeatpumpIOKey, "ab19_control"}} {
				b512HeatpumpIOGot, b512HeatpumpIOFound, b512HeatpumpIOErr := b512HeatpumpIOLookup(b512HeatpumpIOInput, b512HeatpumpIOControl.b512HeatpumpIOKey)
				if !b512HeatpumpIOFound || b512HeatpumpIOErr != nil {
					b512HeatpumpIOT.Fatalf("valid control: %v %v", b512HeatpumpIOFound, b512HeatpumpIOErr)
				}
				b512HeatpumpIOEqual(b512HeatpumpIOT, b512HeatpumpIOObject(b512HeatpumpIOT, b512HeatpumpIOGot)["id"], b512HeatpumpIOControl.b512HeatpumpIOID)
			}
		})
	}
	for _, b512HeatpumpIOAxis := range []string{"secondary", "operation", "group", "register"} {
		b512HeatpumpIOT.Run("ab19AbsentZero_"+b512HeatpumpIOAxis, func(b512HeatpumpIOT *b512HeatpumpIOTesting.T) {
			b512HeatpumpIOFirst := b512HeatpumpIOObject(b512HeatpumpIOT, b512HeatpumpIORaw)
			b512HeatpumpIOSecond := b512HeatpumpIOObject(b512HeatpumpIOT, b512HeatpumpIORaw)
			b512HeatpumpIOFirst["key"].(map[string]any)[b512HeatpumpIOAxis] = nil
			b512HeatpumpIOSecond["key"].(map[string]any)[b512HeatpumpIOAxis] = float64(0)
			b512HeatpumpIOSecond["id"] = "ab19_zero"
			b512HeatpumpIOInput := b512HeatpumpIOMarshal(b512HeatpumpIOT, []any{b512HeatpumpIOFirst, b512HeatpumpIOSecond})
			for _, b512HeatpumpIOControl := range []map[string]any{b512HeatpumpIOFirst, b512HeatpumpIOSecond} {
				var b512HeatpumpIOKey B512HeatpumpIOKey
				if b512HeatpumpIOErr := b512HeatpumpIOJSON.Unmarshal([]byte(b512HeatpumpIOMarshal(b512HeatpumpIOT, b512HeatpumpIOControl["key"])), &b512HeatpumpIOKey); b512HeatpumpIOErr != nil {
					b512HeatpumpIOT.Fatal(b512HeatpumpIOErr)
				}
				b512HeatpumpIOGot, b512HeatpumpIOFound, b512HeatpumpIOErr := b512HeatpumpIOLookup(b512HeatpumpIOInput, b512HeatpumpIOKey)
				if !b512HeatpumpIOFound || b512HeatpumpIOErr != nil {
					b512HeatpumpIOT.Fatalf("absent/zero: %v %v", b512HeatpumpIOFound, b512HeatpumpIOErr)
				}
				b512HeatpumpIOEqual(b512HeatpumpIOT, b512HeatpumpIOObject(b512HeatpumpIOT, b512HeatpumpIOGot)["id"], b512HeatpumpIOControl["id"])
			}
		})
	}
}

func TestAb19ValidationAndIsolation(b512HeatpumpIOT *b512HeatpumpIOTesting.T) {
	b512HeatpumpIOBase, b512HeatpumpIORaw := b512HeatpumpIOGet(b512HeatpumpIOT, 15, "request", "outputs")
	b512HeatpumpIODuplicate := b512HeatpumpIOObject(b512HeatpumpIOT, b512HeatpumpIORaw)
	b512HeatpumpIODuplicate["id"] = "ab19_duplicate"
	b512HeatpumpIOInvalid := []string{"{", "null", "{}", "[]", "[null]", string(b512HeatpumpIORaw), "[" + string(b512HeatpumpIORaw) + "] true", b512HeatpumpIOMarshal(b512HeatpumpIOT, []any{b512HeatpumpIOObject(b512HeatpumpIOT, b512HeatpumpIORaw), b512HeatpumpIODuplicate})}
	for _, b512HeatpumpIOAxis := range []string{"family", "variant", "primary", "secondary", "operation", "direction", "field", "alternative", "group", "instance_scope", "register"} {
		b512HeatpumpIOMissing := b512HeatpumpIOObject(b512HeatpumpIOT, b512HeatpumpIORaw)
		delete(b512HeatpumpIOMissing["key"].(map[string]any), b512HeatpumpIOAxis)
		b512HeatpumpIOInvalid = append(b512HeatpumpIOInvalid, b512HeatpumpIOMarshal(b512HeatpumpIOT, []any{b512HeatpumpIOMissing}))
	}
	for _, b512HeatpumpIOInput := range b512HeatpumpIOInvalid {
		if _, b512HeatpumpIOErr := b512HeatpumpIOEntries(b512HeatpumpIOInput); b512HeatpumpIOErr == nil {
			b512HeatpumpIOT.Fatal("invalid catalog admitted by entries")
		}
		if _, b512HeatpumpIOFound, b512HeatpumpIOErr := b512HeatpumpIOLookup(b512HeatpumpIOInput, b512HeatpumpIOBase.B512HeatpumpIOKey); b512HeatpumpIOErr == nil || b512HeatpumpIOFound {
			b512HeatpumpIOT.Fatal("invalid catalog admitted by lookup")
		}
	}
	b512HeatpumpIOAll, b512HeatpumpIOErr := B512HeatpumpIOEntries()
	if b512HeatpumpIOErr != nil {
		b512HeatpumpIOT.Fatal(b512HeatpumpIOErr)
	}
	b512HeatpumpIOOriginal := append([]byte(nil), b512HeatpumpIOAll...)
	for b512HeatpumpIOIndex := range b512HeatpumpIOAll {
		b512HeatpumpIOAll[b512HeatpumpIOIndex] = 'x'
	}
	for b512HeatpumpIOIndex := range b512HeatpumpIORaw {
		b512HeatpumpIORaw[b512HeatpumpIOIndex] = 'x'
	}
	b512HeatpumpIOAgain, b512HeatpumpIOErr := B512HeatpumpIOEntries()
	if b512HeatpumpIOErr != nil {
		b512HeatpumpIOT.Fatal(b512HeatpumpIOErr)
	}
	b512HeatpumpIOEqual(b512HeatpumpIOT, b512HeatpumpIOBytesPackage.Equal(b512HeatpumpIOAgain, b512HeatpumpIOOriginal), true)
	_, b512HeatpumpIOFresh := b512HeatpumpIOGet(b512HeatpumpIOT, 15, "request", "outputs")
	b512HeatpumpIOLookupOriginal := append([]byte(nil), b512HeatpumpIOFresh...)
	for b512HeatpumpIOIndex := range b512HeatpumpIOFresh {
		b512HeatpumpIOFresh[b512HeatpumpIOIndex] = 'x'
	}
	_, b512HeatpumpIOFresh = b512HeatpumpIOGet(b512HeatpumpIOT, 15, "request", "outputs")
	b512HeatpumpIOEqual(b512HeatpumpIOT, b512HeatpumpIOBytesPackage.Equal(b512HeatpumpIOFresh, b512HeatpumpIOLookupOriginal), true)
}

func TestAb19FieldCoordinates(b512HeatpumpIOT *b512HeatpumpIOTesting.T) {
	for _, b512HeatpumpIOCase := range []struct {
		b512HeatpumpIOOp                             uint8
		b512HeatpumpIODirection, b512HeatpumpIOField string
		b512HeatpumpIOPosition, b512HeatpumpIOWidth  int
		b512HeatpumpIOEncoding                       string
	}{
		{15, "request", "outputs", 3, 8, "bits"}, {15, "request", "valve", 4, 8, "u8"},
		{15, "positive", "temperature_one", 3, 16, "signed-le16"}, {15, "positive", "temperature_two", 6, 16, "signed-le16"},
		{15, "positive", "input", 5, 8, "u8"}, {15, "positive", "pressure", 8, 8, "u8"}, {15, "positive", "backup", 9, 8, "u8"},
		{16, "request", "actuators_one", 3, 8, "bits"}, {16, "request", "actuators_two", 4, 8, "bits"}, {16, "request", "inverse_actuator", 5, 16, "opaque-16"},
		{16, "positive", "status", 3, 8, "bits"}, {16, "positive", "single_phase", 4, 8, "u8"}, {16, "positive", "three_phase", 5, 8, "u8"},
		{16, "positive", "adc_1", 6, 16, "raw-le16"}, {16, "positive", "adc_2", 8, 16, "raw-le16"}, {16, "positive", "adc_3", 10, 16, "raw-le16"},
		{16, "positive", "adc_4", 12, 16, "raw-le16"}, {16, "positive", "adc_5", 14, 16, "raw-le16"}, {16, "positive", "adc_6", 16, 16, "raw-le16"},
		{16, "positive", "board", 18, 8, "u8"},
		{17, "request", "temperature_1", 3, 16, "signed-le16"}, {17, "request", "temperature_2", 5, 16, "signed-le16"}, {17, "request", "temperature_3", 7, 16, "signed-le16"},
		{17, "request", "pressure", 9, 8, "u8"}, {17, "request", "flow", 10, 8, "signed8"}, {17, "request", "pump_status", 11, 8, "raw-u8"}, {17, "request", "compressor", 12, 8, "bits"},
		{17, "positive", "status", 3, 8, "bits"}, {17, "positive", "setpoint", 4, 8, "u8"},
		{18, "request", "emergency", 3, 8, "bits"}, {18, "request", "switch", 4, 8, "raw-u8"}, {18, "request", "co2_1", 5, 16, "raw-le16"}, {18, "request", "co2_2", 7, 16, "raw-le16"},
		{18, "positive", "led", 3, 8, "bits"},
	} {
		b512HeatpumpIOT.Run("ab19_"+b512HeatpumpIOCase.b512HeatpumpIODirection+"_"+b512HeatpumpIOCase.b512HeatpumpIOField+"_"+string(rune('A'+b512HeatpumpIOCase.b512HeatpumpIOOp)), func(b512HeatpumpIOT *b512HeatpumpIOTesting.T) {
			b512HeatpumpIORow, b512HeatpumpIORaw := b512HeatpumpIOGet(b512HeatpumpIOT, b512HeatpumpIOCase.b512HeatpumpIOOp, b512HeatpumpIOCase.b512HeatpumpIODirection, b512HeatpumpIOCase.b512HeatpumpIOField)
			b512HeatpumpIOEqual(b512HeatpumpIOT, b512HeatpumpIORow.B512HeatpumpIOPosition, b512HeatpumpIOCase.b512HeatpumpIOPosition)
			b512HeatpumpIOEqual(b512HeatpumpIOT, b512HeatpumpIORow.B512HeatpumpIOWidth, b512HeatpumpIOCase.b512HeatpumpIOWidth)
			b512HeatpumpIOEqual(b512HeatpumpIOT, b512HeatpumpIORow.B512HeatpumpIOEncoding, b512HeatpumpIOCase.b512HeatpumpIOEncoding)
			if b512HeatpumpIOCase.b512HeatpumpIOOp != 15 || b512HeatpumpIOCase.b512HeatpumpIOField != "pressure" {
				if _, b512HeatpumpIOPresent := b512HeatpumpIOObject(b512HeatpumpIOT, b512HeatpumpIORaw)["default"]; b512HeatpumpIOPresent {
					b512HeatpumpIOT.Fatal("synthesized default")
				}
			}
		})
	}
}

func TestAb19Bits(b512HeatpumpIOT *b512HeatpumpIOTesting.T) {
	for _, b512HeatpumpIOCase := range []struct {
		b512HeatpumpIOOp                             uint8
		b512HeatpumpIODirection, b512HeatpumpIOField string
		b512HeatpumpIOLabels                         []string
		b512HeatpumpIOUnmodeled                      bool
	}{
		{15, "request", "outputs", []string{"output one", "output two", "reserved", "reserved", "reserved", "reserved", "reserved", "reserved"}, false},
		{16, "request", "actuators_one", []string{"cooling valve", "passive mixer one", "passive mixer two", "cooling active", "well pump", "backup safety", "backup one", "backup two"}, false},
		{16, "request", "actuators_two", []string{"backup three", "external valve", "external backup", "reserved", "reserved", "reserved", "reserved", "reserved"}, false},
		{16, "positive", "status", []string{"anode channel one", "anode channel two", "smart-grid", "well", "safety limiter", "safety", "reserved", "reserved"}, false},
		{17, "request", "compressor", []string{"compressor-off", "reserved", "reserved", "reserved", "reserved", "reserved", "reserved", "unmodeled"}, true},
		{17, "positive", "status", []string{"compressor-off", "pump-relay", "reserved", "reserved", "reserved", "reserved", "reserved", "reserved"}, false},
		{18, "request", "emergency", []string{"emergency", "reserved", "reserved", "reserved", "reserved", "reserved", "reserved", "unmodeled"}, true},
		{18, "positive", "led", []string{"LED", "reserved", "reserved", "reserved", "reserved", "reserved", "reserved", "unmodeled"}, true},
	} {
		b512HeatpumpIORow, _ := b512HeatpumpIOGet(b512HeatpumpIOT, b512HeatpumpIOCase.b512HeatpumpIOOp, b512HeatpumpIOCase.b512HeatpumpIODirection, b512HeatpumpIOCase.b512HeatpumpIOField)
		b512HeatpumpIOEqual(b512HeatpumpIOT, len(b512HeatpumpIORow.B512HeatpumpIOBits), 8)
		for b512HeatpumpIOIndex, b512HeatpumpIOLabel := range b512HeatpumpIOCase.b512HeatpumpIOLabels {
			b512HeatpumpIOBit := b512HeatpumpIORow.B512HeatpumpIOBits[b512HeatpumpIOIndex]
			b512HeatpumpIOEqual(b512HeatpumpIOT, b512HeatpumpIOBit.B512HeatpumpIOBit, uint8(b512HeatpumpIOIndex))
			b512HeatpumpIOEqual(b512HeatpumpIOT, b512HeatpumpIOBit.B512HeatpumpIOLabel, b512HeatpumpIOLabel)
			b512HeatpumpIOStatus := "candidate"
			if b512HeatpumpIOLabel == "reserved" || b512HeatpumpIOLabel == "unmodeled" {
				b512HeatpumpIOStatus = b512HeatpumpIOLabel
			}
			b512HeatpumpIOEqual(b512HeatpumpIOT, b512HeatpumpIOBit.B512HeatpumpIOStatus, b512HeatpumpIOStatus)
			for b512HeatpumpIOToggle := 0; b512HeatpumpIOToggle < 8; b512HeatpumpIOToggle++ {
				b512HeatpumpIOSample := uint8(1 << b512HeatpumpIOToggle)
				b512HeatpumpIOEqual(b512HeatpumpIOT, b512HeatpumpIOSample&(1<<b512HeatpumpIOBit.B512HeatpumpIOBit) != 0, b512HeatpumpIOToggle == b512HeatpumpIOIndex)
			}
		}
		if b512HeatpumpIOCase.b512HeatpumpIOUnmodeled {
			b512HeatpumpIOEqual(b512HeatpumpIOT, b512HeatpumpIORow.B512HeatpumpIOBits[7].B512HeatpumpIOStatus, "unmodeled")
		}
	}
}

func b512HeatpumpIOEvaluate(b512HeatpumpIORow b512HeatpumpIORow, b512HeatpumpIORaw []byte) (float64, bool, []byte) {
	b512HeatpumpIORetained := append([]byte(nil), b512HeatpumpIORaw...)
	if len(b512HeatpumpIORaw)*8 != b512HeatpumpIORow.B512HeatpumpIOWidth {
		return 0, false, b512HeatpumpIORetained
	}
	b512HeatpumpIOWord := uint16(b512HeatpumpIORaw[0])
	if len(b512HeatpumpIORaw) == 2 {
		b512HeatpumpIOWord |= uint16(b512HeatpumpIORaw[1]) << 8
	}
	if b512HeatpumpIORow.B512HeatpumpIOReplacement != nil && b512HeatpumpIORow.B512HeatpumpIOReplacement.B512HeatpumpIOSelection == "before_scaling" && b512HeatpumpIOWord == b512HeatpumpIORow.B512HeatpumpIOReplacement.B512HeatpumpIORaw {
		return 0, true, b512HeatpumpIORetained
	}
	b512HeatpumpIOValue := float64(b512HeatpumpIOWord)
	if b512HeatpumpIORow.B512HeatpumpIOEncoding == "signed-le16" {
		b512HeatpumpIOValue = float64(int16(b512HeatpumpIOWord))
	}
	if b512HeatpumpIORow.B512HeatpumpIOEncoding == "signed8" {
		b512HeatpumpIOValue = float64(int8(b512HeatpumpIOWord))
	}
	if b512HeatpumpIORow.B512HeatpumpIOScale != 0 {
		b512HeatpumpIOValue *= b512HeatpumpIORow.B512HeatpumpIOScale
	}
	return b512HeatpumpIOValue, false, b512HeatpumpIORetained
}

func TestAb19NumericCandidates(b512HeatpumpIOT *b512HeatpumpIOTesting.T) {
	for _, b512HeatpumpIOField := range []string{"temperature_one", "temperature_two"} {
		b512HeatpumpIORow, _ := b512HeatpumpIOGet(b512HeatpumpIOT, 15, "positive", b512HeatpumpIOField)
		b512HeatpumpIOEqual(b512HeatpumpIOT, b512HeatpumpIORow.B512HeatpumpIOScale, 0.0625)
		if b512HeatpumpIORow.B512HeatpumpIOReplacement == nil {
			b512HeatpumpIOT.Fatal("missing replacement")
		}
		b512HeatpumpIOEqual(b512HeatpumpIOT, b512HeatpumpIORow.B512HeatpumpIOReplacement.B512HeatpumpIOSelection, "before_scaling")
		b512HeatpumpIOEqual(b512HeatpumpIOT, b512HeatpumpIORow.B512HeatpumpIOReplacement.B512HeatpumpIOMeaning, "replacement")
		b512HeatpumpIOEqual(b512HeatpumpIOT, b512HeatpumpIORow.B512HeatpumpIOReplacement.B512HeatpumpIOUnavailable, false)
		b512HeatpumpIOValue, b512HeatpumpIOReplacement, b512HeatpumpIORaw := b512HeatpumpIOEvaluate(b512HeatpumpIORow, []byte{0, 128})
		b512HeatpumpIOEqual(b512HeatpumpIOT, b512HeatpumpIOValue, float64(0))
		b512HeatpumpIOEqual(b512HeatpumpIOT, b512HeatpumpIOReplacement, true)
		b512HeatpumpIOEqual(b512HeatpumpIOT, b512HeatpumpIORaw, []byte{0, 128})
		b512HeatpumpIOValue, b512HeatpumpIOReplacement, b512HeatpumpIORaw = b512HeatpumpIOEvaluate(b512HeatpumpIORow, []byte{240, 255})
		b512HeatpumpIOEqual(b512HeatpumpIOT, b512HeatpumpIOValue, float64(-1))
		b512HeatpumpIOEqual(b512HeatpumpIOT, b512HeatpumpIOReplacement, false)
		b512HeatpumpIOEqual(b512HeatpumpIOT, b512HeatpumpIORaw, []byte{240, 255})
	}
	b512HeatpumpIOPressure, b512HeatpumpIORaw := b512HeatpumpIOGet(b512HeatpumpIOT, 15, "positive", "pressure")
	b512HeatpumpIOFields := b512HeatpumpIOObject(b512HeatpumpIOT, b512HeatpumpIORaw)
	b512HeatpumpIOEqual(b512HeatpumpIOT, b512HeatpumpIOFields["default"], float64(255))
	b512HeatpumpIOEqual(b512HeatpumpIOT, b512HeatpumpIOFields["default_conflict"], true)
	b512HeatpumpIOEqual(b512HeatpumpIOT, b512HeatpumpIOFields["domain"], []any{float64(0), 25.5})
	b512HeatpumpIOEqual(b512HeatpumpIOT, b512HeatpumpIOFields["raw_domain"], []any{float64(0), float64(255)})
	b512HeatpumpIOValue, b512HeatpumpIOReplacement, b512HeatpumpIOBytes := b512HeatpumpIOEvaluate(b512HeatpumpIOPressure, []byte{255})
	b512HeatpumpIOEqual(b512HeatpumpIOT, b512HeatpumpIOValue, 25.5)
	b512HeatpumpIOEqual(b512HeatpumpIOT, b512HeatpumpIOReplacement, false)
	b512HeatpumpIOEqual(b512HeatpumpIOT, b512HeatpumpIOBytes, []byte{255})
	for _, b512HeatpumpIOField := range []string{"temperature_1", "temperature_2", "temperature_3"} {
		b512HeatpumpIORow, _ := b512HeatpumpIOGet(b512HeatpumpIOT, 17, "request", b512HeatpumpIOField)
		b512HeatpumpIOValue, b512HeatpumpIOReplacement, _ := b512HeatpumpIOEvaluate(b512HeatpumpIORow, []byte{240, 255})
		b512HeatpumpIOEqual(b512HeatpumpIOT, b512HeatpumpIOValue, float64(-1))
		b512HeatpumpIOEqual(b512HeatpumpIOT, b512HeatpumpIOReplacement, false)
		if b512HeatpumpIORow.B512HeatpumpIOReplacement != nil {
			b512HeatpumpIOT.Fatal("synthesized replacement")
		}
	}
	b512HeatpumpIOFlow, _ := b512HeatpumpIOGet(b512HeatpumpIOT, 17, "request", "flow")
	b512HeatpumpIOValue, b512HeatpumpIOReplacement, _ = b512HeatpumpIOEvaluate(b512HeatpumpIOFlow, []byte{255})
	b512HeatpumpIOEqual(b512HeatpumpIOT, b512HeatpumpIOValue, float64(-10))
	b512HeatpumpIOEqual(b512HeatpumpIOT, b512HeatpumpIOReplacement, false)
	for _, b512HeatpumpIOField := range []string{"co2_1", "co2_2"} {
		b512HeatpumpIORow, b512HeatpumpIORaw := b512HeatpumpIOGet(b512HeatpumpIOT, 18, "request", b512HeatpumpIOField)
		b512HeatpumpIOValue, b512HeatpumpIOReplacement, b512HeatpumpIOBytes = b512HeatpumpIOEvaluate(b512HeatpumpIORow, []byte{255, 255})
		b512HeatpumpIOEqual(b512HeatpumpIOT, b512HeatpumpIOValue, float64(65535))
		b512HeatpumpIOEqual(b512HeatpumpIOT, b512HeatpumpIOReplacement, false)
		b512HeatpumpIOEqual(b512HeatpumpIOT, b512HeatpumpIOBytes, []byte{255, 255})
		b512HeatpumpIOEqual(b512HeatpumpIOT, b512HeatpumpIOObject(b512HeatpumpIOT, b512HeatpumpIORaw)["domain"], []any{float64(0), float64(65535)})
	}
}

func TestAb19OpaqueAndPhaseCandidates(b512HeatpumpIOT *b512HeatpumpIOTesting.T) {
	for _, b512HeatpumpIOField := range []string{"adc_1", "adc_2", "adc_3", "adc_4", "adc_5", "adc_6"} {
		b512HeatpumpIORow, _ := b512HeatpumpIOGet(b512HeatpumpIOT, 16, "positive", b512HeatpumpIOField)
		for _, b512HeatpumpIOSample := range []struct {
			b512HeatpumpIOBytes []byte
			b512HeatpumpIOWant  float64
		}{{[]byte{0, 0}, 0}, {[]byte{1, 0}, 1}, {[]byte{0, 1}, 256}, {[]byte{255, 255}, 65535}} {
			b512HeatpumpIOValue, b512HeatpumpIOReplacement, b512HeatpumpIORetained := b512HeatpumpIOEvaluate(b512HeatpumpIORow, b512HeatpumpIOSample.b512HeatpumpIOBytes)
			b512HeatpumpIOEqual(b512HeatpumpIOT, b512HeatpumpIOValue, b512HeatpumpIOSample.b512HeatpumpIOWant)
			b512HeatpumpIOEqual(b512HeatpumpIOT, b512HeatpumpIOReplacement, false)
			b512HeatpumpIOEqual(b512HeatpumpIOT, b512HeatpumpIORetained, b512HeatpumpIOSample.b512HeatpumpIOBytes)
		}
	}
	b512HeatpumpIOInverse, b512HeatpumpIORaw := b512HeatpumpIOGet(b512HeatpumpIOT, 16, "request", "inverse_actuator")
	b512HeatpumpIOEqual(b512HeatpumpIOT, b512HeatpumpIOObject(b512HeatpumpIOT, b512HeatpumpIORaw)["algorithm"], "Unknown")
	for _, b512HeatpumpIOSample := range [][]byte{{255, 255}, {0, 0}, {170, 85}, {1}, {1, 2, 3}} {
		_, b512HeatpumpIOReplacement, b512HeatpumpIORetained := b512HeatpumpIOEvaluate(b512HeatpumpIOInverse, b512HeatpumpIOSample)
		b512HeatpumpIOEqual(b512HeatpumpIOT, b512HeatpumpIOReplacement, false)
		b512HeatpumpIOEqual(b512HeatpumpIOT, b512HeatpumpIORetained, b512HeatpumpIOSample)
	}
	b512HeatpumpIOBoard, b512HeatpumpIORaw := b512HeatpumpIOGet(b512HeatpumpIOT, 16, "positive", "board")
	b512HeatpumpIOEqual(b512HeatpumpIOT, b512HeatpumpIOBoard.B512HeatpumpIOValues, map[string]string{"0": "opaque", "1": "opaque", "2": "opaque", "3": "opaque", "11": "opaque", "12": "opaque"})
	if _, b512HeatpumpIOPresent := b512HeatpumpIOBoard.B512HeatpumpIOValues["13"]; b512HeatpumpIOPresent {
		b512HeatpumpIOT.Fatal("unmapped board code inferred")
	}
	b512HeatpumpIOEqual(b512HeatpumpIOT, b512HeatpumpIOObject(b512HeatpumpIOT, b512HeatpumpIORaw)["unmapped"], "retain_raw")
	for _, b512HeatpumpIOSample := range []byte{12, 13} {
		b512HeatpumpIOValue, _, b512HeatpumpIORetained := b512HeatpumpIOEvaluate(b512HeatpumpIOBoard, []byte{b512HeatpumpIOSample})
		b512HeatpumpIOEqual(b512HeatpumpIOT, b512HeatpumpIOValue, float64(b512HeatpumpIOSample))
		b512HeatpumpIOEqual(b512HeatpumpIOT, b512HeatpumpIORetained, []byte{b512HeatpumpIOSample})
	}
	_, b512HeatpumpIORaw = b512HeatpumpIOGet(b512HeatpumpIOT, 16, "positive", "single_phase")
	b512HeatpumpIOEqual(b512HeatpumpIOT, b512HeatpumpIOObject(b512HeatpumpIOT, b512HeatpumpIORaw)["relations"], []any{map[string]any{"raw": float64(50), "comparison": "approximately", "meaning": "powered"}, map[string]any{"raw": float64(0), "comparison": "exact", "meaning": "absent"}})
	b512HeatpumpIOEqual(b512HeatpumpIOT, b512HeatpumpIOObject(b512HeatpumpIOT, b512HeatpumpIORaw)["tolerance"], "Unknown")
	b512HeatpumpIOPhase, _ := b512HeatpumpIOGet(b512HeatpumpIOT, 16, "positive", "three_phase")
	b512HeatpumpIOEqual(b512HeatpumpIOT, b512HeatpumpIOPhase.B512HeatpumpIOValues, map[string]string{"100": "powered", "0": "all phases absent", "33": "one phase absent", "66": "two phases absent"})
	for _, b512HeatpumpIOField := range []string{"single_phase", "three_phase"} {
		b512HeatpumpIORow, _ := b512HeatpumpIOGet(b512HeatpumpIOT, 16, "positive", b512HeatpumpIOField)
		for _, b512HeatpumpIOSample := range []byte{0, 33, 50, 66, 100} {
			b512HeatpumpIOValue, _, b512HeatpumpIORetained := b512HeatpumpIOEvaluate(b512HeatpumpIORow, []byte{b512HeatpumpIOSample})
			b512HeatpumpIOEqual(b512HeatpumpIOT, b512HeatpumpIOValue, float64(b512HeatpumpIOSample))
			b512HeatpumpIOEqual(b512HeatpumpIOT, b512HeatpumpIORetained, []byte{b512HeatpumpIOSample})
		}
	}
}

func TestAb19EnumsAndPositiveConflict(b512HeatpumpIOT *b512HeatpumpIOTesting.T) {
	b512HeatpumpIOValve, _ := b512HeatpumpIOGet(b512HeatpumpIOT, 15, "request", "valve")
	b512HeatpumpIOEqual(b512HeatpumpIOT, b512HeatpumpIOValve.B512HeatpumpIOValues, map[string]string{"0": "off", "1": "CH", "2": "DHW"})
	b512HeatpumpIOInput, _ := b512HeatpumpIOGet(b512HeatpumpIOT, 15, "positive", "input")
	b512HeatpumpIOEqual(b512HeatpumpIOT, b512HeatpumpIOInput.B512HeatpumpIOValues, map[string]string{"0": "multipurpose", "1": "utility contact"})
	b512HeatpumpIOBackup, _ := b512HeatpumpIOGet(b512HeatpumpIOT, 15, "positive", "backup")
	b512HeatpumpIOEqual(b512HeatpumpIOT, b512HeatpumpIOBackup.B512HeatpumpIOValues, map[string]string{"0": "off", "1": "step one", "2": "step two", "3": "step three"})
	if b512HeatpumpIOBackup.B512HeatpumpIOReplacement == nil {
		b512HeatpumpIOT.Fatal("backup replacement missing")
	}
	b512HeatpumpIOEqual(b512HeatpumpIOT, b512HeatpumpIOBackup.B512HeatpumpIOReplacement.B512HeatpumpIOMeaning, "replacement")
	b512HeatpumpIOEqual(b512HeatpumpIOT, b512HeatpumpIOBackup.B512HeatpumpIOReplacement.B512HeatpumpIOUnavailable, false)
	_, b512HeatpumpIOReplacement, b512HeatpumpIORetained := b512HeatpumpIOEvaluate(b512HeatpumpIOBackup, []byte{255})
	b512HeatpumpIOEqual(b512HeatpumpIOT, b512HeatpumpIOReplacement, true)
	b512HeatpumpIOEqual(b512HeatpumpIOT, b512HeatpumpIORetained, []byte{255})
	for _, b512HeatpumpIOCase := range []struct {
		b512HeatpumpIOOp    uint8
		b512HeatpumpIOField string
	}{{17, "status"}, {17, "setpoint"}, {18, "led"}} {
		b512HeatpumpIORow, b512HeatpumpIORaw := b512HeatpumpIOGet(b512HeatpumpIOT, b512HeatpumpIOCase.b512HeatpumpIOOp, "positive", b512HeatpumpIOCase.b512HeatpumpIOField)
		b512HeatpumpIOEqual(b512HeatpumpIOT, *b512HeatpumpIORow.B512HeatpumpIOKey.B512HeatpumpIOOperation, b512HeatpumpIOCase.b512HeatpumpIOOp)
		b512HeatpumpIOEqual(b512HeatpumpIOT, b512HeatpumpIOObject(b512HeatpumpIOT, b512HeatpumpIORaw)["declared_positive"], map[string]any{"primary": float64(181), "secondary": float64(18), "operation": float64(12), "selector_conflict": true, "decoder_selection": false, "acceptance": false})
		for _, b512HeatpumpIOStructure := range [][]byte{{}, {0}, {0, 0, 0, 0, 0, 0, 0, 0, 0}} {
			b512HeatpumpIOEqual(b512HeatpumpIOT, b512HeatpumpIOObject(b512HeatpumpIOT, b512HeatpumpIORaw)["raw"], map[string]any{"retained": true, "short": "retain_raw", "extra": "retain_raw", "physical_reply_selection": false})
			_, _, b512HeatpumpIORetained := b512HeatpumpIOEvaluate(b512HeatpumpIORow, b512HeatpumpIOStructure)
			b512HeatpumpIOEqual(b512HeatpumpIOT, b512HeatpumpIOBytesPackage.Equal(b512HeatpumpIORetained, b512HeatpumpIOStructure), true)
		}
	}
	_, b512HeatpumpIORaw := b512HeatpumpIOGet(b512HeatpumpIOT, 17, "positive", "setpoint")
	b512HeatpumpIOEqual(b512HeatpumpIOT, b512HeatpumpIOObject(b512HeatpumpIOT, b512HeatpumpIORaw)["domain"], []any{float64(0), float64(100)})
}

func TestAb19UnitsAndScales(b512HeatpumpIOT *b512HeatpumpIOTesting.T) {
	for _, b512HeatpumpIOCase := range []struct {
		b512HeatpumpIOOp                                                 uint8
		b512HeatpumpIODirection, b512HeatpumpIOField, b512HeatpumpIOUnit string
		b512HeatpumpIOScale                                              float64
	}{
		{15, "positive", "temperature_one", "degrees Celsius", 0.0625},
		{15, "positive", "temperature_two", "degrees Celsius", 0.0625},
		{15, "positive", "pressure", "bar", 0.1},
		{16, "positive", "single_phase", "percent", 0},
		{16, "positive", "three_phase", "percent", 0},
		{17, "request", "temperature_1", "degrees Celsius", 0.0625},
		{17, "request", "temperature_2", "degrees Celsius", 0.0625},
		{17, "request", "temperature_3", "degrees Celsius", 0.0625},
		{17, "request", "pressure", "bar", 0.1},
		{17, "request", "flow", "l/h", 10},
		{17, "positive", "setpoint", "percent", 0},
		{18, "request", "co2_1", "ppm", 0},
		{18, "request", "co2_2", "ppm", 0},
	} {
		b512HeatpumpIORow, b512HeatpumpIORaw := b512HeatpumpIOGet(b512HeatpumpIOT, b512HeatpumpIOCase.b512HeatpumpIOOp, b512HeatpumpIOCase.b512HeatpumpIODirection, b512HeatpumpIOCase.b512HeatpumpIOField)
		b512HeatpumpIOEqual(b512HeatpumpIOT, b512HeatpumpIORow.B512HeatpumpIOScale, b512HeatpumpIOCase.b512HeatpumpIOScale)
		b512HeatpumpIOEqual(b512HeatpumpIOT, b512HeatpumpIOObject(b512HeatpumpIOT, b512HeatpumpIORaw)["unit"], b512HeatpumpIOCase.b512HeatpumpIOUnit)
	}
	b512HeatpumpIOPressure, _ := b512HeatpumpIOGet(b512HeatpumpIOT, 17, "request", "pressure")
	b512HeatpumpIOValue, b512HeatpumpIOReplacement, b512HeatpumpIORaw := b512HeatpumpIOEvaluate(b512HeatpumpIOPressure, []byte{255})
	b512HeatpumpIOEqual(b512HeatpumpIOT, b512HeatpumpIOValue, 25.5)
	b512HeatpumpIOEqual(b512HeatpumpIOT, b512HeatpumpIOReplacement, false)
	b512HeatpumpIOEqual(b512HeatpumpIOT, b512HeatpumpIORaw, []byte{255})
}
