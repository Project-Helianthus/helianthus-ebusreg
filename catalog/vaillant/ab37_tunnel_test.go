package vaillant

import (
	"bytes"
	"encoding/json"
	"reflect"
	"strconv"
	"testing"
)

func ab37Decode(ab37T *testing.T, ab37Data []byte) map[string]any {
	ab37T.Helper()
	var ab37Value map[string]any
	if ab37Error := json.Unmarshal(ab37Data, &ab37Value); ab37Error != nil {
		ab37T.Fatal(ab37Error)
	}
	return ab37Value
}

func ab37Expected(ab37T *testing.T) map[string]any {
	ab37T.Helper()
	ab37Destination := `{"name":"destination","logical_position":null,"width_bits":8,"extent":null,"raw_domain":{"min":0,"max":255},"signedness":"unsigned","byte_order":null,"conversion":null,"unit":null,"default":8,"mappings":{},"reserved_regions":null,"grade":"Hypothesis","raw_fallback":"preserve unchanged bytes, including empty data"}`
	ab37Data := `{"name":"data","logical_position":null,"width_bits":null,"extent":{"through":"end-of-PDU","abstract_min":0,"abstract_max":65535,"transport_capacity":null},"raw_domain":"opaque bytes","signedness":"Unknown","byte_order":null,"conversion":null,"unit":null,"default":null,"mappings":{},"reserved_regions":null,"grade":"Hypothesis","raw_fallback":"preserve unchanged bytes, including empty data"}`
	ab37B5 := `{"name":"constant","logical_position":0,"width_bits":8,"extent":null,"raw_domain":{"constant":181},"signedness":"unsigned","byte_order":null,"conversion":null,"unit":null,"default":null,"mappings":{},"reserved_regions":null,"grade":"Hypothesis","raw_fallback":"preserve unchanged bytes, including empty data"}`
	ab37F5 := `{"name":"constant","logical_position":0,"width_bits":8,"extent":null,"raw_domain":{"constant":245},"signedness":"unsigned","byte_order":null,"conversion":null,"unit":null,"default":null,"mappings":{},"reserved_regions":null,"grade":"Hypothesis","raw_fallback":"preserve unchanged bytes, including empty data"}`
	return ab37Decode(ab37T, []byte(`{
 "catalog":"ab37_tunnel",
 "representation_grade":"Proven",
 "coordinate_system":"logical positions include prefix positions; not physical offsets",
 "rows":[
 {"id":"ab37_b517_request","family":"B517","variant":"declared","direction":"request","service":"B5","operation":"17","alternative":"declared_request","body_extent":null,"fields":[`+ab37Destination+`,`+ab37Data+`],"grade":"Hypothesis","raw_fallback":"preserve unchanged bytes, including empty data"},
 {"id":"ab37_b517_positive","family":"B517","variant":"declared","direction":"positive","service":"B5","operation":"18","alternative":"declared_positive","body_extent":null,"fields":[`+ab37Data+`],"grade":"Hypothesis","raw_fallback":"preserve unchanged bytes, including empty data"},
 {"id":"ab37_b517_negative","family":"B517","variant":"declared","direction":"negative","service":null,"operation":null,"alternative":"declared_negative","body_extent":null,"fields":[],"grade":"Hypothesis","raw_fallback":"preserve unchanged bytes, including empty data"},
 {"id":"ab37_b518_request","family":"B518","variant":"partial","direction":"request","service":null,"operation":null,"alternative":"partial_request","body_extent":null,"fields":[`+ab37B5+`],"grade":"Hypothesis","raw_fallback":"preserve unchanged bytes, including empty data"},
 {"id":"ab37_b518_positive","family":"B518","variant":"partial","direction":"positive","service":null,"operation":null,"alternative":"partial_positive","body_extent":null,"fields":[`+ab37F5+`],"grade":"Hypothesis","raw_fallback":"preserve unchanged bytes, including empty data"},
 {"id":"ab37_b518_negative","family":"B518","variant":"partial","direction":"negative","service":null,"operation":null,"alternative":"partial_negative","body_extent":null,"fields":[],"grade":"Hypothesis","raw_fallback":"preserve unchanged bytes, including empty data"}],
 "tunnel_client_policy":{"grade":"Hypothesis","destination_domain":{"min":1,"max":254},"inner_destination":"original destination","opaque_data":"original PDU","returned_data_context":"original service","local_result_codes":[10,6,7],"local_result_meanings":null,"remote_feedback":"separate from local result codes","default_8":"not an installed address"},
 "adaptation_alternative":{"grade":"Hypothesis","ready_state":3,"final_option":"nonzero","length_equation":{"input_byte":4,"add":6,"equals":"available length"},"request_check":{"service":"B5","operation":"17"},"input_check":{"service":"B5","operation":"18"},"universal_prefix_removal":false,"qualification":"all conditions required; length equation alone is insufficient"},
 "limitations":{"physical_relationship":"Unknown","destination_roles":"Unknown","forwarder":"Unknown","correlation":"Unknown","timeout":"Unknown","refusal":"Unknown","fragmentation":"Unknown","transport_support":"Unknown","peer_support":"Unknown","scheduling":"Unknown","write_permission":"Unknown"},
 "admission":{"runtime_registration":false,"wire_builder":false,"automatic_selection":false,"native_codec_promotion":false,"executable_capability":false,"physical_qualification":"Unknown"}
 }`))
}

