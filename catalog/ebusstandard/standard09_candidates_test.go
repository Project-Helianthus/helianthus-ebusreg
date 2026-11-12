package ebusstandard

import (
	standard09CandidateBinary "encoding/binary"
	standard09CandidateReflect "reflect"
	standard09CandidateTesting "testing"
)

func standard09CandidateByte(standard09CandidateValue uint8) *uint8 { return &standard09CandidateValue }

func standard09CandidateKey(standard09CandidateVariant string, standard09CandidateSecondary uint8, standard09CandidateDirection, standard09CandidateField string) Standard09CandidateKey {
	return Standard09CandidateKey{Standard09CandidateFamily: "memory-service", Standard09CandidateVariant: standard09CandidateVariant, Standard09CandidatePrimary: standard09CandidateByte(9), Standard09CandidateSecondary: standard09CandidateByte(standard09CandidateSecondary), Standard09CandidateDirection: standard09CandidateDirection, Standard09CandidateField: standard09CandidateField}
}

func standard09CandidateGet(standard09CandidateT *standard09CandidateTesting.T, standard09CandidateKey Standard09CandidateKey) Standard09CandidateEntry {
	standard09CandidateT.Helper()
	standard09CandidateEntry, standard09CandidateOK := Standard09CandidateLookup(standard09CandidateKey)
	if !standard09CandidateOK {
		standard09CandidateT.Fatalf("missing key: %+v", standard09CandidateKey)
	}
	return standard09CandidateEntry
}

func standard09CandidateInt(standard09CandidateT *standard09CandidateTesting.T, standard09CandidateValue *int, standard09CandidateWant int) {
	standard09CandidateT.Helper()
	if standard09CandidateValue == nil || *standard09CandidateValue != standard09CandidateWant {
		standard09CandidateT.Fatalf("integer %v, want %d", standard09CandidateValue, standard09CandidateWant)
	}
}

func TestAb79ExactKeys(standard09CandidateT *standard09CandidateTesting.T) {
	standard09CandidateControl := standard09CandidateKey("0902", 2, "abstract-positive", "data")
	standard09CandidateGet(standard09CandidateT, standard09CandidateControl)
	standard09CandidateCases := []struct {
		Standard09CandidateName   string
		Standard09CandidateMutate func(*Standard09CandidateKey)
	}{
		{"family", func(standard09CandidateK *Standard09CandidateKey) {
			standard09CandidateK.Standard09CandidateFamily = "unknown"
		}},
		{"variant", func(standard09CandidateK *Standard09CandidateKey) {
			standard09CandidateK.Standard09CandidateVariant = "0901"
		}},
		{"primary", func(standard09CandidateK *Standard09CandidateKey) {
			standard09CandidateK.Standard09CandidatePrimary = standard09CandidateByte(8)
		}},
		{"primary-absent", func(standard09CandidateK *Standard09CandidateKey) {
			standard09CandidateK.Standard09CandidatePrimary = nil
		}},
		{"primary-zero", func(standard09CandidateK *Standard09CandidateKey) {
			standard09CandidateK.Standard09CandidatePrimary = standard09CandidateByte(0)
		}},
		{"secondary", func(standard09CandidateK *Standard09CandidateKey) {
			standard09CandidateK.Standard09CandidateSecondary = standard09CandidateByte(3)
		}},
		{"secondary-absent", func(standard09CandidateK *Standard09CandidateKey) {
			standard09CandidateK.Standard09CandidateSecondary = nil
		}},
		{"secondary-zero", func(standard09CandidateK *Standard09CandidateKey) {
			standard09CandidateK.Standard09CandidateSecondary = standard09CandidateByte(0)
		}},
		{"operation-zero", func(standard09CandidateK *Standard09CandidateKey) {
			standard09CandidateK.Standard09CandidateOperation = standard09CandidateByte(0)
		}},
		{"operation", func(standard09CandidateK *Standard09CandidateKey) {
			standard09CandidateK.Standard09CandidateOperation = standard09CandidateByte(1)
		}},
		{"direction", func(standard09CandidateK *Standard09CandidateKey) {
			standard09CandidateK.Standard09CandidateDirection = "abstract-request"
		}},
		{"field", func(standard09CandidateK *Standard09CandidateKey) {
			standard09CandidateK.Standard09CandidateField = "address"
		}},
	}
	for _, standard09CandidateCase := range standard09CandidateCases {
		standard09CandidateT.Run("ab79-"+standard09CandidateCase.Standard09CandidateName, func(standard09CandidateT *standard09CandidateTesting.T) {
			standard09CandidateK := standard09CandidateControl
			standard09CandidateCase.Standard09CandidateMutate(&standard09CandidateK)
			if _, standard09CandidateOK := Standard09CandidateLookup(standard09CandidateK); standard09CandidateOK {
				standard09CandidateT.Fatal("changed axis matched")
			}
		})
	}
	standard09CandidateZero := standard09CandidateKey("0900", 0, "abstract-positive", "data")
	standard09CandidateGet(standard09CandidateT, standard09CandidateZero)
	standard09CandidateZero.Standard09CandidateSecondary = nil
	if _, standard09CandidateOK := Standard09CandidateLookup(standard09CandidateZero); standard09CandidateOK {
		standard09CandidateT.Fatal("absent secondary matched explicit zero")
	}
}

