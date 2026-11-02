package vaillant

import (
	ab98json "encoding/json"
	ab98reflect "reflect"
	ab98strings "strings"
	ab98testing "testing"
)

func ab98Pointer[ab98Type any](ab98Value ab98Type) *ab98Type { return &ab98Value }

func ab98Get(ab98T *ab98testing.T, ab98Op uint8, ab98Direction, ab98Field string) Ab98Entry {
	ab98T.Helper()
	ab98Entry, ab98Found, ab98Err := Ab98Lookup(Ab98Key{Ab98Family: "vaillant", Ab98Variant: "B51A", Ab98Primary: 181, Ab98Secondary: ab98Pointer(uint8(26)), Ab98Operation: ab98Pointer(ab98Op), Ab98Direction: ab98Direction, Ab98Field: ab98Field})
	if ab98Err != nil || !ab98Found {
		ab98T.Fatalf("lookup %d/%s/%s: %v %v", ab98Op, ab98Direction, ab98Field, ab98Found, ab98Err)
	}
	return ab98Entry
}

func ab98Equal(ab98T *ab98testing.T, ab98Got, ab98Want any) {
	ab98T.Helper()
	if !ab98reflect.DeepEqual(ab98Got, ab98Want) {
		ab98T.Fatalf("got %#v, want %#v", ab98Got, ab98Want)
	}
}

func TestAb98CoordinatesAndGrades(ab98T *ab98testing.T) {
	ab98Entries, ab98Err := Ab98Entries()
	if ab98Err != nil {
		ab98T.Fatal(ab98Err)
	}
	ab98Equal(ab98T, len(ab98Entries), 37)
	for _, ab98Entry := range ab98Entries {
		ab98Equal(ab98T, ab98Entry.Ab98Grade, "Hypothesis")
		ab98Equal(ab98T, ab98Entry.Ab98MetadataPresence, "Proven")
		ab98Equal(ab98T, ab98Entry.Ab98PhysicalQualification, "Unknown")
		ab98Equal(ab98T, ab98Entry.Ab98Invocation, "Unknown")
		ab98Equal(ab98T, ab98Entry.Ab98NativeContract, "https://github.com/Project-Helianthus/helianthus-docs-ebus/blob/main/protocols/ebus-services/ebus-overview.md")
	}
	for _, ab98Op := range []uint8{1, 4, 5, 7, 8} {
		for _, ab98Direction := range []string{"request", "positive"} {
			ab98Entry := ab98Get(ab98T, ab98Op, ab98Direction, "sequence")
			ab98Equal(ab98T, ab98Entry.Ab98Position, ab98Pointer(uint16(3)))
			ab98Equal(ab98T, ab98Entry.Ab98WidthBits, ab98Pointer(uint16(8)))
			ab98Equal(ab98T, ab98Entry.Ab98Domain, []uint64{0, 255})
			ab98Equal(ab98T, ab98Entry.Ab98Default, ab98Pointer(uint64(0)))
		}
	}
}

