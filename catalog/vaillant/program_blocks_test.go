package vaillant

import (
	programBlockBytes "bytes"
	programBlockBinary "encoding/binary"
	programBlockJSON "encoding/json"
	programBlockFmt "fmt"
	programBlockReflect "reflect"
	programBlockTesting "testing"
)

func programBlockRead(programBlockT *programBlockTesting.T) map[string]any {
	programBlockT.Helper()
	var programBlockCatalog map[string]any
	if programBlockErr := programBlockJSON.Unmarshal(ProgramBlocks(), &programBlockCatalog); programBlockErr != nil {
		programBlockT.Fatal(programBlockErr)
	}
	return programBlockCatalog
}

func programBlockField(programBlockName string, programBlockPosition, programBlockWidth, programBlockMaximum, programBlockDefault any, programBlockUnit any, programBlockOrder, programBlockConversion string, programBlockMapping map[string]any) map[string]any {
	programBlockSigned := "unsigned"
	if programBlockWidth == nil {
		programBlockSigned = "Unknown"
	}
	return map[string]any{"name": programBlockName, "position": programBlockPosition, "width": programBlockWidth, "raw_domain": map[string]any{"minimum": 0, "maximum": programBlockMaximum}, "signedness": programBlockSigned, "byte_order": programBlockOrder, "conversion": programBlockConversion, "unit": programBlockUnit, "default": programBlockDefault, "mappings": programBlockMapping, "reserved_regions": nil, "grade": "Hypothesis", "raw_fallback": "preserve raw bytes"}
}

func programBlockExpected() map[string]any {
	programBlockRows := []any{}
	for programBlockOp := 6; programBlockOp <= 14; programBlockOp++ {
		for _, programBlockAlternative := range []string{"request", "positive", "negative"} {
			programBlockFields := []any{}
			var programBlockExtent any
			var programBlockExpression any
			var programBlockPrefix any
			programBlockDirection := "response"
			if programBlockAlternative != "negative" {
				programBlockPrefix = []any{map[string]any{"position": 0, "value": 181}, map[string]any{"position": 1, "value": 27}, map[string]any{"position": 2, "value": programBlockOp}}
			}
			if programBlockAlternative == "request" {
				programBlockDirection = "request"
				programBlockExtent = 3
				switch programBlockOp {
				case 8:
					programBlockExtent = 11
					programBlockFields = append(programBlockFields, programBlockField("block", 3, 4, 4294967295, 0, nil, "little-endian", "identity", map[string]any{}), programBlockField("segment", 7, 4, 4294967295, 0, nil, "little-endian", "identity", map[string]any{}))
				case 9:
					programBlockExtent = nil
					programBlockExpression = "5+k"
					programBlockData := programBlockField("data", 5, nil, nil, nil, nil, "Unknown", "opaque", map[string]any{})
					programBlockData["length"] = map[string]any{"minimum": 0, "maximum": nil}
					programBlockFields = append(programBlockFields, programBlockField("offset", 3, 2, 65535, 0, nil, "little-endian", "identity", map[string]any{}), programBlockData)
				case 13:
					programBlockExtent = 4
					programBlockFields = append(programBlockFields, programBlockField("timeout", 3, 1, 255, 15, "minutes", "not-applicable", "identity", map[string]any{}))
				case 14:
					programBlockExtent = 4
					programBlockFields = append(programBlockFields, programBlockField("direction", 3, 1, 255, 240, nil, "not-applicable", "identity", map[string]any{"15": "boot-loader", "240": "application"}))
				}
			}
			if programBlockAlternative == "positive" {
				programBlockExtent = 4
				switch programBlockOp {
				case 6:
					programBlockExtent = 6
					for programBlockPosition := 3; programBlockPosition <= 5; programBlockPosition++ {
						programBlockFields = append(programBlockFields, programBlockField(programBlockFmt.Sprintf("raw_%d", programBlockPosition), programBlockPosition, 1, 255, 0, nil, "not-applicable", "identity", map[string]any{}))
					}
				case 8:
					programBlockFields = append(programBlockFields, programBlockField("feedback", 3, 1, 255, nil, nil, "not-applicable", "identity", map[string]any{"0": "rejected", "51": "block-unavailable", "85": "accepted", "170": "protected"}))
				case 11:
					programBlockExtent = 5
					programBlockFields = append(programBlockFields, programBlockField("CRC", 3, 2, 65535, nil, nil, "little-endian", "identity", map[string]any{}))
				case 12:
					programBlockFields = append(programBlockFields, programBlockField("state", 3, 1, 255, nil, nil, "not-applicable", "identity", map[string]any{"0": "hardware-init", "1": "wait-programming-request", "2": "check-application-CRC", "3": "preparation", "4": "active", "5": "application-running"}))
				default:
					programBlockMapping := map[string]any{}
					for programBlockValue := 0; programBlockValue < 256; programBlockValue++ {
						programBlockLabel := "rejected"
						if programBlockValue == 1 {
							programBlockLabel = "accepted"
						}
						programBlockMapping[programBlockFmt.Sprint(programBlockValue)] = programBlockLabel
					}
					programBlockFields = append(programBlockFields, programBlockField("feedback", 3, 1, 255, nil, nil, "not-applicable", "identity", programBlockMapping))
				}
			}
			programBlockRows = append(programBlockRows, map[string]any{"id": programBlockFmt.Sprintf("ab41_op%02x_%s", programBlockOp, programBlockAlternative), "family": "Vaillant programming", "variant": "B51B partial", "direction": programBlockDirection, "service": "B51B", "operation": programBlockOp, "alternative": programBlockAlternative, "logical_prefix": programBlockPrefix, "extent": programBlockExtent, "extent_expression": programBlockExpression, "fields": programBlockFields, "reserved_regions": nil, "grade": "Hypothesis", "raw_fallback": "preserve complete native bytes", "mutative_declaration": programBlockAlternative == "request" && (programBlockOp == 7 || programBlockOp == 8 || programBlockOp == 9 || programBlockOp == 10 || programBlockOp == 13 || programBlockOp == 14)})
		}
	}
	programBlockExpectedValue := map[string]any{"representation_grade": "Proven", "candidate_grade": "Hypothesis", "qualification": "Unknown", "coordinates": "logical including prefix; not physical offsets", "limitations": []string{"negative response positions and extents Unknown", "reserved regions Unknown", "full forwarding, segmentation and local CRC variants Unknown", "physical addresses, response adaptation and geometry Unknown", "peer CRC algorithm, coverage and persistence Unknown", "empty OP09 data does not establish peer acceptance", "commit feedback does not prove persistence", "FFFFFFFF block is not a sentinel", "OP08 position 2 is operation, not block", "OP0D zero has no special meaning", "no runtime registration, wire builder, automatic selection or write admission"}, "entries": programBlockRows}
	programBlockEncoded, _ := programBlockJSON.Marshal(programBlockExpectedValue)
	var programBlockNormalized map[string]any
	_ = programBlockJSON.Unmarshal(programBlockEncoded, &programBlockNormalized)
	return programBlockNormalized
}

