package input

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"os"
	"regexp"
	"strconv"
	"strings"
	"testing"
)

type inventoryFile struct {
	SchemaVersion int    `json:"schema_version"`
	Description   string `json:"description"`
	SourceFiles   []struct {
		Path   string `json:"path"`
		Digest string `json:"sha256"`
	} `json:"source_files"`
	Selectors []inventoryEntry `json:"selectors"`
	Gaps      []struct {
		ID   string `json:"id"`
		Text string `json:"text"`
	} `json:"gaps"`
}

type inventoryEntry struct {
	ID              string   `json:"id"`
	Service         string   `json:"service"`
	Geometry        string   `json:"geometry"`
	Operation       string   `json:"operation"`
	GatewayConstant string   `json:"gateway_constant"`
	DefinedIn       string   `json:"defined_in"`
	Primary         string   `json:"primary"`
	Secondary       string   `json:"secondary"`
	Opcode          string   `json:"opcode"`
	OperationByte   string   `json:"operation_byte"`
	Group           string   `json:"group"`
	InstanceFirst   string   `json:"instance_first"`
	InstanceLast    string   `json:"instance_last"`
	Register        string   `json:"register"`
	IdentityID      string   `json:"identity_id"`
	Family          string   `json:"family"`
	Selector        string   `json:"selector"`
	RequestHex      string   `json:"request_hex"`
	RequestHexKind  string   `json:"request_hex_kind"`
	RequestLayout   string   `json:"request_layout"`
	Fields          []string `json:"fields"`
	NativeViewField string   `json:"native_view_field"`
	Reader          string   `json:"reader"`
	Purposes        []string `json:"purposes"`
	CallSites       []string `json:"call_sites"`
	GroupConstant   string   `json:"group_constant"`
	OpcodeConstant  string   `json:"opcode_constant"`
	Target          string   `json:"target"`
	ProgramPairs    string   `json:"program_pairs"`
	Actions         []string `json:"actions"`
	Note            string   `json:"note"`
}

func loadInventory(t *testing.T) inventoryFile {
	t.Helper()
	raw, err := os.ReadFile("testdata/gateway_selectors.json")
	if err != nil {
		t.Fatal(err)
	}
	var inv inventoryFile
	dec := json.NewDecoder(bytes.NewReader(raw))
	dec.DisallowUnknownFields()
	if err := dec.Decode(&inv); err != nil {
		t.Fatal(err)
	}
	return inv
}

func parseByte(t *testing.T, text string) byte {
	t.Helper()
	n, err := strconv.ParseUint(text, 0, 8)
	if err != nil {
		t.Fatalf("byte %q: %v", text, err)
	}
	return byte(n)
}

func parseWord(t *testing.T, text string) uint16 {
	t.Helper()
	n, err := strconv.ParseUint(text, 0, 16)
	if err != nil {
		t.Fatalf("word %q: %v", text, err)
	}
	return uint16(n)
}

func selectorFor(t *testing.T, e inventoryEntry) Selector {
	t.Helper()
	switch e.Geometry {
	case "B524Selector":
		return B524Selector{
			Opcode: parseByte(t, e.Opcode), Operation: parseByte(t, e.OperationByte), Group: parseByte(t, e.Group),
			Instance: parseByte(t, e.InstanceFirst), Register: parseWord(t, e.Register),
		}
	case "B509Selector":
		return B509Selector{Opcode: parseByte(t, e.Opcode), Register: parseWord(t, e.Register)}
	case "InterfaceIdentitySelector":
		return InterfaceIdentitySelector{ID: parseByte(t, e.IdentityID)}
	case "DiagnosticSelector":
		return DiagnosticSelector{Family: parseByte(t, e.Family), Selector: parseByte(t, e.Selector)}
	case "LiveMonitorSelector":
		return LiveMonitorSelector{Action: LiveMonitorRead}
	}
	return nil
}

