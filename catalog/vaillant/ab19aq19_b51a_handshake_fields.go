package vaillant

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
)

func ab19aq19Read(ab19aq19Decoder *json.Decoder) (any, error) {
	ab19aq19Token, ab19aq19Err := ab19aq19Decoder.Token()
	if ab19aq19Err != nil {
		return nil, ab19aq19Err
	}
	if ab19aq19Delimiter, ab19aq19OK := ab19aq19Token.(json.Delim); ab19aq19OK {
		switch ab19aq19Delimiter {
		case '{':
			ab19aq19Object := map[string]any{}
			for ab19aq19Decoder.More() {
				ab19aq19KeyToken, ab19aq19Err := ab19aq19Decoder.Token()
				if ab19aq19Err != nil {
					return nil, ab19aq19Err
				}
				ab19aq19Key := ab19aq19KeyToken.(string)
				if _, ab19aq19Exists := ab19aq19Object[ab19aq19Key]; ab19aq19Exists {
					return nil, fmt.Errorf("duplicate key %q", ab19aq19Key)
				}
				ab19aq19Value, ab19aq19Err := ab19aq19Read(ab19aq19Decoder)
				if ab19aq19Err != nil {
					return nil, ab19aq19Err
				}
				ab19aq19Object[ab19aq19Key] = ab19aq19Value
			}
			_, ab19aq19Err := ab19aq19Decoder.Token()
			return ab19aq19Object, ab19aq19Err
		case '[':
			ab19aq19Array := []any{}
			for ab19aq19Decoder.More() {
				ab19aq19Value, ab19aq19Err := ab19aq19Read(ab19aq19Decoder)
				if ab19aq19Err != nil {
					return nil, ab19aq19Err
				}
				ab19aq19Array = append(ab19aq19Array, ab19aq19Value)
			}
			_, ab19aq19Err := ab19aq19Decoder.Token()
			return ab19aq19Array, ab19aq19Err
		}
		return nil, fmt.Errorf("unexpected delimiter")
	}
	if ab19aq19Number, ab19aq19OK := ab19aq19Token.(json.Number); ab19aq19OK {
		if _, ab19aq19Err := ab19aq19Number.Float64(); ab19aq19Err != nil {
			return nil, fmt.Errorf("nonfinite number")
		}
	}
	return ab19aq19Token, nil
}

func ab19aq19Byte(ab19aq19Value any) (int, bool) {
	ab19aq19Number, ab19aq19OK := ab19aq19Value.(json.Number)
	if !ab19aq19OK {
		return 0, false
	}
	ab19aq19Int, ab19aq19Err := ab19aq19Number.Int64()
	return int(ab19aq19Int), ab19aq19Err == nil && ab19aq19Int >= 0 && ab19aq19Int <= 255
}

