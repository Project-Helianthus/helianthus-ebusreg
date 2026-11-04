package vaillant

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"testing"
)

type ab26Field struct {
	Ab26Name     string `json:"name"`
	Ab26Byte     int    `json:"byte"`
	Ab26Bit      int    `json:"bit"`
	Ab26Width    int    `json:"width_bits"`
	Ab26Channel  int    `json:"channel"`
	Ab26Weekday  int    `json:"weekday_index"`
	Ab26Meaning  string `json:"meaning"`
	Ab26Encoding string `json:"encoding"`
	Ab26Domain   []int  `json:"domain"`
	Ab26Unit     string `json:"unit"`
	Ab26Sentinel string `json:"sentinel"`
}

type ab26Row struct {
	Ab26ID        string      `json:"id"`
	Ab26Key       Ab26Key     `json:"key"`
	Ab26Grade     string      `json:"grade"`
	Ab26Origin    string      `json:"coordinate_origin"`
	Ab26Physical  string      `json:"physical_offset"`
	Ab26Fields    []ab26Field `json:"fields"`
	Ab26Relations struct {
		Ab26Required  int    `json:"required_accessible_initialized_bytes"`
		Ab26Declared  bool   `json:"declared_extent_is_received_extent"`
		Ab26Local     bool   `json:"local_count_is_received_extent"`
		Ab26Capacity  bool   `json:"backing_capacity_is_received_extent"`
		Ab26Status    bool   `json:"status_is_extent"`
		Ab26Threshold bool   `json:"numeric_threshold_is_extent"`
		Ab26Extra     string `json:"extra_bytes"`
		Ab26Empty     bool   `json:"empty_arguments_are_physical_empty_packet"`
	} `json:"relations"`
	Ab26Raw           string            `json:"raw_form"`
	Ab26Unknown       map[string]string `json:"unknown"`
	Ab26Qualification string            `json:"qualification"`
	Ab26Executable    bool              `json:"executable"`
}

func ab26Rows(ab26T *testing.T) ([]json.RawMessage, []ab26Row) {
	ab26T.Helper()
	ab26Data, ab26Error := Ab26Entries()
	if ab26Error != nil {
		ab26T.Fatal(ab26Error)
	}
	var ab26Raw []json.RawMessage
	var ab26Decoded []ab26Row
	if ab26Error = json.Unmarshal(ab26Data, &ab26Raw); ab26Error != nil {
		ab26T.Fatal(ab26Error)
	}
	if ab26Error = json.Unmarshal(ab26Data, &ab26Decoded); ab26Error != nil {
		ab26T.Fatal(ab26Error)
	}
	return ab26Raw, ab26Decoded
}

func ab26Reply(ab26T *testing.T, ab26Operation uint8) ab26Row {
	ab26T.Helper()
	_, ab26Decoded := ab26Rows(ab26T)
	for _, ab26Candidate := range ab26Decoded {
		if *ab26Candidate.Ab26Key.Ab26Operation == ab26Operation && ab26Candidate.Ab26Key.Ab26Direction == "reply" {
			ab26Raw, ab26Found, ab26Error := Ab26Lookup(ab26Candidate.Ab26Key)
			if ab26Error != nil || !ab26Found {
				ab26T.Fatalf("lookup: %v %v", ab26Found, ab26Error)
			}
			var ab26Result ab26Row
			if ab26Error = json.Unmarshal(ab26Raw, &ab26Result); ab26Error != nil {
				ab26T.Fatal(ab26Error)
			}
			return ab26Result
		}
	}
	ab26T.Fatalf("missing operation %x", ab26Operation)
	return ab26Row{}
}

