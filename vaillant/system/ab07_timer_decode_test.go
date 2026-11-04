package system

import (
	"bytes"
	"fmt"
	"testing"
)

func TestAb07TimerCompleteReply(ab07T *testing.T) {
	ab07Payload := []byte{0xA5, 0x06, 0x0C, 0x90, 0x90, 0x8F, 0x90}
	ab07Values := decodeTimerResponse(0, 0, 1, 0, ab07Payload)
	if ab07Got := ab07Values["slot_count"]; !ab07Got.Valid || ab07Got.Value != 3 {
		ab07T.Fatalf("slot_count = %+v; want 3 valid", ab07Got)
	}
	if ab07Got := ab07Values["ab07_parameter_configuration"]; !ab07Got.Valid || ab07Got.Value != byte(0xA5) {
		ab07T.Fatalf("raw parameter configuration = %+v", ab07Got)
	}
	ab07Slots, ab07OK := ab07Values["ab07_slots"].Value.([]Ab07TimerSlot)
	if !ab07Values["ab07_slots"].Valid || !ab07OK || len(ab07Slots) != 3 {
		ab07T.Fatalf("slots = %+v", ab07Values["ab07_slots"])
	}
	ab07Pairs := [][2]byte{{6, 12}, {0x90, 0x90}, {0x8F, 0x90}}
	for ab07Index, ab07Slot := range ab07Slots {
		if ab07Slot.Ab07RawPair != ab07Pairs[ab07Index] || ab07Slot.Ab07Unused != (ab07Index == 1) {
			ab07T.Fatalf("slot %d = %+v", ab07Index, ab07Slot)
		}
		if ab07Index == 1 {
			if ab07Slot.Ab07StartMinutes != nil || ab07Slot.Ab07StopMinutes != nil {
				ab07T.Fatal("unused slot has minute projections")
			}
		} else if ab07Slot.Ab07StartMinutes == nil || ab07Slot.Ab07StopMinutes == nil ||
			*ab07Slot.Ab07StartMinutes != int(ab07Pairs[ab07Index][0])*10 ||
			*ab07Slot.Ab07StopMinutes != int(ab07Pairs[ab07Index][1])*10 {
			ab07T.Fatalf("slot %d minute projections = %+v", ab07Index, ab07Slot)
		}
	}
}

func TestAb07TimerPairBoundaries(ab07T *testing.T) {
	for ab07Position := 0; ab07Position < 3; ab07Position++ {
		for _, ab07Pair := range [][2]byte{{0, 1}, {0, 0x90}, {0x8F, 0x90}, {0x90, 0x90}, {0x91, 0x92}, {0, 0x91}, {0x92, 0x90}, {12, 6}, {6, 6}, {0, 0}} {
			ab07T.Run(fmt.Sprintf("ab07_position_%d_pair_%02x_%02x", ab07Position, ab07Pair[0], ab07Pair[1]), func(ab07T *testing.T) {
				ab07Payload := []byte{0xA5, 1, 2, 3, 4, 5, 6}
				ab07Payload[1+2*ab07Position], ab07Payload[2+2*ab07Position] = ab07Pair[0], ab07Pair[1]
				ab07Values := decodeTimerResponse(0, 0, 1, 0, ab07Payload)
				if ab07Got := ab07Values["slot_count"]; !ab07Got.Valid || ab07Got.Value != 3 {
					ab07T.Fatalf("complete width lost slot count: %+v", ab07Got)
				}
				ab07Slots := ab07Values["ab07_slots"].Value.([]Ab07TimerSlot)
				for ab07Index, ab07Slot := range ab07Slots {
					ab07Start, ab07Stop := ab07Payload[1+2*ab07Index], ab07Payload[2+2*ab07Index]
					if ab07Slot.Ab07RawPair != [2]byte{ab07Start, ab07Stop} || ab07Slot.Ab07Unused != (ab07Start == 0x90 && ab07Stop == 0x90) {
						ab07T.Fatalf("slot %d = %+v", ab07Index, ab07Slot)
					}
					if ab07Start < ab07Stop && ab07Stop <= 0x90 {
						if ab07Slot.Ab07StartMinutes == nil || ab07Slot.Ab07StopMinutes == nil || *ab07Slot.Ab07StartMinutes != int(ab07Start)*10 || *ab07Slot.Ab07StopMinutes != int(ab07Stop)*10 {
							ab07T.Fatalf("valid slot %d lost minutes: %+v", ab07Index, ab07Slot)
						}
					} else if ab07Slot.Ab07StartMinutes != nil || ab07Slot.Ab07StopMinutes != nil {
						ab07T.Fatalf("invalid or unused slot %d invented minutes: %+v", ab07Index, ab07Slot)
					}
				}
			})
		}
	}
}

