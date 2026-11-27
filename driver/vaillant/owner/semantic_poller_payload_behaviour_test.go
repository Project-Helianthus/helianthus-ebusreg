package owner

import (
	"context"
	"encoding/binary"
	"errors"
	"github.com/Project-Helianthus/helianthus-ebusgo/engine"
	"github.com/Project-Helianthus/helianthus-ebusgo/protocol"
	"github.com/Project-Helianthus/helianthus-ebusreg/registry"
	"math"
	"slices"
	"strings"
	"testing"
	"time"
)

func TestVaillantOwnerRefreshDiscoveryRetriesDHWBeforeStartupZoneScan(t *testing.T) {
	if testing.Short() {
		t.Skip("measured long-running coverage; exercised by full CI")
	}
	t.Parallel()

	dhwRequests := 0
	zoneRequestsBeforeDHWSuccess := 0
	dhwSawLiveResponse := false
	poller := &vaillantSemanticPoller{
		publication:    newVaillantSemanticPublication(newRetainedSemanticState(), nil, nil),
		reg:            newTestRegistry(registry.DeviceInfo{Address: 0x15, Manufacturer: "Vaillant", DeviceID: "BASV"}),
		tasks:          newSemanticTaskScheduler(),
		zones:          make(map[byte]*vaillantZoneSnapshot),
		presence:       make(map[byte]*zonePresenceRecord),
		circuits:       make(map[byte]*vaillantCircuitSnapshot),
		radioDevices:   make(map[radioDeviceKey]*vaillantRadioDeviceSnapshot),
		solarCylinders: make(map[byte]*vaillantCylinderSnapshot),
		source:         0x7F,
		requestTimeout: 50 * time.Millisecond,
		b524ProbeFn: func(ctx context.Context, target, opcode, group, instance byte, addr uint16) bool {
			return target == 0x15
		},
	}
	poller.sendFrameFn = func(_ context.Context, frame protocol.Frame) (*protocol.Frame, error) {
		if frame.Primary != vaillantExtRegisterPrimary || frame.Secondary != vaillantExtRegisterSecondary || len(frame.Data) != 6 {
			return &protocol.Frame{Data: []byte{0x00}}, nil
		}
		selector := frame.Data
		group := selector[2]
		instance := selector[3]
		addr := uint16(selector[4]) | uint16(selector[5])<<8
		if group == localZones.group && !dhwSawLiveResponse {
			zoneRequestsBeforeDHWSuccess++
		}
		if group == localDHW.group {
			dhwRequests++
			if dhwRequests <= 2 {
				return nil, errors.New("transient dhw startup miss")
			}
		}
		payload := []byte{0x00}
		switch {
		case group == localDHW.group && instance == dhwInstance && addr == dhwCurrentTemp:
			payload = make([]byte, 4)
			binary.LittleEndian.PutUint32(payload, math.Float32bits(48.5))
			dhwSawLiveResponse = true
		case group == localDHW.group && instance == dhwInstance && addr == dhwOperationMode:
			payload = []byte{0x01, 0x00}
			dhwSawLiveResponse = true
		case group == localZones.group && addr == zoneIndex:
			payload = []byte{0xFF}
		}
		data := []byte{0x01, instance, group, byte(addr), byte(addr >> 8)}
		data = append(data, payload...)
		return &protocol.Frame{Data: data}, nil
	}

	poller.refreshDiscovery(context.Background())

	if dhwRequests < 4 {
		t.Fatalf("DHW startup requests = %d; want retry after first probe pair fails", dhwRequests)
	}
	if zoneRequestsBeforeDHWSuccess != 0 {
		t.Fatalf("zone requests before DHW live response = %d; want 0", zoneRequestsBeforeDHWSuccess)
	}
	if dhw := publishedProvider(poller).DHW(); dhw == nil {
		t.Fatal("provider.DHW() = nil; want retry to prime DHW before startup zone scan")
	}
}

