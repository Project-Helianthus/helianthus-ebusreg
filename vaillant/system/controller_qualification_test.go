package system

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/Project-Helianthus/helianthus-ebusreg/registry"
)

type qualificationFixtureVector struct {
	ID              string `json:"id"`
	Class           string `json:"class"`
	RequestAddress  string `json:"request_address"`
	ResponseAddress string `json:"response_address"`
	PayloadHex      string `json:"payload_hex"`
	Group           string `json:"group"`
	Register        string `json:"register"`
	Expect          struct {
		Outcome         string `json:"outcome"`
		Reason          string `json:"reason"`
		DeviceID        string `json:"device_id"`
		ProductNumber   string `json:"product_number"`
		SoftwareVersion string `json:"software_version"`
		Level           string `json:"level"`
	} `json:"expect"`
}

type qualificationFixture struct {
	Vectors []qualificationFixtureVector `json:"vectors"`
}

func loadQualificationFixture(t *testing.T) qualificationFixture {
	t.Helper()
	var fixture qualificationFixture
	if err := json.Unmarshal(controllerQualificationRuleJSON, &fixture); err != nil {
		t.Fatalf("decode rule vectors: %v", err)
	}
	if len(fixture.Vectors) == 0 {
		t.Fatal("rule carries no vectors")
	}
	return fixture
}

func mustHexBytes(t *testing.T, s string) []byte {
	t.Helper()
	b, err := hex.DecodeString(s)
	if err != nil {
		t.Fatalf("payload %q: %v", s, err)
	}
	return b
}

func parseVectorAddress(t *testing.T, s string) byte {
	t.Helper()
	if s == "" {
		return 0x15
	}
	b, ok := parseHexByte(trimHexPrefix(s))
	if !ok {
		t.Fatalf("address %q", s)
	}
	return b
}

func trimHexPrefix(s string) string {
	if len(s) > 2 && s[:2] == "0x" {
		return s[2:]
	}
	return s
}

// The rule digest is computed from the record bytes, not written as a literal.
func TestControllerQualificationRuleDigestMatchesRecord(t *testing.T) {
	data, err := os.ReadFile("controller_qualification_rule.json")
	if err != nil {
		t.Fatal(err)
	}
	sum := sha256.Sum256(data)
	if got, want := ControllerQualificationRuleDigest(), hex.EncodeToString(sum[:]); got != want {
		t.Fatalf("rule digest = %s; want %s", got, want)
	}
}

// TestControllerQualificationVectors replays every controller vector of the
// record, through the payload classifier.
func TestControllerQualificationVectors(t *testing.T) {
	fixture := loadQualificationFixture(t)
	checked := 0
	for _, v := range fixture.Vectors {
		if v.Class != "controller" {
			continue
		}
		checked++
		t.Run(v.ID, func(t *testing.T) {
			got := QualifyControllerIdentification(parseVectorAddress(t, v.RequestAddress), parseVectorAddress(t, v.ResponseAddress), mustHexBytes(t, v.PayloadHex))
			want := v.Expect
			switch want.Outcome {
			case "qualified":
				if got.Outcome != controllerOutcomeQualified {
					t.Fatalf("outcome = %s (%s); want qualified", got.Outcome, got.Reason)
				}
			case "retire_current_qualification_atomically", "may_retain_existing_direct_proof":
				if got.Outcome != controllerOutcomeQualified {
					t.Fatalf("outcome = %s (%s); want a qualified observation", got.Outcome, got.Reason)
				}
			case "reject_without_mutation":
				if got.Outcome != controllerOutcomeRejected || got.Reason != want.Reason {
					t.Fatalf("outcome = %s reason = %s; want reject %s", got.Outcome, got.Reason, want.Reason)
				}
				return
			case "unknown":
				if want.Reason == "register_not_listed" {
					// The pair is declared; the register is judged by the binding.
					if got.Outcome != controllerOutcomeQualified {
						t.Fatalf("register-level vector needs a qualified pair, got %s (%s)", got.Outcome, got.Reason)
					}
				} else if got.Outcome != controllerOutcomeUnknown || got.Reason != want.Reason {
					t.Fatalf("outcome = %s reason = %s; want unknown %s", got.Outcome, got.Reason, want.Reason)
				}
			default:
				t.Fatalf("unhandled fixture outcome %q", want.Outcome)
			}
			if want.DeviceID != "" && got.DeviceID != want.DeviceID {
				t.Fatalf("device ID = %s; want %s", got.DeviceID, want.DeviceID)
			}
			if want.ProductNumber != "" && got.ProductNumber != want.ProductNumber {
				t.Fatalf("product number = %s; want %s", got.ProductNumber, want.ProductNumber)
			}
			if want.SoftwareVersion != "" && got.SoftwareVersion != want.SoftwareVersion {
				t.Fatalf("software version = %s; want %s", got.SoftwareVersion, want.SoftwareVersion)
			}
			if got.Outcome == controllerOutcomeQualified {
				if got.Level != "declared" || got.RuleDigest != ControllerQualificationRuleDigest() {
					t.Fatalf("level = %s digest = %s", got.Level, got.RuleDigest)
				}
			}
		})
	}
	if checked != 13 {
		t.Fatalf("controller vectors checked = %d; want the 13 of the record", checked)
	}
}

