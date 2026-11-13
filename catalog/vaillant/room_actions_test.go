package vaillant

import (
	roomActionBytes "bytes"
	roomActionJSON "encoding/json"
	roomActionFmt "fmt"
	roomActionReflect "reflect"
	roomActionTesting "testing"
)

func roomActionDecode(roomActionT *roomActionTesting.T, roomActionData []byte) map[string]any {
	roomActionT.Helper()
	var roomActionValue map[string]any
	if roomActionError := roomActionJSON.Unmarshal(roomActionData, &roomActionValue); roomActionError != nil {
		roomActionT.Fatal(roomActionError)
	}
	return roomActionValue
}

func roomActionRows(roomActionT *roomActionTesting.T) []any {
	roomActionT.Helper()
	return roomActionDecode(roomActionT, RoomActions())["rows"].([]any)
}

func roomActionEqual(roomActionT *roomActionTesting.T, roomActionActual, roomActionExpected any) {
	roomActionT.Helper()
	if !roomActionReflect.DeepEqual(roomActionActual, roomActionExpected) {
		roomActionT.Fatalf("got %#v; want %#v", roomActionActual, roomActionExpected)
	}
}

func TestAb29Copy(roomActionT *roomActionTesting.T) {
	roomActionFirst, roomActionControl := RoomActions(), RoomActions()
	roomActionFirst[0] = '!'
	if !roomActionBytes.Equal(RoomActions(), roomActionControl) {
		roomActionT.Fatal("accessor shares bytes")
	}
	roomActionControl[len(roomActionControl)-1] = '!'
	if !roomActionJSON.Valid(RoomActions()) {
		roomActionT.Fatal("second copy aliases catalog")
	}
}

