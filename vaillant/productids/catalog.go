package productids

import (
	_ "embed"
	"encoding/csv"
	"errors"
	"io"
	"strings"
)

//go:embed product_ids.csv
var rawCSV string

var (
	ErrMissingColumns = errors.New("productids: missing required columns")
	requiredColumns   = []string{"brand", "family", "product_model", "part_number", "role"}
)

type Record struct {
	Brand        string
	Family       string
	ProductModel string
	PartNumber   string
	Role         string
	Notes        string
}

type Catalog struct {
	All          []Record
	ByPartNumber map[string]Record
}

func LoadCatalog() (Catalog, error) {
	return parseCatalog(strings.NewReader(rawCSV))
}

func parseCatalog(r io.Reader) (Catalog, error) {
	reader := csv.NewReader(r)
	reader.TrimLeadingSpace = true

	header, err := reader.Read()
	if err != nil {
		return Catalog{}, err
	}
	indexByName := make(map[string]int, len(header))
	for idx, name := range header {
		indexByName[strings.TrimSpace(name)] = idx
	}
	for _, name := range requiredColumns {
		if _, ok := indexByName[name]; !ok {
			return Catalog{}, ErrMissingColumns
		}
	}

	notesIdx, hasNotes := indexByName["notes"]

	catalog := Catalog{
		All:          make([]Record, 0),
		ByPartNumber: make(map[string]Record),
	}
	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return Catalog{}, err
		}
		record := Record{
			Brand:        strings.TrimSpace(valueAt(row, indexByName["brand"])),
			Family:       strings.TrimSpace(valueAt(row, indexByName["family"])),
			ProductModel: strings.TrimSpace(valueAt(row, indexByName["product_model"])),
			PartNumber:   strings.TrimSpace(valueAt(row, indexByName["part_number"])),
			Role:         strings.TrimSpace(valueAt(row, indexByName["role"])),
		}
		if hasNotes {
			record.Notes = strings.TrimSpace(valueAt(row, notesIdx))
		}
		catalog.All = append(catalog.All, record)
		if record.Brand == "" || record.Family == "" || record.ProductModel == "" || record.PartNumber == "" || record.Role == "" {
			continue
		}
		if _, exists := catalog.ByPartNumber[record.PartNumber]; exists {
			continue
		}
		catalog.ByPartNumber[record.PartNumber] = record
	}

	return catalog, nil
}

func valueAt(row []string, idx int) string {
	if idx < 0 || idx >= len(row) {
		return ""
	}
	return row[idx]
}

type ControllerCapability int

const (
	ControllerUnknown ControllerCapability = iota
	ControllerNone
	ControllerPresent
)

func (c ControllerCapability) String() string {
	switch c {
	case ControllerUnknown:
		return "ControllerUnknown"
	case ControllerNone:
		return "ControllerNone"
	case ControllerPresent:
		return "ControllerPresent"
	default:
		return "ControllerCapability(invalid)"
	}
}

func (c Catalog) ControllerCapability(partNumber string) ControllerCapability {
	partNumber = strings.TrimSpace(partNumber)
	if partNumber == "" {
		return ControllerUnknown
	}

	// All is the accepted catalog surface. ByPartNumber is deliberately
	// narrower: it contains only metadata-complete records used for richer
	// identity enrichment. Capability classification must not turn a published
	// part number into Unknown merely because unrelated metadata is incomplete.
	if len(c.All) != 0 {
		found := false
		capability := ControllerUnknown
		for _, record := range c.All {
			if strings.TrimSpace(record.PartNumber) != partNumber {
				continue
			}
			found = true
			role := strings.TrimSpace(record.Role)
			if role == "" {
				return ControllerUnknown
			}
			classified := controllerCapabilityForRole(role)
			if capability == ControllerUnknown {
				capability = classified
				continue
			}
			if capability != classified {
				// Conflicting duplicate catalog rows cannot authorize a positive or
				// negative classification.
				return ControllerUnknown
			}
		}
		if !found {
			return ControllerUnknown
		}
		return capability
	}

	// Preserve programmatic Catalog values created before All became the
	// capability source. Loaded catalogs always take the path above.
	record, found := c.ByPartNumber[partNumber]
	if !found {
		return ControllerUnknown
	}
	if strings.TrimSpace(record.Role) == "" {
		return ControllerUnknown
	}
	return controllerCapabilityForRole(record.Role)
}

func controllerCapabilityForRole(role string) ControllerCapability {
	if strings.EqualFold(strings.TrimSpace(role), "Regulator") ||
		strings.EqualFold(strings.TrimSpace(role), "Thermostat") {
		return ControllerPresent
	}
	return ControllerNone
}
