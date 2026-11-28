package output

// CatalogueEnumeration is one family-codebook value kept beside its raw bytes.
type CatalogueEnumeration struct {
	Group       uint8  `json:"group"`
	Register    uint16 `json:"register"`
	Alternative string `json:"alternative,omitempty"`
	Raw         []byte `json:"raw,omitempty"`
	Code        *int   `json:"code,omitempty"`
	Label       string `json:"label,omitempty"`
	Level       string `json:"level,omitempty"`
}

// CatalogueScalar is one scalar register the family view decodes, with its raw bytes.
type CatalogueScalar struct {
	Name     string   `json:"name,omitempty"`
	Group    uint8    `json:"group"`
	Register uint16   `json:"register"`
	Raw      []byte   `json:"raw,omitempty"`
	Level    string   `json:"level,omitempty"`
	Unit     string   `json:"unit,omitempty"`
	Value    *float64 `json:"value,omitempty"`
	Tier     string   `json:"tier,omitempty"`
}

// VentilationView is the native ventilation view for the family that has the group.
type VentilationView struct {
	DayFanSpeedLimitRaw   []byte  `json:"ventilation_fan_speed_limit_maximum_at_day,omitempty"`
	NightFanSpeedLimitRaw []byte  `json:"ventilation_fan_speed_limit_maximum_at_night,omitempty"`
	DayFanSpeedLimit      *uint16 `json:"ventilation_fan_speed_limit_maximum_at_day_value,omitempty"`
	NightFanSpeedLimit    *uint16 `json:"ventilation_fan_speed_limit_maximum_at_night_value,omitempty"`
	Level                 string  `json:"level,omitempty"`
}

// FamilyCatalogue is the codebook, scalar and ventilation fill for one controller family.
type FamilyCatalogue struct {
	Enumerations  []CatalogueEnumeration `json:"enumerations,omitempty"`
	Scalars       []CatalogueScalar      `json:"scalars,omitempty"`
	Ventilation   *VentilationView       `json:"ventilation,omitempty"`
	RawOperations []RawOperation         `json:"raw_operations,omitempty"`
}

// BoilerSelectorEvidence is the raw selector, its annotation and the codec actually used.
type BoilerSelectorEvidence struct {
	Selector  uint16 `json:"selector"`
	Name      string `json:"name,omitempty"`
	Raw       []byte `json:"raw,omitempty"`
	Label     string `json:"label,omitempty"`
	Level     string `json:"level,omitempty"`
	Note      string `json:"note,omitempty"`
	CodecUsed string `json:"codec_used,omitempty"`
	Tier      string `json:"tier,omitempty"`
}

// RawOperation is a hand-called ebusreg operation. Nothing is sent until it is called.
type RawOperation struct {
	Namespace string `json:"namespace,omitempty"`
	Name      string `json:"name,omitempty"`
	Group     uint8  `json:"group"`
	Register  uint16 `json:"register"`
}

// EnergyDirection is one catalogued energy direction and any value the receiver read.
type EnergyDirection struct {
	Name       string   `json:"name,omitempty"`
	Family     string   `json:"family,omitempty"`
	Variant    string   `json:"variant,omitempty"`
	Field      string   `json:"field,omitempty"`
	Direction  string   `json:"direction,omitempty"`
	Period     string   `json:"period,omitempty"`
	Source     string   `json:"source,omitempty"`
	Usage      string   `json:"usage,omitempty"`
	Freshness  string   `json:"freshness,omitempty"`
	Level      string   `json:"level,omitempty"`
	Label      string   `json:"label,omitempty"`
	ValueKWh   *float64 `json:"value_kwh,omitempty"`
	Raw        []byte   `json:"raw,omitempty"`
	ScopeYear  int      `json:"scope_year,omitempty"`
	ScopeMonth int      `json:"scope_month,omitempty"`
	ObservedAt string   `json:"observed_at,omitempty"`
}

func cloneEnumerations(in []CatalogueEnumeration) []CatalogueEnumeration {
	if len(in) == 0 {
		return nil
	}
	out := make([]CatalogueEnumeration, len(in))
	for i, item := range in {
		out[i] = item
		out[i].Raw = append([]byte(nil), item.Raw...)
		if item.Code != nil {
			code := *item.Code
			out[i].Code = &code
		}
	}
	return out
}

