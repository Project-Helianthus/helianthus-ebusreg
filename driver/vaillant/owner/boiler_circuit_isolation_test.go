package owner

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Project-Helianthus/helianthus-ebusgo/protocol"
	"github.com/Project-Helianthus/helianthus-ebusreg/driver/vaillant/output"
	"os"
	"path/filepath"
	"testing"
)

func TestVaillantOwnerBoilerCircuitIsolationBoilerUnavailableLeavesBoilerFieldsUnknownAndPublishesCircuit(t *testing.T) {
	for _, boiler := range []struct {
		name    string
		address byte
	}{
		{name: "no_boiler_address", address: 0},
		{name: "boiler_reads_fail", address: boilerCircuitIsolationBoilerAddress},
	} {
		for _, circuitRaw := range []uint16{1, 0} {
			name := boiler.name + "/circuit_pump_on"
			stateRaw := uint16(1)
			if circuitRaw == 0 {
				name = boiler.name + "/circuit_pump_off"
				stateRaw = 2
			}
			t.Run(name, func(t *testing.T) {
				poller := boilerCircuitIsolationPoller(boiler.address, circuitRaw, stateRaw)
				var boilerReads int
				poller.sendFrameFn = func(_ context.Context, frame protocol.Frame) (*protocol.Frame, error) {
					if frame.Primary != vaillantB509Primary || frame.Secondary != vaillantB509Secondary {
						t.Fatalf("unexpected non-B509 active read: % x", frame.Data)
					}
					boilerReads++
					return nil, errors.New("boiler B509 unavailable")
				}

				poller.refreshBoilerStatusTier(context.Background(), boilerStatusTierFast)
				poller.publishCircuits(semanticSnapshotSourceLive)

				if boiler.address == 0 && boilerReads != 0 {
					t.Fatalf("B509 reads = %d without a boiler address; want none", boilerReads)
				}
				if boiler.address != 0 && boilerReads == 0 {
					t.Fatal("B509 reads = 0; want the failing boiler path exercised")
				}
				boilerCircuitIsolationAssertSurfaces(t, poller, nil, circuitRaw, stateRaw)
				if got := poller.circuits[0].PumpStatusRaw; got == nil || *got != circuitRaw {
					t.Fatalf("native circuit pump = %v; want preserved %d", got, circuitRaw)
				}
				if got := poller.circuits[0].CircuitStateRaw; got == nil || *got != stateRaw {
					t.Fatalf("native circuit state = %v; want preserved %d", got, stateRaw)
				}
			})
		}
	}
}

func TestVaillantOwnerBoilerCircuitIsolationBoilerPumpComesOnlyFromB509(t *testing.T) {
	for _, boilerByte := range []byte{0xF0, 0x01} {
		boilerActive, ok := decodeOnOff([]byte{boilerByte})
		if !ok {
			t.Fatalf("decodeOnOff(0x%02x) failed", boilerByte)
		}
		for _, circuitRaw := range []uint16{0, 1} {
			name := "boiler_off"
			if boilerActive {
				name = "boiler_on"
			}
			stateRaw := uint16(3)
			if circuitRaw == 0 {
				name += "/circuit_pump_off"
			} else {
				name += "/circuit_pump_on"
				stateRaw = 1
			}
			t.Run(name, func(t *testing.T) {
				poller := boilerCircuitIsolationPoller(boilerCircuitIsolationBoilerAddress, circuitRaw, stateRaw)
				var pumpReads int
				poller.sendFrameFn = func(_ context.Context, frame protocol.Frame) (*protocol.Frame, error) {
					if frame.Primary != vaillantB509Primary || frame.Secondary != vaillantB509Secondary {
						t.Fatalf("unexpected non-B509 active read: % x", frame.Data)
					}
					if len(frame.Data) < 3 || frame.Data[0] != vaillantB509OpcodeRead {
						t.Fatalf("unexpected B509 request % x", frame.Data)
					}
					addr := uint16(frame.Data[1])<<8 | uint16(frame.Data[2])
					value := []byte{0x00, 0x00}
					if addr == boilerB509CentralHeatingPump {
						pumpReads++
						value = []byte{boilerByte}
					}
					return &protocol.Frame{Data: append([]byte{vaillantB509OpcodeRead, frame.Data[1], frame.Data[2]}, value...)}, nil
				}

				poller.refreshBoilerStatusTier(context.Background(), boilerStatusTierFast)
				poller.publishCircuits(semanticSnapshotSourceLive)

				if pumpReads == 0 {
					t.Fatal("boiler pump B509 read not issued")
				}
				boilerCircuitIsolationAssertSurfaces(t, poller, &boilerActive, circuitRaw, stateRaw)
			})
		}
	}
}