func TestAb98KeyAxes(ab98T *ab98testing.T) {
	ab98Base := ab98Get(ab98T, 1, "request", "CRC").Ab98Key
	ab98Changes := []struct {
		ab98Name   string
		ab98Change func(*Ab98Key)
	}{
		{"family", func(ab98Key *Ab98Key) { ab98Key.Ab98Family = "other-family" }},
		{"variant", func(ab98Key *Ab98Key) { ab98Key.Ab98Variant = "B524" }},
		{"primary", func(ab98Key *Ab98Key) { ab98Key.Ab98Primary = 0 }},
		{"secondary", func(ab98Key *Ab98Key) { ab98Key.Ab98Secondary = ab98Pointer(uint8(0)) }},
		{"operation", func(ab98Key *Ab98Key) { ab98Key.Ab98Operation = ab98Pointer(uint8(0)) }},
		{"direction", func(ab98Key *Ab98Key) { ab98Key.Ab98Direction = "positive" }},
		{"field", func(ab98Key *Ab98Key) { ab98Key.Ab98Field = "unlisted" }},
	}
	for _, ab98Case := range ab98Changes {
		ab98T.Run(ab98Case.ab98Name, func(ab98T *ab98testing.T) {
			ab98Key := ab98Base
			ab98Case.ab98Change(&ab98Key)
			_, ab98Found, ab98Err := Ab98Lookup(ab98Key)
			if ab98Found || ab98Err != nil {
				ab98T.Fatalf("nonmatching key: %v %v", ab98Found, ab98Err)
			}
		})
	}
	ab98Original := ab98Data
	defer func() { ab98Data = ab98Original }()
	ab98Synthetic := ab98Get(ab98T, 1, "request", "flags")
	ab98Synthetic.Ab98Key = Ab98Key{Ab98Family: "synthetic-family", Ab98Variant: "B524", Ab98Primary: 181, Ab98Secondary: ab98Pointer(uint8(36)), Ab98Operation: ab98Pointer(uint8(3)), Ab98Direction: "request", Ab98Field: "synthetic-field", Ab98Group: ab98Pointer(uint8(7)), Ab98RequestedInstance: ab98Pointer(uint8(9)), Ab98Register: ab98Pointer(uint16(513)), Ab98CommandRole: ab98Pointer("read")}
	ab98Changes[5].ab98Change = func(ab98Key *Ab98Key) { ab98Key.Ab98Direction = "positive" }
	ab98Changes[1].ab98Change = func(ab98Key *Ab98Key) { ab98Key.Ab98Variant = "B51A" }
	ab98Changes = append(ab98Changes,
		struct {
			ab98Name   string
			ab98Change func(*Ab98Key)
		}{"group", func(ab98Key *Ab98Key) { ab98Key.Ab98Group = ab98Pointer(uint8(0)) }},
		struct {
			ab98Name   string
			ab98Change func(*Ab98Key)
		}{"requested-instance", func(ab98Key *Ab98Key) { ab98Key.Ab98RequestedInstance = ab98Pointer(uint8(0)) }},
		struct {
			ab98Name   string
			ab98Change func(*Ab98Key)
		}{"register", func(ab98Key *Ab98Key) { ab98Key.Ab98Register = ab98Pointer(uint16(257)) }},
		struct {
			ab98Name   string
			ab98Change func(*Ab98Key)
		}{"command-role", func(ab98Key *Ab98Key) { ab98Key.Ab98CommandRole = ab98Pointer("write") }},
	)
	ab98Input, ab98Err := ab98json.Marshal([]Ab98Entry{ab98Synthetic})
	if ab98Err != nil {
		ab98T.Fatal(ab98Err)
	}
	ab98SetData(ab98Input)
	ab98Hit, ab98Found, ab98Err := Ab98Lookup(ab98Synthetic.Ab98Key)
	if ab98Err != nil || !ab98Found {
		ab98T.Fatalf("synthetic baseline: %v %v", ab98Found, ab98Err)
	}
	ab98Equal(ab98T, ab98Hit, ab98Synthetic)
	for _, ab98Case := range ab98Changes {
		ab98T.Run("B524/"+ab98Case.ab98Name, func(ab98T *ab98testing.T) {
			ab98Key := ab98Synthetic.Ab98Key
			ab98Case.ab98Change(&ab98Key)
			_, ab98Found, ab98Err := Ab98Lookup(ab98Key)
			if ab98Found || ab98Err != nil {
				ab98T.Fatalf("axis ignored: %v %v", ab98Found, ab98Err)
			}
		})
	}
	for _, ab98Axis := range []string{"secondary", "operation", "group", "requested-instance", "register", "command-role"} {
		ab98T.Run("absence/"+ab98Axis, func(ab98T *ab98testing.T) {
			ab98Nil := ab98Synthetic
			ab98Zero := ab98Synthetic
			switch ab98Axis {
			case "secondary":
				ab98Nil.Ab98Key.Ab98Secondary = nil
				ab98Zero.Ab98Key.Ab98Secondary = ab98Pointer(uint8(0))
			case "operation":
				ab98Nil.Ab98Key.Ab98Operation = nil
				ab98Zero.Ab98Key.Ab98Operation = ab98Pointer(uint8(0))
			case "group":
				ab98Nil.Ab98Key.Ab98Group = nil
				ab98Zero.Ab98Key.Ab98Group = ab98Pointer(uint8(0))
			case "requested-instance":
				ab98Nil.Ab98Key.Ab98RequestedInstance = nil
				ab98Zero.Ab98Key.Ab98RequestedInstance = ab98Pointer(uint8(0))
			case "register":
				ab98Nil.Ab98Key.Ab98Register = nil
				ab98Zero.Ab98Key.Ab98Register = ab98Pointer(uint16(0))
			case "command-role":
				ab98Nil.Ab98Key.Ab98CommandRole = nil
				ab98Zero.Ab98Key.Ab98CommandRole = ab98Pointer("")
			}
			ab98Nil.Ab98Interpretations = []string{"absent"}
			ab98Zero.Ab98Interpretations = []string{"zero"}
			ab98Input, ab98Err := ab98json.Marshal([]Ab98Entry{ab98Nil, ab98Zero})
			if ab98Err != nil {
				ab98T.Fatal(ab98Err)
			}
			ab98SetData(ab98Input)
			for _, ab98Want := range []Ab98Entry{ab98Nil, ab98Zero} {
				ab98Got, ab98Found, ab98Err := Ab98Lookup(ab98Want.Ab98Key)
				if ab98Err != nil || !ab98Found {
					ab98T.Fatalf("absence: %v %v", ab98Found, ab98Err)
				}
				ab98Equal(ab98T, ab98Got, ab98Want)
			}
		})
	}
}

