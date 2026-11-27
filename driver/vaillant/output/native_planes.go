package output

// HistorySlot is one stock history slot. Empty marks the stock empty sentinel.
type HistorySlot struct {
	Value uint16 `json:"value"`
	Empty bool   `json:"empty,omitempty"`
}

// HistoryRead is one stock error-history or service-history reply.
// Flag and category qualification is not derived from these bytes.
type HistoryRead struct {
	Kind           string        `json:"kind,omitempty"`
	RequestedIndex *uint8        `json:"requested_index,omitempty"`
	Raw            []byte        `json:"raw,omitempty"`
	ReturnedIndex  *uint8        `json:"returned_index,omitempty"`
	Slots          []HistorySlot `json:"slots,omitempty"`
	DecodeError    string        `json:"decode_error,omitempty"`
	FlagsKnown     bool          `json:"flags_known,omitempty"`
	CategoryKnown  bool          `json:"category_known,omitempty"`
	Ambiguity      string        `json:"ambiguity,omitempty"`
}

// HistoryStop is the first failing index of a history list.
type HistoryStop struct {
	Index         uint8         `json:"index"`
	Code          string        `json:"code,omitempty"`
	Detail        string        `json:"detail,omitempty"`
	Raw           []byte        `json:"raw,omitempty"`
	ReturnedIndex *uint8        `json:"returned_index,omitempty"`
	Slots         []HistorySlot `json:"slots,omitempty"`
	DecodedSlot   *uint16       `json:"decoded_slot,omitempty"`
}

// HistoryList is a bounded error-history list that stops at the first failing index.
type HistoryList struct {
	Kind  string        `json:"kind,omitempty"`
	Limit int           `json:"limit,omitempty"`
	Items []HistoryRead `json:"items,omitempty"`
	Stop  *HistoryStop  `json:"stop,omitempty"`
}

// HistoryFlagBit is one bit of an explicitly qualified history flag byte.
type HistoryFlagBit struct {
	Index    int    `json:"index"`
	Name     string `json:"name,omitempty"`
	Set      bool   `json:"set,omitempty"`
	Reserved bool   `json:"reserved,omitempty"`
}

// HistoryFlagReading is a field-level history qualification supplied apart from a stock reply.
type HistoryFlagReading struct {
	Field         string           `json:"field,omitempty"`
	Raw           []byte           `json:"raw,omitempty"`
	Label         string           `json:"label,omitempty"`
	Level         string           `json:"level,omitempty"`
	Unresolved    bool             `json:"unresolved,omitempty"`
	Qualification string           `json:"qualification,omitempty"`
	Bits          []HistoryFlagBit `json:"bits,omitempty"`
}

// HistoryCategoryReading is one text-category alternative for an explicit variant.
type HistoryCategoryReading struct {
	Variant    string `json:"variant,omitempty"`
	Raw        byte   `json:"raw"`
	Label      string `json:"label,omitempty"`
	Named      bool   `json:"named,omitempty"`
	Labelled   bool   `json:"labelled,omitempty"`
	Level      string `json:"level,omitempty"`
	Unresolved bool   `json:"unresolved,omitempty"`
	Ambiguity  string `json:"ambiguity,omitempty"`
}

// HistoryBranchComponent is one calendar cell state from an explicit branch input.
type HistoryBranchComponent struct {
	State string `json:"state,omitempty"`
}

// HistoryBranchReading is the branch selected from an explicit marker and components.
type HistoryBranchReading struct {
	Marker     uint8                    `json:"marker"`
	State      string                   `json:"state,omitempty"`
	Meaning    string                   `json:"meaning,omitempty"`
	Complete   bool                     `json:"complete,omitempty"`
	Components []HistoryBranchComponent `json:"components,omitempty"`
}

// HistoryOwnerRecord is the one history record both surfaces serve.
type HistoryOwnerRecord struct {
	Reads      []HistoryRead            `json:"reads,omitempty"`
	List       *HistoryList             `json:"list,omitempty"`
	Flags      []HistoryFlagReading     `json:"flags,omitempty"`
	Categories []HistoryCategoryReading `json:"categories,omitempty"`
	Branches   []HistoryBranchReading   `json:"branches,omitempty"`
	Ambiguity  string                   `json:"ambiguity,omitempty"`
}

// GraphQL returns the history record served on the GraphQL surface.
func (r HistoryOwnerRecord) GraphQL() HistoryOwnerRecord {
	return cloneHistoryOwner(r)
}