func cloneEnergyDirections(in []EnergyDirection) []EnergyDirection {
	if len(in) == 0 {
		return nil
	}
	out := make([]EnergyDirection, len(in))
	for i, item := range in {
		out[i] = item
		out[i].Raw = append([]byte(nil), item.Raw...)
		if item.ValueKWh != nil {
			value := *item.ValueKWh
			out[i].ValueKWh = &value
		}
	}
	return out
}

func cloneEnergyResiduals(in []EnergyResidual) []EnergyResidual {
	if len(in) == 0 {
		return nil
	}
	out := make([]EnergyResidual, len(in))
	for i, item := range in {
		out[i] = item
		out[i].Raw = append([]byte(nil), item.Raw...)
		if item.Register != nil {
			register := *item.Register
			out[i].Register = &register
		}
	}
	return out
}

func cloneBoilerSelectors(in []BoilerSelectorEvidence) []BoilerSelectorEvidence {
	if len(in) == 0 {
		return nil
	}
	out := make([]BoilerSelectorEvidence, len(in))
	for i, item := range in {
		out[i] = item
		out[i].Raw = append([]byte(nil), item.Raw...)
	}
	return out
}

func cloneRawOperations(in []RawOperation) []RawOperation {
	if len(in) == 0 {
		return nil
	}
	out := make([]RawOperation, len(in))
	copy(out, in)
	return out
}

func cloneFamilyCatalogue(in *FamilyCatalogue) *FamilyCatalogue {
	if in == nil {
		return nil
	}
	out := *in
	out.Enumerations = cloneEnumerations(in.Enumerations)
	out.Scalars = cloneScalars(in.Scalars)
	out.Ventilation = cloneVentilation(in.Ventilation)
	out.RawOperations = cloneRawOperations(in.RawOperations)
	return &out
}

func cloneScalars(in []CatalogueScalar) []CatalogueScalar {
	if len(in) == 0 {
		return nil
	}
	out := make([]CatalogueScalar, len(in))
	for i, item := range in {
		out[i] = item
		out[i].Raw = append([]byte(nil), item.Raw...)
		if item.Value != nil {
			value := *item.Value
			out[i].Value = &value
		}
	}
	return out
}

func cloneVentilation(in *VentilationView) *VentilationView {
	if in == nil {
		return nil
	}
	out := *in
	out.DayFanSpeedLimitRaw = append([]byte(nil), in.DayFanSpeedLimitRaw...)
	out.NightFanSpeedLimitRaw = append([]byte(nil), in.NightFanSpeedLimitRaw...)
	if in.DayFanSpeedLimit != nil {
		value := *in.DayFanSpeedLimit
		out.DayFanSpeedLimit = &value
	}
	if in.NightFanSpeedLimit != nil {
		value := *in.NightFanSpeedLimit
		out.NightFanSpeedLimit = &value
	}
	return &out
}

// EnergyResidual is a residual or undefined energy field kept raw.
type EnergyResidual struct {
	Name     string  `json:"name,omitempty"`
	Family   string  `json:"family,omitempty"`
	Variant  string  `json:"variant,omitempty"`
	Field    string  `json:"field,omitempty"`
	Level    string  `json:"level,omitempty"`
	Raw      []byte  `json:"raw,omitempty"`
	Register *uint16 `json:"register,omitempty"`
	Note     string  `json:"note,omitempty"`
}

// SetEnumerations stores the zone codebook rows beside the zone state.
func (s *ZoneState) SetEnumerations(items []CatalogueEnumeration) {
	if s == nil {
		return
	}
	s.Enumerations = cloneEnumerations(items)
}

// SetEnumerations stores the hot-water codebook rows beside the hot-water state.
func (s *DhwState) SetEnumerations(items []CatalogueEnumeration) {
	if s == nil {
		return
	}
	s.Enumerations = cloneEnumerations(items)
}

// SetEnumerations stores the circuit codebook rows beside the circuit settings.
func (s *CircuitConfig) SetEnumerations(items []CatalogueEnumeration) {
	if s == nil {
		return
	}
	s.Enumerations = cloneEnumerations(items)
}

// SetVentilation stores the native ventilation view on the system record.
func (s *SystemStatus) SetVentilation(view *VentilationView) {
	if s == nil {
		return
	}
	s.Ventilation = cloneVentilation(view)
}

