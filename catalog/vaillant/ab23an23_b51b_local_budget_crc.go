package vaillant

import ab23an23Errors "errors"

type Ab23An23Procedure string

const (
	Ab23An23Grouped1024 Ab23An23Procedure = "ab23an23_grouped_1024"
	Ab23An23Single205   Ab23An23Procedure = "ab23an23_single_205"
)

type Ab23An23Budget struct {
	Ab23An23Procedure Ab23An23Procedure
	Ab23An23Chunks    [][]byte
	Ab23An23Consumed  int
	Ab23An23Unplanned int
}

func Ab23An23Plan(ab23an23Input []byte, ab23an23Procedure Ab23An23Procedure) (Ab23An23Budget, error) {
	ab23an23Result := Ab23An23Budget{Ab23An23Procedure: ab23an23Procedure, Ab23An23Unplanned: len(ab23an23Input)}
	var ab23an23Caps []int
	switch ab23an23Procedure {
	case Ab23An23Grouped1024:
		ab23an23Caps = []int{206, 206, 206, 206, 200}
	case Ab23An23Single205:
		ab23an23Caps = []int{205}
	default:
		return ab23an23Result, ab23an23Errors.New("ab23an23: unknown procedure")
	}
	for _, ab23an23Cap := range ab23an23Caps {
		if ab23an23Result.Ab23An23Unplanned == 0 {
			break
		}
		ab23an23Length := ab23an23Cap
		if ab23an23Length > ab23an23Result.Ab23An23Unplanned {
			ab23an23Length = ab23an23Result.Ab23An23Unplanned
		}
		ab23an23End := ab23an23Result.Ab23An23Consumed + ab23an23Length
		ab23an23Chunk := append([]byte(nil), ab23an23Input[ab23an23Result.Ab23An23Consumed:ab23an23End]...)
		ab23an23Result.Ab23An23Chunks = append(ab23an23Result.Ab23An23Chunks, ab23an23Chunk)
		ab23an23Result.Ab23An23Consumed = ab23an23End
		ab23an23Result.Ab23An23Unplanned -= ab23an23Length
	}
	return ab23an23Result, nil
}

func Ab23An23Progress(ab23an23Prior int, ab23an23Accepted bool, ab23an23Actual int) (int, error) {
	if ab23an23Prior < 0 || ab23an23Actual < 0 {
		return ab23an23Prior, ab23an23Errors.New("ab23an23: negative progress or length")
	}
	if !ab23an23Accepted {
		return ab23an23Prior, nil
	}
	if ab23an23Actual > int(^uint(0)>>1)-ab23an23Prior {
		return ab23an23Prior, ab23an23Errors.New("ab23an23: progress overflow")
	}
	return ab23an23Prior + ab23an23Actual, nil
}

func Ab23An23CRC(ab23an23Bytes []byte) uint16 {
	return Ab23An23CRCWithInitial(ab23an23Bytes, 0xffff)
}

func Ab23An23CRCWithInitial(ab23an23Bytes []byte, ab23an23Initial uint16) uint16 {
	ab23an23CRC := ab23an23Initial
	for _, ab23an23Byte := range ab23an23Bytes {
		ab23an23CRC ^= uint16(ab23an23Byte) << 8
		for ab23an23Bit := 0; ab23an23Bit < 8; ab23an23Bit++ {
			if ab23an23CRC&0x8000 != 0 {
				ab23an23CRC = (ab23an23CRC << 1) ^ 0x1021
			} else {
				ab23an23CRC <<= 1
			}
		}
	}
	return ab23an23CRC
}