// MCP returns the history record served on the MCP surface.
func (r HistoryOwnerRecord) MCP() HistoryOwnerRecord {
	return cloneHistoryOwner(r)
}

// ModuleReading is one raw function-module sensor, actuator or mixer value.
type ModuleReading struct {
	Name       string `json:"name,omitempty"`
	Operation  byte   `json:"operation,omitempty"`
	Raw        []byte `json:"raw,omitempty"`
	Level      string `json:"level,omitempty"`
	Placement  string `json:"placement,omitempty"`
	Unresolved string `json:"unresolved,omitempty"`
}

// ModuleOperation is a hand-called function-module operation in the ebusreg namespace.
type ModuleOperation struct {
	Namespace    string   `json:"namespace,omitempty"`
	Name         string   `json:"name,omitempty"`
	Operation    byte     `json:"operation,omitempty"`
	Circuit      *byte    `json:"circuit,omitempty"`
	Unresolved   []string `json:"unresolved,omitempty"`
	Alternatives []string `json:"alternatives,omitempty"`
}

// FunctionModuleRecord is the native view of one module identified by its identification reply.
type FunctionModuleRecord struct {
	DeviceID          string            `json:"device_id,omitempty"`
	Address           byte              `json:"address,omitempty"`
	Manufacturer      byte              `json:"manufacturer,omitempty"`
	Software          []byte            `json:"software,omitempty"`
	Hardware          []byte            `json:"hardware,omitempty"`
	RawIdentification []byte            `json:"raw_identification,omitempty"`
	Family            string            `json:"family,omitempty"`
	Sensors           []ModuleReading   `json:"sensors,omitempty"`
	Actuators         []ModuleReading   `json:"actuators,omitempty"`
	Mixer             []ModuleReading   `json:"mixer,omitempty"`
	Operations        []ModuleOperation `json:"operations,omitempty"`
	Ambiguity         string            `json:"ambiguity,omitempty"`
}

// GraphQL returns the function-module record served on the GraphQL surface.
func (r FunctionModuleRecord) GraphQL() FunctionModuleRecord {
	return cloneFunctionModule(r)
}

// MCP returns the function-module record served on the MCP surface.
func (r FunctionModuleRecord) MCP() FunctionModuleRecord {
	return cloneFunctionModule(r)
}

// LegacyField is one decoded plane field kept beside its raw bytes.
type LegacyField struct {
	Name  string   `json:"name,omitempty"`
	Raw   []byte   `json:"raw,omitempty"`
	Value *float64 `json:"value,omitempty"`
	Valid bool     `json:"valid,omitempty"`
	Level string   `json:"level,omitempty"`
	Unit  string   `json:"unit,omitempty"`
}

// LegacyRead is one plane method result.
type LegacyRead struct {
	Method     string        `json:"method,omitempty"`
	Raw        []byte        `json:"raw,omitempty"`
	Fields     []LegacyField `json:"fields,omitempty"`
	Unresolved string        `json:"unresolved,omitempty"`
	Sent       bool          `json:"sent,omitempty"`
}

// LegacyPlane is the status, parameters and statistics of one legacy plane.
type LegacyPlane struct {
	Name  string       `json:"name,omitempty"`
	Reads []LegacyRead `json:"reads,omitempty"`
}

// LegacyOperation is a hand-called legacy-plane operation in the ebusreg namespace.
type LegacyOperation struct {
	Namespace  string `json:"namespace,omitempty"`
	Plane      string `json:"plane,omitempty"`
	Name       string `json:"name,omitempty"`
	Unresolved string `json:"unresolved,omitempty"`
}

// LegacyPlaneRecord is the native view of heating, hot water and a standalone solar station.
type LegacyPlaneRecord struct {
	ExtendedRegister bool              `json:"extended_register,omitempty"`
	Planes           []LegacyPlane     `json:"planes,omitempty"`
	Operations       []LegacyOperation `json:"operations,omitempty"`
	Catalogues       []string          `json:"catalogues,omitempty"`
	Ambiguity        string            `json:"ambiguity,omitempty"`
}

// GraphQL returns the legacy-plane record served on the GraphQL surface.
func (r LegacyPlaneRecord) GraphQL() LegacyPlaneRecord {
	return cloneLegacyPlaneRecord(r)
}

// MCP returns the legacy-plane record served on the MCP surface.
func (r LegacyPlaneRecord) MCP() LegacyPlaneRecord {
	return cloneLegacyPlaneRecord(r)
}