func TestAb29Rows(roomActionT *roomActionTesting.T) {
	roomActionDoc := roomActionDecode(roomActionT, RoomActions())
	for roomActionKey, roomActionExpected := range map[string]any{
		"catalog": "ab29_room_actions", "representation_grade": "Proven", "candidate_grade": "Hypothesis", "physical_qualification": "Unknown",
		"coordinate_system": "logical positions include prefixes; not physical offsets",
	} {
		roomActionT.Run("ab29_"+roomActionKey, func(roomActionT *roomActionTesting.T) {
			roomActionEqual(roomActionT, roomActionDoc[roomActionKey], roomActionExpected)
		})
	}
	roomActionRowsValue := roomActionRows(roomActionT)
	roomActionEqual(roomActionT, len(roomActionRowsValue), 18)
	roomActionCases := []struct {
		roomActionFamily, roomActionOperation, roomActionVariant, roomActionDirection string
		roomActionExtent                                                              any
		roomActionNames                                                               []any
		roomActionLimitations                                                         []any
	}{
		{"B50A", "OP00", "room-action-temperature", "read", float64(3), []any{"room_action", "temperature"}, []any{"Unavailable sentinel is Unknown"}},
		{"B50C", "OP00", "keycode-digits", "read", float64(4), []any{"digit_1", "digit_2", "digit_3", "digit_4"}, []any{"Four independent numeric digits; not ASCII, BCD or u32"}},
		{"B50C", "OP01", "signed-value", "read", float64(2), []any{"value"}, []any{"Unit and sentinel are Unknown"}},
		{"B50C", "OP02", "circuit-names", "read", nil, []any{"circuit_name"}, []any{"Independent per-circuit values and failures; retain successful values on partial failure", "Bounded aggregate 1..15 is not installed cardinality", "Echoed circuit, string trimming, B524 namespace equivalence and wildcard semantics of zero are Unknown"}},
		{"B50D", "OP00", "ui-room-action", "UI-write", float64(1), []any{"room_action"}, []any{"UI-direction write is distinct from remote-control writes and completed UI effects"}},
		{"B50D", "OP01", "ui-holiday-status", "UI-write", float64(1), []any{"holiday_status"}, []any{"UI-direction write is distinct from remote-control writes and completed UI effects"}},
	}
	roomActionIDs := []string{"ab29_b50a_op00", "ab29_b50c_op00", "ab29_b50c_op01", "ab29_b50c_op02", "ab29_b50d_op00", "ab29_b50d_op01"}
	for roomActionIndex, roomActionCase := range roomActionCases {
		for roomActionAlternativeIndex, roomActionAlternative := range []string{"request", "positive", "negative"} {
			roomActionRow := roomActionRowsValue[roomActionIndex*3+roomActionAlternativeIndex].(map[string]any)
			roomActionExpected := map[string]any{
				"id": roomActionIDs[roomActionIndex] + "_" + roomActionAlternative, "family": roomActionCase.roomActionFamily, "variant": roomActionCase.roomActionVariant, "direction": roomActionCase.roomActionDirection,
				"service": roomActionCase.roomActionFamily, "operation": roomActionCase.roomActionOperation, "alternative": roomActionAlternative, "grade": "Hypothesis", "physical_qualification": "Unknown",
				"logical_positions": nil, "extent_bytes": nil, "raw_domain": nil, "signedness": "Unknown", "byte_order": "Unknown", "conversion": nil, "unit": nil, "default": nil,
				"finite_mappings": map[string]any{}, "reserved_regions": []any{}, "relative_layout": []any{},
				"limitations":  append(append([]any{}, roomActionCase.roomActionLimitations...), "Absolute logical positions including any prefix are Unknown", "Unknown alternative layouts do not assert a qualified response format"),
				"raw_fallback": "preserve raw bytes and unresolved alternatives",
			}
			roomActionKnown := roomActionAlternative == "positive" && roomActionCase.roomActionDirection == "read" || roomActionAlternative == "request" && roomActionCase.roomActionDirection == "UI-write"
			if roomActionKnown {
				roomActionExpected["extent_bytes"], roomActionExpected["raw_domain"], roomActionExpected["signedness"], roomActionExpected["byte_order"], roomActionExpected["relative_layout"] = roomActionCase.roomActionExtent, "field domains", "per field", "per field", roomActionCase.roomActionNames
			} else if roomActionAlternative == "request" && roomActionCase.roomActionDirection == "read" && roomActionCase.roomActionOperation != "OP02" {
				roomActionExpected["extent_bytes"], roomActionExpected["raw_domain"] = float64(0), "no arguments"
			}
			roomActionT.Run(roomActionIDs[roomActionIndex]+"_"+roomActionAlternative, func(roomActionT *roomActionTesting.T) {
				for roomActionKey, roomActionValue := range roomActionExpected {
					roomActionT.Run("ab29_"+roomActionKey, func(roomActionT *roomActionTesting.T) {
						if _, roomActionPresent := roomActionRow[roomActionKey]; !roomActionPresent {
							roomActionT.Fatal("missing key")
						}
						roomActionEqual(roomActionT, roomActionRow[roomActionKey], roomActionValue)
					})
				}
				roomActionEqual(roomActionT, len(roomActionRow), len(roomActionExpected)+1)
				if !roomActionKnown {
					roomActionEqual(roomActionT, roomActionRow["fields"], []any{})
				}
			})
		}
	}
}

