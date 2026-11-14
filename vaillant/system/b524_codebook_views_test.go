package system

import (
	"encoding/json"
	"os"
	"reflect"
	"testing"
)

type b524CodebookViewVector struct {
	B524CodebookViewName    string `json:"Ab03An03Name"`
	B524CodebookViewFamily  string `json:"Ab03An03Family"`
	B524CodebookViewGG      byte   `json:"Ab03An03GG"`
	B524CodebookViewRR      uint16 `json:"Ab03An03RR"`
	B524CodebookViewWidth   string `json:"Ab03An03Width"`
	B524CodebookViewRaw     []byte `json:"Ab03An03Raw"`
	B524CodebookViewInteger int    `json:"Ab03An03Integer"`
	B524CodebookViewLabel   string `json:"Ab03An03Label"`
}

func b524CodebookViewInput() B524CodebookViewInput {
	b524CodebookViewII := byte(0)
	return B524CodebookViewInput{B524CodebookViewFamily: "controller-700", B524CodebookViewSoftware: "synthetic-software", B524CodebookViewHardware: "synthetic-hardware", B524CodebookViewAPIVersion: "synthetic-api", B524CodebookViewAPIRevision: "synthetic-revision", B524CodebookViewOP: 2, B524CodebookViewGG: 2, B524CodebookViewII: &b524CodebookViewII, B524CodebookViewRR: 2, B524CodebookViewWidth: "u16-le", B524CodebookViewRequestedGG: 2, B524CodebookViewRequestedRR: 2, B524CodebookViewReturnedGG: 2, B524CodebookViewReturnedRR: 2, B524CodebookViewResponseState: "echo-payload", B524CodebookViewRaw: []byte{5, 0}}
}

func TestAb03An03Vectors(b524CodebookViewT *testing.T) {
	b524CodebookViewData, b524CodebookViewErr := os.ReadFile("testdata/ab03an03_b524_codebook_views.json")
	if b524CodebookViewErr != nil {
		b524CodebookViewT.Fatal(b524CodebookViewErr)
	}
	var b524CodebookViewVectors []b524CodebookViewVector
	if b524CodebookViewErr = json.Unmarshal(b524CodebookViewData, &b524CodebookViewVectors); b524CodebookViewErr != nil {
		b524CodebookViewT.Fatal(b524CodebookViewErr)
	}
	if len(b524CodebookViewVectors) != 88 {
		b524CodebookViewT.Fatalf("ab03an03 vector count: %d", len(b524CodebookViewVectors))
	}
	for _, b524CodebookViewVector := range b524CodebookViewVectors {
		b524CodebookViewT.Run(b524CodebookViewVector.B524CodebookViewName, func(b524CodebookViewT *testing.T) {
			b524CodebookViewInput := b524CodebookViewInput()
			b524CodebookViewInput.B524CodebookViewFamily = b524CodebookViewVector.B524CodebookViewFamily
			b524CodebookViewInput.B524CodebookViewGG = b524CodebookViewVector.B524CodebookViewGG
			b524CodebookViewInput.B524CodebookViewRequestedGG = b524CodebookViewVector.B524CodebookViewGG
			b524CodebookViewInput.B524CodebookViewReturnedGG = b524CodebookViewVector.B524CodebookViewGG
			b524CodebookViewInput.B524CodebookViewRR = b524CodebookViewVector.B524CodebookViewRR
			b524CodebookViewInput.B524CodebookViewRequestedRR = b524CodebookViewVector.B524CodebookViewRR
			b524CodebookViewInput.B524CodebookViewReturnedRR = b524CodebookViewVector.B524CodebookViewRR
			b524CodebookViewInput.B524CodebookViewWidth = b524CodebookViewVector.B524CodebookViewWidth
			b524CodebookViewInput.B524CodebookViewRaw = b524CodebookViewVector.B524CodebookViewRaw
			b524CodebookViewResult := B524CodebookViewInterpret(b524CodebookViewInput)
			if !b524CodebookViewResult.B524CodebookViewApplicable || b524CodebookViewResult.B524CodebookViewMalformedWidth || b524CodebookViewResult.B524CodebookViewGrade != "Hypothesis" || b524CodebookViewResult.B524CodebookViewNativeQualification != "Unknown" || b524CodebookViewResult.B524CodebookViewInteger == nil || *b524CodebookViewResult.B524CodebookViewInteger != b524CodebookViewVector.B524CodebookViewInteger || b524CodebookViewResult.B524CodebookViewCandidateLabel == nil || *b524CodebookViewResult.B524CodebookViewCandidateLabel != b524CodebookViewVector.B524CodebookViewLabel {
				b524CodebookViewT.Fatalf("ab03an03 unexpected result: %+v", b524CodebookViewResult)
			}
			if !reflect.DeepEqual(b524CodebookViewResult.B524CodebookViewContext, b524CodebookViewInput) {
				b524CodebookViewT.Fatal("ab03an03 context changed")
			}
		})
	}
}

