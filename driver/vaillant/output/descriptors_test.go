package output

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"reflect"
	"regexp"
	"strings"
	"testing"
)

type schemaField struct {
	Name string
	Type *schemaTypeRef
}

type schemaTypeRef struct {
	Name   string
	OfType *schemaTypeRef
}

func (ref *schemaTypeRef) named() string {
	for ref != nil {
		if ref.Name != "" {
			return ref.Name
		}
		ref = ref.OfType
	}
	return ""
}

func TestVaillantViewDescriptorsMatchServedFields(t *testing.T) {
	fields := readFixture[map[string]map[string]string](t, "graphql-go-fields.json")
	schema := readFixture[struct {
		Response struct {
			Schema struct {
				Types []struct {
					Name   string
					Fields []schemaField
				}
			} `json:"__schema"`
		}
	}](t, "query_schema_introspection.golden.json")
	inventory := readFixture[struct{ Tools []struct{ Name string } }](t, "mcp-tool-inventory.json")
	registered := make(map[string]bool)
	for _, tool := range inventory.Tools {
		registered[tool.Name] = true
	}
	introspection := make(map[string][]schemaField)
	for _, record := range schema.Response.Schema.Types {
		introspection[record.Name] = record.Fields
	}
	roots := frozenRootTools()
	served := reachableGraphQLFields(t, introspection, roots)
	fieldTools := make(map[string]string)
	for field, tool := range roots {
		for local := range reachableGraphQLFields(t, introspection, map[string]string{field: tool}) {
			fieldTools[local] = tool
		}
	}
	mcpTags := frozenMCPRecordTags(t)
	mcpAssigned := frozenMCPAssignments(t)
	seen := make(map[memberKey]bool)
	seenFields := make(map[string]bool)
	seenTags := make(map[string]map[string]string)
	for _, attr := range viewAttributes {
		if attr.Record != "query" {
			if want := descriptorRecordType(attr.Record); attr.GoType != want {
				t.Errorf("%s uses record %s, want %s", attr.Key.Local, attr.GoType, want)
			}
			if tag := mcpTags[attr.GoType][attr.GoField]; tag != "" {
				if seenTags[attr.GoType] == nil {
					seenTags[attr.GoType] = make(map[string]string)
				}
				if _, exists := seenTags[attr.GoType][attr.GoField]; exists {
					t.Errorf("duplicate MCP field %s.%s", attr.GoType, attr.GoField)
				}
				seenTags[attr.GoType][attr.GoField] = tag
			}
		}
		checkDescriptorKey(t, seen, attr.Key, attr.Record+"."+attr.Field)
		if attr.MCPTool != "" && !registered[attr.MCPTool] {
			t.Errorf("unregistered tool %q", attr.MCPTool)
		}
		if attr.Record == "query" {
			tool, ok := roots[attr.Field]
			key := "data"
			if attr.Field == "vaillant_regulator_capability" {
				key = "vaillant_regulator_capability"
			}
			if !ok || attr.MCPTool != tool || attr.MCPKey != key {
				t.Errorf("wrong root descriptor: %+v", attr)
			}
		} else if wantField := outputGoField(fields[attr.Record][attr.Field]); wantField != "" {
			if attr.GoField != wantField {
				t.Errorf("%s maps %s, want %s", attr.Key.Local, attr.GoField, wantField)
			}
			wantKey := mcpTags[attr.GoType][wantField]
			if wantKey != "" && !mcpAssigned[attr.GoType][wantField] {
				t.Errorf("no frozen adapter assignment for %s.%s", attr.GoType, wantField)
			}
			wantTool := ""
			if wantKey != "" {
				wantTool = fieldTools[attr.Key.Local]
			}
			if attr.MCPKey != wantKey || attr.MCPTool != wantTool {
				t.Errorf("%s MCP projection %q/%q, want %q/%q", attr.Key.Local, attr.MCPTool, attr.MCPKey, wantTool, wantKey)
			}
		} else {
			checkWriteAttribute(t, attr, mcpTags, mcpAssigned)
			continue
		}
		if !served[attr.Key.Local] {
			t.Errorf("non-served GraphQL field %s", attr.Key.Local)
		}
		seenFields[attr.Key.Local] = true
	}
	if !reflect.DeepEqual(seenFields, served) {
		for name := range served {
			if !seenFields[name] {
				t.Errorf("missing served field %s", name)
			}
		}
	}
	if !reflect.DeepEqual(seenTags, mcpTags) {
		for record, tags := range mcpTags {
			for field := range tags {
				if seenTags[record][field] == "" {
					t.Errorf("missing MCP field %s.%s", record, field)
				}
			}
		}
	}
	checkDescriptorNodesAndMethods(t, seen, registered, fields)
}

