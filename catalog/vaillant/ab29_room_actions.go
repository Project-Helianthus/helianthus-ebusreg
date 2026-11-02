package vaillant

import (
	_ "embed"
	ab29json "encoding/json"
)

//go:embed data/ab29_room_actions.json
var ab29RoomActionsJSON []byte

func Ab29RoomActions() ab29json.RawMessage {
	return append(ab29json.RawMessage(nil), ab29RoomActionsJSON...)
}
