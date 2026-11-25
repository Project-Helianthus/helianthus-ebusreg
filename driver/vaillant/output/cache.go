package output

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"time"
)

const (
	semanticCacheSchemaVersionV1 = 1
	semanticCacheSchemaVersionV2 = 2
	semanticCacheSchemaVersionV3 = 3
)

// SemanticCacheSnapshot represents the records eligible for persistence and the recorded cache time.
type SemanticCacheSnapshot struct {
	Zones       []Zone
	DHW         *DhwStatus
	Boiler      *BoilerStatus
	PersistedAt time.Time
}

// SemanticCachePersister represents the cache write operation used by the owner.
type SemanticCachePersister interface {
	Save(snapshot SemanticCacheSnapshot) error
}

// SemanticCacheStore represents the versioned cache file at a caller-supplied path.
type SemanticCacheStore struct {
	path string
	logf func(string, ...any)
	now  func() time.Time
}

type semanticCacheEnvelope struct {
	SchemaVersion int `json:"schema_version"`
}

type semanticCacheV3 struct {
	SchemaVersion int                   `json:"schema_version"`
	Metadata      semanticCacheMetadata `json:"metadata"`
	Zones         []semanticCacheZoneV3 `json:"zones,omitempty"`
	DHW           *semanticCacheDHWV3   `json:"dhw,omitempty"`
	Boiler        *semanticCacheBoiler  `json:"boiler,omitempty"`
}

type semanticCacheMetadata struct {
	PersistedAt  time.Time `json:"persisted_at,omitempty"`
	MigratedFrom int       `json:"migrated_from,omitempty"`
}

type semanticCacheZoneV3 struct {
	ID     string                  `json:"id"`
	Name   string                  `json:"name"`
	State  semanticCacheZoneState  `json:"state"`
	Config semanticCacheZoneConfig `json:"config"`
}

type semanticCacheZoneState struct {
	CurrentTempC       *float64 `json:"current_temp_c,omitempty"`
	CurrentHumidityPct *float64 `json:"current_humidity_pct,omitempty"`
	HvacAction         string   `json:"hvac_action,omitempty"`
	SpecialFunction    string   `json:"special_function,omitempty"`
	HeatingDemandPct   *float64 `json:"heating_demand_pct,omitempty"`
	ValvePositionPct   *float64 `json:"valve_position_pct,omitempty"`
}

type semanticCacheZoneConfig struct {
	OperatingMode              string   `json:"operating_mode,omitempty"`
	Preset                     string   `json:"preset,omitempty"`
	TargetTempC                *float64 `json:"target_temp_c,omitempty"`
	AllowedModes               []string `json:"allowed_modes,omitempty"`
	CircuitType                string   `json:"circuit_type,omitempty"`
	AssociatedCircuit          *int     `json:"associated_circuit,omitempty"`
	RoomTemperatureZoneMapping *int     `json:"room_temperature_zone_mapping,omitempty"`
}

type semanticCacheDHWV3 struct {
	State  semanticCacheDhwState  `json:"state"`
	Config semanticCacheDhwConfig `json:"config"`
}

type semanticCacheDhwState struct {
	CurrentTempC     *float64 `json:"current_temp_c,omitempty"`
	SpecialFunction  string   `json:"special_function,omitempty"`
	HeatingDemandPct *float64 `json:"heating_demand_pct,omitempty"`
}

type semanticCacheDhwConfig struct {
	OperatingMode string   `json:"operating_mode,omitempty"`
	Preset        string   `json:"preset,omitempty"`
	TargetTempC   *float64 `json:"target_temp_c,omitempty"`
}

type semanticCacheBoiler struct {
	FlowTemperatureC         *float64 `json:"flow_temperature_c,omitempty"`
	ReturnTemperatureC       *float64 `json:"return_temperature_c,omitempty"`
	CentralHeatingPumpActive *bool    `json:"central_heating_pump_active,omitempty"`
	DhwTemperatureC          *float64 `json:"dhw_temperature_c,omitempty"`
	DhwTargetTemperatureC    *float64 `json:"dhw_target_temperature_c,omitempty"`
	DhwOperatingMode         *string  `json:"dhw_operating_mode,omitempty"`
	HeatingStatusRaw         *int     `json:"heating_status_raw,omitempty"`
	DhwStatusRaw             *int     `json:"dhw_status_raw,omitempty"`
}

