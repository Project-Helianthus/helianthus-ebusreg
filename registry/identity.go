package registry

import "strings"

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
	return strings.Join([]string{
		"triple",
		identity.manufacturer,
		identity.deviceID,
		identity.serialNumber,
	}, "|")
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

func (identity physicalIdentity) withFallbackModelSignature() string {
	if identity.manufacturer == "" || identity.deviceID == "" || identity.softwareVersion == "" || identity.hardwareVersion == "" {
		return ""
	}
	return strings.Join([]string{
		"sig",
		identity.manufacturer,
		identity.deviceID,
		identity.softwareVersion,
		identity.hardwareVersion,
	}, "|")
}

func normalizeIdentityPart(value string) string {
	return strings.ToUpper(strings.TrimSpace(value))
}
