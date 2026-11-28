package writes

import (
	"encoding/json"
	"testing"
)

func TestPrePortWritesHaveOperations(t *testing.T) {
	rows, err := Inventory()
	if err != nil {
		t.Fatalf("Inventory() error = %v", err)
	}
	bySurface := map[string]map[string]Row{}
	byOperation := map[string]Row{}
	for _, row := range rows {
		if row.Surface == "" || row.Operation == "" || row.Selector.Family == "" || row.Selector.Primary == "" || row.Selector.Secondary == "" || row.Selector.Opcode == "" || row.Selector.Group == "" || row.Selector.Instance == "" {
			t.Fatalf("row missing identity: %+v", row)
		}
		if row.Selector.Family != "timer" && row.Selector.Register == "" {
			t.Fatalf("row %s/%s has no register", row.Surface, row.Key)
		}
		if _, ok := bySurface[row.Surface]; !ok {
			bySurface[row.Surface] = map[string]Row{}
		}
		if _, ok := bySurface[row.Surface][row.Key]; ok {
			t.Fatalf("duplicate surface %s key %s", row.Surface, row.Key)
		}
		bySurface[row.Surface][row.Key] = row
		previous, ok := byOperation[row.Operation]
		if ok && !sameOperation(previous, row) {
			t.Fatalf("operation %s disagrees", row.Operation)
		}
		byOperation[row.Operation] = row
		if _, ok := Operation(row.Operation); !ok {
			t.Fatalf("Operation(%s) missing", row.Operation)
		}
	}
	requireKeys(t, bySurface, "set_zone_config", "desiredSetpoint", "operatingMode", "quickVetoTemperature", "quickVetoDuration", "holidaySetpoint")
	requireKeys(t, bySurface, "set_circuit_config", "heatingCurve", "flowTempMaxC", "flowTempMinC", "summerLimitC", "frostProtC", "roomTempControl", "coolingEnabled")
	systemKeys := []string{"dhwBivalencePointC", "maxRoomHumidityPct", "adaptiveHeatingCurve", "hcBivalencePointC", "hcEmergencyTempC", "hwcMaxFlowTempC", "maintenanceDate", "installerName", "installerPhone", "installerMenuCode"}
	requireKeys(t, bySurface, "set_system_config", systemKeys...)
	requireKeys(t, bySurface, "ebus.v1.semantic.system.set_config", systemKeys...)
	boilerKeys := []string{"flowsetHcMaxC", "flowsetHwcMaxC", "partloadHcKW", "partloadHwcKW", "phoneNumber", "installerMenuCode"}
	requireKeys(t, bySurface, "set_boiler_config", boilerKeys...)
	requireKeys(t, bySurface, "ebus.v1.semantic.boiler_status.set_config", boilerKeys...)
	requireKeys(t, bySurface, "set_zone_time_program", "zone", "weekday", "slots")
	requireKeys(t, bySurface, "ebus_v1_semantic_schedules_set_zone_time_program", "zone", "weekday", "slots")
	requireKeys(t, bySurface, "set_dhw_time_program", "weekday", "slots")
	requireKeys(t, bySurface, "ebus_v1_semantic_schedules_set_dhw_time_program", "weekday", "slots")

	setpoint := bySurface["set_zone_config"]["desiredSetpoint"]
	if setpoint.Operation != "thermal.operation.set_temperature" || setpoint.Selector.Group != "03" || setpoint.Selector.Register != "0022" || setpoint.Kind != "kelvin_celsius" || setpoint.Argument != "thermal.setpoint.temperature_setpoint" {
		t.Fatalf("zone setpoint = %+v", setpoint)
	}
	mode := bySurface["thermal.operation.set_mode"][""]
	if mode.Operation != "thermal.operation.set_mode" || mode.Key != "" || mode.Selector.Group != "00" || mode.Selector.Register != "0007" || mode.Kind != "system_switch" || mode.Selector.Instance != "00" {
		t.Fatalf("system mode = %+v", mode)
	}
	if bySurface["ebusreg"]["dhw_operating_mode"].Selector.Register != "0003" || bySurface["ebusreg"]["dhw_temperature_setpoint"].Selector.Register != "0004" || bySurface["ebusreg"]["dhw_temperature_setpoint"].Selector.Group != "01" {
		t.Fatal("group 01 registers 0003 and 0004 missing")
	}
	frost := bySurface["set_circuit_config"]["frostProtC"]
	if frost.Form != "unknown" || frost.Selector.Register != "001d" || frost.Kind != "" {
		t.Fatalf("frost = %+v", frost)
	}
	name := bySurface["set_system_config"]["installerName"]
	phone := bySurface["set_system_config"]["installerPhone"]
	if name.Form != "compound" || name.PairRegister != "006d" || phone.Form != "compound" || phone.PairRegister != "0070" {
		t.Fatalf("installer strings = %+v %+v", name, phone)
	}
	for _, key := range boilerKeys {
		row := bySurface["set_boiler_config"][key]
		if row.Form != "unknown" || row.Selector.Family != "boiler" || row.Selector.Register == "a500" || row.Selector.Opcode != "0e" {
			t.Fatalf("boiler %s = %+v", key, row)
		}
		if !sameOperation(row, bySurface["ebus.v1.semantic.boiler_status.set_config"][key]) {
			t.Fatalf("boiler surfaces disagree for %s", key)
		}
	}
	for _, row := range rows {
		if row.Selector.Group == "03" && row.Selector.Register == "0014" {
			t.Fatalf("effective setpoint is not an inventory write: %+v", row)
		}
	}
	assertTimeProgram(t, bySurface, "set_zone_time_program", "ebus_v1_semantic_schedules_set_zone_time_program", "ebusreg.operation.zone_heating_time_program", "string", "object_array", "required", []string{"zone", "weekday", "slots"})
	assertTimeProgram(t, bySurface, "set_dhw_time_program", "ebus_v1_semantic_schedules_set_dhw_time_program", "ebusreg.operation.dhw_time_program", "string", "object_array", "optional", []string{"weekday", "slots"})
	assertInventoryMutations(t, rows)
}