func ab98SetData(ab98Input []byte) { ab98Data = ab98Input }

func TestAb98FlagsAndConditionalFields(ab98T *ab98testing.T) {
	for _, ab98Case := range []struct {
		ab98Op        uint8
		ab98Direction string
		ab98Labels    map[uint8][]string
	}{
		{1, "request", map[uint8][]string{0: {"power-on/start"}, 1: {"change"}, 2: {"CRC-support"}}},
		{1, "positive", map[uint8][]string{0: {"valid"}, 1: {"invalid"}, 2: {"blocked"}, 3: {"change"}, 4: {"denied"}, 5: {"invalid-range"}, 6: {"not-user-adjustable"}}},
		{5, "positive", map[uint8][]string{0: {"disabled; following data omitted"}, 1: {"writable"}, 2: {"enum-vs-range"}, 3: {"additional-four-byte/no-range", "read-only over writable"}, 6: {"text alternative 6"}, 7: {"text alternative 7"}}},
		{7, "request", map[uint8][]string{0: {"normal-replication"}, 1: {"inverted-replication"}, 2: {"refresh"}}},
	} {
		ab98Entry := ab98Get(ab98T, ab98Case.ab98Op, ab98Case.ab98Direction, "flags")
		ab98Equal(ab98T, len(ab98Entry.Ab98Bits), 8)
		ab98Equal(ab98T, ab98Entry.Ab98Exclusive, ab98Pointer(false))
		for ab98Index := uint8(0); ab98Index < 8; ab98Index++ {
			ab98Bit := ab98Entry.Ab98Bits[ab98Index]
			ab98Labels, ab98Named := ab98Case.ab98Labels[ab98Index]
			ab98Equal(ab98T, ab98Bit.Ab98Bit, ab98Index)
			ab98Equal(ab98T, ab98Bit.Ab98Reserved, !ab98Named)
			ab98Equal(ab98T, ab98Bit.Ab98Default, false)
			if ab98Named {
				ab98Equal(ab98T, ab98Bit.Ab98Labels, ab98Labels)
			} else {
				ab98Equal(ab98T, len(ab98Bit.Ab98Labels), 0)
			}
			ab98Raw := uint8(1) << ab98Index
			ab98Equal(ab98T, ab98Raw&(1<<ab98Bit.Ab98Bit), ab98Raw)
		}
	}
	for ab98Index, ab98Field := range []string{"DSN", "start", "offset"} {
		for _, ab98Direction := range []string{"request", "positive"} {
			ab98Entry := ab98Get(ab98T, 1, ab98Direction, ab98Field)
			ab98Equal(ab98T, ab98Entry.Ab98Position, ab98Pointer(uint16(5+2*ab98Index)))
			ab98Equal(ab98T, ab98Entry.Ab98WidthBits, (*uint16)(nil))
			ab98Equal(ab98T, ab98Entry.Ab98ByteOrder, (*string)(nil))
			if ab98Direction == "positive" {
				ab98Equal(ab98T, ab98Entry.Ab98Conditions, []Ab98Condition{{Ab98Field: "flags", Ab98Bit: 3, Ab98Set: true}})
				for _, ab98Raw := range []uint8{0, 8, 255} {
					ab98Equal(ab98T, ab98Applies(ab98Entry.Ab98Conditions, ab98Raw), ab98Raw&8 != 0)
				}
			} else {
				ab98Equal(ab98T, len(ab98Entry.Ab98Conditions), 0)
			}
		}
	}
	ab98CRC := ab98Get(ab98T, 1, "request", "CRC")
	ab98Equal(ab98T, ab98CRC.Ab98Default, (*uint64)(nil))
	ab98Equal(ab98T, ab98Get(ab98T, 1, "positive", "family-code").Ab98Default, (*uint64)(nil))
	ab98Equal(ab98T, ab98CRC.Ab98Position, ab98Pointer(uint16(10)))
	ab98Equal(ab98T, ab98CRC.Ab98WidthBits, ab98Pointer(uint16(16)))
	ab98Equal(ab98T, ab98CRC.Ab98ByteOrder, ab98Pointer("little-endian"))
	ab98Equal(ab98T, ab98CRC.Ab98Interpretations, []string{"Mirror candidate. Unsupported CRC-support flag is 0."})
	ab98Equal(ab98T, ab98CRC.Ab98Unknown, []string{"replacement number", "physical length", "qualification", "invocation"})
	ab98Equal(ab98T, ab98Get(ab98T, 1, "positive", "family-code").Ab98Position, ab98Pointer(uint16(10)))
}