func programBlockAssert(programBlockT *programBlockTesting.T, programBlockPath string, programBlockGot, programBlockWant any) {
	programBlockT.Helper()
	if programBlockMap, programBlockOk := programBlockWant.(map[string]any); programBlockOk {
		programBlockActual, programBlockOk := programBlockGot.(map[string]any)
		if !programBlockOk || len(programBlockActual) != len(programBlockMap) {
			programBlockT.Fatalf("%s object differs", programBlockPath)
		}
		for programBlockKey, programBlockValue := range programBlockMap {
			programBlockT.Run(programBlockKey, func(programBlockT *programBlockTesting.T) {
				programBlockAssert(programBlockT, programBlockPath+"/"+programBlockKey, programBlockActual[programBlockKey], programBlockValue)
			})
		}
		return
	}
	if programBlockSlice, programBlockOk := programBlockWant.([]any); programBlockOk {
		programBlockActual, programBlockOk := programBlockGot.([]any)
		if !programBlockOk || len(programBlockActual) != len(programBlockSlice) {
			programBlockT.Fatalf("%s array differs", programBlockPath)
		}
		for programBlockIndex, programBlockValue := range programBlockSlice {
			programBlockT.Run(programBlockFmt.Sprint(programBlockIndex), func(programBlockT *programBlockTesting.T) {
				programBlockAssert(programBlockT, programBlockPath, programBlockActual[programBlockIndex], programBlockValue)
			})
		}
		return
	}
	if !programBlockReflect.DeepEqual(programBlockGot, programBlockWant) {
		programBlockT.Fatalf("%s got %#v want %#v", programBlockPath, programBlockGot, programBlockWant)
	}
}

func TestAb41CatalogContract(programBlockT *programBlockTesting.T) {
	programBlockAssert(programBlockT, "ab41", programBlockRead(programBlockT), programBlockExpected())
}

