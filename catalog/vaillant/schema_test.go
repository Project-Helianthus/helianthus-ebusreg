package vaillant

import (
	"encoding/json"
	"os"
	"reflect"
	"testing"
)

func schemaFixture(schemaT *testing.T) []SchemaRecord {
	schemaT.Helper()
	schemaData, schemaError := os.ReadFile("testdata/ab05_synthetic_metadata.json")
	if schemaError != nil {
		schemaT.Fatal(schemaError)
	}
	var schemaRecords []SchemaRecord
	if schemaError = json.Unmarshal(schemaData, &schemaRecords); schemaError != nil {
		schemaT.Fatal(schemaError)
	}
	return schemaRecords
}

func schemaInt(schemaValue int) *int { return &schemaValue }

func TestAb05KeyDimensions(schemaT *testing.T) {
	schemaCases := []struct {
		schemaName   string
		schemaChange func(*SchemaKey)
	}{
		{"ab05_family", func(schemaKey *SchemaKey) { schemaKey.SchemaFamily = "beta" }},
		{"ab05_variant", func(schemaKey *SchemaKey) { schemaKey.SchemaVariant = "ab05_secondary" }},
		{"ab05_direction", func(schemaKey *SchemaKey) { schemaKey.SchemaDirection = "request" }},
		{"ab05_opcode", func(schemaKey *SchemaKey) { schemaKey.SchemaOpcode = 6 }},
		{"ab05_group", func(schemaKey *SchemaKey) { schemaKey.SchemaGroup = 4 }},
		{"ab05_instance", func(schemaKey *SchemaKey) { schemaKey.SchemaInstance = 2 }},
		{"ab05_register", func(schemaKey *SchemaKey) { schemaKey.SchemaRegister = 35 }},
		{"ab05_alternative", func(schemaKey *SchemaKey) { schemaKey.SchemaAlternative = "ab05_other" }},
	}
	for _, schemaCase := range schemaCases {
		schemaT.Run(schemaCase.schemaName, func(schemaT *testing.T) {
			schemaFirst := schemaFixture(schemaT)[0]
			schemaSecond := schemaFirst
			schemaCase.schemaChange(&schemaSecond.SchemaKey)
			schemaSecond.SchemaRawFallback = []byte{42}
			schemaRecords := []SchemaRecord{schemaFirst, schemaSecond}
			if schemaError := SchemaValidate(schemaRecords); schemaError != nil {
				schemaT.Fatal(schemaError)
			}
			for schemaIndex, schemaWant := range schemaRecords {
				schemaGot, schemaFound, schemaError := SchemaQuery(schemaRecords, schemaWant.SchemaKey)
				if schemaError != nil || !schemaFound || !reflect.DeepEqual(schemaGot, schemaWant) {
					schemaT.Fatalf("ab05 record %d: got %#v, found %v, error %v", schemaIndex, schemaGot, schemaFound, schemaError)
				}
			}
			schemaSecond.SchemaKey = schemaFirst.SchemaKey
			if SchemaValidate([]SchemaRecord{schemaFirst, schemaSecond}) == nil {
				schemaT.Fatal("ab05 duplicate accepted")
			}
		})
	}
}

func TestAb05AlternativesAndOverlap(schemaT *testing.T) {
	schemaRecords := schemaFixture(schemaT)
	for _, schemaIndex := range []int{0, 1} {
		schemaGot, schemaFound, schemaError := SchemaQuery(schemaRecords, schemaRecords[schemaIndex].SchemaKey)
		if schemaError != nil || !schemaFound || !reflect.DeepEqual(schemaGot, schemaRecords[schemaIndex]) {
			schemaT.Fatalf("ab05 alternative flattened or overlap lost: %#v %v", schemaGot, schemaError)
		}
	}
	if len(schemaRecords[0].SchemaFields) != 2 || *schemaRecords[0].SchemaFields[0].SchemaWidth != 2 || *schemaRecords[1].SchemaFields[0].SchemaWidth != 1 {
		schemaT.Fatal("ab05 fixture lost distinct overlapping candidates")
	}
}