func TestAb03An03Controls(b524CodebookViewT *testing.T) {
	for _, b524CodebookViewCase := range []struct {
		b524CodebookViewName       string
		b524CodebookViewChange     func(*B524CodebookViewInput)
		b524CodebookViewInteger    *int
		b524CodebookViewApplicable bool
		b524CodebookViewMalformed  bool
	}{
		{"ab03an03-family720-pool", func(b524CodebookViewI *B524CodebookViewInput) {
			b524CodebookViewI.B524CodebookViewFamily = "controller-720"
		}, b524CodebookViewInt(5), true, false},
		{"ab03an03-family-unknown", func(b524CodebookViewI *B524CodebookViewInput) { b524CodebookViewI.B524CodebookViewFamily = "unknown" }, b524CodebookViewInt(5), false, false},
		{"ab03an03-family-missing", func(b524CodebookViewI *B524CodebookViewInput) { b524CodebookViewI.B524CodebookViewFamily = "" }, b524CodebookViewInt(5), false, false},
		{"ab03an03-op06", func(b524CodebookViewI *B524CodebookViewInput) { b524CodebookViewI.B524CodebookViewOP = 6 }, b524CodebookViewInt(5), false, false},
		{"ab03an03-gg", func(b524CodebookViewI *B524CodebookViewInput) { b524CodebookViewI.B524CodebookViewGG = 4 }, b524CodebookViewInt(5), false, false},
		{"ab03an03-rr", func(b524CodebookViewI *B524CodebookViewInput) { b524CodebookViewI.B524CodebookViewRR = 4 }, b524CodebookViewInt(5), false, false},
		{"ab03an03-request-gg", func(b524CodebookViewI *B524CodebookViewInput) { b524CodebookViewI.B524CodebookViewRequestedGG = 4 }, b524CodebookViewInt(5), false, false},
		{"ab03an03-request-rr", func(b524CodebookViewI *B524CodebookViewInput) { b524CodebookViewI.B524CodebookViewRequestedRR = 4 }, b524CodebookViewInt(5), false, false},
		{"ab03an03-return-gg", func(b524CodebookViewI *B524CodebookViewInput) { b524CodebookViewI.B524CodebookViewReturnedGG = 4 }, b524CodebookViewInt(5), false, false},
		{"ab03an03-return-rr", func(b524CodebookViewI *B524CodebookViewInput) { b524CodebookViewI.B524CodebookViewReturnedRR = 4 }, b524CodebookViewInt(5), false, false},
		{"ab03an03-ii-missing", func(b524CodebookViewI *B524CodebookViewInput) { b524CodebookViewI.B524CodebookViewII = nil }, b524CodebookViewInt(5), false, false},
		{"ab03an03-width-missing", func(b524CodebookViewI *B524CodebookViewInput) { b524CodebookViewI.B524CodebookViewWidth = "" }, nil, false, false},
		{"ab03an03-status00", func(b524CodebookViewI *B524CodebookViewInput) {
			b524CodebookViewI.B524CodebookViewResponseState = "status00"
		}, b524CodebookViewInt(5), false, false},
		{"ab03an03-statusFF", func(b524CodebookViewI *B524CodebookViewInput) {
			b524CodebookViewI.B524CodebookViewResponseState = "statusFF"
		}, b524CodebookViewInt(5), false, false},
		{"ab03an03-transport", func(b524CodebookViewI *B524CodebookViewInput) {
			b524CodebookViewI.B524CodebookViewResponseState = "transport-failure"
		}, b524CodebookViewInt(5), false, false},
		{"ab03an03-unknown-code", func(b524CodebookViewI *B524CodebookViewInput) { b524CodebookViewI.B524CodebookViewRaw = []byte{42, 0} }, b524CodebookViewInt(42), true, false},
		{"ab03an03-max-u16", func(b524CodebookViewI *B524CodebookViewInput) {
			b524CodebookViewI.B524CodebookViewRaw = []byte{255, 255}
		}, b524CodebookViewInt(65535), true, false},
		{"ab03an03-little-endian", func(b524CodebookViewI *B524CodebookViewInput) { b524CodebookViewI.B524CodebookViewRaw = []byte{0, 1} }, b524CodebookViewInt(256), true, false},
		{"ab03an03-truncated", func(b524CodebookViewI *B524CodebookViewInput) { b524CodebookViewI.B524CodebookViewRaw = []byte{5} }, nil, false, true},
		{"ab03an03-extra", func(b524CodebookViewI *B524CodebookViewInput) {
			b524CodebookViewI.B524CodebookViewRaw = []byte{5, 0, 0}
		}, nil, false, true},
		{"ab03an03-empty", func(b524CodebookViewI *B524CodebookViewInput) { b524CodebookViewI.B524CodebookViewRaw = nil }, nil, false, true},
	} {
		b524CodebookViewT.Run(b524CodebookViewCase.b524CodebookViewName, func(b524CodebookViewT *testing.T) {
			b524CodebookViewInput := b524CodebookViewInput()
			b524CodebookViewCase.b524CodebookViewChange(&b524CodebookViewInput)
			b524CodebookViewResult := B524CodebookViewInterpret(b524CodebookViewInput)
			b524CodebookViewGrade := "Unknown"
			if b524CodebookViewCase.b524CodebookViewApplicable {
				b524CodebookViewGrade = "Hypothesis"
			}
			if !reflect.DeepEqual(b524CodebookViewResult.B524CodebookViewInteger, b524CodebookViewCase.b524CodebookViewInteger) || b524CodebookViewResult.B524CodebookViewApplicable != b524CodebookViewCase.b524CodebookViewApplicable || b524CodebookViewResult.B524CodebookViewMalformedWidth != b524CodebookViewCase.b524CodebookViewMalformed || b524CodebookViewResult.B524CodebookViewCandidateLabel != nil || b524CodebookViewResult.B524CodebookViewGrade != b524CodebookViewGrade || b524CodebookViewResult.B524CodebookViewNativeQualification != "Unknown" || !reflect.DeepEqual(b524CodebookViewResult.B524CodebookViewContext, b524CodebookViewInput) {
				b524CodebookViewT.Fatalf("ab03an03 result: %+v", b524CodebookViewResult)
			}
		})
	}
}

