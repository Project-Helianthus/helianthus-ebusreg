package owner

import (
	"context"
	"fmt"
	"github.com/Project-Helianthus/helianthus-ebusreg/registry"
	"github.com/Project-Helianthus/helianthus-ebusreg/vaillant/productids"
	"log"
	"strings"
	"sync"
	"testing"
	"time"
)

// TestVaillantSemanticPollerAbsenceStartsUnknown covers the initial absence
// state. The constructor starts in UNKNOWN (not yet evaluated) and publishes it;
// the first regulator recheck moves it to PRESENT when the regulator is
// classified, otherwise to ABSENCE_GRACE with the grace measured from that
// recheck, so ABSENT falls on the same recheck as when the state started at
// PRESENT. An unclassified controller (no serial) follows the not-present path.
func TestVaillantOwnerVaillantSemanticPollerAbsenceStartsUnknown(t *testing.T) {
	catalog, err := productids.LoadCatalog()
	if err != nil {
		t.Fatalf("LoadCatalog() error = %v", err)
	}
	boiler := registry.DeviceInfo{Address: 0x08, Manufacturer: "Vaillant", DeviceID: "BAI00", SerialNumber: controllerlessSerial(controllerlessCatalogPart(t, "Boiler"))}
	regulator := registry.DeviceInfo{Address: 0x15, Manufacturer: "Vaillant", SerialNumber: controllerlessSerial(controllerlessCatalogPartWith(t, "Regulator", productids.ControllerPresent))}
	unclassified := registry.DeviceInfo{Address: 0x15, Manufacturer: "Vaillant"}
	for _, tc := range []struct {
		name    string
		devices []registry.DeviceInfo
		// states at the rechecks after 1, 2, 6 and 7 minutes of polling
		want []regulatorAbsenceState
	}{
		{name: "classified regulator", devices: []registry.DeviceInfo{boiler, regulator},
			want: []regulatorAbsenceState{regulatorPresent, regulatorPresent, regulatorPresent, regulatorPresent}},
		{name: "no regulator", devices: []registry.DeviceInfo{boiler},
			want: []regulatorAbsenceState{regulatorAbsenceGrace, regulatorAbsenceGrace, regulatorAbsent, regulatorAbsent}},
		{name: "unclassified controller", devices: []registry.DeviceInfo{boiler, unclassified},
			want: []regulatorAbsenceState{regulatorAbsenceGrace, regulatorAbsenceGrace, regulatorAbsent, regulatorAbsent}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var mu sync.Mutex
			var out strings.Builder
			original := log.Writer()
			log.SetOutput(fallbackLogWriter{mu: &mu, out: &out})
			t.Cleanup(func() { log.SetOutput(original) })

			cfg := donorDefaultConfig()
			poller := newVaillantSemanticPollerFromDonor(cfg, &donorGateway{Registry: newTestRegistry(tc.devices...)}, nil, nil, nil)
			poller.catalog, poller.catalogErr = catalog, nil
			if got := poller.regAbsenceState; got != regulatorAbsenceUnknown {
				t.Fatalf("initial absence state = %s; want UNKNOWN", got)
			}
			if got := semanticRegulatorState.Value(); got != string(regulatorAbsenceUnknown) {
				t.Fatalf("semantic_regulator_state at construction = %q; want UNKNOWN", got)
			}
			start := time.Unix(0, 0)
			now := start
			poller.nowFn = func() time.Time { return now }
			for i, minute := range []int{1, 2, 6, 7} {
				now = start.Add(time.Duration(minute) * time.Minute)
				poller.refreshRegulatorCapability(context.Background())
				if got := poller.regAbsenceState; got != tc.want[i] {
					t.Fatalf("absence state at the %d min recheck = %s; want %s", minute, got, tc.want[i])
				}
			}
			mu.Lock()
			logs := out.String()
			mu.Unlock()
			first := fmt.Sprintf("semantic_regulator_transition from=UNKNOWN to=%s", tc.want[0])
			if !strings.Contains(logs, first) || strings.Contains(logs, "from=PRESENT to=ABSENCE_GRACE") {
				t.Fatalf("transition log = %q; want %q and no PRESENT->ABSENCE_GRACE", logs, first)
			}
		})
	}
}
