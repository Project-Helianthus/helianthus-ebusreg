package vaillant

import (
	ab90Binary "encoding/binary"
	ab90Fmt "fmt"
	ab90Reflect "reflect"
	ab90Testing "testing"
)

func ab90Ptr[Ab90T any](ab90Value Ab90T) *Ab90T { return &ab90Value }

func ab90ExpectedField(ab90Operation uint8, ab90Direction, ab90Name string, ab90Position, ab90Width int, ab90Maximum int64) Ab90MirrorField {
	ab90Field := Ab90MirrorField{
		Ab90Key:             Ab90MirrorKey{Ab90Family: "vaillant", Ab90Variant: "mirror-diagnosis", Ab90Primary: 181, Ab90Secondary: ab90Ptr(uint8(26)), Ab90Operation: ab90Ptr(ab90Operation), Ab90Direction: ab90Direction, Ab90Field: ab90Name},
		Ab90LogicalPosition: ab90Ptr(ab90Position), Ab90WidthBytes: ab90Ptr(ab90Width), Ab90Signedness: "unsigned", Ab90ByteOrder: "not_applicable", Ab90Conversion: "identity",
		Ab90FormalMin: ab90Ptr(int64(0)), Ab90FormalMax: ab90Ptr(ab90Maximum), Ab90Enumeration: map[int64]string{}, Ab90Bits: []Ab90MirrorBit{}, Ab90Alternatives: []Ab90MirrorAlternative{}, Ab90Reserved: []int{},
		Ab90Grade: Ab05Hypothesis, Ab90RawFallback: "preserve original bytes including unnamed and reserved bits", Ab90Correlation: "Unknown", Ab90Unmapped: "Unknown",
	}
	if ab90Width == 2 {
		ab90Field.Ab90ByteOrder = "little_endian"
	}
	return ab90Field
}

func ab90ExpectedBit(ab90Index int, ab90Name string, ab90Role string) Ab90MirrorBit {
	ab90Bit := Ab90MirrorBit{Ab90Index: ab90Index, Ab90Role: ab90Role, Ab90Default: ab90Ptr(false), Ab90Grade: Ab05Hypothesis}
	if ab90Name != "" {
		ab90Bit.Ab90Name = ab90Ptr(ab90Name)
	}
	if ab90Role == "unnamed" {
		ab90Bit.Ab90Default = nil
		ab90Bit.Ab90Grade = Ab05Unknown
	}
	return ab90Bit
}