func TestAb79Inventory(standard09CandidateT *standard09CandidateTesting.T) {
	standard09CandidateCatalog := Standard09Candidates()
	if len(standard09CandidateCatalog.Standard09CandidateEntries) != 25 || standard09CandidateCatalog.Standard09CandidatePresenceGrade != "Proven" || standard09CandidateCatalog.Standard09CandidateUnknownGrade != "Unknown" || len(standard09CandidateCatalog.Standard09CandidateUnknowns) != 8 {
		standard09CandidateT.Fatal("inventory or grade changed")
	}
	standard09CandidateSeen := make(map[string]bool)
	for _, standard09CandidateEntry := range standard09CandidateCatalog.Standard09CandidateEntries {
		if standard09CandidateEntry.Standard09CandidateGrade != "Hypothesis" || standard09CandidateEntry.Standard09CandidatePhysicalOffset != nil || standard09CandidateEntry.Standard09CandidateUnit != nil || standard09CandidateEntry.Standard09CandidateEnum != nil || standard09CandidateEntry.Standard09CandidateKey.Standard09CandidateOperation != nil {
			standard09CandidateT.Fatal("invented qualification or metadata")
		}
		standard09CandidateFound := standard09CandidateGet(standard09CandidateT, standard09CandidateEntry.Standard09CandidateKey)
		if !standard09CandidateReflect.DeepEqual(standard09CandidateFound, standard09CandidateEntry) {
			standard09CandidateT.Fatal("lookup differs from accessor")
		}
		standard09CandidateID := standard09CandidateEntry.Standard09CandidateKey.Standard09CandidateVariant + standard09CandidateEntry.Standard09CandidateKey.Standard09CandidateDirection + standard09CandidateEntry.Standard09CandidateKey.Standard09CandidateField
		if standard09CandidateSeen[standard09CandidateID] {
			standard09CandidateT.Fatal("duplicate identity")
		}
		standard09CandidateSeen[standard09CandidateID] = true
	}
	for standard09CandidateV, standard09CandidateVariant := range []string{"0900", "0901", "0902", "0903"} {
		for standard09CandidateOffset, standard09CandidateField := range []string{"PB", "SB"} {
			standard09CandidateEntry := standard09CandidateGet(standard09CandidateT, standard09CandidateKey(standard09CandidateVariant, uint8(standard09CandidateV), "abstract-request", standard09CandidateField))
			standard09CandidateInt(standard09CandidateT, standard09CandidateEntry.Standard09CandidateAbstractOffset, standard09CandidateOffset)
			standard09CandidateInt(standard09CandidateT, standard09CandidateEntry.Standard09CandidateWidth, 1)
			if standard09CandidateV != 2 {
				standard09CandidateEntry = standard09CandidateGet(standard09CandidateT, standard09CandidateKey(standard09CandidateVariant, uint8(standard09CandidateV), "abstract-positive", standard09CandidateField))
				standard09CandidateInt(standard09CandidateT, standard09CandidateEntry.Standard09CandidateAbstractOffset, standard09CandidateOffset)
				standard09CandidateInt(standard09CandidateT, standard09CandidateEntry.Standard09CandidateWidth, 1)
			} else if _, standard09CandidateOK := Standard09CandidateLookup(standard09CandidateKey(standard09CandidateVariant, 2, "abstract-positive", standard09CandidateField)); standard09CandidateOK {
				standard09CandidateT.Fatal("0902 invented prefix")
			}
		}
	}
}

