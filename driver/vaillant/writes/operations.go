// Package writes is the driver's write inventory and the scalar executors.
// A scalar operation sends one frame and reads the same selector once.
// A dispatch that is not admitted sends nothing.
package writes

import (
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"sync"
)

//go:embed testdata/pre_port_writes.json
var inventoryFile []byte

// Selector is the native write selector recorded for one operation.
// Instance is "binding" when the dispatch names it, "unknown" when no write
// instance is established, or two hex digits when the selector fixes it.
type Selector struct {
	Family     string `json:"family"`
	Primary    string `json:"primary"`
	Secondary  string `json:"secondary"`
	Opcode     string `json:"opcode"`
	Group      string `json:"group"`
	Register   string `json:"register,omitempty"`
	Instance   string `json:"instance"`
	ReadOpcode string `json:"read_opcode,omitempty"`
}

// Row is one pre-port write surface bound to one operation.
type Row struct {
	Surface         string   `json:"surface"`
	Key             string   `json:"key,omitempty"`
	Arguments       []string `json:"arguments,omitempty"`
	Slots           string   `json:"slots,omitempty"`
	SlotTemperature string   `json:"slot_temperature,omitempty"`
	Selector        Selector `json:"selector"`
	Operation       string   `json:"operation"`
	Form            string   `json:"form"`
	Kind            string   `json:"kind,omitempty"`
	Argument        string   `json:"argument,omitempty"`
	Grade           string   `json:"grade"`
	PairRegister    string   `json:"pair_register,omitempty"`
}

type inventoryDocument struct {
	Rows []Row `json:"rows"`
}

var (
	inventoryOnce  sync.Once
	inventoryRows  []Row
	inventoryError error
)

// ParseInventory checks raw and returns its rows.
// Rows that share an operation share one selector, form, kind, argument, grade
// and pair register. A repeated surface and key is rejected.
func ParseInventory(raw []byte) ([]Row, error) {
	var doc inventoryDocument
	if err := json.Unmarshal(raw, &doc); err != nil {
		return nil, err
	}
	if err := validateRows(doc.Rows); err != nil {
		return nil, err
	}
	return cloneRows(doc.Rows), nil
}

// Inventory returns the embedded pre-port write rows.
func Inventory() ([]Row, error) {
	rows, err := cachedInventory()
	if err != nil {
		return nil, err
	}
	return cloneRows(rows), nil
}

// Operation returns the first row bound to id.
func Operation(id string) (Row, bool) {
	rows, err := cachedInventory()
	if err != nil || id == "" {
		return Row{}, false
	}
	for _, row := range rows {
		if row.Operation == id {
			row.Arguments = append([]string(nil), row.Arguments...)
			return row, true
		}
	}
	return Row{}, false
}

func cachedInventory() ([]Row, error) {
	inventoryOnce.Do(func() {
		inventoryRows, inventoryError = ParseInventory(inventoryFile)
	})
	return inventoryRows, inventoryError
}

func cloneRows(rows []Row) []Row {
	out := make([]Row, len(rows))
	for i, row := range rows {
		row.Arguments = append([]string(nil), row.Arguments...)
		out[i] = row
	}
	return out
}

func validateRows(rows []Row) error {
	if len(rows) == 0 {
		return errors.New("write inventory has no rows")
	}
	seen := map[string]struct{}{}
	identity := map[string]Row{}
	for i, row := range rows {
		if err := validateRow(row); err != nil {
			return fmt.Errorf("rows[%d]: %w", i, err)
		}
		surfaceKey := row.Surface + "\x00" + row.Key
		if _, ok := seen[surfaceKey]; ok {
			return fmt.Errorf("rows[%d]: duplicate surface %q key %q", i, row.Surface, row.Key)
		}
		seen[surfaceKey] = struct{}{}
		previous, ok := identity[row.Operation]
		if ok && !sameOperation(previous, row) {
			return fmt.Errorf("rows[%d]: operation %q disagrees with an earlier row", i, row.Operation)
		}
		identity[row.Operation] = row
	}
	return nil
}