func ab90ExpectedFields() []Ab90MirrorField {
	ab90Fields := []Ab90MirrorField{}
	for _, ab90Operation := range []uint8{0, 1, 2, 3, 6} {
		for _, ab90Direction := range []string{"request", "positive"} {
			ab90Field := ab90ExpectedField(ab90Operation, ab90Direction, "sequence", 3, 1, 255)
			ab90Field.Ab90Default = ab90Ptr(int64(0))
			if ab90Operation == 0 || ab90Operation == 1 {
				ab90Field.Ab90Correlation = "candidate_equal_sequence"
			}
			ab90Fields = append(ab90Fields, ab90Field)
		}
	}
	ab90Field := ab90ExpectedField(0, "positive", "acknowledgement", 4, 1, 255)
	ab90Field.Ab90Default = ab90Ptr(int64(0))
	ab90Field.Ab90Bits = append(ab90Field.Ab90Bits, ab90ExpectedBit(0, "denied", "named"))
	for ab90Index := 1; ab90Index < 8; ab90Index++ {
		ab90Field.Ab90Bits = append(ab90Field.Ab90Bits, ab90ExpectedBit(ab90Index, "", "unnamed"))
	}
	ab90Fields = append(ab90Fields, ab90Field)
	ab90Field = ab90ExpectedField(2, "request", "flags", 4, 1, 255)
	ab90Bit := ab90ExpectedBit(0, "triple-storage", "named")
	ab90Bit.Ab90Default = nil
	ab90Field.Ab90Bits = append(ab90Field.Ab90Bits, ab90Bit)
	for ab90Index := 1; ab90Index < 8; ab90Index++ {
		ab90Field.Ab90Bits = append(ab90Field.Ab90Bits, ab90ExpectedBit(ab90Index, "", "unnamed"))
	}
	ab90Fields = append(ab90Fields, ab90Field, ab90ExpectedField(2, "request", "length", 5, 1, 255))
	ab90Field = ab90ExpectedField(2, "request", "mirror_address", 6, 2, 65535)
	ab90Field.Ab90LogicalMin = ab90Ptr(int64(0))
	ab90Field.Ab90LogicalLimit = ab90Ptr("mirror-length - 1")
	ab90Fields = append(ab90Fields, ab90Field)
	ab90Field = ab90ExpectedField(2, "positive", "positive", 4, 1, 255)
	ab90Field.Ab90Enumeration = map[int64]string{0: "false", 1: "true"}
	ab90Field.Ab90Unmapped = "2..255 unresolved"
	ab90Field.Ab90Alternatives = []Ab90MirrorAlternative{{Ab90Name: "competing out-of-area bit", Ab90WidthBytes: ab90Ptr(1), Ab90Condition: "unresolved alternative to whole-u8 boolean", Ab90Enumeration: map[int64]string{}, Ab90Bits: []Ab90MirrorBit{ab90ExpectedBit(0, "out-of-area", "named")}, Ab90Grade: Ab05Hypothesis}}
	ab90Field.Ab90Alternatives[0].Ab90Bits[0].Ab90Default = nil
	ab90Fields = append(ab90Fields, ab90Field)
	for _, ab90Operation := range []uint8{2, 3} {
		ab90Direction, ab90Position, ab90Limit := "positive", 5, "end-delimited; maximum unset; returned length independent of requested length"
		if ab90Operation == 3 {
			ab90Direction = "request"
			ab90Position = 9
			ab90Limit = "end-delimited; maximum unset; empty data does not prove acceptance"
		}
		ab90Field = ab90ExpectedField(ab90Operation, ab90Direction, "data", ab90Position, 1, 255)
		ab90Field.Ab90WidthBytes = nil
		ab90Field.Ab90FormalMin = nil
		ab90Field.Ab90FormalMax = nil
		ab90Field.Ab90Signedness = "Unknown"
		ab90Field.Ab90Conversion = "opaque"
		ab90Field.Ab90LogicalMin = ab90Ptr(int64(0))
		ab90Field.Ab90LogicalLimit = ab90Ptr(ab90Limit)
		ab90Fields = append(ab90Fields, ab90Field)
	}
	ab90Field = ab90ExpectedField(3, "request", "flags", 4, 1, 255)
	ab90Field.Ab90Default = ab90Ptr(int64(0))
	ab90Field.Ab90Bits = []Ab90MirrorBit{ab90ExpectedBit(0, "triple-storage", "named"), ab90ExpectedBit(1, "direct-update", "named")}
	ab90Field.Ab90Reserved = []int{2, 3, 4, 5, 6, 7}
	for ab90Index := 2; ab90Index < 8; ab90Index++ {
		ab90Field.Ab90Bits = append(ab90Field.Ab90Bits, ab90ExpectedBit(ab90Index, "", "reserved"))
	}
	ab90Fields = append(ab90Fields, ab90Field)
	for ab90Index, ab90Name := range []string{"mirror_address", "symbol_address"} {
		ab90Field = ab90ExpectedField(3, "request", ab90Name, 5+2*ab90Index, 2, 65535)
		ab90Field.Ab90Default = ab90Ptr(int64(0))
		ab90Field.Ab90LogicalLimit = ab90Ptr("independent " + ab90Name + " limit Unknown")
		ab90Fields = append(ab90Fields, ab90Field)
	}
	ab90Field = ab90ExpectedField(3, "positive", "flags", 4, 1, 255)
	ab90Field.Ab90Default = ab90Ptr(int64(0))
	ab90Field.Ab90Reserved = []int{4, 5, 6, 7}
	for ab90Index, ab90Name := range []string{"not-written", "invalid-symbol", "invalid-parameter", "storage-error"} {
		ab90Field.Ab90Bits = append(ab90Field.Ab90Bits, ab90ExpectedBit(ab90Index, ab90Name, "named"))
	}
	for ab90Index := 4; ab90Index < 8; ab90Index++ {
		ab90Field.Ab90Bits = append(ab90Field.Ab90Bits, ab90ExpectedBit(ab90Index, "", "reserved"))
	}
	ab90Fields = append(ab90Fields, ab90Field, ab90ExpectedField(6, "request", "level", 4, 1, 255))
	ab90Field = ab90ExpectedField(6, "request", "point", 5, 1, 255)
	ab90Field.Ab90LogicalMin = ab90Ptr(int64(0))
	ab90Field.Ab90LogicalMax = ab90Ptr(int64(99))
	ab90Fields = append(ab90Fields, ab90Field, ab90ExpectedField(6, "request", "value", 6, 2, 65535))
	ab90Field = ab90ExpectedField(6, "positive", "save_status", 4, 1, 255)
	ab90Field.Ab90Default = ab90Ptr(int64(0))
	ab90Field.Ab90Enumeration = map[int64]string{0: "value-not-saved", 1: "point-not-found", 2: "write-protected", 3: "range-violation", 4: "misconfiguration"}
	ab90Field.Ab90Unmapped = "5..255 unlabeled; no implicit success"
	ab90Fields = append(ab90Fields, ab90Field)
	for ab90Index, ab90Name := range []string{"mirror_address", "symbol_address"} {
		ab90Field = ab90ExpectedField(6, "positive", ab90Name, 5+2*ab90Index, 2, 65535)
		ab90Field.Ab90LogicalLimit = ab90Ptr("independent " + ab90Name + " limit Unknown")
		ab90Fields = append(ab90Fields, ab90Field)
	}
	ab90Field = ab90ExpectedField(6, "positive", "value", 9, 2, 65535)
	ab90Field.Ab90WidthBytes = nil
	for ab90Index, ab90Name := range []string{"formal-u16", "display-one-byte", "display-two-byte"} {
		ab90Field.Ab90Alternatives = append(ab90Field.Ab90Alternatives, Ab90MirrorAlternative{Ab90Name: ab90Name, Ab90WidthBytes: ab90Ptr([]int{2, 1, 2}[ab90Index]), Ab90Condition: []string{"formal description", "display type selects one byte", "display type selects two bytes"}[ab90Index], Ab90Enumeration: map[int64]string{}, Ab90Bits: []Ab90MirrorBit{}, Ab90Grade: Ab05Hypothesis})
	}
	return append(ab90Fields, ab90Field)
}