func descriptorRecordType(record string) string {
	switch record {
	case "circuit_managing_device":
		return "ManagingDevice"
	case "fm5_interpretation":
		return "FunctionModuleInterpretation"
	case "system_config_bound":
		return "ConfigBound"
	}
	parts := strings.Split(record, "_")
	for i, part := range parts {
		parts[i] = strings.ToUpper(part[:1]) + part[1:]
	}
	return strings.Join(parts, "")
}

// The fixtures retain the gateway's Go selectors; descriptors name the moved selectors.
func outputGoField(frozen string) string {
	switch frozen {
	case "ModuleConfigurationVR71":
		return "FunctionModuleConfigurationCode"
	case "OutdoorTemperatureAvg24h":
		return "OutdoorTemperatureDailyAverage"
	default:
		return frozen
	}
}

func outputGoRecord(frozen string) string {
	if frozen == "Fm5Interpretation" {
		return "FunctionModuleInterpretation"
	}
	return frozen
}

func frozenRootTools() map[string]string {
	tools := make(map[string]string)
	for _, name := range strings.Fields("zones dhw circuits radio_devices fm5_interpretation solar cylinders energy_totals boiler_status system schedules") {
		tools[name] = "ebus_v1_semantic_" + name + "_get"
	}
	tools["adapter_hardware_info"] = "ebus_v1_semantic_adapter_info_get"
	tools["fm5_semantic_mode"] = "ebus_v1_semantic_fm5_mode_get"
	tools["vaillant_regulator_capability"] = "ebus_v1_runtime_status_get"
	return tools
}

func reachableGraphQLFields(t *testing.T, schema map[string][]schemaField, roots map[string]string) map[string]bool {
	t.Helper()
	out, walked := make(map[string]bool), make(map[string]bool)
	var walk func(string)
	walk = func(record string) {
		if walked[record] {
			return
		}
		walked[record] = true
		for _, field := range schema[record] {
			out[record+"."+field.Name] = true
			walk(field.Type.named())
		}
	}
	for _, field := range schema["query"] {
		if _, ok := roots[field.Name]; ok {
			out["query."+field.Name] = true
			walk(field.Type.named())
		}
	}
	return out
}

func checkDescriptorKey(t *testing.T, seen map[memberKey]bool, key memberKey, local string) {
	t.Helper()
	if key != (memberKey{Namespace: "ebusreg", Local: local}) || seen[key] {
		t.Errorf("invalid or duplicate member key %+v; want local %q", key, local)
	}
	seen[key] = true
}

func checkWriteAttribute(t *testing.T, row viewField, tags map[string]map[string]string, assigned map[string]map[string]bool) {
	t.Helper()
	switch row.GoType {
	case "TimeProgramSlot":
		if row.Record != "time_program_slot" || row.Field != tags[row.GoType][row.GoField] || row.MCPKey != "" {
			t.Errorf("bad slot input: %+v", row)
		}
	case "TimeProgramWriteResult", "TimeProgramSlotResult":
		if row.Field != tags[row.GoType][row.GoField] || row.MCPKey != row.Field || !assigned[row.GoType][row.GoField] {
			t.Errorf("bad slot output: %+v", row)
		}
	default:
		t.Errorf("non-served descriptor: %+v", row)
	}
}

