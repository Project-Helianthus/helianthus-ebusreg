package output

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestSemanticCacheStoreLoad_DiscardsV1(t *testing.T) {
	logs := &semanticCacheTestLogs{}
	path := filepath.Join(t.TempDir(), "semantic-cache.json")
	legacyV1 := `{
  "zones": [
    {"id":"zone-2","name":"Zone 2"},
    {"id":"zone-1","name":"Zone 1","allowed_modes":["heat","auto"]}
  ],
  "dhw": {"operating_mode":"auto","preset":"schedule"}
}`
	if err := os.WriteFile(path, []byte(legacyV1), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	store := NewSemanticCacheStore(path, logs.Printf)
	if _, ok := store.Load(); ok {
		t.Fatalf("Load() ok = true; want false for discarded V1 cache")
	}
	logs.requireContains(t, "semantic_cache_discarded")
}

func TestSemanticCacheStoreLoad_DiscardsV2(t *testing.T) {
	logs := &semanticCacheTestLogs{}
	path := filepath.Join(t.TempDir(), "semantic-cache.json")
	legacyV2 := `{"schema_version":2,"metadata":{},"zones":[{"id":"zone-1","name":"Zone 1"}]}`
	if err := os.WriteFile(path, []byte(legacyV2), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	store := NewSemanticCacheStore(path, logs.Printf)
	if _, ok := store.Load(); ok {
		t.Fatalf("Load() ok = true; want false for discarded V2 cache")
	}
	logs.requireContains(t, "semantic_cache_discarded")
}

func TestSemanticCacheStoreLoad_HitLogsMarker(t *testing.T) {
	logs := &semanticCacheTestLogs{}
	path := filepath.Join(t.TempDir(), "semantic-cache.json")
	store := NewSemanticCacheStore(path, logs.Printf)

	targetTemp := 22.5
	if err := store.Save(SemanticCacheSnapshot{
		Zones: []Zone{
			{ID: "zone-1", Name: "Zone 1", Config: ZoneConfig{TargetTempC: &targetTemp}},
		},
	}); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	logs.Reset()
	snapshot, ok := store.Load()
	if !ok {
		t.Fatalf("Load() ok = false; want true")
	}
	if len(snapshot.Zones) != 1 {
		t.Fatalf("len(snapshot.Zones) = %d; want 1", len(snapshot.Zones))
	}
	if snapshot.Zones[0].ID != "zone-1" {
		t.Fatalf("snapshot.Zones[0].ID = %q; want zone-1", snapshot.Zones[0].ID)
	}
	logs.requireContains(t, "semantic_cache_hit")
}

func TestSemanticCacheStoreLoad_InvalidJSONLogsMarker(t *testing.T) {
	logs := &semanticCacheTestLogs{}
	path := filepath.Join(t.TempDir(), "semantic-cache.json")
	if err := os.WriteFile(path, []byte("{not-json"), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	store := NewSemanticCacheStore(path, logs.Printf)
	if _, ok := store.Load(); ok {
		t.Fatalf("Load() ok = true; want false for invalid payload")
	}
	logs.requireContains(t, "semantic_cache_invalid")
}

func TestSemanticCacheStoreLoad_MissLogsMarker(t *testing.T) {
	logs := &semanticCacheTestLogs{}
	path := filepath.Join(t.TempDir(), "semantic-cache.json")
	store := NewSemanticCacheStore(path, logs.Printf)

	if _, ok := store.Load(); ok {
		t.Fatalf("Load() ok = true; want false for cache miss")
	}
	logs.requireContains(t, "semantic_cache_miss")
}

func TestSemanticCacheStoreLoad_UnknownSchemaVersionLogsInvalid(t *testing.T) {
	logs := &semanticCacheTestLogs{}
	path := filepath.Join(t.TempDir(), "semantic-cache.json")
	if err := os.WriteFile(path, []byte(`{"schema_version":99,"zones":[]}`), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	store := NewSemanticCacheStore(path, logs.Printf)
	if _, ok := store.Load(); ok {
		t.Fatalf("Load() ok = true; want false for unknown schema")
	}
	logs.requireContains(t, "semantic_cache_invalid")
	logs.requireContains(t, "unknown_schema_version")
}

func TestSemanticCacheStoreSave_SortsZonesAndAllowedModes(t *testing.T) {
	logs := &semanticCacheTestLogs{}
	path := filepath.Join(t.TempDir(), "semantic-cache.json")
	store := NewSemanticCacheStore(path, logs.Printf)

	if err := store.Save(SemanticCacheSnapshot{
		Zones: []Zone{
			{ID: "zone-10", Name: "Zone 10", Config: ZoneConfig{AllowedModes: []string{"heat", "off"}}},
			{ID: "zone-2", Name: "Zone 2", Config: ZoneConfig{AllowedModes: []string{"off", "heat"}}},
			{ID: "zone-1", Name: "Zone 1", Config: ZoneConfig{AllowedModes: []string{"heat", "auto"}}},
		},
	}); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	payload, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}
	var persisted semanticCacheV3
	if err := json.Unmarshal(payload, &persisted); err != nil {
		t.Fatalf("Unmarshal() error = %v", err)
	}
	if len(persisted.Zones) != 3 {
		t.Fatalf("len(persisted.Zones) = %d; want 3", len(persisted.Zones))
	}
	if persisted.Zones[0].ID != "zone-1" || persisted.Zones[1].ID != "zone-2" || persisted.Zones[2].ID != "zone-10" {
		t.Fatalf("zone order = [%s %s %s]; want [zone-1 zone-2 zone-10]", persisted.Zones[0].ID, persisted.Zones[1].ID, persisted.Zones[2].ID)
	}
	if got := strings.Join(persisted.Zones[0].Config.AllowedModes, ","); got != "auto,heat" {
		t.Fatalf("allowed_modes zone-1 = %q; want auto,heat", got)
	}
	if got := strings.Join(persisted.Zones[1].Config.AllowedModes, ","); got != "heat,off" {
		t.Fatalf("allowed_modes zone-2 = %q; want heat,off", got)
	}
	if got := strings.Join(persisted.Zones[2].Config.AllowedModes, ","); got != "heat,off" {
		t.Fatalf("allowed_modes zone-10 = %q; want heat,off", got)
	}
}

func TestSemanticCacheStoreSave_WriteFailureKeepsExistingFile(t *testing.T) {
	logs := &semanticCacheTestLogs{}
	baseDir := t.TempDir()
	path := filepath.Join(baseDir, "semantic-cache.json")
	original := []byte(`{"schema_version":2,"metadata":{"persisted_at":"2026-01-01T00:00:00Z"},"zones":[{"id":"zone-old","name":"Old"}]}` + "\n")
	if err := os.WriteFile(path, original, 0o644); err != nil {
		t.Fatalf("WriteFile(original) error = %v", err)
	}

	if err := os.Chmod(baseDir, 0o500); err != nil {
		t.Fatalf("Chmod(baseDir) error = %v", err)
	}
	defer func() {
		_ = os.Chmod(baseDir, 0o755)
	}()

	store := NewSemanticCacheStore(path, logs.Printf)
	err := store.Save(SemanticCacheSnapshot{
		Zones: []Zone{{ID: "zone-1", Name: "Zone 1"}},
	})
	if err == nil {
		t.Fatalf("Save() error = nil; want failure when target dir is readonly")
	}

	if err := os.Chmod(baseDir, 0o755); err != nil {
		t.Fatalf("Chmod(baseDir restore) error = %v", err)
	}
	current, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile(current) error = %v", err)
	}
	if string(current) != string(original) {
		t.Fatalf("cache file changed on failed write; got %q want %q", string(current), string(original))
	}
	logs.requireContains(t, "semantic_cache_write_error")
}

func TestSemanticCacheStoreSave_WritePermissionErrorFailsSafe(t *testing.T) {
	logs := &semanticCacheTestLogs{}
	baseDir := t.TempDir()
	readonlyDir := filepath.Join(baseDir, "readonly")
	if err := os.Mkdir(readonlyDir, 0o755); err != nil {
		t.Fatalf("Mkdir() error = %v", err)
	}
	if err := os.Chmod(readonlyDir, 0o500); err != nil {
		t.Fatalf("Chmod() error = %v", err)
	}
	defer func() {
		_ = os.Chmod(readonlyDir, 0o755)
	}()

	path := filepath.Join(readonlyDir, "semantic-cache.json")
	store := NewSemanticCacheStore(path, logs.Printf)
	err := store.Save(SemanticCacheSnapshot{
		Zones: []Zone{{ID: "zone-1", Name: "Zone 1"}},
	})
	if err == nil {
		t.Fatalf("Save() error = nil; want permission/write failure")
	}
	logs.requireContains(t, "semantic_cache_write_error")
}

type semanticCacheTestLogs struct {
	entries []string
}

func (logs *semanticCacheTestLogs) Printf(format string, args ...any) {
	logs.entries = append(logs.entries, fmt.Sprintf(format, args...))
}

func (logs *semanticCacheTestLogs) Reset() {
	logs.entries = nil
}

func (logs *semanticCacheTestLogs) requireContains(t *testing.T, fragment string) {
	t.Helper()
	for _, entry := range logs.entries {
		if strings.Contains(entry, fragment) {
			return
		}
	}
	t.Fatalf("logs do not contain %q; logs=%v", fragment, logs.entries)
}

func TestCacheHydrationFromSnapshotKeysZonesAndKeepsLegacyReturn(t *testing.T) {
	ret, current := 33.0, 44.0
	persisted := time.Date(2026, time.March, 1, 6, 0, 0, 0, time.FixedZone("east", 2*3600))
	hydration := CacheHydrationFromSnapshot(SemanticCacheSnapshot{
		Zones: []Zone{
			{ID: "zone-3", Name: "upstairs"},
			{ID: "custom", Name: "fallback"},
		},
		DHW:         &DhwStatus{State: DhwState{CurrentTempC: &current}},
		Boiler:      &BoilerStatus{State: BoilerState{ReturnTemperatureC: &ret}},
		PersistedAt: persisted,
	})
	if len(hydration.Zones) != 2 || hydration.Zones[0].Instance != 0x02 || hydration.Zones[1].Instance != 0x01 {
		t.Fatalf("hydrated zones = %+v; want instances 2 then fallback 1", hydration.Zones)
	}
	if hydration.Zones[0].Name != "upstairs" || hydration.Zones[1].Name != "fallback" {
		t.Fatal("hydrated zones lost cache order")
	}
	if hydration.Dhw == nil || hydration.Dhw.CurrentTempC == nil || *hydration.Dhw.CurrentTempC != current {
		t.Fatalf("hydrated hot water = %+v; want cached current temperature", hydration.Dhw)
	}
	if !hydration.PersistedAt.Equal(persisted) {
		t.Fatalf("persisted time = %v; want %v", hydration.PersistedAt, persisted)
	}
	if hydration.Boiler == nil || hydration.Boiler.ReturnTemperatureC == nil || *hydration.Boiler.ReturnTemperatureC != ret {
		t.Fatalf("hydrated boiler = %+v; want cached legacy return", hydration.Boiler)
	}

	empty := CacheHydrationFromSnapshot(SemanticCacheSnapshot{})
	if empty.Zones != nil || empty.Dhw != nil || empty.Boiler != nil || !empty.PersistedAt.IsZero() {
		t.Fatalf("empty cache hydration = %+v; want the zero value", empty)
	}
}

func TestPersistSemanticSnapshotLaysRetainedZoneValuesOverPublishedZones(t *testing.T) {
	publishedTemp, publishedHumidity := 18.0, 35.0
	publishedMapping, publishedCircuit := 7, 6
	published := func(id string) Zone {
		temp, humidity := publishedTemp, publishedHumidity
		mapping, circuit := publishedMapping, publishedCircuit
		return Zone{
			ID:     id,
			Name:   id,
			State:  ZoneState{CurrentTempC: &temp, CurrentHumidityPct: &humidity},
			Config: ZoneConfig{OperatingMode: "auto", RoomTemperatureZoneMapping: &mapping, AssociatedCircuit: &circuit},
		}
	}
	retainedTemp, retainedHumidity := 21.5, 48.0
	retainedMapping, retainedCircuit := uint16(2), uint16(1)
	zones := CacheZones([]Zone{published("zone-1"), published("zone-2")}, map[byte]NativeZoneCacheOverlay{
		0x01: {CurrentTempC: &retainedTemp, HumidityPct: &retainedHumidity, RoomMappingRaw: &retainedMapping, AssociatedCircuitRaw: &retainedCircuit},
	})
	if len(zones) != 2 {
		t.Fatalf("zones=%d; want two zones", len(zones))
	}
	untouched, overlaid := zones[0], zones[1]
	if untouched.ID != "zone-1" || overlaid.ID != "zone-2" {
		t.Fatalf("saved zone order = %q, %q; want published order", untouched.ID, overlaid.ID)
	}
	if *untouched.State.CurrentTempC != publishedTemp || *untouched.State.CurrentHumidityPct != publishedHumidity ||
		*untouched.Config.RoomTemperatureZoneMapping != publishedMapping || *untouched.Config.AssociatedCircuit != publishedCircuit {
		t.Fatalf("zone without a retained entry changed: %+v", untouched)
	}
	if got := overlaid.State.CurrentTempC; got == nil || *got != retainedTemp {
		t.Fatalf("saved current temperature = %v; want retained %.1f", got, retainedTemp)
	}
	if got := overlaid.State.CurrentHumidityPct; got == nil || *got != retainedHumidity {
		t.Fatalf("saved humidity = %v; want retained %.1f", got, retainedHumidity)
	}
	if got := overlaid.Config.RoomTemperatureZoneMapping; got == nil || *got != int(retainedMapping) {
		t.Fatalf("saved room temperature mapping = %v; want retained %d", got, retainedMapping)
	}
	if got := overlaid.Config.AssociatedCircuit; got == nil || *got != int(retainedCircuit) {
		t.Fatalf("saved associated circuit = %v; want retained %d", got, retainedCircuit)
	}
	if overlaid.State.CurrentTempC == &retainedTemp || overlaid.State.CurrentHumidityPct == &retainedHumidity {
		t.Fatal("saved zone shares a retained value")
	}
	if overlaid.Config.OperatingMode != "auto" || overlaid.Name != "zone-2" {
		t.Fatalf("saved zone lost a published field: %+v", overlaid)
	}
}