func ab90Catalog(ab90T *ab90Testing.T) Ab90MirrorMetadata {
	ab90T.Helper()
	ab90Metadata, ab90Error := Ab90B51AMirrorFields()
	if ab90Error != nil {
		ab90T.Fatal(ab90Error)
	}
	return ab90Metadata
}

func ab90Get(ab90T *ab90Testing.T, ab90Operation uint8, ab90Direction, ab90Name string) Ab90MirrorField {
	ab90T.Helper()
	ab90Key := Ab90MirrorKey{Ab90Family: "vaillant", Ab90Variant: "mirror-diagnosis", Ab90Primary: 181, Ab90Secondary: ab90Ptr(uint8(26)), Ab90Operation: ab90Ptr(ab90Operation), Ab90Direction: ab90Direction, Ab90Field: ab90Name}
	ab90Field, ab90Found, ab90Error := Ab90LookupB51AMirrorField(ab90Key)
	if ab90Error != nil || !ab90Found {
		ab90T.Fatalf("ab90 missing %v: %v", ab90Key, ab90Error)
	}
	return ab90Field
}

func TestAb90ExactProductionMetadata(ab90T *ab90Testing.T) {
	ab90Got := ab90Catalog(ab90T)
	ab90Want := Ab90MirrorMetadata{Ab90RepresentationGrade: Ab05Proven, Ab90SemanticsGrade: Ab05Hypothesis, Ab90PhysicalAdaptation: Ab05Unknown, Ab90StorageDurability: Ab05Unknown, Ab90SuccessSemantics: Ab05Unknown, Ab90Complete: false, Ab90NativeContract: "protocols/vaillant/ebus-vaillant-B51A.md#candidate-mirror-and-diagnosis-layouts", Ab90Limitations: []string{
		"Finite supplemental mirror/diagnosis declarations; configured-register tables remain separate.",
		"Logical positions include service prefix; physical response adaptation is Unknown.",
		"Sequence equality for reset/01 is candidate correlation only; no retry or freshness policy.",
		"No invocation, registration, transport, persistence or execution authority.",
		"Unprovided properties and broader family coverage remain Unknown.",
	}, Ab90Fields: ab90ExpectedFields()}
	if !ab90Reflect.DeepEqual(ab90Got, ab90Want) {
		ab90T.Error("ab90 production metadata differs from independent contract")
	}
	if len(ab90Got.Ab90Fields) != 28 {
		ab90T.Fatalf("ab90 field count %d", len(ab90Got.Ab90Fields))
	}
	for ab90Index, ab90Field := range ab90Got.Ab90Fields {
		ab90T.Run(ab90Fmt.Sprintf("ab90_field_%d", ab90Index), func(ab90T *ab90Testing.T) {
			if !ab90Reflect.DeepEqual(ab90Field, ab90Want.Ab90Fields[ab90Index]) {
				ab90T.Fatalf("ab90 field %s/%s op %d mismatch", ab90Field.Ab90Key.Ab90Direction, ab90Field.Ab90Key.Ab90Field, *ab90Field.Ab90Key.Ab90Operation)
			}
			ab90Found, ab90OK, ab90Error := Ab90LookupB51AMirrorField(ab90Field.Ab90Key)
			if ab90Error != nil || !ab90OK || !ab90Reflect.DeepEqual(ab90Field, ab90Found) {
				ab90T.Fatal("ab90 lookup mismatch")
			}
		})
	}
}