func ab37AssertContract(ab37T *testing.T, ab37Got, ab37Want any) {
	ab37T.Helper()
	switch ab37Expected := ab37Want.(type) {
	case map[string]any:
		ab37Actual, ab37OK := ab37Got.(map[string]any)
		if !ab37OK || len(ab37Actual) != len(ab37Expected) {
			ab37T.Fatalf("ab37 object shape differs: %v", ab37Got)
		}
		for ab37Key, ab37Value := range ab37Expected {
			ab37T.Run(ab37Key, func(ab37T *testing.T) {
				ab37Found, ab37Present := ab37Actual[ab37Key]
				if !ab37Present {
					ab37T.Fatal("ab37 absent property")
				}
				ab37AssertContract(ab37T, ab37Found, ab37Value)
			})
		}
	case []any:
		ab37Actual, ab37OK := ab37Got.([]any)
		if !ab37OK || len(ab37Actual) != len(ab37Expected) {
			ab37T.Fatalf("ab37 alternative or field count differs: %v", ab37Got)
		}
		for ab37Index, ab37Value := range ab37Expected {
			ab37T.Run(strconv.Itoa(ab37Index), func(ab37T *testing.T) { ab37AssertContract(ab37T, ab37Actual[ab37Index], ab37Value) })
		}
	default:
		if !reflect.DeepEqual(ab37Got, ab37Want) {
			ab37T.Fatalf("ab37 got %v; want %v", ab37Got, ab37Want)
		}
	}
}

func TestAb37FiniteCatalog(ab37T *testing.T) {
	ab37AssertContract(ab37T, ab37Decode(ab37T, Ab37TunnelCatalog()), ab37Expected(ab37T))
}

func TestAb37IdentityAxes(ab37T *testing.T) {
	ab37Want := ab37Expected(ab37T)
	for ab37Index := range ab37Want["rows"].([]any) {
		for _, ab37Axis := range []string{"id", "family", "variant", "direction", "service", "operation", "alternative"} {
			ab37T.Run(ab37Want["rows"].([]any)[ab37Index].(map[string]any)["id"].(string)+"/"+ab37Axis, func(ab37T *testing.T) {
				ab37Got := ab37Decode(ab37T, Ab37TunnelCatalog())
				if !reflect.DeepEqual(ab37Got, ab37Want) {
					ab37T.Fatal("ab37 invalid control")
				}
				ab37Row := ab37Got["rows"].([]any)[ab37Index].(map[string]any)
				ab37Row[ab37Axis] = "ab37_changed"
				ab37Encoded, ab37Error := json.Marshal(ab37Got)
				if ab37Error != nil {
					ab37T.Fatal(ab37Error)
				}
				if reflect.DeepEqual(ab37Decode(ab37T, ab37Encoded), ab37Want) {
					ab37T.Fatal("ab37 changed identity collapsed")
				}
				delete(ab37Row, ab37Axis)
				if reflect.DeepEqual(ab37Got, ab37Want) {
					ab37T.Fatal("ab37 dropped identity accepted")
				}
			})
		}
	}
}

func TestAb37CopyOwnership(ab37T *testing.T) {
	ab37First := Ab37TunnelCatalog()
	ab37Control := bytes.Clone(ab37First)
	for ab37Index := range ab37First {
		ab37First[ab37Index] = 0
	}
	if !bytes.Equal(Ab37TunnelCatalog(), ab37Control) {
		ab37T.Fatal("ab37 accessor aliases embedded data")
	}
}

func TestAb37PartialAndLengthOnly(ab37T *testing.T) {
	ab37Catalog := ab37Decode(ab37T, Ab37TunnelCatalog())
	ab37Rows := ab37Catalog["rows"].([]any)
	ab37Positive := ab37Rows[4].(map[string]any)
	ab37Fields := ab37Positive["fields"].([]any)
	if ab37Positive["operation"] != nil || ab37Positive["body_extent"] != nil || len(ab37Fields) != 1 || ab37Fields[0].(map[string]any)["raw_domain"].(map[string]any)["constant"] != float64(245) {
		ab37T.Fatal("ab37 B518 acquired subgroup or body")
	}
	ab37Opaque := ab37Rows[1].(map[string]any)["fields"].([]any)[0].(map[string]any)
	if ab37Opaque["width_bits"] != nil || ab37Opaque["conversion"] != nil || ab37Opaque["extent"].(map[string]any)["transport_capacity"] != nil {
		ab37T.Fatal("ab37 opaque bytes promoted")
	}
	ab37Empty := []byte{}
	if len(ab37Empty) != 0 || ab37Positive["body_extent"] != nil {
		ab37T.Fatal("ab37 empty data ceased to be partial")
	}
	ab37Equation := ab37Catalog["adaptation_alternative"].(map[string]any)["length_equation"].(map[string]any)
	ab37Envelope := make([]byte, 10)
	ab37Envelope[4] = 4
	ab37Position := int(ab37Equation["input_byte"].(float64))
	ab37Add := int(ab37Equation["add"].(float64))
	if int(ab37Envelope[ab37Position])+ab37Add != len(ab37Envelope) {
		ab37T.Fatal("ab37 declared equation control failed")
	}
	ab37Envelope[4] = 5
	if int(ab37Envelope[ab37Position])+ab37Add == len(ab37Envelope) {
		ab37T.Fatal("ab37 changed byte satisfies equation")
	}
	if ab37Catalog["adaptation_alternative"].(map[string]any)["qualification"] != "all conditions required; length equation alone is insufficient" {
		ab37T.Fatal("ab37 equation became qualification")
	}
}