func TestVaillantInputSelectorsInventoryBuildsRecordedRequests(t *testing.T) {
	inv := loadInventory(t)
	if inv.SchemaVersion != 1 || len(inv.Selectors) == 0 {
		t.Fatalf("inventory header: version %d, %d selectors", inv.SchemaVersion, len(inv.Selectors))
	}
	seen := map[string]int{}
	ids := map[string]bool{}
	for _, e := range inv.Selectors {
		if ids[e.ID] {
			t.Errorf("duplicate inventory id %s", e.ID)
		}
		ids[e.ID] = true
		seen[e.Geometry]++
		if e.GatewayConstant == "" || e.DefinedIn == "" || e.NativeViewField == "" || e.Reader == "" || len(e.CallSites) == 0 || len(e.Purposes) == 0 {
			t.Errorf("%s: constant, defining file, native view field, reader, purposes and call sites are all required", e.ID)
		}
		switch e.Geometry {
		case "ScheduleConfigSelector", "ScheduleCountsSelector", "ScheduleTimerSelector", "ScheduleTimerWriteSelector":
			checkScheduleEntry(t, e)
			continue
		}
		sel := selectorFor(t, e)
		if sel == nil {
			t.Errorf("%s: geometry %s has no selector type", e.ID, e.Geometry)
			continue
		}
		svc := sel.Service()
		if svc.Primary != parseByte(t, e.Primary) || svc.Secondary != parseByte(t, e.Secondary) {
			t.Errorf("%s: service %02x/%02x, inventory %s/%s", e.ID, svc.Primary, svc.Secondary, e.Primary, e.Secondary)
		}
		got := hex.EncodeToString(sel.Payload())
		switch e.RequestHexKind {
		case "complete":
			if got != e.RequestHex {
				t.Errorf("%s: payload %s, inventory %s", e.ID, got, e.RequestHex)
			}
		case "header_before_value":
			if got != e.RequestHex {
				t.Errorf("%s: header %s, inventory %s", e.ID, got, e.RequestHex)
			}
		case "body_with_optional_index_byte":
			if got != e.RequestHex {
				t.Errorf("%s: body %s, inventory %s", e.ID, got, e.RequestHex)
			}
		default:
			t.Errorf("%s: request kind %q", e.ID, e.RequestHexKind)
		}
	}
	for _, geometry := range []string{
		"B524Selector", "B509Selector", "InterfaceIdentitySelector", "ScheduleConfigSelector", "ScheduleCountsSelector",
		"ScheduleTimerSelector", "ScheduleTimerWriteSelector", "DiagnosticSelector", "LiveMonitorSelector",
	} {
		if seen[geometry] == 0 {
			t.Errorf("selector type %s has no inventoried selector", geometry)
		}
	}
	if len(seen) != 9 {
		t.Errorf("inventory geometries: %v", seen)
	}
}

func TestVaillantInputSelectorsInventoryBuildsEveryInstanceOfARange(t *testing.T) {
	inv := loadInventory(t)
	for _, e := range inv.Selectors {
		if e.Geometry != "B524Selector" {
			continue
		}
		first, last := parseByte(t, e.InstanceFirst), parseByte(t, e.InstanceLast)
		if last < first {
			t.Fatalf("%s: instance range %02x..%02x", e.ID, first, last)
		}
		for instance := first; ; instance++ {
			sel := B524Selector{
				Opcode: parseByte(t, e.Opcode), Operation: parseByte(t, e.OperationByte), Group: parseByte(t, e.Group),
				Instance: instance, Register: parseWord(t, e.Register),
			}
			p := sel.Payload()
			if len(p) != 6 || p[3] != instance {
				t.Errorf("%s instance %02x: payload %x", e.ID, instance, p)
			}
			if instance == last {
				break
			}
		}
	}
}