func TestAb41IndependentBytes(programBlockT *programBlockTesting.T) {
	programBlockFirst := ProgramBlocks()
	programBlockControl := append([]byte(nil), programBlockFirst...)
	for programBlockIndex := range programBlockFirst {
		programBlockFirst[programBlockIndex] ^= 255
	}
	if !programBlockBytes.Equal(ProgramBlocks(), programBlockControl) {
		programBlockT.Fatal("ab41 accessor aliases embedded storage")
	}
}

func TestAb41IdentityAxes(programBlockT *programBlockTesting.T) {
	for _, programBlockAxis := range []string{"family", "variant", "direction", "service", "operation", "alternative"} {
		programBlockT.Run(programBlockAxis, func(programBlockT *programBlockTesting.T) {
			programBlockControl := programBlockRead(programBlockT)["entries"].([]any)[0].(map[string]any)
			programBlockEncoded, _ := programBlockJSON.Marshal(programBlockControl)
			var programBlockChanged map[string]any
			_ = programBlockJSON.Unmarshal(programBlockEncoded, &programBlockChanged)
			if programBlockAxis == "operation" {
				programBlockChanged[programBlockAxis] = float64(7)
			} else {
				programBlockChanged[programBlockAxis] = "ab41_distinct"
			}
			programBlockEncoded, _ = programBlockJSON.Marshal([]any{programBlockControl, programBlockChanged})
			var programBlockRoundtrip []map[string]any
			_ = programBlockJSON.Unmarshal(programBlockEncoded, &programBlockRoundtrip)
			if len(programBlockRoundtrip) != 2 || programBlockReflect.DeepEqual(programBlockRoundtrip[0], programBlockRoundtrip[1]) || !programBlockReflect.DeepEqual(programBlockRoundtrip[0], programBlockControl) || !programBlockReflect.DeepEqual(programBlockRoundtrip[1], programBlockChanged) {
				programBlockT.Fatal("ab41 identity axis lost")
			}
		})
	}
}

func TestAb41SyntheticInterpretation(programBlockT *programBlockTesting.T) {
	programBlockRows := programBlockRead(programBlockT)["entries"].([]any)
	programBlockBlock := []byte{181, 27, 8, 4, 3, 2, 1, 8, 7, 6, 5}
	programBlockFields := programBlockRows[6].(map[string]any)["fields"].([]any)
	for programBlockIndex, programBlockWant := range []uint32{0x01020304, 0x05060708} {
		programBlockFieldValue := programBlockFields[programBlockIndex].(map[string]any)
		programBlockPosition := int(programBlockFieldValue["position"].(float64))
		if programBlockFieldValue["byte_order"] != "little-endian" || programBlockFieldValue["width"] != float64(4) || programBlockBinary.LittleEndian.Uint32(programBlockBlock[programBlockPosition:]) != programBlockWant {
			programBlockT.Fatal("ab41 block/segment or endian differs")
		}
	}
	programBlockData := []byte{181, 27, 9, 0x34, 0x12, 0, 255}
	programBlockOffset := programBlockRows[9].(map[string]any)["fields"].([]any)[0].(map[string]any)
	if programBlockBinary.LittleEndian.Uint16(programBlockData[int(programBlockOffset["position"].(float64)):]) != 0x1234 || !programBlockBytes.Equal(programBlockData[5:], []byte{0, 255}) {
		programBlockT.Fatal("ab41 opaque offset vector differs")
	}
	programBlockTimeout := programBlockRows[21].(map[string]any)["fields"].([]any)[0].(map[string]any)
	if programBlockTimeout["unit"] != "minutes" || programBlockTimeout["default"] != float64(byte(15)) || programBlockTimeout["conversion"] != "identity" {
		programBlockT.Fatal("ab41 minutes changed")
	}
	programBlockState := programBlockRows[19].(map[string]any)["fields"].([]any)[0].(map[string]any)["mappings"].(map[string]any)
	if _, programBlockFound := programBlockState["6"]; programBlockFound {
		programBlockT.Fatal("ab41 state 6 labeled")
	}
	programBlockSpecial := programBlockRows[7].(map[string]any)["fields"].([]any)[0].(map[string]any)["mappings"].(map[string]any)
	programBlockGeneric := programBlockRows[4].(map[string]any)["fields"].([]any)[0].(map[string]any)["mappings"].(map[string]any)
	if programBlockSpecial["85"] != "accepted" || programBlockSpecial["170"] != "protected" || programBlockGeneric["85"] != "rejected" || programBlockGeneric["170"] != "rejected" {
		programBlockT.Fatal("ab41 feedback tables conflated")
	}
	if programBlockBinary.LittleEndian.Uint16([]byte{0x34, 0x12}) != 0x1234 {
		programBlockT.Fatal("ab41 CRC numeric control")
	}
}