// TestControllerQualificationProductNumberDecoding uses the regulator page's
// worked examples: wire bytes 17 04 give decimal 417, shown as 01A1, and 63 04
// gives 463, shown as 01CF. The record's own pairs follow the same reading.
func TestControllerQualificationProductNumberDecoding(t *testing.T) {
	for _, tc := range []struct {
		low, high byte
		want      string
	}{
		{0x17, 0x04, "01A1"},
		{0x63, 0x04, "01CF"},
		{0x21, 0x03, "0141"},
		{0x77, 0x03, "0179"},
		{0x01, 0x00, "0001"},
	} {
		spn, ok := decodeProductNumber(tc.low, tc.high)
		if !ok {
			t.Fatalf("bytes %02X %02X invalid", tc.low, tc.high)
		}
		if got := fmt.Sprintf("%04X", spn); got != tc.want {
			t.Fatalf("bytes %02X %02X = %s; want %s", tc.low, tc.high, got, tc.want)
		}
	}
}

// TestControllerQualificationNegatives covers the rejections the record names
// without a vector: a register the record does not list, a lower-case or
// odd-length registry identity, and a family the record does not declare.
func TestControllerQualificationNegatives(t *testing.T) {
	base := []byte{0xB5, '7', '0', '0', '0', '0', 0x01, 0x23, 0x21, 0x03}
	q := QualifyControllerIdentification(0x15, 0x15, base)
	if q.Outcome != controllerOutcomeQualified {
		t.Fatalf("base pair outcome = %s (%s)", q.Outcome, q.Reason)
	}
	cases := map[string][]byte{
		"unlisted family": {0xB5, 'Z', 'Z', 'Z', '0', '0', 0x01, 0x23, 0x21, 0x03},
		"wrong manuf":     {0x08, '7', '0', '0', '0', '0', 0x01, 0x23, 0x21, 0x03},
		"nine bytes":      base[:9],
		"eleven bytes":    append(append([]byte{}, base...), 0x00),
	}
	for name, payload := range cases {
		got := QualifyControllerIdentification(0x15, 0x15, payload)
		if got.Outcome == controllerOutcomeQualified {
			t.Fatalf("%s qualified; want not declared", name)
		}
	}
}

// TestControllerDeclarationThroughRegistry checks the declaration applied to
// the registry identity of a current controller result.
func TestControllerQualificationThroughRegistry(t *testing.T) {
	reg := registry.NewDeviceRegistry(nil)
	reg.Register(registry.DeviceInfo{
		Address:         0x15,
		Manufacturer:    "Vaillant",
		DeviceID:        "70000",
		SerialNumber:    "SN-TEST-1",
		SoftwareVersion: "0123",
		HardwareVersion: "2103",
	})
	var got registry.QualifiedController
	if !reg.WithCurrentQualifiedController(0x15, DeclareController, func(q registry.QualifiedController) { got = q }) {
		t.Fatal("declared controller produced no result")
	}
	if got.Declaration.Family != "70000" || got.Declaration.ProductNumber != "0141" || got.ProductNumberRaw != "2103" || got.SoftwareVersion != "0123" {
		t.Fatalf("result = %#v", got)
	}
	if got.Declaration.Level != "declared" || got.Declaration.RuleDigest != ControllerQualificationRuleDigest() {
		t.Fatalf("declaration = %#v", got.Declaration)
	}

	reg2 := registry.NewDeviceRegistry(nil)
	reg2.Register(registry.DeviceInfo{
		Address:         0x15,
		Manufacturer:    "Vaillant",
		DeviceID:        "ZZZ00",
		SerialNumber:    "SN-TEST-2",
		SoftwareVersion: "0123",
		HardwareVersion: "2103",
	})
	if reg2.WithCurrentQualifiedController(0x15, DeclareController, func(registry.QualifiedController) {
		t.Fatal("unlisted identity produced a controller result")
	}) {
		t.Fatal("unlisted identity returned a result")
	}
}