func ab98Applies(ab98Conditions []Ab98Condition, ab98Raw uint8) bool {
	for _, ab98Condition := range ab98Conditions {
		if (ab98Raw&(1<<ab98Condition.Ab98Bit) != 0) != ab98Condition.Ab98Set {
			return false
		}
	}
	return true
}

func TestAb98DiagnosticAlternatives(ab98T *ab98testing.T) {
	ab98Enable := ab98Get(ab98T, 4, "positive", "enable-bits")
	ab98Equal(ab98T, ab98Enable.Ab98Position, ab98Pointer(uint16(4)))
	ab98Equal(ab98T, ab98Enable.Ab98WidthBits, ab98Pointer(uint16(8)))
	ab98Equal(ab98T, ab98Enable.Ab98Count, ab98Pointer(uint16(13)))
	ab98Equal(ab98T, ab98Enable.Ab98LogicalFootprint, ab98Pointer(uint16(17)))
	ab98Equal(ab98T, ab98Enable.Ab98Tail, []uint16{0, 13, 0})
	ab98Equal(ab98T, int(*ab98Enable.Ab98Count)*int(*ab98Enable.Ab98WidthBits), 104)
	for ab98Index := 0; ab98Index < 104; ab98Index++ {
		ab98Raw := make([]byte, 17)
		ab98Raw[4+ab98Index/8] = 1 << (ab98Index % 8)
		ab98Seen := 0
		for ab98Byte := 0; ab98Byte < int(*ab98Enable.Ab98Count); ab98Byte++ {
			for ab98Bit := 0; ab98Bit < int(*ab98Enable.Ab98WidthBits); ab98Bit++ {
				if ab98Raw[int(*ab98Enable.Ab98Position)+ab98Byte]&(1<<ab98Bit) != 0 {
					ab98Equal(ab98T, ab98Byte*8+ab98Bit, ab98Index)
					ab98Seen++
				}
			}
		}
		ab98Equal(ab98T, ab98Seen, 1)
	}
	ab98Equal(ab98T, ab98Enable.Ab98Values, map[uint64]string{0: "disabled", 1: "enabled"})
	ab98Equal(ab98T, ab98Enable.Ab98Default, ab98Pointer(uint64(0)))
	for _, ab98Op := range []uint8{4, 5} {
		ab98Level := ab98Get(ab98T, ab98Op, "request", "level")
		ab98Equal(ab98T, ab98Level.Ab98Position, ab98Pointer(uint16(4)))
		if ab98Op == 4 {
			ab98Equal(ab98T, ab98Level.Ab98WidthBits, ab98Pointer(uint16(8)))
		} else {
			ab98Equal(ab98T, ab98Level.Ab98WidthBits, (*uint16)(nil))
			ab98Equal(ab98T, ab98Level.Ab98Domain, []uint64(nil))
		}
		ab98Equal(ab98T, ab98Level.Ab98Default, ab98Pointer(uint64(0)))
	}
	ab98Point := ab98Get(ab98T, 5, "request", "point")
	ab98Equal(ab98T, ab98Point.Ab98WidthBits, (*uint16)(nil))
	ab98Equal(ab98T, ab98Point.Ab98Position, ab98Pointer(uint16(5)))
	ab98Equal(ab98T, ab98Point.Ab98Domain, []uint64{0, 99})
	ab98Following := ab98Get(ab98T, 5, "positive", "following-data")
	ab98Equal(ab98T, ab98Following.Ab98Conditions, []Ab98Condition{{Ab98Field: "flags", Ab98Bit: 0, Ab98Set: false}})
	ab98Equal(ab98T, len(ab98Following.Ab98Alternatives), 3)
	for ab98Raw := 0; ab98Raw <= 255; ab98Raw++ {
		ab98Equal(ab98T, ab98Applies(ab98Following.Ab98Conditions, uint8(ab98Raw)), ab98Raw&1 == 0)
	}
	ab98Range := ab98Following.Ab98Alternatives[0]
	ab98Equal(ab98T, ab98Range.Ab98Label, "range")
	ab98Equal(ab98T, len(ab98Range.Ab98Fields), 5)
	ab98Equal(ab98T, ab98Range.Ab98Interpretations, []string{"Candidate min+n*step with n*step<=max; enforcement Unknown."})
	for ab98Index, ab98Want := range []struct {
		ab98Name     string
		ab98Position uint16
		ab98Width    *uint16
		ab98Order    *string
	}{
		{"display-type", 5, nil, nil}, {"value", 6, ab98Pointer(uint16(16)), ab98Pointer("little-endian")}, {"minimum", 8, nil, nil}, {"maximum", 10, nil, nil}, {"step", 12, ab98Pointer(uint16(8)), nil},
	} {
		ab98Got := ab98Range.Ab98Fields[ab98Index]
		ab98Equal(ab98T, ab98Got.Ab98Default, ab98Pointer(uint64(0)))
		ab98Equal(ab98T, ab98Got.Ab98Key.Ab98Field, ab98Want.ab98Name)
		ab98Equal(ab98T, ab98Got.Ab98Position, ab98Pointer(ab98Want.ab98Position))
		ab98Equal(ab98T, ab98Got.Ab98WidthBits, ab98Want.ab98Width)
		ab98Equal(ab98T, ab98Got.Ab98ByteOrder, ab98Want.ab98Order)
	}
	ab98EnumRaw := []byte{0, 0, 0, 0, 4, 0, 11, 23, 37, 41, 59, 67}
	ab98Enum := ab98Following.Ab98Alternatives[1]
	ab98Equal(ab98T, ab98Enum.Ab98Label, "enum")
	ab98Equal(ab98T, len(ab98Enum.Ab98Fields), 6)
	ab98Equal(ab98T, ab98Enum.Ab98Interpretations, []string{"Six independent alternatives; repeat last when fewer. No separate current-value field."})
	for ab98Index, ab98Name := range []string{"enum-alternative-1", "enum-alternative-2", "enum-alternative-3", "enum-alternative-4", "enum-alternative-5", "enum-alternative-6"} {
		ab98Field := ab98Enum.Ab98Fields[ab98Index]
		ab98Equal(ab98T, ab98EnumRaw[*ab98Field.Ab98Position], []byte{11, 23, 37, 41, 59, 67}[ab98Index])
		ab98Equal(ab98T, ab98Field.Ab98Key.Ab98Field, ab98Name)
		ab98Equal(ab98T, ab98Field.Ab98Position, ab98Pointer(uint16(6+ab98Index)))
		ab98Equal(ab98T, ab98Field.Ab98WidthBits, ab98Pointer(uint16(8)))
	}
	ab98Four := ab98Following.Ab98Alternatives[2]
	ab98Equal(ab98T, ab98Four.Ab98Label, "four-byte")
	ab98Equal(ab98T, len(ab98Four.Ab98Fields), 4)
	ab98Equal(ab98T, ab98Four.Ab98Fields[0].Ab98Key.Ab98Field, "additional-value")
	ab98Equal(ab98T, ab98Four.Ab98Fields[0].Ab98Position, ab98Pointer(uint16(6)))
	ab98Equal(ab98T, ab98Four.Ab98Fields[0].Ab98WidthBits, ab98Pointer(uint16(32)))
	ab98Equal(ab98T, ab98Four.Ab98Fields[0].Ab98ByteOrder, ab98Pointer("little-endian"))
	for ab98Index, ab98Name := range []string{"reserved-10", "reserved-11", "reserved-12"} {
		ab98Equal(ab98T, ab98Four.Ab98Fields[ab98Index+1].Ab98Key.Ab98Field, ab98Name)
		ab98Equal(ab98T, ab98Four.Ab98Fields[ab98Index+1].Ab98Position, ab98Pointer(uint16(10+ab98Index)))
		ab98Equal(ab98T, ab98Four.Ab98Fields[ab98Index+1].Ab98WidthBits, ab98Pointer(uint16(8)))
	}
	ab98Equal(ab98T, ab98Range.Ab98Conditions, ab98Following.Ab98Conditions)
	ab98Equal(ab98T, ab98Enum.Ab98Conditions, ab98Following.Ab98Conditions)
	ab98Equal(ab98T, ab98Four.Ab98Conditions, []Ab98Condition{{Ab98Field: "flags", Ab98Bit: 0, Ab98Set: false}, {Ab98Field: "flags", Ab98Bit: 3, Ab98Set: true}})
	ab98Equal(ab98T, ab98Applies(ab98Range.Ab98Conditions, 8), true)
	ab98Equal(ab98T, ab98Applies(ab98Four.Ab98Conditions, 8), true)
	for _, ab98Alternative := range ab98Following.Ab98Alternatives {
		ab98Equal(ab98T, ab98Applies(ab98Alternative.Ab98Conditions, 255), false)
	}
	ab98Equal(ab98T, ab98Following.Ab98Unknown, []string{"alternative selection", "step enforcement", "physical length", "qualification", "invocation"})
}