func checkDescriptorNodesAndMethods(t *testing.T, seen map[memberKey]bool, registered map[string]bool, fields map[string]map[string]string) {
	t.Helper()
	want := map[string]bool{"time_program_write_result": true, "time_program_slot_result": true, "time_program_slot": true}
	for name := range fields {
		want[name] = true
	}
	for _, node := range viewTreeNodes {
		checkDescriptorKey(t, seen, node.Key, node.Key.Local)
		if !want[node.Key.Local] {
			t.Errorf("unexpected tree node %+v", node)
		}
		delete(want, node.Key.Local)
	}
	if len(want) != 0 {
		t.Errorf("missing tree nodes: %v", want)
	}
	methods := map[string]bool{"ebus_v1_semantic_schedules_set_zone_time_program": true, "ebus_v1_semantic_schedules_set_dhw_time_program": true}
	for _, method := range viewMethods {
		checkDescriptorKey(t, seen, method.Key, method.MCPTool)
		if !methods[method.MCPTool] || !registered[method.MCPTool] || method.MCPKey != "data" || method.GoType != "TimeProgramWriteResult" {
			t.Errorf("invalid schedule method %+v", method)
		}
		delete(methods, method.MCPTool)
	}
	if len(methods) != 0 {
		t.Errorf("missing methods %v", methods)
	}
}

func fixtureGo(t *testing.T, name string) (*ast.File, string) {
	t.Helper()
	raw, err := os.ReadFile("testdata/" + name)
	if err != nil {
		t.Fatal(err)
	}
	text := "package frozen\n" + string(raw)
	file, err := parser.ParseFile(token.NewFileSet(), name, text, 0)
	if err != nil {
		t.Fatal(err)
	}
	return file, text
}

func frozenMCPRecordTags(t *testing.T) map[string]map[string]string {
	t.Helper()
	file, _ := fixtureGo(t, "mcp-record-tags.txt")
	tags := make(map[string]map[string]string)
	for _, decl := range file.Decls {
		spec := decl.(*ast.GenDecl).Specs[0].(*ast.TypeSpec)
		record := spec.Type.(*ast.StructType)
		name := outputGoRecord(spec.Name.Name)
		tags[name] = make(map[string]string)
		for _, field := range record.Fields.List {
			if field.Tag == nil {
				t.Fatalf("missing frozen tag for %s", field.Names[0].Name)
			}
			tag := reflect.StructTag(strings.Trim(field.Tag.Value, "`")).Get("json")
			tags[name][outputGoField(field.Names[0].Name)] = strings.SplitN(tag, ",", 2)[0]
		}
	}
	return tags
}

func frozenMCPAssignments(t *testing.T) map[string]map[string]bool {
	t.Helper()
	file, text := fixtureGo(t, "mcp-adapter-assignments.txt")
	assigned := make(map[string]map[string]bool)
	ast.Inspect(file, func(n ast.Node) bool {
		lit, ok := n.(*ast.CompositeLit)
		if !ok {
			return true
		}
		typ, ok := lit.Type.(*ast.SelectorExpr)
		if !ok || typ.X.(*ast.Ident).Name != "mcp" {
			return true
		}
		name := outputGoRecord(typ.Sel.Name)
		if assigned[name] == nil {
			assigned[name] = make(map[string]bool)
		}
		for _, entry := range lit.Elts {
			pair := entry.(*ast.KeyValueExpr)
			field := pair.Key.(*ast.Ident).Name
			// Same-name direct selectors or a nested MCP composite for this field.
			value := text[pair.Value.Pos()-1 : pair.Value.End()-1]
			assigned[name][outputGoField(field)] = regexp.MustCompile(`\.`+field+`\b`).MatchString(value) || strings.Contains(value, "mcp.")
		}
		return true
	})
	// The old adapter fills these fields by assignment and slice copy, not a literal.
	for record, fields := range map[string]string{
		"ScheduleStatus": "Programs", "ScheduleProgram": "Config SlotsUsed Days", "ScheduleDayProgram": "Slots",
		"AdapterHardwareInfo": "JumperFlags ResetCauseCode RestartCount LastIdentityQuery LastTelemetryQuery",
		"EnergySeries":        "Yearly Monthly YearlyMeta MonthlyMeta", "Fm5Interpretation": "DegradedReason",
		"TimeProgramWriteResult": "SlotResults",
	} {
		for _, field := range strings.Fields(fields) {
			assigned[outputGoRecord(record)][outputGoField(field)] = regexp.MustCompile(`\.` + field + `\s*=`).MatchString(text)
		}
	}
	return assigned
}