// NewSemanticCacheStore opens a handle to the requested cache path without reading it.
func NewSemanticCacheStore(path string, logf func(string, ...any)) *SemanticCacheStore {
	path = strings.TrimSpace(path)
	if path == "" {
		return nil
	}
	if logf == nil {
		logf = log.Printf
	}
	return &SemanticCacheStore{
		path: path,
		logf: logf,
		now:  time.Now,
	}
}

// Load loads a supported cache snapshot and reports whether it was available.
func (store *SemanticCacheStore) Load() (SemanticCacheSnapshot, bool) {
	if store == nil || store.path == "" {
		return SemanticCacheSnapshot{}, false
	}

	payload, err := os.ReadFile(store.path)
	if err != nil {
		if os.IsNotExist(err) {
			store.logf("semantic_cache_miss path=%q", store.path)
			return SemanticCacheSnapshot{}, false
		}
		store.logf("semantic_cache_invalid path=%q reason=read_error err=%v", store.path, err)
		return SemanticCacheSnapshot{}, false
	}

	var envelope semanticCacheEnvelope
	if err := json.Unmarshal(payload, &envelope); err != nil {
		store.logf("semantic_cache_invalid path=%q reason=json_unmarshal err=%v", store.path, err)
		return SemanticCacheSnapshot{}, false
	}

	switch envelope.SchemaVersion {
	case 0, semanticCacheSchemaVersionV1, semanticCacheSchemaVersionV2:

		store.logf("semantic_cache_discarded path=%q schema_version=%d reason=outdated_schema", store.path, envelope.SchemaVersion)
		return SemanticCacheSnapshot{}, false
	case semanticCacheSchemaVersionV3:
		var cacheV3 semanticCacheV3
		if err := json.Unmarshal(payload, &cacheV3); err != nil {
			store.logf("semantic_cache_invalid path=%q reason=v3_unmarshal err=%v", store.path, err)
			return SemanticCacheSnapshot{}, false
		}
		snapshot := semanticCacheV3ToSnapshot(cacheV3)
		store.logf("semantic_cache_hit path=%q schema_version=%d zones=%d dhw=%t", store.path, semanticCacheSchemaVersionV3, len(snapshot.Zones), snapshot.DHW != nil)
		return snapshot, true
	default:
		store.logf("semantic_cache_invalid path=%q reason=unknown_schema_version schema_version=%d", store.path, envelope.SchemaVersion)
		return SemanticCacheSnapshot{}, false
	}
}

// Save atomically saves the snapshot with the current store clock.
func (store *SemanticCacheStore) Save(snapshot SemanticCacheSnapshot) error {
	if store == nil || store.path == "" {
		return nil
	}
	cacheV3 := semanticCacheSnapshotToV3(snapshot, store.now().UTC())
	return store.writeV3(cacheV3)
}

func (store *SemanticCacheStore) writeV3(cacheV3 semanticCacheV3) error {
	if store == nil || store.path == "" {
		return nil
	}
	payload, err := json.MarshalIndent(cacheV3, "", "  ")
	if err != nil {
		store.logf("semantic_cache_write_error path=%q reason=marshal err=%v", store.path, err)
		return err
	}
	payload = append(payload, '\n')
	if err := writeSemanticCacheFileAtomic(store.path, payload); err != nil {
		store.logf("semantic_cache_write_error path=%q reason=atomic_write err=%v", store.path, err)
		return err
	}
	store.logf(
		"semantic_cache_write path=%q schema_version=%d zones=%d dhw=%t",
		store.path,
		cacheV3.SchemaVersion,
		len(cacheV3.Zones),
		cacheV3.DHW != nil,
	)
	return nil
}