func TestAb07TimerIncompleteWidths(ab07T *testing.T) {
	for _, ab07Length := range []int{0, 1, 2, 3, 4, 5, 6, 8} {
		ab07T.Run(fmt.Sprintf("ab07_width_%d", ab07Length), func(ab07T *testing.T) {
			ab07Payload := []byte{0xA5, 6, 12, 0x90, 0x90, 0x8F, 0x90, 0}[0:ab07Length]
			ab07Values := decodeTimerResponse(0, 0, 1, 0, ab07Payload)
			if ab07Got := ab07Values["payload"]; !ab07Got.Valid || !bytes.Equal(ab07Got.Value.([]byte), ab07Payload) {
				ab07T.Fatalf("payload = %+v", ab07Got)
			}
			if ab07Got := ab07Values["value"]; ab07Got.Valid != (ab07Length > 0) || (ab07Length > 0 && !bytes.Equal(ab07Got.Value.([]byte), ab07Payload)) {
				ab07T.Fatalf("raw value = %+v", ab07Got)
			}
			for _, ab07Key := range []string{"slot_count", "ab07_slots", "ab07_parameter_configuration"} {
				if ab07Values[ab07Key].Valid {
					ab07T.Fatalf("incomplete reply invented %s: %+v", ab07Key, ab07Values[ab07Key])
				}
			}
		})
	}
}

func TestAb07TimerRawOwnership(ab07T *testing.T) {
	ab07Payload := []byte{0xA5, 6, 12, 0x90, 0x90, 0x8F, 0x90}
	ab07Want := append([]byte(nil), ab07Payload...)
	ab07Values := decodeTimerResponse(0, 0, 1, 0, ab07Payload)
	for ab07Index := range ab07Payload {
		ab07Payload[ab07Index] = 0
	}
	for _, ab07Key := range []string{"payload", "value"} {
		if !bytes.Equal(ab07Values[ab07Key].Value.([]byte), ab07Want) {
			ab07T.Fatalf("%s aliases input", ab07Key)
		}
	}
	ab07Values["value"].Value.([]byte)[1] = 0
	if !bytes.Equal(ab07Values["payload"].Value.([]byte), ab07Want) {
		ab07T.Fatal("raw value aliases raw payload")
	}
	ab07Values["payload"].Value.([]byte)[2] = 0
	ab07Slots := ab07Values["ab07_slots"].Value.([]Ab07TimerSlot)
	if ab07Slots[0].Ab07RawPair != [2]byte{6, 12} || *ab07Slots[0].Ab07StartMinutes != 60 || *ab07Slots[0].Ab07StopMinutes != 120 || ab07Values["ab07_parameter_configuration"].Value != byte(0xA5) {
		ab07T.Fatal("structured slot or parameter configuration aliases raw storage")
	}
	ab07Slots[0].Ab07RawPair[0] = 1
	if ab07Values["payload"].Value.([]byte)[1] != 6 {
		ab07T.Fatal("structured raw pair aliases payload")
	}
}
