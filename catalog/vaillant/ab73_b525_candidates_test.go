package vaillant

import (
	ab73bytes "bytes"
	ab73fmt "fmt"
	ab73reflect "reflect"
	ab73testing "testing"
)

func ab73Pointer[ab73T any](ab73Value ab73T) *ab73T { return &ab73Value }

func ab73Key(ab73Op byte, ab73Direction, ab73Field string) Ab73Key {
	return Ab73Key{"Vaillant radio", nil, 0xb5, ab73Pointer(byte(0x25)), ab73Pointer(ab73Op), ab73Direction, ab73Field}
}

func ab73Get(ab73T *ab73testing.T, ab73Op byte, ab73Direction, ab73Field string) Ab73Entry {
	ab73T.Helper()
	ab73Entry, ab73OK := Ab73Lookup(ab73Key(ab73Op, ab73Direction, ab73Field))
	if !ab73OK {
		ab73T.Fatalf("ab73 missing %02x/%s/%s", ab73Op, ab73Direction, ab73Field)
	}
	return ab73Entry
}

func TestAb73ExactAxes(ab73T *ab73testing.T) {
	ab73Control := ab73Key(0, "Unknown", "layer")
	if ab73EqualOptional[byte](nil, ab73Pointer(byte(0))) || ab73EqualOptional[byte](ab73Pointer(byte(0)), nil) || !ab73EqualOptional(ab73Pointer(byte(0)), ab73Pointer(byte(0))) {
		ab73T.Fatal("ab73 optional absence collapsed into zero")
	}
	for ab73Name, ab73Mutate := range map[string]func(*Ab73Key){
		"ab73_family":           func(ab73K *Ab73Key) { ab73K.Ab73Family = "" },
		"ab73_variant":          func(ab73K *Ab73Key) { ab73K.Ab73Variant = ab73Pointer("") },
		"ab73_primary":          func(ab73K *Ab73Key) { ab73K.Ab73Primary = 0 },
		"ab73_secondary_absent": func(ab73K *Ab73Key) { ab73K.Ab73Secondary = nil },
		"ab73_secondary_zero":   func(ab73K *Ab73Key) { ab73K.Ab73Secondary = ab73Pointer(byte(0)) },
		"ab73_operation_absent": func(ab73K *Ab73Key) { ab73K.Ab73Operation = nil },
		"ab73_operation_other":  func(ab73K *Ab73Key) { ab73K.Ab73Operation = ab73Pointer(byte(1)) },
		"ab73_direction":        func(ab73K *Ab73Key) { ab73K.Ab73Direction = "request" },
		"ab73_field":            func(ab73K *Ab73Key) { ab73K.Ab73Field = "" },
	} {
		ab73T.Run(ab73Name, func(ab73T *ab73testing.T) {
			if _, ab73OK := Ab73Lookup(ab73Control); !ab73OK {
				ab73T.Fatal("ab73 invalid control")
			}
			ab73Changed := ab73Control
			ab73Mutate(&ab73Changed)
			if _, ab73OK := Ab73Lookup(ab73Changed); ab73OK {
				ab73T.Fatal("ab73 axis wildcard")
			}
		})
	}
	ab73Catalog := Ab73B525Candidates()
	if ab73Catalog.Ab73RepresentationGrade != "Proven" || ab73Catalog.Ab73PhysicalQualification != "Unknown" {
		ab73T.Fatal("ab73 grade conflation")
	}
	if len(ab73Catalog.Ab73Entries) != 127 {
		ab73T.Fatalf("ab73 finite count %d", len(ab73Catalog.Ab73Entries))
	}
	ab73Seen := map[string]bool{}
	for _, ab73Entry := range ab73Catalog.Ab73Entries {
		ab73Identity := ab73fmt.Sprintf("%s/%v/%d/%d/%d/%s/%s", ab73Entry.Ab73Key.Ab73Family, ab73Entry.Ab73Key.Ab73Variant, ab73Entry.Ab73Key.Ab73Primary, *ab73Entry.Ab73Key.Ab73Secondary, *ab73Entry.Ab73Key.Ab73Operation, ab73Entry.Ab73Key.Ab73Direction, ab73Entry.Ab73Key.Ab73Field)
		if ab73Seen[ab73Identity] {
			ab73T.Fatal("ab73 duplicate")
		}
		ab73Seen[ab73Identity] = true
		ab73Found, ab73OK := Ab73Lookup(ab73Entry.Ab73Key)
		if !ab73OK || !ab73reflect.DeepEqual(ab73Found, ab73Entry) || ab73Entry.Ab73Grade != "Hypothesis" || ab73Entry.Ab73RawPolicy != "preserve_all" {
			ab73T.Fatal("ab73 entry altered")
		}
	}
}