func Ab19Aq19Evaluate(ab19aq19Input json.RawMessage) (json.RawMessage, error) {
	ab19aq19Decoder := json.NewDecoder(bytes.NewReader(ab19aq19Input))
	ab19aq19Decoder.UseNumber()
	ab19aq19Value, ab19aq19Err := ab19aq19Read(ab19aq19Decoder)
	if ab19aq19Err != nil {
		return nil, ab19aq19Err
	}
	if _, ab19aq19Err = ab19aq19Decoder.Token(); ab19aq19Err != io.EOF {
		return nil, fmt.Errorf("trailing JSON")
	}
	ab19aq19Root, ab19aq19OK := ab19aq19Value.(map[string]any)
	if !ab19aq19OK {
		return nil, fmt.Errorf("input must be object")
	}
	ab19aq19Context, ab19aq19OK := ab19aq19Root["context"].(map[string]any)
	if !ab19aq19OK {
		return nil, fmt.Errorf("context must be object")
	}
	ab19aq19Fields, ab19aq19OK := ab19aq19Root["fields"].(map[string]any)
	if !ab19aq19OK {
		return nil, fmt.Errorf("fields must be object")
	}
	for _, ab19aq19Key := range []string{"primary", "secondary", "operation"} {
		if ab19aq19Value, ab19aq19Present := ab19aq19Context[ab19aq19Key]; ab19aq19Present {
			if _, ab19aq19OK := ab19aq19Byte(ab19aq19Value); !ab19aq19OK {
				return nil, fmt.Errorf("invalid %s", ab19aq19Key)
			}
		}
	}
	for _, ab19aq19Key := range []string{"direction", "view"} {
		if ab19aq19Value, ab19aq19Present := ab19aq19Context[ab19aq19Key]; ab19aq19Present {
			if _, ab19aq19OK := ab19aq19Value.(string); !ab19aq19OK {
				return nil, fmt.Errorf("invalid %s", ab19aq19Key)
			}
		}
	}
	if ab19aq19Value, ab19aq19Present := ab19aq19Root["arguments"]; ab19aq19Present {
		ab19aq19Arguments, ab19aq19OK := ab19aq19Value.(map[string]any)
		if !ab19aq19OK || len(ab19aq19Arguments) != 0 {
			return nil, fmt.Errorf("no arguments supported")
		}
	}
	ab19aq19Bytes := map[string][]int{}
	for ab19aq19Key, ab19aq19Value := range ab19aq19Fields {
		ab19aq19Array, ab19aq19OK := ab19aq19Value.([]any)
		if !ab19aq19OK {
			return nil, fmt.Errorf("field %s must be byte array", ab19aq19Key)
		}
		ab19aq19Bytes[ab19aq19Key] = []int{}
		for _, ab19aq19Value := range ab19aq19Array {
			ab19aq19Int, ab19aq19OK := ab19aq19Byte(ab19aq19Value)
			if !ab19aq19OK {
				return nil, fmt.Errorf("invalid byte in %s", ab19aq19Key)
			}
			ab19aq19Bytes[ab19aq19Key] = append(ab19aq19Bytes[ab19aq19Key], ab19aq19Int)
		}
	}
	ab19aq19Primary, _ := ab19aq19Byte(ab19aq19Context["primary"])
	ab19aq19Secondary, _ := ab19aq19Byte(ab19aq19Context["secondary"])
	ab19aq19Operation, ab19aq19HasOperation := ab19aq19Context["operation"]
	ab19aq19Op, _ := ab19aq19Byte(ab19aq19Operation)
	ab19aq19View, _ := ab19aq19Context["view"].(string)
	ab19aq19Direction, _ := ab19aq19Context["direction"].(string)
	ab19aq19Reset := !ab19aq19HasOperation
	ab19aq19Handshake := ab19aq19HasOperation && ab19aq19Op == 1
	ab19aq19Supported := ab19aq19Primary == 181 && ab19aq19Secondary == 26 && ((ab19aq19View == "reset" && ab19aq19Reset || ab19aq19View == "handshake" && ab19aq19Handshake) && (ab19aq19Direction == "request" || ab19aq19Direction == "positive") || ab19aq19View == "sequence-pair" && ab19aq19Direction == "pair" && (ab19aq19Reset || ab19aq19Handshake))
	ab19aq19Results := map[string]any{}
	for ab19aq19Key, ab19aq19Raw := range ab19aq19Fields {
		ab19aq19Results[ab19aq19Key] = map[string]any{"raw": ab19aq19Raw, "state": "Unknown"}
	}
	ab19aq19Root["results"] = ab19aq19Results
	ab19aq19Root["physical"] = "Unknown"
	ab19aq19Root["qualified"] = false
	ab19aq19Root["state"] = "Unknown"
	if !ab19aq19Supported {
		return json.Marshal(ab19aq19Root)
	}
	ab19aq19Root["state"] = "complete"
	ab19aq19Fixed := func(ab19aq19Key string, ab19aq19Width int) (map[string]any, []int, bool) {
		ab19aq19Raw, ab19aq19Present := ab19aq19Fields[ab19aq19Key]
		ab19aq19Result := map[string]any{"state": "incomplete", "grade": "Hypothesis"}
		if ab19aq19Present {
			ab19aq19Result["raw"] = ab19aq19Raw
		}
		ab19aq19Results[ab19aq19Key] = ab19aq19Result
		ab19aq19Data := ab19aq19Bytes[ab19aq19Key]
		if !ab19aq19Present || len(ab19aq19Data) != ab19aq19Width {
			ab19aq19Root["state"] = "partial"
			return ab19aq19Result, nil, false
		}
		ab19aq19Result["state"] = "complete"
		return ab19aq19Result, ab19aq19Data, true
	}
	if ab19aq19View == "sequence-pair" {
		ab19aq19Request, ab19aq19RequestData, ab19aq19RequestOK := ab19aq19Fixed("request_sequence", 1)
		_, ab19aq19ResponseData, ab19aq19ResponseOK := ab19aq19Fixed("response_sequence", 1)
		if ab19aq19RequestOK && ab19aq19ResponseOK {
			ab19aq19Request["flags"] = map[string]bool{"correlated": ab19aq19RequestData[0] == ab19aq19ResponseData[0]}
			ab19aq19Request["calculation"] = "Proven"
		}
		return json.Marshal(ab19aq19Root)
	}
	if ab19aq19Result, ab19aq19Data, ab19aq19OK := ab19aq19Fixed("sequence", 1); ab19aq19OK {
		ab19aq19Result["number"] = ab19aq19Data[0]
		ab19aq19Result["calculation"] = "Proven"
	}
	if ab19aq19Reset {
		if ab19aq19Direction == "positive" {
			if ab19aq19Result, ab19aq19Data, ab19aq19OK := ab19aq19Fixed("ack", 1); ab19aq19OK {
				ab19aq19Result["flags"] = map[string]bool{"denied": ab19aq19Data[0]&1 != 0}
				ab19aq19Result["unmodeled"] = ab19aq19Data[0] & 254
			}
		}
		return json.Marshal(ab19aq19Root)
	}
	ab19aq19Flags, ab19aq19Data, ab19aq19FlagsOK := ab19aq19Fixed("flags", 1)
	ab19aq19Change := false
	if ab19aq19FlagsOK {
		ab19aq19Labels := []string{"power-on/start", "change", "CRC-support"}
		ab19aq19Mask := 248
		if ab19aq19Direction == "positive" {
			ab19aq19Labels = []string{"valid", "invalid", "blocked", "change", "denied", "invalid-range", "not-user-adjustable"}
			ab19aq19Mask = 128
			ab19aq19Change = ab19aq19Data[0]&8 != 0
		}
		ab19aq19Decoded := map[string]bool{}
		for ab19aq19Index, ab19aq19Label := range ab19aq19Labels {
			ab19aq19Decoded[ab19aq19Label] = ab19aq19Data[0]&(1<<ab19aq19Index) != 0
		}
		ab19aq19Flags["flags"] = ab19aq19Decoded
		ab19aq19Flags["reserved"] = ab19aq19Data[0] & ab19aq19Mask
	}
	for _, ab19aq19Key := range []string{"dsn", "start", "offset"} {
		ab19aq19Result := map[string]any{"state": "Unknown", "grade": "Hypothesis"}
		if ab19aq19Raw, ab19aq19Present := ab19aq19Fields[ab19aq19Key]; ab19aq19Present {
			ab19aq19Result["raw"] = ab19aq19Raw
			ab19aq19Result["state"] = "complete"
		}
		if ab19aq19Direction == "positive" && ab19aq19FlagsOK && !ab19aq19Change {
			ab19aq19Result["state"] = "absent-unless-change"
		}
		ab19aq19Results[ab19aq19Key] = ab19aq19Result
	}
	if ab19aq19Direction == "request" {
		if ab19aq19Result, ab19aq19Data, ab19aq19OK := ab19aq19Fixed("crc", 2); ab19aq19OK {
			ab19aq19Result["number"] = ab19aq19Data[0] + 256*ab19aq19Data[1]
			ab19aq19Result["calculation"] = "Proven"
		}
	} else {
		if ab19aq19Result, ab19aq19Data, ab19aq19OK := ab19aq19Fixed("family_code", 1); ab19aq19OK {
			ab19aq19Result["number"] = ab19aq19Data[0]
			ab19aq19Result["calculation"] = "Proven"
		}
	}
	return json.Marshal(ab19aq19Root)
}