func TestAb79AddressAlternatives(standard09CandidateT *standard09CandidateTesting.T) {
	standard09CandidateRaw := []byte{0x34, 0x12}
	for standard09CandidateV, standard09CandidateVariant := range []string{"0900", "0901", "0902", "0903"} {
		if standard09CandidateV < 2 {
			standard09CandidateEntry := standard09CandidateGet(standard09CandidateT, standard09CandidateKey(standard09CandidateVariant, uint8(standard09CandidateV), "abstract-request", "address"))
			if standard09CandidateEntry.Standard09CandidateEncoding == nil || *standard09CandidateEntry.Standard09CandidateEncoding != "unsigned-le16" || standard09CandidateBinary.LittleEndian.Uint16(standard09CandidateRaw) != 4660 {
				standard09CandidateT.Fatal("LE16 address changed")
			}
			standard09CandidateInt(standard09CandidateT, standard09CandidateEntry.Standard09CandidateWidth, 2)
			standard09CandidateInt(standard09CandidateT, standard09CandidateEntry.Standard09CandidateMaximum, 65535)
			standard09CandidateInt(standard09CandidateT, standard09CandidateEntry.Standard09CandidateDefault, 0)
			standard09CandidateInt(standard09CandidateT, standard09CandidateEntry.Standard09CandidateMinimum, 0)
			if standard09CandidateEntry.Standard09CandidateAbstractOffset != nil {
				standard09CandidateT.Fatal("invented address offset")
			}
		} else {
			for standard09CandidateIndex, standard09CandidateField := range []string{"address_low", "address_high"} {
				standard09CandidateEntry := standard09CandidateGet(standard09CandidateT, standard09CandidateKey(standard09CandidateVariant, uint8(standard09CandidateV), "abstract-request", standard09CandidateField))
				if standard09CandidateEntry.Standard09CandidateEncoding == nil || *standard09CandidateEntry.Standard09CandidateEncoding != "unsigned-u8" || standard09CandidateRaw[standard09CandidateIndex] != []byte{52, 18}[standard09CandidateIndex] {
					standard09CandidateT.Fatal("split address changed")
				}
				standard09CandidateInt(standard09CandidateT, standard09CandidateEntry.Standard09CandidateWidth, 1)
				standard09CandidateInt(standard09CandidateT, standard09CandidateEntry.Standard09CandidateMaximum, 255)
				standard09CandidateInt(standard09CandidateT, standard09CandidateEntry.Standard09CandidateDefault, 0)
				standard09CandidateInt(standard09CandidateT, standard09CandidateEntry.Standard09CandidateMinimum, 0)
				if standard09CandidateEntry.Standard09CandidateAbstractOffset != nil {
					standard09CandidateT.Fatal("invented split address offset")
				}
			}
			if _, standard09CandidateOK := Standard09CandidateLookup(standard09CandidateKey(standard09CandidateVariant, uint8(standard09CandidateV), "abstract-request", "address")); standard09CandidateOK {
				standard09CandidateT.Fatal("merged split address")
			}
		}
	}
}

func TestAb79PositiveAlternatives(standard09CandidateT *standard09CandidateTesting.T) {
	for _, standard09CandidateCase := range []struct {
		Standard09CandidateVariant   string
		Standard09CandidateSecondary uint8
		Standard09CandidateOffset    int
		Standard09CandidateRaw       []byte
	}{
		{"0900", 0, 2, []byte{9, 0, 0xaa, 0x55}}, {"0902", 2, 0, []byte{0xaa, 0x55}}, {"0903", 3, 2, []byte{9, 3, 0xaa, 0x55}},
	} {
		standard09CandidateEntry := standard09CandidateGet(standard09CandidateT, standard09CandidateKey(standard09CandidateCase.Standard09CandidateVariant, standard09CandidateCase.Standard09CandidateSecondary, "abstract-positive", "data"))
		standard09CandidateInt(standard09CandidateT, standard09CandidateEntry.Standard09CandidateAbstractOffset, standard09CandidateCase.Standard09CandidateOffset)
		standard09CandidateInt(standard09CandidateT, standard09CandidateEntry.Standard09CandidateMinimum, 0)
		if standard09CandidateEntry.Standard09CandidateWidth != nil || standard09CandidateEntry.Standard09CandidateMaximum != nil || standard09CandidateEntry.Standard09CandidateEncoding == nil || *standard09CandidateEntry.Standard09CandidateEncoding != "opaque-to-end" {
			standard09CandidateT.Fatal("opaque metadata changed")
		}
		if !standard09CandidateReflect.DeepEqual(standard09CandidateCase.Standard09CandidateRaw[*standard09CandidateEntry.Standard09CandidateAbstractOffset:], []byte{0xaa, 0x55}) {
			standard09CandidateT.Fatal("incorrect stripping")
		}
		standard09CandidateEmpty := standard09CandidateCase.Standard09CandidateRaw[:standard09CandidateCase.Standard09CandidateOffset]
		if len(standard09CandidateEmpty[*standard09CandidateEntry.Standard09CandidateAbstractOffset:]) != 0 {
			standard09CandidateT.Fatal("empty opaque data lost")
		}
	}
	for _, standard09CandidateField := range []string{"data", "success"} {
		if _, standard09CandidateOK := Standard09CandidateLookup(standard09CandidateKey("0901", 1, "abstract-positive", standard09CandidateField)); standard09CandidateOK {
			standard09CandidateT.Fatal("prefix-only invented body or success")
		}
	}
}