func requireKeys(t *testing.T, bySurface map[string]map[string]Row, surface string, keys ...string) {
	t.Helper()
	got := bySurface[surface]
	if got == nil {
		t.Fatalf("surface %s missing", surface)
	}
	for _, key := range keys {
		if _, ok := got[key]; !ok {
			t.Fatalf("surface %s key %s missing", surface, key)
		}
	}
}

func assertTimeProgram(t *testing.T, bySurface map[string]map[string]Row, graphqlSurface, mcpSurface, operation, graphqlSlots, mcpSlots, temperature string, keys []string) {
	t.Helper()
	var selector Selector
	var seen bool
	for _, surface := range []string{graphqlSurface, mcpSurface} {
		for _, key := range keys {
			row := bySurface[surface][key]
			if row.Operation != operation || row.Form != "compound" {
				t.Fatalf("%s/%s = %+v", surface, key, row)
			}
			wantSlots := graphqlSlots
			if surface == mcpSurface {
				wantSlots = mcpSlots
			}
			if row.Slots != wantSlots || row.SlotTemperature != temperature {
				t.Fatalf("%s/%s slots = %s %s", surface, key, row.Slots, row.SlotTemperature)
			}
			if !seen {
				selector = row.Selector
				seen = true
			}
			if row.Selector != selector {
				t.Fatalf("%s/%s selector drifted", surface, key)
			}
		}
	}
}

func assertInventoryMutations(t *testing.T, rows []Row) {
	t.Helper()
	base := rows[0]
	copyRow := base
	copyRow.Surface = base.Surface + "-copy"
	copyRow.Key = base.Key + "-copy"
	changes := []struct {
		name   string
		mutate func(*Row)
	}{
		{"family", func(row *Row) { row.Selector.Family = "boiler" }},
		{"primary", func(row *Row) { row.Selector.Primary = "b6" }},
		{"secondary", func(row *Row) { row.Selector.Secondary = "25" }},
		{"opcode", func(row *Row) { row.Selector.Opcode = "06" }},
		{"group", func(row *Row) { row.Selector.Group = "04" }},
		{"register", func(row *Row) { row.Selector.Register = "00ff" }},
		{"instance", func(row *Row) { row.Selector.Instance = "11" }},
		{"form", func(row *Row) { row.Form = "compound"; row.Kind = "" }},
		{"kind", func(row *Row) { row.Kind = "bool_byte" }},
		{"argument", func(row *Row) { row.Argument = "ebusreg.operation.other" }},
		{"grade", func(row *Row) { row.Grade = "unknown" }},
		{"pair", func(row *Row) { row.PairRegister = "0001" }},
	}
	for _, change := range changes {
		mutated := append(cloneRows(rows), copyRow)
		change.mutate(&mutated[len(mutated)-1])
		if err := validateRows(mutated); err == nil {
			t.Fatalf("changed %s was accepted", change.name)
		}
	}
	duplicated := append(cloneRows(rows), base)
	if err := validateRows(duplicated); err == nil {
		t.Fatal("duplicate surface and key was accepted")
	}
	raw, err := json.Marshal(inventoryDocument{Rows: []Row{{Surface: "set_zone_config"}}})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := ParseInventory(raw); err == nil {
		t.Fatal("row without an operation was accepted")
	}
}

