package output

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"
)

func replayCache(t *testing.T, vector recordedCase, planes replayPlanes) {
	t.Helper()
	if raw, ok := vector.Input["loads"]; ok {
		replayCacheLoads(t, raw, vector.Output.Cache["loads"])
		return
	}
	var snapshot SemanticCacheSnapshot
	if raw, ok := vector.Input["snapshot"]; ok {
		snapshot = decodeRecorded[SemanticCacheSnapshot](t, raw)
	} else {
		input := make(map[string]json.RawMessage)
		for key, value := range vector.Input {
			if key != "overlays" {
				input[key] = value
			}
		}
		replayView(t, input, planes)
		snapshot = SemanticCacheSnapshot{
			Zones: CacheZones(planes.graphql["zones"].([]Zone), replayOverlays(t, vector.Input["overlays"])),
			DHW:   planes.graphql["dhw"].(*DhwStatus), Boiler: planes.graphql["boiler_status"].(*BoilerStatus),
		}
	}
	store, logs := replayStore(t)
	store.now = func() time.Time { return vector.Clock.StoreNow }
	if err := store.Save(snapshot); err != nil {
		t.Fatal(err)
	}
	payload, err := os.ReadFile(store.path)
	if err != nil {
		t.Fatal(err)
	}
	assertRecorded(t, "saved cache bytes", string(payload), vector.Output.Cache["saved_bytes"])
	loaded, ok := store.Load()
	if !ok {
		t.Fatal("saved cache did not load")
	}
	assertRecorded(t, "loaded cache", loaded, vector.Output.Cache["loaded_snapshot"])
	if want, ok := vector.Output.Cache["hydration"]; ok {
		assertRecorded(t, "cache hydration", replayHydration(loaded), want)
		planes.add("zones", loaded.Zones, MCPZonesFromGraphQL(loaded.Zones))
		planes.add("dhw", loaded.DHW, MCPDHWFromGraphQL(loaded.DHW))
		planes.add("boiler_status", loaded.Boiler, MCPBoilerStatusFromGraphQL(loaded.Boiler))
		// The recorder loads independently for hydration and for published planes.
		store.Load()
	}
	if !reflect.DeepEqual(*logs, vector.Output.Observations) {
		t.Errorf("cache observations %q, want %q", *logs, vector.Output.Observations)
	}
}

func replayOverlays(t *testing.T, raw json.RawMessage) map[byte]NativeZoneCacheOverlay {
	t.Helper()
	type overlay struct {
		Temperature *float64 `json:"current_temp_c"`
		Humidity    *float64 `json:"humidity_pct"`
		Mapping     *uint16  `json:"room_mapping_raw"`
		Circuit     *uint16  `json:"associated_circuit_raw"`
	}
	input := decodeRecorded[map[byte]overlay](t, raw)
	result := make(map[byte]NativeZoneCacheOverlay, len(input))
	for key, value := range input {
		result[key] = NativeZoneCacheOverlay{CurrentTempC: value.Temperature, HumidityPct: value.Humidity, RoomMappingRaw: value.Mapping, AssociatedCircuitRaw: value.Circuit}
	}
	return result
}

func replayStore(t *testing.T) (*SemanticCacheStore, *[]string) {
	t.Helper()
	var logs []string
	path := filepath.Join(t.TempDir(), "semantic-cache.json")
	return NewSemanticCacheStore(path, func(format string, args ...any) {
		logs = append(logs, strings.ReplaceAll(fmt.Sprintf(format, args...), path, filepath.Base(path)))
	}), &logs
}

func replayHydration(snapshot SemanticCacheSnapshot) map[string]any {
	h := CacheHydrationFromSnapshot(snapshot)
	return map[string]any{"zones": h.Zones, "dhw": h.Dhw, "boiler": h.Boiler, "persisted_at": h.PersistedAt}
}

func replayCacheLoads(t *testing.T, raw, want json.RawMessage) {
	t.Helper()
	loads := decodeRecorded[[]struct {
		Name        string
		FilePresent bool `json:"file_present"`
		Bytes       string
	}](t, raw)
	results := make([]map[string]any, 0, len(loads))
	for _, load := range loads {
		store, logs := replayStore(t)
		if load.FilePresent {
			if err := os.WriteFile(store.path, []byte(load.Bytes), 0600); err != nil {
				t.Fatal(err)
			}
		}
		snapshot, ok := store.Load()
		results = append(results, map[string]any{
			"name": load.Name, "ok": ok, "snapshot": snapshot,
			"hydration": replayHydration(snapshot), "logs": *logs,
		})
	}
	// Translate only the Go type names that encoding/json prints in the
	// recorded diagnostic. Keep the field, failure reason, and other text exact.
	names := strings.NewReplacer(
		"semanticCacheV3.zones", reflect.TypeFor[semanticCacheFile]().Name()+".zones",
		"[]main.semanticCacheZoneV3", fmt.Sprintf("%T", []semanticCacheZone(nil)),
	)
	want = []byte(names.Replace(string(want)))
	assertRecorded(t, "cache loads", results, want)
}
