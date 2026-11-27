package owner

import (
	"context"
	"fmt"
	ebuserrors "github.com/Project-Helianthus/helianthus-ebusgo/errors"
	"github.com/Project-Helianthus/helianthus-ebusreg/driver/vaillant/owner/systemconfig"
	"strings"
	"testing"
)

func TestVaillantOwnerConfigSetterNamesSystemKeysWriteReadBackAndPublish(t *testing.T) {
	cases := cfgKeysSystemCases()
	seen := map[string]bool{}
	for _, tc := range cases {
		seen[tc.key] = true
	}
	if len(seen) != len(systemconfig.FieldSpecs()) {
		t.Fatalf("cases cover %d system keys; specs have %d", len(seen), len(systemconfig.FieldSpecs()))
	}
	for _, tc := range cases {
		t.Run(tc.key, func(t *testing.T) { runCfgKeysCase(t, "set_system_config", tc) })
	}
}

func TestVaillantOwnerConfigSetterNamesBoilerKeysWriteReadBackAndPublish(t *testing.T) {
	cases := cfgKeysBoilerCases()
	seen := map[string]bool{}
	for _, tc := range cases {
		seen[tc.key] = true
	}
	if len(seen) != len(boilerConfigFieldSpecs)+2 {
		t.Fatalf("cases cover %d boiler keys; setter accepts %d", len(seen), len(boilerConfigFieldSpecs)+2)
	}
	for i, tc := range cases {
		t.Run(fmt.Sprintf("%s/%d", tc.key, i), func(t *testing.T) { runCfgKeysCase(t, "set_boiler_config", tc) })
	}
}

// TestConfigSetterNamesSystemPhoneValidationFollowsKey catches the phone
// character check of the pair writer staying on the former key: a letter in
// installer_phone_number is rejected before any bus request.
func TestVaillantOwnerConfigSetterNamesSystemPhoneValidationFollowsKey(t *testing.T) {
	bus := newCfgKeysFakeBus()
	poller, _ := newCfgKeysPoller(t, bus)
	result := poller.SetSystemConfig(context.Background(), "installer_phone_number", "12a")
	if result.Success || !strings.Contains(result.Error, "allowed: digits") {
		t.Fatalf("SetSystemConfig(installer_phone_number, 12a) = %+v; want the phone character rejection", result)
	}
	if requests, _, _ := bus.snapshot(); requests != 0 {
		t.Fatalf("bus requests = %d; want none", requests)
	}
}

// TestConfigSetterNamesFormerKeysRejected sends every former camelCase key,
// every name the setters no longer accept (among them the former boiler key
// installerMenuCode) and spellings of the boiler pump mode key that differ from
// the documentation name to the system and boiler setters.
// Each must fail before any bus request with the exact list of accepted keys.
func TestVaillantOwnerConfigSetterNamesFormerKeysRejected(t *testing.T) {
	cases := []struct {
		mutation, want string
		keys           []string
	}{
		{
			mutation: "set_system_config",
			want:     "unknown system field %q (allowed: " + cfgKeysSystemAllowed + "): " + ebuserrors.ErrInvalidPayload.Error(),
			keys: []string{
				"dhwBivalencePointC", "maxRoomHumidityPct", "adaptiveHeatingCurve", "hcBivalencePointC", "hcEmergencyTempC",
				"hwcMaxFlowTempC", "maintenanceDate", "installerName", "installerPhone", "installerMenuCode",
				"hc_bivalence_point_c", "dhw_bivalence_point_c", "hc_emergency_temp_c", "hwc_max_flow_temp_c",
				"max_room_humidity_pct", "installer_phone", "installer_menu_code",
			},
		},
		{
			mutation: "set_boiler_config",
			want:     "unknown boiler field %q (allowed: " + cfgKeysBoilerAllowed + ")",
			keys: []string{
				"flowsetHcMaxC", "flowsetHwcMaxC", "partloadHcKW", "partloadHwcKW", "phoneNumber", "phone_number", "installer_menu_code",
				"installerMenuCode", "CentralHeatingPumpMode", "central_heating_pump_mode ",
			},
		},
	}
	for _, tc := range cases {
		for _, key := range tc.keys {
			t.Run(tc.mutation+"/"+key, func(t *testing.T) {
				bus := newCfgKeysFakeBus()
				poller, _ := newCfgKeysPoller(t, bus)
				success, message := cfgKeysSet(t, tc.mutation, poller, key, "1")
				if want := fmt.Sprintf(tc.want, key); success || message != want {
					t.Fatalf("success=%v error=%q; want %q", success, message, want)
				}
				if requests, _, _ := bus.snapshot(); requests != 0 {
					t.Fatalf("bus requests = %d; want none", requests)
				}
			})
		}
	}
}