func ab90Read(ab90T *ab90Testing.T, ab90Field Ab90MirrorField, ab90Raw []byte, ab90Alternative string) uint64 {
	ab90T.Helper()
	ab90Width := ab90Field.Ab90WidthBytes
	if ab90Alternative != "" {
		for _, ab90Choice := range ab90Field.Ab90Alternatives {
			if ab90Choice.Ab90Name == ab90Alternative {
				ab90Width = ab90Choice.Ab90WidthBytes
			}
		}
	}
	if ab90Width == nil || ab90Field.Ab90LogicalPosition == nil {
		ab90T.Fatal("ab90 no unconditional scalar layout")
	}
	ab90Position := *ab90Field.Ab90LogicalPosition
	if ab90Position+*ab90Width > len(ab90Raw) {
		ab90T.Fatal("ab90 short synthetic carrier")
	}
	if *ab90Width == 1 {
		return uint64(ab90Raw[ab90Position])
	}
	if *ab90Width == 2 && ab90Field.Ab90ByteOrder == "little_endian" {
		return uint64(ab90Binary.LittleEndian.Uint16(ab90Raw[ab90Position : ab90Position+2]))
	}
	ab90T.Fatal("ab90 unsupported synthetic layout")
	return 0
}

func TestAb90SyntheticSequenceAndReset(ab90T *ab90Testing.T) {
	for _, ab90Operation := range []uint8{0, 1, 2, 3, 6} {
		for _, ab90Direction := range []string{"request", "positive"} {
			ab90Sequence := ab90Get(ab90T, ab90Operation, ab90Direction, "sequence")
			for _, ab90Value := range []byte{0, 255} {
				if ab90Read(ab90T, ab90Sequence, []byte{0, 0, 0, ab90Value}, "") != uint64(ab90Value) {
					ab90T.Fatal("ab90 sequence lost")
				}
			}
			ab90Expected := "Unknown"
			if ab90Operation == 0 || ab90Operation == 1 {
				ab90Expected = "candidate_equal_sequence"
			}
			if ab90Sequence.Ab90Correlation != ab90Expected {
				ab90T.Fatal("ab90 correlation leaked across operations")
			}
		}
		ab90Request := ab90Get(ab90T, ab90Operation, "request", "sequence")
		ab90Positive := ab90Get(ab90T, ab90Operation, "positive", "sequence")
		ab90Matched := ab90Request.Ab90Correlation == "candidate_equal_sequence" && ab90Read(ab90T, ab90Request, []byte{0, 0, 0, 0}, "") == ab90Read(ab90T, ab90Positive, []byte{0, 0, 0, 255}, "")
		if ab90Matched {
			ab90T.Fatal("ab90 mismatched sequence correlated")
		}
	}
	ab90Ack := ab90Get(ab90T, 0, "positive", "acknowledgement")
	ab90Raw := ab90Read(ab90T, ab90Ack, []byte{0, 0, 0, 0, 2}, "")
	if ab90Raw&(1<<ab90Ack.Ab90Bits[0].Ab90Index) != 0 || ab90Raw&(1<<ab90Ack.Ab90Bits[1].Ab90Index) == 0 || ab90Ack.Ab90Bits[1].Ab90Grade != Ab05Unknown || ab90Ack.Ab90Bits[1].Ab90Name != nil {
		ab90T.Fatal("ab90 reset02 is unnamed, not denied")
	}
}

