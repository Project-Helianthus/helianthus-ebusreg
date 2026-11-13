package system

import (
	"bytes"
	"errors"
	"fmt"
	"testing"
)

func eventDecodeFixture(eventDecodeOpcode byte) Event {
	eventDecodeResponse := []byte{0xA0, 0x01, 0x12, 0x23, 0x34, 0x45, 0x56, 0x67}
	if eventDecodeOpcode == 0x0B {
		eventDecodeResponse = []byte{0xB0, 0x78, 0x89, 0x9A, 0xAB, 0xBC, 0xCD, 0xDE}
	}
	return DecodeEvent("synthetic-family", []byte{eventDecodeOpcode, 0, 0xFF, 1, 0xFE}, eventDecodeResponse, nil)
}

func TestAb23ExactRawAndOwnership(eventDecodeT *testing.T) {
	for _, eventDecodeOpcode := range []byte{0x09, 0x0B} {
		eventDecodeT.Run(fmt.Sprintf("ab23-op-%02x", eventDecodeOpcode), func(eventDecodeT *testing.T) {
			eventDecodeExpected := eventDecodeFixture(eventDecodeOpcode)
			eventDecodeRequest := []byte{eventDecodeOpcode, 0, 0xFF, 1, 0xFE}
			eventDecodeResponse := []byte{0xA0, 0x01, 0x12, 0x23, 0x34, 0x45, 0x56, 0x67}
			if eventDecodeOpcode == 0x0B {
				eventDecodeResponse = []byte{0xB0, 0x78, 0x89, 0x9A, 0xAB, 0xBC, 0xCD, 0xDE}
			}
			eventDecodeEvent := DecodeEvent("synthetic-family", eventDecodeRequest, eventDecodeResponse, nil)
			eventDecodeConfiguration, eventDecodeOK := eventDecodeEvent.EventDecodeConfiguration()
			if !eventDecodeOK || eventDecodeConfiguration != eventDecodeResponse[0] {
				eventDecodeT.Fatal("ab23 configuration changed")
			}
			eventDecodeRaw, eventDecodeOK := eventDecodeEvent.EventDecodeRaw()
			if !eventDecodeOK || !bytes.Equal(eventDecodeRaw[:], eventDecodeResponse[1:]) {
				eventDecodeT.Fatal("ab23 raw positions changed")
			}
			if !eventDecodeEvent.EventDecodeValid() || eventDecodeEvent.EventDecodeFailure() != "" || eventDecodeEvent.EventDecodeFamily() != "synthetic-family" {
				eventDecodeT.Fatal("ab23 observation state")
			}
			for eventDecodeIndex := range eventDecodeResponse {
				eventDecodeResponse[eventDecodeIndex] ^= 0xFF
			}
			for eventDecodeIndex := range eventDecodeRequest {
				eventDecodeRequest[eventDecodeIndex] ^= 0xFF
			}
			if eventDecodeEvent != eventDecodeExpected {
				eventDecodeT.Fatal("ab23 input alias")
			}
			eventDecodeOutput := eventDecodeEvent.EventDecodeResponse()
			eventDecodeOutput[0] = 0
			eventDecodeOutput = eventDecodeEvent.EventDecodeRequest()
			eventDecodeOutput[0] = 0
			eventDecodeRaw[0] = 0
			if eventDecodeEvent != eventDecodeExpected {
				eventDecodeT.Fatal("ab23 output alias")
			}
		})
	}
}

func TestAb23PairDimensions(eventDecodeT *testing.T) {
	eventDecodeFirst, eventDecodeSecond := eventDecodeFixture(0x09), eventDecodeFixture(0x0B)
	if !PairEvents(eventDecodeFirst, eventDecodeSecond).EventDecodeComplete() || !PairEvents(eventDecodeSecond, eventDecodeFirst).EventDecodeComplete() {
		eventDecodeT.Fatal("ab23 valid pair")
	}
	for _, eventDecodeDimension := range []string{"family", "system", "instance", "address", "code", "opcode"} {
		eventDecodeT.Run("ab23-"+eventDecodeDimension, func(eventDecodeT *testing.T) {
			eventDecodeFamily, eventDecodeRequest := eventDecodeSecond.EventDecodeFamily(), eventDecodeSecond.EventDecodeRequest()
			switch eventDecodeDimension {
			case "family":
				eventDecodeFamily = "another-synthetic-family"
			case "system":
				eventDecodeRequest[1] = 1
			case "instance":
				eventDecodeRequest[2] = 0xFE
			case "address":
				eventDecodeRequest[3] = 2
			case "code":
				eventDecodeRequest[4] = 0xFD
			case "opcode":
				eventDecodeRequest[0] = 0x09
			}
			eventDecodeChanged := DecodeEvent(eventDecodeFamily, eventDecodeRequest, eventDecodeSecond.EventDecodeResponse(), nil)
			if !eventDecodeChanged.EventDecodeValid() {
				eventDecodeT.Fatal("ab23 cross-pair control must be valid")
			}
			eventDecodePair := PairEvents(eventDecodeFirst, eventDecodeChanged)
			if eventDecodePair.EventDecodeComplete() {
				eventDecodeT.Fatal("ab23 crossed pair dimension")
			}
			if eventDecodePair.EventDecodeFirst() != eventDecodeFirst || eventDecodePair.EventDecodeSecond() != eventDecodeChanged {
				eventDecodeT.Fatal("ab23 member loss")
			}
			if eventDecodeDimension != "opcode" {
				eventDecodeRequest[0] = 0x09
				eventDecodeMatching := DecodeEvent(eventDecodeFamily, eventDecodeRequest, eventDecodeFirst.EventDecodeResponse(), nil)
				if !PairEvents(eventDecodeMatching, eventDecodeChanged).EventDecodeComplete() {
					eventDecodeT.Fatal("ab23 changed identity matching control")
				}
			} else if !PairEvents(eventDecodeSecond, eventDecodeChanged).EventDecodeComplete() {
				eventDecodeT.Fatal("ab23 changed opcode matching control")
			}
		})
	}
}

