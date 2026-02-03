package system

import (
	"fmt"

	"github.com/d3vi1/helianthus-ebusgo/protocol"
	"github.com/d3vi1/helianthus-ebusgo/types"
	"github.com/d3vi1/helianthus-ebusreg/schema"
)

const (
	energyStatsPrimary   = byte(0xB5)
	energyStatsSecondary = byte(0x16)
)

func (plane *plane) DecodeBroadcast(frame protocol.Frame) (map[string]types.Value, bool, error) {
	if frame.Primary != energyStatsPrimary || frame.Secondary != energyStatsSecondary {
		return nil, false, nil
	}

	selector, ok := parseEnergyStatsSelector(frame.Data)
	if !ok {
		return nil, false, nil
	}

	values := map[string]types.Value{
		"payload":         {Value: append([]byte(nil), frame.Data...), Valid: true},
		"selector":        {Value: append([]byte(nil), selector.raw...), Valid: true},
		"period_code":     {Value: selector.period, Valid: true},
		"source_code":     {Value: selector.source, Valid: true},
		"usage_code":      {Value: selector.usage, Valid: true},
		"selector_w":      {Value: selector.w, Valid: true},
		"selector_v":      {Value: selector.v, Valid: true},
		"selector_q":      {Value: selector.q, Valid: true},
		"selector_is_rsp": {Value: selector.isResponse, Valid: true},
	}

	if periodName, ok := energyStatsPeriodName(selector.period); ok {
		values["period"] = types.Value{Value: periodName, Valid: true}
	} else {
		values["period"] = types.Value{Valid: false}
	}
	if sourceName, ok := energyStatsSourceName(selector.source); ok {
		values["source"] = types.Value{Value: sourceName, Valid: true}
	} else {
		values["source"] = types.Value{Valid: false}
	}
	if usageName, ok := energyStatsUsageName(selector.usage); ok {
		values["usage"] = types.Value{Value: usageName, Valid: true}
	} else {
		values["usage"] = types.Value{Valid: false}
	}

	if yearKind, ok := energyStatsYearKind(selector.period, selector.q); ok {
		values["year_kind"] = types.Value{Value: yearKind, Valid: true}
	} else {
		values["year_kind"] = types.Value{Valid: false}
	}

	if month, ok := energyStatsMonth(selector.period, selector.w, selector.q); ok {
		values["month"] = types.Value{Value: month, Valid: true}
	} else {
		values["month"] = types.Value{Valid: false}
	}

	if day, ok := energyStatsDay(selector.period, selector.w, selector.v); ok {
		values["day"] = types.Value{Value: day, Valid: true}
	} else {
		values["day"] = types.Value{Valid: false}
	}

	if selector.isResponse {
		responseSchema := schema.Schema{
			Fields: []schema.SchemaField{
				{Name: "wh", Type: types.EXP{}},
			},
		}
		decoded, err := responseSchema.Decode(frame.Data[len(frame.Data)-4:])
		if err != nil {
			return nil, true, fmt.Errorf("energy stats decode wh: %w", err)
		}
		values["wh"] = decoded["wh"]
	}

	return values, true, nil
}

type energyStatsSelector struct {
	raw        []byte
	period     uint8
	source     uint8
	usage      uint8
	w          uint8
	v          uint8
	q          uint8
	isResponse bool
}

func parseEnergyStatsSelector(payload []byte) (energyStatsSelector, bool) {
	switch len(payload) {
	case 8:
		if payload[0] != 0x10 {
			return energyStatsSelector{}, false
		}
		if payload[2] != 0xFF || payload[3] != 0xFF {
			return energyStatsSelector{}, false
		}
		if payload[7]&0xF0 != 0x30 {
			return energyStatsSelector{}, false
		}

		wv := payload[6]
		return energyStatsSelector{
			raw:        append([]byte(nil), payload...),
			period:     payload[1] & 0x0F,
			source:     payload[4] & 0x0F,
			usage:      payload[5] & 0x0F,
			w:          (wv >> 4) & 0x0F,
			v:          wv & 0x0F,
			q:          payload[7] & 0x0F,
			isResponse: false,
		}, true
	case 11:
		if payload[0] > 0x03 {
			return energyStatsSelector{}, false
		}
		if payload[6]&0xF0 != 0x30 {
			return energyStatsSelector{}, false
		}

		wv := payload[5]
		return energyStatsSelector{
			raw:        append([]byte(nil), payload[:7]...),
			period:     payload[0] & 0x0F,
			source:     payload[3] & 0x0F,
			usage:      payload[4] & 0x0F,
			w:          (wv >> 4) & 0x0F,
			v:          wv & 0x0F,
			q:          payload[6] & 0x0F,
			isResponse: true,
		}, true
	default:
		return energyStatsSelector{}, false
	}
}

func energyStatsPeriodName(code uint8) (string, bool) {
	switch code {
	case 0:
		return "all", true
	case 1:
		return "day", true
	case 2:
		return "month", true
	case 3:
		return "year", true
	default:
		return "", false
	}
}

func energyStatsSourceName(code uint8) (string, bool) {
	switch code {
	case 1:
		return "solar", true
	case 2:
		return "environmental", true
	case 3:
		return "electricity", true
	case 4:
		return "gas", true
	case 9:
		return "heat_pump", true
	default:
		return "", false
	}
}

func energyStatsUsageName(code uint8) (string, bool) {
	switch code {
	case 0:
		return "all", true
	case 3:
		return "heating", true
	case 4:
		return "hot_water", true
	case 5:
		return "cooling", true
	default:
		return "", false
	}
}

func energyStatsYearKind(period, q uint8) (string, bool) {
	switch period {
	case 1, 2, 3:
	default:
		return "", false
	}

	switch q {
	case 0, 1:
		return "previous", true
	case 2, 3:
		return "current", true
	default:
		return "", false
	}
}

func energyStatsMonth(period, w, q uint8) (uint8, bool) {
	if period != 1 && period != 2 {
		return 0, false
	}

	baseW := w
	if period == 1 && baseW%2 == 1 {
		baseW--
	}

	switch q {
	case 0, 2:
		month := baseW / 2
		if month < 1 || month > 7 {
			return 0, false
		}
		return month, true
	case 1, 3:
		month := 8 + baseW/2
		if month < 8 || month > 12 {
			return 0, false
		}
		return month, true
	default:
		return 0, false
	}
}

func energyStatsDay(period, w, v uint8) (uint8, bool) {
	if period != 1 {
		return 0, false
	}

	day := v
	if w%2 == 1 {
		day = 16 + v
	}
	if day < 1 || day > 31 {
		return 0, false
	}
	return day, true
}