func TestAb98ReplicationAndVersionDirections(ab98T *ab98testing.T) {
	ab98Status := ab98Get(ab98T, 7, "positive", "replication-status")
	ab98Equal(ab98T, ab98Status.Ab98Position, ab98Pointer(uint16(4)))
	ab98Equal(ab98T, ab98Status.Ab98WidthBits, ab98Pointer(uint16(8)))
	ab98Equal(ab98T, ab98Status.Ab98Default, ab98Pointer(uint64(0)))
	ab98Equal(ab98T, ab98Status.Ab98Values, map[uint64]string{0: "false", 1: "true"})
	ab98Equal(ab98T, ab98Status.Ab98Interpretations, []string{"Whole byte 0 false / 1 true.", "Nonzero-truth competes with unlabeled 2..255; raw 02 remains unresolved."})
	ab98Equal(ab98T, ab98Status.Ab98Unknown, []string{"truth for 2..255", "replication algorithms", "physical length", "qualification", "invocation"})
	ab98Equal(ab98T, ab98Get(ab98T, 7, "request", "flags").Ab98Unknown, []string{"replication algorithms", "physical length", "qualification", "invocation"})
	ab98Configuration := ab98Get(ab98T, 8, "request", "configuration")
	ab98Equal(ab98T, ab98Configuration.Ab98Default, ab98Pointer(uint64(0)))
	ab98Equal(ab98T, ab98Configuration.Ab98Values, map[uint64]string{1: "initiator-spare", 2: "responder-spare"})
	for _, ab98Case := range []struct {
		ab98Direction, ab98Field string
		ab98Position, ab98Width  uint16
	}{
		{"request", "responder-written", 4, 16}, {"request", "initiator-to-write", 6, 16}, {"request", "configuration", 8, 8}, {"request", "responder-version", 9, 8}, {"request", "initiator-version", 10, 8},
		{"positive", "initiator-written", 4, 0}, {"positive", "responder-to-write", 6, 0}, {"positive", "responder-version", 8, 8}, {"positive", "initiator-version", 9, 8},
	} {
		ab98Entry := ab98Get(ab98T, 8, ab98Case.ab98Direction, ab98Case.ab98Field)
		ab98Equal(ab98T, ab98Entry.Ab98Position, ab98Pointer(ab98Case.ab98Position))
		if ab98Case.ab98Width == 0 {
			ab98Equal(ab98T, ab98Entry.Ab98WidthBits, (*uint16)(nil))
			ab98Equal(ab98T, ab98Entry.Ab98ByteOrder, (*string)(nil))
			ab98Equal(ab98T, ab98Entry.Ab98Domain, []uint64(nil))
		} else {
			ab98Equal(ab98T, ab98Entry.Ab98WidthBits, ab98Pointer(ab98Case.ab98Width))
			ab98Equal(ab98T, ab98Entry.Ab98Domain, []uint64{0, (1 << ab98Case.ab98Width) - 1})
		}
		ab98Equal(ab98T, ab98Entry.Ab98Default, ab98Pointer(uint64(0)))
		if ab98Case.ab98Width == 16 {
			ab98Equal(ab98T, ab98Entry.Ab98ByteOrder, ab98Pointer("little-endian"))
		}
	}
	ab98Request := []byte{0, 0, 0, 0, 0, 0, 0, 0, 2, 41, 53}
	ab98Positive := []byte{0, 0, 0, 0, 0, 0, 0, 0, 67, 79}
	ab98Equal(ab98T, ab98Request[*ab98Get(ab98T, 8, "request", "responder-version").Ab98Position], byte(41))
	ab98Equal(ab98T, ab98Request[*ab98Get(ab98T, 8, "request", "initiator-version").Ab98Position], byte(53))
	ab98Equal(ab98T, ab98Positive[*ab98Get(ab98T, 8, "positive", "responder-version").Ab98Position], byte(67))
	ab98Equal(ab98T, ab98Positive[*ab98Get(ab98T, 8, "positive", "initiator-version").Ab98Position], byte(79))
}

