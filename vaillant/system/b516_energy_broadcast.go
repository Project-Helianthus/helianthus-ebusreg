package system

import (
	"encoding/binary"
	"fmt"
	"time"

	"github.com/Project-Helianthus/helianthus-ebusgo/protocol"
	"github.com/Project-Helianthus/helianthus-ebusgo/types"
	"github.com/Project-Helianthus/helianthus-ebusreg/schema"
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
		"selector_w":      {Value: uint8((selector.dateWord >> 4) & 0x0F), Valid: true},
		"selector_v":      {Value: uint8(selector.dateWord & 0x0F), Valid: true},
		"selector_q":      {Value: uint8((selector.dateWord >> 8) & 0x0F), Valid: true},
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

	// Dates identify absolute windows. Relative labels require a consumer's
	// reference date and cannot be inferred from the selector alone.
	values["year"] = types.Value{Valid: false}
	values["month"] = types.Value{Valid: false}
	values["day"] = types.Value{Valid: false}
	if year, month, day, ok := energyStatsDate(selector.period, selector.dateWord); ok {
		values["year"] = types.Value{Value: year, Valid: true}
		if selector.period == 1 || selector.period == 2 {
			values["month"] = types.Value{Value: month, Valid: true}
		}
		if selector.period == 1 {
			values["day"] = types.Value{Value: day, Valid: true}
		}
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
	dateWord   uint16
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
		return energyStatsSelector{
			raw:        append([]byte(nil), payload...),
			period:     payload[1] & 0x0F,
			source:     payload[4] & 0x0F,
			usage:      payload[5] & 0x0F,
			dateWord:   binary.LittleEndian.Uint16(payload[6:8]),
			isResponse: false,
		}, true
	case 11:
		if payload[0] > 0x03 {
			return energyStatsSelector{}, false
		}
		return energyStatsSelector{
			raw:        append([]byte(nil), payload[:7]...),
			period:     payload[0] & 0x0F,
			source:     payload[3] & 0x0F,
			usage:      payload[4] & 0x0F,
			dateWord:   binary.LittleEndian.Uint16(payload[5:7]),
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

// energyStatsDate validates the components used by the selected time base.
// Total selectors have no date window, even when replies echo a calendar date.
// An invalid date does not invalidate independently decoded energy fields.
func energyStatsDate(period uint8, word uint16) (year uint16, month, day uint8, ok bool) {
	year = 2000 + (word >> 9)
	month = uint8((word >> 5) & 0x0F)
	day = uint8(word & 0x1F)

	switch period {
	case 1:
		if month < 1 || month > 12 || day < 1 {
			return 0, 0, 0, false
		}
		date := time.Date(int(year), time.Month(month), int(day), 0, 0, 0, 0, time.UTC)
		if date.Month() != time.Month(month) || date.Day() != int(day) {
			return 0, 0, 0, false
		}
	case 2:
		if month < 1 || month > 12 || day != 0 {
			return 0, 0, 0, false
		}
	case 3:
		if month != 0 || day != 0 {
			return 0, 0, 0, false
		}
	default:
		return 0, 0, 0, false
	}
	return year, month, day, true
}