func b524CodebookViewInt(b524CodebookViewValue int) *int { return &b524CodebookViewValue }

func TestAb03An03Alternatives(b524CodebookViewT *testing.T) {
	for _, b524CodebookViewCase := range []struct {
		b524CodebookViewName       string
		b524CodebookViewFamily     string
		b524CodebookViewRR         uint16
		b524CodebookViewWidth      string
		b524CodebookViewRaw        []byte
		b524CodebookViewInteger    *int
		b524CodebookViewApplicable bool
		b524CodebookViewMalformed  bool
	}{
		{"ab03an03-720-remote-eight", "controller-720", 19, "u16-le", []byte{9, 0}, b524CodebookViewInt(9), true, false},
		{"ab03an03-720-party", "controller-720", 27, "u8", []byte{3}, b524CodebookViewInt(3), true, false},
		{"ab03an03-no-width-inference", "controller-700", 27, "u16-le", []byte{3}, nil, false, true},
		{"ab03an03-u8-extra", "controller-700", 27, "u8", []byte{3, 0}, nil, false, true},
		{"ab03an03-u8-max", "controller-700", 27, "u8", []byte{255}, b524CodebookViewInt(255), true, false},
		{"ab03an03-signed-max", "controller-700", 28, "signed8", []byte{127}, b524CodebookViewInt(127), true, false},
		{"ab03an03-signed-min", "controller-700", 28, "signed8", []byte{128}, b524CodebookViewInt(-128), true, false},
		{"ab03an03-no-unsigned255", "controller-720", 28, "u8", []byte{255}, b524CodebookViewInt(255), false, false},
		{"ab03an03-signed-extra", "controller-720", 28, "signed8", []byte{255, 0}, nil, false, true},
		{"ab03an03-unknown-selector", "controller-700", 99, "u16-le", []byte{0, 0}, b524CodebookViewInt(0), false, false},
	} {
		b524CodebookViewT.Run(b524CodebookViewCase.b524CodebookViewName, func(b524CodebookViewT *testing.T) {
			b524CodebookViewInput := b524CodebookViewInput()
			b524CodebookViewInput.B524CodebookViewFamily = b524CodebookViewCase.b524CodebookViewFamily
			b524CodebookViewInput.B524CodebookViewGG = 3
			b524CodebookViewInput.B524CodebookViewRequestedGG = 3
			b524CodebookViewInput.B524CodebookViewReturnedGG = 3
			b524CodebookViewInput.B524CodebookViewRR = b524CodebookViewCase.b524CodebookViewRR
			b524CodebookViewInput.B524CodebookViewRequestedRR = b524CodebookViewCase.b524CodebookViewRR
			b524CodebookViewInput.B524CodebookViewReturnedRR = b524CodebookViewCase.b524CodebookViewRR
			b524CodebookViewInput.B524CodebookViewWidth = b524CodebookViewCase.b524CodebookViewWidth
			b524CodebookViewInput.B524CodebookViewRaw = b524CodebookViewCase.b524CodebookViewRaw
			b524CodebookViewResult := B524CodebookViewInterpret(b524CodebookViewInput)
			if !reflect.DeepEqual(b524CodebookViewResult.B524CodebookViewInteger, b524CodebookViewCase.b524CodebookViewInteger) || b524CodebookViewResult.B524CodebookViewApplicable != b524CodebookViewCase.b524CodebookViewApplicable || b524CodebookViewResult.B524CodebookViewMalformedWidth != b524CodebookViewCase.b524CodebookViewMalformed || b524CodebookViewResult.B524CodebookViewCandidateLabel != nil {
				b524CodebookViewT.Fatalf("ab03an03 result: %+v", b524CodebookViewResult)
			}
		})
	}
}