func TestDeclaredWritableKeysHaveOperations(t *testing.T) {
	declared := &Declared{}
	base := RegisterIdentity{Opcode: 0x11, Group: 0x22, Instance: 0x33, Register: 0x4455}
	if _, ok := declared.OperationID(base); ok {
		t.Fatal("unobserved identity has an operation")
	}
	declared.Observe(base.Opcode, base.Group, base.Instance, base.Register, []byte{0x01, base.Group, byte(base.Register), byte(base.Register >> 8)})
	declared.Observe(base.Opcode, base.Group, base.Instance, base.Register, []byte{0x00, base.Group, byte(base.Register), byte(base.Register >> 8)})
	if _, ok := declared.OperationID(base); ok {
		t.Fatal("attribute bit 0 alone admitted an operation")
	}
	declared.Observe(base.Opcode, base.Group, base.Instance, base.Register, []byte{0x02})
	declared.Observe(base.Opcode, base.Group, base.Instance, base.Register, []byte{0x03, base.Group + 1, byte(base.Register), byte(base.Register >> 8)})
	if _, ok := declared.OperationID(base); ok {
		t.Fatal("short or uncorrelated reply admitted an operation")
	}
	header := func(flags, group byte, register uint16) []byte {
		return []byte{flags, group, byte(register), byte(register >> 8), 0x10}
	}
	declared.Observe(base.Opcode, base.Group, base.Instance, base.Register, header(0x02, base.Group, base.Register))
	first, ok := declared.OperationID(base)
	if !ok {
		t.Fatal("attribute bit 1 did not admit an operation")
	}
	grade, ok := declared.Grade(base)
	if !ok || grade != "hypothesis" {
		t.Fatalf("grade = %q ok=%v", grade, ok)
	}
	declared.Observe(base.Opcode, base.Group, base.Instance, base.Register, header(0x03, base.Group, base.Register))
	if again, ok := declared.OperationID(base); !ok || again != first {
		t.Fatalf("attribute 0x03 = %q ok=%v", again, ok)
	}
	fields := []RegisterIdentity{
		{Opcode: base.Opcode + 1, Group: base.Group, Instance: base.Instance, Register: base.Register},
		{Opcode: base.Opcode, Group: base.Group + 1, Instance: base.Instance, Register: base.Register},
		{Opcode: base.Opcode, Group: base.Group, Instance: base.Instance + 1, Register: base.Register},
		{Opcode: base.Opcode, Group: base.Group, Instance: base.Instance, Register: base.Register + 1},
	}
	seen := map[string]struct{}{first: {}}
	for _, id := range fields {
		declared.Observe(id.Opcode, id.Group, id.Instance, id.Register, header(0x02, id.Group, id.Register))
		got, ok := declared.OperationID(id)
		if !ok {
			t.Fatalf("identity %+v has no operation", id)
		}
		if _, dup := seen[got]; dup || got == first {
			t.Fatalf("identity %+v reused %s", id, got)
		}
		seen[got] = struct{}{}
	}
	declared.Observe(base.Opcode, base.Group, base.Instance, base.Register, []byte{0x01, 0x00, 0x00, 0x00})
	if _, ok := declared.OperationID(base); !ok {
		t.Fatal("uncorrelated reply withdrew an admitted operation")
	}
	declared.Observe(base.Opcode, base.Group, base.Instance, base.Register, header(0x01, base.Group, base.Register))
	if _, ok := declared.OperationID(base); ok {
		t.Fatal("cleared attribute bit kept the operation")
	}

	zone := RegisterIdentity{Opcode: 0x02, Group: 0x03, Instance: 0x17, Register: 0x0022}
	other := zone
	other.Instance = 0x05
	missing := zone
	missing.Instance = 0x06
	remote := zone
	remote.Opcode = 0x06
	declared = &Declared{}
	declared.Observe(zone.Opcode, zone.Group, zone.Instance, zone.Register, header(0x02, zone.Group, zone.Register))
	declared.Observe(other.Opcode, other.Group, other.Instance, other.Register, header(0x02, other.Group, other.Register))
	declared.Observe(remote.Opcode, remote.Group, remote.Instance, remote.Register, header(0x02, remote.Group, remote.Register))
	for _, id := range []RegisterIdentity{zone, other} {
		got, ok := declared.OperationID(id)
		if !ok || got != "thermal.operation.set_temperature" {
			t.Fatalf("catalog identity %+v = %q ok=%v", id, got, ok)
		}
	}
	if _, ok := declared.OperationID(missing); ok {
		t.Fatal("unobserved instance inherited the catalog operation")
	}
	if got, ok := declared.OperationForSelector("local", 0x02, 0x03, 0x0022); !ok || got != "thermal.operation.set_temperature" {
		t.Fatalf("selector operation = %q ok=%v", got, ok)
	}
	if got, ok := declared.OperationID(remote); !ok || got == "thermal.operation.set_temperature" {
		t.Fatalf("other opcode = %q ok=%v", got, ok)
	}
	if _, ok := declared.OperationForSelector("local", remote.Opcode, remote.Group, remote.Register); ok {
		t.Fatal("other opcode matched the catalog operation")
	}
	if _, ok := declared.OperationForSelector("boiler", zone.Opcode, zone.Group, zone.Register); ok {
		t.Fatal("boiler family matched a local reply")
	}
}