func TestVaillantInputSelectorsInventoryNamesTheFramesItIsTakenFrom(t *testing.T) {
	inv := loadInventory(t)
	digest := regexp.MustCompile(`^[0-9a-f]{64}$`)
	paths := map[string]bool{}
	for _, f := range inv.SourceFiles {
		if !strings.HasPrefix(f.Path, "cmd/gateway/") || !strings.HasSuffix(f.Path, ".go") || !digest.MatchString(f.Digest) {
			t.Errorf("source file entry %+v", f)
		}
		if paths[f.Path] {
			t.Errorf("duplicate source file %s", f.Path)
		}
		paths[f.Path] = true
	}
	for _, e := range inv.Selectors {
		if !paths[e.DefinedIn] {
			t.Errorf("%s: defining file %s has no digest", e.ID, e.DefinedIn)
		}
		for _, site := range e.CallSites {
			file, _, ok := strings.Cut(site, ": ")
			if !ok || !paths[file] {
				t.Errorf("%s: call site %q names a file without a digest", e.ID, site)
			}
		}
	}
	if len(inv.Gaps) == 0 {
		t.Error("the evidence gaps of the inventory are recorded")
	}
}

func TestVaillantInputSelectorsInventoryCoversEveryServiceTheFilesSend(t *testing.T) {
	inv := loadInventory(t)
	counts := map[string]int{}
	for _, e := range inv.Selectors {
		counts[e.Service]++
	}
	for _, service := range []string{"B524", "B509", "B555", "B503"} {
		if counts[service] == 0 {
			t.Errorf("no %s selector inventoried", service)
		}
	}
	if counts["B555"] != 4 || counts["B503"] != 5 {
		t.Errorf("schedule and diagnostic selector counts: %v", counts)
	}
}