func validateRow(row Row) error {
	if row.Surface == "" || row.Operation == "" {
		return errors.New("surface and operation are required")
	}
	switch row.Form {
	case "scalar":
		if row.Kind == "" || row.Argument == "" || !knownKind(row.Kind) {
			return errors.New("scalar row requires a known kind and argument")
		}
	case "compound", "unknown":
		if row.Kind != "" {
			return errors.New("only a scalar row carries a kind")
		}
	default:
		return errors.New("form must be scalar, compound or unknown")
	}
	switch row.Grade {
	case "hypothesis", "unknown":
	default:
		return errors.New("grade must be hypothesis or unknown")
	}
	if row.Form != "unknown" && row.Selector.Instance == "unknown" {
		return errors.New("instance unknown is only valid when the form is unknown")
	}
	if err := validateSelector(row.Selector, row.Form); err != nil {
		return err
	}
	if row.PairRegister != "" {
		if err := hexDigits(row.PairRegister, 4); err != nil {
			return err
		}
	}
	return nil
}

func knownKind(kind string) bool {
	switch kind {
	case "kelvin_celsius", "system_switch", "room_control", "bool_byte", "native_float32", "native_uint16", "percent_uint16":
		return true
	default:
		return false
	}
}

func validateSelector(sel Selector, form string) error {
	switch sel.Family {
	case "local", "boiler", "timer":
	default:
		return errors.New("selector family")
	}
	if err := hexDigits(sel.Primary, 2); err != nil {
		return err
	}
	if err := hexDigits(sel.Secondary, 2); err != nil {
		return err
	}
	if err := hexDigits(sel.Opcode, 2); err != nil {
		return err
	}
	if err := hexDigits(sel.Group, 2); err != nil {
		return err
	}
	switch sel.Instance {
	case "binding", "unknown":
	default:
		if err := hexDigits(sel.Instance, 2); err != nil {
			return err
		}
	}
	if sel.Family == "timer" {
		if form != "compound" || sel.Register != "" || sel.ReadOpcode != "a5" || sel.Opcode != "a6" {
			return errors.New("timer selector")
		}
		return nil
	}
	if sel.ReadOpcode != "" || sel.Register == "" {
		return errors.New("register selector")
	}
	return hexDigits(sel.Register, 4)
}

func sameOperation(left, right Row) bool {
	return left.Selector == right.Selector && left.Form == right.Form && left.Kind == right.Kind &&
		left.Argument == right.Argument && left.Grade == right.Grade && left.PairRegister == right.PairRegister
}

func hexDigits(value string, width int) error {
	if len(value) != width {
		return fmt.Errorf("hex width of %q", value)
	}
	if _, err := strconv.ParseUint(value, 16, width*4); err != nil {
		return fmt.Errorf("hex %q: %w", value, err)
	}
	for _, char := range value {
		if char < '0' || (char > '9' && char < 'a') || char > 'f' {
			return fmt.Errorf("hex %q", value)
		}
	}
	return nil
}

type selectorNumbers struct {
	opcode, primary, secondary, group byte
	register                          uint16
}

func parseSelector(sel Selector) (selectorNumbers, bool) {
	var numbers selectorNumbers
	var ok bool
	if numbers.opcode, ok = parseByte(sel.Opcode); !ok {
		return selectorNumbers{}, false
	}
	if numbers.primary, ok = parseByte(sel.Primary); !ok {
		return selectorNumbers{}, false
	}
	if numbers.secondary, ok = parseByte(sel.Secondary); !ok {
		return selectorNumbers{}, false
	}
	if numbers.group, ok = parseByte(sel.Group); !ok {
		return selectorNumbers{}, false
	}
	if sel.Register == "" {
		return numbers, true
	}
	register, err := strconv.ParseUint(sel.Register, 16, 16)
	if err != nil {
		return selectorNumbers{}, false
	}
	numbers.register = uint16(register)
	return numbers, true
}

func parseByte(value string) (byte, bool) {
	parsed, err := strconv.ParseUint(value, 16, 8)
	if err != nil || len(value) != 2 {
		return 0, false
	}
	return byte(parsed), true
}
