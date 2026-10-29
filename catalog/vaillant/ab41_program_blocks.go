package vaillant

import (
	_ "embed"
	ab41json "encoding/json"
)

//go:embed data/ab41_program_blocks.json
var ab41ProgramBlocks []byte

func Ab41ProgramBlocks() ab41json.RawMessage {
	return append(ab41json.RawMessage(nil), ab41ProgramBlocks...)
}
