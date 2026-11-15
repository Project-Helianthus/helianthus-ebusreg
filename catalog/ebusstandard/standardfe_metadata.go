package ebusstandard

import _ "embed"

//go:embed testdata/standardfe_metadata.json
var standardFEMetadataFEMetadataJSON []byte

func StandardFEMetadata() []byte {
	return append([]byte(nil), standardFEMetadataFEMetadataJSON...)
}
