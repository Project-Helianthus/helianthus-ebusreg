package productids

import (
	"strings"
	"testing"
)

func TestParseCatalogFiltersIncompleteRows(t *testing.T) {
	input := strings.Join([]string{
		"brand,family,product_model,part_number,role",
		"BrandA,FamilyA,ModelA,PN1,Boiler",
		"BrandB,FamilyB,,PN2,Boiler",
		"BrandC,FamilyC,ModelC,,Heat Pump",
		"BrandA,FamilyA,ModelA,PN1,Duplicate",
		"",
	}, "\n")

	catalog, err := parseCatalog(strings.NewReader(input))
	if err != nil {
		t.Fatalf("parseCatalog error: %v", err)
	}
	if len(catalog.All) != 4 {
		t.Fatalf("expected 4 rows, got %d", len(catalog.All))
	}
	if _, ok := catalog.ByPartNumber["PN1"]; !ok {
		t.Fatalf("expected PN1 in lookup")
	}
	if catalog.ByPartNumber["PN1"].Role != "Boiler" {
		t.Fatalf("expected first PN1 entry to win")
	}
	if _, ok := catalog.ByPartNumber["PN2"]; ok {
		t.Fatalf("expected PN2 to be filtered out")
	}
}

func TestParseCatalogMissingColumns(t *testing.T) {
	input := "brand,family,product_model,part_number\nBrandA,FamilyA,ModelA,PN1"
	if _, err := parseCatalog(strings.NewReader(input)); err != ErrMissingColumns {
		t.Fatalf("expected ErrMissingColumns, got %v", err)
	}
}

func TestControllerCapabilityPresent(t *testing.T) {
	catalog := Catalog{
		ByPartNumber: map[string]Record{
			"PN_REG": {PartNumber: "PN_REG", Role: "Regulator"},
		},
	}

	got := catalog.ControllerCapability("PN_REG")
	if got != ControllerPresent {
		t.Fatalf("expected ControllerPresent, got %v", got)
	}
}

func TestControllerCapabilityNone(t *testing.T) {
	catalog := Catalog{
		ByPartNumber: map[string]Record{
			"PN_BOIL": {PartNumber: "PN_BOIL", Role: "Boiler"},
		},
	}

	got := catalog.ControllerCapability("PN_BOIL")
	if got != ControllerNone {
		t.Fatalf("expected ControllerNone, got %v", got)
	}
}

func TestControllerCapabilityUnknown(t *testing.T) {
	catalog := Catalog{
		ByPartNumber: map[string]Record{
			"PN_REG": {PartNumber: "PN_REG", Role: "Regulator"},
		},
	}

	got := catalog.ControllerCapability("PN_UNKNOWN")
	if got != ControllerUnknown {
		t.Fatalf("expected ControllerUnknown, got %v", got)
	}
}

func TestControllerCapabilityString(t *testing.T) {
	tests := []struct {
		name string
		in   ControllerCapability
		want string
	}{
		{name: "unknown", in: ControllerUnknown, want: "ControllerUnknown"},
		{name: "none", in: ControllerNone, want: "ControllerNone"},
		{name: "present", in: ControllerPresent, want: "ControllerPresent"},
	}

	for _, tt := range tests {
		if got := tt.in.String(); got != tt.want {
			t.Fatalf("%s: expected %q, got %q", tt.name, tt.want, got)
		}
	}
}

func TestControllerCapabilityCaseInsensitive(t *testing.T) {
	catalog := Catalog{
		ByPartNumber: map[string]Record{
			"PN_REG": {PartNumber: "PN_REG", Role: "regulator"},
		},
	}

	got := catalog.ControllerCapability("PN_REG")
	if got != ControllerPresent {
		t.Fatalf("expected ControllerPresent for lowercase role, got %v", got)
	}
}

func TestControllerCapabilityTrimSpace(t *testing.T) {
	catalog := Catalog{
		ByPartNumber: map[string]Record{
			"PN_REG": {PartNumber: "PN_REG", Role: "Regulator"},
		},
	}

	got := catalog.ControllerCapability("  PN_REG  ")
	if got != ControllerPresent {
		t.Fatalf("expected ControllerPresent for padded part number, got %v", got)
	}
}

func TestControllerCapabilityRealCatalog(t *testing.T) {
	catalog, err := LoadCatalog()
	if err != nil {
		t.Fatalf("LoadCatalog error: %v", err)
	}

	if got := catalog.ControllerCapability("0020171314"); got != ControllerPresent {
		t.Fatalf("expected 0020171314 to be ControllerPresent, got %v", got)
	}
	if got := catalog.ControllerCapability("0010016292"); got != ControllerNone {
		t.Fatalf("expected 0010016292 to be ControllerNone, got %v", got)
	}
}