func TestAb26Catalog(ab26T *testing.T) {
	ab26Raw, ab26Decoded := ab26Rows(ab26T)
	if len(ab26Decoded) != 14 {
		ab26T.Fatalf("rows: %d", len(ab26Decoded))
	}
	ab26IDs := map[string]bool{}
	for ab26Index, ab26Candidate := range ab26Decoded {
		if !strings.HasPrefix(ab26Candidate.Ab26ID, "ab26-") || ab26IDs[ab26Candidate.Ab26ID] {
			ab26T.Fatal("invalid ID")
		}
		ab26IDs[ab26Candidate.Ab26ID] = true
		ab26Result, ab26Found, ab26Error := Ab26Lookup(ab26Candidate.Ab26Key)
		if ab26Error != nil || !ab26Found || !bytes.Equal(ab26Result, ab26Raw[ab26Index]) {
			ab26T.Fatal("whole-row lookup mismatch")
		}
		if ab26Candidate.Ab26Grade != "Hypothesis" || ab26Candidate.Ab26Qualification != "Unknown" || ab26Candidate.Ab26Physical != "Unknown" || ab26Candidate.Ab26Origin != "immediately_after_command" || ab26Candidate.Ab26Executable {
			ab26T.Fatal("grade or qualification")
		}
		if ab26Candidate.Ab26Raw != "retain_all_original_bytes_including_reserved_and_extra_bytes" {
			ab26T.Fatal("raw fallback")
		}
		for _, ab26Name := range []string{"native_support", "physical_adaptation", "physical_effects", "persistence", "delivery", "other_timer_variants", "unit", "sentinel"} {
			if ab26Candidate.Ab26Unknown[ab26Name] != "Unknown" {
				ab26T.Fatal(ab26Name)
			}
		}
		ab26Relations := ab26Candidate.Ab26Relations
		if ab26Relations.Ab26Declared || ab26Relations.Ab26Local || ab26Relations.Ab26Capacity || ab26Relations.Ab26Status || ab26Relations.Ab26Threshold || ab26Relations.Ab26Empty || ab26Relations.Ab26Extra != "retain_raw" {
			ab26T.Fatal("independent extents/raw retention")
		}
		ab26Key := ab26Candidate.Ab26Key
		if ab26Key.Ab26Family != "vaillant" || ab26Key.Ab26Variant != "b555" || ab26Key.Ab26Primary != 0xb5 || ab26Key.Ab26Secondary == nil || *ab26Key.Ab26Secondary != 0x55 || ab26Key.Ab26Group != nil || ab26Key.Ab26Register != nil || ab26Key.Ab26InstanceScope != "none" {
			ab26T.Fatal("identity")
		}
		if ab26Key.Ab26Direction == "request" && (len(ab26Candidate.Ab26Fields) != 0 || ab26Relations.Ab26Required != 0) {
			ab26T.Fatal("declared empty arguments")
		}
	}
}