func TestAb29Fields(roomActionT *roomActionTesting.T) {
	roomActionExpected := []struct {
		roomActionRow    int
		roomActionFields string
	}{
		{1, `[
 {"name":"room_action","width_bytes":1,"raw_domain":[0,255],"signedness":"unsigned","byte_order":"not applicable","conversion":{"factor":1,"offset":0},"unit":null,"default":null,"finite_mappings":{"0":"no-action","1":"switch-on","2":"thermostat"},"reserved_regions":[{"range":[3,255],"meaning":"unlabeled; retain raw"}],"sentinel":"Unknown"},
 {"name":"temperature","width_bytes":2,"raw_domain":[-32768,32767],"signedness":"signed","byte_order":"little-endian","conversion":{"factor":0.0625,"offset":0,"candidate_range":[-2048,2047.9375]},"unit":"degrees C","default":null,"finite_mappings":{},"reserved_regions":[],"sentinel":"Unknown; neither 8000 nor FFFF is a qualified unavailable marker"}]`},
		{4, `[
 {"name":"digit_1","width_bytes":1,"raw_domain":[0,255],"signedness":"unsigned","byte_order":"not applicable","conversion":{"factor":1,"offset":0},"unit":null,"default":0,"finite_mappings":{},"reserved_regions":[{"range":[10,255],"meaning":"outside candidate digit domain; retain raw"}],"sentinel":"Unknown"},
 {"name":"digit_2","width_bytes":1,"raw_domain":[0,255],"signedness":"unsigned","byte_order":"not applicable","conversion":{"factor":1,"offset":0},"unit":null,"default":0,"finite_mappings":{},"reserved_regions":[{"range":[10,255],"meaning":"outside candidate digit domain; retain raw"}],"sentinel":"Unknown"},
 {"name":"digit_3","width_bytes":1,"raw_domain":[0,255],"signedness":"unsigned","byte_order":"not applicable","conversion":{"factor":1,"offset":0},"unit":null,"default":0,"finite_mappings":{},"reserved_regions":[{"range":[10,255],"meaning":"outside candidate digit domain; retain raw"}],"sentinel":"Unknown"},
 {"name":"digit_4","width_bytes":1,"raw_domain":[0,255],"signedness":"unsigned","byte_order":"not applicable","conversion":{"factor":1,"offset":0},"unit":null,"default":0,"finite_mappings":{},"reserved_regions":[{"range":[10,255],"meaning":"outside candidate digit domain; retain raw"}],"sentinel":"Unknown"}]`},
		{7, `[{"name":"value","width_bytes":2,"raw_domain":[-32768,32767],"signedness":"signed","byte_order":"little-endian","conversion":{"factor":1,"offset":0},"unit":null,"default":null,"finite_mappings":{},"reserved_regions":[],"sentinel":"Unknown"}]`},
		{10, `[{"name":"circuit_name","width_bytes":null,"raw_domain":null,"signedness":"Unknown","byte_order":"Unknown","conversion":null,"unit":null,"default":null,"finite_mappings":{},"reserved_regions":[],"sentinel":"Unknown"}]`},
		{12, `[{"name":"room_action","width_bytes":1,"raw_domain":[0,255],"signedness":"unsigned","byte_order":"not applicable","conversion":{"factor":1,"offset":0},"unit":null,"default":0,"finite_mappings":{"0":"no-action","1":"switch-on","2":"thermostat"},"reserved_regions":[{"range":[3,255],"meaning":"unlabeled; retain raw"}],"sentinel":"Unknown"}]`},
		{15, `[{"name":"holiday_status","width_bytes":1,"raw_domain":[0,255],"signedness":"unsigned","byte_order":"not applicable","conversion":{"factor":1,"offset":0},"unit":null,"default":0,"finite_mappings":{"0":"0","1":"1"},"reserved_regions":[{"range":[2,255],"meaning":"unlabeled; retain raw"}],"sentinel":"Unknown"}]`},
	}
	roomActionRowsValue := roomActionRows(roomActionT)
	for _, roomActionCase := range roomActionExpected {
		var roomActionFields []any
		if roomActionError := roomActionJSON.Unmarshal([]byte(roomActionCase.roomActionFields), &roomActionFields); roomActionError != nil {
			roomActionT.Fatal(roomActionError)
		}
		roomActionActual := roomActionRowsValue[roomActionCase.roomActionRow].(map[string]any)["fields"].([]any)
		roomActionEqual(roomActionT, len(roomActionActual), len(roomActionFields))
		for roomActionIndex, roomActionField := range roomActionFields {
			roomActionWant := roomActionField.(map[string]any)
			roomActionWant["grade"], roomActionWant["raw_fallback"], roomActionWant["logical_positions"] = "Hypothesis", "preserve raw bytes", nil
			roomActionGot := roomActionActual[roomActionIndex].(map[string]any)
			roomActionEqual(roomActionT, len(roomActionGot), len(roomActionWant))
			for roomActionKey, roomActionValue := range roomActionWant {
				roomActionT.Run(roomActionFmt.Sprintf("ab29_row_%d_field_%d_%s", roomActionCase.roomActionRow, roomActionIndex, roomActionKey), func(roomActionT *roomActionTesting.T) {
					if _, roomActionPresent := roomActionGot[roomActionKey]; !roomActionPresent {
						roomActionT.Fatal("missing key")
					}
					roomActionEqual(roomActionT, roomActionGot[roomActionKey], roomActionValue)
				})
			}
		}
	}
}