func TestAb90SyntheticReadAndWrite(ab90T *ab90Testing.T) {
	ab90Length := ab90Get(ab90T, 2, "request", "length")
	ab90Data := ab90Get(ab90T, 2, "positive", "data")
	ab90Requested := ab90Read(ab90T, ab90Length, []byte{0, 0, 0, 0, 0, 3}, "")
	ab90Returned := []byte{0, 0, 0, 0, 1, 0x12, 0x34}
	ab90Actual := len(ab90Returned[*ab90Data.Ab90LogicalPosition:])
	if ab90Requested != 3 || ab90Actual != 2 || ab90Data.Ab90WidthBytes != nil || ab90Data.Ab90FormalMax != nil || *ab90Data.Ab90LogicalMin != 0 {
		ab90T.Fatal("ab90 requested/actual lengths conflated")
	}
	ab90Positive := ab90Get(ab90T, 2, "positive", "positive")
	for ab90Raw, ab90Want := range map[byte]string{0: "false", 1: "true", 2: ""} {
		ab90Value := ab90Read(ab90T, ab90Positive, []byte{0, 0, 0, 0, ab90Raw}, "")
		if ab90Positive.Ab90Enumeration[int64(ab90Value)] != ab90Want {
			ab90T.Fatal("ab90 whole-u8 alternatives collapsed")
		}
	}
	ab90Mirror := ab90Get(ab90T, 3, "request", "mirror_address")
	ab90Symbol := ab90Get(ab90T, 3, "request", "symbol_address")
	ab90Write := []byte{0, 0, 0, 0, 0, 0x34, 0x12, 0x78, 0x56}
	if ab90Read(ab90T, ab90Mirror, ab90Write, "") != 0x1234 || ab90Read(ab90T, ab90Symbol, ab90Write, "") != 0x5678 {
		ab90T.Fatal("ab90 mirror/symbol swapped")
	}
	ab90Opaque := ab90Get(ab90T, 3, "request", "data")
	if len(ab90Write[*ab90Opaque.Ab90LogicalPosition:]) != 0 || ab90Opaque.Ab90Conversion != "opaque" || ab90Opaque.Ab90WidthBytes != nil || *ab90Opaque.Ab90LogicalMin != 0 || ab90Catalog(ab90T).Ab90SuccessSemantics != Ab05Unknown {
		ab90T.Fatal("ab90 empty write inferred acceptance")
	}
}

func TestAb90SyntheticSave(ab90T *ab90Testing.T) {
	ab90Point := ab90Get(ab90T, 6, "request", "point")
	for ab90Raw, ab90Want := range map[byte]bool{99: true, 100: false} {
		ab90Value := int64(ab90Read(ab90T, ab90Point, []byte{0, 0, 0, 7, 9, ab90Raw}, ""))
		ab90Fits := ab90Value >= *ab90Point.Ab90LogicalMin && ab90Value <= *ab90Point.Ab90LogicalMax
		if ab90Fits != ab90Want {
			ab90T.Fatal("ab90 point domain changed")
		}
	}
	ab90Status := ab90Get(ab90T, 6, "positive", "save_status")
	for ab90Raw, ab90Want := range map[byte]string{0: "value-not-saved", 1: "point-not-found", 5: ""} {
		if ab90Status.Ab90Enumeration[int64(ab90Read(ab90T, ab90Status, []byte{0, 0, 0, 0, ab90Raw}, ""))] != ab90Want {
			ab90T.Fatal("ab90 save status invented success")
		}
	}
	ab90Value := ab90Get(ab90T, 6, "positive", "value")
	ab90Carrier := []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0x34, 0x12}
	if ab90Value.Ab90WidthBytes != nil || ab90Read(ab90T, ab90Value, ab90Carrier, "display-one-byte") != 0x34 || ab90Read(ab90T, ab90Value, ab90Carrier, "display-two-byte") != 0x1234 || ab90Read(ab90T, ab90Value, ab90Carrier, "formal-u16") != 0x1234 {
		ab90T.Fatal("ab90 conditional widths collapsed")
	}
	ab90Request := ab90Get(ab90T, 6, "request", "value")
	if ab90Read(ab90T, ab90Request, []byte{0, 0, 0, 0, 0, 0, 0x34, 0x12}, "") != 0x1234 {
		ab90T.Fatal("ab90 save request narrowed")
	}
}