func TestAb73FreshNestedCopies(ab73T *ab73testing.T) {
	ab73Before := Ab73B525Candidates()
	ab73Changed := Ab73B525Candidates()
	ab73Changed.Ab73Limitations[0] = "ab73_changed"
	for ab73I := range ab73Changed.Ab73Entries {
		ab73E := &ab73Changed.Ab73Entries[ab73I]
		ab73E.Ab73Key.Ab73Family = "ab73_changed"
		*ab73E.Ab73Key.Ab73Secondary = 0
		*ab73E.Ab73Key.Ab73Operation = 255
		for _, ab73P := range []*int{ab73E.Ab73Position, ab73E.Ab73Width, ab73E.Ab73Minimum, ab73E.Ab73Maximum, ab73E.Ab73Default, ab73E.Ab73Count, ab73E.Ab73StartIndex} {
			if ab73P != nil {
				*ab73P = 999
			}
		}
		ab73E.Ab73Mappings[0] = "ab73_changed"
		ab73E.Ab73Bits[0] = "ab73_changed"
		if len(ab73E.Ab73Unknowns) > 0 {
			ab73E.Ab73Unknowns[0] = "ab73_changed"
		}
	}
	if !ab73reflect.DeepEqual(ab73Before, Ab73B525Candidates()) {
		ab73T.Fatal("ab73 accessor shared backing")
	}
	ab73Lookup := ab73Get(ab73T, 0, "Unknown", "sender-role")
	ab73Lookup.Ab73Mappings[1] = "ab73_changed"
	*ab73Lookup.Ab73Key.Ab73Operation = 255
	if !ab73reflect.DeepEqual(ab73Before, Ab73B525Candidates()) || ab73Get(ab73T, 0, "Unknown", "sender-role").Ab73Mappings[1] != "remote-eight" {
		ab73T.Fatal("ab73 lookup shared backing")
	}
}