func TestAb03An03ContextOwnership(b524CodebookViewT *testing.T) {
	for _, b524CodebookViewField := range []string{"B524CodebookViewSoftware", "B524CodebookViewHardware", "B524CodebookViewAPIVersion", "B524CodebookViewAPIRevision"} {
		b524CodebookViewT.Run("ab03an03-"+b524CodebookViewField, func(b524CodebookViewT *testing.T) {
			b524CodebookViewInput := b524CodebookViewInput()
			reflect.ValueOf(&b524CodebookViewInput).Elem().FieldByName(b524CodebookViewField).SetString("synthetic-alternative")
			b524CodebookViewResult := B524CodebookViewInterpret(b524CodebookViewInput)
			if !reflect.DeepEqual(b524CodebookViewResult.B524CodebookViewContext, b524CodebookViewInput) || !b524CodebookViewResult.B524CodebookViewApplicable || b524CodebookViewResult.B524CodebookViewCandidateLabel == nil || *b524CodebookViewResult.B524CodebookViewCandidateLabel != "pool" {
				b524CodebookViewT.Fatal("ab03an03 independent identity lost or used as inferred selector")
			}
		})
	}
	for _, b524CodebookViewII := range []byte{1, 255} {
		b524CodebookViewInput := b524CodebookViewInput()
		b524CodebookViewInput.B524CodebookViewII = &b524CodebookViewII
		if b524CodebookViewResult := B524CodebookViewInterpret(b524CodebookViewInput); !b524CodebookViewResult.B524CodebookViewApplicable || *b524CodebookViewResult.B524CodebookViewContext.B524CodebookViewII != b524CodebookViewII {
			b524CodebookViewT.Fatal("ab03an03 explicit dynamic II lost")
		}
		b524CodebookViewInput.B524CodebookViewGG = 1
		b524CodebookViewInput.B524CodebookViewRequestedGG = 1
		b524CodebookViewInput.B524CodebookViewReturnedGG = 1
		b524CodebookViewInput.B524CodebookViewRR = 3
		b524CodebookViewInput.B524CodebookViewRequestedRR = 3
		b524CodebookViewInput.B524CodebookViewReturnedRR = 3
		if b524CodebookViewResult := B524CodebookViewInterpret(b524CodebookViewInput); b524CodebookViewResult.B524CodebookViewApplicable {
			b524CodebookViewT.Fatal("ab03an03 fixed II00 broadened")
		}
	}
	for _, b524CodebookViewState := range []string{"echo-payload", "transport-failure"} {
		b524CodebookViewInput := b524CodebookViewInput()
		b524CodebookViewInput.B524CodebookViewResponseState = b524CodebookViewState
		b524CodebookViewFirst := B524CodebookViewInterpret(b524CodebookViewInput)
		b524CodebookViewSecond := B524CodebookViewInterpret(b524CodebookViewInput)
		b524CodebookViewInput.B524CodebookViewRaw[0] = 99
		*b524CodebookViewInput.B524CodebookViewII = 99
		if b524CodebookViewFirst.B524CodebookViewContext.B524CodebookViewRaw[0] != 5 || *b524CodebookViewFirst.B524CodebookViewContext.B524CodebookViewII != 0 {
			b524CodebookViewT.Fatal("ab03an03 input aliases output")
		}
		b524CodebookViewFirst.B524CodebookViewContext.B524CodebookViewRaw[0] = 88
		*b524CodebookViewFirst.B524CodebookViewContext.B524CodebookViewII = 88
		*b524CodebookViewFirst.B524CodebookViewInteger = 88
		if b524CodebookViewFirst.B524CodebookViewCandidateLabel != nil {
			*b524CodebookViewFirst.B524CodebookViewCandidateLabel = "ab03an03-mutated"
		}
		if b524CodebookViewInput.B524CodebookViewRaw[0] != 99 || *b524CodebookViewInput.B524CodebookViewII != 99 {
			b524CodebookViewT.Fatal("ab03an03 output aliases input")
		}
		if b524CodebookViewSecond.B524CodebookViewContext.B524CodebookViewRaw[0] != 5 || *b524CodebookViewSecond.B524CodebookViewContext.B524CodebookViewII != 0 || *b524CodebookViewSecond.B524CodebookViewInteger != 5 {
			b524CodebookViewT.Fatal("ab03an03 outputs alias")
		}
		if b524CodebookViewSecond.B524CodebookViewCandidateLabel != nil && *b524CodebookViewSecond.B524CodebookViewCandidateLabel != "pool" {
			b524CodebookViewT.Fatal("ab03an03 label aliases")
		}
	}
}

func TestAb03An03CandidateSurface(b524CodebookViewT *testing.T) {
	b524CodebookViewType := reflect.TypeOf(B524CodebookViewResult{})
	b524CodebookViewFields := []string{"B524CodebookViewContext", "B524CodebookViewInteger", "B524CodebookViewCandidateLabel", "B524CodebookViewApplicable", "B524CodebookViewMalformedWidth", "B524CodebookViewGrade", "B524CodebookViewNativeQualification"}
	if b524CodebookViewType.NumField() != len(b524CodebookViewFields) {
		b524CodebookViewT.Fatal("ab03an03 candidate surface gained an unsupported policy or permission field")
	}
	for b524CodebookViewIndex, b524CodebookViewField := range b524CodebookViewFields {
		if b524CodebookViewType.Field(b524CodebookViewIndex).Name != b524CodebookViewField {
			b524CodebookViewT.Fatal("ab03an03 candidate surface changed")
		}
	}
}

