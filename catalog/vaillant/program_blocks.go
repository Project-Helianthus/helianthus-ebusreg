package vaillant

import (
	_ "embed"
	programBlockJSON "encoding/json"
)

//go:embed data/program_blocks.json
var programBlockProgramBlocks []byte

func ProgramBlocks() programBlockJSON.RawMessage {
	return append(programBlockJSON.RawMessage(nil), programBlockProgramBlocks...)
}
