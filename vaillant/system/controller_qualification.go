package system

import (
	"crypto/sha256"
	_ "embed"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/Project-Helianthus/helianthus-ebusreg/registry"
)

// controllerQualificationRuleJSON is the docs-ebus controller qualification
// record, a byte-identical copy of its merged fixture.
//
//go:embed controller_qualification_rule.json
var controllerQualificationRuleJSON []byte

const (
	controllerQualificationManufacturer = 0xB5

	controllerOutcomeQualified = "qualified"
	controllerOutcomeUnknown   = "unknown"
	controllerOutcomeRejected  = "reject_without_mutation"
)

// Reasons named by the controller declaration. Malformed or incomplete input
// is rejected without mutation; a well-formed input the record does not
// declare is unknown.
const (
	ControllerReasonIncomplete        = "incomplete"
	ControllerReasonAddressMismatch   = "address_mismatch"
	ControllerReasonOtherManufacturer = "other_manufacturer"
	ControllerReasonInvalidBCD        = "invalid_bcd"
	ControllerReasonFFReplacement     = "ff_replacement"
	ControllerReasonUnlistedPair      = "unlisted_pair"
)

type controllerRuleDocument struct {
	Identities []struct {
		DeviceID      string `json:"device_id"`
		ProductNumber string `json:"product_number"`
		Level         string `json:"level"`
		EvidenceGrade string `json:"evidence_grade"`
	} `json:"identities"`
	DeclaredRows []struct {
		Group         string   `json:"group"`
		Register      string   `json:"register"`
		SemregElement *string  `json:"semreg_element"`
		EvidenceGrade string   `json:"evidence_grade"`
		Level         string   `json:"level"`
		DeviceIDs     []string `json:"device_ids"`
	} `json:"declared_rows"`
}

var (
	controllerRule       controllerRuleDocument
	controllerRuleDigest string
)

func init() {
	if err := json.Unmarshal(controllerQualificationRuleJSON, &controllerRule); err != nil {
		panic(fmt.Sprintf("controller qualification rule: %v", err))
	}
	sum := sha256.Sum256(controllerQualificationRuleJSON)
	controllerRuleDigest = hex.EncodeToString(sum[:])
}

// ControllerQualificationRuleDigest returns the SHA-256 of the embedded
// controller qualification record, in lower-case hexadecimal.
func ControllerQualificationRuleDigest() string {
	return controllerRuleDigest
}

// ControllerQualification is the driver's answer for one identification
// payload. Outcome is qualified, unknown or reject_without_mutation; Reason
// names the cause when it is not qualified. Level and EvidenceGrade are the
// recorded level and grade of the (device ID, product number) pair and decide
// nothing.
type ControllerQualification struct {
	Outcome         string
	Reason          string
	DeviceID        string
	ProductNumber   string
	SoftwareVersion string
	Level           string
	EvidenceGrade   string
	Rows            []registry.DeclaredRow
	RuleDigest      string
}

// QualifyControllerIdentification decides the controller declaration of a
// ten-byte Identification (0x07 0x04) payload received from responseAddress
// in reply to a request to requestAddress. Byte 0 is the manufacturer, bytes
// 1-5 the device ID, bytes 6-7 the software version and bytes 8-9 the product
// number. The product number is BCD, least-significant byte first: the value
// is BCD(byte 8) + 100 x BCD(byte 9), printed as four hexadecimal digits.
func QualifyControllerIdentification(requestAddress, responseAddress byte, payload []byte) ControllerQualification {
	result := ControllerQualification{RuleDigest: controllerRuleDigest}
	if len(payload) != 10 {
		result.Outcome, result.Reason = controllerOutcomeRejected, ControllerReasonIncomplete
		return result
	}
	if requestAddress != responseAddress {
		result.Outcome, result.Reason = controllerOutcomeRejected, ControllerReasonAddressMismatch
		return result
	}
	deviceID := string(payload[1:6])
	result.DeviceID = deviceID
	result.SoftwareVersion = fmt.Sprintf("%02X%02X", payload[6], payload[7])
	if payload[0] != controllerQualificationManufacturer {
		result.Outcome, result.Reason = controllerOutcomeUnknown, ControllerReasonOtherManufacturer
		return result
	}
	if payload[8] == 0xFF || payload[9] == 0xFF {
		result.Outcome, result.Reason = controllerOutcomeUnknown, ControllerReasonFFReplacement
		return result
	}
	spn, ok := decodeProductNumber(payload[8], payload[9])
	if !ok {
		result.Outcome, result.Reason = controllerOutcomeRejected, ControllerReasonInvalidBCD
		return result
	}
	result.ProductNumber = fmt.Sprintf("%04X", spn)

	identity, listed := controllerIdentity(deviceID, result.ProductNumber)
	if !listed {
		result.Outcome, result.Reason = controllerOutcomeUnknown, ControllerReasonUnlistedPair
		return result
	}
	result.Outcome = controllerOutcomeQualified
	result.Level = identity.Level
	result.EvidenceGrade = identity.EvidenceGrade
	result.Rows = controllerDeclaredRows(deviceID)
	return result
}

