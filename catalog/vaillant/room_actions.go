package vaillant

import (
	_ "embed"
	roomActionJSON "encoding/json"
)

//go:embed data/ab29_room_actions.json
var roomActionRoomActionsJSON []byte

func RoomActions() roomActionJSON.RawMessage {
	return append(roomActionJSON.RawMessage(nil), roomActionRoomActionsJSON...)
}
