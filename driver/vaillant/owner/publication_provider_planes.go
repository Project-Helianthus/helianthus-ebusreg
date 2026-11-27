package owner

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/Project-Helianthus/helianthus-ebusreg/driver/vaillant/output"
)

func (provider *retainedSemanticState) System() *output.SystemStatus {
	if provider == nil {
		return nil
	}
	provider.mu.RLock()
	defer provider.mu.RUnlock()
	return cloneSystemStatus(provider.system)
}

func (provider *retainedSemanticState) SetSystem(status *output.SystemStatus) {
	if provider == nil {
		return
	}
	provider.mu.Lock()
	defer provider.mu.Unlock()
	provider.system = cloneSystemStatus(status)
}

func (provider *retainedSemanticState) SetSystemFromCache(status *output.SystemStatus) {
	provider.SetSystem(status)
}

func (provider *retainedSemanticState) Schedules() *output.ScheduleStatus {
	if provider == nil {
		return nil
	}
	provider.mu.RLock()
	defer provider.mu.RUnlock()
	return cloneScheduleStatus(provider.schedules)
}

func (provider *retainedSemanticState) SetSchedules(status *output.ScheduleStatus) {
	if provider == nil {
		return
	}
	provider.mu.Lock()
	defer provider.mu.Unlock()
	provider.schedules = cloneScheduleStatus(status)
}

func (provider *retainedSemanticState) AdapterHardwareInfo() *output.AdapterHardwareInfo {
	if provider == nil {
		return nil
	}
	provider.mu.RLock()
	defer provider.mu.RUnlock()
	return cloneAdapterHardwareInfo(provider.adapter)
}

func (provider *retainedSemanticState) SetAdapterHardwareInfo(info *output.AdapterHardwareInfo) {
	if provider == nil {
		return
	}
	provider.mu.Lock()
	defer provider.mu.Unlock()
	provider.adapter = cloneAdapterHardwareInfo(info)
}

// semanticRegistryLeaf identifies one materialized leaf owned by the
// semantic owner. It carries no value payload.
type semanticRegistryLeaf struct {
	Path           string
	PromotionState string
	Source         string
}

// semanticRegistrySnapshot is a read-only inventory from the semantic owner.
type semanticRegistrySnapshot struct {
	Authority string
	Leaves    []semanticRegistryLeaf
}

// SemanticRegistryState snapshots only owner-materialized semantic state.
func (provider *retainedSemanticState) SemanticRegistryState() (semanticRegistrySnapshot, error) {
	state := semanticRegistrySnapshot{
		Authority: "ebus.promoted",
		Leaves:    []semanticRegistryLeaf{},
	}
	if provider == nil {
		return state, nil
	}

	provider.mu.RLock()
	values := []struct {
		path  string
		value any
	}{
		{path: "/zones", value: provider.zonesForInventoryLocked()},
		{path: "/dhw", value: cloneDhwStatus(provider.dhw)},
		{path: "/circuits", value: provider.circuitsForInventoryLocked()},
		{path: "/radio", value: provider.radioForInventoryLocked()},
		{path: "/solar", value: cloneSolarStatus(provider.solar)},
		{path: "/cylinders", value: cloneCylinderStatuses(provider.cylinders)},
		{path: "/boiler", value: provider.boilerForInventoryLocked()},
		{path: "/system", value: cloneSystemStatus(provider.system)},
		{path: "/schedules", value: cloneScheduleStatus(provider.schedules)},
		{path: "/adapter", value: cloneAdapterHardwareInfo(provider.adapter)},
	}
	if provider.functionModuleMode != output.FunctionModuleSemanticModeAbsent && provider.functionModuleMode != "" {
		values = append(values, struct {
			path  string
			value any
		}{path: "/fm5", value: provider.functionModuleMode})
	}
	provider.mu.RUnlock()

	energyRevision := uint64(0)
	var energy *output.EnergyTotals
	if provider.energyMerge != nil {
		energyRevision = provider.energyMerge.Revision()
		if energyRevision > 0 {
			energy = provider.EnergyTotals()
		}
		if provider.energyMerge.Revision() != energyRevision {
			return semanticRegistrySnapshot{}, fmt.Errorf("inventory semantic owner energy changed during snapshot")
		}
	}
	if energyRevision > 0 {
		values = append(values, struct {
			path  string
			value any
		}{path: "/energy", value: energy})
	}

	for _, item := range values {
		paths, err := semanticMaterializedLeafPaths(item.path, item.value)
		if err != nil {
			return semanticRegistrySnapshot{}, fmt.Errorf("inventory semantic owner %s: %w", item.path, err)
		}
		for _, path := range paths {
			state.Leaves = append(state.Leaves, semanticRegistryLeaf{
				Path: path, PromotionState: "PROMOTED", Source: "ebus",
			})
		}
	}
	return state, nil
}

func (provider *retainedSemanticState) zonesForInventoryLocked() []output.Zone {
	if !provider.zonePublished {
		return nil
	}
	zones := make([]output.Zone, len(provider.zones))
	for i, zone := range provider.zones {
		zones[i] = cloneZone(zone)
	}
	return zones
}

func (provider *retainedSemanticState) circuitsForInventoryLocked() []output.CircuitStatus {
	if !provider.circuitPublished {
		return nil
	}
	circuits := make([]output.CircuitStatus, len(provider.circuits))
	for i, circuit := range provider.circuits {
		circuits[i] = cloneCircuitStatus(circuit)
	}
	return circuits
}

func (provider *retainedSemanticState) radioForInventoryLocked() []output.RadioDevice {
	if !provider.radioPublished {
		return nil
	}
	radio := make([]output.RadioDevice, len(provider.radio))
	for i, device := range provider.radio {
		radio[i] = cloneRadioDevice(device)
	}
	return radio
}

func (provider *retainedSemanticState) boilerForInventoryLocked() *output.BoilerStatus {
	if !provider.boilerPublished || provider.boiler == nil {
		return nil
	}
	boiler := *provider.boiler
	boiler.State = cloneBoilerState(boiler.State)
	boiler.Config = cloneBoilerConfig(boiler.Config)
	boiler.Diagnostics = cloneBoilerDiagnostics(boiler.Diagnostics)
	return &boiler
}

func semanticMaterializedLeafPaths(root string, value any) ([]string, error) {
	if value == nil {
		return nil, nil
	}
	raw, err := json.Marshal(value)
	if err != nil {
		return nil, err
	}
	decoder := json.NewDecoder(bytes.NewReader(raw))
	decoder.UseNumber()
	var decoded any
	if err := decoder.Decode(&decoded); err != nil {
		return nil, err
	}
	paths := []string{}
	var walk func(string, any)
	walk = func(path string, current any) {
		switch typed := current.(type) {
		case nil:
			return
		case map[string]any:
			keys := make([]string, 0, len(typed))
			for key := range typed {
				keys = append(keys, key)
			}
			sort.Strings(keys)
			for _, key := range keys {
				escaped := strings.ReplaceAll(strings.ReplaceAll(key, "~", "~0"), "/", "~1")
				walk(path+"/"+escaped, typed[key])
			}
		case []any:
			for index, item := range typed {
				walk(path+"/"+strconv.Itoa(index), item)
			}
		default:
			paths = append(paths, path)
		}
	}
	walk(root, decoded)
	return paths, nil
}
