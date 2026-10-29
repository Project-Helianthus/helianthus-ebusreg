package ebus_standard_catalog

import (
	ab77reflect "reflect"
	ab77testing "testing"
)

func ab77Byte(ab77Value byte) *byte { return &ab77Value }

func ab77IdentifierKey() Ab77Key {
	return Ab77Key{Ab77Family: "standard", Ab77Variant: "identification", Ab77Primary: 7, Ab77Secondary: ab77Byte(4), Ab77Direction: "response", Ab77Field: "identifier"}
}

func TestAb77ExactKey(ab77T *ab77testing.T) {
	ab77Control := ab77IdentifierKey()
	ab77Entry, ab77OK := Ab77Lookup(ab77Control)
	if !ab77OK || ab77Entry.Ab77LocalParser == nil {
		ab77T.Fatal("valid independent key missing")
	}
	for _, ab77Axis := range []string{"family", "variant", "primary", "secondary", "secondary_absent", "secondary_zero", "operation_zero", "direction", "field"} {
		ab77T.Run(ab77Axis, func(ab77T *ab77testing.T) {
			ab77Key := ab77IdentifierKey()
			switch ab77Axis {
			case "family":
				ab77Key.Ab77Family = ""
			case "variant":
				ab77Key.Ab77Variant = "existence_inquiry"
			case "primary":
				ab77Key.Ab77Primary = 0
			case "secondary":
				ab77Key.Ab77Secondary = ab77Byte(254)
			case "secondary_absent":
				ab77Key.Ab77Secondary = nil
			case "secondary_zero":
				ab77Key.Ab77Secondary = ab77Byte(0)
			case "operation_zero":
				ab77Key.Ab77Operation = ab77Byte(0)
			case "direction":
				ab77Key.Ab77Direction = "request"
			case "field":
				ab77Key.Ab77Field = "missing"
			}
			if _, ab77OK := Ab77Lookup(ab77Key); ab77OK {
				ab77T.Fatal("changed axis matched")
			}
			if _, ab77OK := Ab77Lookup(ab77Control); !ab77OK {
				ab77T.Fatal("control lost")
			}
		})
	}
	if ab77EqualByte(nil, ab77Byte(0)) || ab77EqualByte(ab77Byte(0), nil) || !ab77EqualByte(ab77Byte(0), ab77Byte(0)) {
		ab77T.Fatal("absent/zero equality")
	}
}

func TestAb77Metadata(ab77T *ab77testing.T) {
	ab77Data := Ab77CandidateMetadata()
	if ab77Data.Ab77UnknownGrade != "Unknown" || ab77Data.Ab77PresenceGrade != "Proven" || len(ab77Data.Ab77Entries) != 7 || len(ab77Data.Ab77Unknown) != 5 || ab77Data.Ab77Scope == "" {
		ab77T.Fatal("catalog boundary")
	}
	for ab77Index, ab77Field := range []string{"identifier", "major", "minor", "product_number", "version_control", "revision_control", "existence_inquiry"} {
		ab77Key := ab77IdentifierKey()
		ab77Key.Ab77Field = ab77Field
		if ab77Index == 6 {
			ab77Key.Ab77Variant = "existence_inquiry"
			ab77Key.Ab77Secondary = ab77Byte(254)
			ab77Key.Ab77Direction = "request"
		}
		ab77Entry, ab77OK := Ab77Lookup(ab77Key)
		if !ab77OK || !ab77reflect.DeepEqual(ab77Entry, ab77Data.Ab77Entries[ab77Index]) || ab77Entry.Ab77Grade != "Hypothesis" || ab77Entry.Ab77Width != nil {
			ab77T.Fatalf("field %s", ab77Field)
		}
		if ab77Index == 4 || ab77Index == 5 {
			if ab77Entry.Ab77LogicalPosition == nil || *ab77Entry.Ab77LogicalPosition != ab77Index+4 || ab77Entry.Ab77Control == nil {
				ab77T.Fatal("logical overlap missing")
			}
			ab77Control := ab77Entry.Ab77Control
			if ab77Control.Ab77NotOkayMinimum != 0 || ab77Control.Ab77NotOkayMaximum != 98 || ab77Control.Ab77Default != "not_okay" || ab77Control.Ab77OkayValues != nil {
				ab77T.Fatal("unsupported okay claim")
			}
		} else if ab77Entry.Ab77LogicalPosition != nil || ab77Entry.Ab77Control != nil {
			ab77T.Fatal("invented native coordinates")
		}
		if ab77Index == 6 {
			if ab77Entry.Ab77DirectReplyRequired == nil || *ab77Entry.Ab77DirectReplyRequired {
				ab77T.Fatal("existence direct reply required")
			}
		} else if ab77Entry.Ab77DirectReplyRequired != nil {
			ab77T.Fatal("invented answer policy")
		}
		if (ab77Entry.Ab77LocalParser != nil) != (ab77Index == 0) {
			ab77T.Fatal("local parser attached to wrong field")
		}
	}
}