func TestVaillantInputSelectorsVaryingAnyFieldChangesExactlyItsWireBytes(t *testing.T) {
	cases := []struct {
		name   string
		base   Selector
		other  Selector
		offset int
	}{
		{"b524 opcode", B524Selector{0x02, 0, 3, 1, 0x1234}, B524Selector{0x06, 0, 3, 1, 0x1234}, 0},
		{"b524 operation", B524Selector{0x02, 0, 3, 1, 0x1234}, B524Selector{0x02, 1, 3, 1, 0x1234}, 1},
		{"b524 group", B524Selector{0x02, 0, 3, 1, 0x1234}, B524Selector{0x02, 0, 4, 1, 0x1234}, 2},
		{"b524 instance", B524Selector{0x02, 0, 3, 1, 0x1234}, B524Selector{0x02, 0, 3, 2, 0x1234}, 3},
		{"b524 register low", B524Selector{0x02, 0, 3, 1, 0x1234}, B524Selector{0x02, 0, 3, 1, 0x1235}, 4},
		{"b524 register high", B524Selector{0x02, 0, 3, 1, 0x1234}, B524Selector{0x02, 0, 3, 1, 0x1334}, 5},
		{"b509 opcode", B509Selector{0x0d, 0x4400}, B509Selector{0x0e, 0x4400}, 0},
		{"b509 register high", B509Selector{0x0d, 0x4400}, B509Selector{0x0d, 0x4500}, 1},
		{"b509 register low", B509Selector{0x0d, 0x4400}, B509Selector{0x0d, 0x4401}, 2},
		{"identity id", InterfaceIdentitySelector{0x06}, InterfaceIdentitySelector{0x07}, 1},
		{"schedule config zone", ScheduleConfigSelector{1, 2}, ScheduleConfigSelector{2, 2}, 1},
		{"schedule config circuit", ScheduleConfigSelector{1, 2}, ScheduleConfigSelector{1, 3}, 2},
		{"schedule counts zone", ScheduleCountsSelector{1, 2}, ScheduleCountsSelector{2, 2}, 1},
		{"schedule counts circuit", ScheduleCountsSelector{1, 2}, ScheduleCountsSelector{1, 3}, 2},
		{"schedule timer zone", ScheduleTimerSelector{1, 2, 3, 4}, ScheduleTimerSelector{2, 2, 3, 4}, 1},
		{"schedule timer circuit", ScheduleTimerSelector{1, 2, 3, 4}, ScheduleTimerSelector{1, 3, 3, 4}, 2},
		{"schedule timer weekday", ScheduleTimerSelector{1, 2, 3, 4}, ScheduleTimerSelector{1, 2, 4, 4}, 3},
		{"schedule timer index", ScheduleTimerSelector{1, 2, 3, 4}, ScheduleTimerSelector{1, 2, 3, 5}, 4},
		{"schedule write zone", ScheduleTimerWriteSelector{1, 2, 3, 4, 5, 6, 7, 8, 9, 0x1234}, ScheduleTimerWriteSelector{2, 2, 3, 4, 5, 6, 7, 8, 9, 0x1234}, 1},
		{"schedule write circuit", ScheduleTimerWriteSelector{1, 2, 3, 4, 5, 6, 7, 8, 9, 0x1234}, ScheduleTimerWriteSelector{1, 3, 3, 4, 5, 6, 7, 8, 9, 0x1234}, 2},
		{"schedule write weekday", ScheduleTimerWriteSelector{1, 2, 3, 4, 5, 6, 7, 8, 9, 0x1234}, ScheduleTimerWriteSelector{1, 2, 4, 4, 5, 6, 7, 8, 9, 0x1234}, 3},
		{"schedule write index", ScheduleTimerWriteSelector{1, 2, 3, 4, 5, 6, 7, 8, 9, 0x1234}, ScheduleTimerWriteSelector{1, 2, 3, 5, 5, 6, 7, 8, 9, 0x1234}, 4},
		{"schedule write count", ScheduleTimerWriteSelector{1, 2, 3, 4, 5, 6, 7, 8, 9, 0x1234}, ScheduleTimerWriteSelector{1, 2, 3, 4, 6, 6, 7, 8, 9, 0x1234}, 5},
		{"schedule write start hour", ScheduleTimerWriteSelector{1, 2, 3, 4, 5, 6, 7, 8, 9, 0x1234}, ScheduleTimerWriteSelector{1, 2, 3, 4, 5, 7, 7, 8, 9, 0x1234}, 6},
		{"schedule write start minute", ScheduleTimerWriteSelector{1, 2, 3, 4, 5, 6, 7, 8, 9, 0x1234}, ScheduleTimerWriteSelector{1, 2, 3, 4, 5, 6, 8, 8, 9, 0x1234}, 7},
		{"schedule write end hour", ScheduleTimerWriteSelector{1, 2, 3, 4, 5, 6, 7, 8, 9, 0x1234}, ScheduleTimerWriteSelector{1, 2, 3, 4, 5, 6, 7, 9, 9, 0x1234}, 8},
		{"schedule write end minute", ScheduleTimerWriteSelector{1, 2, 3, 4, 5, 6, 7, 8, 9, 0x1234}, ScheduleTimerWriteSelector{1, 2, 3, 4, 5, 6, 7, 8, 10, 0x1234}, 9},
		{"schedule write temperature low", ScheduleTimerWriteSelector{1, 2, 3, 4, 5, 6, 7, 8, 9, 0x1234}, ScheduleTimerWriteSelector{1, 2, 3, 4, 5, 6, 7, 8, 9, 0x1235}, 10},
		{"schedule write temperature high", ScheduleTimerWriteSelector{1, 2, 3, 4, 5, 6, 7, 8, 9, 0x1234}, ScheduleTimerWriteSelector{1, 2, 3, 4, 5, 6, 7, 8, 9, 0x1334}, 11},
		{"diagnostic family", DiagnosticSelector{Family: 1, Selector: 1}, DiagnosticSelector{Family: 0, Selector: 1}, 0},
		{"diagnostic selector", DiagnosticSelector{Family: 1, Selector: 1}, DiagnosticSelector{Family: 1, Selector: 2}, 1},
	}
	for _, c := range cases {
		a, b := c.base.Payload(), c.other.Payload()
		if len(a) != len(b) {
			t.Errorf("%s: lengths %d and %d", c.name, len(a), len(b))
			continue
		}
		for i := range a {
			if (a[i] != b[i]) != (i == c.offset) {
				t.Errorf("%s: byte %d of %x and %x", c.name, i, a, b)
			}
		}
	}
}