func roomActionCandidate(roomActionField map[string]any, roomActionRaw []byte) float64 {
	roomActionCarrier := int(roomActionRaw[0])
	if roomActionField["width_bytes"] == float64(2) {
		roomActionCarrier += int(roomActionRaw[1]) << 8
	}
	if roomActionField["signedness"] == "signed" && roomActionCarrier >= 32768 {
		roomActionCarrier -= 65536
	}
	roomActionConversion := roomActionField["conversion"].(map[string]any)
	return float64(roomActionCarrier)*roomActionConversion["factor"].(float64) + roomActionConversion["offset"].(float64)
}

func TestAb29Synthetic(roomActionT *roomActionTesting.T) {
	roomActionRowsValue := roomActionRows(roomActionT)
	roomActionRoom := roomActionRowsValue[1].(map[string]any)["fields"].([]any)
	for _, roomActionCase := range []struct {
		roomActionRaw   []byte
		roomActionLabel string
		roomActionValue float64
	}{
		{[]byte{2, 128, 255}, "thermostat", -8},
		{[]byte{3, 0, 0}, "", 0},
		{[]byte{0, 0, 128}, "no-action", -2048},
		{[]byte{0, 255, 255}, "no-action", -0.0625},
	} {
		roomActionT.Run(roomActionFmt.Sprintf("ab29_%v", roomActionCase.roomActionRaw), func(roomActionT *roomActionTesting.T) {
			roomActionMapping := roomActionRoom[0].(map[string]any)["finite_mappings"].(map[string]any)
			roomActionLabel, _ := roomActionMapping[roomActionFmt.Sprint(roomActionCase.roomActionRaw[0])].(string)
			roomActionEqual(roomActionT, roomActionLabel, roomActionCase.roomActionLabel)
			roomActionEqual(roomActionT, roomActionCandidate(roomActionRoom[1].(map[string]any), roomActionCase.roomActionRaw[1:]), roomActionCase.roomActionValue)
		})
	}
	roomActionDigits := roomActionRowsValue[4].(map[string]any)["fields"].([]any)
	roomActionEqual(roomActionT, len(roomActionDigits), 4)
	for roomActionIndex, roomActionRaw := range []byte{1, 2, 3, 4} {
		roomActionEqual(roomActionT, roomActionCandidate(roomActionDigits[roomActionIndex].(map[string]any), []byte{roomActionRaw}), float64(roomActionIndex+1))
	}
	roomActionNumeric := roomActionRowsValue[7].(map[string]any)["fields"].([]any)[0].(map[string]any)
	roomActionEqual(roomActionT, roomActionCandidate(roomActionNumeric, []byte{254, 255}), float64(-2))
	roomActionEqual(roomActionT, roomActionNumeric["unit"], nil)
	for _, roomActionIndex := range []int{12, 15} {
		roomActionMapping := roomActionRowsValue[roomActionIndex].(map[string]any)["fields"].([]any)[0].(map[string]any)["finite_mappings"].(map[string]any)
		if _, roomActionPresent := roomActionMapping["255"]; roomActionPresent {
			roomActionT.Fatal("unknown byte gained a label")
		}
	}
	// This catalog declares independent circuit outcomes; it supplies no aggregate decoder.
	roomActionCircuitPolicy := roomActionRowsValue[10].(map[string]any)["limitations"].([]any)
	roomActionEqual(roomActionT, roomActionCircuitPolicy[0], "Independent per-circuit values and failures; retain successful values on partial failure")
}
