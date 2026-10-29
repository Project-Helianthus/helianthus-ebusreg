package vaillant

import (
	_ "embed"
	ab53JSON "encoding/json"
	ab53Reflect "reflect"
)

type Ab53Key struct {
	Ab53Family    string `json:"family"`
	Ab53Variant   string `json:"variant"`
	Ab53Primary   uint8  `json:"primary"`
	Ab53Secondary *uint8 `json:"secondary"`
	Ab53Operation *uint8 `json:"operation"`
	Ab53Direction string `json:"direction"`
	Ab53Field     string `json:"field"`
}

type Ab53Domain struct {
	Ab53Values  []int `json:"values"`
	Ab53Minimum *int  `json:"minimum"`
	Ab53Maximum *int  `json:"maximum"`
	Ab53Step    *int  `json:"step"`
}

type Ab53Entry struct {
	Ab53Key           Ab53Key        `json:"key"`
	Ab53Grade         Ab05Grade      `json:"grade"`
	Ab53Position      *int           `json:"position"`
	Ab53Width         *int           `json:"width_bytes"`
	Ab53Signedness    string         `json:"signedness"`
	Ab53ByteOrder     string         `json:"byte_order"`
	Ab53Divisor       *int           `json:"divisor"`
	Ab53Unit          *string        `json:"unit"`
	Ab53Labels        map[int]string `json:"labels"`
	Ab53PositiveLabel *string        `json:"positive_label"`
	Ab53Domain        *Ab53Domain    `json:"domain"`
	Ab53Notes         []string       `json:"notes"`
}

type Ab53Catalog struct {
	Ab53MetadataPresence   Ab05Grade   `json:"metadata_presence"`
	Ab53PhysicalAdaptation Ab05Grade   `json:"physical_adaptation"`
	Ab53InstalledSupport   Ab05Grade   `json:"installed_support"`
	Ab53PollingPermission  Ab05Grade   `json:"polling_permission"`
	Ab53RoutableInvocation Ab05Grade   `json:"routable_invocation"`
	Ab53Limitations        []string    `json:"limitations"`
	Ab53Entries            []Ab53Entry `json:"entries"`
}

//go:embed data/ab53_b504_candidates.json
var ab53B504JSON []byte

func Ab53B504Candidates() (Ab53Catalog, error) {
	var ab53Catalog Ab53Catalog
	ab53Error := ab53JSON.Unmarshal(ab53B504JSON, &ab53Catalog)
	return ab53Catalog, ab53Error
}

func Ab53LookupB504Candidate(ab53Key Ab53Key) (Ab53Entry, bool, error) {
	ab53Catalog, ab53Error := Ab53B504Candidates()
	if ab53Error != nil {
		return Ab53Entry{}, false, ab53Error
	}
	for _, ab53Entry := range ab53Catalog.Ab53Entries {
		if ab53Reflect.DeepEqual(ab53Entry.Ab53Key, ab53Key) {
			return ab53Entry, true, nil
		}
	}
	return Ab53Entry{}, false, nil
}