func TestAb77FreshCopies(ab77T *ab77testing.T) {
	ab77Before := Ab77CandidateMetadata()
	ab77Changed := Ab77CandidateMetadata()
	ab77Changed.Ab77Unknown[0] = "changed"
	*ab77Changed.Ab77Entries[0].Ab77Key.Ab77Secondary = 0
	ab77Changed.Ab77Entries[0].Ab77LocalParser.Ab77NumericOffsets[0] = 0
	ab77Changed.Ab77Entries[0].Ab77LocalParser.Ab77ProductOffsets[0] = 0
	*ab77Changed.Ab77Entries[4].Ab77LogicalPosition = 0
	ab77Changed.Ab77Entries[4].Ab77Control.Ab77Default = "changed"
	*ab77Changed.Ab77Entries[6].Ab77DirectReplyRequired = true
	ab77Changed.Ab77Entries[0].Ab77Key.Ab77Family = "changed"
	if !ab77reflect.DeepEqual(ab77Before, Ab77CandidateMetadata()) {
		ab77T.Fatal("shared accessor storage")
	}
	ab77Lookup, ab77OK := Ab77Lookup(ab77IdentifierKey())
	if !ab77OK {
		ab77T.Fatal("missing lookup")
	}
	ab77Lookup.Ab77LocalParser.Ab77NumericOffsets[0] = 0
	*ab77Lookup.Ab77Key.Ab77Secondary = 0
	ab77Again, ab77OK := Ab77Lookup(ab77IdentifierKey())
	if !ab77OK || !ab77reflect.DeepEqual(ab77Again, ab77Before.Ab77Entries[0]) {
		ab77T.Fatal("shared lookup storage")
	}
}

func ab77Evaluate(ab77Parser Ab77LocalParser, ab77Bytes []byte, ab77Count int) (string, int, int, int, bool) {
	if ab77Count < ab77Parser.Ab77MinimumDeclaredCount || len(ab77Bytes) < ab77Parser.Ab77MinimumAvailableCount || (ab77Parser.Ab77MaximumDeclaredCount != nil && ab77Count > *ab77Parser.Ab77MaximumDeclaredCount) {
		return "", 0, 0, 0, false
	}
	for ab77Index := 0; ab77Index < ab77Parser.Ab77IdentifierWidth; ab77Index++ {
		ab77Byte := ab77Bytes[ab77Parser.Ab77IdentifierOffset+ab77Index]
		if ab77Index == 0 && ab77Parser.Ab77FirstIdentifierNonzero && ab77Byte == 0 {
			return "", 0, 0, 0, false
		}
		if (ab77Byte != 0 || !ab77Parser.Ab77CharacterZeroAllowed) && (ab77Byte < ab77Parser.Ab77CharacterMinimum || ab77Byte > ab77Parser.Ab77CharacterMaximum) {
			return "", 0, 0, 0, false
		}
	}
	for _, ab77Offset := range ab77Parser.Ab77NumericOffsets {
		ab77Byte := ab77Bytes[ab77Offset]
		if ab77Byte&15 > ab77Parser.Ab77LowNibbleMaximum || (ab77Parser.Ab77HighNibbleMaximum != nil && ab77Byte>>4 > *ab77Parser.Ab77HighNibbleMaximum) {
			return "", 0, 0, 0, false
		}
	}
	ab77Decimal := func(ab77Byte byte) int { return ab77Parser.Ab77DecimalMultiplier*int(ab77Byte>>4) + int(ab77Byte&15) }
	return string(ab77Bytes[ab77Parser.Ab77IdentifierOffset : ab77Parser.Ab77IdentifierOffset+ab77Parser.Ab77IdentifierWidth]), ab77Decimal(ab77Bytes[ab77Parser.Ab77MajorOffset]), ab77Decimal(ab77Bytes[ab77Parser.Ab77MinorOffset]), ab77Decimal(ab77Bytes[ab77Parser.Ab77ProductOffsets[0]]) + ab77Parser.Ab77ProductMultiplier*ab77Decimal(ab77Bytes[ab77Parser.Ab77ProductOffsets[1]]), true
}