func TestAb98IndependentStorage(ab98T *ab98testing.T) {
	ab98Before, ab98Err := Ab98Entries()
	if ab98Err != nil {
		ab98T.Fatal(ab98Err)
	}
	ab98Mutated, ab98Err := Ab98Entries()
	if ab98Err != nil {
		ab98T.Fatal(ab98Err)
	}
	for ab98Index := range ab98Mutated {
		ab98Entry := &ab98Mutated[ab98Index]
		*ab98Entry.Ab98Key.Ab98Secondary = 0
		*ab98Entry.Ab98Key.Ab98Operation = 0
		if ab98Entry.Ab98Position != nil {
			*ab98Entry.Ab98Position = 0
		}
		if len(ab98Entry.Ab98Domain) > 0 {
			ab98Entry.Ab98Domain[0] = 99
		}
		if len(ab98Entry.Ab98Bits) > 0 {
			ab98Entry.Ab98Bits[0].Ab98Labels[0] = "mutated"
		}
		if len(ab98Entry.Ab98Conditions) > 0 {
			ab98Entry.Ab98Conditions[0].Ab98Bit = 7
		}
		if len(ab98Entry.Ab98Alternatives) > 0 {
			ab98Entry.Ab98Alternatives[0].Ab98Fields[0].Ab98Key.Ab98Field = "mutated"
			*ab98Entry.Ab98Alternatives[0].Ab98Fields[0].Ab98Position = 99
			ab98Entry.Ab98Alternatives[0].Ab98Interpretations[0] = "mutated"
		}
		if len(ab98Entry.Ab98Values) > 0 {
			ab98Entry.Ab98Values[0] = "mutated"
		}
		ab98Entry.Ab98Unknown[0] = "mutated"
	}
	ab98After, ab98Err := Ab98Entries()
	if ab98Err != nil {
		ab98T.Fatal(ab98Err)
	}
	ab98Equal(ab98T, ab98After, ab98Before)
	ab98One := ab98Get(ab98T, 5, "positive", "following-data")
	ab98Two := ab98Get(ab98T, 5, "positive", "following-data")
	*ab98One.Ab98Alternatives[0].Ab98Fields[0].Ab98Position = 99
	ab98One.Ab98Alternatives[0].Ab98Conditions[0].Ab98Bit = 7
	ab98Three := ab98Get(ab98T, 5, "positive", "following-data")
	ab98Equal(ab98T, ab98Two, ab98Three)
	ab98Flags := ab98Get(ab98T, 1, "request", "flags")
	ab98Flags.Ab98Bits[0].Ab98Labels[0] = "mutated"
	ab98Equal(ab98T, ab98Get(ab98T, 1, "request", "flags").Ab98Bits[0].Ab98Labels[0], "power-on/start")
}