func TestVaillantOwnerBoilerCircuitIsolationCachedHeatingStatusRawWithoutBoilerSourceIsDiscarded(t *testing.T) {
	for _, cachedRaw := range []int{0, 1, 3} {
		t.Run(DecodeCircuitStateToken(boilerCircuitIsolationUint16(uint16(cachedRaw))), func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "semantic-cache.json")
			dhwTemp := 48.5
			body := []byte(fmt.Sprintf("{\"schema_version\":3,\"metadata\":{\"persisted_at\":\"2026-01-01T00:00:00Z\"},\"boiler\":{\"dhw_temperature_c\":48.5,\"heating_status_raw\":%d}}\n", cachedRaw))
			if err := os.WriteFile(path, body, 0o644); err != nil {
				t.Fatal(err)
			}
			store := output.NewSemanticCacheStore(path, t.Logf)

			poller := boilerCircuitIsolationPoller(0, 1, uint16(cachedRaw))
			poller.publication = newVaillantSemanticPublication(newRetainedSemanticState(), publishedHub(poller), publishedCache(poller))
			snapshot, ok := preloadSemanticCache(publishedProvider(poller), store)
			if !ok {
				t.Fatal("legacy cache was discarded wholesale")
			}
			status := publishedProvider(poller).BoilerStatus()
			if status == nil {
				t.Fatal("cached boiler status missing; want other cached fields restored")
			}
			if got := status.State.DhwTemperatureC; got == nil || *got != dhwTemp {
				t.Fatalf("cached DHW temperature = %v; want restored %.1f", got, dhwTemp)
			}
			if got := status.Diagnostics.HeatingStatusRaw; got != nil {
				t.Fatalf("cached heating status = %d; want unset without a boiler-qualified source", *got)
			}
			poller.HydrateFromCache(snapshot)
			if poller.boiler == nil || poller.boiler.HeatingStatusRaw != nil {
				t.Fatal("hydrated boiler heating status must be unset")
			}

			legacy := cachedRaw
			if err := store.Save(output.SemanticCacheSnapshot{Boiler: &output.BoilerStatus{
				State:       output.BoilerState{DhwTemperatureC: &dhwTemp},
				Diagnostics: output.BoilerDiagnostics{HeatingStatusRaw: &legacy},
			}}); err != nil {
				t.Fatal(err)
			}
			written, err := os.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}
			var persisted struct {
				Boiler *struct {
					DhwTemperatureC  *float64 `json:"dhw_temperature_c"`
					HeatingStatusRaw *int     `json:"heating_status_raw"`
				} `json:"boiler"`
			}
			if err := json.Unmarshal(written, &persisted); err != nil {
				t.Fatal(err)
			}
			if persisted.Boiler == nil || persisted.Boiler.DhwTemperatureC == nil || *persisted.Boiler.DhwTemperatureC != dhwTemp {
				t.Fatalf("persisted boiler = %+v; want other boiler fields written", persisted.Boiler)
			}
			if persisted.Boiler.HeatingStatusRaw != nil {
				t.Fatalf("persisted heating status = %d; want not written without a boiler source", *persisted.Boiler.HeatingStatusRaw)
			}
			cached, ok := store.Load()
			if !ok || cached.Boiler == nil {
				t.Fatal("saved cache did not load")
			}
			if cached.Boiler.Diagnostics.HeatingStatusRaw != nil {
				t.Fatal("cache must not persist the boiler heating status without a boiler source")
			}
		})
	}
}

func TestVaillantOwnerBoilerCircuitIsolationBoilerRegisterDefinitionsNeverReadCircuitPumpOrState(t *testing.T) {
	t.Parallel()

	for _, tier := range []boilerStatusTier{
		boilerStatusTierFast,
		boilerStatusTierMedium,
		boilerStatusTierSlow,
	} {
		for _, register := range boilerStatusRegisterDefinitionsForTier(tier) {
			if register.group != localCircuits.group {
				continue
			}
			switch register.addr {
			case circuitPumpStatus, circuitCircuitState:
				t.Errorf("tier %v maps circuit GG=0x%02x II=0x%02x RR=0x%04x into the boiler status; circuit registers belong to the circuit surface only",
					tier, register.group, register.instance, register.addr)
			}
		}
	}
}