func TestAb73RolesAndHeader(ab73T *ab73testing.T) {
	ab73Codes := []int{1, 15, 16, 21, 22, 23, 48, 55, 112, 119, 240, 247, 253, 254}
	ab73Labels := []string{"remote-eight", "outdoor-sensor", "primary-control", "base-station", "base-outdoor-sensor", "remote-four", "remote-one", "remote-five", "remote-two", "remote-six", "remote-three", "remote-seven", "diagnosis-tool", "range-extender"}
	for ab73I, ab73Name := range []string{"sender-role", "receiver-role"} {
		ab73E := ab73Get(ab73T, 0, "Unknown", ab73Name)
		if *ab73E.Ab73Position != 7+ab73I || *ab73E.Ab73Width != 1 || *ab73E.Ab73Default != 1 || len(ab73E.Ab73Mappings) != 14 {
			ab73T.Fatal("ab73 role extent")
		}
		for ab73J, ab73Code := range ab73Codes {
			if ab73E.Ab73Mappings[ab73Code] != ab73Labels[ab73J] {
				ab73T.Fatal("ab73 role reversal")
			}
		}
		if _, ab73OK := ab73E.Ab73Mappings[0]; ab73OK {
			ab73T.Fatal("ab73 unmapped 00")
		}
	}
	ab73Sender := ab73Get(ab73T, 0, "Unknown", "sender-role")
	ab73Receiver := ab73Get(ab73T, 0, "Unknown", "receiver-role")
	if ab73Sender.Ab73Mappings[1] != "remote-eight" || ab73Receiver.Ab73Mappings[254] != "range-extender" {
		ab73T.Fatal("ab73 roles collapsed")
	}
	for _, ab73Case := range []struct {
		ab73Name                    string
		ab73Pos, ab73Width, ab73Max int
		ab73Codec                   string
	}{{"layer", 3, 1, 255, "u8"}, {"radio-domain", 4, 2, 65535, "u16-le"}, {"request-id", 9, 1, 255, "u8"}} {
		ab73E := ab73Get(ab73T, 0, "Unknown", ab73Case.ab73Name)
		if *ab73E.Ab73Position != ab73Case.ab73Pos || *ab73E.Ab73Width != ab73Case.ab73Width || *ab73E.Ab73Maximum != ab73Case.ab73Max || *ab73E.Ab73Minimum != 0 || *ab73E.Ab73Default != 0 || ab73E.Ab73Codec != ab73Case.ab73Codec {
			ab73T.Fatal("ab73 OP00 metadata")
		}
	}
	ab73Status := ab73Get(ab73T, 0, "Unknown", "status")
	for ab73I, ab73Name := range []string{"repeat", "group", "encrypted", "frame-type", "outdoor-sensor-pairing", "reserved", "reserved", "reserved"} {
		if ab73Status.Ab73Bits[ab73I] != ab73Name || *ab73Status.Ab73Default != 0 {
			ab73T.Fatal("ab73 status bits")
		}
	}
}

func ab73Signed(ab73E Ab73Entry, ab73Raw byte) int {
	if ab73E.Ab73Codec == "signed8" {
		return int(int8(ab73Raw))
	}
	return int(ab73Raw)
}

func TestAb73RSSIAndFeedback(ab73T *ab73testing.T) {
	ab73RSSI := ab73Get(ab73T, 3, "Unknown", "RSSI")
	if ab73Signed(ab73RSSI, 255) != -1 || ab73RSSI.Ab73Unit != nil || *ab73RSSI.Ab73Minimum != -128 || *ab73RSSI.Ab73Maximum != 127 || *ab73RSSI.Ab73Default != 0 {
		ab73T.Fatal("ab73 RSSI signed carrier or unit")
	}
	ab73Link := ab73Get(ab73T, 3, "Unknown", "link-quality")
	if len(ab73Link.Ab73Bits) != 8 || *ab73Link.Ab73Default != 0 {
		ab73T.Fatal("ab73 link bits")
	}
	for ab73I := 0; ab73I < 8; ab73I++ {
		ab73Want := ab73fmt.Sprintf("link-quality-%d", ab73I)
		if ab73I == 7 {
			ab73Want = "CRC-okay"
		}
		if ab73Link.Ab73Bits[ab73I] != ab73Want {
			ab73T.Fatal("ab73 link bit independence")
		}
	}
	ab73Feedback := ab73Get(ab73T, 4, "positive", "feedback")
	ab73Flags := []string{}
	for ab73I := 0; ab73I < 8; ab73I++ {
		if byte(7)&(1<<ab73I) != 0 {
			ab73Flags = append(ab73Flags, ab73Feedback.Ab73Bits[ab73I])
		}
	}
	if !ab73reflect.DeepEqual(ab73Flags, []string{"success", "too-short", "too-long"}) || ab73Feedback.Ab73Relation != "independent-bits" {
		ab73T.Fatal("ab73 feedback collapsed")
	}
	for _, ab73D := range []string{"request", "positive"} {
		ab73E := ab73Get(ab73T, 4, ab73D, "bytes")
		if *ab73E.Ab73Count != 16 || *ab73E.Ab73Width != 16 || *ab73E.Ab73Default != 0 || ab73E.Ab73Relation != "Unknown" || ab73E.Ab73Position != nil {
			ab73T.Fatal("ab73 array contract/echo")
		}
		for ab73I := 0; ab73I < 16; ab73I++ {
			ab73Raw := make([]byte, 16)
			ab73Raw[ab73I] = 0xaa
			for ab73J, ab73B := range ab73Raw {
				ab73Want := 0
				if ab73I == ab73J {
					ab73Want = 170
				}
				if ab73Signed(ab73E, ab73B) != ab73Want {
					ab73T.Fatal("ab73 independent raw byte")
				}
			}
		}
	}
}