func TestVaillantOwnerParseB509ReadPayload(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		payload []byte
		addr    uint16
		want    []byte
		ok      bool
	}{
		{
			name:    "accepts single byte zero payload",
			payload: []byte{0x00},
			addr:    0x4400,
			want:    []byte{0x00},
			ok:      true,
		},
		{
			name:    "accepts single byte nonzero payload",
			payload: []byte{0x01},
			addr:    0x7B00,
			want:    []byte{0x01},
			ok:      true,
		},
		{
			name:    "strips opcode header",
			payload: []byte{0x0D, 0x44, 0x00, 0x00},
			addr:    0x4400,
			want:    []byte{0x00},
			ok:      true,
		},
		{
			name:    "strips compact address header",
			payload: []byte{0x44, 0x00, 0x00},
			addr:    0x4400,
			want:    []byte{0x00},
			ok:      true,
		},
		{
			name:    "rejects empty payload",
			payload: nil,
			addr:    0x4400,
			ok:      false,
		},
		{
			name:    "rejects header only opcode payload",
			payload: []byte{0x0D, 0x44, 0x00},
			addr:    0x4400,
			ok:      false,
		},
		{
			name:    "rejects header only compact payload",
			payload: []byte{0x44, 0x00},
			addr:    0x4400,
			ok:      false,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			got, ok := parseB509ReadPayload(test.payload, test.addr)
			if ok != test.ok {
				t.Fatalf("ok = %v; want %v", ok, test.ok)
			}
			if !test.ok {
				return
			}
			if !slices.Equal(got, test.want) {
				t.Fatalf("parseB509ReadPayload(%x, 0x%04x) = %x; want %x", test.payload, test.addr, got, test.want)
			}
		})
	}
}

func TestVaillantOwnerEncodeBoilerConfigPayload(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		fieldName   string
		spec        boilerConfigFieldSpec
		value       float64
		wantPayload []byte
		wantValue   float64
		wantErr     string
	}{
		{
			name:        "temperature field uses DATA2c",
			fieldName:   "maximum_dhw_flow_temperature_setpoint",
			spec:        boilerConfigFieldSpecs["maximum_dhw_flow_temperature_setpoint"],
			value:       59,
			wantPayload: []byte{0xB0, 0x03},
			wantValue:   59,
		},
		{
			name:        "temperature field normalizes to wire precision",
			fieldName:   "maximum_heating_flow_temperature_setpoint",
			spec:        boilerConfigFieldSpecs["maximum_heating_flow_temperature_setpoint"],
			value:       55.1,
			wantPayload: []byte{0x71, 0x03},
			wantValue:   55.0625,
		},
		{
			name:        "partload field uses UCH",
			fieldName:   "heating_partial_load_power",
			spec:        boilerConfigFieldSpecs["heating_partial_load_power"],
			value:       18,
			wantPayload: []byte{0x12},
			wantValue:   18,
		},
		{
			name:      "fractional partload rejected",
			fieldName: "dhw_partial_load_power",
			spec:      boilerConfigFieldSpecs["dhw_partial_load_power"],
			value:     18.5,
			wantErr:   "whole-number kW required",
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			payload, normalizedValue, err := encodeBoilerConfigPayload(test.fieldName, test.spec, test.value)
			if test.wantErr != "" {
				if err == nil {
					t.Fatalf("encodeBoilerConfigPayload(%q, %.4g) error = nil; want %q", test.fieldName, test.value, test.wantErr)
				}
				if !strings.Contains(err.Error(), test.wantErr) {
					t.Fatalf("encodeBoilerConfigPayload(%q, %.4g) error = %q; want substring %q", test.fieldName, test.value, err.Error(), test.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("encodeBoilerConfigPayload(%q, %.4g) error = %v", test.fieldName, test.value, err)
			}
			if !slices.Equal(payload, test.wantPayload) {
				t.Fatalf("encodeBoilerConfigPayload(%q, %.4g) payload = %x; want %x", test.fieldName, test.value, payload, test.wantPayload)
			}
			if normalizedValue != test.wantValue {
				t.Fatalf("encodeBoilerConfigPayload(%q, %.4g) value = %.4g; want %.4g", test.fieldName, test.value, normalizedValue, test.wantValue)
			}
		})
	}
}

