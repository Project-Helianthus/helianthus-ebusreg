package vaillant

import (
	"bytes"
	"encoding/json"
	"reflect"
	"strconv"
	"testing"
)

func tunnelDecode(tunnelT *testing.T, tunnelData []byte) map[string]any {
	tunnelT.Helper()
	var tunnelValue map[string]any
	if tunnelError := json.Unmarshal(tunnelData, &tunnelValue); tunnelError != nil {
		tunnelT.Fatal(tunnelError)
	}
	return tunnelValue
}

func tunnelExpected(tunnelT *testing.T) map[string]any {
	tunnelT.Helper()
	tunnelDestination := `{"name":"destination","logical_position":null,"width_bits":8,"extent":null,"raw_domain":{"min":0,"max":255},"signedness":"unsigned","byte_order":null,"conversion":null,"unit":null,"default":8,"mappings":{},"reserved_regions":null,"grade":"Hypothesis","raw_fallback":"preserve unchanged bytes, including empty data"}`
	tunnelData := `{"name":"data","logical_position":null,"width_bits":null,"extent":{"through":"end-of-PDU","abstract_min":0,"abstract_max":65535,"transport_capacity":null},"raw_domain":"opaque bytes","signedness":"Unknown","byte_order":null,"conversion":null,"unit":null,"default":null,"mappings":{},"reserved_regions":null,"grade":"Hypothesis","raw_fallback":"preserve unchanged bytes, including empty data"}`
	tunnelB5 := `{"name":"constant","logical_position":0,"width_bits":8,"extent":null,"raw_domain":{"constant":181},"signedness":"unsigned","byte_order":null,"conversion":null,"unit":null,"default":null,"mappings":{},"reserved_regions":null,"grade":"Hypothesis","raw_fallback":"preserve unchanged bytes, including empty data"}`
	tunnelF5 := `{"name":"constant","logical_position":0,"width_bits":8,"extent":null,"raw_domain":{"constant":245},"signedness":"unsigned","byte_order":null,"conversion":null,"unit":null,"default":null,"mappings":{},"reserved_regions":null,"grade":"Hypothesis","raw_fallback":"preserve unchanged bytes, including empty data"}`
	return tunnelDecode(tunnelT, []byte(`{
 "catalog":"ab37_tunnel",
 "representation_grade":"Proven",
 "coordinate_system":"logical positions include prefix positions; not physical offsets",
 "rows":[
 {"id":"ab37_b517_request","family":"B517","variant":"declared","direction":"request","service":"B5","operation":"17","alternative":"declared_request","body_extent":null,"fields":[`+tunnelDestination+`,`+tunnelData+`],"grade":"Hypothesis","raw_fallback":"preserve unchanged bytes, including empty data"},
 {"id":"ab37_b517_positive","family":"B517","variant":"declared","direction":"positive","service":"B5","operation":"18","alternative":"declared_positive","body_extent":null,"fields":[`+tunnelData+`],"grade":"Hypothesis","raw_fallback":"preserve unchanged bytes, including empty data"},
 {"id":"ab37_b517_negative","family":"B517","variant":"declared","direction":"negative","service":null,"operation":null,"alternative":"declared_negative","body_extent":null,"fields":[],"grade":"Hypothesis","raw_fallback":"preserve unchanged bytes, including empty data"},
 {"id":"ab37_b518_request","family":"B518","variant":"partial","direction":"request","service":null,"operation":null,"alternative":"partial_request","body_extent":null,"fields":[`+tunnelB5+`],"grade":"Hypothesis","raw_fallback":"preserve unchanged bytes, including empty data"},
 {"id":"ab37_b518_positive","family":"B518","variant":"partial","direction":"positive","service":null,"operation":null,"alternative":"partial_positive","body_extent":null,"fields":[`+tunnelF5+`],"grade":"Hypothesis","raw_fallback":"preserve unchanged bytes, including empty data"},
 {"id":"ab37_b518_negative","family":"B518","variant":"partial","direction":"negative","service":null,"operation":null,"alternative":"partial_negative","body_extent":null,"fields":[],"grade":"Hypothesis","raw_fallback":"preserve unchanged bytes, including empty data"}],
 "tunnel_client_policy":{"grade":"Hypothesis","destination_domain":{"min":1,"max":254},"inner_destination":"original destination","opaque_data":"original PDU","returned_data_context":"original service","local_result_codes":[10,6,7],"local_result_meanings":null,"remote_feedback":"separate from local result codes","default_8":"not an installed address"},
 "adaptation_alternative":{"grade":"Hypothesis","ready_state":3,"final_option":"nonzero","length_equation":{"input_byte":4,"add":6,"equals":"available length"},"request_check":{"service":"B5","operation":"17"},"input_check":{"service":"B5","operation":"18"},"universal_prefix_removal":false,"qualification":"all conditions required; length equation alone is insufficient"},
 "limitations":{"physical_relationship":"Unknown","destination_roles":"Unknown","forwarder":"Unknown","correlation":"Unknown","timeout":"Unknown","refusal":"Unknown","fragmentation":"Unknown","transport_support":"Unknown","peer_support":"Unknown","scheduling":"Unknown","write_permission":"Unknown"},
 "admission":{"runtime_registration":false,"wire_builder":false,"automatic_selection":false,"native_codec_promotion":false,"executable_capability":false,"physical_qualification":"Unknown"}
 }`))
}

