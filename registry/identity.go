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
	if identity.manufacturer == "" || identity.deviceID == "" {
		return identity.withFallbackModelSignature()
	}
	if identity.serialNumber != "" {
		return strings.Join([]string{
			"sn",
			identity.manufacturer,
			identity.deviceID,
			identity.serialNumber,
		}, "|")
	}
	if identity.macAddress != "" {
		return strings.Join([]string{
			"mac",
			identity.manufacturer,
			identity.deviceID,
			identity.macAddress,
		}, "|")
	}
	return identity.withFallbackModelSignature()
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