func semanticCacheV3ToSnapshot(cacheV3 semanticCacheV3) SemanticCacheSnapshot {
	out := SemanticCacheSnapshot{
		Zones:       make([]Zone, 0, len(cacheV3.Zones)),
		DHW:         nil,
		PersistedAt: cacheV3.Metadata.PersistedAt,
	}
	for _, zone := range normalizeSemanticCacheZonesV3(cacheV3.Zones) {
		out.Zones = append(out.Zones, Zone{
			ID:   zone.ID,
			Name: zone.Name,
			State: ZoneState{
				CurrentTempC:       CloneFloatPtr(zone.State.CurrentTempC),
				CurrentHumidityPct: CloneFloatPtr(zone.State.CurrentHumidityPct),
				HvacAction:         zone.State.HvacAction,
				SpecialFunction:    zone.State.SpecialFunction,
				HeatingDemandPct:   CloneFloatPtr(zone.State.HeatingDemandPct),
				ValvePositionPct:   CloneFloatPtr(zone.State.ValvePositionPct),
			},
			Config: ZoneConfig{
				OperatingMode:              zone.Config.OperatingMode,
				Preset:                     zone.Config.Preset,
				TargetTempC:                CloneFloatPtr(zone.Config.TargetTempC),
				AllowedModes:               append([]string(nil), zone.Config.AllowedModes...),
				CircuitType:                zone.Config.CircuitType,
				AssociatedCircuit:          CloneIntPtr(zone.Config.AssociatedCircuit),
				RoomTemperatureZoneMapping: CloneIntPtr(zone.Config.RoomTemperatureZoneMapping),
			},
		})
	}
	if cacheV3.DHW != nil {
		out.DHW = &DhwStatus{
			State: DhwState{
				CurrentTempC:     CloneFloatPtr(cacheV3.DHW.State.CurrentTempC),
				SpecialFunction:  cacheV3.DHW.State.SpecialFunction,
				HeatingDemandPct: CloneFloatPtr(cacheV3.DHW.State.HeatingDemandPct),
			},
			Config: DhwConfig{
				OperatingMode: cacheV3.DHW.Config.OperatingMode,
				Preset:        cacheV3.DHW.Config.Preset,
				TargetTempC:   CloneFloatPtr(cacheV3.DHW.Config.TargetTempC),
			},
		}
	}
	if cacheV3.Boiler != nil {
		out.Boiler = cacheBoilerToGraphQL(cacheV3.Boiler)
	}
	return out
}

func semanticCacheSnapshotToV3(snapshot SemanticCacheSnapshot, persistedAt time.Time) semanticCacheV3 {
	normalizedSnapshot := normalizeSemanticCacheSnapshot(snapshot)
	cacheV3 := semanticCacheV3{
		SchemaVersion: semanticCacheSchemaVersionV3,
		Metadata: semanticCacheMetadata{
			PersistedAt: persistedAt,
		},
		Zones: make([]semanticCacheZoneV3, 0, len(normalizedSnapshot.Zones)),
	}
	for _, zone := range normalizedSnapshot.Zones {
		cacheV3.Zones = append(cacheV3.Zones, semanticCacheZoneV3{
			ID:   zone.ID,
			Name: zone.Name,
			State: semanticCacheZoneState{
				CurrentTempC:       CloneFloatPtr(zone.State.CurrentTempC),
				CurrentHumidityPct: CloneFloatPtr(zone.State.CurrentHumidityPct),
				HvacAction:         zone.State.HvacAction,
				SpecialFunction:    zone.State.SpecialFunction,
				HeatingDemandPct:   CloneFloatPtr(zone.State.HeatingDemandPct),
				ValvePositionPct:   CloneFloatPtr(zone.State.ValvePositionPct),
			},
			Config: semanticCacheZoneConfig{
				OperatingMode:              zone.Config.OperatingMode,
				Preset:                     zone.Config.Preset,
				TargetTempC:                CloneFloatPtr(zone.Config.TargetTempC),
				AllowedModes:               append([]string(nil), zone.Config.AllowedModes...),
				CircuitType:                zone.Config.CircuitType,
				AssociatedCircuit:          CloneIntPtr(zone.Config.AssociatedCircuit),
				RoomTemperatureZoneMapping: CloneIntPtr(zone.Config.RoomTemperatureZoneMapping),
			},
		})
	}
	if normalizedSnapshot.DHW != nil {
		cacheV3.DHW = &semanticCacheDHWV3{
			State: semanticCacheDhwState{
				CurrentTempC:     CloneFloatPtr(normalizedSnapshot.DHW.State.CurrentTempC),
				SpecialFunction:  normalizedSnapshot.DHW.State.SpecialFunction,
				HeatingDemandPct: CloneFloatPtr(normalizedSnapshot.DHW.State.HeatingDemandPct),
			},
			Config: semanticCacheDhwConfig{
				OperatingMode: normalizedSnapshot.DHW.Config.OperatingMode,
				Preset:        normalizedSnapshot.DHW.Config.Preset,
				TargetTempC:   CloneFloatPtr(normalizedSnapshot.DHW.Config.TargetTempC),
			},
		}
	}
	if normalizedSnapshot.Boiler != nil {
		cacheV3.Boiler = graphQLBoilerToCache(normalizedSnapshot.Boiler)
	}
	return cacheV3
}