func TestAb90IdentityAxes(ab90T *ab90Testing.T) {
	ab90Base := ab90Get(ab90T, 2, "request", "length").Ab90Key
	for ab90Index := 0; ab90Index < ab90Reflect.TypeOf(ab90Base).NumField(); ab90Index++ {
		ab90Axis := ab90Reflect.TypeOf(ab90Base).Field(ab90Index).Name
		ab90T.Run("ab90_"+ab90Axis, func(ab90T *ab90Testing.T) {
			ab90Changed := ab90Base
			ab90Value := ab90Reflect.ValueOf(&ab90Changed).Elem().Field(ab90Index)
			switch ab90Value.Kind() {
			case ab90Reflect.String:
				ab90Value.SetString("ab90_distinct")
			case ab90Reflect.Uint8:
				ab90Value.SetUint(0)
			case ab90Reflect.Pointer:
				ab90New := ab90Reflect.New(ab90Value.Type().Elem())
				ab90New.Elem().SetUint(99)
				ab90Value.Set(ab90New)
			default:
				ab90T.Fatal("ab90 unchecked axis")
			}
			ab90Field, ab90Found, ab90Error := Ab90LookupB51AMirrorField(ab90Changed)
			if ab90Error != nil || ab90Found || !ab90Reflect.DeepEqual(ab90Field, Ab90MirrorField{}) {
				ab90T.Fatal("ab90 changed identity matched")
			}
			if ab90MirrorKeysEqual(ab90Base, ab90Changed) {
				ab90T.Fatal("ab90 ignored key dimension")
			}
			ab90Value.Set(ab90Reflect.ValueOf(ab90Base).Field(ab90Index))
			if _, ab90Found, ab90Error = Ab90LookupB51AMirrorField(ab90Changed); ab90Error != nil || !ab90Found {
				ab90T.Fatal("ab90 invalid restored control")
			}
		})
	}
	for _, ab90Field := range ab90Catalog(ab90T).Ab90Fields {
		for _, ab90Axis := range []string{"Ab90Secondary", "Ab90Operation", "Ab90Opcode", "Ab90Group", "Ab90Instance", "Ab90Register"} {
			ab90Changed := ab90Field.Ab90Key
			ab90Value := ab90Reflect.ValueOf(&ab90Changed).Elem().FieldByName(ab90Axis)
			if ab90Value.IsNil() {
				ab90Value.Set(ab90Reflect.New(ab90Value.Type().Elem()))
			} else {
				ab90Value.SetZero()
			}
			if _, ab90Found, ab90Error := Ab90LookupB51AMirrorField(ab90Changed); ab90Error != nil || ab90Found {
				ab90T.Fatalf("ab90 absent versus supplied %s wildcarded", ab90Axis)
			}
		}
	}
}

