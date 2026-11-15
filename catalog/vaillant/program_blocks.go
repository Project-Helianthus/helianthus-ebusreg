package vaillant

import (
	_ "embed"
	programBlockJSON "encoding/json"
)

//go:embed data/program_blocks.json
var programBlockProgramBlocks []byte

// ProgramBlocks returns a caller-owned copy of the embedded program block catalog JSON.
func ProgramBlocks() programBlockJSON.RawMessage {
	return append(programBlockJSON.RawMessage(nil), programBlockProgramBlocks...)
}