func normalizeSemanticCacheSnapshot(snapshot SemanticCacheSnapshot) SemanticCacheSnapshot {
	out := SemanticCacheSnapshot{
		Zones: make([]Zone, 0, len(snapshot.Zones)),
		DHW:   nil,
	}
	for _, zone := range snapshot.Zones {
		zoneCopy := zone
		zoneCopy.State = ZoneState{
			CurrentTempC:       CloneFloatPtr(zone.State.CurrentTempC),
			CurrentHumidityPct: CloneFloatPtr(zone.State.CurrentHumidityPct),
			HvacAction:         zone.State.HvacAction,
			SpecialFunction:    zone.State.SpecialFunction,
			HeatingDemandPct:   CloneFloatPtr(zone.State.HeatingDemandPct),
			ValvePositionPct:   CloneFloatPtr(zone.State.ValvePositionPct),
		}
		zoneCopy.Config = ZoneConfig{
			OperatingMode:              zone.Config.OperatingMode,
			Preset:                     zone.Config.Preset,
			TargetTempC:                CloneFloatPtr(zone.Config.TargetTempC),
			AllowedModes:               append([]string(nil), zone.Config.AllowedModes...),
			CircuitType:                zone.Config.CircuitType,
			AssociatedCircuit:          CloneIntPtr(zone.Config.AssociatedCircuit),
			RoomTemperatureZoneMapping: CloneIntPtr(zone.Config.RoomTemperatureZoneMapping),
		}
		slices.Sort(zoneCopy.Config.AllowedModes)
		out.Zones = append(out.Zones, zoneCopy)
	}
	slices.SortFunc(out.Zones, func(a, b Zone) int {
		if compare := compareSemanticZoneID(a.ID, b.ID); compare != 0 {
			return compare
		}
		return strings.Compare(a.Name, b.Name)
	})
	if snapshot.DHW != nil {
		out.DHW = &DhwStatus{
			State: DhwState{
				CurrentTempC:     CloneFloatPtr(snapshot.DHW.State.CurrentTempC),
				SpecialFunction:  snapshot.DHW.State.SpecialFunction,
				HeatingDemandPct: CloneFloatPtr(snapshot.DHW.State.HeatingDemandPct),
			},
			Config: DhwConfig{
				OperatingMode: snapshot.DHW.Config.OperatingMode,
				Preset:        snapshot.DHW.Config.Preset,
				TargetTempC:   CloneFloatPtr(snapshot.DHW.Config.TargetTempC),
			},
		}
	}
	out.Boiler = snapshot.Boiler
	return out
}

func normalizeSemanticCacheZonesV3(zones []semanticCacheZoneV3) []semanticCacheZoneV3 {
	out := make([]semanticCacheZoneV3, 0, len(zones))
	for _, zone := range zones {
		zoneCopy := zone
		zoneCopy.State.CurrentTempC = CloneFloatPtr(zone.State.CurrentTempC)
		zoneCopy.State.CurrentHumidityPct = CloneFloatPtr(zone.State.CurrentHumidityPct)
		zoneCopy.State.HeatingDemandPct = CloneFloatPtr(zone.State.HeatingDemandPct)
		zoneCopy.State.ValvePositionPct = CloneFloatPtr(zone.State.ValvePositionPct)
		zoneCopy.Config.TargetTempC = CloneFloatPtr(zone.Config.TargetTempC)
		zoneCopy.Config.AllowedModes = append([]string(nil), zone.Config.AllowedModes...)
		slices.Sort(zoneCopy.Config.AllowedModes)
		zoneCopy.Config.AssociatedCircuit = CloneIntPtr(zone.Config.AssociatedCircuit)
		zoneCopy.Config.RoomTemperatureZoneMapping = CloneIntPtr(zone.Config.RoomTemperatureZoneMapping)
		out = append(out, zoneCopy)
	}
	slices.SortFunc(out, func(a, b semanticCacheZoneV3) int {
		if compare := compareSemanticZoneID(a.ID, b.ID); compare != 0 {
			return compare
		}
		return strings.Compare(a.Name, b.Name)
	})
	return out
}