// TestConfigSetterNamesFixtureComparisonCoversEveryField changes each setter
// field of each system and boiler row in turn (key, opcode, group, instance,
// registers, value type, enum) and requires the comparison to fail.
func TestVaillantOwnerConfigSetterNamesFixtureComparisonCoversEveryField(t *testing.T) {
	mutations := map[string]func(*cfgKeysRowDoc){
		"key":       func(row *cfgKeysRowDoc) { row.Key += "_x" },
		"opcode":    func(row *cfgKeysRowDoc) { row.Selector.Opcode = "0x7E" },
		"group":     func(row *cfgKeysRowDoc) { row.Selector.Group = "0x7F" },
		"instance":  func(row *cfgKeysRowDoc) { row.Selector.Instance = "index" },
		"register":  func(row *cfgKeysRowDoc) { row.Selector.Registers = []string{"0xFFFF"} },
		"registers": func(row *cfgKeysRowDoc) { row.Selector.Registers = append(row.Selector.Registers, "0xFFFF") },
		"value":     func(row *cfgKeysRowDoc) { row.ValueType = "uint8" },
		"enum":      func(row *cfgKeysRowDoc) { row.Enum = map[string]uint16{"other": 9} },
	}
	for name, want := range map[string][]cfgKeysRowDoc{"set_system_config": cfgKeysSystemRows(t), "set_boiler_config": cfgKeysBoilerRows()} {
		clone := func() []cfgKeysRowDoc {
			rows := make([]cfgKeysRowDoc, len(want))
			for i, row := range want {
				row.Selector.Registers = append([]string(nil), row.Selector.Registers...)
				rows[i] = row
			}
			return rows
		}
		if err := cfgKeysDiff(name, want, clone()); err != nil {
			t.Fatalf("unchanged rows rejected: %v", err)
		}
		for i := range want {
			for field, mutate := range mutations {
				got := clone()
				mutate(&got[i])
				if err := cfgKeysDiff(name, want, got); err == nil {
					t.Fatalf("%s row %s: changed %s not detected", name, want[i].Key, field)
				}
			}
		}
		if err := cfgKeysDiff(name, want, clone()[1:]); err == nil {
			t.Fatalf("%s: missing key not detected", name)
		}
	}
}

// TestConfigSetterNamesBoilerPumpModeKeepsItsBound catches a change of the
// value type or bound of the renamed 0x4904 key: a value outside one unsigned
// byte or a non-integer is rejected before any bus request.
func TestVaillantOwnerConfigSetterNamesBoilerPumpModeKeepsItsBound(t *testing.T) {
	for value, want := range map[string]string{
		"256": "value 256 out of range [0, 255]",
		"-1":  "value -1 out of range [0, 255]",
		"1.5": `invalid integer "1.5"`,
	} {
		t.Run(value, func(t *testing.T) {
			bus := newCfgKeysFakeBus()
			poller, _ := newCfgKeysPoller(t, bus)
			success, message := cfgKeysSet(t, "set_boiler_config", poller, "central_heating_pump_mode", value)
			if success || !strings.HasPrefix(message, want) {
				t.Fatalf("success=%v error=%q; want prefix %q", success, message, want)
			}
			if requests, _, _ := bus.snapshot(); requests != 0 {
				t.Fatalf("bus requests = %d; want none", requests)
			}
		})
	}
}