func TestAb73SignedBlocks(ab73T *ab73testing.T) {
	for ab73Op := byte(16); ab73Op <= 23; ab73Op++ {
		ab73E := ab73Get(ab73T, ab73Op, "request", "values")
		if *ab73E.Ab73Position != 3 || *ab73E.Ab73Width != 32 || *ab73E.Ab73Count != 32 || *ab73E.Ab73StartIndex != int(ab73Op-16)*32 || ab73E.Ab73Unit != nil || len(ab73E.Ab73Mappings) != 0 || !ab73reflect.DeepEqual(ab73E.Ab73Unknowns, []string{"frequency-mapping", "scanning-effects"}) || *ab73E.Ab73Minimum != -128 || *ab73E.Ab73Maximum != 127 || *ab73E.Ab73Default != 0 {
			ab73T.Fatal("ab73 block alignment/frequency")
		}
		for _, ab73Length := range []int{31, 32, 33} {
			if (ab73Length == *ab73E.Ab73Count) != (ab73Length == 32) {
				ab73T.Fatal("ab73 incomplete/extra accepted")
			}
		}
		for ab73I := 0; ab73I < 32; ab73I++ {
			for _, ab73Case := range []struct {
				ab73Raw  byte
				ab73Want int
			}{{128, -128}, {0, 0}, {127, 127}} {
				ab73Raw := make([]byte, 32)
				ab73Raw[ab73I] = ab73Case.ab73Raw
				for ab73J, ab73B := range ab73Raw {
					ab73Want := 0
					if ab73J == ab73I {
						ab73Want = ab73Case.ab73Want
					}
					if ab73Signed(ab73E, ab73B) != ab73Want {
						ab73T.Fatal("ab73 signed array")
					}
				}
			}
		}
		for _, ab73Row := range Ab73B525Candidates().Ab73Entries {
			if *ab73Row.Ab73Key.Ab73Operation == ab73Op && ab73Row.Ab73Key.Ab73Direction == "positive" && ab73Row.Ab73Key.Ab73Field != "primary" && ab73Row.Ab73Key.Ab73Field != "secondary" && ab73Row.Ab73Key.Ab73Field != "selector" {
				ab73T.Fatal("ab73 positive is prefix-only")
			}
		}
	}
}

