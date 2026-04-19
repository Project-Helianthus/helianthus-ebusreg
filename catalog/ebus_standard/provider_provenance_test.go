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

// TestCompareIdentity_InvalidDescriptorDoesNotPromoteEmptyExisting confirms
// that when existingVal is empty and the descriptor is marked invalid
// (desc.Valid == false), the malformed descriptor value MUST NOT be
// promoted to Preferred. This guards against surfacing failed-decode
// identity as the canonical preferred value for new devices, which would
// mislabel them until manually corrected (see docs/architecture/
// ebus_standard/07-identity-provenance.md §"Deterministic Precedence":
// an invalid descriptor is not a confident identity source).
func TestCompareIdentity_InvalidDescriptorDoesNotPromoteEmptyExisting(t *testing.T) {
	existing := DeviceInfo{} // no prior device_info evidence
	desc := IdentificationDescriptor{
		Manufacturer:    "??",
		DeviceID:        "??",
		SoftwareVersion: "??",
		HardwareVersion: "??",
		Source:          string(SourceEBUSStandardIdent),
		CatalogVersion:  CatalogVersion,
		Valid:           false, // malformed / failed decode
	}
	rec := CompareIdentity(existing, desc)

	for name, fp := range map[string]FieldProvenance{
		"Manufacturer":    rec.Manufacturer,
		"DeviceID":        rec.DeviceID,
		"SoftwareVersion": rec.SoftwareVersion,
		"HardwareVersion": rec.HardwareVersion,
	} {
		if fp.Preferred != "" {
			t.Fatalf("%s.Preferred=%q, want empty (invalid descriptor must not be promoted when existing is empty)",
				name, fp.Preferred)
		}
		// Evidence trail must still carry the invalid descriptor value
		// with Valid=false so downstream consumers can see it.
		if len(fp.Sources) != 1 {
			t.Fatalf("%s.Sources len=%d, want 1 (invalid descriptor retained in evidence trail)",
				name, len(fp.Sources))
		}
		if fp.Sources[0].Valid {
			t.Fatalf("%s.Sources[0].Valid=true, want false (descriptor was marked invalid)", name)
		}
	}
}

// TestCompareIdentity_ValidDescriptorPromotesEmptyExisting is the happy-path
// sanity check for the empty-existing path: a VALID descriptor MUST be
// promoted to Preferred when there is no prior device_info evidence.
func TestCompareIdentity_ValidDescriptorPromotesEmptyExisting(t *testing.T) {
	existing := DeviceInfo{} // no prior device_info evidence
	desc := IdentificationDescriptor{
		Manufacturer:    "Vaillant",
		DeviceID:        "BAI00",
		SoftwareVersion: "0407",
		HardwareVersion: "9001",
		Source:          string(SourceEBUSStandardIdent),
		CatalogVersion:  CatalogVersion,
		Valid:           true,
	}
	rec := CompareIdentity(existing, desc)

	if rec.Manufacturer.Preferred != "Vaillant" {
		t.Fatalf("Manufacturer.Preferred=%q, want %q", rec.Manufacturer.Preferred, "Vaillant")
	}
	if rec.DeviceID.Preferred != "BAI00" {
		t.Fatalf("DeviceID.Preferred=%q, want %q", rec.DeviceID.Preferred, "BAI00")
	}
	if rec.SoftwareVersion.Preferred != "0407" {
		t.Fatalf("SoftwareVersion.Preferred=%q, want %q", rec.SoftwareVersion.Preferred, "0407")
	}
	if rec.HardwareVersion.Preferred != "9001" {
		t.Fatalf("HardwareVersion.Preferred=%q, want %q", rec.HardwareVersion.Preferred, "9001")
	}
}

// TestCompareIdentity_InvalidDescriptorPreservesExistingPrecedence confirms
// that the existing precedence rule (device_info wins over descriptor) is
// preserved when the descriptor is invalid AND existing is non-empty.
// The invalid descriptor's value must be retained as evidence with
// Valid=false, but Preferred must be the existing value.
func TestCompareIdentity_InvalidDescriptorPreservesExistingPrecedence(t *testing.T) {
	existing := DeviceInfo{
		Manufacturer:    "Vaillant",
		DeviceID:        "BAI00",
		SoftwareVersion: "0407",
		HardwareVersion: "9001",
	}
	desc := IdentificationDescriptor{
		Manufacturer:    "??",
		DeviceID:        "??",
		SoftwareVersion: "??",
		HardwareVersion: "??",
		Source:          string(SourceEBUSStandardIdent),
		CatalogVersion:  CatalogVersion,
		Valid:           false,
	}
	rec := CompareIdentity(existing, desc)

	if rec.Manufacturer.Preferred != "Vaillant" {
		t.Fatalf("Manufacturer.Preferred=%q, want %q (device_info precedence preserved)",
			rec.Manufacturer.Preferred, "Vaillant")
	}
	if len(rec.Manufacturer.Sources) != 2 {
		t.Fatalf("Manufacturer.Sources len=%d, want 2 (both sources retained as evidence)",
			len(rec.Manufacturer.Sources))
	}
	// Locate the descriptor source and confirm Valid=false propagated.
	foundInvalid := false
	for _, s := range rec.Manufacturer.Sources {
		if s.Source == SourceEBUSStandardIdent {
			if s.Valid {
				t.Fatalf("descriptor source Valid=true, want false (descriptor marked invalid)")
			}
			foundInvalid = true
		}
	}
	if !foundInvalid {
		t.Fatalf("descriptor source not present in Sources; got %+v", rec.Manufacturer.Sources)
	}
}
