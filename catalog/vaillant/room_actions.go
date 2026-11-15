package vaillant

import (
	_ "embed"
	roomActionJSON "encoding/json"
)

//go:embed data/room_actions.json
var roomActionRoomActionsJSON []byte

// RoomActions returns a caller-owned copy of the embedded room action catalog JSON.
func RoomActions() roomActionJSON.RawMessage {
	return append(roomActionJSON.RawMessage(nil), roomActionRoomActionsJSON...)
}
