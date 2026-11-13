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
		{B524CalendarViewCell{}, map[string]string{"B524CalendarViewRaw": "Ab01An01Raw", "B524CalendarViewState": "Ab01An01State"}},
		{B524CalendarViewInput{}, map[string]string{"B524CalendarViewFamily": "Ab01An01Family", "B524CalendarViewSoftware": "Ab01An01Software", "B524CalendarViewHardware": "Ab01An01Hardware", "B524CalendarViewAPIVersion": "Ab01An01APIVersion", "B524CalendarViewAPIRevision": "Ab01An01APIRevision", "B524CalendarViewOP": "Ab01An01OP", "B524CalendarViewRequestedGG": "Ab01An01RequestedGG", "B524CalendarViewII": "Ab01An01II", "B524CalendarViewRequestedRR": "Ab01An01RequestedRR", "B524CalendarViewAlternative": "Ab01An01Alternative", "B524CalendarViewReturnedGG": "Ab01An01ReturnedGG", "B524CalendarViewReturnedRR": "Ab01An01ReturnedRR", "B524CalendarViewResponseState": "Ab01An01ResponseState", "B524CalendarViewCells": "Ab01An01Cells", "B524CalendarViewRaw": "Ab01An01Raw"}},
		{B524CalendarViewComponent{}, map[string]string{"B524CalendarViewCell": "Ab01An01Cell", "B524CalendarViewState": "Ab01An01State", "B524CalendarViewCandidate": "Ab01An01Candidate"}},
		{B524CalendarViewResult{}, map[string]string{"B524CalendarViewContext": "Ab01An01Context", "B524CalendarViewComponents": "Ab01An01Components", "B524CalendarViewCompleteCandidate": "Ab01An01CompleteCandidate", "B524CalendarViewGrade": "Ab01An01Grade", "B524CalendarViewContractGrade": "Ab01An01ContractGrade", "B524CalendarViewNativeQualification": "Ab01An01NativeQualification", "B524CalendarViewCombinedTimestamp": "Ab01An01CombinedTimestamp", "B524CalendarViewEditValidation": "Ab01An01EditValidation"}},
		{DiscoveryCatalogInstanceInterval{}, map[string]string{"DiscoveryCatalogMinimum": "Ab03Minimum", "DiscoveryCatalogMaximum": "Ab03Maximum"}},
		{DiscoveryWindow{}, map[string]string{"DiscoveryCatalogFamily": "Ab03Family", "DiscoveryCatalogOpcode": "Ab03Opcode", "DiscoveryCatalogGroup": "Ab03Group", "DiscoveryCatalogInstances": "Ab03Instances", "DiscoveryCatalogRegisterCeiling": "Ab03RegisterCeiling", "DiscoveryCatalogAvailabilityOnly": "Ab03AvailabilityOnly", "DiscoveryCatalogPredicateRegister": "Ab03PredicateRegister", "DiscoveryCatalogVirtual": "Ab03Virtual"}},
		{B524CodebookViewInput{}, map[string]string{"B524CodebookViewFamily": "Ab03An03Family", "B524CodebookViewSoftware": "Ab03An03Software", "B524CodebookViewHardware": "Ab03An03Hardware", "B524CodebookViewAPIVersion": "Ab03An03APIVersion", "B524CodebookViewAPIRevision": "Ab03An03APIRevision", "B524CodebookViewOP": "Ab03An03OP", "B524CodebookViewGG": "Ab03An03GG", "B524CodebookViewII": "Ab03An03II", "B524CodebookViewRR": "Ab03An03RR", "B524CodebookViewWidth": "Ab03An03Width", "B524CodebookViewRequestedGG": "Ab03An03RequestedGG", "B524CodebookViewRequestedRR": "Ab03An03RequestedRR", "B524CodebookViewReturnedGG": "Ab03An03ReturnedGG", "B524CodebookViewReturnedRR": "Ab03An03ReturnedRR", "B524CodebookViewResponseState": "Ab03An03ResponseState", "B524CodebookViewRaw": "Ab03An03Raw"}},
		{B524CodebookViewResult{}, map[string]string{"B524CodebookViewContext": "Ab03An03Context", "B524CodebookViewInteger": "Ab03An03Integer", "B524CodebookViewCandidateLabel": "Ab03An03CandidateLabel", "B524CodebookViewApplicable": "Ab03An03Applicable", "B524CodebookViewMalformedWidth": "Ab03An03MalformedWidth", "B524CodebookViewGrade": "Ab03An03Grade", "B524CodebookViewNativeQualification": "Ab03An03NativeQualification"}},
		{B509ConstraintWordViewInput{}, map[string]string{"B509ConstraintWordViewFamily": "Ab07An07Family", "B509ConstraintWordViewVariant": "Ab07An07Variant", "B509ConstraintWordViewOperation": "Ab07An07Operation", "B509ConstraintWordViewDirection": "Ab07An07Direction", "B509ConstraintWordViewRequestedAddress": "Ab07An07RequestedAddress", "B509ConstraintWordViewReturnedAddress": "Ab07An07ReturnedAddress", "B509ConstraintWordViewMinimum": "Ab07An07Minimum", "B509ConstraintWordViewMaximum": "Ab07An07Maximum", "B509ConstraintWordViewDefault": "Ab07An07Default", "B509ConstraintWordViewStep": "Ab07An07Step", "B509ConstraintWordViewFactor": "Ab07An07Factor"}},
		{B509ConstraintWord{}, map[string]string{"B509ConstraintWordViewRaw": "Ab07An07Raw", "B509ConstraintWordViewState": "Ab07An07State", "B509ConstraintWordViewValue": "Ab07An07Value"}},
		{B509ConstraintWordView{}, map[string]string{"B509ConstraintWordViewContext": "Ab07An07Context", "B509ConstraintWordViewSelected": "Ab07An07Selected", "B509ConstraintWordViewRawAddressEqual": "Ab07An07RawAddressEqual", "B509ConstraintWordViewRequestedEditBoundsAdmitted": "Ab07An07RequestedEditBoundsAdmitted", "B509ConstraintWordViewMinimum": "Ab07An07Minimum", "B509ConstraintWordViewMaximum": "Ab07An07Maximum", "B509ConstraintWordViewDefault": "Ab07An07Default", "B509ConstraintWordViewStep": "Ab07An07Step", "B509ConstraintWordViewFactor": "Ab07An07Factor", "B509ConstraintWordViewSoftwareContract": "Ab07An07SoftwareContract", "B509ConstraintWordViewInterpretation": "Ab07An07Interpretation", "B509ConstraintWordViewNativeEchoRules": "Ab07An07NativeEchoRules", "B509ConstraintWordViewPhysicalQualification": "Ab07An07PhysicalQualification", "B509ConstraintWordViewPacketAssociation": "Ab07An07PacketAssociation", "B509ConstraintWordViewFamilyApplicability": "Ab07An07FamilyApplicability", "B509ConstraintWordViewFactorAlgorithm": "Ab07An07FactorAlgorithm", "B509ConstraintWordViewWritability": "Ab07An07Writability", "B509ConstraintWordViewPeerEnforcement": "Ab07An07PeerEnforcement", "B509ConstraintWordViewFreshness": "Ab07An07Freshness", "B509ConstraintWordViewExecution": "Ab07An07Execution"}},
		{TimerIdentity{}, map[string]string{"TimerPlanFamily": "family", "TimerPlanChannel": "channel", "TimerPlanSystemType": "system_type", "TimerPlanInstance": "instance", "TimerPlanAddress": "address", "TimerPlanDay": "day"}},
		{TimerBaseline{}, map[string]string{"TimerPlanIdentity": "identity", "TimerPlanRaw": "raw"}},
		{TimerSlot{}, map[string]string{"TimerPlanStart": "start", "TimerPlanStop": "stop"}},
		{B509QuantityViewSelection{}, map[string]string{"B509QuantityViewFamily": "family", "B509QuantityViewVariant": "variant", "B509QuantityViewOperation": "operation", "B509QuantityViewDirection": "direction"}},
		{B509QuantityViewFields{}, map[string]string{"B509QuantityViewTemperature": "temperature", "B509QuantityViewFlow": "flow", "B509QuantityViewVolume": "volume", "B509QuantityViewMinimumVolume": "minimum_volume", "B509QuantityViewMaximumVolume": "maximum_volume"}},
		{B509QuantityViewCandidate{}, map[string]string{"B509QuantityViewRaw": "raw", "B509QuantityViewState": "state", "B509QuantityViewValue": "value", "B509QuantityViewUnit": "unit", "B509QuantityViewSemanticGrade": "semantic_grade", "B509QuantityViewNativeQualification": "native_qualification"}},
		{B509QuantityView{}, map[string]string{"B509QuantityViewIdentity": "identity", "B509QuantityViewSelection": "selection", "B509QuantityViewSelected": "selected", "B509QuantityViewContractGrade": "contract_grade", "B509QuantityViewTemperature": "temperature", "B509QuantityViewFlow": "flow", "B509QuantityViewVolume": "volume", "B509QuantityViewMinimumVolume": "minimum_volume", "B509QuantityViewMaximumVolume": "maximum_volume"}},
		{ValueCodecValue{}, map[string]string{"ValueCodecRaw": "Ab11Raw", "ValueCodecGrade": "Ab11Grade", "ValueCodecInterpretation": "Ab11Interpretation", "ValueCodecQualified": "Ab11Qualified"}},
		{ValueCodecDate{}, map[string]string{"ValueCodecDay": "Ab11Day", "ValueCodecMonth": "Ab11Month", "ValueCodecYear": "Ab11Year"}},
		{ValueCodecTime{}, map[string]string{"ValueCodecHour": "Ab11Hour", "ValueCodecMinute": "Ab11Minute", "ValueCodecSecond": "Ab11Second"}},
		{ValueCodecString{}, map[string]string{"ValueCodecBytes": "Ab11Bytes", "ValueCodecLatin1": "Ab11Latin1"}},
		{StatusRequestContext{}, map[string]string{"StatusFamily": "Ab22Family", "StatusOpcode": "Ab22Opcode"}},
		{StatusStatusFields{}, map[string]string{"StatusBindingZone": "Ab22BindingZone", "SpecialFunctionStatus": "Ab22SpecialFunctionStatus", "StatusHeatingOperatingMode": "Ab22HeatingOperatingMode", "StatusCoolingOperatingMode": "Ab22CoolingOperatingMode", "StatusInfo": "Ab22StatusInfo", "StatusFrostProtection": "Ab22FrostProtection", "StatusHeatingTemperatureRaw": "Ab22HeatingTemperatureRaw", "StatusCoolingTemperatureRaw": "Ab22CoolingTemperatureRaw"}},
		{StatusStatus{}, map[string]string{"StatusContext": "Ab22Context", "StatusPayload": "Ab22Payload", "StatusFields": "Ab22Fields"}},
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
		{B524SystemScalarField[any]{}, map[string]string{"B524SystemScalarPresent": "Ab58Present", "B524SystemScalarValue": "Ab58Value"}},
		{B524SystemScalarContext{}, map[string]string{"B524SystemScalarFamily": "Ab58Family", "B524SystemScalarSoftware": "Ab58Software", "B524SystemScalarHardware": "Ab58Hardware", "B524SystemScalarAPIVersion": "Ab58APIVersion", "B524SystemScalarAPIRevision": "Ab58APIRevision", "B524SystemScalarOpcode": "Ab58Opcode", "B524SystemScalarGroup": "Ab58Group", "B524SystemScalarInstance": "Ab58Instance", "B524SystemScalarRegister": "Ab58Register"}},
		{B524SystemScalarInput{}, map[string]string{"B524SystemScalarContext": "Ab58Context", "B524SystemScalarReplyKind": "Ab58ReplyKind", "B524SystemScalarReturnedGroup": "Ab58ReturnedGroup", "B524SystemScalarReturnedRegister": "Ab58ReturnedRegister", "B524SystemScalarCodecGrade": "Ab58CodecGrade", "B524SystemScalarRaw": "Ab58Raw"}},
		{B524SystemScalar{}, map[string]string{"B524SystemScalarFloat32": "Ab58Float32", "B524SystemScalarUint16": "Ab58Uint16"}},
		{B524SystemScalarResult{}, map[string]string{"B524SystemScalarInput": "Ab58Input", "B524SystemScalarState": "Ab58State", "B524SystemScalarCodec": "Ab58Codec", "B524SystemScalarDeclaredUnit": "Ab58DeclaredUnit", "B524SystemScalarGrade": "Ab58Grade", "B524SystemScalarPresence": "Ab58Presence", "B524SystemScalarFreshness": "Ab58Freshness", "B524SystemScalarExecution": "Ab58Execution", "B524SystemScalarSupport": "Ab58Support", "B524SystemScalarEffects": "Ab58Effects", "B524SystemScalarPersistence": "Ab58Persistence", "B524SystemScalarDelivery": "Ab58Delivery", "B524SystemScalarValue": "Ab58Value"}},
		{B524DHWVentilationScalarField[any]{}, map[string]string{"B524DHWVentilationScalarPresent": "Ab60Present", "B524DHWVentilationScalarValue": "Ab60Value"}},
		{B524DHWVentilationScalarContext{}, map[string]string{"B524DHWVentilationScalarFamily": "Ab60Family", "B524DHWVentilationScalarSoftware": "Ab60Software", "B524DHWVentilationScalarHardware": "Ab60Hardware", "B524DHWVentilationScalarAPIVersion": "Ab60APIVersion", "B524DHWVentilationScalarAPIRevision": "Ab60APIRevision", "B524DHWVentilationScalarOpcode": "Ab60Opcode", "B524DHWVentilationScalarGroup": "Ab60Group", "B524DHWVentilationScalarInstance": "Ab60Instance", "B524DHWVentilationScalarRegister": "Ab60Register"}},
		{B524DHWVentilationScalarInput{}, map[string]string{"B524DHWVentilationScalarContext": "Ab60Context", "B524DHWVentilationScalarReplyKind": "Ab60ReplyKind", "B524DHWVentilationScalarReturnedGroup": "Ab60ReturnedGroup", "B524DHWVentilationScalarReturnedRegister": "Ab60ReturnedRegister", "B524DHWVentilationScalarCodecGrade": "Ab60CodecGrade", "B524DHWVentilationScalarRaw": "Ab60Raw"}},
		{B524DHWVentilationScalar{}, map[string]string{"B524DHWVentilationScalarFloat32": "Ab60Float32", "B524DHWVentilationScalarUint16": "Ab60Uint16"}},
		{B524DHWVentilationScalarResult{}, map[string]string{"B524DHWVentilationScalarInput": "Ab60Input", "B524DHWVentilationScalarState": "Ab60State", "B524DHWVentilationScalarCodec": "Ab60Codec", "B524DHWVentilationScalarDeclaredUnit": "Ab60DeclaredUnit", "B524DHWVentilationScalarGrade": "Ab60Grade", "B524DHWVentilationScalarPresence": "Ab60Presence", "B524DHWVentilationScalarFreshness": "Ab60Freshness", "B524DHWVentilationScalarExecution": "Ab60Execution", "B524DHWVentilationScalarSupport": "Ab60Support", "B524DHWVentilationScalarEffects": "Ab60Effects", "B524DHWVentilationScalarPersistence": "Ab60Persistence", "B524DHWVentilationScalarDelivery": "Ab60Delivery", "B524DHWVentilationScalarValue": "Ab60Value"}},
		{B524CircuitScalarField[any]{}, map[string]string{"B524CircuitScalarPresent": "Ab62Present", "B524CircuitScalarValue": "Ab62Value"}},
		{B524CircuitScalarContext{}, map[string]string{"B524CircuitScalarFamily": "Ab62Family", "B524CircuitScalarSoftware": "Ab62Software", "B524CircuitScalarHardware": "Ab62Hardware", "B524CircuitScalarAPIVersion": "Ab62APIVersion", "B524CircuitScalarAPIRevision": "Ab62APIRevision", "B524CircuitScalarOpcode": "Ab62Opcode", "B524CircuitScalarGroup": "Ab62Group", "B524CircuitScalarInstance": "Ab62Instance", "B524CircuitScalarRegister": "Ab62Register"}},
		{B524CircuitScalarInput{}, map[string]string{"B524CircuitScalarContext": "Ab62Context", "B524CircuitScalarReplyKind": "Ab62ReplyKind", "B524CircuitScalarReturnedGroup": "Ab62ReturnedGroup", "B524CircuitScalarReturnedRegister": "Ab62ReturnedRegister", "B524CircuitScalarCodecGrade": "Ab62CodecGrade", "B524CircuitScalarRaw": "Ab62Raw"}},
		{B524CircuitScalarResult{}, map[string]string{"B524CircuitScalarInput": "Ab62Input", "B524CircuitScalarState": "Ab62State", "B524CircuitScalarGrade": "Ab62Grade", "B524CircuitScalarDeclaredUnit": "Ab62DeclaredUnit", "B524CircuitScalarCodec": "Ab62Codec", "B524CircuitScalarValue": "Ab62Value"}},
		{B524ZoneScalarField[any]{}, map[string]string{"B524ZoneScalarPresent": "Ab64Present", "B524ZoneScalarValue": "Ab64Value"}},
		{B524ZoneScalarContext{}, map[string]string{"B524ZoneScalarFamily": "Ab64Family", "B524ZoneScalarSoftware": "Ab64Software", "B524ZoneScalarHardware": "Ab64Hardware", "B524ZoneScalarAPIVersion": "Ab64APIVersion", "B524ZoneScalarAPIRevision": "Ab64APIRevision", "B524ZoneScalarOpcode": "Ab64Opcode", "B524ZoneScalarGroup": "Ab64Group", "B524ZoneScalarInstance": "Ab64Instance", "B524ZoneScalarRegister": "Ab64Register"}},
		{B524ZoneScalarInput{}, map[string]string{"B524ZoneScalarContext": "Ab64Context", "B524ZoneScalarReplyKind": "Ab64ReplyKind", "B524ZoneScalarReturnedGroup": "Ab64ReturnedGroup", "B524ZoneScalarReturnedRegister": "Ab64ReturnedRegister", "B524ZoneScalarCodecGrade": "Ab64CodecGrade", "B524ZoneScalarRaw": "Ab64Raw"}},
		{B524ZoneScalarResult{}, map[string]string{"B524ZoneScalarInput": "Ab64Input", "B524ZoneScalarState": "Ab64State", "B524ZoneScalarValue": "Ab64Value", "B524ZoneScalarGrade": "Ab64Grade", "B524ZoneScalarUnit": "Ab64Unit", "B524ZoneScalarCodec": "Ab64Codec", "B524ZoneScalarInstanceRule": "Ab64InstanceRule", "B524ZoneScalarNativeSupport": "Ab64NativeSupport", "B524ZoneScalarQualification": "Ab64Qualification", "B524ZoneScalarPresence": "Ab64Presence", "B524ZoneScalarFreshness": "Ab64Freshness", "B524ZoneScalarExecution": "Ab64Execution", "B524ZoneScalarEffects": "Ab64Effects", "B524ZoneScalarPersistence": "Ab64Persistence", "B524ZoneScalarDelivery": "Ab64Delivery"}},
		{B524UnscopedHintB524HintKey{}, map[string]string{"B524UnscopedHintOpcode": "Ab66Opcode", "B524UnscopedHintGroup": "Ab66Group", "B524UnscopedHintRegister": "Ab66Register"}},
		{B524UnscopedHintProperties{}, map[string]string{"B524UnscopedHintType": "Ab66Type", "B524UnscopedHintMin": "Ab66Min", "B524UnscopedHintMax": "Ab66Max", "B524UnscopedHintStep": "Ab66Step"}},
		{B524UnscopedHintResult{}, map[string]string{"B524UnscopedHintKey": "Ab66Key", "B524UnscopedHintStatus": "Ab66Status", "B524UnscopedHintQualified": "Ab66Qualified", "B524UnscopedHintApplicability": "Ab66Applicability", "B524UnscopedHint": "Ab66Hint"}},
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