// SetHistory stores the shared history record on the system record.
func (s *SystemStatus) SetHistory(record *HistoryOwnerRecord) {
	if s == nil {
		return
	}
	s.history = cloneHistoryPtr(record)
}

// History returns the shared history record.
func (s SystemStatus) History() *HistoryOwnerRecord {
	return cloneHistoryPtr(s.history)
}

// SetFunctionModule stores the shared function-module record on the system record.
func (s *SystemStatus) SetFunctionModule(record *FunctionModuleRecord) {
	if s == nil {
		return
	}
	s.functionModule = cloneFunctionModulePtr(record)
}

// FunctionModule returns the shared function-module record.
func (s SystemStatus) FunctionModule() *FunctionModuleRecord {
	return cloneFunctionModulePtr(s.functionModule)
}

// SetLegacyPlanes stores the shared legacy-plane record on the system record.
func (s *SystemStatus) SetLegacyPlanes(record *LegacyPlaneRecord) {
	if s == nil {
		return
	}
	s.legacyPlanes = cloneLegacyPlanePtr(record)
}

// LegacyPlanes returns the shared legacy-plane record.
func (s SystemStatus) LegacyPlanes() *LegacyPlaneRecord {
	return cloneLegacyPlanePtr(s.legacyPlanes)
}

// SetHistory stores the shared history record on the MCP system record.
func (s *MCPSystemStatus) SetHistory(record *HistoryOwnerRecord) {
	if s == nil {
		return
	}
	s.history = cloneHistoryPtr(record)
}

// History returns the MCP shared history record.
func (s MCPSystemStatus) History() *HistoryOwnerRecord {
	return cloneHistoryPtr(s.history)
}

// SetFunctionModule stores the shared function-module record on the MCP system record.
func (s *MCPSystemStatus) SetFunctionModule(record *FunctionModuleRecord) {
	if s == nil {
		return
	}
	s.functionModule = cloneFunctionModulePtr(record)
}

// FunctionModule returns the MCP shared function-module record.
func (s MCPSystemStatus) FunctionModule() *FunctionModuleRecord {
	return cloneFunctionModulePtr(s.functionModule)
}

// SetLegacyPlanes stores the shared legacy-plane record on the MCP system record.
func (s *MCPSystemStatus) SetLegacyPlanes(record *LegacyPlaneRecord) {
	if s == nil {
		return
	}
	s.legacyPlanes = cloneLegacyPlanePtr(record)
}

// LegacyPlanes returns the MCP shared legacy-plane record.
func (s MCPSystemStatus) LegacyPlanes() *LegacyPlaneRecord {
	return cloneLegacyPlanePtr(s.legacyPlanes)
}

func cloneHistoryPtr(in *HistoryOwnerRecord) *HistoryOwnerRecord {
	if in == nil {
		return nil
	}
	out := cloneHistoryOwner(*in)
	return &out
}

func cloneHistoryOwner(in HistoryOwnerRecord) HistoryOwnerRecord {
	out := in
	out.Reads = cloneHistoryReads(in.Reads)
	out.List = cloneHistoryList(in.List)
	out.Flags = cloneHistoryFlags(in.Flags)
	out.Categories = cloneHistoryCategories(in.Categories)
	out.Branches = cloneHistoryBranches(in.Branches)
	return out
}

func cloneHistoryReads(in []HistoryRead) []HistoryRead {
	if len(in) == 0 {
		return nil
	}
	out := make([]HistoryRead, len(in))
	for i, item := range in {
		out[i] = item
		out[i].Raw = append([]byte(nil), item.Raw...)
		out[i].RequestedIndex = cloneByte(item.RequestedIndex)
		out[i].ReturnedIndex = cloneByte(item.ReturnedIndex)
		out[i].Slots = cloneHistorySlots(item.Slots)
	}
	return out
}

func cloneHistorySlots(in []HistorySlot) []HistorySlot {
	if len(in) == 0 {
		return nil
	}
	out := make([]HistorySlot, len(in))
	copy(out, in)
	return out
}

func cloneHistoryList(in *HistoryList) *HistoryList {
	if in == nil {
		return nil
	}
	out := *in
	out.Items = cloneHistoryReads(in.Items)
	out.Stop = cloneHistoryStop(in.Stop)
	return &out
}

func cloneHistoryStop(in *HistoryStop) *HistoryStop {
	if in == nil {
		return nil
	}
	out := *in
	out.Raw = append([]byte(nil), in.Raw...)
	out.ReturnedIndex = cloneByte(in.ReturnedIndex)
	out.Slots = cloneHistorySlots(in.Slots)
	if in.DecodedSlot != nil {
		value := *in.DecodedSlot
		out.DecodedSlot = &value
	}
	return &out
}

