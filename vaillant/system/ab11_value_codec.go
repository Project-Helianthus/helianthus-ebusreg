package system

import (
	"bytes"
	"encoding/binary"
	"math"
)

type Ab11Codec string

const (
	Ab11BOOL Ab11Codec = "BOOL"
	Ab11UCH  Ab11Codec = "UCH"
	Ab11UIN  Ab11Codec = "UIN"
	Ab11I8   Ab11Codec = "I8"
	Ab11I16  Ab11Codec = "I16"
	Ab11U32  Ab11Codec = "U32"
	Ab11I32  Ab11Codec = "I32"
	Ab11EXP  Ab11Codec = "EXP"
	Ab11DATE Ab11Codec = "DATE_DMY_U8X3"
	Ab11TIME Ab11Codec = "TIME_HMS_U8X3"
	Ab11FW   Ab11Codec = "FW"
	Ab11FWU  Ab11Codec = "FWU"
	Ab11STR  Ab11Codec = "STR"
	Ab11HEX  Ab11Codec = "HEX"
)

type Ab11Grade string

const (
	Ab11Unknown    Ab11Grade = "Unknown"
	Ab11Hypothesis Ab11Grade = "Hypothesis"
	Ab11Proven     Ab11Grade = "Proven"
)

// Ab11Value keeps native evidence independently of an explicitly requested
// interpretation. Ab11Qualified describes software decoding with a Proven
// caller-supplied grade, never physical support or register qualification.
type Ab11Value struct {
	Ab11Raw            []byte
	Ab11Grade          Ab11Grade
	Ab11Interpretation any
	Ab11Qualified      bool
}

type Ab11Date struct {
	Ab11Day   uint8
	Ab11Month uint8
	Ab11Year  uint8
}

type Ab11Time struct {
	Ab11Hour   uint8
	Ab11Minute uint8
	Ab11Second uint8
}

// Ab11String retains every byte; Ab11Latin1 maps each byte to its Unicode
// code point without trimming or interpreting termination.
type Ab11String struct {
	Ab11Bytes  []byte
	Ab11Latin1 string
}

// Ab11DecodeValue requires separate codec, byte width and evidence grade.
// For STR:n and HEX:n select Ab11STR or Ab11HEX and pass n as ab11Width.
// Failed admission returns owned raw evidence and no interpretation.
// Hypothesis interpretations remain Hypothesis and are never qualified.
func Ab11DecodeValue(ab11Codec Ab11Codec, ab11Width int, ab11Grade Ab11Grade, ab11Raw []byte) Ab11Value {
	ab11Result := Ab11Value{Ab11Raw: append([]byte(nil), ab11Raw...), Ab11Grade: ab11Grade}
	if ab11Grade != Ab11Proven && ab11Grade != Ab11Hypothesis {
		return ab11Result
	}
	ab11Expected := 0
	switch ab11Codec {
	case Ab11BOOL, Ab11UCH, Ab11I8:
		ab11Expected = 1
	case Ab11UIN, Ab11I16:
		ab11Expected = 2
	case Ab11U32, Ab11I32, Ab11EXP:
		ab11Expected = 4
	case Ab11DATE, Ab11TIME, Ab11FW, Ab11FWU:
		ab11Expected = 3
	case Ab11STR, Ab11HEX:
		ab11Expected = ab11Width
	default:
		return ab11Result
	}
	if ab11Width <= 0 || ab11Width != ab11Expected || len(ab11Raw) != ab11Width {
		return ab11Result
	}
	switch ab11Codec {
	case Ab11BOOL:
		if ab11Raw[0] > 1 {
			return ab11Result
		}
		ab11Result.Ab11Interpretation = ab11Raw[0] == 1
	case Ab11UCH:
		ab11Result.Ab11Interpretation = ab11Raw[0]
	case Ab11UIN:
		ab11Result.Ab11Interpretation = binary.LittleEndian.Uint16(ab11Raw)
	case Ab11I8:
		ab11Result.Ab11Interpretation = int8(ab11Raw[0])
	case Ab11I16:
		ab11Result.Ab11Interpretation = int16(binary.LittleEndian.Uint16(ab11Raw))
	case Ab11U32:
		ab11Result.Ab11Interpretation = binary.LittleEndian.Uint32(ab11Raw)
	case Ab11I32:
		ab11Result.Ab11Interpretation = int32(binary.LittleEndian.Uint32(ab11Raw))
	case Ab11EXP:
		ab11Float := math.Float32frombits(binary.LittleEndian.Uint32(ab11Raw))
		if math.IsNaN(float64(ab11Float)) || math.IsInf(float64(ab11Float), 0) {
			return ab11Result
		}
		ab11Result.Ab11Interpretation = ab11Float
	case Ab11DATE:
		ab11Result.Ab11Interpretation = Ab11Date{ab11Raw[0], ab11Raw[1], ab11Raw[2]}
	case Ab11TIME:
		ab11Result.Ab11Interpretation = Ab11Time{ab11Raw[0], ab11Raw[1], ab11Raw[2]}
	case Ab11FW, Ab11FWU:
		if ab11Codec == Ab11FWU && bytes.Equal(ab11Raw, []byte{255, 255, 255}) {
			// All-0xFF is the unavailable-firmware sentinel; keep raw only.
			return ab11Result
		}
		ab11Parts := [3]uint8{}
		for ab11Index, ab11Byte := range ab11Raw {
			if ab11Codec == Ab11FW {
				if ab11Byte>>4 > 9 || ab11Byte&15 > 9 {
					return ab11Result
				}
				ab11Parts[ab11Index] = (ab11Byte>>4)*10 + ab11Byte&15
			} else {
				ab11Parts[ab11Index] = ab11Byte
			}
		}
		ab11Result.Ab11Interpretation = ab11Parts
	case Ab11STR:
		ab11Runes := make([]rune, len(ab11Raw))
		for ab11Index, ab11Byte := range ab11Raw {
			ab11Runes[ab11Index] = rune(ab11Byte)
		}
		ab11Result.Ab11Interpretation = Ab11String{append([]byte(nil), ab11Raw...), string(ab11Runes)}
	case Ab11HEX:
		ab11Result.Ab11Interpretation = append([]byte(nil), ab11Raw...)
	}
	ab11Result.Ab11Qualified = ab11Grade == Ab11Proven
	return ab11Result
}