func TestVaillantOwnerParseBoilerConfigValue(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		rawValue  string
		spec      boilerConfigFieldSpec
		wantValue float64
		wantErr   string
	}{
		{
			name:      "accepts finite in-range value",
			rawValue:  "59",
			spec:      boilerConfigFieldSpecs["maximum_dhw_flow_temperature_setpoint"],
			wantValue: 59,
		},
		{
			name:     "rejects nan",
			rawValue: "NaN",
			spec:     boilerConfigFieldSpecs["maximum_heating_flow_temperature_setpoint"],
			wantErr:  "finite number required",
		},
		{
			name:     "rejects infinity",
			rawValue: "+Inf",
			spec:     boilerConfigFieldSpecs["maximum_heating_flow_temperature_setpoint"],
			wantErr:  "finite number required",
		},
		{
			name:     "rejects out of range value",
			rawValue: "85",
			spec:     boilerConfigFieldSpecs["maximum_heating_flow_temperature_setpoint"],
			wantErr:  "out of range",
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			value, err := parseBoilerConfigValue(test.rawValue, test.spec)
			if test.wantErr != "" {
				if err == nil {
					t.Fatalf("parseBoilerConfigValue(%q) error = nil; want %q", test.rawValue, test.wantErr)
				}
				if !strings.Contains(err.Error(), test.wantErr) {
					t.Fatalf("parseBoilerConfigValue(%q) error = %q; want substring %q", test.rawValue, err.Error(), test.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("parseBoilerConfigValue(%q) error = %v", test.rawValue, err)
			}
			if value != test.wantValue {
				t.Fatalf("parseBoilerConfigValue(%q) value = %.4g; want %.4g", test.rawValue, value, test.wantValue)
			}
		})
	}
}

func TestVaillantOwnerSetBoilerConfigWriteConfirmForcesLiveReadAfterWrite(t *testing.T) {
	t.Parallel()

	const (
		boilerAddress = byte(0x08)
		fieldName     = "maximum_dhw_flow_temperature_setpoint"
	)

	spec := boilerConfigFieldSpecs[fieldName]
	registerAddr := spec.addrs[0]

	stalePayload, _, err := encodeBoilerConfigPayload(fieldName, spec, 58)
	if err != nil {
		t.Fatalf("encode stale payload error = %v", err)
	}
	expectedPayload, _, err := encodeBoilerConfigPayload(fieldName, spec, 59)
	if err != nil {
		t.Fatalf("encode expected payload error = %v", err)
	}

	key := engine.NewB509WatchKey(boilerAddress, registerAddr)
	catalog, err := engine.NewWatchCatalog([]engine.WatchDescriptor{
		{
			Key:               key,
			SemanticClass:     engine.WatchSemanticClassState,
			FreshnessProfile:  engine.WatchFreshnessProfileStateFast,
			DecoderID:         "test.semantic.b509.write.confirm",
			CorrelationPolicy: engine.WatchCorrelationPolicyRequestResponse,
			DirectApplyPolicy: engine.WatchDirectApplyPolicyStateDefault,
		},
	}, EngineCollaborator{})
	if err != nil {
		t.Fatalf("NewWatchCatalog() error = %v", err)
	}
	activations := engine.NewWatchActivationSet(catalog)
	if err := activations.Activate(engine.WatchActivationSourcePoller, key); err != nil {
		t.Fatalf("Activate(poller) error = %v", err)
	}

	now := time.Unix(100, 0)
	shadow := engine.NewShadowCache(engine.ShadowCacheOptions{
		Catalog:               catalog,
		Activations:           activations,
		Capacity:              8,
		PinnedCapacity:        4,
		WriteConfirmPinnedCap: 2,
		Now:                   func() time.Time { return now },
	})
	scheduler := engine.NewSemanticReadScheduler()
	scheduler.SetShadowCache(shadow)

	bus := newB509MutationTestBus(registerAddr, stalePayload)
	poller := &vaillantSemanticPoller{
		scheduler:       scheduler,
		shadow:          shadow,
		sendFrameFn:     bus.Send,
		source:          0x31,
		requestTimeout:  250 * time.Millisecond,
		boilerAddress:   boilerAddress,
		watchObserver:   nil,
		transportConfig: engine.TransportConfig{},
	}

	_ = poller.prepareSemanticReadWatch(key)
	staleWrite := shadow.Write(engine.ShadowWrite{
		Key:        key,
		Source:     engine.ShadowWriteSourcePassive,
		Confidence: engine.ShadowConfidenceHigh,
		Value:      stalePayload,
		ObservedAt: now,
	})
	if !staleWrite.Accepted {
		t.Fatalf("shadow stale seed write rejected: %s", staleWrite.Reason)
	}

	result := poller.SetBoilerConfig(context.Background(), fieldName, "59")
	if !result.Success {
		t.Fatalf("SetBoilerConfig() = %+v; want success with live post-write readback", result)
	}

	reads, writes := bus.Counters()
	if writes != 1 {
		t.Fatalf("write calls = %d; want 1", writes)
	}
	if reads < 1 {
		t.Fatalf("read calls = %d; want at least one live read for write confirm", reads)
	}
	trace := bus.OperationTrace()
	writeAt := strings.IndexByte(trace, 'W')
	if writeAt < 0 {
		t.Fatalf("operation trace = %q; want a write operation", trace)
	}
	if writeAt+1 >= len(trace) || !strings.Contains(trace[writeAt+1:], "R") {
		t.Fatalf("operation trace = %q; want read-after-write for live confirm (no stale cache confirm)", trace)
	}
	if got := bus.Value(); !slices.Equal(got, expectedPayload) {
		t.Fatalf("bus value = %x; want written payload %x", got, expectedPayload)
	}
}