// SetCatalogue stores the family catalogue on the system record.
func (s *SystemStatus) SetCatalogue(view *FamilyCatalogue) {
	if s == nil {
		return
	}
	s.Catalogue = cloneFamilyCatalogue(view)
}

// SetRegisterDirections stores the register energy directions.
func (s *EnergyTotals) SetRegisterDirections(items []EnergyDirection) {
	if s == nil {
		return
	}
	s.RegisterDirections = cloneEnergyDirections(items)
}

// SetCatalogueDirections stores the catalogue energy directions.
func (s *EnergyTotals) SetCatalogueDirections(items []EnergyDirection) {
	if s == nil {
		return
	}
	s.CatalogueDirections = cloneEnergyDirections(items)
}

// SetResidualFields stores the residual energy fields.
func (s *EnergyTotals) SetResidualFields(items []EnergyResidual) {
	if s == nil {
		return
	}
	s.ResidualFields = cloneEnergyResiduals(items)
}

// SetCandidateFields stores the candidate energy fields.
func (s *EnergyTotals) SetCandidateFields(items []EnergyResidual) {
	if s == nil {
		return
	}
	s.CandidateFields = cloneEnergyResiduals(items)
}

// SetSelectors stores the boiler selector evidence.
func (s *BoilerStatus) SetSelectors(items []BoilerSelectorEvidence) {
	if s == nil {
		return
	}
	s.Selectors = cloneBoilerSelectors(items)
}

// SetOperations stores the admitted raw operations.
func (s *BoilerStatus) SetOperations(items []RawOperation) {
	if s == nil {
		return
	}
	s.Operations = cloneRawOperations(items)
}

// SetEnumerations stores the zone codebook rows on the MCP zone state.
func (s *MCPZoneState) SetEnumerations(items []CatalogueEnumeration) {
	if s == nil {
		return
	}
	s.Enumerations = cloneEnumerations(items)
}

// SetEnumerations stores the hot-water codebook rows on the MCP hot-water state.
func (s *MCPDhwState) SetEnumerations(items []CatalogueEnumeration) {
	if s == nil {
		return
	}
	s.Enumerations = cloneEnumerations(items)
}

// SetEnumerations stores the circuit codebook rows on the MCP circuit settings.
func (s *MCPCircuitConfig) SetEnumerations(items []CatalogueEnumeration) {
	if s == nil {
		return
	}
	s.Enumerations = cloneEnumerations(items)
}

// SetVentilation stores the native ventilation view on the MCP system record.
func (s *MCPSystemStatus) SetVentilation(view *VentilationView) {
	if s == nil {
		return
	}
	s.Ventilation = cloneVentilation(view)
}

// SetCatalogue stores the family catalogue on the MCP system record.
func (s *MCPSystemStatus) SetCatalogue(view *FamilyCatalogue) {
	if s == nil {
		return
	}
	s.Catalogue = cloneFamilyCatalogue(view)
}

// SetRegisterDirections stores the MCP register energy directions.
func (s *MCPEnergyTotals) SetRegisterDirections(items []EnergyDirection) {
	if s == nil {
		return
	}
	s.RegisterDirections = cloneEnergyDirections(items)
}

// SetCatalogueDirections stores the MCP catalogue energy directions.
func (s *MCPEnergyTotals) SetCatalogueDirections(items []EnergyDirection) {
	if s == nil {
		return
	}
	s.CatalogueDirections = cloneEnergyDirections(items)
}

// SetResidualFields stores the MCP residual energy fields.
func (s *MCPEnergyTotals) SetResidualFields(items []EnergyResidual) {
	if s == nil {
		return
	}
	s.ResidualFields = cloneEnergyResiduals(items)
}

// SetCandidateFields stores the MCP candidate energy fields.
func (s *MCPEnergyTotals) SetCandidateFields(items []EnergyResidual) {
	if s == nil {
		return
	}
	s.CandidateFields = cloneEnergyResiduals(items)
}

// SetSelectors stores the MCP boiler selector evidence.
func (s *MCPBoilerStatus) SetSelectors(items []BoilerSelectorEvidence) {
	if s == nil {
		return
	}
	s.Selectors = cloneBoilerSelectors(items)
}

// SetOperations stores the MCP admitted raw operations.
func (s *MCPBoilerStatus) SetOperations(items []RawOperation) {
	if s == nil {
		return
	}
	s.Operations = cloneRawOperations(items)
}