func TestAb05AbsentExtents(schemaT *testing.T) {
	schemaRecords := schemaFixture(schemaT)
	schemaHeader, schemaFound, schemaError := SchemaQuery(schemaRecords, schemaRecords[2].SchemaKey)
	if schemaError != nil || !schemaFound || schemaHeader.SchemaBodyWidth != nil || len(schemaHeader.SchemaFields) != 0 || !reflect.DeepEqual(schemaHeader.SchemaRawFallback, []byte{171, 205}) {
		schemaT.Fatalf("ab05 header ceased to be partial: %#v %v", schemaHeader, schemaError)
	}
	schemaOpaque, schemaFound, schemaError := SchemaQuery(schemaRecords, schemaRecords[3].SchemaKey)
	if schemaError != nil || !schemaFound || schemaOpaque.SchemaFields[0].SchemaPosition != nil || schemaOpaque.SchemaFields[0].SchemaWidth != nil {
		schemaT.Fatalf("ab05 unknown extent became explicit: %#v %v", schemaOpaque, schemaError)
	}
	schemaRecords[2].SchemaBodyWidth = schemaInt(0)
	schemaZero, _, schemaError := SchemaQuery(schemaRecords, schemaRecords[2].SchemaKey)
	if schemaError != nil || schemaZero.SchemaBodyWidth == nil || *schemaZero.SchemaBodyWidth != 0 {
		schemaT.Fatal("ab05 explicit empty body lost")
	}
}

func TestAb05UnknownGradeAndCodec(schemaT *testing.T) {
	schemaRecords := schemaFixture(schemaT)
	schemaHeader, _, schemaError := SchemaQuery(schemaRecords, schemaRecords[2].SchemaKey)
	if schemaError != nil || schemaHeader.SchemaGrade != SchemaUnknown {
		schemaT.Fatal("ab05 missing grade promoted")
	}
	schemaOpaque, _, schemaError := SchemaQuery(schemaRecords, schemaRecords[3].SchemaKey)
	if schemaError != nil || schemaOpaque.SchemaGrade != SchemaUnknown || schemaOpaque.SchemaFields[0].SchemaGrade != SchemaUnknown || schemaOpaque.SchemaFields[0].SchemaCodec != "ab05_unrecognized" || !reflect.DeepEqual(schemaOpaque.SchemaFields[0].SchemaRawFallback, []byte{171, 205}) || !reflect.DeepEqual(schemaOpaque.SchemaRawFallback, []byte{171, 205}) {
		schemaT.Fatalf("ab05 opaque evidence changed: %#v %v", schemaOpaque, schemaError)
	}
	if schemaRecords[2].SchemaGrade != "" || schemaRecords[3].SchemaFields[0].SchemaGrade != "" {
		schemaT.Fatal("ab05 supplied evidence mutated")
	}
	for _, schemaGrade := range []SchemaGrade{SchemaUnknown, SchemaHypothesis, SchemaProven} {
		schemaRecords[3].SchemaGrade = schemaGrade
		schemaRecords[3].SchemaFields[0].SchemaGrade = schemaGrade
		schemaGot, _, schemaError := SchemaQuery(schemaRecords, schemaRecords[3].SchemaKey)
		if schemaError != nil || schemaGot.SchemaGrade != schemaGrade || schemaGot.SchemaFields[0].SchemaGrade != schemaGrade || schemaGot.SchemaFields[0].SchemaCodec != "ab05_unrecognized" {
			schemaT.Fatal("ab05 grade conflated with codec")
		}
	}
}

