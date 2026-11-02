package system

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
)

func TestDecodeExtRegisterResponseMarksStaticConstraintsUnscoped(t *testing.T) {
	t.Parallel()

	mapped := decodeExtRegisterResponse(0x02, B524OpcodeLocal, 0x01, 0x01, 0x0004, []byte{0x01, 0x01, 0x04, 0x00, 0x2A})
	if _, ok := mapped["constraints"]; !ok {
		t.Fatal("constraints absent for a catalog entry; want present")
	}
	if got := mapped["constraint_qualified"]; !got.Valid || got.Value != false {
		t.Fatalf("constraint_qualified = %+v; want valid false", got)
	}
	if got := mapped["constraint_scope"]; !got.Valid || got.Value != "unscoped_catalog" {
		t.Fatalf("constraint_scope = %+v; want valid unscoped_catalog", got)
	}
	// Existing keys stay.
	for _, key := range []string{"constraint_record", "constraint_record_hex", "constraint_type", "constraint_min", "constraint_max", "constraint_step"} {
		if got := mapped[key]; !got.Valid {
			t.Errorf("%s = %+v; want valid (additive change must keep existing keys)", key, got)
		}
	}

	unmapped := decodeExtRegisterResponse(0x02, B524OpcodeLocal, 0x00, 0x00, 0x5C00, []byte{0x00})
	for _, key := range []string{"constraints", "constraint_qualified", "constraint_scope"} {
		if got, ok := unmapped[key]; ok {
			t.Errorf("%s = %+v; want absent without a catalog entry", key, got)
		}
	}
}

