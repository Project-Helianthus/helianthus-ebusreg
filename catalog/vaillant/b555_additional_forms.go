package vaillant

import (
	"embed"
	"encoding/json"
	"fmt"
	"reflect"
)

type B555AdditionalFormKey struct {
	B555AdditionalFormFamily        string  `json:"family"`
	B555AdditionalFormVariant       string  `json:"variant"`
	B555AdditionalFormPrimary       uint8   `json:"primary"`
	B555AdditionalFormSecondary     *uint8  `json:"secondary"`
	B555AdditionalFormOperation     *uint8  `json:"operation"`
	B555AdditionalFormDirection     string  `json:"direction"`
	B555AdditionalFormField         string  `json:"field"`
	B555AdditionalFormAlternative   string  `json:"alternative"`
	B555AdditionalFormGroup         *uint8  `json:"group"`
	B555AdditionalFormInstanceScope string  `json:"instance_scope"`
	B555AdditionalFormRegister      *uint16 `json:"register"`
}

//go:embed data/ab26_b555_additional_forms.json
var b555AdditionalFormFiles embed.FS

func b555AdditionalFormValidate(b555AdditionalFormData []byte) ([]json.RawMessage, []B555AdditionalFormKey, error) {
	var b555AdditionalFormRows []json.RawMessage
	if b555AdditionalFormError := json.Unmarshal(b555AdditionalFormData, &b555AdditionalFormRows); b555AdditionalFormError != nil {
		return nil, nil, b555AdditionalFormError
	}
	if b555AdditionalFormRows == nil {
		return nil, nil, fmt.Errorf("ab26: expected JSON array")
	}
	b555AdditionalFormKeys := make([]B555AdditionalFormKey, len(b555AdditionalFormRows))
	for b555AdditionalFormIndex, b555AdditionalFormRow := range b555AdditionalFormRows {
		var b555AdditionalFormEnvelope struct {
			B555AdditionalFormKey B555AdditionalFormKey `json:"key"`
		}
		if b555AdditionalFormError := json.Unmarshal(b555AdditionalFormRow, &b555AdditionalFormEnvelope); b555AdditionalFormError != nil {
			return nil, nil, b555AdditionalFormError
		}
		var b555AdditionalFormObject map[string]json.RawMessage
		if b555AdditionalFormError := json.Unmarshal(b555AdditionalFormRow, &b555AdditionalFormObject); b555AdditionalFormError != nil {
			return nil, nil, b555AdditionalFormError
		}
		var b555AdditionalFormIdentity map[string]json.RawMessage
		if b555AdditionalFormError := json.Unmarshal(b555AdditionalFormObject["key"], &b555AdditionalFormIdentity); b555AdditionalFormError != nil {
			return nil, nil, b555AdditionalFormError
		}
		for _, b555AdditionalFormAxis := range []string{"family", "variant", "primary", "secondary", "operation", "direction", "field", "alternative", "group", "instance_scope", "register"} {
			b555AdditionalFormValue, b555AdditionalFormPresent := b555AdditionalFormIdentity[b555AdditionalFormAxis]
			if !b555AdditionalFormPresent || (string(b555AdditionalFormValue) == "null" && b555AdditionalFormAxis != "secondary" && b555AdditionalFormAxis != "operation" && b555AdditionalFormAxis != "group" && b555AdditionalFormAxis != "register") {
				return nil, nil, fmt.Errorf("ab26: missing identity axis %s", b555AdditionalFormAxis)
			}
		}
		b555AdditionalFormKey := b555AdditionalFormEnvelope.B555AdditionalFormKey
		if b555AdditionalFormKey.B555AdditionalFormFamily == "" || b555AdditionalFormKey.B555AdditionalFormVariant == "" || b555AdditionalFormKey.B555AdditionalFormDirection == "" || b555AdditionalFormKey.B555AdditionalFormField == "" || b555AdditionalFormKey.B555AdditionalFormAlternative == "" || b555AdditionalFormKey.B555AdditionalFormInstanceScope == "" {
			return nil, nil, fmt.Errorf("ab26: incomplete identity at row %d", b555AdditionalFormIndex)
		}
		for _, b555AdditionalFormPrevious := range b555AdditionalFormKeys[:b555AdditionalFormIndex] {
			if reflect.DeepEqual(b555AdditionalFormPrevious, b555AdditionalFormKey) {
				return nil, nil, fmt.Errorf("ab26: duplicate identity at row %d", b555AdditionalFormIndex)
			}
		}
		b555AdditionalFormKeys[b555AdditionalFormIndex] = b555AdditionalFormKey
	}
	return b555AdditionalFormRows, b555AdditionalFormKeys, nil
}

func B555AdditionalFormEntries() (json.RawMessage, error) {
	b555AdditionalFormData, b555AdditionalFormError := b555AdditionalFormFiles.ReadFile("data/ab26_b555_additional_forms.json")
	if b555AdditionalFormError != nil {
		return nil, b555AdditionalFormError
	}
	if _, _, b555AdditionalFormError = b555AdditionalFormValidate(b555AdditionalFormData); b555AdditionalFormError != nil {
		return nil, b555AdditionalFormError
	}
	return append(json.RawMessage(nil), b555AdditionalFormData...), nil
}

func B555AdditionalFormLookup(b555AdditionalFormKey B555AdditionalFormKey) (json.RawMessage, bool, error) {
	b555AdditionalFormData, b555AdditionalFormError := B555AdditionalFormEntries()
	if b555AdditionalFormError != nil {
		return nil, false, b555AdditionalFormError
	}
	b555AdditionalFormRows, b555AdditionalFormKeys, b555AdditionalFormError := b555AdditionalFormValidate(b555AdditionalFormData)
	if b555AdditionalFormError != nil {
		return nil, false, b555AdditionalFormError
	}
	for b555AdditionalFormIndex, b555AdditionalFormCandidate := range b555AdditionalFormKeys {
		if reflect.DeepEqual(b555AdditionalFormCandidate, b555AdditionalFormKey) {
			return append(json.RawMessage(nil), b555AdditionalFormRows[b555AdditionalFormIndex]...), true, nil
		}
	}
	return nil, false, nil
}