// decodeProductNumber returns BCD(low) + 100 x BCD(high) and false when a
// nibble is not a decimal digit.
func decodeProductNumber(low, high byte) (int, bool) {
	lowDigits, lowValid := decodeBCDByte(low)
	highDigits, highValid := decodeBCDByte(high)
	if !lowValid || !highValid {
		return 0, false
	}
	return lowDigits + 100*highDigits, true
}

func decodeBCDByte(b byte) (int, bool) {
	if b>>4 > 9 || b&0x0F > 9 {
		return 0, false
	}
	return int(b>>4)*10 + int(b&0x0F), true
}

type controllerIdentityEntry struct {
	Level         string
	EvidenceGrade string
}

func controllerIdentity(deviceID, productNumber string) (controllerIdentityEntry, bool) {
	for _, identity := range controllerRule.Identities {
		if identity.DeviceID == deviceID && identity.ProductNumber == productNumber {
			return controllerIdentityEntry{Level: identity.Level, EvidenceGrade: identity.EvidenceGrade}, true
		}
	}
	return controllerIdentityEntry{}, false
}

// controllerDeclaredRows returns the rows the record lists for deviceID, in
// record order. Each row keeps its grade and level.
func controllerDeclaredRows(deviceID string) []registry.DeclaredRow {
	var rows []registry.DeclaredRow
	for _, row := range controllerRule.DeclaredRows {
		if !containsString(row.DeviceIDs, deviceID) {
			continue
		}
		group, groupOK := parseHexByte(row.Group)
		register, registerOK := parseHexUint16(row.Register)
		if !groupOK || !registerOK {
			continue
		}
		element := ""
		if row.SemregElement != nil {
			element = *row.SemregElement
		}
		rows = append(rows, registry.DeclaredRow{
			Group:         group,
			Register:      register,
			Element:       element,
			EvidenceGrade: row.EvidenceGrade,
			Level:         row.Level,
		})
	}
	return rows
}

// DeclareController is the controller declaration rule for the registry. It
// rebuilds the Identification payload from the registry's identity fields and
// qualifies it. It declares a controller only when the record lists its exact
// pair; an identity the registry cannot express as a ten-byte payload is
// incomplete.
func DeclareController(info registry.DeviceInfo) registry.ControllerDeclaration {
	payload, ok := identificationPayloadOf(info)
	if !ok {
		return registry.ControllerDeclaration{Reason: ControllerReasonIncomplete, RuleDigest: controllerRuleDigest}
	}
	q := QualifyControllerIdentification(info.Address, info.Address, payload)
	declaration := registry.ControllerDeclaration{
		Declared:      q.Outcome == controllerOutcomeQualified,
		Reason:        q.Reason,
		ProductNumber: q.ProductNumber,
		RuleDigest:    q.RuleDigest,
	}
	if declaration.Declared {
		declaration.Family = q.DeviceID
		declaration.Level = q.Level
		declaration.EvidenceGrade = q.EvidenceGrade
		declaration.Rows = q.Rows
	}
	return declaration
}

// identificationPayloadOf rebuilds the ten-byte Identification payload that a
// registry identity was parsed from. The manufacturer is "Vaillant" for 0xB5
// or "0xNN" otherwise; the device ID is five characters; the software and
// product-number fields are four hexadecimal digits each.
func identificationPayloadOf(info registry.DeviceInfo) ([]byte, bool) {
	manufacturer, ok := manufacturerByteOf(info.Manufacturer)
	if !ok {
		return nil, false
	}
	if len(info.DeviceID) != 5 {
		return nil, false
	}
	software, ok := decodeHexPair(info.SoftwareVersion)
	if !ok {
		return nil, false
	}
	product, ok := decodeHexPair(info.HardwareVersion)
	if !ok {
		return nil, false
	}
	payload := make([]byte, 0, 10)
	payload = append(payload, manufacturer)
	payload = append(payload, info.DeviceID...)
	payload = append(payload, software...)
	payload = append(payload, product...)
	return payload, true
}

func decodeHexPair(s string) ([]byte, bool) {
	if len(s) != 4 {
		return nil, false
	}
	b, err := hex.DecodeString(s)
	if err != nil || len(b) != 2 {
		return nil, false
	}
	return b, true
}

func parseHexByte(s string) (byte, bool) {
	value, err := strconv.ParseUint(strings.TrimPrefix(s, "0x"), 16, 8)
	return byte(value), err == nil
}

func parseHexUint16(s string) (uint16, bool) {
	value, err := strconv.ParseUint(strings.TrimPrefix(s, "0x"), 16, 16)
	return uint16(value), err == nil
}

func containsString(values []string, want string) bool {
	for _, value := range values {
		if value == want {
			return true
		}
	}
	return false
}