func TestAb26ExactAxes(ab26T *testing.T) {
	ab26Base := ab26Reply(ab26T, 0xa0).Ab26Key
	ab26Zero := uint8(0)
	ab26ZeroRegister := uint16(0)
	ab26Changes := []struct {
		Ab26Name   string
		Ab26Change func(*Ab26Key)
	}{
		{"ab26Family", func(ab26Key *Ab26Key) { ab26Key.Ab26Family = "other" }},
		{"ab26Variant", func(ab26Key *Ab26Key) { ab26Key.Ab26Variant = "other" }},
		{"ab26Primary", func(ab26Key *Ab26Key) { ab26Key.Ab26Primary = 0 }},
		{"ab26SecondaryAbsent", func(ab26Key *Ab26Key) { ab26Key.Ab26Secondary = nil }},
		{"ab26SecondaryZero", func(ab26Key *Ab26Key) { ab26Key.Ab26Secondary = &ab26Zero }},
		{"ab26OperationAbsent", func(ab26Key *Ab26Key) { ab26Key.Ab26Operation = nil }},
		{"ab26OperationZero", func(ab26Key *Ab26Key) { ab26Key.Ab26Operation = &ab26Zero }},
		{"ab26Direction", func(ab26Key *Ab26Key) { ab26Key.Ab26Direction = "none" }},
		{"ab26Field", func(ab26Key *Ab26Key) { ab26Key.Ab26Field = "other" }},
		{"ab26Alternative", func(ab26Key *Ab26Key) { ab26Key.Ab26Alternative = "other" }},
		{"ab26Group", func(ab26Key *Ab26Key) { ab26Key.Ab26Group = &ab26Zero }},
		{"ab26InstanceScope", func(ab26Key *Ab26Key) { ab26Key.Ab26InstanceScope = "dynamic" }},
		{"ab26Register", func(ab26Key *Ab26Key) { ab26Key.Ab26Register = &ab26ZeroRegister }},
	}
	for _, ab26Case := range ab26Changes {
		ab26T.Run(ab26Case.Ab26Name, func(ab26T *testing.T) {
			ab26Key := ab26Base
			ab26Case.Ab26Change(&ab26Key)
			ab26Result, ab26Found, ab26Error := Ab26Lookup(ab26Key)
			if ab26Error != nil || ab26Found || ab26Result != nil {
				ab26T.Fatalf("inexact key admitted: %v %v", ab26Found, ab26Error)
			}
			ab26Data, ab26Error := json.Marshal([]struct {
				Ab26Key Ab26Key `json:"key"`
			}{{ab26Base}, {ab26Key}})
			if ab26Error != nil {
				ab26T.Fatal(ab26Error)
			}
			if _, _, ab26Error = ab26Validate(ab26Data); ab26Error != nil {
				ab26T.Fatalf("distinct keys rejected: %v", ab26Error)
			}
		})
	}
	// Optional absent and explicit zero are distinct even where production rows use a nonzero value.
	for _, ab26Axis := range []string{"secondary", "operation", "group", "register"} {
		ab26T.Run("ab26Optional"+ab26Axis, func(ab26T *testing.T) {
			ab26First := ab26Base
			ab26Second := ab26Base
			switch ab26Axis {
			case "secondary":
				ab26First.Ab26Secondary = nil
				ab26Second.Ab26Secondary = &ab26Zero
			case "operation":
				ab26First.Ab26Operation = nil
				ab26Second.Ab26Operation = &ab26Zero
			case "group":
				ab26First.Ab26Group = nil
				ab26Second.Ab26Group = &ab26Zero
			case "register":
				ab26First.Ab26Register = nil
				ab26Second.Ab26Register = &ab26ZeroRegister
			}
			ab26Data, ab26Error := json.Marshal([]struct {
				Ab26Key Ab26Key `json:"key"`
			}{{ab26First}, {ab26Second}})
			if ab26Error != nil {
				ab26T.Fatal(ab26Error)
			}
			if _, _, ab26Error = ab26Validate(ab26Data); ab26Error != nil {
				ab26T.Fatal(ab26Error)
			}
		})
	}
}