func TestAb23IncompleteMembers(eventDecodeT *testing.T) {
	eventDecodeGood := eventDecodeFixture(0x09)
	eventDecodeOther := eventDecodeFixture(0x0B)
	for eventDecodeLength := 0; eventDecodeLength <= 9; eventDecodeLength++ {
		if eventDecodeLength == 8 {
			continue
		}
		eventDecodeT.Run(fmt.Sprintf("ab23-length-%d", eventDecodeLength), func(eventDecodeT *testing.T) {
			eventDecodeResponse := bytes.Repeat([]byte{0xA5}, eventDecodeLength)
			eventDecodeBad := DecodeEvent(eventDecodeOther.EventDecodeFamily(), eventDecodeOther.EventDecodeRequest(), eventDecodeResponse, nil)
			eventDecodeCheckIncomplete(eventDecodeT, eventDecodeGood, eventDecodeBad)
			if !bytes.Equal(eventDecodeBad.EventDecodeResponse(), eventDecodeResponse) {
				eventDecodeT.Fatal("ab23 malformed bytes lost")
			}
		})
	}
	eventDecodeFailed := DecodeEvent(eventDecodeOther.EventDecodeFamily(), eventDecodeOther.EventDecodeRequest(), eventDecodeOther.EventDecodeResponse(), errors.New("ab23 synthetic failure"))
	eventDecodeCheckIncomplete(eventDecodeT, eventDecodeGood, eventDecodeFailed)
	if eventDecodeFailed.EventDecodeFailure() != "ab23 synthetic failure" || !bytes.Equal(eventDecodeFailed.EventDecodeResponse(), eventDecodeOther.EventDecodeResponse()) {
		eventDecodeT.Fatal("ab23 failure evidence lost")
	}
	eventDecodeCheckIncomplete(eventDecodeT, eventDecodeGood, Event{})
	eventDecodeCheckIncomplete(eventDecodeT, Event{}, Event{})
}

func eventDecodeCheckIncomplete(eventDecodeT *testing.T, eventDecodeFirst, eventDecodeSecond Event) {
	eventDecodeT.Helper()
	for _, eventDecodePair := range []EventPair{PairEvents(eventDecodeFirst, eventDecodeSecond), PairEvents(eventDecodeSecond, eventDecodeFirst)} {
		if eventDecodePair.EventDecodeComplete() {
			eventDecodeT.Fatal("ab23 incomplete pair admitted")
		}
	}
	eventDecodePair := PairEvents(eventDecodeFirst, eventDecodeSecond)
	if eventDecodePair.EventDecodeFirst() != eventDecodeFirst || eventDecodePair.EventDecodeSecond() != eventDecodeSecond {
		eventDecodeT.Fatal("ab23 independent member lost")
	}
	if !eventDecodeSecond.EventDecodeValid() {
		if _, eventDecodeOK := eventDecodeSecond.EventDecodeConfiguration(); eventDecodeOK {
			eventDecodeT.Fatal("ab23 invalid configuration admitted")
		}
		if _, eventDecodeOK := eventDecodeSecond.EventDecodeRaw(); eventDecodeOK {
			eventDecodeT.Fatal("ab23 invalid raw admitted")
		}
	}
}

func TestAb23ContextAdmissionAndUnknownFamily(eventDecodeT *testing.T) {
	eventDecodeReply := eventDecodeFixture(0x09).EventDecodeResponse()
	for _, eventDecodeSystem := range []byte{0, 1, 3, 2, 0xFF} {
		for eventDecodeAddress := byte(0); eventDecodeAddress <= 4; eventDecodeAddress++ {
			eventDecodeRequest := []byte{0x09, eventDecodeSystem, 0xFF, eventDecodeAddress, 0xFF}
			eventDecodeEvent := DecodeEvent("unknown-family", eventDecodeRequest, eventDecodeReply, nil)
			eventDecodeWant := (eventDecodeSystem == 0 && eventDecodeAddress >= 1 && eventDecodeAddress <= 3) || ((eventDecodeSystem == 1 || eventDecodeSystem == 3) && eventDecodeAddress >= 1 && eventDecodeAddress <= 2)
			if eventDecodeEvent.EventDecodeValid() != eventDecodeWant || !bytes.Equal(eventDecodeEvent.EventDecodeRequest(), eventDecodeRequest) || !bytes.Equal(eventDecodeEvent.EventDecodeResponse(), eventDecodeReply) || eventDecodeEvent.EventDecodeFamily() != "unknown-family" {
				eventDecodeT.Fatalf("ab23 context %v", eventDecodeRequest)
			}
		}
	}
	for _, eventDecodeRequest := range [][]byte{nil, {0x09}, {0x09, 0, 0xFF, 1}, {0x09, 0, 0xFF, 1, 0xFF, 0}, {0x03, 0, 0xFF, 1, 0xFF}} {
		eventDecodeEvent := DecodeEvent("unknown-family", eventDecodeRequest, eventDecodeReply, nil)
		if eventDecodeEvent.EventDecodeValid() || !bytes.Equal(eventDecodeEvent.EventDecodeRequest(), eventDecodeRequest) {
			eventDecodeT.Fatal("ab23 malformed request admission")
		}
		eventDecodeCheckIncomplete(eventDecodeT, eventDecodeFixture(0x0B), eventDecodeEvent)
	}
}