func TestVaillantOwnerSetBoilerConfigWriteConfirmInvalidatesShadowBeforeReadback(t *testing.T) {
	t.Parallel()

	const (
		boilerAddress = byte(0x08)
		fieldName     = "maximum_dhw_flow_temperature_setpoint"
	)

	spec := boilerConfigFieldSpecs[fieldName]
	registerAddr := spec.addrs[0]

	stalePayload, _, err := encodeBoilerConfigPayload(fieldName, spec, 58)
	if err != nil {
		t.Fatalf("encode stale payload error = %v", err)
	}
	expectedPayload, _, err := encodeBoilerConfigPayload(fieldName, spec, 59)
	if err != nil {
		t.Fatalf("encode expected payload error = %v", err)
	}

	key := engine.NewB509WatchKey(boilerAddress, registerAddr)
	catalog, err := engine.NewWatchCatalog([]engine.WatchDescriptor{
		{
			Key:               key,
			SemanticClass:     engine.WatchSemanticClassState,
			FreshnessProfile:  engine.WatchFreshnessProfileStateFast,
			DecoderID:         "test.semantic.b509.write.confirm.fence",
			CorrelationPolicy: engine.WatchCorrelationPolicyRequestResponse,
			DirectApplyPolicy: engine.WatchDirectApplyPolicyStateDefault,
		},
	}, EngineCollaborator{})
	if err != nil {
		t.Fatalf("NewWatchCatalog() error = %v", err)
	}
	activations := engine.NewWatchActivationSet(catalog)
	if err := activations.Activate(engine.WatchActivationSourcePoller, key); err != nil {
		t.Fatalf("Activate(poller) error = %v", err)
	}

	now := time.Unix(100, 0)
	shadow := engine.NewShadowCache(engine.ShadowCacheOptions{
		Catalog:               catalog,
		Activations:           activations,
		Capacity:              8,
		PinnedCapacity:        4,
		WriteConfirmPinnedCap: 2,
		Now:                   func() time.Time { return now },
	})
	scheduler := engine.NewSemanticReadScheduler()
	scheduler.SetShadowCache(shadow)

	var seenGenerations []uint64
	var seenEligibility []bool
	bus := newB509MutationTestBus(registerAddr, stalePayload)
	bus.onRead = func() {
		snapshot := shadow.SnapshotEligibility(key)
		seenGenerations = append(seenGenerations, snapshot.Generation)
		seenEligibility = append(seenEligibility, snapshot.Eligible)
	}
	poller := &vaillantSemanticPoller{
		scheduler:       scheduler,
		shadow:          shadow,
		sendFrameFn:     bus.Send,
		source:          0x31,
		requestTimeout:  250 * time.Millisecond,
		boilerAddress:   boilerAddress,
		watchObserver:   nil,
		transportConfig: engine.TransportConfig{},
	}

	_ = poller.prepareSemanticReadWatch(key)
	staleWrite := shadow.Write(engine.ShadowWrite{
		Key:        key,
		Source:     engine.ShadowWriteSourcePassive,
		Confidence: engine.ShadowConfidenceHigh,
		Value:      stalePayload,
		ObservedAt: now,
	})
	if !staleWrite.Accepted {
		t.Fatalf("shadow stale seed write rejected: %s", staleWrite.Reason)
	}
	if staleWrite.Generation == 0 {
		t.Fatal("stale seed generation = 0; want non-zero generation")
	}

	result := poller.SetBoilerConfig(context.Background(), fieldName, "59")
	if !result.Success {
		t.Fatalf("SetBoilerConfig() = %+v; want success", result)
	}

	if len(seenGenerations) == 0 {
		t.Fatal("expected readback generation snapshots to be captured")
	}
	sawAdvancedGeneration := false
	for _, generation := range seenGenerations {
		if generation > staleWrite.Generation {
			sawAdvancedGeneration = true
			break
		}
	}
	if !sawAdvancedGeneration {
		t.Fatalf("seen generations = %v; want at least one generation > stale generation %d after write-confirm fence invalidation", seenGenerations, staleWrite.Generation)
	}
	if len(seenEligibility) == 0 {
		t.Fatal("expected readback eligibility snapshots to be captured")
	}
	sawIneligible := false
	for _, eligible := range seenEligibility {
		if !eligible {
			sawIneligible = true
			break
		}
	}
	if !sawIneligible {
		t.Fatalf("seen eligibility = %v; want at least one ineligible snapshot before confirm read", seenEligibility)
	}

	if got := bus.Value(); !slices.Equal(got, expectedPayload) {
		t.Fatalf("bus value = %x; want written payload %x", got, expectedPayload)
	}
}