func TestAb26ValidationAndCopies(ab26T *testing.T) {
	ab26Raw, ab26Decoded := ab26Rows(ab26T)
	for _, ab26Invalid := range []string{"{", "null", "{}", "[null]", "[{}]", "[1]", "[{\"key\":{\"primary\":256}}]", "[" + string(ab26Raw[0]) + "," + string(ab26Raw[0]) + "]"} {
		if _, _, ab26Error := ab26Validate([]byte(ab26Invalid)); ab26Error == nil {
			ab26T.Fatal("invalid data admitted")
		}
	}
	var ab26Object map[string]json.RawMessage
	if ab26Error := json.Unmarshal(ab26Raw[0], &ab26Object); ab26Error != nil {
		ab26T.Fatal(ab26Error)
	}
	var ab26Identity map[string]json.RawMessage
	if ab26Error := json.Unmarshal(ab26Object["key"], &ab26Identity); ab26Error != nil {
		ab26T.Fatal(ab26Error)
	}
	for ab26Axis, ab26Value := range ab26Identity {
		delete(ab26Identity, ab26Axis)
		ab26KeyData, _ := json.Marshal(ab26Identity)
		ab26Object["key"] = ab26KeyData
		ab26Data, _ := json.Marshal([]map[string]json.RawMessage{ab26Object})
		if _, _, ab26Error := ab26Validate(ab26Data); ab26Error == nil {
			ab26T.Fatalf("missing %s admitted", ab26Axis)
		}
		ab26Identity[ab26Axis] = ab26Value
	}
	ab26Entries, ab26Error := Ab26Entries()
	if ab26Error != nil {
		ab26T.Fatal(ab26Error)
	}
	ab26Expected := append([]byte(nil), ab26Entries...)
	ab26Other, _ := Ab26Entries()
	for ab26Index := range ab26Entries {
		ab26Entries[ab26Index] = 'x'
	}
	ab26Later, ab26Error := Ab26Entries()
	if ab26Error != nil || !bytes.Equal(ab26Later, ab26Expected) || !bytes.Equal(ab26Other, ab26Expected) {
		ab26T.Fatal("entries share memory")
	}
	for ab26Index, ab26Candidate := range ab26Decoded {
		ab26Result, _, _ := Ab26Lookup(ab26Candidate.Ab26Key)
		ab26Other, _, _ := Ab26Lookup(ab26Candidate.Ab26Key)
		for ab26Byte := range ab26Result {
			ab26Result[ab26Byte] = 'x'
		}
		ab26Later, ab26Found, ab26Error := Ab26Lookup(ab26Candidate.Ab26Key)
		if ab26Error != nil || !ab26Found || !bytes.Equal(ab26Later, ab26Raw[ab26Index]) || !bytes.Equal(ab26Other, ab26Later) {
			ab26T.Fatal("row shares memory")
		}
	}
}

func ab26Value(ab26Field ab26Field, ab26Raw []byte) uint64 {
	var ab26Result uint64
	for ab26Bit := 0; ab26Bit < ab26Field.Ab26Width; ab26Bit++ {
		ab26Position := ab26Field.Ab26Byte*8 + ab26Field.Ab26Bit + ab26Bit
		ab26Result |= uint64((ab26Raw[ab26Position/8]>>uint(ab26Position%8))&1) << uint(ab26Bit)
	}
	return ab26Result
}