func TestAb73OpaqueAndConfiguration(ab73T *ab73testing.T) {
	for _, ab73Op := range []byte{1, 2, 24, 25, 26, 27} {
		ab73Directions := []string{"Unknown"}
		if ab73Op >= 24 {
			ab73Directions = []string{"request", "positive"}
		}
		for _, ab73D := range ab73Directions {
			ab73E := ab73Get(ab73T, ab73Op, ab73D, "data")
			if ab73E.Ab73Codec != "opaque" || *ab73E.Ab73Minimum != 0 || ab73E.Ab73Maximum != nil || ab73E.Ab73Width != nil || ab73E.Ab73Relation != "Unknown" {
				ab73T.Fatal("ab73 opaque bounds")
			}
			for _, ab73Raw := range [][]byte{{}, {0, 170, 255}} {
				ab73Copy := ab73bytes.Clone(ab73Raw)
				if !ab73bytes.Equal(ab73Raw, ab73Copy) || len(ab73Copy) < *ab73E.Ab73Minimum {
					ab73T.Fatal("ab73 opaque representation")
				}
			}
			if ab73Op < 24 {
				if *ab73E.Ab73Position != 3 {
					ab73T.Fatal("ab73 opaque position")
				}
				continue
			}
			if ab73E.Ab73Position != nil {
				ab73T.Fatal("ab73 invented layout")
			}
			ab73Config := ab73Get(ab73T, ab73Op, ab73D, "configuration")
			if *ab73Config.Ab73Position != 3 || *ab73Config.Ab73Default != 0 || ab73Config.Ab73Bits[0] != "configuration" || len(ab73Config.Ab73Bits) != 1 || ab73Config.Ab73Mappings[0] != "read" || len(ab73Config.Ab73Mappings) != 1 || ab73Config.Ab73RawPolicy != "preserve_all" {
				ab73T.Fatal("ab73 configuration metadata")
			}
			for _, ab73Raw := range []byte{0, 1, 0xaa, 0xff} {
				ab73Known := ab73Raw & 1
				ab73Unknown := ab73Raw & 0xfe
				if ab73Known|ab73Unknown != ab73Raw {
					ab73T.Fatal("ab73 unknown bits lost")
				}
			}
		}
	}
}

func TestAb73PrefixAndUnknowns(ab73T *ab73testing.T) {
	for _, ab73E := range Ab73B525Candidates().Ab73Entries {
		for ab73I, ab73Name := range []string{"primary", "secondary", "selector"} {
			if ab73E.Ab73Key.Ab73Field != ab73Name {
				continue
			}
			ab73Value := []int{181, 37, int(*ab73E.Ab73Key.Ab73Operation)}[ab73I]
			if *ab73E.Ab73Position != ab73I || *ab73E.Ab73Width != 1 || *ab73E.Ab73Minimum != ab73Value || *ab73E.Ab73Maximum != ab73Value || ab73E.Ab73Default != nil {
				ab73T.Fatal("ab73 prefix changed")
			}
		}
	}
	for _, ab73C := range []struct {
		ab73Op       byte
		ab73D, ab73F string
		ab73Unknowns []string
	}{
		{0, "Unknown", "request-id", []string{"lifecycle"}},
		{1, "Unknown", "data", []string{"cipher", "padding", "sender-identity"}},
		{2, "Unknown", "data", []string{"cipher", "padding", "sender-identity"}},
		{3, "Unknown", "RSSI", []string{"physical-unit"}},
		{3, "Unknown", "link-quality", []string{"CRC-algorithm"}},
		{4, "request", "bytes", []string{"active-key", "persistent-key", "echo"}},
		{4, "positive", "bytes", []string{"active-key", "persistent-key", "echo"}},
		{4, "positive", "feedback", []string{"exclusivity"}},
	} {
		if !ab73reflect.DeepEqual(ab73Get(ab73T, ab73C.ab73Op, ab73C.ab73D, ab73C.ab73F).Ab73Unknowns, ab73C.ab73Unknowns) {
			ab73T.Fatal("ab73 unknown became claim")
		}
	}
	for ab73Op := byte(24); ab73Op <= 27; ab73Op++ {
		for _, ab73D := range []string{"request", "positive"} {
			if !ab73reflect.DeepEqual(ab73Get(ab73T, ab73Op, ab73D, "configuration").Ab73Unknowns, []string{"bits-1..7", "effects"}) || !ab73reflect.DeepEqual(ab73Get(ab73T, ab73Op, ab73D, "data").Ab73Unknowns, []string{"layout", "units", "correlation", "effects"}) {
				ab73T.Fatal("ab73 configuration claims")
			}
		}
	}
}
