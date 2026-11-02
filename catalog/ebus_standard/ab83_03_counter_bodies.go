package ebus_standard_catalog

import (
	_ "embed"
	ab83json "encoding/json"
)

type Ab83Key struct {
	Ab83Family    string  `json:"family"`
	Ab83Variant   string  `json:"variant"`
	Ab83Primary   *uint16 `json:"primary"`
	Ab83Secondary *uint8  `json:"secondary"`
	Ab83Operation *uint8  `json:"operation"`
	Ab83Direction string  `json:"direction"`
	Ab83Field     string  `json:"field"`
	Ab83Opcode    *uint8  `json:"opcode"`
	Ab83Group     *uint8  `json:"group"`
	Ab83Instance  *uint8  `json:"instance"`
	Ab83Register  *uint16 `json:"register"`
}

type Ab83Field struct {
	Ab83Identity          string            `json:"identity"`
	Ab83LogicalPosition   *int              `json:"logical_position"`
	Ab83NativePosition    *int              `json:"native_position"`
	Ab83WidthBits         *int              `json:"width_bits"`
	Ab83Signedness        *string           `json:"signedness"`
	Ab83ByteOrder         *string           `json:"byte_order"`
	Ab83Encoding          *string           `json:"encoding"`
	Ab83Conversion        *string           `json:"conversion"`
	Ab83Unit              *string           `json:"unit"`
	Ab83Default           *int              `json:"default"`
	Ab83Minimum           *int              `json:"minimum"`
	Ab83Maximum           *int              `json:"maximum"`
	Ab83Enumeration       map[string]string `json:"enumeration"`
	Ab83ReplacementValues []int             `json:"replacement_values"`
	Ab83ReservedValues    []int             `json:"reserved_values"`
	Ab83ClaimGrade        string            `json:"claim_grade"`
	Ab83RawFallback       bool              `json:"raw_fallback"`
}

type Ab83Body struct {
	Ab83Key               Ab83Key     `json:"key"`
	Ab83BodyState         string      `json:"body_state"`
	Ab83ClaimGrade        string      `json:"claim_grade"`
	Ab83PhysicalDirection string      `json:"physical_direction"`
	Ab83Callable          bool        `json:"callable"`
	Ab83RawFallback       bool        `json:"raw_fallback"`
	Ab83Fields            []Ab83Field `json:"fields"`
}

//go:embed testdata/ab83_03_counter_bodies.json
var ab83CounterJSON string

// Ab83CounterBodies returns detached finite supplemental declarations, not invocation support.
func Ab83CounterBodies() []Ab83Body {
	var ab83Bodies []Ab83Body
	if ab83Err := ab83json.Unmarshal([]byte(ab83CounterJSON), &ab83Bodies); ab83Err != nil {
		panic("ab83 counter metadata: " + ab83Err.Error())
	}
	return ab83Bodies
}

func ab83OptionalEqual[ab83T comparable](ab83Left, ab83Right *ab83T) bool {
	if ab83Left == nil || ab83Right == nil {
		return ab83Left == nil && ab83Right == nil
	}
	return *ab83Left == *ab83Right
}

func ab83EqualKey(ab83Left, ab83Right Ab83Key) bool {
	return ab83Left.Ab83Family == ab83Right.Ab83Family &&
		ab83Left.Ab83Variant == ab83Right.Ab83Variant &&
		ab83OptionalEqual(ab83Left.Ab83Primary, ab83Right.Ab83Primary) &&
		ab83OptionalEqual(ab83Left.Ab83Secondary, ab83Right.Ab83Secondary) &&
		ab83OptionalEqual(ab83Left.Ab83Operation, ab83Right.Ab83Operation) &&
		ab83Left.Ab83Direction == ab83Right.Ab83Direction &&
		ab83Left.Ab83Field == ab83Right.Ab83Field &&
		ab83OptionalEqual(ab83Left.Ab83Opcode, ab83Right.Ab83Opcode) &&
		ab83OptionalEqual(ab83Left.Ab83Group, ab83Right.Ab83Group) &&
		ab83OptionalEqual(ab83Left.Ab83Instance, ab83Right.Ab83Instance) &&
		ab83OptionalEqual(ab83Left.Ab83Register, ab83Right.Ab83Register)
}

// Ab83LookupBody requires every axis to match, including an empty field axis.
func Ab83LookupBody(ab83Key Ab83Key) (Ab83Body, bool) {
	for _, ab83Body := range Ab83CounterBodies() {
		if ab83EqualKey(ab83Body.Ab83Key, ab83Key) {
			return ab83Body, true
		}
	}
	return Ab83Body{}, false
}

// Ab83LookupField requires the complete body identity and the exact field identity.
func Ab83LookupField(ab83Key Ab83Key) (Ab83Field, bool) {
	for _, ab83Body := range Ab83CounterBodies() {
		for _, ab83Field := range ab83Body.Ab83Fields {
			ab83FieldKey := ab83Body.Ab83Key
			ab83FieldKey.Ab83Field = ab83Field.Ab83Identity
			if ab83EqualKey(ab83FieldKey, ab83Key) {
				return ab83Field, true
			}
		}
	}
	return Ab83Field{}, false
}