func tunnelAssertContract(tunnelT *testing.T, tunnelGot, tunnelWant any) {
	tunnelT.Helper()
	switch tunnelExpected := tunnelWant.(type) {
	case map[string]any:
		tunnelActual, tunnelOK := tunnelGot.(map[string]any)
		if !tunnelOK || len(tunnelActual) != len(tunnelExpected) {
			tunnelT.Fatalf("ab37 object shape differs: %v", tunnelGot)
		}
		for tunnelKey, tunnelValue := range tunnelExpected {
			tunnelT.Run(tunnelKey, func(tunnelT *testing.T) {
				tunnelFound, tunnelPresent := tunnelActual[tunnelKey]
				if !tunnelPresent {
					tunnelT.Fatal("ab37 absent property")
				}
				tunnelAssertContract(tunnelT, tunnelFound, tunnelValue)
			})
		}
	case []any:
		tunnelActual, tunnelOK := tunnelGot.([]any)
		if !tunnelOK || len(tunnelActual) != len(tunnelExpected) {
			tunnelT.Fatalf("ab37 alternative or field count differs: %v", tunnelGot)
		}
		for tunnelIndex, tunnelValue := range tunnelExpected {
			tunnelT.Run(strconv.Itoa(tunnelIndex), func(tunnelT *testing.T) { tunnelAssertContract(tunnelT, tunnelActual[tunnelIndex], tunnelValue) })
		}
	default:
		if !reflect.DeepEqual(tunnelGot, tunnelWant) {
			tunnelT.Fatalf("ab37 got %v; want %v", tunnelGot, tunnelWant)
		}
	}
}

func TestAb37FiniteCatalog(tunnelT *testing.T) {
	tunnelAssertContract(tunnelT, tunnelDecode(tunnelT, TunnelCatalog()), tunnelExpected(tunnelT))
}

func TestAb37IdentityAxes(tunnelT *testing.T) {
	tunnelWant := tunnelExpected(tunnelT)
	for tunnelIndex := range tunnelWant["rows"].([]any) {
		for _, tunnelAxis := range []string{"id", "family", "variant", "direction", "service", "operation", "alternative"} {
			tunnelT.Run(tunnelWant["rows"].([]any)[tunnelIndex].(map[string]any)["id"].(string)+"/"+tunnelAxis, func(tunnelT *testing.T) {
				tunnelGot := tunnelDecode(tunnelT, TunnelCatalog())
				if !reflect.DeepEqual(tunnelGot, tunnelWant) {
					tunnelT.Fatal("ab37 invalid control")
				}
				tunnelRow := tunnelGot["rows"].([]any)[tunnelIndex].(map[string]any)
				tunnelRow[tunnelAxis] = "ab37_changed"
				tunnelEncoded, tunnelError := json.Marshal(tunnelGot)
				if tunnelError != nil {
					tunnelT.Fatal(tunnelError)
				}
				if reflect.DeepEqual(tunnelDecode(tunnelT, tunnelEncoded), tunnelWant) {
					tunnelT.Fatal("ab37 changed identity collapsed")
				}
				delete(tunnelRow, tunnelAxis)
				if reflect.DeepEqual(tunnelGot, tunnelWant) {
					tunnelT.Fatal("ab37 dropped identity accepted")
				}
			})
		}
	}
}

func TestAb37CopyOwnership(tunnelT *testing.T) {
	tunnelFirst := TunnelCatalog()
	tunnelControl := bytes.Clone(tunnelFirst)
	for tunnelIndex := range tunnelFirst {
		tunnelFirst[tunnelIndex] = 0
	}
	if !bytes.Equal(TunnelCatalog(), tunnelControl) {
		tunnelT.Fatal("ab37 accessor aliases embedded data")
	}
}

func TestAb37PartialAndLengthOnly(tunnelT *testing.T) {
	tunnelCatalog := tunnelDecode(tunnelT, TunnelCatalog())
	tunnelRows := tunnelCatalog["rows"].([]any)
	tunnelPositive := tunnelRows[4].(map[string]any)
	tunnelFields := tunnelPositive["fields"].([]any)
	if tunnelPositive["operation"] != nil || tunnelPositive["body_extent"] != nil || len(tunnelFields) != 1 || tunnelFields[0].(map[string]any)["raw_domain"].(map[string]any)["constant"] != float64(245) {
		tunnelT.Fatal("ab37 B518 acquired subgroup or body")
	}
	tunnelOpaque := tunnelRows[1].(map[string]any)["fields"].([]any)[0].(map[string]any)
	if tunnelOpaque["width_bits"] != nil || tunnelOpaque["conversion"] != nil || tunnelOpaque["extent"].(map[string]any)["transport_capacity"] != nil {
		tunnelT.Fatal("ab37 opaque bytes promoted")
	}
	tunnelEmpty := []byte{}
	if len(tunnelEmpty) != 0 || tunnelPositive["body_extent"] != nil {
		tunnelT.Fatal("ab37 empty data ceased to be partial")
	}
	tunnelEquation := tunnelCatalog["adaptation_alternative"].(map[string]any)["length_equation"].(map[string]any)
	tunnelEnvelope := make([]byte, 10)
	tunnelEnvelope[4] = 4
	tunnelPosition := int(tunnelEquation["input_byte"].(float64))
	tunnelAdd := int(tunnelEquation["add"].(float64))
	if int(tunnelEnvelope[tunnelPosition])+tunnelAdd != len(tunnelEnvelope) {
		tunnelT.Fatal("ab37 declared equation control failed")
	}
	tunnelEnvelope[4] = 5
	if int(tunnelEnvelope[tunnelPosition])+tunnelAdd == len(tunnelEnvelope) {
		tunnelT.Fatal("ab37 changed byte satisfies equation")
	}
	if tunnelCatalog["adaptation_alternative"].(map[string]any)["qualification"] != "all conditions required; length equation alone is insufficient" {
		tunnelT.Fatal("ab37 equation became qualification")
	}
}
