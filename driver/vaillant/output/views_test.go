package output

import (
	"testing"
)

func TestVaillantZoneViewsDeriveRadioStateOnlyForMissingZoneValues(t *testing.T) {
	zoneAssignment := uint8(2)
	radioTemp := 20.5
	radioHumidity := 44.0
	ownTemp := 22.0
	ownHumidity := 51.0
	ownMapping := uint16(0)
	wantMapping := uint16(zoneAssignment) + 1
	views := []NativeZoneView{
		{Instance: 0x00, Snapshot: &NativeZoneSnapshot{Instance: 0x00, Present: true, CurrentTempC: &ownTemp, HumidityPct: &ownHumidity, ConfigurationRoomTemperatureZoneMappingRaw: &ownMapping}},
		{Instance: 0x01, Snapshot: &NativeZoneSnapshot{Instance: 0x01, Present: true, Name: "Upstairs", AllowedModes: []string{"auto"}}, RadioRoomTempC: &radioTemp, RadioRoomHumidityPct: &radioHumidity, InferredRoomMappingRaw: &wantMapping},
	}

	zone := ZoneFromView(views[1])
	if zone.ID != "zone-2" || zone.Name != "Upstairs" {
		t.Fatalf("zone identity = %q/%q; want zone-2/Upstairs", zone.ID, zone.Name)
	}
	if zone.State.CurrentTempC == nil || *zone.State.CurrentTempC != 20.5 || zone.State.CurrentHumidityPct == nil || *zone.State.CurrentHumidityPct != radioHumidity {
		t.Fatalf("zone state = %#v; want radio room fallback", zone.State)
	}
	wantRoomMapping := decodeRoomTemperatureZoneMapping(&wantMapping)
	if zone.Config.RoomTemperatureZoneMapping == nil || wantRoomMapping == nil || *zone.Config.RoomTemperatureZoneMapping != *wantRoomMapping {
		t.Fatalf("room mapping = %v; want %v", zone.Config.RoomTemperatureZoneMapping, wantRoomMapping)
	}
	associated := uint16(ResolveAssociatedCircuitInstance(&wantMapping, 0x01))
	wantAssociated := decodeAssociatedCircuit(&associated)
	if zone.Config.AssociatedCircuit == nil || wantAssociated == nil || *zone.Config.AssociatedCircuit != *wantAssociated {
		t.Fatalf("associated circuit = %v; want %v", zone.Config.AssociatedCircuit, wantAssociated)
	}

	own := ZoneFromView(views[0])
	if own.Name != "Zone 1" || own.State.CurrentTempC == nil || *own.State.CurrentTempC != 22.0 {
		t.Fatalf("zone with its own values = %#v", own)
	}
}