func cloneHistoryFlags(in []HistoryFlagReading) []HistoryFlagReading {
	if len(in) == 0 {
		return nil
	}
	out := make([]HistoryFlagReading, len(in))
	for i, item := range in {
		out[i] = item
		out[i].Raw = append([]byte(nil), item.Raw...)
		if len(item.Bits) > 0 {
			out[i].Bits = append([]HistoryFlagBit(nil), item.Bits...)
		}
	}
	return out
}

func cloneHistoryCategories(in []HistoryCategoryReading) []HistoryCategoryReading {
	if len(in) == 0 {
		return nil
	}
	out := make([]HistoryCategoryReading, len(in))
	copy(out, in)
	return out
}

func cloneHistoryBranches(in []HistoryBranchReading) []HistoryBranchReading {
	if len(in) == 0 {
		return nil
	}
	out := make([]HistoryBranchReading, len(in))
	for i, item := range in {
		out[i] = item
		if len(item.Components) > 0 {
			out[i].Components = append([]HistoryBranchComponent(nil), item.Components...)
		}
	}
	return out
}

func cloneByte(in *uint8) *uint8 {
	if in == nil {
		return nil
	}
	value := *in
	return &value
}

func cloneFunctionModulePtr(in *FunctionModuleRecord) *FunctionModuleRecord {
	if in == nil {
		return nil
	}
	out := cloneFunctionModule(*in)
	return &out
}

func cloneFunctionModule(in FunctionModuleRecord) FunctionModuleRecord {
	out := in
	out.Software = append([]byte(nil), in.Software...)
	out.Hardware = append([]byte(nil), in.Hardware...)
	out.RawIdentification = append([]byte(nil), in.RawIdentification...)
	out.Sensors = cloneModuleReadings(in.Sensors)
	out.Actuators = cloneModuleReadings(in.Actuators)
	out.Mixer = cloneModuleReadings(in.Mixer)
	out.Operations = cloneModuleOperations(in.Operations)
	return out
}

func cloneModuleReadings(in []ModuleReading) []ModuleReading {
	if len(in) == 0 {
		return nil
	}
	out := make([]ModuleReading, len(in))
	for i, item := range in {
		out[i] = item
		out[i].Raw = append([]byte(nil), item.Raw...)
	}
	return out
}

func cloneModuleOperations(in []ModuleOperation) []ModuleOperation {
	if len(in) == 0 {
		return nil
	}
	out := make([]ModuleOperation, len(in))
	for i, item := range in {
		out[i] = item
		if item.Circuit != nil {
			value := *item.Circuit
			out[i].Circuit = &value
		}
		if len(item.Unresolved) > 0 {
			out[i].Unresolved = append([]string(nil), item.Unresolved...)
		}
		if len(item.Alternatives) > 0 {
			out[i].Alternatives = append([]string(nil), item.Alternatives...)
		}
	}
	return out
}

func cloneLegacyPlanePtr(in *LegacyPlaneRecord) *LegacyPlaneRecord {
	if in == nil {
		return nil
	}
	out := cloneLegacyPlaneRecord(*in)
	return &out
}

func cloneLegacyPlaneRecord(in LegacyPlaneRecord) LegacyPlaneRecord {
	out := in
	if len(in.Planes) > 0 {
		out.Planes = make([]LegacyPlane, len(in.Planes))
		for i, plane := range in.Planes {
			out.Planes[i] = plane
			out.Planes[i].Reads = cloneLegacyReads(plane.Reads)
		}
	}
	if len(in.Operations) > 0 {
		out.Operations = append([]LegacyOperation(nil), in.Operations...)
	}
	if len(in.Catalogues) > 0 {
		out.Catalogues = append([]string(nil), in.Catalogues...)
	}
	return out
}

func cloneLegacyReads(in []LegacyRead) []LegacyRead {
	if len(in) == 0 {
		return nil
	}
	out := make([]LegacyRead, len(in))
	for i, item := range in {
		out[i] = item
		out[i].Raw = append([]byte(nil), item.Raw...)
		out[i].Fields = cloneLegacyFields(item.Fields)
	}
	return out
}

func cloneLegacyFields(in []LegacyField) []LegacyField {
	if len(in) == 0 {
		return nil
	}
	out := make([]LegacyField, len(in))
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
