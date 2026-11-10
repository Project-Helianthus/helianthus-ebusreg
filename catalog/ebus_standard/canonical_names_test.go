package ebus_standard_catalog

import (
	"encoding/json"
	"testing"
)

func TestCanonicalOutdoorNames(t *testing.T) {
	for _, secondary := range []uint8{0, 2} {
		key := ab28TestKey(secondary, "request", "outdoor_temperature")
		row := ab28TestRow(t, key)
		ab28Assert(t, row, "unit", `"degrees Celsius"`)
		ab28Assert(t, row, "width", `2`)
		ab28Assert(t, row, "encoding", `"signed LE16 /256"`)
		for _, old := range []string{"outside_temperature", "temperature"} {
			key.Ab28Field = old
			if raw, found, err := Ab28Lookup(key); err != nil || found || raw != nil {
				t.Fatalf("documentary key resolved: %s %v", old, err)
			}
		}
	}
	raw, err := Ab28Entries()
	if err != nil {
		t.Fatal(err)
	}
	var rows []struct {
		Key Ab28Key `json:"key"`
	}
	if err := json.Unmarshal(raw, &rows); err != nil {
		t.Fatal(err)
	}
	count := 0
	for _, row := range rows {
		if row.Key.Ab28Secondary != nil && (*row.Key.Ab28Secondary == 0 || *row.Key.Ab28Secondary == 2) && row.Key.Ab28Field == "outdoor_temperature" {
			count++
		}
		if row.Key.Ab28Field == "outside_temperature" {
			t.Fatal("old outdoor name emitted")
		}
	}
	if count != 2 {
		t.Fatalf("canonical quantity count: %d", count)
	}
	duration := ab28TestRow(t, ab28TestKey(2, "request", "duration"))
	ab28Assert(t, duration, "default", `0`)
	ab28Assert(t, duration, "zero_permanent_interpretation", `"Unknown"`)
}
