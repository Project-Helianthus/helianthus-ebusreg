package ebus_standard_catalog

import "testing"

// TestCompareIdentity_RetainsBothOnDisagreement confirms the non-overwrite
// policy from docs/architecture/ebus_standard/07-identity-provenance.md:
// when DeviceInfo and the ebus_standard Identification disagree, both
// values MUST be retained with source labels and Agreement=false.
func TestCompareIdentity_RetainsBothOnDisagreement(t *testing.T) {
	existing := DeviceInfo{
		Manufacturer:    "Vaillant",
		DeviceID:        "BAI00",
		SoftwareVersion: "0407",
		HardwareVersion: "9001",
	}
	desc := IdentificationDescriptor{
		Manufacturer:    "0xb5",
		DeviceID:        "BAI00",
		SoftwareVersion: "0407",
		HardwareVersion: "9001",
		Source:          string(SourceEBUSStandardIdent),
		CatalogVersion:  CatalogVersion,
		Valid:           true,
	}
	rec := CompareIdentity(existing, desc)

	if rec.Manufacturer.Agreement {
		t.Fatalf("Manufacturer.Agreement=true, want false (Vaillant != 0xb5)")
	}
	if rec.Manufacturer.Preferred != "Vaillant" {
		t.Fatalf("Manufacturer.Preferred=%q, want %q (deterministic precedence: device_info wins)",
			rec.Manufacturer.Preferred, "Vaillant")
	}
	if len(rec.Manufacturer.Sources) != 2 {
		t.Fatalf("Manufacturer.Sources len=%d, want 2 (both sources retained)", len(rec.Manufacturer.Sources))
	}
	// Agreement rows for the identical fields must still carry both sources
	// so consumers can see the evidence trail.
	if !rec.DeviceID.Agreement {
		t.Fatalf("DeviceID.Agreement=false, want true")
	}
}

// TestCompareIdentity_DoesNotMutateExisting confirms the non-overwrite rule:
// the existing DeviceInfo MUST remain byte-identical after the comparison.
func TestCompareIdentity_DoesNotMutateExisting(t *testing.T) {
	original := DeviceInfo{
		Manufacturer: "Vaillant", DeviceID: "BAI00",
		SoftwareVersion: "0407", HardwareVersion: "9001",
	}
	copy := original
	desc := IdentificationDescriptor{
		Manufacturer: "0xb5", DeviceID: "BAI00",
		SoftwareVersion: "0407", HardwareVersion: "9001",
		Source:         string(SourceEBUSStandardIdent),
		CatalogVersion: CatalogVersion, Valid: true,
	}
	_ = CompareIdentity(copy, desc)
	if copy != original {
		t.Fatalf("CompareIdentity mutated existing DeviceInfo: before=%+v after=%+v", original, copy)
	}
}
