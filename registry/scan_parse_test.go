package registry

import "testing"

func TestParseDeviceInfo_VaillantBAI00(t *testing.T) {
	t.Parallel()

	info, err := parseDeviceInfo(0x08, []byte{
		0xB5, 'B', 'A', 'I', '0', '0', // MF + ID(5)
		0x07, 0x04, // SW(PIN=2)
		0x76, 0x03, // HW(PIN=2)
	})
	if err != nil {
		t.Fatalf("parseDeviceInfo error = %v", err)
	}

	if info.Manufacturer != "Vaillant" {
		t.Fatalf("Manufacturer = %q, want %q", info.Manufacturer, "Vaillant")
	}
	if info.DeviceID != "BAI00" {
		t.Fatalf("DeviceID = %q, want %q", info.DeviceID, "BAI00")
	}
	if info.SoftwareVersion != "0704" {
		t.Fatalf("SoftwareVersion = %q, want %q", info.SoftwareVersion, "0704")
	}
	if info.HardwareVersion != "7603" {
		t.Fatalf("HardwareVersion = %q, want %q", info.HardwareVersion, "7603")
	}
}