func TestAb77LocalBoundaries(ab77T *ab77testing.T) {
	ab77Entry, ab77OK := Ab77Lookup(ab77IdentifierKey())
	if !ab77OK || ab77Entry.Ab77LocalParser == nil {
		ab77T.Fatal("missing parser metadata")
	}
	ab77Parser := *ab77Entry.Ab77LocalParser
	ab77Expected := Ab77LocalParser{Ab77MinimumDeclaredCount: 10, Ab77MinimumAvailableCount: 10, Ab77IdentifierOffset: 1, Ab77IdentifierWidth: 5, Ab77FirstIdentifierNonzero: true, Ab77CharacterZeroAllowed: true, Ab77CharacterMinimum: 32, Ab77CharacterMaximum: 126, Ab77NumericOffsets: []int{6, 7, 8, 9}, Ab77LowNibbleMaximum: 9, Ab77DecimalMultiplier: 10, Ab77MajorOffset: 6, Ab77MinorOffset: 7, Ab77ProductOffsets: []int{8, 9}, Ab77ProductMultiplier: 100, Ab77Grade: "Hypothesis", Ab77Scope: "partial local parser only; not physical BCD qualification"}
	if !ab77reflect.DeepEqual(ab77Parser, ab77Expected) {
		ab77T.Fatal("local predicate differs from supplied boundary")
	}
	for _, ab77Case := range []string{"valid", "count9", "count11", "short", "first_zero", "internal_zero", "low_character", "high_character", "FA", "byte_zero"} {
		ab77T.Run(ab77Case, func(ab77T *ab77testing.T) {
			ab77Bytes := []byte{0xA1, 0x41, 0x42, 0x43, 0x44, 0x45, 0xF9, 0x99, 0x17, 0x04}
			ab77Count := 10
			ab77WantOK := true
			ab77WantID := "ABCDE"
			switch ab77Case {
			case "count9":
				ab77Count = 9
				ab77WantOK = false
			case "count11":
				ab77Count = 11
			case "short":
				ab77Bytes = ab77Bytes[:9]
				ab77Count = 11
				ab77WantOK = false
			case "first_zero":
				ab77Bytes[1] = 0
				ab77WantOK = false
			case "internal_zero":
				ab77Bytes[2] = 0
				ab77WantID = "A\x00CDE"
			case "low_character":
				ab77Bytes[3] = 0x1F
				ab77WantOK = false
			case "high_character":
				ab77Bytes[3] = 0x7F
				ab77WantOK = false
			case "FA":
				ab77Bytes[6] = 0xFA
				ab77WantOK = false
			case "byte_zero":
				ab77Bytes[0] = 0
			}
			ab77ID, ab77Major, ab77Minor, ab77Product, ab77OK := ab77Evaluate(ab77Parser, ab77Bytes, ab77Count)
			if ab77OK != ab77WantOK {
				ab77T.Fatal("admission")
			}
			if ab77OK && (ab77ID != ab77WantID || ab77Major != 159 || ab77Minor != 99 || ab77Product != 417) {
				ab77T.Fatalf("result %q %d %d %d", ab77ID, ab77Major, ab77Minor, ab77Product)
			}
		})
	}
	for _, ab77Offset := range []int{6, 7, 8, 9} {
		ab77Bytes := []byte{0xA1, 0x41, 0x42, 0x43, 0x44, 0x45, 0xF9, 0x99, 0x17, 0x04}
		ab77Bytes[ab77Offset] = 0xFA
		if _, _, _, _, ab77OK := ab77Evaluate(ab77Parser, ab77Bytes, 10); ab77OK {
			ab77T.Fatalf("low nibble accepted at %d", ab77Offset)
		}
	}
	for ab77Length := 0; ab77Length < 10; ab77Length++ {
		if _, _, _, _, ab77OK := ab77Evaluate(ab77Parser, make([]byte, ab77Length), 10); ab77OK {
			ab77T.Fatal("short buffer admitted")
		}
	}
}
