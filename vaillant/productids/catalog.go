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
	record, found := c.ByPartNumber[strings.TrimSpace(partNumber)]
	if !found {
		return ControllerUnknown
	}
	if strings.EqualFold(record.Role, "Regulator") {
		return ControllerPresent
	}
	return ControllerNone
}