func TestVaillantInputSelectorsAbsentAndZeroHistoryIndexAreDistinct(t *testing.T) {
	zero := byte(0)
	absent := DiagnosticSelector{Family: 1, Selector: 1}.Payload()
	explicit := DiagnosticSelector{Family: 1, Selector: 1, Index: &zero}.Payload()
	if hex.EncodeToString(absent) != "0101" || hex.EncodeToString(explicit) != "010100" {
		t.Fatalf("absent %x, explicit zero %x", absent, explicit)
	}
	other := DiagnosticSelector{Family: 1, Selector: 2, Index: &zero}.Payload()
	if hex.EncodeToString(other) != "010200" {
		t.Fatalf("service history index zero %x", other)
	}
}

func TestVaillantInputSelectorsLiveMonitorActionsShareOneBodyAndKeepTheirOwnKey(t *testing.T) {
	want := map[LiveMonitorAction]string{
		LiveMonitorEnable: "live_monitor_enable", LiveMonitorRead: "live_monitor_read", LiveMonitorDisable: "live_monitor_disable",
	}
	for action, key := range want {
		sel := LiveMonitorSelector{Action: action}
		if hex.EncodeToString(sel.Payload()) != "0003" || sel.Key() != key {
			t.Errorf("%s: body %x key %s", action, sel.Payload(), sel.Key())
		}
		if sel.Service() != (Service{Primary: 0xb5, Secondary: 0x03}) {
			t.Errorf("%s: service %+v", action, sel.Service())
		}
	}
}

func TestVaillantInputSelectorsWriteHeadersTakeTheValueTail(t *testing.T) {
	if got := hex.EncodeToString(B524Selector{Opcode: 0x02, Operation: 0x01, Group: 0, Instance: 0, Register: 0x006c}.WithValue([]byte{0x41, 0x00})); got != "02010000"+"6c00"+"4100" {
		t.Errorf("b524 write %s", got)
	}
	if got := hex.EncodeToString(B509Selector{Opcode: 0x0e, Register: 0x0e04}.WithValue([]byte{0x80, 0x01})); got != "0e0e04"+"8001" {
		t.Errorf("b509 write %s", got)
	}
	header := B509Selector{Opcode: 0x0e, Register: 0x0e04}
	_ = header.WithValue([]byte{1})
	if hex.EncodeToString(header.Payload()) != "0e0e04" {
		t.Errorf("payload changed after WithValue: %x", header.Payload())
	}
}

func checkScheduleEntry(t *testing.T, e inventoryEntry) {
	t.Helper()
	types := map[string]Selector{
		"ScheduleConfigSelector": ScheduleConfigSelector{}, "ScheduleCountsSelector": ScheduleCountsSelector{},
		"ScheduleTimerSelector": ScheduleTimerSelector{}, "ScheduleTimerWriteSelector": ScheduleTimerWriteSelector{},
	}
	sel := types[e.Geometry]
	payload := sel.Payload()
	extra := 1
	if e.Geometry == "ScheduleTimerWriteSelector" {
		extra = 2
	}
	if payload[0] != parseByte(t, e.Opcode) || len(payload) != len(e.Fields)+extra {
		t.Errorf("%s: payload %x for opcode %s with %d fields", e.ID, payload, e.Opcode, len(e.Fields))
	}
	if sel.Service() != (Service{Primary: parseByte(t, e.Primary), Secondary: parseByte(t, e.Secondary)}) {
		t.Errorf("%s: service %+v", e.ID, sel.Service())
	}
}