func TestVaillantOwnerBoilerSnapshotWithConfigValueClonesExistingSnapshot(t *testing.T) {
	t.Parallel()

	existingHc := 75.0
	existingPartload := 18.0
	existing := &vaillantBoilerSnapshot{
		FlowsetHcMaxC: cloneFloat64Ptr(&existingHc),
		PartloadHcKW:  cloneFloat64Ptr(&existingPartload),
	}

	updated := boilerSnapshotWithConfigValue(existing, "maximum_heating_flow_temperature_setpoint", 55.0625)
	if updated == existing {
		t.Fatal("boilerSnapshotWithConfigValue returned the original snapshot pointer")
	}
	if existing.FlowsetHcMaxC == nil || *existing.FlowsetHcMaxC != 75 {
		t.Fatalf("existing FlowsetHcMaxC = %v; want 75", existing.FlowsetHcMaxC)
	}
	if updated.FlowsetHcMaxC == nil || math.Abs(*updated.FlowsetHcMaxC-55.0625) > 0 {
		t.Fatalf("updated FlowsetHcMaxC = %v; want 55.0625", updated.FlowsetHcMaxC)
	}
	if updated.PartloadHcKW == nil || *updated.PartloadHcKW != 18 {
		t.Fatalf("updated PartloadHcKW = %v; want 18", updated.PartloadHcKW)
	}
}