func TestRenamedFieldJSONKeys(t *testing.T) {
	cases := []struct {
		value any
		keys  map[string]string
	}{
		{B524CalendarViewCell{}, map[string]string{"B524CalendarViewRaw": "raw", "B524CalendarViewState": "state"}},
		{B524CalendarViewInput{}, map[string]string{"B524CalendarViewFamily": "family", "B524CalendarViewSoftware": "software", "B524CalendarViewHardware": "hardware", "B524CalendarViewAPIVersion": "api_version", "B524CalendarViewAPIRevision": "api_revision", "B524CalendarViewOP": "op", "B524CalendarViewRequestedGG": "requested_gg", "B524CalendarViewII": "ii", "B524CalendarViewRequestedRR": "requested_rr", "B524CalendarViewAlternative": "alternative", "B524CalendarViewReturnedGG": "returned_gg", "B524CalendarViewReturnedRR": "returned_rr", "B524CalendarViewResponseState": "response_state", "B524CalendarViewCells": "cells", "B524CalendarViewRaw": "raw"}},
		{B524CalendarViewComponent{}, map[string]string{"B524CalendarViewCell": "cell", "B524CalendarViewState": "state", "B524CalendarViewCandidate": "candidate"}},
		{B524CalendarViewResult{}, map[string]string{"B524CalendarViewContext": "context", "B524CalendarViewComponents": "components", "B524CalendarViewCompleteCandidate": "complete_candidate", "B524CalendarViewGrade": "grade", "B524CalendarViewContractGrade": "contract_grade", "B524CalendarViewNativeQualification": "native_qualification", "B524CalendarViewCombinedTimestamp": "combined_timestamp", "B524CalendarViewEditValidation": "edit_validation"}},
		{DiscoveryCatalogInstanceInterval{}, map[string]string{"DiscoveryCatalogMinimum": "minimum", "DiscoveryCatalogMaximum": "maximum"}},
		{DiscoveryWindow{}, map[string]string{"DiscoveryCatalogFamily": "family", "DiscoveryCatalogOpcode": "opcode", "DiscoveryCatalogGroup": "group", "DiscoveryCatalogInstances": "instances", "DiscoveryCatalogRegisterCeiling": "register_ceiling", "DiscoveryCatalogAvailabilityOnly": "availability_only", "DiscoveryCatalogPredicateRegister": "predicate_register", "DiscoveryCatalogVirtual": "virtual"}},
		{B524CodebookViewInput{}, map[string]string{"B524CodebookViewFamily": "family", "B524CodebookViewSoftware": "software", "B524CodebookViewHardware": "hardware", "B524CodebookViewAPIVersion": "api_version", "B524CodebookViewAPIRevision": "api_revision", "B524CodebookViewOP": "op", "B524CodebookViewGG": "gg", "B524CodebookViewII": "ii", "B524CodebookViewRR": "rr", "B524CodebookViewWidth": "width", "B524CodebookViewRequestedGG": "requested_gg", "B524CodebookViewRequestedRR": "requested_rr", "B524CodebookViewReturnedGG": "returned_gg", "B524CodebookViewReturnedRR": "returned_rr", "B524CodebookViewResponseState": "response_state", "B524CodebookViewRaw": "raw"}},
		{B524CodebookViewResult{}, map[string]string{"B524CodebookViewContext": "context", "B524CodebookViewInteger": "integer", "B524CodebookViewCandidateLabel": "candidate_label", "B524CodebookViewApplicable": "applicable", "B524CodebookViewMalformedWidth": "malformed_width", "B524CodebookViewGrade": "grade", "B524CodebookViewNativeQualification": "native_qualification"}},
		{B509ConstraintWordViewInput{}, map[string]string{"B509ConstraintWordViewFamily": "family", "B509ConstraintWordViewVariant": "variant", "B509ConstraintWordViewOperation": "operation", "B509ConstraintWordViewDirection": "direction", "B509ConstraintWordViewRequestedAddress": "requested_address", "B509ConstraintWordViewReturnedAddress": "returned_address", "B509ConstraintWordViewMinimum": "minimum", "B509ConstraintWordViewMaximum": "maximum", "B509ConstraintWordViewDefault": "default", "B509ConstraintWordViewStep": "step", "B509ConstraintWordViewFactor": "factor"}},
		{B509ConstraintWord{}, map[string]string{"B509ConstraintWordViewRaw": "raw", "B509ConstraintWordViewState": "state", "B509ConstraintWordViewValue": "value"}},
		{B509ConstraintWordView{}, map[string]string{"B509ConstraintWordViewContext": "context", "B509ConstraintWordViewSelected": "selected", "B509ConstraintWordViewRawAddressEqual": "raw_address_equal", "B509ConstraintWordViewRequestedEditBoundsAdmitted": "requested_edit_bounds_admitted", "B509ConstraintWordViewMinimum": "minimum", "B509ConstraintWordViewMaximum": "maximum", "B509ConstraintWordViewDefault": "default", "B509ConstraintWordViewStep": "step", "B509ConstraintWordViewFactor": "factor", "B509ConstraintWordViewSoftwareContract": "software_contract", "B509ConstraintWordViewInterpretation": "interpretation", "B509ConstraintWordViewNativeEchoRules": "native_echo_rules", "B509ConstraintWordViewPhysicalQualification": "physical_qualification", "B509ConstraintWordViewPacketAssociation": "packet_association", "B509ConstraintWordViewFamilyApplicability": "family_applicability", "B509ConstraintWordViewFactorAlgorithm": "factor_algorithm", "B509ConstraintWordViewWritability": "writability", "B509ConstraintWordViewPeerEnforcement": "peer_enforcement", "B509ConstraintWordViewFreshness": "freshness", "B509ConstraintWordViewExecution": "execution"}},
		{TimerIdentity{}, map[string]string{"TimerPlanFamily": "family", "TimerPlanChannel": "channel", "TimerPlanSystemType": "system_type", "TimerPlanInstance": "instance", "TimerPlanAddress": "address", "TimerPlanDay": "day"}},
		{TimerBaseline{}, map[string]string{"TimerPlanIdentity": "identity", "TimerPlanRaw": "raw"}},
		{TimerSlot{}, map[string]string{"TimerPlanStart": "start", "TimerPlanStop": "stop"}},
		{B509QuantityViewSelection{}, map[string]string{"B509QuantityViewFamily": "family", "B509QuantityViewVariant": "variant", "B509QuantityViewOperation": "operation", "B509QuantityViewDirection": "direction"}},
		{B509QuantityViewFields{}, map[string]string{"B509QuantityViewTemperature": "temperature", "B509QuantityViewFlow": "flow", "B509QuantityViewVolume": "volume", "B509QuantityViewMinimumVolume": "minimum_volume", "B509QuantityViewMaximumVolume": "maximum_volume"}},
		{B509QuantityViewCandidate{}, map[string]string{"B509QuantityViewRaw": "raw", "B509QuantityViewState": "state", "B509QuantityViewValue": "value", "B509QuantityViewUnit": "unit", "B509QuantityViewSemanticGrade": "semantic_grade", "B509QuantityViewNativeQualification": "native_qualification"}},
		{B509QuantityView{}, map[string]string{"B509QuantityViewIdentity": "identity", "B509QuantityViewSelection": "selection", "B509QuantityViewSelected": "selected", "B509QuantityViewContractGrade": "contract_grade", "B509QuantityViewTemperature": "temperature", "B509QuantityViewFlow": "flow", "B509QuantityViewVolume": "volume", "B509QuantityViewMinimumVolume": "minimum_volume", "B509QuantityViewMaximumVolume": "maximum_volume"}},
		{ValueCodecValue{}, map[string]string{"ValueCodecRaw": "raw", "ValueCodecGrade": "grade", "ValueCodecInterpretation": "interpretation", "ValueCodecQualified": "qualified"}},
		{ValueCodecDate{}, map[string]string{"ValueCodecDay": "day", "ValueCodecMonth": "month", "ValueCodecYear": "year"}},
		{ValueCodecTime{}, map[string]string{"ValueCodecHour": "hour", "ValueCodecMinute": "minute", "ValueCodecSecond": "second"}},
		{ValueCodecString{}, map[string]string{"ValueCodecBytes": "bytes", "ValueCodecLatin1": "latin1"}},
		{StatusRequestContext{}, map[string]string{"StatusFamily": "family", "StatusOpcode": "opcode"}},
		{StatusStatusFields{}, map[string]string{"StatusBindingZone": "binding_zone", "SpecialFunctionStatus": "special_function_status", "StatusHeatingOperatingMode": "heating_operating_mode", "StatusCoolingOperatingMode": "cooling_operating_mode", "StatusInfo": "status_info", "StatusFrostProtection": "frost_protection", "StatusHeatingTemperatureRaw": "heating_temperature_raw", "StatusCoolingTemperatureRaw": "cooling_temperature_raw"}},
		{StatusStatus{}, map[string]string{"StatusContext": "context", "StatusPayload": "payload", "StatusFields": "fields"}},
		{EventIdentity{}, map[string]string{"EventPlanFamily": "family", "EventPlanSystemType": "system_type", "EventPlanInstance": "instance", "EventPlanAddress": "address", "EventPlanCode": "code"}},
		{EventBaseline{}, map[string]string{"EventPlanIdentity": "identity", "EventPlanOpcode": "opcode", "EventPlanSuccess": "success", "EventPlanRaw": "raw"}},
		{EventEdit{}, map[string]string{"EventPlanIdentity": "identity", "EventPlanOpcode": "opcode", "EventPlanValues": "values", "EventPlanOP09": "op09", "EventPlanOP0B": "op0b"}},
		{EventPlan{}, map[string]string{"EventPlanIdentity": "identity", "EventPlanPayload": "payload", "EventPlanValues": "values", "EventPlanOP09": "op09", "EventPlanOP0B": "op0b"}},
		{DescriptionKey{}, map[string]string{"DescriptionIndexFamily": "family", "DescriptionIndexSoftwareIdentity": "software_identity", "DescriptionIndexHardwareIdentity": "hardware_identity", "DescriptionIndexAPIVersion": "api_version", "DescriptionIndexAPIRevision": "api_revision", "DescriptionIndexOpcode": "opcode", "DescriptionIndexGroup": "group", "DescriptionIndexRequestedInstance": "requested_instance", "DescriptionIndexRegister": "register"}},
		{DescriptionRecord{}, map[string]string{"DescriptionIndexKey": "key", "DescriptionIndexState": "state", "RawDescription": "raw_description", "DescriptionIndexCodec": "codec", "DescriptionIndexWidth": "width"}},
		{B524SystemInformationInformationContext{}, map[string]string{"B524SystemInformationFamily": "family", "B524SystemInformationSoftware": "software", "B524SystemInformationHardware": "hardware", "B524SystemInformationAPIVersion": "api_version", "B524SystemInformationAPIRevision": "api_revision", "B524SystemInformationIdentifier": "identifier"}},
		{B524SystemInformationResult{}, map[string]string{"B524SystemInformationContext": "context", "B524SystemInformationRaw": "raw", "B524SystemInformationValid": "valid", "B524SystemInformationInvalidReason": "invalid_reason", "B524SystemInformationValue": "value", "B524SystemInformationMeaning": "meaning", "B524SystemInformationCapacity": "capacity", "B524SystemInformationHeatPumpCount": "heat_pump_count", "B524SystemInformationRepresentationGrade": "representation_grade", "B524SystemInformationLayoutGrade": "layout_grade", "B524SystemInformationQualificationGrade": "qualification_grade", "B524SystemInformationPresenceGrade": "presence_grade", "B524SystemInformationEligibilityGrade": "eligibility_grade"}},
		{B524DescriptionSpanContext{}, map[string]string{"B524DescriptionSpanFamily": "family", "B524DescriptionSpanSoftware": "software", "B524DescriptionSpanHardware": "hardware", "B524DescriptionSpanAPIVersion": "api_version", "B524DescriptionSpanAPIRevision": "api_revision", "B524DescriptionSpanOP": "op", "B524DescriptionSpanGG": "gg", "B524DescriptionSpanII": "ii", "B524DescriptionSpanRR": "rr"}},
		{B524DescriptionSpanLimits{}, map[string]string{"B524DescriptionSpanMinimum": "minimum", "B524DescriptionSpanMaximum": "maximum", "B524DescriptionSpanStep": "step"}},
		{B524Description{}, map[string]string{"B524DescriptionSpanContext": "context", "B524DescriptionSpanCodec": "codec", "B524DescriptionSpanWidth": "width", "B524DescriptionSpanGrade": "grade", "B524DescriptionSpanRaw": "raw", "B524DescriptionSpanMinimumRaw": "minimum_raw", "B524DescriptionSpanMaximumRaw": "maximum_raw", "B524DescriptionSpanStepRaw": "step_raw", "B524DescriptionSpanMetadataOP": "metadata_op", "B524DescriptionSpanLimits": "limits", "B524DescriptionSpanImplementationValid": "implementation_valid", "B524DescriptionSpanNumericValid": "numeric_valid", "B524DescriptionSpanStepMeaning": "step_meaning", "B524DescriptionSpanReason": "reason"}},
		{B524SystemScalarField[any]{}, map[string]string{"B524SystemScalarPresent": "present", "B524SystemScalarValue": "value"}},
		{B524SystemScalarContext{}, map[string]string{"B524SystemScalarFamily": "family", "B524SystemScalarSoftware": "software", "B524SystemScalarHardware": "hardware", "B524SystemScalarAPIVersion": "api_version", "B524SystemScalarAPIRevision": "api_revision", "B524SystemScalarOpcode": "opcode", "B524SystemScalarGroup": "group", "B524SystemScalarInstance": "instance", "B524SystemScalarRegister": "register"}},
		{B524SystemScalarInput{}, map[string]string{"B524SystemScalarContext": "context", "B524SystemScalarReplyKind": "reply_kind", "B524SystemScalarReturnedGroup": "returned_group", "B524SystemScalarReturnedRegister": "returned_register", "B524SystemScalarCodecGrade": "codec_grade", "B524SystemScalarRaw": "raw"}},
		{B524SystemScalar{}, map[string]string{"B524SystemScalarFloat32": "float32", "B524SystemScalarUint16": "uint16"}},
		{B524SystemScalarResult{}, map[string]string{"B524SystemScalarInput": "input", "B524SystemScalarState": "state", "B524SystemScalarCodec": "codec", "B524SystemScalarDeclaredUnit": "declared_unit", "B524SystemScalarGrade": "grade", "B524SystemScalarPresence": "presence", "B524SystemScalarFreshness": "freshness", "B524SystemScalarExecution": "execution", "B524SystemScalarSupport": "support", "B524SystemScalarEffects": "effects", "B524SystemScalarPersistence": "persistence", "B524SystemScalarDelivery": "delivery", "B524SystemScalarValue": "value"}},
		{B524DHWVentilationScalarField[any]{}, map[string]string{"B524DHWVentilationScalarPresent": "present", "B524DHWVentilationScalarValue": "value"}},
		{B524DHWVentilationScalarContext{}, map[string]string{"B524DHWVentilationScalarFamily": "family", "B524DHWVentilationScalarSoftware": "software", "B524DHWVentilationScalarHardware": "hardware", "B524DHWVentilationScalarAPIVersion": "api_version", "B524DHWVentilationScalarAPIRevision": "api_revision", "B524DHWVentilationScalarOpcode": "opcode", "B524DHWVentilationScalarGroup": "group", "B524DHWVentilationScalarInstance": "instance", "B524DHWVentilationScalarRegister": "register"}},
		{B524DHWVentilationScalarInput{}, map[string]string{"B524DHWVentilationScalarContext": "context", "B524DHWVentilationScalarReplyKind": "reply_kind", "B524DHWVentilationScalarReturnedGroup": "returned_group", "B524DHWVentilationScalarReturnedRegister": "returned_register", "B524DHWVentilationScalarCodecGrade": "codec_grade", "B524DHWVentilationScalarRaw": "raw"}},
		{B524DHWVentilationScalar{}, map[string]string{"B524DHWVentilationScalarFloat32": "float32", "B524DHWVentilationScalarUint16": "uint16"}},
		{B524DHWVentilationScalarResult{}, map[string]string{"B524DHWVentilationScalarInput": "input", "B524DHWVentilationScalarState": "state", "B524DHWVentilationScalarCodec": "codec", "B524DHWVentilationScalarDeclaredUnit": "declared_unit", "B524DHWVentilationScalarGrade": "grade", "B524DHWVentilationScalarPresence": "presence", "B524DHWVentilationScalarFreshness": "freshness", "B524DHWVentilationScalarExecution": "execution", "B524DHWVentilationScalarSupport": "support", "B524DHWVentilationScalarEffects": "effects", "B524DHWVentilationScalarPersistence": "persistence", "B524DHWVentilationScalarDelivery": "delivery", "B524DHWVentilationScalarValue": "value"}},
		{B524CircuitScalarField[any]{}, map[string]string{"B524CircuitScalarPresent": "present", "B524CircuitScalarValue": "value"}},
		{B524CircuitScalarContext{}, map[string]string{"B524CircuitScalarFamily": "family", "B524CircuitScalarSoftware": "software", "B524CircuitScalarHardware": "hardware", "B524CircuitScalarAPIVersion": "api_version", "B524CircuitScalarAPIRevision": "api_revision", "B524CircuitScalarOpcode": "opcode", "B524CircuitScalarGroup": "group", "B524CircuitScalarInstance": "instance", "B524CircuitScalarRegister": "register"}},
		{B524CircuitScalarInput{}, map[string]string{"B524CircuitScalarContext": "context", "B524CircuitScalarReplyKind": "reply_kind", "B524CircuitScalarReturnedGroup": "returned_group", "B524CircuitScalarReturnedRegister": "returned_register", "B524CircuitScalarCodecGrade": "codec_grade", "B524CircuitScalarRaw": "raw"}},
		{B524CircuitScalarResult{}, map[string]string{"B524CircuitScalarInput": "input", "B524CircuitScalarState": "state", "B524CircuitScalarGrade": "grade", "B524CircuitScalarDeclaredUnit": "declared_unit", "B524CircuitScalarCodec": "codec", "B524CircuitScalarValue": "value"}},
		{B524ZoneScalarField[any]{}, map[string]string{"B524ZoneScalarPresent": "present", "B524ZoneScalarValue": "value"}},
		{B524ZoneScalarContext{}, map[string]string{"B524ZoneScalarFamily": "family", "B524ZoneScalarSoftware": "software", "B524ZoneScalarHardware": "hardware", "B524ZoneScalarAPIVersion": "api_version", "B524ZoneScalarAPIRevision": "api_revision", "B524ZoneScalarOpcode": "opcode", "B524ZoneScalarGroup": "group", "B524ZoneScalarInstance": "instance", "B524ZoneScalarRegister": "register"}},
		{B524ZoneScalarInput{}, map[string]string{"B524ZoneScalarContext": "context", "B524ZoneScalarReplyKind": "reply_kind", "B524ZoneScalarReturnedGroup": "returned_group", "B524ZoneScalarReturnedRegister": "returned_register", "B524ZoneScalarCodecGrade": "codec_grade", "B524ZoneScalarRaw": "raw"}},
		{B524ZoneScalarResult{}, map[string]string{"B524ZoneScalarInput": "input", "B524ZoneScalarState": "state", "B524ZoneScalarValue": "value", "B524ZoneScalarGrade": "grade", "B524ZoneScalarUnit": "unit", "B524ZoneScalarCodec": "codec", "B524ZoneScalarInstanceRule": "instance_rule", "B524ZoneScalarNativeSupport": "native_support", "B524ZoneScalarQualification": "qualification", "B524ZoneScalarPresence": "presence", "B524ZoneScalarFreshness": "freshness", "B524ZoneScalarExecution": "execution", "B524ZoneScalarEffects": "effects", "B524ZoneScalarPersistence": "persistence", "B524ZoneScalarDelivery": "delivery"}},
		{B524UnscopedHintB524HintKey{}, map[string]string{"B524UnscopedHintOpcode": "opcode", "B524UnscopedHintGroup": "group", "B524UnscopedHintRegister": "register"}},
		{B524UnscopedHintProperties{}, map[string]string{"B524UnscopedHintType": "type", "B524UnscopedHintMin": "min", "B524UnscopedHintMax": "max", "B524UnscopedHintStep": "step"}},
		{B524UnscopedHintResult{}, map[string]string{"B524UnscopedHintKey": "key", "B524UnscopedHintStatus": "status", "B524UnscopedHintQualified": "qualified", "B524UnscopedHintApplicability": "applicability", "B524UnscopedHint": "hint"}},
	}
	for _, tc := range cases {
		data, err := json.Marshal(tc.value)
		if err != nil {
			t.Fatal(err)
		}
		var object map[string]json.RawMessage
		if err := json.Unmarshal(data, &object); err != nil {
			t.Fatal(err)
		}
		typ := reflect.TypeOf(tc.value)
		for field, key := range tc.keys {
			if _, ok := object[key]; !ok {
				t.Errorf("%s.%s missing JSON key %s", typ, field, key)
			}
			if _, ok := object[field]; ok {
				t.Errorf("%s exposes renamed key %s", typ, field)
			}
			if member, ok := typ.FieldByName(field); !ok || member.Tag.Get("json") != key {
				t.Errorf("%s.%s does not decode JSON key %s", typ, field, key)
			}
		}
	}
}