func TestAb90B524IdentityAxes(ab90T *ab90Testing.T) {
	ab90Base := Ab90MirrorKey{Ab90Family: "ab90_synthetic_family", Ab90Variant: "ab90_synthetic_variant", Ab90Primary: 181, Ab90Secondary: ab90Ptr(uint8(36)), Ab90Operation: ab90Ptr(uint8(0)), Ab90Direction: "request", Ab90Field: "ab90_synthetic_field", Ab90Opcode: ab90Ptr(uint8(0)), Ab90Group: ab90Ptr(uint8(0)), Ab90Instance: ab90Ptr(uint8(0)), Ab90Register: ab90Ptr(uint16(0))}
	for ab90Index := 0; ab90Index < ab90Reflect.TypeOf(ab90Base).NumField(); ab90Index++ {
		ab90Axis := ab90Reflect.TypeOf(ab90Base).Field(ab90Index).Name
		ab90T.Run("ab90_"+ab90Axis, func(ab90T *ab90Testing.T) {
			ab90Changed := ab90Base
			ab90Value := ab90Reflect.ValueOf(&ab90Changed).Elem().Field(ab90Index)
			switch ab90Value.Kind() {
			case ab90Reflect.String:
				ab90Value.SetString("ab90_other")
			case ab90Reflect.Uint8:
				ab90Value.SetUint(0)
			case ab90Reflect.Pointer:
				ab90New := ab90Reflect.New(ab90Value.Type().Elem())
				ab90New.Elem().SetUint(ab90Value.Elem().Uint() + 1)
				ab90Value.Set(ab90New)
			}
			if ab90MirrorKeysEqual(ab90Base, ab90Changed) {
				ab90T.Fatal("ab90 B524 axis ignored")
			}
			if ab90Value.Kind() == ab90Reflect.Pointer {
				ab90Value.SetZero()
				if ab90MirrorKeysEqual(ab90Base, ab90Changed) {
					ab90T.Fatal("ab90 B524 absence wildcarded")
				}
				ab90Value.Set(ab90Reflect.New(ab90Value.Type().Elem()))
				ab90Zero := ab90Changed
				ab90Value.SetZero()
				if ab90MirrorKeysEqual(ab90Zero, ab90Changed) {
					ab90T.Fatal("ab90 B524 absent equals zero")
				}
			}
			ab90Value.Set(ab90Reflect.ValueOf(ab90Base).Field(ab90Index))
			if !ab90MirrorKeysEqual(ab90Base, ab90Changed) {
				ab90T.Fatal("ab90 B524 restored control differs")
			}
		})
	}
	if _, ab90Found, ab90Error := Ab90LookupB51AMirrorField(ab90Base); ab90Error != nil || ab90Found {
		ab90T.Fatal("ab90 synthetic B524 acquired production declaration")
	}
}

func ab90Mutate(ab90Value ab90Reflect.Value) {
	switch ab90Value.Kind() {
	case ab90Reflect.Struct:
		for ab90Index := 0; ab90Index < ab90Value.NumField(); ab90Index++ {
			ab90Mutate(ab90Value.Field(ab90Index))
		}
	case ab90Reflect.Pointer:
		if !ab90Value.IsNil() {
			ab90Mutate(ab90Value.Elem())
		}
	case ab90Reflect.Slice:
		for ab90Index := 0; ab90Index < ab90Value.Len(); ab90Index++ {
			ab90Mutate(ab90Value.Index(ab90Index))
		}
	case ab90Reflect.Map:
		for _, ab90Key := range ab90Value.MapKeys() {
			ab90Value.SetMapIndex(ab90Key, ab90Reflect.ValueOf("ab90_mutated"))
		}
	case ab90Reflect.String:
		ab90Value.SetString("ab90_mutated")
	case ab90Reflect.Bool:
		ab90Value.SetBool(!ab90Value.Bool())
	case ab90Reflect.Int, ab90Reflect.Int64:
		ab90Value.SetInt(ab90Value.Int() + 1)
	case ab90Reflect.Uint8, ab90Reflect.Uint16:
		ab90Value.SetUint(ab90Value.Uint() + 1)
	}
}

func TestAb90DeepCopyIsolation(ab90T *ab90Testing.T) {
	ab90Before := ab90Catalog(ab90T)
	ab90Copy := ab90Catalog(ab90T)
	ab90Mutate(ab90Reflect.ValueOf(&ab90Copy).Elem())
	if ab90Reflect.DeepEqual(ab90Before, ab90Copy) {
		ab90T.Fatal("ab90 mutation control did not change")
	}
	if !ab90Reflect.DeepEqual(ab90Before, ab90Catalog(ab90T)) {
		ab90T.Fatal("ab90 accessor shares nested storage")
	}
	for _, ab90Expected := range ab90Before.Ab90Fields {
		ab90Result, ab90Found, ab90Error := Ab90LookupB51AMirrorField(ab90Expected.Ab90Key)
		if ab90Error != nil || !ab90Found {
			ab90T.Fatal("ab90 lookup control missing")
		}
		ab90Mutate(ab90Reflect.ValueOf(&ab90Result).Elem())
		ab90Fresh, ab90Found, ab90Error := Ab90LookupB51AMirrorField(ab90Expected.Ab90Key)
		if ab90Error != nil || !ab90Found || !ab90Reflect.DeepEqual(ab90Fresh, ab90Expected) {
			ab90T.Fatal("ab90 lookup shares nested storage")
		}
	}
}