func TestVaillantOwnerDecodeBoilerConfigRaw(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		spec      boilerConfigFieldSpec
		raw       []byte
		wantValue float64
		wantOK    bool
	}{
		{
			name:      "temperature register decodes data2c",
			spec:      boilerConfigFieldSpecs["maximum_dhw_flow_temperature_setpoint"],
			raw:       []byte{0xB0, 0x03},
			wantValue: 59,
			wantOK:    true,
		},
		{
			name:      "power register decodes uch",
			spec:      boilerConfigFieldSpecs["dhw_partial_load_power"],
			raw:       []byte{0x16},
			wantValue: 22,
			wantOK:    true,
		},
		{
			name:   "temperature register rejects one byte payload",
			spec:   boilerConfigFieldSpecs["maximum_heating_flow_temperature_setpoint"],
			raw:    []byte{0xF0},
			wantOK: false,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			value, ok := decodeBoilerConfigRaw(test.spec, test.raw)
			if ok != test.wantOK {
				t.Fatalf("decodeBoilerConfigRaw(%q, %x) ok = %v; want %v", test.name, test.raw, ok, test.wantOK)
			}
			if ok && value != test.wantValue {
				t.Fatalf("decodeBoilerConfigRaw(%q, %x) value = %.4g; want %.4g", test.name, test.raw, value, test.wantValue)
			}
		})
	}
}

func TestVaillantOwnerFlowsetHcMaxSpecIncludesFallbackAddress(t *testing.T) {
	t.Parallel()

	spec := boilerConfigFieldSpecs["maximum_heating_flow_temperature_setpoint"]
	want := []uint16{boilerB509FlowsetHcMaxC, boilerB509FlowsetHcMaxCFallback}
	if !slices.Equal(spec.addrs, want) {
		t.Fatalf("maximum_heating_flow_temperature_setpoint addrs = %#v; want %#v", spec.addrs, want)
	}
}

func TestVaillantOwnerComposeZoneName(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		primary string
		prefix  string
		suffix  string
		want    string
	}{
		{
			name:    "prefers primary name",
			primary: "Parter",
			prefix:  "ignored",
			suffix:  "ignored",
			want:    "Parter",
		},
		{
			name:   "joins prefix and suffix",
			prefix: "Etaj",
			suffix: "2",
			want:   "Etaj 2",
		},
		{
			name:   "handles single part",
			prefix: "Parter",
			want:   "Parter",
		},
		{
			name: "returns empty when all missing",
			want: "",
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			if got := composeZoneName(test.primary, test.prefix, test.suffix); got != test.want {
				t.Fatalf("composeZoneName(%q, %q, %q) = %q; want %q", test.primary, test.prefix, test.suffix, got, test.want)
			}
		})
	}
}

func TestVaillantOwnerDecodeB524Uint16(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		payload []byte
		want    uint16
		ok      bool
	}{
		{name: "empty", payload: nil, want: 0, ok: false},
		{name: "single byte", payload: []byte{0x02}, want: 0x0002, ok: true},
		{name: "little endian", payload: []byte{0x02, 0x01}, want: 0x0102, ok: true},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			got, ok := decodeB524Uint16(test.payload)
			if ok != test.ok {
				t.Fatalf("decodeB524Uint16(%v) ok = %v; want %v", test.payload, ok, test.ok)
			}
			if got != test.want {
				t.Fatalf("decodeB524Uint16(%v) = 0x%04x; want 0x%04x", test.payload, got, test.want)
			}
		})
	}
}

func TestVaillantOwnerDecodeB524FirmwareVersion(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		payload []byte
		wantNil bool
		want    string
	}{
		{name: "short payload", payload: []byte{0x08, 0x05}, wantNil: true},
		{name: "ff ff ff is absent", payload: []byte{0xFF, 0xFF, 0xFF}, wantNil: true},
		{name: "formats byte decimal", payload: []byte{0x08, 0x05, 0x00}, want: "08.05.00"},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			got := decodeB524FirmwareVersion(test.payload)
			if test.wantNil {
				if got != nil {
					t.Fatalf("decodeB524FirmwareVersion(%v) = %q; want nil", test.payload, *got)
				}
				return
			}
			if got == nil || *got != test.want {
				t.Fatalf("decodeB524FirmwareVersion(%v) = %v; want %q", test.payload, got, test.want)
			}
		})
	}
}