func TestAb79CountConflict(standard09CandidateT *standard09CandidateTesting.T) {
	for _, standard09CandidateV := range []uint8{0, 2} {
		standard09CandidateVariant := map[uint8]string{0: "0900", 2: "0902"}[standard09CandidateV]
		standard09CandidateEntry := standard09CandidateGet(standard09CandidateT, standard09CandidateKey(standard09CandidateVariant, standard09CandidateV, "abstract-request", "count"))
		standard09CandidateInt(standard09CandidateT, standard09CandidateEntry.Standard09CandidateMaximum, 255)
		standard09CandidateInt(standard09CandidateT, standard09CandidateEntry.Standard09CandidateMinimum, 0)
		standard09CandidateInt(standard09CandidateT, standard09CandidateEntry.Standard09CandidateDefault, 0)
		standard09CandidateInt(standard09CandidateT, standard09CandidateEntry.Standard09CandidateWidth, 1)
		if standard09CandidateEntry.Standard09CandidateAbstractOffset != nil || standard09CandidateEntry.Standard09CandidateEncoding == nil || *standard09CandidateEntry.Standard09CandidateEncoding != "unsigned-u8" {
			standard09CandidateT.Fatal("count metadata changed")
		}
		if standard09CandidateV == 2 {
			standard09CandidateInt(standard09CandidateT, standard09CandidateEntry.Standard09CandidateAlternativeMaximum, 10)
		} else if standard09CandidateEntry.Standard09CandidateAlternativeMaximum != nil {
			standard09CandidateT.Fatal("conflict transferred")
		}
		for _, standard09CandidateRaw := range []byte{0, 10, 11, 255} {
			if int(standard09CandidateRaw) < *standard09CandidateEntry.Standard09CandidateMinimum || int(standard09CandidateRaw) > *standard09CandidateEntry.Standard09CandidateMaximum {
				standard09CandidateT.Fatalf("raw count %d clipped", standard09CandidateRaw)
			}
		}
	}
}

func TestAb79DeepCopies(standard09CandidateT *standard09CandidateTesting.T) {
	standard09CandidateBaseline := Standard09Candidates()
	standard09CandidateChanged := Standard09Candidates()
	standard09CandidateChanged.Standard09CandidateUnknowns[0] = "changed"
	for standard09CandidateIndex := range standard09CandidateChanged.Standard09CandidateEntries {
		standard09CandidateEntry := &standard09CandidateChanged.Standard09CandidateEntries[standard09CandidateIndex]
		*standard09CandidateEntry.Standard09CandidateKey.Standard09CandidatePrimary = 0
		*standard09CandidateEntry.Standard09CandidateKey.Standard09CandidateSecondary = 99
		for _, standard09CandidatePointer := range []*int{standard09CandidateEntry.Standard09CandidateAbstractOffset, standard09CandidateEntry.Standard09CandidateWidth, standard09CandidateEntry.Standard09CandidateDefault, standard09CandidateEntry.Standard09CandidateMinimum, standard09CandidateEntry.Standard09CandidateMaximum, standard09CandidateEntry.Standard09CandidateAlternativeMaximum} {
			if standard09CandidatePointer != nil {
				*standard09CandidatePointer = 99
			}
		}
		if standard09CandidateEntry.Standard09CandidateEncoding != nil {
			*standard09CandidateEntry.Standard09CandidateEncoding = "changed"
		}
		if len(standard09CandidateEntry.Standard09CandidateNotes) > 0 {
			standard09CandidateEntry.Standard09CandidateNotes[0] = "changed"
		}
	}
	standard09CandidateLookup := standard09CandidateGet(standard09CandidateT, standard09CandidateKey("0902", 2, "abstract-request", "count"))
	*standard09CandidateLookup.Standard09CandidateKey.Standard09CandidateSecondary = 0
	*standard09CandidateLookup.Standard09CandidateAlternativeMaximum = 99
	standard09CandidateLookup.Standard09CandidateNotes[0] = "changed"
	if !standard09CandidateReflect.DeepEqual(Standard09Candidates(), standard09CandidateBaseline) {
		standard09CandidateT.Fatal("shared nested storage")
	}
	standard09CandidateFresh := standard09CandidateGet(standard09CandidateT, standard09CandidateKey("0902", 2, "abstract-request", "count"))
	standard09CandidateInt(standard09CandidateT, standard09CandidateFresh.Standard09CandidateAlternativeMaximum, 10)
	if standard09CandidateFresh.Standard09CandidateNotes[0] == "changed" {
		standard09CandidateT.Fatal("lookup shares notes")
	}
}