func TestAb98MalformedData(ab98T *ab98testing.T) {
	ab98Original := ab98Data
	defer func() { ab98Data = ab98Original }()
	ab98Entry := ab98Get(ab98T, 1, "request", "sequence")
	ab98Duplicate, ab98Err := ab98json.Marshal([]Ab98Entry{ab98Entry, ab98Entry})
	if ab98Err != nil {
		ab98T.Fatal(ab98Err)
	}
	for _, ab98Input := range []string{"{", "null", "[]", "[{}]", string(ab98Original) + " []", string(ab98Duplicate), ab98strings.Replace(string(ab98Original), `"Hypothesis"`, `"qualified"`, 1), ab98strings.Replace(string(ab98Original), `"position": 3`, `"position": -1`, 1), ab98strings.Replace(string(ab98Original), `"grade":`, `"unexpected":`, 1),
		ab98strings.Replace(string(ab98Original), `"bit": 0`, `"bit": 8`, 1),
		ab98strings.Replace(string(ab98Original), `"width_bits": 8`, `"width_bits": 0`, 1),
		ab98strings.Replace(string(ab98Original), `"little-endian"`, `"guessed-endian"`, 1)} {
		ab98Data = []byte(ab98Input)
		ab98Entries, ab98Err := Ab98Entries()
		if ab98Err == nil || ab98Entries != nil {
			ab98T.Fatalf("malformed catalog accepted")
		}
		ab98Got, ab98Found, ab98Err := Ab98Lookup(ab98Entry.Ab98Key)
		if ab98Err == nil || ab98Found || !ab98reflect.DeepEqual(ab98Got, Ab98Entry{}) {
			ab98T.Fatal("malformed lookup succeeded")
		}
	}
}
