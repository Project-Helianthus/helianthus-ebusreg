package system

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Project-Helianthus/helianthus-ebusreg/registry"
)

// Outcomes of the boiler-target declaration.
const (
	BoilerOutcomePublish             = "publish"
	BoilerOutcomeHeatPumpCollision   = "reject_heat_pump_collision"
	BoilerOutcomeUnknown             = "unknown"
	BoilerOutcomeRejectWithoutChange = "reject_without_mutation"
)

// Reasons named by the boiler-target declaration.
const (
	BoilerReasonOtherManufacturer = "other_manufacturer"
	BoilerReasonNotBoilerClass    = "not_boiler_class"
)

// BoilerFlowRowOpcode and BoilerFlowRowSelector name the declared boiler flow
// row: B509 opcode 0x0D with selector 18 00.
const (
	BoilerFlowRowOpcode   byte   = 0x0D
	BoilerFlowRowSelector uint16 = 0x1800
)

type boilerTargetRecord struct {
	BoilerTarget struct {
		Manufacturer     string `json:"manufacturer"`
		DeviceIDPrefix   string `json:"device_id_prefix"`
		HeatPumpDeviceID string `json:"heat_pump_device_id"`
		HeatPumpPrefix   string `json:"heat_pump_prefix"`
		Level            string `json:"level"`
		EvidenceGrade    string `json:"evidence_grade"`
		FlowRow          struct {
			Opcode                   string `json:"opcode"`
			Selector                 string `json:"selector"`
			SemregElement            string `json:"semreg_element"`
			EvidenceGrade            string `json:"evidence_grade"`
			QualifiedBindingRequired bool   `json:"qualified_binding_required"`
		} `json:"flow_row"`
	} `json:"boiler_target"`
}

var boilerRule boilerTargetRecord

func init() {
	if err := json.Unmarshal(controllerQualificationRuleJSON, &boilerRule); err != nil {
		panic(fmt.Sprintf("boiler target rule: %v", err))
	}
}

// BoilerQualification is the driver's answer for one Identification payload
// of a bus participant that may be a boiler target.
type BoilerQualification struct {
	Outcome    string
	Reason     string
	DeviceID   string
	Level      string
	RuleDigest string
}

// QualifyBoilerIdentification decides whether a ten-byte Identification
// payload from responseAddress, in reply to requestAddress, is a declared
// boiler target. A participant is a boiler target when its manufacturer is
// 0xB5 and its device ID has the boiler class prefix at the address it answers
// on. A heat-pump identity is a collision and is not a boiler target.
func QualifyBoilerIdentification(requestAddress, responseAddress byte, payload []byte) BoilerQualification {
	result := BoilerQualification{RuleDigest: controllerRuleDigest}
	if len(payload) != 10 {
		result.Outcome, result.Reason = BoilerOutcomeRejectWithoutChange, ControllerReasonIncomplete
		return result
	}
	if requestAddress != responseAddress {
		result.Outcome, result.Reason = BoilerOutcomeRejectWithoutChange, ControllerReasonAddressMismatch
		return result
	}
	return classifyBoiler(payload[0], string(payload[1:6]), result)
}

func classifyBoiler(manufacturer byte, deviceID string, result BoilerQualification) BoilerQualification {
	result.DeviceID = deviceID
	if manufacturer != controllerQualificationManufacturer {
		result.Outcome, result.Reason = BoilerOutcomeUnknown, BoilerReasonOtherManufacturer
		return result
	}
	if deviceID == boilerRule.BoilerTarget.HeatPumpDeviceID || strings.HasPrefix(deviceID, boilerRule.BoilerTarget.HeatPumpPrefix) {
		result.Outcome = BoilerOutcomeHeatPumpCollision
		return result
	}
	if !strings.HasPrefix(deviceID, boilerRule.BoilerTarget.DeviceIDPrefix) {
		result.Outcome, result.Reason = BoilerOutcomeUnknown, BoilerReasonNotBoilerClass
		return result
	}
	result.Outcome = BoilerOutcomePublish
	result.Level = boilerRule.BoilerTarget.Level
	return result
}

// DeclareBoilerTarget is the boiler-target declaration rule for the registry.
// It needs no serial and no software or product-number field: the class is
// decided by manufacturer and device ID, which the registry must hold in the
// exact five-character and 0xB5 forms.
func DeclareBoilerTarget(info registry.DeviceInfo) registry.TargetDeclaration {
	manufacturer, ok := manufacturerByteOf(info.Manufacturer)
	if !ok || len(info.DeviceID) != 5 {
		return registry.TargetDeclaration{Reason: ControllerReasonIncomplete, RuleDigest: controllerRuleDigest}
	}
	q := classifyBoiler(manufacturer, info.DeviceID, BoilerQualification{RuleDigest: controllerRuleDigest})
	declaration := registry.TargetDeclaration{
		Reason:     q.Reason,
		RuleDigest: q.RuleDigest,
	}
	if q.Outcome != BoilerOutcomePublish {
		return declaration
	}
	declaration.Declared = true
	declaration.Class = "boiler_target"
	declaration.Level = q.Level
	declaration.FlowOpcode = BoilerFlowRowOpcode
	declaration.FlowSelector = BoilerFlowRowSelector
	declaration.FlowElement = boilerRule.BoilerTarget.FlowRow.SemregElement
	declaration.FlowEvidenceGrade = boilerRule.BoilerTarget.FlowRow.EvidenceGrade
	declaration.FlowQualifiedBinding = boilerRule.BoilerTarget.FlowRow.QualifiedBindingRequired
	return declaration
}

// manufacturerByteOf maps a registry manufacturer string to its byte: "Vaillant"
// is 0xB5 and "0xNN" is NN.
func manufacturerByteOf(manufacturer string) (byte, bool) {
	if manufacturer == "Vaillant" {
		return controllerQualificationManufacturer, true
	}
	if strings.HasPrefix(manufacturer, "0x") && len(manufacturer) == 4 {
		return parseHexByte(manufacturer[2:])
	}
	return 0, false
}