func compareSemanticZoneID(left, right string) int {
	leftOrdinal, leftOK := parseSemanticZoneOrdinal(left)
	rightOrdinal, rightOK := parseSemanticZoneOrdinal(right)
	if leftOK && rightOK {
		switch {
		case leftOrdinal < rightOrdinal:
			return -1
		case leftOrdinal > rightOrdinal:
			return 1
		}
	}
	return strings.Compare(left, right)
}

func parseSemanticZoneOrdinal(id string) (int, bool) {
	id = strings.TrimSpace(strings.ToLower(id))
	if !strings.HasPrefix(id, "zone-") {
		return 0, false
	}
	ordinalRaw := strings.TrimPrefix(id, "zone-")
	ordinal, err := strconv.Atoi(ordinalRaw)
	if err != nil || ordinal <= 0 {
		return 0, false
	}
	return ordinal, true
}

func cacheBoilerToGraphQL(cache *semanticCacheBoiler) *BoilerStatus {
	if cache == nil {
		return nil
	}

	out := &BoilerStatus{
		State: BoilerState{
			ReturnTemperatureC:    CloneFloatPtr(cache.ReturnTemperatureC),
			DhwTemperatureC:       CloneFloatPtr(cache.DhwTemperatureC),
			DhwTargetTemperatureC: CloneFloatPtr(cache.DhwTargetTemperatureC),
		},
		Diagnostics: BoilerDiagnostics{
			DhwStatusRaw: CloneIntPtr(cache.DhwStatusRaw),
		},
	}

	if cache.DhwOperatingMode != nil {
		v := *cache.DhwOperatingMode
		out.Config.DhwOperatingMode = &v
	}
	return out
}

func graphQLBoilerToCache(status *BoilerStatus) *semanticCacheBoiler {
	if status == nil {
		return nil
	}
	out := &semanticCacheBoiler{
		FlowTemperatureC:      CloneFloatPtr(status.State.FlowTemperatureC),
		ReturnTemperatureC:    CloneFloatPtr(status.State.ReturnTemperatureC),
		DhwTemperatureC:       CloneFloatPtr(status.State.DhwTemperatureC),
		DhwTargetTemperatureC: CloneFloatPtr(status.State.DhwTargetTemperatureC),
		DhwOperatingMode:      CloneStringPtr(status.Config.DhwOperatingMode),
		DhwStatusRaw:          CloneIntPtr(status.Diagnostics.DhwStatusRaw),
	}
	return out
}

func writeSemanticCacheFileAtomic(path string, payload []byte) error {
	if strings.TrimSpace(path) == "" {
		return fmt.Errorf("semantic cache path empty")
	}
	path = filepath.Clean(path)
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}

	tempFile, err := os.CreateTemp(dir, ".semantic-cache-*.tmp")
	if err != nil {
		return err
	}
	tempPath := tempFile.Name()
	removeTemp := true
	defer func() {
		if removeTemp {
			_ = os.Remove(tempPath)
		}
	}()

	if _, err := tempFile.Write(payload); err != nil {
		_ = tempFile.Close()
		return err
	}
	if err := tempFile.Sync(); err != nil {
		_ = tempFile.Close()
		return err
	}
	if err := tempFile.Close(); err != nil {
		return err
	}
	if err := os.Chmod(tempPath, 0o644); err != nil {
		return err
	}
	if err := os.Rename(tempPath, path); err != nil {
		return err
	}
	removeTemp = false
	return nil
}