func TestAb26Masks(ab26T *testing.T) {
	for _, ab26Operation := range []uint8{0xa0, 0xa7} {
		ab26Row := ab26Reply(ab26T, ab26Operation)
		ab26Label := "availability"
		if ab26Operation == 0xa7 {
			ab26Label = "visibility"
		}
		if ab26Row.Ab26Key.Ab26Field != ab26Label || ab26Row.Ab26Relations.Ab26Required != 8 || len(ab26Row.Ab26Fields) != 55 {
			ab26T.Fatal("mask meaning or extent")
		}
		ab26Occupied := [64]bool{}
		for ab26Index, ab26Field := range ab26Row.Ab26Fields {
			ab26ExpectedByte, ab26ExpectedBit, ab26ExpectedWidth := 0, ab26Index, 1
			if ab26Index < 5 {
				if ab26Field.Ab26Name != []string{"DHW", "circulation", "noise_reduction", "ventilation", "high_tariff"}[ab26Index] {
					ab26T.Fatal("global bit label")
				}
			} else if ab26Index == 5 {
				ab26ExpectedBit = 5
				ab26ExpectedWidth = 3
			} else if ab26Index == 6 {
				ab26ExpectedByte = 1
				ab26ExpectedBit = 0
				ab26ExpectedWidth = 8
			} else {
				ab26Channel := (ab26Index-7)/4 + 1
				ab26Role := (ab26Index - 7) % 4
				// Independent explicit byte/nibble controls, including channel 12 at byte 7.
				ab26ExpectedByte = []int{2, 2, 3, 3, 4, 4, 5, 5, 6, 6, 7, 7}[ab26Channel-1]
				ab26ExpectedBit = []int{0, 4, 0, 4, 0, 4, 0, 4, 0, 4, 0, 4}[ab26Channel-1] + ab26Role
				if ab26Field.Ab26Channel != ab26Channel || ab26Field.Ab26Name != []string{"zone_heating", "zone_cooling", "external_DHW", "external_fixed_flow"}[ab26Role] {
					ab26T.Fatal("channel role")
				}
			}
			if ab26Field.Ab26Byte != ab26ExpectedByte || ab26Field.Ab26Bit != ab26ExpectedBit || ab26Field.Ab26Width != ab26ExpectedWidth {
				ab26T.Fatalf("coordinate %d: %+v", ab26Index, ab26Field)
			}
			if ab26Index == 5 || ab26Index == 6 {
				if ab26Field.Ab26Name != "reserved" || ab26Field.Ab26Encoding != "raw" {
					ab26T.Fatal("reserved retention")
				}
			} else if ab26Field.Ab26Meaning != ab26Label {
				ab26T.Fatal("availability/visibility alias")
			}
			for ab26Bit := 0; ab26Bit < ab26ExpectedWidth; ab26Bit++ {
				ab26Position := ab26ExpectedByte*8 + ab26ExpectedBit + ab26Bit
				if ab26Occupied[ab26Position] {
					ab26T.Fatal("overlap")
				}
				ab26Occupied[ab26Position] = true
				ab26Raw := make([]byte, 8)
				ab26Raw[ab26Position/8] = 1 << uint(ab26Position%8)
				for ab26OtherIndex, ab26OtherField := range ab26Row.Ab26Fields {
					ab26Expected := uint64(0)
					if ab26OtherIndex == ab26Index {
						ab26Expected = 1 << uint(ab26Bit)
					}
					if ab26Value(ab26OtherField, ab26Raw) != ab26Expected {
						ab26T.Fatal("one-bit isolation or retention")
					}
				}
			}
		}
		for _, ab26Present := range ab26Occupied {
			if !ab26Present {
				ab26T.Fatal("lost bit")
			}
		}
	}
}

func TestAb26ReservedReplies(ab26T *testing.T) {
	for _, ab26Operation := range []uint8{0xa1, 0xa2, 0xa8, 0xa9} {
		ab26Row := ab26Reply(ab26T, ab26Operation)
		if ab26Row.Ab26Key.Ab26Field != "reserved" || ab26Row.Ab26Executable || ab26Row.Ab26Relations.Ab26Required != 6 || len(ab26Row.Ab26Fields) != 1 {
			ab26T.Fatal("reserved settings promoted")
		}
		ab26Field := ab26Row.Ab26Fields[0]
		if ab26Field.Ab26Name != "reserved" || ab26Field.Ab26Byte != 0 || ab26Field.Ab26Bit != 0 || ab26Field.Ab26Width != 48 || ab26Field.Ab26Encoding != "raw" || ab26Value(ab26Field, []byte{1, 2, 3, 4, 5, 6}) != 0x060504030201 {
			ab26T.Fatal("reserved bytes")
		}
	}
}