func TestAb01ReplyCorrelation(ab01T *testing.T) {
	ab01Cases := []struct {
		ab01Name                            string
		ab01Opcode, ab01Group, ab01Instance byte
		ab01Addr                            uint16
		ab01Payload                         []byte
		ab01Value, ab01Bound                bool
	}{
		{"ab01 local control", 2, 1, 1, 4, []byte{1, 1, 4, 0, 42}, true, true},
		{"ab01 OP only remote control", 6, 1, 1, 4, []byte{1, 1, 4, 0, 42}, true, false},
		{"ab01 II only request context", 2, 1, 9, 4, []byte{1, 1, 4, 0, 42}, true, true},
		{"ab01 FLAGS only zero", 2, 1, 1, 4, []byte{0, 1, 4, 0, 42}, true, true},
		{"ab01 FLAGS only ff", 2, 1, 1, 4, []byte{255, 1, 4, 0, 42}, true, true},
		{"ab01 GG only known bound mismatch", 2, 1, 1, 4, []byte{1, 2, 4, 0, 42}, false, false},
		{"ab01 RR low only known bound mismatch", 2, 1, 1, 4, []byte{1, 1, 5, 0, 42}, false, false},
		{"ab01 RR high only mismatch", 2, 1, 1, 4, []byte{1, 1, 4, 1, 42}, false, false},
		{"ab01 high byte correlated control", 2, 1, 1, 0x0104, []byte{1, 1, 4, 1, 42}, true, false},
		{"ab01 empty", 2, 1, 1, 4, nil, false, false},
		{"ab01 status zero", 2, 1, 1, 4, []byte{0}, false, false},
		{"ab01 status ff", 2, 1, 1, 4, []byte{255}, false, false},
		{"ab01 short two", 2, 1, 1, 4, []byte{1, 1}, false, false},
		{"ab01 short three", 2, 1, 1, 4, []byte{1, 1, 4}, false, false},
		{"ab01 short four correlated", 2, 1, 1, 4, []byte{1, 1, 4, 0}, false, false},
		{"ab01 short four known mismatch", 2, 1, 1, 4, []byte{1, 3, 5, 0}, false, false},
	}
	for _, ab01Case := range ab01Cases {
		ab01T.Run(ab01Case.ab01Name, func(ab01T *testing.T) {
			ab01Input := append([]byte(nil), ab01Case.ab01Payload...)
			ab01Values := decodeExtRegisterResponse(2, ab01Case.ab01Opcode, ab01Case.ab01Group, ab01Case.ab01Instance, ab01Case.ab01Addr, ab01Case.ab01Payload)
			for ab01Key, ab01Want := range map[string]any{"opcode": ab01Case.ab01Opcode, "group": ab01Case.ab01Group, "instance": ab01Case.ab01Instance, "addr": ab01Case.ab01Addr, "addr_hex": fmt.Sprintf("%04X", ab01Case.ab01Addr)} {
				if ab01Got := ab01Values[ab01Key]; !ab01Got.Valid || ab01Got.Value != ab01Want {
					ab01T.Fatalf("%s = %+v; want %v", ab01Key, ab01Got, ab01Want)
				}
			}
			if ab01Got := ab01Values["value"]; ab01Got.Valid != ab01Case.ab01Value {
				ab01T.Fatalf("value = %+v; want valid %v", ab01Got, ab01Case.ab01Value)
			} else if ab01Case.ab01Value {
				if !bytes.Equal(ab01Got.Value.([]byte), ab01Input[4:]) {
					ab01T.Fatalf("value = %+v; want payload suffix", ab01Got)
				}
			} else if ab01Got.Value != nil {
				ab01T.Fatalf("invalid value retains usable bytes: %+v", ab01Got)
			}
			if _, ab01OK := ab01Values["constraints"]; ab01OK != ab01Case.ab01Bound {
				ab01T.Fatalf("constraints present = %v; want %v", ab01OK, ab01Case.ab01Bound)
			}
			for ab01Key := range ab01Values {
				if !ab01Case.ab01Bound && strings.HasPrefix(ab01Key, "constraint") {
					ab01T.Fatalf("unexpected constraint association %s", ab01Key)
				}
			}
			if ab01Case.ab01Bound {
				if ab01Values["constraint_scope"].Value != "unscoped_catalog" || ab01Values["constraint_qualified"].Value != false || ab01Values["constraint_record"].Value != ab01Case.ab01Addr {
					ab01T.Fatalf("unscoped markers = %+v", ab01Values)
				}
			}
			if len(ab01Input) > 0 {
				if ab01Got := ab01Values["reply_kind"]; !ab01Got.Valid || ab01Got.Value != ab01Input[0] {
					ab01T.Fatalf("FLAGS = %+v", ab01Got)
				}
			} else if ab01Values["reply_kind"].Valid {
				ab01T.Fatal("empty payload has FLAGS")
			}
			if len(ab01Input) >= 4 {
				for ab01Key, ab01Want := range map[string]any{"reply_group": ab01Input[1], "reply_addr": uint16(ab01Input[2]) | uint16(ab01Input[3])<<8, "reply_addr_hex": fmt.Sprintf("%04X", uint16(ab01Input[2])|uint16(ab01Input[3])<<8)} {
					if ab01Got := ab01Values[ab01Key]; !ab01Got.Valid || ab01Got.Value != ab01Want {
						ab01T.Fatalf("%s = %+v; want %v", ab01Key, ab01Got, ab01Want)
					}
				}
			} else {
				if ab01Values["prefix"].Valid {
					ab01T.Fatal("short reply has valid prefix")
				}
				for _, ab01Key := range []string{"reply_group", "reply_addr", "reply_addr_hex"} {
					if _, ab01OK := ab01Values[ab01Key]; ab01OK {
						ab01T.Fatalf("short reply invents %s", ab01Key)
					}
				}
			}
			if !ab01Values["payload"].Valid || !bytes.Equal(ab01Values["payload"].Value.([]byte), ab01Input) {
				ab01T.Fatal("raw payload not preserved")
			}
			for _, ab01Key := range []string{"payload", "prefix", "value"} {
				if ab01Raw, ab01OK := ab01Values[ab01Key].Value.([]byte); ab01OK && len(ab01Raw) > 0 {
					ab01Raw[0] ^= 255
					if !bytes.Equal(ab01Case.ab01Payload, ab01Input) {
						ab01T.Fatalf("%s aliases input", ab01Key)
					}
				}
			}
		})
	}
}
