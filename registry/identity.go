package registry

import (
	"strconv"
	"strings"
)

const stableIdentityKeyPrefix = "triple:"

type physicalIdentity struct {
	manufacturer    string
	deviceID        string
	serialNumber    string
	macAddress      string
	softwareVersion string
	hardwareVersion string
}

func canonicalPhysicalIdentity(info DeviceInfo) physicalIdentity {
	return physicalIdentity{
		manufacturer:    normalizeIdentityPart(info.Manufacturer),
		deviceID:        normalizeIdentityPart(info.DeviceID),
		serialNumber:    normalizeIdentityPart(info.SerialNumber),
		macAddress:      normalizeIdentityPart(info.MacAddress),
		softwareVersion: normalizeIdentityPart(info.SoftwareVersion),
		hardwareVersion: normalizeIdentityPart(info.HardwareVersion),
	}
}

func (identity physicalIdentity) key() string {
	if !identity.isQualified() {
		return ""
	}
	return stableIdentityKeyPrefix + encodeIdentityMember(identity.manufacturer) +
		encodeIdentityMember(identity.deviceID) + encodeIdentityMember(identity.serialNumber)
}

// encodeIdentityMember preserves a normalized identity member exactly while
// making each member boundary unambiguous. Internal punctuation is contract
// data, not a separator, so delimiter concatenation cannot represent a triple.
func encodeIdentityMember(member string) string {
	return strconv.Itoa(len(member)) + ":" + member
}

// isQualified reports whether identity has the complete normalized triple
// required to join independently observed eBUS addresses. The serial sentinel
// check intentionally accepts only optional 0x prefix and leading-zero
// spelling variations of the rejected hexadecimal values; it does not parse or
// reinterpret ordinary product serial formats.
func (identity physicalIdentity) isQualified() bool {
	return identity.manufacturer != "" && identity.deviceID != "" &&
		identity.serialNumber != "" && !isInvalidSerialSentinel(identity.serialNumber)
}

func isInvalidSerialSentinel(serial string) bool {
	normalized := normalizeIdentityPart(serial)
	digits := strings.TrimPrefix(normalized, "0X")
	if digits == "" {
		return false
	}
	for _, digit := range digits {
		if (digit < '0' || digit > '9') && (digit < 'A' || digit > 'F') {
			return false
		}
	}
	digits = strings.TrimLeft(digits, "0")
	if digits == "" {
		return true
	}
	return digits == "FFFFFFFF" || digits == "7FFFFFFF"
}

// preservesCurrentIdentityAuthority reports whether a partial same-address
// observation can retain an entry's already-published identity. Every
// supplied triple member must agree with the entry's current normalized
// triple. In particular, a model-only change is not a sparse refresh: its
// retained serial may remain in local state, but cannot keep or create a
// cross-address identity authority for the composite.
func preservesCurrentIdentityAuthority(incoming DeviceInfo, entry *deviceEntry) bool {
	if entry == nil || entry.identityKey == "" || entry.physical.key() != entry.identityKey {
		return false
	}
	incomingPhysical := canonicalPhysicalIdentity(incoming)
	return identityPartAgreesOrMissing(incomingPhysical.manufacturer, entry.physical.manufacturer) &&
		identityPartAgreesOrMissing(incomingPhysical.deviceID, entry.physical.deviceID) &&
		identityPartAgreesOrMissing(incomingPhysical.serialNumber, entry.physical.serialNumber)
}

func identityPartAgreesOrMissing(incoming, current string) bool {
	return incoming == "" || incoming == current
}

func normalizeIdentityPart(value string) string {
	return strings.ToUpper(strings.TrimSpace(value))
}