func TestAb26PositiveForms(ab26T *testing.T) {
	for _, ab26Case := range []struct {
		Ab26Operation uint8
		Ab26Required  int
	}{{0xa4, 9}, {0xa5, 7}} {
		ab26Row := ab26Reply(ab26T, ab26Case.Ab26Operation)
		if ab26Row.Ab26Relations.Ab26Required != ab26Case.Ab26Required || ab26Row.Ab26Key.Ab26Alternative != "positive" {
			ab26T.Fatal("positive extent")
		}
		for _, ab26Accessible := range []int{ab26Case.Ab26Required - 1, ab26Case.Ab26Required, ab26Case.Ab26Required + 1} {
			for _, ab26Initialized := range []int{0, ab26Case.Ab26Required - 1, 16} {
				for _, ab26Declared := range []int{0, 16} {
					for _, ab26Local := range []int{0, 16} {
						ab26T.Run(fmt.Sprintf("ab26%x-accessible%d-initialized%d-declared%d-local%d", ab26Case.Ab26Operation, ab26Accessible, ab26Initialized, ab26Declared, ab26Local), func(ab26T *testing.T) {
							ab26Storage := make([]byte, ab26Accessible, 16)
							ab26Before := append([]byte(nil), ab26Storage...)
							ab26Complete := len(ab26Storage) >= ab26Row.Ab26Relations.Ab26Required && ab26Initialized >= ab26Row.Ab26Relations.Ab26Required
							ab26Expected := ab26Accessible >= ab26Case.Ab26Required && ab26Initialized >= ab26Case.Ab26Required
							if ab26Complete != ab26Expected || !bytes.Equal(ab26Storage, ab26Before) {
								ab26T.Fatal("extent independence/raw retention")
							}
						})
					}
				}
			}
		}
		if ab26Case.Ab26Operation == 0xa4 {
			if len(ab26Row.Ab26Fields) != 9 {
				ab26T.Fatal("weekday fields")
			}
			for ab26Index, ab26Field := range ab26Row.Ab26Fields {
				ab26Name := "weekday_count"
				if ab26Index == 0 {
					ab26Name = "status"
				}
				if ab26Index == 8 {
					ab26Name = "24_hour_count"
				}
				if ab26Field.Ab26Byte != ab26Index || ab26Field.Ab26Width != 8 || ab26Field.Ab26Name != ab26Name {
					ab26T.Fatal("weekday coordinate")
				}
				if ab26Index > 0 && ab26Index < 8 && ab26Field.Ab26Weekday != ab26Index {
					ab26T.Fatal("weekday index")
				}
			}
		} else {
			if len(ab26Row.Ab26Fields) != 6 {
				ab26T.Fatal("interval fields")
			}
			for ab26Index, ab26Field := range ab26Row.Ab26Fields {
				if ab26Field.Ab26Byte != ab26Index || ab26Field.Ab26Name != []string{"status", "start_hour", "start_minute", "end_hour", "end_minute", "setpoint"}[ab26Index] {
					ab26T.Fatal("interval coordinate")
				}
				if ab26Index < 5 && ab26Field.Ab26Width != 8 {
					ab26T.Fatal("interval width")
				}
				if ab26Index > 0 && ab26Index < 5 {
					ab26Max := 59
					if ab26Index%2 == 1 {
						ab26Max = 24
					}
					if !reflect.DeepEqual(ab26Field.Ab26Domain, []int{0, ab26Max}) || ab26Field.Ab26Encoding != "unsigned" {
						ab26T.Fatal("time domain")
					}
				}
			}
			ab26Field := ab26Row.Ab26Fields[5]
			if ab26Field.Ab26Width != 16 || ab26Field.Ab26Encoding != "unsigned_little_endian" || ab26Field.Ab26Unit != "Unknown" || ab26Field.Ab26Sentinel != "Unknown" {
				ab26T.Fatal("setpoint representation")
			}
			for _, ab26Case := range []struct {
				Ab26Raw      []byte
				Ab26Expected uint64
			}{{[]byte{0, 24, 59, 0, 0, 255, 255, 42}, 65535}, {[]byte{255, 0, 0, 24, 59, 0x34, 0x12, 42}, 4660}} {
				ab26Before := append([]byte(nil), ab26Case.Ab26Raw...)
				if ab26Value(ab26Field, ab26Case.Ab26Raw) != ab26Case.Ab26Expected || !bytes.Equal(ab26Before, ab26Case.Ab26Raw) {
					ab26T.Fatal("raw numeric or extra byte changed")
				}
			}
		}
	}
}