func TestAb05ContradictoryExtents(schemaT *testing.T) {
	schemaCases := []struct {
		schemaName   string
		schemaChange func(*SchemaRecord)
	}{
		{"ab05_negative_body", func(schemaRecord *SchemaRecord) { schemaRecord.SchemaBodyWidth = schemaInt(-1) }},
		{"ab05_negative_position", func(schemaRecord *SchemaRecord) { schemaRecord.SchemaFields[0].SchemaPosition = schemaInt(-1) }},
		{"ab05_zero_field_width", func(schemaRecord *SchemaRecord) { schemaRecord.SchemaFields[0].SchemaWidth = schemaInt(0) }},
		{"ab05_negative_field_width", func(schemaRecord *SchemaRecord) { schemaRecord.SchemaFields[0].SchemaWidth = schemaInt(-1) }},
		{"ab05_end_exceeds_body", func(schemaRecord *SchemaRecord) { schemaRecord.SchemaFields[0].SchemaPosition = schemaInt(1) }},
		{"ab05_position_exceeds_body", func(schemaRecord *SchemaRecord) {
			schemaRecord.SchemaFields[0].SchemaPosition = schemaInt(3)
			schemaRecord.SchemaFields[0].SchemaWidth = nil
		}},
		{"ab05_position_at_body_end", func(schemaRecord *SchemaRecord) {
			schemaRecord.SchemaFields[0].SchemaPosition = schemaInt(2)
			schemaRecord.SchemaFields[0].SchemaWidth = nil
		}},
		{"ab05_width_without_position", func(schemaRecord *SchemaRecord) {
			schemaRecord.SchemaFields[0].SchemaPosition = nil
			schemaRecord.SchemaFields[0].SchemaWidth = schemaInt(3)
		}},
		{"ab05_overflow_extent", func(schemaRecord *SchemaRecord) {
			schemaRecord.SchemaFields[0].SchemaPosition = schemaInt(1)
			schemaRecord.SchemaFields[0].SchemaWidth = schemaInt(int(^uint(0) >> 1))
		}},
		{"ab05_invalid_record_grade", func(schemaRecord *SchemaRecord) { schemaRecord.SchemaGrade = "ab05_invalid" }},
		{"ab05_invalid_field_grade", func(schemaRecord *SchemaRecord) { schemaRecord.SchemaFields[0].SchemaGrade = "ab05_invalid" }},
	}
	for _, schemaCase := range schemaCases {
		schemaT.Run(schemaCase.schemaName, func(schemaT *testing.T) {
			schemaRecord := schemaFixture(schemaT)[0]
			if SchemaValidate([]SchemaRecord{schemaRecord}) != nil {
				schemaT.Fatal("ab05 invalid control")
			}
			schemaCase.schemaChange(&schemaRecord)
			schemaGot, schemaFound, schemaError := SchemaQuery([]SchemaRecord{schemaRecord}, schemaRecord.SchemaKey)
			if schemaError == nil || schemaFound || !reflect.DeepEqual(schemaGot, SchemaRecord{}) {
				schemaT.Fatal("ab05 invalid evidence admitted")
			}
		})
	}
}

func TestAb05QueryIsolation(schemaT *testing.T) {
	schemaRecords := schemaFixture(schemaT)
	schemaBefore := schemaFixture(schemaT)
	schemaGot, _, schemaError := SchemaQuery(schemaRecords, schemaRecords[0].SchemaKey)
	if schemaError != nil {
		schemaT.Fatal(schemaError)
	}
	*schemaGot.SchemaBodyWidth = 99
	*schemaGot.SchemaFields[0].SchemaPosition = 99
	*schemaGot.SchemaFields[0].SchemaWidth = 99
	schemaGot.SchemaRawFallback[0] = 0
	schemaGot.SchemaFields[0].SchemaRawFallback[0] = 0
	schemaGot.SchemaFields[0].SchemaSentinel[0] = 0
	schemaGot.SchemaFields[0].SchemaEnum[0] = "ab05_changed"
	schemaGot.SchemaFields[1].SchemaName = "ab05_changed"
	if !reflect.DeepEqual(schemaRecords, schemaBefore) {
		schemaT.Fatal("ab05 query changed prior field evidence")
	}
	schemaMissing := schemaRecords[0].SchemaKey
	schemaMissing.SchemaFamily = "ab05_absent"
	schemaGot, schemaFound, schemaError := SchemaQuery(schemaRecords, schemaMissing)
	if schemaError != nil || schemaFound || !reflect.DeepEqual(schemaGot, SchemaRecord{}) {
		schemaT.Fatal("ab05 missing key admitted")
	}
}

