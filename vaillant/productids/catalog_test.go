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
