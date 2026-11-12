package ebusstandard

import _ "embed"

//go:embed testdata/ab49_fe_metadata.json
var standardFEMetadataFEMetadataJSON []byte

func StandardFEMetadata() []byte {
	return append([]byte(nil), standardFEMetadataFEMetadataJSON...)
}
