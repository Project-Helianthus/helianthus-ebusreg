package input

import (
	"encoding/hex"
	"sort"
	"strings"
	"testing"
)

type completenessRow struct {
	id, service, primary, secondary string
}

type completenessInventory struct {
	Selectors []struct {
		ID        string `json:"id"`
		Service   string `json:"service"`
		Geometry  string `json:"geometry"`
		Primary   string `json:"primary"`
		Secondary string `json:"secondary"`
	} `json:"selectors"`
}

func completenessRows(t *testing.T) map[string]completenessRow {
	t.Helper()
	rows := map[string]completenessRow{}
	add := func(id string, s Service, name string) {
		t.Helper()
		if _, dup := rows[id]; dup {
			t.Fatalf("row id %s is defined twice", id)
		}
		rows[id] = completenessRow{id: id, service: name,
			primary: "0x" + strings.ToUpper(hex.EncodeToString([]byte{s.Primary})), secondary: "0x" + strings.ToUpper(hex.EncodeToString([]byte{s.Secondary}))}
	}
	for _, r := range ExtendedRows() {
		add(r.ID, B524Selector{}.Service(), "B524")
	}
	for _, r := range InterfaceIdentityRows {
		add(r.ID, r.Selector().Service(), "B524")
	}
	for _, r := range BoilerRegisterRows {
		add(r.ID, r.Selector().Service(), "B509")
	}
	for _, r := range BoilerWriteRows {
		add(r.ID, r.Selector().Service(), "B509")
	}
	for _, r := range DiagnosticRows {
		add(r.ID, DiagnosticSelector{}.Service(), "B503")
	}
	for _, r := range ScheduleRows {
		add(r.ID, ScheduleConfigSelector{}.Service(), "B555")
	}
	return rows
}

func TestVaillantInputInventory(t *testing.T) {
	inv := boilerRowsLoad[completenessInventory](t, "testdata/gateway_selectors.json")
	rows := completenessRows(t)
	listed := map[string]bool{}
	perService := map[string]int{}
	for _, s := range inv.Selectors {
		if listed[s.ID] {
			t.Errorf("inventoried selector %s is listed twice", s.ID)
		}
		listed[s.ID] = true
		perService[s.Service]++
		row, ok := rows[s.ID]
		if !ok {
			t.Errorf("inventoried selector %s has no row", s.ID)
			continue
		}
		if row.service != s.Service || row.primary != s.Primary || row.secondary != s.Secondary {
			t.Errorf("row %s is %s %s/%s, the inventory says %s %s/%s", s.ID, row.service, row.primary, row.secondary, s.Service, s.Primary, s.Secondary)
		}
	}
	var orphans []string
	for id := range rows {
		if !listed[id] {
			orphans = append(orphans, id)
		}
	}
	sort.Strings(orphans)
	if len(orphans) > 0 {
		t.Errorf("rows without an inventoried selector: %v", orphans)
	}
	if len(inv.Selectors) != len(rows) {
		t.Errorf("%d inventoried selectors, %d rows", len(inv.Selectors), len(rows))
	}
	want := map[string]int{"B524": 129, "B509": 43, "B503": 5, "B555": 4}
	for service, n := range want {
		if perService[service] != n {
			t.Errorf("service %s has %d inventoried selectors, want %d", service, perService[service], n)
		}
	}
	t.Logf("inventoried selectors with a row and rows with an inventoried selector: %d (B524 %d, B509 %d, B503 %d, B555 %d)",
		len(rows), perService["B524"], perService["B509"], perService["B503"], perService["B555"])
}

func TestVaillantInputInventoryGeometries(t *testing.T) {
	inv := boilerRowsLoad[completenessInventory](t, "testdata/gateway_selectors.json")
	geometry := map[string]string{}
	for _, s := range inv.Selectors {
		geometry[s.ID] = s.Geometry
	}
	check := func(id, want string) {
		t.Helper()
		if geometry[id] != want {
			t.Errorf("row %s has geometry %q in the inventory, want %s", id, geometry[id], want)
		}
	}
	for _, r := range ExtendedRows() {
		check(r.ID, "B524Selector")
	}
	for _, r := range InterfaceIdentityRows {
		check(r.ID, "InterfaceIdentitySelector")
	}
	for _, r := range BoilerRegisterRows {
		check(r.ID, "B509Selector")
	}
	for _, r := range BoilerWriteRows {
		check(r.ID, "B509Selector")
	}
	for _, r := range DiagnosticRows {
		if r.LiveMonitor {
			check(r.ID, "LiveMonitorSelector")
			continue
		}
		check(r.ID, "DiagnosticSelector")
	}
	for _, r := range ScheduleRows {
		switch r.Operation {
		case ScheduleConfigRead:
			check(r.ID, "ScheduleConfigSelector")
		case ScheduleCountsRead:
			check(r.ID, "ScheduleCountsSelector")
		case ScheduleTimerRead:
			check(r.ID, "ScheduleTimerSelector")
		default:
			check(r.ID, "ScheduleTimerWriteSelector")
		}
	}
}