func TestAb05JSONNumericKeyPresence(schemaT *testing.T) {
	schemaBase := map[string]any{
		"ab05_family": "ab05_synthetic", "ab05_variant": "ab05_primary",
		"ab05_direction": "request", "ab05_alternative": "ab05_candidate",
		"ab05_opcode": 0, "ab05_group": 0, "ab05_instance": 0, "ab05_register": 0,
	}
	schemaDecode := func(schemaT *testing.T, schemaKey map[string]any) (SchemaRecord, error) {
		schemaT.Helper()
		schemaData, schemaError := json.Marshal(map[string]any{"ab05_key": schemaKey})
		if schemaError != nil {
			schemaT.Fatal(schemaError)
		}
		var schemaRecord SchemaRecord
		schemaError = json.Unmarshal(schemaData, &schemaRecord)
		return schemaRecord, schemaError
	}
	schemaZero, schemaError := schemaDecode(schemaT, schemaBase)
	if schemaError != nil {
		schemaT.Fatal(schemaError)
	}
	schemaWant := SchemaKey{SchemaFamily: "ab05_synthetic", SchemaVariant: "ab05_primary", SchemaDirection: "request", SchemaAlternative: "ab05_candidate"}
	schemaGot, schemaFound, schemaError := SchemaQuery([]SchemaRecord{schemaZero}, schemaWant)
	if schemaError != nil || !schemaFound || schemaGot.SchemaKey != schemaWant {
		schemaT.Fatalf("ab05 explicit zero selectors rejected: %#v %v %v", schemaGot, schemaFound, schemaError)
	}
	for _, schemaName := range []string{"ab05_opcode", "ab05_group", "ab05_instance", "ab05_register"} {
		for _, schemaMode := range []string{"absent", "null", "misspelled"} {
			schemaT.Run(schemaName+"/"+schemaMode, func(schemaT *testing.T) {
				schemaKey := make(map[string]any, len(schemaBase))
				for schemaName, schemaValue := range schemaBase {
					schemaKey[schemaName] = schemaValue
				}
				delete(schemaKey, schemaName)
				switch schemaMode {
				case "null":
					schemaKey[schemaName] = nil
				case "misspelled":
					schemaKey[schemaName+"_typo"] = 0
				}
				if _, schemaError := schemaDecode(schemaT, schemaKey); schemaError == nil {
					schemaT.Fatal("ab05 missing numeric selector accepted")
				}
			})
		}
	}
	schemaBase["ab05_unknown"] = 0
	if _, schemaError := schemaDecode(schemaT, schemaBase); schemaError == nil {
		schemaT.Fatal("ab05 unknown key property accepted")
	}
}

func TestAb05InvalidStringKeys(schemaT *testing.T) {
	schemaCases := []struct {
		schemaName   string
		schemaChange func(*SchemaKey)
	}{
		{"ab05_family", func(schemaKey *SchemaKey) { schemaKey.SchemaFamily = "" }},
		{"ab05_variant", func(schemaKey *SchemaKey) { schemaKey.SchemaVariant = "" }},
		{"ab05_direction", func(schemaKey *SchemaKey) { schemaKey.SchemaDirection = "" }},
		{"ab05_alternative", func(schemaKey *SchemaKey) { schemaKey.SchemaAlternative = "" }},
		{"ab05_invalid_direction", func(schemaKey *SchemaKey) { schemaKey.SchemaDirection = "ab05_invalid" }},
	}
	for _, schemaCase := range schemaCases {
		schemaT.Run(schemaCase.schemaName, func(schemaT *testing.T) {
			schemaRecord := schemaFixture(schemaT)[0]
			schemaCase.schemaChange(&schemaRecord.SchemaKey)
			if _, schemaFound, schemaError := SchemaQuery([]SchemaRecord{schemaRecord}, schemaRecord.SchemaKey); schemaError == nil || schemaFound {
				schemaT.Fatal("ab05 invalid string dimension accepted")
			}
		})
	}
}

func TestAb05QueryEmptyBytes(schemaT *testing.T) {
	for _, schemaBytes := range [][]byte{nil, {}} {
		schemaRecord := schemaFixture(schemaT)[0]
		schemaRecord.SchemaRawFallback = schemaBytes
		schemaRecord.SchemaFields[0].SchemaRawFallback = schemaBytes
		schemaRecord.SchemaFields[0].SchemaSentinel = schemaBytes
		schemaGot, schemaFound, schemaError := SchemaQuery([]SchemaRecord{schemaRecord}, schemaRecord.SchemaKey)
		if schemaError != nil || !schemaFound || !reflect.DeepEqual(schemaGot, schemaRecord) {
			schemaT.Fatalf("ab05 empty byte presence lost: %#v %v", schemaGot, schemaError)
		}
	}
}
